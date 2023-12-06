/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
    ERR_JSON_CREATION_ERR,
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
