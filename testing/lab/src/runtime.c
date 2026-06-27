/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <ctype.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "runtime.h"
#include "log.h"
#include "util.h"
#include "ulab.h"

static int run_script(runtime_t *rt,
                      const char *name,
                      const char *args,
                      ulab_error_t *err) {

    char cmd[ULAB_MAX_QUERY];
    char script[ULAB_MAX_PATH];
    char log_path[ULAB_MAX_PATH];
    int rc;
    int n;

    n = snprintf(script, sizeof(script), "%s/%s", rt->script_dir, name);
    if (n < 0 || (size_t)n >= sizeof(script)) {
        snprintf(err->msg, sizeof(err->msg), "script path too long");
        return ULAB_ERR;
    }

    n = snprintf(log_path, sizeof(log_path), "%s/%s.log",
                 rt->run_dir, name);
    if (n < 0 || (size_t)n >= sizeof(log_path)) {
        snprintf(err->msg, sizeof(err->msg), "script log path too long");
        return ULAB_ERR;
    }

    /*
     * Keep script output in its run log. Do not print scrolling RUNNING lines.
     * When stderr is interactive, show a tiny single-line spinner instead.
     */
    n = snprintf(cmd, sizeof(cmd),
                 "mkdir -p '%s' >/dev/null 2>&1; "
                 ": > '%s'; "
                 "('%s' %s >> '%s' 2>&1) & "
                 "pid=$!; "
                 "start=$(date +%%s); "
                 "interval=${ULAB_SCRIPT_PROGRESS_INTERVAL:-1}; "
                 "spin='|/-\\'; "
                 "idx=0; "
                 "printed=0; "
                 "while kill -0 $pid 2>/dev/null; do "
                 "sleep $interval; "
                 "if kill -0 $pid 2>/dev/null && [ -t 2 ]; then "
                 "now=$(date +%%s); "
                 "elapsed=$((now - start)); "
                 "idx=$((idx + 1)); "
                 "case $((idx %% 4)) in "
                 "0) ch='|';; 1) ch='/';; 2) ch='-';; *) ch='\\';; "
                 "esac; "
                 "printf '\\r%%s %s %%ss' \"$ch\" \"$elapsed\" >&2; "
                 "printed=1; "
                 "fi; "
                 "done; "
                 "if [ \"$printed\" = 1 ]; then printf '\\r\\033[K' >&2; fi; "
                 "wait $pid",
                 rt->run_dir,
                 log_path,
                 script,
                 args ? args : "",
                 log_path,
                 name);
    if (n < 0 || (size_t)n >= sizeof(cmd)) {
        snprintf(err->msg, sizeof(err->msg), "script command too long");
        return ULAB_ERR;
    }

    if (rt->logf) {
        fprintf(rt->logf, "--- script %s ---\n", name);
        fprintf(rt->logf, "%s\n", cmd);
        fprintf(rt->logf, "log=%s\n", log_path);
        fflush(rt->logf);
    }

    rc = system(cmd);
    if (rc != 0) {
        snprintf(err->msg, sizeof(err->msg),
                 "script failed: %s; see runtime log", name);

        if (rt->logf) {
            fprintf(rt->logf, "script failed: %s\n", name);
            fprintf(rt->logf, "script log: %s\n", log_path);
            fflush(rt->logf);
        }

        return ULAB_ERR;
    }

    return ULAB_OK;
}

static void safe_name(const char *in, char *out, size_t out_len) {
    size_t i;
    size_t j;
    unsigned char ch;

    if (out_len == 0) {
        return;
    }

    j = 0;
    for (i = 0; in != NULL && in[i] != '\0' && j + 1 < out_len; i++) {
        ch = (unsigned char)in[i];
        if (isalnum(ch) || ch == '_' || ch == '.' || ch == '-') {
            out[j++] = (char)ch;
        } else {
            out[j++] = '-';
        }
    }

    out[j] = '\0';
}

