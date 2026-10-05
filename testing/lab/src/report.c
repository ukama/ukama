/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "report.h"
#include "log.h"
#include "util.h"

#include <stdio.h>
#include <string.h>
#include <time.h>
#include <stdlib.h>
#include <sys/stat.h>
#include <unistd.h>

static void json_str(FILE *f, const char *key, const char *value,
                     int comma) {
    char *encoded;
    json_t *text;

    text = json_string(value ? value : "");
    encoded = text ? json_dumps(text, JSON_ENCODE_ANY) : NULL;
    fprintf(f, "  \"%s\": %s%s\n", key, encoded ? encoded : "null", comma ? "," : "");
    free(encoded); json_decref(text);
}

static void json_result_prefix(report_t *r) {
    if (r->json_results) {
        fprintf(r->json, ",\n");
    }
    r->json_results = 1;
}

int report_open(report_t *r,
                const scenario_t *scenario,
                const char *run_id,
                const char *run_dir) {
    char path[ULAB_MAX_PATH];

    memset(r, 0, sizeof(*r));
    ulab_copy(r->run_id, sizeof(r->run_id), run_id);
    ulab_copy(r->run_dir, sizeof(r->run_dir), run_dir);

    if (scenario != NULL) {
        ulab_copy(r->scenario, sizeof(r->scenario), scenario->name);
        ulab_copy(r->description, sizeof(r->description),
                  scenario->description);
        ulab_copy(r->suite, sizeof(r->suite), scenario->suite);
        ulab_copy(r->priority, sizeof(r->priority), scenario->priority);
        ulab_copy(r->status, sizeof(r->status), scenario->status);
        ulab_copy(r->tags, sizeof(r->tags), scenario->tags);
    }

    r->started_at = time(NULL);

    snprintf(path, sizeof(path), "%s/report.json", run_dir);
    r->json = fopen(path, "w");
    if (r->json == NULL) {
        return ULAB_ERR;
    }

    snprintf(path, sizeof(path), "%s/report.txt", run_dir);
    r->txt = fopen(path, "w");
    if (r->txt == NULL) {
        fclose(r->json);
        r->json = NULL;
        return ULAB_ERR;
    }
    if (scenario && scenario->version == ULAB_WEBAPP_SCHEMA_VER &&
        (fchmod(fileno(r->json), 0600) || fchmod(fileno(r->txt), 0600))) {
        fclose(r->json); fclose(r->txt);
        r->json = r->txt = NULL;
        return ULAB_ERR;
    }

    fprintf(r->json, "{\n");
    json_str(r->json, "run_id", r->run_id, 1);
    json_str(r->json, "scenario", r->scenario, 1);
    json_str(r->json, "description", r->description, 1);
    json_str(r->json, "suite", r->suite, 1);
    json_str(r->json, "priority", r->priority, 1);
    json_str(r->json, "status", r->status, 1);
    json_str(r->json, "tags", r->tags, 1);
    fprintf(r->json, "  \"started_at\": %ld,\n", (long)r->started_at);
    fprintf(r->json, "  \"results\": [\n");
    fflush(r->json);

    fprintf(r->txt, "run_id: %s\n", r->run_id);
    fprintf(r->txt, "scenario: %s\n", r->scenario);
    fprintf(r->txt, "description: %s\n", r->description);
    fprintf(r->txt, "suite: %s\n", r->suite);
    fprintf(r->txt, "priority: %s\n", r->priority);
    fprintf(r->txt, "status: %s\n", r->status);
    fprintf(r->txt, "tags: %s\n\n", r->tags);
    fflush(r->txt);

    return ULAB_OK;
}

