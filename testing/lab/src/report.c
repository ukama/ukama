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

static void json_str(FILE *f, const char *key, const char *value,
                     int comma) {
    char esc[ULAB_MAX_LINE * 2];

    ulab_json_escape(value ? value : "", esc, sizeof(esc));
    fprintf(f, "  \"%s\": \"%s\"%s\n", key, esc, comma ? "," : "");
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

    fprintf(r->json, "{\n");
    json_str(r->json, "run_id", r->run_id, 1);
    json_str(r->json, "scenario", r->scenario, 1);
    json_str(r->json, "suite", r->suite, 1);
    json_str(r->json, "priority", r->priority, 1);
    json_str(r->json, "status", r->status, 1);
    json_str(r->json, "tags", r->tags, 1);
    fprintf(r->json, "  \"started_at\": %ld,\n", (long)r->started_at);
    fprintf(r->json, "  \"results\": [\n");
    fflush(r->json);

    fprintf(r->txt, "run_id: %s\n", r->run_id);
    fprintf(r->txt, "scenario: %s\n", r->scenario);
    fprintf(r->txt, "suite: %s\n", r->suite);
    fprintf(r->txt, "priority: %s\n", r->priority);
    fprintf(r->txt, "status: %s\n", r->status);
    fprintf(r->txt, "tags: %s\n\n", r->tags);
    fflush(r->txt);

    return ULAB_OK;
}

void report_close(report_t *r) {
    int passed;

    r->ended_at = time(NULL);
    passed = r->final_rc == ULAB_OK && r->failed == 0 &&
        r->event_failed == 0 && r->cleanup_failed == 0;

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
        fprintf(r->json, "  \"artifacts\": {"
                "\"run_dir\": \"%s\", "
                "\"world\": \"%s/world.json\", "
                "\"model\": \"%s/model.json\", "
                "\"created\": \"%s/created.json\", "
                "\"created_final\": \"%s/created.final.json\"},\n",
                r->run_dir, r->run_dir, r->run_dir, r->run_dir, r->run_dir);
        fprintf(r->json, "  \"passed\": %s,\n", passed ? "true" : "false");
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
        fprintf(r->txt, "  result: %s\n", passed ? "PASS" : "FAIL");
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
        json_result_prefix(r);
        fprintf(r->json,
                "    {\"kind\":\"event\",\"phase\":\"%s\","
                "\"name\":\"%s\",\"state\":\"%s\","
                "\"detail\":\"%s\"}",
                phase ? phase : "", scenario_event_name(event->type), state,
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

void report_set_final_rc(report_t *r, int rc) {
    if (r != NULL) {
        r->final_rc = rc;
    }
}

void report_result(report_t *r) {
    if (r == NULL || (r->json == NULL && r->txt == NULL)) {
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
