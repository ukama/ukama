/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <jansson.h>
#include <ulfius.h>
#include <string.h>
#include <time.h>
#include <errno.h>

#include "usys_log.h"

#include "mesh.h"
#include "work.h"
#include "jserdes.h"
#include "data.h"
#include "map.h"
#include "config.h"

#include "static.h"

extern WorkList *Transmit;
extern MapTable *ClientTable;
extern State    *state;
extern int start_websocket_client(Config *config,
                                  struct _websocket_client_handler *handler);

static 	pthread_mutex_t websocketMutex;
static	pthread_cond_t  websocketFail;

#include <sys/time.h>

static const char *ts_now(void) {
    static char buf[32];
    struct timeval tv;
    struct tm tm;

    gettimeofday(&tv, NULL);
    gmtime_r(&tv.tv_sec, &tm);

    snprintf(buf, sizeof(buf),
             "%04d-%02d-%02d %02d:%02d:%02d.%03ldZ",
             tm.tm_year + 1900,
             tm.tm_mon + 1,
             tm.tm_mday,
             tm.tm_hour,
             tm.tm_min,
             tm.tm_sec,
             tv.tv_usec / 1000);

    return buf;
}

STATIC void clear_response(MResponse **resp) {

	if (*resp==NULL) return;

	free((*resp)->reqType);
	free((*resp)->serviceInfo);
	if ((*resp)->data) {
		free((*resp)->data);
	}

	free(*resp);
}

STATIC bool is_websocket_client_valid(struct _websocket_client_handler *handler,
                                      char *port) {

    if (handler == NULL) return USYS_FALSE;

    if (ulfius_websocket_client_connection_status(handler) == U_WEBSOCKET_STATUS_OPEN) {
        return USYS_TRUE;
    } else {
        handler->websocket = NULL;
        usys_log_debug("Websocket connection is closed with cloud at: %s", port);
        return USYS_FALSE;
    }

    return USYS_FALSE;
}

STATIC bool is_websocket_valid(WSManager *manager, char *port) {

    if (manager == NULL) return USYS_FALSE;

    if (ulfius_websocket_status(manager) == U_WEBSOCKET_STATUS_OPEN) {
        return USYS_TRUE;
    } else {
        manager = NULL;
        usys_log_debug("Websocket connection is closed with cloud at: %s", port);
        return USYS_FALSE;
    }

    return USYS_FALSE;
}

void* monitor_websocket(void *args){

    int ret;
    struct timespec ts;
    Config *config=NULL;
    struct _websocket_client_handler *handler=NULL;
    ThreadArgs *threadArgs;

    threadArgs = (ThreadArgs *)args;

    config  = (Config *)threadArgs->config;
    handler = threadArgs->handler;

    while (TRUE) {

        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += MESH_LOCK_TIMEOUT;

        /* Wait or timed out until the socket closes */
        ret = pthread_cond_timedwait(&websocketFail, &websocketMutex, &ts);
        if (ret == ETIMEDOUT) {
            pthread_mutex_unlock(&websocketMutex);
            if (!is_websocket_client_valid(handler, config->remoteConnect)) {
                usys_log_error("Trying to reconnect ...");
                /* Connect again */
                while (start_websocket_client(config, handler) == FALSE) {
                    usys_log_error("Remote websocket connect failure. Retrying: %d",
                              MESH_LOCK_TIMEOUT);
                    sleep(MESH_LOCK_TIMEOUT);
                }
            } else {
                continue;
            }
        }
    }

    return NULL;
}

#define WDBG(fmt, ...) usys_log_debug("[%s] " fmt, ts_now(), ##__VA_ARGS__)
#define WERR(fmt, ...) usys_log_error("[%s] " fmt, ts_now(), ##__VA_ARGS__)

