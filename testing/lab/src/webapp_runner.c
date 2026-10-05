/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 */
#include "webapp.h"
#include "log.h"
#include "util.h"
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static volatile sig_atomic_t stopped;
static void stop_handler(int signal_number) { stopped = signal_number; }

static int absolute_path(char *out, size_t len, const char *path, ulab_error_t *err) {
    char cwd[ULAB_MAX_PATH];
    if (path[0] == '/') {
        if (strlen(path) >= len) return webapp_error(err, "web-app path too long");
        ulab_copy(out, len, path);
        return ULAB_OK;
    }
    if (!getcwd(cwd, sizeof(cwd))) return webapp_error(err, "cannot resolve lab working directory");
    return webapp_path(out, len, cwd, path, err);
}
int webapp_event_inputs(const event_spec_t *event, world_t *world,
                         json_t **inputs, ulab_error_t *err) {
    network_t *network = NULL;
    site_t *site;
    node_t *node;
    const char *ref = NULL;
    const char *id = NULL;
    const char *text = NULL;
    const char *network_ref = NULL;
    *inputs = json_object();
    if (!*inputs) return webapp_error(err, "cannot allocate browser event");
    if (event->type == EVT_WEB_RELOAD) return ULAB_OK;
    if (event->type == EVT_WEB_TAB) {
        if (json_object_set_new(*inputs, "tab", json_string(event->profile))) goto memory;
        return ULAB_OK;
    }
    if (event->networks.kind != SEL_NONE) {
        network = world_network_by_ref(world, event->networks.value);
        if (event->networks.kind != SEL_REF || !network || !network->bff_id[0]) goto unresolved;
        if (json_object_set_new(*inputs, "network_name", json_string(network->name))) goto memory;
    }
    if (event->type == EVT_WEB_COMMERCE) {
        if (webapp_commerce_inputs(event, world, *inputs, err)) { json_decref(*inputs); *inputs = NULL; return ULAB_ERR; }
        return ULAB_OK;
    }
    if (event->type == EVT_WEB_SELECT_NETWORK) {
        if (!network) goto unresolved;
        return ULAB_OK;
    }
    if (event->type != EVT_WEB_OPEN && event->type != EVT_WEB_ACTION) goto unresolved;
    if (json_object_set_new(*inputs, "view", json_string(event->view))) goto memory;
    if (!strcmp(event->view, "network_site_detail")) {
        site = world_site_by_ref(world, event->sites.value);
        if (!site || !site->bff_id[0]) goto unresolved;
        ref = site->ref; id = site->bff_id; text = site->name; network_ref = site->network_ref;
    } else if (!strcmp(event->view, "network_node_detail")) {
        node = world_node_by_ref(world, event->nodes.value);
        if (!node || !node->bff_id[0]) goto unresolved;
        ref = node->ref; id = node->bff_id; text = node->bff_id; network_ref = node->network_ref;
    }
    if (ref) {
        if (!network || strcmp(network->ref, network_ref)) goto unresolved;
        if (json_object_set_new(*inputs, "entity", json_pack("{s:s,s:s,s:s}", "ref", ref, "id", id, "text", text))) goto memory;
    }
    if (event->type == EVT_WEB_ACTION) {
        if (!ref) goto unresolved;
        if (json_object_set_new(*inputs, "action", json_string(event->target))) goto memory;
        if ((event->web_fields & (1u << 6)) || event->variant[0])
            if (json_object_set_new(*inputs, "value", json_string(event->variant[0] ? text : event->status))) goto memory;
        if (event->app[0] && json_object_set_new(*inputs, "app", json_string(event->app))) goto memory;
        if (event->tag[0] && json_object_set_new(*inputs, "tag", json_string(event->tag))) goto memory;
    }
    return ULAB_OK;
unresolved:
    json_decref(*inputs); *inputs = NULL;
    return webapp_error(err, "browser event requires resolved world identities in the selected network");
memory:
    json_decref(*inputs); *inputs = NULL;
    return webapp_error(err, "cannot encode browser event");
}
/* Match JavaScript /\s+/ normalization without a locale-dependent byte
 * classification of UTF-8. Numeric formatting/units remain exact. */
