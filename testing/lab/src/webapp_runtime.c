/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "runtime.h"
#include "log.h"
#include "bff.h"
#include "sim_factory.h"
#include "util.h"
#include <ctype.h>
#include <errno.h>
#include <fcntl.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
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
    int teardown;
    bff_client_t bff;
    int bff_open;
    int replay;
    json_t *resources;
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
/* These reads are used only to establish ownership and reconcile teardown.
 * No backend result can satisfy a browser event or assertion. */
static const char *field(json_t *obj, const char *key) {
    return json_string_value(json_object_get(obj, key));
}

static int open_cleanup_bff(local_t *local, ulab_error_t *err) {
    char path[ULAB_MAX_PATH];
    if (local->bff_open) return ULAB_OK;
    if (webapp_path(path, sizeof(path), local->run_dir, "cleanup-bff-XXXXXX", err) || !mkdtemp(path))
        return webapp_error(err, "cannot open cleanup BFF log directory");
    if (bff_init(&local->bff, local->opts->bff_url, path)) {
        bff_close(&local->bff);
        return webapp_error(err, "cannot initialize BFF teardown authentication");
    }
    local->bff_open = 1;
    return ULAB_OK;
}

static int inventory(local_t *local, const char *target, network_t *network,
                     json_t **rows, ulab_error_t *err) {
    if (open_cleanup_bff(local, err)) return ULAB_ERR;
    return bff_cleanup_list(&local->bff, target, network, local->opts->sim_type, rows, err);
}

/* Match every supplied field. A duplicate identity or malformed ID is an
 * error, even if the first entry would have been a plausible candidate. */
static int lookup_id(json_t *rows, json_t *match, const char *key, const char **id,
                     ulab_error_t *err) {
    size_t i;
    json_t *row;
    json_t *expected;
    const char *field_name;
    *id = NULL;
    if (!match || !json_object_size(match)) return webapp_error(err, "empty cleanup identity");
    json_array_foreach(rows, i, row) {
        json_object_foreach(match, field_name, expected) {
            if (!json_object_get(row, field_name))
                return webapp_error(err, "backend cleanup inventory is missing %s", field_name);
            if (!json_equal(json_object_get(row, field_name), expected)) break;
        }
        if (field_name) continue;
        if (*id) return webapp_error(err, "multiple backend rows match the planned resource");
        *id = field(row, key);
        if (!*id || !**id) return webapp_error(err, "backend resource has no cleanup ID");
    }
    return ULAB_OK;
}

static int check_absent(json_t *rows, const char *name, const char *key, ulab_error_t *err) {
    const char *id;
    json_t *match = json_pack("{s:s}", "name", name);
    int rc = lookup_id(rows, match, key, &id, err);
    json_decref(match);
    if (!rc && id) return webapp_error(err, "planned resource already exists: %s", name);
    return rc;
}

static int check_global_names(local_t *local, webapp_journal_t *journal, ulab_error_t *err) {
    world_t *w = local->world;
    json_t *rows = NULL;
    size_t i;
    int rc = inventory(local, "networks", NULL, &rows, err);
    for (i = 0; !rc && i < w->network_count; i++)
        rc = check_absent(rows, w->networks[i].name, "id", err);
    json_decref(rows);
    rows = NULL;
    for (i = 0; !rc && i < w->package_count; i++) {
        if (w->packages[i].network_ref[0]) continue;
        if (!rows) rc = inventory(local, "packages", NULL, &rows, err);
        if (!rc) rc = check_absent(rows, w->packages[i].name, "uuid", err);
    }
    json_decref(rows);
    if (rc) return rc;
    json_object_set_new(journal->root, "global_names_absent", json_true());
    return webapp_journal_save(journal, err);
}

static int bind_lookup(webapp_journal_t *j, const char *kind, const char *ref,
                       const char *name, const char *id, unsigned int command, ulab_error_t *err) {
    json_t *bindings = json_pack("[{s:s,s:s,s:s,s:s,s:s}]", "kind", kind, "ref", ref,
                                 "name", name, "id", id, "observed_via", "cleanup_lookup");
    int rc = webapp_journal_bind(j, bindings, 3, command, err);
    json_decref(bindings);
    return rc;
}

