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

#define ENV_LOOKOUT_APP_MANAGER    "LOOKOUT_APP_MANAGER"
#define LOOKOUT_MANAGER_STARTER    "starter"
#define LOOKOUT_MANAGER_SUPERVISOR "supervisor"

#define LOOKOUT_STATUS_NA      "not-available"
#define LOOKOUT_GPS_COORD_NA   "-999.000000,-999.000000"
#define LOOKOUT_GPS_TIME_NA    "not-available"
#define LOOKOUT_SCHEMA_VERSION "1.0"

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

    char   *status;
    pid_t  pid;
    int    memory;
    int    disk;
    double cpu;
} CappRuntime;

typedef struct _capp {

    char        *space;
    char        *name;
    char        *tag;
    CappRuntime *runtime;
} Capp;

typedef struct _cappList {

    Capp             *capp;
    struct _cappList *next;
} CappList;

typedef struct {

    bool  available;
    char *state;
    bool  updateInProgress;
    bool  switchRequested;
    bool  terminateRequested;
    int   exitCode;
} StarterStatusData;

typedef struct {

    bool  available;
    bool  ok;
    char *board;
    char *reason;
    double totalWatts;
    double temperatureC;
} PowerStatusData;

typedef struct {

    bool  available;
    bool  lock;
    char *coordinates;
    char *time;
} GPSClientData;

typedef struct {

    bool available;
    char *state;
} RadioStatusData;

typedef struct {

    bool available;
    char *service;
    char *error;
} CellularStatusData;

typedef struct {

    bool available;
    JsonObj *status;
    JsonObj *policy;
    JsonObj *ports;
} SwitchStatusData;

typedef struct {

    bool available;
    JsonObj *status;
} ControllerStatusData;

typedef struct {

    bool available;
    JsonObj *status;
} BackhaulStatusData;

typedef struct {

    bool available;
    JsonObj *status;
} FemStatusData;

typedef struct {

    StarterStatusData   starter;
    PowerStatusData     power;
    GPSClientData       gps;
    RadioStatusData     radio;
    CellularStatusData  cellular;
    SwitchStatusData    sw;
    ControllerStatusData controller;
    BackhaulStatusData  backhaul;
    FemStatusData       fem;
} LookoutStatusData;

#endif /* LOOKOUT_H_ */
