/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef REQUEST_H_
#define REQUEST_H_

#include "ctrl.h"

#define CTRL_REQ_ID_LEN                    64

typedef enum {
    CtrlMsgUnknown = 0,
    CtrlMsgPing,
    CtrlMsgGetStatus,
    CtrlMsgScan,
    CtrlMsgGetInfo,
    CtrlMsgGetAlarmStatus,
    CtrlMsgClearActiveAlarms,
    CtrlMsgAlarmSubscribe,
    CtrlMsgSelfTest,
    CtrlMsgSendConfigurationData,
    CtrlMsgCalibrate,
    CtrlMsgGetTilt,
    CtrlMsgSetTilt,
    CtrlMsgGetDeviceData,
    CtrlMsgResetSoftware
} CtrlMsgType;

typedef struct {
    char id[CTRL_REQ_ID_LEN];
    CtrlMsgType type;
    JsonObj *payload;
} CtrlRequest;

void ctrl_request_init(CtrlRequest *request);
void ctrl_request_free(CtrlRequest *request);
bool ctrl_request_from_json(JsonObj *json, CtrlRequest *request);
const char *ctrl_msg_type_str(CtrlMsgType type);
CtrlMsgType ctrl_msg_type_from_str(const char *type);

#endif /* REQUEST_H_ */
