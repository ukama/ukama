/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>

#include "usys_log.h"

#include "map.h"
#include "mesh.h"

void init_map_table(MapTable **table) {

	(*table)->first = NULL;
	(*table)->last  = NULL;

	pthread_mutex_init(&(*table)->mutex, NULL);
}

static MapItem *create_map_item(char *name, char *port) {

	MapItem *map;

	/* Sanity check */
	if (name == NULL || port == NULL) return NULL;

	map = (MapItem *)malloc(sizeof(MapItem));
	if (map == NULL) {
		usys_log_error("Error allocating memory: %d", sizeof(MapItem));
		return NULL;
	}

    map->serviceInfo = (ServiceInfo *)calloc(1, sizeof(ServiceInfo));
    if (map->serviceInfo == NULL) {
        usys_log_error("Error allocating memory: %d", sizeof(ServiceInfo));
        free(map);
        return NULL;
    }

    map->serviceInfo->name = strdup(name);
    map->serviceInfo->port = strdup(port);
    map->transmit = NULL;
    map->receive  = NULL;

	pthread_mutex_init(&map->mutex, NULL);
	pthread_cond_init(&map->hasResp, NULL);

	map->next = NULL;

	return map;
}

void free_map_item(MapItem *map) {

	if (!map) {
		return;
	}

    if (map->serviceInfo) {
        free(map->serviceInfo->name);
        free(map->serviceInfo->port);
        free(map->serviceInfo);
    }

	free(map);
}

MapItem *is_existing_item(MapTable *table, char *name, char *port) {

	MapItem *item=NULL;

	if (table == NULL || name == NULL || port == NULL) {
		return NULL;
	}

	/* is empty */
	if (table->first == NULL) {
		return NULL;
	}

	for (item = table->first; item; item=item->next) {
        if (strcmp(item->serviceInfo->name, name) == 0 &&
            strcmp(item->serviceInfo->port, port) == 0) {
            return item;
        }
    }

    return NULL;
}

MapItem *add_map_to_table(MapTable **table, char *name, char *port) {

	MapItem *map=NULL;

    if (*table == NULL || name == NULL || port == NULL) return NULL;

    /* An existing mapping? */
    map = is_existing_item(*table, name, port);
    if (map != NULL) {
        return map;
    }

    map = create_map_item(name, port);
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

    log_debug("Added new mapping entry in the table. Name: %s port: %s",
              name, port);

    return map;
}

void remove_map_item_from_table(MapTable *table, char *name, char *port) {

    MapItem *current=NULL, *previous=NULL;

    pthread_mutex_lock(&table->mutex);

    current  = table->first;
    previous = NULL;

    while (current != NULL) {
        if (strcmp(current->serviceInfo->name, name) == 0 &&
            strcmp(current->serviceInfo->port, port) == 0) {
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
    }

    pthread_mutex_unlock(&table->mutex);
}
