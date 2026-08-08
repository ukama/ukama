/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <math.h>
#include <stdio.h>
#include <string.h>
#include <strings.h>
#include <time.h>
#include <unistd.h>

#include "check.h"
#include "selector.h"
#include "util.h"

#define ULAB_MAX_KPI_BUCKETS 64

static int bool_text(int value) {
    return value ? 1 : 0;
}

static int compare_number(double actual, const check_spec_t *check) {
    double tolerance;

    tolerance = check->tolerance_value > 0 ?
        check->tolerance_value : 0.000001;
    if (ulab_streq(check->comparator, "greater_than") ||
        ulab_streq(check->comparator, "gt")) {
        return actual > check->expected_value;
    }
    if (ulab_streq(check->comparator, "greater_or_equal") ||
        ulab_streq(check->comparator, "gte")) {
        return actual + tolerance >= check->expected_value;
    }
    if (ulab_streq(check->comparator, "less_than") ||
        ulab_streq(check->comparator, "lt")) {
        return actual < check->expected_value;
    }
    if (ulab_streq(check->comparator, "less_or_equal") ||
        ulab_streq(check->comparator, "lte")) {
        return actual - tolerance <= check->expected_value;
    }
    if (ulab_streq(check->comparator, "not_equals")) {
        return fabs(actual - check->expected_value) > tolerance;
    }
    return fabs(actual - check->expected_value) <= tolerance;
}

static int check_list_count(check_ctx_t *ctx,
                            const check_spec_t *check,
                            check_result_t *res,
                            ulab_error_t *err) {
    selector_result_t networks;
    size_t actual;
    size_t i;
    uint32_t expected;

    if (check->target[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "list_count_equals requires target");
        return ULAB_ERR;
    }
    if (check->has_expected_count) {
        expected = check->expected_count;
    } else if (ulab_parse_u32(check->expected, &expected) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "list_count_equals requires expected_count");
        return ULAB_ERR;
    }

    actual = 0;
    if (ulab_streq(check->target, "site_nodes")) {
        selector_result_t sites;

        if (selector_resolve_sites(ctx->world, &check->sites,
                                   &sites, err)) {
            return ULAB_ERR;
        }
        for (i = 0; i < sites.count; i++) {
            size_t site_count;

            site_count = 0;
            if (bff_get_site_node_count(
                    ctx->bff, &ctx->world->sites[sites.idx[i]],
                    check->nodes.kind == SEL_NODE_TYPE ||
                    check->nodes.kind ==
                        SEL_NODE_TYPE_COUNT_PER_NETWORK ?
                        check->nodes.value : NULL,
                    check->view, &site_count, err)) {
                selector_result_free(&sites);
                return ULAB_ERR;
            }
            actual += site_count;
        }
        selector_result_free(&sites);
    } else if (ulab_streq(check->target, "networks")) {
        if (bff_get_list_count(ctx->bff, check->target, NULL,
                               &actual, err)) {
            return ULAB_ERR;
        }
    } else {
        if (selector_resolve_networks(ctx->world, &check->networks,
                                      &networks, err)) {
            return ULAB_ERR;
        }
        if (networks.count != 1) {
            snprintf(err->msg, sizeof(err->msg),
                     "list_count_equals target %s requires exactly one "
                     "network, got %zu",
                     check->target, networks.count);
            selector_result_free(&networks);
            return ULAB_ERR;
        }
        i = networks.idx[0];
        if (ulab_streq(check->target, "nodes")) {
            const char *node_type;

            node_type = check->nodes.kind == SEL_NODE_TYPE ||
                check->nodes.kind == SEL_NODE_TYPE_COUNT_PER_NETWORK ?
                check->nodes.value : NULL;
            if (bff_get_node_list_count(
                    ctx->bff, &ctx->world->networks[i],
                    node_type, check->view, &actual, err)) {
                selector_result_free(&networks);
                return ULAB_ERR;
            }
        } else if (ulab_streq(check->target, "sites")) {
            if (bff_get_site_list_count(
                    ctx->bff, &ctx->world->networks[i],
                    check->view, &actual, err)) {
                selector_result_free(&networks);
                return ULAB_ERR;
            }
        } else if (bff_get_list_count(
                       ctx->bff, check->target,
                       &ctx->world->networks[i], &actual, err)) {
            selector_result_free(&networks);
            return ULAB_ERR;
        }
        selector_result_free(&networks);
    }

    res->passed = actual == (size_t)expected;
    snprintf(res->detail, sizeof(res->detail),
             "target=%s expected=%u actual=%zu",
             check->target, expected, actual);
    return ULAB_OK;
}

