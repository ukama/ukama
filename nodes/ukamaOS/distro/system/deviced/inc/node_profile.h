/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef NODE_PROFILE_H_
#define NODE_PROFILE_H_

#include "config.h"
#include "control.h"
#include "deviced.h"

typedef int (*NodeEndpointCb)(const URequest *request,
                              UResponse *response,
                              void *epConfig);

typedef struct {
    const char     *method;
    const char     *resource;
    NodeEndpointCb callback;
} NodeEndpoint;

typedef struct {
    const char         *nodeType;
    const NodeEndpoint *endpoints;

    void (*init_control)(Config *config);
    int  (*build_state)(Config *config, JsonObj *json);
    bool (*supports)(ControlSubsystem subsystem);
    int  (*apply)(Config *config,
                  ControlSubsystem subsystem,
                  ControlState desired);
    int  (*before_restart)(Config *config);
} NodeProfile;

extern const NodeProfile node_profile_client;
extern const NodeProfile node_profile_tower;
extern const NodeProfile node_profile_amplifier;
extern const NodeProfile node_profile_controller;

const NodeProfile *node_profile_get(Config *config);

void node_profile_init_control(Config *config);
int node_profile_build_state(Config *config, JsonObj *json);
int node_profile_apply(Config *config,
                       ControlSubsystem subsystem,
                       ControlState desired);
int node_profile_before_restart(Config *config);
bool node_profile_has_subsystem(Config *config, ControlSubsystem subsystem);

#endif /* NODE_PROFILE_H_ */
