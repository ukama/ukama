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
#include "amqp.h"

extern WorkList *Transmit;
extern MapTable *IDTable;

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
			      void *data);
extern void websocket_incoming_message(const URequest *request,
				       WSManager *manager, WSMessage *message,
				       void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager,
			       void *data);

/*
 *
 */

static int is_valid_request(URequest *request) {

  return TRUE;
}

/*
 * Ulfius main callback function, send AMQP msg and calls the websocket
 * manager and closes.
 */
int callback_websocket (const URequest *request, UResponse *response,
			void *data) {
  int ret;
  char *idStr=NULL;
  Config *config = (Config *)data;
  uuid_t uuid;

  idStr = u_map_get(request->map_header, "User-Agent");
  if (idStr == NULL) {
    log_error("Missing UUID as User-Agent");
    return U_CALLBACK_ERROR;
  }

  if (uuid_parse(idStr, uuid)==-1) {
    log_error("Error parsing the UUID into binary: %s", idStr);
    return U_CALLBACK_ERROR;
  }

  if (config->deviceInfo) {
    if (uuid_compare(config->deviceInfo->uuid, uuid) != 0) {
      /* Only accept one device at a time until the socket is closed. */
      log_error("Only accept one device at a time. Ignoring");
      free(idStr);
      return U_CALLBACK_ERROR;
    }
  } else {
    config->deviceInfo = (DeviceInfo *)malloc(sizeof(DeviceInfo));
    if (config->deviceInfo == NULL) {
      log_error("Error allocating memory: %d", sizeof(DeviceInfo));
      free(idStr);
      return U_CALLBACK_ERROR;
    }
    uuid_copy(config->deviceInfo->uuid, uuid);
  }

  /* Publish device (uuid) 'connect' event to AMQP exchange */
  if (publish_amqp_event(config->conn, config->amqpExchange, CONN_CONNECT,
			 config->mode, uuid) == FALSE) {
    log_error("Error publishing device connect msg on AMQP exchange");
    free(idStr);
    return U_CALLBACK_ERROR;
  } else {
    log_debug("AMQP device connect msg successfull for UUID: %s", idStr);
  }

  if ((ret = ulfius_set_websocket_response(response, NULL, NULL,
					   &websocket_manager,
					   data,
					   &websocket_incoming_message,
					   data,
					   &websocket_onclose,
					   data)) == U_OK) {
    ulfius_add_websocket_deflate_extension(response);
    free(idStr);
    return U_CALLBACK_CONTINUE;
  }

  free(idStr);
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
 * callback_webservice --
 *
 */
int callback_webservice(const URequest *request, UResponse *response,
			void *data) {

  json_t *jReq=NULL;
  Config *config;
  int ret, statusCode=200;
  char *str;
  char ip[INET_ADDRSTRLEN];
  unsigned short port;
  MapItem *map=NULL;
  struct sockaddr_in *sin;

  config = (Config *)data;
  
  /* For every incoming request, do following:
   *
   * 1. Sanity check.
   * 2. Convert request into JSON.
   * 3. Send request to Ukama proxy via websocket.
   * 4. Process websocket response.
   * 5. Wait for the response from server.
   * 6. Process response.
   * 7. Send response back to the client.
   * 8. Done
   */

  if (is_valid_request(request)==FALSE) {
    statusCode=400;
    goto done;
  }

  sin = (struct sockaddr_in *)request->client_address;
  inet_ntop(AF_INET, &sin->sin_addr, &ip[0], INET_ADDRSTRLEN);
  port = sin->sin_port;

  map = add_map_to_table(&IDTable, &ip[0], port);
  if (map == NULL) {
    statusCode = 500;
    goto done;
  }

  ret = serialize_forward_request(request, &jReq, config, map->uuid);
  if (ret == FALSE && jReq == NULL) {
    log_error("Failed to convert request to JSON");
    statusCode = 400;
    goto done;
  } else {
    str = json_dumps(jReq, 0);
    log_debug("Forward request JSON: %s", str);
    free(str);
  }

  /* Add work for the websocket for transmission. */
  if (jReq != NULL) {
    /* No pre/post transmission func. This will block. */
    add_work_to_queue(&Transmit, (Packet)jReq, NULL, 0, NULL, 0);
  }

  /* Wait for the response back. The cond is set by the websocket thread */
  pthread_mutex_lock(&(map->mutex));
  log_debug("Waiting for response back from the server ...");
  pthread_cond_wait(&(map->hasResp), &(map->mutex));
  pthread_mutex_unlock(&(map->mutex));

  log_debug("Got response back from server. Len: %d Response: %s",
	    map->size, (char *)map->data);

  /* Send response back. */
  if (map->size == 0) {
    statusCode = 402;
    goto done;
  }
  
 done:
  /* Send response back to the callee */
  if (statusCode != 200) {
    ulfius_set_string_body_response(response, statusCode,
				    "Something went wrong! What you up to?");
  } else {
    ulfius_set_string_body_response(response, statusCode, map->data);
  }

  if (map->size)
    free(map->data);

  return U_CALLBACK_CONTINUE;
}
