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

int check_backend_count(check_ctx_t *ctx, const check_spec_t *check,
                        check_result_t *res, ulab_error_t *err);
int check_list(check_ctx_t *ctx, const check_spec_t *check,
               check_result_t *res, ulab_error_t *err);
int check_status(check_ctx_t *ctx, const check_spec_t *check,
                 check_result_t *res, ulab_error_t *err);
int check_runtime(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err);
int check_usage(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err);
int check_package(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err);
int check_dashboard(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err);
int check_payment(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err);
int check_analytics(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err);
int check_data_package(check_ctx_t *ctx, const check_spec_t *check,
                       check_result_t *res, ulab_error_t *err);
int check_business(check_ctx_t *ctx, const check_spec_t *check,
                   check_result_t *res, ulab_error_t *err);
int check_console(check_ctx_t *ctx, const check_spec_t *check,
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
    case CHECK_BACKEND_COUNT:
        return check_backend_count(ctx, check, res, err);
    case CHECK_LIST_CONTAINS:
    case CHECK_LIST_EXCLUDES:
        return check_list(ctx, check, res, err);
    case CHECK_STATUS_EQUALS:
        return check_status(ctx, check, res, err);
    case CHECK_NODE_READY:
    case CHECK_UE_ATTACHED:
    case CHECK_NODE_STATE:
    case CHECK_TRAFFIC_ALLOWED:
    case CHECK_TRAFFIC_BLOCKED:
    case CHECK_TRAFFIC_UNAVAILABLE:
    case CHECK_NODE_VERSION_EQUALS:
    case CHECK_NODE_HEALTH_OK:
    case CHECK_RELEASE_UNAVAILABLE:
        return check_runtime(ctx, check, res, err);
    case CHECK_LIST_COUNT_EQUALS:
    case CHECK_ENTITY_FIELDS_EQUAL:
    case CHECK_ENTITY_RECONCILES:
    case CHECK_NODE_STATUS_EQUALS:
    case CHECK_SOFTWARE_STATUS_EQUALS:
    case CHECK_SOFTWARE_COUNT_EQUALS:
    case CHECK_NODE_OPERATION_STATUS_EQUALS:
    case CHECK_SITE_OPERATION_STATUS_EQUALS:
    case CHECK_SITE_NODE_COUNTS_EQUALS:
    case CHECK_KPI_STATE_EQUALS:
    case CHECK_KPI_TIMESERIES:
        return check_console(ctx, check, res, err);
    case CHECK_USAGE_PER_SIM:
    case CHECK_USAGE_SAMPLE:
        return check_usage(ctx, check, res, err);
    case CHECK_PACKAGE_ACTIVE:
    case CHECK_PACKAGE_REMAINING:
    case CHECK_PACKAGE_STATE:
    case CHECK_PACKAGE_ASSIGNMENT_COUNT:
    case CHECK_BALANCE_NON_NEGATIVE:
        return check_package(ctx, check, res, err);
    case CHECK_PACKAGE_ASSIGNMENT_CHAIN:
    case CHECK_PACKAGE_FIELDS_EQUAL:
    case CHECK_PACKAGE_VISIBLE:
    case CHECK_PACKAGE_HIDDEN:
    case CHECK_PACKAGE_NAME_AVAILABLE:
    case CHECK_PACKAGE_BUSINESS_METRICS:
    case CHECK_SIM_UNALLOCATED:
        return check_data_package(ctx, check, res, err);
    case CHECK_PAYMENT_EQUALS:
    case CHECK_PAYMENT_COUNT:
        return check_payment(ctx, check, res, err);
    case CHECK_KPI_VALUE:
    case CHECK_KPI_TREND:
    case CHECK_KPI_CONTRACT:
    case CHECK_KPI_ROLLUP_CONSISTENCY:
    case CHECK_PERFORMANCE_REPORT_CELL:
    case CHECK_PERFORMANCE_REPORT_ROW:
        return check_analytics(ctx, check, res, err);
    case CHECK_REVENUE_SUMMARY:
    case CHECK_SUBSCRIBER_BILLING_SUMMARY:
    case CHECK_PAYMENT_ENTITLEMENT_RECONCILES:
    case CHECK_PACKAGE_DASHBOARD_METRIC:
    case CHECK_NETWORK_SUMMARY_METRIC:
    case CHECK_CONSOLE_INVENTORY_RECONCILES:
        return check_business(ctx, check, res, err);
    case CHECK_USAGE_AGGREGATE:
        return check_usage(ctx, check, res, err);
    case CHECK_DASHBOARD_LOADS:
    case CHECK_DASHBOARD_SECTION_OK:
        return check_dashboard(ctx, check, res, err);
    case CHECK_HISTORY_PRESERVED:
    case CHECK_AUDIT_EVENT_EXISTS:
    case CHECK_RELATIONSHIP_EXISTS:
    case CHECK_RELATIONSHIP_ENDED:
        res->passed = check->required ? 0 : 1;
        res->skipped = check->required ? 0 : 1;
        snprintf(res->detail, sizeof(res->detail),
                 "%s has no BFF-backed implementation%s",
                 scenario_check_name(check->type),
                 check->required ? " (required)" : "");
        (void)ctx;
        (void)err;
        return ULAB_OK;
    default:
        snprintf(err->msg, sizeof(err->msg), "unsupported check");
        return ULAB_ERR;
    }
}

void check_list_supported(void) {
    scenario_list_checks();
}
