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
#include "http_status.h"

#include "usys_log.h"

static void tower_init_control(Config *config) {

    if (!config || !config->control) {
        return;
    }

    config->control->Service.Current = CONTROL_STATE_OFF;
    config->control->Service.Desired = CONTROL_STATE_OFF;
    config->control->Radio.Current   = CONTROL_STATE_ON;
    config->control->Radio.Desired   = CONTROL_STATE_ON;
}

static int tower_build_state(Config *config, JsonObj *json) {

    char state[32];

    if (!config || !config->control || !json) {
        return STATUS_NOK;
    }

    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_SERVICE,
                                        state,
                                        sizeof(state)) != STATUS_OK) {
        return STATUS_NOK;
    }
    json_object_set_new(json, "service", json_string(state));

    if (control_get_subsys_public_state(config->control,
                                        CONTROL_SUBSYS_RADIO,
                                        state,
                                        sizeof(state)) == STATUS_OK) {
        json_object_set_new(json, "radio", json_string(state));
    }

    return STATUS_OK;
}

static bool tower_supports(ControlSubsystem subsystem) {

    return subsystem == CONTROL_SUBSYS_SERVICE ||
           subsystem == CONTROL_SUBSYS_RADIO   ||
           subsystem == CONTROL_SUBSYS_RESTART;
}

static int tower_apply_service(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) {
        return STATUS_NOK;
    }

    usys_log_info("service: %s", desired == CONTROL_STATE_ON ? "on" : "off");

    return wc_post_service_to_pcrf(config, desired, &retCode);
}

static int tower_apply_radio(Config *config, ControlState desired) {

    int retCode = -1;

    if (!config) {
        return STATUS_NOK;
    }

    usys_log_info("radio: %s (tower trx client)",
                  desired == CONTROL_STATE_ON ? "on" : "off");

    return wc_send_radio_to_client(config, desired, &retCode) == USYS_OK ?
           STATUS_OK : STATUS_NOK;
}

static int tower_apply(Config *config,
                       ControlSubsystem subsystem,
                       ControlState desired) {

    if (subsystem == CONTROL_SUBSYS_SERVICE) {
        return tower_apply_service(config, desired);
    }

    if (subsystem == CONTROL_SUBSYS_RADIO) {
        return tower_apply_radio(config, desired);
    }

    usys_log_error("tower: unsupported action");
    return STATUS_NOK;
}

static int tower_before_restart(Config *config) {

    int retCode = -1;

    if (!config) {
        return STATUS_NOK;
    }

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

static const NodeEndpoint tower_endpoints[] = {
    { "GET",  API_RES_EP("state"),   web_service_cb_state },
    { "POST", API_RES_EP("restart"), web_service_cb_post_restart },
    { "POST", API_RES_EP("service"), web_service_cb_post_service },
    { "POST", API_RES_EP("radio"),   web_service_cb_post_radio },
    { NULL, NULL, NULL }
};

const NodeProfile node_profile_tower = {
    .nodeType       = UKAMA_TOWER_NODE,
    .endpoints      = tower_endpoints,
    .init_control   = tower_init_control,
    .build_state    = tower_build_state,
    .supports       = tower_supports,
    .apply          = tower_apply,
    .before_restart = tower_before_restart,
};
