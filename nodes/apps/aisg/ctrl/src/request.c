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

static const struct {
    CtrlMsgType type;
    const char *name;
} gMsgNames[] = {
    { CtrlMsgPing,                  "ping" },
    { CtrlMsgGetStatus,             "get_status" },
    { CtrlMsgScan,                  "scan" },
    { CtrlMsgGetInfo,               "get_info" },
    { CtrlMsgGetAlarmStatus,        "get_alarm_status" },
    { CtrlMsgClearActiveAlarms,     "clear_active_alarms" },
    { CtrlMsgAlarmSubscribe,        "alarm_subscribe" },
    { CtrlMsgSelfTest,              "self_test" },
    { CtrlMsgSendConfigurationData, "send_configuration_data" },
    { CtrlMsgCalibrate,             "calibrate" },
    { CtrlMsgGetTilt,               "get_tilt" },
    { CtrlMsgSetTilt,               "set_tilt" },
    { CtrlMsgGetDeviceData,         "get_device_data" },
    { CtrlMsgResetSoftware,         "reset_software" },
    { CtrlMsgUnknown,               NULL }
};

void ctrl_request_init(CtrlRequest *request) {
    memset(request, 0, sizeof(CtrlRequest));
}

void ctrl_request_free(CtrlRequest *request) {
    json_decref(request->payload);
    memset(request, 0, sizeof(CtrlRequest));
}

const char *ctrl_msg_type_str(CtrlMsgType type) {
    int i;

    for (i = 0; gMsgNames[i].name != NULL; i++) {
        if (gMsgNames[i].type == type) return gMsgNames[i].name;
    }

    return "unknown";
}

CtrlMsgType ctrl_msg_type_from_str(const char *type) {
    int i;

    if (type == NULL) return CtrlMsgUnknown;

    for (i = 0; gMsgNames[i].name != NULL; i++) {
        if (!strcmp(type, gMsgNames[i].name)) return gMsgNames[i].type;
    }

    return CtrlMsgUnknown;
}

bool ctrl_request_from_json(JsonObj *json, CtrlRequest *request) {
    JsonObj *value;

    if (json == NULL || request == NULL || !json_is_object(json)) return false;

    ctrl_request_init(request);

    value = json_object_get(json, "id");
    snprintf(request->id,
             sizeof(request->id),
             "%s",
             json_is_string(value) ? json_string_value(value) : "req");

    value = json_object_get(json, "type");
    request->type = ctrl_msg_type_from_str(json_is_string(value) ?
                                           json_string_value(value) : NULL);
    if (request->type == CtrlMsgUnknown) return false;

    value = json_object_get(json, "payload");
    request->payload = value ? json_deep_copy(value) : json_object();

    return request->payload != NULL;
}
