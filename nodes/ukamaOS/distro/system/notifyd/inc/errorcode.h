/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_ERRCODE_H_
#define INC_ERRCODE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_error.h"

/* Error codes returned by NodeD service */
typedef enum {
    ERR_APP_BASE    = (USYS_ERR_APP_BASE_CODE+1),
    ERR_JSON_PARSER,
    ERR_JSON_CRETATION_ERR,
    ERR_JSON_NO_VAL_TO_ENCODE,
    ERR_JSON_INVALID,
    ERR_JSON_UNEXPECTED_TAG,
    ERR_JSON_BAD_REQ,
    ERR_APP_MAX
} ErrorCode;

#ifdef __cplusplus
}
#endif

#endif /* INC_ERRCODE_H_*/
