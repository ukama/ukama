/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
