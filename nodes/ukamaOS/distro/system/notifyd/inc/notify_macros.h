/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef NOTIFY_MACROS_H_
#define NOTIFY_MACROS_H_

#include "usys_services.h"

#define SERVICE_NAME           SERVICE_NOTIFY
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define MAX_SERVICE_COUNT      (32)

#define NOTIFICATION_ALERT     "alert"
#define NOTIFICATION_EVENT     "event"

#define DEF_LOG_LEVEL          "TRACE"
#define NOTIFY_VERSION         "0.0.0"

#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_EP           "/noded/v1/nodeinfo"
#define DEF_REMOTE_EP          "/notification"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_MAP_FILE           "status.map"

#define ENV_NOTIFY_DEBUG_MODE  "NOTIFYD_DEBUG_MODE"

#endif /* INC_NOTIFY_MACROS_H_ */
