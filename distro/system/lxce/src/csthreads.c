/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* Functions related to CSpace threads and threads list
 */

#include <uuid/uuid.h>
#include <sys/wait.h>
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/ipc.h>
#include <sys/types.h>
#include <sys/shm.h>
#include <string.h>
#include <errno.h>
#include <sys/resource.h>
#include <sys/time.h>
#include <sys/socket.h>
#include <limits.h>

#include "cspace.h"
#include "csthreads.h"
#include "log.h"
#include "capp.h"
#include "utils.h"
#include "ipnet.h"

static CSThreadsList *threadsList=NULL;
static CSThreadsList *cPtr=NULL;
static void *shMem_global=NULL;
static int seqno=1;

static void process_parent_request(CSpaceThread *thread, PacketList *txList);
static int log_request(CSpaceThread *thread, int seqno, char *cmd,
		       char *params);
static void log_response(CSpaceThread *thread, int seq, char *resp);

/* From shmem.c */
extern int create_shared_memory(int *shmId, char *memFile, size_t size,
				ThreadShMem **shmem);

/*
 * init_cspace_thread --
 *
 */
CSpaceThread *init_cspace_thread(char *name, CSpace *space) {

  CSpaceThread *thread;

  thread = (CSpaceThread *)calloc(1, sizeof(CSpaceThread));
  if (thread == NULL) {
    return NULL;
  }

  uuid_generate(thread->uuid);

  thread->name  = strdup(name);
  thread->state = CSPACE_THREAD_STATE_CREATE;
  thread->space = space;

  thread->actionList = NULL;
  thread->shMem      = NULL;

  return thread;
}

/*
 * init_cspace_thread_list --
 */
int init_cspace_thread_list(void) {

  /* sanity check */
  if (threadsList) {
    log_error("Trying to re-initialize an existing thread list");
    return FALSE;
  }

  threadsList = (CSThreadsList *)calloc(1, sizeof(CSThreadsList));
  if (threadsList == NULL) {
    log_error("Error allocating memory of size: %d", sizeof(CSThreadsList));
    return FALSE;
  }

  cPtr = threadsList; /* set current pointer */

  return TRUE;
}

/*
 * free_cspace_thread --
 *
 */
void free_cspace_thread(CSpaceThread *thread) {

  free(thread->name);
  free(thread);
}

/*
 * free_cspace_thread_list --
 */
void free_cspace_thread_list(void) {

  CSThreadsList *ptr = threadsList;
  CSThreadsList *tmp;

  while (ptr) {

    free_cspace_thread(ptr->thread);
    tmp = ptr->next;
    free(ptr);

    ptr = tmp;
  }

  threadsList = NULL;
  cPtr        = NULL;
}

/*
 * add_to_cspace_thread_list --
 *
 */
int add_to_cspace_thread_list(CSpaceThread *thread) {

  FILE *fp=NULL;
  ShMemInfo *shmemInfo=NULL;

  if (thread == NULL) return FALSE;

  if (threadsList == NULL) {
    if (init_cspace_thread_list()!=TRUE) {
      log_error("Error adding to thread list: init failed");
      return FALSE;
    }
  } else {
    /* Add to the list */
    cPtr->next = (CSThreadsList *)calloc(1, sizeof(CSThreadsList));
    if (cPtr->next == NULL) {
      log_error("Error allocating memory of size: %d", sizeof(CSThreadsList));
      return FALSE;
    }

    cPtr->shmemInfo = (ShMemInfo *)malloc(sizeof(ShMemInfo));
    if (cPtr->shmemInfo==NULL) {
      log_error("Error allocating memory. size: %s", sizeof(ShMemInfo));
      return FALSE;
    }

    shmemInfo = cPtr->shmemInfo;

    /* create the memfile. */
    fp = fopen(CSPACE_MEMFILE, "w+");
    if (fp==NULL) {
      log_error("Error creating memory file: %s. Error: %s", CSPACE_MEMFILE,
		  strerror(errno));
      return FALSE;
    } else {
      fclose(fp);
    }

    if (!create_shared_memory(&(shmemInfo->shmId), CSPACE_MEMFILE,
			      sizeof(ThreadShMem), &(cPtr->shMem))) {
      log_error("Error creating shared memory for thread. Name: %s",
		thread->space->name);
      goto failure;
    }

    pthread_mutex_init(&(cPtr->shMem->txMutex), NULL);
    pthread_cond_init(&(cPtr->shMem->hasTX), NULL);
    pthread_mutex_lock(&(cPtr->shMem->txMutex));

    thread->shMem = cPtr->shMem;
  }

  cPtr->thread = thread;
  cPtr->next   = NULL;

  return TRUE;

 failure:
  free(cPtr->shmemInfo->memFile);
  free(cPtr->shmemInfo);

  return FALSE;
}

