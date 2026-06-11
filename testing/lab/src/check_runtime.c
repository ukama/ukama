/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "check.h"
#include "selector.h"
#include "util.h"
#include <stdio.h>

int check_runtime(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    selector_result_t sel;
    size_t i;
    size_t ok = 0;

    if (check->type == CHECK_UE_ATTACHED) {
        if (selector_resolve_ues(ctx->world, &check->ues, &sel, err)) {
            return ULAB_ERR;
        }
        for (i = 0; i < sel.count; i++) {
            if (ctx->world->ues[sel.idx[i]].attached) ok++;
        }
        res->passed = ok == sel.count;
        snprintf(res->detail, sizeof(res->detail), "attached=%zu/%zu", ok,
                 sel.count);
        selector_result_free(&sel);
        return ULAB_OK;
    }
    if (selector_resolve_nodes(ctx->world, &check->nodes, &sel, err)) {
        return ULAB_ERR;
    }
    if (check->type == CHECK_NODE_READY) {
        res->passed = 1;
        snprintf(res->detail, sizeof(res->detail),
                 "runtime nodes=%zu", sel.count);
        selector_result_free(&sel);
        return ULAB_OK;
    }
    for (i = 0; i < sel.count; i++) {
        bff_node_state_t st = {0};
        if (bff_get_node_state(ctx->bff, &ctx->world->nodes[sel.idx[i]],
            &st, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
        if (check->expected[0] == '\0' || ulab_streq(st.state,
            check->expected)) {
            ok++;
        }
    }
    res->passed = ok == sel.count;
    snprintf(res->detail, sizeof(res->detail), "node_state=%zu/%zu", ok,
             sel.count);
    selector_result_free(&sel);
    return ULAB_OK;
}
