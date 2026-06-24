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

#include "deviced.h"
#include "config.h"
#include "nodes.h"
#include "web_service.h"

static int start_framework(Config *config, UInst *instance) {

    (void)config;

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting the webservice/websocket.");

        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice sucessfully started.");
    return USYS_TRUE;
}

void node_add_unsupported_methods(UInst *instance,
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

static void setup_common_webservice_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, config);
    node_add_unsupported_methods(instance, "GET",
                                 URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, config);
    node_add_unsupported_methods(instance, "GET",
                                 URL_PREFIX, API_RES_EP("version"));
}

static int setup_webservice_endpoints(Config *config, UInst *serviceInst) {

    setup_common_webservice_endpoints(config, serviceInst);

    if (config->clientMode) {
        node_client_setup_endpoints(config, serviceInst);
        ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);
        return USYS_TRUE;
    }

    if (node_is_tower(config)) {
        node_tower_setup_endpoints(config, serviceInst);
        ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);
        return USYS_TRUE;
    }

    if (node_is_amplifier(config)) {
        node_amplifier_setup_endpoints(config, serviceInst);
        ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);
        return USYS_TRUE;
    }

    if (node_is_controller(config)) {
        node_controller_setup_endpoints(config, serviceInst);
        ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);
        return USYS_TRUE;
    }

    usys_log_error("Unable to setup web services for: %s",
                   config->nodeType ? config->nodeType : "unknown");
    return USYS_FALSE;
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

    /* Setup endpoints and methods callback. */
    if (setup_webservice_endpoints(config, serviceInst) != USYS_TRUE) {
        usys_log_error("Unable to setup endpoints: %s",
                       config->nodeType ? config->nodeType : "unknown");
        return USYS_FALSE;
    }

    return start_framework(config, serviceInst);
}
