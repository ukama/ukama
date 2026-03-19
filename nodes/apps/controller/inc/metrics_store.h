/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef METRICS_STORE_H
#define METRICS_STORE_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "driver.h"

/*
 * Alarm severity levels
 */
typedef enum {
    SEVERITY_OK = 0,
    SEVERITY_WARN,
    SEVERITY_CRITICAL
} Severity;

/*
 * Alarm types
 */
typedef enum {
    ALARM_NONE = 0,
    ALARM_LOW_BATTERY_VOLTAGE,
    ALARM_HIGH_TEMPERATURE,
    ALARM_CONTROLLER_FAULT,
    ALARM_COMMUNICATION_LOST,
    ALARM_PV_OVERVOLTAGE,
    ALARM_OVERCURRENT,
    ALARM_MAX
} AlarmType;

/*
 * Single alarm record
 */
typedef struct {
    AlarmType   type;
    Severity    severity;
    uint64_t    timestamp_ms;
    char        message[128];
    bool        active;
} AlarmRecord;

/*
 * Metrics snapshot - thread-safe copy of current state
 */
typedef struct {
    /* Controller data */
    ControllerData  data;

    /* Derived metrics */
    double          efficiency_pct;     /* (batt_power / pv_power) * 100 */

    /* Alarm state */
    Severity        overall_severity;
    AlarmRecord     alarms[ALARM_MAX];
    int             active_alarm_count;

    /* Communication health */
    uint64_t        last_successful_read_ms;
    uint32_t        consecutive_errors;
    bool            comm_healthy;

    /* Statistics */
    uint64_t        total_samples;
    uint64_t        failed_samples;
} MetricsSnapshot;

/*
 * Metrics store - thread-safe storage
 */
typedef struct {
    pthread_mutex_t lock;
    MetricsSnapshot snapshot;

    /* Alarm history (circular buffer) */
    AlarmRecord     alarm_history[64];
    int             alarm_history_head;
    int             alarm_history_count;
} MetricsStore;

/* Lifecycle */
int  metrics_store_init(MetricsStore *store);
void metrics_store_free(MetricsStore *store);

/* Update with new controller data */
void metrics_store_update(MetricsStore *store, const ControllerData *data);

/* Update with error (failed read) */
void metrics_store_set_error(MetricsStore *store, int err_code, const char *msg);

/* Get thread-safe copy of current state */
void metrics_store_get(MetricsStore *store, MetricsSnapshot *out);

/* Alarm management */
void metrics_store_set_alarm(MetricsStore *store, AlarmType type,
                             Severity severity, const char *msg);
void metrics_store_clear_alarm(MetricsStore *store, AlarmType type);
int  metrics_store_get_alarm_history(MetricsStore *store, AlarmRecord *out,
                                     int max_records);

/* Helper functions */
const char *severity_str(Severity sev);
const char *alarm_type_str(AlarmType type);

#endif /* METRICS_STORE_H */
