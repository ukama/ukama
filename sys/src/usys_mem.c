/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_mem.h"

void* usys_malloc(size_t size) {
    return malloc(size);
}

void usys_free(void* ptr) {
    free(ptr);
}

void* usys_realloc( void *ptr, size_t new_size ) {
    return realloc(ptr, new_size);
}

void* usys_calloc( size_t num, size_t size ) {
    return calloc(num, size);
}



