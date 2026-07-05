/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "event.h"
#include "selector.h"
#include "util.h"

int event_runtime(event_ctx_t *ctx,
                  const event_spec_t *event,
                  ulab_error_t *err) {

    selector_result_t res;
    int rc;

    switch (event->type) {
    case EVT_TOGGLE_SERVICE:
        return runtime_set_service(ctx->runtime,
                                   !ulab_streq(event->status, "off") &&
                                   !ulab_streq(event->status, "inactive") &&
                                   !ulab_streq(event->status, "disabled"),
                                   err);

    case EVT_TOGGLE_RADIO:
        return runtime_set_radio(ctx->runtime,
                                 !ulab_streq(event->status, "off") &&
                                 !ulab_streq(event->status, "inactive") &&
                                 !ulab_streq(event->status, "disabled"),
                                 err);

    case EVT_MARK_NODE_OFFLINE:
        return runtime_mark_node_offline(ctx->runtime, err);

    case EVT_RESTORE_NODE:
        return runtime_restore_nodes(ctx->runtime, err);

    case EVT_SOFTWARE_UPDATE:
        return runtime_set_node_version(ctx->runtime, event->status, err);

    case EVT_RESTART_SITE:
        rc = selector_resolve_nodes(ctx->world, &event->nodes, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_restart_nodes(ctx->runtime,
                                       ctx->world,
                                       &res,
                                       err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_START_UES:
        rc = selector_resolve_ues(ctx->world, &event->ues, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_ensure_media(ctx->runtime, err);
        }
        if (rc == ULAB_OK) {
            rc = runtime_build_and_start_ues(NULL,
                                             ctx->runtime,
                                             ctx->world,
                                             &res,
                                             err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_WAIT_UES_ATTACHED:
        rc = selector_resolve_ues(ctx->world, &event->ues, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_wait_ues_attached(ctx->runtime,
                                           ctx->world,
                                           &res,
                                           err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_RESTART_NODES:
        rc = selector_resolve_nodes(ctx->world, &event->nodes, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_restart_nodes(ctx->runtime,
                                       ctx->world,
                                       &res,
                                       err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_WAIT_NODES_READY:
        rc = selector_resolve_nodes(ctx->world, &event->nodes, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_wait_nodes_ready(ctx->runtime,
                                          ctx->world,
                                          &res,
                                          err);
        }
        selector_result_free(&res);
        return rc;

    default:
        snprintf(err->msg, sizeof(err->msg), "not a runtime event");
        return ULAB_ERR;
    }
}
