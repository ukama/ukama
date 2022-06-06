/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_NOTIFY_H_
#define INC_NOTIFY_H_

#include <usys_types.h>

typedef struct {
    char* serviceName;
    char* notificationType;
    char* nodeId;
    char* nodeType;
    char* severity;
    char* description;
    char* deviceAttr;
    uint32_t epcohTime;
} Notification;

typedef struct {
   char* serviceName;
   char* severity;
   uint32_t epcohTime;
   char* moduleID;
   char* deviceName;
   char* deviceDesc;
   char* deviceAttr;
   char* dataType;
   double* deviceAttrValue;
   char* units;
} NodedNotifDetails;





#endif /* INC_NOTIFY_H_ */
