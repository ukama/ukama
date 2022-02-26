/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "usys_error.h"

const char *usysErrorCodes[] = {
    [USYS_ERROR_CODE(ERR_SOCK_CREATION)] = "Failed to create socket",
    [USYS_ERROR_CODE(ERR_SOCK_CONNECT)] = "Failed to connect to socket",
    [USYS_ERROR_CODE(ERR_SOCK_SEND)] = "Failed to send on socket",
    [USYS_ERROR_CODE(ERR_SOCK_RECV)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_OBJ_NULL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_ATTR_INIT_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_ATTR_SET_PROTO_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_ATTR_SET_TYPE_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_INIT_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_LOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_TRYLOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_TIMEDLOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_UNLOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_MUTEX_DESTROY_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_OBJ_NULL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_INIT_FAILURE)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_WAIT_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_TRYWAIT_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_TIMEDWAIT_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_POST_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SEM_DESTROY_FAIL)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SPIN_LOCK_INIT_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SPIN_LOCK_LOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SPIN_LOCK_UNLOCK_FAILED)] = "Failed to read from socket",
    [USYS_ERROR_CODE(ERR_SPIN_LOCK_DESTROY_FAILED)] = "Failed to read from socket",
};

/**
 * @brief Read error description
 *
 * @param err
 * @return char*
 */
const char* usys_error(int err) {
    if (err < USYS_BASE_ERROR_CODE) {
        return strerror(err);
    } else {
        if (err < ERR_MAX_ERROR_CODE) {
         return usysErrorCodes[err-USYS_BASE_ERROR_CODE];
        }
        return NULL;
    }
}
