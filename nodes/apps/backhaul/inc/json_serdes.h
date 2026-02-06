/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_SERDES_H
#define JSON_SERDES_H

#include "jansson.h"
#include "backhauld.h"

#define JTAG_NEAR_URL "nearUrl"
#define JTAG_FAR_URL  "farUrl"

/* Parse: {"nearUrl":"...", "farUrl":"..."} (keys may also match json_types.h tags if you standardize there) */
int json_parse_reflector_set(const char *json_str, ReflectorSet *out);

/* Build a JSON object for diagnostics/status */
json_t* json_build_reflector_set(const ReflectorSet *set);

/* /v1/status */
json_t* json_build_status(const BackhaulMetrics *m);

/* /v1/metrics (snapshot + status) */
json_t* json_build_metrics_snapshot(MetricsStore *store);

/* Optional: build a generic {"queued":true,"name":"..."} response */
json_t* json_build_queued(const char *name);

/* Safely dumps JSON to string. Caller must free() the returned char*. */
char* json_dump_compact(json_t *obj);

/* Reads a JSON string field; returns NULL if missing/not string. */
const char* json_get_string(json_t *obj, const char *key);

/* Reads a JSON int field; returns def if missing/not integer. */
int json_get_int(json_t *obj, const char *key, int def);

json_t* json_backhaul_status(MetricsStore *store);

json_t* json_backhaul_metrics(MetricsStore *store);

#endif /* JSON_SERDES_H */
