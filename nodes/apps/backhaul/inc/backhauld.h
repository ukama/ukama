/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef BACKHAULD_H_
#define BACKHAULD_H_

#include "ulfius.h"
#include "jansson.h"
#include "usys_types.h"
#include "usys_services.h"
#include "usys_log.h"

#define SERVICE_NAME             "backhaul"
#define STATUS_OK                (0)
#define STATUS_NOK               (-1)

#define DEF_LOG_LEVEL            "TRACE"

#define EP_BS                    "/"
#define REST_API_VERSION         "v1"
#define URL_PREFIX               EP_BS REST_API_VERSION
#define API_RES_EP(RES)          EP_BS RES

#define ENV_BACKHAULD_DEBUG_MODE "BACKHAULD_DEBUG_MODE"

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

#endif /* BACKHAULD_H_ */
