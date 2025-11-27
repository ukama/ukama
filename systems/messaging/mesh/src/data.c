/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

/* Functions related to the handling of received data on websockets. */

#include <mesh.h>
#include <curl/curl.h>
#include <curl/easy.h>
#include <string.h>
#include <jansson.h>

#include "mesh.h"
#include "work.h"
#include "data.h"
#include "jserdes.h"
#include "initClient.h"
#include "httpStatus.h"

typedef struct _response {
	char *buffer;
	size_t size;
} Response;

extern WorkList *Transmit; /* global */

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

static int get_substring_after_index(char **out,
                                     const char *str,
                                     int index,
                                     char delimiter) {
    int count = 0;

    while (*str) {
        if (*str == delimiter) {
            count++;
            if (count == index) {
                *out = strdup(str + 1); // Skip the delimiter itself
                return 1;
            }
        }
        str++;
    }

    return 0;
}

static int extract_system_path(char *str, char **name, char **path) {

    char *ptr = NULL;
    int len = 0;
    
    // Check if the string starts with "http://"
    if (strncmp(str, "http://", 7) == 0) {
        // Skip the "http://" part and move past the domain/IP and port
        str = strchr(str + 7, '/');
        if (!str) return FALSE; // No path found after domain/IP
    }

    // Extract the path and name
    if (!get_substring_after_index(&ptr, str, 2, '/')) return FALSE;

    len = strlen(str) - strlen(ptr) - 2; /* -2 to skip the /s */
    *name = (char *)calloc(1, len + 1);
    strncpy(*name, str + 1, len);

    *path = strdup(ptr);
    free(ptr);

    return TRUE;
}

void clear_request(MRequest **data) {

	free((*data)->reqType);
	free((*data)->nodeInfo);
	free((*data)->serviceInfo);

	ulfius_clean_request_full((*data)->requestInfo);

	free(*data);
}

static int send_data_to_system(URequest *data, char *ep,
                               char *ip, int port,
                               int *retCode, char **retStr) {

	int i;
	CURL *curl=NULL;
	CURLcode res;
	struct curl_slist *headers=NULL;
	char url[MAX_BUFFER] = {0};
	UMap *map;
	Response response = {NULL, 0};
    long responseCode;

	*retCode = 0;

	/* Sanity check */
	if (data == NULL && ip == NULL && !port) {
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

	sprintf(url, "http://%s:%d/%s", ip, port, ep);
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
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &responseCode);
        *retCode = (int)responseCode;
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

static URequest* create_http_request(char *jStr) {

    URequest *request;
	json_t *json, *jMethod, *jURL, *jPath, *jRaw, *obj, *jData;

	if (jStr == NULL) return FALSE;

    json = json_loads(jStr, JSON_DECODE_ANY, NULL);
    if (json == NULL) return FALSE;

    jMethod = json_object_get(json, JSON_METHOD);
    jURL    = json_object_get(json, JSON_URL);
    jPath   = json_object_get(json, JSON_PATH);
	jRaw    = json_object_get(json, JSON_RAW_DATA);

    if (jMethod == NULL || jURL == NULL || jPath == NULL || jRaw == NULL) {
        log_error("Missing json parameter in recevied requests");
        json_decref(json);
        return NULL;
    }

	request = (URequest *)calloc(1, sizeof(URequest));
	if (request == NULL) {
        json_decref(json);
        log_error("Error allocating memory of size: %lu", sizeof(URequest));
		return NULL;
    }

    if (ulfius_init_request(request)) {
        log_error("Error initializing new http request.");
        json_decref(json);
        return NULL;
    }

    ulfius_set_request_properties(request,
                                  U_OPT_HTTP_VERB, json_string_value(jMethod),
                                  U_OPT_HTTP_URL,  json_string_value(jURL),
                                  U_OPT_HEADER_PARAMETER, "User-Agent", "mesh",
                                  U_OPT_TIMEOUT, 20,
                                  U_OPT_NONE);
    if (jPath) {
        request->url_path = strdup(json_string_value(jPath));
        if (request->url_path == NULL) {
            log_error("Error allocating memory for URL path");
            ulfius_clean_request(request);
            free(request);
            json_decref(json);
            return NULL;
        }
    }

    /* Get the actual data now */
    jData = json_object_get(jRaw, JSON_DATA);
    if (jData) {
        const char *str = json_string_value(jData);
        request->binary_body        = strdup(str);
        request->binary_body_length = strlen(str);
    }

    json_decref(json);

	return request;
}

int process_incoming_websocket_message(Message *message, char **responseRemote){

    /*
     * 1. Find system info from init
     * 2. create thread, make connection with system, send/recv
     * 3. Put the response back on the websocket via outgoing queue
     */
	int retCode=0;
	URequest *request;
	char *responseLocal=NULL, *jStr=NULL;
    char *systemName=NULL, *systemEP=NULL;
	char *systemHost=NULL;
    int systemPort=0;

    log_debug("Recevied message from mesh-host: %s", message->data);

    request = create_http_request(message->data);
    if (request == NULL) {
        log_error("Unable to deser the request on websocket");
        retCode = HttpStatus_BadRequest;
    }

    if (!extract_system_path(request->url_path, &systemName, &systemEP)) {
        log_error("Unable to extract system name and path: %s",
                  request->url_path);
        retCode = HttpStatus_BadRequest;
        responseLocal = HttpStatusStr(retCode);
    }

	if (!get_systemInfo_from_initClient(systemName, &systemHost, &systemPort)) {
		/* No match. Ignore. */
		log_error("No matching server found for system: %s", systemName);
        retCode = HttpStatus_InternalServerError;
        responseLocal = HttpStatusStr(retCode);
	} else {
    
        log_debug("Matching server found for system: %s host: %s port: %d",
                  systemName, systemHost, systemPort);

        send_data_to_system(request, systemEP,
                            systemHost, systemPort,
                            &retCode, &responseLocal);
        log_debug("Return code from system %s:%d: code: %d Response: %s",
                  systemHost, systemPort, retCode, responseLocal);
    }

    serialize_system_response(responseRemote, message, retCode,
                              strlen(responseLocal), responseLocal);
    log_debug("Sending response back: %s", *responseRemote);

    ulfius_clean_request(request);
    if (request)       free(request);
    if (responseLocal) free(responseLocal);
    if (systemName)    free(systemName);
    if (systemEP)      free(systemEP);
    if (systemHost)    free(systemHost);

    return retCode;
}
