/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "ops.h"

#define CTRL_GET_STATUS       "get_status"
#define CTRL_SCAN             "scan"
#define CTRL_GET_INFO         "get_info"
#define CTRL_GET_ALARMS       "get_alarm_status"
#define CTRL_CLEAR_ALARMS     "clear_active_alarms"
#define CTRL_SUBSCRIBE_ALARMS "alarm_subscribe"
#define CTRL_SELF_TEST        "self_test"
#define CTRL_CONFIGURE        "send_configuration_data"
#define CTRL_CALIBRATE        "calibrate"
#define CTRL_GET_TILT         "get_tilt"
#define CTRL_SET_TILT         "set_tilt"
#define CTRL_GET_DEVICE_DATA  "get_device_data"
#define CTRL_RESET_SOFTWARE   "reset_software"

#define OP_SCAN               "scan"
#define OP_CONFIG             "configure"
#define OP_CALIBRATE          "calibrate"
#define OP_GET_TILT           "get-tilt"
#define OP_SET_TILT           "set-tilt"
#define OP_SELF_TEST          "self-test"

static JsonObj *empty_payload(void) {
    return json_object();
}

static bool call_controller(AisgdContext *ctx,
                            const char *type,
                            JsonObj *payload,
                            JsonObj **response) {

    CtrlResponse ctrlResp;

    if (ctx == NULL || type == NULL || response == NULL) {
        json_decref(payload);
        return false;
    }

    if (!controller_client_call(&ctx->controller, type, payload, &ctrlResp)) {
        status_mark_controller_down(ctx->status, "controller unavailable");
        return false;
    }

    if (!ctrlResp.ok) {
        status_set(ctx->status, AisgdStateDegraded, ctrlResp.reason);
        ctrl_response_free(&ctrlResp);
        return false;
    }

    *response = ctrl_response_steal_payload(&ctrlResp);
    ctrl_response_free(&ctrlResp);

    return true;
}

static bool controller_get_status(AisgdContext *ctx) {

    JsonObj *payload = NULL;

    status_set(ctx->status,
               AisgdStateConnectController,
               "connect-controller");

    if (!call_controller(ctx, CTRL_GET_STATUS, empty_payload(), &payload)) {
        return false;
    }

    status_update_from_controller(ctx->status, payload);
    json_decref(payload);

    return true;
}

static bool reconcile_scan(AisgdContext *ctx)
{
    JsonObj *payload = NULL;

    status_set(ctx->status, AisgdStateScanDevice, "scan-device");

    if (!aisgd_ops_scan(ctx, &payload)) {
        return false;
    }

    json_decref(payload);

    return true;
}

static bool reconcile_alarm_subscribe(AisgdContext *ctx)
{
    JsonObj *payload = NULL;

    status_set(ctx->status,
               AisgdStateSubscribeAlarms,
               "subscribe-alarms");

    if (!aisgd_ops_subscribe_alarms(ctx, &payload)) {
        return false;
    }

    json_decref(payload);

    return true;
}

static JsonObj *build_reconcile_response(void)
{
    JsonObj *json    = NULL;
    JsonObj *actions = NULL;

    actions = json_array();
    if (actions == NULL) {
        return NULL;
    }

    json_array_append_new(actions, json_string("controller-ok"));
    json_array_append_new(actions, json_string("scan-complete"));
    json_array_append_new(actions, json_string("alarm-subscribed"));

    json = json_object();
    if (json == NULL) {
        json_decref(actions);
        return NULL;
    }

    json_object_set_new(json, "state", json_string("ready"));
    json_object_set_new(json, "actions", actions);

    return json;
}

static JsonObj *build_device_response(AisgdContext *ctx)
{
    JsonObj *json = NULL;

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "deviceId", json_string(AISGD_DEVICE_ID));
    json_object_set_new(json,
                        "present",
                        json_boolean(ctx->status->devicePresent));
    json_object_set_new(json, "model", json_string(ctx->status->model));
    json_object_set_new(json, "protocol", json_string("AISGv2+3GPP"));

    return json;
}

