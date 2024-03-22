/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#include "nodeInfo.h"
#include "jserdes.h"

#define NODED_PATH    "noded/v1"
#define NODE_INFO_EP  "nodeinfo"

struct Response {
    char *buffer;
    size_t size;
};

static char *create_noded_url(char *host, int port);
static size_t response_callback(void *contents, size_t size, size_t nmemb, void *userp);
static long send_request_to_noded(char *nodedURL, struct Response *response);
static int process_response_from_noded(char *response,	char **uuid);

static char *create_noded_url(char *host, int port) {

    char *url=NULL;

    if (host == NULL || port == 0) return NULL;

    url = (char *)malloc(MAX_URL_LEN);
    if (url) {
        sprintf(url, "%s:%d/%s/%s", host, port, NODED_PATH, NODE_INFO_EP);
    }

    return url;
}

static size_t response_callback(void *contents, size_t size, size_t nmemb,
								void *userp) {

    size_t realsize = size * nmemb;
    struct Response *response = (struct Response *)userp;

    response->buffer = realloc(response->buffer, response->size + realsize + 1);
  
    if(response->buffer == NULL) {
        usys_log_error("Not enough memory to realloc of size: %d",
                       response->size + realsize + 1);
        return 0;
    }

    memcpy(&(response->buffer[response->size]), contents, realsize);
    response->size += realsize;
    response->buffer[response->size] = 0;
  
    return realsize;
}

static int process_response_from_noded(char *response,	char **uuid) {

    int ret=USYS_FALSE;
    json_t *json=NULL;
    NodeInfo *nodeInfo=NULL;

    if (response == NULL) return USYS_FALSE;

    json = json_loads(response, JSON_DECODE_ANY, NULL);

    if (!json) {
        usys_log_error("Can not load str into JSON object. Str: %s", response);
        goto done;
    }

    ret = deserialize_node_info(&nodeInfo, json);
    if (ret == USYS_FALSE) {
        usys_log_error("Deserialization failed for response: %s", response);
        goto done;
    }

    *uuid = strdup(nodeInfo->uuid);
    ret = USYS_TRUE;
	
done:
    json_decref(json);
    free_node_info(nodeInfo);
    return ret;
}

static long send_request_to_noded(char *nodedURL, struct Response *response) {

    long resCode=0;
    CURL *curl=NULL;
    CURLcode res;
    struct curl_slist *headers=NULL;

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (curl == NULL) {
        return resCode;
    }

    response->buffer = malloc(1);
    response->size   = 0;
  
    /* Add to the header. */
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, "charset: utf-8");

    curl_easy_setopt(curl, CURLOPT_URL, nodedURL);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)response);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "rlog.d/0.1");

    res = curl_easy_perform(curl);

    if (res != CURLE_OK) {
        usys_log_error("Error sending request to node.d at URL %s: %s", nodedURL,
                       curl_easy_strerror(res));
    } else {
        /* get status code. */
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    }

    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    curl_global_cleanup();

    return resCode;
}

int get_nodeID_from_noded(char **nodeID, char *host, int port) {

    int ret=USYS_FALSE;
    char *nodedURL=NULL;
    struct Response response;

    if (host == NULL || port == 0) return USYS_FALSE;

    *nodeID = NULL;
    nodedURL = create_noded_url(host, port);

    if (send_request_to_noded(nodedURL, &response) == 200) {
        if (process_response_from_noded(response.buffer, nodeID)) {
            usys_log_debug("Recevied NodeID (UUID) from noded: %s", *nodeID);
        } else {
            usys_log_error("Unable to receive proper NodeID from noded");
            goto done;
        }
    } else {
        usys_log_error("Unable to send request to noded");
        goto done;
    }

    ret = USYS_TRUE;

done:
    if (nodedURL) free(nodedURL);
    if (response.buffer) free(response.buffer);

    return ret;
}

void free_node_info(NodeInfo *nodeInfo) {

    NodeInfo *ptr = NULL;

    if (nodeInfo == NULL) return;

    ptr = nodeInfo;
	
    usys_free(ptr->uuid);
    usys_free(ptr->name);
    usys_free(ptr->partNumber);
    usys_free(ptr->skew);
    usys_free(ptr->mac);
    usys_free(ptr->assemblyDate);
    usys_free(ptr->oem);
    
    usys_free(nodeInfo);
}
