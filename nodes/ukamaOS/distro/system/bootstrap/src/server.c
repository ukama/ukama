/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <curl/curl.h>
#include <curl/easy.h>
#include <string.h>
#include <stdlib.h>
#include <jansson.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include "usys_log.h"

#include "httpStatus.h"
#include "jserdes.h"
#include "nodeInfo.h"
#include "server.h"

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

static long send_request_to_server(char *url,
                                   Response *response,
                                   char **retStr) {

	long resCode=0;
	CURL *curl=NULL;
	CURLcode res;
	struct curl_slist *headers=NULL;

	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) {
		return resCode;
	}

	headers = curl_slist_append(headers, "Accept: application/json");
	headers = curl_slist_append(headers, "charset: utf-8");

	curl_easy_setopt(curl, CURLOPT_URL, url);
	curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "GET");
	curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
	curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
	curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)response);
	curl_easy_setopt(curl, CURLOPT_USERAGENT, "bootstrap/0.1");

	res = curl_easy_perform(curl);

	if (res != CURLE_OK) {
		log_error("Error sending request to server at URL %s: %s", url,
				  curl_easy_strerror(res));
        *retStr = NULL;
	} else {
		/* get status code */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
        *retStr = strdup(response->buffer);
	}

	curl_slist_free_all(headers);
	curl_easy_cleanup(curl);
	curl_global_cleanup();

	return resCode;
}

static int process_response_from_server(char *response, ServerInfo *server) {

	int ret=FALSE;
	json_t *json=NULL;

	if (response == NULL) return FALSE;

	json = json_loads(response, JSON_DECODE_ANY, NULL);

	if (!json) {
		log_error("Can not load str into JSON object. Str: %s", response);
		goto done;
	}

	ret = deserialize_server_info(server, json);
	if (ret==FALSE) {
		log_error("Deserialization failed for response: %s", response);
		goto done;
	}

	ret = TRUE;

 done:
	json_decref(json);
	return ret;
}

static int send_request_to_init(char *remoteServer, int remotePort,
                                char *uuid, ServerInfo *server,
                                char **responseStr) {

	int ret=FALSE, respCode=0;
	Response response = {NULL, 0};
	char url[MAX_GET_URL_LEN] = {0};

	if (remoteServer == NULL || uuid == NULL) return FALSE;

    /* Create URL + request */
	sprintf(url, "http://%s:%d/%s/%s/%s",
            remoteServer, remotePort, API_VERSION, EP_NODES, uuid);

    respCode = send_request_to_server(&url[0], &response, responseStr);

    switch(respCode) {
    case HttpStatus_OK:
		if (process_response_from_server(response.buffer, server)) {
			log_debug_server(server);
			log_debug("Recevied server IP: %s", server->IP);
		} else {
			log_error("Unable to receive proper server info from bootstrap: %s",
					  remoteServer);
			goto done;
		}
        ret=TRUE;
        break;
    case HttpStatus_NotFound:
        log_debug("NodeID: %s not found on server: %s", uuid, remoteServer);
        ret=-1;
        break;
    default:
        if (respCode) {
            log_error("Error sending request to init %s. Error: %s",
                      url, HttpStatusStr(respCode));
        } else {
            log_error("Error sending request to init %s", url);
        }
        ret=FALSE;
        break;
    }

 done:
	if (response.buffer) free(response.buffer);

	return ret;
}

int send_request_to_init_with_exponential_backoff(char *remoteServer,
                                                  int remotePort,
                                                  char *uuid,
                                                  ServerInfo *server) {

    int backoffTime=1;
    int maxBackoff, backoffInterval;
    char *responseStr = NULL;

    srand(time(NULL));

    do {
        if (send_request_to_init(remoteServer,
                                 remotePort,
                                 uuid,
                                 server,
                                 &responseStr) == TRUE) {

            if (responseStr) free(responseStr);
            return TRUE;
        }

        // Calculate exponential backoff time
        maxBackoff = (1 << backoffTime) - 1;
        backoffInterval = rand() % maxBackoff + 1;

        printf("Error: %s. Retrying the boostrap in  %d seconds.\n",
               responseStr, backoffInterval);

        if (responseStr) free(responseStr);
        sleep(backoffInterval);

        backoffTime = (backoffTime < MAX_BACKOFF) ? backoffTime+1 : MAX_BACKOFF;
    } while (TRUE);

    return FALSE;
}

void free_server_info(ServerInfo *server) {

	if (server == NULL) return;

	if (server->IP)   free(server->IP);
	if (server->cert) free(server->cert);
	if (server->org)  free(server->org);
}

void log_debug_server(ServerInfo *server) {

	if (server == NULL) return;

	if (server->IP) {
		log_debug("server IP: %s", server->IP);
	}

	if (server->cert) {
		log_debug("server cert: %s", server->cert);
	}

	if (server->org) {
		log_debug("node org: %s", server->org);
	}
}
