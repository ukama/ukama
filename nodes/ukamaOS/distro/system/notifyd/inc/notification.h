/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_NOTIFICATION_H_
#define INC_NOTIFICATION_H_

#include "config.h"
#include "jserdes.h"
#include "notify/notify.h"

typedef int (*ServiceHandler)(JsonObj* json, char* type, Config *config);

typedef struct {
    char *service;
    ServiceHandler alertHandler;
    ServiceHandler eventHandler;
} NotifyHandler;


int notify_process_incoming_notification(const char* service, char* notif,
                                         JsonObj* json, Config *config);

int notify_process_incoming_generic_notification(JsonObj *json, char *type,
                                                 Config *config);

void free_notification(Notification *ptr);

#endif /* INC_NOTIFICATION_H_ */
