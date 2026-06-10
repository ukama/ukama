/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_REQUEST_H_
#define AISG_EMU_REQUEST_H_

#include <stdbool.h>

#include "emu.h"

typedef enum {
    EmuMsgUnknown = 0,
    EmuMsgPing,
    EmuMsgGetStatus,
    EmuMsgScan,
    EmuMsgGetInfo,
    EmuMsgGetAlarmStatus,
    EmuMsgClearActiveAlarms,
    EmuMsgAlarmSubscribe,
    EmuMsgSelfTest,
    EmuMsgSendConfigurationData,
    EmuMsgCalibrate,
    EmuMsgGetTilt,
    EmuMsgSetTilt,
    EmuMsgGetDeviceData,
    EmuMsgResetSoftware,
} EmuMsgType;

typedef struct {
    char       id[64];
    EmuMsgType type;
    JsonObj   *payload;
} EmuRequest;

void emu_request_init(EmuRequest *request);
void emu_request_free(EmuRequest *request);

bool emu_request_from_json(JsonObj *json, EmuRequest *request);

const char *emu_msg_type_str(EmuMsgType type);
EmuMsgType emu_msg_type_from_str(const char *type);

#endif /* AISG_EMU_REQUEST_H_ */