static int resolve_submitted(local_t *local, webapp_journal_t *j, json_t *intent,
                             ulab_error_t *err) {
    world_t *w = local->world;
    const char *kind = field(intent, "kind");
    const char *ref = field(intent, "ref");
    const char *name = field(intent, "name");
    const char *target;
    const char *id_key = "id";
    const char *network_ref = NULL;
    const char *id = NULL;
    network_t *net = NULL;
    json_t *match;
    json_t *rows = NULL;
    int rc;
    int attempt;

    if (!kind || !ref || !name) return webapp_error(err, "invalid cleanup intent");
    if (ulab_streq(kind, "network") && world_network_by_ref(w, ref)) {
        target = "networks";
        match = json_pack("{s:s}", "name", name);
    } else if (ulab_streq(kind, "site") && world_site_by_ref(w, ref)) {
        target = "sites";
        network_ref = world_site_by_ref(w, ref)->network_ref;
        match = json_pack("{s:s}", "name", name);
    } else if (ulab_streq(kind, "package") && world_package_by_ref(w, ref)) {
        package_t *p = world_package_by_ref(w, ref);
        target = "packages";
        id_key = "uuid";
        if (p->network_ref[0]) network_ref = p->network_ref;
        match = json_pack("{s:s}", "name", name);
    } else if (ulab_streq(kind, "subscriber") && world_subscriber_by_ref(w, ref)) {
        subscriber_t *s = world_subscriber_by_ref(w, ref);
        target = "subscribers";
        id_key = "uuid";
        network_ref = s->network_ref;
        match = json_pack("{s:s,s:s}", "name", name, "email", s->email);
    } else if (ulab_streq(kind, "sim") && world_ue_by_ref(w, ref)) {
        ue_t *ue = world_ue_by_ref(w, ref);
        subscriber_t *s = world_subscriber_by_ref(w, ue->subscriber_ref);
        if (!s || !s->bff_id[0]) return webapp_error(err, "unresolved SIM subscriber");
        target = "sims";
        network_ref = ue->network_ref;
        /* Auto-assignment lets the backend choose the ICCID. Its allocation
         * belongs to our subscriber; its original pool entry may be shared. */
        match = json_pack("{s:s}", "subscriberId", s->bff_id);
        if (!json_is_true(json_object_get(intent, "auto_assign")))
            json_object_set_new(match, "iccid", json_string(name));
    } else {
        return webapp_error(err, "unresolved %s creation for %s", kind, ref);
    }
    if (network_ref) net = world_network_by_ref(w, network_ref);
    if (network_ref && net && net->bff_id[0]) {
        json_object_set_new(match, "networkId", json_string(net->bff_id));
    } else if (network_ref || !json_is_true(json_object_get(j->root, "global_names_absent"))) {
        json_decref(match);
        return webapp_error(err, "no ownership scope for %s %s", kind, ref);
    }
    /* A lost response can precede visibility in the inventory. */
    for (attempt = 0; attempt < 3; attempt++) {
        if (attempt) sleep(1);
        rc = inventory(local, target, net, &rows, err);
        if (rc) break;
        rc = lookup_id(rows, match, id_key, &id, err);
        if (!rc && id)
            rc = bind_lookup(j, kind, ref, name, id, json_integer_value(json_object_get(intent, "command_id")), err);
        if (rc || id) break;
        json_decref(rows);
        rows = NULL;
    }
    if (!rc && !id && !json_is_true(json_object_get(intent, "validation_rejected")))
        rc = webapp_error(err, "submitted %s %s has no confirmed outcome", kind, ref);
    if (!rc) json_object_set_new(intent, "state", json_string(id ? "reconciled" : "not_found"));
    json_decref(rows);
    json_decref(match);
    return rc;
}

