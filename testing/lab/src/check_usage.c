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

int check_usage(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok = 0;
    uint32_t tol = check->tolerance_percent ? check->tolerance_percent : 2;

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue = &ctx->world->ues[ues.idx[i]];
        model_ue_t *mu = model_ue(ctx->model, ue->ref);
        uint64_t actual = 0;
        uint64_t expected = check->expected_used_mb;

        if (mu == NULL) continue;
        if (expected == 0 || ulab_streq(check->expected, "from_model")) {
            expected = mu->used_mb;
        }
        if (bff_get_sim_usage(ctx->bff, ue, &actual, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        if (ulab_within_pct(expected, actual, tol)) {
            ok++;
        }
    }
    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail), "usage=%zu/%zu tol=%u%%", ok,
             ues.count, tol);
    selector_result_free(&ues);
    return ULAB_OK;
}
