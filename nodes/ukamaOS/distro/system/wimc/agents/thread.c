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
#define DEFAULT_PATH "/tmp"
#define MAX_ARGS 10

/* For shared memory object and map */
static char *memFile=NULL;
static void *shMem=NULL;
static int shmId=0;

static int is_valid_folder(char *folder);
static void log_wait_status(int status);
static void configure_runtime_args(WFetch *fetch, char **arg);
static void *execute_agent(void *data);
static void copy_fetch_request(WFetch **dest, WFetch *src);
void request_handler(WFetch *fetch);

/* from shmem. */
extern void *create_shared_memory(int *shmId, char *memFile, size_t size);
extern void delete_shared_memory(int shmId, void *shMem);
extern void read_stats_and_update_wimc(void *args);

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
 * create_working_dir_file --
 *
 */
static void create_working_dir_file(WFetch *fetch) {

  FILE *fp;
  char folder[WIMC_MAX_PATH_LEN]={0};
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  WContent *content;

  if (fetch==NULL || uuid_is_null(fetch->uuid))
    return;

  content = fetch->content;

  uuid_unparse(fetch->uuid, idStr);
  sprintf(folder, "%s/%s", DEFAULT_PATH, idStr);

  /* check if directory exists */
  if (!is_valid_folder(&folder[0])) {
    log_debug("Creating default folder for Agent at: %s", folder);
    if (mkdir(folder, 0700) == -1) {
      log_error("Error creating dir: %s. Error: %s", folder, strerror(errno));
    }
  }

  /* Create shared memory file. */
  sprintf(memFile, "%s/%s/shared.mem", DEFAULT_PATH, idStr);
  log_debug("Creating default shared memory file for Agent at: %s", memFile);
  fp = fopen(memFile, "a");
  if (fp) {
    fclose(fp);
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
 * configure_runtime_args -- setup runtime argument for the agent. 
 *
 */
static void configure_runtime_args(WFetch *fetch, char **args) {

  WContent *content=NULL;
  char folder[WIMC_MAX_PATH_LEN];
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  
  /* sanity check. */
  if (fetch==NULL) return;

  content = fetch->content;
  
  if (content->indexURL==NULL && content->storeURL==NULL) return;

  memset(&folder[0], 0, WIMC_MAX_PATH_LEN);

  /* Setup download path */
  if (!uuid_is_null(fetch->uuid)) {
    uuid_unparse(fetch->uuid, idStr);
    sprintf(folder, "%s/%s/%s_%s", DEFAULT_PATH, idStr, content->name,
	    content->tag);
  } else {
    goto done;
  }

  /* Put everything together */
  args[0] = strdup(AGENT_EXEC);
  args[1] = strdup("extract");
  args[2] = strdup(content->indexURL);
  args[3] = strdup("--store");
  args[4] = strdup(content->storeURL);
  args[5] = strdup(folder);
  args[6] = NULL; /* Null terminate for execv */

  return;
 done:
  args[0] = NULL;
  return;
}

/*
 * execute_agent --
 *
 */
static void *execute_agent(void *data) {

  WFetch *fetch;
  char *args[MAX_ARGS];
  TStats *stats=NULL;
  TParams *params;
  pid_t pid=0;
  int ret=0, i=0;
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  char buffer[1024] = {0};
  FILE *fp;

  fetch = (WFetch *)data;

  params = (TParams *)malloc(sizeof(TParams));
  memFile = (char *)calloc(1, WIMC_MAX_PATH_LEN);
  if (params==NULL || memFile==NULL) {
    log_error("Memory allocation error. size: %s", sizeof(TParams));
    pthread_exit(&ret);
  }

  /* create working directory anf file (for shared memory) */
  create_working_dir_file(fetch);

  /* Step 1. configure runtime argument for agent. */
  configure_runtime_args(fetch, args);
  if (args[0] == NULL) {
    log_error("Can not setup runtime argument for the Agent");
    pthread_exit(&ret);
  } else {
    for (i=0; args[i] != NULL && i < MAX_ARGS; i++) {
      sprintf(buffer, "%s %s", buffer, args[i]);
    }
    log_debug("Agent runtime arguments: %s", buffer);
  }

  /* Step 2. configure shared memory */
  shMem = create_shared_memory(&shmId, memFile, sizeof(TStats));
  if (shMem == MAP_FAILED || shMem == NULL) {
    log_error("Error creating shared memory of size: %d. Error: %s",
	      sizeof(TStats), strerror(errno));
    goto failure;
  }
  reset_stat_counter(shMem);
  
  /* Step 3. Fork and exec */
  pid = fork();
  if (pid < 0) {
    log_error("Failed to fork for agent");
    goto failure;
  }
  
  if (pid==0) { /* Child process. */
    execv(AGENT_EXEC, args);
    _exit(127);
  } else {

    params->stats = shMem;
    copy_fetch_request((WFetch **)&params->fetch, fetch);
    
    /* Step 4. process status chanage and update wimc.d */
    /* Thread to read the update status from agent and send to WIMC */
    read_stats_and_update_wimc(params);
    ret = 1; /* We are good. */
  }

failure:
  delete_shared_memory(shmId, shMem);
  shMem = NULL;
  free(params);
  free(memFile);
  for (i=0; args[i] != NULL && i < MAX_ARGS; i++) {
    if (args[i]) free(args[i]);
  }

  pthread_exit(&ret);
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
  int status;
  pid_t pid, w;

  /* Step1-4: Thread which will fork/exec the agent and send status via CB */
  ret = pthread_create(&tid, NULL, execute_agent, (void *)fetch);
  if (ret) { /* Some error. */
    log_error("Error creating agent thread. Return code: %s", ret);
    return;
  }

  log_debug("Waiting for agent thread to finish it work ...");

  pthread_join(tid, &status);

  if (status == 0) {
    log_error("Error executing agent for request handler.");
  } else if (status == 1) {
    log_debug("Successfully executed request handler. Forward and upward");
  }
}
