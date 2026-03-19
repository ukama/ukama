/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <time.h>

#include "alarms.h"
#include "notify.h"

static void set_alarm(EmuAlarm *alarm, int active, int code,
                      const char *source, const char *severity,
                      const char *message) {
    alarm->active = active;
    alarm->code   = code;
    alarm->raisedAt = active ? time(NULL) : 0;
    snprintf(alarm->source, sizeof(alarm->source), "%s", source);
    snprintf(alarm->severity, sizeof(alarm->severity), "%s", severity);
    snprintf(alarm->message, sizeof(alarm->message), "%s", message);
}

void alarms_refresh(EmuModel *model) {
    if (model->alarmCount < 2) {
        model->alarmCount = 2;
    }

    set_alarm(&model->alarms[0], model->info.alarmLinkFailure,
              1001, SERVICE_NAME, "ERROR", "Link failure detected");
    set_alarm(&model->alarms[1], model->info.alarmPoeFailure,
              1002, SERVICE_NAME, "ERROR", "PoE failure detected");

    if (model->alarms[0].active) {
        notify_send_alarm(&model->cfg, &model->alarms[0], &model->info);
    }
    if (model->alarms[1].active) {
        notify_send_alarm(&model->cfg, &model->alarms[1], &model->info);
    }
}
