/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WIMC_H
#define WIMC_H

#define WIMC_MAX_NAME_LEN   256
#define WIMC_MAX_PATH_LEN   256
#define WIMC_MAX_ERR_STR    1024

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
  char statusStr[WIMC_MAX_ERR_STR];
} TStats;

void reset_stat_counter(void *ptr);
void *create_shared_memory(char *memFile, size_t size);
void update_shmem_counters(CaSync *s, void *shMem);
void flag_end_shared_memory(void *ptr);
void flag_start_shared_memory(void *ptr);


#endif /* WIMC_H */