void websocket_manager(const URequest *request, WSManager *manager,
					   void *data) {

    int ret;
	WorkList *list=NULL;
	WorkItem *work=NULL;
	WorkList **transmit = &Transmit;
    json_t *jData=NULL;
    struct timespec ts;
    Config *config=NULL;

	if (*transmit == NULL)
		return;

	list   = *transmit;
    config = (Config *)data;

    pthread_mutex_init(&websocketMutex, NULL);
    pthread_cond_init(&websocketFail, NULL);

	while (TRUE) {

		pthread_mutex_lock(&(list->mutex));

        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += MESH_LOCK_TIMEOUT;

		if (list->exit) { /* Likely we are closing the socket. */
			break;
		}

        if (list->exit) break;

		if (list->first == NULL) { /* Empty. Wait. */
			ret = pthread_cond_timedwait(&(list->hasWork), &(list->mutex), &ts);
            if (ret == ETIMEDOUT) {
                pthread_mutex_unlock(&(list->mutex));
                if (!is_websocket_valid(manager, config->remoteConnect)) {
                    pthread_cond_broadcast(&websocketFail);
                    return; /* Close the websocket */
                } else {
                    continue;
                }
            }
        }

		/* We have some packet to transmit. */
		work = get_work_to_transmit(list);

		/* Unlock. */
		pthread_mutex_unlock(&(list->mutex));

		if (work == NULL) {
			continue;
		}

		/* 1. Any pre-processing. */
		if (work->preFunc) {
			work->preFunc(work->data, work->preArgs);
		}

        /* 2. Send data over the wire. */
        if (ulfius_websocket_wait_close(manager, 1) == U_WEBSOCKET_STATUS_OPEN) {

            json_error_t jerr;
            jData = json_loads(work->data, 0, &jerr);

            if (jData == NULL) {
                WERR("json_loads failed: line=%d col=%d pos=%d text=%s",
                     jerr.line, jerr.column, jerr.position, jerr.text);
                WERR("payload (first 256): %.256s", work->data ? work->data : "(null)");
                destroy_work_item(work);
                continue;
            }

            /* Optional: dump outbound JSON (bounded) */
            char *dump = json_dumps(jData, JSON_COMPACT);
            if (dump) {
                WDBG("WS TX -> %zu bytes: %.512s%s",
                     strlen(dump), dump, (strlen(dump) > 512) ? "..." : "");
                free(dump);
            }

            int rc = ulfius_websocket_send_json_message(manager, jData);
            if (rc != U_OK) {
                WERR("ulfius_websocket_send_json_message failed rc=%d", rc);
            } else {
                WDBG("WS TX OK");
            }

            json_decref(jData);
            jData = NULL;

        } else {
            WERR("websocket not open; dropping message");
        }

		/* 3. Any post-processing. */
		if (work->postFunc) {
			work->postFunc(work->data, work->postArgs);
		}

		/* Free up the memory */
		destroy_work_item(work);
	}

	return;
}

void websocket_incoming_message(const URequest *request,
                                WSManager *manager,
                                const WSMessage *message,
                                void *config) {

	MResponse *rcvdResp=NULL;
    Message *rcvdMessage=NULL;
    char *data = NULL;
	json_t *json;
	int ret;

    data = (char *)calloc(1, message->data_len+1);
    if (data == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                  message->data_len+1);
        return;
    }
    strncpy(data, message->data, message->data_len);
    
	usys_log_debug("Packet recevied. Data: %s", data);

	json = json_loads(data, JSON_DECODE_ANY, NULL);
	if (json == NULL) {
		usys_log_error("Error loading recevied data into JSON format: %s",
				  data);
		goto done;
	}

	ret = deserialize_websocket_message(&rcvdMessage, json);
	if (ret==FALSE) {
		if (rcvdResp != NULL) free(rcvdResp);
		goto done;
	}

    if (strcmp(rcvdMessage->reqType, MESH_SERVICE_REQUEST) == 0 ) {
        process_incoming_websocket_message(rcvdMessage, (Config *)config);
    } else if (strcmp(rcvdMessage->reqType, MESH_NODE_RESPONSE) == 0) {
        process_incoming_websocket_response(rcvdMessage, config);
    } else {
        usys_log_error("Invalid incoming message on the websocket. Ignored");
    }

done:
    if (data) free(data);
	if (json) json_decref(json);
	clear_response(&rcvdResp);
	return;
}

void  websocket_onclose(const URequest *request, WSManager *manager,
						void *data) {

	return;
}

