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

#include "log.h"
#include "csthreads.h"

#define SECRET_ID 46504650 /* Id for ftok() */

/*
 * create_shared_memory --
 *
 */
int create_shared_memory(int *shmId, char *memFile, size_t size,
			 ThreadShMem **shmem) {

  key_t key;

  /* Sanity check. */
  if (memFile == NULL) return NULL;

  key = ftok(memFile, SECRET_ID);
  if (key == -1) {
    log_error("Error generating key token for shared memory. Error: %s",
	      strerror(errno));
    return FALSE;
  }

  *shmId = shmget(key, size, 0644 | IPC_CREAT);
  if (shmId == -1) {
    log_error("Error creating shared memory of size %d. Error: %s",
	      (int)size, strerror(errno));
    return FALSE;
  }

  *shmem = shmat(*shmId, NULL, 0);
  if (*shmem == MAP_FAILED || *shmem == NULL) {
    log_error("Error creating shared memory of size: %d. Error: %s", size,
	      strerror(errno));
    return FALSE;
  }

  return TRUE;
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
