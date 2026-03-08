/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

typedef struct App App;

typedef struct Space {
    char *name;
    App *appList;
    struct Space *next;
} Space;

Space* space_find(Space *spaceList, const char *name);
