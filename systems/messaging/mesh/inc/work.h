/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
void destroy_work_item(WorkItem *work);

#endif /* MESH_WORK_H */
