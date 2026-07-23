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

#include <stdint.h>
#include <math.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#define ULAB_BYTES_PER_MB              1048576ULL
#define ULAB_DEFAULT_OVERHEAD_PERCENT  15u

static uint64_t mb_to_bytes(uint64_t mb) {
    return mb * ULAB_BYTES_PER_MB;
}

static uint64_t percent_of(uint64_t value, uint32_t percent) {
    return (value * percent) / 100u;
}

static uint32_t max_u32(uint32_t a, uint32_t b) {
    return a > b ? a : b;
}

static int usage_within_range(uint64_t expected_bytes,
                              uint64_t actual_bytes,
                              uint32_t lower_tol_percent,
                              uint32_t upper_tol_percent,
                              uint64_t *lower_out,
                              uint64_t *upper_out) {
    uint64_t lower_delta;
    uint64_t upper_delta;
    uint64_t lower;
    uint64_t upper;

    lower_delta = percent_of(expected_bytes, lower_tol_percent);
    upper_delta = percent_of(expected_bytes, upper_tol_percent);

    lower = expected_bytes > lower_delta ? expected_bytes - lower_delta : 0;
    upper = expected_bytes + upper_delta;

    if (lower_out != NULL) {
        *lower_out = lower;
    }

    if (upper_out != NULL) {
        *upper_out = upper;
    }

    return actual_bytes >= lower && actual_bytes <= upper;
}

static double usage_value_for_unit(uint64_t expected_mb, const char *unit) {
    if (unit != NULL &&
        (ulab_streq(unit, "MB") || ulab_streq(unit, "mb") ||
         ulab_streq(unit, "megabytes"))) {
        return (double)expected_mb;
    }
    if (unit != NULL &&
        (ulab_streq(unit, "GB") || ulab_streq(unit, "gb") ||
         ulab_streq(unit, "gigabytes"))) {
        return (double)expected_mb / 1024.0;
    }
    return (double)mb_to_bytes(expected_mb);
}

static uint64_t model_subscriber_usage(const check_ctx_t *ctx,
                                       const char *subscriber_ref) {
    uint64_t total;
    size_t i;

    total = 0;
    for (i = 0; i < ctx->world->ue_count; i++) {
        const ue_t *ue;
        const model_ue_t *model_ue_value;

        ue = &ctx->world->ues[i];
        if (!ulab_streq(ue->subscriber_ref, subscriber_ref)) continue;
        model_ue_value = model_ue_const(ctx->model, ue->ref);
        if (model_ue_value != NULL) total += model_ue_value->used_mb;
    }
    return total;
}

static uint64_t model_package_usage(const check_ctx_t *ctx,
                                    const package_t *pkg) {
    uint64_t total;
    size_t i;

    total = 0;
    for (i = 0; i < ctx->world->ue_count; i++) {
        const ue_t *ue;
        const model_ue_t *model_ue_value;

        ue = &ctx->world->ues[i];
        if (!ulab_streq(ue->package_ref, pkg->ref)) continue;
        model_ue_value = model_ue_const(ctx->model, ue->ref);
        if (model_ue_value != NULL) total += model_ue_value->used_mb;
    }
    return total;
}