static int pool_inventory(local_t *local, webapp_journal_t *j, int cleanup, ulab_error_t *err) {
    json_t *rows = NULL;
    json_t *sim_inventory = json_object_get(j->root, "sim_inventory");
    size_t i;
    int failed = 0;

    if (cleanup && local->replay && !json_is_true(json_object_get(sim_inventory, "absent_before_import")))
        return webapp_error(err, "legacy SIM pool journal lacks a pre-import ownership check");
    if (inventory(local, "pool", NULL, &rows, err)) return ULAB_ERR;
    for (i = 0; i < local->world->ue_count; i++) {
        ue_t *ue = &local->world->ues[i];
        json_t *match = json_pack("{s:s}", "iccid", ue->iccid);
        const char *id;
        int rc = lookup_id(rows, match, "id", &id, err);
        json_decref(match);
        if (!rc && id && !cleanup) rc = webapp_error(err, "factory ICCID already in the SIM pool: %s", ue->iccid);
        if (!rc && !id && cleanup && !ue->pool_sim_id[0])
            rc = webapp_error(err, "could not confirm SIM pool import outcome for %s", ue->iccid);
        if (!rc && id && cleanup) rc = bind_lookup(j, "pool_sim", ue->ref, ue->iccid, id, 0, err);
        if (rc) failed = 1;
    }
    json_decref(rows);
    if (cleanup) json_object_set_new(sim_inventory, "cleanup", json_string(failed ? "failed" : "pending"));
    return failed ? ULAB_ERR : ULAB_OK;
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
        if (local->replay && (ulab_streq(field(intent, "state"), "reconciled") ||
            ulab_streq(field(intent, "state"), "identified") || ulab_streq(field(intent, "state"), "not_submitted"))) continue;
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
            failed = 1;
        } else if (ulab_streq(state, "identified") && json_array_size(bindings) == 1 &&
                   json_equal(json_object_get(b, "kind"), json_object_get(intent, "kind")) &&
                   json_equal(json_object_get(b, "ref"), json_object_get(intent, "ref")) &&
                   !webapp_journal_bind(journal, bindings, 1, (unsigned int)json_integer_value(json_object_get(intent, "command_id")), err)) {
            json_object_set_new(intent, "state", json_string("identified"));
        } else if ((ulab_streq(state, "prepared") || ulab_streq(state,"intercepted")) && json_is_array(bindings) && !json_array_size(bindings)) {
            json_object_set_new(intent, "state", json_string("not_submitted"));
        } else if (local->teardown && ulab_streq(state, "submitted") &&
                   json_is_array(bindings) && !json_array_size(bindings)) {
            if (resolve_submitted(local, journal, intent, err)) failed = 1;
        } else { failed = 1; }
        if (failed && !err->msg[0]) webapp_error(err, "could not reconcile %s %s (command %lld)",
            field(intent, "kind"), field(intent, "ref"), (long long)json_integer_value(json_object_get(intent, "command_id")));
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
    if (local->teardown && json_object_get(journal->root, "sim_inventory") &&
        pool_inventory(local, journal, 1, err)) {
        json_object_set_new(json_object_get(journal->root, "sim_inventory"), "cleanup", json_string("failed"));
        failed = 1;
    }
    json_object_set_new(journal->root, "uncertain_creation", json_boolean(failed));
    if (webapp_journal_save(journal, err)) failed = 1;
    if (failed && !err->msg[0]) webapp_error(err, "creation ownership is unresolved; retain journal and investigate planned names before manual cleanup");
    return failed ? ULAB_ERR : ULAB_OK;
}
static int recover_cleanup(void *ctx, webapp_journal_t *journal, ulab_error_t *err) {
    local_t *local = ctx;
    int rc;
    /* A long browser run can outlive its original BFF session. */
    if (local->bff_open) { bff_close(&local->bff); local->bff_open = 0; }
    local->teardown = 1;
    local->resources = json_object_get(journal->root, "resources");
    rc = recover(local, journal, err);
    if (json_array_size(json_object_get(journal->root, "resources")) && open_cleanup_bff(local, err))
        rc = ULAB_ERR;
    return rc;
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
    ulab_status("WEBAPP", "%s %s", action, json_string_value(json_object_get(inputs, "name")));
    rc = webapp_call(local->client, action, inputs, 900, &reply, err);
    if (rc && !local->client->broken &&
        ulab_streq(field(json_object_get(reply, "error"), "code"), "SITE_CREATE_FAILED") &&
        ulab_streq(field(json_object_get(reply, "error"), "message"), "Argument Validation Error"))
        json_object_set_new(intent, "validation_rejected", json_true());
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
    if (pool_inventory(local, journal, 0, err)) return ULAB_ERR;
    iccids = json_array();
    for (i = 0; i < local->world->ue_count; i++) json_array_append_new(iccids, json_string(local->world->ues[i].iccid));
    json_object_set_new(journal->root, "sim_inventory", json_pack("{s:O,s:s,s:s}", "iccids", iccids,
        "state", "import_pending", "cleanup", "pending"));
    json_object_set_new(json_object_get(journal->root, "sim_inventory"), "absent_before_import", json_true());
    if (webapp_journal_save(journal, err)) { json_decref(iccids); return ULAB_ERR; }
    inputs = json_pack("{s:s,s:o}", "csv_path", csv, "iccids", iccids);
    rc = webapp_call(local->client, "web_import_sims", inputs, 900, &reply, err);
    if (!rc && (!json_is_true(json_object_get(json_object_get(reply, "actual"), "executed")) ||
        json_array_size(json_object_get(reply, "bindings")))) rc = webapp_error(err, "invalid SIM import acknowledgement");
    json_object_set_new(json_object_get(journal->root, "sim_inventory"), "state", json_string(rc ? "uncertain" : "visible"));
    if (webapp_journal_save(journal, err)) rc = ULAB_ERR;
    json_decref(inputs); json_decref(reply); return rc;
}
static int wait_site_anchors(void *ctx, ulab_error_t *err) {
    local_t *local = ctx;
    bff_client_t bff;
    size_t i;
    int rc;

    rc = bff_init(&bff, local->opts->bff_url, local->run_dir);
    if (rc) webapp_error(err, "cannot initialize BFF for site readiness");
    for (i = 0; !rc && i < local->world->site_count; i++) {
        site_t *site = &local->world->sites[i];
        ulab_status("BACKEND", "wait for site %s tower %s in registry", site->ref, site->tnode_id);
        rc = bff_wait_site_anchor_online(&bff, site, err);
    }
    bff_close(&bff);
    return rc;
}

