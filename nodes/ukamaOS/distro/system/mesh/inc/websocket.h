/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef WEBSOCKET_H
#define WEBSCOKET_H

#include <ulfius.h>

void* monitor_websocket(void *args);
void websocket_manager_cb(const struct _u_request *request,
						  struct _websocket_manager *manager,
						  void *user_data);
void websocket_incoming_message_cb(const struct _u_request *request,
								   struct _websocket_manager *manager,
								   const struct _websocket_message *message,
								   void *user_data);
void  websocket_onclose_cb(const struct _u_request *request,
						   struct _websocket_manager *manager,
						   void *user_data);

#endif /* WEBSCOKET_H */
