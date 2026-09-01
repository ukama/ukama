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
    bff_component_t *components;
    size_t count;
    size_t wrong_category;
    size_t i;

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

    components = calloc(ULAB_MAX_BFF_COMPONENTS, sizeof(*components));
    if (components == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory reading component inventory");
        return ULAB_ERR;
    }

    if (bff_get_components_by_category(ctx->bff, check->category,
                                       components,
                                       ULAB_MAX_BFF_COMPONENTS,
                                       &count, err)) {
        free(components);
        return ULAB_ERR;
    }

    wrong_category = 0;
    for (i = 0; i < count; i++) {
        if (strcasecmp(components[i].category, check->category) != 0) {
            wrong_category++;
        }
    }

    res->passed = count == (size_t)check->expected_count &&
        wrong_category == 0;
    snprintf(res->detail, sizeof(res->detail),
             "category=%.32s expected=%u actual=%zu wrong_category=%zu",
             check->category, check->expected_count, count, wrong_category);
    free(components);

    return ULAB_OK;
}

static int component_for_part(const bff_component_t *components,
                              size_t count,
                              const char *part_number,
                              const bff_component_t **out) {
    size_t i;

    for (i = 0; i < count; i++) {
        if (ulab_streq(components[i].part_number, part_number)) {
            if (out != NULL) {
                *out = &components[i];
            }
            return 1;
        }
    }

    return 0;
}

/*
 * The inventory component service polls the node factory and turns every
 * provisioned node into an access component keyed by the factory node id.
 * The lab reserves its node bundle from the same factory, so each node it
 * started must surface as a component within a few sync cycles.
 */
static int check_node_component(check_ctx_t *ctx,
                                const check_spec_t *check,
                                check_result_t *res,
                                ulab_error_t *err) {
    selector_result_t nodes;
    bff_component_t *components;
    const char *category;
    time_t deadline;
    unsigned int poll;
    size_t count;
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

    for (i = 0; i < nodes.count; i++) {
        if (ctx->world->nodes[nodes.idx[i]].bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "node %s has no factory id; the scenario must start "
                     "nodes before checking the component inventory",
                     ctx->world->nodes[nodes.idx[i]].ref);
            selector_result_free(&nodes);
            return ULAB_ERR;
        }
    }

    components = calloc(ULAB_MAX_BFF_COMPONENTS, sizeof(*components));
    if (components == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory reading component inventory");
        selector_result_free(&nodes);
        return ULAB_ERR;
    }

    deadline = time(NULL) + (time_t)check->timeout_seconds;
    poll = check->poll_seconds ? check->poll_seconds : 5u;
    count = 0;
    matched = 0;
    wrong_type = 0;
    missing[0] = '\0';

    do {
        matched = 0;
        wrong_type = 0;
        missing[0] = '\0';

        if (bff_get_components_by_category(ctx->bff, category, components,
                                           ULAB_MAX_BFF_COMPONENTS,
                                           &count, err)) {
            free(components);
            selector_result_free(&nodes);
            return ULAB_ERR;
        }

        for (i = 0; i < nodes.count; i++) {
            const node_t *node;
            const bff_component_t *component;

            node = &ctx->world->nodes[nodes.idx[i]];
            component = NULL;
            if (!component_for_part(components, count, node->bff_id,
                                    &component)) {
                if (missing[0] == '\0') {
                    ulab_copy(missing, sizeof(missing), node->bff_id);
                }
                continue;
            }
            if (!ulab_streq(component->type, node_kind(node))) {
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
             "category=%.32s matched=%zu/%zu components=%zu wrong_type=%zu "
             "first_missing_part=%.96s",
             category, matched, nodes.count, count, wrong_type,
             missing[0] != '\0' ? missing : "-");
    free(components);
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
