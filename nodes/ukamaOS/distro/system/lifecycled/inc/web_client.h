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
#include "fsm.h"
#include "lifecycled.h"

bool starter_client_get_status(const Config *config,
                               StarterSnapshot *snapshot);

bool notify_client_send_event(const Config *config,
                              const LifecycleEvent *event);

