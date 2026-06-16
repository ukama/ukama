/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "check.h"
#include "selector.h"
#include <stdio.h>

int check_dashboard(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err) {
    selector_result_t nets;
    size_t i;
    size_t ok = 0;

    if (selector_resolve_networks(ctx->world, &check->networks, &nets, err)) {
        return ULAB_ERR;
    }
    for (i = 0; i < nets.count; i++) {
        if (bff_network_overview_loads(ctx->bff,
            &ctx->world->networks[nets.idx[i]], err) == ULAB_OK) {
            ok++;
        } else {
            selector_result_free(&nets);
            return ULAB_ERR;
        }
    }
    res->passed = ok == nets.count;
    snprintf(res->detail, sizeof(res->detail), "dashboard_loads=%zu/%zu", ok,
             nets.count);
    selector_result_free(&nets);
    return ULAB_OK;
}
