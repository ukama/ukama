/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "usys_log.h"

#include "callback.h"
#include "mesh.h"
#include "websocket.h"
#include "config.h"

#include "static.h"

#define WEB_SOCKETS 1
#define WEB_SERVICE 0
#define FWD_SERVICE 2

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager, void *data);
extern void websocket_incoming_message(const URequest *request,
                                       WSManager *manager,
                                       const WSMessage *message,
                                       void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager, void *data);

STATIC int init_framework(UInst *inst, int port) {

  if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
    usys_log_error("Error initializing instance for websocket remote port %d", port);
    return FALSE;
  }

  /* Set few params. */
  u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

  return TRUE;
}

STATIC void setup_webservice_endpoints(Config *config, UInst *instance) {

  /* Endpoint list declaration. */
  ulfius_add_endpoint_by_val(instance, "GET", PREFIX_FWDSERVICE, NULL, 0,
			     &callback_forward_service, config);
  ulfius_add_endpoint_by_val(instance, "POST", PREFIX_FWDSERVICE, NULL, 0,
			     &callback_forward_service, config);
  ulfius_add_endpoint_by_val(instance, "PUT", PREFIX_FWDSERVICE, NULL, 0,
			     &callback_forward_service, config);
  ulfius_add_endpoint_by_val(instance, "DELETE", PREFIX_FWDSERVICE, NULL, 0,
			     &callback_forward_service, config);

  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_not_allowed, config);
}

STATIC void setup_websocket_endpoints(Config *config, UInst *instance) {

  /* Endpoint list declaration. */
  ulfius_add_endpoint_by_val(instance, "GET", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_websocket, config);
  ulfius_add_endpoint_by_val(instance, "POST", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_websocket, config);
  ulfius_add_endpoint_by_val(instance, "PUT", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_not_allowed, config);
  ulfius_add_endpoint_by_val(instance, "DELETE", PREFIX_WEBSOCKET, NULL, 0,
			     &callback_not_allowed, config);
  
  /* default endpoint. */
  ulfius_set_default_endpoint(instance, &callback_not_allowed, config);
}

STATIC int start_framework(Config *config, UInst *instance, int flag) {
  
  if (ulfius_start_framework(instance) != U_OK) {
      usys_log_error("Error starting the webservice/websocket.");
    
      /* clean up. */
      ulfius_stop_framework(instance); /* don't think need this. XXX */
      ulfius_clean_instance(instance);
    
      return FALSE;
  }

  if (flag == WEB_SOCKETS) {
      usys_log_debug("Websocket succesfully started.");
  } else if (flag == FWD_SERVICE) {
      usys_log_debug("Forward service sucessfully started.");
  } else {
      usys_log_debug("Webservice sucessfully started.");
  }

  return TRUE;
}

int start_websocket_client(Config *config,
                           struct _websocket_client_handler *handler) {

    int ret = FALSE;
    struct _u_request request;
    struct _u_response response;

    if (ulfius_init_request(&request) != U_OK) goto done;
    if (ulfius_init_response(&response) != U_OK) goto done;

    if (ulfius_set_websocket_request(&request, config->remoteConnect,
                                     "protocol", "permessage-deflate") == U_OK) {

        if (config->deviceInfo && config->deviceInfo->nodeID) {
            u_map_put(request.map_header, "User-Agent", config->deviceInfo->nodeID);
            u_map_put(request.map_header, "X-node-id",  config->deviceInfo->nodeID);
        }

        ulfius_add_websocket_client_deflate_extension(handler);
        request.check_server_certificate = FALSE;

        /* Open websocket connection to remote host. */
        ret = ulfius_open_websocket_client_connection(&request,
                          &websocket_manager, (void *)config,
                          &websocket_incoming_message, (void *)config,
                          &websocket_onclose, (void*)config,
                          handler, &response);

        if (ret == U_OK) {
            ret = TRUE;
            goto done;
        } else {
            usys_log_error("Unable to open websocket connect to: %s",
                           config->remoteConnect);
            handler->websocket = NULL;
            ret = FALSE;
            goto done;
        }
    } else {
        usys_log_error("Error initializing the websocket client request");
        ret = FALSE;
        goto done;
    }

done:
    ulfius_clean_request(&request);
    ulfius_clean_response(&response);
    return ret;
}

int start_forward_services(Config *config, UInst *clientInst) {

    /* Initialize the admin and client webservices framework. */
    if (init_framework(clientInst, config->forwardPort) != TRUE){
        usys_log_error("Error initializing webservice framework");
        return FALSE;
    }

    /* setup endpoints and methods callback. */
    setup_webservice_endpoints(config, clientInst);

    /* open connection for both admin and client webservices */
    if (!start_framework(config, clientInst, FWD_SERVICE)) {
        usys_log_error("Failed to start webservices for client: %d",
                  config->forwardPort);
        return FALSE;
    }

    usys_log_debug("Forward service on port: %d started.", config->forwardPort);

    return TRUE;
}

int start_web_services(Config *config, UInst *clientInst) {

    if (init_framework(clientInst, config->servicePort) != TRUE){
        usys_log_error("Error initializing webservice framework");
        return FALSE;
    }

    ulfius_add_endpoint_by_val(clientInst, "GET", "/v1/", "ping", 0,
                               &web_service_cb_ping, config);
    ulfius_add_endpoint_by_val(clientInst, "GET", "/v1/", "version", 0,
                               &web_service_cb_version, config);
    ulfius_set_default_endpoint(clientInst, &web_service_cb_default, config);

    if (!start_framework(config, clientInst, WEB_SERVICE)) {
        usys_log_error("Failed to start webservices for client: %d",
                  config->servicePort);
        return FALSE;
    }

    usys_log_debug("Web service on port: %d started.", config->servicePort);

    return TRUE;
}
