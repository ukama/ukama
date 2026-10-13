/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#ifndef ULAB_HARDWARE_OBSERVE_H_
#define ULAB_HARDWARE_OBSERVE_H_
#include "check.h"

int hardware_wait_nodes(bff_client_t *bff, runtime_t *runtime, const world_t *world,
                        const selector_result_t *nodes, unsigned int timeout, ulab_error_t *err);
int hardware_check_runtime(check_ctx_t *ctx, const check_spec_t *check, check_result_t *result,
                           ulab_error_t *err);
#endif
