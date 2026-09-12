/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <limits.h>
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

    if (!buffer || bytes == 0) {
        return bytes;
    }
    if (buffer->length + bytes >= RESPONSE_LIMIT) {
        return 0;
    }

    grown = realloc(buffer->data, buffer->length + bytes + 1);
    if (!grown) {
        return 0;
    }

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

    if (!method || !url || timeoutSec <= 0 || !status) {
        return false;
    }

    curl = curl_easy_init();
    if (!curl) {
        return false;
    }

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

    if (!json_is_object(object) || !key) {
        return NULL;
    }

    value = json_object_get(object, key);
    return json_is_string(value) ? json_string_value(value) : NULL;
}

static void copy_text(char *dst, size_t size, const char *src) {

    if (!dst || size == 0) {
        return;
    }
    snprintf(dst, size, "%s", src ? src : "");
}

static bool parse_starter_status(const char *body, StarterSnapshot *snapshot) {

    json_t *root;
    json_t *starter;
    json_t *readiness;
    const char *state;
    bool valid;

    root = json_loads(body, JSON_REJECT_DUPLICATES, NULL);
    if (root == NULL) {
        return false;
    }

    starter = json_object_get(root, "starterd");
    readiness = json_object_get(starter, "readiness");
    state = json_string_or_null(readiness, "state");
    valid = starter_aggregate_parse(state, &snapshot->aggregate);
    copy_text(snapshot->aggregateReason, sizeof(snapshot->aggregateReason),
               json_string_or_null(readiness, "reason"));
    json_decref(root);
    return valid;
}

static bool read_config_text(json_t *root, const char *key,
                              char *value, size_t size) {

    json_t *entry = json_object_get(root, key);
    const char *text;

    if (!json_is_string(entry)) {
        return false;
    }
    text = json_string_value(entry);
    if (json_string_length(entry) >= size ||
        json_string_length(entry) != strlen(text)) {
        return false;
    }
    snprintf(value, size, "%s", text);
    return true;
}

static bool parse_config_status(const char *body, ConfigSnapshot *snapshot) {

    json_t *root;
    json_t *generation;
    json_t *revision;
    json_t *version;
    char phase[16];
    bool valid = false;

    root = json_loads(body, JSON_REJECT_DUPLICATES, NULL);
    if (root == NULL) {
        return false;
    }
    version = json_object_get(root, "schemaVersion");
    generation = json_object_get(root, "generation");
    revision = json_object_get(root, "revision");
    if (!json_is_integer(version) || json_integer_value(version) != 1 ||
        !json_is_integer(generation) || json_integer_value(generation) < 0 ||
        !json_is_integer(revision) || json_integer_value(revision) < 0 ||
        json_integer_value(revision) > INT_MAX) {
        goto done;
    }
    if (!read_config_text(root, "phase", phase, sizeof(phase)) ||
        !read_config_text(root, "mode", snapshot->mode, sizeof(snapshot->mode)) ||
        !read_config_text(root, "requestId", snapshot->requestId, sizeof(snapshot->requestId)) ||
        !read_config_text(root, "error", snapshot->error, sizeof(snapshot->error))) {
        goto done;
    }
    snapshot->generation = json_integer_value(generation);
    snapshot->revision = json_integer_value(revision);

    if (strcmp(phase, "failed") == 0) {
        snapshot->phase = CONFIG_PHASE_FAILED;
        valid = true;
        goto done;
    }
    if (snapshot->error[0] != '\0') {
        goto done;
    }
    if (strcmp(snapshot->mode, "NONE") == 0) {
        snapshot->phase = CONFIG_PHASE_AWAITING;
        valid = strcmp(phase, "awaiting") == 0 && snapshot->revision == 0 &&
            ((snapshot->generation == 0 && snapshot->requestId[0] == '\0') ||
             (snapshot->generation > 0 && snapshot->requestId[0] != '\0'));
        goto done;
    }
    if (snapshot->generation == 0 || snapshot->requestId[0] == '\0') {
        goto done;
    }
    if (strcmp(snapshot->mode, "NOCONFIG") == 0) {
        if (snapshot->revision != 0) {
            goto done;
        }
    } else if (strcmp(snapshot->mode, "CONFIG") != 0 || snapshot->revision == 0) {
        goto done;
    }
    if (strcmp(phase, "pending") == 0) {
        snapshot->phase = CONFIG_PHASE_IN_PROGRESS;
        valid = true;
    } else if (strcmp(phase, "completed") == 0) {
        snapshot->phase = CONFIG_PHASE_APPLIED;
        valid = true;
    }

done:
    json_decref(root);
    return valid;
}

bool config_client_get_status(const Config *config, ConfigSnapshot *snapshot) {

    char url[URL_LEN];
    char *body = NULL;
    long status = 0;
    bool received;

    memset(snapshot, 0, sizeof(*snapshot));
    snapshot->phase = CONFIG_PHASE_UNKNOWN;
    snprintf(url, sizeof(url), "http://%s:%d/v1/config/status",
             config->configHost, config->configPort);

    received = http_request("GET", url, NULL, config->requestTimeoutSec,
                            &status, &body);
    if (received && status == HttpStatus_OK && body != NULL) {
        snapshot->available = parse_config_status(body, snapshot);
    }
    free(body);
    return snapshot->available;
}

bool starter_client_get_status(const Config *config,
                               StarterSnapshot *snapshot) {

    char url[URL_LEN];
    char *body;
    long status;
    bool ok;

    if (!config || !snapshot) {
        return false;
    }

    memset(snapshot, 0, sizeof(*snapshot));
    snapshot->aggregate = STARTER_AGGREGATE_UNKNOWN;

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
        body == NULL || !parse_starter_status(body, snapshot)) {
        free(body);
        return false;
    }

    snapshot->available = true;
    free(body);
    return true;
}

static char *event_details(const LifecycleEvent *event) {

    json_t *json;
    char *text;

    json = json_pack("{s:i,s:s,s:I,s:s,s:s,s:I,s:s}",
                     "schemaVersion", 1,
                     "bootId", event->bootId,
                     "sequence", (json_int_t)event->sequence,
                     "requestId", event->requestId,
                     "configMode", event->configMode,
                     "configGeneration", (json_int_t)event->configGeneration,
                     "reason", event->reason);
    if (json == NULL) {
        return NULL;
    }
    text = json_dumps(json, JSON_COMPACT);
    json_decref(json);
    return text;
}

bool notify_client_send_event(const Config *config,
                              const LifecycleEvent *event) {

    char url[URL_LEN];
    char *body;
    json_t *json;
    char *details;
    const char *value;
    long status;
    bool ok;

    if (!config || !event) {
        return false;
    }

    snprintf(url,
             sizeof(url),
             "http://%s:%d/v1/event/%s",
             config->notifyHost,
             config->notifyPort,
             SERVICE_LIFECYCLE);

    details = event_details(event);
    if (details == NULL) {
        return false;
    }
    json = json_object();
    if (json == NULL) {
        free(details);
        return false;
    }
    value = lifecycle_state_str(event->state);
    if (event->state == LIFECYCLE_STATE_STARTING) {
        value = "INIT";
    }

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
                        json_string(value));
    json_object_set_new(json, "units", json_string(""));
    json_object_set_new(json,
                        "details",
                        json_string(details));
    free(details);

    body = json_dumps(json, JSON_COMPACT);
    json_decref(json);
    if (!body) {
        return false;
    }

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
