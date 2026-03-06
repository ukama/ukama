/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "actions.h"

#include <stdlib.h>
#include <string.h>

static char* a_strdup(const char *s) {

    if (!s || !*s) return NULL;
    return strdup(s);
}

void actions_init(ActionQueue *q) {

    if (!q) return;
    q->head = NULL;
    q->tail = NULL;
}

static void action_free(Action *a) {

    if (!a) return;

    free(a->space);
    free(a->name);
    free(a->tag);
    free(a);
}

void actions_free(ActionQueue *q) {

    Action *a;
    Action *n;

    if (!q) return;

    a = q->head;
    while (a) {
        n = a->next;
        action_free(a);
        a = n;
    }

    q->head = NULL;
    q->tail = NULL;
}

Action* action_new(ActionType type, const char *space, const char *name, const char *tag) {

    Action *a;

    a = calloc(1, sizeof(*a));
    if (!a) return NULL;

    a->type  = type;
    a->space = a_strdup(space);
    a->name  = a_strdup(name);
    a->tag   = a_strdup(tag);
    a->next  = NULL;

    return a;
}

bool actions_enqueue(ActionQueue *q, Action *a) {

    if (!q || !a) return false;

    if (!q->head) {
        q->head = a;
        q->tail = a;
        return true;
    }

    q->tail->next = a;
    q->tail = a;
    return true;
}

Action* actions_dequeue(ActionQueue *q) {

    Action *a;

    if (!q || !q->head) return NULL;

    a = q->head;
    q->head = a->next;
    if (!q->head) q->tail = NULL;

    a->next = NULL;
    return a;
}
