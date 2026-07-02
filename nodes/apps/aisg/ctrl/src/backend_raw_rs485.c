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

#define RAW_CONFIG_FILE_MAX_BYTES          (256U * 1024U)

typedef struct {
    SerialPort serial;
    AisgBus bus;
    AisgDevice device;

    bool identified;
    bool configured;
    bool calibrated;
    bool tiltKnown;
    bool targetTiltKnown;
    int16_t tiltTenthsDeg;
    int16_t targetTiltTenthsDeg;
    char productNumber[64];
    char serialNumber[64];
    char hardwareVersion[64];
    char softwareVersion[64];
} RawRs485Context;

static bool read_file_alloc(const char *path,
                            uint8_t **data,
                            size_t *len)
{
    FILE *file = NULL;
    uint8_t *buf = NULL;
    long fileLen;
    size_t n;

    if (path == NULL || data == NULL || len == NULL) {
        return false;
    }

    *data = NULL;
    *len = 0;

    file = fopen(path, "rb");
    if (file == NULL) {
        return false;
    }

    if (fseek(file, 0, SEEK_END) != 0) {
        fclose(file);
        return false;
    }

    fileLen = ftell(file);
    if (fileLen <= 0 || (unsigned long)fileLen > RAW_CONFIG_FILE_MAX_BYTES) {
        fclose(file);
        return false;
    }

    if (fseek(file, 0, SEEK_SET) != 0) {
        fclose(file);
        return false;
    }

    buf = malloc((size_t)fileLen);
    if (buf == NULL) {
        fclose(file);
        return false;
    }

    n = fread(buf, 1, (size_t)fileLen, file);
    if (ferror(file) || n != (size_t)fileLen) {
        free(buf);
        fclose(file);
        return false;
    }

    fclose(file);

    *data = buf;
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
                        "identified",
                        json_boolean(ctx != NULL && ctx->identified));
    json_object_set_new(json,
                        "configured",
                        json_boolean(ctx != NULL && ctx->configured));
    json_object_set_new(json,
                        "calibrated",
                        json_boolean(ctx != NULL && ctx->calibrated));
    json_object_set_new(json, "powerManaged", json_boolean(false));
    json_object_set_new(json, "transport", json_string("raw-rs485"));

    if (ctx != NULL) {
        json_object_set_new(json,
                            "model",
                            json_string(ctx->productNumber));
        json_object_set_new(json,
                            "productNumber",
                            json_string(ctx->productNumber));
        json_object_set_new(json,
                            "serialNumber",
                            json_string(ctx->serialNumber));
        json_object_set_new(json,
                            "hardwareVersion",
                            json_string(ctx->hardwareVersion));
        json_object_set_new(json,
                            "softwareVersion",
                            json_string(ctx->softwareVersion));
        json_object_set_new(json, "tiltKnown", json_boolean(ctx->tiltKnown));
        if (ctx->tiltKnown) {
            json_object_set_new(json,
                                "currentTiltDeg",
                                json_real(ctx->tiltTenthsDeg / 10.0));
        } else {
            json_object_set_new(json, "currentTiltDeg", json_null());
        }
        json_object_set_new(json,
                            "linkState",
                            json_string(aisg_v2_l2_state_str(ctx->bus.state)));
        json_object_set_new(json,
                            "lastLinkError",
                            json_string(aisg_v2_error_str(ctx->bus.lastError)));
        json_object_set_new(json,
                            "hdlcMaxInfoLen",
                            json_integer((json_int_t)ctx->bus.maxInfoLen));
        json_object_set_new(json,
                            "targetTiltKnown",
                            json_boolean(ctx->targetTiltKnown));
        if (ctx->targetTiltKnown) {
            json_object_set_new(json,
                                "targetTiltDeg",
                                json_real(ctx->targetTiltTenthsDeg / 10.0));
        } else {
            json_object_set_new(json, "targetTiltDeg", json_null());
        }
    }

    return json;
}

