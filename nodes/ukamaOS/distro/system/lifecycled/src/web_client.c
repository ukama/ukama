/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "http_status.h"
#include "web_client.h"

#include "usys_log.h"
#include "usys_services.h"

#define RESPONSE_LIMIT (256 * 1024)
#define URL_LEN        512

typedef struct {
    char *data;
    size_t length;
} ResponseBuffer;

static size_t response_write(void *contents,
                             size_t size,
                             size_t count,
                             void *userData) {

    ResponseBuffer *buffer;
    size_t bytes;
    char *grown;

    buffer = (ResponseBuffer *)userData;
    bytes = size * count;

    if (!buffer || bytes == 0) return bytes;
    if (buffer->length + bytes >= RESPONSE_LIMIT) return 0;

    grown = realloc(buffer->data, buffer->length + bytes + 1);
    if (!grown) return 0;

    buffer->data = grown;
    memcpy(buffer->data + buffer->length, contents, bytes);
    buffer->length += bytes;
    buffer->data[buffer->length] = '\0';
    return bytes;
}

static bool http_request(const char *method,
                         const char *url,
                         const char *body,
                         int timeoutSec,
                         long *status,
                         char **responseBody) {

    CURL *curl;
    CURLcode result;
    struct curl_slist *headers;
    ResponseBuffer response;
    bool ok;

    if (!method || !url || timeoutSec <= 0 || !status) return false;

    curl = curl_easy_init();
    if (!curl) return false;

    memset(&response, 0, sizeof(response));
    headers = NULL;
    ok = false;
    *status = 0;

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, (long)timeoutSec);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_write);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

    if (body) {
        headers = curl_slist_append(headers,
                                    "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body);
    }

    result = curl_easy_perform(curl);
    if (result == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, status);
        ok = true;
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    if (responseBody) {
        *responseBody = response.data;
    } else {
        free(response.data);
    }

    return ok;
}

static const char *json_string_or_null(json_t *object, const char *key) {

    json_t *value;

    if (!json_is_object(object) || !key) return NULL;

    value = json_object_get(object, key);
    return json_is_string(value) ? json_string_value(value) : NULL;
}

static void copy_text(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) return;
    snprintf(dst, size, "%s", src ? src : "");
}

static void parse_config_entry(json_t *entry,
                               StarterSnapshot *snapshot) {

    const char *reason;
    const char *requestId;

    reason = json_string_or_null(entry, "reason");
    requestId = json_string_or_null(entry, "requestId");

    copy_text(snapshot->configReason,
              sizeof(snapshot->configReason),
              reason);
    copy_text(snapshot->configRequestId,
              sizeof(snapshot->configRequestId),
              requestId);

    snapshot->configPhase = config_phase_from_reason(reason);
}

static bool parse_starter_status(const char *body,
                                 StarterSnapshot *snapshot) {

    json_t *root;
    json_t *starter;
    json_t *readiness;
    json_t *apps;
    json_t *entry;
    json_error_t error;
    const char *state;
    const char *reason;
    const char *name;
    size_t index;

    if (!body || !snapshot) return false;

    root = json_loads(body, 0, &error);
    if (!root) return false;

    starter = json_object_get(root, "starterd");
    readiness = json_is_object(starter) ?
        json_object_get(starter, "readiness") : NULL;

    state = json_string_or_null(readiness, "state");
    reason = json_string_or_null(readiness, "reason");

    if (!starter_aggregate_parse(state, &snapshot->aggregate)) {
        json_decref(root);
        return false;
    }

    copy_text(snapshot->aggregateReason,
              sizeof(snapshot->aggregateReason),
              reason);
    snapshot->configPhase = CONFIG_PHASE_ABSENT;

    apps = json_object_get(readiness, "apps");
    if (json_is_array(apps)) {
        json_array_foreach(apps, index, entry) {
            name = json_string_or_null(entry, "name");
            if (name && strcmp(name, SERVICE_CONFIG) == 0) {
                parse_config_entry(entry, snapshot);
                break;
            }
        }
    }

    json_decref(root);
    return true;
}

bool starter_client_get_status(const Config *config,
                               StarterSnapshot *snapshot) {

    char url[URL_LEN];
    char *body;
    long status;
    bool ok;

    if (!config || !snapshot) return false;

    memset(snapshot, 0, sizeof(*snapshot));
    snapshot->aggregate = STARTER_AGGREGATE_UNKNOWN;
    snapshot->configPhase = CONFIG_PHASE_UNKNOWN;

    snprintf(url,
             sizeof(url),
             "http://%s:%d/v1/status",
             config->starterHost,
             config->starterPort);

    body = NULL;
    status = 0;
    ok = http_request("GET",
                      url,
                      NULL,
                      config->requestTimeoutSec,
                      &status,
                      &body);

    if (!ok || status != HttpStatus_OK ||
        !parse_starter_status(body, snapshot)) {
        free(body);
        return false;
    }

    snapshot->available = true;
    free(body);
    return true;
}

bool notify_client_send_event(const Config *config,
                              const LifecycleEvent *event) {

    char url[URL_LEN];
    char *body;
    json_t *json;
    long status;
    bool ok;

    if (!config || !event) return false;

    snprintf(url,
             sizeof(url),
             "http://%s:%d/v1/event/%s",
             config->notifyHost,
             config->notifyPort,
             SERVICE_LIFECYCLE);

    json = json_object();
    if (!json) return false;

    json_object_set_new(json,
                        "service_name",
                        json_string(SERVICE_LIFECYCLE));
    json_object_set_new(json,
                        "severity",
                        json_string(event->state ==
                                    LIFECYCLE_STATE_FAULTY ?
                                    "high" : "low"));
    json_object_set_new(json,
                        "time",
                        json_integer(event->occurredAt));
    json_object_set_new(json, "module", json_string("node"));
    json_object_set_new(json, "name", json_string("state"));
    json_object_set_new(json,
                        "value",
                        json_string(lifecycle_state_str(event->state)));
    json_object_set_new(json, "units", json_string(""));
    json_object_set_new(json,
                        "details",
                        json_string(event->reason));

    body = json_dumps(json, JSON_COMPACT);
    json_decref(json);
    if (!body) return false;

    status = 0;
    ok = http_request("POST",
                      url,
                      body,
                      config->requestTimeoutSec,
                      &status,
                      NULL);
    free(body);

    if (!ok || status != HttpStatus_Accepted) {
        usys_log_warn("notify: event %s sequence %llu not accepted",
                      lifecycle_state_str(event->state),
                      (unsigned long long)event->sequence);
        return false;
    }

    return true;
}
