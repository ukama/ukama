/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "check.h"
#include "selector.h"
#include "util.h"

static int value_matches(double actual, const check_spec_t *check) {
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

static network_t *first_network(check_ctx_t *ctx,
                                const check_spec_t *check,
                                ulab_error_t *err) {
    selector_result_t selected;
    network_t *network;

    network = NULL;
    memset(&selected, 0, sizeof(selected));
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
    if (ctx->world->network_count == 1) {
        return &ctx->world->networks[0];
    }
    snprintf(err->msg, sizeof(err->msg),
             "%s requires a network selector",
             scenario_check_name(check->type));
    return NULL;
}

static int check_revenue(check_ctx_t *ctx, const check_spec_t *check,
                         check_result_t *res, ulab_error_t *err) {
    network_t *network;
    bff_revenue_summary_t summary;
    double actual;

    network = first_network(ctx, check, err);
    if (network == NULL) return ULAB_ERR;
    if (bff_get_revenue_summary(ctx->bff, network, &summary, err)) {
        return ULAB_ERR;
    }
    if (ulab_streq(check->column, "total_paid")) {
        actual = summary.total_paid;
    } else if (ulab_streq(check->column, "total_pending")) {
        actual = summary.total_pending;
    } else if (ulab_streq(check->column, "month_paid")) {
        actual = summary.month_paid;
    } else if (ulab_streq(check->column, "previous_month_paid")) {
        actual = summary.previous_month_paid;
    } else if (ulab_streq(check->column, "mom_pct")) {
        actual = summary.month_over_month_percent;
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown revenue summary column %.64s", check->column);
        return ULAB_ERR;
    }
    res->passed = value_matches(actual, check);
    snprintf(res->detail, sizeof(res->detail),
             "network=%.32s column=%.32s actual=%.6g expected=%.6g",
             network->ref, check->column, actual, check->expected_value);
    return ULAB_OK;
}

static int check_package_dashboard(check_ctx_t *ctx,
                                   const check_spec_t *check,
                                   check_result_t *res,
                                   ulab_error_t *err) {
    network_t *network;
    double actual;

    network = first_network(ctx, check, err);
    if (network == NULL) return ULAB_ERR;
    actual = 0;
    if (ulab_streq(check->column, "mrr") ||
        ulab_streq(check->column, "arpu")) {
        bff_package_kpis_t kpis;

        if (bff_get_package_kpis(ctx->bff, network, &kpis, err)) {
            return ULAB_ERR;
        }
        if (ulab_streq(check->column, "mrr")) {
            if (!kpis.has_mrr) {
                snprintf(err->msg, sizeof(err->msg),
                         "MRR KPI is unavailable");
                return ULAB_ERR;
            }
            actual = kpis.mrr;
        } else {
            if (!kpis.has_arpu) {
                snprintf(err->msg, sizeof(err->msg),
                         "ARPU KPI is unavailable");
                return ULAB_ERR;
            }
            actual = kpis.arpu;
        }
    } else {
        package_t *pkg;
        bff_package_metrics_t metrics;
        int found;

        if (check->package_ref[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "plan metric requires package");
            return ULAB_ERR;
        }
        pkg = world_package_for_network(ctx->world, check->package_ref,
                                        network->ref);
        if (pkg == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "unknown package %.64s", check->package_ref);
            return ULAB_ERR;
        }
        found = 0;
        if (bff_get_package_metrics(ctx->bff, pkg, network, &metrics,
                                    &found, err)) {
            return ULAB_ERR;
        }
        if (!found) {
            snprintf(err->msg, sizeof(err->msg),
                     "package performance is missing package %.64s",
                     check->package_ref);
            return ULAB_ERR;
        }
        if (ulab_streq(check->column, "revenue")) {
            actual = metrics.revenue;
        } else if (ulab_streq(check->column, "attach_count")) {
            if (!metrics.has_attach_count) {
                snprintf(err->msg, sizeof(err->msg),
                         "active package assignment count is unavailable");
                return ULAB_ERR;
            }
            actual = (double)metrics.attach_count;
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "unknown package metric column %.64s",
                     check->column);
            return ULAB_ERR;
        }
    }
    res->passed = value_matches(actual, check);
    snprintf(res->detail, sizeof(res->detail),
             "network=%.32s column=%.32s actual=%.6g expected=%.6g",
             network->ref, check->column, actual, check->expected_value);
    return ULAB_OK;
}

