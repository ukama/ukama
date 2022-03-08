/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_SHM_H_
#define USYS_SHM_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     int usys_shm_open(const char*, int, mode_t)
 * @brief  creates and opens a new, or opens an existing, POSIX
 *         shared memory object.
 *
 * @param  name
 * @param  oflag
 * @param  mode
 * @return On success a nonnegative file descriptor and On failure, returns -1
 */
static inline int usys_shm_open(const char *name, int oflag, mode_t mode) {
    return shm_open(name, oflag, mode);
}

/**
 * @fn     int usys_shm_unlink(const char*)
 * @brief  Removing an object previously created by usys_shm_open()
 *
 * @param  name
 * @return 0 on success, or -1 on error.
 */
static inline int usys_shm_unlink(const char *name) {
    return shm_unlink(name);
}

/**
 * @fn     int usys_ftruncate(int, off_t)
 * @brief  Cause regular file referenced by fd to be truncated
 *         to a size of precisely length bytes
 *
 * @param  fd
 * @param  length
 * @return On success, zero is returned. On error, -1 is returned.
 */
static inline int usys_ftruncate(int fd, off_t length) {
    return ftruncate(fd, length);
}

/**
 * @fn     void usys_mmap*(void*, size_t, int, int, int, off_t)
 * @brief  Creates a new mapping in the virtual address space of the
 *         calling process
 *
 * @param  addr
 * @param  length
 * @param  prot
 * @param  flags
 * @param  fd
 * @param  offset
 * @return On success returns a pointer to the mapped area.
 *         Onerror, the value MAP_FAILED (that is, (void *) -1) is returned
 */
static inline void *usys_mmap(void *addr, size_t length, int prot, int flags,
                              int fd, off_t offset) {
    return mmap(addr, length, prot, flags, fd, offset);
}

/**
 * @fn     int usys_munmap(void*, size_t)
 * @brief  Deletes the mappings for the specified
 *         address range, and causes further references to addresses within
 *         the range to generate invalid memory references.
 *
 * @param  addr
 * @param  length
 * @return On success, zero is returned. On error, -1 is returned.
 */
static inline int usys_munmap(void *addr, size_t length) {
    return munmap(addr, length);
}

#ifdef __GNU_FLAG
/**
 * @fn     void mremap*(void*, size_t, size_t, int)
 * @brief  expands (or shrinks) an existing memory mapping,
 *         potentially moving it at the same time
 * @param  address
 * @param  length
 * @param  new_length
 * @param  flag
 * @return On success returns a pointer to the new virtual memory area.
 *         On error, the value MAP_FAILED (that is, (void *) -1) is returned.
 */
static inline void *usys_mremap(void *address, size_t length, size_t new_length,
                                int flag) {
    return mremap(address, length, new_length, flag);
}
#endif
#endif /* USYS_SHM_H_ */
