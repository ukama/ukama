/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "config.h"
#include "config_macros.h"
#include "web.h"
#include "json_types.h"

int wc_read_node_info(Config* config);
bool wc_send_app_restart_request(Config *config, char *app);
int get_nodeid_from_noded(Config *config);

#endif /* INC_WEB_CLIENT_H_ */