static CtrlCode retap_response_code(RetapResponse *response)
{
    if (response == NULL) {
        return CtrlCodeTransportError;
    }

    return retap_failure_to_ctrl_code(response->failureReason);
}


static CtrlCode aisg_error_to_ctrl_code(AisgError error)
{
    switch (error) {
    case AISG_ERROR_NONE:
        return CtrlCodeOk;
    case AISG_ERROR_MULTIPLE_DEVICES:
        return CtrlCodeMultipleDevices;
    case AISG_ERROR_UNSUPPORTED_DEVICE_TYPE:
        return CtrlCodeUnsupportedDeviceType;
    case AISG_ERROR_UNSUPPORTED_PROTOCOL_VERSION:
        return CtrlCodeUnsupportedProtocolVersion;
    case AISG_ERROR_LINK_NOT_CONNECTED:
        return CtrlCodeLinkNotConnected;
    case AISG_ERROR_FRAME_REJECT:
        return CtrlCodeFrameReject;
    case AISG_ERROR_RECEIVER_NOT_READY:
        return CtrlCodeReceiverNotReady;
    case AISG_ERROR_PROTOCOL:
        return CtrlCodeProtocolError;
    case AISG_ERROR_TIMEOUT:
        return CtrlCodeTimeout;
    case AISG_ERROR_TRANSPORT:
    default:
        return CtrlCodeTransportError;
    }
}

static bool ctrl_response_set_aisg_error(CtrlResponse *response,
                                         AisgError error,
                                         const char *fallback)
{
    CtrlCode code;
    char reason[CTRL_REASON_LEN];

    code = aisg_error_to_ctrl_code(error);
    if (code == CtrlCodeOk) {
        code = CtrlCodeTransportError;
    }

    snprintf(reason,
             sizeof(reason),
             "%s%s%s",
             fallback ? fallback : ctrl_code_str(code),
             error == AISG_ERROR_NONE ? "" : ": ",
             error == AISG_ERROR_NONE ? "" : aisg_v2_error_str(error));

    return ctrl_response_set_error(response, code, reason);
}

static bool execute_retap(RawRs485Context *ctx,
                          RetapRequest *request,
                          RetapResponse *response,
                          CtrlResponse *ctrlResp)
{
    CtrlCode code;

    if (!aisg_v2_send_retap(&ctx->bus, request, response)) {
        return ctrl_response_set_aisg_error(ctrlResp,
                                            ctx->bus.lastError,
                                            "failed to execute RETAP");
    }

    if (response->returnCode != RETAP_RETURN_FAIL) {
        return true;
    }

    code = retap_response_code(response);

    return ctrl_response_set_error(ctrlResp, code, ctrl_code_str(code));
}

static bool raw_require_connected(RawRs485Context *ctx, CtrlResponse *response)
{
    if (ctx != NULL && ctx->device.present &&
        ctx->bus.state == AISG_L2_CONNECTED) {
        return true;
    }

    if (ctx != NULL) {
        ctx->bus.lastError = AISG_ERROR_LINK_NOT_CONNECTED;
    }

    return ctrl_response_set_error(response,
                                   CtrlCodeLinkNotConnected,
                                   "device not connected; run scan first");
}

static bool raw_require_configured(RawRs485Context *ctx, CtrlResponse *response)
{
    if (!raw_require_connected(ctx, response)) {
        return false;
    }

    if (ctx->configured) {
        return true;
    }

    return ctrl_response_set_error(response,
                                   CtrlCodeNotConfigured,
                                   "configuration must be loaded first");
}

static bool raw_require_calibrated(RawRs485Context *ctx, CtrlResponse *response)
{
    if (!raw_require_configured(ctx, response)) {
        return false;
    }

    if (ctx->calibrated) {
        return true;
    }

    return ctrl_response_set_error(response,
                                   CtrlCodeNotCalibrated,
                                   "calibration must be completed first");
}

