/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "nodes.h"

#include "http_status.h"
#include "web_client.h"
#include "web_service.h"

#include "usys_log.h"

void node_tower_init_control(Config *config) {

    if (!config || !config->control) return;

    config->control->Service.Current = CONTROL_STATE_OFF;
    config->control->Service.Desired = CONTROL_STATE_OFF;

    config->control->Radio.Current   = CONTROL_STATE_ON;
    config->control->Radio.Desired   = CONTROL_STATE_ON;
}

void node_tower_setup_endpoints(Config *config, UInst *instance) {

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

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("service"), 0,
                               &web_service_cb_post_service, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("service"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("radio"), 0,
                               &web_service_cb_post_radio, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("radio"));
}

int node_tower_build_state(Config *config, JsonObj *json) {

    char state[32] = {0};

    if (!config || !config->control || !json) return STATUS_NOK;

    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_SERVICE,
                                        state,
                                        sizeof(state)) != STATUS_OK) {
        return STATUS_NOK;
    }
    json_object_set_new(json, "service", json_string(state));

    memset(state, 0, sizeof(state));
    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_RADIO,
                                        state,
                                        sizeof(state)) != STATUS_OK) {
        return STATUS_NOK;
    }
    json_object_set_new(json, "radio", json_string(state));

    return STATUS_OK;
}

int node_tower_apply_service(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) return STATUS_NOK;

    usys_log_info("service: %s", desired == CONTROL_STATE_ON ? "on" : "off");

    return wc_post_service_to_pcrf(config, desired, &retCode);
}

int node_tower_apply_radio(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) return STATUS_NOK;

    usys_log_info("radio: %s (tower trx client)",
                  desired == CONTROL_STATE_ON ? "on" : "off");

    return wc_send_radio_to_client(config, desired, &retCode) == USYS_OK ?
           STATUS_OK : STATUS_NOK;
}

int node_tower_before_restart(Config *config) {

    int retCode = -1;

    if (!config) return STATUS_NOK;

    if (wc_send_reboot_to_client(config, &retCode) != USYS_OK) {
        usys_log_error("Remote client reboot failed");
        return STATUS_NOK;
    }

    if (retCode != HttpStatus_Accepted) {
        usys_log_error("Remote client reboot not accepted: %d (%s)",
                       retCode,
                       HttpStatusStr(retCode));
        return STATUS_NOK;
    }

    return STATUS_OK;
}
