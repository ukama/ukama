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

#include "web_service.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_file.h"

static int init_framework(UInst *inst, int port) {

	if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
		usys_log_error("Error initializing instance for webservice port %d",
                       port);
		return USYS_FALSE;
	}

	u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

	return USYS_TRUE;
}

static int start_framework(UInst *instance) {

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

static void setup_admin_webservice_endpoints(UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("version"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, NULL);
}

int start_admin_web_service(UInst *adminInst) {

    int adminPort;

    adminPort = usys_find_service_port(SERVICE_METRICS_ADMIN);
    if (adminPort == 0) {
        usys_log_error("Unable to determine admin port for: %s",
                       SERVICE_METRICS_ADMIN);
        return USYS_FALSE;
    }
    
	if (init_framework(adminInst, adminPort) != USYS_TRUE){
		usys_log_error("Error initializing admin framework on port: %d",
                       adminPort);
		return USYS_FALSE;
	}

	setup_admin_webservice_endpoints(adminInst);

	if (!start_framework(adminInst)) {
		usys_log_error("Failed to start admin webservice on port: %d",
                       adminPort);
		return USYS_FALSE;
	}

	usys_log_debug("Webservice started on port: %d", adminPort);

	return USYS_TRUE;
}
