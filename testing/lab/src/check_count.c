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

int check_backend_count(check_ctx_t *ctx, const check_spec_t *check,
                        check_result_t *res, ulab_error_t *err) {
    size_t actual = 0;
    uint32_t expected = 0;

    if (check->target[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "backend_count missing target");
        return ULAB_ERR;
    }

    if (ulab_parse_u32(check->expected, &expected) != ULAB_OK &&
        !ulab_streq(check->expected, "from_world")) {
        snprintf(err->msg, sizeof(err->msg), "backend_count missing expected");
        return ULAB_ERR;
    }

    if (bff_backend_count(ctx->bff, check->target, ctx->world,
                          &actual, err)) {
        return ULAB_ERR;
    }

    if (ulab_streq(check->expected, "from_world")) {
        if (ulab_streq(check->target, "networks")) {
            expected = (uint32_t)ctx->world->network_count;
        } else if (ulab_streq(check->target, "sites")) {
            expected = (uint32_t)ctx->world->site_count;
        } else if (ulab_streq(check->target, "nodes")) {
            expected = (uint32_t)ctx->world->node_count;
        } else if (ulab_streq(check->target, "sims") ||
                   ulab_streq(check->target, "ues")) {
            expected = (uint32_t)ctx->world->ue_count;
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "from_world unsupported for target: %s",
                     check->target);
            return ULAB_ERR;
        }
    }

    res->passed = actual == expected;
    snprintf(res->detail, sizeof(res->detail),
             "target=%s expected=%u actual=%zu", check->target,
             expected, actual);

    return ULAB_OK;
}
