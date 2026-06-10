/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_OPS_H_
#define AISG_EMU_OPS_H_

#include "model.h"
#include "request.h"
#include "response.h"

bool emu_ops_handle(EmuModel *model,
                    const EmuRequest *request,
                    EmuResponse *response);

#endif /* AISG_EMU_OPS_H_ */
