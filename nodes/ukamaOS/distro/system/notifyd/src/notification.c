/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "notification.h"

#include "notify_macros.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

NotifyHandler handler[MAX_SERVICE_COUNT] = {
    { /* First entry is always the default */
        .service = "default",
        .alertHandler = &notify_process_incoming_generic_notification,
        .eventHandler = &notify_process_incoming_generic_notification,
    },
    {
        .service = "noded",
        .alertHandler = &notify_process_incoming_generic_notification,
        .eventHandler = &notify_process_incoming_generic_notification,
    },
    {
        .service = "core",
        .alertHandler = &notify_process_incoming_generic_notification,
        .eventHandler = &notify_process_incoming_generic_notification,
    },
    {
        .service = "stack",
        .alertHandler = &notify_process_incoming_generic_notification,
        .eventHandler = &notify_process_incoming_generic_notification,
    },
    {
        .service = "deviced",
        .alertHandler = &notify_process_incoming_generic_notification,
        .eventHandler = &notify_process_incoming_generic_notification,
    },
};

ServiceHandler find_handler(const char* service, char* type) {

    int idx;
    
    for (idx = 1; idx <= MAX_SERVICE_COUNT ; idx++) {
        if (usys_strcmp(service, handler[idx].service) == 0) {

            if (usys_strcmp(type, NOTIFICATION_ALERT) == 0) {
                return handler[idx].alertHandler;
            }

            if (usys_strcmp(type, NOTIFICATION_EVENT) == 0) {
                return handler[idx].eventHandler;
            }

            break;
        }
    }

    /* return default */
    if (usys_strcmp(type, NOTIFICATION_ALERT) == 0) {
        return handler[0].alertHandler;
    } else if (usys_strcmp(type, NOTIFICATION_EVENT) == 0) {
        return handler[0].eventHandler;
    } else {
        usys_log_error("Failed to forward notification from %s to remote "
                       "server as no agent to process the notification found "
                       "for event type: %s", service, type);
    }

    return NULL;
}

int notify_send_notification(JsonObj* jNotify, Config *config) {
    
    char urlWithEp[MAX_URL_LENGTH] = {0};
    
    usys_sprintf(urlWithEp, "%s%s", config->remoteServer, DEF_REMOTE_EP);
    return wc_forward_notification(urlWithEp, "POST", jNotify);
}

void free_notification(Notification *ptr) {

    if (ptr == NULL) return;
    
    usys_free(ptr->serviceName);
    usys_free(ptr->severity);
    usys_free(ptr->module);
    usys_free(ptr->device);
    usys_free(ptr->propertyName);
    usys_free(ptr->propertyValue);
    usys_free(ptr->propertyUnit);
    usys_free(ptr->details);

    usys_free(ptr);
}

int getCode(const Entry* entries, int numEntries, char *type,
            Notification *notif) {

    int code=-1, i;

    for (i = 0; i < numEntries; i++) {
        if (strcmp(notif->serviceName, entries[i].serviceName) == 0 &&
            strcmp(notif->module, entries[i].moduleName) == 0 &&
            strcmp(notif->propertyName, entries[i].propertyName) == 0 &&
            strcmp(type, entries[i].type) == 0) {
            code = entries[i].code;
            break;
        }
    }

    return code;
}

int notify_process_incoming_generic_notification(JsonObj *json, char *type,
                                                 Config *config) {

    int statusCode=-1;
    JsonObj *jNotify=NULL;
    Notification *notification=NULL;

    /* Deserialize incoming message from noded */
    if (!json_deserialize_notification(json, &notification)) {
        return STATUS_NOK;
    }

    statusCode =
        getCode(config->entries, config->numEntries, type, notification);
    if (statusCode == -1) {
        usys_log_error("Unable to process incoming event/alarm. Ignoring"
                       "name: %s module: %s property: %s type: %s",
                       notification->serviceName,
                       notification->module,
                       notification->propertyName,
                       type);
        return STATUS_NOK;
    }

    if (json_serialize_notification(&jNotify, notification, type,
                                    config->nodeID, statusCode) == USYS_FALSE) {
        usys_log_error("Unable to serialize the JSON object");
        return STATUS_NOK;
    }
    json_log(jNotify);

    if (notify_send_notification(jNotify, config) != STATUS_OK) {
        usys_log_error("Failed to forward notification. Service: %s",
                       notification->serviceName);
        return STATUS_NOK;
    }

    json_free(&jNotify);

    return STATUS_OK;
}

int notify_process_incoming_notification(const char *service, char *type,
                                         JsonObj *json, Config *config){

    int ret;
    ServiceHandler handler = NULL;

    handler = find_handler(service, type);
    if (handler) {
        ret = handler(json, type, config);
    }

    return ret;
}
