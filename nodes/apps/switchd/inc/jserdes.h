/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSERDES_H_
#define JSERDES_H_

#include "switchd.h"

bool json_deserialize_bool_request(const URequest *request,
                                   const char *key,
                                   bool *value);

bool json_deserialize_int_request(const URequest *request,
                                  const char *key,
                                  int *value);

bool json_deserialize_firmware_stage_request(const URequest *request,
                                             char *path,
                                             size_t pathLen,
                                             char *version,
                                             size_t versionLen,
                                             char *sha256,
                                             size_t sha256Len);

bool json_serialize_alarm_notification(JsonObj **json,
                                       const char *serviceName,
                                       const SwitchAlarm *alarm,
                                       bool clear);

JsonObj *json_serialize_metrics(const SwitchdContext *ctx);
JsonObj *json_serialize_status(const SwitchdContext *ctx);
JsonObj *json_serialize_switch_info(const SwitchdContext *ctx);
JsonObj *json_serialize_switch_health(const SwitchdContext *ctx);
JsonObj *json_serialize_switch_capabilities(const SwitchdContext *ctx);
JsonObj *json_serialize_switch_kpis(const SwitchdContext *ctx);
JsonObj *json_serialize_port(const SwitchPortState *port);
JsonObj *json_serialize_port_with_policy(const SwitchdContext *ctx,
                                         const SwitchPortState *port);
JsonObj *json_serialize_ports(const SwitchdContext *ctx);
JsonObj *json_serialize_firmware(const SwitchdContext *ctx);
JsonObj *json_serialize_active_alarms(const SwitchdContext *ctx);

void json_free(JsonObj **json);

#endif /* JSERDES_H_ */
