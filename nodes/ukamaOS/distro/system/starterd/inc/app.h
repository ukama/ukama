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
#include <sys/types.h>
#include "starterd.h"

#include "space.h"

typedef struct App {
    char *space;
    char *name;
    char *tag;

    char *cmd;
    char **argv;
    int argc;

    char **envp;
    int envc;

    char *workdir;

    int port;

    pid_t pid;
    pid_t pgid;

    AppState state;
    InstallState installState;

    int lastExitCode;
    int lastExitSignal;
    time_t lastStartTime;
    time_t lastExitTime;

    int restartCount;
    time_t restartWindowStart;
    int nextBackoffSec;

    char *lastGoodTag;

    struct App *next;
} App;

App* app_find(Space *spaceList, const char *space, const char *name);