static void raw_clear_device_runtime_state(RawRs485Context *ctx)
{
    if (ctx == NULL) {
        return;
    }

    ctx->device.present = false;
    ctx->identified = false;
    ctx->configured = false;
    ctx->calibrated = false;
    ctx->tiltKnown = false;
    ctx->targetTiltKnown = false;
    ctx->tiltTenthsDeg = 0;
    ctx->targetTiltTenthsDeg = 0;
    ctx->productNumber[0] = '\0';
    ctx->serialNumber[0] = '\0';
    ctx->hardwareVersion[0] = '\0';
    ctx->softwareVersion[0] = '\0';
    aisg_v2_bus_reset_link(&ctx->bus);
}

static void raw_update_identity(RawRs485Context *ctx, RetapInfo *info)
{
    if (ctx == NULL || info == NULL) {
        return;
    }

    snprintf(ctx->productNumber,
             sizeof(ctx->productNumber),
             "%s",
             info->productNumber);
    snprintf(ctx->serialNumber,
             sizeof(ctx->serialNumber),
             "%s",
             info->serialNumber);
    snprintf(ctx->hardwareVersion,
             sizeof(ctx->hardwareVersion),
             "%s",
             info->hardwareVersion);
    snprintf(ctx->softwareVersion,
             sizeof(ctx->softwareVersion),
             "%s",
             info->softwareVersion);

    ctx->identified = true;
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
        return ctrl_response_set_aisg_error(response,
                                            ctx->bus.lastError,
                                            "scan failed");
    }

    return ctrl_response_set_ok(response, build_status_payload(ctx));
}

static bool raw_handle_get_info(RawRs485Context *ctx, CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    RetapInfo info;

    if (!raw_require_connected(ctx, response)) {
        return false;
    }

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

    raw_update_identity(ctx, &info);

    return ctrl_response_set_ok(response, build_info_payload(&info));
}

static bool raw_handle_get_alarms(RawRs485Context *ctx,
                                  CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    RetapAlarmList alarms;

    if (!raw_require_connected(ctx, response)) {
        return false;
    }

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

    if (!raw_require_connected(ctx, response)) {
        json_decref(payload);
        return false;
    }

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
                             uint8_t **data,
                             size_t *len,
                             CtrlResponse *response)
{
    JsonObj *value = NULL;
    const char *path = NULL;

    if (request == NULL || data == NULL || len == NULL) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "invalid config request");
    }

    value = json_object_get(request->payload, "configPath");
    path = json_is_string(value) ? json_string_value(value) : NULL;

    if (path == NULL || path[0] == '\0') {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "missing configPath");
    }

    if (!read_file_alloc(path, data, len)) {
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
    uint8_t *data = NULL;
    size_t len = 0;
    size_t off = 0;
    size_t chunkLen;
    size_t chunks = 0;
    size_t totalChunks;
    JsonObj *payload = NULL;

    if (!raw_require_connected(ctx, response)) {
        return false;
    }

    if (!read_config_blob(request, &data, &len, response)) {
        return false;
    }

    totalChunks = (len + RETAP_CONFIG_SEGMENT_MAX - 1) /
                  RETAP_CONFIG_SEGMENT_MAX;

    while (off < len) {
        chunkLen = len - off;
        if (chunkLen > RETAP_CONFIG_SEGMENT_MAX) {
            chunkLen = RETAP_CONFIG_SEGMENT_MAX;
        }

        if (!retap_build_send_configuration_data(&retapReq,
                                                 &data[off],
                                                 chunkLen)) {
            free(data);
            return ctrl_response_set_error(response,
                                           CtrlCodeInvalidRequest,
                                           "failed to build config segment request");
        }

        usys_log_debug("aisg: send config segment %zu/%zu offset=%zu len=%zu",
                       chunks + 1,
                       totalChunks,
                       off,
                       chunkLen);

        if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
            free(data);
            ctx->configured = false;
            ctx->calibrated = false;
            return false;
        }

        off += chunkLen;
        chunks++;
    }

    free(data);

    ctx->configured = true;
    ctx->calibrated = false;
    ctx->tiltKnown = false;
    ctx->targetTiltKnown = false;

    payload = build_ok_payload();
    if (payload == NULL) {
        return false;
    }

    json_object_set_new(payload, "configured", json_true());
    json_object_set_new(payload, "calibrated", json_false());
    json_object_set_new(payload, "bytes", json_integer((json_int_t)len));
    json_object_set_new(payload, "chunks", json_integer((json_int_t)chunks));
    json_object_set_new(payload,
                        "totalChunks",
                        json_integer((json_int_t)totalChunks));
    json_object_set_new(payload,
                        "segmentMaxBytes",
                        json_integer(RETAP_CONFIG_SEGMENT_MAX));

    return ctrl_response_set_ok(response, payload);
}