void report_close(report_t *r) {
    int passed;
    char *directory;
    json_t *value;

    r->ended_at = time(NULL);
    passed = r->final_rc == ULAB_OK && r->failed == 0 &&
        r->event_failed == 0 && r->cleanup_failed == 0 && !r->scenario_skipped;

    if (r->json != NULL) {
        fprintf(r->json, "\n  ],\n");
        fprintf(r->json, "  \"ended_at\": %ld,\n", (long)r->ended_at);
        fprintf(r->json, "  \"duration_sec\": %ld,\n",
                (long)(r->ended_at - r->started_at));
        fprintf(r->json, "  \"events\": {\"total\": %zu, \"passed\": %zu, \"failed\": %zu},\n",
                r->events, r->events - r->event_failed, r->event_failed);
        fprintf(r->json, "  \"checks\": {\"total\": %zu, \"passed\": %zu, \"failed\": %zu},\n",
                r->checks, r->checks - r->failed, r->failed);
        fprintf(r->json, "  \"cleanup\": \"%s\",\n",
                r->cleanup_failed ? "failed" : "ok");
        value = json_string(r->run_dir);
        directory = value ? json_dumps(value, JSON_ENCODE_ANY) : NULL;
        fprintf(r->json, "  \"artifacts\": {\"run_dir\": %s", directory ? directory : "null");
        free(directory); json_decref(value);
        /* Paths are JSON encoded, including spaces, quotes and backslashes. */
        {
            const char *names[] = {"world", "model", "created", "created_final"};
            const char *files[] = {"world.json", "model.json", "created.json", "created.final.json"};
            char path[ULAB_MAX_PATH + 32];
            size_t i;
            for (i = 0; i < sizeof(names) / sizeof(names[0]); i++) {
                /* Browser runs do not produce backend-created/model files. */
                if ((r->webapp_artifacts[0] || r->scenario_skipped) && i > 0) continue;
                snprintf(path, sizeof(path), "%s/%s", r->run_dir, files[i]);
                value = json_string(path);
                directory = value ? json_dumps(value, JSON_ENCODE_ANY) : NULL;
                fprintf(r->json, ", \"%s\": %s", names[i], directory ? directory : "null");
                free(directory); json_decref(value);
            }
        }
        fprintf(r->json, "},\n");
        fprintf(r->json, "  \"passed\": %s,\n", passed ? "true" : "false");
        json_str(r->json, "outcome", r->scenario_skipped ? "SKIP" : passed ? "PASS" : "FAIL", 1);
        if (r->webapp_artifacts[0]) json_str(r->json, "webapp_artifacts", r->webapp_artifacts, 1);
        if (r->error[0]) json_str(r->json, "error", r->error, 1);
        fprintf(r->json, "  \"final_rc\": %d\n", r->final_rc);
        fprintf(r->json, "}\n");
        fclose(r->json);
        r->json = NULL;
    }

    if (r->txt != NULL) {
        fprintf(r->txt, "\nsummary:\n");
        fprintf(r->txt, "  events: %zu passed, %zu failed, %zu total\n",
                r->events - r->event_failed, r->event_failed, r->events);
        fprintf(r->txt, "  checks: %zu passed, %zu failed, %zu total\n",
                r->checks - r->failed, r->failed, r->checks);
        fprintf(r->txt, "  cleanup: %s\n", r->cleanup_failed ? "failed" : "ok");
        if (r->error[0]) fprintf(r->txt, "  error: %s\n", r->error);
        fprintf(r->txt, "  result: %s\n", r->scenario_skipped ? "SKIP" : passed ? "PASS" : "FAIL");
        fclose(r->txt);
        r->txt = NULL;
    }
}

void report_world(const world_t *w) {
    ulab_status("WORLD", "networks=%zu sites=%zu nodes=%zu ues=%zu",
                w->network_count, w->site_count, w->node_count, w->ue_count);
    ulab_status("WORLD", "subscribers=%zu packages=%zu", w->subscriber_count,
                w->package_count);
}

void report_event(report_t *r,
                  const char *phase,
                  const event_spec_t *event,
                  int passed,
                  const char *detail) {
    char esc_detail[ULAB_MAX_ERR * 2];
    char esc_phase[ULAB_MAX_NAME * 2];
    const char *state;

    if (r == NULL || event == NULL) {
        return;
    }

    state = passed ? "PASS" : "FAIL";
    r->events++;
    if (!passed) {
        r->event_failed++;
    }

    ulab_status(state, "event %s/%s: %s", phase ? phase : "",
                scenario_event_name(event->type), detail ? detail : "ok");

    if (r->txt != NULL) {
        fprintf(r->txt, "%s event %s/%s: %s\n", state,
                phase ? phase : "", scenario_event_name(event->type),
                detail ? detail : "ok");
        fflush(r->txt);
    }

    if (r->json != NULL) {
        ulab_json_escape(detail ? detail : "ok", esc_detail, sizeof(esc_detail));
        ulab_json_escape(phase ? phase : "", esc_phase, sizeof(esc_phase));
        json_result_prefix(r);
        fprintf(r->json,
                "    {\"kind\":\"event\",\"phase\":\"%s\","
                "\"name\":\"%s\",\"state\":\"%s\","
                "\"detail\":\"%s\"}",
                esc_phase, scenario_event_name(event->type), state,
                esc_detail);
        fflush(r->json);
    }
}

