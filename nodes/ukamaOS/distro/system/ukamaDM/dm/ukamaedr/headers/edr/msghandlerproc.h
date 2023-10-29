/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef INC_MSGHANDLERPROC_H_
#define INC_MSGHANDLERPROC_H_

#include "headers/edr/ifmsg.h"

char*  msghandler_proc_read_req(MsgFrame* req, size_t* size);
char*  msghandler_proc_write_req(MsgFrame* req, size_t* size);
char*  msghandler_proc_exec_req(MsgFrame* req, size_t* size);
char*  msghandler_proc_alert_resp(MsgFrame* req, size_t* size);
char* msghandler_proc_unkown_req(MsgFrame* req, size_t* size);
char* msghandler_proc(char* req, size_t* size);

#endif /* INC_MSGHANDLERPROC_H_ */
