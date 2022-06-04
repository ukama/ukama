/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_NOTIFICATION_H_
#define INC_NOTIFICATION_H_

#include "notify.h"
#include "jserdes.h"
#include "json_types.h"

typedef int (*ServiceHandler)(JsonObj* json, char* type);

typedef struct {
    char* service;
    ServiceHandler alertHandler;
    ServiceHandler eventHandler;
} NotifyHandler;

/**
 * @fn      int notify_init(char*, char*)
 * @brief   Set some parameters required for notification like node ID and
 *          type.
 *
 * @param   nodeID
 * @param   nodeType
 * @return  On Success, USYS_OK (0)
 *          On Failure, USYS_NOK (-1)
 */
int notify_init(char* nodeID, char* nodeType);
/**
 * @fn      int notify_process_incoming_notification(char*, char*, JsonObj*)
 * @brief   Receive all incoming notification from services and choose the
 *          appropriate handler to parse and process request.
 *
 * @param   service
 * @param   notif
 * @param   json
 * @return  On Success, USYS_OK (0)
 *          On Failure, USYS_NOK (-1)
 */
int notify_process_incoming_notification(char* service, char* notif, JsonObj* json);

/**
 * @fn      int notify_process_incoming_noded_event(JsonObj*, char*)
 * @brief   Parses and process all noded service events.
 *
 * @param   json
 * @param   notif
 * @return  On Success, USYS_OK (0)
 *          On Failure, USYS_NOK (-1)
 */
int notify_process_incoming_noded_event(JsonObj* json, char* notif);

/**
 * @fn      int notify_process_incoming_noded_alert(JsonObj*, char*)
 * @brief   Parses and process all noded service events.
 *
 * @param   json
 * @param   notif
 * @return  On Success, USYS_OK (0)
 *          On Failure, USYS_NOK (-1)
 */
int notify_process_incoming_noded_alert(JsonObj* json, char* notif);

#endif /* INC_NOTIFICATION_H_ */
