/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef WEB_SOCKET_H_
#define WEB_SOCKET_H_

#include <ulfius.h>

#include "rlogd.h"

int web_socket_cb_ping(const URequest *request, UResponse *response, void *data);
int web_socket_cb_post_log(const URequest *request, UResponse *response, void *data);
int web_socket_cb_default(const URequest *request, UResponse *response, void *data);

void websocket_manager(const URequest *request, WSManager *manager, void *data);
void websocket_incoming_message(const URequest *request, WSManager *manager,
                                const WSMessage *message, void *nodeID);
void websocket_onclose(const URequest *request, WSManager *manager, void *data);

#endif /* WEB_SOCKET_H_ */
