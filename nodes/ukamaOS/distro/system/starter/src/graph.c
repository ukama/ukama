/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "starter.h"
#include "graph.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"

static int get_program_index(Graph* graph, char* name) {

    int i;

    for (i = 0; i < graph->programCount; i++) {
        if (strcmp(graph->programs[i].name, name) == 0) {
            return i;
        }
    }

    strcpy(graph->programs[graph->programCount].name, name);
    strcpy(graph->programs[graph->programCount].state, STATE_NONE);

    return graph->programCount++;
}

static AdjListNode* create_node(int programIndex) {

    AdjListNode* newNode = NULL;

    newNode = (AdjListNode*)calloc(1, sizeof(AdjListNode));
    if (newNode == NULL) {
        usys_log_error("Unable to allocate memory of size: %s",
                       sizeof(AdjListNode));
        return NULL;
    }

    newNode->programIndex = programIndex;
    newNode->next         = NULL;

    return newNode;
}

Graph* create_graph(int vertices) {

    Graph *graph = NULL;

    graph = (Graph*)calloc(1, sizeof(Graph));
    if (graph == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                       sizeof(Graph));
        return NULL;
    }

    graph->programCount = 0;
    graph->numVertices  = vertices;
    graph->adjLists     = (AdjListNode**)calloc(1, vertices * sizeof(AdjListNode*));
    graph->visited      = (int*)calloc(vertices, sizeof(int));

    if (graph->adjLists == NULL ||
        graph->visited  == NULL) {
        usys_log_error("Unable to allocate memory of sizes: %d %s",
                       vertices * sizeof(AdjListNode*),
                       sizeof(int));

        usys_free(graph->adjLists);
        usys_free(graph->visited);

        return NULL;
    }

    return graph;
}

void add_program(Graph* graph, char* name, const char* state) {

    int index;

    index = get_program_index(graph, name);
    strcpy(graph->programs[index].state, state);
}

void add_edge(Graph* graph, char* src, char* dest) {

    int srcIndex, destIndex;
    AdjListNode *newNode;

    srcIndex  = get_program_index(graph, src);
    destIndex = get_program_index(graph, dest);

    newNode = create_node(destIndex);
    if (newNode == NULL) return;

    newNode->next             = graph->adjLists[srcIndex];
    graph->adjLists[srcIndex] = newNode;
}

void dfs(Graph* graph, int vertex, int* topologicalOrder, int* index) {

    AdjListNode *temp;

    graph->visited[vertex] = 1;
    temp = graph->adjLists[vertex];

    while (temp != NULL) {

        int connectedVertex = temp->programIndex;

        if (graph->visited[connectedVertex] == 0) {
            dfs(graph, connectedVertex, topologicalOrder, index);
        }
        temp = temp->next;
    }

    topologicalOrder[(*index)++] = vertex;
}

void topological_sort(Graph* graph, int* topologicalOrder) {

    int index = 0, i = 0;

    for (i = 0; i < graph->numVertices; i++) {
        graph->visited[i] = 0;
    }

    for (i = 0; i < graph->numVertices; i++) {
        if (!graph->visited[i] && strcmp(graph->programs[i].state, STATE_RUN) == 0) {
            dfs(graph, i, topologicalOrder, &index);
        }
    }
}

void print_execution_order(Graph* graph, const int* topologicalOrder) {

    int i;

    usys_log_debug("Capps execution order: ");

    for (i = 0; i < graph->numVertices; i++) {
        if (strcmp(graph->programs[i].state, STATE_NONE) == 0) {
            usys_log_debug("%s ", graph->programs[i].name);
        }
    }

    for (i = graph->numVertices - 1; i >= 0; i--) {
        if (topologicalOrder[i] != -1) {
            usys_log_debug("%s ", graph->programs[topologicalOrder[i]].name);
        }
    }
}
