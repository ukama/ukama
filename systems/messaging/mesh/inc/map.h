/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef MESH_MAP_H
#define MESH_MAP_H

#include <pthread.h>

#include "mesh.h"
#include "work.h"
#include "client.h"

/* keep track of ip:port to UUID mapping along with various mutex and
 * conditional variable. The thread will wait until its unlocked by the
 * response or timeout.
 */
typedef struct map_item_t {

    NodeInfo  *nodeInfo;
    WorkList  *transmit;
    WorkList  *receive;
    WSManager *wsManager;
    void      *configData;
    UInst     *forwardInst;

    ForwardList *forwardList;  /* services list */

    pthread_mutex_t   mutex;   /* Client thread waiting on response */
    pthread_cond_t    hasResp; /* Conditional wait for response */

    int               code;
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
MapItem *is_existing_item(MapTable *table, char *nodeID);
MapItem *is_existing_item_by_port(MapTable *table, int port);
MapItem *add_map_to_table(MapTable **table,
                          char *nodeID, UInst **instance,
                          char *nodeIP, int nodePort,
                          char *meshIP, int meshPort);
#endif /* MESH_MAP_H */
