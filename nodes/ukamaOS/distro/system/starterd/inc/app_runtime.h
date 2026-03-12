/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <sys/types.h>
#include "config.h"
#include "app.h"

bool app_runtime_start(Config *config, App *app, const char *execPath);
bool app_runtime_stop(Config *config, App *app);
void app_runtime_note_exit(App *app, int status);
