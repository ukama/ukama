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

typedef struct {
    const Config    *config;
    MetricsStore    *store;

    int             low_volt_count;
    int             high_temp_count;
    int             comm_fail_count;
    int             fault_count;

    int             debounce_samples;
} AlarmChecker;

int alarms_init(AlarmChecker *checker, const Config *config, MetricsStore *store);

void alarms_check(AlarmChecker *checker, const ControllerData *data);

void alarms_check_comm_failure(AlarmChecker *checker, bool comm_ok);

int alarms_send_notification(const Config *config, AlarmType type,
                             Severity severity, const char *message);

#endif /* ALARMS_H */
