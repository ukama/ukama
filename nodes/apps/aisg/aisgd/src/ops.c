/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

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
#define OP_GET_INFO           "get-info"
#define OP_GET_ALARMS         "get-error-status"
#define OP_CLEAR_ALARMS       "clear-active-alarms"
#define OP_SUBSCRIBE_ALARMS   "alarm-subscribe"
#define OP_CONFIG             "configure"
#define OP_CALIBRATE          "calibrate"
#define OP_GET_TILT           "get-tilt"
#define OP_SET_TILT           "set-tilt"
#define OP_SELF_TEST          "self-test"
#define OP_GET_DEVICE_DATA    "get-device-data"
#define OP_RESET              "reset"

static JsonObj *empty_payload(void)
{
    return json_object();
}

static bool call_controller(AisgdContext *ctx,
                            const char *type,
                            JsonObj *payload,
                            JsonObj **response)
{
    CtrlResponse ctrlResp;

    if (ctx == NULL || type == NULL || response == NULL) {
        json_decref(payload);
        return false;
    }

    if (!controller_client_call(&ctx->controller, type, payload, &ctrlResp)) {
        status_mark_controller_down(ctx->status, "controller unavailable");
        return false;
    }

    status_mark_controller_up(ctx->status, "controller connected");

    if (!ctrlResp.ok) {
        status_set(ctx->status, AisgdStateDegraded, ctrlResp.reason);
        ctrl_response_free(&ctrlResp);
        return false;
    }

    *response = ctrl_response_steal_payload(&ctrlResp);
    ctrl_response_free(&ctrlResp);

    return true;
}

static bool status_has(AisgdContext *ctx,
                       bool (*predicate)(const AppStatusSnapshot *snapshot),
                       AisgdState state,
                       const char *reason)
{
    AppStatusSnapshot snapshot;

    if (ctx == NULL || ctx->status == NULL || predicate == NULL) {
        return false;
    }

    if (!status_snapshot(ctx->status, &snapshot)) {
        return false;
    }

    if (snapshot.operationActive) {
        status_set(ctx->status, AisgdStateOperationRunning, "operation already running");
        return false;
    }

    if (predicate(&snapshot)) {
        return true;
    }

    status_set(ctx->status, state, reason);
    return false;
}

static bool has_device(const AppStatusSnapshot *snapshot)
{
    return snapshot != NULL && snapshot->controllerConnected && snapshot->devicePresent;
}

static bool has_identified_device(const AppStatusSnapshot *snapshot)
{
    return has_device(snapshot) && snapshot->identified;
}

static bool has_configured_device(const AppStatusSnapshot *snapshot)
{
    return has_identified_device(snapshot) && snapshot->configured;
}

static bool has_calibrated_device(const AppStatusSnapshot *snapshot)
{
    return has_configured_device(snapshot) && snapshot->calibrated;
}

static bool require_device(AisgdContext *ctx)
{
    return status_has(ctx,
                      has_device,
                      AisgdStateScanDevice,
                      "device not connected; run scan or reconcile");
}

static bool require_identified(AisgdContext *ctx)
{
    return status_has(ctx,
                      has_identified_device,
                      AisgdStateConnected,
                      "device not identified; run get-info or reconcile");
}

static bool require_configured(AisgdContext *ctx)
{
    return status_has(ctx,
                      has_configured_device,
                      AisgdStateIdentified,
                      "device not configured");
}

static bool require_calibrated(AisgdContext *ctx)
{
    return status_has(ctx,
                      has_calibrated_device,
                      AisgdStateConfigured,
                      "device not calibrated");
}

