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

int check_count(check_ctx_t *ctx, const check_spec_t *check,
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
    if (bff_query_count(ctx->bff, check->target, ctx->world, &actual, err)) {
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
