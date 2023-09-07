/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/ubsp/devices.h"
#include "headers/globalheader.h"
#include "headers/edr/ifmsg.h"
#include "inc/reghelper.h"
#include "inc/registry.h"
#include "headers/edr/msghandlerproc.h"
#include "headers/utils/log.h"

#include <stdint.h>
#include <string.h>

char *msghandler_proc_read_req(MsgFrame *req, size_t *size) {
    int ret = 0;
    char *resp = NULL;
    /* Send the read request */
    /*if ResId = 0 read or write all otherwise follow the resource number. */
    ret = reg_read_dev(req);
    req->msgtype = MSG_TYPE_READ_RESP;
    req->response = ret;

    /* serialize the response */
    resp = msgframe_serialize(req, size);
    return resp;
}

char *msghandler_proc_read_inst_count_req(MsgFrame *req, size_t *size) {
    int ret = 0;
    char *resp = NULL;
    /* Send the read request */
    /*if ResId = 0 read or write all otherwise follow the resource number. */
    ret = reg_read_inst_count(req);
    req->msgtype = MSG_TYPE_READ_INST_COUNT_RESP;
    req->response = ret;

    /* serialize the response */
    resp = msgframe_serialize(req, size);
    return resp;
}

char *msghandler_proc_write_req(MsgFrame *req, size_t *size) {
    int ret = 0;
    char *resp = NULL;
    /* Send the read request */
    /*if ResId = 0 read or write all otherwise follow the resource number. */
    ret = reg_write_dev(req);
    req->msgtype = MSG_TYPE_WRITE_RESP;
    req->response = ret;

    /* serialize the response */
    resp = msgframe_serialize(req, size);

    return resp;
}

char *msghandler_proc_exec_req(MsgFrame *req, size_t *size) {
    int ret = 0;
    char *resp = NULL;
    /* Send the read request */
    /*if ResId = 0 read or write all otherwise follow the resource number. */
    ret = reg_exec_dev(req);
    req->msgtype = MSG_TYPE_EXEC_RESP;
    req->response = ret;

    /* serialize the response */
    resp = msgframe_serialize(req, size);
    return resp;
}

char *msghandler_proc_alert_resp(MsgFrame *req, size_t *size) {
    int ret = 0;
    char *resp = NULL;
    /* TODO: Need to put some re-sending of alert in case of failure. */
    if (req->response == MSG_RESP_SUCCESS) {
        log_debug(
            "MSGHANDLERPROC:: Alert responded as received for Msg Token %d",
            req->reqtoken);
    } else {
        log_error(
            "MSGHANDLERPROC:: Alert responded with failure (%d) for Msg Token %d",
            req->response, req->reqtoken);
    }

    return resp;
}

char *msghandler_proc_unkown_req(MsgFrame *req, size_t *size) {
    req->response = MSG_RESP_INAVLID_REQMSG;
    char *resp = msgframe_serialize(req, size);
    if (resp) {
        UKAMA_FREE(req->data);
        UKAMA_FREE(req);
    }
    return resp;
}

char *msghandler_proc(char *req, size_t *size) {
    char *resp = NULL;
    MsgFrame *msg = msgframe_deserialize(req);
    if (msg) {
        switch (msg->msgtype) {
        case MSG_TYPE_READ_REQ:
            resp = msghandler_proc_read_req(msg, size);
            break;
        case MSG_TYPE_WRITE_REQ:
            resp = msghandler_proc_write_req(msg, size);
            break;
        case MSG_TYPE_EXEC_REQ:
            resp = msghandler_proc_exec_req(msg, size);
            break;
        case MSG_TYPE_ALERT_RESP:
            resp = msghandler_proc_alert_resp(msg, size);
            break;
        case MSG_TYPE_READ_INST_COUNT:
            resp = msghandler_proc_read_inst_count_req(msg, size);
            break;
        default:
            resp = msghandler_proc_unkown_req(msg, size);
            break;
        }
    }

    if (resp) {
        UKAMA_FREE(msg->data);
        UKAMA_FREE(msg);
    }
    return resp;
}
