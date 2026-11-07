/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "service_api.h"
#include "service_client.h"
#include "usys_file.h"
#include "usys_services.h"

#define RESPONSE_LIMIT    (1024 * 1024)
#define HTTP_CONNECT_MSEC 1000L
#define HTTP_REQUEST_MSEC 2000L
#define HTTP_JSON_HEADER  "Content-Type: application/json"

typedef struct {
    char *data;
    size_t size;
} Response;

static pthread_once_t curlOnce = PTHREAD_ONCE_INIT;
static CURLcode curlStatus = CURLE_FAILED_INIT;

static void init_curl(void) {

    curlStatus = curl_global_init(CURL_GLOBAL_DEFAULT);
}

static size_t receive_body(char *data, size_t size, size_t count, void *arg) {

    Response *response;
    size_t bytes;
    char *buffer;

    response = arg;

    if (size != 0 && count > RESPONSE_LIMIT / size) {
        return 0;
    }

    bytes = size * count;

    if (bytes > RESPONSE_LIMIT - response->size) {
        return 0;
    }

    buffer = realloc(response->data, response->size + bytes + 1);

    if (buffer == NULL) {
        return 0;
    }

    response->data = buffer;
    memcpy(buffer + response->size, data, bytes);

    response->size += bytes;
    buffer[response->size] = '\0';

    return bytes;
}

int service_http(const char *method, const char *url, json_t *body,
                 json_t **reply) {

    CURL *curl;
    struct curl_slist *headers = NULL;
    Response response          = {0};
    char *text                 = NULL;
    long code                  = HTTP_TRANSPORT_ERROR;

    if (reply != NULL) {
        *reply = NULL;
    }

    pthread_once(&curlOnce, init_curl);

    if (curlStatus != CURLE_OK) {
        return HTTP_TRANSPORT_ERROR;
    }

    curl = curl_easy_init();

    if (curl == NULL) {
        return HTTP_TRANSPORT_ERROR;
    }

    if (body != NULL) {
        text = json_dumps(body, JSON_COMPACT);

        if (text == NULL) {
            goto done;
        }

        headers = curl_slist_append(NULL, HTTP_JSON_HEADER);

        if (headers == NULL) {
            goto done;
        }

        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, text);
    }

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS, HTTP_CONNECT_MSEC);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, HTTP_REQUEST_MSEC);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, receive_body);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

    if (curl_easy_perform(curl) == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);

        if (reply != NULL && response.data != NULL) {
            *reply = json_loadb(response.data, response.size,
                                JSON_REJECT_DUPLICATES, NULL);
        }
    }

done:
    free(response.data);
    free(text);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    return (int)code;
}

bool service_peers_init(ServicePeer peers[SERVICE_MEMBER_COUNT]) {

    const char *names[]    = {ENV_NETWORK_URL, ENV_PCRF_URL, ENV_EPC_URL};
    char *services[]       = {SERVICE_INIT_NETWORK, SERVICE_PCRF,
                              SERVICE_EPCEMU};
    const char *value;
    const char *kind;
    char defaultUrl[SERVICE_URL_LEN];
    int port;
    size_t i;

    kind = getenv(ENV_EPC_KIND);

    if (kind != NULL && strcmp(kind, EPC_KIND_REAL) != 0 &&
        strcmp(kind, EPC_KIND_EMULATED) != 0) {
        return false;
    }

    for (i = 0; i < SERVICE_MEMBER_COUNT; i++) {
        value = getenv(names[i]);

        if (value == NULL) {
            port = usys_find_service_port(services[i]);

            if (port < MIN_SERVICE_PORT || port > MAX_SERVICE_PORT) {
                fprintf(stderr, "Cannot resolve service port: %s\n",
                        services[i]);

                return false;
            }

            snprintf(defaultUrl, sizeof(defaultUrl), DEF_SERVICE_URL_FORMAT,
                     port);
            value = defaultUrl;
        }

        if (strncmp(value, HTTP_SCHEME, sizeof(HTTP_SCHEME) - 1) != 0 ||
            strlen(value) >= sizeof(peers[i].url) - SERVICE_PATH_LEN) {
            return false;
        }

        snprintf(peers[i].url, sizeof(peers[i].url), "%s", value);
        peers[i].emulated =
            i == SERVICE_MEMBER_EPC &&
            (kind == NULL || strcmp(kind, EPC_KIND_EMULATED) == 0);
    }

    /* A real EPC must be selected explicitly with its own API address. */

    return kind == NULL || strcmp(kind, EPC_KIND_REAL) != 0 ||
           getenv(ENV_EPC_URL) != NULL;
}

bool service_peer_set(const ServicePeer *peer, bool enabled) {

    char url[SERVICE_URL_LEN + SERVICE_PATH_LEN];
    json_t *body;
    int code;

    snprintf(url, sizeof(url), "%s%s", peer->url, SERVICE_API_SERVICE);

    body = json_pack("{s:s}", JTAG_STATE,
                     enabled ? SERVICE_STATE_ON : SERVICE_STATE_OFF);

    if (body == NULL) {
        return false;
    }

    code = service_http(HTTP_POST, url, body, NULL);
    json_decref(body);

    return code == HttpStatus_OK || code == HttpStatus_Accepted;
}

ServiceObservation service_peer_read(const ServicePeer *peer) {

    ServiceObservation observed = {0};
    char url[SERVICE_URL_LEN + SERVICE_PATH_LEN];
    json_t *body = NULL;
    json_t *service;
    json_t *count;
    const char *state;
    const char *endpoint;
    int code;

    endpoint = peer->emulated ? SERVICE_API_STATUS : SERVICE_API_SERVICE;
    snprintf(url, sizeof(url), "%s%s", peer->url, endpoint);

    code               = service_http(HTTP_GET, url, NULL, &body);
    observed.reachable = code == HttpStatus_OK && json_is_object(body);

    if (!observed.reachable) {
        goto done;
    }

    service = peer->emulated ? json_object_get(body, JTAG_SERVICE) : body;
    state   = json_string_value(json_object_get(service, JTAG_STATE));

    if (state == NULL || (strcmp(state, SERVICE_STATE_ON) != 0 &&
                          strcmp(state, SERVICE_STATE_OFF) != 0)) {
        observed.reachable = false;
        goto done;
    }

    observed.on    = strcmp(state, SERVICE_STATE_ON) == 0;
    observed.ready = json_is_true(json_object_get(body, JTAG_READY));

    if (peer->emulated) {
        observed.ready = observed.ready &&
                         json_is_true(json_object_get(
                             json_object_get(body, JTAG_PCRF), JTAG_READY));
    }

    count             = json_object_get(service, JTAG_PDN_COUNT);
    observed.pdnCount = json_is_integer(count) ? json_integer_value(count) : 0;
    observed.confirmed =
        !json_is_true(json_object_get(service, JTAG_CLEANUP_PENDING));

    if (json_object_get(service, JTAG_CONFIRMED) != NULL) {
        observed.confirmed =
            observed.confirmed &&
            json_is_true(json_object_get(service, JTAG_CONFIRMED));
    } else if (!json_is_integer(count)) {
        observed.confirmed = false;
    }

    if (observed.pdnCount < 0) {
        observed.confirmed = false;
    }

    if (!observed.on && observed.pdnCount != 0) {
        observed.confirmed = false;
    }

done:
    json_decref(body);

    return observed;
}
