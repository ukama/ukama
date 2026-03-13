/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSERDES_H
#define JSERDES_H

#include <jansson.h>

#include "metrics_store.h"
#include "driver.h"
#include "config.h"

/*
 * JSON serialization for controller data
 */

/* Serialize controller status to JSON */
json_t *json_serialize_status(const MetricsSnapshot *snap);

/* Serialize metrics for metricsd (Prometheus-style) */
json_t *json_serialize_metrics(const MetricsSnapshot *snap, const char *node_id);

/* Serialize alarm history to JSON */
json_t *json_serialize_alarms(const AlarmRecord *alarms, int count);

/* Serialize single alarm notification (for notify.d) */
json_t *json_serialize_alarm_notification(const Config *config,
                                          AlarmType type,
                                          Severity severity,
                                          const char *message);

/* Deserialize control request (absorption/float voltage, relay) */
int json_deserialize_voltage_request(json_t *json, double *voltage_v);
int json_deserialize_relay_request(json_t *json, bool *state);

#endif /* JSERDES_H */
