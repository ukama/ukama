/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "nodes.h"
#include "web_service.h"

#include "usys_log.h"

void node_client_setup_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("reboot"), 0,
                               &web_service_cb_post_client_reboot,
                               config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("reboot"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("radio"), 0,
                               &web_service_cb_post_radio_client, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("radio"));
}

int node_client_apply_radio(Config *config, ControlState desired) {

    FILE *fp;

    if (!config) return STATUS_NOK;

    usys_log_info("radio: %s (client emu)",
                  desired == CONTROL_STATE_ON ? "on" : "off");

    fp = fopen(DEF_RADIO_EMU_FILE, "w");
    if (!fp) {
        usys_log_error("radio: failed to open emu state file: %s",
                       DEF_RADIO_EMU_FILE);
        return STATUS_NOK;
    }

    fprintf(fp, "%s\n", desired == CONTROL_STATE_ON ? "on" : "off");
    fclose(fp);

    return STATUS_OK;
}