void report_check(report_t *r, const check_result_t *res) {
    char esc_detail[ULAB_MAX_ERR * 2];
    const char *state;

    state = res->skipped ? "SKIP" : (res->passed ? "PASS" : "FAIL");

    r->checks++;
    if (!res->passed && !res->skipped) {
        r->failed++;
    }

    ulab_status(state, "%s: %s", res->name, res->detail);

    if (r->txt != NULL) {
        fprintf(r->txt, "%s check %s: %s\n", state, res->name, res->detail);
        fflush(r->txt);
    }

    if (r->json != NULL) {
        ulab_json_escape(res->detail, esc_detail, sizeof(esc_detail));
        json_result_prefix(r);
        fprintf(r->json,
                "    {\"kind\":\"check\",\"name\":\"%s\","
                "\"state\":\"%s\",\"detail\":\"%s\"}",
                res->name, state, esc_detail);
        fflush(r->json);
    }
}

void report_set_cleanup(report_t *r, int failed) {
    if (r != NULL) {
        r->cleanup_failed = failed ? 1 : 0;
    }
}

int report_web_check(report_t *r, const char *phase, const check_spec_t *check,
                      json_t *response, int passed, const char *error) {
    json_t *record;
    json_t *expected;
    json_t *actual;
    json_t *artifacts;
    char *want;
    char *got;
    int rc = ULAB_OK;
    expected = json_object_get(response, "expected");
    actual = json_object_get(response, "actual");
    artifacts = json_object_get(response, "artifacts");
    want = json_dumps(expected ? expected : json_null(), JSON_COMPACT | JSON_ENCODE_ANY);
    got = json_dumps(actual ? actual : json_null(), JSON_COMPACT | JSON_ENCODE_ANY);
    r->checks++;
    if (!passed) r->failed++;
    ulab_status(passed ? "PASS" : "FAIL", "%s/%s [%s] %s expected=%s actual=%s%s%s",
                phase, scenario_check_name(check->type), check->requirement, check->label,
                want ? want : "null", got ? got : "null", error && *error ? ": " : "", error ? error : "");
    if (r->txt) {
        if (fprintf(r->txt, "%s check %s/%s [%s] %s expected=%s actual=%s %s\n",
                    passed ? "PASS" : "FAIL", phase, scenario_check_name(check->type), check->requirement,
                    check->label, want ? want : "null", got ? got : "null", error ? error : "") < 0 || fflush(r->txt)) rc = ULAB_ERR;
    }
    free(want); free(got);
    record = json_pack("{s:s,s:s,s:s,s:s,s:s,s:s,s:s,s:O,s:O}",
                        "kind", "check", "phase", phase, "name", scenario_check_name(check->type),
                        "state", passed ? "PASS" : "FAIL", "requirement", check->requirement,
                        "label", check->label, "detail", error ? error : "",
                        "expected", expected ? expected : json_null(), "actual", actual ? actual : json_null());
    if (!record) return ULAB_ERR;
    if (json_object_set_new(record, "match", json_string(check->variant[0] ? check->variant : "equals"))) rc = ULAB_ERR;
    if (check->app[0] && json_object_set_new(record, "app", json_string(check->app))) rc = ULAB_ERR;
    if (check->package_ref[0] && json_object_set_new(record, "package", json_string(check->package_ref))) rc = ULAB_ERR;
    if (check->ues.kind == SEL_REF && json_object_set_new(record, "ues", json_string(check->ues.value))) rc = ULAB_ERR;
    if (check->type == CHECK_WEB_COMMERCE_EQUALS && check->key[0] && json_object_set_new(record, "expected_property", json_string(check->key))) rc = ULAB_ERR;
    if (json_object_set(record, "artifacts", artifacts ? artifacts : json_null())) rc = ULAB_ERR;
    if (r->json) {
        json_result_prefix(r);
        if (json_dumpf(record, r->json, JSON_COMPACT) || fflush(r->json)) rc = ULAB_ERR;
    }
    json_decref(record);
    return rc;
}

void report_set_final_rc(report_t *r, int rc) {
    if (r != NULL) {
        r->final_rc = rc;
    }
}

void report_result(report_t *r) {
    if (r == NULL || (r->json == NULL && r->txt == NULL)) {
        return;
    }

    if (r->scenario_skipped) {
        ulab_status("SKIP", "%s status=%s; no execution or coverage credit", r->scenario, r->status);
        return;
    }

    if (r->failed || r->event_failed || r->cleanup_failed ||
        r->final_rc != ULAB_OK) {
        ulab_status("FAIL", "events=%zu failed=%zu checks=%zu failed=%zu artifacts=%s",
                    r->events, r->event_failed, r->checks, r->failed,
                    r->run_dir);
    } else {
        ulab_status("PASS", "events=%zu checks=%zu artifacts=%s",
                    r->events, r->checks, r->run_dir);
    }
}
