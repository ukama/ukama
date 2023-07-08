/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>

#include "deviced.h"
#include "config.h"
#include "web_service.h"

static int start_framework(Config *config, UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
		usys_log_error("Error starting the webservice/websocket.");
    
		/* clean up. */
		ulfius_stop_framework(instance); /* don't think need this. XXX */
		ulfius_clean_instance(instance);
    
		return USYS_FALSE;
	}

    usys_log_debug("Webservice sucessfully started.");
	return USYS_TRUE;
}

static void setup_webservice_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("reboot/:id"), 0,
                               &web_service_cb_post_reboot, config);

    if (config->clientMode == USYS_FALSE) {
        ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                                   API_RES_EP("restart/:id"), 0,
                                   &web_service_cb_post_restart, config);
    }

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