static JsonObj *build_config_payload(const char *profile, const char *path)
{
    JsonObj *payload = NULL;

    payload = json_object();
    if (payload == NULL) {
        return NULL;
    }

    json_object_set_new(payload,
                        "profile",
                        json_string(profile ? profile : ""));
    json_object_set_new(payload,
                        "configPath",
                        json_string(path ? path : ""));

    return payload;
}

static JsonObj *build_tilt_payload(double targetTiltDeg)
{
    JsonObj *payload = NULL;

    payload = json_object();
    if (payload == NULL) {
        return NULL;
    }

    json_object_set_new(payload,
                        "targetTiltDeg",
                        json_real(targetTiltDeg));

    return payload;
}

static bool json_number_at(JsonObj *json, const char *key)
{
    JsonObj *value = NULL;

    if (json == NULL || key == NULL) {
        return false;
    }

    value = json_object_get(json, key);
    return json_is_number(value);
}

static JsonObj *json_object_child(JsonObj *json, const char *key)
{
    JsonObj *value = NULL;

    if (json == NULL || key == NULL) {
        return NULL;
    }

    value = json_object_get(json, key);
    if (!json_is_object(value)) {
        return NULL;
    }

    return value;
}

static bool payload_has_current_tilt(JsonObj *payload)
{
    static const char *keys[] = {
        "currentTiltDeg",
        "tiltDeg",
        "electricalTiltDeg",
        "tilt"
    };
    JsonObj *child = NULL;
    size_t i;

    if (payload == NULL) {
        return false;
    }

    for (i = 0; i < sizeof(keys) / sizeof(keys[0]); i++) {
        if (json_number_at(payload, keys[i])) {
            return true;
        }
    }

    child = json_object_child(payload, "device");
    if (child != NULL) {
        for (i = 0; i < sizeof(keys) / sizeof(keys[0]); i++) {
            if (json_number_at(child, keys[i])) {
                return true;
            }
        }
    }

    child = json_object_child(payload, "tilt");
    if (child != NULL) {
        for (i = 0; i < sizeof(keys) / sizeof(keys[0]); i++) {
            if (json_number_at(child, keys[i])) {
                return true;
            }
        }
    }

    return false;
}

static JsonObj *build_device_data_payload(int field)
{
    JsonObj *payload = NULL;

    payload = json_object();
    if (payload == NULL) {
        return NULL;
    }

    json_object_set_new(payload, "field", json_integer(field));

    return payload;
}

bool aisgd_ops_refresh_status(AisgdContext *ctx) {

    JsonObj *payload = NULL;

    if (ctx == NULL) {
        return false;
    }

    if (!call_controller(ctx, CTRL_GET_STATUS, empty_payload(), &payload)) {
        return false;
    }

    status_update_from_controller(ctx->status, payload);
    json_decref(payload);

    status_set_ready_if_idle(ctx->status, "ready");

    return true;
}

bool aisgd_ops_reconcile(AisgdContext *ctx, JsonObj **response) {

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!controller_get_status(ctx)) {
        return false;
    }

    if (!reconcile_scan(ctx)) {
        return false;
    }

    if (!reconcile_alarm_subscribe(ctx)) {
        return false;
    }

    *response = build_reconcile_response();
    if (*response == NULL) {
        return false;
    }

    status_set(ctx->status, AisgdStateReady, "ready");

    return true;
}

