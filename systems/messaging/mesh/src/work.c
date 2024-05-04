/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>

#include "work.h"

void init_work_list(WorkList **list) {

	(*list)->first = NULL;
	(*list)->last  = NULL;

	pthread_mutex_init(&(*list)->mutex, NULL);
	pthread_cond_init(&(*list)->hasWork, NULL);

	(*list)->exit = FALSE;
}

static WorkItem *create_work_item(char *data, thread_func_t pre, void *preArgs,
								  thread_func_t post, void *postArgs) {

	WorkItem *work;

	if (data == NULL) return NULL;

	work = (WorkItem *)malloc(sizeof(WorkItem));
	if (!work) {
		log_error("Error allocating memory: %d", sizeof(WorkItem));
		return NULL;
	}

	work->preFunc  = pre;
	work->postFunc = post;
	work->preArgs  = preArgs;
	work->postArgs = postArgs;

	work->data = strdup(data);
	work->next = NULL;

	return work;
}

void free_work_item(WorkItem *work) {

	if (work == NULL) return;

	free(work->data);
	free(work);
}

void free_work_list(WorkList *workList) {

    WorkItem *current, *temp;

    if (workList == NULL) return;

    pthread_mutex_lock(&workList->mutex);

    current = workList->first;
    while (current != NULL) {
        temp = current->next;
        free(current->data);
        free(current);
        current = temp;
    }

    workList->first = NULL;
    workList->last  = NULL;
    pthread_mutex_unlock(&workList->mutex);

    pthread_mutex_destroy(&workList->mutex);
    pthread_cond_destroy(&workList->hasWork);

    free(workList);
}

/*
 * add_work_to_queue -- blocking call to add the work item to the queue
 *                      for websocket.
 *
 */
int add_work_to_queue(WorkList **list, char *data, thread_func_t pre,
					  void *preArgs, thread_func_t post, void *postArgs) {

	WorkItem *work=NULL;

	if (data == NULL && *list == NULL) return FALSE;

	work = create_work_item(data, pre, preArgs, post, postArgs);
	if (work == NULL) return FALSE;

	/* Try to get lock. */
	pthread_mutex_lock(&(*list)->mutex);

	/* Got the lock. Add to the list and unlock. */
	if ((*list)->first == NULL) {
		(*list)->first = work;
		(*list)->last  = work;
	} else {
		(*list)->last->next = work;
	}

	/* Update pointer to last entry. */
	(*list)->last = work;
	(*list)->last->next = NULL;

	/* Broadcast new work item is available in the queue. */
	pthread_cond_broadcast(&((*list)->hasWork));

	/* Unlock */
	pthread_mutex_unlock(&((*list)->mutex));
    log_debug("Work added on the queue. Len: %d Data: %s", strlen(data), data);

	return TRUE;
}

/*
 * get_work_to_transmit -- remove the first work item from the queue.
 *                         callee is responsible for memory free.
 *
 */
WorkItem *get_work_to_transmit(WorkList *list){

	WorkItem *item=NULL;

	/* Is empty. */
	if (list->first == NULL) {
		return NULL;
	}

	/* Is the only item. i.e., first == last */
	if (list->first == list->last) {
		item = list->first;
		list->first = NULL;
		list->last  = NULL;
	} else { /* General case. */
		/* FIFO, always return the first entry in */
		item = list->first;
		list->first = item->next;
		item->next = NULL;
	}

	return item;
}
