/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <unistd.h>

#include "event.h"
#include "log.h"
#include "selector.h"
#include "util.h"

#define SIM_STATUS_WAIT_DEFAULT_SEC 120u
#define SIM_STATUS_WAIT_MAX_SEC     900u
#define SIM_STATUS_POLL_DEFAULT_SEC 2u
#define SIM_STATUS_POLL_MAX_SEC     30u

static package_t *event_package(event_ctx_t *ctx,
                                const event_spec_t *event,
                                ue_t *ue,
                                ulab_error_t *err) {
    package_t *pkg;

    if (event->package_ref[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing package", scenario_event_name(event->type));
        return NULL;
    }

    pkg = world_package_for_network(ctx->world, event->package_ref,
                                    ue->network_ref);
    if (pkg == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown package %.128s for UE %.128s",
                 event->package_ref, ue->ref);
        return NULL;
    }

    return pkg;
}

static int event_add_package_to_sim(event_ctx_t *ctx,
                                    const event_spec_t *event,
                                    ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];
        package_t *pkg = event_package(ctx, event, ue, err);

        if (pkg == NULL) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }

        if (bff_add_package_to_sim(ctx->bff, ue, pkg, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int event_purchase_package(event_ctx_t *ctx,
                                  const event_spec_t *event,
                                  ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    if (event->idempotency_key[0] != '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "BFF addPayment does not expose an idempotency key yet");
        return ULAB_ERR;
    }
    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue;
        package_t *pkg;
        subscriber_t *subscriber;
        subscriber_t payer;
        bff_payment_t payment;

        ue = &ctx->world->ues[sel.idx[i]];
        pkg = event_package(ctx, event, ue, err);
        subscriber = world_subscriber_by_ref(ctx->world,
                                             ue->subscriber_ref);
        if (pkg == NULL || subscriber == NULL) {
            selector_result_free(&sel);
            if (subscriber == NULL) {
                snprintf(err->msg, sizeof(err->msg),
                         "unknown subscriber for UE %.128s", ue->ref);
            }
            return ULAB_ERR;
        }

        payer = *subscriber;
        if (event->payer_email[0] != '\0') {
            ulab_copy(payer.email, sizeof(payer.email),
                      event->payer_email);
        }
        if (event->payer_phone[0] != '\0') {
            ulab_copy(payer.phone, sizeof(payer.phone),
                      event->payer_phone);
        }

        memset(&payment, 0, sizeof(payment));
        if (bff_record_cash_package_sale(ctx->bff, ue, pkg, &payer,
                                         event->amount,
                                         event->currency,
                                         &payment, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
        ulab_status("PAYMENT", "id=%s sim=%s package=%s status=%s",
                    payment.id, ue->ref, pkg->ref, payment.status);
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int event_set_package_active(event_ctx_t *ctx,
                                    const event_spec_t *event,
                                    ulab_error_t *err) {
    size_t i;
    size_t changed;

    if (event->package_ref[0] == '\0' || !event->has_active) {
        snprintf(err->msg, sizeof(err->msg),
                 "set_package_active requires package and active");
        return ULAB_ERR;
    }

    changed = 0;
    for (i = 0; i < ctx->world->package_count; i++) {
        package_t *pkg = &ctx->world->packages[i];

        if (!ulab_streq(pkg->base_ref, event->package_ref) &&
            !ulab_streq(pkg->ref, event->package_ref)) {
            continue;
        }
        if (bff_set_package_active(ctx->bff, pkg, event->active, err)) {
            return ULAB_ERR;
        }
        changed++;
    }

    if (changed == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "set_package_active unknown package %.128s",
                 event->package_ref);
        return ULAB_ERR;
    }
    return ULAB_OK;
}

static int event_remove_package_from_sim(event_ctx_t *ctx,
                                         const event_spec_t *event,
                                         ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    (void)event;

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];

        if (bff_clear_sim_packages(ctx->bff, ue, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int event_set_sim_status(event_ctx_t *ctx,
                                const event_spec_t *event,
                                ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    if (event->status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "set_sim_status missing status");
        return ULAB_ERR;
    }

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];

        if (bff_toggle_sim_status(ctx->bff, ue, event->status, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static unsigned int sim_status_poll_seconds(void) {
    const char *value;
    uint32_t seconds;

    value = getenv("ULAB_SIM_STATUS_POLL_SEC");
    if (value == NULL || value[0] == '\0' ||
        ulab_parse_u32(value, &seconds) != ULAB_OK || seconds == 0) {
        return SIM_STATUS_POLL_DEFAULT_SEC;
    }
    if (seconds > SIM_STATUS_POLL_MAX_SEC) {
        return SIM_STATUS_POLL_MAX_SEC;
    }
    return seconds;
}

static int event_wait_sim_status(event_ctx_t *ctx,
                                 const event_spec_t *event,
                                 ulab_error_t *err) {
    selector_result_t sel;
    unsigned char *matched;
    unsigned int timeout;
    unsigned int poll;
    unsigned int elapsed;
    size_t matched_count;
    size_t i;
    char last_sim[ULAB_MAX_ID];
    char last_status[ULAB_MAX_REF];
    char last_error[ULAB_MAX_ERR];

    memset(&sel, 0, sizeof(sel));
    memset(last_sim, 0, sizeof(last_sim));
    memset(last_status, 0, sizeof(last_status));
    memset(last_error, 0, sizeof(last_error));

    if (event->status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status missing status");
        return ULAB_ERR;
    }
    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }
    if (sel.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status selected no SIMs");
        selector_result_free(&sel);
        return ULAB_ERR;
    }

    timeout = event->amount_mb > 0 ?
        (unsigned int)event->amount_mb : SIM_STATUS_WAIT_DEFAULT_SEC;
    if (timeout > SIM_STATUS_WAIT_MAX_SEC) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status exceeds %u seconds",
                 SIM_STATUS_WAIT_MAX_SEC);
        selector_result_free(&sel);
        return ULAB_ERR;
    }

    matched = calloc(sel.count, sizeof(*matched));
    if (matched == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status allocation failed");
        selector_result_free(&sel);
        return ULAB_ERR;
    }

    poll = sim_status_poll_seconds();
    elapsed = 0;
    matched_count = 0;
    ulab_status("SIM", "wait status=%s sims=%zu timeout=%us",
                event->status, sel.count, timeout);

    for (;;) {
        for (i = 0; i < sel.count; i++) {
            ue_t *ue;
            ulab_error_t query_err;
            char actual[ULAB_MAX_REF];

            if (matched[i]) {
                continue;
            }

            ue = &ctx->world->ues[sel.idx[i]];
            memset(&query_err, 0, sizeof(query_err));
            memset(actual, 0, sizeof(actual));
            ulab_copy(last_sim, sizeof(last_sim), ue->bff_id);

            if (bff_get_sim_status(ctx->bff, ue, actual,
                                   sizeof(actual), &query_err)) {
                ulab_copy(last_error, sizeof(last_error), query_err.msg);
                continue;
            }

            last_error[0] = '\0';
            ulab_copy(last_status, sizeof(last_status), actual);
            if (strcasecmp(actual, event->status) == 0) {
                matched[i] = 1;
                matched_count++;
                ulab_status("SIM", "%s status=%s", ue->bff_id, actual);
            }
        }

        if (matched_count == sel.count) {
            free(matched);
            selector_result_free(&sel);
            return ULAB_OK;
        }
        if (elapsed >= timeout) {
            break;
        }

        sleep(poll);
        if (timeout - elapsed < poll) {
            elapsed = timeout;
        } else {
            elapsed += poll;
        }
    }

    if (last_error[0] != '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status expected=%.64s matched=%zu/%zu "
                 "sim=%.192s timeout=%us error=%.512s",
                 event->status, matched_count, sel.count, last_sim,
                 timeout, last_error);
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_sim_status expected=%.64s matched=%zu/%zu "
                 "sim=%.192s actual=%.64s timeout=%us",
                 event->status, matched_count, sel.count, last_sim,
                 last_status, timeout);
    }

    free(matched);
    selector_result_free(&sel);
    return ULAB_ERR;
}

static int event_promote_release(event_ctx_t *ctx,
                                 const event_spec_t *event,
                                 ulab_error_t *err) {
    if (event->app[0] == '\0' || event->tag[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "promote_release requires app and tag");
        return ULAB_ERR;
    }

    ulab_status("RELEASE", "promote app=%s version=%s",
                event->app, event->tag);

    if (bff_promote_release(ctx->bff, event->app, "app",
                            event->tag, err)) {
        return ULAB_ERR;
    }

    ulab_status("RELEASE", "promoted app=%s desired=%s",
                event->app, event->tag);
    return ULAB_OK;
}

int event_bff(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err) {
    switch (event->type) {
    case EVT_CREATE_UES:
        (void)ctx;
        snprintf(err->msg, sizeof(err->msg),
                 "create_ues is defined but not enabled in v1.0 runtime path");
        return ULAB_ERR;

    case EVT_ADD_PACKAGE_TO_SIM:
        return event_add_package_to_sim(ctx, event, err);

    case EVT_PURCHASE_PACKAGE:
        return event_purchase_package(ctx, event, err);

    case EVT_SET_PACKAGE_ACTIVE:
        return event_set_package_active(ctx, event, err);

    case EVT_REMOVE_PACKAGE_FROM_SIM:
        return event_remove_package_from_sim(ctx, event, err);

    case EVT_SET_SIM_STATUS:
        return event_set_sim_status(ctx, event, err);

    case EVT_WAIT_SIM_STATUS:
        return event_wait_sim_status(ctx, event, err);

    case EVT_PROMOTE_RELEASE:
        return event_promote_release(ctx, event, err);

    default:
        snprintf(err->msg, sizeof(err->msg), "not a BFF event");
        return ULAB_ERR;
    }
}
