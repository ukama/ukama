/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_SERVER_H_
#define AISG_EMU_SERVER_H_

#include "config.h"
#include "model.h"

bool emu_server_run(EmuConfig *config, EmuModel *model, volatile bool *running);

#endif /* AISG_EMU_SERVER_H_ */
