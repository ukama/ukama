/* SPDX-License-Identifier: MPL-2.0
 * Controlled infrastructure fixtures; no console-app acceptance credit.
 */
#include "webapp.h"
#include "util.h"
#include "workload.h"
#include <assert.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>

static const char *mode;
static int record(void *ctx, const char *kind, const char *id, ulab_error_t *err) {
    FILE *f = fopen(ctx, "a");
    assert(f);
    fprintf(f, "%s %s\n", kind, id); fclose(f);
    if (!strcmp(mode, "timeout")) {
        signal(SIGTERM, SIG_IGN);
        for (;;) pause();
    }
    return !strcmp(mode, "partial") && !strcmp(kind, "site") ? webapp_error(err, "fixture delete failure") : ULAB_OK;
}
static int runtime(void *ctx, const event_spec_t *event, ulab_error_t *err) {
    assert(event->type == EVT_DISCONNECT_NODES);
    if (record(ctx, "runtime", "disconnect", err)) return ULAB_ERR;
    return !strcmp(mode, "runtime-fail") ? webapp_error(err, "fixture runtime failure") : ULAB_OK;
}
static int final_cleanup(void *ctx, ulab_error_t *err) {
    int rc;
    rc = record(ctx, "runtime", "cleanup", err);
    if (!strcmp(mode, "runtime-cancel-cleanup")) kill(getppid(), SIGTERM);
    return rc;
}
static json_t *binding(const char *kind, const char *ref, const char *id,
                        const char *name, int created) {
    return json_pack("[{s:s,s:s,s:s,s:s,s:s}]", "kind", kind, "ref", ref,
                       "id", id, "name", name, "observed_via", created ? "ui_result" : "url");
}
int main(int argc, char **argv) {
    scenario_t *s;
    world_t world;
    webapp_journal_t journal;
    webapp_hooks_t hooks;
    runner_opts_t opts;
    report_t report;
    ulab_error_t err = {{0}};
    json_t *b;
    json_t *inputs = NULL;
    json_t *root;
    json_error_t jerr;
    char log[ULAB_MAX_PATH];
    char run_id[ULAB_MAX_ID];
    double start;
    int rc;
    int want_failure;
    event_spec_t event;
    check_spec_t check;
    if (argc < 4) return 2;
    memset(&check, 0, sizeof(check)); check.type = CHECK_WEB_FIELD_EQUALS;
    strcpy(check.expected, " \t10\xc2\xa0 GB\n");
    b = webapp_check_inputs(&check); assert(b);
    assert(!strcmp(json_string_value(json_object_get(b, "expected")), "10 GB")); json_decref(b);
    mode = argv[1];
    s = calloc(1, sizeof(*s)); assert(s);
    memset(&world, 0, sizeof(world)); memset(&hooks, 0, sizeof(hooks));
    assert(!webapp_path(log, sizeof(log), argv[2], "cleanup.log", &err));
    hooks.ctx = log; hooks.cleanup_resource = record;
    if (!strcmp(mode, "classify")) {
        assert(workload_is_file(argv[3]) == (argc > 4 ? atoi(argv[4]) : 0));
        free(s); return 0;
    }
    if (!strncmp(mode, "runtime-", 8)) {
        assert(argc == 5);
        assert(!scenario_load(argv[3], s, &err));
        /* Inject a controlled runtime step below the public Patch 4 gate to
         * exercise ordering and failure cleanup. It is not a runnable product scenario. */
        s->phases[1].events[1] = s->phases[1].events[0];
        memset(&s->phases[1].events[0], 0, sizeof(event_spec_t));
        s->phases[1].events[0].type = EVT_DISCONNECT_NODES;
        s->phases[1].event_count = 2;
        hooks.runtime_event = runtime; hooks.cleanup_runtime = final_cleanup;
        memset(&opts, 0, sizeof(opts)); ulab_copy(opts.webapp_worker, sizeof(opts.webapp_worker), argv[4]);
        assert(!world_generate(s, "fixture-runtime", &world, &err));
        assert(!report_open(&report, s, world.run_id, argv[2]));
        rc = webapp_execute(&opts, s, &world, &report, argv[2], &hooks, &err);
        report_set_final_rc(&report, rc); report_close(&report);
        assert((rc != 0) == (!strcmp(mode, "runtime-fail") || !strcmp(mode, "runtime-cancel-cleanup")));
        world_free(&world); free(s); return 0;
    }
    s->version = ULAB_WEBAPP_SCHEMA_VER; s->world.networks = 1;
    s->world.sites_per_network = 1; s->world.tower_per_site = 1;
    memset(run_id, 'a', 80); run_id[80] = '\0';
    assert(!world_generate(s, run_id, &world, &err));
    assert(strlen(world.networks[0].name) <= 40 && strlen(world.sites[0].name) <= 40);
    assert(!webapp_journal_open(&journal, &world, argv[2], &err));
    b = binding("network", "net-001", "network-id", world.networks[0].name, 1);
    assert(!webapp_journal_bind(&journal, b, 1, 1, &err)); json_decref(b);
    b = binding("site", "site-001", "site-id", world.sites[0].name, 1);
    assert(!webapp_journal_bind(&journal, b, 1, 2, &err));
    assert(!webapp_journal_bind(&journal, b, 1, 2, &err)); json_decref(b);
    assert(json_array_size(json_object_get(journal.root, "resources")) == 2);
    assert(!strcmp(world.sites[0].bff_id, "site-id"));
    b = binding("site", "site-001", "different-id", world.sites[0].name, 1);
    assert(webapp_journal_bind(&journal, b, 1, 3, &err)); json_decref(b);
    b = binding("site", "site-999", "foreign-id", "foreign", 1);
    assert(webapp_journal_bind(&journal, b, 1, 3, &err)); json_decref(b);
    b = binding("node", world.nodes[0].ref, world.nodes[0].id, world.nodes[0].name, 0);
    assert(!webapp_journal_bind(&journal, b, 0, 4, &err)); json_decref(b);
    assert(json_array_size(json_object_get(journal.root, "observations")) == 1);
    assert(!world.nodes[0].bff_id[0]);
    b = binding("network", "net-001", "network-id", world.networks[0].name, 0);
    assert(webapp_journal_bind(&journal, b, 1, 5, &err)); json_decref(b);
    memset(&event, 0, sizeof(event)); event.type = EVT_WEB_OPEN;
    strcpy(event.view, "network_site_detail"); event.networks.kind = SEL_REF;
    strcpy(event.networks.value, "net-001"); strcpy(event.sites.value, "site-001");
    assert(!webapp_event_inputs(&event, &world, &inputs, &err));
    assert(!strcmp(json_string_value(json_object_get(inputs, "network_name")), world.networks[0].name));
    assert(!strcmp(json_string_value(json_object_get(json_object_get(inputs, "entity"), "id")), "site-id"));
    json_decref(inputs);
    strcpy(world.sites[0].network_ref, "net-002");
    assert(webapp_event_inputs(&event, &world, &inputs, &err));
    assert(!inputs);
    start = webapp_now();
    rc = webapp_journal_cleanup(&journal, &hooks, start + (!strcmp(mode, "timeout") ? 0.15 : 5), &err);
    want_failure = !strcmp(mode, "partial") || !strcmp(mode, "timeout");
    assert((rc != 0) == want_failure);
    assert(webapp_now() - start < 4);
    root = json_load_file(journal.path, 0, &jerr); assert(root);
    assert(!strcmp(json_string_value(json_object_get(root, "cleanup")), want_failure ? "failed" : "complete"));
    json_decref(root);
    webapp_journal_close(&journal); world_free(&world);
    s->version = 1;
    assert(!world_generate(s, "legacy-run", &world, &err));
    assert(!strcmp(world.networks[0].name, "legacy-run-net-001"));
    assert(!strcmp(world.sites[0].name, "legacy-run-site-001"));
    world_free(&world); free(s);
    return 0;
}
