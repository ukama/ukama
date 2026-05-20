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

#include <curl/curl.h>

#include "epcemu.h"
#include "http_client.h"

#define HTTP_TIMEOUT_SEC 5

typedef struct {
    char *data;
    size_t size;
} ResponseBuf;

static size_t write_cb(void *contents, size_t size, size_t nmemb,
                       void *userp) {

    size_t realSize;
    ResponseBuf *mem;
    char *ptr;

    realSize = size * nmemb;
    mem = (ResponseBuf *)userp;

    ptr = realloc(mem->data, mem->size + realSize + 1);
    if (ptr == NULL) return 0;

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), contents, realSize);
    mem->size += realSize;
    mem->data[mem->size] = '\0';

    return realSize;
}

static int parse_json_body(ResponseBuf *buf, JsonObj **outJson) {

    json_error_t error;

    if (outJson == NULL) return USYS_TRUE;

    *outJson = NULL;

    if (buf == NULL || buf->data == NULL || buf->size == 0) {
        return USYS_TRUE;
    }

    *outJson = json_loads(buf->data, 0, &error);
    if (*outJson == NULL) {
        usys_log_error("failed to parse JSON response: %s", error.text);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int http_send_json(const char *method, const char *url, JsonObj *body,
                   JsonObj **outJson, long *httpCode) {

    CURL *curl;
    CURLcode res;
    struct curl_slist *headers;
    ResponseBuf response;
    char *payload;
    int ret;

    if (method == NULL || url == NULL) return USYS_FALSE;

    response.data = NULL;
    response.size = 0;
    headers = NULL;
    payload = NULL;
    ret = USYS_FALSE;

    curl = curl_easy_init();
    if (curl == NULL) return USYS_FALSE;

    if (body != NULL) {
        payload = json_dumps(body, JSON_COMPACT);
        if (payload == NULL) goto done;
    }

    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "Accept: application/json");

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, HTTP_TIMEOUT_SEC);

    if (!strcmp(method, "POST")) {
        curl_easy_setopt(curl, CURLOPT_POST, 1L);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload ? payload : "{}");
    } else if (!strcmp(method, "DELETE")) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload ? payload : "{}");
    } else if (strcmp(method, "GET")) {
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, method);
        if (payload != NULL) curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload);
    }

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("curl %s %s failed: %s", method, url,
                       curl_easy_strerror(res));
        goto done;
    }

    if (httpCode != NULL) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, httpCode);
    }

    if (!parse_json_body(&response, outJson)) goto done;

    ret = USYS_TRUE;

done:
    if (headers != NULL) curl_slist_free_all(headers);
    if (payload != NULL) free(payload);
    if (response.data != NULL) free(response.data);
    curl_easy_cleanup(curl);

    return ret;
}

int http_get_json(const char *url, JsonObj **outJson, long *httpCode) {

    return http_send_json("GET", url, NULL, outJson, httpCode);
}
