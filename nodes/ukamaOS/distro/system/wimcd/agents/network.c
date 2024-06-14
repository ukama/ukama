/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>

#include "log.h"
#include "callback.h"
#include "agent.h"
#include "wimc.h"
#include "agent/network.h"
#include "agent/jserdes.h"
#include "agent/callback.h"

#include "usys_types.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_api.h"

static int get_agent_port(char *method) {

    char buffer[128] = {0};

    sprintf(buffer, "wimc-agent-%s", method);

    return usys_find_service_port(buffer);
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

static void setup_endpoints(struct _u_instance *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &agent_web_service_cb_ping, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &agent_web_service_cb_version, NULL);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("version"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("apps/:name/:tag"), 0,
                               &agent_web_service_cb_post_capp, NULL);
    setup_unsupported_methods(instance, "POST",
                              URL_PREFIX, API_RES_EP("apps/:name/:tag"));

    ulfius_set_default_endpoint(instance,
                                &agent_web_service_cb_default, NULL);
}

bool start_web_service(char *method, struct _u_instance *webInstance) {

    int servicePort = 0;

    servicePort = get_agent_port(method);
    if (ulfius_init_instance(webInstance, servicePort, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for port %d", servicePort);
        return USYS_FALSE;
    }

    u_map_put(webInstance->default_headers, "Access-Control-Allow-Origin", "*");
    webInstance->max_post_body_size = 1024;

    setup_endpoints(webInstance);

    if (ulfius_start_framework(webInstance) != U_OK) {
        usys_log_error("Failed to start webservices at port:%s", servicePort);

        ulfius_stop_framework(webInstance); 
        ulfius_clean_instance(webInstance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", servicePort);

    return USYS_TRUE;
}
