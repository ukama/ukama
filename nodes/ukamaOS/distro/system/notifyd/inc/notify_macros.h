/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef NOTIFY_MACROS_H_
#define NOTIFY_MACROS_H_

#include <pthread.h>
#include "usys_services.h"

#define SERVICE_NAME   SERVICE_NOTIFY
#define STATUS_OK      (0)
#define STATUS_NOK     (-1)

#define MAX_SERVICE_COUNT   (5)

#define STDOUT        1
#define STDERR        2
#define LOG_FILE      3
#define UKAMA_SERVICE 4

/* exit stages */
#define NORMAL_EXIT      1
#define WEB_ADMIN_FAIL   2
#define WEB_SERVICE_FAIL 3
#define NODED_FAIL       4
#define ADDR_FAIL        5

#define DEF_OUTPUT    UKAMA_SERVICE
#define DEF_LOG_FILE  "/ukama/notification.log"
#define DEF_LOG_LEVEL "TRACE"

#define NOTIFICATION_ALERT  "alert"
#define NOTIFICATION_EVENT  "event"

#define LOG_ELEMENTS   6
#define LOG_FORMAT     "%s %s %s %[^:]:%d: %[^\n]"

#define DEF_NODED_HOST "localhost"
#define DEF_NODED_EP   "/v1/nodeinfo"
#define DEF_REMOTE_EP  "node/v1/notify"
#define DEF_NODE_ID    "ukama-aaa-bbbb-ccc-dddd"
#define DEF_MAP_FILE   "status.map"

#define ENV_NOTIFY_DEBUG_MODE  "NOTIFYD_DEBUG_MODE"

typedef struct {

    int             output; /* STDOUT, STDERR, LOG_FILE, UKAMA */
    int             count;  /* successful notification */
    pthread_mutex_t mutex;
} ThreadData;

#endif /* INC_NOTIFY_MACROS_H_ */
