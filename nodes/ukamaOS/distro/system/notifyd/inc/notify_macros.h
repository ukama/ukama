/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef NOTIFY_MACROS_H_
#define NOTIFY_MACROS_H_

#define SERVICE_NAME           "notifyD"
#define STATUS_OK              (0)
#define STATUS_NOK             (-1)

#define MAX_SERVICE_COUNT      (32)

#define NOTIFICATION_ALERT     "alert"
#define NOTIFICATION_EVENT     "event"

#define DEF_LOG_LEVEL          "TRACE"
#define DEF_SERVICE_PORT       "8085"
#define NOTIFY_VERSION         "0.0.0"

#define DEF_NODED_HOST         "localhost"
#define DEF_NODED_PORT         "8095"
#define DEF_NODED_EP           "/noded/v1/nodeinfo"
#define DEF_REMOTE_SERVER      "http://localhost:8091"
#define DEF_REMOTE_EP          "/notification"
#define DEF_NODE_ID            "ukama-aaa-bbbb-ccc-dddd"
#define DEF_MAP_FILE           "status.map"

#define ENV_NOTIFY_DEBUG_MODE  "NOTIFY_DEBUG_MODE"

#endif /* INC_NOTIFY_MACROS_H_ */
