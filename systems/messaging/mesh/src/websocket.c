/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <jansson.h>
#include <ulfius.h>
#include <string.h>
#include <time.h>
#include <errno.h>

#include "mesh.h"
#include "log.h"
#include "work.h"
#include "jserdes.h"
#include "data.h"
#include "map.h"
#include "config.h"

extern WorkList *Transmit;
extern MapTable *IDTable;

/*
 * clear_response -- free up memory from MResponse.
 *
 */
static void clear_response(MResponse **resp) {

	if (*resp==NULL) return;

	free((*resp)->reqType);
	free((*resp)->serviceInfo);
	if ((*resp)->data) {
		free((*resp)->data);
	}

	free(*resp);
}

/*
 * clear_websocket_data -- free up memory from websocket
 *
 */
static void clear_websocket_data(WebsocketData **ptr) {

    if (*ptr==NULL) return;

    free((*ptr)->deviceInfo->nodeID);
    free((*ptr)->deviceInfo);
    free(*ptr);
}
/*
 * websocket_status --
 *
 */
static int websocket_valid(WSManager *manager, void *data) {

    WebsocketData *websocketData = NULL;
    Config *config=NULL;

    websocketData = (WebsocketData *)data;
    config = (Config *)websocketData->data;

    if (manager == NULL || data == NULL) return FALSE;

    if (ulfius_websocket_status(manager) ==
        U_WEBSOCKET_STATUS_CLOSE) { /* connection close */
        log_debug("Websocket connection is closed");

        /* publish event on AMQP */
        if (publish_amqp_event(config->conn, config->amqpExchange, CONN_CLOSE,
						   websocketData->deviceInfo->nodeID) == FALSE) {
            log_error("Error publishing device connect msg on AMQP exchange");
        } else {
            log_debug("Send AMQP offline msg for NodeID: %s",
                      websocketData->deviceInfo->nodeID);
        }

        return FALSE;
    }

    return TRUE;
}

/*
 * websocket related callback functions.
 */
void websocket_manager(const URequest *request, WSManager *manager,
					   void *data) {

	WorkList *list;
	WorkItem *work;
	WorkList **transmit = &Transmit;
    struct timespec ts;
    int ret;

	if (*transmit == NULL)
		return;

	list = *transmit;

	while (TRUE) {

		pthread_mutex_lock(&(list->mutex));

        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += MESH_LOCK_TIMEOUT;

		if (list->exit) { /* Likely we are closing the socket. */
			break;
		}

		if (list->first == NULL) { /* Empty. Wait. */
			log_debug("Waiting for work to be available ...");
			ret = pthread_cond_timedwait(&(list->hasWork), &(list->mutex), &ts);
            if (ret == ETIMEDOUT) {
                /* Check if this connection is still valid, otherwise
                 * report and close
                 */
                pthread_mutex_unlock(&(list->mutex));

                if (!websocket_valid(manager, data)) {
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

		/* We have valid work to do. yaay. */

		/* 1. Any pre-processing. */
		if (work->preFunc) {
			work->preFunc(work->data, work->preArgs);
		}

		/* 2. Send data over the wire. */
		/* Currently, Packet is JSON string. Send it over. */
		if (ulfius_websocket_wait_close(manager, 2000) ==
			U_WEBSOCKET_STATUS_OPEN) {
			if (ulfius_websocket_send_json_message(manager, work->data)
				!= U_OK) {
				log_error("Error sending JSON message.");
			}
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

/*
 * websocket_incoming_message -- handle incoming message over websocket.
 *
 */
void websocket_incoming_message(const URequest *request,
								WSManager *manager, WSMessage *message,
								void *data) {
	MRequest *rcvdData=NULL;
	MResponse *rcvdResp=NULL;
    Config *config=NULL;
    WebsocketData *websocketData=NULL;
	json_t *json;
	int ret;

    websocketData = (WebsocketData *)data;
    config = (Config *)websocketData->data;

	log_debug("Packet recevied. Data: %s", message->data);

	/* If we recevied a packet and our proxy is disable log and reject*/
	if (config->proxy == FALSE) {
		log_error("Recevie packet while reverse-proxy is disabled ignored");
		goto done;
	}

	/* Ignore the rest, for now. */
	if (message->opcode == U_WEBSOCKET_OPCODE_TEXT) {
		/* Convert to JSON and deserialize it. */
		json = json_loads(message->data, JSON_DECODE_ANY, NULL);

		if (json==NULL) {
			log_error("Error loading recevied data to JSON format. Str: %s",
					  message->data);
			goto done;
		}

		/* Convert JSON into request. */
		ret = deserialize_forward_request(&rcvdData, json);
		if (ret==FALSE) goto done;

		handle_recevied_data(rcvdData, config);

		/* Free up the memory from deser. */
		clear_request(&rcvdData);
	}

 done:
	if (json) json_decref(json);
	clear_response(&rcvdResp);
	return;
}

/*
 * websocket_onclose -- is called when the websocket is closed.
 *
 */

void  websocket_onclose(const URequest *request, WSManager *manager,
						void *data) {

    WebsocketData *websocketData=NULL;
    Config *config=NULL;

    websocketData = (WebsocketData *)data;
    config = (Config *)websocketData->data;

	if (config == NULL)
		return;

	if (config->deviceInfo) {
		if (publish_amqp_event(config->conn, config->amqpExchange, CONN_CLOSE,
                               config->deviceInfo->nodeID) == FALSE) {
			log_error("Error publish device close msg on AMQP exchange: %s",
					  config->deviceInfo->nodeID);
		} else {
			log_debug("AMQP device close msg successfull for NodeID: %s",
					  config->deviceInfo->nodeID);
		}
	}

    clear_websocket_data(&websocketData);

	return;
}
