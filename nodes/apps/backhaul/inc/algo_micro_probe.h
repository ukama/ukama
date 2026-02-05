/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ALGO_MICRO_PROBE_H_
#define ALGO_MICRO_PROBE_H_

#include "config.h"
#include "metrics_store.h"
#include "worker.h"

int algo_micro_probe_run(Config *config, MetricsStore *store, Worker *worker);

#endif /* ALGO_MICRO_PROBE_H_ */
