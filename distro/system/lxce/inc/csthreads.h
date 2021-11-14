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

#ifndef CSTHREAD_H
#define CSTHREAD_H

#include <uuid/uuid.h>

#include "cspace.h"
#include "capp_packet.h"

/* thread/cspace states. */
#define CSPACE_THREAD_STATE_CREATE   0x01
#define CSPACE_THREAD_STATE_ACTIVE   0x02
#define CSPACE_THREAD_STATE_ABORT    0x03
#define CSPACE_THREAD_STATE_DELETED  0x04

/* cspace exit status. */
#define CSPACE_THREAD_EXIT_NORMAL 0x01
#define CSPACE_THREAD_EXIT_TERM   0x02
#define CSPACE_THREAD_EXIT_STOP   0x03

#define CSPACE_COND_WAIT 1 /* Cond wait time in sec. */

#define CSPACE_MEMFILE "lxce_memfile"

#define CSPACE_MAX_BUFFER 1024

#define CSPACE_READ_ERROR   1
#define CSPACE_READ_TIMEOUT 2


/* Logging request/response from thread to cspace */
typedef struct action_list_t {

  int seqno;    /* sequence id for request/response mapping */
  int state;    /* state of this action. */

  char *cmd;    /* CAPP_CMD_* */
  char *params; /* Parameters for the command */
  char *resp;   /* Response if done otherwise is NULL */

  struct action_list_t *next; /* Next one */
} ActionList;

/* Shared memory related info. */
typedef struct {

  char *memFile;
  int  shmId;
} ShMemInfo;

/* Shared memory between thread and parent process */
typedef struct {

  PacketList *txList; /* List of TX packets */
  PacketList *rxList; /* List of RX packets */

  /* Mutex for TX and RX packets */
  pthread_mutex_t txMutex;
  pthread_mutex_t rxMutex;
  pthread_cond_t  hasTX;
  pthread_cond_t  hasRX;
} ThreadShMem;

/* Thread for a single contained space */
typedef struct cspace_thread_t {

  pthread_t tid;
  pid_t     pid;
  
  uuid_t    uuid;        /* UUID assigned to the space. */
  char      *name;       /* Name of the space -- from config file */
  int       sockets[2];  /* socket pair between parent and cspace */
  int       state;       /* state of the space. */
  int       exit_status; /* only if state is ABORT or DELETED, otherwise is 0*/
  CSpace    *space;      /* cspace associated with this thread. */

  ActionList  *actionList; /* On-going action associated with this thread */
  ThreadShMem *shMem;      /* shared between parent and thread */
} CSpaceThread;

/* List of all contained space threads */
typedef struct cspace_thread_list_t {

  CSpaceThread *thread;    /* Thread def for the cspace */
  ThreadShMem  *shMem;     /* Shared memory between process and thread */
  ShMemInfo    *shmemInfo; /* shared memory related variables */

  struct cspace_thread_list_t *next;
} CSThreadsList;

CSpaceThread *init_cspace_thread(char *name, CSpace *space);
int init_cspace_thread_list(void);
int add_to_cspace_thread_list(CSpaceThread *thread);
void* cspace_thread_start(void *args);

#endif /* CSTHREAD_H */
