/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "provider_virtual.h"
#include "ulab.h"

int virtual_provider_build(runtime_t *rt, world_t *w, ulab_error_t *err) {
    (void)rt;
    (void)w;
    (void)err;

    /* Current virtual site script builds images as part of start. */
    return ULAB_OK;
}

int virtual_provider_start(runtime_t *rt, world_t *w, ulab_error_t *err) {
    return runtime_virtual_build_and_start_sites(rt->repo, rt, w, err);
}

int virtual_provider_wait_ready(runtime_t *rt, const world_t *w,
                                const selector_result_t *nodes,
                                ulab_error_t *err) {
    return runtime_virtual_wait_nodes_ready(rt, w, nodes, err);
}

int virtual_provider_restart(runtime_t *rt, const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err) {
    return runtime_virtual_restart_nodes(rt, w, nodes, err);
}

int virtual_provider_status(runtime_t *rt, const world_t *w,
                            const selector_result_t *nodes,
                            ulab_error_t *err) {
    (void)rt;
    (void)w;
    (void)nodes;
    (void)err;

    return ULAB_OK;
}

int virtual_provider_stop(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    return runtime_virtual_cleanup_infra(rt, w, err);
}

int virtual_provider_cleanup(runtime_t *rt, const world_t *w,
                             ulab_error_t *err) {
    return runtime_virtual_cleanup_infra(rt, w, err);
}
