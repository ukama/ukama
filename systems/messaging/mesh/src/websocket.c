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
#include "u_amqp.h"

extern MapTable *IDsTable;

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
 * is_websocket_valid --
 *
 */
static int is_websocket_valid(WSManager *manager, MapItem *map) {

    if (manager == NULL || map == NULL) return FALSE;

    if (ulfius_websocket_status(manager) == U_WEBSOCKET_STATUS_CLOSE) {
        log_debug("Websocket is closed with node: %s", map->nodeInfo->nodeID);

        /* publish event on AMQP */
        if (publish_event(CONN_CLOSE,
                          map->nodeInfo->nodeID,
                          map->nodeInfo->nodeIP,
                          map->nodeInfo->nodePort,
                          map->nodeInfo->meshIP,
                          map->nodeInfo->meshPort) == FALSE) {
            log_error("Error publishing device connect msg on AMQP exchange");
        } else {
            log_debug("Send AMQP offline msg for NodeID: %s",
                      map->nodeInfo->nodeID);
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

    MapItem *map=NULL;
	WorkList *list;
	WorkItem *work;
    struct timespec ts;
    int ret;
    json_t *jData;

    map = is_existing_item(IDsTable, (char *)data);
    if (map == NULL) {
        log_error("Websocket error NodeID: %s not found in table",
                  (char *)data);
        return;
    }

    /* Setup transmit and receiving queues for the websocket */
    map->transmit = (WorkList *)calloc(1, sizeof(WorkList));
    map->receive  = (WorkList *)calloc(1, sizeof(WorkList));

    if (map->transmit == NULL || map->receive == NULL) {
        log_error("Memory allocation failure: %d", sizeof(WorkList));
        return;
    }

    /* Initializa the transmit and receive list for the websocket. */
    init_work_list(&map->transmit);
    init_work_list(&map->receive);
    list = map->transmit;

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

                if (!is_websocket_valid(manager, map)) {
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
		/* Currently, packet is JSON string. Send it over. */
		if (ulfius_websocket_wait_close(manager, 2000) ==
			U_WEBSOCKET_STATUS_OPEN) {
            jData = json_loads(work->data, JSON_DECODE_ANY, NULL);
			if (ulfius_websocket_send_json_message(manager, jData) != U_OK) {
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
    Message *rcvdMessage=NULL;
    char *responseRemote=NULL;
	MRequest *rcvdData=NULL;
	int ret;
    MapItem *map=NULL;

    map = is_existing_item(IDsTable, (char *)data);
    if (map == NULL) {
        log_error("Websocket error NodeID: %s not found in table",
                  (char *)data);
        return;
    }

	log_debug("Packet recevied. Data: %s", message->data);

	/* Ignore the rest, for now. */
	if (message->opcode == U_WEBSOCKET_OPCODE_TEXT) {

		ret = deserialize_websocket_message(&rcvdMessage, message->data);
		if (ret==FALSE) goto done;

        /* process the incoming and response back on the queue */
        if (process_incoming_websocket_message(rcvdMessage, &responseRemote)) {
            add_work_to_queue(&map->transmit, responseRemote, NULL, 0, NULL, 0);
        }
	}

 done:
	return;
}

/*
 * websocket_onclose -- is called when the websocket is closed.
 *
 */
void websocket_onclose(const URequest *request, WSManager *manager,
                       void *data) {

    MapItem *map=NULL;

    map = is_existing_item(IDsTable, (char *)data);
    if (map == NULL) {
        log_error("Websocket error NodeID: %s not found in table",
                  (char *)data);
        return;
    }

	if (map->nodeInfo) {
        if (publish_event(CONN_CLOSE,
                          map->nodeInfo->nodeID,
                          map->nodeInfo->nodeIP,
                          map->nodeInfo->nodePort,
                          map->nodeInfo->meshIP,
                          map->nodeInfo->meshPort) == FALSE) {
			log_error("Error publish device close msg on AMQP exchange: %s",
					  map->nodeInfo->nodeID);
		} else {
			log_debug("AMQP device close msg successfull for NodeID: %s",
					  map->nodeInfo->nodeID);
		}
	}

    remove_map_item_from_table(IDsTable, map->nodeInfo->nodeID);

	return;
}
