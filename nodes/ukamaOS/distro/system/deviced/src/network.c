/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "config.h"
#include "deviced.h"
#include "node_profile.h"
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

static void setup_unsupported_methods(UInst *instance,
                                      const char *allowedMethod,
                                      const char *prefix,
                                      const char *resource) {

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
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, config);
    setup_unsupported_methods(instance, "GET",
                              URL_PREFIX, API_RES_EP("version"));
}

static int setup_profile_endpoints(Config *config,
                                   UInst *instance,
                                   const NodeProfile *profile) {

    const NodeEndpoint *ep = NULL;

    if (!config || !instance || !profile || !profile->endpoints) {
        return USYS_FALSE;
    }

    for (ep = profile->endpoints; ep->method != NULL; ep++) {
        ulfius_add_endpoint_by_val(instance,
                                   ep->method,
                                   URL_PREFIX,
                                   ep->resource,
                                   0,
                                   ep->callback,
                                   config);
        setup_unsupported_methods(instance,
                                  ep->method,
                                  URL_PREFIX,
                                  ep->resource);
    }

    return USYS_TRUE;
}

static int setup_webservice_endpoints(Config *config, UInst *serviceInst) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (!profile) {
        usys_log_error("Unable to setup web services for: %s",
                       config && config->nodeType ? config->nodeType : "unknown");
        return USYS_FALSE;
    }

    setup_common_webservice_endpoints(config, serviceInst);

    if (!setup_profile_endpoints(config, serviceInst, profile)) {
        return USYS_FALSE;
    }

    ulfius_set_default_endpoint(serviceInst, &web_service_cb_default, config);

    return USYS_TRUE;
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
    if (!setup_webservice_endpoints(config, serviceInst)) {
        return USYS_FALSE;
    }

    if (!start_framework(config, serviceInst)) {
        usys_log_error("Failed to start webservices on port: %d",
                       config->servicePort);
        return USYS_FALSE;
    }

    usys_log_debug("Webservice started on port: %d", config->servicePort);

    return USYS_TRUE;
}
