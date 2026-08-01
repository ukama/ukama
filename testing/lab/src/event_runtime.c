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
#include <strings.h>
#include <unistd.h>

#define NODE_CONNECTIVITY_WAIT_DEFAULT_SEC 180u
#define NODE_CONNECTIVITY_WAIT_MAX_SEC     900u
#define NODE_CONNECTIVITY_POLL_DEFAULT_SEC 2u
#define NODE_CONNECTIVITY_POLL_MAX_SEC     30u

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

static unsigned int node_connectivity_poll_seconds(void) {
    const char *value;
    uint32_t seconds;

    value = getenv("ULAB_NODE_CONNECTIVITY_POLL_SEC");
    if (value == NULL || value[0] == '\0' ||
        ulab_parse_u32(value, &seconds) != ULAB_OK || seconds == 0) {
        return NODE_CONNECTIVITY_POLL_DEFAULT_SEC;
    }

    if (seconds > NODE_CONNECTIVITY_POLL_MAX_SEC) {
        return NODE_CONNECTIVITY_POLL_MAX_SEC;
    }

    return seconds;
}

static int event_wait_node_connectivity(event_ctx_t *ctx,
                                        const event_spec_t *event,
                                        ulab_error_t *err) {
    selector_result_t res;
    bff_node_status_t status;
    ulab_error_t query_err;
    unsigned char *matched;
    unsigned int timeout;
    unsigned int poll;
    unsigned int elapsed;
    size_t matched_count;
    size_t i;
    char last_node[ULAB_MAX_ID];
    char last_connectivity[ULAB_MAX_REF];
    char last_state[ULAB_MAX_REF];
    char last_error[ULAB_MAX_ERR];

    memset(&res, 0, sizeof(res));
    memset(last_node, 0, sizeof(last_node));
    memset(last_connectivity, 0, sizeof(last_connectivity));
    memset(last_state, 0, sizeof(last_state));
    memset(last_error, 0, sizeof(last_error));

    if (event->status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity requires connectivity");
        return ULAB_ERR;
    }

    if (selector_resolve_nodes(ctx->world, &event->nodes, &res, err)) {
        return ULAB_ERR;
    }

    if (res.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity selected no nodes");
        selector_result_free(&res);
        return ULAB_ERR;
    }

    timeout = event->amount_mb > 0 ?
              (unsigned int)event->amount_mb :
              NODE_CONNECTIVITY_WAIT_DEFAULT_SEC;
    if (timeout > NODE_CONNECTIVITY_WAIT_MAX_SEC) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity exceeds %u seconds",
                 NODE_CONNECTIVITY_WAIT_MAX_SEC);
        selector_result_free(&res);
        return ULAB_ERR;
    }

    matched = calloc(res.count, sizeof(*matched));
    if (matched == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity allocation failed");
        selector_result_free(&res);
        return ULAB_ERR;
    }

    poll = node_connectivity_poll_seconds();
    elapsed = 0;
    matched_count = 0;

    ulab_status("NODE",
                "wait BFF connectivity=%s nodes=%zu timeout=%us",
                event->status, res.count, timeout);

    for (;;) {
        for (i = 0; i < res.count; i++) {
            node_t *node;

            if (matched[i]) {
                continue;
            }

            node = &ctx->world->nodes[res.idx[i]];
            memset(&status, 0, sizeof(status));
            memset(&query_err, 0, sizeof(query_err));

            ulab_copy(last_node, sizeof(last_node), node->bff_id);
            if (bff_get_node_status(ctx->bff, node, &status,
                                    &query_err)) {
                ulab_copy(last_error, sizeof(last_error), query_err.msg);
                continue;
            }

            last_error[0] = '\0';
            ulab_copy(last_connectivity, sizeof(last_connectivity),
                      status.connectivity);
            ulab_copy(last_state, sizeof(last_state), status.state);

            if (strcasecmp(status.connectivity, event->status) == 0) {
                matched[i] = 1;
                matched_count++;
                ulab_status("NODE",
                            "%s connectivity=%s state=%s",
                            node->bff_id, status.connectivity,
                            status.state);
            }
        }

        if (matched_count == res.count) {
            free(matched);
            selector_result_free(&res);
            return ULAB_OK;
        }

        if (elapsed >= timeout) {
            break;
        }

        sleep(poll);
        if (timeout - elapsed < poll) {
            elapsed = timeout;
        } else {
            elapsed += poll;
        }
    }

    if (last_error[0] != '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity expected=%.64s "
                 "matched=%zu/%zu node=%.192s timeout=%us "
                 "error=%.512s",
                 event->status, matched_count, res.count, last_node,
                 timeout, last_error);
    } else {
        snprintf(err->msg, sizeof(err->msg),
                 "wait_node_connectivity expected=%.64s "
                 "matched=%zu/%zu node=%.192s connectivity=%.64s "
                 "state=%.64s timeout=%us",
                 event->status, matched_count, res.count, last_node,
                 last_connectivity, last_state, timeout);
    }

    free(matched);
    selector_result_free(&res);
    return ULAB_ERR;
}

static int selector_result_add_unique(selector_result_t *res,
                                      size_t idx) {
    size_t *next;
    size_t i;

    for (i = 0; i < res->count; i++) {
        if (res->idx[i] == idx) {
            return ULAB_OK;
        }
    }

    next = realloc(res->idx, (res->count + 1) * sizeof(*next));
    if (next == NULL) {
        return ULAB_ERR;
    }
    res->idx = next;
    res->idx[res->count++] = idx;
    return ULAB_OK;
}