static int check_usage_aggregate(check_ctx_t *ctx,
                                 const check_spec_t *check,
                                 check_result_t *res,
                                 ulab_error_t *err) {
    selector_result_t selected;
    network_t *network;
    const char *key;
    const char *scope_key;
    char scope_value[ULAB_MAX_ID];
    uint64_t expected_mb;
    bff_kpi_value_t value;
    time_t deadline;
    unsigned int poll;
    int found;
    int matched;

    memset(&selected, 0, sizeof(selected));
    memset(&value, 0, sizeof(value));
    scope_value[0] = '\0';
    expected_mb = 0;
    network = NULL;
    key = check->key;
    scope_key = NULL;

    if (ulab_streq(check->target, "sim")) {
        ue_t *ue;
        const model_ue_t *model_ue_value;

        if (selector_resolve_ues(ctx->world, &check->ues, &selected, err) ||
            selected.count != 1) {
            selector_result_free(&selected);
            snprintf(err->msg, sizeof(err->msg),
                     "SIM usage aggregate requires exactly one UE");
            return ULAB_ERR;
        }
        ue = &ctx->world->ues[selected.idx[0]];
        network = world_network_by_ref(ctx->world, ue->network_ref);
        model_ue_value = model_ue_const(ctx->model, ue->ref);
        expected_mb = model_ue_value ? model_ue_value->used_mb : 0;
        ulab_copy(scope_value, sizeof(scope_value), ue->bff_id);
        scope_key = "sim_id";
        if (key[0] == '\0') key = "USAGE_BY_SIM";
    } else if (ulab_streq(check->target, "subscriber")) {
        subscriber_t *subscriber;

        subscriber = check->ref[0] ?
            world_subscriber_by_ref(ctx->world, check->ref) : NULL;
        if (subscriber == NULL && check->ues.kind != SEL_NONE &&
            selector_resolve_ues(ctx->world, &check->ues,
                                 &selected, err) == ULAB_OK &&
            selected.count > 0) {
            subscriber = world_subscriber_by_ref(
                ctx->world,
                ctx->world->ues[selected.idx[0]].subscriber_ref);
        }
        if (subscriber == NULL) {
            selector_result_free(&selected);
            snprintf(err->msg, sizeof(err->msg),
                     "subscriber usage aggregate requires ref or UE");
            return ULAB_ERR;
        }
        network = world_network_by_ref(ctx->world,
                                       subscriber->network_ref);
        expected_mb = model_subscriber_usage(ctx, subscriber->ref);
        ulab_copy(scope_value, sizeof(scope_value), subscriber->bff_id);
        scope_key = "subscriber_id";
        if (key[0] == '\0') key = "USAGE_BY_SUBSCRIBER";
    } else if (ulab_streq(check->target, "package")) {
        package_t *pkg;

        if (check->networks.kind != SEL_NONE &&
            selector_resolve_networks(ctx->world, &check->networks,
                                      &selected, err) == ULAB_OK &&
            selected.count > 0) {
            network = &ctx->world->networks[selected.idx[0]];
        } else if (ctx->world->network_count == 1) {
            network = &ctx->world->networks[0];
        }
        pkg = network ? world_package_for_network(ctx->world,
                                                   check->package_ref,
                                                   network->ref) : NULL;
        if (pkg == NULL) {
            selector_result_free(&selected);
            snprintf(err->msg, sizeof(err->msg),
                     "package usage aggregate requires package and network");
            return ULAB_ERR;
        }
        expected_mb = model_package_usage(ctx, pkg);
        ulab_copy(scope_value, sizeof(scope_value), pkg->bff_id);
        scope_key = "package_id";
        if (key[0] == '\0') key = "USAGE_BY_PACKAGE";
    } else if (ulab_streq(check->target, "site")) {
        site_t *site;

        if (selector_resolve_sites(ctx->world, &check->sites,
                                   &selected, err) || selected.count != 1) {
            selector_result_free(&selected);
            snprintf(err->msg, sizeof(err->msg),
                     "site usage aggregate requires exactly one site");
            return ULAB_ERR;
        }
        site = &ctx->world->sites[selected.idx[0]];
        network = world_network_by_ref(ctx->world, site->network_ref);
        expected_mb = model_site_usage(ctx->model, ctx->world, site->ref);
        ulab_copy(scope_value, sizeof(scope_value), site->bff_id);
        scope_key = "site_id";
        if (key[0] == '\0') key = "USAGE_BY_SITE";
    } else if (ulab_streq(check->target, "network")) {
        if (selector_resolve_networks(ctx->world, &check->networks,
                                      &selected, err) ||
            selected.count != 1) {
            selector_result_free(&selected);
            snprintf(err->msg, sizeof(err->msg),
                     "network usage aggregate requires exactly one network");
            return ULAB_ERR;
        }
        network = &ctx->world->networks[selected.idx[0]];
        expected_mb = model_network_usage(ctx->model, ctx->world,
                                          network->ref);
        ulab_copy(scope_value, sizeof(scope_value), network->bff_id);
        scope_key = "network_id";
        if (key[0] == '\0') key = "USAGE_BY_NETWORK";
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported usage aggregate target %.64s", check->target);
        return ULAB_ERR;
    }
    selector_result_free(&selected);
    if (network == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "usage aggregate could not resolve network");
        return ULAB_ERR;
    }

    deadline = time(NULL) + (time_t)(check->timeout_seconds ?
                                    check->timeout_seconds : 300u);
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    found = 0;
    matched = 0;
    do {
        double expected;
        double tolerance;

        if (bff_get_kpi_value(ctx->bff, key,
                              check->span[0] ? check->span : "daily",
                              check->op[0] ? check->op : "SUM",
                              network->bff_id, scope_key, scope_value,
                              &value, &found, err)) {
            return ULAB_ERR;
        }
        expected = usage_value_for_unit(expected_mb, value.unit);
        tolerance = expected *
            (double)(check->tolerance_percent ?
                     check->tolerance_percent : 15u) / 100.0;
        if (tolerance < check->tolerance_value) {
            tolerance = check->tolerance_value;
        }
        matched = found && fabs(value.value - expected) <= tolerance;
        if (!matched && time(NULL) < deadline) {
            sleep(poll > 60u ? 60u : poll);
        }
    } while (!matched && time(NULL) < deadline);

    res->passed = matched;
    snprintf(res->detail, sizeof(res->detail),
             "target=%.16s kpi=%.40s found=%s actual=%.6g unit=%.16s "
             "expected_mb=%llu scope=%.24s:%.48s",
             check->target, key, found ? "true" : "false", value.value,
             value.unit, (unsigned long long)expected_mb, scope_key,
             scope_value);
    return ULAB_OK;
}

