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

int report_open(report_t *r, const char *run_id, const char *run_dir) {
    char path[ULAB_MAX_PATH];

    memset(r, 0, sizeof(*r));
    ulab_copy(r->run_id, sizeof(r->run_id), run_id);
    ulab_copy(r->run_dir, sizeof(r->run_dir), run_dir);
    snprintf(path, sizeof(path), "%s/report.json", run_dir);
    r->json = fopen(path, "w");
    if (r->json == NULL) {
        return ULAB_ERR;
    }
    fprintf(r->json, "{\n  \"run_id\": \"%s\",\n", run_id);
    fprintf(r->json, "  \"checks\": [\n");
    return ULAB_OK;
}

void report_close(report_t *r) {
    if (r->json == NULL) {
        return;
    }
    fprintf(r->json, "\n  ],\n");
    fprintf(r->json, "  \"passed\": %s,\n", r->failed ? "false" : "true");
    fprintf(r->json, "  \"failed\": %zu\n}\n", r->failed);
    fclose(r->json);
    r->json = NULL;
}

void report_world(const world_t *w) {
    ulab_status("WORLD", "networks=%zu sites=%zu nodes=%zu ues=%zu",
                w->network_count, w->site_count, w->node_count, w->ue_count);
    ulab_status("WORLD", "subscribers=%zu packages=%zu", w->subscriber_count,
                w->package_count);
}

void report_check(report_t *r, const check_result_t *res) {
    const char *state = res->skipped ? "SKIP" : (res->passed ? "PASS" : "FAIL");

    r->checks++;
    if (!res->passed && !res->skipped) {
        r->failed++;
    }
    ulab_status(state, "%s: %s", res->name, res->detail);
    if (r->json) {
        if (r->checks > 1) {
            fprintf(r->json, ",\n");
        }
        fprintf(r->json,
                "    {\"name\":\"%s\",\"state\":\"%s\","
                "\"detail\":\"%s\"}", res->name, state, res->detail);
        fflush(r->json);
    }
}

void report_result(report_t *r) {
    if (r->failed) {
        ulab_status("FAIL", "checks=%zu failed=%zu artifacts=%s",
                    r->checks, r->failed, r->run_dir);
    } else {
        ulab_status("PASS", "checks=%zu artifacts=%s", r->checks,
                    r->run_dir);
    }
}
