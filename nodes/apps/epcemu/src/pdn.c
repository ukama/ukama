/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <time.h>

#include "netutil.h"
#include "pcrf.h"
#include "pdn.h"
#include "ue.h"

#define PDN_POLL_BATCH 8

static uint64_t monotonic_sec(void) {

    struct timespec now;

    if (clock_gettime(CLOCK_MONOTONIC, &now) != 0) return 0;
    return (uint64_t)now.tv_sec;
}

static int error_status(int code, const char *text, char *reason,
                        size_t reasonLen) {

    snprintf(reason, reasonLen, "%s", text);
    return code;
}

static int copy_field(JsonObj *body, const char *name, char *out,
                      size_t size, bool required) {

    JsonObj *value;
    const char *text;
    size_t len;

    value = json_object_get(body, name);
    if (value == NULL && !required) {
        out[0] = '\0';
        return USYS_TRUE;
    }
    if (!json_is_string(value)) return USYS_FALSE;
    text = json_string_value(value);
    len = json_string_length(value);
    if (len >= size || strlen(text) != len || (required && len == 0)) {
        return USYS_FALSE;
    }
    memcpy(out, text, len + 1);
    return USYS_TRUE;
}

int pdn_init(ServiceContext *ctx) {

    if (ctx == NULL || ctx->config == NULL) return USYS_FALSE;
    if (!apn_normalize(EPCEMU_DEF_APN, strlen(EPCEMU_DEF_APN),
                       ctx->config->apn)) return USYS_FALSE;
    return identity_create(ctx->instanceId);
}

static int mark_closing(const UeEntry *ue, void *arg) {

    (void)arg;
    ue_detach_start(ue->imsi, NULL);
    return USYS_TRUE;
}

static int mark_unknown(const UeEntry *ue, void *arg) {

    if (ue->state == UeStateAttached) {
        ue_set_policy(ue->imsi, "unknown", (const char *)arg, false);
    }
    return USYS_TRUE;
}

static int cleanup_one(ServiceContext *ctx, const UeEntry *ue) {

    int code;

    code = pcrf_delete_session(ctx->config, ue);
    if (code != 202) {
        usys_log_debug("PDN cleanup pending imsi=%s pdn=%s http=%d",
                       ue->imsi, ue->pdnId, code);
        return USYS_FALSE;
    }
    ue_detach_complete(ue->imsi);
    return USYS_TRUE;
}

static int cleanup_cb(const UeEntry *ue, void *arg) {

    ServiceContext *ctx;

    ctx = (ServiceContext *)arg;
    if (ue->state != UeStateDetaching) return USYS_TRUE;
    return cleanup_one(ctx, ue);
}

/* Called with serviceMutex held. Do not reuse old PDNs after recovery. */
static int sync_owner(ServiceContext *ctx) {

    PcrfEpc epc;
    int code;

    ctx->config->pcrfReady = false;
    if (ctx->ownerLost) return USYS_FALSE;

    memset(&epc, 0, sizeof(epc));
    code = pcrf_get_epc(ctx->config, &epc);
    if (code != 200) {
        ue_for_each(mark_unknown, "pcrf_unavailable");
        return USYS_FALSE;
    }

    if (ctx->pcrfRegistered &&
        (strcmp(epc.instanceId, ctx->instanceId) ||
         strcmp(epc.state, "ready"))) {
        ctx->pcrfRegistered = false;
        ctx->recovering = true;
        if (epc.instanceId[0] != '\0' &&
            strcmp(epc.instanceId, ctx->instanceId)) {
            ctx->ownerLost = true;
            usys_log_error("PCRF owner changed; restart EPCEMU");
        }
        ue_for_each(mark_closing, NULL);
    }

    if (ctx->recovering) {
        ue_for_each(cleanup_cb, ctx);
        if (ue_count_total() != 0 || ctx->ownerLost ||
            strcmp(epc.state, "ready")) return USYS_FALSE;
        if (!identity_create(ctx->instanceId)) return USYS_FALSE;
        ctx->recovering = false;
        ctx->expectedSet = false;
    }

    if (strcmp(epc.state, "ready")) return USYS_FALSE;
    if (!ctx->pcrfRegistered) {
        if (!ctx->expectedSet) {
            snprintf(ctx->expectedInstanceId,
                     sizeof(ctx->expectedInstanceId), "%s", epc.instanceId);
            ctx->expectedSet = true;
        }
        code = pcrf_register_epc(ctx->config, ctx->instanceId,
                                 ctx->expectedInstanceId);
        if (code == 409) {
            ctx->expectedSet = false;
        } else if (code == 410) {
            /* A timed-out startup may have been retired by a PCRF restart. */
            ctx->expectedSet = false;
            if (!identity_create(ctx->instanceId)) return USYS_FALSE;
        }
        if (code != 202) return USYS_FALSE;
        ctx->pcrfRegistered = true;
    }
    ctx->config->pcrfReady = true;
    return USYS_TRUE;
}

