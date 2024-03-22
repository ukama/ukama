/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef RLOGD_H_
#define RLOGD_H_

#include "ulfius.h"
#include "jansson.h"

#define SERVICE_NAME           SERVICE_RLOG
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_CLIENT_HOST "localhost"
#define RLOGD_VERSION           "0.0.1"

#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_EP           "/noded/v1/nodeinfo"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE          "tower"

#define EP_BS                  "/"
#define REST_API_VERSION       "v1"
#define URL_PREFIX             EP_BS REST_API_VERSION
#define API_RES_EP(RES)        EP_BS RES

/* various Ukama nodes */
#define UKAMA_TOWER_NODE     "tower"
#define UKAMA_AMPLIFIER_NODE "amplifier"
#define UKAMA_POWER_NODE     "power"

typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

#endif /* RLOGD_H_ */
