/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>

#include "status.h"

static void copy_str(char *dst, size_t size, const char *src)
{
    if (dst == NULL || size == 0) {
        return;
    }

    snprintf(dst, size, "%s", src ? src : "");
}

static JsonObj *json_child(JsonObj *json, const char *key)
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

static bool json_bool_at(JsonObj *json, const char *key, bool *out)
{
    JsonObj *value = NULL;

    if (json == NULL || key == NULL || out == NULL) {
        return false;
    }

    value = json_object_get(json, key);
    if (!json_is_boolean(value)) {
        return false;
    }

    *out = json_is_true(value);
    return true;
}

static bool json_number_at(JsonObj *json, const char *key, double *out)
{
    JsonObj *value = NULL;

    if (json == NULL || key == NULL || out == NULL) {
        return false;
    }

    value = json_object_get(json, key);
    if (!json_is_number(value)) {
        return false;
    }

    *out = json_number_value(value);
    return true;
}

static bool json_find_number(JsonObj *json,
                             const char **keys,
                             size_t keyCount,
                             double *out)
{
    JsonObj *child = NULL;
    size_t i;

    if (json == NULL || keys == NULL || out == NULL) {
        return false;
    }

    for (i = 0; i < keyCount; i++) {
        if (json_number_at(json, keys[i], out)) {
            return true;
        }
    }

    child = json_child(json, "device");
    if (child != NULL) {
        for (i = 0; i < keyCount; i++) {
            if (json_number_at(child, keys[i], out)) {
                return true;
            }
        }
    }

    child = json_child(json, "tilt");
    if (child != NULL) {
        for (i = 0; i < keyCount; i++) {
            if (json_number_at(child, keys[i], out)) {
                return true;
            }
        }
    }

    return false;
}

static const char *json_string_value_or(JsonObj *json, const char *key)
{
    JsonObj *value = NULL;

    if (json == NULL || key == NULL) {
        return NULL;
    }

    value = json_object_get(json, key);
    if (!json_is_string(value)) {
        return NULL;
    }

    return json_string_value(value);
}

const char *status_state_name(AisgdState state)
{
    switch (state) {
    case AisgdStateStarting:          return "starting";
    case AisgdStateDisconnected:      return "disconnected";
    case AisgdStateConnectController: return "connect-controller";
    case AisgdStateScanDevice:        return "scan-device";
    case AisgdStateConnected:         return "connected";
    case AisgdStateIdentified:        return "identified";
    case AisgdStateSubscribeAlarms:   return "subscribe-alarms";
    case AisgdStateVerifyConfig:      return "verify-config";
    case AisgdStateConfigured:        return "configured";
    case AisgdStateCalibrated:        return "calibrated";
    case AisgdStateReady:             return "ready";
    case AisgdStateOperationRunning:  return "operation-running";
    case AisgdStateDegraded:          return "degraded";
    case AisgdStateFailed:            return "failed";
    default:                          return "unknown";
    }
}

static void recompute_locked(AppStatus *status, const char *reason)
{
    if (status == NULL || status->operationActive) {
        return;
    }

    status->ready = false;

    if (!status->controllerConnected) {
        status->state = AisgdStateDisconnected;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "controller unavailable");
        return;
    }

    if (!status->devicePresent) {
        status->state = AisgdStateScanDevice;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "device scan required");
        return;
    }

    if (!status->identified) {
        status->state = AisgdStateConnected;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "device connected; identification required");
        return;
    }

    if (!status->configured) {
        status->state = AisgdStateIdentified;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "device identified; configuration required");
        return;
    }

    if (!status->calibrated) {
        status->state = AisgdStateConfigured;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "device configured; calibration required");
        return;
    }

    status->state = AisgdStateReady;
    status->ready = true;
    copy_str(status->reason,
             sizeof(status->reason),
             reason ? reason : "ready");
}

void status_init(AppStatus *status)
{
    if (status == NULL) {
        return;
    }

    memset(status, 0, sizeof(AppStatus));
    pthread_mutex_init(&status->mutex, NULL);

    status->state = AisgdStateStarting;
    copy_str(status->reason,  sizeof(status->reason),  "starting");
    copy_str(status->backend, sizeof(status->backend), "unknown");
    copy_str(status->mode,    sizeof(status->mode),    "unknown");
}

void status_destroy(AppStatus *status)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_destroy(&status->mutex);
}

void status_set(AppStatus *status, AisgdState state, const char *reason)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->state = state;
    status->ready = (state == AisgdStateReady);
    copy_str(status->reason, sizeof(status->reason), reason);
    pthread_mutex_unlock(&status->mutex);
}