static int event_resolve_sites(event_ctx_t *ctx,
                               const event_spec_t *event,
                               selector_result_t *sites,
                               ulab_error_t *err) {
    selector_result_t nodes;
    size_t i;
    size_t j;

    memset(sites, 0, sizeof(*sites));
    if (event->sites.kind != SEL_NONE) {
        return selector_resolve_sites(ctx->world, &event->sites,
                                      sites, err);
    }

    if (event->nodes.kind == SEL_NONE) {
        snprintf(err->msg, sizeof(err->msg),
                 "%s requires a site selector",
                 scenario_event_name(event->type));
        return ULAB_ERR;
    }

    memset(&nodes, 0, sizeof(nodes));
    if (selector_resolve_nodes(ctx->world, &event->nodes, &nodes, err)) {
        return ULAB_ERR;
    }

    for (i = 0; i < nodes.count; i++) {
        const node_t *node;
        int found;

        node = &ctx->world->nodes[nodes.idx[i]];
        found = 0;
        for (j = 0; j < ctx->world->site_count; j++) {
            if (!ulab_streq(ctx->world->sites[j].ref, node->site_ref)) {
                continue;
            }
            if (selector_result_add_unique(sites, j)) {
                selector_result_free(&nodes);
                selector_result_free(sites);
                snprintf(err->msg, sizeof(err->msg),
                         "failed to collect selected sites");
                return ULAB_ERR;
            }
            found = 1;
            break;
        }
        if (!found) {
            selector_result_free(&nodes);
            selector_result_free(sites);
            snprintf(err->msg, sizeof(err->msg),
                     "node %.128s has unknown site %.128s",
                     node->ref, node->site_ref);
            return ULAB_ERR;
        }
    }

    selector_result_free(&nodes);
    return ULAB_OK;
}

static int event_restart_site(event_ctx_t *ctx,
                              const event_spec_t *event,
                              ulab_error_t *err) {
    selector_result_t sites;
    size_t i;

    if (event_resolve_sites(ctx, event, &sites, err)) {
        return ULAB_ERR;
    }
    if (sites.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "restart_site selected no sites");
        selector_result_free(&sites);
        return ULAB_ERR;
    }

    for (i = 0; i < sites.count; i++) {
        site_t *site;
        network_t *network;

        site = &ctx->world->sites[sites.idx[i]];
        network = world_network_by_ref(ctx->world, site->network_ref);
        if (network == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "restart_site unknown network %.128s for site %.128s",
                     site->network_ref, site->ref);
            selector_result_free(&sites);
            return ULAB_ERR;
        }

        ulab_status("SITE", "restart site=%s network=%s",
                    site->bff_id, network->bff_id);
        if (bff_restart_site(ctx->bff, site, network, err)) {
            selector_result_free(&sites);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sites);
    return ULAB_OK;
}

static int event_toggle_internet_switch(event_ctx_t *ctx,
                                        const event_spec_t *event,
                                        ulab_error_t *err) {
    selector_result_t sites;
    size_t i;
    int enabled;

    if (!event->has_port || event->port == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "toggle_internet_switch requires port > 0");
        return ULAB_ERR;
    }
    if (event_state_enabled(event, &enabled, err)) {
        return ULAB_ERR;
    }
    if (event_resolve_sites(ctx, event, &sites, err)) {
        return ULAB_ERR;
    }
    if (sites.count == 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "toggle_internet_switch selected no sites");
        selector_result_free(&sites);
        return ULAB_ERR;
    }

    for (i = 0; i < sites.count; i++) {
        site_t *site;

        site = &ctx->world->sites[sites.idx[i]];
        if (bff_toggle_internet_switch(ctx->bff, site, event->port,
                                       enabled, err)) {
            selector_result_free(&sites);
            return ULAB_ERR;
        }
    }

    selector_result_free(&sites);
    return ULAB_OK;
}

static int event_failure_control(event_ctx_t *ctx,
                                 const event_spec_t *event,
                                 ulab_error_t *err) {
    int enabled;

    if (event->target[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "failure_control requires target");
        return ULAB_ERR;
    }
    if (event_state_enabled(event, &enabled, err)) {
        return ULAB_ERR;
    }

    return runtime_set_failure_control(ctx->runtime, event->target,
                                       enabled, err);
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

    case EVT_TOGGLE_INTERNET_SWITCH:
        return event_toggle_internet_switch(ctx, event, err);

    case EVT_MARK_NODE_OFFLINE:
        return runtime_mark_node_offline(ctx->runtime, err);

    case EVT_RESTORE_NODE:
        return runtime_restore_nodes(ctx->runtime, err);

    case EVT_FAILURE_CONTROL:
        return event_failure_control(ctx, event, err);

    case EVT_SOFTWARE_UPDATE:
        return event_software_update(ctx, event, err);

    case EVT_DISCONNECT_NODES:
        rc = selector_resolve_nodes(ctx->world, &event->nodes, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_disconnect_nodes(ctx->runtime,
                                          ctx->world,
                                          &res,
                                          err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_RECONNECT_NODES:
        rc = selector_resolve_nodes(ctx->world, &event->nodes, &res, err);
        if (rc == ULAB_OK) {
            rc = runtime_reconnect_nodes(ctx->runtime,
                                         ctx->world,
                                         &res,
                                         err);
        }
        selector_result_free(&res);
        return rc;

    case EVT_RESTART_SITE:
        return event_restart_site(ctx, event, err);

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

    case EVT_WAIT_NODE_CONNECTIVITY:
        return event_wait_node_connectivity(ctx, event, err);

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
