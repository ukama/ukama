/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef DEVICED_H_
#define DEVICED_H_

#include "ulfius.h"
#include "usys_types.h"
#include "usys_services.h"
#include "usys_log.h"
#include "jansson.h"

#define SERVICE_NAME           SERVICE_DEVICE
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_CLIENT_HOST "localhost"

#define DEF_NODED_HOST         "localhost"
#define DEF_NOTIFY_HOST        "localhost"
#define DEF_NODED_EP           "/v1/nodeinfo"
#define DEF_NOTIFY_EP          "/notify/v1/event/"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE          "tower"
#define ENV_DEVICED_DEBUG_MODE "DEVICED_DEBUG_MODE"

#define EP_BS                  "/"
#define REST_API_VERSION       "v1"
#define URL_PREFIX             EP_BS REST_API_VERSION
#define API_RES_EP(RES)        EP_BS RES
#define WAIT_BEFORE_REBOOT     5 /* seconds */

/* for json de/ser */
#define ALARM_HIGH             "high"
#define ALARM_NODE             "node"
#define ALARM_REBOOT           "reboot"
#define ALARM_REBOOT_DESCRP    "Rebooting the node"
#define EMPTY_STRING           ""
#define MODULE_NONE            "none"

/* various Ukama nodes */
#define UKAMA_TOWER_NODE     "tower"
#define UKAMA_AMPLIFIER_NODE "amplifier"
#define UKAMA_POWER_NODE     "power"

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct {

    void* config;
    int   retCode;
} ThreadArgs;

#endif /* DEVICED_H_ */
