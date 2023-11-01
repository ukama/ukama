/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

/*
 * Network related stuff based on ulfius framework
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "initClient.h"
#include "config.h"
#include "callback.h"
#include "log.h"

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

	/* Add ping EP for liveliness */
	ulfius_add_endpoint_by_val(instance, "GET", EP_PING, NULL, 0,
							   &callback_ping, config);

	/* EP for client to get system info */
	ulfius_add_endpoint_by_val(instance, "GET", EP_SYSTEMS, NULL, 0,
							   &callback_get_systems, config);

	/* default endpoint. */
	ulfius_set_default_endpoint(instance, &callback_default_webservice, config);
}

/* 
 * start_framework --
 *
 */
static int start_framework(UInst *instance) {

	if (ulfius_start_framework(instance) != U_OK) {
		log_error("Error starting the webservice/websocket.");

		/* clean up. */
		ulfius_stop_framework(instance); /* don't think need this. XXX */
		ulfius_clean_instance(instance);

		return FALSE;
	}

	log_debug("Webservice sucessfully started.");

	return TRUE;
}

/*
 * start_web_services -- start accepting REST clients on 127.0.0.1:port
 *
 */
int start_web_services(Config *config, UInst *clientInst) {

	/* Initialize the webservices framework. */
	if (init_framework(clientInst, atoi(config->port)) != TRUE){
		log_error("Error initializing webservice framework");
		return FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_webservice_endpoints(config, clientInst);

	/* open connection for both admin and client webservices */
	if (!start_framework(clientInst)) {
		log_error("Failed to start webservices for client: %s",
				  config->port);
		return FALSE;
	}

	log_debug("Webservice on client port: %s started.", config->port);

	return TRUE;
}
