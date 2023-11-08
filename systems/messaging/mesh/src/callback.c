/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
#include <uuid/uuid.h>

#include "callback.h"
#include "mesh.h"
#include "log.h"
#include "work.h"
#include "jserdes.h"
#include "map.h"
#include "u_amqp.h"
#include "httpStatus.h"

extern MapTable *NodesTable;

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
							  void *data);
extern void websocket_incoming_message(const URequest *request,
									   WSManager *manager, WSMessage *message,
									   void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager,
							   void *data);

/*
 * Ulfius main callback function, send AMQP msg and calls the websocket
 * manager and closes.
 */
int callback_websocket(const URequest *request, UResponse *response,
                       void *data) {
	int ret, forwardPort;
	char *nodeID=NULL;
	Config *config=NULL;
    MapItem *map=NULL;
    char ip[INET_ADDRSTRLEN]={0};
    struct sockaddr_in *sin = NULL;
    UInst *forwardInst      = NULL;

    config = (Config *)data;

    sin = (struct sockaddr_in *)request->client_address;
    inet_ntop(AF_INET, &sin->sin_addr, &ip[0], INET_ADDRSTRLEN);

	nodeID = u_map_get(request->map_header, "User-Agent");
	if (nodeID == NULL) {
		log_error("Missing NodeID as User-Agent");
		return U_CALLBACK_ERROR;
	}

    /* Open up forwarding web instance for services */
    forwardPort = start_forward_service(config, &forwardInst);
    if (forwardPort <= 0 ) {
        log_error("Unable to start forwarding serice");
        return U_CALLBACK_ERROR;
    }

    map = add_map_to_table(&NodesTable,
                           nodeID,
                           &forwardInst,
                           &ip[0], sin->sin_port,
                           config->bindingIP,
                           forwardPort);
	if (map == NULL) {
        ulfius_stop_framework(forwardInst);
        ulfius_clean_instance(forwardInst);
        return U_CALLBACK_ERROR;
	}

    map->configData = data;

	/* Publish device (nodeID) 'connect' event to AMQP exchange */
	if (publish_event(CONN_CONNECT,
                      nodeID,
                      &ip[0], sin->sin_port,
                      config->bindingIP,
                      forwardPort) == FALSE) {
		log_error("Error publishing device connect msg on AMQP exchange");
        remove_map_item_from_table(&NodesTable, nodeID);
        ulfius_stop_framework(forwardInst);
        ulfius_clean_instance(forwardInst);
        return U_CALLBACK_ERROR;
	} else {
		log_debug("AMQP device connect msg successfull for NodeID: %s", nodeID);
	}

	if ((ret = ulfius_set_websocket_response(response, NULL, NULL,
											 &websocket_manager,
											 map->nodeInfo->nodeID,
											 &websocket_incoming_message,
											 map->nodeInfo->nodeID,
											 &websocket_onclose,
											 map->nodeInfo->nodeID)) == U_OK) {
		ulfius_add_websocket_deflate_extension(response);
		return U_CALLBACK_CONTINUE;
	}

	return U_CALLBACK_CONTINUE;
}

int callback_default_websocket(const URequest *request,
                               UResponse *response,
							   void *user_data) {

	ulfius_set_string_body_response(response,
                                    HttpStatus_Forbidden,
                                    HttpStatusStr(HttpStatus_Forbidden));
	return U_CALLBACK_CONTINUE;
}

int callback_default_webservice(const URequest *request,
                                UResponse *response,
								void *data) {

	ulfius_set_string_body_response(response,
                                    HttpStatus_Forbidden,
                                    HttpStatusStr(HttpStatus_Forbidden));
	return U_CALLBACK_CONTINUE;
}

int callback_get_ping(const URequest *request,
                      UResponse *response,
                      void *data) {

	ulfius_set_string_body_response(response, 200, "");
	return U_CALLBACK_CONTINUE;
}

int callback_default_forward(const URequest *request,
                             UResponse *response,
                             void *user_data) {

    MapItem *map=NULL;
    char *host=NULL, *port=NULL, *url=NULL;
    char *requestStr=NULL;
    char *responseStr=NULL;
    int statusCode;
    Forward *forward = NULL;
    char uuidStr[36+1];
    uuid_t uuid;

    url   = u_map_get(request->map_header, "Host");
    split_strings(url, &host, &port, ":");

    if (host == NULL || port == NULL) {
        ulfius_set_string_body_response(response,
                                        HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    map = is_existing_item_by_port(NodesTable, atoi(port));
    if (map == NULL) { /* No matching node connected. */
        log_error("No matching node on port: %s", port);
        ulfius_set_string_body_response(response,
                                        HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        free(port);
        free(host);
        return U_CALLBACK_CONTINUE;
    }

    uuid_generate(uuid);
    uuid_unparse(uuid, uuidStr);
    forward = add_client_to_list(&map->forwardList, uuidStr);
    if (forward == NULL) {
        log_error("Error adding to the forward list");
        statusCode  = HttpStatus_InternalServerError;
        responseStr = HttpStatusStr(statusCode);
        goto done;
    }

    if (serialize_websocket_message(&requestStr,
                                    request,
                                    uuidStr) == FALSE) {
        log_error("Failed to convert request to JSON");
        statusCode  = HttpStatus_InternalServerError;
        responseStr = HttpStatusStr(statusCode);
        goto done;
    } else {
        log_debug("Forward request JSON: %s", requestStr);
    }

    /* Add work for the websocket for transmission. */
    add_work_to_queue(&map->transmit, requestStr, NULL, 0, NULL, 0);
    free(requestStr);

    /* Wait for the response back. The cond is set by the websocket thread */
    pthread_mutex_lock(&(forward->mutex));
    log_debug("Waiting for response back from the node...");

    pthread_cond_wait(&forward->hasData, &forward->mutex);
	pthread_mutex_unlock(&forward->mutex);

    log_debug("Response from System Code: %d len: %d Data: %s",
              forward->httpCode,
              forward->size,
              (char *)forward->data);

    statusCode  = forward->httpCode;
    responseStr = (char *)forward->data;

done:
    ulfius_set_string_body_response(response,
                                    statusCode,
                                    responseStr);

    remove_item_from_list(map->forwardList, uuidStr);
    free(forward->data);
    free(host);
    free(port);

    return U_CALLBACK_CONTINUE;
}
