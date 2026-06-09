/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>

#include "model.h"

static void model_add_motor_alarm(EmuModel *model)
{
    JsonObj *alarm = NULL;

    if (model == NULL || model->alarms == NULL) {
        return;
    }

    alarm = json_object();
    if (alarm == NULL) {
        return;
    }

    json_object_set_new(alarm, "code", json_string("MotorJam"));
    json_object_set_new(alarm, "state", json_string("raised"));
    json_array_append_new(model->alarms, alarm);
}

void emu_model_init(EmuModel *model)
{
    if (model == NULL) {
        return;
    }

    memset(model, 0, sizeof(EmuModel));

    model->present         = true;
    model->configured      = false;
    model->calibrated      = false;
    model->busy            = false;
    model->alarmSubscribed = false;
    model->tiltTenthsDeg   = 30;
    model->alarms          = json_array();
}

void emu_model_free(EmuModel *model)
{
    if (model == NULL) {
        return;
    }

    json_decref(model->alarms);
    memset(model, 0, sizeof(EmuModel));
}

bool emu_model_load_scenario(EmuModel *model, const char *scenario)
{
    if (model == NULL || scenario == NULL) {
        return false;
    }

    if (!strcmp(scenario, "missing")) {
        model->present = false;
        return true;
    }

    if (!strcmp(scenario, "not_configured")) {
        model->configured = false;
        model->calibrated = false;
        return true;
    }

    if (!strcmp(scenario, "not_calibrated")) {
        model->configured = true;
        model->calibrated = false;
        return true;
    }

    if (!strcmp(scenario, "busy")) {
        model->busy = true;
        return true;
    }

    if (!strcmp(scenario, "alarm_active")) {
        model_add_motor_alarm(model);
        return true;
    }

    return true;
}

JsonObj *emu_model_status(EmuModel *model)
{
    JsonObj *json = NULL;

    if (model == NULL) {
        return NULL;
    }

    json = json_object();
    if (json == NULL) {
        return NULL;
    }

    json_object_set_new(json, "mode",
                        json_string("operating"));
    json_object_set_new(json, "busy",
                        json_boolean(model->busy));
    json_object_set_new(json, "present",
                        json_boolean(model->present));
    json_object_set_new(json, "configured",
                        json_boolean(model->configured));
    json_object_set_new(json, "calibrated",
                        json_boolean(model->calibrated));
    json_object_set_new(json, "powerManaged",
                        json_boolean(false));
    json_object_set_new(json, "transport",
                        json_string("emu"));

    return json;
}
