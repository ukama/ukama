/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <errno.h>
#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/stat.h>
#include <time.h>
#include <unistd.h>

#include "session.h"

typedef struct {
    char body[8192];
    size_t length;
    long code;
    CURLcode result;
} Response;

static const char *string_value(json_t *body, const char *key) {

    json_t *value;
    const char *text;

    value = json_object_get(body, key);
    text = json_string_value(value);
    if (!text || strlen(text) != json_string_length(value)) return NULL;

    return text;
}

static int status_write(UeSession *session) {

    json_t *body;
    int rc;

    body = json_pack("{s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:I,s:I}",
                     "state", session->state, "imsi", session->imsi,
                     "iccid", session->iccid, "ip", session->ip,
                     "requested_apn", session->requestedApn,
                     "apn", session->apn,
                     "epc_instance_id", session->instanceId,
                     "pdn_id", session->pdnId,
                     "data_state", session->dataState,
                     "reason", session->reason,
                     "source", "ue-agent",
                     "http_status", (json_int_t)session->httpStatus,
                     "updated_at", (json_int_t)time(NULL));
    if (!body) return -1;

    rc = json_dump_file(body, UE_STATUS_FILE ".tmp", JSON_COMPACT);
    json_decref(body);
    if (rc == 0) rc = rename(UE_STATUS_FILE ".tmp", UE_STATUS_FILE);
    if (rc != 0) {
        fprintf(stderr, "cannot write UE status: %s\n", strerror(errno));
    }

    return rc;
}

static void state_set(UeSession *session, const char *state,
                      const char *reason) {

    snprintf(session->state, sizeof(session->state), "%s", state);
    if (reason) {
        snprintf(session->reason, sizeof(session->reason), "%s", reason);
    }
    (void)status_write(session);
}

static size_t receive_body(void *ptr, size_t size, size_t count,
                            void *data) {

    Response *response = data;
    size_t length;

    if (size && count > (sizeof(response->body) - 1) / size) return 0;
    length = size * count;
    if (length >= sizeof(response->body) - response->length) return 0;
    memcpy(response->body + response->length, ptr, length);
    response->length += length;
    response->body[response->length] = '\0';

    return length;
}

static void request(UeSession *session, const char *method,
                     const char *path, json_t *body, Response *response) {

    CURL *curl;
    struct curl_slist *headers;
    char url[1024];
    char *encoded;

    memset(response, 0, sizeof(*response));
    response->result = CURLE_FAILED_INIT;
    encoded = body ? json_dumps(body, JSON_COMPACT) : NULL;
    if (body && !encoded) return;
    curl = curl_easy_init();
    if (!curl) {
        free(encoded);
        return;
    }

    headers = NULL;
    snprintf(url, sizeof(url), "%s%s", session->url, path);
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10L);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, receive_body);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, response);
    if (body) {
        headers = curl_slist_append(NULL, "Content-Type: application/json");
        if (!headers) {
            free(encoded);
            curl_easy_cleanup(curl);
            return;
        }
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, encoded);
    }

    response->result = curl_easy_perform(curl);
    if (response->result == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response->code);
    } else {
        fprintf(stderr, "%s %s: %s\n", method, path,
                curl_easy_strerror(response->result));
    }
    session->httpStatus = response->code;
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    free(encoded);
}

static bool uuid_valid(const char *text) {

    size_t i;

    if (!text || strlen(text) != 36) return false;
    for (i = 0; i < 36; i++) {
        if (i == 8 || i == 13 || i == 18 || i == 23) {
            if (text[i] != '-') return false;
        } else if (!strchr("0123456789abcdef", text[i])) {
            return false;
        }
    }

    return true;
}

static bool matches(json_t *body, const char *key, const char *expected) {

    const char *value = string_value(body, key);

    return value && strcmp(value, expected) == 0;
}

/* Pin identity only after the complete UE tuple has been checked. */
static int binding_read(UeSession *session, Response *response) {

    json_t *body;
    const char *instance;
    const char *id;
    const char *apn;
    const char *state;
    const char *data;
    const char *reason;
    int rc;

    rc = -1;
    body = json_loadb(response->body, response->length,
                      JSON_REJECT_DUPLICATES, NULL);
    instance = string_value(body, "epc_instance_id");
    id = string_value(body, "pdn_id");
    apn = string_value(body, "apn");
    state = string_value(body, "state");
    data = string_value(body, "data_state");
    reason = string_value(body, "reason");

    if (!matches(body, "imsi", session->imsi) ||
        !matches(body, "iccid", session->iccid) ||
        !matches(body, "ip", session->ip) ||
        !uuid_valid(instance) || !uuid_valid(id) ||
        !apn || strlen(apn) >= sizeof(session->apn) ||
        strcasecmp(apn, session->requestedApn) ||
        !state || !data || !reason ||
        strlen(reason) >= sizeof(session->reason)) goto done;

    if (strcmp(state, "attached") && strcmp(state, "attaching") &&
        strcmp(state, "detaching")) goto done;
    if (strcmp(data, "allowed") && strcmp(data, "blocked") &&
        strcmp(data, "unknown")) goto done;
    if (session->pdnId[0] &&
        (strcmp(id, session->pdnId) ||
         strcmp(instance, session->instanceId))) goto done;

    snprintf(session->instanceId, sizeof(session->instanceId), "%s",
             instance);
    snprintf(session->pdnId, sizeof(session->pdnId), "%s", id);
    snprintf(session->apn, sizeof(session->apn), "%s", apn);
    snprintf(session->dataState, sizeof(session->dataState), "%s", data);
    state_set(session, state, reason);
    session->uncertain = false;
    session->probe = true;
    rc = strcmp(state, "attached") == 0 ? 1 : 0;

done:
    json_decref(body);
    return rc;
}