void status_set_operation(AppStatus *status, const char *type, const char *id)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->operationActive = true;
    copy_str(status->operationType, sizeof(status->operationType), type);
    copy_str(status->operationId,   sizeof(status->operationId),   id);
    status->state = AisgdStateOperationRunning;
    status->ready = false;
    pthread_mutex_unlock(&status->mutex);
}

void status_clear_operation(AppStatus *status)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->operationActive  = false;
    status->operationType[0] = '\0';
    status->operationId[0]   = '\0';
    recompute_locked(status, NULL);
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_controller_up(AppStatus *status, const char *reason)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->controllerConnected = true;
    recompute_locked(status, reason);
    pthread_mutex_unlock(&status->mutex);
}

static void copy_identity_locked(AppStatus *status, JsonObj *payload)
{
    const char *s = NULL;

    if (status == NULL || payload == NULL) {
        return;
    }

    s = json_string_value_or(payload, "productNumber");
    if (s != NULL) {
        copy_str(status->productNumber, sizeof(status->productNumber), s);
        if (status->model[0] == '\0') {
            copy_str(status->model, sizeof(status->model), s);
        }
    }

    s = json_string_value_or(payload, "serialNumber");
    if (s != NULL) {
        copy_str(status->serialNumber, sizeof(status->serialNumber), s);
    }

    s = json_string_value_or(payload, "hardwareVersion");
    if (s != NULL) {
        copy_str(status->hardwareVersion, sizeof(status->hardwareVersion), s);
    }

    s = json_string_value_or(payload, "softwareVersion");
    if (s != NULL) {
        copy_str(status->softwareVersion, sizeof(status->softwareVersion), s);
    }

    s = json_string_value_or(payload, "model");
    if (s != NULL) {
        copy_str(status->model, sizeof(status->model), s);
    }
}

void status_update_from_controller(AppStatus *status, JsonObj *payload)
{
    JsonObj *value = NULL;
    bool b;

    if (status == NULL || payload == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);

    value = json_object_get(payload, "transport");
    if (json_is_string(value)) {
        copy_str(status->backend,
                 sizeof(status->backend),
                 json_string_value(value));
    }

    value = json_object_get(payload, "mode");
    if (json_is_string(value)) {
        copy_str(status->mode,
                 sizeof(status->mode),
                 json_string_value(value));
    }

    if (json_bool_at(payload, "powerManaged", &b)) {
        status->powerManaged = b;
    }

    if (json_bool_at(payload, "present", &b)) {
        status->devicePresent = b;
        if (!b) {
            status->identified = false;
            status->configured = false;
            status->calibrated = false;
            status->tiltKnown = false;
            status->targetTiltKnown = false;
            status->model[0] = '\0';
            status->productNumber[0] = '\0';
            status->serialNumber[0] = '\0';
            status->hardwareVersion[0] = '\0';
            status->softwareVersion[0] = '\0';
        }
    }

    if (json_bool_at(payload, "configured", &b)) {
        status->configured = b;
        if (!b) {
            status->calibrated = false;
        }
    }

    if (json_bool_at(payload, "calibrated", &b)) {
        status->calibrated = b;
    }

    if (json_bool_at(payload, "busy", &b)) {
        status->busy = b;
    }

    copy_identity_locked(status, payload);
    if (status->productNumber[0] != '\0' ||
        status->serialNumber[0] != '\0' ||
        status->model[0] != '\0') {
        status->identified = true;
    }

    {
        static const char *currentKeys[] = {
            "currentTiltDeg",
            "tiltDeg",
            "electricalTiltDeg",
            "tilt"
        };
        static const char *targetKeys[] = {
            "targetTiltDeg",
            "requestedTiltDeg"
        };
        double tilt = 0.0;

        if (json_find_number(payload,
                             currentKeys,
                             sizeof(currentKeys) / sizeof(currentKeys[0]),
                             &tilt)) {
            status->currentTiltDeg = tilt;
            status->tiltKnown = true;
        }

        if (json_find_number(payload,
                             targetKeys,
                             sizeof(targetKeys) / sizeof(targetKeys[0]),
                             &tilt)) {
            status->targetTiltDeg = tilt;
            status->targetTiltKnown = true;
        }
    }

    status->controllerConnected = true;
    recompute_locked(status, NULL);

    pthread_mutex_unlock(&status->mutex);
}

void status_update_tilt_from_controller(AppStatus *status, JsonObj *payload)
{
    if (status == NULL || payload == NULL) {
        return;
    }

    status_update_from_controller(status, payload);
}

