/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <time.h>
#include <unistd.h>

#include "check.h"
#include "selector.h"
#include "util.h"

/* Upper bound on the SIM pool page the lab is willing to scan in one call. */
#define ULAB_MAX_POOL_SCAN 4096

static const char *node_kind(const node_t *node) {
    if (node == NULL) {
        return "";
    }
    if (ulab_streq(node->type, ULAB_NODE_TOWER)) {
        return ULAB_NODE_KIND_TOWER;
    }
    if (ulab_streq(node->type, ULAB_NODE_AMPLIFIER)) {
        return ULAB_NODE_KIND_AMPLIFIER;
    }
    if (ulab_streq(node->type, ULAB_NODE_CONTROLLER)) {
        return ULAB_NODE_KIND_CONTROLLER;
    }
    return node->type;
}

/*
 * Components in the power, backhaul, switch, and spectrum categories are
 * seeded from the networks repository, so their counts are fixed by the
 * fixture rather than by anything the scenario provisions.
 */
static int check_component_count(check_ctx_t *ctx,
                                 const check_spec_t *check,
                                 check_result_t *res,
                                 ulab_error_t *err) {
    uint32_t total;
    uint32_t wrong_category;

    if (check->category[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "component_count_by_category requires category");
        return ULAB_ERR;
    }
    if (!check->has_expected_count) {
        snprintf(err->msg, sizeof(err->msg),
                 "component_count_by_category requires expected_count");
        return ULAB_ERR;
    }

    total = 0;
    wrong_category = 0;
    if (bff_probe_components_by_category(ctx->bff, check->category, NULL, 0,
                                         &total, &wrong_category, err)) {
        return ULAB_ERR;
    }

    res->passed = total == check->expected_count && wrong_category == 0;
    snprintf(res->detail, sizeof(res->detail),
             "category=%.32s expected=%u actual=%u wrong_category=%u",
             check->category, check->expected_count, total, wrong_category);

    return ULAB_OK;
}

/*
 * The inventory component service polls the node factory and turns every
 * provisioned node allocated to the org into an access component keyed by the
 * node id. Every node this run put through the factory must surface as a
 * component within a few sync cycles.
 */
