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
#include <time.h>
#include <unistd.h>

#include "check.h"
#include "log.h"
#include "selector.h"
#include "util.h"

static package_t *resolve_package(world_t *world,
                                  const char *ref,
                                  const network_t *network) {
    package_t *pkg;

    pkg = network != NULL ?
        world_package_for_network(world, ref, network->ref) : NULL;
    if (pkg == NULL) {
        pkg = world_package_by_base_ref(world, ref);
    }
    return pkg;
}

static int resolve_networks(world_t *world,
                            const selector_t *selector,
                            selector_result_t *result,
                            ulab_error_t *err) {
    selector_t effective;

    effective = *selector;
    if (effective.kind == SEL_NONE) {
        effective.kind = SEL_ALL;
    }
    return selector_resolve_networks(world, &effective, result, err);
}

static int package_catalog_equals(check_ctx_t *ctx,
                                  const check_spec_t *check,
                                  check_result_t *res,
                                  ulab_error_t *err) {
    selector_result_t networks;
    size_t i;
    size_t matched;
    char detail[512];

    memset(&networks, 0, sizeof(networks));
    if (resolve_networks(ctx->world, &check->networks, &networks, err)) {
        return ULAB_ERR;
    }
    matched = 0;
    detail[0] = '\0';
    for (i = 0; i < networks.count; i++) {
        network_t *network;
        package_t *pkg;
        network_t *owner;
        bff_package_t actual;
        uint32_t expected_duration;
        int ok;

        network = &ctx->world->networks[networks.idx[i]];
        pkg = resolve_package(ctx->world, check->package_ref, network);
        if (pkg == NULL) {
            continue;
        }
        owner = pkg->network_ref[0] != '\0' ?
            world_network_by_ref(ctx->world, pkg->network_ref) : NULL;
        memset(&actual, 0, sizeof(actual));
        if (bff_get_package(ctx->bff, pkg, &actual, err)) {
            selector_result_free(&networks);
            return ULAB_ERR;
        }
        expected_duration = pkg->duration_minutes > 0 ?
            pkg->duration_minutes : pkg->duration_days * 1440u;
        ok = ulab_streq(actual.uuid, pkg->bff_id) &&
             ulab_streq(actual.name, pkg->name) &&
             actual.data_volume == pkg->data_mb &&
             actual.duration_minutes == expected_duration &&
             fabs(actual.amount - pkg->amount) <= 0.001 &&
             ulab_streq(actual.data_unit, "MB") &&
             ulab_streq(actual.currency, pkg->currency) &&
             ulab_streq(actual.country, pkg->country) &&
             actual.active == pkg->active &&
             ((owner == NULL && actual.network_id[0] == '\0') ||
              (owner != NULL &&
               ulab_streq(actual.network_id, owner->bff_id)));
        if (ok) {
            matched++;
        } else if (detail[0] == '\0') {
            snprintf(detail, sizeof(detail),
                     "pkg=%.96s actual(data=%llu duration=%u amount=%.2f "
                     "currency=%.24s active=%d network=%.96s)",
                     pkg->ref, (unsigned long long)actual.data_volume,
                     actual.duration_minutes, actual.amount,
                     actual.currency, actual.active, actual.network_id);
        }
    }
    res->passed = matched == networks.count;
    snprintf(res->detail, sizeof(res->detail),
             "package=%s catalog_match=%zu/%zu%s%s",
             check->package_ref, matched, networks.count,
             detail[0] ? " " : "", detail);
    selector_result_free(&networks);
    return ULAB_OK;
}

