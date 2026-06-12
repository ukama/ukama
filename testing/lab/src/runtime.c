/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "runtime.h"
#include "log.h"
#include "util.h"
#include "ulab.h"

static int run_script(runtime_t *rt,
                      const char *script,
                      const char *args,
                      ulab_error_t *err) {
    char cmd[8192];
    char out[4096];
    int n;

    n = snprintf(cmd, sizeof(cmd), "%s/%s %s",
                 rt->script_dir, script, args);
    if (n < 0 || (size_t)n >= sizeof(cmd)) {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime command too long: %s", script);
        return ULAB_ERR;
    }

    if (rt->logf) {
        fprintf(rt->logf, "$ %s\n", cmd);
        fflush(rt->logf);
    }

    memset(out, 0, sizeof(out));
    if (ulab_run_cmd(cmd, out, sizeof(out)) != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "script failed: %s", script);
        if (rt->logf) {
            fprintf(rt->logf, "FAILED: %s\n", out);
            fflush(rt->logf);
        }
        return ULAB_ERR;
    }

    if (rt->logf) {
        fprintf(rt->logf, "%s\n", out);
        fflush(rt->logf);
    }

    return ULAB_OK;
}

int runtime_init(runtime_t *rt,
                 const char *script_dir,
                 const char *run_dir) {

    char path[ULAB_MAX_PATH];

    memset(rt, 0, sizeof(*rt));
    ulab_copy(rt->script_dir, sizeof(rt->script_dir), script_dir);
    ulab_copy(rt->run_dir, sizeof(rt->run_dir), run_dir);

    snprintf(path, sizeof(path), "%s/runtime.log", run_dir);
    rt->logf = fopen(path, "w");

    return ULAB_OK;
}

void runtime_close(runtime_t *rt) {
    if (rt && rt->logf) {
        fclose(rt->logf);
        rt->logf = NULL;
    }
}

int runtime_build_and_start_nodes(const char *repo,
                                  runtime_t *rt,
                                  const world_t *w,
                                  const selector_result_t *nodes,
                                  ulab_error_t *err) {

    size_t i;
    const node_t *n;
    char args[ULAB_MAX_ARGS];
    int rc;

    for (i = 0; i < nodes->count; i++) {
        n = &w->nodes[nodes->idx[i]];

        memset(args, 0, sizeof(args));
        rc = snprintf(args,
                      sizeof(args),
                      "%s %s %s %s %s %s %llu",
                      repo,
                      n->id,
                      n->type,
                      n->site_ref,
                      n->network_ref,
                      rt->run_dir,
                      (unsigned long long)i);

        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "start-node args too long for node %s", n->id);
            return ULAB_ERR;
        }

        ulab_status("NODE", "build/start %s type=%s", n->id, n->type);

        if (run_script(rt, "start-node.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_wait_nodes_ready(runtime_t *rt,
                             const world_t *w,
                             const selector_result_t *nodes,
                             ulab_error_t *err) {
    size_t i;
    const node_t *n;
    char args[4096];
    int rc;

    for (i = 0; i < nodes->count; i++) {
        n = &w->nodes[nodes->idx[i]];

        memset(args, 0, sizeof(args));
        rc = snprintf(args, sizeof(args), "%s %s", n->id, rt->run_dir);
        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "wait-node args too long for node %s", n->id);
            return ULAB_ERR;
        }

        ulab_status("NODE", "wait ready %s", n->id);

        if (run_script(rt, "wait-nodes-ready.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_build_and_start_ues(const char *repo,
                                runtime_t *rt,
                                const world_t *w,
                                const selector_result_t *ues,
                                ulab_error_t *err) {
    size_t i;
    (void)repo;

    for (i = 0; i < ues->count; i++) {
        const ue_t *ue = &w->ues[ues->idx[i]];
        char args[4096];
        int rc;

        rc = snprintf(args, sizeof(args), "%s %s %s %s %s",
                      ue->id,
                      ue->imsi,
                      ue->iccid,
                      ue->site_ref,
                      rt->run_dir);

        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "start-ue args too long for ue %s", ue->id);
            return ULAB_ERR;
        }

        if (run_script(rt, "start-ue.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_wait_ues_attached(runtime_t *rt,
                              world_t *w,
                              const selector_result_t *ues,
                              ulab_error_t *err) {
    size_t i;

    for (i = 0; i < ues->count; i++) {
        ue_t *ue = &w->ues[ues->idx[i]];
        char args[4096];
        int rc;

        rc = snprintf(args, sizeof(args), "%s %s", ue->id, rt->run_dir);
        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "wait-ue args too long for ue %s", ue->id);
            return ULAB_ERR;
        }

        if (run_script(rt, "wait-ues-attached.sh", args, err)) {
            return ULAB_ERR;
        }

        ue->attached = 1;
    }

    return ULAB_OK;
}

int runtime_generate_traffic(runtime_t *rt,
                             const world_t *w,
                             const selector_result_t *ues,
                             uint64_t amount_mb,
                             ulab_error_t *err) {
    size_t i;

    for (i = 0; i < ues->count; i++) {
        const ue_t *ue = &w->ues[ues->idx[i]];
        char args[4096];
        int rc;

        rc = snprintf(args, sizeof(args), "%s %llu %s",
                      ue->id,
                      (unsigned long long)amount_mb,
                      rt->run_dir);

        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "traffic args too long for ue %s", ue->id);
            return ULAB_ERR;
        }

        if (run_script(rt, "traffic.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_restart_nodes(runtime_t *rt,
                          const world_t *w,
                          const selector_result_t *nodes,
                          ulab_error_t *err) {
    size_t i;

    for (i = 0; i < nodes->count; i++) {
        const node_t *n = &w->nodes[nodes->idx[i]];
        char args[4096];
        int rc;

        rc = snprintf(args, sizeof(args), "%s %s", n->id, rt->run_dir);
        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "restart-node args too long for node %s", n->id);
            return ULAB_ERR;
        }

        if (run_script(rt, "restart-node.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}
