/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef GPSD_H_
#define GPSD_H_

#include <pthread.h>

#include "config.h"

#include "ulfius.h"
#include "usys_types.h"
#include "usys_services.h"
#include "usys_log.h"
#include "jansson.h"

#define MAX_BUFFER        512
#define MAX_LINE_LENGTH   256
#define MAX_URL_LENGTH    128
#define MAX_LAT_LONG_SIZE 32

#define SERVICE_NAME           SERVICE_GPS
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define GPS_SCRIPT             "/sbin/process_gps_data.sh"
#define GPS_LOC_FILE           "/tmp/gps_loc.log"
#define GPS_RAW_FILE           "/tmp/gps_raw.log"
#define GPS_WAIT_TIME          5

#define DEF_LOG_LEVEL          "TRACE"
#define DEF_GPS_MODULE_HOST    "localhost"

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

typedef struct _u_instance  UInst;
typedef struct _u_instance  UInst;
typedef struct _u_request   URequest;
typedef struct _u_response  UResponse;
typedef json_t              JsonObj;
typedef json_error_t        JsonErrObj;

typedef struct {

    bool gpsLock;
    char *time;
    char *latitude;
    char *longitude;

    pthread_mutex_t mutex;
} GPSData;

bool start_gps_data_collection_and_processing(Config *config, pthread_t *tid);
void stop_gps_data_collection_and_processing(pthread_t tid);

#endif /* GPSD_H_ */
