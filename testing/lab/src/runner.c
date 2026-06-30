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
#include <time.h>
#include <unistd.h>
#include <jansson.h>

#include "runner.h"
#include "bff.h"
#include "event.h"
#include "runtime.h"
#include "log.h"
#include "selector.h"
#include "sim_factory.h"
#include "util.h"

static void make_run_id(char *out, size_t len, const scenario_t *scenario,
                        const runner_opts_t *opts) {

    time_t now;

    if (opts != NULL && opts->run_id[0] != '\0') {
        snprintf(out, len, "%s", opts->run_id);
        return;
    }

    now = time(NULL);
    snprintf(out, len, "lab-%s-%u-%ld", scenario->name,
             scenario->seed, (long)now);
}

static int prepare_run(const runner_opts_t *opts,
                       scenario_t *scenario,
                       world_t *world,
                       model_t *model,
                       char *runDir,
                       size_t runDirLen,
                       ulab_error_t *err) {

    char run_id[ULAB_MAX_ID];

    if (scenario_load(opts->scenario_path, scenario, err)) {
        return ULAB_ESCENARIO;
    }

    if (opts->has_seed_override) {
        scenario->seed = opts->seed_override;
    }

    if (scenario_validate(scenario, err)) {
        return ULAB_ESCENARIO;
    }

    make_run_id(run_id, sizeof(run_id), scenario, opts);
    snprintf(runDir, runDirLen, "%s/%s", opts->out_dir, run_id);

    if (ulab_mkdir_p(runDir)) {
        snprintf(err->msg, sizeof(err->msg), "failed to create run dir");
        return ULAB_EINTERNAL;
    }

    if (world_generate(scenario, run_id, world, err)) {
        return ULAB_EINTERNAL;
    }

    if (model_init(model, world)) {
        snprintf(err->msg, sizeof(err->msg), "model init failed");
        return ULAB_EINTERNAL;
    }

    return ULAB_OK;
}

static void created_path(const char *run_dir, char *path, size_t path_len) {
    snprintf(path, path_len, "%s/created.json", run_dir);
}

static int created_exists(const char *run_dir) {
    char path[ULAB_MAX_PATH];

    created_path(run_dir, path, sizeof(path));
    return access(path, F_OK) == 0;
}

static void write_json_string_array(FILE *f,
                                    const char *name,
                                    const char ids[][ULAB_MAX_ID],
                                    size_t count,
                                    int comma) {
    size_t i;
    int first;

    fprintf(f, "  \"%s\": [", name);
    first = 1;
    for (i = 0; i < count; i++) {
        char esc[ULAB_MAX_ID * 2];

        if (!first) {
            fprintf(f, ", ");
        }
        ulab_json_escape(ids[i], esc, sizeof(esc));
        fprintf(f, "\"%s\"", esc);
        first = 0;
    }
    fprintf(f, "]%s\n", comma ? "," : "");
}

