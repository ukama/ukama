/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "nodes.h"

#include "web_service.h"

void node_controller_init_control(Config *config) {

    (void)config;
}

void node_controller_setup_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("state"), 0,
                               &web_service_cb_state, config);
    node_add_unsupported_methods(instance, "GET", URL_PREFIX,
                                 API_RES_EP("state"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("restart"), 0,
                               &web_service_cb_post_restart, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("restart"));

    /* CNode owns no service/radio endpoint in device.d. */
}

int node_controller_build_state(Config *config, JsonObj *json) {

    (void)config;
    (void)json;

    /* CNode exposes only common state fields added by web_service.c. */
    return STATUS_OK;
}
