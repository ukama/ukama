/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "aisg_v2.h"
#include "backend_raw_rs485.h"
#include "retap_ops.h"
#include "usys_log.h"

typedef struct {
    SerialPort serial;
    AisgBus bus;
    AisgDevice device;

    bool configured;
    bool calibrated;
    int16_t tiltTenthsDeg;
} RawRs485Context;

static bool read_file_bytes(const char *path,
                            uint8_t *buf,
                            size_t size,
                            size_t *len)
{
    FILE *file = NULL;
    size_t n;

    if (path == NULL || buf == NULL || len == NULL) {
        return false;
    }

    file = fopen(path, "rb");
    if (file == NULL) {
        return false;
    }

    n = fread(buf, 1, size, file);
    if (ferror(file)) {
        fclose(file);
        return false;
    }

    fclose(file);

    *len = n;

    return true;
}

static JsonObj *build_status_payload(RawRs485Context *ctx)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "mode", json_string("operating"));
    json_object_set_new(json, "busy", json_boolean(false));
    json_object_set_new(json,
                        "present",
                        json_boolean(ctx != NULL && ctx->device.present));
    json_object_set_new(json,
                        "configured",
                        json_boolean(ctx != NULL && ctx->configured));
    json_object_set_new(json,
                        "calibrated",
                        json_boolean(ctx != NULL && ctx->calibrated));
    json_object_set_new(json, "powerManaged", json_boolean(false));
    json_object_set_new(json, "transport", json_string("raw-rs485"));

    return json;
}

static CtrlCode retap_response_code(RetapResponse *response)
{
    if (response == NULL) {
        return CtrlCodeTransportError;
    }

    return retap_failure_to_ctrl_code(response->failureReason);
}

static bool execute_retap(RawRs485Context *ctx,
                          RetapRequest *request,
                          RetapResponse *response,
                          CtrlResponse *ctrlResp)
{
    CtrlCode code;

    if (!aisg_v2_send_retap(&ctx->bus, request, response)) {
        return ctrl_response_set_error(ctrlResp,
                                       CtrlCodeTransportError,
                                       "failed to execute RETAP");
    }

    if (response->returnCode != RETAP_RETURN_FAIL) {
        return true;
    }

    code = retap_response_code(response);

    return ctrl_response_set_error(ctrlResp, code, ctrl_code_str(code));
}

static JsonObj *build_ok_payload(void)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "ok", json_true());

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

static JsonObj *build_info_payload(RetapInfo *info)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json,
                        "productNumber",
                        json_string(info->productNumber));
    json_object_set_new(json, "serialNumber", json_string(info->serialNumber));
    json_object_set_new(json,
                        "hardwareVersion",
                        json_string(info->hardwareVersion));
    json_object_set_new(json,
                        "softwareVersion",
                        json_string(info->softwareVersion));

    return json;
}

static JsonObj *build_alarm_payload(RetapAlarmList *alarms)
{
    JsonObj *json = NULL;
    JsonObj *array = NULL;
    JsonObj *alarm = NULL;
    int i;

    array = json_array();
    if (array == NULL) {
        return NULL;
    }

    for (i = 0; i < alarms->count; i++) {
        alarm = json_object();
        if (alarm == NULL) {
            json_decref(array);
            return NULL;
        }

        json_object_set_new(alarm, "code", json_integer(alarms->codes[i]));
        json_object_set_new(alarm,
                            "name",
                            json_string(retap_return_code_str(alarms->codes[i])));
        json_object_set_new(alarm, "state", json_string("raised"));
        json_array_append_new(array, alarm);
    }

    json = json_object();
    if (json == NULL) {
        json_decref(array);
        return NULL;
    }

    json_object_set_new(json, "active", array);

    return json;
}

static JsonObj *build_tilt_payload(int16_t tilt)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "currentTiltDeg", json_real(tilt / 10.0));
    json_object_set_new(json, "rawTiltTenthsDeg", json_integer(tilt));

    return json;
}

static JsonObj *build_device_data_payload(int field)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "field", json_integer(field));
    json_object_set_new(json, "dataHex", json_string(""));

    return json;
}

static bool raw_handle_ping(RawRs485Context *ctx, CtrlResponse *response)
{
    (void)ctx;

    return ctrl_response_set_ok(response, json_object());
}

