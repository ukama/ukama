/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "ops.h"
#include "version.h"

#define CTRL_CODE_OK              "OK"
#define CTRL_CODE_BUSY            "Busy"
#define CTRL_CODE_NOT_CONFIGURED  "NotConfigured"
#define CTRL_CODE_NOT_CALIBRATED  "NotCalibrated"
#define CTRL_CODE_UNSUPPORTED     "UnsupportedProcedure"

static JsonObj *new_ok_payload(void)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "ok", json_true());

    return json;
}

static JsonObj *build_info_payload(void)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "productNumber", json_string("RET1T1"));
    json_object_set_new(json, "serialNumber", json_string("EMU-0001"));
    json_object_set_new(json, "hardwareVersion", json_string("A1"));
    json_object_set_new(json, "softwareVersion", json_string(VERSION));

    return json;
}

static JsonObj *build_operation_payload(const char *id, const char *type)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "operationId", json_string(id));
    json_object_set_new(json, "state", json_string("running"));
    json_object_set_new(json, "type", json_string(type));

    return json;
}

static JsonObj *build_tilt_payload(EmuModel *model)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json,
                        "currentTiltDeg",
                        json_real(model->tiltTenthsDeg / 10.0));
    json_object_set_new(json,
                        "rawTiltTenthsDeg",
                        json_integer(model->tiltTenthsDeg));

    return json;
}

static JsonObj *build_device_data_payload(JsonObj *payload)
{
    JsonObj *json = NULL;
    JsonObj *value = NULL;
    int field;

    value = json_object_get(payload, "field");
    field = (int)json_integer_value(value);

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "field", json_integer(field));
    json_object_set_new(json, "dataHex", json_string(""));

    return json;
}

static bool fail_if_busy(EmuModel *model,
                         const EmuRequest *request,
                         EmuResponse *response)
{
    if (!model->busy) {
        return false;
    }

    if (request->type == EmuMsgPing || request->type == EmuMsgGetStatus) {
        return false;
    }

    emu_response_set_error(response, CTRL_CODE_BUSY, "device busy");

    return true;
}

static bool handle_scan(EmuModel *model, EmuResponse *response)
{
    if (!model->present) {
        return emu_response_set_error(response,
                                      "DeviceMissing",
                                      "device not present");
    }

    return emu_response_set_ok(response, emu_model_status(model));
}

static bool handle_alarm_status(EmuModel *model, EmuResponse *response)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return false;
    }

    json_object_set_new(json, "active", json_deep_copy(model->alarms));

    return emu_response_set_ok(response, json);
}

static bool handle_clear_alarms(EmuModel *model, EmuResponse *response)
{
    json_decref(model->alarms);
    model->alarms = json_array();

    return emu_response_set_ok(response, new_ok_payload());
}

static bool handle_alarm_subscribe(EmuModel *model, EmuResponse *response)
{
    JsonObj *json = NULL;

    model->alarmSubscribed = true;

    json = new_ok_payload();
    if (json == NULL) {
        return false;
    }

    json_object_set_new(json, "subscribed", json_true());

    return emu_response_set_ok(response, json);
}

static bool handle_configure(EmuModel *model, EmuResponse *response)
{
    JsonObj *json = NULL;

    model->configured = true;

    json = new_ok_payload();
    if (json == NULL) {
        return false;
    }

    json_object_set_new(json, "configured", json_true());

    return emu_response_set_ok(response, json);
}

static bool handle_calibrate(EmuModel *model, EmuResponse *response)
{
    if (!model->configured) {
        return emu_response_set_error(response,
                                      CTRL_CODE_NOT_CONFIGURED,
                                      "device is not configured");
    }

    model->calibrated = true;

    return emu_response_set_ok(
        response,
        build_operation_payload("op-cal-001", "calibrate"));
}

static bool handle_set_tilt(EmuModel *model,
                            const EmuRequest *request,
                            EmuResponse *response)
{
    JsonObj *value = NULL;
    JsonObj *json = NULL;
    double target;

    if (!model->calibrated) {
        return emu_response_set_error(response,
                                      CTRL_CODE_NOT_CALIBRATED,
                                      "device is not calibrated");
    }

    value = json_object_get(request->payload, "targetTiltDeg");
    target = json_number_value(value);
    model->tiltTenthsDeg = (int16_t)(target * 10.0);

    json = build_operation_payload("op-tilt-001", "set-tilt");
    if (json == NULL) {
        return false;
    }

    json_object_set_new(json, "targetTiltDeg", json_real(target));

    return emu_response_set_ok(response, json);
}

static bool handle_reset(EmuModel *model, EmuResponse *response)
{
    model->configured = false;
    model->calibrated = false;

    return emu_response_set_ok(response, new_ok_payload());
}

bool emu_ops_handle(EmuModel *model,
                    const EmuRequest *request,
                    EmuResponse *response)
{
    if (model == NULL || request == NULL || response == NULL) {
        return false;
    }

    if (fail_if_busy(model, request, response)) {
        return true;
    }

    switch (request->type) {
    case EmuMsgPing:
        return emu_response_set_ok(response, json_object());

    case EmuMsgGetStatus:
        return emu_response_set_ok(response, emu_model_status(model));

    case EmuMsgScan:
        return handle_scan(model, response);

    case EmuMsgGetInfo:
        return emu_response_set_ok(response, build_info_payload());

    case EmuMsgGetAlarmStatus:
        return handle_alarm_status(model, response);

    case EmuMsgClearActiveAlarms:
        return handle_clear_alarms(model, response);

    case EmuMsgAlarmSubscribe:
        return handle_alarm_subscribe(model, response);

    case EmuMsgSelfTest:
        return emu_response_set_ok(
            response,
            build_operation_payload("op-selftest-001", "self-test"));

    case EmuMsgSendConfigurationData:
        return handle_configure(model, response);

    case EmuMsgCalibrate:
        return handle_calibrate(model, response);

    case EmuMsgGetTilt:
        return emu_response_set_ok(response, build_tilt_payload(model));

    case EmuMsgSetTilt:
        return handle_set_tilt(model, request, response);

    case EmuMsgGetDeviceData:
        return emu_response_set_ok(
            response,
            build_device_data_payload(request->payload));

    case EmuMsgResetSoftware:
        return handle_reset(model, response);

    default:
        return emu_response_set_error(response,
                                      CTRL_CODE_UNSUPPORTED,
                                      "unsupported request type");
    }
}
