/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JSON_SERDES_H
#define JSON_SERDES_H

#include <jansson.h>

#include "json_types.h"

json_t *json_serdes_power_metrics_to_json(const PowerMetrics *m);

#endif
