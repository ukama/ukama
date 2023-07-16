/**
 * Copyright (c) 2022-present, Ukama Inc.
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
#include "u_amqp.h"
#include "httpStatus.h"

extern MapTable *IDsTable;

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
	int ret;
	char *nodeID=NULL;
	Config *config=NULL;
    MapItem *map=NULL;
    char ip[INET_ADDRSTRLEN]={0};
    struct sockaddr_in *sin=NULL;

    config = (Config *)data;

    sin = (struct sockaddr_in *)request->client_address;
    inet_ntop(AF_INET, &sin->sin_addr, &ip[0], INET_ADDRSTRLEN);

	nodeID = u_map_get(request->map_header, "User-Agent");
	if (nodeID == NULL) {
		log_error("Missing NodeID as User-Agent");
		return U_CALLBACK_ERROR;
	}

    map = add_map_to_table(&IDsTable,
                           nodeID,
                           &ip[0], sin->sin_port,
                           &ip[0], sin->sin_port);
	if (map == NULL) {
        return U_CALLBACK_CONTINUE; // XXX
	}

    map->configData = data;

	/* Publish device (nodeID) 'connect' event to AMQP exchange */
	if (publish_event(CONN_CONNECT,
                      nodeID,
                      &ip[0], sin->sin_port,
                      &ip[0], sin->sin_port) == FALSE) {
		log_error("Error publishing device connect msg on AMQP exchange");
        //		return U_CALLBACK_ERROR; xxx
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
 * callback_ping --
 *
 */
int callback_ping(const URequest *request, UResponse *response,
						void *data) {

	ulfius_set_string_body_response(response, 200, "ok");
	return U_CALLBACK_CONTINUE;
}

/*
 * callback_webservice --
 *
 */
int callback_webservice(const URequest *request, UResponse *response,
						void *data) {

	int ret, statusCode=200;
	char *str;
	MapItem *map=NULL;
	struct sockaddr_in *sin;
    char *nodeID=NULL, *nodePort=NULL, *url=NULL, *agent=NULL;
    char *packetStr=NULL;
    char *requestStr=NULL, *responseStr=NULL;

    /*
     * Find the NodeID from the URL.
     * For the given NodeID, look up the Map item in the IDstable.
     *   If match found, add the task to the transmit list and set cond
     *   otherwise return 503 (Service Unavailable)
     */

    url   = u_map_get(request->map_header, "Host");
    agent = u_map_get(request->map_header, "User-Agent");
    split_strings(url, &nodeID, &nodePort, ":");

    if (nodeID == NULL || nodePort == NULL) {
        ulfius_set_string_body_response(response, HttpStatus_BadRequest,
                                        HttpStatusStr(HttpStatus_BadRequest));
        return U_CALLBACK_CONTINUE;
    }

    map = is_existing_item(IDsTable, nodeID);
    if (map == NULL) { /* No matching node connected. */
        log_error("For agent: %s No matching node with nodeID: %s",
                  agent, nodeID);
        ulfius_set_string_body_response(response, HttpStatus_NotFound,
                                        HttpStatusStr(HttpStatus_NotFound));
        return U_CALLBACK_CONTINUE;
    }

	ret = serialize_websocket_message(&requestStr, request, nodeID, nodePort,
                                      agent);
	if (ret == FALSE && requestStr == NULL) {
		log_error("Failed to convert request to JSON");
		statusCode = 400;
		goto done;
	} else {
		log_debug("Forward request JSON: %s", requestStr);
	}

	/* Add work for the websocket for transmission. */
    add_work_to_queue(&map->transmit, requestStr, NULL, 0, NULL, 0);

	/* Wait for the response back. The cond is set by the websocket thread */
	pthread_mutex_lock(&(map->mutex));
	log_debug("Waiting for response back from the server ...");
	pthread_cond_wait(&(map->receive->hasWork), &(map->mutex));
	pthread_mutex_unlock(&(map->mutex));

    // xxx
    responseStr = map->receive->first->data;
	log_debug("Got response back from server. Response: %s", responseStr);

	/* Send response back. */
	if (map->size == 0) {
		statusCode = 402;
		goto done;
	}
  
 done:
    // xxx
	/* Send response back to the callee */
	// if (statusCode != 200) {
	//	ulfius_set_string_body_response(response, statusCode, "");
	//} else {
	//	ulfius_set_string_body_response(response, statusCode, map->data);
	//}
    ulfius_set_string_body_response(response, 200, responseStr);
    //    free(packetStr);
    // free(requestStr);
    
	// if (map->size)
	//	free(map->data);

	return U_CALLBACK_CONTINUE;
}
