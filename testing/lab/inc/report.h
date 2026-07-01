/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_REPORT_H_
#define ULAB_REPORT_H_

#include <stdio.h>
#include <time.h>

#include "check.h"
#include "event.h"

typedef struct {
    char   run_id[ULAB_MAX_ID];
    char   run_dir[ULAB_MAX_PATH];
    char   scenario[ULAB_MAX_NAME];
    char   suite[ULAB_MAX_REF];
    char   priority[ULAB_MAX_REF];
    char   status[ULAB_MAX_REF];
    char   tags[ULAB_MAX_LINE];

    size_t events;
    size_t event_failed;
    size_t checks;
    size_t failed;

    int cleanup_failed;
    int final_rc;

    time_t started_at;
    time_t ended_at;

    FILE *json;
    FILE *txt;
    int   json_results;
} report_t;

int  report_open(report_t *r,
                 const scenario_t *scenario,
                 const char *run_id,
                 const char *run_dir);
void report_close(report_t *r);
void report_world(const world_t *w);
void report_event(report_t *r,
                  const char *phase,
                  const event_spec_t *event,
                  int passed,
                  const char *detail);
void report_check(report_t *r, const check_result_t *res);
void report_set_cleanup(report_t *r, int failed);
void report_set_final_rc(report_t *r, int rc);
void report_result(report_t *r);

#endif /* ULAB_REPORT_H_ */
