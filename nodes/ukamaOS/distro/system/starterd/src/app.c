/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "app.h"
#include "space.h"

#include <string.h>

App* app_find(Space *spaceList, const char *space, const char *name) {

    Space *s;
    App *a;

    if (!space || !name) return NULL;

    s = space_find(spaceList, space);
    if (!s) return NULL;

    a = s->appList;
    while (a) {
        if (a->name && strcmp(a->name, name) == 0) return a;
        a = a->next;
    }

    return NULL;
}
