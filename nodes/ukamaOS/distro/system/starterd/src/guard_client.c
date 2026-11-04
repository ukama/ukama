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
#include "guard_client.h"

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

int guard_http(const char *method, const char *url, json_t *body,
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
