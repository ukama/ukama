/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "runtime.h"
#include "bff.h"
#include "sim_factory.h"
#include "util.h"
#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/wait.h>
#include <unistd.h>

typedef struct {
    const runner_opts_t *opts;
    const scenario_t *scenario;
    world_t *world;
    runtime_t runtime;
    webapp_client_t *client;
    int started;
    int child;
    int uncertain;
    const char *run_dir;
} local_t;
typedef struct { local_t *local; const char *name; const char *args; } script_t;
static int safe_args(const char *s) {
    const unsigned char *p;
    for (p = (const unsigned char *)s; *p; p++)
        if (!isalnum(*p) && !strchr("/_-. ", *p)) return 0;
    return 1;
}
static int script_job(void *ctx, ulab_error_t *err) {
    script_t *job = ctx;
    char path[ULAB_MAX_PATH];
    char log[ULAB_MAX_PATH];
    char args[ULAB_MAX_ARGS];
    char *argv[32];
    char *save = NULL;
    char *part;
    int fd;
    int status;
    size_t n = 1;
    pid_t pid;
    if (webapp_path(path, sizeof(path), job->local->runtime.script_dir, job->name, err)) return ULAB_ERR;
    if (snprintf(log, sizeof(log), "%s/runtime-%ld-%s.log", job->local->run_dir, (long)getpid(), job->name) >= (int)sizeof(log))
        return webapp_error(err, "runtime log path too long");
    if (!safe_args(job->args) || ulab_copy(args, sizeof(args), job->args)) return webapp_error(err, "invalid runtime arguments");
    argv[0] = path;
    for (part = strtok_r(args, " ", &save); part; part = strtok_r(NULL, " ", &save)) {
        if (n >= 31) return webapp_error(err, "too many runtime arguments");
        argv[n++] = part;
    }
    argv[n] = NULL;
    fd = open(log, O_WRONLY | O_CREAT | O_APPEND | O_CLOEXEC, 0600);
    if (fd < 0) return webapp_error(err, "cannot open runtime log");
    pid = fork();
    if (!pid) {
        dup2(fd, STDOUT_FILENO); dup2(fd, STDERR_FILENO); close(fd);
        setenv("ULAB_WEBAPP_RUN", "1", 1);
        if (job->local->opts->factory_url[0]) setenv("FACTORY_URL", job->local->opts->factory_url, 1);
        execv(path, argv); _exit(127);
    }
    close(fd);
    if (pid < 0) return webapp_error(err, "cannot launch runtime script");
    while (waitpid(pid, &status, 0) < 0) if (errno != EINTR) return webapp_error(err, "cannot wait for runtime script");
    if (!WIFEXITED(status) || WEXITSTATUS(status)) return webapp_error(err, "runtime script %s failed; see private runtime log", job->name);
    return ULAB_OK;
}
static int execute_script(void *ctx, const char *name, const char *args, ulab_error_t *err) {
    local_t *local = ctx;
    script_t job = {local, name, args};
    double end;
    if (local->child) return script_job(&job, err); /* already in a bounded process group */
    end = webapp_now() + 900;
    if (end > local->client->deadline) end = local->client->deadline;
    return webapp_bounded_job(script_job, &job, end, local->client->cancel, err);
}
static int recover(void *ctx, webapp_journal_t *journal, ulab_error_t *err) {
    local_t *local = ctx;
    json_t *intents;
    json_t *intent;
    json_t *receipt;
    json_t *bindings;
    json_t *b;
    size_t i;
    size_t n;
    int failed = 0;
    char path[ULAB_MAX_PATH];
    char suffix[200];
    const char *state;
    if (!journal->root) return ULAB_OK;
    intents = json_object_get(journal->root, "creation_intents");
    json_array_foreach(intents, i, intent) {
        snprintf(suffix, sizeof(suffix), "browser/%.80s/creation-%lld.json", local->world->run_id,
                 (long long)json_integer_value(json_object_get(intent, "command_id")));
        if (webapp_path(path, sizeof(path), local->run_dir, suffix, err)) return ULAB_ERR;
        receipt = json_load_file(path, JSON_REJECT_DUPLICATES, NULL);
        state = receipt ? json_string_value(json_object_get(receipt, "state")) : NULL;
        bindings = receipt ? json_object_get(receipt, "bindings") : NULL;
        b = json_array_get(bindings, 0);
        if (!receipt || !json_equal(json_object_get(receipt, "command_id"), json_object_get(intent, "command_id")) ||
            !json_equal(json_object_get(receipt, "kind"), json_object_get(intent, "kind")) ||
            !json_equal(json_object_get(receipt, "ref"), json_object_get(intent, "ref")) ||
            !json_equal(json_object_get(receipt, "name"), json_object_get(intent, "name"))) {
            local->uncertain = failed = 1;
        } else if (ulab_streq(state, "identified") && json_array_size(bindings) == 1 &&
                   json_equal(json_object_get(b, "kind"), json_object_get(intent, "kind")) &&
                   json_equal(json_object_get(b, "ref"), json_object_get(intent, "ref")) &&
                   !webapp_journal_bind(journal, bindings, 1, (unsigned int)json_integer_value(json_object_get(intent, "command_id")), err)) {
            json_object_set_new(intent, "state", json_string("identified"));
        } else if (ulab_streq(state, "prepared") && json_is_array(bindings) && !json_array_size(bindings)) {
            json_object_set_new(intent, "state", json_string("not_submitted"));
        } else { local->uncertain = failed = 1; }
        json_decref(receipt);
    }
    if (local->started) {
        if (runtime_load_workload_sites(&local->runtime, local->world, err)) failed = 1;
        for (n = 0; n < local->world->node_count; n++) {
            node_t *node = &local->world->nodes[n];
            if (!node->bff_id[0]) continue;
            bindings = json_pack("[{s:s,s:s,s:s,s:s,s:s}]", "kind", "node", "ref", node->ref,
                                 "id", node->bff_id, "name", node->name, "observed_via", "runtime_claim");
            if (!bindings || webapp_journal_bind(journal, bindings, 2, 0, err)) failed = 1;
            json_decref(bindings);
        }
    }
    json_object_set_new(journal->root, "uncertain_creation", json_boolean(local->uncertain));
    if (webapp_journal_save(journal, err)) failed = 1;
    if (failed && !err->msg[0]) webapp_error(err, "creation ownership is unresolved; retain journal and investigate planned names before manual cleanup");
    return failed ? ULAB_ERR : ULAB_OK;
}
static int create(local_t *local, webapp_journal_t *journal, const char *kind, json_t *inputs, ulab_error_t *err) {
    json_t *reply = NULL;
    json_t *intent;
    json_t *intents;
    int rc;
    const char *action = !strcmp(kind, "network") ? "web_create_network" : "web_create_site";
    intents = json_object_get(journal->root, "creation_intents");
    intent = json_pack("{s:i,s:s,s:O,s:O,s:s}", "command_id", local->client->sequence + 1, "kind", kind,
                       "ref", json_object_get(inputs, "ref"), "name", json_object_get(inputs, "name"), "state", "pending");
    if (!intent || json_array_append_new(intents, intent) || webapp_journal_save(journal, err)) { json_decref(inputs); return ULAB_ERR; }
    rc = webapp_call(local->client, action, inputs, 900, &reply, err);
    if (!rc) {
        json_t *bindings = json_object_get(reply, "bindings");
        json_t *binding = json_array_get(bindings, 0);
        if (json_array_size(bindings) != 1 ||
            !json_equal(json_object_get(binding, "kind"), json_object_get(intent, "kind")) ||
            !json_equal(json_object_get(binding, "ref"), json_object_get(intent, "ref")) ||
            !json_equal(json_object_get(binding, "name"), json_object_get(intent, "name")) ||
            !json_is_true(json_object_get(json_object_get(reply, "actual"), "visible")))
            rc = webapp_error(err, "creation response differs from the planned entity or lacks visible acceptance");
        else rc = webapp_journal_bind(journal, bindings, 1, local->client->sequence, err);
    }
    json_decref(inputs); json_decref(reply);
    return rc;
}
static int prepare_sims(void *ctx, ulab_error_t *err) {
    local_t *local = ctx;
    char csv[ULAB_MAX_PATH];
    return sim_factory_prepare_world(local->opts, local->world, local->run_dir, csv, sizeof(csv), err);
}
static int import_sims(local_t *local, webapp_journal_t *journal, ulab_error_t *err) {
    char csv[ULAB_MAX_PATH];
    json_t *inputs;
    json_t *iccids;
    json_t *reply = NULL;
    size_t i;
    double end = webapp_now() + 900;
    int rc;
    if (!local->world->ue_count) return ULAB_OK;
    if (end > local->client->deadline) end = local->client->deadline;
    if (webapp_bounded_job(prepare_sims, local, end, local->client->cancel, err) ||
        webapp_path(csv, sizeof(csv), local->run_dir, "factory-sims.csv", err) ||
        sim_factory_load_world_csv(local->world, csv, err)) return ULAB_ERR;
    iccids = json_array();
    for (i = 0; i < local->world->ue_count; i++) json_array_append_new(iccids, json_string(local->world->ues[i].iccid));
    json_object_set_new(journal->root, "sim_inventory", json_pack("{s:O,s:s,s:s}", "iccids", iccids,
        "state", "import_pending", "cleanup", "retained_factory_pool"));
    if (webapp_journal_save(journal, err)) { json_decref(iccids); return ULAB_ERR; }
    inputs = json_pack("{s:s,s:o}", "csv_path", csv, "iccids", iccids);
    rc = webapp_call(local->client, "web_import_sims", inputs, 900, &reply, err);
    if (!rc && (!json_is_true(json_object_get(json_object_get(reply, "actual"), "executed")) ||
        json_array_size(json_object_get(reply, "bindings")))) rc = webapp_error(err, "invalid SIM import acknowledgement");
    json_object_set_new(json_object_get(journal->root, "sim_inventory"), "state", json_string(rc ? "uncertain" : "visible"));
    if (webapp_journal_save(journal, err)) rc = ULAB_ERR;
    json_decref(inputs); json_decref(reply); return rc;
}
static int provision(void *ctx, webapp_client_t *client, webapp_journal_t *journal, ulab_error_t *err) {
    local_t *local = ctx;
    world_t *world = local->world;
    const webapp_spec_t *profile = &local->scenario->webapp;
    size_t i;
    size_t n;
    json_t *inputs;
    json_t *components;
    selector_result_t selected = {0};
    int rc;
    local->client = client;
    json_object_set_new(journal->root, "creation_intents", json_array());
    for (i = 0; i < world->network_count; i++) {
        inputs = json_pack("{s:s,s:s}", "ref", world->networks[i].ref, "name", world->networks[i].name);
        if (!inputs || create(local, journal, "network", inputs, err)) return ULAB_ERR;
    }
    if (import_sims(local, journal, err)) return ULAB_ERR;
    if (!world->site_count) return ULAB_OK;
    local->started = 1;
    if (runtime_ensure_network(&local->runtime, err)) return ULAB_ERR;
    for (i = 0; i < world->site_count; i++) {
        site_t *site = &world->sites[i];
        network_t *network = world_network_by_ref(world, site->network_ref);
        if (runtime_start_selected_site(&local->runtime, world, i, err)) return ULAB_ERR;
        selected.idx = calloc(world->node_count, sizeof(size_t));
        if (!selected.idx) return webapp_error(err, "cannot allocate runtime selection");
        selected.count = 0;
        for (n = 0; n < world->node_count; n++) if (!strcmp(world->nodes[n].site_ref, site->ref)) selected.idx[selected.count++] = n;
        rc = runtime_wait_nodes_ready(&local->runtime, world, &selected, err);
        selector_result_free(&selected);
        if (rc) return rc;
        components = json_pack("{s:s,s:s,s:s}", "switch", profile->switch_component, "backhaul", profile->backhaul_component, "power", profile->power_component);
        inputs = json_pack("{s:s,s:s,s:s,s:s,s:s,s:o}", "ref", site->ref, "name", site->name, "network_name", network->name,
                           "network_id", network->bff_id, "tower_id", site->tnode_id, "components", components);
        if (!inputs || create(local, journal, "site", inputs, err)) return ULAB_ERR;
    }
    return recover(local, journal, err);
}
static int runtime_event(void *ctx, const event_spec_t *event, ulab_error_t *err) {
    local_t *local = ctx;
    selector_result_t nodes = {0};
    int rc;
    local->child = 1;
    if (event->type == EVT_START_UES || event->type == EVT_TRAFFIC) {
        size_t i;
        if (selector_resolve_ues(local->world, &event->ues, &nodes, err)) return ULAB_ERR;
        for (i = 0; i < nodes.count; i++) {
            ue_t *ue = &local->world->ues[nodes.idx[i]];
            if (!ue->bff_id[0] || !ue->site_ref[0] || (event->type == EVT_TRAFFIC && !ue->started)) { selector_result_free(&nodes); return webapp_error(err, "UE runtime requires an owned, UI-allocated site SIM"); }
        }
        if (event->type == EVT_START_UES) {
            for (i = 0; i < nodes.count; i++) if (sim_factory_wait_asr(local->opts, &local->world->ues[nodes.idx[i]], err)) { selector_result_free(&nodes); return ULAB_ERR; }
            rc = runtime_ensure_media(&local->runtime, err);
            if (!rc) rc = runtime_build_and_start_ues(local->opts->repo, &local->runtime, local->world, &nodes, err);
            if (!rc) rc = runtime_wait_ues_attached(&local->runtime, local->world, &nodes, err);
        } else rc = runtime_generate_traffic(&local->runtime, local->world, &nodes, event->amount_mb, err);
        selector_result_free(&nodes); return rc;
    }
    if (selector_resolve_nodes(local->world, &event->nodes, &nodes, err)) return ULAB_ERR;
    rc = event->type == EVT_DISCONNECT_NODES ? runtime_disconnect_nodes(&local->runtime, local->world, &nodes, err) :
         event->type == EVT_RECONNECT_NODES ? runtime_reconnect_nodes(&local->runtime, local->world, &nodes, err) :
         webapp_error(err, "unsupported webapp runtime event");
    selector_result_free(&nodes); return rc;
}
static int reconnect(void *ctx, ulab_error_t *err) {
    local_t *local = ctx;
    selector_result_t nodes = {0};
    size_t i;
    int rc;
    if (!local->started) return ULAB_OK;
    local->child = 1;
    nodes.idx = calloc(local->world->node_count, sizeof(size_t));
    if (!nodes.idx) return webapp_error(err, "cannot allocate cleanup selection");
    for (i = 0; i < local->world->node_count; i++) if (local->world->nodes[i].bff_id[0]) nodes.idx[nodes.count++] = i;
    rc = runtime_reconnect_nodes(&local->runtime, local->world, &nodes, err);
    selector_result_free(&nodes); return rc;
}
static int cleanup_runtime(void *ctx, ulab_error_t *err) {
    local_t *local = ctx;
    if (!local->started) return ULAB_OK;
    local->child = 1;
    if (local->world->ue_count && runtime_cleanup_ues(&local->runtime, local->world, err)) return ULAB_ERR;
    return runtime_cleanup_infra(&local->runtime, local->world, err);
}
static int cleanup_resource(void *ctx, const char *kind, const char *id, ulab_error_t *err) {
    local_t *local = ctx;
    bff_client_t bff;
    char path[ULAB_MAX_PATH];
    char suffix[ULAB_MAX_ID + 40];
    int rc;
    int linked = 0;
    size_t i;
    if (local->uncertain) return webapp_error(err, "unresolved creation: retain resources for reconciliation using the journal");
    snprintf(suffix, sizeof(suffix), "cleanup-%s-%s", kind, id);
    if (webapp_path(path, sizeof(path), local->run_dir, suffix, err) || mkdir(path, 0700)) return webapp_error(err, "cannot open cleanup log directory");
    for (i = 0; i < local->world->node_count; i++) if (!strcmp(local->world->nodes[i].bff_id, id)) {
        site_t *site = world_site_by_ref(local->world, local->world->nodes[i].site_ref);
        linked = site && site->bff_id[0];
    }
    rc = bff_init(&bff, local->opts->bff_url, path);
    if (!rc) rc = bff_cleanup_resource(&bff, kind, id, linked, err);
    else webapp_error(err, "cannot initialize BFF teardown authentication");
    bff_close(&bff); return rc;
}
int webapp_run_local(const runner_opts_t *opts, const scenario_t *scenario,
                     world_t *world, report_t *report, const char *run_dir, ulab_error_t *err) {
    local_t local;
    webapp_hooks_t hooks;
    int rc;
    if (!world->network_count) return webapp_execute(opts, scenario, world, report, run_dir, NULL, err);
    /* Existing runtime APIs serialize arguments as whitespace-delimited strings. */
    if (!safe_args(opts->repo) || strchr(opts->repo, ' ') || !safe_args(run_dir) || strchr(run_dir, ' '))
        return webapp_error(err, "local runtime repo/run paths must use letters, numbers, slash, dot, dash or underscore");
    memset(&local, 0, sizeof(local)); memset(&hooks, 0, sizeof(hooks));
    local.opts = opts; local.scenario = scenario; local.world = world; local.run_dir = run_dir;
    runtime_init(&local.runtime, "virtual", opts->script_dir, run_dir, opts->repo);
    local.runtime.execute = execute_script; local.runtime.execute_ctx = &local;
    hooks.ctx = &local; hooks.provision = provision; hooks.recover = recover; hooks.prepare_cleanup = reconnect;
    hooks.runtime_event = runtime_event; hooks.cleanup_resource = cleanup_resource; hooks.cleanup_runtime = cleanup_runtime;
    rc = webapp_execute(opts, scenario, world, report, run_dir, &hooks, err);
    runtime_close(&local.runtime); return rc;
}
