/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "event.h"
#include "selector.h"
#include "util.h"

static package_t *event_package(event_ctx_t *ctx,
                                const event_spec_t *event,
                                ue_t *ue,
                                ulab_error_t *err) {
    package_t *pkg;

    if (event->package_ref[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "%s missing package", scenario_event_name(event->type));
        return NULL;
    }

    pkg = world_package_for_network(ctx->world, event->package_ref,
                                    ue->network_ref);
    if (pkg == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "unknown package %.128s for UE %.128s",
                 event->package_ref, ue->ref);
        return NULL;
    }

    return pkg;
}

static int event_add_package_to_sim(event_ctx_t *ctx,
                                    const event_spec_t *event,
                                    ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];
        package_t *pkg = event_package(ctx, event, ue, err);

        if (pkg == NULL) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }

        if (bff_add_package_to_sim(ctx->bff, ue, pkg, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int event_remove_package_from_sim(event_ctx_t *ctx,
                                         const event_spec_t *event,
                                         ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    (void)event;

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];

        if (bff_clear_sim_packages(ctx->bff, ue, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

static int event_set_sim_status(event_ctx_t *ctx,
                                const event_spec_t *event,
                                ulab_error_t *err) {
    selector_result_t sel;
    size_t i;

    if (event->status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg), "set_sim_status missing status");
        return ULAB_ERR;
    }

    if (selector_resolve_ues(ctx->world, &event->ues, &sel, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < sel.count; i++) {
        ue_t *ue = &ctx->world->ues[sel.idx[i]];

        if (bff_toggle_sim_status(ctx->bff, ue, event->status, err)) {
            selector_result_free(&sel);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sel);
    return ULAB_OK;
}

int event_bff(event_ctx_t *ctx, const event_spec_t *event,
              ulab_error_t *err) {
    switch (event->type) {
    case EVT_CREATE_UES:
        (void)ctx;
        snprintf(err->msg, sizeof(err->msg),
                 "create_ues is defined but not enabled in v1.0 runtime path");
        return ULAB_ERR;

    case EVT_ADD_PACKAGE_TO_SIM:
        return event_add_package_to_sim(ctx, event, err);

    case EVT_REMOVE_PACKAGE_FROM_SIM:
        return event_remove_package_from_sim(ctx, event, err);

    case EVT_SET_SIM_STATUS:
        return event_set_sim_status(ctx, event, err);

    default:
        snprintf(err->msg, sizeof(err->msg), "not a BFF event");
        return ULAB_ERR;
    }
}
