/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef DEVICED_H_
#define DEVICED_H_

#include "ulfius.h"
#include "usys_types.h"
#include "usys_log.h"
#include "jansson.h"

#define SERVICE_NAME           "deviced"
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_PORT        "8086"
#define DEF_SERVICE_CLIENT_PORT "8087"
#define DEF_SERVICE_CLIENT_HOST "localhost"
#define DEVICED_VERSION         "0.0.1"

#define DEF_NODED_HOST         "localhost"
#define DEF_NOTIFY_HOST        "localhost"
#define DEF_NOTIFY_PORT        "8085"
#define DEF_NODED_PORT         "8095"
#define DEF_NODED_EP           "/noded/v1/nodeinfo"
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

#endif /* DEVICED_H_ */