static int created_write(const char *run_dir, const world_t *world) {
    char path[ULAB_MAX_PATH];
    char tmp[ULAB_MAX_PATH];
    char (*ids)[ULAB_MAX_ID];
    const char *suffix = ".tmp";
    size_t suffix_len = 4;
    size_t path_len;
    FILE *f;
    size_t max;
    size_t i;

    if (run_dir == NULL || world == NULL) {
        return ULAB_ERR;
    }

    created_path(run_dir, path, sizeof(path));

    path_len = strlen(path);
    if (path_len + suffix_len + 1 > sizeof(tmp)) {
        fprintf(stderr, "created path too long: %s\n", path);
        return ULAB_ERR;
    }

    memcpy(tmp, path, path_len);
    memcpy(tmp + path_len, suffix, suffix_len + 1);

    max = world->network_count;
    if (world->site_count > max) max = world->site_count;
    if (world->node_count > max) max = world->node_count;
    if (world->subscriber_count > max) max = world->subscriber_count;
    if (world->package_count > max) max = world->package_count;
    if (max == 0) max = 1;

    ids = calloc(max, sizeof(*ids));
    if (ids == NULL) {
        return ULAB_ERR;
    }

    f = fopen(tmp, "w");
    if (f == NULL) {
        free(ids);
        return ULAB_ERR;
    }

    fprintf(f, "{\n");
    fprintf(f, "  \"run_id\": \"%s\",\n", world->run_id);

    memset(ids, 0, max * sizeof(*ids));
    for (i = 0; i < world->network_count; i++) {
        ulab_copy(ids[i], ULAB_MAX_ID, world->networks[i].bff_id);
    }
    write_json_string_array(f, "networks", ids, world->network_count, 1);

    memset(ids, 0, max * sizeof(*ids));
    for (i = 0; i < world->site_count; i++) {
        ulab_copy(ids[i], ULAB_MAX_ID, world->sites[i].bff_id);
    }
    write_json_string_array(f, "sites", ids, world->site_count, 1);

    memset(ids, 0, max * sizeof(*ids));
    for (i = 0; i < world->node_count; i++) {
        ulab_copy(ids[i], ULAB_MAX_ID, world->nodes[i].bff_id);
    }
    write_json_string_array(f, "nodes", ids, world->node_count, 1);

    memset(ids, 0, max * sizeof(*ids));
    for (i = 0; i < world->subscriber_count; i++) {
        ulab_copy(ids[i], ULAB_MAX_ID, world->subscribers[i].bff_id);
    }
    write_json_string_array(f, "subscribers", ids, world->subscriber_count, 1);

    memset(ids, 0, max * sizeof(*ids));
    for (i = 0; i < world->package_count; i++) {
        ulab_copy(ids[i], ULAB_MAX_ID, world->packages[i].bff_id);
    }
    write_json_string_array(f, "packages", ids, world->package_count, 1);

    fprintf(f, "  \"sims\": [");
    for (i = 0; i < world->ue_count; i++) {
        char id[ULAB_MAX_ID * 2];
        char pool[ULAB_MAX_ID * 2];
        char pkg[ULAB_MAX_ID * 2];

        if (i > 0) fprintf(f, ", ");
        ulab_json_escape(world->ues[i].bff_id, id, sizeof(id));
        ulab_json_escape(world->ues[i].pool_sim_id, pool, sizeof(pool));
        ulab_json_escape(world->ues[i].sim_package_id, pkg, sizeof(pkg));
        fprintf(f, "{\"id\":\"%s\",\"pool_sim_id\":\"%s\","
                   "\"package_id\":\"%s\"}", id, pool, pkg);
    }
    fprintf(f, "]\n");
    fprintf(f, "}\n");

    if (fclose(f) != 0) {
        free(ids);
        unlink(tmp);
        return ULAB_ERR;
    }

    free(ids);

    if (rename(tmp, path) != 0) {
        unlink(tmp);
        return ULAB_ERR;
    }

    return ULAB_OK;
}

static void copy_array_value(json_t *root,
                             const char *name,
                             size_t idx,
                             char *out,
                             size_t out_len) {
    json_t *arr;
    json_t *v;

    arr = json_object_get(root, name);
    if (arr == NULL || !json_is_array(arr) || idx >= json_array_size(arr)) {
        return;
    }

    v = json_array_get(arr, idx);
    if (v != NULL && json_is_string(v) && json_string_value(v) != NULL) {
        ulab_copy(out, out_len, json_string_value(v));
    }
}

