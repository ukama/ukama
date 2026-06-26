/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_SIM_FACTORY_H_
#define ULAB_SIM_FACTORY_H_

#include "runner.h"
#include "world.h"
#include "ulab.h"

int sim_factory_prepare_world(const runner_opts_t *opts,
                              world_t *world,
                              const char *run_dir,
                              char *csv_path,
                              size_t csv_path_len,
                              ulab_error_t *err);

int sim_factory_wait_asr(const runner_opts_t *opts,
                         const ue_t *ue,
                         ulab_error_t *err);

#endif /* ULAB_SIM_FACTORY_H_ */
