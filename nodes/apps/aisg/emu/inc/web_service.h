/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_WEB_SERVICE_H_
#define AISG_EMU_WEB_SERVICE_H_

#include <ulfius.h>

#include "config.h"
#include "model.h"

bool start_web_service(UInst *instance, EmuConfig *config, EmuModel *model);
void stop_web_service(UInst *instance);

#endif /* AISG_EMU_WEB_SERVICE_H_ */
