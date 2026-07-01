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

static const char *backend_id_for_ref(world_t *w, const char *view,
                                      const char *ref) {
    network_t *net;
    site_t *site;
    node_t *node;
    ue_t *ue;

    if (ulab_streq(view, "networks")) {
        net = world_network_by_ref(w, ref);
        return net ? net->bff_id : NULL;
    }
    if (ulab_streq(view, "sites")) {
        site = world_site_by_ref(w, ref);
        return site ? site->bff_id : NULL;
    }
    if (ulab_streq(view, "nodes")) {
        node = world_node_by_ref(w, ref);
        return node ? node->bff_id : NULL;
    }
    if (ulab_streq(view, "sims") || ulab_streq(view, "ues")) {
        ue = world_ue_by_ref(w, ref);
        return ue ? ue->bff_id : NULL;
    }

    return NULL;
}

int check_list(check_ctx_t *ctx, const check_spec_t *check,
               check_result_t *res, ulab_error_t *err) {
    const char *id;
    int found;
    int want_found;

    if (check->view[0] == '\0' || check->ref[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing view/ref", scenario_check_name(check->type));
        return ULAB_ERR;
    }

    id = backend_id_for_ref(ctx->world, check->view, check->ref);
    if (id == NULL || id[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "cannot resolve %s ref=%s", check->view, check->ref);
        return ULAB_ERR;
    }

    found = 0;
    if (bff_backend_contains(ctx->bff, check->view, id, ctx->world,
                             &found, err)) {
        return ULAB_ERR;
    }

    want_found = check->type == CHECK_LIST_CONTAINS;
    res->passed = found == want_found;
    snprintf(res->detail, sizeof(res->detail),
             "view=%s ref=%s backend_id=%s found=%s expected=%s",
             check->view, check->ref, id,
             found ? "true" : "false",
             want_found ? "true" : "false");

    return ULAB_OK;
}