int check_usage(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok;
    uint32_t lower_tol;
    uint32_t upper_tol;
    char first_fail[512];

    if (check->type == CHECK_USAGE_AGGREGATE) {
        return check_usage_aggregate(ctx, check, res, err);
    }

    ok = 0;
    first_fail[0] = '\0';

    lower_tol = check->tolerance_percent ? check->tolerance_percent : 2;
    upper_tol = max_u32(lower_tol, ULAB_DEFAULT_OVERHEAD_PERCENT);

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        network_t *network;
        model_ue_t *mu;
        uint64_t expected_mb;
        uint64_t expected_bytes;
        uint64_t actual_bytes;
        uint64_t lower_bytes;
        uint64_t upper_bytes;

        ue = &ctx->world->ues[ues.idx[i]];
        network = world_network_by_ref(ctx->world, ue->network_ref);
        mu = model_ue(ctx->model, ue->ref);

        if (mu == NULL) {
            continue;
        }

        expected_mb = check->expected_used_mb;
        if (expected_mb == 0 || ulab_streq(check->expected, "from_model")) {
            expected_mb = mu->used_mb;
        }

        expected_bytes = mb_to_bytes(expected_mb);
        actual_bytes = 0;
        lower_bytes = 0;
        upper_bytes = 0;

        if (network == NULL ||
            bff_get_sim_usage(ctx->bff, ue, network, &actual_bytes, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }

        if (usage_within_range(expected_bytes, actual_bytes,
                               lower_tol, upper_tol,
                               &lower_bytes, &upper_bytes)) {
            ok++;
        } else if (first_fail[0] == '\0') {
            snprintf(first_fail, sizeof(first_fail),
                     " first_fail=%.64s iccid=%.32s expected=%lluB actual=%lluB "
                     "range=%llu..%lluB",
                     ue->ref,
                     ue->iccid,
                     (unsigned long long)expected_bytes,
                     (unsigned long long)actual_bytes,
                     (unsigned long long)lower_bytes,
                     (unsigned long long)upper_bytes);
        }
    }

    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail),
             "usage=%zu/%zu tol=%u%% upper=%u%%%s",
             ok, ues.count, lower_tol, upper_tol, first_fail);

    selector_result_free(&ues);
    return ULAB_OK;
}
