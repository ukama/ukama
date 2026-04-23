/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include <stdlib.h>
#include <string.h>

#include "network.h"
#include "web_service.h"

#include "usys_file.h"
#include "usys_log.h"
#include "usys_services.h"

static int init_framework(UInst *instance, int port) {

    if (ulfius_init_instance(instance, port, NULL, NULL) != U_OK) {
        usys_log_error("failed to initialize web service on port %d", port);
        return 0;
    }

    u_map_put(instance->default_headers, "Access-Control-Allow-Origin", "*");
    return 1;
}

static int start_framework(UInst *instance) {

    if (ulfius_start_framework(instance) != U_OK) {
        usys_log_error("failed to start web service framework");
        ulfius_stop_framework(instance);
        ulfius_clean_instance(instance);
        return 0;
    }

    return 1;
}

static void setup_unsupported_methods(UInst *instance,
                                      char *allowedMethod,
                                      char *prefix,
                                      char *resource) {

    if (strcmp(allowedMethod, "GET") != 0) {
        ulfius_add_endpoint_by_val(instance, "GET", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "POST") != 0) {
        ulfius_add_endpoint_by_val(instance, "POST", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "PUT") != 0) {
        ulfius_add_endpoint_by_val(instance, "PUT", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
    if (strcmp(allowedMethod, "DELETE") != 0) {
        ulfius_add_endpoint_by_val(instance, "DELETE", prefix, resource, 0,
                                   &web_service_cb_not_allowed, NULL);
    }
}

int start_metrics_web_service(UInst *metricsInst, AppState *state) {

    int port = 0;

    port = usys_find_service_port(SERVICE_NAME);
    if (port <= 0) {
        usys_log_error("unable to determine metrics port from services");
        return RETURN_NOTOK;
    }

    if (!init_framework(metricsInst, port)) {
        return RETURN_NOTOK;
    }

    ulfius_add_endpoint_by_val(metricsInst, "GET", NULL, "/metrics", 0,
                               &web_service_cb_metrics, state);
    setup_unsupported_methods(metricsInst, "GET", NULL, "/metrics");
    ulfius_set_default_endpoint(metricsInst, &web_service_cb_default, NULL);

    if (!start_framework(metricsInst)) {
        return RETURN_NOTOK;
    }

    usys_log_info("metrics web service started on port %d", port);
    return RETURN_OK;
}

int start_admin_web_service(UInst *adminInst, AppState *state) {

    int port = 0;

    port = usys_find_service_port(SERVICE_NAME_ADMIN);
    if (port <= 0) {
        usys_log_error("unable to determine admin port from services");
        return RETURN_NOTOK;
    }

    if (!init_framework(adminInst, port)) {
        return RETURN_NOTOK;
    }

    ulfius_add_endpoint_by_val(adminInst, "GET", URL_PREFIX,
                               API_RES_EP("ping"), 0,
                               &web_service_cb_ping, state);
    setup_unsupported_methods(adminInst, "GET", URL_PREFIX,
                              API_RES_EP("ping"));

    ulfius_add_endpoint_by_val(adminInst, "GET", URL_PREFIX,
                               API_RES_EP("version"), 0,
                               &web_service_cb_version, state);
    setup_unsupported_methods(adminInst, "GET", URL_PREFIX,
                              API_RES_EP("version"));

    ulfius_add_endpoint_by_val(adminInst, "GET", URL_PREFIX,
                               API_RES_EP("status"), 0,
                               &web_service_cb_status, state);
    setup_unsupported_methods(adminInst, "GET", URL_PREFIX,
                              API_RES_EP("status"));

    ulfius_set_default_endpoint(adminInst, &web_service_cb_default, NULL);

    if (!start_framework(adminInst)) {
        return RETURN_NOTOK;
    }

    usys_log_info("admin web service started on port %d", port);
    return RETURN_OK;
}

void stop_web_service(UInst *inst) {

    if (inst == NULL) {
        return;
    }

    ulfius_stop_framework(inst);
    ulfius_clean_instance(inst);
}
