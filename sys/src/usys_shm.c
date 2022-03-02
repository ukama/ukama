/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_shm.h"

#include "usys_error.h"
#include "usys_log.h"

/**
 * @fn     void usys_allocate_shared_mem*(const char*, uint32_t)
 * @brief  Creates a shared memory using usys_shm_open, usys_ftruncate and usys_mmap api's.
 *
 * @param  name
 * @param  size
 * @return On Success memmory address for shared memory
 *         On failure NULL
 */
void* usys_allocate_shared_mem(const char* name, uint32_t size)
{
    int fd, ret;
    void* mem = NULL;

    if ((fd = usys_shm_open(name, O_CREAT|O_RDWR, S_IRWXU|S_IRWXG|S_IRWXO)) < 0)
    {
        usys_log_error("Shared memory open failed. Error %s", usys_error(errno));
        return NULL;
    }

    if ((ret = ftruncate(fd, size)) == -1)
    {
        usys_log_error("Shared memory ftruncate failed. Error %s", usys_error(errno));
        return NULL;
    }

    if ((mem = mmap(0, size, PROT_EXEC|PROT_READ|PROT_WRITE, MAP_SHARED, fd, 0)) == MAP_FAILED)
    {
        usys_log_error("Shared memory mmap failed. Error %s", usys_error(errno));
        return NULL;
    }

    usys_log_trace(" Allocated shared memory size %lu (%s)",(long unsigned int)size, name);
    return mem;
}

/**
 * @fn     int usys_free_shared_mem(const char*, void*, uint32_t)
 * @brief  Use functions usys_munmap and usys_shm_unlink to free shared memory
 *
 * @param  name
 * @param  ptr
 * @param  size
 * @return On Success return 0.
 *         On failure -1
 */
int usys_free_shared_mem(const char* name, void* ptr, uint32_t size)
{
    int ret = 0;

    if ((ret = usys_munmap(ptr, size)) < 0)
    {
        usys_log_error("Could not munmap shared memory %s, ptr %p. Error: %s", name, ptr, usys_error(errno));
        return ret;
    }

    if ((ret = usys_shm_unlink(name)) < 0)
    {
        usys_log_error("Could not unlink shared memory %s. Error: %s", name, usys_error(errno));
        return ret;
    }
    return ret;
}

/**
 * @fn     void usys_map_shared_mem*(const char*, uint32_t)
 * @brief  This function usese usys_shm_open and usys_mmap to map a shared memory created
 *
 * @param  name
 * @param  size
 * @return On Success memmory address for shared memory
 *         On failure NULL
 */
void* usys_map_shared_mem(const char* name, uint32_t size)
{
    int fd;
    void* mem = NULL;

    if ((fd = usys_shm_open(name, O_RDWR, S_IRWXU|S_IRWXG|S_IRWXO)) < 0)
    {
        usys_log_error("Could not open shared memory %s. Error %s", name, usys_error(errno));
        return NULL;
    }

    if ((mem = usys_mmap(0, size, PROT_EXEC|PROT_READ|PROT_WRITE, MAP_SHARED, fd, 0)) == MAP_FAILED)
    {
        usys_log_info("Could not map shared memory %s",name, usys_error(errno));
        return NULL;
    }

    return mem;
}

#ifdef __GNU_FLAG
/**
 * @fn     void usys_remap_shared_mem*(void*, size_t, size_t)
 * @brief  This function usys_mremap to remap the shared memory mapped using usys_map_shared_mem
 *         with a different size
 *
 * @param  old_address
 * @param  old_size
 * @param  new_size
 * @return On success returns a pointer to the new virtual memory area.
 *         On error, the value MAP_FAILED (that is, (void *) -1) is returned.
 */
static inline void* usys_remap_shared_mem(void* old_address, size_t old_size, size_t new_size)
{
    void*  mem;

    if ((mem = usys_mremap(old_address, old_size, new_size, MREMAP_MAYMOVE)) == MAP_FAILED)
    {
        usys_log_info("Could not remap shared memory (old size %d, new size %d). Error: %s", (int)old_size,
               (int)new_size, usys_error(errno));
        return NULL;
    }

    return mem;
}
#endif