static bool raw_handle_status(RawRs485Context *ctx, CtrlResponse *response)
{
    return ctrl_response_set_ok(response, build_status_payload(ctx));
}

static bool raw_handle_scan(RawRs485Context *ctx, CtrlResponse *response)
{
    if (!aisg_v2_scan(&ctx->bus, &ctx->device)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeTransportError,
                                       "scan failed");
    }

    return ctrl_response_set_ok(response, build_status_payload(ctx));
}

static bool raw_handle_get_info(RawRs485Context *ctx, CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    RetapInfo info;

    retap_build_get_information(&request);

    if (!execute_retap(ctx, &request, &retapResp, response)) {
        return false;
    }

    memset(&info, 0, sizeof(info));
    if (!retap_parse_get_information(&retapResp, &info)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeFormatError,
                                       "failed to parse information");
    }

    return ctrl_response_set_ok(response, build_info_payload(&info));
}

static bool raw_handle_get_alarms(RawRs485Context *ctx,
                                  CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    RetapAlarmList alarms;

    retap_build_get_error_status(&request);

    if (!execute_retap(ctx, &request, &retapResp, response)) {
        return false;
    }

    memset(&alarms, 0, sizeof(alarms));
    if (!retap_parse_return_code_list(&retapResp, &alarms)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeFormatError,
                                       "failed to parse error status");
    }

    return ctrl_response_set_ok(response, build_alarm_payload(&alarms));
}

static bool raw_handle_simple(RawRs485Context *ctx,
                              CtrlResponse *response,
                              bool (*build)(RetapRequest *request),
                              JsonObj *payload)
{
    RetapRequest request;
    RetapResponse retapResp;

    if (!build(&request)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "failed to build RETAP request");
    }

    if (!execute_retap(ctx, &request, &retapResp, response)) {
        json_decref(payload);
        return false;
    }

    return ctrl_response_set_ok(response, payload ? payload : json_object());
}

static bool raw_handle_self_test(RawRs485Context *ctx,
                                 CtrlResponse *response)
{
    return raw_handle_simple(
        ctx,
        response,
        retap_build_self_test,
        build_operation_payload("op-selftest-001", "self-test"));
}

static bool read_config_blob(CtrlRequest *request,
                             uint8_t *data,
                             size_t size,
                             size_t *len,
                             CtrlResponse *response)
{
    JsonObj *value = NULL;
    const char *path = NULL;

    value = json_object_get(request->payload, "configPath");
    path = json_is_string(value) ? json_string_value(value) : NULL;

    if (path == NULL || path[0] == '\0') {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "missing configPath");
    }

    if (!read_file_bytes(path, data, size, len)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "failed to read config blob");
    }

    return true;
}

static bool raw_handle_send_config(RawRs485Context *ctx,
                                   CtrlRequest *request,
                                   CtrlResponse *response)
{
    RetapRequest retapReq;
    RetapResponse retapResp;
    uint8_t data[RETAP_MAX_PAYLOAD];
    size_t len;
    JsonObj *payload = NULL;

    if (!read_config_blob(request, data, sizeof(data), &len, response)) {
        return false;
    }

    if (!retap_build_send_configuration_data(&retapReq, data, len)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "failed to build config request");
    }

    if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
        return false;
    }

    ctx->configured = true;

    payload = build_ok_payload();
    if (payload == NULL) {
        return false;
    }

    json_object_set_new(payload, "configured", json_true());

    return ctrl_response_set_ok(response, payload);
}

static bool raw_handle_calibrate(RawRs485Context *ctx,
                                 CtrlResponse *response)
{
    bool ok;

    ok = raw_handle_simple(
        ctx,
        response,
        retap_build_calibrate,
        build_operation_payload("op-cal-001", "calibrate"));

    if (ok) {
        ctx->calibrated = true;
    }

    return ok;
}

static bool raw_handle_get_tilt(RawRs485Context *ctx,
                                CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    int16_t tilt;

    retap_build_get_tilt(&request);

    if (!execute_retap(ctx, &request, &retapResp, response)) {
        return false;
    }

    if (!retap_parse_get_tilt(&retapResp, &tilt)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeFormatError,
                                       "failed to parse tilt");
    }

    ctx->tiltTenthsDeg = tilt;

    return ctrl_response_set_ok(response, build_tilt_payload(tilt));
}

