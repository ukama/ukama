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

#include "work.h"

/*
 * init_work_list -- 
 */
void init_work_list(WorkList **list) {

	(*list)->first = NULL;
	(*list)->last  = NULL;

	pthread_mutex_init(&(*list)->mutex, NULL);
	pthread_cond_init(&(*list)->hasWork, NULL);

	(*list)->exit = FALSE;
}

/*
 * create_work_item --
 *
 */
static WorkItem *create_work_item(char *data, thread_func_t pre, void *preArgs,
								  thread_func_t post, void *postArgs) {

	WorkItem *work;

	/* Sanity check */
	if (data == NULL)
		return NULL;

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

/*
 * destroy_work_item --
 *
 */
void destroy_work_item(WorkItem *work) {

	if (!work) {
		return;
	}

	free(work->data);
	free(work);
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
