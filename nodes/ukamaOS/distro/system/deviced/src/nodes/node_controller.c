/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "node_profile.h"
#include "web_service.h"

#include "usys_log.h"

static void controller_init_control(Config *config) {

    (void)config;
}

static int controller_build_state(Config *config, JsonObj *json) {

    (void)config;
    (void)json;

    /* CNode does not own service/radio. */
    return STATUS_OK;
}

static bool controller_supports(ControlSubsystem subsystem) {

    return subsystem == CONTROL_SUBSYS_RESTART;
}

static int controller_apply(Config *config,
                            ControlSubsystem subsystem,
                            ControlState desired) {

    (void)config;
    (void)subsystem;
    (void)desired;

    usys_log_error("controller: unsupported action");
    return STATUS_NOK;
}

static const NodeEndpoint controller_endpoints[] = {
    { "GET",  API_RES_EP("state"),   web_service_cb_state },
    { "POST", API_RES_EP("restart"), web_service_cb_post_restart },
    { NULL, NULL, NULL }
};

const NodeProfile node_profile_controller = {
    .nodeType       = UKAMA_CONTROLLER_NODE,
    .endpoints      = controller_endpoints,
    .init_control   = controller_init_control,
    .build_state    = controller_build_state,
    .supports       = controller_supports,
    .apply          = controller_apply,
    .before_restart = NULL,
};
