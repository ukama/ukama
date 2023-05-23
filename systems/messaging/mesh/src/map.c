/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <pthread.h>
#include <string.h>

#include "map.h"

/*
 * init_map_table -- 
 */
void init_map_table(MapTable **table) {

	(*table)->first = NULL;
	(*table)->last  = NULL;

	pthread_mutex_init(&(*table)->mutex, NULL);
}

/*
 * create_map_item --
 *
 */
static MapItem *create_map_item(char *ip, unsigned short port) { 

	MapItem *map;

	/* Sanity check */
	if (ip == NULL && !port)
		return NULL;

	map = (MapItem *)malloc(sizeof(MapItem));
	if (!map) {
		log_error("Error allocating memory: %d", sizeof(MapItem));
		return NULL;
	}

	map->port = port;
	map->ip   = strdup(ip);

	pthread_mutex_init(&map->mutex, NULL);
	pthread_cond_init(&map->hasResp, NULL);

	/* Assign a new NodeID. */
    map->nodeID = strdup("uk-xx1234-tnode-m0-xxxx");
	map->next = NULL;

	return map;
}

/*
 * destroy_work_item --
 *
 */
void destroy_map_item(MapItem *map) {

	if (!map) {
		return;
	}

	free(map->ip);
	free(map);
}

/*
 * is_existing_item -- Check if the ip:port mapping already exist. 
 *
 */
static MapItem *is_existing_item(MapTable *table, char *ip,
								 unsigned short port) {

	MapItem *items;

	if (table == NULL && ip == NULL && port == 0) {
		return NULL;
	}

	/* Is empty. */
	if (table->first == NULL) {
		return NULL;
	}

	for (items = table->first; items; items=items->next) {
		if (strcmp(items->ip, ip)==0 && port == items->port && /* Match found */
            items->nodeID != NULL) { /* have valid NodeID. */
			return items;
		}
	}

	return NULL;
}

/*
 * add_map_to_table -- Add new ip:port into mapping table against NodeID.
 *
 */
MapItem *add_map_to_table(MapTable **table, char *ip, unsigned short port) {

	MapItem *map=NULL;

	if (ip == NULL && *table == NULL && !port)
		return NULL;

	/* An existing mapping? */
	map = is_existing_item(*table, ip, port);
	if (map != NULL) {
		return map;
	}

	map = create_map_item(ip, port);
	if (map == NULL) {
		return NULL;
	}

	/* Try to get lock. */
	pthread_mutex_lock(&(*table)->mutex);

	/* Got the lock. Add to the list and unlock. */
	if ((*table)->first == NULL) {
		(*table)->first = map;
		(*table)->last  = map;
	} else {
		(*table)->last->next = map;
	}

	/* Update pointer to last entry. */
	(*table)->last = map;
	(*table)->last->next = NULL;

	/* Unlock */
	pthread_mutex_unlock(&((*table)->mutex));

	log_debug("Added new mapping entry in the table. IP: %s port: %d",
			  ip, port);

	return map;
}

/*
 * lookup_item -- find the matching item by nodeID
 *
 */
MapItem *lookup_item(MapTable *table, char *nodeID) {

	MapItem *items;

	if (table == NULL && nodeID == NULL) {
		return NULL;
	}

	/* Is empty. */
	if (table->first == NULL) {
		return NULL;
	}

	for (items = table->first; items; items=items->next) {
        if (strcmp(nodeID, items->nodeID)==0) {
			return items;
		}
	}

	return NULL;
}

/*
 * remove_item -- remove the matching item from the table and free()
 *
 */
void remove_item(MapTable **table, char *nodeID) {


}
