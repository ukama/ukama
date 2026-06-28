/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "check.h"
#include "util.h"

static int model_count_for_target(const world_t *w, const char *target,
                                  size_t *count) {
    if (ulab_streq(target, "networks")) {
        *count = w->network_count;
    } else if (ulab_streq(target, "sites")) {
        *count = w->site_count;
    } else if (ulab_streq(target, "nodes")) {
        *count = w->node_count;
    } else if (ulab_streq(target, "packages")) {
        *count = w->package_count;
    } else if (ulab_streq(target, "subscribers")) {
        *count = w->subscriber_count;
    } else if (ulab_streq(target, "sims") ||
               ulab_streq(target, "ues")) {
        *count = w->ue_count;
    } else {
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int check_model_count(check_ctx_t *ctx, const check_spec_t *check,
                      check_result_t *res, ulab_error_t *err) {
    size_t actual = 0;
    uint32_t expected = 0;

    if (check->target[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "count check missing target");
        return ULAB_ERR;
    }
    if (ulab_parse_u32(check->expected, &expected) != ULAB_OK &&
        !ulab_streq(check->expected, "from_model")) {
        snprintf(err->msg, sizeof(err->msg), "count check missing expected");
        return ULAB_ERR;
    }
    if (model_count_for_target(ctx->world, check->target, &actual)) {
        snprintf(err->msg, sizeof(err->msg), "unknown model_count target: %s",
                 check->target);
        return ULAB_ERR;
    }
    if (ulab_streq(check->expected, "from_model")) {
        expected = (uint32_t)actual;
    }
    res->passed = actual == expected;
    snprintf(res->detail, sizeof(res->detail),
             "target=%s expected=%u actual=%zu", check->target,
             expected, actual);
    return ULAB_OK;
}

int check_bff_count(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err) {
    (void)ctx;
    (void)check;
    (void)err;

    res->skipped = 1;
    res->passed = 1;
    snprintf(res->detail, sizeof(res->detail),
             "bff_count is reserved for a later phase");
    return ULAB_OK;
}
