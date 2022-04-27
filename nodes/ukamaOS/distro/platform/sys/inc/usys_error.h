/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_SYS_ERROR_H
#define USYS_SYS_ERROR_H

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

#define USYS_ERR_BASE_PLTF_CODE  (1000)
#define USYS_ERR_PLTF_CODE(code) ((code)-USYS_ERR_BASE_PLTF_CODE)
#define USYS_ERR_APP_BASE_CODE   (2000)

/*
*
* Error codes = USYS_ERR_PLTF_CODE_IDX_BASE + USysErrCodeIdx
*/
typedef enum {
    ERR_PLTF_NONE = 0,
    /* Sample error code */
    ERR_PLTF_SOCK_CREATION = (USYS_ERR_BASE_PLTF_CODE + 1),
    ERR_PLTF_SOCK_CONNECT,
    ERR_PLTF_SOCK_SEND,
    ERR_PLTF_SOCK_RECV,

    /* POSIX ERROR */
    ERR_PLTF_MUTEX_OBJ_NULL,
    ERR_PLTF_MUTEX_ATTR_INIT_FAIL,
    ERR_PLTF_MUTEX_ATTR_SET_PROTO_FAIL,
    ERR_PLTF_MUTEX_ATTR_SET_TYPE_FAIL,
    ERR_PLTF_MUTEX_INIT_FAILED,
    ERR_PLTF_MUTEX_LOCK_FAILED,
    ERR_PLTF_MUTEX_TRYLOCK_FAILED,
    ERR_PLTF_MUTEX_TIMEDLOCK_FAILED,
    ERR_PLTF_MUTEX_UNLOCK_FAILED,
    ERR_PLTF_MUTEX_DESTROY_FAILED,

    ERR_PLTF_SEM_OBJ_NULL,
    ERR_PLTF_SEM_INIT_FAILURE,
    ERR_PLTF_SEM_WAIT_FAIL,
    ERR_PLTF_SEM_TRYWAIT_FAIL,
    ERR_PLTF_SEM_TIMEDWAIT_FAIL,
    ERR_PLTF_SEM_POST_FAIL,
    ERR_PLTF_SEM_DESTROY_FAIL,

    ERR_PLTF_SPIN_LOCK_INIT_FAILED,
    ERR_PLTF_SPIN_LOCK_LOCK_FAILED,
    ERR_PLTF_SPIN_LOCK_UNLOCK_FAILED,
    ERR_PLTF_SPIN_LOCK_DESTROY_FAILED,
    ERR_PLTF_MAX_ERR_CODE
} USysErrCodeIdx;

typedef USysErrCodeIdx USysError;

#define ERR_NONE    ERR_PLTF_NONE

const char *usys_error(int err);

#ifdef __cplusplus
extern "C" {
#endif
#endif /* USYS_SYS_ERROR_H */
