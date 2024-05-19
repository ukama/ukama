/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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
#include "client.h"

extern MapTable *NodesTable;

static void free_message(Message *message) {

    if (message == NULL) return;
    
    free(message->reqType);
    free(message->seqNo);
    free(message->data);

    free(message);
}

static int is_websocket_valid(WSManager *manager, MapItem *map) {

    Config *config = NULL;

    if (manager == NULL || map == NULL) return FALSE;

    if (ulfius_websocket_status(manager) == U_WEBSOCKET_STATUS_CLOSE) {
        log_debug("Websocket is closed with node: %s", map->nodeInfo->nodeID);
        return FALSE;
    }

    return TRUE;
}

void websocket_manager(const URequest *request, WSManager *manager,
					   void *data) {

    MapItem *map=NULL;
	WorkList *list;
	WorkItem *work;
    struct timespec ts;
    int ret;
    json_t *jData;

    map = is_existing_item(NodesTable, (char *)data);
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
		if (ulfius_websocket_wait_close(manager, 1) ==
			U_WEBSOCKET_STATUS_OPEN) {
            jData = json_loads(work->data, JSON_DECODE_ANY, NULL);
			if (ulfius_websocket_send_json_message(manager, jData) != U_OK) {
				log_error("Error sending JSON message.");
			}
            json_decref(jData);
		}

		/* 3. Any post-processing. */
		if (work->postFunc) {
			work->postFunc(work->data, work->postArgs);
		}

		/* Free up the memory */
		free_work_item(work);
	}

	return;
}

void websocket_incoming_message(const URequest *request,
								WSManager *manager, WSMessage *message,
								void *data) {
    Message *rcvdMessage=NULL;
    char *responseRemote=NULL;
	int ret;
    MapItem *map=NULL;
    Forward *forward=NULL;
    char *rcvdDataStr=NULL;

    map = is_existing_item(NodesTable, (char *)data);
    if (map == NULL) {
        log_error("Websocket error NodeID: %s not found in table",
                  (char *)data);
        return;
    }

    rcvdDataStr = (char *)calloc(1, message->data_len + 1);
    strncpy(rcvdDataStr, message->data, message->data_len);
    log_debug("Packet recevied. data: %s", rcvdDataStr);

	/* Ignore the rest, for now. */
	if (message->opcode == U_WEBSOCKET_OPCODE_TEXT) {

		ret = deserialize_websocket_message(&rcvdMessage, rcvdDataStr);
		if (ret==FALSE) goto done;

        if (strcmp(rcvdMessage->reqType, UKAMA_NODE_REQUEST) == 0) {
            /* process the incoming and response back on the queue */
            if (process_incoming_websocket_message(rcvdMessage, &responseRemote)) {
                add_work_to_queue(&map->transmit, responseRemote, NULL, 0, NULL, 0);
            }
        }

        else if (strcmp(rcvdMessage->reqType, UKAMA_SERVICE_RESPONSE) == 0) {

            forward = is_existing_item_in_list(map->forwardList,
                                                rcvdMessage->seqNo);

            if (forward == NULL) {
                log_error("No matching uuid in the list. uuid: %s",
                          rcvdMessage->seqNo);
                goto done;
            }

            forward->size     = rcvdMessage->dataSize;
            forward->data     = strdup(rcvdMessage->data);
            forward->httpCode = rcvdMessage->code;

            pthread_cond_broadcast(&forward->hasData);

        } else {
            log_error("Invalid request type on websocket");
            goto done;
        }
	}

done:
    free_message(rcvdMessage);
    free(rcvdDataStr);
	return;
}

void websocket_onclose(const URequest *request,
                       WSManager *manager,
                       void *data) {

    MapItem *map=NULL;
    Config *config=NULL;

    map = is_existing_item(NodesTable, (char *)data);
    if (map == NULL) {
        log_error("Websocket error NodeID: %s not found in table",
                  (char *)data);
        return;
    }

    config = (Config *)map->configData;

	if (map->nodeInfo) {
        if (publish_event(CONN_CLOSE,
                          config->orgName,
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

    remove_map_item_from_table(NodesTable, map->nodeInfo->nodeID);

	return;
}