static int start_nodes(local_t *local, ulab_error_t *err) {
    selector_t all = { .kind = SEL_ALL };
    selector_result_t nodes = {0};
    int rc;

    if (!local->world->site_count) return ULAB_OK;
    local->started = 1;
    if (runtime_ensure_network(&local->runtime, err) ||
        runtime_build_and_start_sites(local->opts->repo, &local->runtime, local->world, err))
        return ULAB_ERR;
    rc = selector_resolve_nodes(local->world, &all, &nodes, err);
    if (!rc) rc = runtime_wait_nodes_ready(&local->runtime, local->world, &nodes, err);
    selector_result_free(&nodes);
    return rc;
}

static int provision(void *ctx, webapp_client_t *client, webapp_journal_t *journal, ulab_error_t *err) {
    local_t *local = ctx;
    world_t *world = local->world;
    const webapp_spec_t *profile = &local->scenario->webapp;
    size_t i;
    json_t *inputs;
    json_t *components;
    local->client = client;
    json_object_set_new(journal->root, "bff_url", json_string(local->opts->bff_url));
    if (check_global_names(local, journal, err)) return ULAB_ERR;
    json_object_set_new(journal->root, "creation_intents", json_array());
    if (start_nodes(local, err)) return ULAB_ERR;
    if (world->site_count && !ulab_streq(profile->session_mode, "onboarding") &&
        webapp_bounded_job(wait_site_anchors, local, client->deadline, client->cancel, err)) return ULAB_ERR;
    for (i = 0; i < world->network_count && !ulab_streq(profile->session_mode,"onboarding"); i++) {
        inputs = json_pack("{s:s,s:s}", "ref", world->networks[i].ref, "name", world->networks[i].name);
        if (!inputs || create(local, journal, "network", inputs, err)) return ULAB_ERR;
    }
    if (import_sims(local, journal, err)) return ULAB_ERR;
    if (!world->site_count) return ULAB_OK;
    for (i = 0; i < world->site_count && !ulab_streq(profile->session_mode,"onboarding"); i++) {
        site_t *site = &world->sites[i];
        network_t *network = world_network_by_ref(world, site->network_ref);
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
    int rc;
    if (!local->started) return ULAB_OK;
    local->child = 1;
    rc = runtime_cleanup_ues(&local->runtime, local->world, err);
    if (runtime_cleanup_infra(&local->runtime, local->world, err)) rc = ULAB_ERR;
    return rc;
}
/* Read by journaled ID so teardown can verify a delete or safely replay it. */
static int recorded_present(local_t *local, const char *kind, const char *id,
                            int *present, ulab_error_t *err) {
    const char *target;
    const char *key = "id";
    const char *ref = NULL;
    const char *network_ref = NULL;
    network_t *network = NULL;
    json_t *r;
    json_t *rows = NULL;
    size_t i;
    int rc;

    *present = 0;
    json_array_foreach(local->resources, i, r) {
        if (ulab_streq(field(r, "kind"), kind) && ulab_streq(field(r, "id"), id)) {
            ref = field(r, "ref");
            break;
        }
    }
    if (!ref) return webapp_error(err, "missing journal reference for %s", id);
    if (ulab_streq(kind, "network")) target = "networks";
    else if (ulab_streq(kind, "package")) {
        package_t *p = world_package_by_ref(local->world, ref);
        target = "packages";
        key = "uuid";
        if (p->network_ref[0]) network_ref = p->network_ref;
    } else if (ulab_streq(kind, "pool_sim")) target = "pool";
    else if (ulab_streq(kind, "site")) {
        target = "sites";
        network_ref = world_site_by_ref(local->world, ref)->network_ref;
    } else if (ulab_streq(kind, "subscriber")) {
        target = "subscribers";
        key = "uuid";
        network_ref = world_subscriber_by_ref(local->world, ref)->network_ref;
    } else if (ulab_streq(kind, "sim")) {
        target = "sims";
        network_ref = world_ue_by_ref(local->world, ref)->network_ref;
    } else return webapp_error(err, "unknown cleanup resource %s", kind);
    if (network_ref) network = world_network_by_ref(local->world, network_ref);
    rc = inventory(local, target, network, &rows, err);
    if (rc) return rc;
    json_array_foreach(rows, i, r) {
        if (!field(r, key)) { rc = webapp_error(err, "missing ID in %s", target); break; }
        if (ulab_streq(field(r, key), id)) { *present = 1; break; }
    }
    json_decref(rows);
    return rc;
}

static int cleanup_resource(void *ctx, const char *kind, const char *id, ulab_error_t *err) {
    local_t *local = ctx;
    int rc;
    int present;
    int assigned;
    int attempt;
    size_t i;
    const char *site_id = NULL;

    if (open_cleanup_bff(local, err)) return ULAB_ERR;
    if (local->replay && !ulab_streq(kind, "node")) {
        rc = recorded_present(local, kind, id, &present, err);
        if (rc || !present) return rc;
    }
    if (ulab_streq(kind, "node")) {
        for (i = 0; i < local->world->node_count; i++) {
            node_t *node = &local->world->nodes[i];
            if (ulab_streq(node->bff_id, id)) {
                site_t *site = world_site_by_ref(local->world, node->site_ref);
                if (site) site_id = site->bff_id;
                break;
            }
        }
    }
    if (bff_cleanup_resource(&local->bff, kind, id, site_id, err)) return ULAB_ERR;
    for (attempt = 0; attempt < 3; attempt++) {
        if (attempt) sleep(1);
        if (ulab_streq(kind, "node"))
            rc = bff_cleanup_node_state(&local->bff, id, site_id, &present, &assigned, err);
        else rc = recorded_present(local, kind, id, &present, err);
        if (!rc && !present) return ULAB_OK;
    }
    if (rc) return rc;
    return webapp_error(err, "%s %s is still in backend inventory after deletion", kind, id);
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
    hooks.ctx = &local; hooks.provision = provision; hooks.recover = recover_cleanup; hooks.prepare_cleanup = reconnect;
    hooks.runtime_event = runtime_event; hooks.cleanup_resource = cleanup_resource; hooks.cleanup_runtime = cleanup_runtime;
    rc = webapp_execute(opts, scenario, world, report, run_dir, &hooks, err);
    if (local.bff_open) bff_close(&local.bff);
    runtime_close(&local.runtime); return rc;
}

/* Retry only backend teardown for a finished run. The scenario regenerates
 * world references; the journal and creation receipts supply backend IDs. */
int webapp_cleanup_backend(const runner_opts_t *opts, ulab_error_t *err) {
    scenario_t *scenario = NULL;
    world_t world = {0};
    webapp_journal_t j = {0};
    webapp_hooks_t hooks = {0};
    local_t local = {0};
    char run_dir[ULAB_MAX_PATH];
    char path[ULAB_MAX_PATH];
    json_t *report = NULL;
    json_t *pool;
    json_t *recorded;
    const char *endpoint;
    size_t i;
    int fd = -1;
    int rc = ULAB_ERR;

    if (!webapp_valid_run_id(opts->run_id)) return webapp_error(err, "cleanup requires the original --run-id");
    if (webapp_path(run_dir, sizeof(run_dir), opts->out_dir, opts->run_id, err) ||
        webapp_path(j.path, sizeof(j.path), run_dir, "webapp-resources.json", err) ||
        webapp_path(path, sizeof(path), run_dir, "report.json", err)) return ULAB_ERR;
    fd = open(run_dir, O_RDONLY | O_DIRECTORY | O_CLOEXEC);
    if (fd < 0) return webapp_error(err, "cannot open the original run directory");
    if (flock(fd, LOCK_EX | LOCK_NB)) {
        webapp_error(err, "another cleanup is active for this run");
        goto done;
    }
    j.root = json_load_file(j.path, JSON_REJECT_DUPLICATES, NULL);
    report = json_load_file(path, JSON_REJECT_DUPLICATES, NULL);
    scenario = calloc(1, sizeof(*scenario));
    if (!scenario) { webapp_error(err, "cannot allocate scenario"); goto done; }
    if (scenario_load(opts->scenario_path, scenario, err)) goto done;
    if (scenario->version != ULAB_WEBAPP_SCHEMA_VER || json_integer_value(json_object_get(j.root, "version")) != 1 ||
        !field(j.root, "run_result") || !json_is_array(json_object_get(j.root, "resources")) ||
        !ulab_streq(field(j.root, "run_id"), opts->run_id) || !ulab_streq(field(report, "run_id"), opts->run_id) ||
        !ulab_streq(field(report, "scenario"), scenario->name)) {
        webapp_error(err, "cleanup requires the matching scenario and a completed webapp run");
        goto done;
    }
    endpoint = field(j.root, "bff_url");
    if (endpoint && !ulab_streq(endpoint, opts->bff_url)) {
        webapp_error(err, "BFF endpoint differs from the recorded run");
        goto done;
    }
    if (ulab_streq(field(j.root, "cleanup"), "complete") && (!json_object_get(j.root, "sim_inventory") ||
        ulab_streq(field(json_object_get(j.root, "sim_inventory"), "cleanup"), "complete"))) {
        rc = ULAB_OK;
        goto done;
    }
    /* Version-one journals predate the validation_rejected flag. The final
     * provisioning error identifies the rejected site command in these runs. */
    if (!endpoint && ulab_streq(field(report, "error"), "web-app: Argument Validation Error")) {
        json_t *intents = json_object_get(j.root, "creation_intents");
        json_t *last = json_array_get(intents, json_array_size(intents) - 1);
        if (ulab_streq(field(last, "kind"), "site"))
            json_object_set_new(last, "validation_rejected", json_true());
    }
    if (opts->has_seed_override) scenario->seed = opts->seed_override;
    if (world_generate(scenario, opts->run_id, &world, err)) goto done;
    j.world = &world;
    pool = json_object_get(j.root, "sim_inventory");
    if (pool) {
        if (webapp_path(path, sizeof(path), run_dir, "factory-sims.csv", err) ||
            sim_factory_load_world_csv(&world, path, err)) goto done;
        recorded = json_object_get(pool, "iccids");
        for (i = 0; i < world.ue_count; i++)
            if (!ulab_streq(json_string_value(json_array_get(recorded, i)), world.ues[i].iccid)) break;
        if (i != world.ue_count || i != json_array_size(recorded)) {
            webapp_error(err, "factory batch differs from the journaled ICCIDs");
            goto done;
        }
    }
    if (webapp_journal_restore(&j, err)) goto done;
    local.opts = opts; local.world = &world; local.run_dir = run_dir;
    local.replay = 1;
    hooks.ctx = &local; hooks.cleanup_resource = cleanup_resource;
    ulab_status("CLEANUP", "retry backend cleanup for %s", opts->run_id);
    rc = recover_cleanup(&local, &j, err);
    if (webapp_journal_cleanup(&j, &hooks, 0, err)) rc = ULAB_ERR;
done:
    if (!rc) ulab_status("CLEANUP", "backend cleanup complete for %s", opts->run_id);
    if (local.bff_open) bff_close(&local.bff);
    webapp_journal_close(&j);
    json_decref(report);
    world_free(&world);
    free(scenario);
    close(fd);
    return rc;
}
