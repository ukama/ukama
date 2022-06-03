/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "notification.h"

#include "notify_macros.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

//TODO: try runtime update
NotifyHandler handler[MAX_SERVICE_COUNT] = {
                {
                    .service = "noded",
                    .alertHandler = &notify_process_incoming_noded_alert,
                    .eventHandler = &notify_process_incoming_noded_event,
                },
};

ServiceHandler* find_handler(char* service, char* notif) {
    for (uint8_t idx = 0; idx <= MAX_SERVICE_COUNT ; idx++) {

        if (!usys_strcmp(service, handler[idx].service)) {

            if (!usys_strcmp(notif,NOTIFICATION_ALERT)) {
                return handler[idx].alertHandler;
            }

            if (!usys_strcmp(notif,NOTIFICATION_EVENT)) {
                return handler[idx].eventHandler;
            }

            break;
        }
    }

    return NULL;
}

int notify_process_incoming_notification(char* service, char* notif, JsonObj* json){
    int ret = STATUS_OK;
    ServiceHandler handler = find_handler(service, notif);
    if (handler) {
       ret =  handler(json, notif);
    }

    return ret;
}

int notify_process_incoming_noded_alert(JsonObj* json, char* notif) {
    int ret = STATUS_OK;

    return ret;
}

int notify_process_incoming_noded_event(JsonObj* json, char* notif) {
    int ret = STATUS_OK;

     return ret;
}
