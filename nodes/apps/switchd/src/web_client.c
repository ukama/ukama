/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <string.h>

#include "jserdes.h"
#include "web_client.h"

#include "usys_log.h"

struct CurlBuffer {
    char *buf;
    size_t len;
    size_t cap;
};

static size_t curl_write_cb(char *ptr,
                            size_t size,
                            size_t nmemb,
                            void *userdata) {
    struct CurlBuffer *cb;
    size_t total;
    size_t copyLen;

    cb = (struct CurlBuffer *)userdata;
    total = size * nmemb;
    copyLen = (cb->len + total < cb->cap - 1) ? total : (cb->cap - 1 - cb->len);

    if (copyLen > 0) {
        memcpy(cb->buf + cb->len, ptr, copyLen);
        cb->len += copyLen;
        cb->buf[cb->len] = '\0';
    }

    return total;
}

int web_client_post_json(const char *url,
                         const char *json,
                         long timeoutMs,
                         long *status,
                         char *response,
                         size_t responseLen) {
    CURL *curl;
    CURLcode rc;
    struct curl_slist *headers;
    struct CurlBuffer cb;

    if (url == NULL || json == NULL || response == NULL || responseLen == 0) {
        return SWITCHD_ERR_INVAL;
    }

    if (status) {
        *status = 0;
    }
    response[0] = '\0';

    curl = curl_easy_init();
    if (curl == NULL) {
        return SWITCHD_ERR_INTERNAL;
    }

    headers = NULL;
    cb.buf = response;
    cb.len = 0;
    cb.cap = responseLen;

    headers = curl_slist_append(headers, "Content-Type: application/json");
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, json);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS, timeoutMs);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, curl_write_cb);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &cb);

    rc = curl_easy_perform(curl);
    if (rc != CURLE_OK) {
        usys_log_error("HTTP POST failed for notify.d: %s",
                       curl_easy_strerror(rc));
        curl_slist_free_all(headers);
        curl_easy_cleanup(curl);
        return SWITCHD_ERR_IO;
    }

    if (status) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, status);
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    return SWITCHD_OK;
}

int web_client_notify_alarm(const SwitchdConfig *cfg,
                            const SwitchAlarm *alarm,
                            bool clear) {
    JsonObj *json;
    char *payload;
    char response[256];
    long status;
    int ret;

    if (cfg == NULL || alarm == NULL) {
        return SWITCHD_ERR_INVAL;
    }

    json = NULL;
    payload = NULL;
    status = 0;

    if (json_serialize_alarm_notification(&json,
                                          SERVICE_NAME,
                                          alarm,
                                          clear) == false) {
        return SWITCHD_ERR_INTERNAL;
    }

    payload = json_dumps(json, JSON_COMPACT);
    json_free(&json);
    if (payload == NULL) {
        return SWITCHD_ERR_NOMEM;
    }

    ret = web_client_post_json(cfg->notifyUrl,
                               payload,
                               cfg->notifyTimeoutMs,
                               &status,
                               response,
                               sizeof(response));
    free(payload);
    if (ret != SWITCHD_OK) {
        return ret;
    }

    if (status < 200 || status >= 300) {
        usys_log_error("notify.d returned HTTP %ld", status);
        return SWITCHD_ERR_IO;
    }

    return SWITCHD_OK;
}
