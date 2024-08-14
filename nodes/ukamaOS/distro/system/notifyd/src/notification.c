/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <pthread.h>
#include <jansson.h>
#include <string.h>
#include <stdlib.h>

#include "notification.h"
#include "notify_macros.h"
#include "web_client.h"

#include "usys_api.h"
#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

extern ThreadData *gData;

/* Mutexs to ensure thread-safe writes for various dest */
static pthread_mutex_t logFileMutex = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t stdoutMutex  = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t stderrMutex  = PTHREAD_MUTEX_INITIALIZER;

static int write_to_log_file(JsonObj *json);
static int write_to_stdout(JsonObj *json);
static int write_to_stderr(JsonObj *json);
static int notify_process_incoming_generic_notification(JsonObj *json, char *type,
                                                        Config *config);

static NotifyHandler handler[MAX_SERVICE_COUNT] = {
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

static ServiceHandler find_handler(const char* service, char* type) {

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

static int write_to_log_file(JsonObj *jBuffer) {

    FILE *fPtr = NULL;
    char *str  = NULL;

    if (jBuffer == NULL) return USYS_FALSE;

    str=json_dumps(jBuffer, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY));
    if (str == NULL) return USYS_FALSE;

    pthread_mutex_lock(&logFileMutex);

    fPtr = fopen(DEF_LOG_FILE, "a+");
    if (fPtr == NULL) {
        usys_log_error("Unable to open file: %s error: %s",
                       DEF_LOG_FILE,
                       strerror(errno));
        usys_free(str);
        return USYS_FALSE;
    } else {
        fputs(str, fPtr);
        fclose(fPtr);
    }

    pthread_mutex_unlock(&logFileMutex);
    usys_free(str);

    return USYS_TRUE;
}

static int write_to_stdout(JsonObj *jBuffer) {

    char *str = NULL;

    if (jBuffer == NULL) return USYS_FALSE;

    str=json_dumps(jBuffer, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY));
    if (str == NULL) return USYS_FALSE;

    pthread_mutex_lock(&stdoutMutex);
    fprintf(stdout, "%s\n", str);
    pthread_mutex_unlock(&stdoutMutex);

    usys_free(str);

    return USYS_TRUE;
}

static int write_to_stderr(JsonObj *jBuffer) {

    char *str = NULL;

    if (jBuffer == NULL)                      return USYS_FALSE;

    str=json_dumps(jBuffer, (JSON_INDENT(4)|JSON_COMPACT|JSON_ENCODE_ANY));
    if (str == NULL) return USYS_FALSE;

    pthread_mutex_lock(&stderrMutex);
    fprintf(stderr, "%s\n", str);
    pthread_mutex_unlock(&stderrMutex);

    usys_free(str);

    return USYS_TRUE;
}

static int send_notification(JsonObj* jNotify, Config *config) {

    if (gData->output == STDOUT) {
        return write_to_stdout(jNotify);
    } else if (gData->output == STDERR) {
        return write_to_stderr(jNotify);
    } else if (gData->output == LOG_FILE) {
        return write_to_log_file(jNotify);
    } else if (gData->output == UKAMA_SERVICE) {
        return wc_forward_notification(config->remoteServer,
                                       DEF_REMOTE_EP,
                                       "POST",
                                       jNotify);
    }

    return USYS_FALSE;
}

static int getCode(const Entry* entries, int numEntries, char *type,
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

static int notify_process_incoming_generic_notification(JsonObj *json, char *type,
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
        usys_log_error("No matching code found for received event/alarm");
        usys_log_error("Unable to process incoming event/alarm. Ignoring.");
        usys_log_error("type: %s service: %s name: %s value: %s",
                       type,
                       notification->serviceName,
                       notification->propertyName,
                       notification->propertyValue);
        free_notification(notification);
        return STATUS_NOK;
    }

    if (json_serialize_notification(&jNotify, notification, type,
                                    config->nodeID, statusCode) == USYS_FALSE) {
        usys_log_error("Unable to serialize the JSON object");
        free_notification(notification);
        json_free(&jNotify);
        return STATUS_NOK;
    }
    json_log(jNotify);

    if (send_notification(jNotify, config) != STATUS_OK) {
        usys_log_error("Failed to forward notification. Service: %s",
                       notification->serviceName);
        free_notification(notification);
        json_free(&jNotify);
        return STATUS_NOK;
    }

    /* increment counter */
    pthread_mutex_lock(&gData->mutex);
    gData->count++;
    pthread_mutex_unlock(&gData->mutex);

    json_free(&jNotify);
    free_notification(notification);

    return STATUS_OK;
}

int process_incoming_notification(const char *service, char *type,
                                  JsonObj *json, Config *config){

    int ret;
    ServiceHandler handler = NULL;

    handler = find_handler(service, type);
    if (handler) {
        ret = handler(json, type, config);
    }

    return ret;
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
