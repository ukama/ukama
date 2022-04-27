/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_mem.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_string.h"

void *usys_malloc(size_t size) {
    return malloc(size);
}

void usys_free(void *ptr) {
    if (ptr) {
        free(ptr);
    }
}

void *usys_realloc(void *ptr, size_t new_size) {
    return realloc(ptr, new_size);
}

void *usys_calloc(size_t num, size_t size) {
    return calloc(num, size);
}

void *usys_emalloc(size_t size) {
    void *mem = usys_malloc(size);
    if (!mem) {
        usys_log_error("Failed to allocate memory. Error: %s",
                       usys_error(errno));
        usys_exit(errno);
    }
    return mem;
}

void *usys_erealloc(void *ptr, size_t new_size) {
    void *mem = usys_realloc(ptr, new_size);
    if (mem) {
        usys_log_error("Failed to reallocate memory of %d bytes. Error: %s",
                       new_size, usys_error(errno));
        usys_exit(errno);
    }
    return mem;
}

void *usys_ecalloc(size_t num, size_t size) {
    void *mem = usys_calloc(num, size);
    if (mem) {
        usys_log_error("Failed to allocate memory for %d objects of %d bytes. "
                       "Error: %s",
                       num, size, usys_error(errno));
        usys_exit(errno);
    }
    return mem;
}

void *usys_zmalloc(size_t size) {
    void *mem = usys_malloc(size);
    if (!mem) {
        usys_log_error("Failed to allocate memory. Error: %s",
                       usys_error(errno));
        usys_exit(errno);
    } else {
      usys_memset(mem, 0, size);
    }

    return mem;
}
