/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <jansson.h>
#include <ulfius.h>
#include <string.h>

#include "websocket.h"
#include "usys_log.h"

/* logger.c */
extern void process_logs(void *nodeID, const char *log);

void websocket_manager(const URequest *request, WSManager *manager,
					   void *data) {

    do {
        sleep(DEF_FLUSH_TIME);
        if (ulfius_websocket_status(manager) == U_WEBSOCKET_STATUS_CLOSE) {
            return;
        }
    } while (1);

	return;
}

void websocket_incoming_message(const URequest *request,
								WSManager *manager,
                                WSMessage *message,
								void *nodeID) {

    usys_log_debug("Recevied message: %s", message->data);
    process_logs(nodeID, message->data);

	return;
}

void websocket_onclose(const URequest *request,
                       WSManager *manager,
                       void *data) {
	return;
}