static int read_state_value(const char *path,
                            const char *key,
                            char *out,
                            size_t out_len) {
    FILE *f;
    char line[ULAB_MAX_LINE];
    size_t key_len;
    char *v;

    f = fopen(path, "r");
    if (f == NULL) {
        return ULAB_ERR;
    }

    key_len = strlen(key);
    while (fgets(line, sizeof(line), f) != NULL) {
        if (strncmp(line, key, key_len) != 0 || line[key_len] != '=') {
            continue;
        }

        v = ulab_trim(line + key_len + 1);
        if (ulab_copy(out, out_len, v)) {
            fclose(f);
            return ULAB_ERR;
        }

        fclose(f);
        return ULAB_OK;
    }

    fclose(f);
    return ULAB_ERR;
}


static int runtime_site_state_path(runtime_t *rt,
                                   const site_t *site,
                                   char *path,
                                   size_t path_len,
                                   ulab_error_t *err) {
    char safe[ULAB_MAX_REF];
    int rc;

    safe_name(site->ref, safe, sizeof(safe));
    rc = snprintf(path, path_len, "%s/runtime-sites/%s.env",
                  rt->run_dir, safe);
    if (rc < 0 || (size_t)rc >= path_len) {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime site state path too long");
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static int write_runtime_node_state(runtime_t *rt,
                                    const site_t *site,
                                    node_t *node,
                                    const char *factory_id,
                                    const char *node_kind,
                                    const char *container,
                                    ulab_error_t *err) {
    char safe[ULAB_MAX_REF];
    char path[ULAB_MAX_PATH];
    FILE *f;
    int rc;

    safe_name(node->id, safe, sizeof(safe));
    rc = snprintf(path, sizeof(path), "%s/runtime-nodes/%s.env",
                  rt->run_dir, safe);
    if (rc < 0 || (size_t)rc >= sizeof(path)) {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime node state path too long for %s", node->id);
        return ULAB_ERR;
    }

    f = fopen(path, "w");
    if (f == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "failed to write runtime node state for %s", node->id);
        return ULAB_ERR;
    }

    fprintf(f, "LOGICAL_NODE_ID=%s\n", node->id);
    fprintf(f, "FACTORY_NODE_ID=%s\n", factory_id);
    fprintf(f, "NODE_TYPE=%s\n", node_kind);
    fprintf(f, "NODE_KIND=%s\n", node_kind);
    fprintf(f, "SITE_REF=%s\n", site->ref);
    fprintf(f, "NETWORK_REF=%s\n", site->network_ref);
    fprintf(f, "CONTAINER_NAME=%s\n", container);
    fprintf(f, "IMAGE=testing/virtualnode:%s\n", factory_id);
    fclose(f);

    ulab_copy(node->bff_id, sizeof(node->bff_id), factory_id);

    return ULAB_OK;
}

static int map_runtime_site_nodes(runtime_t *rt,
                                  world_t *w,
                                  site_t *site,
                                  ulab_error_t *err) {
    char site_state[ULAB_MAX_PATH];
    char tcontainer[ULAB_MAX_ID];
    char ccontainer[ULAB_MAX_ID];
    char acontainer[ULAB_MAX_ID];
    size_t i;

    if (runtime_site_state_path(rt, site, site_state,
        sizeof(site_state), err)) {
        return ULAB_ERR;
    }

    if (read_state_value(site_state, "TNODE_CONTAINER", tcontainer,
        sizeof(tcontainer))) {
        snprintf(err->msg, sizeof(err->msg),
                 "TNODE_CONTAINER missing for site %s", site->ref);
        return ULAB_ERR;
    }

    if (read_state_value(site_state, "CNODE_CONTAINER", ccontainer,
        sizeof(ccontainer))) {
        snprintf(err->msg, sizeof(err->msg),
                 "CNODE_CONTAINER missing for site %s", site->ref);
        return ULAB_ERR;
    }

    if (read_state_value(site_state, "ANODE_CONTAINER", acontainer,
        sizeof(acontainer))) {
        snprintf(err->msg, sizeof(err->msg),
                 "ANODE_CONTAINER missing for site %s", site->ref);
        return ULAB_ERR;
    }

    for (i = 0; i < w->node_count; i++) {
        node_t *node = &w->nodes[i];

        if (!ulab_streq(node->site_ref, site->ref)) {
            continue;
        }

        if (ulab_streq(node->type, ULAB_NODE_TOWER)) {
            if (write_runtime_node_state(rt, site, node, site->tnode_id,
                ULAB_NODE_KIND_TOWER, tcontainer, err)) {
                return ULAB_ERR;
            }
        } else if (ulab_streq(node->type, ULAB_NODE_AMPLIFIER)) {
            if (write_runtime_node_state(rt, site, node, site->anode_id,
                ULAB_NODE_KIND_AMPLIFIER, acontainer, err)) {
                return ULAB_ERR;
            }
        } else if (ulab_streq(node->type, ULAB_NODE_CONTROLLER)) {
            if (write_runtime_node_state(rt, site, node, site->cnode_id,
                ULAB_NODE_KIND_CONTROLLER, ccontainer, err)) {
                return ULAB_ERR;
            }
        }
    }

    return ULAB_OK;
}

static int load_runtime_site_state(runtime_t *rt,
                                   site_t *site,
                                   ulab_error_t *err) {
    char path[ULAB_MAX_PATH];

    if (runtime_site_state_path(rt, site, path, sizeof(path), err)) {
        return ULAB_ERR;
    }

    if (read_state_value(path, "TNODE_ID", site->tnode_id,
        sizeof(site->tnode_id))) {
        snprintf(err->msg, sizeof(err->msg),
                 "TNODE_ID missing for site %s", site->ref);
        return ULAB_ERR;
    }

    if (read_state_value(path, "CNODE_ID", site->cnode_id,
        sizeof(site->cnode_id))) {
        snprintf(err->msg, sizeof(err->msg),
                 "CNODE_ID missing for site %s", site->ref);
        return ULAB_ERR;
    }

    if (read_state_value(path, "ANODE_ID", site->anode_id,
        sizeof(site->anode_id))) {
        snprintf(err->msg, sizeof(err->msg),
                 "ANODE_ID missing for site %s", site->ref);
        return ULAB_ERR;
    }

    if (rt->logf) {
        fprintf(rt->logf,
                "runtime-site site=%s tnode=%s cnode=%s anode=%s state=%s\n",
                site->ref, site->tnode_id, site->cnode_id,
                site->anode_id, path);
        fflush(rt->logf);
    }

    return ULAB_OK;
}

int runtime_init(runtime_t *rt,
                 const char *script_dir,
                 const char *run_dir,
                 const char *repo) {

    char path[ULAB_MAX_PATH];

    memset(rt, 0, sizeof(*rt));
    ulab_copy(rt->script_dir, sizeof(rt->script_dir), script_dir);
    ulab_copy(rt->run_dir, sizeof(rt->run_dir), run_dir);
    ulab_copy(rt->repo, sizeof(rt->repo), repo ? repo : "");

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


int runtime_ensure_network(runtime_t *rt, ulab_error_t *err) {

    char args[ULAB_MAX_ARGS];
    int rc;

    rc = snprintf(args, sizeof(args), "%s", rt->run_dir);
    if (rc < 0 || (size_t)rc >= sizeof(args)) {
        snprintf(err->msg, sizeof(err->msg),
                 "ensure-network args too long");
        return ULAB_ERR;
    }

    ulab_status("NET", "ensure lab podman network");
    if (run_script(rt, "ensure-network.sh", args, err)) {
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runtime_build_and_start_sites(const char *repo,
                                  runtime_t *rt,
                                  world_t *w,
                                  ulab_error_t *err) {
    size_t i;
    site_t *site;
    char args[ULAB_MAX_ARGS];
    int rc;

    for (i = 0; i < w->site_count; i++) {
        site = &w->sites[i];

        memset(args, 0, sizeof(args));
        rc = snprintf(args, sizeof(args), "%s %s %s %s %llu",
                      repo,
                      site->ref,
                      site->network_ref,
                      rt->run_dir,
                      (unsigned long long)i);
        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "start-site args too long for site %s", site->ref);
            return ULAB_ERR;
        }

        ulab_status("SITE", "factory/build/start %s", site->ref);

        if (run_script(rt, "build-and-start-site.sh", args, err)) {
            return ULAB_ERR;
        }

        if (load_runtime_site_state(rt, site, err)) {
            return ULAB_ERR;
        }

        if (map_runtime_site_nodes(rt, w, site, err)) {
            return ULAB_ERR;
        }

        ulab_status("SITE", "%s tnode=%s cnode=%s anode=%s",
                    site->ref, site->tnode_id, site->cnode_id,
                    site->anode_id);
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

int runtime_enable_pcrf_service(runtime_t *rt, const world_t *w,
                                ulab_error_t *err) {
    size_t i;
    int rc;
    char args[ULAB_MAX_ARGS];

    if (rt == NULL || w == NULL) {
        return ULAB_OK;
    }

    for (i = 0; i < w->node_count; i++) {
        const node_t *n = &w->nodes[i];

        if (!ulab_streq(n->type, ULAB_NODE_TOWER)) {
            continue;
        }

        memset(args, 0, sizeof(args));
        rc = snprintf(args, sizeof(args), "%s %s", n->id, rt->run_dir);
        if (rc < 0 || (size_t)rc >= sizeof(args)) {
            snprintf(err->msg, sizeof(err->msg),
                     "enable-pcrf args too long for node %s", n->id);
            return ULAB_ERR;
        }

        ulab_status("PCRF", "enable service %s", n->id);
        if (run_script(rt, "enable-pcrf-service.sh", args, err)) {
            return ULAB_ERR;
        }
    }

    return ULAB_OK;
}

int runtime_ensure_media(runtime_t *rt, ulab_error_t *err) {

    char args[ULAB_MAX_ARGS];
    int rc;

    if (rt->repo[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime media requires --repo");
        return ULAB_ERR;
    }

    rc = snprintf(args, sizeof(args), "%s %s", rt->repo, rt->run_dir);
    if (rc < 0 || (size_t)rc >= sizeof(args)) {
        snprintf(err->msg, sizeof(err->msg),
                 "start-media args too long");
        return ULAB_ERR;
    }

    ulab_status("MEDIA", "start/ensure media target");
    if (run_script(rt, "start-media.sh", args, err)) {
        return ULAB_ERR;
    }

    memset(args, 0, sizeof(args));
    rc = snprintf(args, sizeof(args), "%s", rt->run_dir);
    if (rc < 0 || (size_t)rc >= sizeof(args)) {
        snprintf(err->msg, sizeof(err->msg),
                 "wait-media args too long");
        return ULAB_ERR;
    }

    ulab_status("MEDIA", "wait ready");
    if (run_script(rt, "wait-media-ready.sh", args, err)) {
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runtime_build_and_start_ues(const char *repo,
                                runtime_t *rt,
                                const world_t *w,
                                const selector_result_t *ues,
                                ulab_error_t *err) {
    size_t i;
    const char *repo_path;

    repo_path = (repo != NULL && repo[0] != '\0') ? repo : rt->repo;
    if (repo_path == NULL || repo_path[0] == '\0') {
        snprintf(err->msg, sizeof(err->msg),
                 "start-ue requires --repo");
        return ULAB_ERR;
    }

    for (i = 0; i < ues->count; i++) {
        const ue_t *ue = &w->ues[ues->idx[i]];
        char args[4096];
        int rc;

        if (ue->imsi[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "start-ue missing IMSI for ue %s", ue->id);
            return ULAB_ERR;
        }

        if (ue->iccid[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "start-ue missing ICCID for ue %s", ue->id);
            return ULAB_ERR;
        }

        if (ue->ip[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "start-ue missing IP for ue %s", ue->id);
            return ULAB_ERR;
        }

        ulab_status("UE", "start %s imsi=%s iccid=%s ip=%s site=%s",
                    ue->ref, ue->imsi, ue->iccid, ue->ip, ue->site_ref);

        rc = snprintf(args, sizeof(args), "%s %s %s %s %s %s %s %s",
                      repo_path,
                      ue->ref,
                      ue->id,
                      ue->imsi,
                      ue->iccid,
                      ue->ip,
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


static int cleanup_script(runtime_t *rt,
                          const char *script,
                          const char *args) {
    ulab_error_t err;

    memset(&err, 0, sizeof(err));
    if (run_script(rt, script, args, &err)) {
        if (rt->logf) {
            fprintf(rt->logf, "cleanup warning: %s %s: %s\n",
                    script, args ? args : "", err.msg);
            fflush(rt->logf);
        }
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runtime_stop_ues(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    char args[ULAB_MAX_ARGS];
    size_t i;
    int failures;
    int rc;

    failures = 0;

    if (rt == NULL || rt->run_dir[0] == '\0' || w == NULL) {
        return ULAB_OK;
    }

    for (i = 0; i < w->ue_count; i++) {
        rc = snprintf(args, sizeof(args), "%s %s",
                      w->ues[i].id, rt->run_dir);
        if (rc >= 0 && (size_t)rc < sizeof(args)) {
            if (cleanup_script(rt, "stop-ue.sh", args)) {
                failures++;
            }
        }
    }

    if (failures > 0 && err != NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "UE cleanup had %d failed step(s)", failures);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runtime_cleanup_infra(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    char args[ULAB_MAX_ARGS];
    size_t i;
    int failures;
    int rc;

    failures = 0;

    if (rt == NULL || rt->run_dir[0] == '\0') {
        return ULAB_OK;
    }

    rc = snprintf(args, sizeof(args), "%s", rt->run_dir);
    if (rc >= 0 && (size_t)rc < sizeof(args)) {
        if (cleanup_script(rt, "stop-media.sh", args)) {
            failures++;
        }
    }

    if (w != NULL) {
        for (i = 0; i < w->node_count; i++) {
            rc = snprintf(args, sizeof(args), "%s %s",
                          w->nodes[i].id, rt->run_dir);
            if (rc >= 0 && (size_t)rc < sizeof(args)) {
                if (cleanup_script(rt, "stop-node.sh", args)) {
                    failures++;
                }
            }
        }
    }

    rc = snprintf(args, sizeof(args), "%s", rt->run_dir);
    if (rc >= 0 && (size_t)rc < sizeof(args)) {
        if (cleanup_script(rt, "cleanup-network.sh", args)) {
            failures++;
        }
    }

    if (failures > 0 && err != NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime infra cleanup had %d failed step(s)", failures);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runtime_cleanup(runtime_t *rt, const world_t *w, ulab_error_t *err) {
    ulab_error_t tmp;
    int failures;

    failures = 0;
    memset(&tmp, 0, sizeof(tmp));

    if (runtime_stop_ues(rt, w, &tmp)) {
        failures++;
    }

    memset(&tmp, 0, sizeof(tmp));
    if (runtime_cleanup_infra(rt, w, &tmp)) {
        failures++;
    }

    if (failures > 0 && err != NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "runtime cleanup had %d failed section(s)", failures);
        return ULAB_ERR;
    }

    return ULAB_OK;
}
