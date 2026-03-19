/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <time.h>

#include "alarms.h"
#include "web_client.h"

#include "usys_log.h"

static SwitchAlarm *find_alarm(SwitchdContext *ctx,
                               SwitchdAlarmCode code,
                               const char *resource) {
    size_t i;

    for (i = 0; i < ctx->alarmCount; i++) {
        if (ctx->alarms[i].code == code &&
            strcmp(ctx->alarms[i].resource, resource) == 0) {
            return &ctx->alarms[i];
        }
    }

    return NULL;
}

static void upsert_alarm(SwitchdContext *ctx,
                         SwitchdAlarmCode code,
                         SwitchdAlarmSeverity severity,
                         const char *resource,
                         const char *text,
                         bool active) {
    SwitchAlarm *alarm;
    time_t now;

    now = time(NULL);
    alarm = find_alarm(ctx, code, resource);
    if (alarm == NULL) {
        if (ctx->alarmCount >= (sizeof(ctx->alarms) / sizeof(ctx->alarms[0]))) {
            return;
        }

        alarm = &ctx->alarms[ctx->alarmCount++];
        memset(alarm, 0, sizeof(*alarm));
        alarm->code = code;
        snprintf(alarm->resource, sizeof(alarm->resource), "%s", resource);
        alarm->firstSeen = now;
    }

    alarm->severity = severity;
    snprintf(alarm->text, sizeof(alarm->text), "%s", text);
    alarm->lastSeen = now;

    if (active && !alarm->active) {
        alarm->active = true;
        alarm->latched = false;
    } else if (!active && alarm->active) {
        alarm->active = false;
        alarm->latched = false;
    }
}

static void maybe_send_alarm(SwitchdContext *ctx, SwitchAlarm *alarm) {
    if (alarm->latched) {
        return;
    }

    if (web_client_notify_alarm(&ctx->config,
                                alarm,
                                !alarm->active) == SWITCHD_OK) {
        alarm->latched = true;
        alarm->lastSent = time(NULL);
    } else {
        usys_log_error("Failed to send alarm %d for %s",
                       alarm->code,
                       alarm->resource);
    }
}

int alarms_scan(SwitchdContext *ctx) {
    uint32_t i;
    char resource[32];
    char text[128];

    pthread_mutex_lock(&ctx->alarmMutex);

    upsert_alarm(ctx,
                 SWITCHD_ALARM_SWITCH_UNREACHABLE,
                 SWITCHD_ALARM_SEV_CRITICAL,
                 "switch",
                 "Switch is unreachable",
                 !ctx->info.reachable);

    upsert_alarm(ctx,
                 SWITCHD_ALARM_SWITCH_RECOVERED,
                 SWITCHD_ALARM_SEV_INFO,
                 "switch",
                 "Switch recovered",
                 ctx->info.reachable);

    upsert_alarm(ctx,
                 SWITCHD_ALARM_HIGH_SYSTEM_TEMP,
                 SWITCHD_ALARM_SEV_WARNING,
                 "system",
                 "System temperature above threshold",
                 ctx->kpis.systemTemperatureC > 70.0);

    upsert_alarm(ctx,
                 SWITCHD_ALARM_HIGH_AMBIENT_TEMP,
                 SWITCHD_ALARM_SEV_WARNING,
                 "ambient",
                 "Ambient temperature above threshold",
                 ctx->kpis.ambientTemperatureC > 60.0);

    for (i = 0; i < ctx->portCount; i++) {
        snprintf(resource, sizeof(resource), "port%u", ctx->ports[i].id);

        snprintf(text, sizeof(text), "Port %u link is down", ctx->ports[i].id);
        upsert_alarm(ctx,
                     SWITCHD_ALARM_PORT_LINK_DOWN,
                     SWITCHD_ALARM_SEV_WARNING,
                     resource,
                     text,
                     ctx->ports[i].present &&
                     ctx->ports[i].adminUp &&
                     !ctx->ports[i].linkUp &&
                     ctx->config.strictLinkAlarms);

        snprintf(text, sizeof(text), "Port %u PoE is off", ctx->ports[i].id);
        upsert_alarm(ctx,
                     SWITCHD_ALARM_PORT_POE_OFF,
                     SWITCHD_ALARM_SEV_WARNING,
                     resource,
                     text,
                     ctx->ports[i].poeSupported &&
                     ctx->ports[i].poeEnabled &&
                     !ctx->ports[i].poeOperational);

        upsert_alarm(ctx,
                     SWITCHD_ALARM_PORT_POE_FAULT,
                     SWITCHD_ALARM_SEV_CRITICAL,
                     resource,
                     ctx->ports[i].fault[0] ? ctx->ports[i].fault : "Port PoE fault",
                     ctx->ports[i].fault[0] != '\0');
    }

    upsert_alarm(ctx,
                 SWITCHD_ALARM_FIRMWARE_FAILED,
                 SWITCHD_ALARM_SEV_CRITICAL,
                 "firmware",
                 ctx->fw.detail[0] ? ctx->fw.detail : "Firmware update failed",
                 ctx->fw.state == SWITCHD_FW_FAILED);

    upsert_alarm(ctx,
                 SWITCHD_ALARM_FIRMWARE_DONE,
                 SWITCHD_ALARM_SEV_INFO,
                 "firmware",
                 "Firmware update completed",
                 ctx->fw.state == SWITCHD_FW_DONE);

    for (i = 0; i < ctx->alarmCount; i++) {
        maybe_send_alarm(ctx, &ctx->alarms[i]);
    }

    pthread_mutex_unlock(&ctx->alarmMutex);
    return SWITCHD_OK;
}
