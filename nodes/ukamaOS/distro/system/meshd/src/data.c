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
#include <pthread.h>

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

static pthread_once_t curl_once = PTHREAD_ONCE_INIT;

STATIC void free_umap(UMap *map) {
    if (!map) return;

    for (int i = 0; i < map->nb_values; i++) {
        SAFE_FREE(map->keys[i]);
        SAFE_FREE(map->values[i]);
    }

    SAFE_FREE(map->keys);
    SAFE_FREE(map->values);
    SAFE_FREE(map->lengths);

    free(map);
}

STATIC void curl_init_once(void) {
    curl_global_init(CURL_GLOBAL_ALL);
}

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
    char *start, *separator;

    if (input[0] != '/') {
        *name = strdup("");
        *ep   = strdup("");
        return;
    }

    start     = &input[1];
    separator = strchr(start, '/');

    if (separator != NULL) {
        size_t name_len = separator - start;
        if (name_len > 0) {
            *name = (char *)calloc(name_len + 1, sizeof(char));
            strncpy(*name, start, name_len);
        } else {
            *name = strdup("");
        }

        const char *ep_start = separator + 1;
        if (*ep_start != '\0') {
            *ep = strdup(ep_start);
        } else {
            *ep = strdup("");
        }
    } else {
        if (strlen(start) > 0) {
            *name = strdup(start);
        } else {
            *name = strdup("");
        }
        *ep = strdup("");
    }

    /* metrics is a special case - prom ep is hard coded */
    if (strcasecmp(*name, SERVICE_METRICS) == 0) {
        free(*ep);
        *ep = strdup(*name);
    }
}

STATIC int send_data_to_local_service(URequest *data,
                                      char *hostname,
                                      int *httpStatus,
                                      char **retStr) {

    CURL *curl = NULL;
    CURLcode res;
    struct curl_slist *headers = NULL;
    char *serviceName = NULL;
    char *serviceEP   = NULL;
    int servicePort = 0;

    char url[MAX_BUFFER] = {0};
    Response response = {NULL, 0};

    if (!httpStatus || !retStr) return FALSE;
    *retStr = NULL;

    if (!data || !hostname) {
        *httpStatus = HttpStatus_InternalServerError;
        return FALSE;
    }

    pthread_once(&curl_once, curl_init_once);

    curl = curl_easy_init();
    if (!curl) {
        *httpStatus = HttpStatus_InternalServerError;
        return FALSE;
    }

    /* Build headers properly: each entry must be "Key: Value" */
    if (data->map_header && data->map_header->nb_values > 0) {
        UMap *map = data->map_header;
        for (int i = 0; i < map->nb_values; i++) {
            const char *k = map->keys[i];
            const char *v = map->values[i];
            if (!k || !v) continue;

            const char *use_v = v;
            if (strcasecmp(k, "Host") == 0) {
                use_v = hostname;
            }

            char line[512];
            snprintf(line, sizeof(line), "%s: %s", k, use_v);
            headers = curl_slist_append(headers, line);
        }
    }

    find_service_name_and_ep(data->http_url, &serviceName, &serviceEP);
    if (!serviceName || !serviceEP) {
        usys_log_error("Unable to extract service name/EP. input=%s", data->http_url);
        *httpStatus = HttpStatus_InternalServerError;
        goto cleanup;
    }

    servicePort = usys_find_service_port(serviceName);
    if (servicePort <= 0) {
        usys_log_error("Unable to find service in /etc/services: %s", serviceName);
        *httpStatus = HttpStatus_ServiceUnavailable;
        *retStr     = strdup(HttpStatusStr(*httpStatus));
        goto cleanup;
    }

    snprintf(url, sizeof(url), "http://localhost:%d/%s", servicePort, serviceEP);

    curl_easy_setopt(curl, CURLOPT_URL, url);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, data->http_verb ? data->http_verb : "GET");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response);

    if (data->binary_body_length > 0 && data->binary_body) {
        /* binary-safe */
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data->binary_body);
        curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, (long)data->binary_body_length);
    }

    res = curl_easy_perform(curl);

    if (res != CURLE_OK) {
        usys_log_error("Error sending request to %s Error: %s", url, curl_easy_strerror(res));
        *httpStatus = HttpStatus_ServiceUnavailable;
        *retStr     = strdup(HttpStatusStr(*httpStatus));
    } else {
        long code = 0;
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
        *httpStatus = (int)code;

        if (response.size && response.buffer) {
            usys_log_debug("Response received from server: %s", response.buffer);
            *retStr = strdup(response.buffer);
        } else {
            *retStr = strdup("");
        }
    }

cleanup:
    SAFE_FREE(response.buffer);
    SAFE_FREE(serviceName);
    SAFE_FREE(serviceEP);

    if (headers) curl_slist_free_all(headers);
    if (curl) curl_easy_cleanup(curl);

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
    SAFE_FREE(item->data); /* free old one, if any*/
    item->data = strdup(message->data);
    item->code = message->code;

    pthread_cond_broadcast(&item->hasResp);
}

int process_incoming_websocket_message(Message *message, Config *config) {

    int httpStatus = HttpStatus_InternalServerError;
    int ret;
    char *responseLocal  = NULL;
    char *responseRemote = NULL;
    URequest *request = NULL;

    if (!message || !config || !message->data) {
        usys_log_error("process_incoming_websocket_message: invalid args");
        return FALSE;
    }

    if (deserialize_request_info(&request, message->data) == FALSE) {
        usys_log_error("Unable to deserialize request on websocket");
        return FALSE;
    }

    ret = send_data_to_local_service(request,
                                     config->localHostname,
                                     &httpStatus,
                                     &responseLocal);

    if (ret) {
        usys_log_debug("Received response from local server Code: %d Response: %s",
                       httpStatus, responseLocal ? responseLocal : "(null)");

        serialize_local_service_response(&responseRemote,
                                         message,
                                         httpStatus,
                                         responseLocal ? (int)strlen(responseLocal) : 0,
                                         responseLocal ? responseLocal : "");
    } else {
        usys_log_error("Error sending message to local service. Error: %d", httpStatus);

        serialize_local_service_response(&responseRemote,
                                         message,
                                         httpStatus,
                                         (int)strlen(HttpStatusStr(httpStatus)),
                                         (char *)HttpStatusStr(httpStatus));
    }

    if (responseRemote) {
        usys_log_debug("Adding response to websocket queue: %s", responseRemote);
        add_work_to_queue(&Transmit, responseRemote, NULL, 0, NULL, 0);
    }

    SAFE_FREE(responseLocal);
    SAFE_FREE(responseRemote);

    if (request) {
        free_umap(request->map_url);
        free_umap(request->map_header);
        free_umap(request->map_post_body);
        free_umap(request->map_cookie);

        SAFE_FREE(request->http_protocol);
        SAFE_FREE(request->http_verb);
        SAFE_FREE(request->http_url);
        SAFE_FREE(request->url_path);
        SAFE_FREE(request->binary_body);

        free(request);
        request = NULL;
    }

    return TRUE;
}
