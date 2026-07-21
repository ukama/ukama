/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <math.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "check.h"
#include "selector.h"
#include "util.h"

static int compare_value(double actual, const check_spec_t *check) {
    double tolerance;

    tolerance = check->tolerance_value > 0 ?
        check->tolerance_value : 0.000001;
    if (ulab_streq(check->comparator, "greater_than") ||
        ulab_streq(check->comparator, "gt")) {
        return actual > check->expected_value;
    }
    if (ulab_streq(check->comparator, "greater_or_equal") ||
        ulab_streq(check->comparator, "gte")) {
        return actual + tolerance >= check->expected_value;
    }
    if (ulab_streq(check->comparator, "less_than") ||
        ulab_streq(check->comparator, "lt")) {
        return actual < check->expected_value;
    }
    if (ulab_streq(check->comparator, "less_or_equal") ||
        ulab_streq(check->comparator, "lte")) {
        return actual - tolerance <= check->expected_value;
    }
    if (ulab_streq(check->comparator, "not_equals")) {
        return fabs(actual - check->expected_value) > tolerance;
    }
    return fabs(actual - check->expected_value) <= tolerance;
}

static network_t *selected_network(check_ctx_t *ctx,
                                   const check_spec_t *check,
                                   ulab_error_t *err) {
    selector_result_t selected;
    network_t *network;

    network = NULL;
    if (check->networks.kind != SEL_NONE) {
        if (selector_resolve_networks(ctx->world, &check->networks,
                                      &selected, err)) {
            return NULL;
        }
        if (selected.count > 0) {
            network = &ctx->world->networks[selected.idx[0]];
        }
        selector_result_free(&selected);
        return network;
    }
    if (check->ues.kind != SEL_NONE) {
        if (selector_resolve_ues(ctx->world, &check->ues, &selected, err)) {
            return NULL;
        }
        if (selected.count > 0) {
            ue_t *ue;

            ue = &ctx->world->ues[selected.idx[0]];
            network = world_network_by_ref(ctx->world, ue->network_ref);
        }
        selector_result_free(&selected);
        return network;
    }
    if (check->sites.kind != SEL_NONE) {
        if (selector_resolve_sites(ctx->world, &check->sites,
                                  &selected, err)) {
            return NULL;
        }
        if (selected.count > 0) {
            site_t *site;

            site = &ctx->world->sites[selected.idx[0]];
            network = world_network_by_ref(ctx->world, site->network_ref);
        }
        selector_result_free(&selected);
    }
    return network;
}

static void resolve_scope(check_ctx_t *ctx, const check_spec_t *check,
                          network_t *network, char *scope_value,
                          size_t scope_value_len, ulab_error_t *err) {
    selector_result_t selected;

    scope_value[0] = '\0';
    if (check->scope_value[0] != '\0') {
        network_t *net;
        site_t *site;
        ue_t *ue;
        package_t *pkg;

        net = world_network_by_ref(ctx->world, check->scope_value);
        site = world_site_by_ref(ctx->world, check->scope_value);
        ue = world_ue_by_ref(ctx->world, check->scope_value);
        pkg = network ? world_package_for_network(ctx->world,
                                                   check->scope_value,
                                                   network->ref) : NULL;
        if (net != NULL) {
            ulab_copy(scope_value, scope_value_len, net->bff_id);
        } else if (site != NULL) {
            ulab_copy(scope_value, scope_value_len, site->bff_id);
        } else if (ue != NULL) {
            ulab_copy(scope_value, scope_value_len, ue->bff_id);
        } else if (pkg != NULL) {
            ulab_copy(scope_value, scope_value_len, pkg->bff_id);
        } else {
            ulab_copy(scope_value, scope_value_len, check->scope_value);
        }
        return;
    }
    if (ulab_streq(check->scope_key, "network_id") && network != NULL) {
        ulab_copy(scope_value, scope_value_len, network->bff_id);
        return;
    }
    if (ulab_streq(check->scope_key, "package_id") &&
        check->package_ref[0] != '\0' && network != NULL) {
        package_t *pkg;

        pkg = world_package_for_network(ctx->world, check->package_ref,
                                        network->ref);
        if (pkg != NULL) {
            ulab_copy(scope_value, scope_value_len, pkg->bff_id);
        }
        return;
    }
    if (ulab_streq(check->scope_key, "sim_id") &&
        check->ues.kind != SEL_NONE) {
        if (selector_resolve_ues(ctx->world, &check->ues, &selected, err) ==
            ULAB_OK) {
            if (selected.count > 0) {
                ulab_copy(scope_value, scope_value_len,
                          ctx->world->ues[selected.idx[0]].bff_id);
            }
            selector_result_free(&selected);
        }
        return;
    }
    if (ulab_streq(check->scope_key, "site_id") &&
        check->sites.kind != SEL_NONE) {
        if (selector_resolve_sites(ctx->world, &check->sites, &selected, err) ==
            ULAB_OK) {
            if (selected.count > 0) {
                ulab_copy(scope_value, scope_value_len,
                          ctx->world->sites[selected.idx[0]].bff_id);
            }
            selector_result_free(&selected);
        }
    }
}

