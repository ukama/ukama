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
#include <uuid/uuid.h>

#include "mesh.h"

/* keep track of ip:port to UUID mapping along with various mutex and
 * conditional variable. The thread will wait until its unlocked by the
 * response or timeout.
 */
typedef struct map_item_t {

	unsigned short    port;    /* Client port number */
	char              *ip;     /* Client IP in sting format */
	uuid_t            uuid;    /* Mapped UUID */
	pthread_mutex_t   mutex;   /* Client thread waiting on response. 
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

/* Functions. */
void init_map_table(MapTable **table);
void destroy_map_item(MapItem *map);
MapItem *add_map_to_table(MapTable **table, char *ip, unsigned short port);
MapItem *lookup_item(MapTable *table, uuid_t uuid);

#endif /* MESH_MAP_H */
