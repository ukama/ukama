/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "node_profile.h"
#include "web_client.h"
#include "web_service.h"

#include "usys_log.h"

static void amplifier_init_control(Config *config) {

    if (!config || !config->control) {
        return;
    }

    config->control->Radio.Current = CONTROL_STATE_ON;
    config->control->Radio.Desired = CONTROL_STATE_ON;
}

static int amplifier_build_state(Config *config, JsonObj *json) {

    char state[32];

    if (!config || !config->control || !json) {
        return STATUS_NOK;
    }

    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_RADIO,
                                        state,
                                        sizeof(state)) != STATUS_OK) {
        return STATUS_NOK;
    }

    json_object_set_new(json, "radio", json_string(state));
    return STATUS_OK;
}

static bool amplifier_supports(ControlSubsystem subsystem) {

    return subsystem == CONTROL_SUBSYS_RADIO ||
           subsystem == CONTROL_SUBSYS_RESTART;
}

static int amplifier_apply_radio(Config *config, ControlState desired) {

    int retCode;

    if (!config) {
        return STATUS_NOK;
    }

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

static int amplifier_apply(Config *config,
                           ControlSubsystem subsystem,
                           ControlState desired) {

    if (subsystem == CONTROL_SUBSYS_RADIO) {
        return amplifier_apply_radio(config, desired);
    }

    usys_log_error("amplifier: unsupported action");
    return STATUS_NOK;
}

static const NodeEndpoint amplifier_endpoints[] = {
    { "GET",  API_RES_EP("state"),   web_service_cb_state },
    { "POST", API_RES_EP("restart"), web_service_cb_post_restart },
    { "POST", API_RES_EP("radio"),   web_service_cb_post_radio },
    { NULL, NULL, NULL }
};

const NodeProfile node_profile_amplifier = {
    .nodeType       = UKAMA_AMPLIFIER_NODE,
    .endpoints      = amplifier_endpoints,
    .init_control   = amplifier_init_control,
    .build_state    = amplifier_build_state,
    .supports       = amplifier_supports,
    .apply          = amplifier_apply,
    .before_restart = NULL,
};
