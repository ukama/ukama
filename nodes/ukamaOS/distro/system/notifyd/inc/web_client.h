/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "config.h"
#include "notify_macros.h"
#include "web.h"
#include "json_types.h"

int wc_forward_notification(char* url, char *path, char* method, JsonObj* body );
int wc_read_node_info(Config* config);
int web_client_init(char* nodeID, Config* config);
int get_nodeid_from_noded(Config *config);

#endif /* INC_WEB_CLIENT_H_ */
