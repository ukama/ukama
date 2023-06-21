/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to the handling of received data on websockets. */

#include <mesh.h>
#include <curl/curl.h>
#include <curl/easy.h>
#include <string.h>

#include "mesh.h"
#include "work.h"
#include "data.h"
#include "jserdes.h"
#include "initClient.h"

typedef struct _response {
	char *buffer;
	size_t size;
} Response;

extern WorkList *Transmit; /* global */

/*
 * response_callback --
 */
static size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

	size_t realsize = size * nmemb;
	Response *response = (Response *)userp;

	response->buffer = realloc(response->buffer, response->size + realsize + 1);

	if(response->buffer == NULL) {
		log_error("Not enough memory to realloc of size: %s",
				  response->size + realsize + 1);
		return 0;
	}

	memcpy(&(response->buffer[response->size]), contents, realsize);
	response->size += realsize;
	response->buffer[response->size] = 0; /* Null terminate. */

	return realsize;
}

/*
 * extract_system_path --
 */
static int extract_system_path(char *str, char **name, char **path) {

    char *ptr=NULL;
    int len=0;

    if (!get_substring_after_index(&ptr, str, 2, '/')) return FALSE;

    len = strlen(str) - strlen(ptr) - 2; /* -2 to skip the /s */
    *name = (char *)calloc(1, len+1);
    strncpy(*name, str+1, len);

    *path = strdup(ptr);

    return TRUE;
}

/*
 * clear_request -- free up memory from MRequest.
 *
 */
void clear_request(MRequest **data) {

	free((*data)->reqType);
	free((*data)->nodeInfo);
	free((*data)->serviceInfo);

	ulfius_clean_request_full((*data)->requestInfo);

	free(*data);
}

/*
 * send_data_to_system -- Forward recevied data to the system
 *
 */
static long send_data_to_system(URequest *data, char *ip, char *port,
								int *retCode, char **retStr) {
  
	int i;
	long code=0;
	CURL *curl=NULL;
	CURLcode res;
	struct curl_slist *headers=NULL;
	char url[MAX_BUFFER] = {0};
	UMap *map;
	Response response = {NULL, 0};

	*retCode = 0;

	/* Sanity check */
	if (data == NULL && ip == NULL && port == NULL) {
		return FALSE;
	}

	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) return FALSE;

	/* Add to the header if exists. */
	if (data->map_header) {
		map = data->map_header;
		for (i=0; i < map->nb_values; i++) {
			headers = curl_slist_append(headers, map->keys[i]);
			headers = curl_slist_append(headers,": ");
			headers = curl_slist_append(headers, map->values[i]);
		}
	}

	sprintf(url, "http://%s:%s/%s", ip, port, data->http_url);
	curl_easy_setopt(curl, CURLOPT_URL, url);

	curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, data->http_verb);
	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

	if (data->binary_body_length > 0 && data->binary_body) {
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data->binary_body);
	}

	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
	curl_easy_setopt(curl, CURLOPT_USERAGENT, "mesh");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to server at %s Error: %s",
				  url, curl_easy_strerror(res));
		*retStr = strdup("Target service is not available. Try again.");
        *retCode = 0;
	} else {
		/* get status code. */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &retCode);
		if (response.size) {
			log_debug("Response recevied from server: %s", response.buffer);
			*retStr = strdup(response.buffer);
		}
	}

	if (response.buffer) 
		free(response.buffer);

	curl_slist_free_all(headers);
	curl_easy_cleanup(curl);
	curl_global_cleanup();

	return TRUE;
}

/*
 * process_incoming_websocket_message --
 *
 */
int process_incoming_websocket_message(Message *message, char **responseRemote){

    /*
     * 1. Find system info from init
     * 2. create thread, make connection with system, send/recv
     * 3. Put the response back on the websocket via outgoing queue
     */
	int ret=FALSE, retCode=0;
	URequest *request;
	char *responseLocal=NULL, *jStr=NULL;
    char *systemName=NULL, *systemEP=NULL;
	char *systemHost=NULL, *systemPort=NULL;
	json_t *jResp=NULL;

    if (strcmp(message->reqType, MESH_NODE_REQUEST) != 0) {
        log_error("Invalid request type. ignoring.");
        return FALSE;
    }

    if (deserialize_request_info(&request, message->data) == FALSE) {
        log_error("Unable to deser the request on websocket");
        return FALSE;
    }

    if (!extract_system_path(request->url_path, &systemName, &systemEP)) {
        log_error("Unable to extract system name and path: %s",
                  request->url_path);
        return FALSE;
    }

	if (!get_systemInfo_from_initClient(systemName, &systemHost, &systemPort)) {
		/* No match. Ignore. */
		log_error("No matching server found for system: %s", systemName);
        return FALSE;
	}

    log_debug("Matching server found for system: %s host: %s port: %s",
              systemName, systemHost, systemPort);

    ret = send_data_to_system(request,
                              systemHost, systemPort,
                              &retCode, &responseLocal);
    log_debug("Return code from system %s:%s: code: %d Response: %s",
              systemHost, systemPort, retCode, responseLocal);

    serialize_system_response(responseRemote, message,
                              strlen(responseLocal), responseLocal);

    if (responseRemote) {
        log_debug("Sending response back: %s", *responseRemote);
    } else {
        log_error("Invalid response type (expected JSON)");
        return FALSE;
    }

    return TRUE;
}
