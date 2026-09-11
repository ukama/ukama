/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
#pragma once

#include "lifecycled.h"

/* The FSM and its pending observations are committed in one atomic file. */
bool state_store_open(LifecycleContext *ctx);
void state_store_close(LifecycleContext *ctx);

bool state_store_load(LifecycleContext *ctx);
bool state_store_save(const LifecycleContext *ctx);