static int valid_address(EpcemuConfig *config, const char *ip) {

    char host[EPCEMU_MAX_STR];
    char *slash;
    uint32_t value;

    if (!ip_to_uint32(ip, &value) ||
        !ip_in_cidr(ip, config->ueCidr)) return USYS_FALSE;
    if ((value >> 24) == 0 || (value >> 24) == 127 ||
        (value >> 28) >= 14 || (value & 0xffff0000u) == 0xa9fe0000u) {
        return USYS_FALSE;
    }
    /* A UE cannot reserve the provisioned TUN or bridge address. */
    snprintf(host, sizeof(host), "%s", config->tunAddr);
    slash = strchr(host, '/');
    if (slash != NULL) *slash = '\0';
    if (!strcmp(host, ip)) return USYS_FALSE;
    snprintf(host, sizeof(host), "%s", config->bridgeCidr);
    slash = strchr(host, '/');
    if (slash != NULL) *slash = '\0';
    return strcmp(host, ip) != 0;
}

static int read_attach(ServiceContext *ctx, JsonObj *body, UeEntry *ue,
                        char *reason, size_t reasonLen) {

    JsonObj *apn;

    memset(ue, 0, sizeof(*ue));
    if (!json_is_object(body) ||
        !copy_field(body, "imsi", ue->imsi, sizeof(ue->imsi), true) ||
        !imsi_valid(ue->imsi)) {
        return error_status(400, "invalid_imsi", reason, reasonLen);
    }
    if (!copy_field(body, "iccid", ue->iccid, sizeof(ue->iccid), false)) {
        return error_status(400, "invalid_iccid", reason, reasonLen);
    }
    if (!copy_field(body, "ip", ue->ip, sizeof(ue->ip), true) ||
        !valid_address(ctx->config, ue->ip)) {
        return error_status(400, "invalid_ue_address", reason, reasonLen);
    }
    apn = json_object_get(body, "apn");
    if (apn == NULL) {
        snprintf(ue->apn, sizeof(ue->apn), "%s", ctx->config->apn);
    } else if (!json_is_string(apn) ||
               !apn_normalize(json_string_value(apn), json_string_length(apn),
                              ue->apn)) {
        return error_status(400, "invalid_apn", reason, reasonLen);
    }
    if (strcmp(ue->apn, ctx->config->apn)) {
        return error_status(403, "unsupported_apn", reason, reasonLen);
    }
    return 0;
}

static int same_request(const UeEntry *a, const UeEntry *b) {

    return !strcmp(a->ip, b->ip) && !strcmp(a->apn, b->apn) &&
           !strcmp(a->iccid, b->iccid);
}

static int refresh_one(ServiceContext *ctx, const UeEntry *ue) {

    PcrfPdn pdn;
    int code;

    code = pcrf_get_pdn(ctx->config, ue, &pdn);
    if (code == 404 || code == 410 ||
        (code == 200 && !strcmp(pdn.state, "closing"))) {
        ue_detach_start(ue->imsi, NULL);
        cleanup_one(ctx, ue);
        return USYS_FALSE;
    }
    if (code != 200) {
        ue_set_policy(ue->imsi, "unknown", pdn.reason, false);
        return USYS_FALSE;
    }
    ue_set_policy(ue->imsi, pdn.dataState, pdn.reason, true);
    return USYS_TRUE;
}

static int attach_locked(ServiceContext *ctx, const UeEntry *request,
                           JsonObj **reply, char *reason, size_t reasonLen) {

    UeEntry ue;
    PcrfPdn pdn;
    int code;
    int responseCode;
    bool exists;

    if (ctx->stopping || !sync_owner(ctx)) {
        return error_status(503, "pcrf_registration_pending",
                             reason, reasonLen);
    }

    exists = ue_get(request->imsi, &ue);
    if (exists) {
        if (!same_request(&ue, request)) {
            return error_status(409, "pdn_conflict", reason, reasonLen);
        }
        if (ue.state == UeStateDetaching) {
            return error_status(503, "cleanup_pending", reason, reasonLen);
        }
        if (ue.state == UeStateAttaching &&
            monotonic_sec() - ue.startedAt >= PDN_ATTACH_TIMEOUT_SEC) {
            ue_detach_start(ue.imsi, NULL);
            cleanup_one(ctx, &ue);
            return error_status(503, "attach_expired", reason, reasonLen);
        }
        if (ue.state == UeStateAttached) {
            if (!refresh_one(ctx, &ue)) {
                return error_status(503, "pdn_state_unavailable",
                                     reason, reasonLen);
            }
            *reply = ue_get_json(ue.imsi);
            return *reply != NULL ? 200 : 500;
        }
    } else {
        memcpy(&ue, request, sizeof(ue));
        snprintf(ue.instanceId, sizeof(ue.instanceId), "%s", ctx->instanceId);
        if (!identity_create(ue.pdnId)) {
            return error_status(500, "identity_failed", reason, reasonLen);
        }
        ue.startedAt = monotonic_sec();
        if (!ue_attach_start(&ue, reason, reasonLen)) return 409;
    }

    code = pcrf_create_session(ctx->config, &ue, &pdn);
    if (code == 202 && !strcmp(pdn.state, "live") &&
        (!strcmp(pdn.dataState, "allowed") ||
         !strcmp(pdn.dataState, "blocked"))) {
        ue_attach_complete(ue.imsi);
        ue_set_policy(ue.imsi, pdn.dataState, pdn.reason, true);
        *reply = ue_get_json(ue.imsi);
        return *reply != NULL ? 201 : 500;
    }

    /* Explicit rejection requires rollback, even if DELETE overtakes POST. */
    if ((code >= 400 && code < 500) || code == 202 ||
        (code == 503 && !strcmp(pdn.reason, "pdn_setup_failed"))) {
        ue_detach_start(ue.imsi, NULL);
        cleanup_one(ctx, &ue);
    }
    responseCode = (code == 400 || code == 403 || code == 409) ? code : 503;
    return error_status(responseCode, pdn.reason, reason, reasonLen);
}

