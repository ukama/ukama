/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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


int process_incoming_notification(const char* service, char* notif,
                                  JsonObj* json, Config *config);
void free_notification(Notification *ptr);

#endif /* INC_NOTIFICATION_H_ */