static int package_visibility(check_ctx_t *ctx,
                              const check_spec_t *check,
                              check_result_t *res,
                              ulab_error_t *err) {
    selector_result_t networks;
    size_t i;
    size_t matched;
    int expected_visible;

    memset(&networks, 0, sizeof(networks));
    if (resolve_networks(ctx->world, &check->networks, &networks, err)) {
        return ULAB_ERR;
    }
    expected_visible = check->type == CHECK_PACKAGE_VISIBLE;
    matched = 0;
    for (i = 0; i < networks.count; i++) {
        network_t *network;
        package_t *pkg;
        int visible;

        network = &ctx->world->networks[networks.idx[i]];
        pkg = resolve_package(ctx->world, check->package_ref, network);
        if (pkg == NULL) {
            selector_result_free(&networks);
            snprintf(err->msg, sizeof(err->msg),
                     "unknown package %.128s", check->package_ref);
            return ULAB_ERR;
        }
        visible = 0;
        if (bff_package_visible_for_network(ctx->bff, pkg, network,
                                            &visible, err)) {
            selector_result_free(&networks);
            return ULAB_ERR;
        }
        if (visible == expected_visible) {
            matched++;
        }
    }
    res->passed = matched == networks.count;
    snprintf(res->detail, sizeof(res->detail),
             "package=%s visibility=%s matched=%zu/%zu",
             check->package_ref,
             expected_visible ? "visible" : "hidden",
             matched, networks.count);
    selector_result_free(&networks);
    return ULAB_OK;
}

static int package_name_available(check_ctx_t *ctx,
                                  const check_spec_t *check,
                                  check_result_t *res,
                                  ulab_error_t *err) {
    package_t *pkg;
    int available;

    pkg = world_package_by_base_ref(ctx->world, check->package_ref);
    if (pkg == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown package %.128s", check->package_ref);
        return ULAB_ERR;
    }
    available = 0;
    if (bff_invalid_package_name_available(ctx->bff, pkg, check->variant,
                                           &available, err)) {
        return ULAB_ERR;
    }
    res->passed = available;
    snprintf(res->detail, sizeof(res->detail),
             "package=%s invalid_variant=%s name_available=%s",
             check->package_ref, check->variant,
             available ? "true" : "false");
    return ULAB_OK;
}

static int package_business_metrics(check_ctx_t *ctx,
                                    const check_spec_t *check,
                                    check_result_t *res,
                                    ulab_error_t *err) {
    selector_result_t networks;
    time_t deadline;
    unsigned int poll;
    size_t i;
    size_t matched;
    bff_package_metrics_t last;
    int last_found;

    memset(&networks, 0, sizeof(networks));
    memset(&last, 0, sizeof(last));
    last_found = 0;
    if (resolve_networks(ctx->world, &check->networks, &networks, err)) {
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)(check->timeout_seconds ?
                                    check->timeout_seconds : 300u);
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    if (check->has_expected_value && check->has_expected_count) {
        ulab_status("CHECK",
                    "package_business_metrics package=%s "
                    "expected_revenue=%.2f expected_attach=%u "
                    "timeout=%us poll=%us",
                    check->package_ref, check->expected_value,
                    check->expected_count,
                    check->timeout_seconds ? check->timeout_seconds : 300u,
                    poll);
    } else if (check->has_expected_value) {
        ulab_status("CHECK",
                    "package_business_metrics package=%s "
                    "expected_revenue=%.2f timeout=%us poll=%us",
                    check->package_ref, check->expected_value,
                    check->timeout_seconds ? check->timeout_seconds : 300u,
                    poll);
    } else if (check->has_expected_count) {
        ulab_status("CHECK",
                    "package_business_metrics package=%s "
                    "expected_attach=%u timeout=%us poll=%us",
                    check->package_ref, check->expected_count,
                    check->timeout_seconds ? check->timeout_seconds : 300u,
                    poll);
    } else {
        ulab_status("CHECK",
                    "package_business_metrics package=%s "
                    "timeout=%us poll=%us",
                    check->package_ref,
                    check->timeout_seconds ? check->timeout_seconds : 300u,
                    poll);
    }
    do {
        matched = 0;
        for (i = 0; i < networks.count; i++) {
            network_t *network;
            package_t *pkg;
            bff_package_metrics_t metrics;
            int found;
            int ok;

            network = &ctx->world->networks[networks.idx[i]];
            pkg = resolve_package(ctx->world, check->package_ref, network);
            if (pkg == NULL) {
                continue;
            }
            memset(&metrics, 0, sizeof(metrics));
            found = 0;
            if (bff_get_package_metrics(ctx->bff, pkg, network, &metrics,
                                        &found, err)) {
                selector_result_free(&networks);
                return ULAB_ERR;
            }
            last = metrics;
            last_found = found;
            ok = found;
            if (ok && check->has_expected_value) {
                ok = fabs(metrics.revenue - check->expected_value) <=
                    (check->tolerance_value > 0 ?
                     check->tolerance_value : 0.001);
            }
            if (ok && check->has_expected_count) {
                ok = metrics.has_attach_count &&
                    metrics.attach_count == check->expected_count;
            }
            if (ok) {
                matched++;
            }
        }
        if (matched != networks.count && time(NULL) < deadline) {
            sleep(poll > 60u ? 60u : poll);
        }
    } while (matched != networks.count && time(NULL) < deadline);

    res->passed = matched == networks.count;
    if (check->has_expected_value && check->has_expected_count) {
        snprintf(res->detail, sizeof(res->detail),
                 "package=%s found=%s revenue=%.2f "
                 "expected_revenue=%.2f attach=%u expected_attach=%u "
                 "matched=%zu/%zu",
                 check->package_ref, last_found ? "true" : "false",
                 last.revenue, check->expected_value, last.attach_count,
                 check->expected_count, matched, networks.count);
    } else if (check->has_expected_value) {
        snprintf(res->detail, sizeof(res->detail),
                 "package=%s found=%s revenue=%.2f "
                 "expected_revenue=%.2f matched=%zu/%zu",
                 check->package_ref, last_found ? "true" : "false",
                 last.revenue, check->expected_value, matched,
                 networks.count);
    } else {
        snprintf(res->detail, sizeof(res->detail),
                 "package=%s found=%s revenue=%.2f attach=%u "
                 "expected_attach=%u matched=%zu/%zu",
                 check->package_ref, last_found ? "true" : "false",
                 last.revenue, last.attach_count, check->expected_count,
                 matched, networks.count);
    }
    selector_result_free(&networks);
    return ULAB_OK;
}

