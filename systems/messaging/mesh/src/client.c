/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <pthread.h>
#include <string.h>

#include "client.h"
#include "mesh.h"

void init_forward_list(ForwardList **list) {

	(*list)->first = NULL;
	(*list)->last  = NULL;

	pthread_mutex_init(&(*list)->mutex, NULL);
}

void free_forward_item(Forward *item) {

	if (!item) return;

    free(item->uuid);
	free(item);
}

Forward *is_existing_item_in_list(ForwardList *list, char *uuid) {

	Forward *item;

	if (list == NULL || uuid == NULL) return NULL;
	if (list->first == NULL)          return NULL;

	for (item=list->first; item; item=item->next) {
		if (strcmp(item->uuid, uuid) == 0) {
			return item;
		}
	}

	return NULL;
}

Forward *add_client_to_list(ForwardList **list,
                            char *uuid) {

	Forward *item=NULL;

	if (*list == NULL || uuid == NULL)
        return NULL;

	item = is_existing_item_in_list(*list, uuid);
	if (item != NULL) {
		return item;
	}

    item = (Forward *)calloc(1, sizeof(Forward));
    if (item == NULL) {
        log_error("Unable to allocate memory of size: %d",
                  sizeof(Forward));
        return NULL;
    }
    item->uuid = strdup(uuid);

	pthread_mutex_lock(&(*list)->mutex);

	if ((*list)->first == NULL) {
		(*list)->first = item;
		(*list)->last  = item;
	} else {
		(*list)->last->next = item;
	}

	/* Update pointer to last entry. */
	(*list)->last = item;
	(*list)->last->next = NULL;

	/* Unlock */
	pthread_mutex_unlock(&((*list)->mutex));

	log_debug("Added new mapping entry in forward list with uuid: %s", uuid);

	return item;
}

void remove_item_from_list(ForwardList *list, char *uuid) {

    Forward *current, *previous;

    pthread_mutex_lock(&list->mutex);

    current  = list->first;
    previous = NULL;

    while (current != NULL) {
        if (strcmp(current->uuid, uuid) == 0) {
            if (previous != NULL) {
                previous->next = current->next;
                if (current == list->last) {
                    list->last = previous;
                }
            } else {
                list->first = current->next;
                if (current == list->last) {
                    list->last = NULL;
                }
            }

            pthread_mutex_unlock(&list->mutex);
            pthread_mutex_destroy(&current->mutex);
            pthread_cond_destroy(&current->hasData);
            free_forward_item(current);

            return;
        }

        previous = current;
        current = current->next;
    }

    pthread_mutex_unlock(&list->mutex);
}