void status_mark_identified(AppStatus *status, JsonObj *payload)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->controllerConnected = true;
    status->devicePresent = true;
    status->identified = true;
    copy_identity_locked(status, payload);
    recompute_locked(status, "device identified");
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_configured(AppStatus *status, JsonObj *payload)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->controllerConnected = true;
    status->devicePresent = true;
    status->identified = true;
    status->configured = true;
    status->calibrated = false;
    status->tiltKnown = false;
    status->targetTiltKnown = false;
    copy_identity_locked(status, payload);
    recompute_locked(status, "device configured; calibration required");
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_calibrated(AppStatus *status, JsonObj *payload)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->controllerConnected = true;
    status->devicePresent = true;
    status->identified = true;
    status->configured = true;
    status->calibrated = true;
    copy_identity_locked(status, payload);
    recompute_locked(status, "ready");
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_reset(AppStatus *status, const char *reason)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);

    status->devicePresent = false;
    status->identified = false;
    status->configured = false;
    status->calibrated = false;
    status->busy = false;
    status->tiltKnown = false;
    status->targetTiltKnown = false;
    status->currentTiltDeg = 0.0;
    status->targetTiltDeg = 0.0;
    status->model[0] = '\0';
    status->productNumber[0] = '\0';
    status->serialNumber[0] = '\0';
    status->hardwareVersion[0] = '\0';
    status->softwareVersion[0] = '\0';

    recompute_locked(status, reason ? reason : "device reset; scan required");

    pthread_mutex_unlock(&status->mutex);
}

void status_set_tilt(AppStatus *status, double currentTiltDeg)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->currentTiltDeg = currentTiltDeg;
    status->tiltKnown = true;
    pthread_mutex_unlock(&status->mutex);
}

void status_set_target_tilt(AppStatus *status, double targetTiltDeg)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    status->targetTiltDeg = targetTiltDeg;
    status->targetTiltKnown = true;
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_controller_down(AppStatus *status, const char *reason)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);

    status->controllerConnected = false;
    status->powerManaged = false;
    status->devicePresent = false;
    status->identified = false;
    status->configured = false;
    status->calibrated = false;
    status->busy = false;
    status->tiltKnown = false;
    status->targetTiltKnown = false;
    status->currentTiltDeg = 0.0;
    status->targetTiltDeg = 0.0;

    copy_str(status->mode,  sizeof(status->mode), "unknown");
    copy_str(status->model, sizeof(status->model), "");
    status->productNumber[0] = '\0';
    status->serialNumber[0] = '\0';
    status->hardwareVersion[0] = '\0';
    status->softwareVersion[0] = '\0';

    status->operationActive = false;
    status->operationType[0] = '\0';
    status->operationId[0] = '\0';

    recompute_locked(status,
                     reason ? reason : "controller unavailable");

    pthread_mutex_unlock(&status->mutex);
}

void status_recompute_if_idle(AppStatus *status, const char *reason)
{
    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);
    recompute_locked(status, reason);
    pthread_mutex_unlock(&status->mutex);
}

void status_set_ready_if_idle(AppStatus *status, const char *reason)
{
    status_recompute_if_idle(status, reason);
}

bool status_is_ready(AppStatus *status)
{
    bool ready;

    if (status == NULL) {
        return false;
    }

    pthread_mutex_lock(&status->mutex);
    ready = status->ready;
    pthread_mutex_unlock(&status->mutex);

    return ready;
}

bool status_snapshot(AppStatus *status, AppStatusSnapshot *snapshot)
{
    if (status == NULL || snapshot == NULL) {
        return false;
    }

    pthread_mutex_lock(&status->mutex);

    memset(snapshot, 0, sizeof(*snapshot));
    snapshot->state = status->state;
    snapshot->ready = status->ready;
    snapshot->controllerConnected = status->controllerConnected;
    snapshot->devicePresent = status->devicePresent;
    snapshot->identified = status->identified;
    snapshot->configured = status->configured;
    snapshot->calibrated = status->calibrated;
    snapshot->busy = status->busy;
    snapshot->operationActive = status->operationActive;
    snapshot->tiltKnown = status->tiltKnown;
    snapshot->targetTiltKnown = status->targetTiltKnown;
    snapshot->currentTiltDeg = status->currentTiltDeg;
    snapshot->targetTiltDeg = status->targetTiltDeg;
    copy_str(snapshot->reason, sizeof(snapshot->reason), status->reason);
    copy_str(snapshot->model, sizeof(snapshot->model), status->model);
    copy_str(snapshot->operationType,
             sizeof(snapshot->operationType),
             status->operationType);
    copy_str(snapshot->operationId,
             sizeof(snapshot->operationId),
             status->operationId);

    pthread_mutex_unlock(&status->mutex);

    return true;
}

