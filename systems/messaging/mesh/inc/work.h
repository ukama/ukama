/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef MESH_WORK_H
#define MESH_WORK_H

#include <pthread.h>

#include "mesh.h"

typedef void (*thread_func_t)(char *, void *arg);

/*
 * Link list of object websocket is waiting to handle.
 */
typedef struct work_item_t {

	thread_func_t preFunc;    /* Packet pre-process function, can be NULL */
	thread_func_t postFunc;   /* Packet post-processing function, can be NULL */
	void          *preArgs;   /* Args for pre-function. */
	void          *postArgs;  /* Args for post-funciion. */
	char          *data;       /* Data packet to process. */
	struct work_item_t *next; /* Link to next item in the queue. */
}WorkItem;

/*
 * WorkList - Mutex-ed work items.
 */
typedef struct {

	/* Transmit work list and mutex */
	WorkItem       *first;   /* First item in the TX queue. */
	WorkItem       *last;    /* Pointer to last item in the TX queue. */
	
	pthread_mutex_t mutex;    /* Mutex for insert and remove */
	pthread_cond_t  hasWork;  /* Cond to signal when there is work. */
	int             exit;     /* if TX thread is to exit or exited. */
}WorkList;

/* Functions. */
int add_work_to_queue(WorkList **list, char *data, thread_func_t pre,
					  void *preArgs, thread_func_t post, void *postArgs);
WorkItem *get_work_to_transmit(WorkList *list);
void init_work_list(WorkList **list);
void free_work_item(WorkItem *work);
void free_work_list(WorkList *workList);

#endif /* MESH_WORK_H */
