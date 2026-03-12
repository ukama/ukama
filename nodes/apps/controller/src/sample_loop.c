/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "sample_loop.h"
#include "alarms.h"
#include "time_util.h"
#include "usys_log.h"

static void *sample_thread(void *arg) {
    SampleLoop *loop = (SampleLoop *)arg;
    ControllerData data;
    AlarmChecker alarmChecker;
    uint64_t next_sample_time;
    int ret;

    if (!loop || !loop->config || !loop->driver || !loop->store) {
        usys_log_error("sample_loop: invalid context");
        return NULL;
    }

    /* Initialize alarm checker */
    alarms_init(&alarmChecker, loop->config, loop->store);

    usys_log_info("sample_loop: started (interval=%d ms)", loop->config->sampleMs);

    next_sample_time = time_now_ms();

    while (loop->running) {
        uint64_t now = time_now_ms();

        /* Wait until next sample time */
        if (now < next_sample_time) {
            usleep((next_sample_time - now) * 1000);
            continue;
        }

        /* Read data from controller */
        memset(&data, 0, sizeof(data));
        ret = loop->driver->read_data(loop->driver_ctx, &data);

        if (ret == 0) {
            /* Success - update metrics store */
            metrics_store_update(loop->store, &data);
            loop->samples_taken++;

            /* Check for alarms */
            alarms_check(&alarmChecker, &data);

            usys_log_trace("sample_loop: V=%.2f I=%.2f PV=%.1fW state=%s",
                           data.batt_voltage_v, data.batt_current_a,
                           data.pv_power_w, charge_state_str(data.charge_state));
        } else {
            /* Error - update error state */
            loop->samples_failed++;
            metrics_store_set_error(loop->store, ret, "Failed to read from controller");
            alarms_check_comm_failure(&alarmChecker, false);

            usys_log_warn("sample_loop: read failed (total_fails=%lu)",
                          (unsigned long)loop->samples_failed);
        }

        /* Schedule next sample */
        next_sample_time += loop->config->sampleMs;

        /* If we've fallen behind, catch up */
        if (next_sample_time < now) {
            next_sample_time = now + loop->config->sampleMs;
        }
    }

    usys_log_info("sample_loop: stopped (samples=%lu, failed=%lu)",
                  (unsigned long)loop->samples_taken,
                  (unsigned long)loop->samples_failed);

    return NULL;
}

int sample_loop_start(SampleLoop *loop, const Config *config,
                      const ControllerDriver *driver, void *driver_ctx,
                      MetricsStore *store) {
    if (!loop || !config || !driver || !store) {
        return -1;
    }

    memset(loop, 0, sizeof(*loop));
    loop->config = config;
    loop->driver = driver;
    loop->driver_ctx = driver_ctx;
    loop->store = store;
    loop->running = true;

    if (pthread_create(&loop->thread, NULL, sample_thread, loop) != 0) {
        usys_log_error("sample_loop: pthread_create failed");
        loop->running = false;
        return -1;
    }

    return 0;
}

void sample_loop_stop(SampleLoop *loop) {
    if (!loop || !loop->running) return;

    loop->running = false;
    pthread_join(loop->thread, NULL);
}

void sample_loop_get_stats(SampleLoop *loop, uint64_t *taken, uint64_t *failed) {
    if (!loop) return;
    if (taken)  *taken  = loop->samples_taken;
    if (failed) *failed = loop->samples_failed;
}
