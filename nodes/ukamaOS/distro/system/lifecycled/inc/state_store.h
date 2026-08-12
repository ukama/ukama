/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

#include "fsm.h"

bool state_store_load(const char *path,
                      const char *bootId,
                      LifecycleFsm *fsm);

bool state_store_save(const char *path,
                      const char *bootId,
                      const LifecycleFsm *fsm);

