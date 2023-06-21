/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Callback functions for various endpoints and REST methods.
 */

#include <ulfius.h>
#include <string.h>
#include <jansson.h>
#include <pthread.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <arpa/inet.h>

#include "callback.h"
#include "mesh.h"
#include "log.h"
#include "work.h"
#include "jserdes.h"
#include "map.h"
#include "httpStatus.h"

extern WorkList *Transmit;
extern MapTable *IDTable;
extern pthread_mutex_t mutex;
extern pthread_cond_t hasData;
extern char *queue;

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
							  void *data);
extern void websocket_incoming_message(const URequest *request,
									   WSManager *manager, WSMessage *message,
									   void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager,
							   void *data);

/*
 * Ulfius main callback function, calls the websocket manager and closes.
 */
int callback_websocket (const URequest *request, UResponse *response,
						void *data) {
	int ret;
	char *nodeID=NULL;
	Config *config = (Config *)data;

	nodeID = u_map_get(request->map_header, "User-Agent");
	if (nodeID == NULL) {
		log_error("Missing NodeID as User-Agent");
		return U_CALLBACK_ERROR;
	}

	if (config->deviceInfo) {
		if (strcmp(config->deviceInfo->nodeID, nodeID) != 0) {
			/* Only accept one device at a time until the socket is closed. */
			log_error("Only accept one device at a time. Ignoring");
			return U_CALLBACK_ERROR;
		}
	} else {
		config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
		if (config->deviceInfo == NULL) {
			log_error("Error allocating memory: %d", sizeof(DeviceInfo));
            free(nodeID);
			return U_CALLBACK_ERROR;
		}
        config->deviceInfo->nodeID = strdup(nodeID);
	}

	if ((ret = ulfius_set_websocket_response(response, NULL, NULL,
											 &websocket_manager,
											 data,
											 &websocket_incoming_message,
											 data,
											 &websocket_onclose,
											 data)) == U_OK) {
		ulfius_add_websocket_deflate_extension(response);
		return U_CALLBACK_CONTINUE;
	}

	return U_CALLBACK_CONTINUE;
}

/*
 * callback_not_allowed -- 
 *
 */
int callback_not_allowed(const URequest *request, UResponse *response,
						 void *user_data) {
  
	ulfius_set_string_body_response(response, 403, "Operation not allowed\n");
	return U_CALLBACK_CONTINUE;
}

/*
 * callback_default_websocket -- default callback for no-match
 *
 */
int callback_default_websocket(const URequest *request, UResponse *response,
							   void *user_data) {

	ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
	return U_CALLBACK_CONTINUE;
}

/*
 * callback_default -- default callback for no-match
 *
 */
int callback_default_webservice(const URequest *request, UResponse *response,
								void *data) {

	ulfius_set_string_body_response(response, 404, "You are clearly high!\n");
	return U_CALLBACK_CONTINUE;
}

/*
 * split_strings --
 *
 */
static void split_strings(char *input, char **str1, char **str2,
                          char *delimiter) {

    char *token=NULL;

    token = strtok(input, delimiter);

    if (token != NULL && str1) {
        *str1 = strdup(token);

        token = strtok(NULL, delimiter);
        if (token != NULL && str2) {
            *str2 = strdup(token);
        }
    }
}

/*
 * callback_webservice --
 *
 */
int callback_webservice(const URequest *request, UResponse *response,
						void *data) {

	int ret, statusCode=200;
	char *str, *destHost=NULL, *destPort=NULL, *service=NULL;
    char *requestStr=NULL, *url=NULL;

    url      = u_map_get(request->map_header, "Host");
    service  = u_map_get(request->map_header, "User-Agent");
    split_strings(url, &destHost, &destPort, ":");
    if (destHost == NULL || destPort == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    ret = serialize_websocket_message(&requestStr, request, destHost, destPort,
                                      service);
	if (ret == FALSE && requestStr == NULL) {
		log_error("Failed to convert request to JSON");
		statusCode = 400;
		goto done;
	} else {
		log_debug("Forward request JSON: %s", requestStr);
	}

	/* Add work for the websocket for transmission. */
    add_work_to_queue(&Transmit, requestStr, NULL, 0, NULL, 0);

	/* Wait for the response back. The cond is set by the websocket thread */
	pthread_mutex_lock(&mutex);
	log_debug("Waiting for response back from the server ...");
	pthread_cond_wait(&hasData, &mutex);
	pthread_mutex_unlock(&mutex);

    log_debug("Got response back from the system in the cloud.");

 done:
    ulfius_set_string_body_response(response, 200, queue);

	return U_CALLBACK_CONTINUE;
}
