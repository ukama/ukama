/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_SELECTOR_H_
#define ULAB_SELECTOR_H_

#include "world.h"

typedef enum {
    SEL_OBJ_UE = 0,
    SEL_OBJ_NODE,
    SEL_OBJ_SITE,
    SEL_OBJ_NETWORK
} selector_obj_t;

typedef struct {
    size_t *idx;
    size_t count;
} selector_result_t;

int selector_resolve_ues(const world_t *w,
                         const selector_t *sel,
                         selector_result_t *out,
                         ulab_error_t *err);
int selector_resolve_nodes(const world_t *w,
                           const selector_t *sel,
                           selector_result_t *out,
                           ulab_error_t *err);
int selector_resolve_sites(const world_t *w,
                           const selector_t *sel,
                           selector_result_t *out,
                           ulab_error_t *err);
int selector_resolve_networks(const world_t *w,
                              const selector_t *sel,
                              selector_result_t *out,
                              ulab_error_t *err);
void selector_result_free(selector_result_t *r);

#endif /* ULAB_SELECTOR_H_ */
