/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_NOTIFY_H_
#define INC_NOTIFY_H_

#include <usys_types.h>

/* JSON:
   {
     "serviceName": "abc",
     "time": 1654566750,
     "severity": "high",
     "property-name": ActiveUser
     "property-value": "64"
     "property-units" : "integer"
     "reason" : "Too many users"
     "details" : "Use exeeding limits"
   }
*/

typedef struct {

    char     *serviceName;
    char     *severity;      /* low, medium, high */
    int      epochTime;
    char     *module;
    char     *propertyName;
    char     *propertyValue;
    char     *propertyUnit;
    char     *details;
} Notification;

#endif /* INC_NOTIFY_H_ */
