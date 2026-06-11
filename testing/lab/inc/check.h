/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_CHECK_H_
#define ULAB_CHECK_H_

#include "bff.h"
#include "model.h"
#include "runtime.h"

typedef struct check_ctx {
    const scenario_t *scenario;
    world_t      *world;
    model_t      *model;
    bff_client_t *bff;
    runtime_t    *runtime;
} check_ctx_t;

typedef struct {
    char name[ULAB_MAX_NAME];
    int  passed;
    int  skipped;
    char detail[ULAB_MAX_ERR];
} check_result_t;

int check_run(check_ctx_t *ctx,
              const check_spec_t *check,
              check_result_t *res,
              ulab_error_t *err);
void check_list_supported(void);

#endif /* ULAB_CHECK_H_ */
