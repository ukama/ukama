/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef MESH_MAP_H
#define MESH_MAP_H

#include <pthread.h>

#include "mesh.h"
#include "work.h"

typedef struct map_item_t {

    char              *serviceName;
    char              *servicePort;
    char              *uuid;
    WorkList          *transmit;
    WorkList          *receive;
    void              *configData;

    pthread_mutex_t   mutex;   /* Client thread waiting on response */
	pthread_cond_t    hasResp; /* Conditional wait for response */

    int               code;
	int               size;    /* size of data packet. */
	void              *data;   /* response data recevied. */

	struct map_item_t *next;   /* Link to next item in the table */
} MapItem;

typedef struct {

	MapItem *first;        /* First item in the mapping table */
	MapItem *last;         /* Last item in the mapping table */

	pthread_mutex_t mutex;    /* Mutex for insert and remove */
} MapTable;

/* Functions. */
void init_map_table(MapTable **table);
void destroy_map_item(MapItem *map);
MapItem *is_existing_item(MapTable *table, char *uuid);
void remove_map_item_from_table(MapTable *table, char *uuid);
MapItem *add_map_to_table(MapTable **table, char *name, char *port, char *uuid);

#endif /* MESH_MAP_H */
