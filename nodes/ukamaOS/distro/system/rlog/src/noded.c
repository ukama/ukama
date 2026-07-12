/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "jserdes.h"
#include "nodeInfo.h"
#include "rlogd.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#define NODE_INFO_EP "/v1/nodeinfo"
#define NODED_CONNECT_TIMEOUT_MS 1000L
#define NODED_REQUEST_TIMEOUT_MS 2000L

struct Response {
    char *buffer;
    size_t size;
};

static char *create_noded_url(const char *host, int port) {
    char *url;
    int length;

    if (!host || !*host || port <= 0) return NULL;

    url = calloc(MAX_URL_LEN, 1);
    if (!url) return NULL;

    length = snprintf(url, MAX_URL_LEN, "http://%s:%d%s",
                      host, port, NODE_INFO_EP);
    if (length < 0 || length >= MAX_URL_LEN) {
        free(url);
        return NULL;
    }

    return url;
}

static size_t response_callback(void *contents, size_t size,
                                size_t count, void *userData) {
    struct Response *response;
    char *buffer;
    size_t bytes;

    response = userData;
    bytes = size * count;
    if (bytes == 0) return 0;

    buffer = realloc(response->buffer,
                     response->size + bytes + 1U);
    if (!buffer) {
        usys_log_error("Unable to grow noded response buffer");
        return 0;
    }

    response->buffer = buffer;
    memcpy(response->buffer + response->size, contents, bytes);
    response->size += bytes;
    response->buffer[response->size] = '\0';
    return bytes;
}

static int process_response(const char *response, char **uuid) {
    json_error_t error;
    json_t *json;
    NodeInfo *nodeInfo;
    int rc;

    if (!response || !uuid) return USYS_FALSE;

    nodeInfo = NULL;
    json = json_loads(response, 0, &error);
    if (!json) {
        usys_log_error("Unable to parse noded response: %s", error.text);
        return USYS_FALSE;
    }

    rc = deserialize_node_info(&nodeInfo, json);
    json_decref(json);
    if (rc == USYS_FALSE || !nodeInfo || !nodeInfo->uuid) {
        free_node_info(nodeInfo);
        return USYS_FALSE;
    }

    *uuid = strdup(nodeInfo->uuid);
    free_node_info(nodeInfo);
    return *uuid ? USYS_TRUE : USYS_FALSE;
}

static long send_request(const char *url, struct Response *response) {
    struct curl_slist *headers;
    CURL *curl;
    CURLcode result;
    long status;

    if (!url || !response) return 0;

    status = 0;
    headers = NULL;
    curl = curl_easy_init();
    if (!curl) return 0;

    headers = curl_slist_append(headers, "Accept: application/json");
    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "rlog.d/1");
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT_MS,
                     NODED_CONNECT_TIMEOUT_MS);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT_MS,
                     NODED_REQUEST_TIMEOUT_MS);

    result = curl_easy_perform(curl);
    if (result == CURLE_OK) {
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &status);
    } else {
        usys_log_warn("Unable to query noded: %s",
                      curl_easy_strerror(result));
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    return status;
}

int get_nodeID_from_noded(char **nodeId, char *host, int port) {
    struct Response response;
    char *url;
    int rc;

    if (!nodeId || !host || port <= 0) return USYS_FALSE;

    *nodeId = NULL;
    memset(&response, 0, sizeof(response));
    url = create_noded_url(host, port);
    if (!url) return USYS_FALSE;

    if (curl_global_init(CURL_GLOBAL_DEFAULT) != CURLE_OK) {
        free(url);
        return USYS_FALSE;
    }

    rc = USYS_FALSE;
    if (send_request(url, &response) == 200 &&
        process_response(response.buffer, nodeId) == USYS_TRUE) {
        rc = USYS_TRUE;
    }

    curl_global_cleanup();
    free(response.buffer);
    free(url);
    return rc;
}

void free_node_info(NodeInfo *nodeInfo) {
    if (!nodeInfo) return;

    usys_free(nodeInfo->uuid);
    usys_free(nodeInfo->name);
    usys_free(nodeInfo->partNumber);
    usys_free(nodeInfo->skew);
    usys_free(nodeInfo->mac);
    usys_free(nodeInfo->assemblyDate);
    usys_free(nodeInfo->oem);
    usys_free(nodeInfo);
}
