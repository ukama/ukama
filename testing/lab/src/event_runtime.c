/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "event.h"
#include "log.h"
#include "selector.h"
#include "util.h"

#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static int event_state_enabled(const event_spec_t *event, int *enabled,
                               ulab_error_t *err) {
    if (ulab_streq(event->status, "on") ||
        ulab_streq(event->status, "active") ||
        ulab_streq(event->status, "enabled") ||
        ulab_streq(event->status, "true")) {
        *enabled = 1;
        return ULAB_OK;
    }

    if (ulab_streq(event->status, "off") ||
        ulab_streq(event->status, "inactive") ||
        ulab_streq(event->status, "disabled") ||
        ulab_streq(event->status, "false")) {
        *enabled = 0;
        return ULAB_OK;
    }

    snprintf(err->msg, sizeof(err->msg), "%s invalid state: %s",
             scenario_event_name(event->type), event->status);
    return ULAB_ERR;
}

static int event_wait_ues_attached(event_ctx_t *ctx,
                                    const event_spec_t *event,
                                    ulab_error_t *err) {
    selector_result_t res;
    const char *current;
    char *previous;
    char timeout[32];
    int rc;

    previous = NULL;
    if (event->amount_mb > 300) {
        snprintf(err->msg, sizeof(err->msg),
                 "UE attach wait exceeds 300 seconds");
        return ULAB_ERR;
    }

    if (event->amount_mb > 0) {
        current = getenv("ULAB_UE_ATTACH_TIMEOUT");
        if (current != NULL) {
            previous = strdup(current);
            if (previous == NULL) {
                snprintf(err->msg, sizeof(err->msg),
                         "failed to save UE attach timeout");
                return ULAB_ERR;
            }
        }

        snprintf(timeout, sizeof(timeout), "%llu",
                 (unsigned long long)event->amount_mb);
        if (setenv("ULAB_UE_ATTACH_TIMEOUT", timeout, 1) != 0) {
            free(previous);
            snprintf(err->msg, sizeof(err->msg),
                     "failed to set UE attach timeout");
            return ULAB_ERR;
        }
    }

    rc = selector_resolve_ues(ctx->world, &event->ues, &res, err);
    if (rc == ULAB_OK) {
        rc = runtime_wait_ues_attached(ctx->runtime,
                                       ctx->world,
                                       &res,
                                       err);
    }
    selector_result_free(&res);

    if (event->amount_mb > 0) {
        if (previous != NULL) {
            setenv("ULAB_UE_ATTACH_TIMEOUT", previous, 1);
        } else {
            unsetenv("ULAB_UE_ATTACH_TIMEOUT");
        }
    }
    free(previous);

    return rc;
}

static int event_restart_nodes(event_ctx_t *ctx,
                               const event_spec_t *event,
                               ulab_error_t *err) {
    selector_result_t res;
    size_t i;

    if (selector_resolve_nodes(ctx->world, &event->nodes, &res, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < res.count; i++) {
        node_t *node;

        node = &ctx->world->nodes[res.idx[i]];
        if (bff_restart_node(ctx->bff, node, err)) {
            selector_result_free(&res);
            return ULAB_ERR;
        }
    }

    selector_result_free(&res);
    return ULAB_OK;
}

static int event_software_update(event_ctx_t *ctx,
                                 const event_spec_t *event,
                                 ulab_error_t *err) {
    selector_result_t res;
    size_t i;

    memset(&res, 0, sizeof(res));

    if (event->app[0] == '\0' || event->tag[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "software_update requires app and tag");
        return ULAB_ERR;
    }

    if (selector_resolve_nodes(ctx->world, &event->nodes, &res, err)) {
        return ULAB_ERR;
    }

    if (res.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "software_update selected no nodes");
        selector_result_free(&res);
        return ULAB_ERR;
    }

    for (i = 0; i < res.count; i++) {
        node_t *node;

        node = &ctx->world->nodes[res.idx[i]];
        ulab_status("SOFTWARE", "update node=%s app=%s tag=%s",
                    node->bff_id, event->app, event->tag);

        if (bff_update_software(ctx->bff, node, event->app,
                                event->tag, err)) {
            selector_result_free(&res);
            return ULAB_ERR;
        }
    }

    selector_result_free(&res);
    return ULAB_OK;
}

static int event_toggle_service(event_ctx_t *ctx,
                                const event_spec_t *event,
                                ulab_error_t *err) {
    selector_result_t res;
    size_t i;
    int enabled;

    if (event_state_enabled(event, &enabled, err)) {
        return ULAB_ERR;
    }

    if (selector_resolve_sites(ctx->world, &event->sites, &res, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < res.count; i++) {
        site_t *site;

        site = &ctx->world->sites[res.idx[i]];
        if (bff_toggle_site_service(ctx->bff, site, enabled, err)) {
            selector_result_free(&res);
            return ULAB_ERR;
        }
    }

    selector_result_free(&res);
    return runtime_set_service(ctx->runtime, enabled, err);
}

static int event_toggle_radio(event_ctx_t *ctx,
                              const event_spec_t *event,
                              ulab_error_t *err) {
    selector_result_t res;
    size_t i;
    int enabled;

    if (event_state_enabled(event, &enabled, err)) {
        return ULAB_ERR;
    }

    if (selector_resolve_sites(ctx->world, &event->sites, &res, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < res.count; i++) {
        site_t *site;

        site = &ctx->world->sites[res.idx[i]];
        if (bff_toggle_site_radio(ctx->bff, site, enabled, err)) {
            selector_result_free(&res);
            return ULAB_ERR;
        }
    }

    selector_result_free(&res);
    return runtime_set_radio(ctx->runtime, enabled, err);
}

int event_runtime(event_ctx_t *ctx,
                  const event_spec_t *event,
                  ulab_error_t *err) {
    selector_result_t res;
    int rc;

    switch (event->type) {
    case EVT_TOGGLE_SERVICE:
        return event_toggle_service(ctx, event, err);

    case EVT_TOGGLE_RADIO:
        return event_toggle_radio(ctx, event, err);

    case EVT_MARK_NODE_OFFLINE:
        return runtime_mark_node_offline(ctx->runtime, err);

    case EVT_RESTORE_NODE:
        return runtime_restore_nodes(ctx->runtime, err);

    case EVT_SOFTWARE_UPDATE:
        return event_software_update(ctx, event, err);

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
        return event_wait_ues_attached(ctx, event, err);

    case EVT_WAIT:
        if (event->amount_mb > 300) {
            snprintf(err->msg, sizeof(err->msg),
                     "wait exceeds 300 seconds");
            return ULAB_ERR;
        }
        ulab_status("WAIT", "%llu sec",
                    (unsigned long long)event->amount_mb);
        sleep((unsigned int)event->amount_mb);
        return ULAB_OK;

    case EVT_RESTART_NODES:
        return event_restart_nodes(ctx, event, err);

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