JsonObj *status_to_json(AppStatus *status)
{
    JsonObj *root;
    JsonObj *controller;
    JsonObj *device;
    JsonObj *operation;
    JsonObj *identity;

    char reason[STATUS_REASON_LEN];
    char backend[STATUS_MAX_STR];
    char mode[STATUS_MAX_STR];
    char model[STATUS_MAX_STR];
    char productNumber[STATUS_MAX_STR];
    char serialNumber[STATUS_MAX_STR];
    char hardwareVersion[STATUS_MAX_STR];
    char softwareVersion[STATUS_MAX_STR];
    char opType[STATUS_MAX_STR];
    char opId[STATUS_MAX_STR];

    AisgdState state;

    bool ready;
    bool connected;
    bool powerManaged;
    bool present;
    bool identified;
    bool configured;
    bool calibrated;
    bool busy;
    bool tiltKnown;
    bool targetTiltKnown;
    bool opActive;
    double currentTiltDeg;
    double targetTiltDeg;

    if (status == NULL) {
        return NULL;
    }

    pthread_mutex_lock(&status->mutex);

    state        = status->state;
    ready        = status->ready;
    connected    = status->controllerConnected;
    powerManaged = status->powerManaged;
    present      = status->devicePresent;
    identified   = status->identified;
    configured   = status->configured;
    calibrated   = status->calibrated;
    busy         = status->busy;
    tiltKnown    = status->tiltKnown;
    targetTiltKnown = status->targetTiltKnown;
    currentTiltDeg  = status->currentTiltDeg;
    targetTiltDeg   = status->targetTiltDeg;
    opActive     = status->operationActive;
    copy_str(reason,  sizeof(reason),  status->reason);
    copy_str(backend, sizeof(backend), status->backend);
    copy_str(mode,    sizeof(mode),    status->mode);
    copy_str(model,   sizeof(model),   status->model);
    copy_str(productNumber,   sizeof(productNumber),   status->productNumber);
    copy_str(serialNumber,    sizeof(serialNumber),    status->serialNumber);
    copy_str(hardwareVersion, sizeof(hardwareVersion), status->hardwareVersion);
    copy_str(softwareVersion, sizeof(softwareVersion), status->softwareVersion);
    copy_str(opType,  sizeof(opType),  status->operationType);
    copy_str(opId,    sizeof(opId),    status->operationId);

    pthread_mutex_unlock(&status->mutex);

    root       = json_object();
    controller = json_object();
    device     = json_object();
    operation  = json_object();
    identity   = json_object();

    if (root == NULL || controller == NULL || device == NULL ||
        operation == NULL || identity == NULL) {
        json_decref(root);
        json_decref(controller);
        json_decref(device);
        json_decref(operation);
        json_decref(identity);
        return NULL;
    }

    json_object_set_new(root, "state", json_string(status_state_name(state)));
    json_object_set_new(root, "ready", json_boolean(ready));
    json_object_set_new(root, "reason", json_string(reason));

    json_object_set_new(controller, "connected", json_boolean(connected));
    json_object_set_new(controller, "backend", json_string(backend));
    json_object_set_new(controller, "powerManaged", json_boolean(powerManaged));
    json_object_set_new(controller, "mode", json_string(mode));

    json_object_set_new(identity, "model", json_string(model));
    json_object_set_new(identity, "productNumber", json_string(productNumber));
    json_object_set_new(identity, "serialNumber", json_string(serialNumber));
    json_object_set_new(identity, "hardwareVersion", json_string(hardwareVersion));
    json_object_set_new(identity, "softwareVersion", json_string(softwareVersion));

    json_object_set_new(device, "present", json_boolean(present));
    json_object_set_new(device, "identified", json_boolean(identified));
    json_object_set_new(device, "model", json_string(model));
    json_object_set_new(device, "identity", identity);
    json_object_set_new(device, "configured", json_boolean(configured));
    json_object_set_new(device, "calibrated", json_boolean(calibrated));
    json_object_set_new(device, "busy", json_boolean(busy));
    json_object_set_new(device, "tiltKnown", json_boolean(tiltKnown));
    if (tiltKnown) {
        json_object_set_new(device, "currentTiltDeg", json_real(currentTiltDeg));
    } else {
        json_object_set_new(device, "currentTiltDeg", json_null());
    }
    json_object_set_new(device, "targetTiltKnown", json_boolean(targetTiltKnown));
    if (targetTiltKnown) {
        json_object_set_new(device, "targetTiltDeg", json_real(targetTiltDeg));
    } else {
        json_object_set_new(device, "targetTiltDeg", json_null());
    }

    json_object_set_new(operation, "active", json_boolean(opActive));
    json_object_set_new(operation, "type", json_string(opType));
    json_object_set_new(operation, "id", json_string(opId));

    json_object_set_new(root, "controller", controller);
    json_object_set_new(root, "device", device);
    json_object_set_new(root, "operation", operation);

    return root;
}
