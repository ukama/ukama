/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_MSGHANDLER_H_
#define INC_MSGHANDLER_H_

#include "inc/ukamadr.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "ifmsg.h"
#include "headers/ubsp/devices.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"

#include <pthread.h>

#define UKAMAREGPORT 			7832
#define UKAMALWM2MPORT 			7932
#define RECV_MSG_TIMEOUT		10

void msghandler_init();
void msghandler_exit();
void msghandler_start();
void msghandler_stop();

/* Client for the Asynchronous messages like event and  Alerts.*/
int msghandler_client(char* resp, size_t respsize);

/* Server for the request/response messages from the client.*/
int msghandler_server();
void msghandler_create_client();
void * msghandler_service(void* data);

MsgFrame* msghandler_client_send(MsgFrame *msg, size_t* size, int* sflag);
#endif /* INC_MSGHANDLER_H_ */
