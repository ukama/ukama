/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef MESH_MAP_H
#define MESH_MAP_H

#include <pthread.h>

#include "mesh.h"
#include "work.h"

/* keep track of ip:port to UUID mapping along with various mutex and
 * conditional variable. The thread will wait until its unlocked by the
 * response or timeout.
 */
typedef struct map_item_t {

    NodeInfo *nodeInfo;
    WorkList *transmit;
    WorkList *receive;
    void     *configData;

    pthread_mutex_t   mutex;   /* Client thread waiting on response
								* This mutex is released by websocket */
    pthread_cond_t    hasResp; /* Conditional wait for response */

    int               size;    /* size of data packet. */
    void              *data;   /* response data recevied. */

	struct map_item_t *next;   /* Link to next item in the table */
} MapItem;

/*
 * MapTable - Mutex-ed table elements. 
 */
typedef struct {

	MapItem *first;        /* First item in the mapping table */
	MapItem *last;         /* Last item in the mapping table */

	pthread_mutex_t mutex;    /* Mutex for insert and remove */
} MapTable;

/* Functions */
void init_map_table(MapTable **table);
void free_map_item(MapItem *map);
void remove_map_item_from_table(MapTable *table, char *nodeID);
MapItem *add_map_to_table(MapTable **table, char *nodeID);
MapItem *is_existing_item(MapTable *table, char *nodeID);

#endif /* MESH_MAP_H */
