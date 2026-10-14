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
#include "hardware.h"
#include <signal.h>
#include <stdatomic.h>

typedef struct {
    char scenario_path[ULAB_MAX_PATH];
    char bff_url[ULAB_MAX_PATH];
    char hub_url[ULAB_MAX_URL];
    char out_dir[ULAB_MAX_PATH];
    char script_dir[ULAB_MAX_PATH];
    char repo[ULAB_MAX_PATH];
    char sim_type[ULAB_MAX_REF];
    char warehouse_url[ULAB_MAX_URL];
    char factory_url[ULAB_MAX_URL];
    char asr_url[ULAB_MAX_URL];
    char run_id[ULAB_MAX_ID];
    char suite_filter[ULAB_MAX_REF];
    char priority_filter[ULAB_MAX_REF];
    char tag_filter[ULAB_MAX_REF];
    char workload_assets[ULAB_MAX_PATH];
    char webapp_worker[ULAB_MAX_PATH];
    char provider[ULAB_MAX_REF]; /* Empty means use scenario provider.type. */
    char hardware_config[ULAB_MAX_PATH];
    _Atomic sig_atomic_t *workload_cancel; /* NULL for legacy scenarios */

    uint32_t seed_override;

    int has_seed_override;
    int quiet;
    int verbose;
    int generated_only;
} runner_opts_t;

int runner_validate(const runner_opts_t *opts);
int runner_load_scenario(const runner_opts_t *opts, scenario_t *scenario,
                          hardware_config_t *hardware, ulab_error_t *err);
int runner_preflight(const runner_opts_t *opts);

#endif /* ULAB_RUNNER_H_ */
