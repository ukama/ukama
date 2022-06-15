/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WEBSOCKET_H
#define WEBSCOKET_H

#include <ulfius.h>

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
