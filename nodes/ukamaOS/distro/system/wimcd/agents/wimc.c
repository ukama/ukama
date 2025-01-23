/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/* Functions related to wimc. */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "wimc.h"
#include "err.h"
#include "http_status.h"
#include "common/utils.h"
#include "agent/jserdes.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_api.h"
#include "usys_mem.h"

#define AGENT_CB_EP "app"
#define WIMC_EP     "v1/agents"

struct Response {
    char *buffer;
    size_t size;
};

static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

    size_t realsize = size * nmemb;
    struct Response *response = (struct Response *)userp;

    response->buffer = realloc(response->buffer, response->size + realsize + 1);

    if(response->buffer == NULL) {
        log_error("Not enough memory to realloc of size: %s",
                  response->size + realsize + 1);
        return 0;
    }

    memcpy(&(response->buffer[response->size]), contents, realsize);
    response->size += realsize;
    response->buffer[response->size] = 0;

    return realsize;
}

static long send_request_to_wimc(char *wimcURL, json_t *json) {

    long code=0;
    char *jsonStr = NULL;
    CURL *curl = NULL;
    CURLcode res;
    
    struct curl_slist *headers=NULL;
    struct Response response;
  
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) return 0;

    response.buffer = malloc(1);
    response.size   = 0;
    if (json) jsonStr = json_dumps(json, 0);
  
    /* Add to the header. */
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, wimcURL);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "PUT");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    if (jsonStr) curl_easy_setopt(curl, CURLOPT_POSTFIELDS, jsonStr);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

    curl_easy_setopt(curl, CURLOPT_USERAGENT, "agent/0.1");

    res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        usys_log_error("Error sending request to WIMC at URL %s: %s",
                       wimcURL,
                       curl_easy_strerror(res));
    } else {
        /* get status code. */
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    }

    usys_free(jsonStr);
    usys_free(response.buffer);

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();
    
    return code;
}

long communicate_with_wimc(int reqType,
                           char *wimcURL,
                           char *cbURL,
                           void *data) {

    long code=0;
    json_t   *json=NULL;
    TStats   *stats=NULL;

    if (reqType == REQUEST_UPDATE) {
        stats = (TStats *)data;
    }

    code = send_request_to_wimc(wimcURL, json);
    if (code == HttpStatus_OK) {
        usys_log_debug("Communication with wimc: %s code: %d", 
                       wimcURL, code);
    } else {
        usys_log_error("Communication with WIMC %s: failed. Code: %d",
                       wimcURL, code);
    }

    if (json) json_decref(json);
    
    return code;
}
