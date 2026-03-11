/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include "config.h"
#include "app.h"

bool installer_ensure_installed(Config *config,
                                App *app,
                                const char *hub);
bool installer_switch_current(Config *config, App *app);
bool installer_revert_to_last_good(Config *config, App *app);