bool aisgd_ops_scan(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    status_set_operation(ctx->status, OP_SCAN, "op-scan-001");
    ok = call_controller(ctx, CTRL_SCAN, empty_payload(), response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_get_device(AisgdContext *ctx, JsonObj **response)
{
    if (ctx == NULL || response == NULL) {
        return false;
    }

    *response = build_device_response(ctx);

    return *response != NULL;
}

bool aisgd_ops_get_info(AisgdContext *ctx, JsonObj **response) {

    return call_controller(ctx,
                           CTRL_GET_INFO,
                           empty_payload(),
                           response);
}

bool aisgd_ops_get_alarms(AisgdContext *ctx, JsonObj **response) {

    return call_controller(ctx,
                           CTRL_GET_ALARMS,
                           empty_payload(),
                           response);
}

bool aisgd_ops_clear_alarms(AisgdContext *ctx, JsonObj **response) {

    return call_controller(ctx,
                           CTRL_CLEAR_ALARMS,
                           empty_payload(),
                           response);
}

bool aisgd_ops_subscribe_alarms(AisgdContext *ctx, JsonObj **response) {

    return call_controller(ctx,
                           CTRL_SUBSCRIBE_ALARMS,
                           empty_payload(),
                           response);
}

bool aisgd_ops_self_test(AisgdContext *ctx, JsonObj **response) {

    bool ok;

    status_set_operation(ctx->status, OP_SELF_TEST, "op-selftest-001");
    ok = call_controller(ctx,
                         CTRL_SELF_TEST,
                         empty_payload(),
                         response);

    if (!ok) {
        status_clear_operation(ctx->status);
    }

    return ok;
}

bool aisgd_ops_configure(AisgdContext *ctx,
                         const char *profile,
                         const char *configPath,
                         JsonObj **response) {

    JsonObj *payload = NULL;
    bool ok;

    payload = build_config_payload(profile, configPath);
    if (payload == NULL) {
        return false;
    }

    status_set(ctx->status, AisgdStateVerifyConfig, "verify-config");
    status_set_operation(ctx->status, OP_CONFIG, "op-config-001");

    ok = call_controller(ctx, CTRL_CONFIGURE, payload, response);

    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_calibrate(AisgdContext *ctx, JsonObj **response) {

    bool ok;

    if (ctx->config->requireConfigBeforeCalibrate && !ctx->status->configured) {
        status_set(ctx->status,
                   AisgdStateDegraded,
                   "device not configured");
        return false;
    }

    status_set_operation(ctx->status, OP_CALIBRATE, "op-cal-001");
    ok = call_controller(ctx, CTRL_CALIBRATE, empty_payload(), response);

    if (!ok) {
        status_clear_operation(ctx->status);
    }

    return ok;
}

bool aisgd_ops_get_tilt(AisgdContext *ctx, JsonObj **response) {

    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    status_set_operation(ctx->status, OP_GET_TILT, "op-get-tilt-001");

    ok = call_controller(ctx,
                         CTRL_GET_TILT,
                         empty_payload(),
                         response);

    if (ok) {
        status_update_tilt_from_controller(ctx->status, *response);
    }

    status_clear_operation(ctx->status);
    if (ok) {
        status_set_ready_if_idle(ctx->status, "ready");
    }

    return ok;
}

bool aisgd_ops_set_tilt(AisgdContext *ctx,
                        double targetTiltDeg,
                        JsonObj **response) {

    JsonObj *payload = NULL;
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (ctx->config->requireCalibrateBeforeSetTilt &&
        !ctx->status->calibrated) {
        status_set(ctx->status,
                   AisgdStateDegraded,
                   "device not calibrated");
        return false;
    }

    payload = build_tilt_payload(targetTiltDeg);
    if (payload == NULL) {
        return false;
    }

    status_set_operation(ctx->status, OP_SET_TILT, "op-set-tilt-001");

    ok = call_controller(ctx, CTRL_SET_TILT, payload, response);

    if (ok) {
        status_set_target_tilt(ctx->status, targetTiltDeg);
        status_update_tilt_from_controller(ctx->status, *response);

        if (!payload_has_current_tilt(*response)) {
            /*
             * TS 25.463 SetTilt is a blocking Class 1 operation. If the
             * controller returned OK but did not echo the current tilt, keep
             * aisgd status useful by reflecting the requested final position.
             */
            status_set_tilt(ctx->status, targetTiltDeg);
        }
    }

    status_clear_operation(ctx->status);
    if (ok) {
        status_set_ready_if_idle(ctx->status, "ready");
    }

    return ok;
}

bool aisgd_ops_get_device_data(AisgdContext *ctx,
                               int field,
                               JsonObj **response) {

    JsonObj *payload = NULL;

    payload = build_device_data_payload(field);
    if (payload == NULL) {
        return false;
    }

    return call_controller(ctx,
                           CTRL_GET_DEVICE_DATA,
                           payload,
                           response);
}

bool aisgd_ops_reset(AisgdContext *ctx, JsonObj **response) {

    return call_controller(ctx,
                           CTRL_RESET_SOFTWARE,
                           empty_payload(),
                           response);
}