static int check_network_overview(check_ctx_t *ctx,
                                  const check_spec_t *check,
                                  check_result_t *res,
                                  ulab_error_t *err) {
    network_t *network;
    bff_network_summary_t summary;
    double actual;

    network = first_network(ctx, check, err);
    if (network == NULL) return ULAB_ERR;
    if (bff_get_network_summary(ctx->bff, network, &summary, err)) {
        return ULAB_ERR;
    }
    if (ulab_streq(check->column, "subscribers_total")) {
        actual = (double)summary.subscribers_total;
    } else if (ulab_streq(check->column, "subscribers_active")) {
        actual = (double)summary.subscribers_active;
    } else if (ulab_streq(check->column, "subscribers_inactive")) {
        actual = (double)summary.subscribers_inactive;
    } else if (ulab_streq(check->column, "sites_total")) {
        actual = (double)summary.sites_total;
    } else if (ulab_streq(check->column, "nodes_total")) {
        actual = (double)summary.nodes_total;
    } else if (ulab_streq(check->column, "nodes_online")) {
        actual = (double)summary.nodes_online;
    } else if (ulab_streq(check->column, "nodes_offline")) {
        actual = (double)summary.nodes_offline;
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown network summary column %.64s", check->column);
        return ULAB_ERR;
    }
    res->passed = value_matches(actual, check);
    snprintf(res->detail, sizeof(res->detail),
             "network=%.32s column=%.32s actual=%.6g expected=%.6g",
             network->ref, check->column, actual, check->expected_value);
    return ULAB_OK;
}

static int check_inventory(check_ctx_t *ctx, const check_spec_t *check,
                           check_result_t *res, ulab_error_t *err) {
    if (ulab_streq(check->target, "component")) {
        bff_inventory_summary_t inventory;

        if (bff_get_inventory_summary(ctx->bff, ctx->sim_type,
                                      &inventory, err)) {
            return ULAB_ERR;
        }
        res->passed = inventory.component_total ==
            inventory.component_category_total;
        snprintf(res->detail, sizeof(res->detail),
                 "components=%u category_sum=%u",
                 inventory.component_total,
                 inventory.component_category_total);
        return ULAB_OK;
    }
    if (ulab_streq(check->target, "sim")) {
        bff_inventory_summary_t inventory;

        if (bff_get_inventory_summary(ctx->bff, ctx->sim_type,
                                      &inventory, err)) {
            return ULAB_ERR;
        }
        res->passed = inventory.sim_total == inventory.sim_pool_total &&
            inventory.sim_available == inventory.sim_pool_available &&
            inventory.sim_consumed == inventory.sim_pool_consumed &&
            inventory.sim_total == inventory.sim_available +
                inventory.sim_consumed;
        snprintf(res->detail, sizeof(res->detail),
                 "inventory=%u/%u/%u pool=%u/%u/%u",
                 inventory.sim_total, inventory.sim_available,
                 inventory.sim_consumed, inventory.sim_pool_total,
                 inventory.sim_pool_available, inventory.sim_pool_consumed);
        return ULAB_OK;
    }
    if (ulab_streq(check->target, "node")) {
        network_t *network;
        bff_network_summary_t summary;
        uint32_t nodes_list_count;

        network = first_network(ctx, check, err);
        if (network == NULL) return ULAB_ERR;
        if (bff_get_network_summary(ctx->bff, network, &summary, err) ||
            bff_get_nodes_count(ctx->bff, network,
                                &nodes_list_count, err)) {
            return ULAB_ERR;
        }
        res->passed = summary.nodes_total == nodes_list_count &&
            summary.nodes_total == summary.nodes_online +
                summary.nodes_offline;
        snprintf(res->detail, sizeof(res->detail),
                 "nodes_total=%u online=%u offline=%u nodes_list=%u",
                 summary.nodes_total, summary.nodes_online,
                 summary.nodes_offline, nodes_list_count);
        return ULAB_OK;
    }
    snprintf(err->msg, sizeof(err->msg),
             "unsupported inventory target %.64s", check->target);
    return ULAB_ERR;
}

static int subscriber_seen(const selector_result_t *ues,
                           const world_t *world,
                           size_t position) {
    size_t i;
    const char *subscriber_ref;

    subscriber_ref = world->ues[ues->idx[position]].subscriber_ref;
    for (i = 0; i < position; i++) {
        if (ulab_streq(world->ues[ues->idx[i]].subscriber_ref,
                       subscriber_ref)) {
            return 1;
        }
    }
    return 0;
}

