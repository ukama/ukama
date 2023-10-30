/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
    char     *device;
    char     *propertyName;
    char     *propertyValue;
    char     *propertyUnit;
    char     *details;
} Notification;

#endif /* INC_NOTIFY_H_ */