static size_t whitespace_width(const unsigned char *p) {
    static const char *const unicode[] = {
        "\xc2\xa0", "\xe1\x9a\x80", "\xe2\x80\x80", "\xe2\x80\x81", "\xe2\x80\x82",
        "\xe2\x80\x83", "\xe2\x80\x84", "\xe2\x80\x85", "\xe2\x80\x86", "\xe2\x80\x87",
        "\xe2\x80\x88", "\xe2\x80\x89", "\xe2\x80\x8a", "\xe2\x80\xa8", "\xe2\x80\xa9",
        "\xe2\x80\xaf", "\xe2\x81\x9f", "\xe3\x80\x80", "\xef\xbb\xbf"
    };
    size_t i;
    size_t len;
    if (*p == ' ' || (*p >= '\t' && *p <= '\r')) return 1;
    for (i = 0; i < sizeof(unicode) / sizeof(unicode[0]); i++) {
        len = strlen(unicode[i]);
        if (!strncmp((const char *)p, unicode[i], len)) return len;
    }
    return 0;
}
json_t *webapp_check_inputs(const check_spec_t *check) {
    json_t *inputs;
    const unsigned char *p;
    char normalized[ULAB_MAX_REF];
    size_t n = 0;
    size_t width;
    int space = 0;
    inputs = json_pack("{s:s,s:s,s:s}", "view", check->view, "label", check->label, "requirement", check->requirement);
    if (!inputs) return NULL;
    if (check->app[0] && json_object_set_new(inputs, "app", json_string(check->app))) goto fail;
    if (check->variant[0] && json_object_set_new(inputs, "match", json_string(check->variant))) goto fail;
    if (check->type == CHECK_WEB_ACTION_AVAILABLE) {
        if (json_object_set_new(inputs, "available", json_boolean(check->expected_value != 0))) goto fail;
    } else if (check->type == CHECK_WEB_TABLE_COUNT_EQUALS) {
        if (json_object_set_new(inputs, "expected_count", json_integer(check->expected_count))) goto fail;
    } else {
        for (p = (const unsigned char *)check->expected; *p; p++) {
            width = whitespace_width(p);
            if (width) { if (n) space = 1; p += width - 1; continue; }
            if (space) normalized[n++] = ' ';
            normalized[n++] = (char)*p; space = 0;
        }
        normalized[n] = '\0';
        if (json_object_set_new(inputs, "expected", json_string(normalized))) goto fail;
    }
    return inputs;
fail:
    json_decref(inputs);
    return NULL;
}
static int actual_matches(const check_spec_t *check, json_t *expected, json_t *actual) {
    if (!strcmp(check->variant, "contains"))
        return json_is_string(expected) && json_is_string(actual) &&
            strstr(json_string_value(actual), json_string_value(expected)) != NULL;
    return json_equal(expected, actual);
}
static int check_one(webapp_client_t *client, world_t *world, report_t *report,
                       const char *phase, const check_spec_t *check,
                       ulab_error_t *err) {
    json_t *inputs;
    json_t *reply = NULL;
    json_t *expected;
    int rc;
    check_spec_t resolved;
    node_t *node;
    site_t *site;
    resolved = *check;
    if (check->ref[0]) {
        node = world_node_by_ref(world, check->ref);
        if (!node || !node->bff_id[0]) return webapp_error(err, "expected_ref has no provisioned node identity");
        site = world_site_by_ref(world, node->site_ref);
        if (!site) return webapp_error(err, "expected_ref site missing");
        if (ulab_copy(resolved.expected, sizeof(resolved.expected), !strcmp(check->key, "id") ? node->bff_id :
                  !strcmp(check->key, "site_name") ? site->name : !strcmp(node->type, "tower") ? "Tower node" :
                  !strcmp(node->type, "amplifier") ? "Amplifier node" : "Controller node"))
            return webapp_error(err, "resolved browser expectation is too long");
    }
    if (check->type == CHECK_WEB_COMMERCE_EQUALS) {
        if (webapp_commerce_check(check, &resolved, world, &inputs, err)) return ULAB_ERR;
    } else inputs = webapp_check_inputs(&resolved);
    if (inputs && check->nodes.kind != SEL_NONE) {
        node = world_node_by_ref(world, check->nodes.value);
        if (!node || !node->bff_id[0]) { json_decref(inputs); return webapp_error(err, "node card identity is unresolved"); }
        json_object_set_new(inputs, "node_id", json_string(node->bff_id));
    }
    if (!inputs) return webapp_error(err, "cannot encode browser assertion");
    expected = json_object_get(inputs, check->type == CHECK_WEB_ACTION_AVAILABLE ? "available" :
                               check->type == CHECK_WEB_TABLE_COUNT_EQUALS ? "expected_count" : "expected");
    ulab_status("CHECK", "%s [%s] %s", phase, check->requirement, check->label);
    rc = webapp_call(client, scenario_check_name(check->type), inputs,
                      check->timeout_seconds, &reply, err);
    if (!rc && (!json_equal(expected, json_object_get(reply, "expected")) ||
                !actual_matches(check, expected, json_object_get(reply, "actual")) ||
                json_array_size(json_object_get(reply, "bindings")))) {
        client->broken = 1;
        rc = webapp_error(err, "worker PASS does not match the requested visible expectation");
    }
    if (reply) json_object_set(reply, "expected", expected);
    if (!reply) reply = json_pack("{s:O,s:n,s:[]}", "expected", expected, "actual", "artifacts");
    if (report_web_check(report, phase, check, reply, !rc, rc ? err->msg : ""))
        rc = webapp_error(err, "cannot write browser assertion report");
    json_decref(reply); json_decref(inputs);
    return rc;
}
typedef struct { const webapp_hooks_t *hooks; const event_spec_t *event; } runtime_job_t;
static int runtime_event_job(void *ctx, ulab_error_t *err) {
    runtime_job_t *job = ctx;
    return job->hooks->runtime_event(job->hooks->ctx, job->event, err);
}
static int event_one(webapp_client_t *client, webapp_journal_t *journal,
                       const webapp_hooks_t *hooks, const event_spec_t *event,
                       unsigned int default_timeout, ulab_error_t *err) {
    json_t *inputs = NULL;
    json_t *reply = NULL;
    json_t *bindings;
    json_t *entity;
    json_t *binding;
    json_t *intent = NULL;
    const char *kind = webapp_commerce_kind(event);
    double deadline;
    runtime_job_t job;
    int rc;
    if (!scenario_is_web_event(event->type)) {
        if (!hooks || !hooks->runtime_event) return webapp_error(err, "runtime event requires a provisioned runtime handler");
        job.hooks = hooks; job.event = event;
        deadline = webapp_now() + (event->timeout_seconds ? event->timeout_seconds : default_timeout);
        if (deadline > client->deadline) deadline = client->deadline;
        rc = webapp_bounded_job(runtime_event_job, &job, deadline, client->cancel, err);
        if (!rc && event->type == EVT_START_UES) {
            size_t i;
            for (i = 0; i < journal->world->ue_count; i++) {
                ue_t *ue = &journal->world->ues[i];
                if (event->ues.kind == SEL_ALL || !strcmp(event->ues.value, ue->ref)) ue->started = ue->attached = 1;
            }
        }
        return rc;
    }
    rc = webapp_event_inputs(event, journal->world, &inputs, err);
    if (rc) return rc;
    if (kind) {
        json_t *intents = json_object_get(journal->root, "creation_intents");
        if (!intents) { intents = json_array(); json_object_set_new(journal->root, "creation_intents", intents); }
        intent = json_deep_copy(json_object_get(inputs, "creation"));
        json_object_set_new(intent, "command_id", json_integer(client->sequence + 1));
        json_object_set_new(intent, "state", json_string("pending"));
        if (json_array_append_new(intents, intent) || webapp_journal_save(journal, err)) { json_decref(inputs); return ULAB_ERR; }
    }
    rc = webapp_call(client, scenario_event_name(event->type), inputs,
                      event->timeout_seconds, &reply, err);
    if (reply && !client->broken) {
        bindings = json_object_get(reply, "bindings");
        entity = json_object_get(inputs, "entity");
        binding = json_array_get(bindings, 0);
        if (kind) {
            if (!rc && (json_array_size(bindings) != 1 ||
                !json_equal(json_object_get(intent, "kind"), json_object_get(binding, "kind")) ||
                !json_equal(json_object_get(intent, "ref"), json_object_get(binding, "ref")) ||
                !json_equal(json_object_get(intent, "name"), json_object_get(binding, "name")) ||
                !json_is_true(json_object_get(json_object_get(reply, "actual"), "executed"))))
                rc = webapp_error(err, "commerce creation acknowledgement does not match intent");
            else if (!rc && webapp_journal_bind(journal, bindings, 1, client->sequence, err)) rc = ULAB_ERR;
        } else if (!rc && (event->type == EVT_WEB_ACTION || event->type == EVT_WEB_TAB || event->type == EVT_WEB_COMMERCE) &&
            (json_array_size(bindings) || !json_is_true(json_object_get(json_object_get(reply, "actual"), "executed")))) {
            rc = webapp_error(err, "worker operation acknowledgement is invalid");
        } else if (json_array_size(bindings) && (!entity || json_array_size(bindings) != 1 ||
            !json_equal(json_object_get(entity, "ref"), json_object_get(binding, "ref")) ||
            !json_equal(json_object_get(entity, "id"), json_object_get(binding, "id")))) {
            rc = webapp_error(err, "worker binding does not match the requested detail entity");
        } else if (webapp_journal_bind(journal, bindings, 0, client->sequence, err)) rc = ULAB_ERR;
    }
    json_decref(inputs); json_decref(reply);
    return rc;
}
static int initialize(webapp_client_t *client, const scenario_t *scenario,
                        const char *run_dir, ulab_error_t *err) {
    char auth[ULAB_MAX_PATH];
    char artifacts[ULAB_MAX_PATH];
    json_t *profile;
    json_t *inputs;
    json_t *reply = NULL;
    int rc;
    if (absolute_path(auth, sizeof(auth), scenario->webapp.auth_state, err) ||
        webapp_path(artifacts, sizeof(artifacts), run_dir, "browser", err)) return ULAB_ERR;
    profile = json_pack("{s:s,s:s,s:s,s:b,s:i,s:i,s:i}",
                        "base_url", scenario->webapp.base_url, "auth_state", auth,
                        "browser", scenario->webapp.browser, "headless", scenario->webapp.headless,
                        "action_timeout_seconds", scenario->webapp.action_timeout_seconds,
                        "check_timeout_seconds", scenario->webapp.check_timeout_seconds,
                        "scenario_timeout_seconds", scenario->webapp.scenario_timeout_seconds);
    inputs = profile ? json_pack("{s:o,s:s}", "profile", profile, "artifacts_dir", artifacts) : NULL;
    if (!inputs) return webapp_error(err, "cannot encode worker initialization");
    rc = webapp_call(client, "init", inputs, scenario->webapp.action_timeout_seconds, &reply, err);
    if (!rc && (!json_is_true(json_object_get(json_object_get(reply, "actual"), "authenticated")) ||
                json_array_size(json_object_get(reply, "bindings")))) rc = webapp_error(err, "worker initialization did not establish an authenticated session");
    json_decref(inputs); json_decref(reply);
    return rc;
}
int webapp_execute(const runner_opts_t *opts, const scenario_t *scenario,
                    world_t *world, report_t *report, const char *run_dir,
                    const webapp_hooks_t *hooks, ulab_error_t *err) {
    webapp_client_t client;
    webapp_journal_t journal;
    struct sigaction action;
    struct sigaction old_int;
    struct sigaction old_term;
    char absolute[ULAB_MAX_PATH];
    char world_path[ULAB_MAX_PATH];
    const char *worker;
    size_t p;
    size_t i;
    int rc = ULAB_ERR;
    int cleanup_failed = 0;
    int worker_cleanup_failed = 0;
    int runtime_cleanup_failed = 0;
    int signals = 0;
    int browser_started = 0;
    double deadline;
    ulab_error_t cleanup_error;

    memset(&client, 0, sizeof(client)); client.input = client.output = -1;
    memset(&journal, 0, sizeof(journal));
    memset(&action, 0, sizeof(action));
    stopped = 0; action.sa_handler = stop_handler; sigemptyset(&action.sa_mask);
    if (sigaction(SIGINT, &action, &old_int)) return webapp_error(err, "cannot install web-app cancellation handler");
    if (sigaction(SIGTERM, &action, &old_term)) { sigaction(SIGINT, &old_int, NULL); return webapp_error(err, "cannot install web-app cancellation handler"); }
    signals = 1;
    deadline = webapp_now() + scenario->webapp.scenario_timeout_seconds;
    if (absolute_path(absolute, sizeof(absolute), run_dir, err)) goto done;
    ulab_copy(report->webapp_artifacts, sizeof(report->webapp_artifacts), absolute);
    if (webapp_journal_open(&journal, world, absolute, err)) goto done;
    worker = opts->webapp_worker[0] ? opts->webapp_worker : "utils/webapp-worker.sh";
    ulab_status("WEBAPP", "start local browser worker");
    if (webapp_client_start(&client, worker, world->run_id, absolute, deadline, &stopped, err)) goto done;
    browser_started = 1;
    if (initialize(&client, scenario, absolute, err)) goto done;
    if (hooks && hooks->provision && hooks->provision(hooks->ctx, &client, &journal, err)) goto done;
    for (p = 0; p < scenario->phase_count; p++) {
        const phase_spec_t *phase;
        phase = &scenario->phases[p];
        ulab_status("PHASE", "%s", phase->name);
        for (i = 0; i < phase->event_count; i++) {
            rc = event_one(&client, &journal, hooks, &phase->events[i], scenario->webapp.action_timeout_seconds, err);
            report_event(report, phase->name, &phase->events[i], !rc, rc ? err->msg : "ok");
            if (rc) goto done;
        }
        for (i = 0; i < phase->check_count; i++)
            if ((rc = check_one(&client, world, report, phase->name, &phase->checks[i], err))) goto done;
    }
    for (i = 0; i < scenario->final_check_count; i++)
        if ((rc = check_one(&client, world, report, "final", &scenario->final_checks[i], err))) goto done;
    rc = ULAB_OK;
done:
    if (stopped || webapp_now() > deadline) {
        if (!err->msg[0]) webapp_error(err, stopped ? "web-app scenario cancelled" : "web-app scenario deadline exceeded");
        rc = ULAB_ERR;
    }
    memset(&cleanup_error, 0, sizeof(cleanup_error));
    if (browser_started && webapp_client_stop(&client, rc != ULAB_OK, &cleanup_error)) worker_cleanup_failed = cleanup_failed = 1;
    ulab_status("CLEANUP", "web-app worker and test-owned resources");
    if (hooks && hooks->recover && hooks->recover(hooks->ctx, &journal, &cleanup_error)) cleanup_failed = 1;
    if (hooks && hooks->prepare_cleanup && webapp_bounded_job(hooks->prepare_cleanup, hooks->ctx,
        webapp_now() + 30, NULL, &cleanup_error)) cleanup_failed = 1;
    if (webapp_journal_cleanup(&journal, hooks, webapp_now() + 120, &cleanup_error)) cleanup_failed = 1;
    if (hooks && hooks->cleanup_runtime && webapp_bounded_job(hooks->cleanup_runtime, hooks->ctx,
        webapp_now() + 15, NULL, &cleanup_error)) runtime_cleanup_failed = cleanup_failed = 1;
    if (world && (webapp_path(world_path, sizeof(world_path), run_dir, "world.json", &cleanup_error) ||
        world_write_json(world, world_path))) {
        cleanup_failed = 1;
        webapp_error(&cleanup_error, "cannot persist final world artifact");
    }
    if ((report->json && (fflush(report->json) || ferror(report->json))) ||
        (report->txt && (fflush(report->txt) || ferror(report->txt))))
        rc = webapp_error(err, "cannot persist web-app report");
    if (stopped && rc == ULAB_OK)
        rc = webapp_error(err, "web-app scenario cancelled during cleanup");
    if (journal.root) {
        json_object_set_new(journal.root, "run_result", json_string(rc || cleanup_failed ? "failed" : "passed"));
        json_object_set_new(journal.root, "worker_cleanup", json_string(worker_cleanup_failed ? "unconfirmed" : browser_started ? "complete" : "not_started"));
        json_object_set_new(journal.root, "runtime_cleanup", json_string(runtime_cleanup_failed ? "failed" : hooks && hooks->cleanup_runtime ? "complete" : "not_requested"));
        if (webapp_journal_save(&journal, &cleanup_error)) cleanup_failed = 1;
    }
    if (cleanup_failed) {
        report_set_cleanup(report, 1);
        if (rc == ULAB_OK) { *err = cleanup_error; rc = ULAB_ERR; }
        ulab_log_error("web-app cleanup: %s", cleanup_error.msg);
    }
    if (rc) ulab_copy(report->error, sizeof(report->error), err->msg);
    webapp_journal_close(&journal);
    if (signals) { sigaction(SIGINT, &old_int, NULL); sigaction(SIGTERM, &old_term, NULL); }
    return rc;
}
