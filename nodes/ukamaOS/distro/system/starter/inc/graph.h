/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef GRAPH_H
#define GRAPH_H

#define MAX_PROGRAMS 100
#define MAX_NAME_LEN 50

#define STATE_NONE "none"
#define STATE_DONE "done"
#define STATE_RUN  "run"

typedef struct Program {

    char   name[MAX_NAME_LEN];
    int    index;
    char   state[5]; /* none, run, done */
} Program;

typedef struct AdjListNode {

    int programIndex;

    struct AdjListNode* next;
} AdjListNode;

typedef struct Graph {

    int          numVertices;
    int          *visited;
    int          programCount;
    AdjListNode  **adjLists;
    Program      programs[MAX_PROGRAMS];
} Graph;

#endif /* GRAPH_H */
