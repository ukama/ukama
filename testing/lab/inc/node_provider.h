/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_NODE_PROVIDER_H_
#define ULAB_NODE_PROVIDER_H_

#include "runtime.h"

int node_provider_build(runtime_t *rt, world_t *w, ulab_error_t *err);
int node_provider_start(runtime_t *rt, world_t *w, ulab_error_t *err);
int node_provider_wait_ready(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err);
int node_provider_restart(runtime_t *rt, const world_t *w,
                          const selector_result_t *nodes,
                          ulab_error_t *err);
int node_provider_status(runtime_t *rt, const world_t *w,
                         const selector_result_t *nodes,
                         ulab_error_t *err);
int node_provider_stop(runtime_t *rt, const world_t *w, ulab_error_t *err);
int node_provider_cleanup(runtime_t *rt, const world_t *w, ulab_error_t *err);

#endif /* ULAB_NODE_PROVIDER_H_ */
