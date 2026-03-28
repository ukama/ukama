/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include <curl/curl.h>

#include "notify.h"

int notify_send_alarm(const EmuConfig *cfg,
                      const EmuAlarm *alarm,
                      const EmuSwitchInfo *info) {
    CURL *curl       = NULL;
    CURLcode rc      = CURLE_OK;
    char url[256]    = {0};
    char payload[512] = {0};
    struct curl_slist *headers = NULL;

    if (cfg == NULL || alarm == NULL || info == NULL || !cfg->notifyEnabled) {
        return STATUS_OK;
    }

    snprintf(url, sizeof(url), "http://%s:%d%s",
             cfg->notifyHost, cfg->notifyPort, cfg->notifyPath);
    snprintf(payload, sizeof(payload),
             "{\"source\":\"%s\",\"severity\":\"%s\","
             "\"code\":%d,\"message\":\"%s\",\"serial\":\"%s\"}",
             alarm->source, alarm->severity, alarm->code,
             alarm->message, info->serial);

    curl = curl_easy_init();
    if (curl == NULL) {
        return STATUS_NOK;
    }

    headers = curl_slist_append(headers, "Content-Type: application/json");

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_POST, 1L);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 2L);

    rc = curl_easy_perform(curl);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);

    return (rc == CURLE_OK) ? STATUS_OK : STATUS_NOK;
}
