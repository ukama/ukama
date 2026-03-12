/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>

typedef enum {
    ACTION_NONE = 0,
    ACTION_RUN_BOOT,
    ACTION_RUN_ALL,
    ACTION_TERMINATE_APP,
    ACTION_UPDATE_APP
} ActionType;

typedef struct Action {
    ActionType type;
    char *space;
    char *name;
    char *tag;
    struct Action *next;
} Action;

typedef struct ActionQueue {
    Action *head;
    Action *tail;
} ActionQueue;

void actions_init(ActionQueue *q);
void actions_free(ActionQueue *q);
bool actions_enqueue(ActionQueue *q, Action *a);
Action* actions_dequeue(ActionQueue *q);
Action* action_new(ActionType type, const char *space, const char *name, const char *tag);