static int check_subscriber_billing(check_ctx_t *ctx,
                                    const check_spec_t *check,
                                    check_result_t *res,
                                    ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t checked;
    size_t matched;
    double observed;

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    checked = 0;
    matched = 0;
    observed = 0;
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        subscriber_t *subscriber;
        bff_subscriber_billing_t billing;
        int count_ok;
        int amount_ok;

        if (subscriber_seen(&ues, ctx->world, i)) continue;
        ue = &ctx->world->ues[ues.idx[i]];
        subscriber = world_subscriber_by_ref(ctx->world,
                                             ue->subscriber_ref);
        if (subscriber == NULL ||
            bff_get_subscriber_payment_summary(ctx->bff, subscriber,
                                               &billing, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        checked++;
        observed += billing.settled_amount;
        count_ok = !check->has_expected_count ||
            billing.settled_count == check->expected_count;
        amount_ok = !check->has_expected_value ||
            value_matches(billing.settled_amount, check);
        if (count_ok && amount_ok) matched++;
    }
    res->passed = checked > 0 && matched == checked;
    snprintf(res->detail, sizeof(res->detail),
             "subscribers=%zu/%zu settled_amount=%.6g expected_count=%u "
             "expected_amount=%.6g",
             matched, checked, observed, check->expected_count,
             check->expected_value);
    selector_result_free(&ues);
    return ULAB_OK;
}

static int settled_status(const char *status) {
    return ulab_streq(status, "success") || ulab_streq(status, "paid") ||
        ulab_streq(status, "completed") ||
        ulab_streq(status, "processed");
}

static int check_payment_entitlements(check_ctx_t *ctx,
                                      const check_spec_t *check,
                                      check_result_t *res,
                                      ulab_error_t *err) {
    selector_result_t ues;
    package_t *pkg;
    bff_payment_t payments[ULAB_MAX_BFF_PAYMENTS];
    size_t payment_count;
    size_t settled_count;
    size_t entitlement_count;
    size_t i;

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    if (ues.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "payment reconciliation selected no SIMs");
        selector_result_free(&ues);
        return ULAB_ERR;
    }
    pkg = world_package_for_network(
        ctx->world, check->package_ref,
        ctx->world->ues[ues.idx[0]].network_ref);
    if (pkg == NULL ||
        bff_get_package_payments(ctx->bff, pkg, payments,
                                 ULAB_MAX_BFF_PAYMENTS,
                                 &payment_count, err)) {
        selector_result_free(&ues);
        return ULAB_ERR;
    }
    settled_count = 0;
    for (i = 0; i < payment_count; i++) {
        if (settled_status(payments[i].status)) settled_count++;
    }
    entitlement_count = 0;
    for (i = 0; i < ues.count; i++) {
        bff_sim_package_t assignments[ULAB_MAX_BFF_SIM_PACKAGES];
        size_t assignment_count;
        size_t j;

        if (bff_get_sim_packages(ctx->bff,
                                 &ctx->world->ues[ues.idx[i]],
                                 assignments, ULAB_MAX_BFF_SIM_PACKAGES,
                                 &assignment_count, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        for (j = 0; j < assignment_count; j++) {
            if (ulab_streq(assignments[j].package_id, pkg->bff_id)) {
                entitlement_count++;
            }
        }
    }
    res->passed = settled_count == entitlement_count &&
        entitlement_count > 0;
    snprintf(res->detail, sizeof(res->detail),
             "package=%.48s settled_payments=%zu entitlements=%zu",
             check->package_ref, settled_count, entitlement_count);
    selector_result_free(&ues);
    return ULAB_OK;
}

int check_business(check_ctx_t *ctx, const check_spec_t *check,
                   check_result_t *res, ulab_error_t *err) {
    switch (check->type) {
    case CHECK_REVENUE_SUMMARY:
        return check_revenue(ctx, check, res, err);
    case CHECK_SUBSCRIBER_BILLING_SUMMARY:
        return check_subscriber_billing(ctx, check, res, err);
    case CHECK_PAYMENT_ENTITLEMENT_RECONCILES:
        return check_payment_entitlements(ctx, check, res, err);
    case CHECK_PACKAGE_DASHBOARD_METRIC:
        return check_package_dashboard(ctx, check, res, err);
    case CHECK_NETWORK_OVERVIEW_METRIC:
        return check_network_overview(ctx, check, res, err);
    case CHECK_CONSOLE_INVENTORY_RECONCILES:
        return check_inventory(ctx, check, res, err);
    default:
        snprintf(err->msg, sizeof(err->msg),
                 "not a business GraphQL check");
        return ULAB_ERR;
    }
}