/*
 * cspace_exit_check --
 *
 */
static void cspace_exit_check(CSpaceThread *thread) {

  int status;
  pid_t w;

  /* check if the cspace exited */
  w = waitpid(thread->pid, &status, WNOHANG | WUNTRACED | WCONTINUED);

  if (w == 0) /* No status change. */
    return;

  if (w == -1) {
    log_error("waitpid failed for space: %s", thread->space->name);
    exit(EXIT_FAILURE);
  }

  if (WIFEXITED(status)) {
    log_debug("Space exited. Space name: %s status: %d\n",
	      thread->space->name, WEXITSTATUS(status));
    process_cspace_thread_exit(thread, CSPACE_THREAD_EXIT_NORMAL);
  } else if (WIFSIGNALED(status)) {
    printf("Space killed. Space name: %s signal: %d\n",
	   thread->space->name, WTERMSIG(status));
    process_cspace_thread_exit(thread, CSPACE_THREAD_EXIT_TERM);
  } else if (WIFSTOPPED(status)) {
    printf("space stopped. Space name: %s signal %d\n",
	   thread->space->name, WSTOPSIG(status));
    process_cspace_thread_exit(thread, CSPACE_THREAD_EXIT_STOP);
  }
}

/*
 * cspace_thread_start -- Thread routine to create cspaces
 */
void* cspace_thread_start(void *args) {

  CSpaceThread *thread  = (CSpaceThread *)args;
  ThreadShMem *shMem = NULL;
  struct timespec ts;
  struct timeval  tv;
  int status, ret;
  pid_t w;
  char idStr[36+1];

  if (!create_cspace(thread->space, &thread->pid)) {
    log_error("Error creating cspace: %s using config file: %s. Exiting",
	      thread->space->name, thread->space->configFile);
    exit(1);
  }

  /* setup networking using veth and bridge */
  if (ipnet_setup(IPNET_DEV_TYPE_CSPACE, DEF_BRIDGE, DEF_IFACE,
		  thread->space->name, thread->pid) != TRUE) {
    log_error("Error setting up networking for cspace. %s %s %ld",
	      DEF_BRIDGE, thread->space->name, (long)thread->pid);
    exit(1);
  }

  /* Test cspace network setup */
  if (!ipnet_test(thread->space->name)) {
    log_error("Failed test for cspace network setup: %s", thread->space->name);
  } else {
    log_debug("Passed test for cspace network setup: %s", thread->space->name);
  }

  uuid_unparse(thread->uuid, &idStr[0]);
  log_debug("Successfully created cspace. Name: %s UUID: %s PID: %d",
	    thread->space->name, idStr, thread->pid);

  /* set proper state for the cspace thread*/
  thread->state = CSPACE_THREAD_STATE_ACTIVE;

  shMem = thread->shMem;

  /* Lock on the RX mutex. */
  pthread_mutex_init(&(thread->shMem->rxMutex), NULL);
  pthread_cond_init(&(thread->shMem->hasRX), NULL);

  /* thread main loop */
  while(TRUE) {

    /* Check cspace exit status, if any. */
    cspace_exit_check(thread);

    gettimeofday(&tv, NULL);
    ts.tv_sec   = time(NULL) + CSPACE_COND_WAIT / 1000;
    ts.tv_nsec  = tv.tv_usec * 1000 + 1000 * 1000 * (CSPACE_COND_WAIT % 1000);
    ts.tv_sec  += ts.tv_nsec / (1000 * 1000 * 1000);
    ts.tv_nsec %= (1000 * 1000 * 1000);

    /* Timed wait on the capp packet from parent process. */
    //    ret = pthread_cond_timedwait(&(shMem->hasTX), &(shMem->txMutex), &ts);
    ret = pthread_cond_wait(&(shMem->hasTX), &(shMem->txMutex));
    if (ret == ETIMEDOUT) {
      continue;
    }

    if (ret == -1) {
      log_error("thread conditional/mutex wait error: %s. Exiting.",
		thread->space->name);
      exit(1);
    }

    /* Valid capp packet is waiting for us. Process it */
    if (ret == 0) {
      process_parent_request(thread, shMem->txList);
      pthread_mutex_unlock(&(shMem->rxMutex));
    }

    /* Check if there is anything on the socket to process. */
    process_response_packet(thread);
  }

  return (void *)0;
}

