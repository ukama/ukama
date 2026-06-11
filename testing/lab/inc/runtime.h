/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_RUNTIME_H_
#define ULAB_RUNTIME_H_

#include "selector.h"
#include <stdio.h>

typedef struct {
    char script_dir[ULAB_MAX_PATH];
    char run_dir[ULAB_MAX_PATH];
    FILE *logf;
} runtime_t;

int runtime_init(runtime_t *rt, const char *script_dir,
                 const char *run_dir);
void runtime_close(runtime_t *rt);
int runtime_start_nodes(runtime_t *rt, const world_t *w,
                        const selector_result_t *nodes,
                        ulab_error_t *err);
int runtime_wait_nodes_ready(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err);
int runtime_start_ues(runtime_t *rt, const world_t *w,
                      const selector_result_t *ues,
                      ulab_error_t *err);
int runtime_wait_ues_attached(runtime_t *rt, world_t *w,
                              const selector_result_t *ues,
                              ulab_error_t *err);
int runtime_generate_traffic(runtime_t *rt, const world_t *w,
                             const selector_result_t *ues,
                             uint64_t amount_mb, ulab_error_t *err);
int runtime_restart_nodes(runtime_t *rt, const world_t *w,
                          const selector_result_t *nodes,
                          ulab_error_t *err);

#endif /* ULAB_RUNTIME_H_ */
