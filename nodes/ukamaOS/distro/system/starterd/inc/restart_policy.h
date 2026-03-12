/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <time.h>
#include "config.h"
#include "app.h"

int restart_policy_next_delay(Config *config, App *app, time_t now);
void restart_policy_on_start(Config *config, App *app, time_t now);
void restart_policy_on_exit(Config *config, App *app, time_t now);