static int created_load(const char *run_dir, world_t *world) {
    char path[ULAB_MAX_PATH];
    json_error_t jerr;
    json_t *root;
    json_t *arr;
    size_t i;

    created_path(run_dir, path, sizeof(path));
    root = json_load_file(path, 0, &jerr);
    if (root == NULL || !json_is_object(root)) {
        if (root != NULL) json_decref(root);
        return ULAB_ERR;
    }

    for (i = 0; i < world->network_count; i++) {
        copy_array_value(root, "networks", i, world->networks[i].bff_id,
                         sizeof(world->networks[i].bff_id));
    }
    for (i = 0; i < world->site_count; i++) {
        copy_array_value(root, "sites", i, world->sites[i].bff_id,
                         sizeof(world->sites[i].bff_id));
    }
    for (i = 0; i < world->node_count; i++) {
        copy_array_value(root, "nodes", i, world->nodes[i].bff_id,
                         sizeof(world->nodes[i].bff_id));
    }
    for (i = 0; i < world->subscriber_count; i++) {
        copy_array_value(root, "subscribers", i, world->subscribers[i].bff_id,
                         sizeof(world->subscribers[i].bff_id));
    }
    for (i = 0; i < world->package_count; i++) {
        copy_array_value(root, "packages", i, world->packages[i].bff_id,
                         sizeof(world->packages[i].bff_id));
    }

    arr = json_object_get(root, "sims");
    if (arr != NULL && json_is_array(arr)) {
        for (i = 0; i < json_array_size(arr) && i < world->ue_count; i++) {
            json_t *obj = json_array_get(arr, i);
            json_t *id;
            json_t *pool;
            json_t *pkg;

            if (obj == NULL || !json_is_object(obj)) {
                continue;
            }
            id = json_object_get(obj, "id");
            pool = json_object_get(obj, "pool_sim_id");
            pkg = json_object_get(obj, "package_id");
            if (id != NULL && json_is_string(id)) {
                ulab_copy(world->ues[i].bff_id, sizeof(world->ues[i].bff_id),
                          json_string_value(id));
            }
            if (pool != NULL && json_is_string(pool)) {
                ulab_copy(world->ues[i].pool_sim_id,
                          sizeof(world->ues[i].pool_sim_id),
                          json_string_value(pool));
            }
            if (pkg != NULL && json_is_string(pkg)) {
                ulab_copy(world->ues[i].sim_package_id,
                          sizeof(world->ues[i].sim_package_id),
                          json_string_value(pkg));
            }
        }
    }

    json_decref(root);
    return ULAB_OK;
}

static void created_clear_ids(world_t *world) {
    size_t i;

    for (i = 0; i < world->network_count; i++) world->networks[i].bff_id[0] = '\0';
    for (i = 0; i < world->site_count; i++) world->sites[i].bff_id[0] = '\0';
    for (i = 0; i < world->node_count; i++) world->nodes[i].bff_id[0] = '\0';
    for (i = 0; i < world->subscriber_count; i++) world->subscribers[i].bff_id[0] = '\0';
    for (i = 0; i < world->package_count; i++) world->packages[i].bff_id[0] = '\0';
    for (i = 0; i < world->ue_count; i++) {
        world->ues[i].bff_id[0] = '\0';
        world->ues[i].sim_package_id[0] = '\0';
        world->ues[i].pool_sim_id[0] = '\0';
    }
}

static void created_finalize(const char *run_dir) {
    char path[ULAB_MAX_PATH];
    char final_path[ULAB_MAX_PATH];

    if (run_dir == NULL || run_dir[0] == '\0') {
        return;
    }

    created_path(run_dir, path, sizeof(path));
    snprintf(final_path, sizeof(final_path), "%s/created.final.json", run_dir);
    unlink(final_path);
    rename(path, final_path);
}

