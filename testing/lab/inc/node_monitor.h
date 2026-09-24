/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_NODE_MONITOR_H_
#define ULAB_NODE_MONITOR_H_

#include "bff.h"

typedef struct node_monitor node_monitor_t;

int node_monitor_start(node_monitor_t **monitor, const bff_client_t *bff,
                       const world_t *world, const selector_t *nodes,
                       const char *connectivity, ulab_error_t *err);
int node_monitor_status(node_monitor_t *monitor, ulab_error_t *err);
int node_monitor_stop(node_monitor_t **monitor, ulab_error_t *err);

#endif /* ULAB_NODE_MONITOR_H_ */
