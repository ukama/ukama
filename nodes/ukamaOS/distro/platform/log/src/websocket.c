/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <ulfius.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <pthread.h>

#include "usys_file.h"
#include "usys_types.h"
#include "usys_services.h"

#define PREFIX_WEBSOCKET "v1/logit/"
#define ENV_BINDING_IP   "ENV_BINDING_IP"
#define DEF_BINDING_IP   "127.0.0.1"
#define MAX_LOG_LEN      512

#define WEBSOCKET_MONITOR_TIMEOUT 5

typedef struct _u_instance UInst;
typedef struct _u_request  URequest;
typedef struct _u_response UResponse;
typedef struct _websocket_manager WSManager;
typedef struct _websocket_message WSMessage;
typedef struct {
    char *serviceName;
    int  port;
} ThreadArgs;

static struct _websocket_client_handler handler = {NULL, NULL};
static pthread_mutex_t mutex          = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t websocketMutex = PTHREAD_MUTEX_INITIALIZER;
static pthread_mutex_t initMutex      = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t  hasData        = PTHREAD_COND_INITIALIZER;
static pthread_cond_t  websocketFail  = PTHREAD_COND_INITIALIZER;
static pthread_t monitorThread = NULL;
static char dataToSend[MAX_LOG_LEN] = {0};
static ThreadArgs gThreadArgs;
static int logRemoteInitialized = 0;

int start_websocket_client(char *serviceName, int rlogdPort);

/* log.c */
extern void log_enable_rlogd(int flag);

static bool is_websocket_valid(struct _websocket_client_handler *handler) {

    if (handler == NULL) return USYS_FALSE;

    if (ulfius_websocket_client_connection_status(handler) == U_WEBSOCKET_STATUS_OPEN) {
        return USYS_TRUE;
    } else {
        handler->websocket = NULL;
        return USYS_FALSE;
    }

    return USYS_FALSE;
}

static void* monitor_websocket(void *args) {

    struct timespec ts;

    while (USYS_TRUE) {

        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += WEBSOCKET_MONITOR_TIMEOUT;

        /* timeout or socket failure */
        if (pthread_cond_timedwait(&websocketFail, &websocketMutex, &ts) == ETIMEDOUT) {
            pthread_mutex_unlock(&websocketMutex);
            if (!is_websocket_valid(&handler)) {
                while (start_websocket_client(gThreadArgs.serviceName,
                                              gThreadArgs.port) == USYS_FALSE) {
                    sleep(WEBSOCKET_MONITOR_TIMEOUT);
                }
            } else {
                continue;
            }
        }
    }

    return NULL;
}

static void websocket_manager(const URequest *request,
                              WSManager *manager,
                              void *data) {
    int ret;

    while (USYS_TRUE) {

        pthread_mutex_lock(&mutex);

        if (pthread_cond_wait(&hasData, &mutex) != 0) {
            pthread_mutex_unlock(&mutex);
            return;
        }

        if (ulfius_websocket_status(manager) != U_WEBSOCKET_STATUS_OPEN) {
            handler.websocket = NULL;
            pthread_cond_broadcast(&websocketFail);
            pthread_mutex_unlock(&mutex);
            return;
        }

        /* send the message */
        if (ulfius_websocket_send_message(manager,
                                          U_WEBSOCKET_OPCODE_TEXT,
                                          sizeof(dataToSend),
                                          dataToSend) != U_OK) {
            printf("Unable to send message: %s", dataToSend);
        }

        memset(&dataToSend[0], 0, MAX_LOG_LEN);
        pthread_mutex_unlock(&mutex);
    }

    return;
}

static void websocket_incoming_message(const URequest *request,
                                       WSManager *manager,
                                       const WSMessage *message,
                                       void *data) {
	return;
}

static void  websocket_onclose(const URequest *request,
                               WSManager *manager,
                               void *data) {

	return;
}

int start_websocket_client(char *serviceName, int rlogdPort) {

    int ret = USYS_FALSE;
    char *hostname = NULL;
    char url[128]  = {0};

    struct _u_request request;
    struct _u_response response;

    if (ulfius_init_request(&request) != U_OK) goto done;

    if (ulfius_init_response(&response) != U_OK) goto done;

    if (getenv(ENV_BINDING_IP) == NULL)
        hostname = DEF_BINDING_IP;
    else 
        hostname = getenv(ENV_BINDING_IP);

    sprintf(url, "ws://%s:%d/%s", hostname, rlogdPort, PREFIX_WEBSOCKET);
    if (ulfius_set_websocket_request(&request, url, "protocol",
                                     "permessage-deflate") == U_OK) {
        /* Setup request parameters */
        u_map_put(request.map_header, "User-Agent", serviceName);
        ulfius_add_websocket_client_deflate_extension(&handler);
        request.check_server_certificate = USYS_FALSE;

        ret = ulfius_open_websocket_client_connection(&request,
                                                      &websocket_manager, NULL,
                                                      &websocket_incoming_message, NULL,
                                                      &websocket_onclose, NULL,
                                                      &handler, &response);
        if ( ret == U_OK) {
            ret = USYS_TRUE;
        } else {
            handler.websocket = NULL;
            ret = USYS_FALSE;
        }
    } else {
        ret = USYS_FALSE;
    }

done:
    ulfius_clean_request(&request);
    ulfius_clean_response(&response);

    return ret;
}

void log_remote_init(char *serviceName) {

    int rlogdPort = 0;

    pthread_mutex_lock(&initMutex);

    if (logRemoteInitialized) {
        pthread_mutex_unlock(&initMutex);
        return;
    }

    rlogdPort = usys_find_service_port(SERVICE_RLOG);

    if (handler.websocket) return;
    if (rlogdPort == 0)    return;
    if (strcmp(serviceName, SERVICE_RLOG) == 0)    return;
    if (strcmp(serviceName, SERVICE_STARTER) == 0) return;

    pthread_mutex_init(&mutex, NULL);
    pthread_cond_init(&hasData, NULL);

    if (start_websocket_client(serviceName, rlogdPort) == USYS_FALSE) {
        handler.websocket = NULL;
    }

    gThreadArgs.serviceName = strdup(serviceName);
    gThreadArgs.port        = rlogdPort;
    if (pthread_create(&monitorThread,
                       NULL,
                       monitor_websocket,
                       NULL) != 0) {
        return;
    }

    pthread_detach(monitorThread);
    log_enable_rlogd(USYS_TRUE);
    logRemoteInitialized = USYS_TRUE;
    pthread_mutex_unlock(&initMutex);
}

int log_rlogd(char *message) {

    if (strlen(message) > MAX_LOG_LEN -1) return USYS_FALSE;

    if (handler.websocket == NULL) return USYS_FALSE;

    pthread_mutex_lock(&mutex);
    strncpy(&dataToSend[0], message, strlen(message));
    pthread_cond_broadcast(&hasData);
    pthread_mutex_unlock(&mutex);

    return USYS_TRUE;
}

int is_connect_with_rlogd() {

    if (handler.websocket == NULL) return USYS_FALSE;

    return USYS_TRUE;
}