/*
 * process_cspace_thread_exit --
 *
 */
void process_cspace_thread_exit(CSpaceThread *thread, int status) {

  thread->exit_status = status;
  thread->state       = CSPACE_THREAD_STATE_ABORT;

}

/*
 * find_matching_thread_shmem -- for a given cspace name find the matching
 *                               shared memory.
 *
 */
ThreadShMem *find_matching_thread_shmem(char *name) {

  CSThreadsList *ptr=NULL;

  if (name == NULL) return NULL;

  if (threadsList == NULL) return NULL;

  for (ptr=threadsList; ptr; ptr=ptr->next) {
    if (strcmp(ptr->thread->name, name)==0) {
      return ptr->shMem;
    }
  }

  return NULL;
}

/*
 * process_parent_request --
 *
 */
static void process_parent_request(CSpaceThread *thread, PacketList *txList) {

  PacketList *list=NULL;
  CAppPacket *packet=NULL;
  char *cmd=NULL, *params=NULL;
  int size;

  if (txList == NULL) return;

  /* Steps are:
   * For each capp request:
   *   check the request type.
   *   convert the request into proper command
   *   send it over socket
   *   wait for response.
   *   process response.
   */

  for (list=txList; list; list=list->next) {

    packet = list->packet;
    if (!packet) continue;

    switch(packet->reqType) {
    case CAPP_TYPE_REQ_CREATE:

      if (!packet->name || !packet->tag || !packet->path) continue;

      cmd    = CAPP_CMD_CREATE;
      size   = 3 + strlen(packet->name) + strlen(packet->tag) +
	strlen(packet->path);
      params = (char *)malloc(size);
      sprintf(params, "%s:%s:%s", packet->name, packet->tag, packet->path);
      break;

    case CAPP_TYPE_REQ_RUN:
      cmd    = CAPP_CMD_RUN;
      params = (char *)malloc(36+1);
      uuid_unparse(packet->uuid, params);
      break;

    case CAPP_TYPE_REQ_STATUS:
      cmd    = CAPP_CMD_STATUS;
      params = (char *)malloc(36+1);
      uuid_unparse(packet->uuid, params);
      break;

    default:
      cmd    = NULL;
      params = NULL;
      log_error("Invalid request type recevied: %d", packet->reqType);
      break;
    }

    log_debug("Request recevied from parent process. Cmd: %s Params: %s",
	      cmd, params);

    if (cmd && params) {
      send_request_packet(thread, cmd, params);
      free(params);
    }
  }
}

/*
 * send_request_packet --
 *
 */
