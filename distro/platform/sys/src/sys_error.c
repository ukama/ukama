/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sys_error.h"

const char *usysErrorCodes[] = {
    [USYS_ERROR_CODE(ERR_SOCK_CREATION)] = "Failed to create socket",
    [USYS_ERROR_CODE(ERR_SOCK_CONNECT)] = "Failed to connect to socket",
    [USYS_ERROR_CODE(ERR_SOCK_SEND)] = "Failed to send on socket",
    [USYS_ERROR_CODE(ERR_SOCK_RECV)] = "Failed to read from socket",
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
         return usysErrorCodes[err];
        }
        return NULL;
    }
}
