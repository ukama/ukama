/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "agent.h"
#include "agent/callback.h"
#include "agent/jserdes.h"
#include "agent/network.h"
#include "callback.h"
#include "log.h"
#include "wimc.h"

#include "usys_api.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_services.h"
#include "usys_types.h"

static int get_agent_port(char *method) {

    char buffer[128];

    if (method == NULL || *method == '\0') {
        return 0;
    }

    if (snprintf(buffer, sizeof(buffer), "wimc-agent-%s", method) >=
        (int)sizeof(buffer)) {
        return 0;
    }

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

    int servicePort;

    servicePort = get_agent_port(method);
    if (servicePort <= 0) {
        usys_log_error("Unable to find service port for wimc-agent-%s",
                       method ? method : "(null)");
        return USYS_FALSE;
    }

    if (ulfius_init_instance(webInstance, servicePort, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing instance for port %d",
                       servicePort);
        return USYS_FALSE;
    }

    u_map_put(webInstance->default_headers, "Access-Control-Allow-Origin", "*");
    webInstance->max_post_body_size = 4096;

    setup_endpoints(webInstance);

    if (ulfius_start_framework(webInstance) != U_OK) {
        usys_log_error("Failed to start webservices at port:%d",
                       servicePort);

        ulfius_stop_framework(webInstance);
        ulfius_clean_instance(webInstance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", servicePort);
    return USYS_TRUE;
}
