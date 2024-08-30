/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "config.h"
#include "web.h"
#include "usys_types.h"
#include "usys_log.h"
#include "web_service.h"

static int init_framework(UInst *inst, int port) {

	if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
		usys_log_error("Error initializing instance for webservice port %d",
                       port);
		return USYS_FALSE;
	}

	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	return USYS_TRUE;
}

static int start_framework(Config *config, UInst *instance) {

	int ret;
  
    ret = ulfius_start_framework(instance);
	if (ret != U_OK) {
		usys_log_error("Error starting the webservice/websocket.");
    
		ulfius_stop_framework(instance);
		ulfius_clean_instance(instance);

		return USYS_FALSE;
	}

    usys_log_debug("Webservice sucessfully started.");
	return USYS_TRUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix,
                                   resource, 0,
                                   &web_service_cb_not_allowed,
                                   (void *)allowedMethod);
    }
}

static void setup_webservice_endpoints(Config *config, UInst *instance) {
    
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));


    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("config"), 0,
                               &web_service_cb_post_config, config);
    setup_unsupported_methods(instance, "POST",
                              URL_PREFIX, API_RES_EP("config"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, config);
}

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


