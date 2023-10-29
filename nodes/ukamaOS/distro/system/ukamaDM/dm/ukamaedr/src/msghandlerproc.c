/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
