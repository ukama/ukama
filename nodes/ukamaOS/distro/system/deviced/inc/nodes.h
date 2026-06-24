/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef NODES_H_
#define NODES_H_

#include "config.h"
#include "control.h"
#include "deviced.h"

bool node_is_tower(Config *config);
bool node_is_amplifier(Config *config);
bool node_is_controller(Config *config);

void node_add_unsupported_methods(UInst *instance,
                                  char *allowedMethod,
                                  char *prefix,
                                  char *resource);

void node_tower_init_control(Config *config);
void node_amplifier_init_control(Config *config);
void node_controller_init_control(Config *config);

void node_tower_setup_endpoints(Config *config, UInst *instance);
void node_amplifier_setup_endpoints(Config *config, UInst *instance);
void node_controller_setup_endpoints(Config *config, UInst *instance);
void node_client_setup_endpoints(Config *config, UInst *instance);

int node_tower_build_state(Config *config, JsonObj *json);
int node_amplifier_build_state(Config *config, JsonObj *json);
int node_controller_build_state(Config *config, JsonObj *json);

int node_tower_apply_service(Config *config, ControlState desired);
int node_tower_apply_radio(Config *config, ControlState desired);
int node_tower_before_restart(Config *config);

int node_amplifier_apply_radio(Config *config, ControlState desired);

int node_client_apply_radio(Config *config, ControlState desired);

#endif /* NODES_H_ */
