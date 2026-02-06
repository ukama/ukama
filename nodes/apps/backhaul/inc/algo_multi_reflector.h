/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ALGO_MULTI_REFLECTOR_H
#define ALGO_MULTI_REFLECTOR_H

#include "config.h"
#include "metrics_store.h"

int algo_multi_reflector_run(Config *config,
                             MetricsStore *store,
                             void *unused);

#endif /* ALGO_MULTI_REFLECTOR_H_ */
