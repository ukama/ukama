/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <ulfius.h>

#include "example.h"
#include "web_service.h"

#include "usys_log.h"
#include "usys_types.h"

static int init_framework(UInst *inst, int port) {

    if (ulfius_init_instance(inst, port, NULL, NULL) != U_OK) {
        usys_log_error("Error initializing webservice on port %d", port);
        return USYS_FALSE;
    }

    u_map_put(inst->default_headers, "Access-Control-Allow-Origin", "*");

    return USYS_TRUE;
}

static int start_framework(UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("Error starting webservice");

        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);

        return USYS_FALSE;
    }

    usys_log_debug("Webservice successfully started.");

    return USYS_TRUE;
}

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   allowedMethod);
    }

    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   allowedMethod);
    }

    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   allowedMethod);
    }

    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed,
                                   allowedMethod);
    }
}

static void setup_webservice_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping,
                               config);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version,
                               config);
    setup_unsupported_methods(instance, "GET", URL_PREFIX,
                              API_RES_EP("version"));

    ulfius_set_default_endpoint(instance, &web_service_cb_default, config);
}

int start_web_services(Config *config, UInst *serviceInst) {

    if (!config || !serviceInst) {
        return USYS_FALSE;
    }

    if (init_framework(serviceInst, config->servicePort) != USYS_TRUE) {
        usys_log_error("Error initializing webservice on port %d",
                       config->servicePort);
        return USYS_FALSE;
    }

    setup_webservice_endpoints(config, serviceInst);

    if (start_framework(serviceInst) != USYS_TRUE) {
        usys_log_error("Failed to start webservice on port %d",
                       config->servicePort);
        return USYS_FALSE;
    }

    usys_log_info("Webservice started on port %d", config->servicePort);

    return USYS_TRUE;
}
