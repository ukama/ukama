/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef LOOKOUT_H_
#define LOOKOUT_H_

#include "ulfius.h"
#include "usys_types.h"
#include "usys_services.h"
#include "usys_log.h"
#include "jansson.h"

#define SERVICE_NAME           SERVICE_LOOKOUT
#define SERVICE_NODED          1
#define SERVICE_STARTERD       2
#define SERVICE_NOTIFYD        3

#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_REPORT_INTERVAL     30

#define DEF_NODED_HOST         "localhost"
#define DEF_STARTERD_HOST      "localhost"
#define DEF_NODE_SYSTEM_EP     "v1/nodes/%s/status"
#define DEF_NODED_EP           "/v1/nodeinfo"
#define DEF_STARTERD_EP        "v1/status"
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

    char        *space;    /* space capp belongs */
    char        *name;     /* Name of the cApp */
    char        *tag;      /* cApp tag/version */
    CappRuntime *runtime;  /* runtime of capp */
} Capp;

typedef struct _cappList {

    Capp             *capp;
    struct _cappList *next;
} CappList;


typedef struct {

    bool  gpsLock;
    char *coordinates; /* "lon,lat" */
    char *gpsTime;     /* string */
} GPSClientData;

#endif /* LOOKOUT_H_ */
