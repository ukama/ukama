/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#ifndef ULAB_RUNTIME_HARDWARE_H_
#define ULAB_RUNTIME_HARDWARE_H_
#include "runtime.h"
#include "hardware.h"

int runtime_is_hardware(const runtime_t *rt);
int hardware_runtime_open(runtime_t *rt, const hardware_config_t *config, const world_t *world,
                          ulab_error_t *err);
int hardware_runtime_healthy(runtime_t *rt, ulab_error_t *err);
int hardware_runtime_power(runtime_t *rt, int enabled, ulab_error_t *err);
int hardware_runtime_media(runtime_t *rt, int enabled, ulab_error_t *err);
int hardware_runtime_ues(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                         const char *action, uint64_t amount_mb, ulab_error_t *err);
int hardware_runtime_wait(runtime_t *rt, world_t *world, const selector_result_t *ues, int attached,
                          unsigned int timeout, ulab_error_t *err);
int hardware_runtime_sessions(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                              ulab_error_t *err);
int hardware_runtime_cleanup_ues(runtime_t *rt, const world_t *world, const selector_result_t *ues,
                                 ulab_error_t *err);
int hardware_runtime_cleanup(runtime_t *rt, ulab_error_t *err);
void hardware_runtime_close(runtime_t *rt);

#endif
