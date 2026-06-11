/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_EVENT_H_
#define ULAB_EVENT_H_

#include "bff.h"
#include "model.h"
#include "runtime.h"

typedef struct check_ctx check_ctx_t;

typedef struct {
    scenario_t *scenario;
    world_t *world;
    model_t *model;
    bff_client_t *bff;
    runtime_t *runtime;
    const char *phaseName;
    check_ctx_t *checks;
} event_ctx_t;

int event_run(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err);
void event_list_supported(void);

#endif /* ULAB_EVENT_H_ */