static bool raw_handle_set_tilt(RawRs485Context *ctx,
                                CtrlRequest *request,
                                CtrlResponse *response)
{
    RetapRequest retapReq;
    RetapResponse retapResp;
    JsonObj *value = NULL;
    JsonObj *payload = NULL;
    double target;
    int16_t tilt;

    value = json_object_get(request->payload, "targetTiltDeg");
    target = json_number_value(value);
    tilt = (int16_t)(target * 10.0);

    retap_build_set_tilt(&retapReq, tilt);

    if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
        return false;
    }

    ctx->tiltTenthsDeg = tilt;

    payload = build_operation_payload("op-tilt-001", "set-tilt");
    if (payload == NULL) {
        return false;
    }

    json_object_set_new(payload, "targetTiltDeg", json_real(target));

    return ctrl_response_set_ok(response, payload);
}

static bool raw_handle_get_device_data(RawRs485Context *ctx,
                                       CtrlRequest *request,
                                       CtrlResponse *response)
{
    RetapRequest retapReq;
    RetapResponse retapResp;
    JsonObj *value = NULL;
    int field;

    value = json_object_get(request->payload, "field");
    field = (int)json_integer_value(value);

    retap_build_get_device_data(&retapReq, (uint8_t)field);

    if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
        return false;
    }

    return ctrl_response_set_ok(response, build_device_data_payload(field));
}

static bool raw_execute(Backend *backend,
                        CtrlRequest *request,
                        CtrlResponse *response)
{
    RawRs485Context *ctx = NULL;

    if (backend == NULL || request == NULL || response == NULL) {
        return false;
    }

    ctx = backend->priv;

    switch (request->type) {
    case CtrlMsgPing:
        return raw_handle_ping(ctx, response);
    case CtrlMsgGetStatus:
        return raw_handle_status(ctx, response);
    case CtrlMsgScan:
        return raw_handle_scan(ctx, response);
    case CtrlMsgGetInfo:
        return raw_handle_get_info(ctx, response);
    case CtrlMsgGetAlarmStatus:
        return raw_handle_get_alarms(ctx, response);
    case CtrlMsgClearActiveAlarms:
        return raw_handle_simple(
            ctx, response, retap_build_clear_active_alarms, NULL);
    case CtrlMsgAlarmSubscribe:
        return raw_handle_simple(
            ctx, response, retap_build_alarm_subscribe, NULL);
    case CtrlMsgSelfTest:
        return raw_handle_self_test(ctx, response);
    case CtrlMsgSendConfigurationData:
        return raw_handle_send_config(ctx, request, response);
    case CtrlMsgCalibrate:
        return raw_handle_calibrate(ctx, response);
    case CtrlMsgGetTilt:
        return raw_handle_get_tilt(ctx, response);
    case CtrlMsgSetTilt:
        return raw_handle_set_tilt(ctx, request, response);
    case CtrlMsgGetDeviceData:
        return raw_handle_get_device_data(ctx, request, response);
    case CtrlMsgResetSoftware:
        return raw_handle_simple(
            ctx, response, retap_build_reset_software, NULL);
    default:
        return ctrl_response_set_error(response,
                                       CtrlCodeUnsupportedProcedure,
                                       "unsupported request type");
    }
}

static bool raw_open(Backend *backend)
{
    RawRs485Context *ctx = NULL;

    ctx = calloc(1, sizeof(RawRs485Context));
    if (ctx == NULL) {
        return false;
    }

    if (!serial_open(&ctx->serial,
                     backend->config->rawDevice,
                     backend->config->rawBaud)) {
        free(ctx);
        return false;
    }

    aisg_v2_bus_init(&ctx->bus, &ctx->serial);

    ctx->tiltTenthsDeg = 0;
    backend->priv = ctx;

    return true;
}

static void raw_close(Backend *backend)
{
    RawRs485Context *ctx = NULL;

    if (backend == NULL || backend->priv == NULL) {
        return;
    }

    ctx = backend->priv;
    serial_close(&ctx->serial);
    free(ctx);
    backend->priv = NULL;
}

bool backend_raw_rs485_init(Backend *backend, Config *config)
{
    static BackendOps ops = {
        .open    = raw_open,
        .close   = raw_close,
        .execute = raw_execute,
    };

    if (backend == NULL || config == NULL) {
        return false;
    }

    memset(backend, 0, sizeof(Backend));

    backend->config = config;
    backend->ops    = ops;

    return true;
}
