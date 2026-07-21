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
#include <stdio.h>

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

int check_usage(check_ctx_t *ctx, const check_spec_t *check,
                check_result_t *res, ulab_error_t *err) {
    selector_result_t ues;
    size_t i;
    size_t ok;
    uint32_t lower_tol;
    uint32_t upper_tol;
    char first_fail[512];

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
