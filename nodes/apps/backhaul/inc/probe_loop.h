/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef PROBE_LOOP_H
#define PROBE_LOOP_H

#include <pthread.h>

#include "config.h"
#include "metrics_store.h"

int probe_loop_start(pthread_t *thread,
                     Config *config,
                     MetricsStore *store,
                     volatile int *stopFlag);

#endif /* PROBE_LOOP_H */
