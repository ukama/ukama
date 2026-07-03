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
#include <stdlib.h>

#define ULAB_BYTES_PER_MIB 1048576ULL
#define ULAB_DEFAULT_WIRE_OVERHEAD_PCT 15u

static uint64_t sat_add_u64(uint64_t a, uint64_t b) {
    if (UINT64_MAX - a < b) {
        return UINT64_MAX;
    }

    return a + b;
}

static uint64_t sat_sub_u64(uint64_t a, uint64_t b) {
    if (a < b) {
        return 0;
    }

    return a - b;
}

static uint64_t mib_to_bytes(uint64_t mib) {
    if (mib > UINT64_MAX / ULAB_BYTES_PER_MIB) {
        return UINT64_MAX;
    }

    return mib * ULAB_BYTES_PER_MIB;
}

static uint64_t pct_of_u64(uint64_t value, uint32_t pct) {
    uint64_t q;
    uint64_t r;

    q = value / 100u;
    r = value % 100u;

    if (pct == 0) {
        return 0;
    }

    if (q > UINT64_MAX / pct) {
        return UINT64_MAX;
    }

    return (q * pct) + ((r * pct) / 100u);
}

static uint32_t usage_upper_tolerance_percent(uint32_t lower_tol) {
    const char *env;
    uint32_t configured;

    env = getenv("ULAB_USAGE_WIRE_OVERHEAD_PERCENT");
    configured = 0;

    if (env != NULL && env[0] != '\0' &&
        ulab_parse_u32(env, &configured) == ULAB_OK) {
        return configured < lower_tol ? lower_tol : configured;
    }

    return lower_tol > ULAB_DEFAULT_WIRE_OVERHEAD_PCT ?
           lower_tol : ULAB_DEFAULT_WIRE_OVERHEAD_PCT;
}

static uint64_t normalize_usage_to_bytes(uint64_t raw, uint64_t expected_mib) {
    uint64_t mb_like_limit;

    if (raw == 0) {
        return 0;
    }

    /*
     * The current ASR/CDR usage API returns bytes. Older/stubbed test
     * responses may return MiB. If the number is small enough to look like
     * an MiB value for this expectation, normalize it to bytes here so the
     * check stays compatible with both forms.
     */
    if (expected_mib > 0) {
        if (expected_mib > UINT64_MAX / 1024ULL) {
            mb_like_limit = UINT64_MAX;
        } else {
            mb_like_limit = expected_mib * 1024ULL;
        }

        if (raw < mb_like_limit) {
            return mib_to_bytes(raw);
        }
    }

    return raw;
}

static int usage_in_range(uint64_t expected_bytes,
                          uint64_t actual_bytes,
                          uint32_t lower_tol,
                          uint32_t upper_tol,
                          uint64_t *lower_out,
                          uint64_t *upper_out) {
    uint64_t lower_delta;
    uint64_t upper_delta;
    uint64_t lower;
    uint64_t upper;

    lower_delta = pct_of_u64(expected_bytes, lower_tol);
    upper_delta = pct_of_u64(expected_bytes, upper_tol);
    lower = sat_sub_u64(expected_bytes, lower_delta);
    upper = sat_add_u64(expected_bytes, upper_delta);

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
    char first_fail[768];
    size_t i;
    size_t ok;
    uint32_t lower_tol;
    uint32_t upper_tol;

    ok = 0;
    first_fail[0] = '\0';
    lower_tol = check->tolerance_percent ? check->tolerance_percent : 2;
    upper_tol = usage_upper_tolerance_percent(lower_tol);

    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < ues.count; i++) {
        ue_t *ue;
        model_ue_t *mu;
        uint64_t actual_raw;
        uint64_t actual_bytes;
        uint64_t actual_mib;
        uint64_t expected_mib;
        uint64_t expected_bytes;
        uint64_t lower_bytes;
        uint64_t upper_bytes;

        ue = &ctx->world->ues[ues.idx[i]];
        mu = model_ue(ctx->model, ue->ref);
        actual_raw = 0;
        expected_mib = check->expected_used_mb;
        lower_bytes = 0;
        upper_bytes = 0;

        if (mu == NULL) {
            continue;
        }

        if (expected_mib == 0 || ulab_streq(check->expected, "from_model")) {
            expected_mib = mu->used_mb;
        }

        if (bff_get_sim_usage(ctx->bff, ue, &actual_raw, err)) {
            selector_result_free(&ues);
            return ULAB_ERR;
        }

        expected_bytes = mib_to_bytes(expected_mib);
        actual_bytes = normalize_usage_to_bytes(actual_raw, expected_mib);
        actual_mib = actual_bytes / ULAB_BYTES_PER_MIB;

        if (usage_in_range(expected_bytes, actual_bytes, lower_tol,
                           upper_tol, &lower_bytes, &upper_bytes)) {
            ok++;
        } else if (first_fail[0] == '\0') {
            snprintf(first_fail, sizeof(first_fail),
                     " first_fail=%.64s iccid=%.32s expected=%lluMiB actual=%lluMiB "
                     "actual_bytes=%llu raw=%llu range=%llu..%lluB",
                     ue->ref, ue->iccid,
                     (unsigned long long)expected_mib,
                     (unsigned long long)actual_mib,
                     (unsigned long long)actual_bytes,
                     (unsigned long long)actual_raw,
                     (unsigned long long)lower_bytes,
                     (unsigned long long)upper_bytes);
        }
    }

    res->passed = ok == ues.count;
    snprintf(res->detail, sizeof(res->detail),
             "usage=%zu/%zu tol=%u%% overhead=%u%%%s", ok,
             ues.count, lower_tol, upper_tol, first_fail);

    selector_result_free(&ues);
    return ULAB_OK;
}
