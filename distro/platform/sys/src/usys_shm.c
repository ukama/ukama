/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_shm.h"

#include "usys_error.h"
#include "usys_log.h"

void *usys_allocate_shared_mem(const char *name, uint32_t size) {
    int fd, ret;
    void *mem = NULL;

    if ((fd = usys_shm_open(name, O_CREAT | O_RDWR,
                            S_IRWXU | S_IRWXG | S_IRWXO)) < 0) {
        usys_log_error("Shared memory open failed. Error %s",
                       usys_error(errno));
        return NULL;
    }

    if ((ret = ftruncate(fd, size)) == -1) {
        usys_log_error("Shared memory ftruncate failed. Error %s",
                       usys_error(errno));
        return NULL;
    }

    if ((mem = mmap(0, size, PROT_EXEC | PROT_READ | PROT_WRITE, MAP_SHARED, fd,
                    0)) == MAP_FAILED) {
        usys_log_error("Shared memory mmap failed. Error %s",
                       usys_error(errno));
        return NULL;
    }

    usys_log_trace(" Allocated shared memory size %lu (%s)",
                   (long unsigned int)size, name);
    return mem;
}

int usys_free_shared_mem(const char *name, void *ptr, uint32_t size) {
    int ret = 0;

    if ((ret = usys_munmap(ptr, size)) < 0) {
        usys_log_error("Could not munmap shared memory %s, ptr %p. Error: %s",
                       name, ptr, usys_error(errno));
        return ret;
    }

    if ((ret = usys_shm_unlink(name)) < 0) {
        usys_log_error("Could not unlink shared memory %s. Error: %s", name,
                       usys_error(errno));
        return ret;
    }
    return ret;
}

void *usys_map_shared_mem(const char *name, uint32_t size) {
    int fd;
    void *mem = NULL;

    if ((fd = usys_shm_open(name, O_RDWR, S_IRWXU | S_IRWXG | S_IRWXO)) < 0) {
        usys_log_error("Could not open shared memory %s. Error %s", name,
                       usys_error(errno));
        return NULL;
    }

    if ((mem = usys_mmap(0, size, PROT_EXEC | PROT_READ | PROT_WRITE,
                         MAP_SHARED, fd, 0)) == MAP_FAILED) {
        usys_log_info("Could not map shared memory %s", name,
                      usys_error(errno));
        return NULL;
    }

    return mem;
}

#ifdef __GNU_FLAG
static inline void *usys_remap_shared_mem(void *old_address, size_t old_size,
                                          size_t new_size) {
    void *mem;

    if ((mem = usys_mremap(old_address, old_size, new_size, MREMAP_MAYMOVE)) ==
        MAP_FAILED) {
        usys_log_info(
            "Could not remap shared memory (old size %d, new size %d). Error: %s",
            (int)old_size, (int)new_size, usys_error(errno));
        return NULL;
    }

    return mem;
}
#endif