static bool raw_handle_calibrate(RawRs485Context *ctx,
                                 CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    JsonObj *payload = NULL;

    if (!raw_require_configured(ctx, response)) {
        return false;
    }

    if (!retap_build_calibrate(&request)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "failed to build calibrate request");
    }

    if (!execute_retap(ctx, &request, &retapResp, response)) {
        ctx->calibrated = false;
        return false;
    }

    ctx->calibrated = true;

    payload = build_operation_payload("op-cal-001", "calibrate");
    if (payload == NULL) {
        return false;
    }

    json_object_set_new(payload, "state", json_string("completed"));
    json_object_set_new(payload, "configured", json_true());
    json_object_set_new(payload, "calibrated", json_true());

    return ctrl_response_set_ok(response, payload);
}

static bool raw_handle_get_tilt(RawRs485Context *ctx,
                                CtrlResponse *response)
{
    RetapRequest request;
    RetapResponse retapResp;
    int16_t tilt;

    if (!raw_require_connected(ctx, response)) {
        return false;
    }

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
    ctx->tiltKnown = true;

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

    if (!raw_require_calibrated(ctx, response)) {
        return false;
    }

    value = json_object_get(request->payload, "targetTiltDeg");
    if (!json_is_number(value)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "missing targetTiltDeg");
    }

    target = json_number_value(value);
    tilt = (int16_t)(target * 10.0);

    retap_build_set_tilt(&retapReq, tilt);

    if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
        return false;
    }

    ctx->tiltTenthsDeg = tilt;
    ctx->targetTiltTenthsDeg = tilt;
    ctx->tiltKnown = true;
    ctx->targetTiltKnown = true;

    payload = build_operation_payload("op-tilt-001", "set-tilt");
    if (payload == NULL) {
        return false;
    }

    json_object_set_new(payload, "targetTiltDeg", json_real(target));
    json_object_set_new(payload, "currentTiltDeg", json_real(tilt / 10.0));
    json_object_set_new(payload, "rawTiltTenthsDeg", json_integer(tilt));

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

    if (!raw_require_connected(ctx, response)) {
        return false;
    }

    value = json_object_get(request->payload, "field");
    if (!json_is_integer(value)) {
        return ctrl_response_set_error(response,
                                       CtrlCodeInvalidRequest,
                                       "missing field");
    }
    field = (int)json_integer_value(value);

    retap_build_get_device_data(&retapReq, (uint8_t)field);

    if (!execute_retap(ctx, &retapReq, &retapResp, response)) {
        return false;
    }

    return ctrl_response_set_ok(response, build_device_data_payload(field));
}


static bool raw_handle_reset(RawRs485Context *ctx, CtrlResponse *response)
{
    bool ok;

    ok = raw_handle_simple(ctx, response, retap_build_reset_software, NULL);
    if (ok) {
        raw_clear_device_runtime_state(ctx);
    }

    return ok;
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
        return raw_handle_reset(ctx, response);
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