static int check_node_component(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    selector_result_t nodes;
    bff_component_probe_t *probes;
    const char *category;
    time_t deadline;
    unsigned int poll;
    uint32_t total;
    size_t matched;
    size_t wrong_type;
    size_t i;
    char missing[ULAB_MAX_ID];

    category = check->category[0] != '\0' ? check->category : "access";

    if (selector_resolve_nodes(ctx->world, &check->nodes, &nodes, err)) {
        return ULAB_ERR;
    }
    if (nodes.count == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail), "selected_nodes=0");
        selector_result_free(&nodes);
        return ULAB_OK;
    }

    probes = calloc(nodes.count, sizeof(*probes));
    if (probes == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory probing component inventory");
        selector_result_free(&nodes);
        return ULAB_ERR;
    }

    for (i = 0; i < nodes.count; i++) {
        const node_t *node;

        node = &ctx->world->nodes[nodes.idx[i]];
        if (node->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "node %s has no factory id; the scenario must start "
                     "nodes before checking the component inventory",
                     node->ref);
            free(probes);
            selector_result_free(&nodes);
            return ULAB_ERR;
        }
        ulab_copy(probes[i].part_number, sizeof(probes[i].part_number),
                  node->bff_id);
        ulab_copy(probes[i].type, sizeof(probes[i].type), node_kind(node));
    }

    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    total = 0;
    matched = 0;
    wrong_type = 0;
    missing[0] = '\0';

    do {
        matched = 0;
        wrong_type = 0;
        missing[0] = '\0';

        if (bff_probe_components_by_category(ctx->bff, category, probes,
                                             nodes.count, &total, NULL,
                                             err)) {
            free(probes);
            selector_result_free(&nodes);
            return ULAB_ERR;
        }

        for (i = 0; i < nodes.count; i++) {
            if (!probes[i].found) {
                if (missing[0] == '\0') {
                    ulab_copy(missing, sizeof(missing),
                              probes[i].part_number);
                }
                continue;
            }
            if (!probes[i].type_matches) {
                wrong_type++;
                continue;
            }
            matched++;
        }

        if (matched == nodes.count || time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched == nodes.count;
    snprintf(res->detail, sizeof(res->detail),
             "category=%.32s matched=%zu/%zu components=%u wrong_type=%zu "
             "first_missing_part=%.96s",
             category, matched, nodes.count, total, wrong_type,
             missing[0] != '\0' ? missing : "-");
    free(probes);
    selector_result_free(&nodes);

    return ULAB_OK;
}

static int pool_contains(char (*iccids)[ULAB_MAX_ID],
                         size_t count,
                         const char *iccid) {
    size_t i;

    for (i = 0; i < count; i++) {
        if (ulab_streq(iccids[i], iccid)) {
            return 1;
        }
    }

    return 0;
}

/*
 * Every SIM the lab pushes into the warehouse is exported by the SIM factory
 * and uploaded to the pool during setup. This check reads the pool back and
 * requires each of those ICCIDs to be there and unassigned, then confirms the
 * pool statistics account for them.
 */
static int check_sim_pool_contains(check_ctx_t *ctx,
                                   const check_spec_t *check,
                                   check_result_t *res,
                                   ulab_error_t *err) {
    selector_result_t ues;
    bff_inventory_summary_t stats;
    char (*iccids)[ULAB_MAX_ID];
    time_t deadline;
    unsigned int poll;
    size_t pool_count;
    size_t expected;
    size_t matched;
    size_t i;
    int stats_consistent;
    char missing[ULAB_MAX_ID];

    if (ctx->sim_type == NULL || ctx->sim_type[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "sim_pool_contains_sims requires a SIM type");
        return ULAB_ERR;
    }
    if (selector_resolve_ues(ctx->world, &check->ues, &ues, err)) {
        return ULAB_ERR;
    }

    expected = 0;
    for (i = 0; i < ues.count; i++) {
        if (ctx->world->ues[ues.idx[i]].iccid[0] != '\0') {
            expected++;
        }
    }
    if (expected == 0) {
        res->passed = 0;
        snprintf(res->detail, sizeof(res->detail),
                 "selected_sims=%zu with_iccid=0", ues.count);
        selector_result_free(&ues);
        return ULAB_OK;
    }

    iccids = calloc(ULAB_MAX_POOL_SCAN, sizeof(*iccids));
    if (iccids == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory reading SIM pool");
        selector_result_free(&ues);
        return ULAB_ERR;
    }

    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    pool_count = 0;
    matched = 0;
    stats_consistent = 0;
    memset(&stats, 0, sizeof(stats));
    missing[0] = '\0';

    do {
        matched = 0;
        missing[0] = '\0';

        if (bff_get_sims_from_pool(ctx->bff, ctx->sim_type, iccids, NULL,
                                   ULAB_MAX_POOL_SCAN, &pool_count, err) ||
            bff_get_sim_pool_summary(ctx->bff, ctx->sim_type, &stats, err)) {
            free(iccids);
            selector_result_free(&ues);
            return ULAB_ERR;
        }

        for (i = 0; i < ues.count; i++) {
            const ue_t *ue;

            ue = &ctx->world->ues[ues.idx[i]];
            if (ue->iccid[0] == '\0') {
                continue;
            }
            if (pool_contains(iccids, pool_count, ue->iccid)) {
                matched++;
                continue;
            }
            if (missing[0] == '\0') {
                ulab_copy(missing, sizeof(missing), ue->iccid);
            }
        }

        stats_consistent = stats.sim_total == stats.sim_available +
            stats.sim_consumed + stats.sim_failed;

        if ((matched == expected && stats_consistent) ||
            time(NULL) >= deadline) {
            break;
        }
        sleep(poll > 60u ? 60u : poll);
    } while (1);

    res->passed = matched == expected && stats_consistent &&
        stats.sim_available >= (uint32_t)expected;
    snprintf(res->detail, sizeof(res->detail),
             "type=%.32s in_pool=%zu/%zu unassigned=%zu total=%u "
             "available=%u consumed=%u failed=%u first_missing_iccid=%.32s",
             ctx->sim_type, matched, expected, pool_count, stats.sim_total,
             stats.sim_available, stats.sim_consumed, stats.sim_failed,
             missing[0] != '\0' ? missing : "-");
    free(iccids);
    selector_result_free(&ues);

    return ULAB_OK;
}

int check_inventory(check_ctx_t *ctx, const check_spec_t *check,
                    check_result_t *res, ulab_error_t *err) {
    switch (check->type) {
    case CHECK_COMPONENT_COUNT_BY_CATEGORY:
        return check_component_count(ctx, check, res, err);
    case CHECK_NODE_COMPONENT_REGISTERED:
        return check_node_component(ctx, check, res, err);
    case CHECK_SIM_POOL_CONTAINS_SIMS:
        return check_sim_pool_contains(ctx, check, res, err);
    default:
        snprintf(err->msg, sizeof(err->msg), "not an inventory check");
        return ULAB_ERR;
    }
}
