/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <sys/time.h>

#include "metrics_store.h"
#include "time_util.h"
#include "usys_log.h"

int metrics_store_init(MetricsStore *store) {
    if (!store) return -1;

    memset(store, 0, sizeof(*store));

    if (pthread_mutex_init(&store->lock, NULL) != 0) {
        usys_log_error("metrics_store: mutex init failed");
        return -1;
    }

    /* Initialize alarms to inactive */
    for (int i = 0; i < ALARM_MAX; i++) {
        store->snapshot.alarms[i].type = (AlarmType)i;
        store->snapshot.alarms[i].active = false;
        store->snapshot.alarms[i].severity = SEVERITY_OK;
    }

    store->snapshot.overall_severity = SEVERITY_OK;
    store->snapshot.comm_healthy = false;

    usys_log_debug("metrics_store: initialized");
    return 0;
}

void metrics_store_free(MetricsStore *store) {
    if (!store) return;
    pthread_mutex_destroy(&store->lock);
}

void metrics_store_update(MetricsStore *store, const ControllerData *data) {
    if (!store || !data) return;

    pthread_mutex_lock(&store->lock);

    /* Copy controller data */
    memcpy(&store->snapshot.data, data, sizeof(ControllerData));

    /* Update statistics */
    store->snapshot.total_samples++;

    if (data->comm_ok) {
        store->snapshot.last_successful_read_ms = time_now_ms();
        store->snapshot.consecutive_errors = 0;
        store->snapshot.comm_healthy = true;
    } else {
        store->snapshot.consecutive_errors++;
        store->snapshot.failed_samples++;
        if (store->snapshot.consecutive_errors > 5) {
            store->snapshot.comm_healthy = false;
        }
    }

    /* Calculate efficiency if we have valid power readings */
    if (data->pv_power_w > 0 && data->batt_voltage_v > 0 && data->batt_current_a > 0) {
        double batt_power = data->batt_voltage_v * data->batt_current_a;
        store->snapshot.efficiency_pct = (batt_power / data->pv_power_w) * 100.0;
        if (store->snapshot.efficiency_pct > 100.0) {
            store->snapshot.efficiency_pct = 100.0;
        }
    } else {
        store->snapshot.efficiency_pct = 0.0;
    }

    pthread_mutex_unlock(&store->lock);
}

void metrics_store_set_error(MetricsStore *store, int err_code, const char *msg) {
    if (!store) return;

    pthread_mutex_lock(&store->lock);

    store->snapshot.data.comm_ok = false;
    store->snapshot.data.comm_errors++;
    store->snapshot.consecutive_errors++;
    store->snapshot.failed_samples++;

    if (store->snapshot.consecutive_errors > 5) {
        store->snapshot.comm_healthy = false;
    }

    usys_log_warn("metrics_store: error %d: %s", err_code, msg ? msg : "unknown");

    pthread_mutex_unlock(&store->lock);
}

void metrics_store_get(MetricsStore *store, MetricsSnapshot *out) {
    if (!store || !out) return;

    pthread_mutex_lock(&store->lock);
    memcpy(out, &store->snapshot, sizeof(MetricsSnapshot));
    pthread_mutex_unlock(&store->lock);
}

void metrics_store_set_alarm(MetricsStore *store, AlarmType type,
                             Severity severity, const char *msg) {
    if (!store || type >= ALARM_MAX) return;

    pthread_mutex_lock(&store->lock);

    AlarmRecord *alarm = &store->snapshot.alarms[type];

    /* Only update if alarm is new or severity changed */
    if (!alarm->active || alarm->severity != severity) {
        alarm->type = type;
        alarm->severity = severity;
        alarm->timestamp_ms = time_now_ms();
        alarm->active = true;
        if (msg) {
            snprintf(alarm->message, sizeof(alarm->message), "%s", msg);
        }

        /* Add to history */
        int idx = store->alarm_history_head;
        memcpy(&store->alarm_history[idx], alarm, sizeof(AlarmRecord));
        store->alarm_history_head = (store->alarm_history_head + 1) % 64;
        if (store->alarm_history_count < 64) {
            store->alarm_history_count++;
        }

        usys_log_warn("metrics_store: alarm set - type=%s severity=%s msg=%s",
                      alarm_type_str(type), severity_str(severity), msg ? msg : "");
    }

    /* Recalculate overall severity */
    store->snapshot.overall_severity = SEVERITY_OK;
    store->snapshot.active_alarm_count = 0;
    for (int i = 0; i < ALARM_MAX; i++) {
        if (store->snapshot.alarms[i].active) {
            store->snapshot.active_alarm_count++;
            if (store->snapshot.alarms[i].severity > store->snapshot.overall_severity) {
                store->snapshot.overall_severity = store->snapshot.alarms[i].severity;
            }
        }
    }

    pthread_mutex_unlock(&store->lock);
}

void metrics_store_clear_alarm(MetricsStore *store, AlarmType type) {
    if (!store || type >= ALARM_MAX) return;

    pthread_mutex_lock(&store->lock);

    AlarmRecord *alarm = &store->snapshot.alarms[type];
    if (alarm->active) {
        alarm->active = false;
        alarm->severity = SEVERITY_OK;

        usys_log_info("metrics_store: alarm cleared - type=%s", alarm_type_str(type));

        /* Recalculate overall severity */
        store->snapshot.overall_severity = SEVERITY_OK;
        store->snapshot.active_alarm_count = 0;
        for (int i = 0; i < ALARM_MAX; i++) {
            if (store->snapshot.alarms[i].active) {
                store->snapshot.active_alarm_count++;
                if (store->snapshot.alarms[i].severity > store->snapshot.overall_severity) {
                    store->snapshot.overall_severity = store->snapshot.alarms[i].severity;
                }
            }
        }
    }

    pthread_mutex_unlock(&store->lock);
}

int metrics_store_get_alarm_history(MetricsStore *store, AlarmRecord *out,
                                    int max_records) {
    if (!store || !out || max_records <= 0) return 0;

    pthread_mutex_lock(&store->lock);

    int count = store->alarm_history_count;
    if (count > max_records) count = max_records;

    /* Copy from circular buffer (newest first) */
    for (int i = 0; i < count; i++) {
        int idx = (store->alarm_history_head - 1 - i + 64) % 64;
        memcpy(&out[i], &store->alarm_history[idx], sizeof(AlarmRecord));
    }

    pthread_mutex_unlock(&store->lock);
    return count;
}

const char *severity_str(Severity sev) {
    switch (sev) {
    case SEVERITY_OK:       return "ok";
    case SEVERITY_WARN:     return "warning";
    case SEVERITY_CRITICAL: return "critical";
    default:                return "unknown";
    }
}

const char *alarm_type_str(AlarmType type) {
    switch (type) {
    case ALARM_NONE:                return "none";
    case ALARM_LOW_BATTERY_VOLTAGE: return "low_battery_voltage";
    case ALARM_HIGH_TEMPERATURE:    return "high_temperature";
    case ALARM_CONTROLLER_FAULT:    return "controller_fault";
    case ALARM_COMMUNICATION_LOST:  return "communication_lost";
    case ALARM_PV_OVERVOLTAGE:      return "pv_overvoltage";
    case ALARM_OVERCURRENT:         return "overcurrent";
    default:                        return "unknown";
    }
}