static int setup_bff_networks(bff_client_t *bff,
                              const scenario_t *scenario,
                              world_t *world,
                              ulab_error_t *err) {

    size_t i;

    if (!scenario->setup.create_networks) {
        return ULAB_OK;
    }

    for (i = 0; i < world->network_count; i++) {
        ulab_status("NETWORK", "add network %s", world->networks[i].ref);
        if (bff_add_network(bff, &world->networks[i], err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static int setup_bff_sites(bff_client_t *bff,
                           const scenario_t *scenario,
                           world_t *world,
                           ulab_error_t *err) {

    network_t *network;
    size_t i;

    if (!scenario->setup.create_sites) {
        return ULAB_OK;
    }

    for (i = 0; i < world->site_count; i++) {
        network = world_network_by_ref(world, world->sites[i].network_ref);
        if (network == NULL) {
            snprintf(err->msg, sizeof(err->msg),
                     "site %s has invalid network ref",
                     world->sites[i].ref);
            return ULAB_EBFF;
        }

        if (bff_wait_site_anchor_online(bff, &world->sites[i], err)) {
            return ULAB_EBFF;
        }

        ulab_status("SITE", "add site %s", world->sites[i].ref);
        if (bff_add_site(bff, &world->sites[i], network, err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static int setup_bff_packages(bff_client_t *bff,
                              const scenario_t *scenario,
                              world_t *world,
                              ulab_error_t *err) {

    network_t *network;
    package_t *package;
    size_t i;

    if (!scenario->setup.create_packages) {
        return ULAB_OK;
    }

    for (i = 0; i < world->package_count; i++) {
        package = &world->packages[i];
        network = world_network_by_ref(world, package->network_ref);
        if (network == NULL || network->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "package %s has invalid network id",
                     package->ref);
            return ULAB_EBFF;
        }

        ulab_status("PACKAGE", "add package %s", package->ref);
        if (bff_add_package(bff, package, network, err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static int setup_bff_subscribers(bff_client_t *bff,
                                 const scenario_t *scenario,
                                 world_t *world,
                                 ulab_error_t *err) {

    network_t *network;
    subscriber_t *sub;
    size_t i;

    if (!scenario->setup.create_subscribers) {
        return ULAB_OK;
    }

    for (i = 0; i < world->subscriber_count; i++) {
        sub = &world->subscribers[i];
        network = world_network_by_ref(world, sub->network_ref);
        if (network == NULL || network->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "subscriber %s has invalid network id",
                     sub->ref);
            return ULAB_EBFF;
        }

        ulab_status("SUBSCRIBER", "add subscriber %s", sub->ref);
        if (bff_add_subscriber(bff, sub, network, err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static subscriber_t *find_subscriber(world_t *world,
                                     const char *ref) {

    size_t i;

    for (i = 0; i < world->subscriber_count; i++) {
        if (ulab_streq(world->subscribers[i].ref, ref)) {
            return &world->subscribers[i];
        }
    }

    return NULL;
}

static int setup_bff_sim_pool(bff_client_t *bff,
                              world_t *world,
                              const runner_opts_t *opts,
                              const char *run_dir,
                              ulab_error_t *err) {

    char (*pool_iccids)[ULAB_MAX_ID];
    char (*pool_ids)[ULAB_MAX_ID];
    char factory_csv[ULAB_MAX_PATH];
    size_t pool_count;
    size_t max_pool;
    size_t i;
    size_t j;

    if (world->ue_count == 0) {
        return ULAB_OK;
    }

    memset(factory_csv, 0, sizeof(factory_csv));
    if (sim_factory_prepare_world(opts, world, run_dir, factory_csv,
                                  sizeof(factory_csv), err)) {
        return ULAB_EBFF;
    }

    ulab_status("SIMPOOL", "upload %s type=%s", factory_csv,
                opts->sim_type);
    if (bff_upload_sims_from_csv(bff, factory_csv, opts->sim_type, err)) {
        return ULAB_EBFF;
    }

    max_pool = world->ue_count * 32;
    if (max_pool < 1024) {
        max_pool = 1024;
    }

    pool_iccids = calloc(max_pool, sizeof(*pool_iccids));
    pool_ids = calloc(max_pool, sizeof(*pool_ids));
    if (pool_iccids == NULL || pool_ids == NULL) {
        snprintf(err->msg, sizeof(err->msg),
                 "out of memory reading SIM pool");
        free(pool_iccids);
        free(pool_ids);
        return ULAB_EINTERNAL;
    }

    ulab_status("SIMPOOL", "get unassigned sims type=%s", opts->sim_type);
    if (bff_get_sims_from_pool(bff, opts->sim_type, pool_iccids, pool_ids,
                               max_pool, &pool_count, err)) {
        free(pool_iccids);
        free(pool_ids);
        return ULAB_EBFF;
    }

    for (i = 0; i < world->ue_count; i++) {
        int found;

        found = 0;
        for (j = 0; j < pool_count; j++) {
            if (!ulab_streq(world->ues[i].iccid, pool_iccids[j])) {
                continue;
            }

            ulab_copy(world->ues[i].pool_sim_id,
                      sizeof(world->ues[i].pool_sim_id), pool_ids[j]);
            found = 1;
            break;
        }

        if (!found) {
            snprintf(err->msg, sizeof(err->msg),
                     "prepared factory SIM iccid=%s for ue=%s is not "
                     "available as UNASSIGNED in SIM pool",
                     world->ues[i].iccid, world->ues[i].ref);
            free(pool_iccids);
            free(pool_ids);
            return ULAB_EBFF;
        }

        ulab_status("SIMPOOL", "match ue %s iccid=%s imsi=%s pool=%s",
                    world->ues[i].ref, world->ues[i].iccid,
                    world->ues[i].imsi, world->ues[i].pool_sim_id);
    }

    free(pool_iccids);
    free(pool_ids);

    return ULAB_OK;
}

static int setup_bff_sims(bff_client_t *bff,
                          const scenario_t *scenario,
                          world_t *world,
                          const runner_opts_t *opts,
                          ulab_error_t *err) {

    subscriber_t *sub;
    network_t *network;
    package_t *package;
    ue_t *ue;
    int active;
    size_t i;

    if (!scenario->setup.create_sims) {
        return ULAB_OK;
    }

    for (i = 0; i < world->ue_count; i++) {
        ue = &world->ues[i];

        sub = find_subscriber(world, ue->subscriber_ref);
        network = world_network_by_ref(world, ue->network_ref);
        package = world_package_for_network(world, ue->package_ref,
                                            ue->network_ref);

        if (sub == NULL || !ulab_streq(sub->ref, ue->subscriber_ref)) {
            snprintf(err->msg, sizeof(err->msg),
                     "ue %s has invalid subscriber ref", ue->ref);
            return ULAB_EBFF;
        }

        if (network == NULL || network->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "ue %s has invalid network id", ue->ref);
            return ULAB_EBFF;
        }

        if (package == NULL || package->bff_id[0] == '\0') {
            snprintf(err->msg, sizeof(err->msg),
                     "ue %s has invalid package id", ue->ref);
            return ULAB_EBFF;
        }

        /*
         * allocateSim is authoritative for this setup path.
         * We pass package_id into allocateSim, so it assigns the SIM to the
         * subscriber/network and binds the package. Do not clear, re-add, or
         * explicitly activate the SIM here; doing so creates duplicate package
         * rows and can fail when the SIM is already active.
         */
        ulab_status("SIM", "allocate sim %s iccid=%s", ue->ref,
                    ue->iccid);
        if (bff_allocate_sim_from_pool(bff, ue, sub, network, package,
                                       opts->sim_type, err)) {
            return ULAB_EBFF;
        }

        active = 0;
        ulab_status("SIM", "verify sim package %s package=%s",
                    ue->ref, package->ref);
        if (bff_get_packages_for_sim(bff, ue, package->bff_id, &active,
                                     err)) {
            return ULAB_EBFF;
        }

        if (!active) {
            snprintf(err->msg, sizeof(err->msg),
                     "ue %s sim allocated but package %s is not active",
                     ue->ref, package->bff_id);
            return ULAB_EBFF;
        }

        ulab_copy(ue->sim_package_id,
                  sizeof(ue->sim_package_id),
                  package->bff_id);

        if (sim_factory_wait_asr(opts, ue, err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static int setup_bff_world(bff_client_t *bff,
                           const scenario_t *scenario,
                           world_t *world,
                           const runner_opts_t *opts,
                           const char *run_dir,
                           ulab_error_t *err) {

    if (setup_bff_networks(bff, scenario, world, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    if (setup_bff_sites(bff, scenario, world, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    if (setup_bff_sim_pool(bff, world, opts, run_dir, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    if (setup_bff_packages(bff, scenario, world, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    if (setup_bff_subscribers(bff, scenario, world, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    if (setup_bff_sims(bff, scenario, world, opts, err)) {
        return ULAB_EBFF;
    }
    created_write(run_dir, world);

    return ULAB_OK;
}

static int start_runtime_sites(const char *repo,
                               const scenario_t *scenario,
                               world_t *world,
                               runtime_t *runtime,
                               ulab_error_t *err) {

    if (!scenario->runtime.start_nodes) {
        return ULAB_OK;
    }

    return runtime_build_and_start_sites(repo, runtime, world, err) ?
        ULAB_ERUNTIME : ULAB_OK;
}

static int wait_runtime_nodes(const scenario_t *scenario,
                              world_t *world,
                              runtime_t *runtime,
                              ulab_error_t *err) {

    selector_t all;
    selector_result_t nodes;
    int rc;

    if (!scenario->runtime.wait_nodes_ready) {
        return ULAB_OK;
    }

    memset(&all, 0, sizeof(all));
    memset(&nodes, 0, sizeof(nodes));
    all.kind = SEL_ALL;

    rc = selector_resolve_nodes(world, &all, &nodes, err);
    if (rc != ULAB_OK) {
        return ULAB_ERUNTIME;
    }

    ulab_status("NODE", "wait nodes ready");
    rc = runtime_wait_nodes_ready(runtime, world, &nodes, err);
    selector_result_free(&nodes);

    return rc == ULAB_OK ? ULAB_OK : ULAB_ERUNTIME;
}

static int runtime_all_ues(const scenario_t *scenario,
                           world_t *world,
                           runtime_t *runtime,
                           ulab_error_t *err) {

    selector_t all;
    selector_result_t ues;
    int rc;

    memset(&all, 0, sizeof(all));
    memset(&ues, 0, sizeof(ues));
    all.kind = SEL_ALL;

    if (!scenario->runtime.start_ues &&
        !scenario->runtime.wait_ues_attached) {
        return ULAB_OK;
    }

    rc = selector_resolve_ues(world, &all, &ues, err);
    if (rc != ULAB_OK) {
        return ULAB_ERUNTIME;
    }

    if (scenario->runtime.start_ues) {
        ulab_status("MEDIA", "start media");
        rc = runtime_ensure_media(runtime, err);
        if (rc != ULAB_OK) {
            selector_result_free(&ues);
            return ULAB_ERUNTIME;
        }

        ulab_status("UE", "start ues");
        rc = runtime_build_and_start_ues(NULL, runtime, world, &ues, err);
        if (rc != ULAB_OK) {
            selector_result_free(&ues);
            return ULAB_ERUNTIME;
        }
    }

    if (scenario->runtime.wait_ues_attached) {
        ulab_status("UE", "wait ues attached");
        rc = runtime_wait_ues_attached(runtime, world, &ues, err);
        if (rc != ULAB_OK) {
            selector_result_free(&ues);
            return ULAB_ERUNTIME;
        }
    }

    selector_result_free(&ues);

    return ULAB_OK;
}

static int run_checks(check_ctx_t *ctx,
                      const check_spec_t *checks,
                      size_t count,
                      report_t *report,
                      ulab_error_t *err) {

    check_result_t result;
    size_t i;
    int failed;

    failed = 0;

    for (i = 0; i < count; i++) {
        if (check_run(ctx, &checks[i], &result, err)) {
            return ULAB_ERR;
        }

        report_check(report, &result);
        if (!result.passed && !result.skipped) {
            failed = 1;
        }
    }

    return failed ? ULAB_ERR : ULAB_OK;
}

static void init_check_ctx(check_ctx_t *ctx,
                           const scenario_t *scenario,
                           world_t *world,
                           model_t *model,
                           bff_client_t *bff,
                           runtime_t *runtime) {

    memset(ctx, 0, sizeof(*ctx));
    ctx->scenario = scenario;
    ctx->world    = world;
    ctx->model    = model;
    ctx->bff      = bff;
    ctx->runtime  = runtime;
}

static void init_event_ctx(event_ctx_t *ctx,
                           scenario_t *scenario,
                           world_t *world,
                           model_t *model,
                           bff_client_t *bff,
                           runtime_t *runtime,
                           const char *phaseName) {

    memset(ctx, 0, sizeof(*ctx));
    ctx->scenario   = scenario;
    ctx->world      = world;
    ctx->model      = model;
    ctx->bff        = bff;
    ctx->runtime    = runtime;
    ctx->phaseName  = phaseName;
}

static int run_phase(scenario_t *scenario,
                     world_t *world,
                     model_t *model,
                     bff_client_t *bff,
                     runtime_t *runtime,
                     report_t *report,
                     phase_spec_t *phase,
                     ulab_error_t *err) {

    event_ctx_t event_ctx;
    check_ctx_t check_ctx;
    size_t i;
    int rc;

    ulab_status("PHASE", "%s", phase->name);
    init_event_ctx(&event_ctx, scenario, world, model, bff,
                   runtime, phase->name);

    for (i = 0; i < phase->event_count; i++) {
        rc = event_run(&event_ctx, &phase->events[i], err);
        report_event(report, phase->name, &phase->events[i],
                     rc == ULAB_OK, rc == ULAB_OK ? "ok" : err->msg);
        if (rc != ULAB_OK) {
            return rc;
        }
    }

    init_check_ctx(&check_ctx, scenario, world, model, bff, runtime);
    rc = run_checks(&check_ctx, phase->checks, phase->check_count,
                    report, err);

    return rc;
}

static void write_world_artifact(const world_t *world,
                                 const char *runDir) {

    char path[ULAB_MAX_PATH * 2];

    snprintf(path, sizeof(path), "%s/world.json", runDir);
    world_write_json(world, path);
}

static void write_model_artifact(const model_t *model,
                                 const char *runDir) {

    char path[ULAB_MAX_PATH * 2];

    snprintf(path, sizeof(path), "%s/model.json", runDir);
    model_write_json(model, path);
}

static int cleanup_run(bff_client_t *bff,
                       runtime_t *runtime,
                       world_t *world,
                       const char *run_dir) {

    ulab_error_t cleanup_err;
    int failures;

    failures = 0;

    memset(&cleanup_err, 0, sizeof(cleanup_err));
    ulab_status("CLEANUP", "stop UE runtime");
    if (runtime_stop_ues(runtime, world, &cleanup_err)) {
        failures++;
        ulab_log_error("%s", cleanup_err.msg);
    }

    memset(&cleanup_err, 0, sizeof(cleanup_err));
    ulab_status("CLEANUP", "delete backend resources");
    if (bff_cleanup_world(bff, world, &cleanup_err)) {
        failures++;
        ulab_log_error("%s", cleanup_err.msg);
    }

    memset(&cleanup_err, 0, sizeof(cleanup_err));
    ulab_status("CLEANUP", "stop media/nodes/network");
    if (runtime_cleanup_infra(runtime, world, &cleanup_err)) {
        failures++;
        ulab_log_error("%s", cleanup_err.msg);
    }

    if (failures == 0) {
        created_finalize(run_dir);
    }

    return failures == 0 ? ULAB_OK : ULAB_ERR;
}

static int preclean_existing_run(bff_client_t *bff,
                                 runtime_t *runtime,
                                 world_t *world,
                                 const char *run_dir,
                                 ulab_error_t *err) {
    int rc;

    if (!created_exists(run_dir)) {
        return ULAB_OK;
    }

    ulab_status("CLEANUP", "existing created.json found; cleaning run first");

    if (created_load(run_dir, world)) {
        snprintf(err->msg, sizeof(err->msg),
                 "failed to load existing created.json for cleanup");
        return ULAB_ERR;
    }

    rc = cleanup_run(bff, runtime, world, run_dir);
    created_clear_ids(world);

    if (rc != ULAB_OK) {
        snprintf(err->msg, sizeof(err->msg),
                 "pre-cleanup failed for existing run_id");
        return ULAB_ERR;
    }

    return ULAB_OK;
}

int runner_validate(const runner_opts_t *opts) {

    scenario_t *scenario;
    world_t world;
    model_t model;
    bff_client_t bff;
    runtime_t runtime;
    report_t report;
    ulab_error_t err;
    char runDir[ULAB_MAX_PATH];
    check_ctx_t check_ctx;
    size_t i;
    int rc;
    int cleanup_rc;
    int skip_cleanup;

    scenario = NULL;
    rc = ULAB_OK;
    cleanup_rc = ULAB_OK;
    skip_cleanup = 0;
    memset(&world,   0, sizeof(world));
    memset(&model,   0, sizeof(model));
    memset(&bff,     0, sizeof(bff));
    memset(&runtime, 0, sizeof(runtime));
    memset(&report,  0, sizeof(report));
    memset(&err,     0, sizeof(err));
    memset(runDir,   0, sizeof(runDir));

    scenario = calloc(1, sizeof(*scenario));
    if (scenario == NULL) {
        return ULAB_EINTERNAL;
    }

    rc = prepare_run(opts, scenario, &world, &model, runDir,
                     sizeof(runDir), &err);
    if (rc != ULAB_OK) {
        skip_cleanup = 1;
        goto done;
    }

    if (report_open(&report, scenario, world.run_id, runDir)) {
        snprintf(err.msg, sizeof(err.msg), "failed to open report files");
        rc = ULAB_EINTERNAL;
        goto done;
    }
    report_world(&world);
    write_world_artifact(&world, runDir);

    if (ulab_streq(scenario->status, "skip") ||
        ulab_streq(scenario->status, "wip")) {
        ulab_status("SKIP", "%s status=%s", scenario->name,
                    scenario->status);
        skip_cleanup = 1;
        rc = ULAB_OK;
        goto done;
    }

    rc = runtime_init(&runtime, scenario->provider.type, opts->script_dir, runDir, opts->repo);
    if (rc != ULAB_OK) {
        rc = ULAB_ERUNTIME;
        goto done;
    }

    rc = bff_init(&bff, opts->bff_url, runDir);
    if (rc != ULAB_OK) {
        rc = ULAB_EBFF;
        goto done;
    }

    rc = preclean_existing_run(&bff, &runtime, &world, runDir, &err);
    if (rc != ULAB_OK) {
        rc = ULAB_EBFF;
        goto done;
    }

    rc = runtime_ensure_network(&runtime, &err);
    if (rc != ULAB_OK) {
        rc = ULAB_ERUNTIME;
        goto done;
    }

    ulab_status("SITE", "factory/build/start site node bundles");
    rc = start_runtime_sites(opts->repo, scenario, &world, &runtime, &err);
    if (rc != ULAB_OK) {
        goto done;
    }

    rc = wait_runtime_nodes(scenario, &world, &runtime, &err);
    if (rc != ULAB_OK) {
        goto done;
    }

    rc = runtime_enable_pcrf_service(&runtime, &world, &err);
    if (rc != ULAB_OK) {
        rc = ULAB_ERUNTIME;
        goto done;
    }

    ulab_status("BACKEND", "creating backend world resources");
    rc = setup_bff_world(&bff, scenario, &world, opts, runDir, &err);
    if (rc != ULAB_OK) {
        goto done;
    }

    model_sync_world(&model, &world);
    write_world_artifact(&world, runDir);

    rc = runtime_all_ues(scenario, &world, &runtime, &err);
    if (rc != ULAB_OK) {
        goto done;
    }

    for (i = 0; i < scenario->phase_count; i++) {
        rc = run_phase(scenario, &world, &model, &bff, &runtime,
                       &report, &scenario->phases[i], &err);
        if (rc != ULAB_OK) {
            goto done;
        }
    }

    init_check_ctx(&check_ctx, scenario, &world, &model, &bff, &runtime);
    rc = run_checks(&check_ctx, scenario->final_checks,
                    scenario->final_check_count, &report, &err);

    write_model_artifact(&model, runDir);

    if (rc == ULAB_OK && (report.failed || report.event_failed)) {
        rc = ULAB_ERR;
    }

done:
    if (err.msg[0] != '\0') {
        ulab_log_error("%s", err.msg);
    }

    if (!skip_cleanup) {
        cleanup_rc = cleanup_run(&bff, &runtime, &world, runDir);
        report_set_cleanup(&report, cleanup_rc != ULAB_OK);
        if (cleanup_rc != ULAB_OK && rc == ULAB_OK) {
            rc = ULAB_ERR;
        }
    }

    if (scenario != NULL && ulab_streq(scenario->status, "xfail") &&
        rc != ULAB_OK) {
        ulab_status("XFAIL", "%s failed as expected", scenario->name);
        rc = ULAB_OK;
    }

    report_set_final_rc(&report, rc);
    report_result(&report);
    report_close(&report);
    runtime_close(&runtime);
    bff_close(&bff);
    world_free(&world);
    model_free(&model);
    free(scenario);

    return rc;
}
