/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>

#include "map.h"
#include "work.h"
#include "mesh.h"

void init_map_table(MapTable **table) {

	(*table)->first = NULL;
	(*table)->last  = NULL;

	pthread_mutex_init(&(*table)->mutex, NULL);
}

static MapItem *create_map_item(char *nodeID,
                                UInst **instance,
                                char *nodeIP, int nodePort,
                                char *meshIP, int meshPort) {

	MapItem *map=NULL;

	/* Sanity check */
	if (nodeID == NULL)
		return NULL;

	map = (MapItem *)calloc(1, sizeof(MapItem));
	if (!map) {
		log_error("Error allocating memory: %d", sizeof(MapItem));
		return NULL;
	}

    map->forwardList = (ForwardList *)calloc(1, sizeof(ForwardList));
    if (map->forwardList == NULL) {
        log_error("Error allocating memory: %s", sizeof(ForwardList));
        free(map);
        return NULL;
    }

	map->nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
    if (map->nodeInfo == NULL) {
        log_error("Error allocating memory: %s", sizeof(NodeInfo));
        free(map);
        return NULL;
    }

    map->forwardInst        = *instance;
    map->nodeInfo->nodeID   = strdup(nodeID);
    map->nodeInfo->nodeIP   = strdup(nodeIP);
    map->nodeInfo->nodePort = nodePort;
    map->nodeInfo->meshIP   = strdup(meshIP);
    map->nodeInfo->meshPort = meshPort;
    map->transmit           = NULL;
    map->receive            = NULL;

	pthread_mutex_init(&map->mutex, NULL);
	pthread_cond_init(&map->hasResp, NULL);

	map->next = NULL;

	return map;
}

void free_map_item(MapItem *map) {

	if (!map) {
		return;
	}

    if (map->nodeInfo) {
        free(map->nodeInfo->nodeID);
        free(map->nodeInfo->nodeIP);
        free(map->nodeInfo->meshIP);
        free(map->nodeInfo);
    }

    free_work_list(map->transmit);
    free_work_list(map->receive);
    free_forward_list(map->forwardList);

    free(map);
}

MapItem *is_existing_item(MapTable *table, char *nodeID) {

	MapItem *item;

	if (table == NULL || nodeID == NULL) {
		return NULL;
	}

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

MapItem *is_existing_item_by_port(MapTable *table, int port) {

	MapItem *item = NULL;

	if (table == NULL || port == 0) {
		return NULL;
	}

	if (table->first == NULL) {
		return NULL;
	}

    return table->first;
}

MapItem *add_map_to_table(MapTable **table,
                          char *nodeID,
                          UInst **instance,
                          char *nodeIP, int nodePort,
                          char *meshIP, int meshPort) {

	MapItem *map=NULL;

	if (*table == NULL || nodeID == NULL)
		return NULL;

	/* An existing mapping? */
	map = is_existing_item(*table, nodeID);
	if (map != NULL) {
		return map;
	}

	map = create_map_item(nodeID, instance,
                          nodeIP, nodePort,
                          meshIP, meshPort);
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

void remove_map_item_from_table(MapTable *table, char *nodeID) {

    MapItem *current, *previous;

    pthread_mutex_lock(&table->mutex);

    current  = table->first;
    previous = NULL;

    while (current != NULL) {
        if (strcmp(current->nodeInfo->nodeID, nodeID) == 0) {
            if (previous != NULL) {
                previous->next = current->next;
                if (current == table->last) {
                    table->last = previous;
                }
            } else {
                table->first = current->next;
                if (current == table->last) {
                    table->last = NULL;
                }
            }

            pthread_mutex_unlock(&table->mutex);
            pthread_mutex_destroy(&current->mutex);
            pthread_cond_destroy(&current->hasResp);
            free_map_item(current);

            return;
        }

        previous = current;
        current = current->next;
        free(previous);
    }

    pthread_mutex_unlock(&table->mutex);
}
