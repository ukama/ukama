/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "check.h"
#include "util.h"

static const char *wanted_status(const check_spec_t *check) {
    if (check->status[0] != '\0') {
        return check->status;
    }
    return check->expected;
}

static int check_node_status(check_ctx_t *ctx, const check_spec_t *check,
                             const char *want, check_result_t *res,
                             ulab_error_t *err) {
    node_t *node;
    bff_node_state_t st;

    node = world_node_by_ref(ctx->world, check->ref);
    if (node == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "status_equals unknown node ref=%s", check->ref);
        return ULAB_ERR;
    }

    memset(&st, 0, sizeof(st));
    if (bff_get_node_state(ctx->bff, node, &st, err)) {
        return ULAB_ERR;
    }

    res->passed = ulab_streq(st.state, want) ||
                  ulab_streq(st.connectivity, want);
    snprintf(res->detail, sizeof(res->detail),
             "entity=node ref=%s expected=%s state=%s connectivity=%s",
             check->ref, want, st.state, st.connectivity);

    return ULAB_OK;
}

static int check_sim_status(check_ctx_t *ctx, const check_spec_t *check,
                            const char *want, check_result_t *res,
                            ulab_error_t *err) {
    ue_t *ue;
    int active;

    ue = world_ue_by_ref(ctx->world, check->ref);
    if (ue == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "status_equals unknown sim/ue ref=%s", check->ref);
        return ULAB_ERR;
    }

    if (ulab_streq(want, "started")) {
        res->passed = ue->started != 0;
    } else if (ulab_streq(want, "attached")) {
        res->passed = ue->attached != 0;
    } else if (ulab_streq(want, "active")) {
        active = 0;
        if (bff_get_packages_for_sim(ctx->bff, ue, NULL, &active, err)) {
            return ULAB_ERR;
        }
        res->passed = active != 0;
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported sim status: %s", want);
        return ULAB_ERR;
    }

    snprintf(res->detail, sizeof(res->detail),
             "entity=sim ref=%s expected=%s passed=%s",
             check->ref, want, res->passed ? "true" : "false");

    return ULAB_OK;
}

int check_status(check_ctx_t *ctx, const check_spec_t *check,
                 check_result_t *res, ulab_error_t *err) {
    const char *want;

    want = wanted_status(check);
    if (check->entity[0] == '\0' || check->ref[0] == '\0' ||
        want == NULL || want[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "status_equals missing entity/ref/status");
        return ULAB_ERR;
    }

    if (ulab_streq(check->entity, "node")) {
        return check_node_status(ctx, check, want, res, err);
    }
    if (ulab_streq(check->entity, "sim") ||
        ulab_streq(check->entity, "ue")) {
        return check_sim_status(ctx, check, want, res, err);
    }

    snprintf(err->msg, sizeof(err->msg),
             "unsupported status entity: %s", check->entity);
    return ULAB_ERR;
}
