/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <ulfius.h>

#include "mesh.h"
#include "log.h"
#include "work.h"

extern WorkList **get_transmit(void);

/*
 * websocket related callback functions.
 */

void websocket_manager(const URequest *request, WSManager *manager,
		       void *data) {

  WorkList *list;
  WorkItem *work;
  WorkList **transmit = get_transmit();

  if (*transmit == NULL)
    return;

  list = *transmit;

  while (TRUE) {

    pthread_mutex(&(list->mutex));

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
    pthread_unlock_mutex(&(list->mutex));

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

void websocket_incoming_message(const URequest *request,
				WSManager *manager, WSMessage *message,
				void *data) {
  return;
}

void  websocket_onclose(const URequest *request, WSManager *manager,
			void *data) {

  return;
}