int send_request_packet(CSpaceThread *thread, char *cmd, char *params) {

  char *data;
  int size;

  if (thread == NULL || cmd == NULL || params == NULL) return FALSE;

  if (thread->space == NULL) return FALSE;

  if (thread->space->sockets[PARENT_SOCKET] <= 0) {
    log_error("Socket pair is closed between thread and cspace. Name: %s",
	      thread->name);
    return FALSE;
  }

  size = strlen(cmd) + strlen(params) +
    (3*sizeof(int)+2); /* for integer */

  data = (char *)malloc(size + 3); /* 2 space + 1 null */
  if (data == NULL) {
    log_error("Memory allocation error. Size: %d", size);
    return FALSE;
  }

  sprintf(data, "%s %d %s", cmd, seqno, params);

  log_debug("Sending data to thread: %s", data);

  if (write(thread->space->sockets[PARENT_SOCKET], data, strlen(data)) <0) {
    log_error("Error sending request packet to cspace via thread. Error: %s",
	      strerror(errno));
    free(data);
    return FALSE;
  }

  /* log request for the thread. This will help with matching the resp */
  log_request(thread, seqno, cmd, params);

  if (seqno == INT_MAX-1) {
    seqno = 0; /* reset */
  }
  seqno++;

  free(data);
  return TRUE;
}

/*
 * process_response_packet --
 *
 */
int process_response_packet(CSpaceThread *thread) {

  struct timeval tv;
  int count;
  char buffer[CSPACE_MAX_BUFFER] = {0};
  char seq[CSPACE_MAX_BUFFER]    = {0};
  char resp[CSPACE_MAX_BUFFER]   = {0};

  /* time-out socket */
  tv.tv_sec  = 5; /* XXX - check on this. */
  tv.tv_usec = 0;
  setsockopt(thread->space->sockets[PARENT_SOCKET], SOL_SOCKET, SO_RCVTIMEO,
	     (const char*)&tv, sizeof tv);

  count = recv(thread->space->sockets[PARENT_SOCKET], buffer,
	       CSPACE_MAX_BUFFER, 0);

  if (count <=0  && errno != EAGAIN) {
    log_error("Error reading packet from cspace socket. Name: %s",
	      thread->name);
    return CSPACE_READ_ERROR;
  }

  if (count == 0 && errno == EAGAIN) {
    return CSPACE_READ_TIMEOUT;
  }

  /* we have some packet. Let's see what we got
   * packet format is [seq_id some_text]
   */
  sscanf(buffer, "%s %s", seq, resp);

  log_response(thread, atoi(seq), resp);

  return TRUE;
}

/*
 * init_action --
 *
 */
static void init_action(ActionList *ptr, int seqno, char *cmd, char *params) {

  ptr->seqno  = seqno;
  ptr->state  = CAPP_CMD_STATE_WAIT;
  ptr->cmd    = strdup(cmd);
  ptr->params = strdup(params);
  ptr->resp   = NULL;
  ptr->next   = NULL;
}

/*
 * log_request --
 *
 */
static int log_request(CSpaceThread *thread, int seqno, char *cmd,
			char *params) {

  ActionList *ptr;

  if (!thread || !cmd || !params) return;

  /* base case */
  if (thread->actionList==NULL) {
    thread->actionList = (ActionList *)calloc(1, sizeof(ActionList));
    if (!thread->actionList) {
      log_error("Memory allocation error. Size: %d", sizeof(ActionList));
      return FALSE;
    }

    init_action(thread->actionList, seqno, cmd, params);
    return TRUE;
  }

  for (ptr = thread->actionList; ptr; ptr=ptr->next) {
    if (ptr->seqno == seqno) { /* already in the queue. */
      return FALSE;
    }
  }

  ptr = (ActionList *)malloc(sizeof(ActionList));
  if (!ptr) {
    log_error("Memory allocation error. Size: %d", sizeof(ActionList));
    return FALSE;
  }

  init_action(ptr, seqno, cmd, params);

  return TRUE;
}

/*
 * log_response --
 *
 */
static void log_response(CSpaceThread *thread, int seq, char *resp) {

  ActionList *ptr;

  /* Find the matching action from the list, update its state
   * and update the response
   */

  if (thread == NULL || resp == NULL) return;

  for (ptr = thread->actionList; ptr; ptr=ptr->next) {
    if (ptr->seqno == seq) {
      ptr->state = CAPP_CMD_STATE_DONE;
      ptr->resp  = strdup(resp);
      return;
    }
  }

  log_error("No mathing seq %d found in the action list. Ignoring. %s",
	    seq, thread->name);
  return;
}
