/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef WEB_CLIENT_H_
#define WEB_CLIENT_H_

#include "notify_macros.h"
#include "web.h"
#include "json_types.h"

/**
 * @fn      int wc_forward_notification(char*, char*, char*, JsonObj*)
 * @brief   Forward the node notifications to the remote server
 *
 * @param   url
 * @param   ep
 * @param   method
 * @param   body
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int wc_forward_notification(char* url, char* ep, char* method,
                JsonObj* body );

/**
 * @fn      int wc_read_node_info(char*, char*, char*, int)
 * @brief   Read node UUID and Type form the node info provided by noded
 *          service.
 *
 * @param   nodeID
 * @param   nodeType
 * @param   host
 * @param   port
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int wc_read_node_info(char* nodeID, char* nodeType, char* host, int port);

/**
 * @fn      int web_client_init()
 * @brief   Connected to Noded for reading Unit info.
 *
 * @return  On success, STATUS_OK
 *          On Failure, STATUS_NOK
 */
int web_client_init();

#endif /* INC_WEB_CLIENT_H_ */
