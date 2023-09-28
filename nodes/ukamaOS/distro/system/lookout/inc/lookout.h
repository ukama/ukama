/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef LOOKOUT_H_
#define LOOKOUT_H_

#include "ulfius.h"
#include "usys_types.h"
#include "usys_log.h"
#include "jansson.h"

#define SERVICE_NAME           "lookout.d"
#define SERVICE_NODED          1
#define SERVICE_STARTERD       2
#define SERVICE_NOTIFYD        3

#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_PORT        "8091"
#define LOOKOUT_VERSION         "0.0.1"
#define DEF_REPORT_INTERVAL     30

#define DEF_NODED_HOST         "localhost"
#define DEF_STARTERD_HOST      "localhost"
#define DEF_NODE_SYSTEM_HOST   "localhost"
#define DEF_NODED_PORT         "8095"
#define DEF_STARTERD_PORT      "8086"
#define DEF_NODE_SYSTEM_PORT   "7001"
#define DEF_NODE_SYSTEM_EP     "v1/nodes/%s/status"
#define DEF_NODED_EP           "noded/v1/nodeinfo"
#define DEF_STARTERD_EP        "status"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define ENV_LOOKOUT_DEBUG_MODE "LOOKOUT_DEBUG_MODE"

#define EP_BS              "/"
#define REST_API_VERSION   "v1"
#define URL_PREFIX         EP_BS REST_API_VERSION
#define API_RES_EP(RES)    EP_BS RES

#define MAX_BUFFER         256

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct _runtime {

    char   *status; /* Current status 0: not executed, 1: running, 2: done */
    pid_t  pid;     /* process ID */
    int    memory;  /* memory usage, in bytes */
    int    disk;    /* disk usage, in bytes */
    double cpu;     /* CPU usage, in percentage */
} CappRuntime;

typedef struct _capp {

    char        *name;     /* Name of the cApp */
    char        *tag;      /* cApp tag/version */
    CappRuntime *runtime;  /* runtime of capp */
} Capp;

typedef struct _cappList {

    Capp             *capp;
    struct _cappList *next;
} CappList;

#endif /* LOOKOUT_H_ */
