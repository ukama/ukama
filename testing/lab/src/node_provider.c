/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "node_provider.h"
#include "provider_virtual.h"
#include "util.h"

static int unsupported(runtime_t *rt, ulab_error_t *err) {
    snprintf(err->msg, sizeof(err->msg),
             "unsupported node provider: %s",
             rt && rt->provider[0] ? rt->provider : "");
    return ULAB_ERR;
}

static int is_virtual(runtime_t *rt) {
    return rt == NULL || rt->provider[0] == '\0' ||
           ulab_streq(rt->provider, "virtual");
}

int node_provider_build(runtime_t *rt, world_t *w, ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_build(rt, w, err);
    return unsupported(rt, err);
}

int node_provider_start(runtime_t *rt, world_t *w, ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_start(rt, w, err);
    return unsupported(rt, err);
}

int node_provider_wait_ready(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_wait_ready(rt, w, nodes, err);
    return unsupported(rt, err);
}

int node_provider_restart(runtime_t *rt, const world_t *w,
                          const selector_result_t *nodes,
                          ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_restart(rt, w, nodes, err);
    return unsupported(rt, err);
}

int node_provider_status(runtime_t *rt, const world_t *w,
                         const selector_result_t *nodes,
                         ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_status(rt, w, nodes, err);
    return unsupported(rt, err);
}

int node_provider_stop(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_stop(rt, w, err);
    return unsupported(rt, err);
}

int node_provider_cleanup(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    if (is_virtual(rt)) return virtual_provider_cleanup(rt, w, err);
    return unsupported(rt, err);
}
