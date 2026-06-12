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

#include "runner.h"
#include "event.h"
#include "log.h"
#include "selector.h"
#include "util.h"

static void make_run_id(char *out, size_t len, const scenario_t *scenario) {

    time_t now;

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

    make_run_id(run_id, sizeof(run_id), scenario);
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

static int setup_bff(bff_client_t *bff,
                     world_t *world,
                     ulab_error_t *err) {

    network_t *network;
    site_t *site;
    package_t *package;
    subscriber_t *subscriber;
    size_t i;

    for (i = 0; i < world->network_count; i++) {
        if (bff_add_network(bff, &world->networks[i], err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->site_count; i++) {
        network = world_network_by_ref(world, world->sites[i].network_ref);
        if (bff_add_site(bff, &world->sites[i], network, err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->node_count; i++) {
        if (bff_add_node(bff, &world->nodes[i], err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->node_count; i++) {
        site = world_site_by_ref(world, world->nodes[i].site_ref);
        network = world_network_by_ref(world, world->nodes[i].network_ref);

        if (bff_add_node_to_site(bff, &world->nodes[i], site,
                                 network, err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->package_count; i++) {
        if (bff_add_package(bff, &world->packages[i], err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->subscriber_count; i++) {
        network = world_network_by_ref(world,
                                       world->subscribers[i].network_ref);
        if (bff_add_subscriber(bff, &world->subscribers[i],
                               network, err)) {
            return ULAB_EBFF;
        }
    }

    for (i = 0; i < world->ue_count; i++) {
        subscriber = &world->subscribers[i];
        network = world_network_by_ref(world, world->ues[i].network_ref);
        package = world_package_by_ref(world, world->ues[i].package_ref);

        if (bff_allocate_sim(bff, &world->ues[i], subscriber,
                             network, package, err)) {
            return ULAB_EBFF;
        }
    }

    return ULAB_OK;
}

static int start_runtime(const char *repo,
                         const scenario_t *scenario,
                         world_t *world,
                         runtime_t *runtime,
                         ulab_error_t *err) {

    selector_result_t result;
    selector_t selector;
    
    memset(&selector, 0, sizeof(selector));
    selector.kind = SEL_ALL;

    if (scenario->runtime.start_nodes) {

        if (selector_resolve_nodes(world, &selector, &result, err)) {
            return ULAB_ERUNTIME;
        }

        if (runtime_build_and_start_nodes(repo,runtime, world, &result, err)) {
            selector_result_free(&result);
            return ULAB_ERUNTIME;
        }

        selector_result_free(&result);
    }

    if (scenario->runtime.wait_nodes_ready) {

        if (selector_resolve_nodes(world, &selector, &result, err)) {
            return ULAB_ERUNTIME;
        }

        if (runtime_wait_nodes_ready(runtime, world, &result, err)) {
            selector_result_free(&result);
            return ULAB_ERUNTIME;
        }

        selector_result_free(&result);
    }

    if (scenario->runtime.start_ues) {

        if (selector_resolve_ues(world, &selector, &result, err)) {
            return ULAB_ERUNTIME;
        }

        if (runtime_build_and_start_ues(repo, runtime, world, &result, err)) {
            selector_result_free(&result);
            return ULAB_ERUNTIME;
        }

        selector_result_free(&result);
    }

    if (scenario->runtime.wait_ues_attached) {

        if (selector_resolve_ues(world, &selector, &result, err)) {
            return ULAB_ERUNTIME;
        }

        if (runtime_wait_ues_attached(runtime, world, &result, err)) {
            selector_result_free(&result);
            return ULAB_ERUNTIME;
        }

        selector_result_free(&result);
    }

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

    scenario = NULL;
    rc = ULAB_OK;
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
        goto done;
    }

    report_open(&report, world.run_id, runDir);
    report_world(&world);
    write_world_artifact(&world, runDir);

    rc = bff_init(&bff, opts->bff_url, runDir);
    if (rc != ULAB_OK) {
        rc = ULAB_EBFF;
        goto done;
    }

    rc = runtime_init(&runtime, opts->script_dir, runDir);
    if (rc != ULAB_OK) {
        rc = ULAB_ERUNTIME;
        goto done;
    }

    ulab_status("SETUP", "creating product world via BFF");
    rc = setup_bff(&bff, &world, &err);
    if (rc != ULAB_OK) {
        goto done;
    }

    model_sync_world(&model, &world);
    if (!opts->setup_only) {
        ulab_status("RUNTIME", "starting real nodes and UEs");
        rc = start_runtime(opts->repo,
                           scenario,
                           &world,
                           &runtime,
                           &err);
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
    }

    init_check_ctx(&check_ctx, scenario, &world, &model, &bff, &runtime);
    rc = run_checks(&check_ctx, scenario->final_checks,
                    scenario->final_check_count, &report, &err);

    write_model_artifact(&model, runDir);
    report_result(&report);

    if (rc == ULAB_OK && report.failed) {
        rc = ULAB_ERR;
    }

done:
    if (err.msg[0] != '\0') {
        ulab_log_error("%s", err.msg);
    }

    report_close(&report);
    runtime_close(&runtime);
    bff_close(&bff);
    world_free(&world);
    model_free(&model);
    free(scenario);

    return rc;
}

int runner_dry_run(const runner_opts_t *opts) {

    scenario_t *scenario;
    world_t world;
    model_t model;
    ulab_error_t err;
    char runDir[ULAB_MAX_PATH];
    int rc;

    scenario = NULL;
    memset(&world, 0, sizeof(world));
    memset(&model, 0, sizeof(model));
    memset(&err,   0, sizeof(err));
    memset(runDir, 0, sizeof(runDir));

    scenario = calloc(1, sizeof(*scenario));
    if (scenario == NULL) {
        return ULAB_EINTERNAL;
    }

    rc = prepare_run(opts, scenario, &world, &model, runDir,
                     sizeof(runDir), &err);
    if (rc != ULAB_OK) {
        ulab_log_error("%s", err.msg);
        goto done;
    }

    ulab_status("DRY-RUN", "%s", scenario->name);
    report_world(&world);

    if (opts->print_world) {
        world_print(&world);
    }

done:
    world_free(&world);
    model_free(&model);
    free(scenario);

    return rc;
}
