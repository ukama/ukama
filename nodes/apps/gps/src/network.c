/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "gpsd.h"
#include "config.h"
#include "web_service.h"

static int start_framework(Config *config, UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
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

    /* ping */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));

    /* version */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("version"));

    /* lock */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("lock"), 0,
                               &web_service_cb_lock, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("lock"));

    /* coordinates */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("coordinates"), 0,
                               &web_service_cb_coordinates, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("coordinates"));

    /* time */
    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("time"), 0,
                               &web_service_cb_time, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("time"));

    /* default */
    ulfius_set_default_endpoint(instance, &web_service_cb_default, config);
}

int start_web_service(Config *config, UInst *serviceInst) {

    if (ulfius_init_instance(serviceInst,
                             config->servicePort,
                             NULL,
                             NULL) != U_OK) {
		usys_log_error("Error initializing instance for webservice port %d",
                       config->servicePort);
		return USYS_FALSE;
	}

	/* Set few params. */
	u_map_put(serviceInst->default_headers, "Access-Control-Allow-Origin", "*");

	/* setup endpoints and methods callback. */
	setup_webservice_endpoints(config, serviceInst);

	if (!start_framework(config, serviceInst)) {
		usys_log_error("Failed to start webservices on port: %d",
                       config->servicePort);
		return USYS_FALSE;
	}

	usys_log_debug("Webservice started on port: %d", config->servicePort);

	return USYS_TRUE;
}
