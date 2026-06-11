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

int check_package(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok = 0;

    if (check->type == CHECK_BALANCE_NON_NEGATIVE) {
        res->passed = model_balance_non_negative(ctx->model);
        snprintf(res->detail, sizeof(res->detail),
                 "model balance non-negative");
        return ULAB_OK;
    }
    if (check->type == CHECK_PACKAGE_REMAINING) {
        res->skipped = 1;
        snprintf(res->detail, sizeof(res->detail),
                 "BFF remaining balance query not available yet");
        return ULAB_OK;
    }
    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < ues.count; i++) {
        ue_t *ue = &ctx->world->ues[ues.idx[i]];
        package_t *pkg = world_package_by_ref(ctx->world,
            check->package_ref[0] ? check->package_ref : ue->package_ref);
        int active = 0;
        if (pkg == NULL) continue;
        if (bff_get_packages_for_sim(ctx->bff, ue, pkg->bff_id, &active,
            err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        if (active) ok++;
    }
    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail), "package_active=%zu/%zu", ok,
             ues.count);
    selector_result_free(&ues);
    return ULAB_OK;
}
