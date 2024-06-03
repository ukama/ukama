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

#define SERVICE_NAME  SERVICE_RLOG
#define STATUS_OK     (0)
#define STATUS_NOK    (-1)

#define STDOUT        1
#define STDERR        2
#define LOG_FILE      3
#define UKAMA_SERVICE 4

#define NORMAL_EXIT      1
#define WEB_SOCKET_FAIL  2
#define WEB_SERVICE_FAIL 3
#define NODED_FAIL       4

#define DEF_LOG_LEVEL           "TRACE"
#define DEF_SERVICE_CLIENT_HOST "localhost"

#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_EP           "/v1/nodeinfo"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_NODE_TYPE          "tower"
#define DEF_OUTPUT             LOG_FILE
#define DEF_LOG_FILE           "/ukama/apps.log"
#define DEF_FLUSH_TIME         5 /*seconds */

#define EP_BS                  "/"
#define API_VERSION            "v1"
#define URL_PREFIX             EP_BS API_VERSION
#define API_RES_EP(RES)        EP_BS RES

#define ENV_BINDING_IP       "ENV_BINDING_IP"
#define DEF_BINDING_IP       "127.0.0.1"

/* various Ukama nodes */
#define UKAMA_TOWER_NODE     "tower"
#define UKAMA_AMPLIFIER_NODE "amplifier"
#define UKAMA_POWER_NODE     "power"

#define LOG_ELEMENTS  6
#define LOG_FORMAT    "%s %s %s %[^:]:%d: %[^\n]"

#define MAX_SIZE        128
#define MAX_URL_LEN     128
#define MAX_MSG_LEN     256
#define MAX_LOG_LEN     512
#define MAX_LOG_BUFFER 8096

#define JTAG_LOGS      "logs"
#define JTAG_APP_NAME  "app_name"
#define JTAG_TIME      "time"
#define JTAG_LEVEL     "level"
#define JTAG_MESSAGE   "message"

typedef struct _u_instance        UInst;
typedef struct _u_request         URequest;
typedef struct _u_response        UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;
typedef json_t                    JsonObj;
typedef json_error_t              JsonErrObj;

typedef struct {

    int             output;    /* STDOUT, STDERR, LOG_FILE, UKAMA */
    int             level;     /* INFO, DEBUG, ERROR */
    int             flushTime; /* flush to output interval */
    time_t          lastWriteTime;
    json_t          *jOutputBuffer;
    int             bufferSize;
    pthread_mutex_t bufferMutex; /* thread-safe the buffer */
} ThreadData;

#endif /* RLOGD_H_ */
