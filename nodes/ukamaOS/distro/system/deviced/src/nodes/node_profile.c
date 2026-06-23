/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "node_profile.h"
#include "deviced.h"

#include "usys_log.h"

const NodeProfile *node_profile_get(Config *config) {

    if (!config) {
        return NULL;
    }

    if (config->clientMode) {
        return &node_profile_client;
    }

    if (!config->nodeType) {
        return NULL;
    }

    if (strcmp(config->nodeType, UKAMA_TOWER_NODE) == 0) {
        return &node_profile_tower;
    }

    if (strcmp(config->nodeType, UKAMA_AMPLIFIER_NODE) == 0) {
        return &node_profile_amplifier;
    }

    if (strcmp(config->nodeType, UKAMA_CONTROLLER_NODE) == 0) {
        return &node_profile_controller;
    }

    return NULL;
}

void node_profile_init_control(Config *config) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (profile && profile->init_control) {
        profile->init_control(config);
    }
}

int node_profile_build_state(Config *config, JsonObj *json) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (!profile || !profile->build_state) {
        return STATUS_NOK;
    }

    return profile->build_state(config, json);
}

int node_profile_apply(Config *config,
                       ControlSubsystem subsystem,
                       ControlState desired) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (!profile || !profile->apply) {
        usys_log_error("node profile: unsupported action");
        return STATUS_NOK;
    }

    return profile->apply(config, subsystem, desired);
}

int node_profile_before_restart(Config *config) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (!profile || !profile->before_restart) {
        return STATUS_OK;
    }

    return profile->before_restart(config);
}

bool node_profile_has_subsystem(Config *config, ControlSubsystem subsystem) {

    const NodeProfile *profile = NULL;

    profile = node_profile_get(config);
    if (!profile) {
        return false;
    }

    if (subsystem == CONTROL_SUBSYS_RESTART) {
        return true;
    }

    if (!profile->supports) {
        return false;
    }

    return profile->supports(subsystem);
}
