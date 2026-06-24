/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "nodes.h"

#include "web_client.h"
#include "web_service.h"

#include "usys_log.h"

void node_amplifier_init_control(Config *config) {

    if (!config || !config->control) return;

    config->control->Radio.Current = CONTROL_STATE_ON;
    config->control->Radio.Desired = CONTROL_STATE_ON;
}

void node_amplifier_setup_endpoints(Config *config, UInst *instance) {

    ulfius_add_endpoint_by_val(instance, "GET", URL_PREFIX,
                               API_RES_EP("state"), 0,
                               &web_service_cb_state, config);
    node_add_unsupported_methods(instance, "GET", URL_PREFIX,
                                 API_RES_EP("state"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("reboot"), 0,
                               &web_service_cb_post_reboot, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("reboot"));

    ulfius_add_endpoint_by_val(instance, "POST", URL_PREFIX,
                               API_RES_EP("radio"), 0,
                               &web_service_cb_post_radio, config);
    node_add_unsupported_methods(instance, "POST", URL_PREFIX,
                                 API_RES_EP("radio"));
}

int node_amplifier_build_state(Config *config, JsonObj *json) {

    char state[32] = {0};

    if (!config || !config->control || !json) return STATUS_NOK;

    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_RADIO,
                                        state,
                                        sizeof(state)) != STATUS_OK) {
        return STATUS_NOK;
    }
    json_object_set_new(json, "radio", json_string(state));

    return STATUS_OK;
}

int node_amplifier_apply_radio(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) return STATUS_NOK;

    usys_log_info("radio: %s (via femd:%d)",
                  desired == CONTROL_STATE_ON ? "on" : "off",
                  config->femPort);

    retCode = -1;
    if (wc_put_gpio_to_femd(config, 1, desired, &retCode) != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=1 http=%d", retCode);
        return STATUS_NOK;
    }

    retCode = -1;
    if (wc_put_gpio_to_femd(config, 2, desired, &retCode) != STATUS_OK) {
        usys_log_error("radio: femd gpio apply failed fem=2 http=%d", retCode);
        return STATUS_NOK;
    }

    return STATUS_OK;
}
