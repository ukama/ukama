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

#include "parse-util.h"
#include "casync.h"
#include "wimc.h"

/*
 * reset_stat_counter --
 *
 */

void reset_stat_counter(void *ptr) {

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
 * create_shared_memory --
 *
 */

void *create_shared_memory(char *memFile, size_t size) {

  int memFd;
  int prot=0, flags=0, ret;

  if (memFile == NULL) {
    return NULL;
  }

  memFd = shm_open(memFile, O_RDWR, S_IRWXU);
  if (memFd == -1) {
    log_error("Error creating shared memory object. Error: %s",
              strerror(errno));
    goto fail;
  }

  /* Truncate to the right size. */
  ret = ftruncate(memFd, size);
  if (ret == -1) {
    log_error("Error truncating the shared memory file to size: %d. Error: %s",
              (int)size, strerror(errno));
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
 * update_shmem_counters --
 *
 */

void update_shmem_counters(CaSync *s, void *shMem) {
  
  char buffer[FORMAT_BYTES_MAX];
  TStats *stats=NULL;
  uint64_t n_bytes, n_requests;
  uint64_t n_local_requests = UINT64_MAX, n_seed_requests = UINT64_MAX;
  uint64_t n_remote_requests = UINT64_MAX, n_local_bytes = UINT64_MAX;
  uint64_t n_seed_bytes = UINT64_MAX, n_remote_bytes = UINT64_MAX;
  uint64_t total_requests = 0, total_bytes = 0;
  uint64_t nsec = 0, runtime_nsec = 0;
  int r;

  stats = (TStats *)shMem;

  ca_sync_get_punch_holes_bytes(s, &n_bytes);
  ca_sync_get_reflink_bytes(s, &n_bytes);
  ca_sync_get_hardlink_bytes(s, &n_bytes);

  r = ca_sync_get_local_requests(s, &n_requests);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_local_requests = n_requests;
    total_requests += n_requests;
  }

  r = ca_sync_get_local_request_bytes(s, &n_bytes);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_local_bytes = n_bytes;
    total_bytes += n_bytes;
  }

  r = ca_sync_get_seed_requests(s, &n_requests);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_seed_requests = n_requests;
    total_requests += n_requests;
  }

  r = ca_sync_get_seed_request_bytes(s, &n_bytes);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_seed_bytes = n_bytes;
    total_bytes += n_bytes;
  }

  r = ca_sync_get_remote_requests(s, &n_requests);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_remote_requests = n_requests;
    total_requests += n_requests;
  }

  r = ca_sync_get_remote_request_bytes(s, &n_bytes);
  if (!IN_SET(r, -ENODATA, -ENOTTY)) {
    if (r < 0) return;

    n_remote_bytes = n_bytes;
    total_bytes += n_bytes;
  }

#if 0  
  if (n_local_requests != UINT64_MAX)
    log_info("Chunk requests fulfilled from local store: %" PRIu64 " (%" PRIu64 "%%)",
	     n_local_requests,
	     total_requests > 0 ? n_local_requests * 100U / total_requests : 0);
  if (n_local_bytes != UINT64_MAX)
    log_info("Bytes used from local store: %s (%" PRIu64 "%%)",
	     format_bytes(buffer, sizeof(buffer), n_local_bytes),
	     total_bytes > 0 ? n_local_bytes * 100U / total_bytes : 0);
  if (n_seed_requests != UINT64_MAX)
    log_info("Chunk requests fulfilled from local seed: %" PRIu64 " (%" PRIu64 "%%)",
	     n_seed_requests,
	     total_requests > 0 ? n_seed_requests * 100U / total_requests : 0);
  if (n_seed_bytes != UINT64_MAX)
    log_info("Bytes used from local seed: %s (%" PRIu64 "%%)",
	     format_bytes(buffer, sizeof(buffer), n_seed_bytes),
	     total_bytes > 0 ? n_seed_bytes * 100U / total_bytes : 0);
  if (n_remote_requests != UINT64_MAX)
    log_info("Chunk requests fulfilled from remote store: %" PRIu64 " (%" PRIu64 "%%)",
	     n_remote_requests,
	     total_requests > 0 ? n_remote_requests * 100U / total_requests : 0);
  if (n_remote_bytes != UINT64_MAX)
    log_info("Bytes used from remote store: %s (%" PRIu64 "%%)",
	     format_bytes(buffer, sizeof(buffer), n_remote_bytes),
	     total_bytes > 0 ? n_remote_bytes * 100U / total_bytes : 0);
  
  r = ca_sync_get_runtime_nsec(s, &runtime_nsec);
  if (!IN_SET(r, -ENODATA)) {
    if (r < 0)
      return log_error_errno(r, "Failed to determine runtime: %m");
  }
#endif

  stats->n_bytes=n_bytes;
  stats->n_requests=n_requests;
  stats->n_local_requests=n_local_requests;
  stats->n_seed_requests=n_seed_requests;
  stats->n_remote_requests=n_remote_requests;
  stats->n_local_bytes=n_local_bytes;
  stats->n_seed_bytes=n_seed_bytes;
  stats->n_remote_bytes=n_remote_bytes;
  stats->total_requests=total_requests;
  stats->total_bytes=total_bytes;
  stats->nsec=nsec;
  stats->runtime_nsec=runtime_nsec;

  stats->status=WSTATUS_RUNNING;
  memset(stats->statusStr, 0, WIMC_MAX_ERR_STR);

}
