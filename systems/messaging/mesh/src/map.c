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

	pthread_mutex_init(&map->mutex, NULL);
	pthread_cond_init(&map->hasResp, NULL);

    map->forwardList = (ForwardList *)calloc(1, sizeof(ForwardList));
    if (map->forwardList == NULL) {
        log_error("Error allocating memory: %s", sizeof(ForwardList));
        goto error;
    }
    init_forward_list(&map->forwardList);

	map->nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
    if (map->nodeInfo == NULL) {
        log_error("Error allocating memory: %s", sizeof(NodeInfo));
        goto error;
    }

    map->forwardInst        = *instance;
    map->nodeInfo->nodeID   = strdup(nodeID);
    map->nodeInfo->nodeIP   = strdup(nodeIP);
    map->nodeInfo->nodePort = nodePort;
    map->nodeInfo->meshIP   = strdup(meshIP);
    map->nodeInfo->meshPort = meshPort;
    map->transmit = (WorkList *)calloc(1, sizeof(WorkList));
    map->receive  = (WorkList *)calloc(1, sizeof(WorkList));
    if (map->transmit) init_work_list(&map->transmit);
    if (map->receive)  init_work_list(&map->receive);
    if (!map->nodeInfo->nodeID || !map->nodeInfo->nodeIP ||
        !map->nodeInfo->meshIP || !map->transmit || !map->receive) {
        goto error;
    }

	map->next = NULL;

	return map;

error:
    free_map_item(map);
    return NULL;
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

    pthread_mutex_destroy(&map->mutex);
    pthread_cond_destroy(&map->hasResp);

    free(map);
}

MapItem *is_existing_item(MapTable *table, char *nodeID) {

	MapItem *item;

	if (table == NULL || nodeID == NULL) {
		return NULL;
	}

	pthread_mutex_lock(&table->mutex);
	for (item=table->first; item; item=item->next) {
		if (strcmp(item->nodeInfo->nodeID, nodeID) == 0) {
			item->references++;
			break;
		}
	}
	pthread_mutex_unlock(&table->mutex);

	return item;
}

MapItem *is_existing_item_by_port(MapTable *table, int port) {

	MapItem *item = NULL;

	if (table == NULL || port == 0) {
		return NULL;
	}

    pthread_mutex_lock(&table->mutex);
    item = table->first;
    if (item && item->online) {
        item->references++;
    } else {
        item = NULL;
    }
    pthread_mutex_unlock(&table->mutex);

    return item;
}

MapItem *add_map_to_table(MapTable **table,
                          char *nodeID,
                          UInst **instance,
                          char *nodeIP, int nodePort,
                          char *meshIP, int meshPort) {

	MapItem *map=NULL;

	if (*table == NULL || nodeID == NULL)
		return NULL;

    pthread_mutex_lock(&(*table)->mutex);
    for (map=(*table)->first; map; map=map->next) {
        if (strcmp(map->nodeInfo->nodeID, nodeID) == 0) {
            /* Its close callback owns retirement; the client will retry. */
            if (map->wsManager) {
                ulfius_websocket_send_close_signal(map->wsManager);
            }
            pthread_mutex_unlock(&(*table)->mutex);
            return NULL;
        }
    }

	map = create_map_item(nodeID, instance,
                          nodeIP, nodePort,
                          meshIP, meshPort);
	if (map == NULL) {
		pthread_mutex_unlock(&(*table)->mutex);
		return NULL;
	}

    map->references = 2; /* Table and websocket callback ownership. */

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

void remove_map_item_from_table(MapTable *table, MapItem *map) {

    MapItem *current, *previous;

    pthread_mutex_lock(&table->mutex);

    current  = table->first;
    previous = NULL;

    while (current != NULL) {
        if (current == map) {
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

            current->next = NULL;
            pthread_mutex_unlock(&table->mutex);
            release_map_item(table, current);

            return;
        }

        previous = current;
        current = current->next;
    }

    pthread_mutex_unlock(&table->mutex);
}

void release_map_item(MapTable *table, MapItem *map) {

    int unused;

    if (map == NULL) return;

    pthread_mutex_lock(&table->mutex);
    unused = (--map->references == 0);
    pthread_mutex_unlock(&table->mutex);

    if (unused) free_map_item(map);
}