static void observe(UeSession *session, Response *response) {

    char path[512];

    snprintf(path, sizeof(path), "/v1/ue/%s", session->imsi);
    request(session, "GET", path, NULL, response);
}

static void binding_clear(UeSession *session) {

    session->instanceId[0] = '\0';
    session->pdnId[0] = '\0';
    session->apn[0] = '\0';
    snprintf(session->dataState, sizeof(session->dataState), "unknown");
    session->probe = false;
}

int session_init(UeSession *session) {

    if (mkdir(UE_STATUS_DIR, 0755) != 0 && errno != EEXIST) return -1;
    snprintf(session->state, sizeof(session->state), "starting");
    snprintf(session->dataState, sizeof(session->dataState), "unknown");

    return status_write(session);
}

int session_attach(UeSession *session) {

    Response response;
    json_t *body;
    const char *reason;
    int rc;

    if (session->terminal) return 0;
    if (session->probe) {
        observe(session, &response);
        if (response.code == 200) {
            rc = binding_read(session, &response);
            if (rc >= 0) return rc;
            session->terminal = true;
            state_set(session, "error", "binding_mismatch");
            return 0;
        }
        if (response.code != 404 || session->uncertain) {
            state_set(session, "attaching", "attach_outcome_unknown");
            return 0;
        }
        binding_clear(session);
    }

    state_set(session, "attaching", "attach_pending");
    body = json_pack("{s:s,s:s,s:s,s:s,s:s,s:i}",
                     "imsi", session->imsi, "iccid", session->iccid,
                     "ip", session->ip, "apn", session->requestedApn,
                     "userPlaneHost", session->dataHost,
                     "userPlanePort", session->dataPort);
    if (!body) return 0;
    request(session, "POST", "/v1/ue/attach", body, &response);
    json_decref(body);
    if (response.code == 200 || response.code == 201) {
        rc = binding_read(session, &response);
        if (rc >= 0) return rc;
    }
    if (response.code == 400 || response.code == 403) {
        body = json_loadb(response.body, response.length,
                          JSON_REJECT_DUPLICATES, NULL);
        reason = string_value(body, "error");
        session->terminal = true;
        state_set(session, "rejected", reason ? reason : "invalid_response");
        json_decref(body);
        return 0;
    }

    /* A lost POST response is not permission to create another PDN. */
    session->probe = true;
    session->uncertain = response.code == 0 || response.code == 200 ||
                         response.code == 201;
    if (response.result == CURLE_COULDNT_CONNECT ||
        response.result == CURLE_COULDNT_RESOLVE_HOST ||
        response.result == CURLE_COULDNT_RESOLVE_PROXY ||
        response.result == CURLE_FAILED_INIT) {
        session->probe = false;
        session->uncertain = false;
    }
    state_set(session, "attaching", session->uncertain ?
              "attach_outcome_unknown" : "attach_pending");

    return 0;
}

int session_poll(UeSession *session) {

    Response response;
    int rc;

    observe(session, &response);
    if (response.code == 404) {
        binding_clear(session);
        state_set(session, "attaching", "pdn_absent");
        return 0;
    }
    if (response.code != 200) {
        snprintf(session->dataState, sizeof(session->dataState), "unknown");
        state_set(session, "attached", "observation_unavailable");
        return 1;
    }
    rc = binding_read(session, &response);
    if (rc < 0) {
        session->terminal = true;
        state_set(session, "error", "binding_mismatch");
        return 0;
    }

    return rc;
}

int session_detach(UeSession *session) {

    Response response;
    json_t *body;

    if (!strcmp(session->state, "detached")) return 1;
    state_set(session, "detaching", "cleanup_pending");
    if (!session->pdnId[0] && session->probe) {
        observe(session, &response);
        if (response.code == 200) {
            if (binding_read(session, &response) < 0) {
                state_set(session, "error", "binding_mismatch");
                return 0;
            }
        } else if (response.code != 404 || session->uncertain) {
            state_set(session, "detaching", "attach_outcome_unknown");
            return 0;
        }
    }
    if (session->pdnId[0]) {
        body = json_pack("{s:s,s:s,s:s}", "imsi", session->imsi,
                         "epc_instance_id", session->instanceId,
                         "pdn_id", session->pdnId);
        if (!body) return 0;
        request(session, "DELETE", "/v1/ue/detach", body, &response);
        json_decref(body);
        if (response.code != 200) {
            state_set(session, "detaching", response.code == 409 ?
                      "stale_pdn_identity" : "cleanup_pending");
            return 0;
        }
    }

    binding_clear(session);
    session->uncertain = false;
    session->terminal = true;
    state_set(session, "detached", "explicit_detach");

    return 1;
}
