/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <pthread.h>
#include <string.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include <sys/wait.h>

#include "agent.h"
#include "wimc.h"

#define AGENT_EXEC "/usr/bin/casync"
#define DEFAULT_PATH "/tmp/"

/* For shared memory object and map */
static char *memFile=NULL;
static void *shMem=NULL;

static int is_valid_folder(char *folder);
static void log_wait_status(int status);
static void configure_runtime_args(WFetch *fetch, char **arg);
static void *execute_agent(void *data);
static void copy_fetch_request(WFetch **dest, WFetch *src);
void request_handler(WFetch *fetch);

/* from shmem. */
extern void *create_shared_memory(char *memFile, size_t size);
extern void delete_shared_memory(char *memFile, void *shMem, size_t size);
extern void read_stats_and_update_wimc(TStats *stats, WFetch *fetch);

/*
 * is_valid_folder -- If 'folder' is valid directory
 *
 */
static int is_valid_folder(char *folder) {

  struct stat sb;
  
  if (stat(folder, &sb) == -1) {
    return FALSE;
  }
  
  /* Check to see if it was file. */
  if (S_ISDIR(sb.st_mode)) {
    return TRUE;
  } else {
    return FALSE;
  }
}

/*
 * reset_stat_counter --
 *
 */
static void reset_stat_counter(void *ptr) {

  TStats *stat = (TStats *)ptr;

  stat->start = 0;
  stat->stop = 0;
  stat->exitStatus =0;

  stat->n_bytes=UINT64_MAX;
  stat->n_requests=UINT64_MAX;
  stat->n_local_requests=UINT64_MAX;
  stat->n_seed_requests=UINT64_MAX;
  stat->n_remote_requests=UINT64_MAX;
  stat->n_local_bytes=UINT64_MAX;
  stat->n_seed_bytes=UINT64_MAX;
  stat->n_remote_bytes=UINT64_MAX;
  stat->total_requests=UINT64_MAX;
  stat->total_bytes=UINT64_MAX;
  stat->nsec=UINT64_MAX;
  stat->runtime_nsec=UINT64_MAX;

  stat->status=WSTATUS_PEND;
  memset(stat->statusStr, 0, WIMC_MAX_ERR_STR);
}

/*
 * copy_fetch_request --
 *
 */
static void copy_fetch_request(WFetch **dest, WFetch *src) {

  WContent *content;
  WFetch *df;

  if (src == NULL)
    return;

  *dest = (WFetch *)calloc(1, sizeof(WFetch));
  if (*dest == NULL) return;

  df = *dest;

  uuid_copy(df->uuid, src->uuid);
  df->cbURL = strdup(src->cbURL);
  df->interval = src->interval;

  df->content = (WContent *)calloc(1, sizeof(WContent));
  if (df->content == NULL) {
    goto fail;
  }

  content = df->content;
  content->name = strdup(src->content->name);
  content->tag  = strdup(src->content->tag);
  content->method = strdup(src->content->method);
  content->providerURL = strdup(src->content->providerURL);
  content->indexURL = strdup(src->content->indexURL);
  content->storeURL = strdup(src->content->storeURL);

  return;

 fail:
  free(df->cbURL);
  free(df);
}

/*
 * free_fetch_request --
 *
 */
void free_fetch_request(WFetch *ptr) {

  WContent *cPtr;

  if (!ptr) return;

  free(ptr->cbURL);
  cPtr = ptr->content;

  if (!cPtr) return;

  free(cPtr->name);
  free(cPtr->tag);
  free(cPtr->method);
  free(cPtr->providerURL);
  free(cPtr->indexURL);
  free(cPtr->storeURL);
  free(cPtr);
}

/*
 * log_wait_status --
 *
 */
static void log_wait_status(int status) {

  if (WIFEXITED(status)) {
    log_debug("agent exited, status=%d\n", WEXITSTATUS(status));
  } else if (WIFSIGNALED(status)) {
    log_debug("agent killed by signal %d (%s)",
	      WTERMSIG(status), strsignal(WTERMSIG(status)));
    if (WCOREDUMP(status)) {
      printf("Reason: core dump");
    }
  } else if (WIFSTOPPED(status)) {
    log_debug("agent stopped by signal %d (%s)\n",
                WSTOPSIG(status), strsignal(WSTOPSIG(status)));
  } else {
    log_debug("agent stopped for unknown reasons.");
  }
}


/*
 * exit_thread -- exit thread, inform whomever, etc. 
 *
 */
static void exit_thread(void *retVal) {

  pthread_exit(retVal);
}

/*
 * configure_runtime_args -- setup runtime argument for the agent. 
 *
 */
static void configure_runtime_args(WFetch *fetch, char **arg) {

  char *argv=NULL;
  WContent *content=NULL;
  char folder[WIMC_MAX_PATH_LEN];
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  
  /* sanity check. */
  if (fetch==NULL) return;

  content = fetch->content;
  
  if (content->indexURL==NULL && content->storeURL==NULL) return;

  *arg = (char *)calloc(1, WIMC_MAX_ARGS_LEN);
  if (*arg==NULL) 
    return;

  argv = *arg;
  memset(&folder[0], 0, WIMC_MAX_PATH_LEN);

  /* Setup download path */
  if (!uuid_is_null(fetch->uuid)) {
    uuid_unparse(fetch->uuid, idStr);
    sprintf(folder, "%s/%s/%s_%s", DEFAULT_PATH, idStr, content->name,
	    content->tag);
  } else {
    goto done;
  }
  
  /* check if directory exists */
  if (!is_valid_folder(&folder[0])) {
    log_debug("Creating default folder for Agent at: %s", folder);
    mkdir(folder, 0700);
  }

  /* Put everything together */
  sprintf(argv, "%s extract %s --store %s --shmem %s %s", AGENT_EXEC,
	  content->indexURL, content->storeURL, memFile, folder);
  
  return;
  
 done:
  free(*arg);
  *arg=NULL;
  return;
}

