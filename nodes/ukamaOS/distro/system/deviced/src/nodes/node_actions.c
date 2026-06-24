/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "actions.h"
#include "nodes.h"

#include "usys_log.h"

int actions_service_apply(Config *config, ControlState desired) {

    if (!config) return STATUS_NOK;

    if (node_is_tower(config)) {
        return node_tower_apply_service(config, desired);
    }

    usys_log_error("service: unsupported on this device.d instance");
    return STATUS_NOK;
}

int actions_radio_apply(Config *config, ControlState desired) {

    if (!config) return STATUS_NOK;

    if (config->clientMode) {
        return node_client_apply_radio(config, desired);
    }

    if (node_is_tower(config)) {
        return node_tower_apply_radio(config, desired);
    }

    if (node_is_amplifier(config)) {
        return node_amplifier_apply_radio(config, desired);
    }

    usys_log_error("radio: unsupported node type: %s",
                   config->nodeType ? config->nodeType : "unknown");
    return STATUS_NOK;
}
