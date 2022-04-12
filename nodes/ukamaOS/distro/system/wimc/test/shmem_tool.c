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
#include <stdint.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>

#define SECRET_ID 46504650
#define USAGE "%s: memFile\n"
#define MAX_ERR_STR 1024

typedef enum {
  WSTATUS_PEND=1,
  WSTATUS_START,
  WSTATUS_RUNNING,
  WSTATUS_DONE,
  WSTATUS_ERROR
} TaskStatus;

typedef struct tStats {

  int start;
  int stop;
  int exitStatus;  /* as returned by waitpid */

  /* various stats */
  uint64_t n_bytes;
  uint64_t n_requests;
  uint64_t n_local_requests;
  uint64_t n_seed_requests;
  uint64_t n_remote_requests;
  uint64_t n_local_bytes;
  uint64_t n_seed_bytes;
  uint64_t n_remote_bytes;
  uint64_t total_requests;
  uint64_t total_bytes;
  uint64_t nsec;
  uint64_t runtime_nsec;

  TaskStatus status;
  char statusStr[MAX_ERR_STR];
} TStats;

static void *sharedMem=NULL;

int main(int argc, char **argv) {

  char *memFile=NULL;
  TStats *stats=NULL;
  key_t key;
  int shmid;

  if (argc!=2) {
    fprintf(stderr, USAGE, argv[0]);
    exit(1);
  }

  memFile = argv[1];

  key = ftok(memFile, SECRET_ID);
  if (key == -1) {
    fprintf(stderr, "Error generating key token for shared memory. Error: %s\n",
	    strerror(errno));
    exit(1);
  }

  shmid = shmget(key, sizeof(TStats), 0644|IPC_CREAT);
  if (shmid == -1) {
    fprintf(stderr, "Error creating shared memory of size %d. Error: %s\n",
	    (int)sizeof(TStats), strerror(errno));
    exit(1);
  }

  sharedMem = shmat(shmid, NULL, 0);

  stats = (TStats *)sharedMem;

  fprintf(stdout, "Press any char to start reading from shared memory...\n");
  getchar();

  fprintf(stdout, "Started reading from shared memory ... \n");

  do {

    fprintf(stdout, "\r ... Recevied: %ju (bytes)", stats->n_bytes);
  } while (stats->stop != 1);

  fprintf(stdout, "\n ... Done\n");

  /* cleanup */
  if (shmdt(sharedMem) == -1) {
    fprintf(stderr, "Error deattaching. Error: %s \n", strerror(errno));
    exit(1);
  }

  /* Remove this ID. ca-sync doesn't remove this. */
  if (shmctl(shmid, IPC_RMID, 0) == -1) {
    fprintf(stderr, "Error removing shared memory id. Error: %s\n", 
	    strerror(errno));
    exit(1);
  }

  return 0;
}
