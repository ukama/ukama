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
#include "mesh.h"

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
static MapItem *create_map_item(char *nodeID) {

	MapItem *map=NULL;

	/* Sanity check */
	if (nodeID == NULL)
		return NULL;

	map = (MapItem *)calloc(1, sizeof(MapItem));
	if (!map) {
		log_error("Error allocating memory: %d", sizeof(MapItem));
		return NULL;
	}

	map->nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
    if (map->nodeInfo == NULL) {
        log_error("Error allocating memory: %s", sizeof(NodeInfo));
        free(map);
        return NULL;
    }

    map->nodeInfo->nodeID = strdup(nodeID);
    map->transmit = NULL;
    map->receive  = NULL;

	pthread_mutex_init(&map->mutex, NULL);
	pthread_cond_init(&map->hasResp, NULL);

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

    if (map->nodeInfo) {
        free(map->nodeInfo->nodeID);
        free(map->nodeInfo);
    }
	free(map);
}

/*
 * is_existing_item --
 *
 */
MapItem *is_existing_item(MapTable *table, char *nodeID) {

	MapItem *item;

	if (table == NULL && nodeID == NULL) {
		return NULL;
	}

	/* Is empty. */
	if (table->first == NULL) {
		return NULL;
	}

	for (item=table->first; item; item=item->next) {
		if (strcmp(item->nodeInfo->nodeID, nodeID) == 0) {
			return item;
		}
	}

	return NULL;
}

/*
 * add_map_to_table --
 *
 */
MapItem *add_map_to_table(MapTable **table, char *nodeID) {

	MapItem *map=NULL;

	if (*table == NULL || nodeID == NULL)
		return NULL;

	/* An existing mapping? */
	map = is_existing_item(*table, nodeID);
	if (map != NULL) {
		return map;
	}

	map = create_map_item(nodeID);
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

	log_debug("Added new mapping entry in the table. NodeID: %s", nodeID);

	return map;
}

/*
 * lookup_item -- find the matching item by nodeID
 *
 */
MapItem *lookup_item(MapTable *table, char *nodeID) {

	MapItem *item;

	if (table == NULL && nodeID == NULL) {
		return NULL;
	}

	/* Is empty. */
	if (table->first == NULL) {
		return NULL;
	}

	for (item=table->first; item; item=item->next) {
        if (strcmp(item->nodeInfo->nodeID, nodeID) == 0) {
			return item;
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
