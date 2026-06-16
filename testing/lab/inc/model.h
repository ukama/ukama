/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_MODEL_H_
#define ULAB_MODEL_H_

#include "world.h"

typedef struct {
    char     ue_ref[ULAB_MAX_REF];
    char     sim_id[ULAB_MAX_ID];
    char     package_ref[ULAB_MAX_REF];

    uint64_t package_mb;
    uint64_t used_mb;
    uint64_t remaining_mb;
} model_ue_t;

typedef struct {
    model_ue_t *ues;
    size_t     ue_count;
} model_t;

int model_init(model_t *m, const world_t *w);
void model_free(model_t *m);
int model_sync_world(model_t *m, const world_t *w);
int model_add_usage(model_t *m, const char *ue_ref, uint64_t amount_mb);
model_ue_t *model_ue(model_t *m, const char *ue_ref);
const model_ue_t *model_ue_const(const model_t *m, const char *ue_ref);
uint64_t model_site_usage(const model_t *m, const world_t *w,
                          const char *site_ref);
uint64_t model_network_usage(const model_t *m, const world_t *w,
                             const char *network_ref);
int model_balance_non_negative(const model_t *m);
int model_write_json(const model_t *m, const char *path);

#endif /* ULAB_MODEL_H_ */
