/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/msghandlerproc.h"
#include "headers/globalheader.h"
#include "inc/ereg.h"

#include <stdint.h>
#include <string.h>

char *msghandler_proc_alert_rep(void* ctx, MsgFrame *rep, size_t *size) {
	int ret = 0;
	char *resp = NULL;
	if (rep->response == MSG_RESP_SUCCESS) {
		fprintf(stdout,
				"MSGHANDLERPROC:: Alert received for X/%d/%d with Token %d and size %d\r\n",
				rep->instance, rep->resourceId, rep->reqtoken, rep->datasize);
		/* TODO: Add Alarm handling here */
		if (rep->data) {
		ret = ereg_handle_alarm(ctx, rep->data, size);
		}
	} else {
		fprintf(stderr,
				"MSGHANDLERPROC:: Alert responded with failure (%d) for Msg Token %d\r\n",
				rep->response, rep->reqtoken);
	}
    rep->msgtype = MSG_TYPE_ALERT_RESP;
    rep->response = ret;
	/* serialize the response */
	resp = msgframe_serialize(rep, size);

	return resp;
}

char *msghandler_proc_unkown_req(MsgFrame *req) {
    char *resp = NULL;
    return resp;
}

char *msghandler_proc(void* ctx, char *req, size_t *size) {
    char *resp = NULL;
    MsgFrame *msg = msgframe_deserialize(req);
    if (msg) {
    	switch (msg->msgtype) {
    	case MSG_TYPE_ALERT_REP:
    		resp = msghandler_proc_alert_rep(ctx, msg, size);
    		break;
    	default:
            resp = msghandler_proc_unkown_req(msg);
            break;
        }
    }
	if (msg) {
		UKAMA_FREE(msg->data);
		UKAMA_FREE(msg);
	}
    return resp;
}
