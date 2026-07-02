/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_RET_MODE_H_
#define AISG_EMU_RET_MODE_H_

#include <signal.h>

#include "config.h"

bool ret_mode_run(const EmuConfig *config, volatile sig_atomic_t *running);

#endif /* AISG_EMU_RET_MODE_H_ */
