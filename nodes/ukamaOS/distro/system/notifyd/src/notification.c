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

char *gRemoteServer;
char* gNodeID;
char* gNodeType;

NotifyHandler handler[MAX_SERVICE_COUNT] = {
                {
                    .service = "noded",
                    .alertHandler = &notify_process_incoming_noded_notification,
                    .eventHandler = &notify_process_incoming_noded_notification,
                },
                {
                    .service = "core",
                    .alertHandler =
                               &notify_process_incoming_generic_notification,
                    .eventHandler =
                               &notify_process_incoming_generic_notification,
                },
                {
                    .service = "stack",
                    .alertHandler =
                               &notify_process_incoming_generic_notification,
                    .eventHandler =
                               &notify_process_incoming_generic_notification,
                }
};


ServiceHandler find_handler(const char* service, char* notif) {
    for (uint8_t idx = 0; idx <= MAX_SERVICE_COUNT ; idx++) {

        if (handler[idx].service && !usys_strcmp(service, handler[idx].service) ) {

            if (!usys_strcmp(notif,NOTIFICATION_ALERT)) {
                return handler[idx].alertHandler;
            }

            if (!usys_strcmp(notif,NOTIFICATION_EVENT)) {
                return handler[idx].eventHandler;
            }

            break;
        }
    }

    usys_log_error("Failed to forward notification from %s to remote "
                        "server as no agent to process the notification found",
                        service);
    return NULL;
}

int notification_init(char* nodeID, char* nodeType, Config* config) {
    gNodeID = usys_strdup(nodeID);
    gNodeType = usys_strdup(nodeType);
    gRemoteServer = usys_strdup(config->remoteServer);
    return STATUS_OK;
}

int notification_exit() {
    usys_free(gRemoteServer);
    usys_free(gNodeID);
    usys_free(gNodeType);
    return STATUS_OK;
}

int notify_send_notification(JsonObj* jNotify) {
    char urlWithEp[MAX_URL_LENGTH] = {0};
    usys_sprintf(urlWithEp, "%s%s", gRemoteServer, DEF_REMOTE_EP);
    return wc_forward_notification(urlWithEp, "POST", jNotify);
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

        if(notif->deviceAttr) {
            usys_free(notif->deviceAttr);
            notif->deviceAttr = NULL;
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

void free_generic_notif_details(ServiceNotifDetails* notif) {

    if (notif) {
        usys_free(notif->serviceName);
        notif->serviceName = NULL;

        usys_free(notif->severity);
        notif->severity = NULL;

        usys_free(notif->description);
        notif->description = NULL;

        usys_free(notif->reason);
        notif->reason = NULL;

        usys_free(notif->details);
        notif->details = NULL;

        if(notif->attr) {
            usys_free(notif->attr->name);
            notif->attr->name = NULL;

            usys_free(notif->attr->value);
            notif->attr->value = NULL;

            usys_free(notif->attr->units);
            notif->attr->units = NULL;

            usys_free(notif->attr);
            notif->attr = NULL;
        }
    }
}

Notification* notify_new_message_from_noded_notif(NodedNotifDetails* noded,
                char* notifType) {

    Notification *envlp = usys_calloc(1, sizeof(Notification));
    if (!envlp) {
        return NULL;
    }

    envlp->serviceName = usys_strdup(noded->serviceName);

    envlp->severity = usys_strdup(noded->severity);

    envlp->notificationType = usys_strdup(notifType);

    envlp->epochTime = noded->epochTime;

    envlp->nodeId = usys_strdup(gNodeID);

    envlp->nodeType = usys_strdup(gNodeType);

    envlp->description = usys_strdup(noded->deviceDesc);

    return envlp;

}

Notification* notify_new_message_from_generic_notification(
                ServiceNotifDetails* notif, char* notifType) {

    Notification *envlp = usys_calloc(1, sizeof(Notification));
    if (!envlp) {
        return NULL;
    }

    envlp->serviceName = usys_strdup(notif->serviceName);

    envlp->severity = usys_strdup(notif->severity);

    envlp->notificationType = usys_strdup(notifType);

    envlp->epochTime = notif->epochTime;

    envlp->nodeId = usys_strdup(gNodeID);

    envlp->nodeType = usys_strdup(gNodeType);

    envlp->description = usys_strdup(notif->description);

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

int notify_process_incoming_noded_notification(JsonObj* json, char* notifType) {
    int ret = STATUS_NOK;
    JsonObj* jDetails = NULL;
    JsonObj* jNotify = NULL;
    NodedNotifDetails details = {0};

    /* Deserialize incoming message from noded */
    if (!json_deserialize_noded_notif(json, &details)) {
        goto cleanup;
    }

    Notification *envlp =
                    notify_new_message_from_noded_notif(&details, notifType);
    if (!envlp) {
        goto cleanup;
    }

    /* Serialize details */
    if(json_serialize_noded_notif_details(&jDetails, &details)){
        goto cleanup;
    }

    /* Serialize Notification */
    if(json_serialize_notification(&jNotify,jDetails, envlp)){
        goto cleanup;
    }

    ret = notify_send_notification(jNotify);
    if (ret != STATUS_OK) {
        usys_log_error("Failed to forward notification %s from %s to remote "
                        "server", envlp->description, envlp->serviceName);
    }

    cleanup:
    free_notification(envlp);
    free_noded_details(&details);
    json_free(&jNotify);
    return ret;
}

int notify_process_incoming_noded_event(JsonObj* json, char* notifType) {
    int ret = STATUS_OK;

    return ret;
}

int notify_process_incoming_generic_notification(JsonObj* json,
                char* notifType) {
    int ret = STATUS_NOK;
    JsonObj* jDetails;
    JsonObj* jNotify;
    ServiceNotifDetails details = {0};

    /* Deserialize incoming message from noded */
    if (!json_deserialize_generic_notification(json, &details)) {
        free_generic_notif_details(&details);
        return ret;
    }

    Notification *envlp =
                    notify_new_message_from_generic_notification(&details,
                                    notifType );
    if (!envlp) {
        return ret;
    }

    /* Serialize details */
    if(json_serialize_generic_details(&jDetails, &details)){
        free_notification(envlp);
        return ret;
    }

    /* Serialize Notification */
    if(json_serialize_notification(&jNotify, jDetails, envlp)){
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
    free_generic_notif_details(&details);
    json_free(&jNotify);

    return ret;
}
