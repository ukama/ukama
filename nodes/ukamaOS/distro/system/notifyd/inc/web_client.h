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

/**
 * @fn      int wc_forward_notification(char*, char*, char*, JsonObj*)
 * @brief   Forward the node notifications to the remote server
 *
 * @param   url
 * @param   method
 * @param   body
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int wc_forward_notification(char* url, char* method,
                JsonObj* body );

/**
 * @fn      int wc_read_node_info(char*, char*, char*, int)
 * @brief   Read node UUID and Type form the node info provided by noded
 *          service.
 *
 * @param   nodeID
 * @param   config
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int wc_read_node_info(Config* config);

/**
 * @fn      int web_client_init(char*, char*)
 * @brief   Connected to Noded for reading Unit info.
 *
 * @param   nodeID
 * @param   config
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int web_client_init(char* nodeID, Config* config);

int get_nodeid_from_noded(Config *config);

#endif /* INC_WEB_CLIENT_H_ */
