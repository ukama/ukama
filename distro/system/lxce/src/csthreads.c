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

#include "cspace.h"
#include "csthreads.h"
#include "log.h"

static CSThreadsList *threadsList=NULL;
static CSThreadsList *cPtr=NULL;

/*
 * init_cspace_thread --
 *
 */
CSpaceThread *init_cspace_thread(char *name, CSpace *space) {

  CSpaceThread *thread;

  thread = (CSpaceThread *) malloc(sizeof(CSpaceThread));
  if (thread == NULL) {
    return NULL;
  }

  uuid_generate(thread->uuid);

  thread->name  = strdup(name);
  thread->state = CSPACE_THREAD_STATE_CREATE;
  thread->space = space;

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
    cPtr = cPtr->next;
  }

  cPtr->thread = thread;
  cPtr->next   = NULL;

  return TRUE;
}

/*
 * cspace_thread_start -- Thread routine to create cspaces
 */
void* cspace_thread_start(void *args) {

  CSpaceThread *thread  = (CSpaceThread *)args;
  int status;
  pid_t w;

  if (!create_cspace(thread->space, &thread->pid)) {
    log_error("Error creating cspace: %s using config file: %s. Exiting",
	      thread->space->name, thread->space->configFile);
    exit(1);
  } else {
    log_debug("Successfully created cspace: %s", thread->space->name);
  }

  /* set proper state for the cspace thread*/
  thread->state = CSPACE_THREAD_STATE_ACTIVE;

  /* Wait for the child to exit, aka space abort. */
  do {
    w = waitpid(thread->pid, &status, WUNTRACED | WCONTINUED);

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
  } while (!WIFEXITED(status) && !WIFSIGNALED(status));

  return;
}

/*
 * process_cspace_thread_exit --
 *
 */
void process_cspace_thread_exit(CSpaceThread *thread, int status) {

  thread->exit_status = status;
  thread->state       = CSPACE_THREAD_STATE_ABORT;

  return;
}
