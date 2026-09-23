/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <unistd.h>

#include "node_monitor.h"
#include "selector.h"
#include "util.h"
#include "log.h"

struct node_monitor {
    pthread_t thread;
    pthread_mutex_t lock;
    bff_client_t bff;
    node_t *nodes;
    size_t count;
    char connectivity[ULAB_MAX_REF];
    ulab_error_t failure;
    int stopping;
};

/* Keep the first failed observation, even if the node later recovers. */
static int monitor_sample(node_monitor_t *monitor) {
    size_t i;

    for (i = 0; i < monitor->count; i++) {
        bff_node_status_t status;
        ulab_error_t query_err;
        ulab_error_t failure;
        const node_t *node;

        node = &monitor->nodes[i];
        memset(&status, 0, sizeof(status));
        memset(&query_err, 0, sizeof(query_err));
        memset(&failure, 0, sizeof(failure));
        if (bff_get_node_status(&monitor->bff, node, &status, &query_err)) {
            snprintf(failure.msg, sizeof(failure.msg),
                     "node monitor node=%.128s query failed: %.700s",
                     node->bff_id, query_err.msg);
        } else if (strcasecmp(status.connectivity,
                              monitor->connectivity) != 0) {
            snprintf(failure.msg, sizeof(failure.msg),
                     "node monitor node=%.128s expected=%.32s "
                     "observed=%.32s state=%.32s",
                     node->bff_id, monitor->connectivity,
                     status.connectivity, status.state);
        }
        if (failure.msg[0] != '\0') {
            pthread_mutex_lock(&monitor->lock);
            if (monitor->failure.msg[0] == '\0') {
                monitor->failure = failure;
            }
            pthread_mutex_unlock(&monitor->lock);
            return ULAB_ERR;
        }
    }
    return ULAB_OK;
}

static void *monitor_run(void *arg) {
    node_monitor_t *monitor;
    int stopping;

    monitor = arg;
    for (;;) {
        sleep(1);
        pthread_mutex_lock(&monitor->lock);
        stopping = monitor->stopping;
        pthread_mutex_unlock(&monitor->lock);
        if (stopping || monitor_sample(monitor) != ULAB_OK) {
            return NULL;
        }
    }
}

int node_monitor_status(node_monitor_t *monitor, ulab_error_t *err) {
    int failed;

    if (monitor == NULL) {
        return ULAB_OK;
    }
    pthread_mutex_lock(&monitor->lock);
    failed = monitor->failure.msg[0] != '\0';
    if (failed) {
        *err = monitor->failure;
    }
    pthread_mutex_unlock(&monitor->lock);
    return failed ? ULAB_ERR : ULAB_OK;
}

int node_monitor_start(node_monitor_t **out, const bff_client_t *bff,
                       const world_t *world, const selector_t *nodes,
                       const char *connectivity, ulab_error_t *err) {
    selector_result_t selected;
    node_monitor_t *monitor;
    size_t i;
    int rc;

    if (out == NULL || *out != NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "node monitor context missing or already active");
        return ULAB_ERR;
    }
    if (!ulab_streq(connectivity, "Online") &&
        !ulab_streq(connectivity, "Offline")) {
        snprintf(err->msg, sizeof(err->msg),
                 "node monitor connectivity must be Online or Offline");
        return ULAB_ERR;
    }
    if (selector_resolve_nodes(world, nodes, &selected, err)) {
        selector_result_free(&selected);
        return ULAB_ERR;
    }
    if (selected.count == 0) {
        selector_result_free(&selected);
        snprintf(err->msg, sizeof(err->msg), "node monitor selected no nodes");
        return ULAB_ERR;
    }
    monitor = calloc(1, sizeof(*monitor));
    if (monitor == NULL) {
        selector_result_free(&selected);
        snprintf(err->msg, sizeof(err->msg), "node monitor allocation failed");
        return ULAB_ERR;
    }
    monitor->nodes = calloc(selected.count, sizeof(*monitor->nodes));
    if (monitor->nodes == NULL) {
        free(monitor);
        selector_result_free(&selected);
        snprintf(err->msg, sizeof(err->msg), "node monitor allocation failed");
        return ULAB_ERR;
    }
    monitor->count = selected.count;
    for (i = 0; i < selected.count; i++) {
        monitor->nodes[i] = world->nodes[selected.idx[i]];
    }
    selector_result_free(&selected);
    /* Private client and node snapshots: no shared token/log/world writes. */
    monitor->bff = *bff;
    monitor->bff.logf = NULL;
    ulab_copy(monitor->connectivity, sizeof(monitor->connectivity),
              connectivity);
    rc = pthread_mutex_init(&monitor->lock, NULL);
    if (rc != 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "node monitor mutex: %s", strerror(rc));
        free(monitor->nodes);
        free(monitor);
        return ULAB_ERR;
    }
    /* Establish the baseline synchronously, before the next event runs. */
    if (monitor_sample(monitor) != ULAB_OK) {
        node_monitor_status(monitor, err);
        goto failed;
    }
    rc = pthread_create(&monitor->thread, NULL, monitor_run, monitor);
    if (rc != 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "node monitor thread: %s", strerror(rc));
        goto failed;
    }
    *out = monitor;
    ulab_status("NODE", "monitor connectivity=%s nodes=%zu poll=1s",
                connectivity, monitor->count);
    return ULAB_OK;

failed:
    pthread_mutex_destroy(&monitor->lock);
    free(monitor->nodes);
    free(monitor);
    return ULAB_ERR;
}

int node_monitor_stop(node_monitor_t **out, ulab_error_t *err) {
    node_monitor_t *monitor;
    int rc;

    if (out == NULL || *out == NULL) {
        return ULAB_OK;
    }
    monitor = *out;
    pthread_mutex_lock(&monitor->lock);
    monitor->stopping = 1;
    pthread_mutex_unlock(&monitor->lock);
    pthread_join(monitor->thread, NULL);
    rc = node_monitor_status(monitor, err);
    if (rc == ULAB_OK) {
        monitor_sample(monitor);
        rc = node_monitor_status(monitor, err);
    }
    pthread_mutex_destroy(&monitor->lock);
    free(monitor->nodes);
    free(monitor);
    *out = NULL;
    return rc;
}
