/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef STARTER_H_
#define STARTER_H_

#include "ulfius.h"
#include "usys_types.h"
#include "usys_log.h"
#include "jansson.h"
#include "capp_config.h"

#define SERVICE_NAME           "starter.d"
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_PORT        "8086"
#define STARTER_VERSION         "0.0.1"

#define DEF_NODED_HOST         "localhost"
#define DEF_NOTIFY_HOST        "localhost"
#define DEF_WIMC_HOST          "localhost"
#define DEF_NOTIFY_PORT        "8085"
#define DEF_NODED_PORT         "8095"
#define DEF_WIMC_PORT          "8075"
#define DEF_NODED_EP           "/noded/v1/nodeinfo"
#define DEF_NOTIFY_EP          "/notify/v1/event/"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE          "tower"
#define DEF_MANIFEST_FILE      "manifest.json"
#define ENV_STARTER_DEBUG_MODE "STARTER_DEBUG_MODE"

/* for spaces and capps */
#define DEF_CAPP_PATH          "/capps/pkgs"
#define DEF_SPACE_ROOTFS_PATH  "/capps/rootfs"
#define DEF_CAPP_CONFIG_FILE   "config.json"

/* runtime status */
#define CAPP_RUNTIME_NO_EXEC   0
#define CAPP_RUNTIME_EXEC      1
#define CAPP_RUNTIME_DONE      2

#define EP_BS                  "/"
#define REST_API_VERSION       "v1"
#define URL_PREFIX             EP_BS REST_API_VERSION
#define API_RES_EP(RES)        EP_BS RES

/* various Ukama nodes */
#define UKAMA_TOWER_NODE     "tower"
#define UKAMA_AMPLIFIER_NODE "amplifier"
#define UKAMA_POWER_NODE     "power"

#define MAX_BUFFER          3072
#define SPACE_MAX_BUFFER    1024
#define CAPP_MAX_BUFFER     1024

#define SPACE_BOOT   "boot"
#define SPACE_REBOOT "reboot"

#define CAPP_PKG_NOT_FOUND 0
#define CAPP_PKG_FOUND     1

/* number of second to wait and retry on the capp pkg */
#define FETCH_AND_UPDATE_RETRY 10

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct _runtime {

    char  *cmd;
    char  **argv;
    char  **env;
    
    int   status; /* Current status 0: not executed, 1: running, 2: done */ 
    pid_t pid;    /* process ID */
} CappRuntime;

typedef struct _capp {

    char        *name;     /* Name of the cApp */
    char        *tag;      /* cApp tag/version */
    char        *rootfs;   /* Location where the rootfs is at */
    char        *space;    /* group it belongs to */
    int         restart;   /* 1: yes, always restart. 0: No */
    int         fetch;     /* fetch from hub? */
    CappRuntime *runtime;  /* runtime of capp */
    CappConfig  *config;   /* configuration of the capp */
} Capp;

typedef struct _cappList {

    Capp             *capp;
    struct _cappList *next;
} CappList;

typedef struct _space {

    char     *name;
    char     *rootfs;
    CappList *cappList;
} Space;

typedef struct _spaceList {

    Space             *space;
    struct _spaceList *next;
} SpaceList;

#endif /* STARTER_H_ */
