/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "request.h"

typedef struct {
    EmuMsgType type;
    const char *name;
} EmuMsgName;

static const EmuMsgName g_msg_names[] = {
    { EmuMsgPing,                  "ping" },
    { EmuMsgGetStatus,             "get_status" },
    { EmuMsgScan,                  "scan" },
    { EmuMsgGetInfo,               "get_info" },
    { EmuMsgGetAlarmStatus,        "get_alarm_status" },
    { EmuMsgClearActiveAlarms,     "clear_active_alarms" },
    { EmuMsgAlarmSubscribe,        "alarm_subscribe" },
    { EmuMsgSelfTest,              "self_test" },
    { EmuMsgSendConfigurationData, "send_configuration_data" },
    { EmuMsgCalibrate,             "calibrate" },
    { EmuMsgGetTilt,               "get_tilt" },
    { EmuMsgSetTilt,               "set_tilt" },
    { EmuMsgGetDeviceData,         "get_device_data" },
    { EmuMsgResetSoftware,         "reset_software" },
    { EmuMsgUnknown,               NULL }
};

void emu_request_init(EmuRequest *request)
{
    if (request == NULL) {
        return;
    }

    memset(request, 0, sizeof(EmuRequest));
}

void emu_request_free(EmuRequest *request)
{
    if (request == NULL) {
        return;
    }

    json_decref(request->payload);
    memset(request, 0, sizeof(EmuRequest));
}

const char *emu_msg_type_str(EmuMsgType type)
{
    int i;

    for (i = 0; g_msg_names[i].name != NULL; i++) {
        if (g_msg_names[i].type == type) {
            return g_msg_names[i].name;
        }
    }

    return "unknown";
}

EmuMsgType emu_msg_type_from_str(const char *type)
{
    int i;

    if (type == NULL) {
        return EmuMsgUnknown;
    }

    for (i = 0; g_msg_names[i].name != NULL; i++) {
        if (!strcmp(type, g_msg_names[i].name)) {
            return g_msg_names[i].type;
        }
    }

    return EmuMsgUnknown;
}

bool emu_request_from_json(JsonObj *json, EmuRequest *request)
{
    JsonObj *value = NULL;

    if (json == NULL || request == NULL || !json_is_object(json)) {
        return false;
    }

    emu_request_init(request);

    value = json_object_get(json, "id");
    snprintf(request->id,
             sizeof(request->id),
             "%s",
             json_is_string(value) ? json_string_value(value) : "req");

    value = json_object_get(json, "type");
    request->type = emu_msg_type_from_str(
        json_is_string(value) ? json_string_value(value) : NULL);

    if (request->type == EmuMsgUnknown) {
        return false;
    }

    value = json_object_get(json, "payload");
    request->payload = value ? json_deep_copy(value) : json_object();

    return request->payload != NULL;
}
