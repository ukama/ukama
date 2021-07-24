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
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>

#include "wimc.h"

/*
 * create_shared_memory --
 *
 */

void *create_shared_memory(char *memFile, size_t size) {

  int memFd;
  int prot=0, flags=0, ret;
  char *path=NULL;

  if (memFile == NULL) {
    path = DEFAULT_SHMEM;
  } else {
    path = memFile;
  }

  memFd = shm_open(path, O_CREAT | O_RDWR, S_IRWXU);
  if (memFd == -1) {
    log_error("Error creating shared memory object. Error: %s",
	      strerror(errno));
    goto fail;
  }

  /* Truncate to the right size. */
  ret = ftruncate(memFd, size);
  if (ret == -1) {
    log_error("Error truncating the shared memory file to size: %d. Error: %s",
	      size, strerror(errno));
    goto fail;
  }

  /* readable and writeable */
  prot = PROT_READ | PROT_WRITE;

  /* only visible for parent/child and none other. */
  flags = MAP_SHARED | MAP_ANONYMOUS;

  return mmap(NULL, size, prot, flags, memFd, 0);

 fail:
  return NULL;
}

/*
 * delete_shared_memory --
 *
 */
void delete_shared_memory(char *memFile, void *shMem, size_t size) {

  int ret;

  ret = munmap(shMem, size);
  if (ret==-1) {
    log_error("Failed to unmap shared memory of size: %d. Error: %s", size,
	      strerror(errno));
    return;
  }
  free(shMem);
  shMem = NULL;

  ret = shm_unlink(memFile);
  if (ret==-1) {
    log_error("Failed to close shared memory object. Error: %s",
	      strerror(errno));
    return;
  }
}

/*
 * read_data_and_update_wimc -- Read data available at the shmem after 
 *                              certain interval and send it back to wimc 
 *                              callback URL.
 */
void read_stats_and_update_wimc(TStats *stats, WFetch *fetch) {

  int interval;
  long code;
  
  /* sanity check. */
  if (fetch == NULL)
    return;

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
    code = communicate_with_wimc(REQ_UPDATE, fetch->cbURL, NULL, NULL,
				 fetch->uuid, (void *)stats);
    if (!code || code == 400) {
      log_error("Failed to send update to the wimc.d. Exit");
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
    wait(interval);
  } while(TRUE);

 cleanup:

  return;
}
