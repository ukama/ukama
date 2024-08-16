/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
#include "httpStatus.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_services.h"
#include "usys_log.h"

#include "static.h"

typedef struct _response {
	char *buffer;
	size_t size;
} Response;

extern WorkList *Transmit; /* global */
extern MapTable *ClientTable;


STATIC size_t response_callback(void *contents, size_t size, size_t nmemb,
                                void *userp) {

	size_t realsize = size * nmemb;
	Response *response = (Response *)userp;

	response->buffer = realloc(response->buffer, response->size + realsize + 1);

	if(response->buffer == NULL) {
		usys_log_error("Not enough memory to realloc of size: %s",
				  response->size + realsize + 1);
		return 0;
	}

	memcpy(&(response->buffer[response->size]), contents, realsize);
	response->size += realsize;
	response->buffer[response->size] = 0; /* Null terminate. */

	return realsize;
}

void clear_request(MRequest **data) {

	free((*data)->reqType);
	free((*data)->deviceInfo);
	free((*data)->serviceInfo);

	ulfius_clean_request_full((*data)->requestInfo);

	free(*data);
}

STATIC void find_service_name_and_ep(char *input,
                                     char **name,
                                     char **ep) {

    const char *separator = strchr(&input[1], '/');

    if (separator != NULL) {
        size_t length = separator - input;

        *name = (char *)calloc(length, sizeof(char));
        strncpy(*name, &input[1], length-1);
        *ep = strdup(separator);
    } else {
        *name = strdup(input);
        *ep   = strdup("");
    }
}

STATIC int send_data_to_local_service(URequest *data,
                                      char *hostname,
                                      int *httpStatus,
                                      char **retStr) {
  
	CURL *curl=NULL;
	CURLcode res;
	struct curl_slist *headers=NULL;
    char *serviceName = NULL;
    char *serviceEP   = NULL;
    int  servicePort = 0, i;

	char url[MAX_BUFFER] = {0};
	UMap *map = NULL;
	Response response = {NULL, 0};

	if (data == NULL && hostname == NULL) {
		return FALSE;
	}
     
	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) {
        *httpStatus = HttpStatus_InternalServerError;
		return FALSE;
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

    find_service_name_and_ep(data->http_url, &serviceName, &serviceEP);
    if (serviceName == NULL || serviceEP == NULL) {
        usys_log_error("Unable to extract service namd and EP. input",
                  data->http_url);
        *httpStatus = HttpStatus_InternalServerError;
        return FALSE;
    }

    servicePort = usys_find_service_port(serviceName);
    if (servicePort <= 0) {
        usys_log_error("Unable to find service name in /etc/services: %s",
                  serviceName);
        *httpStatus = HttpStatus_ServiceUnavailable;
        return FALSE;
    }

    sprintf(url,
            "http://localhost:%d/%s",
            servicePort,
            serviceEP);

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, data->http_verb);
	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

	if (data->binary_body_length > 0 && data->binary_body) {
		curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data->binary_body);
	}

	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		usys_log_error("Error sending request to server at %s Error: %s",
				  url, curl_easy_strerror(res));
        *httpStatus = HttpStatus_ServiceUnavailable;
        *retStr     = strdup(HttpStatusStr(*httpStatus));
	} else {
		/* get status code. */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, httpStatus);
		if (response.size) {
			usys_log_debug("Response recevied from server: %s", response.buffer);
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

void process_incoming_websocket_response(Message *message, void *data) {

    MapItem *item=NULL;

    usys_log_debug("Response to local service Code: %d Data: %s",
              message->code, message->data);

    item = is_existing_item(ClientTable, message->seqNo);
    if (item == NULL) {
        usys_log_error("Matching client not found. Port: %s. Ignoring",
                  message->serviceInfo->port);
        return;
    }

    item->size = message->dataSize;
    item->data = strdup(message->data);
    item->code = message->code;

    pthread_cond_broadcast(&item->hasResp);
}

int process_incoming_websocket_message(Message *message, Config *config) {

	int httpStatus, ret;
	char *responseLocal  = NULL;
    char *responseRemote = NULL;
    json_t *jResp=NULL;
    URequest *request=NULL;

    if (deserialize_request_info(&request, message->data) == FALSE) {
        usys_log_error("Unable to deser the request on websocket");
        return FALSE;
    }

    ret = send_data_to_local_service(request,
                                     config->localHostname,
                                     &httpStatus,
                                     &responseLocal);

    if (ret) {
        usys_log_debug("Recevied response from local servier Code: %d Response: %s",
                  httpStatus, responseLocal);

        /* Convert the response into proper format and return. */
        serialize_local_service_response(&responseRemote,
                                         message,
                                         httpStatus,
                                         strlen(responseLocal),
                                         responseLocal);
    } else {
        usys_log_error("Error sending message to local service. Error: %d",
                  httpStatus);

        /* Convert the response into proper format and return. */
        serialize_local_service_response(&responseRemote,
                                         message,
                                         httpStatus,
                                         strlen(HttpStatusStr(httpStatus)),
                                         HttpStatusStr(httpStatus));
    }

    usys_log_debug("Adding response to the websocket queue: %s", responseRemote);
    add_work_to_queue(&Transmit, responseRemote, NULL, 0, NULL, 0);

    return TRUE;
}
