/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WEB_H_
#define WEB_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "ulfius.h"

/* RESPONSE CODE */
#define RESP_CODE_SUCCESS               200
#define RESP_CODE_CREATED               201
#define RESP_CODE_ACCEPTED              202
#define RESP_CODE_INVALID_REQUEST       400
#define RESP_CODE_RESOURCE_NOT_FOUND    404
#define RESP_CODE_SERVER_FAILURE        500

#define METHOD_LENGTH                   7
#define URL_EXT_LENGTH                  64
#define MAX_END_POINTS                  64
#define MAX_URL_LENGTH                  128

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;

#ifdef __cplusplus
}
#endif
#endif /* INC_WEB_H_ */
