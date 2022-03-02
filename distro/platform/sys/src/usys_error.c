/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_error.h"

const char *usysErrCodes[] = {
    [USYS_ERR_CODE(ERR_SOCK_CREATION)] = "failed to create socket",
    [USYS_ERR_CODE(ERR_SOCK_CONNECT)] = "failed to connect to socket",
    [USYS_ERR_CODE(ERR_SOCK_SEND)] = "failed to send on socket",
    [USYS_ERR_CODE(ERR_SOCK_RECV)] = "failed to read from socket",
    [USYS_ERR_CODE(ERR_MUTEX_OBJ_NULL)] = "failed to create mutex",
    [USYS_ERR_CODE(ERR_MUTEX_ATTR_INIT_FAIL)] =
        "failed to initialize mutex attributes",
    [USYS_ERR_CODE(ERR_MUTEX_ATTR_SET_PROTO_FAIL)] =
        "failed to set protocol for mutex",
    [USYS_ERR_CODE(ERR_MUTEX_ATTR_SET_TYPE_FAIL)] =
        "failed to set mutex attributes",
    [USYS_ERR_CODE(ERR_MUTEX_INIT_FAILED)] = "failed to initialize mutex",
    [USYS_ERR_CODE(ERR_MUTEX_LOCK_FAILED)] = "mutex lock failed",
    [USYS_ERR_CODE(ERR_MUTEX_TRYLOCK_FAILED)] = "mutex try lock failed",
    [USYS_ERR_CODE(ERR_MUTEX_TIMEDLOCK_FAILED)] = "mutex timed lock failed",
    [USYS_ERR_CODE(ERR_MUTEX_UNLOCK_FAILED)] = "mutex unlock failed",
    [USYS_ERR_CODE(ERR_MUTEX_DESTROY_FAILED)] = "mutex destroy failed",
    [USYS_ERR_CODE(ERR_SEM_OBJ_NULL)] = "failed to create semaphore",
    [USYS_ERR_CODE(ERR_SEM_INIT_FAILURE)] = "failed to initialize semaphore",
    [USYS_ERR_CODE(ERR_SEM_WAIT_FAIL)] = "semaphore wait failed",
    [USYS_ERR_CODE(ERR_SEM_TRYWAIT_FAIL)] = "semaphore trywait failed",
    [USYS_ERR_CODE(ERR_SEM_TIMEDWAIT_FAIL)] = "semaphore timedwait failed",
    [USYS_ERR_CODE(ERR_SEM_POST_FAIL)] = "semaphore post failed",
    [USYS_ERR_CODE(ERR_SEM_DESTROY_FAIL)] = "semaphore destroy failed",
    [USYS_ERR_CODE(ERR_SPIN_LOCK_INIT_FAILED)] = "spinlock init failed",
    [USYS_ERR_CODE(ERR_SPIN_LOCK_LOCK_FAILED)] = "spinlock lock failed",
    [USYS_ERR_CODE(ERR_SPIN_LOCK_UNLOCK_FAILED)] = "spinlock unlock failed",
    [USYS_ERR_CODE(ERR_SPIN_LOCK_DESTROY_FAILED)] = "spinlock destroy failed",
};

/**
 * @brief  Read error description
 *
 * @param  err
 * @return char*
 *
 */
const char *usys_error(int err) {
    if (err < USYS_BASE_ERR_CODE) {
        /* TBU: check if we can use some thread safe api like strerror_r */
        return strerror(err);
    } else {
        if (err < ERR_MAX_ERR_CODE) {
            return usysErrCodes[err - USYS_BASE_ERR_CODE];
        }
        return NULL;
    }
}
