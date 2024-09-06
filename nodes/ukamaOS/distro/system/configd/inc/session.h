/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef INC_SESSION_H_
#define INC_SESSION_H_

#include "config_macros.h"
#include "usys_types.h"

typedef enum {

    CONFIG_ADD    = 1,
    CONFIG_DELETE = 2,
    CONFIG_UPDATE = 3
} Reason;

typedef struct {

    char *name;
    char *version;
    char *fileName;
    char *data;
    int  reason; /* added, deleted, updated */
} AppData;

typedef struct {

    char *fileName;
    char *app;
    char *data;
    char *version;
    int timestamp;
    int reason;
    int fileCount;
} SessionData;

typedef struct {

    AppData  *apps[MAX_APPS];
    int      timestamp;
    int      expectedCount;
    int      receviedCount;
} ConfigSession;

#endif
