/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

/*
 * Network related stuff based on ulfius framework.
 */
#include <ulfius.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#include "callback.h"
#include "mesh.h"
#include "websocket.h"
#include "config.h"
#include "jserdes.h"

/* define in websocket.c */
extern void websocket_manager(const URequest *request,
                              WSManager *manager,
							  void *data);
extern void websocket_incoming_message(const URequest *request,
									   WSManager *manager,
                                       WSMessage *message,
									   void *data);
extern void  websocket_onclose(const URequest *request,
                               WSManager *manager,
							   void *data);

static int find_first_available_port(int start, int end) {

    int port = -1, sockfd;
    struct sockaddr_in addr;

    for (port = start; port <= end; port++) {

        sockfd = socket(AF_INET, SOCK_STREAM, 0);
        if (sockfd < 0) return -1;

        addr.sin_family      = AF_INET;
        addr.sin_addr.s_addr = INADDR_ANY;
        addr.sin_port        = htons(port);

        if (bind(sockfd, (struct sockaddr *)&addr, sizeof(addr)) == 0) {
            close(sockfd);
            return port;
        }

        close(sockfd);
    }

    return 0;
}

static int init_framework(UInst *inst,
                          struct sockaddr_in *bindAddr,
                          int bindPort) {

    if (ulfius_init_instance(inst,
                             bindPort,
                             bindAddr,
                             NULL)!= U_OK) {
		log_error("Error initializing instance for websocket: %d",
				  bindPort);
		return FALSE;
	}

	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	return TRUE;
}

static void setup_webservice_endpoints(Config *config, UInst *instance) {

	ulfius_add_endpoint_by_val(instance, "GET",
                               EP_WEBSERVICE_PING, NULL, 0,
							   &callback_get_ping, config);

	ulfius_set_default_endpoint(instance,
                                &callback_default_webservice,
                                config);
}

static void setup_websocket_endpoints(Config *config, UInst *instance) {

	ulfius_add_endpoint_by_val(instance, "GET", EP_WEBSOCKET, NULL, 0,
							   &callback_websocket, config);
	ulfius_add_endpoint_by_val(instance, "POST", EP_WEBSOCKET, NULL, 0,
							   &callback_websocket, config);
	ulfius_add_endpoint_by_val(instance, "PUT", EP_WEBSOCKET, NULL, 0,
							   &callback_websocket, config);
	ulfius_add_endpoint_by_val(instance, "DELETE", EP_WEBSOCKET, NULL, 0,
							   &callback_websocket, config);

	ulfius_set_default_endpoint(instance,
                                &callback_default_websocket,
                                config);
}

static int start_framework(Config *config, UInst *instance, int flag) {

	int ret;
  
    ret = ulfius_start_framework(instance);

	if (ret != U_OK) {
		log_error("Error starting the webservice/websocket.");
    
		/* clean up. */
		ulfius_stop_framework(instance); /* don't think need this. XXX */
		ulfius_clean_instance(instance);
		return FALSE;
	}

	if (flag == WEBSOCKET) {
		log_debug("Websocket succesfully started.");
	} else if (flag == SERVICE) {
		log_debug("Webservice sucessfully started.");
	} else if (flag == FORWARD) {
        log_debug("Forward service sucessfully started.");
    }
  
	return TRUE;
}

int start_forward_service(Config *config, UInst **forwardInst) {

    struct sockaddr_in bindAddr;
    int port;

    port = find_first_available_port(START_PORT, END_PORT);
    if (port <= 0) {
        log_error("Unable to find empty port to bind");
        return FALSE;
    }

    memset(&bindAddr, 0, sizeof(bindAddr));
    bindAddr.sin_family      = AF_INET;
    bindAddr.sin_port        = htons(port);
    bindAddr.sin_addr.s_addr = inet_addr(config->bindingIP);

    *forwardInst = (UInst *)calloc(1, sizeof(UInst));
    if (*forwardInst == NULL) {
        log_error("Error allocating memory of size: %d",
                  sizeof(UInst));
        return FALSE;
    }

	if (init_framework(*forwardInst,
                       &bindAddr,
                       port) != TRUE) {
		log_error("Error initializing forward framework");
		return FALSE;
	}

	/* setup endpoint */
    ulfius_add_endpoint_by_val(*forwardInst,
                               "GET",
                               "*", NULL, 0,
							   &callback_default_forward, config);

    ulfius_set_default_endpoint(*forwardInst,
                                &callback_default_forward, config);

	if (start_framework(config,
                        *forwardInst,
                        FORWARD) == FALSE) {
		log_error("Failed to start forward service at port %d",
                  port);
		return FALSE;
	}

	log_debug("Forward service accepting on port: %d", port);

    return port;
}

int start_websocket_server(Config *config, UInst *websocketInst) {

    struct sockaddr_in bindAddr;

    memset(&bindAddr, 0, sizeof(bindAddr));
    bindAddr.sin_family = AF_INET;
    bindAddr.sin_port   = htons(atoi(config->websocketPort));
    bindAddr.sin_addr.s_addr = inet_addr(config->bindingIP);

	/* Initialize the admin and client webservices framework. */
	if (init_framework(websocketInst,
                       &bindAddr,
                       atoi(config->websocketPort)) != TRUE) {
		log_error("Error initializing websocket framework");
		return FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_websocket_endpoints(config, websocketInst);

	if (start_framework(config, websocketInst, WEBSOCKET)==FALSE) {
		log_error("Failed to start websocket at remote port %s",
				  config->websocketPort);
		return FALSE;
	}
	log_debug("Websocket accepting on port: %s", config->websocketPort);

 	return TRUE;
}

int start_web_services(Config *config, UInst *clientInst) {

	/* Initialize the admin and client webservices framework. */
	if (init_framework(clientInst, NULL, atoi(config->servicesPort)) != TRUE){
		log_error("Error initializing webservice framework");
		return FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_webservice_endpoints(config, clientInst);

	/* open connection for both admin and client webservices */
	if (!start_framework(config, clientInst, SERVICE)) {
		log_error("Failed to start webservices for client: %s",
                  config->servicesPort);
		return FALSE;
	}

	log_debug("Service accepting on port: %s", config->servicesPort);

	return TRUE;
}
