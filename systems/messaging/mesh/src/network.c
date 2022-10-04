/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Network related stuff based on ulfius framework.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "callback.h"
#include "mesh.h"
#include "websocket.h"
#include "config.h"
#include "jserdes.h"

#define WEB_SOCKETS 1
#define WEB_SERVICE 0

/* define in websocket.c */
extern void websocket_manager(const URequest *request, WSManager *manager,
							  void *data);
extern void websocket_incoming_message(const URequest *request,
									   WSManager *manager, WSMessage *message,
									   void *data);
extern void  websocket_onclose(const URequest *request, WSManager *manager,
							   void *data);
/*
 * init_framework -- initializa ulfius framework.
 *
 */
static int init_framework(UInst *inst, int port) {

	if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
		log_error("Error initializing instance for websocket remote port %d",
				  port);
		return FALSE;
	}

	/* Set few params. */
	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	return TRUE;
}

/*
 * setup_webservice_endpoints --
 *
 */
static void setup_webservice_endpoints(Config *config, UInst *instance) {

	/* Endpoint list declaration. */
	ulfius_add_endpoint_by_val(instance, "GET", PREFIX_WEBSERVICE, NULL, 0,
							   &callback_webservice, config);
	ulfius_add_endpoint_by_val(instance, "POST", PREFIX_WEBSERVICE, NULL, 0,
							   &callback_webservice, config);
	ulfius_add_endpoint_by_val(instance, "PUT", PREFIX_WEBSERVICE, NULL, 0,
							   &callback_webservice, config);
	ulfius_add_endpoint_by_val(instance, "DELETE", PREFIX_WEBSERVICE, NULL, 0,
							   &callback_webservice, config);

	/* Add ping EP for liveliness */
	ulfius_add_endpoint_by_val(instance, "GET", EP_PING, NULL, 0,
							   &callback_ping, config);

	/* default endpoint. */
	ulfius_set_default_endpoint(instance, &callback_default_webservice, config);
}

/*
 * setup_websocket_endpoints -- 
 *
 */
static void setup_websocket_endpoints(Config *config, UInst *instance) {

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
	ulfius_set_default_endpoint(instance, &callback_default_websocket, config);
}

/* 
 * start_framework --
 *
 */
static int start_framework(Config *config, UInst *instance, int flag) {

	int ret;
  
	/* open HTTPS/HTTP connection. */
	if (config->secure && flag == WEB_SOCKETS) {
		ret = ulfius_start_secure_framework(instance, config->keyFile,
											config->certFile);
	} else {
		ret = ulfius_start_framework(instance);
	}

	if (ret != U_OK) {
		log_error("Error starting the webservice/websocket.");
    
		/* clean up. */
		ulfius_stop_framework(instance); /* don't think need this. XXX */
		ulfius_clean_instance(instance);
    
		return FALSE;
	}

	if (flag == WEB_SOCKETS) {
		log_debug("Websocket succesfully started.");
	} else {
		log_debug("Webservice sucessfully started.");
	}
  
	return TRUE;
}

/*
 * start_websocket_server -- start websocket server on the server port.
 *
 */

int start_websocket_server(Config *cfg, UInst *serverInst) {

	/* Initialize the admin and client webservices framework. */
	if (init_framework(serverInst, atoi(cfg->remoteAccept)) != TRUE) {
		log_error("Error initializing webservice framework");
		return FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_websocket_endpoints(cfg, serverInst);
  
	/* open connection for both admin and client webservices */
	if (start_framework(cfg, serverInst, WEB_SOCKETS)==FALSE) {
		log_error("Failed to start websocket at remote port %s",
				  cfg->remoteAccept);
		return FALSE;
	}
  
	log_debug("Websocket on remote port %s: started.", cfg->remoteAccept);

 	return TRUE;
}

/*
 * add_device_info_to_request -- Add device related information to the
 *                               request
 *
 */
static int add_device_info_to_request(struct _u_request *request,
									  Config *config) {
	json_t *json=NULL;
	char *jStr=NULL;

	if (serialize_device_info(&json, config->deviceInfo) == FALSE) {
		log_error("Failed to serialize device info for request");
		return FALSE;
	}

	/* Add the json into request body. */
	jStr = json_dumps(json, 0);
	if (jStr == NULL) {
		json_decref(json);
		return FALSE;
	}

	request->binary_body_length = strlen(jStr);
	request->binary_body = strdup(jStr);

	free(jStr);
	json_decref(json);

	return TRUE;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */
int start_web_services(Config *config, UInst *clientInst) {

	/* Initialize the admin and client webservices framework. */
	if (init_framework(clientInst, atoi(config->localAccept)) != TRUE){
		log_error("Error initializing webservice framework");
		return FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_webservice_endpoints(config, clientInst);

	/* open connection for both admin and client webservices */
	if (!start_framework(config, clientInst, WEB_SERVICE)) {
		log_error("Failed to start webservices for client: %s",
				  config->localAccept);
		return FALSE;
	}

	log_debug("Webservice on client port: %s started.", config->localAccept);

	return TRUE;
}
