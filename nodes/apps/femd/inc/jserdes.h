/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <stdbool.h>
#include <jansson.h>

#include "femd.h"
#include "snapshot.h"
#include "jobs.h"
#include "config.h"

void json_log(json_t *json);

bool json_serialize_op_id(json_t **json, uint64_t opId);
bool json_serialize_op_status(json_t **json, const OpStatus *st);

bool json_serialize_ctrl_snapshot(json_t **json, const CtrlSnapshot *s);
bool json_serialize_fem_snapshot(json_t **json, FemUnit unit, const FemSnapshot *s);

bool json_deserialize_node_info(char **data,
                                int  *iData,
                                char *tag,
                                json_type type,
                                json_t *json);

bool json_serialize_pa_alarm_notification(json_t **json,
                                          Config *config,
                                          int type);

void json_free(json_t **json);

#endif /* JSERDES_H */
