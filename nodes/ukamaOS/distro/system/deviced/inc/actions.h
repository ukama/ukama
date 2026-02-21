/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ACTIONS_H_
#define ACTIONS_H_

#include "config.h"
#include "control.h"

int actions_service_apply(Config *config, ControlState desired);
int actions_radio_apply(Config *config, ControlState desired);
int actions_restart_apply(Config *config);

#endif
