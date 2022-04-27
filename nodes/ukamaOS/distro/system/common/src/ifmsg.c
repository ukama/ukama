/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "ifmsg.h"

#include "headers/errorcode.h"

#include <time.h>
#include <stdlib.h>
#include <string.h>

void ifmsg_init() {
    srand(time(NULL));
}

MsgFrame *create_msgframe(MsgType msgtype, uint16_t inst, uint16_t rid,
                          uint16_t objtype, size_t size, bool token,
                          void *data) {
    MsgFrame *msg = malloc(sizeof(MsgFrame));
    if (msg) {
        memset(msg, '\0', sizeof(MsgFrame));
        msg->datasize = size;
        msg->objecttype = objtype;
        msg->instance = inst;
        msg->msgtype = msgtype;
        if (token) {
            msg->reqtoken = rand();
        }
        msg->resourceId = rid;
        msg->response = 0;
        if (size > 0) {
            msg->data = malloc(size);
            if (msg->data) {
                memset(msg->data, '\0', size);
                memcpy(msg->data, data, size);
            }
        } else {
            msg->data = NULL;
        }
    }
    return msg;
}

void free_msgframe(MsgFrame **msg) {
    if (*msg) {
        if ((*msg)->data) {
            free((*msg)->data);
            (*msg)->data = NULL;
        }
        free(*msg);
        *msg = NULL;
    }
}

char *msgframe_serialize(MsgFrame *msg, size_t *sz) {
    char *data = NULL;
    if (msg) {
        *sz = sizeof(MsgFrame) + msg->datasize;
        data = malloc(*sz);
        if (data) {
            memset(data, '\0', *sz);
            memcpy(data, msg, sizeof(MsgFrame));
            memcpy(&data[sizeof(MsgFrame)], msg->data, msg->datasize);
        } else {
            *sz = 0;
        }
    }
    return data;
}

MsgFrame *msgframe_deserialize(char *data) {
    MsgFrame *msg = NULL;
    if (data) {
        msg = malloc(sizeof(MsgFrame));
        if (msg) {
            memset(msg, '\0', sizeof(MsgFrame));
            memcpy(msg, data, sizeof(MsgFrame));
            if (msg->datasize > 0) {
                msg->data = NULL;
                msg->data = malloc(msg->datasize);
                if (msg->data) {
                    memcpy(msg->data, &data[sizeof(MsgFrame)], msg->datasize);
                }
            }
        }
    }
    return msg;
}

int msgframe_validate(MsgFrame *rmsg, MsgFrame *smsg) {
    int ret = 0;

    if (rmsg->reqtoken != smsg->reqtoken) {
        ret = ERR_IFMSG_MISMAT_TOKEN;
    }

    if ((rmsg->objecttype != smsg->objecttype) && (rmsg->misc != smsg->misc)) {
        ret = ERR_IFMSG_MISMAT_MSG_REQ;
    }

    if (rmsg->instance != smsg->instance) {
        ret = ERR_IFMSG_MISMAT_INST;
    }

    if (rmsg->resourceId != smsg->resourceId) {
        ret = ERR_IFMSG_MISMAT_RSRC_ID;
    }

    return ret;
}
