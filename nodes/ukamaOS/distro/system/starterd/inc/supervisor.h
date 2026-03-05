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
#include "space.h"
#include "actions.h"

typedef struct Supervisor Supervisor;

Supervisor* supervisor_start(Config *config, Space *spaceList, ActionQueue *queue);
void supervisor_stop(Supervisor *s);

bool supervisor_signal(Supervisor *s);
