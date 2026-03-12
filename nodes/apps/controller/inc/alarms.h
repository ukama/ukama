/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ALARMS_H
#define ALARMS_H

#include "config.h"
#include "metrics_store.h"
#include "driver.h"

/*
 * Alarm checker - evaluates controller data against thresholds
 * and updates the metrics store with alarm state
 */
typedef struct {
    const Config    *config;
    MetricsStore    *store;

    /* Hysteresis tracking */
    int             low_volt_count;
    int             high_temp_count;
    int             comm_fail_count;
    int             fault_count;

    /* Debounce threshold (number of consecutive samples) */
    int             debounce_samples;
} AlarmChecker;

/* Initialize alarm checker */
int alarms_init(AlarmChecker *checker, const Config *config, MetricsStore *store);

/* Check controller data and update alarms */
void alarms_check(AlarmChecker *checker, const ControllerData *data);

/* Check for communication failure */
void alarms_check_comm_failure(AlarmChecker *checker, bool comm_ok);

/* Send alarm notification to notify.d */
int alarms_send_notification(const Config *config, AlarmType type,
                             Severity severity, const char *message);

#endif /* ALARMS_H */