/*
 * execute_agent --
 *
 */
static void *execute_agent(void *data) {

  WFetch *fetch;
  char *args=NULL;
  TStats *stats=NULL;
  TParams *params;
  pid_t pid=0;
  int ret;
  pthread_t tid;
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  
  fetch = (WFetch *)data;

  params = (TParams *)malloc(sizeof(TParams));

  if (params==NULL) {
    log_error("Memory allocation error. size: %s", sizeof(TParams));
    return (void *)0;
  }
      
  /* de-attach itself from the parent. */
  // pthread_detach(pthread_self()); XXX remove me.

  /* Step 1. configure runtime argument for agent. */
  if (!uuid_is_null(fetch->uuid)) {
    memFile = (char *)calloc(1, WIMC_MAX_PATH_LEN);
    if (memFile==NULL) {
      log_error("Memory allocation error. size: %s", WIMC_MAX_PATH_LEN);
      return (void *)0;
    }
    uuid_unparse(fetch->uuid, idStr);
    sprintf(memFile, "%s.shmem", idStr);
  } else {
    memFile = strdup(DEFAULT_SHMEM);
  }

  configure_runtime_args(fetch, &args);
  if (args == NULL) {
    log_error("Can not setup runtime argument for the Agent");
    exit_thread(NULL);
  } else {
    log_debug("Agent runtime arguments: %s", args);
  }
  
  /* Step 2. configure shared memory */
  shMem = create_shared_memory(memFile, sizeof(TStats)); /* use default file */
  if (shMem == MAP_FAILED || shMem == NULL) {
    log_error("Error creating shared memory of size: %d. Error: %s",
	      sizeof(TStats), strerror(errno));
    return (void *)0;
    exit_thread(NULL);
  }
  reset_stat_counter(shMem);
  
  /* Step 3. Fork and exec */
  pid = fork();
  if (pid < 0) {
    log_error("Failed to fork for agent");
    return (void *) 0;
  }
  
  if (pid==0) { /* Child process. */
    execv(AGENT_EXEC, args);
    _exit(127);
  } else {

    params->stats = shMem;
    copy_fetch_request((WFetch **)&params->fetch, fetch);
    
    /* Step 4. process status chanage and update wimc.d */
    /* Thread to read the update status from agent and send to WIMC */
    ret = pthread_create(&tid, NULL, read_stats_and_update_wimc, params);
    if (ret) { /* Some error. */
      log_error("Error creating agent thread. Return code: %s", ret);
      goto failure;
    }
    free(args);
    free(params);
    return (void *) pid;
  }

 failure:
  delete_shared_memory(memFile, shMem, sizeof(TStats));
  free(params);
  free(memFile);
  free(args);

  return (void *)0;
}

/*
 * request_handler --
 *
 */

void request_handler(WFetch *fetch) {

  /* Flow is as follows:
   * 1. create thread.
   * 2. Setup runtime argument for CA-Sync.
   * 3. setup shared memory space for status update etc.
   * 3. Fork and run ca-sysnc.
   * 4. update the wimc.d, on the callback URL, transfer status after 'interval'
   * 5. monitor child process exit and its status.
   *
   */
  int ret, wstatus;
  pthread_t tid;
  void *status;
  pid_t pid, w;
  
#if 0
  /* Step1-4: Thread which will fork/exec the agent and send status via CB */
  ret = pthread_create(&tid, NULL, execute_agent, (void *)fetch);
  if (ret) { /* Some error. */
    log_error("Error creating agent thread. Return code: %s", ret);
    return;
  }
#endif

  execute_agent((void *)fetch);

#if 0
  pthread_join(tid, &status);

  pid = (pid_t)status;

  if (!pid) { /* XXX */
    return;
  }

  /* Step 5. wait for the agent to exit. */
  do {
    w = waitpid(pid, &wstatus, WUNTRACED | WCONTINUED);
    
    /* Once child exit, kill the stats updating thread. 
       XXXX
      */
    if (w == -1) {
      perror("waitpid");
      exit(EXIT_FAILURE);
    }
    
    if (WIFEXITED(wstatus)) {
      printf("exited, status=%d\n", WEXITSTATUS(wstatus));
    } else if (WIFSIGNALED(wstatus)) {
      printf("killed by signal %d\n", WTERMSIG(wstatus));
    } else if (WIFSTOPPED(wstatus)) {
      printf("stopped by signal %d\n", WSTOPSIG(wstatus));
    } else if (WIFCONTINUED(wstatus)) {
      printf("continued\n");
    }
  } while (!WIFEXITED(wstatus) && !WIFSIGNALED(wstatus));

  /* Agent is done. clearup shared the shared memory object and mapping */
  delete_shared_memory(memFile, shMem, sizeof(TStats));
#endif
  
  free(memFile);
  memFile = NULL;
  shMem = NULL;
}
