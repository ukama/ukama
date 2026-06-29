/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_RUNNER_H_
#define ULAB_RUNNER_H_

#include "report.h"

typedef struct {
    char scenario_path[ULAB_MAX_PATH];
    char bff_url[ULAB_MAX_PATH];
    char out_dir[ULAB_MAX_PATH];
    char script_dir[ULAB_MAX_PATH];
    char repo[ULAB_MAX_PATH];
    char sim_type[ULAB_MAX_REF];
    char warehouse_url[ULAB_MAX_URL];
    char factory_url[ULAB_MAX_URL];
    char asr_url[ULAB_MAX_URL];
    char run_id[ULAB_MAX_ID];

    uint32_t seed_override;

    int has_seed_override;
    int quiet;
    int verbose;
} runner_opts_t;

int runner_validate(const runner_opts_t *opts);

#endif /* ULAB_RUNNER_H_ */