static bool controller_get_status(AisgdContext *ctx)
{
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

static JsonObj *build_reconcile_response(AisgdContext *ctx, JsonObj *actions)
{
    JsonObj *json = NULL;
    JsonObj *statusJson = NULL;
    AppStatusSnapshot snapshot;

    if (ctx == NULL || actions == NULL) {
        json_decref(actions);
        return NULL;
    }

    if (!status_snapshot(ctx->status, &snapshot)) {
        json_decref(actions);
        return NULL;
    }

    statusJson = status_to_json(ctx->status);
    if (statusJson == NULL) {
        json_decref(actions);
        return NULL;
    }

    json = json_object();
    if (json == NULL) {
        json_decref(actions);
        json_decref(statusJson);
        return NULL;
    }

    json_object_set_new(json, "state", json_string(status_state_name(snapshot.state)));
    json_object_set_new(json, "ready", json_boolean(snapshot.ready));
    json_object_set_new(json, "reason", json_string(snapshot.reason));
    json_object_set_new(json, "actions", actions);
    json_object_set_new(json, "status", statusJson);

    return json;
}

static bool add_action(JsonObj *actions, const char *action)
{
    if (actions == NULL || action == NULL) {
        return false;
    }

    return json_array_append_new(actions, json_string(action)) == 0;
}

static JsonObj *build_device_response(AisgdContext *ctx)
{
    JsonObj *json = NULL;
    AppStatusSnapshot snapshot;

    if (ctx == NULL || !status_snapshot(ctx->status, &snapshot)) {
        return NULL;
    }

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "deviceId", json_string(AISGD_DEVICE_ID));
    json_object_set_new(json, "present", json_boolean(snapshot.devicePresent));
    json_object_set_new(json, "identified", json_boolean(snapshot.identified));
    json_object_set_new(json, "configured", json_boolean(snapshot.configured));
    json_object_set_new(json, "calibrated", json_boolean(snapshot.calibrated));
    json_object_set_new(json, "model", json_string(snapshot.model));
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

bool aisgd_ops_refresh_status(AisgdContext *ctx)
{
    JsonObj *payload = NULL;

    if (ctx == NULL) {
        return false;
    }

    if (!call_controller(ctx, CTRL_GET_STATUS, empty_payload(), &payload)) {
        return false;
    }

    status_update_from_controller(ctx->status, payload);
    json_decref(payload);
    status_recompute_if_idle(ctx->status, "status refreshed");

    return true;
}

bool aisgd_ops_reconcile(AisgdContext *ctx, JsonObj **response)
{
    JsonObj *actions = NULL;
    JsonObj *payload = NULL;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    actions = json_array();
    if (actions == NULL) {
        return false;
    }

    if (!controller_get_status(ctx)) {
        json_decref(actions);
        return false;
    }
    add_action(actions, "controller-ok");

    if (!aisgd_ops_scan(ctx, &payload)) {
        json_decref(actions);
        return false;
    }
    json_decref(payload);
    payload = NULL;
    add_action(actions, "scan-complete");

    if (!aisgd_ops_get_info(ctx, &payload)) {
        json_decref(actions);
        return false;
    }
    json_decref(payload);
    payload = NULL;
    add_action(actions, "device-identified");

    if (!aisgd_ops_subscribe_alarms(ctx, &payload)) {
        json_decref(actions);
        return false;
    }
    json_decref(payload);
    payload = NULL;
    add_action(actions, "alarm-subscribed");

    status_recompute_if_idle(ctx->status, "reconcile complete");

    *response = build_reconcile_response(ctx, actions);
    return *response != NULL;
}

bool aisgd_ops_scan(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    status_set(ctx->status, AisgdStateScanDevice, "scan-device");
    status_set_operation(ctx->status, OP_SCAN, "op-scan-001");

    ok = call_controller(ctx, CTRL_SCAN, empty_payload(), response);
    if (ok) {
        status_update_from_controller(ctx->status, *response);
    }

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

bool aisgd_ops_get_info(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!require_device(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_GET_INFO, "op-get-info-001");
    ok = call_controller(ctx, CTRL_GET_INFO, empty_payload(), response);
    if (ok) {
        status_mark_identified(ctx->status, *response);
    }
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_get_alarms(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (!require_identified(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_GET_ALARMS, "op-get-alarms-001");
    ok = call_controller(ctx, CTRL_GET_ALARMS, empty_payload(), response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_clear_alarms(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (!require_identified(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_CLEAR_ALARMS, "op-clear-alarms-001");
    ok = call_controller(ctx, CTRL_CLEAR_ALARMS, empty_payload(), response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_subscribe_alarms(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (!require_identified(ctx)) {
        return false;
    }

    status_set(ctx->status, AisgdStateSubscribeAlarms, "subscribe-alarms");
    status_set_operation(ctx->status, OP_SUBSCRIBE_ALARMS, "op-alarm-sub-001");
    ok = call_controller(ctx, CTRL_SUBSCRIBE_ALARMS, empty_payload(), response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_self_test(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (!require_identified(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_SELF_TEST, "op-selftest-001");
    ok = call_controller(ctx,
                         CTRL_SELF_TEST,
                         empty_payload(),
                         response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_configure(AisgdContext *ctx,
                         const char *profile,
                         const char *configPath,
                         JsonObj **response)
{
    JsonObj *payload = NULL;
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!require_identified(ctx)) {
        return false;
    }

    payload = build_config_payload(profile, configPath);
    if (payload == NULL) {
        return false;
    }

    status_set(ctx->status, AisgdStateVerifyConfig, "verify-config");
    status_set_operation(ctx->status, OP_CONFIG, "op-config-001");

    ok = call_controller(ctx, CTRL_CONFIGURE, payload, response);
    if (ok) {
        status_mark_configured(ctx->status, *response);
    }

    status_clear_operation(ctx->status);
    return ok;
}

bool aisgd_ops_calibrate(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (ctx->config->requireConfigBeforeCalibrate && !require_configured(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_CALIBRATE, "op-cal-001");
    ok = call_controller(ctx, CTRL_CALIBRATE, empty_payload(), response);
    if (ok) {
        status_mark_calibrated(ctx->status, *response);
    }
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_get_tilt(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!require_device(ctx)) {
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
    return ok;
}

bool aisgd_ops_set_tilt(AisgdContext *ctx,
                        double targetTiltDeg,
                        JsonObj **response)
{
    JsonObj *payload = NULL;
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!require_device(ctx)) {
        return false;
    }

    if (ctx->config->requireCalibrateBeforeSetTilt && !require_calibrated(ctx)) {
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
            status_set_tilt(ctx->status, targetTiltDeg);
        }
    }

    status_clear_operation(ctx->status);
    return ok;
}

bool aisgd_ops_get_device_data(AisgdContext *ctx,
                               int field,
                               JsonObj **response)
{
    JsonObj *payload = NULL;
    bool ok;

    if (!require_identified(ctx)) {
        return false;
    }

    payload = build_device_data_payload(field);
    if (payload == NULL) {
        return false;
    }

    status_set_operation(ctx->status, OP_GET_DEVICE_DATA, "op-get-data-001");
    ok = call_controller(ctx,
                         CTRL_GET_DEVICE_DATA,
                         payload,
                         response);
    status_clear_operation(ctx->status);

    return ok;
}

bool aisgd_ops_reset(AisgdContext *ctx, JsonObj **response)
{
    bool ok;

    if (ctx == NULL || response == NULL) {
        return false;
    }

    if (!require_device(ctx)) {
        return false;
    }

    status_set_operation(ctx->status, OP_RESET, "op-reset-001");
    ok = call_controller(ctx,
                         CTRL_RESET_SOFTWARE,
                         empty_payload(),
                         response);
    if (ok) {
        status_mark_reset(ctx->status, "device reset; scan required");
    }
    status_clear_operation(ctx->status);

    return ok;
}