static int check_entity(check_ctx_t *ctx,
                        const check_spec_t *check,
                        check_result_t *res,
                        ulab_error_t *err) {
    int matched;

    if (check->entity[0] == '\0' || check->ref[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "%s requires entity and ref",
                 scenario_check_name(check->type));
        return ULAB_ERR;
    }

    matched = 0;
    if (check->type == CHECK_ENTITY_FIELDS_EQUAL) {
        if (bff_entity_fields_match_world(ctx->bff, check->entity,
                                          check->ref, ctx->world,
                                          check->view,
                                          &matched, res->detail,
                                          sizeof(res->detail), err)) {
            return ULAB_ERR;
        }
    } else {
        if (bff_entity_list_detail_reconciles(ctx->bff, check->entity,
                                              check->ref, ctx->world,
                                              check->view,
                                              &matched, res->detail,
                                              sizeof(res->detail), err)) {
            return ULAB_ERR;
        }
    }
    res->passed = matched;
    return ULAB_OK;
}

static int check_node_status_equals(check_ctx_t *ctx,
                                    const check_spec_t *check,
                                    check_result_t *res,
                                    ulab_error_t *err) {
    selector_result_t nodes;
    time_t deadline;
    unsigned int poll;
    size_t matched;
    size_t i;
    bff_node_status_t last;
    char last_node[ULAB_MAX_ID];

    if (check->connectivity[0] == '\0' &&
        check->lifecycle_state[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "node_status_equals requires connectivity and/or state");
        return ULAB_ERR;
    }
    if (selector_resolve_nodes(ctx->world, &check->nodes, &nodes, err)) {
        return ULAB_ERR;
    }
    if (nodes.count == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail), "selected_nodes=0");
        selector_result_free(&nodes);
        return ULAB_OK;
    }

    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 2u;
    memset(&last, 0, sizeof(last));
    last_node[0] = '\0';
    do {
        matched = 0;
        for (i = 0; i < nodes.count; i++) {
            node_t *node;
            network_t *network;
            bff_node_status_t status;
            int ok;

            node = &ctx->world->nodes[nodes.idx[i]];
            network = world_network_by_ref(ctx->world,
                                           node->network_ref);
            memset(&status, 0, sizeof(status));
            if (bff_get_node_status_for_view(ctx->bff, network, node,
                                             check->view, &status, err)) {
                selector_result_free(&nodes);
                return ULAB_ERR;
            }
            ok = 1;
            if (check->connectivity[0] != '\0') {
                ok = strcasecmp(status.connectivity,
                                check->connectivity) == 0;
            }
            if (ok && check->lifecycle_state[0] != '\0') {
                ok = strcasecmp(status.state,
                                check->lifecycle_state) == 0;
            }
            if (ok) {
                matched++;
            }
            last = status;
            ulab_copy(last_node, sizeof(last_node), node->bff_id);
        }
        if (matched == nodes.count || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched == nodes.count;
    snprintf(res->detail, sizeof(res->detail),
             "matched=%zu/%zu node=%.96s expected_connectivity=%.32s "
             "actual_connectivity=%.32s expected_state=%.32s "
             "actual_state=%.32s",
             matched, nodes.count, last_node, check->connectivity,
             last.connectivity, check->lifecycle_state, last.state);
    selector_result_free(&nodes);
    return ULAB_OK;
}

static int software_matches(const bff_software_t *software,
                            const check_spec_t *check) {
    if (check->status[0] != '\0') {
        if (ulab_streq(check->status, "terminal")) {
            if (strcasecmp(software->status, "up_to_date") != 0 &&
                strcasecmp(software->status, "update_failed") != 0) {
                return 0;
            }
        } else if (strcasecmp(software->status, check->status) != 0) {
            return 0;
        }
    }
    if (check->current_version[0] != '\0' &&
        !ulab_streq(software->current_version,
                    check->current_version)) {
        return 0;
    }
    if (check->desired_version[0] != '\0' &&
        !ulab_streq(software->desired_version,
                    check->desired_version)) {
        return 0;
    }
    return 1;
}

static int software_required_fields(const bff_software_t *software) {
    return software->name[0] != '\0' &&
        software->current_version[0] != '\0' &&
        software->release_date[0] != '\0' &&
        software->status[0] != '\0';
}

static int check_software_status(check_ctx_t *ctx,
                                 const check_spec_t *check,
                                 check_result_t *res,
                                 ulab_error_t *err) {
    selector_result_t nodes;
    time_t deadline;
    unsigned int poll;
    size_t matched;
    size_t i;
    bff_software_t last;
    int last_found;

    if (check->app[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "software_status_equals requires app");
        return ULAB_ERR;
    }
    if (check->status[0] == '\0' &&
        check->current_version[0] == '\0' &&
        check->desired_version[0] == '\0' && !check->required) {
        snprintf(err->msg, sizeof(err->msg),
                 "software_status_equals requires expected fields");
        return ULAB_ERR;
    }
    if (selector_resolve_nodes(ctx->world, &check->nodes, &nodes, err)) {
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 2u;
    memset(&last, 0, sizeof(last));
    last_found = 0;
    do {
        matched = 0;
        for (i = 0; i < nodes.count; i++) {
            node_t *node;
            int ok;

            node = &ctx->world->nodes[nodes.idx[i]];
            ok = 0;
            if (ulab_streq(check->app, "all")) {
                bff_software_t rows[ULAB_MAX_LIST];
                size_t count;
                size_t j;

                memset(rows, 0, sizeof(rows));
                count = 0;
                if (bff_get_software_list(ctx->bff, node, check->view, rows,
                                          ULAB_MAX_LIST, &count, err)) {
                    selector_result_free(&nodes);
                    return ULAB_ERR;
                }
                ok = count > 0;
                for (j = 0; ok && j < count; j++) {
                    ok = software_matches(&rows[j], check) &&
                        (!check->required ||
                         software_required_fields(&rows[j]));
                    last = rows[j];
                }
                last_found = count > 0;
            } else {
                bff_software_t software;
                int found;

                memset(&software, 0, sizeof(software));
                found = 0;
                if (bff_get_software(ctx->bff, node, check->app,
                                     check->view,
                                     &software, &found, err)) {
                    selector_result_free(&nodes);
                    return ULAB_ERR;
                }
                ok = found && software_matches(&software, check) &&
                    (!check->required ||
                     software_required_fields(&software));
                last = software;
                last_found = found;
            }
            if (ok) matched++;
        }
        if (matched == nodes.count || time(NULL) >= deadline) break;
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = nodes.count > 0 && matched == nodes.count;
    snprintf(res->detail, sizeof(res->detail),
             "app=%.96s matched=%zu/%zu found=%s status=%.32s "
             "current=%.64s release=%.48s required=%s",
             check->app, matched, nodes.count,
             last_found ? "true" : "false", last.status,
             last.current_version, last.release_date,
             check->required ? "true" : "false");
    selector_result_free(&nodes);
    return ULAB_OK;
}

static int check_software_count(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    selector_result_t nodes;
    time_t deadline;
    unsigned int poll;
    uint32_t expected;
    size_t matched;
    size_t i;
    size_t last_count;

    if (check->has_expected_count) {
        expected = check->expected_count;
    } else if (ulab_parse_u32(check->expected, &expected) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "software_count_equals requires expected_count");
        return ULAB_ERR;
    }
    if (selector_resolve_nodes(ctx->world, &check->nodes, &nodes, err)) {
        return ULAB_ERR;
    }
    if (nodes.count == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail), "selected_nodes=0");
        selector_result_free(&nodes);
        return ULAB_OK;
    }

    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 2u;
    matched = 0;
    last_count = 0;
    do {
        matched = 0;
        for (i = 0; i < nodes.count; i++) {
            size_t count;

            count = 0;
            if (bff_get_software_count(
                    ctx->bff, &ctx->world->nodes[nodes.idx[i]],
                    &count, err)) {
                selector_result_free(&nodes);
                return ULAB_ERR;
            }
            if (count == (size_t)expected) {
                matched++;
            }
            last_count = count;
        }
        if (matched == nodes.count || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched == nodes.count;
    snprintf(res->detail, sizeof(res->detail),
             "matched=%zu/%zu expected_count=%u last_count=%zu",
             matched, nodes.count, expected, last_count);
    selector_result_free(&nodes);
    return ULAB_OK;
}

static int node_operation_matches(
    const bff_node_operation_status_t *status,
    const check_spec_t *check) {
    if (check->has_expected_busy &&
        bool_text(status->busy) != bool_text(check->expected_busy)) {
        return 0;
    }
    if (check->operation_type[0] != '\0') {
        if (!status->has_operation ||
            strcasecmp(status->operation.type,
                       check->operation_type) != 0) {
            return 0;
        }
    }
    if (check->operation_status[0] != '\0') {
        if (!status->has_operation ||
            strcasecmp(status->operation.status,
                       check->operation_status) != 0) {
            return 0;
        }
    }
    return 1;
}

static int check_node_operation(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    selector_result_t nodes;
    time_t deadline;
    unsigned int poll;
    size_t matched;
    size_t i;
    bff_node_operation_status_t last;

    if (!check->has_expected_busy && check->operation_type[0] == '\0' &&
        check->operation_status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "node_operation_status_equals has no expected fields");
        return ULAB_ERR;
    }
    if (selector_resolve_nodes(ctx->world, &check->nodes, &nodes, err)) {
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 2u;
    memset(&last, 0, sizeof(last));
    do {
        matched = 0;
        for (i = 0; i < nodes.count; i++) {
            bff_node_operation_status_t status;

            memset(&status, 0, sizeof(status));
            if (bff_get_node_operation_status(
                    ctx->bff, &ctx->world->nodes[nodes.idx[i]],
                    &status, err)) {
                selector_result_free(&nodes);
                return ULAB_ERR;
            }
            if (node_operation_matches(&status, check)) {
                matched++;
            }
            last = status;
        }
        if (matched == nodes.count || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = nodes.count > 0 && matched == nodes.count;
    snprintf(res->detail, sizeof(res->detail),
             "matched=%zu/%zu busy=%s has_operation=%s type=%.48s "
             "status=%.32s",
             matched, nodes.count, last.busy ? "true" : "false",
             last.has_operation ? "true" : "false",
             last.operation.type, last.operation.status);
    selector_result_free(&nodes);
    return ULAB_OK;
}

static int contains_nocase(const char *text, const char *needle) {
    size_t len;

    if (needle == NULL || needle[0] == '\0') return 1;
    if (ulab_streq(needle, "present") ||
        ulab_streq(needle, "nonempty")) {
        return text != NULL && text[0] != '\0';
    }
    if (ulab_streq(needle, "absent") || ulab_streq(needle, "empty")) {
        return text == NULL || text[0] == '\0';
    }
    if (text == NULL) return 0;
    len = strlen(needle);
    while (*text != '\0') {
        if (strncasecmp(text, needle, len) == 0) return 1;
        text++;
    }
    return 0;
}

static int site_operation_matches(
    const bff_site_operation_status_t *status,
    const check_spec_t *check) {
    size_t i;
    int operation_matches;

    if (check->has_expected_busy &&
        bool_text(status->busy) != bool_text(check->expected_busy)) {
        return 0;
    }
    if (check->has_expected_degraded &&
        bool_text(status->degraded) !=
        bool_text(check->expected_degraded)) {
        return 0;
    }
    if (check->has_expected_restart_available &&
        bool_text(status->restart_site.available) !=
        bool_text(check->expected_restart_available)) {
        return 0;
    }
    if (check->has_expected_rf_available &&
        bool_text(status->rf.available) !=
        bool_text(check->expected_rf_available)) {
        return 0;
    }
    if (check->has_expected_service_available &&
        bool_text(status->service.available) !=
        bool_text(check->expected_service_available)) {
        return 0;
    }
    if (!contains_nocase(status->restart_site.reason,
                         check->restart_reason) ||
        !contains_nocase(status->rf.reason, check->rf_reason) ||
        !contains_nocase(status->service.reason,
                         check->service_reason)) {
        return 0;
    }

    if (check->operation_type[0] == '\0' &&
        check->operation_status[0] == '\0') {
        return 1;
    }
    operation_matches = 0;
    for (i = 0; i < status->node_count; i++) {
        const bff_node_operation_status_t *node;

        node = &status->nodes[i];
        if (!node->has_operation) {
            continue;
        }
        if (check->operation_type[0] != '\0' &&
            strcasecmp(node->operation.type,
                       check->operation_type) != 0) {
            continue;
        }
        if (check->operation_status[0] != '\0' &&
            strcasecmp(node->operation.status,
                       check->operation_status) != 0) {
            continue;
        }
        operation_matches = 1;
        break;
    }
    return operation_matches;
}

static int check_site_operation(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    selector_result_t sites;
    time_t deadline;
    unsigned int poll;
    size_t matched;
    size_t i;
    bff_site_operation_status_t last;

    if (!check->has_expected_busy && !check->has_expected_degraded &&
        !check->has_expected_restart_available &&
        !check->has_expected_rf_available &&
        !check->has_expected_service_available &&
        check->restart_reason[0] == '\0' &&
        check->rf_reason[0] == '\0' &&
        check->service_reason[0] == '\0' &&
        check->operation_type[0] == '\0' &&
        check->operation_status[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "site_operation_status_equals has no expected fields");
        return ULAB_ERR;
    }
    if (selector_resolve_sites(ctx->world, &check->sites, &sites, err)) {
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 2u;
    memset(&last, 0, sizeof(last));
    do {
        matched = 0;
        for (i = 0; i < sites.count; i++) {
            bff_site_operation_status_t status;

            memset(&status, 0, sizeof(status));
            if (bff_get_site_operation_status(
                    ctx->bff, &ctx->world->sites[sites.idx[i]],
                    &status, err)) {
                selector_result_free(&sites);
                return ULAB_ERR;
            }
            if (site_operation_matches(&status, check)) {
                matched++;
            }
            last = status;
        }
        if (matched == sites.count || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = sites.count > 0 && matched == sites.count;
    snprintf(res->detail, sizeof(res->detail),
             "matched=%zu/%zu busy=%s degraded=%s restart=%s rf=%s "
             "service=%s nodes=%zu reasons={restart=%.96s rf=%.96s "
             "service=%.96s}",
             matched, sites.count, last.busy ? "true" : "false",
             last.degraded ? "true" : "false",
             last.restart_site.available ? "true" : "false",
             last.rf.available ? "true" : "false",
             last.service.available ? "true" : "false",
             last.node_count, last.restart_site.reason, last.rf.reason,
             last.service.reason);
    selector_result_free(&sites);
    return ULAB_OK;
}

static network_t *one_network(check_ctx_t *ctx,
                              const check_spec_t *check,
                              ulab_error_t *err) {
    selector_result_t selected;
    network_t *network;

    network = NULL;
    if (check->networks.kind != SEL_NONE) {
        if (selector_resolve_networks(ctx->world, &check->networks,
                                      &selected, err)) {
            return NULL;
        }
        if (selected.count == 1) {
            network = &ctx->world->networks[selected.idx[0]];
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "%s requires exactly one selected network, got %zu",
                     scenario_check_name(check->type), selected.count);
        }
        selector_result_free(&selected);
        return network;
    }
    if (check->sites.kind != SEL_NONE) {
        if (selector_resolve_sites(ctx->world, &check->sites,
                                   &selected, err)) {
            return NULL;
        }
        if (selected.count == 1) {
            site_t *site;

            site = &ctx->world->sites[selected.idx[0]];
            network = world_network_by_ref(ctx->world, site->network_ref);
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "%s requires one site to infer the network, got %zu",
                     scenario_check_name(check->type), selected.count);
        }
        selector_result_free(&selected);
        return network;
    }
    if (check->nodes.kind != SEL_NONE) {
        if (selector_resolve_nodes(ctx->world, &check->nodes,
                                   &selected, err)) {
            return NULL;
        }
        if (selected.count == 1) {
            node_t *node;

            node = &ctx->world->nodes[selected.idx[0]];
            network = world_network_by_ref(ctx->world, node->network_ref);
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "%s requires one node to infer the network, got %zu",
                     scenario_check_name(check->type), selected.count);
        }
        selector_result_free(&selected);
        return network;
    }
    if (check->ues.kind != SEL_NONE) {
        if (selector_resolve_ues(ctx->world, &check->ues,
                                 &selected, err)) {
            return NULL;
        }
        if (selected.count == 1) {
            ue_t *ue;

            ue = &ctx->world->ues[selected.idx[0]];
            network = world_network_by_ref(ctx->world, ue->network_ref);
        } else {
            snprintf(err->msg, sizeof(err->msg),
                     "%s requires one UE to infer the network, got %zu",
                     scenario_check_name(check->type), selected.count);
        }
        selector_result_free(&selected);
        return network;
    }
    if (ctx->world->network_count == 1) {
        return &ctx->world->networks[0];
    }
    snprintf(err->msg, sizeof(err->msg),
             "%s requires a network selector in a multi-network world",
             scenario_check_name(check->type));
    return NULL;
}

static int resolve_scope_value(check_ctx_t *ctx,
                               const check_spec_t *check,
                               const network_t *network,
                               char *out,
                               size_t out_len,
                               ulab_error_t *err) {
    selector_result_t selected;

    out[0] = '\0';
    if (check->scope_value[0] != '\0') {
        network_t *net;
        site_t *site;
        node_t *node;
        ue_t *ue;
        subscriber_t *subscriber;
        package_t *package;

        net = world_network_by_ref(ctx->world, check->scope_value);
        site = world_site_by_ref(ctx->world, check->scope_value);
        node = world_node_by_ref(ctx->world, check->scope_value);
        ue = world_ue_by_ref(ctx->world, check->scope_value);
        subscriber = world_subscriber_by_ref(ctx->world,
                                             check->scope_value);
        package = network ?
            world_package_for_network(ctx->world, check->scope_value,
                                      network->ref) :
            world_package_by_base_ref(ctx->world, check->scope_value);
        if (net != NULL) {
            return ulab_copy(out, out_len, net->bff_id);
        }
        if (site != NULL) {
            return ulab_copy(out, out_len, site->bff_id);
        }
        if (node != NULL) {
            return ulab_copy(out, out_len, node->bff_id);
        }
        if (ue != NULL) {
            return ulab_copy(out, out_len, ue->bff_id);
        }
        if (subscriber != NULL) {
            return ulab_copy(out, out_len, subscriber->bff_id);
        }
        if (package != NULL) {
            return ulab_copy(out, out_len, package->bff_id);
        }
        return ulab_copy(out, out_len, check->scope_value);
    }
    if (ulab_streq(check->scope_key, "network_id") && network != NULL) {
        return ulab_copy(out, out_len, network->bff_id);
    }
    if (ulab_streq(check->scope_key, "package_id") &&
        check->package_ref[0] != '\0' && network != NULL) {
        package_t *package;

        package = world_package_for_network(ctx->world, check->package_ref,
                                            network->ref);
        if (package == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "cannot resolve package scope ref=%s",
                     check->package_ref);
            return ULAB_ERR;
        }
        return ulab_copy(out, out_len, package->bff_id);
    }
    if (ulab_streq(check->scope_key, "sim_id") &&
        check->ues.kind != SEL_NONE) {
        if (selector_resolve_ues(ctx->world, &check->ues,
                                 &selected, err)) {
            return ULAB_ERR;
        }
        if (selected.count != 1) {
            snprintf(err->msg, sizeof(err->msg),
                     "sim_id scope requires exactly one UE, got %zu",
                     selected.count);
            selector_result_free(&selected);
            return ULAB_ERR;
        }
        ulab_copy(out, out_len,
                  ctx->world->ues[selected.idx[0]].bff_id);
        selector_result_free(&selected);
        return ULAB_OK;
    }
    if (ulab_streq(check->scope_key, "site_id") &&
        check->sites.kind != SEL_NONE) {
        if (selector_resolve_sites(ctx->world, &check->sites,
                                   &selected, err)) {
            return ULAB_ERR;
        }
        if (selected.count != 1) {
            snprintf(err->msg, sizeof(err->msg),
                     "site_id scope requires exactly one site, got %zu",
                     selected.count);
            selector_result_free(&selected);
            return ULAB_ERR;
        }
        ulab_copy(out, out_len,
                  ctx->world->sites[selected.idx[0]].bff_id);
        selector_result_free(&selected);
    }
    return ULAB_OK;
}

static int resolve_timeseries_site(check_ctx_t *ctx,
                                   const check_spec_t *check,
                                   char *site_id,
                                   size_t site_id_len,
                                   ulab_error_t *err) {
    selector_result_t sites;

    site_id[0] = '\0';
    if (check->sites.kind != SEL_NONE) {
        if (selector_resolve_sites(ctx->world, &check->sites,
                                   &sites, err)) {
            return ULAB_ERR;
        }
        if (sites.count != 1) {
            snprintf(err->msg, sizeof(err->msg),
                     "kpi_timeseries supports exactly one site, got %zu",
                     sites.count);
            selector_result_free(&sites);
            return ULAB_ERR;
        }
        ulab_copy(site_id, site_id_len,
                  ctx->world->sites[sites.idx[0]].bff_id);
        selector_result_free(&sites);
        return ULAB_OK;
    }
    if (ulab_streq(check->scope_key, "site_id") &&
        check->scope_value[0] != '\0') {
        site_t *site;

        site = world_site_by_ref(ctx->world, check->scope_value);
        if (site != NULL) {
            return ulab_copy(site_id, site_id_len, site->bff_id);
        }
        return ulab_copy(site_id, site_id_len, check->scope_value);
    }
    return ULAB_OK;
}

static long long days_from_civil(int year, unsigned int month,
                                 unsigned int day) {
    int era;
    int adjusted_month;
    unsigned int yoe;
    unsigned int doy;
    unsigned int doe;

    year -= month <= 2;
    era = (year >= 0 ? year : year - 399) / 400;
    yoe = (unsigned int)(year - era * 400);
    adjusted_month = (int)month + (month > 2 ? -3 : 9);
    doy = (153u * (unsigned int)adjusted_month + 2u) / 5u + day - 1u;
    doe = yoe * 365u + yoe / 4u - yoe / 100u + doy;
    return (long long)era * 146097LL + (long long)doe - 719468LL;
}

static int parse_rfc3339_epoch(const char *value, time_t *out) {
    int year;
    unsigned int month;
    unsigned int day;
    unsigned int hour;
    unsigned int minute;
    unsigned int second;
    int consumed;
    const char *suffix;
    int offset_sign;
    unsigned int offset_hour;
    unsigned int offset_minute;
    long long epoch;
    long long offset;

    consumed = 0;
    if (value == NULL || out == NULL ||
        sscanf(value, "%d-%u-%uT%u:%u:%u%n", &year, &month, &day,
               &hour, &minute, &second, &consumed) != 6 ||
        month < 1 || month > 12 || day < 1 || day > 31 ||
        hour > 23 || minute > 59 || second > 60) {
        return ULAB_ERR;
    }
    suffix = value + consumed;
    if (*suffix == '.') {
        suffix++;
        while (*suffix >= '0' && *suffix <= '9') {
            suffix++;
        }
    }
    offset = 0;
    if (*suffix == 'Z' || *suffix == 'z' || *suffix == '\0') {
        if (*suffix != '\0') {
            suffix++;
        }
    } else if (*suffix == '+' || *suffix == '-') {
        offset_sign = *suffix == '+' ? 1 : -1;
        suffix++;
        if (sscanf(suffix, "%2u:%2u", &offset_hour,
                   &offset_minute) != 2 || offset_hour > 23 ||
            offset_minute > 59) {
            return ULAB_ERR;
        }
        offset = (long long)offset_sign *
            ((long long)offset_hour * 3600LL +
             (long long)offset_minute * 60LL);
        suffix += 5;
    } else {
        return ULAB_ERR;
    }
    if (*suffix != '\0') {
        return ULAB_ERR;
    }
    epoch = days_from_civil(year, month, day) * 86400LL +
        (long long)hour * 3600LL + (long long)minute * 60LL + second;
    *out = (time_t)(epoch - offset);
    return ULAB_OK;
}

static int valid_kpi_value_state(const char *state) {
    return ulab_streq(state, "missing") ||
        ulab_streq(state, "unavailable") ||
        ulab_streq(state, "no_data") ||
        ulab_streq(state, "present") ||
        ulab_streq(state, "available") ||
        ulab_streq(state, "zero") ||
        ulab_streq(state, "non_zero") ||
        ulab_streq(state, "fresh") ||
        ulab_streq(state, "stale");
}

static int valid_kpi_series_state(const char *state) {
    return ulab_streq(state, "empty") ||
        ulab_streq(state, "present") ||
        ulab_streq(state, "all_zero") ||
        ulab_streq(state, "non_zero");
}

static int kpi_state_matches(const bff_kpi_value_t *value,
                             int found,
                             const check_spec_t *check) {
    const char *state;
    double tolerance;

    state = check->value_state[0] ? check->value_state : check->expected;
    tolerance = check->tolerance_value > 0 ?
        check->tolerance_value : 0.000001;
    if (ulab_streq(state, "missing") ||
        ulab_streq(state, "unavailable")) {
        return !found;
    }
    if (ulab_streq(state, "no_data")) {
        return !found || value->computed_at[0] == '\0';
    }
    if (ulab_streq(state, "present")) {
        return found;
    }
    if (ulab_streq(state, "available")) {
        return found && value->computed_at[0] != '\0';
    }
    if (ulab_streq(state, "zero")) {
        return found && fabs(value->value) <= tolerance;
    }
    if (ulab_streq(state, "non_zero")) {
        return found && fabs(value->value) > tolerance;
    }
    if (ulab_streq(state, "fresh") || ulab_streq(state, "stale")) {
        time_t computed;
        time_t now;
        double age;
        uint32_t max_age;

        if (!found || value->computed_at[0] == '\0' ||
            parse_rfc3339_epoch(value->computed_at, &computed)) {
            return 0;
        }
        now = time(NULL);
        age = difftime(now, computed);
        max_age = check->max_age_seconds ? check->max_age_seconds : 120u;
        if (ulab_streq(state, "fresh")) {
            return age >= 0 && age <= (double)max_age;
        }
        return age > (double)max_age;
    }
    return 0;
}

static int check_kpi_state(check_ctx_t *ctx,
                           const check_spec_t *check,
                           check_result_t *res,
                           ulab_error_t *err) {
    network_t *network;
    bff_kpi_value_t value;
    char scope_value[ULAB_MAX_ID];
    const char *state;
    time_t deadline;
    unsigned int poll;
    int found;
    int matched;

    if (check->key[0] == '\0' ||
        (check->value_state[0] == '\0' && check->expected[0] == '\0')) {
        snprintf(err->msg, sizeof(err->msg),
                 "kpi_state_equals requires key and value_state");
        return ULAB_ERR;
    }
    network = one_network(ctx, check, err);
    if (network == NULL) {
        return ULAB_ERR;
    }
    if (resolve_scope_value(ctx, check, network, scope_value,
                            sizeof(scope_value), err)) {
        return ULAB_ERR;
    }
    state = check->value_state[0] ? check->value_state : check->expected;
    if (!valid_kpi_value_state(state)) {
        snprintf(err->msg, sizeof(err->msg),
                 "kpi_state_equals unsupported value_state %.32s", state);
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    found = 0;
    matched = 0;
    memset(&value, 0, sizeof(value));
    do {
        found = 0;
        memset(&value, 0, sizeof(value));
        if (bff_get_kpi_value(ctx->bff, check->key,
                              check->span[0] ? check->span : "daily",
                              check->op, network->bff_id,
                              check->scope_key,
                              scope_value[0] ? scope_value : NULL,
                              &value, &found, err)) {
            return ULAB_ERR;
        }
        matched = kpi_state_matches(&value, found, check);
        if (matched || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched;
    snprintf(res->detail, sizeof(res->detail),
             "kpi=%.48s expected_state=%.24s found=%s value=%.6g "
             "computed_at=%.40s",
             check->key, state, found ? "true" : "false", value.value,
             value.computed_at);
    return ULAB_OK;
}

static int series_metadata_valid(const bff_kpi_value_t values[],
                                 size_t count,
                                 const check_spec_t *check) {
    size_t i;
    time_t previous_from;
    time_t requested_from;
    time_t requested_to;
    int has_requested_from;
    int has_requested_to;

    previous_from = 0;
    has_requested_from = check->from[0] != '\0' &&
        parse_rfc3339_epoch(check->from, &requested_from) == ULAB_OK;
    has_requested_to = check->to[0] != '\0' &&
        parse_rfc3339_epoch(check->to, &requested_to) == ULAB_OK;
    for (i = 0; i < count; i++) {
        time_t from;
        time_t to;

        if (!ulab_streq(values[i].kpi, check->key) ||
            values[i].from[0] == '\0' || values[i].to[0] == '\0' ||
            parse_rfc3339_epoch(values[i].from, &from) != ULAB_OK ||
            parse_rfc3339_epoch(values[i].to, &to) != ULAB_OK ||
            from >= to || (i > 0 && previous_from >= from)) {
            return 0;
        }
        if (check->span[0] != '\0' && values[i].span[0] != '\0' &&
            !ulab_streq(values[i].span, check->span)) {
            return 0;
        }
        if (has_requested_from && from < requested_from) {
            return 0;
        }
        if (has_requested_to && to > requested_to) {
            return 0;
        }
        previous_from = from;
    }
    return 1;
}

static int check_kpi_timeseries(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    network_t *network;
    bff_kpi_value_t values[ULAB_MAX_KPI_BUCKETS];
    char site_id[ULAB_MAX_ID];
    time_t deadline;
    unsigned int poll;
    size_t count;
    size_t i;
    double sum;
    int matched;
    int metadata_ok;
    const char *state;

    if (check->key[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "kpi_timeseries requires key");
        return ULAB_ERR;
    }
    network = one_network(ctx, check, err);
    if (network == NULL) {
        return ULAB_ERR;
    }
    if (resolve_timeseries_site(ctx, check, site_id,
                                sizeof(site_id), err)) {
        return ULAB_ERR;
    }
    if (check->from[0] != '\0') {
        time_t parsed_from;

        if (parse_rfc3339_epoch(check->from, &parsed_from)) {
            snprintf(err->msg, sizeof(err->msg),
                     "kpi_timeseries invalid from timestamp %.64s",
                     check->from);
            return ULAB_ERR;
        }
    }
    if (check->to[0] != '\0') {
        time_t parsed_to;

        if (parse_rfc3339_epoch(check->to, &parsed_to)) {
            snprintf(err->msg, sizeof(err->msg),
                     "kpi_timeseries invalid to timestamp %.64s",
                     check->to);
            return ULAB_ERR;
        }
    }
    if (check->from[0] != '\0' && check->to[0] != '\0') {
        time_t parsed_from;
        time_t parsed_to;

        if (parse_rfc3339_epoch(check->from, &parsed_from) ||
            parse_rfc3339_epoch(check->to, &parsed_to) ||
            parsed_from >= parsed_to) {
            snprintf(err->msg, sizeof(err->msg),
                     "kpi_timeseries requires from before to");
            return ULAB_ERR;
        }
    }
    state = check->value_state[0] ? check->value_state : "present";
    if (!valid_kpi_series_state(state)) {
        snprintf(err->msg, sizeof(err->msg),
                 "kpi_timeseries unsupported value_state %.32s", state);
        return ULAB_ERR;
    }
    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    count = 0;
    sum = 0;
    matched = 0;
    metadata_ok = 0;
    do {
        memset(values, 0, sizeof(values));
        count = 0;
        if (bff_get_kpi_timeseries(
                ctx->bff, check->key,
                check->span[0] ? check->span : "daily",
                check->op, check->from, check->to,
                network->bff_id, site_id, values,
                ULAB_MAX_KPI_BUCKETS, &count, err)) {
            return ULAB_ERR;
        }
        sum = 0;
        metadata_ok = series_metadata_valid(values, count, check);
        for (i = 0; i < count; i++) {
            sum += values[i].value;
            if (check->require_computed_at &&
                values[i].computed_at[0] == '\0') {
                metadata_ok = 0;
            }
        }

        matched = metadata_ok;
        if (ulab_streq(state, "empty")) {
            matched = count == 0;
        } else if (ulab_streq(state, "present")) {
            matched = matched && count > 0;
        } else if (ulab_streq(state, "all_zero")) {
            matched = matched && count > 0 && fabs(sum) <=
                (check->tolerance_value > 0 ?
                 check->tolerance_value : 0.000001);
        } else if (ulab_streq(state, "non_zero")) {
            matched = matched && count > 0 && fabs(sum) >
                (check->tolerance_value > 0 ?
                 check->tolerance_value : 0.000001);
        }
        if (matched && check->has_expected_count) {
            matched = count == check->expected_count;
        }
        if (matched && check->has_expected_value) {
            matched = compare_number(sum, check);
        }
        if (matched || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched;
    snprintf(res->detail, sizeof(res->detail),
             "kpi=%.48s state=%.24s buckets=%zu expected_count=%u "
             "sum=%.6g expected=%.6g metadata=%s site=%.64s",
             check->key, state, count, check->expected_count, sum,
             check->expected_value, metadata_ok ? "true" : "false",
             site_id);
    return ULAB_OK;
}

int check_console(check_ctx_t *ctx, const check_spec_t *check,
                  check_result_t *res, ulab_error_t *err) {
    switch (check->type) {
    case CHECK_LIST_COUNT_EQUALS:
        return check_list_count(ctx, check, res, err);
    case CHECK_ENTITY_FIELDS_EQUAL:
    case CHECK_ENTITY_RECONCILES:
        return check_entity(ctx, check, res, err);
    case CHECK_NODE_STATUS_EQUALS:
        return check_node_status_equals(ctx, check, res, err);
    case CHECK_SOFTWARE_STATUS_EQUALS:
        return check_software_status(ctx, check, res, err);
    case CHECK_SOFTWARE_COUNT_EQUALS:
        return check_software_count(ctx, check, res, err);
    case CHECK_NODE_OPERATION_STATUS_EQUALS:
        return check_node_operation(ctx, check, res, err);
    case CHECK_SITE_OPERATION_STATUS_EQUALS:
        return check_site_operation(ctx, check, res, err);
    case CHECK_KPI_STATE_EQUALS:
        return check_kpi_state(ctx, check, res, err);
    case CHECK_KPI_TIMESERIES:
        return check_kpi_timeseries(ctx, check, res, err);
    default:
        snprintf(err->msg, sizeof(err->msg),
                 "unsupported reusable console check");
        return ULAB_ERR;
    }
}
