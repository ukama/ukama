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
    char subscriber_network_id[ULAB_MAX_ID];
    char sim_csv_path[ULAB_MAX_PATH];
    char sim_type[ULAB_MAX_REF];
    char warehouse_url[ULAB_MAX_URL];
    char factory_url[ULAB_MAX_URL];
    char asr_url[ULAB_MAX_URL];
    char sim_org[ULAB_MAX_REF];
    char sim_vendor[ULAB_MAX_REF];
    char sim_profile[ULAB_MAX_REF];
    char sim_form_factor[ULAB_MAX_REF];
    char sim_batch_prefix[ULAB_MAX_REF];
    char run_id[ULAB_MAX_ID];

    uint32_t seed_override;

    int has_seed_override;
    int dry_run;
    int setup_only;
    int subscriber_only;
    int print_world;
    int print_plan;
    int cleanup;
    int keep;
    int keep_on_failure;
    int quiet;
    int verbose;
} runner_opts_t;

int runner_validate(const runner_opts_t *opts);
int runner_dry_run(const runner_opts_t *opts);

#endif /* ULAB_RUNNER_H_ */