static int assignment_compare(const void *a, const void *b) {
    const bff_sim_package_t *left;
    const bff_sim_package_t *right;

    left = a;
    right = b;
    return strcmp(left->start_date, right->start_date);
}

static int package_assignment_chain(check_ctx_t *ctx,
                                    const check_spec_t *check,
                                    check_result_t *res,
                                    ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t matched_ues;
    size_t last_target_count;

    memset(&ues, 0, sizeof(ues));
    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    matched_ues = 0;
    last_target_count = 0;
    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        package_t *target;
        package_t *before;
        bff_sim_package_t assignments[ULAB_MAX_BFF_SIM_PACKAGES];
        bff_sim_package_t target_rows[ULAB_MAX_BFF_SIM_PACKAGES];
        bff_sim_package_t before_rows[ULAB_MAX_BFF_SIM_PACKAGES];
        size_t assignment_count;
        size_t target_count;
        size_t before_count;
        size_t j;
        int ok;

        ue = &ctx->world->ues[ues.idx[i]];
        target = resolve_package(ctx->world, check->package_ref,
                                 world_network_by_ref(ctx->world,
                                                      ue->network_ref));
        before = check->other_package_ref[0] != '\0' ?
            resolve_package(ctx->world, check->other_package_ref,
                            world_network_by_ref(ctx->world,
                                                 ue->network_ref)) : NULL;
        if (target == NULL ||
            (check->other_package_ref[0] != '\0' && before == NULL)) {
            continue;
        }
        if (bff_get_sim_packages(ctx->bff, ue, assignments,
                                 ULAB_MAX_BFF_SIM_PACKAGES,
                                 &assignment_count, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        target_count = 0;
        before_count = 0;
        for (j = 0; j < assignment_count; j++) {
            if ((ulab_streq(check->expected, "all") ||
                 ulab_streq(assignments[j].package_id, target->bff_id)) &&
                target_count < ULAB_MAX_BFF_SIM_PACKAGES) {
                target_rows[target_count++] = assignments[j];
            }
            if (before != NULL &&
                ulab_streq(assignments[j].package_id, before->bff_id) &&
                before_count < ULAB_MAX_BFF_SIM_PACKAGES) {
                before_rows[before_count++] = assignments[j];
            }
        }
        qsort(target_rows, target_count, sizeof(target_rows[0]),
              assignment_compare);
        qsort(before_rows, before_count, sizeof(before_rows[0]),
              assignment_compare);
        ok = target_count > 0;
        if (ok && check->has_expected_count) {
            ok = target_count == check->expected_count;
        }
        for (j = 1; ok && j < target_count; j++) {
            ok = target_rows[j - 1].id[0] != '\0' &&
                 target_rows[j].id[0] != '\0' &&
                 !ulab_streq(target_rows[j - 1].id, target_rows[j].id) &&
                 strcmp(target_rows[j].start_date,
                        target_rows[j - 1].end_date) >= 0;
            if (ok && ulab_streq(check->expected, "contiguous")) {
                ok = ulab_streq(target_rows[j].start_date,
                                target_rows[j - 1].end_date);
            }
        }
        if (ok && before != NULL && !ulab_streq(check->expected, "all")) {
            ok = before_count > 0 &&
                 strcmp(target_rows[0].start_date,
                        before_rows[before_count - 1].end_date) >= 0;
            if (ok && ulab_streq(check->expected, "contiguous")) {
                ok = ulab_streq(target_rows[0].start_date,
                                before_rows[before_count - 1].end_date);
            }
        }
        if (ok) {
            matched_ues++;
        }
        last_target_count = target_count;
    }
    res->passed = matched_ues == ues.count;
    snprintf(res->detail, sizeof(res->detail),
             "package=%s after=%s target_count=%zu chain=%zu/%zu",
             check->package_ref,
             check->other_package_ref[0] ? check->other_package_ref : "none",
             last_target_count, matched_ues, ues.count);
    selector_result_free(&ues);
    return ULAB_OK;
}

static int sim_unallocated(check_ctx_t *ctx,
                           const check_spec_t *check,
                           check_result_t *res,
                           ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t matched;

    memset(&ues, 0, sizeof(ues));
    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }
    matched = 0;
    for (i = 0; i < ues.count; i++) {
        int unallocated;

        unallocated = 0;
        if (bff_sim_is_unallocated(ctx->bff,
                                   &ctx->world->ues[ues.idx[i]],
                                   ctx->sim_type, &unallocated, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }
        if (unallocated) {
            matched++;
        }
    }
    res->passed = matched == ues.count;
    snprintf(res->detail, sizeof(res->detail),
             "unallocated=%zu/%zu sim_type=%s", matched, ues.count,
             ctx->sim_type ? ctx->sim_type : "");
    selector_result_free(&ues);
    return ULAB_OK;
}

int check_data_package(check_ctx_t *ctx, const check_spec_t *check,
                       check_result_t *res, ulab_error_t *err) {
    switch (check->type) {
    case CHECK_PACKAGE_CATALOG_EQUALS:
        return package_catalog_equals(ctx, check, res, err);
    case CHECK_PACKAGE_VISIBLE:
    case CHECK_PACKAGE_HIDDEN:
        return package_visibility(ctx, check, res, err);
    case CHECK_PACKAGE_NAME_AVAILABLE:
        return package_name_available(ctx, check, res, err);
    case CHECK_PACKAGE_BUSINESS_METRICS:
        return package_business_metrics(ctx, check, res, err);
    case CHECK_PACKAGE_ASSIGNMENT_CHAIN:
        return package_assignment_chain(ctx, check, res, err);
    case CHECK_SIM_UNALLOCATED:
        return sim_unallocated(ctx, check, res, err);
    default:
        snprintf(err->msg, sizeof(err->msg),
                 "not a data-package check");
        return ULAB_ERR;
    }
}
