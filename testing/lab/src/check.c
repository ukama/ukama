/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "check.h"
#include <string.h>
#include <stdio.h>

int check_count(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err);
int check_runtime(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err);
int check_usage(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err);
int check_package(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err);
int check_dashboard(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err);

static void res_init(check_result_t *res, const check_spec_t *check) {
    memset(res, 0, sizeof(*res));
    snprintf(res->name, sizeof(res->name), "%s",
             scenario_check_name(check->type));
}

int check_run(check_ctx_t *ctx, const check_spec_t *check,
              check_result_t *res, ulab_error_t *err) {
    res_init(res, check);
    switch (check->type) {
    case CHECK_COUNT:
        return check_count(ctx, check, res, err);
    case CHECK_NODE_READY:
    case CHECK_UE_ATTACHED:
    case CHECK_NODE_STATE:
        return check_runtime(ctx, check, res, err);
    case CHECK_USAGE_PER_SIM:
    case CHECK_USAGE_SAMPLE:
        return check_usage(ctx, check, res, err);
    case CHECK_PACKAGE_ACTIVE:
    case CHECK_PACKAGE_REMAINING:
    case CHECK_BALANCE_NON_NEGATIVE:
        return check_package(ctx, check, res, err);
    case CHECK_DASHBOARD_LOADS:
        return check_dashboard(ctx, check, res, err);
    default:
        snprintf(err->msg, sizeof(err->msg), "unsupported check");
        return ULAB_ERR;
    }
}

void check_list_supported(void) {
    scenario_list_checks();
}
