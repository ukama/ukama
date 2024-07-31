/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <ulfius.h>
#include <string.h>
#include <jansson.h>
#include <pthread.h>

#include <sys/socket.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <arpa/inet.h>

#include "usys_log.h"

#include "callback.h"
#include "mesh.h"
#include "work.h"
#include "jserdes.h"
#include "map.h"
#include "httpStatus.h"
#include "config.h"

#include "version.h"
#include "static.h"

extern WorkList *Transmit;
extern MapTable *ClientTable;
extern pthread_mutex_t mutex;
extern pthread_cond_t hasData;
extern char *queue;

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
							  void *data);
extern void websocket_incoming_message(const URequest *request,
									   WSManager *manager, WSMessage *message,
									   void *data);
extern void websocket_onclose(const URequest *request, WSManager *manager,
                              void *data);

int callback_websocket (const URequest *request, UResponse *response,
						void *data) {
	int ret;
	char *nodeID=NULL;
	Config *config = (Config *)data;

	nodeID = u_map_get(request->map_header, "User-Agent");
	if (nodeID == NULL) {
		usys_log_error("Missing NodeID as User-Agent");
		return U_CALLBACK_ERROR;
	}

	if (config->deviceInfo) {
		if (strcmp(config->deviceInfo->nodeID, nodeID) != 0) {
			/* Only accept one device at a time until the socket is closed. */
			usys_log_error("Only accept one device at a time. Ignoring");
			return U_CALLBACK_ERROR;
		}
	} else {
		config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
		if (config->deviceInfo == NULL) {
			usys_log_error("Error allocating memory: %d", sizeof(DeviceInfo));
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

int callback_not_allowed(const URequest *request,
                         UResponse *response,
						 void *user_data) {
  
	ulfius_set_string_body_response(response,
                                    HttpStatus_Forbidden,
                                    HttpStatusStr(HttpStatus_Forbidden));
	return U_CALLBACK_CONTINUE;
}

int callback_forward_service(const URequest *request,
                             UResponse *response,
                             void *data) {

	int ret;
	char *destHost=NULL, *destPort=NULL, *service=NULL;
    char *requestStr=NULL, *url=NULL;
    char ip[INET_ADDRSTRLEN]={0}, sourcePort[MAX_BUFFER]={0};
    struct sockaddr_in *sin=NULL;
    MapItem *map=NULL;

    sin = (struct sockaddr_in *)request->client_address;
    inet_ntop(AF_INET, &sin->sin_addr, &ip[0], INET_ADDRSTRLEN);
    sprintf(sourcePort, "%d",sin->sin_port);

    url      = u_map_get(request->map_header, "Host");
    service  = u_map_get(request->map_header, "User-Agent");
    split_strings(url, &destHost, &destPort, ":");
    if (destHost == NULL || destPort == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    ret = serialize_websocket_message(&requestStr, request, destHost, destPort,
                                      service, sourcePort);
	if (ret == FALSE && requestStr == NULL) {
		usys_log_error("Failed to convert request to JSON");
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
	} else {
		usys_log_debug("Forward request JSON: %s", requestStr);
	}

    /* map it */
    map = add_map_to_table(&ClientTable, service, sourcePort);

	/* Add work for the websocket for transmission. */
    add_work_to_queue(&Transmit, requestStr, NULL, 0, NULL, 0);

	/* Wait for the response back. The cond is set by the websocket thread */
	pthread_mutex_lock(&map->mutex);
	usys_log_debug("Waiting for response back from the server ...");
	pthread_cond_wait(&map->hasResp, &map->mutex);
	pthread_mutex_unlock(&map->mutex);

    usys_log_debug("Response from System Code: %d len: %d Data: %s",
                   map->code, map->size, map->data);

    ulfius_set_string_body_response(response, map->code, (char *)map->data);

	return U_CALLBACK_CONTINUE;
}

int web_service_cb_ping(const URequest *request,
                        UResponse *response,
                        void *epConfig) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    HttpStatusStr(HttpStatus_OK));

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_version(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_OK,
                                    VERSION);

    return U_CALLBACK_CONTINUE;
}

int web_service_cb_default(const URequest *request,
                           UResponse *response,
                           void *epConfig) {

    ulfius_set_string_body_response(response,
                                    HttpStatus_NotFound,
                                    HttpStatusStr(HttpStatus_NotFound));

    return U_CALLBACK_CONTINUE;
}
