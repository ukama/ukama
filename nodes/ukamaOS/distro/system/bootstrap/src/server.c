/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to interactions with bootstrap server */

#include <curl/curl.h>
#include <curl/easy.h>
#include <string.h>
#include <stdlib.h>
#include <jansson.h>

#include "jserdes.h"
#include "server.h"
#include "log.h"

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
 * send_request_to_server
 *
 */
static long send_request_to_server(char *url, char *uuid, Response *response) {

	long resCode=0;
	CURL *curl=NULL;
	CURLcode res;
	struct curl_slist *headers=NULL;

	curl_global_init(CURL_GLOBAL_ALL);
	curl = curl_easy_init();
	if (curl == NULL) {
		return resCode;
	}

	headers = curl_slist_append(headers, KEY_UUID);
	headers = curl_slist_append(headers,": ");
	headers = curl_slist_append(headers, uuid);

	headers = curl_slist_append(headers, KEY_LOOKING_FOR);
	headers = curl_slist_append(headers,": ");
	headers = curl_slist_append(headers, VALUE_VALIDATION);

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
	} else {
		/* get status code */
		curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
	}

	curl_slist_free_all(headers);
	curl_easy_cleanup(curl);
	curl_global_cleanup();

	return resCode;
}

/*
 * process_response_from_server --
 *
 */
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

/*
 * register_to_server --
 *
 */
int register_to_server(char *bootstrapServer, char *uuid, ServerInfo *server) {

	int ret=FALSE;
	Response response = {NULL, 0};

	if (bootstrapServer == NULL || uuid == NULL) return FALSE;

	if (send_request_to_server(bootstrapServer, uuid, &response) == 200) {
		if (process_response_from_server(response.buffer, server)) {
			log_debug_server(server);
			log_debug("Recevied server IP: %s", server->IP);
		} else {
			log_error("Unable to receive proper server info from bootstrap: %s",
					  bootstrapServer);
			goto done;
		}
	} else {
		log_error("Unable to send request to bootstrap server at: %s",
				  bootstrapServer);
		goto done;
	}

	ret = TRUE;

 done:
	if (response.buffer) free(response.buffer);

	return ret;
}

/*
 * clear_server --
 *
 */
void free_server_info(ServerInfo *server) {

	if (server == NULL) return;

	if (server->IP)   free(server->IP);
	if (server->cert) free(server->cert);
	if (server->org)  free(server->org);
}

/*
 * log_debug_server --
 *
 */
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
