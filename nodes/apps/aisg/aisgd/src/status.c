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

static void copy_str(char *dst, size_t size, const char *src) {
    if (dst == NULL || size == 0) return;
    snprintf(dst, size, "%s", src ? src : "");
}

static const char *state_str(AisgdState state) {
    switch (state) {
    case AisgdStateStarting:          return "starting";
    case AisgdStateConnectController: return "connect-controller";
    case AisgdStateScanDevice:        return "scan-device";
    case AisgdStateSubscribeAlarms:   return "subscribe-alarms";
    case AisgdStateVerifyConfig:      return "verify-config";
    case AisgdStateReady:             return "ready";
    case AisgdStateOperationRunning:  return "operation-running";
    case AisgdStateDegraded:          return "degraded";
    case AisgdStateFailed:            return "failed";
    default:                          return "unknown";
    }
}

void status_init(AppStatus *status) {
    if (status == NULL) return;

    memset(status, 0, sizeof(AppStatus));
    pthread_mutex_init(&status->mutex, NULL);
    status->state = AisgdStateStarting;
    copy_str(status->reason,  sizeof(status->reason),  "starting");
    copy_str(status->backend, sizeof(status->backend), "unknown");
    copy_str(status->mode,    sizeof(status->mode),    "unknown");
}

void status_destroy(AppStatus *status) {
    if (status == NULL) return;
    pthread_mutex_destroy(&status->mutex);
}

void status_set(AppStatus *status, AisgdState state, const char *reason) {
    if (status == NULL) return;

    pthread_mutex_lock(&status->mutex);
    status->state = state;
    status->ready = (state == AisgdStateReady);
    copy_str(status->reason, sizeof(status->reason), reason);
    pthread_mutex_unlock(&status->mutex);
}

void status_set_operation(AppStatus *status, const char *type, const char *id) {
    if (status == NULL) return;

    pthread_mutex_lock(&status->mutex);
    status->operationActive = true;
    copy_str(status->operationType, sizeof(status->operationType), type);
    copy_str(status->operationId,   sizeof(status->operationId),   id);
    status->state = AisgdStateOperationRunning;
    status->ready = false;
    pthread_mutex_unlock(&status->mutex);
}

void status_clear_operation(AppStatus *status) {
    if (status == NULL) return;

    pthread_mutex_lock(&status->mutex);
    status->operationActive  = false;
    status->operationType[0] = '\0';
    status->operationId[0]   = '\0';
    pthread_mutex_unlock(&status->mutex);
}

void status_update_from_controller(AppStatus *status, JsonObj *payload) {
    JsonObj *value;

    if (status == NULL || payload == NULL) return;

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

    value = json_object_get(payload, "powerManaged");
    if (json_is_boolean(value))
        status->powerManaged = json_is_true(value);

    value = json_object_get(payload, "present");
    if (json_is_boolean(value))
        status->devicePresent = json_is_true(value);

    value = json_object_get(payload, "configured");
    if (json_is_boolean(value))
        status->configured = json_is_true(value);

    value = json_object_get(payload, "calibrated");
    if (json_is_boolean(value))
        status->calibrated = json_is_true(value);

    value = json_object_get(payload, "busy");
    if (json_is_boolean(value))
        status->busy = json_is_true(value);

    value = json_object_get(payload, "model");
    if (json_is_string(value)) {
        copy_str(status->model,
                 sizeof(status->model),
                 json_string_value(value));
    }

    status->controllerConnected = true;
    pthread_mutex_unlock(&status->mutex);
}

void status_mark_controller_down(AppStatus *status, const char *reason) {

    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);

    status->state = AisgdStateDegraded;
    status->ready = false;

    copy_str(status->reason,
             sizeof(status->reason),
             reason ? reason : "controller unavailable");

    status->controllerConnected = false;

    /*
     * Clear stale controller/device state. Once the controller is gone,
     * we should not report the last known device as still present.
     */
    status->powerManaged  = false;
    status->devicePresent = false;
    status->configured    = false;
    status->calibrated    = false;
    status->busy          = false;

    copy_str(status->mode,  sizeof(status->mode), "unknown");
    copy_str(status->model, sizeof(status->model), "");

    status->operationActive = false;
    status->operationType[0] = '\0';
    status->operationId[0] = '\0';

    pthread_mutex_unlock(&status->mutex);
}

void status_set_ready_if_idle(AppStatus *status, const char *reason) {

    if (status == NULL) {
        return;
    }

    pthread_mutex_lock(&status->mutex);

    if (!status->operationActive) {
        status->state = AisgdStateReady;
        status->ready = true;
        copy_str(status->reason,
                 sizeof(status->reason),
                 reason ? reason : "ready");
    }

    pthread_mutex_unlock(&status->mutex);
}

bool status_is_ready(AppStatus *status) {

    bool ready;

    if (status == NULL) return false;

    pthread_mutex_lock(&status->mutex);
    ready = status->ready;
    pthread_mutex_unlock(&status->mutex);

    return ready;
}

JsonObj *status_to_json(AppStatus *status) {
    JsonObj *root;
    JsonObj *controller;
    JsonObj *device;
    JsonObj *operation;

    char reason[STATUS_REASON_LEN];
    char backend[STATUS_MAX_STR];
    char mode[STATUS_MAX_STR];
    char model[STATUS_MAX_STR];
    char opType[STATUS_MAX_STR];
    char opId[STATUS_MAX_STR];

    AisgdState state;

    bool ready;
    bool connected;
    bool powerManaged;
    bool present;
    bool configured;
    bool calibrated;
    bool busy;
    bool opActive;

    if (status == NULL) return NULL;

    pthread_mutex_lock(&status->mutex);

    state        = status->state;
    ready        = status->ready;
    connected    = status->controllerConnected;
    powerManaged = status->powerManaged;
    present      = status->devicePresent;
    configured   = status->configured;
    calibrated   = status->calibrated;
    busy         = status->busy;
    opActive     = status->operationActive;
    copy_str(reason,  sizeof(reason),  status->reason);
    copy_str(backend, sizeof(backend), status->backend);
    copy_str(mode,    sizeof(mode),    status->mode);
    copy_str(model,   sizeof(model),   status->model);
    copy_str(opType,  sizeof(opType),  status->operationType);
    copy_str(opId,    sizeof(opId),    status->operationId);

    pthread_mutex_unlock(&status->mutex);

    root       = json_object();
    controller = json_object();
    device     = json_object();
    operation  = json_object();

    json_object_set_new(root, "state", json_string(state_str(state)));
    json_object_set_new(root, "ready", json_boolean(ready));
    json_object_set_new(root, "reason", json_string(reason));

    json_object_set_new(controller, "connected", json_boolean(connected));
    json_object_set_new(controller, "backend", json_string(backend));
    json_object_set_new(controller, "powerManaged", json_boolean(powerManaged));
    json_object_set_new(controller, "mode", json_string(mode));

    json_object_set_new(device, "present", json_boolean(present));
    json_object_set_new(device, "model", json_string(model));
    json_object_set_new(device, "configured", json_boolean(configured));
    json_object_set_new(device, "calibrated", json_boolean(calibrated));
    json_object_set_new(device, "busy", json_boolean(busy));

    json_object_set_new(operation, "active", json_boolean(opActive));
    json_object_set_new(operation, "type", json_string(opType));
    json_object_set_new(operation, "id", json_string(opId));

    json_object_set_new(root, "controller", controller);
    json_object_set_new(root, "device", device);
    json_object_set_new(root, "operation", operation);

    return root;
}
