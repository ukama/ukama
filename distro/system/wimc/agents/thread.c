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

#define AGENT_EXEC "ca-sync"
#define DEFAULT_PATH "/tmp/"

/* For passing thread arguments */
typedef struct {

  TStats *stats;
  WFetch  *fetch;
} TParams;

static int is_valid_folder(char *folder);
static void log_wait_status(int status);
static void configure_runtime_args(WFetch *fetch, char *memFile, char **arg);
static void *execute_agent(void *data);
void request_handler(WFetch *fetch);

/* from shmem. */
extern void *create_shared_memory(char *memFile, size_t size);
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
static void configure_runtime_args(WFetch *fetch, char *memFile, char **arg) {

  char *argv=NULL;
  WContent *content=NULL;
  char folder[WIMC_MAX_PATH_LEN];
  char idStr[36+1]; /* 36-bytes for UUID + trailing `\0` */
  char *shMem=NULL;
  
  /* sanity check. */
  if (content->indexURL==NULL && content->storeURL==NULL) 
    return;

  content = fetch->content;
  
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

  if (memFile) {
    shMem = memFile;
  } else {
    shMem = DEFAULT_SHMEM;
  }
  
  /* Put everything together */
  sprintf(argv, "%s extract %s --store %s --shmem %s %s", AGENT_EXEC,
	  content->indexURL, content->storeURL, shMem, folder);
  log_debug("Agent runtime arguments: %s", argv);
  
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
  char *args=NULL, *memFile=NULL;
  TStats *stats=NULL;
  TParams *params;
  void *shMem=NULL;
  pid_t pid=0;
  int ret;
  pthread_t tid;
  
  fetch = (WFetch *)data;

  stats = (TStats *)calloc(1, sizeof(TStats));
  if (stats==NULL) {
    log_error("Memory allocation error. size: %s", sizeof(TStats));
    exit_thread(NULL);
  }
      
  /* de-attach itself from the parent. */
  pthread_detach(pthread_self());

  /* 1. configure runtime argument for agent. */
  configure_runtime_args(fetch, shMem, &args); /* XXX - free args. */
  if (args == NULL) {
    log_error("Can not setup runtime argument for the Agent");
    exit_thread(NULL);
  }
  
  /* 2. configure shared memory */
  shMem = create_shared_memory(memFile, sizeof(stats)); /* use default file */
  if (shMem == MAP_FAILED || shMem == NULL) {
    log_error("Error creating shared memory of size: %d. Error: %s",
	      sizeof(stats), strerror(errno));
    exit_thread(NULL);
  }

  stats = (TStats *)shMem;
  
  /* 3. Fork and exec */
  pid = fork();
  if (pid < 0) {
    log_error("Failed to fork for agent");
    return (void *) 0;
  }
  
  if (pid==0) { /* Child process. */
    execv(AGENT_EXEC, args);
    _exit(127);
  } else {

    params = (TParams *)malloc(sizeof(TParams));
    params->stats = stats;
    params->fetch = fetch;
    
    /* 4. process status chanage and update wimc.d */
    /* Thread to read the update status from agent and send to WIMC */
    ret = pthread_create(&tid, NULL, read_stats_and_update_wimc, params);
    if (ret) { /* Some error. */
      log_error("Error creating agent thread. Return code: %s", ret);
      return (void *) 0;
    }
    return (void *) pid;
  }
  
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
   * 4. monitor child process exit and its status.
   * 5. update the wimc.d, on the callback URL, transfer status after 'interval'
   */

  int ret, wstatus;
  pthread_t tid;
  void *status;
  pid_t pid, w;
  
  /* Thread which will fork/exec the agent. */
  ret = pthread_create(&tid, NULL, execute_agent, (void *)fetch);
  if (ret) { /* Some error. */
    log_error("Error creating agent thread. Return code: %s", ret);
    return;
  }

  pthread_join(tid, &status);

  pid = (pid_t)status;

  if (!pid) { /* XXX */
    return;
  }
  
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
  
}
