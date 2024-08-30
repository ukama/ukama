/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
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
