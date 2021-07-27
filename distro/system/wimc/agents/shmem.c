/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/ipc.h>
#include <sys/types.h>
#include <sys/shm.h>

#include "wimc.h"

#define SECRET_ID 46504650 /* Id for ftok() */

/*
 * create_shared_memory --
 *
 */
void *create_shared_memory(int *shmId, char *memFile, size_t size) {

  key_t key;

  /* Sanity check. */
  if (memFile == NULL) return NULL;

  key = ftok(memFile, SECRET_ID);
  if (key == -1) {
    log_error("Error generating key token for shared memory. Error: %s",
	      strerror(errno));
    return NULL;
  }

  *shmId = shmget(key, size, 0644|IPC_CREAT);
  if (*shmId == -1) {
    log_error("Error creating shared memory of size %d. Error: %s",
	      (int)size, strerror(errno));
    return NULL;
  }

  return shmat(*shmId, NULL, 0);
}

/*
 * delete_shared_memory --
 *
 */
void delete_shared_memory(int shmid, void *shMem) {

  if (shmdt(shMem) == -1) {
    log_error("Error deattaching. Error: %s \n", strerror(errno));
    return;
  }

  /*Mark it to be destroyed */
  if (shmctl(shmid, IPC_RMID, 0) == -1) {
    log_error("Error destroying shared memory. Error: %s", strerror(errno));
    return;
  }
}

/*
 * read_data_and_update_wimc -- Read data available at the shmem after 
 *                              certain interval and send it back to wimc 
 *                              callback URL.
 */
void read_stats_and_update_wimc(void *args) {

  TParams *params;
  TStats *stats;
  WFetch *fetch;
  int interval;
  long code;

  params = (TParams *)args;
  
  /* sanity check. */
  if (params == NULL)
    return;

  stats = (TStats *)params->stats;
  fetch = (WFetch *)params->fetch;

  if (stats == NULL) {
    /* This can happen when 1. not yet created or 2. is freed */
    log_error("Trying to access invalid shared memory object");
    return;
  }
  
  if (fetch->interval==0) {
    interval=DEFAULT_INTERVAL;
  } else {
    interval=fetch->interval;
  }

  do {
    log_debug("Sending update to wimc.d ...");
    code = communicate_with_wimc(REQ_UPDATE, fetch->cbURL, NULL, NULL,
				 fetch->uuid, (void *)stats);
    if (!code || code == 400 || code == 404) {
      log_error("Failed to send update to the wimc.d. Thread Exit");
      goto cleanup;
    }

    /* We exit if agent is done. */
    if (stats->stop == TRUE) {
      goto cleanup;
    }

    /* Also exit if task is done or if there is an error. */
    if (stats->status == (TaskStatus)WSTATUS_DONE ||
	stats->status == (TaskStatus)WSTATUS_ERROR) {
      goto cleanup;
    }

    /* Otherwise sleep for 'interval' and repeat again. */
    sleep(interval);
  } while(TRUE);

 cleanup:
  free_fetch_request(fetch);
  free(fetch);
  return;
}
