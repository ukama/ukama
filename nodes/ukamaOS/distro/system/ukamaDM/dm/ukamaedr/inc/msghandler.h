/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef INC_MSGHANDLER_H_
#define INC_MSGHANDLER_H_

#include "headers/ubsp/devices.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "headers/edr/ifmsg.h"

#include <pthread.h>

#define UKAMAREGPORT            7832
#define UKAMALWM2MPORT          7932
#define RECV_MSG_TIMEOUT        10

void msghandler_init();
pthread_t msghandler_start(void* args);

/* Client for the Asynchronous messages like event and  Alerts.*/
int msghandler_client(char* resp, size_t respsize);

/* Server for the request/response messages from the client.*/
int msghandler_server(void* args);
void msghandler_create_client();
void * msghandler_service(void* args);

MsgFrame* msghandler_client_send(MsgFrame *msg, size_t* size, int* sflag);
#endif /* INC_MSGHANDLER_H_ */
