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

#include "check.h"

typedef struct {
    char   run_id[ULAB_MAX_ID];
    char   run_dir[ULAB_MAX_PATH];
    size_t checks;
    size_t failed;
    FILE   *json;
} report_t;

int  report_open(report_t *r, const char *run_id, const char *run_dir);
void report_close(report_t *r);
void report_world(const world_t *w);
void report_check(report_t *r, const check_result_t *res);
void report_result(report_t *r);

#endif /* ULAB_REPORT_H_ */
