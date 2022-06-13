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

#ifdef __cplusplus
extern "C" {
#endif

#include <usys_types.h>

typedef struct {
    char* name;
    double* value;
    char* units;
} ServiceAttr;

/* Generic notification sent from the NotifyD to remote server */
typedef struct {
    char* serviceName;
    char* notificationType;
    char* nodeId;
    char* nodeType;
    char* severity;
    char* description;
    char* deviceAttr;
    uint32_t epochTime;
} Notification;

/* Noded Specific notifications */
typedef struct {
   char* serviceName;
   char* severity;
   uint32_t epochTime;
   char* moduleID;
   char* deviceName;
   char* deviceDesc;
   char* deviceAttr;
   int   dataType;
   double* deviceAttrValue;
   char* units;
} NodedNotifDetails;

/* Generic notification data */
typedef struct {
   char* serviceName;   /* Service name */
   char* severity;      /* Importance */
   uint32_t epochTime;  /* Time */
   char* description;   /* short description for the notification */
   ServiceAttr* attr;   /* for any value related property */
   char* reason;        /* Service error or Service task accomplishment */
   char* details;       /* Additional details if required */
} ServiceNotifDetails;

#ifdef __cplusplus
}
#endif
#endif /* INC_NOTIFY_H_ */
