/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef MESH_DATA_H
#define MESH_DATA_H

#include "mesh.h"
#include "config.h"

void clear_request(MRequest **data);
void handle_recevied_data(MRequest *data, Config *config);
int process_incoming_websocket_message(Message *message, Config *config);
void process_incoming_websocket_response(Message *message, void *data);

#endif /* MESH_DATA_H */
