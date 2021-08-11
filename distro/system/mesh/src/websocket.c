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

#include "mesh.h"
#include "log.h"
#include "work.h"

extern WorkList *Transmit;
extern void handle_recevied_data(MRequest *data, void *cfg);
extern void clear_data(MRequest **data);

/*
 * websocket related callback functions.
 */

void websocket_manager(const URequest *request, WSManager *manager,
		       void *data) {

  WorkList *list;
  WorkItem *work;
  WorkList **transmit = &Transmit;

  if (*transmit == NULL)
    return;

  list = *transmit;

  while (TRUE) {

    pthread_mutex_lock(&(list->mutex));

    if (list->exit) { /* Likely we are closing the socket. */
      break;
    }

    if (list->first == NULL) { /* Empty. Wait. */
      log_debug("Waiting for work to be available ...");
      pthread_cond_wait(&(list->hasWork), &(list->mutex)); /* unlock mutex. */
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
    if (ulfius_websocket_wait_close(manager, 2000) == U_WEBSOCKET_STATUS_OPEN) {
      if (ulfius_websocket_send_json_message(manager, work->data) != U_OK) {
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
  MRequest *rcvdData;
  json_t *json;
  char *str;
  int ret;
  Config *config = (Config *)data;

  /* If we recevied a packet and our proxy is disable, log error and reject*/
  if (config->proxy == FALSE) {
    log_error("Recevied packet while reverse-proxy is disabled. Ignoring");
    goto done;
  }

  /* Ignore the rest, for now. */
  if (message->opcode == U_WEBSOCKET_OPCODE_TEXT) {

    log_debug("Packet received. Data: %s", message->data);

    /* Convert to JSON and deserialize it. */
    json = json_loads(message->data, JSON_DECODE_ANY, NULL);

    if (json==NULL) {
      log_error("Error loading recevied data into JSON format. Str: %s",
		message->data);
      goto done;
    }

    /* Convert JSON into request. */
    ret = deserialize_forward_request(&rcvdData, json);
    if (ret==FALSE)
      goto done;

    handle_recevied_data(rcvdData, config);

    /* Free up the memory from deser. */
    clear_request(&rcvdData);
  }

 done:
  json_decref(json);
  return;
}

void  websocket_onclose(const URequest *request, WSManager *manager,
			void *data) {

  return;
}