int pdn_attach(ServiceContext *ctx, JsonObj *body, JsonObj **reply,
                char *reason, size_t reasonLen) {

    UeEntry ue;
    int code;

    *reply = NULL;
    pthread_mutex_lock(&ctx->serviceMutex);
    code = read_attach(ctx, body, &ue, reason, reasonLen);
    if (code == 0) code = attach_locked(ctx, &ue, reply, reason, reasonLen);
    pthread_mutex_unlock(&ctx->serviceMutex);
    return code;
}

int pdn_detach(ServiceContext *ctx, JsonObj *body,
                char *reason, size_t reasonLen) {

    char imsi[UE_IMSI_LEN];
    char instance[EPCEMU_ID_LEN];
    char id[EPCEMU_ID_LEN];
    UeEntry ue;
    int code;

    if (!copy_field(body, "imsi", imsi, sizeof(imsi), true) ||
        !imsi_valid(imsi) ||
        !copy_field(body, "epc_instance_id", instance, sizeof(instance),
                     true) || !identity_valid(instance) ||
        !copy_field(body, "pdn_id", id, sizeof(id), true) ||
        !identity_valid(id)) {
        return error_status(400, "invalid_pdn_identity", reason, reasonLen);
    }
    pthread_mutex_lock(&ctx->serviceMutex);
    code = 200;
    if (ue_get(imsi, &ue)) {
        if (strcmp(ue.instanceId, instance) || strcmp(ue.pdnId, id)) {
            code = error_status(409, "pdn_conflict", reason, reasonLen);
        } else {
            ue_detach_start(imsi, NULL);
            if (!cleanup_one(ctx, &ue)) {
                code = error_status(503, "cleanup_pending", reason, reasonLen);
            }
        }
    }
    pthread_mutex_unlock(&ctx->serviceMutex);
    return code;
}

static int poll_one(const UeEntry *ue, void *arg) {

    ServiceContext *ctx;
    uint64_t now;

    ctx = (ServiceContext *)arg;
    now = monotonic_sec();
    if (ue->state == UeStateAttaching && now >= ue->startedAt &&
        now - ue->startedAt >= PDN_ATTACH_TIMEOUT_SEC) {
        ue_detach_start(ue->imsi, NULL);
        return cleanup_one(ctx, ue);
    }
    if (ue->state == UeStateDetaching) return cleanup_one(ctx, ue);
    if (ue->state == UeStateAttached && ctx->config->pcrfReady) {
        return refresh_one(ctx, ue);
    }
    return USYS_TRUE;
}

void pdn_poll(ServiceContext *ctx) {

    UeEntry entries[EPCEMU_MAX_UES];
    size_t count;
    size_t i;
    size_t index;

    pthread_mutex_lock(&ctx->serviceMutex);
    if (!ctx->stopping) {
        sync_owner(ctx);
        count = ue_snapshot(entries, EPCEMU_MAX_UES);
        /* Bound polling work and advance past failed requests on retry. */
        for (i = 0; i < count && i < PDN_POLL_BATCH; i++) {
            index = ctx->pollOffset++ % count;
            if (!poll_one(&entries[index], ctx)) break;
        }
    }
    pthread_mutex_unlock(&ctx->serviceMutex);
}

void pdn_shutdown(ServiceContext *ctx) {

    pthread_mutex_lock(&ctx->serviceMutex);
    ctx->stopping = true;
    ctx->config->pcrfReady = false;
    ue_for_each(mark_closing, NULL);
    ue_for_each(cleanup_cb, ctx);
    if (ue_count_total() != 0) {
        usys_log_error("PDN cleanup incomplete; next EPC registration "
                        "must retire instance %s", ctx->instanceId);
    }
    pthread_mutex_unlock(&ctx->serviceMutex);
}
