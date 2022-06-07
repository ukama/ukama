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
#include "web_client.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

//TODO
char *remote_url = "localhost:8091";
char* gNodeID;
char* gNodeType;

//TODO: try runtime update
NotifyHandler handler[MAX_SERVICE_COUNT] = {
                {
                    .service = "noded",
                    .alertHandler = &notify_process_incoming_noded_alert,
                    .eventHandler = &notify_process_incoming_noded_event,
                },
};


ServiceHandler find_handler(const char* service, char* notif) {
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

int notify_init(char* nodeID, char* nodeType) {
    gNodeID = usys_strdup(nodeID);
    gNodeType = usys_strdup(nodeType);
    return STATUS_OK;
}

int notify_send_notification(JsonObj* jNotify) {

    return wc_forward_notification(remote_url, "POST", jNotify);

}

void free_notification(Notification* notif) {

    if (notif) {

        if (notif->nodeId) {
            usys_free(notif->nodeId);
            notif->nodeId = NULL;
        }

        if(notif->serviceName) {
            usys_free(notif->serviceName);
            notif->serviceName = NULL;
        }

        if(notif->notificationType) {
            usys_free(notif->notificationType);
            notif->notificationType = NULL;
        }

        if(notif->severity) {
            usys_free(notif->severity);
            notif->severity = NULL;
        }

        if(notif->description) {
            usys_free(notif->description);
            notif->description = NULL;
        }

        if(notif->nodeType) {
            usys_free(notif->nodeType);
            notif->nodeType = NULL;
        }

        usys_free(notif);
        notif = NULL;
    }
}

void free_noded_details(NodedNotifDetails* notif) {

    if (notif) {

        if(notif->serviceName) {
            usys_free(notif->serviceName);
            notif->serviceName = NULL;
        }

        if(notif->moduleID) {
            usys_free(notif->moduleID);
            notif->moduleID = NULL;
        }

        if(notif->severity) {
            usys_free(notif->severity);
            notif->severity = NULL;
        }

        if(notif->deviceName) {
            usys_free(notif->deviceName);
            notif->deviceName = NULL;
        }

        if(notif->deviceDesc) {
            usys_free(notif->deviceDesc);
            notif->deviceDesc = NULL;
        }

        if(notif->deviceAttr) {
            usys_free(notif->deviceAttr);
            notif->deviceAttr = NULL;
        }

        if(notif->dataType) {
            usys_free(notif->dataType);
            notif->dataType = NULL;
        }

        if(notif->deviceAttrValue) {
            usys_free(notif->deviceAttrValue);
            notif->deviceAttrValue = NULL;
        }
        if(notif->units) {
            usys_free(notif->units);
            notif->units = NULL;
        }

    }
}

Notification* notify_new_message_from_noded_alert(NodedNotifDetails* noded) {

    Notification *envlp = usys_malloc(sizeof(Notification));
    if (!envlp) {
        return NULL;
    }

    envlp->serviceName = usys_strdup(noded->serviceName);

    envlp->severity = usys_strdup(noded->severity);

    envlp->notificationType = usys_strdup(NOTIFICATION_ALERT);

    envlp->epcohTime = noded->epcohTime;

    envlp->nodeId = usys_strdup(gNodeID);

    envlp->nodeType = usys_strdup(gNodeType);

    envlp->description = usys_strdup(noded->deviceDesc);

    return envlp;

}

int notify_process_incoming_notification(const char* service, char* notif,
                JsonObj* json){
    int ret = STATUS_OK;
    ServiceHandler handler = find_handler(service, notif);
    if (handler) {
       ret =  handler(json, notif);
    }

    return ret;
}

int notify_process_incoming_noded_alert(JsonObj* json, char* notifType) {
    int ret = STATUS_NOK;
    JsonObj* jDetails;
    JsonObj* jNotify;
    NodedNotifDetails details;

    /* Deserialize incoming message from noded */
    if (!json_deserialize_noded_alerts(json, &details)) {
        return ret;
    }

    Notification *envlp =
                    notify_new_message_from_noded_alert(&details);
    if (!envlp) {
        return ret;
    }

    /* Serialize details */
    if(json_serialize_noded_alert_details(&jDetails, &details)){
        free_notification(envlp);
        return ret;
    }

    /* Serialize Notification */
    if(json_serialize_notification(&jNotify,jDetails, envlp)){
        free_notification(envlp);
        json_decref(jDetails);
        return ret;
    }

    ret = notify_send_notification(jNotify);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to forward notification %s from %s to remote "
                        "server", envlp->description, envlp->serviceName);
    }

    free_notification(envlp);
    free_noded_details(&details);
    json_decref(jDetails);
    json_decref(jNotify);

    return ret;
}

int notify_process_incoming_noded_event(JsonObj* json, char* notifType) {
    int ret = STATUS_OK;

    return ret;
}
