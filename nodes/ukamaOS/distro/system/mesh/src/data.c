/**
 * Copyright (c) 2021-present, Ukama Inc.
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
#include "map.h"
#include "work.h"
#include "data.h"
#include "jserdes.h"

typedef struct _response {
	char *buffer;
	size_t size;
} Response;

extern WorkList *Transmit; /* global */
extern MapTable *ClientTable;

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
 * clear_request -- free up memory from MRequest.
 *
 */
void clear_request(MRequest **data) {

	free((*data)->reqType);
	free((*data)->deviceInfo);
	free((*data)->serviceInfo);

	ulfius_clean_request_full((*data)->requestInfo);

	free(*data);
}

/*
 * send_data_to_server -- Forward recevied data to the local server.
 *
 */
static long send_data_to_local_service(URequest *data, char *hostname,
                                       char *port, int *retCode,
                                       char **retStr) {
  
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
	if (data == NULL && hostname == NULL && port == NULL) {
		return code;
	}
     
	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) {
		return code;
	}

	/* Add to the header if exists. */
	if (data->map_header) {
		map = data->map_header;
		for (i=0; i < map->nb_values; i++) {
			headers = curl_slist_append(headers, map->keys[i]);
			headers = curl_slist_append(headers,": ");
            if (strcmp(map->keys[i], "Host") == 0) {
                headers = curl_slist_append(headers, hostname);
            } else {
                headers = curl_slist_append(headers, map->values[i]);
            }
		}
	}

	sprintf(url, "http://%s:%d/%s", hostname, atoi(port)+1, data->http_url);
	curl_easy_setopt(curl, CURLOPT_URL, url);

	curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, data->http_verb);
	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

	if (data->binary_body_length > 0 && data->binary_body) {
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data->binary_body);
	}

	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);
    //	curl_easy_setopt(curl, CURLOPT_USERAGENT, "mesh/0.1");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to server at %s Error: %s",
				  url, curl_easy_strerror(res));
		*retStr = strdup("Target service is not available. Try again");
	} else {
		/* get status code. */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
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

	return code;
}

/*
 * process_incoming_websocket_response --
 *
 */
void process_incoming_websocket_response(Message *message, void *data) {

    MapItem *item=NULL;

    log_debug("Response to local service Code: %d Data: %s",
              message->code, message->data);

    item = is_existing_item(ClientTable,message->serviceInfo->name,
                            message->serviceInfo->port);
    if (item == NULL) {
        log_error("Matching client not found. Port: %s. Ignoring",
                  message->serviceInfo->port);
        return;
    }

    item->size = message->dataSize;
    item->data = strdup(message->data);
    item->code = message->code;

    pthread_cond_broadcast(&item->hasResp);
}

/*
 * process_incoming_websocket_message --
 *
 */
int process_incoming_websocket_message(Message *message, Config *config) {

    /*
     * 1. replace the "nodeID" with hostname (e.g., localhost)
     * 2. create thread,  make connection, send msg, recv response.
     * 3. put the response back on the websocket outgoing queue.
     */

	int retCode=0, ret;
	char *responseLocal=NULL, *responseRemote=NULL;
    json_t *jResp=NULL;
    URequest *request=NULL;

    if (deserialize_request_info(&request, message->data) == FALSE) {
        log_error("Unable to deser the request on websocket");
        return FALSE;
    }

    ret = send_data_to_local_service(request,
                                     config->localHostname,
                                     message->nodeInfo->port,
                                     &retCode, &responseLocal);

    if (ret == 200) {
        log_debug("Command success. CURL return code: %d. Return code: %d",
                  ret, retCode);
    } else {
        log_debug("Command failed. CURL return code: %d. Return code: %d",
                  ret, retCode);
    }

    /* Convert the response into proper format and return. */
    serialize_local_service_response(&responseRemote, message,
                                     retCode,
                                     strlen(responseLocal),
                                     responseLocal);

    log_debug("Sending response back: %s", responseRemote);
    add_work_to_queue(&Transmit, responseRemote, NULL, 0, NULL, 0);

    return TRUE;
}
