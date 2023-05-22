/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to calling initClient */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "initClient.h"
#include "jserdes.h"
#include "config.h"
#include "log.h"

#define EP_SYSTEMS "systems"
#define QUERY_KEY  "name"

struct Response {
	char *buffer;
	size_t size;
};

/* Function def. */
static char *create_url(char *host, char *port, char *systemName);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
								void *userp);
static long send_request_to_initClient(char *url, struct Response *response);
static int process_response_from_initClient(char *response,
											char **host, char **port);

/*
 * create_url --
 *
 */
static char *create_url(char *host, char *port, char *systemName) {

	char *url=NULL;

	if (host == NULL || port == NULL || systemName == NULL) return NULL;

	url = (char *)malloc(MAX_URL_LEN);
	if (url) {
		sprintf(url, "%s:%s/%s?%s=%s", host, port, EP_SYSTEMS,
				QUERY_KEY, systemName);
	}

	return url;
}

/*
 * response_callback --
 *
 */
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

/*
 * process_response_from_initClient --
 *
 */
static int process_response_from_initClient(char *response,
											char **host, char **port) {

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
	*port = strdup(systemInfo->port);
	ret = TRUE;
	
 done:
	json_decref(json);
	free_system_info(systemInfo);
	return ret;
}

/*
 * send_request_to_initClient --
 *
 */
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

/*
 *
 * get_systemInfo_from_initClient --
 *
 */
int get_systemInfo_from_initClient(Config *config, char *systemName,
								   char **host, char **port) {

	int ret=FALSE;
	char *url=NULL;
	struct Response response;

	if (systemName == NULL) return FALSE;

	*host = NULL;
	*port = NULL;

	url = create_url(config->initClientHost, config->initClientPort,
					 systemName);

	if (send_request_to_initClient(url, &response) == 200) {
		if (process_response_from_initClient(response.buffer, host, port)) {
			log_debug("Recevied info from initClient: host %s port %s",
					  *host, *port);
		} else {
			log_error("Unable to receive proper NodeID from noded");
			goto done;
		}
	} else {
		log_error("Unable to send request to noded");
		goto done;
	}

	ret = TRUE;
 done:
	if (url)             free(url);
	if (response.buffer) free(response.buffer);

	return ret;
}

/*
 * free_system_info --
 *
 */
void free_system_info(SystemInfo *systemInfo) {

	if (systemInfo == NULL) return;

	if (systemInfo->systemName)  free(systemInfo->systemName);
	if (systemInfo->systemId)    free(systemInfo->systemId);
	if (systemInfo->certificate) free(systemInfo->certificate);
	if (systemInfo->ip)          free(systemInfo->ip);
	if (systemInfo->port)        free(systemInfo->port);
	if (systemInfo->health)      free(systemInfo->health);
	
	free(systemInfo);
}
