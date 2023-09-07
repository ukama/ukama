/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_MEM_H
#define USYS_MEM_H

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     void usys_malloc*(size_t)
 * @brief  Allocate memory of size bytes
 *
 * @param  size
 * @return On Success start of the memory address.
 *         On failure NULL
 */
void *usys_malloc(size_t size);

/**
  * @fn    void usys_free(void*)
  * @brief Free memory allocated by usys_malloc
  *
  * @param ptr
  */
void usys_free(void *ptr);

/**
 * @fn     void usys_realloc*(void*, size_t)
 * @brief  Reallocates the given area of memory. It must be previously
 *         allocated by usys_malloc(), usys_calloc() or usys_realloc()
 *         and not yet freed with a call to free or realloc.
 *         Otherwise, the results are undefined.
 *
 * @param  ptr
 * @param  new_size
 * @return On success, returns the pointer to the beginning of newly
 *         allocated memory.
 *         On failure, returns a null pointer.
 */
void *usys_realloc(void *ptr, size_t new_size);

/**
 * @fn     void calloc*(size_t, size_t)
 * @brief  Allocates memory for an array of num objects of size and
 *         initializes all bytes in the allocated storage to zero.
 *
 * @param  num
 * @param  size
 * @return On success, returns the pointer to the beginning of newly
 *         allocated memory.
 *         On failure, returns a null pointer.
 */
void *usys_calloc(size_t num, size_t size);

/**
 * @fn     void usys_emalloc*(size_t)
 * @brief  Wrapper on usys_calloc function which exits the calling process
 * 	       on failure.
 *
 * @param  size
 * @return On success base address of memory allocated.
 * 		   On error call usys_exit(errno)
 */
void *usys_emalloc(size_t size);

/**
 * @fn     void usys_erealloc*(void*, size_t)
 * @brief  Wrapper on usys_realloc function which exits the process
 * 	       on failure.
 *
 * @param  ptr
 * @param  new_size
 * @return On success base address of memory allocated.
 * 		   On error call usys_exit(errno)
 */
void *usys_erealloc(void *ptr, size_t new_size);

/**
 * @fn     void usys_ecalloc*(size_t, size_t)
 * @brief  Wrapper on usys_realloc function which exits the process
 * 	       on failure.
 * @param  num
 * @param  size
 * @return On success base address of memory allocated.
 * 		   On error call usys_exit(errno)
 */
void *usys_ecalloc(size_t num, size_t size);

/**
 * @fn     void usys_zmalloc*(size_t)
 * @brief  Wrapper on usys_malloc function which initializes all allocated
 *         bytes to zero.
 *
 * @param  size
 * @return On success base address of memory allocated.
 * 		   On error call usys_exit(errno)
 */
void *usys_zmalloc(size_t size);

#ifdef __cplusplus
}
#endif

#endif /* USYS_MEM_H */
