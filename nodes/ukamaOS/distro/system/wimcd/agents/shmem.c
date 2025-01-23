/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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

#include "usys_log.h"
#include "usys_mem.h"

/* wimc.c */
extern long communicate_with_wimc(int reqType,
                                  char *wimcURL,
                                  char *cbURL,
                                  void *data);
/* thread.c */
void free_fetch_request(WFetch *ptr);
 
#define SECRET_ID 46504650 /* Id for ftok() */

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
    WContent *content;
    int interval;
    long code;
    char idStr[36+1] = {0};
    char folder[WIMC_MAX_PATH_LEN] = {0};

    if (args == NULL) return;

    params = (TParams *)args;
    stats  = (TStats *)params->stats;
    fetch  = (WFetch *)params->fetch;

    if (stats == NULL) {
        /* This can happen when 1. not yet created or 2. is freed */
        usys_log_error("Trying to access invalid shared memory object");
        return;
    }

    if (fetch->interval==0) {
        interval = DEFAULT_INTERVAL;
    } else {
        interval = fetch->interval;
    }

  do {
      
      usys_log_debug("Sending update to wimc.d ...");
      code = communicate_with_wimc(WREQ_UPDATE,
                                   fetch->cbURL,
                                   NULL,
                                   (void *)stats);
      /* code 404 means ID is not yet register. We try again. */
      if (!code || code == 400 ) {
          usys_log_error("Failed to send update to the wimc.d. Thread Exit");
          goto cleanup;
      }

      /* We exit if agent is done. */
      if (stats->stop == TRUE) {
          goto last;
      }

      /* Also exit if task is done or if there is an error. */
      if (stats->status == (TaskStatus)WSTATUS_DONE ||
          stats->status == (TaskStatus)WSTATUS_ERROR) {
          goto last;
      }

      /* Otherwise sleep for 'interval' and repeat again. */
      sleep(interval);
  } while(TRUE);

last:
  if (stats->stop == TRUE) {

      content = fetch->content;
      uuid_unparse(fetch->uuid, idStr);
      sprintf(folder, "%s/%s/%s_%s", DEFAULT_PATH, idStr, content->name,
              content->tag);

      /* Update the path location and notify WIMC of it */
      stats->status = WSTATUS_DONE;
      strcpy(stats->statusStr, folder);
      communicate_with_wimc(WREQ_UPDATE,
                            fetch->cbURL,
                            NULL,
                            (void *)stats);
  }

cleanup:
  free_fetch_request(fetch);
  free(fetch);
}
