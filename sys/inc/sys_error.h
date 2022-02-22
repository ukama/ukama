/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_SYS_ERROR_H
#define USYS_SYS_ERROR_H

#ifdef __cplusplus
extern "C" {
#endif

#include "sys_types.h"

#define   USYS_BASE_ERROR_CODE      (1000)
#define   USYS_ERROR_CODE(code)     ((code) - USYS_BASE_ERROR_CODE)

/*
*
* Error codes = USYS_ERROR_CODE_IDX_BASE + USysErrorCodeIdx
*/
typedef enum {
    ERR_NONE = 0,
    /* Sample error code */
    ERR_SOCK_CREATION = (USYS_BASE_ERROR_CODE+1),
    ERR_SOCK_CONNECT,
    ERR_SOCK_SEND,
    ERR_SOCK_RECV,

    /* POSIX ERROR */
    ERR_MUTEX_OBJ_NULL,
    ERR_MUTEX_ATTR_INIT_FAIL,
    ERR_MUTEX_ATTR_SET_PROTO_FAIL,
    ERR_MUTEX_ATTR_SET_TYPE_FAIL,
    ERR_MUTEX_INIT_FAILED,
    ERR_MUTEX_LOCK_FAILED,
    ERR_MUTEX_TRYLOCK_FAILED,
    ERR_MUTEX_TIMEDLOCK_FAILED,
    ERR_MUTEX_UNLOCK_FAILED,
    ERR_MUTEX_DESTROY_FAILED,

    ERR_SEM_OBJ_NULL,
    ERR_SEM_INIT_FAILURE,
    ERR_SEM_WAIT_FAIL,
    ERR_SEM_TRYWAIT_FAIL,
    ERR_SEM_TIMEDWAIT_FAIL,
    ERR_SEM_POST_FAIL,
    ERR_SEM_DESTROY_FAIL,

    ERR_SPIN_LOCK_INIT_FAILED,
    ERR_SPIN_LOCK_LOCK_FAILED,
    ERR_SPIN_LOCK_UNLOCK_FAILED,
    ERR_SPIN_LOCK_DESTROY_FAILED,
    ERR_MAX_ERROR_CODE
} USysErrorCodeIdx;

typedef  USysErrorCodeIdx USysError;

const char* usys_error(int err);

#ifdef __cplusplus
extern "C" {
#endif
#endif /* USYS_SYS_ERROR_H */
