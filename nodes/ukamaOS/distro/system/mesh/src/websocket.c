/**
 * Copyright (c) 2021-present, Ukama Inc.
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

extern WorkList *Transmit;
extern MapTable *IDTable;
extern State    *state;

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
static int is_websocket_valid(WSManager *manager, char *port) {

    if (manager == NULL) return FALSE;

    if (ulfius_websocket_status(manager) == U_WEBSOCKET_STATUS_CLOSE) {
        log_debug("Websocket connection is closed with cloud at: %s", port);
        return FALSE;
    }

    return TRUE;
}

/*
 * websocket related callback functions.
 */
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

	while (TRUE) {

		pthread_mutex_lock(&(list->mutex));

        clock_gettime(CLOCK_REALTIME, &ts);
        ts.tv_sec += MESH_LOCK_TIMEOUT;

		if (list->exit) { /* Likely we are closing the socket. */
			break;
		}

        if (list->exit) break;

		if (list->first == NULL) { /* Empty. Wait. */
			log_debug("Waiting for work to be available ...");
			ret = pthread_cond_timedwait(&(list->hasWork), &(list->mutex), &ts);
            if (ret == ETIMEDOUT) {
                pthread_mutex_unlock(&(list->mutex));
                if (!is_websocket_valid(manager, config->remoteConnect)) {
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
		/* Currently, Packet is JSON string. Send it over. */
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
	MResponse *rcvdResp=NULL;
    Message *rcvdMessage=NULL;
	json_t *json;
	int ret;

	log_debug("Packet recevied. Data: %s", message->data);

	json = json_loads(message->data, JSON_DECODE_ANY, NULL);
	if (json==NULL) {
		log_error("Error loading recevied data into JSON format. Str: %s",
				  message->data);
		goto done;
	}

	ret = deserialize_websocket_message(&rcvdMessage, json);
	if (ret==FALSE) {
		if (rcvdResp != NULL) free(rcvdResp);
		goto done;
	}

    process_incoming_websocket_message(rcvdMessage, (Config *)data);

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

	Config *config = (Config *)data;

	return;
}

