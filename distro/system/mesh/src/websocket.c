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

/*
 * websocket related callback functions.
 */

void websocket_manager(const URequest *request, WSManager *manager,
		       void *data) {

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
