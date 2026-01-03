/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "httpStatus.h"
#include "initClient.h"
#include "jserdes.h"
#include "config.h"
#include "log.h"

struct Response {
	char *buffer;
	size_t size;
};

static char *create_url(char *systemName);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
								void *userp);
static long send_request_to_initClient(char *url, struct Response *response);
static int process_response_from_initClient(char *response,
											char **host,
                                            int *port);

static char *create_url(char *systemName) {

	char *url=NULL;
    char *initSystemHost=NULL;
    char *initSystemPort=NULL;
    char *orgName=NULL;

    initSystemHost = getenv(ENV_INIT_SYSTEM_ADDR);
    initSystemPort = getenv(ENV_INIT_SYSTEM_PORT);
    orgName        = getenv(ENV_SYSTEM_ORG);

    if (initSystemHost == NULL ||
        initSystemPort == NULL ||
        orgName        == NULL ||
        systemName     == NULL) return;

	url = (char *)calloc(MAX_URL_LEN, sizeof(char));
	if (url) {
		sprintf(url, "http://%s:%s/v1/orgs/%s/systems/%s",
                initSystemHost,
                initSystemPort,
                orgName,
                systemName);
	}

	return url;
}

static size_t response_callback(void *contents, size_t size, size_t nmemb,
								void *userp) {

	size_t realsize = size * nmemb;
	struct Response *response = (struct Response *)userp;

	response->buffer = realloc(response->buffer, response->size + realsize + 1);

	if (response->buffer == NULL) {
		log_error("Not enough memory to realloc of size: %d",
				  response->size + realsize + 1);
		return 0;
	}

	memcpy(&(response->buffer[response->size]), contents, realsize);
	response->size += realsize;
	response->buffer[response->size] = 0; /* Null terminate. */
  
	return realsize;
}

static int process_response_from_initClient(char *response,
											char **host,
                                            int *port) {

	int ret=FALSE;
	json_t *json=NULL;
	SystemInfo *systemInfo=NULL;

	if (response == NULL) return FALSE;

	json = json_loads(response, JSON_DECODE_ANY, NULL);

	if (!json) {
		log_error("Can not load str into JSON object. Str: %s", response);
		goto done;
	}

	ret = deserialize_system_info(&systemInfo, json);

	if (ret==FALSE) {
		log_error("Deserialization failed for response: %s", response);
		goto done;
	}

	*host = strdup(systemInfo->ip);
	*port = systemInfo->port;
	ret = TRUE;
	
 done:
	json_decref(json);
	free_system_info(systemInfo);
	return ret;
}

static long send_request_to_initClient(char *url, struct Response *response) {

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

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)response);
	curl_easy_setopt(curl, CURLOPT_USERAGENT, "mesh/0.1");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to initClient at URL %s: %s", url,
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

int get_systemInfo_from_initClient(char *systemName,
                                   char **systemHost,
                                   int *systemPort) {

	int ret=FALSE;
	char *url=NULL;
	struct Response response;

    *systemHost = NULL;
    *systemPort = 0;

	if (systemName == NULL) return FALSE;

    url = create_url(systemName);
    if (url == NULL) return FALSE;

	if (send_request_to_initClient(url, &response) == HttpStatus_OK) {
		if (process_response_from_initClient(response.buffer,
                                             systemHost, systemPort)) {
			log_debug("Recevied info from initClient: host %s port %d",
					  *systemHost, *systemPort);
		} else {
			log_error("Unable to receive info from init");
			goto done;
		}
	} else {
		log_error("Unable to send request to init");
		goto done;
	}

	ret = TRUE;
 done:
	if (url)             free(url);
	if (response.buffer) free(response.buffer);

	return ret;
}

/*
 * 1) Ask init system for "inventory" system credentials (host + port)
 * 2) Call inventory verify endpoint:
 *      http://{ip}:{port}/v1/orgs/{ORGNAME}/components/verify/{node_id}
 *    (ORGNAME comes from ENV_SYSTEM_ORG)
 * 3) Return TRUE only if HTTP 200, otherwise FALSE (including if unreachable).
 */
int verify_nodeid_with_inventory_system(char *nodeID) {

    int  ret      = FALSE;
    char *invHost = NULL;
    int  invPort  = 0;
    char *orgName = NULL;
    char *url     = NULL;

    struct Response response;
    memset(&response, 0, sizeof(response));

    if (nodeID == NULL || nodeID[0] == '\0') {
        return FALSE;
    }

    /* Step 1: get inventory system credentials from init. */
    if (get_systemInfo_from_initClient("inventory", &invHost, &invPort) == FALSE) {
        log_error("Unable to get inventory system info from init system");
        return FALSE; /* includes unreachable/init failure */
    }

    orgName = getenv(ENV_SYSTEM_ORG);
    if (orgName == NULL || orgName[0] == '\0') {
        log_error("Missing env var %s", ENV_SYSTEM_ORG);
        goto done;
    }

    /* Step 2: build verify endpoint URL. */
    url = (char *)calloc(MAX_URL_LEN, sizeof(char));
    if (url == NULL) {
        log_error("Out of memory allocating inventory verify URL");
        goto done;
    }

    snprintf(url, MAX_URL_LEN,
             "http://%s:%d/v1/components/verify/%s",
             invHost, invPort, nodeID);

    /* Step 3: call inventory verify EP; only 200 => TRUE. */
    if (send_request_to_initClient(url, &response) == HttpStatus_OK) {
        ret = TRUE;
        log_debug("Inventory verify OK for nodeID=%s", nodeID);
    } else {
        /* Any failure => FALSE (including unreachable). */
        ret = FALSE;
        log_error("Inventory verify failed for nodeID=%s (url=%s)", nodeID, url);
    }

done:
    if (url) free(url);
    if (response.buffer) free(response.buffer);
    if (invHost) free(invHost);

    return ret;
}

void free_system_info(SystemInfo *systemInfo) {

	if (systemInfo == NULL) return;

	free(systemInfo->systemName);
    free(systemInfo->systemId);
    free(systemInfo->certificate);
    free(systemInfo->ip);

	free(systemInfo);
}
