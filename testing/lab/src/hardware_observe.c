/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc. */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <time.h>
#include <unistd.h>
#include "hardware_observe.h"
#include "runtime_hardware.h"
#include "util.h"

static uint64_t now_ms(void) {
    struct timespec now;
    clock_gettime(CLOCK_MONOTONIC, &now);
    return (uint64_t)now.tv_sec * 1000u + (uint64_t)now.tv_nsec / 1000000u;
}

static unsigned int control_wait(void) {
    uint32_t seconds = 30;
    const char *value = getenv("ULAB_CONTROL_WAIT_SEC");
    if (value && (ulab_parse_u32(value, &seconds) || seconds > 300))
        return 30;
    return seconds;
}

static int nodes_ready(bff_client_t *bff, runtime_t *runtime, const world_t *world,
                       const selector_result_t *nodes, size_t *ready, ulab_error_t *err) {
    size_t i;
    *ready = 0;
    if (hardware_runtime_healthy(runtime, err))
        return ULAB_ERR;
    if (nodes->count == 0) {
        snprintf(err->msg, sizeof(err->msg), "hardware readiness selected no nodes");
        return ULAB_ERR;
    }
    for (i = 0; i < nodes->count; i++) {
        const node_t *node;
        bff_node_status_t status = {0};
        if (nodes->idx[i] >= world->node_count)
            return ULAB_ERR;
        node = &world->nodes[nodes->idx[i]];
        if (node->bff_id[0] == '\0' || !ulab_streq(node->id, node->bff_id)) {
            snprintf(err->msg, sizeof(err->msg), "hardware node identity is not bound");
            return ULAB_ERR;
        }
        if (bff_get_node_status(bff, node, &status, err))
            return ULAB_ERR;
        if (!ulab_streq(status.id, node->id)) {
            snprintf(err->msg, sizeof(err->msg), "backend returned a different physical node ID");
            return ULAB_ERR;
        }
        if (strcasecmp(status.connectivity, "online") == 0 &&
            strcasecmp(status.state, "operational") == 0) {
            (*ready)++;
        } else {
            snprintf(err->msg, sizeof(err->msg), "node %.128s connectivity=%.64s state=%.64s",
                     node->id, status.connectivity, status.state);
        }
    }
    return hardware_runtime_healthy(runtime, err);
}

static int wait_nodes(bff_client_t *bff, runtime_t *runtime, const world_t *world,
                       const selector_result_t *nodes, unsigned int timeout, int *matched,
                       ulab_error_t *err) {
    uint64_t deadline;
    size_t ready;
    int rc;
    *matched = 0;
    if (nodes->count == 0) {
        snprintf(err->msg, sizeof(err->msg), "hardware readiness selected no nodes");
        return ULAB_ERR;
    }
    if (timeout > 900) {
        snprintf(err->msg, sizeof(err->msg), "hardware readiness timeout exceeds 900 seconds");
        return ULAB_ERR;
    }
    deadline = now_ms() + (uint64_t)timeout * 1000u;
    for (;;) {
        rc = nodes_ready(bff, runtime, world, nodes, &ready, err);
        if (rc == ULAB_OK && ready == nodes->count) {
            *matched = 1;
            ulab_error_clear(err);
            return ULAB_OK;
        }
        if (hardware_runtime_healthy(runtime, err)) return ULAB_ERR;
        if (now_ms() >= deadline) return rc;
        sleep(1);
    }
}

int hardware_wait_nodes(bff_client_t *bff, runtime_t *runtime, const world_t *world,
                         const selector_result_t *nodes, unsigned int timeout,
                         ulab_error_t *err) {
    int matched;
    int rc = wait_nodes(bff, runtime, world, nodes, timeout, &matched, err);
    return rc == ULAB_OK && matched ? ULAB_OK : ULAB_ERR;
}

static int traffic_check(check_ctx_t *ctx, const check_spec_t *check, const selector_result_t *ues,
                         check_result_t *result, ulab_error_t *err) {
    hardware_traffic_result_t observed;
    unsigned int attempts = 0;
    uint64_t deadline = now_ms() + (uint64_t)control_wait() * 1000u;
    uint64_t amount = check->expected_used_mb ? check->expected_used_mb : 1;
    size_t matched;
    for (;;) {
        attempts++;
        if (hardware_runtime_traffic(ctx->runtime, ctx->world, ues, amount, &observed, err))
            return ULAB_ERR;
        matched = check->type == CHECK_TRAFFIC_ALLOWED ? observed.passed : observed.unavailable;
        result->passed = matched == observed.selected && observed.selected > 0;
        snprintf(result->detail, sizeof(result->detail),
                 "hardware ues=%zu passed=%zu unavailable=%zu amount_mb=%llu attempts=%u",
                 observed.selected, observed.passed, observed.unavailable,
                 (unsigned long long)amount, attempts);
        if (result->passed || now_ms() >= deadline)
            return ULAB_OK;
        sleep(1);
    }
}

int hardware_check_runtime(check_ctx_t *ctx, const check_spec_t *check, check_result_t *result,
                           ulab_error_t *err) {
    selector_result_t selected = {0};
    size_t attached;
    int rc;
    if (hardware_runtime_healthy(ctx->runtime, err))
        return ULAB_ERR;
    if (check->type == CHECK_TRAFFIC_BLOCKED) {
        snprintf(err->msg, sizeof(err->msg),
                 "hardware traffic_blocked requires real EPC policy evidence");
        return ULAB_ERR;
    }
    if (check->type == CHECK_NODE_READY) {
        if (selector_resolve_nodes(ctx->world, &check->nodes, &selected, err))
            return ULAB_ERR;
        rc = wait_nodes(ctx->bff, ctx->runtime, ctx->world, &selected,
                        check->timeout_seconds, &result->passed, err);
        snprintf(result->detail, sizeof(result->detail), "hardware nodes=%zu ready=%s %.768s",
                 selected.count, result->passed ? "true" : "false", err->msg);
        selector_result_free(&selected);
        return rc;
    }
    if (selector_resolve_ues(ctx->world, &check->ues, &selected, err))
        return ULAB_ERR;
    if (selected.count == 0) {
        snprintf(err->msg, sizeof(err->msg), "hardware check selected no UEs");
        selector_result_free(&selected);
        return ULAB_ERR;
    }
    if (check->type == CHECK_UE_ATTACHED) {
        rc = hardware_runtime_attached(ctx->runtime, ctx->world, &selected, &attached, err);
        result->passed = rc == ULAB_OK && attached == selected.count;
        snprintf(result->detail, sizeof(result->detail), "hardware attached=%zu/%zu", attached,
                 selected.count);
    } else {
        rc = traffic_check(ctx, check, &selected, result, err);
    }
    selector_result_free(&selected);
    if (rc == ULAB_OK)
        rc = hardware_runtime_healthy(ctx->runtime, err);
    return rc;
}