static int check_kpi(check_ctx_t *ctx, const check_spec_t *check,
                     check_result_t *res, ulab_error_t *err) {
    network_t *network;
    char scope_value[ULAB_MAX_ID];
    bff_kpi_value_t value;
    time_t deadline;
    unsigned int poll;
    int found;
    int matched;

    network = selected_network(ctx, check, err);
    resolve_scope(ctx, check, network, scope_value, sizeof(scope_value), err);
    deadline = time(NULL) + (time_t)(check->timeout_seconds ?
                                    check->timeout_seconds : 300u);
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    matched = 0;
    memset(&value, 0, sizeof(value));

    do {
        found = 0;
        if (bff_get_kpi_value(ctx->bff, check->key,
                              check->span[0] ? check->span : "daily",
                              check->op,
                              network ? network->bff_id : NULL,
                              check->scope_key, scope_value,
                              &value, &found, err)) {
            return ULAB_ERR;
        }
        if (found) {
            if (check->type == CHECK_KPI_TREND) {
                matched = check->trend_direction[0] == '\0' ||
                    ulab_streq(value.trend.direction,
                               check->trend_direction);
                if (matched && check->has_expected_value) {
                    matched = compare_value(value.value, check);
                }
            } else {
                matched = check->has_expected_value &&
                    compare_value(value.value, check);
            }
        }
        if (!matched && time(NULL) < deadline) {
            sleep(poll > 60u ? 60u : poll);
        }
    } while (!matched && time(NULL) < deadline);

    res->passed = matched;
    snprintf(res->detail, sizeof(res->detail),
             "kpi=%.48s found=%s value=%.6g expected=%.6g op=%.16s "
             "scope=%.24s:%.64s computed_at=%.40s trend=%.16s",
             check->key, found ? "true" : "false", value.value,
             check->expected_value, check->op, check->scope_key, scope_value,
             value.computed_at, value.trend.direction);
    return ULAB_OK;
}

static int check_report_cell(check_ctx_t *ctx, const check_spec_t *check,
                             check_result_t *res, ulab_error_t *err) {
    network_t *network;
    package_t *pkg;
    time_t deadline;
    unsigned int poll;
    char unit[ULAB_MAX_REF];
    double actual;
    int found;
    int matched;

    network = selected_network(ctx, check, err);
    pkg = network ? world_package_for_network(ctx->world,
                                               check->package_ref,
                                               network->ref) : NULL;
    if (network == NULL || pkg == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "performance report check requires network and package");
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)(check->timeout_seconds ?
                                    check->timeout_seconds : 300u);
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    actual = 0;
    found = 0;
    matched = 0;
    unit[0] = '\0';

    do {
        if (bff_get_performance_report_cell(ctx->bff, check->report,
                                            check->span[0] ?
                                                check->span : "daily",
                                            network->bff_id, pkg->bff_id,
                                            check->column, &actual,
                                            unit, sizeof(unit),
                                            &found, err)) {
            return ULAB_ERR;
        }
        matched = found && check->has_expected_value &&
            compare_value(actual, check);
        if (!matched && time(NULL) < deadline) {
            sleep(poll > 60u ? 60u : poll);
        }
    } while (!matched && time(NULL) < deadline);

    res->passed = matched;
    snprintf(res->detail, sizeof(res->detail),
             "report=%.48s package=%.48s column=%.48s found=%s "
             "value=%.6g%.16s expected=%.6g",
             check->report, check->package_ref, check->column,
             found ? "true" : "false", actual, unit,
             check->expected_value);
    return ULAB_OK;
}

int check_analytics(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err) {
    if (check->type == CHECK_PERFORMANCE_REPORT_CELL) {
        return check_report_cell(ctx, check, res, err);
    }
    return check_kpi(ctx, check, res, err);
}
