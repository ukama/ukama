/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef POLICY_H_
#define POLICY_H_

#include "switchd.h"

#define SWITCHD_POLICY_SOURCE_SITE_CONTROLLER "site-controller"
#define SWITCHD_POLICY_PATH_DEFAULT "/ukama/configs/switch/policy.json"

const char *policy_state_to_str(SwitchPolicyState state);
const char *policy_type_to_str(SwitchPortPolicyType type);
const char *policy_action_to_str(SwitchPolicyAction action);

int policy_load(SwitchdContext *ctx);
int policy_apply_body(SwitchdContext *ctx,
                      const char *body,
                      size_t bodyLen,
                      char *err,
                      size_t errLen);
int policy_check_action(SwitchdContext *ctx,
                        uint32_t portId,
                        SwitchPolicyAction action,
                        const char *source,
                        char *err,
                        size_t errLen);
const SwitchPortPolicy *policy_get_port(const SwitchdContext *ctx,
                                        uint32_t portId);
JsonObj *policy_serialize(const SwitchdContext *ctx);
JsonObj *policy_serialize_overlay(const SwitchdContext *ctx,
                                  uint32_t portId);

#endif /* POLICY_H_ */
