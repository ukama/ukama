/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to node.d */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <jansson.h>
#include <curl/curl.h>

#include "nodeInfo.h"
#include "jserdes.h"
#include "log.h"

#define NODED_PATH    "noded/v1"
#define NODE_INFO_EP  "nodeinfo"

struct Response {
	char *buffer;
	size_t size;
};

/* Function def. */
static char *create_noded_url(char *host, char *port);
static size_t response_callback(void *contents, size_t size, size_t nmemb,
								void *userp);
static long send_request_to_noded(char *nodedURL, struct Response *response);
static int process_response_from_noded(char *response,	char **uuid);

/*
 * create_noded_url --
 *
 */
static char *create_noded_url(char *host, char *port) {

	char *url=NULL;

	if (host == NULL || port == NULL) return NULL;

	url = (char *)malloc(MAX_URL_LEN);
	if (url) {
		sprintf(url, "%s:%s/%s/%s", host, port, NODED_PATH, NODE_INFO_EP);
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
  
	if(response->buffer == NULL) {
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
 * process_response_from_noded --
 *
 */
static int process_response_from_noded(char *response,	char **uuid) {

	int ret=FALSE;
	json_t *json=NULL;
	NodeInfo *nodeInfo=NULL;

	if (response == NULL) return FALSE;

	json = json_loads(response, JSON_DECODE_ANY, NULL);

	if (!json) {
		log_error("Can not load str into JSON object. Str: %s", response);
		goto done;
	}

	ret = deserialize_node_info(&nodeInfo, json);

	if (ret==FALSE) {
		log_error("Deserialization failed for response: %s", response);
		goto done;
	}

	*uuid = strdup(nodeInfo->uuid);
	ret = TRUE;
	
 done:
	json_decref(json);
	free_node_info(nodeInfo);
	return ret;
}

/*
 * send_request_to_noded --
 *
 */
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
	curl_easy_setopt(curl, CURLOPT_USERAGENT, "bootstrap/0.1");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to node.d at URL %s: %s", nodedURL,
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
 * get_nodeID_from_noded -- get NodeID
 *
 */
int get_nodeID_from_noded(char **nodeID, char *host, char *port) {

	int ret=FALSE;
	char *nodedURL=NULL;
	struct Response response;

	if (host == NULL || port == NULL) return FALSE;

	*nodeID = NULL;
	nodedURL = create_noded_url(host, port);

	if (send_request_to_noded(nodedURL, &response) == 200) {
		if (process_response_from_noded(response.buffer, nodeID)) {
			log_debug("Recevied NodeID (UUID) from noded: %s", *nodeID);
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
	if (nodedURL) free(nodedURL);
	if (response.buffer) free(response.buffer);
	
	return ret;
}

/*
 * free_node_info --
 *
 */
void free_node_info(NodeInfo *nodeInfo) {

	NodeInfo *ptr;

	if (nodeInfo == NULL) return;

	ptr = nodeInfo;
	
	if (ptr->uuid)          free(ptr->uuid);
	if (ptr->name)          free(ptr->name);
	if (ptr->partNumber)    free(ptr->partNumber);
	if (ptr->skew)          free(ptr->skew);
	if (ptr->mac)           free(ptr->mac);
	if (ptr->assemblyDate)  free(ptr->assemblyDate);
	if (ptr->oem)           free(ptr->oem);
		
	free(nodeInfo);
	nodeInfo=NULL;
}
