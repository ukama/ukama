/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * csthreads.h
 */

#include <uuid/uuid.h>

#include "cspace.h"

/* thread/cspace states. */
#define CSPACE_THREAD_STATE_CREATE   0x01
#define CSPACE_THREAD_STATE_ACTIVE   0x02
#define CSPACE_THREAD_STATE_ABORT    0x03
#define CSPACE_THREAD_STATE_DELETED  0x04

/* cspace exit status. */
#define CSPACE_THREAD_EXIT_NORMAL 0x01
#define CSPACE_THREAD_EXIT_TERM   0x02
#define CSPACE_THREAD_EXIT_STOP   0x03

/* Thread for a single contained space */
typedef struct _cspace_thread {

  pthread_t tid;
  pid_t     pid;
  
  uuid_t    uuid;        /* UUID assigned to the space. */
  char      *name;       /* Name of the space -- from config file */
  int       sockets[2];  /* socket pair between parent and cspace */
  int       state;       /* state of the space. */
  int       exit_status; /* only if state is ABORT or DELETED, otherwise is 0*/
  CSpace    *space;      /* cspace associated with this thread. */
} CSpaceThread;

/* List of all contained space threads */
typedef struct _cspace_thread_list {

  CSpaceThread *thread;
  struct _cspace_thread_list *next;
} CSThreadsList;

CSpaceThread *init_cspace_thread(char *name, CSpace *space);
int init_cspace_thread_list(void);
int add_to_cspace_thread_list(CSpaceThread *thread);
void* cspace_thread_start(void *args);
