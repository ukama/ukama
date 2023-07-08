/**
 * Copyright (c) 2023-present, Ukama Inc.
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

#include "config.h"
#include "web.h"
#include "usys_types.h"
#include "usys_log.h"
#include "web_service.h"

/*
 * init_framework -- initializa ulfius framework.
 *
 */
static int init_framework(UInst *inst, int port) {

	if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
		usys_log_error("Error initializing instance for webservice port %d",
                       port);
		return USYS_FALSE;
	}

	/* Set few params. */
	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	return USYS_TRUE;
}

/* 
 * start_framework --
 *
 */
static int start_framework(Config *config, UInst *instance) {

	int ret;
  
	/* open HTTPS/HTTP connection. */
    ret = ulfius_start_framework(instance);
	if (ret != U_OK) {
		usys_log_error("Error starting the webservice/websocket.");
    
		/* clean up. */
		ulfius_stop_framework(instance); /* don't think need this. XXX */
		ulfius_clean_instance(instance);
    
		return USYS_FALSE;
	}

    usys_log_debug("Webservice sucessfully started.");
	return USYS_TRUE;
}

/**
 * @fn      int web_service_start()
 * @brief   Add API endpoints and start the REST HTTP server
 *
 */
static void setup_webservice_endpoints(Config *config, UInst *instance) {
    
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);
    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("event/:service"), 0,
                               &web_service_cb_post_event, config);
    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("alert/:service"), 0,
                               &web_service_cb_post_event, config);
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("alert/:service"), 0,
                               &web_service_cb_ping, config);

    /* default */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, config);
}

/*
 * start_web_services -- start accepting REST clients
 *
 */
int start_web_services(Config *config, UInst *serviceInst) {

	if (init_framework(serviceInst, config->servicePort) != USYS_TRUE){
		usys_log_error("Error initializing webservice framework on port: %d",
                       config->servicePort);
		return USYS_FALSE;
	}

	/* setup endpoints and methods callback. */
	setup_webservice_endpoints(config, serviceInst);

	/* open connection for both admin and client webservices */
	if (!start_framework(config, serviceInst)) {
		usys_log_error("Failed to start webservices for client on port: %d",
                       config->servicePort);
		return USYS_FALSE;
	}

	usys_log_debug("Webservice started on port: %d", config->servicePort);

	return USYS_TRUE;
}


