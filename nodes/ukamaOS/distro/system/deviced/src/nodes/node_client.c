/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "node_profile.h"
#include "web_service.h"

#include "usys_log.h"

static int client_build_state(Config *config, JsonObj *json) {

    (void)config;
    (void)json;

    return STATUS_NOK;
}

static bool client_supports(ControlSubsystem subsystem) {

    return subsystem == CONTROL_SUBSYS_RADIO ||
           subsystem == CONTROL_SUBSYS_RESTART;
}

static int client_apply_radio(Config *config, ControlState desired) {

    FILE *fp;

    (void)config;

    usys_log_info("radio: %s (client emu)",
                  desired == CONTROL_STATE_ON ? "on" : "off");

    fp = fopen(DEF_RADIO_EMU_FILE, "w");
    if (!fp) {
        usys_log_error("radio: failed to open emu state file: %s", DEF_RADIO_EMU_FILE);
        return STATUS_NOK;
    }

    fprintf(fp, "%s\n", desired == CONTROL_STATE_ON ? "on" : "off");
    fclose(fp);

    return STATUS_OK;
}

static int client_apply(Config *config,
                        ControlSubsystem subsystem,
                        ControlState desired) {

    if (subsystem == CONTROL_SUBSYS_RADIO) {
        return client_apply_radio(config, desired);
    }

    usys_log_error("client: unsupported action");
    return STATUS_NOK;
}

static const NodeEndpoint client_endpoints[] = {
    { "POST", API_RES_EP("reboot"), web_service_cb_post_reboot },
    { "POST", API_RES_EP("radio"),  web_service_cb_post_radio_client },
    { NULL, NULL, NULL }
};

const NodeProfile node_profile_client = {
    .nodeType       = "client",
    .endpoints      = client_endpoints,
    .init_control   = NULL,
    .build_state    = client_build_state,
    .supports       = client_supports,
    .apply          = client_apply,
    .before_restart = NULL,
};
