/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <sys/time.h>

#include "lanes.h"
#include "i2c_controller.h"

#include "usys_log.h"

typedef struct {
    LaneId         lane;
    FemUnit        femUnit;
    Jobs          *jobs;
    SnapshotStore *snap;
    GpioController *gpio;
    I2cBus        *bus;
} LaneCtx;

static uint32_t now_ms(void) {

    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (uint32_t)(tv.tv_sec * 1000UL + tv.tv_usec / 1000UL);
}

static int fem_sample(LaneCtx *ctx, FemSnapshot *out) {

    FemSnapshot s;
    uint32_t ts;

    if (!ctx || !out) return STATUS_NOK;

    ts = now_ms();
    memset(&s, 0, sizeof(s));

    s.sampleTsMs = ts;
    s.present    = true;

    if (ctx->gpio) {
        if (gpio_read_all(ctx->gpio, ctx->femUnit, &s.gpio) == STATUS_OK) {
            s.haveGpio = true;
        }
    }

    if (temp_sensor_read(ctx->bus, ctx->femUnit, &s.tempC) == STATUS_OK) {
        s.haveTemp = true;
    }

    {
        float v = 0.0f;

        if (adc_read_reverse_power(ctx->bus, ctx->femUnit, &v) == STATUS_OK) {
            s.reversePowerDbm = v;
            s.haveAdc = true;
        }

        if (adc_read_forward_power(ctx->bus, ctx->femUnit, &v) == STATUS_OK) {
            s.forwardPowerDbm = v;
            s.haveAdc = true;
        }

        if (adc_read_pa_current(ctx->bus, ctx->femUnit, &v) == STATUS_OK) {
            s.paCurrentA = v;
            s.haveAdc = true;
        }
    }

    {
        float carrier = 0.0f;
        float peak = 0.0f;

        if (dac_get_config(ctx->bus, ctx->femUnit, &carrier, &peak) == STATUS_OK) {
            s.carrierVoltage = carrier;
            s.peakVoltage    = peak;
            s.haveDac        = true;
        }
    }

    *out = s;
    return STATUS_OK;
}

static int ctrl_sample(LaneCtx *ctx, CtrlSnapshot *out) {

    CtrlSnapshot s;
    uint32_t ts;

    if (!ctx || !out) return STATUS_NOK;

    ts = now_ms();
    memset(&s, 0, sizeof(s));

    s.sampleTsMs = ts;
    s.present    = true;

    if (ctrl_temp_read_tmp10x(ctx->bus, &s.tempC) == STATUS_OK) {
        s.haveTemp = true;
    }

    *out = s;
    return STATUS_OK;
}

static int exec_fem_job(LaneCtx *ctx, const Job *job) {

    int rc = STATUS_NOK;

    if (!ctx || !job) return STATUS_NOK;

    switch (job->cmd) {
    case JobCmdSampleFem: {
        FemSnapshot s;
        if (fem_sample(ctx, &s) == STATUS_OK) {
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &s);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdGpioReadAll: {
        FemSnapshot cur;
        if (snapshot_get_fem(ctx->snap, ctx->femUnit, &cur) != STATUS_OK) {
            memset(&cur, 0, sizeof(cur));
            cur.sampleTsMs = now_ms();
            cur.present = true;
        }
        if (ctx->gpio && gpio_read_all(ctx->gpio, ctx->femUnit, &cur.gpio) == STATUS_OK) {
            cur.haveGpio = true;
            cur.sampleTsMs = now_ms();
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &cur);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdGpioApply:
        if (ctx->gpio && gpio_apply(ctx->gpio, ctx->femUnit, &job->arg.gpioApply.gpio) == STATUS_OK) {
            rc = STATUS_OK;
        }
        break;

    case JobCmdGpioDisablePa:
        if (ctx->gpio && gpio_disable_pa(ctx->gpio, ctx->femUnit) == STATUS_OK) {
            rc = STATUS_OK;
        }
        break;

    case JobCmdDacInit:
        rc = dac_init(ctx->bus, ctx->femUnit);
        break;

    case JobCmdDacSetCarrier:
        rc = dac_set_carrier_voltage(ctx->bus, ctx->femUnit, job->arg.voltage.voltage);
        break;

    case JobCmdDacSetPeak:
        rc = dac_set_peak_voltage(ctx->bus, ctx->femUnit, job->arg.voltage.voltage);
        break;

    case JobCmdDacDisablePa:
        rc = dac_disable_pa(ctx->bus, ctx->femUnit);
        break;

    case JobCmdTempInit:
        rc = temp_sensor_init(ctx->bus, ctx->femUnit);
        break;

    case JobCmdTempRead: {
        FemSnapshot cur;
        if (snapshot_get_fem(ctx->snap, ctx->femUnit, &cur) != STATUS_OK) {
            memset(&cur, 0, sizeof(cur));
            cur.sampleTsMs = now_ms();
            cur.present = true;
        }
        if (temp_sensor_read(ctx->bus, ctx->femUnit, &cur.tempC) == STATUS_OK) {
            cur.haveTemp = true;
            cur.sampleTsMs = now_ms();
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &cur);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdTempSetThreshold:
        rc = temp_sensor_set_threshold(ctx->bus, ctx->femUnit, job->arg.tempThr.thresholdC);
        break;

    case JobCmdAdcInit:
        rc = adc_init(ctx->bus, ctx->femUnit);
        break;

    case JobCmdAdcReadAll:
        rc = adc_read_all_channels(ctx->bus, ctx->femUnit);
        break;

    case JobCmdAdcReadReversePower: {
        float v = 0.0f;
        if (adc_read_reverse_power(ctx->bus, ctx->femUnit, &v) == STATUS_OK) {
            FemSnapshot cur;
            if (snapshot_get_fem(ctx->snap, ctx->femUnit, &cur) != STATUS_OK) {
                memset(&cur, 0, sizeof(cur));
                cur.present = true;
            }
            cur.reversePowerDbm = v;
            cur.haveAdc = true;
            cur.sampleTsMs = now_ms();
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &cur);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdAdcReadPaCurrent: {
        float v = 0.0f;
        if (adc_read_pa_current(ctx->bus, ctx->femUnit, &v) == STATUS_OK) {
            FemSnapshot cur;
            if (snapshot_get_fem(ctx->snap, ctx->femUnit, &cur) != STATUS_OK) {
                memset(&cur, 0, sizeof(cur));
                cur.present = true;
            }
            cur.paCurrentA = v;
            cur.haveAdc = true;
            cur.sampleTsMs = now_ms();
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &cur);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdEepromReadSerial: {
        char serial[SNAPSHOT_SERIAL_MAX_LEN] = {0};
        if (eeprom_read_serial(ctx->bus, ctx->femUnit, serial, sizeof(serial)) == STATUS_OK) {
            FemSnapshot cur;
            if (snapshot_get_fem(ctx->snap, ctx->femUnit, &cur) != STATUS_OK) {
                memset(&cur, 0, sizeof(cur));
                cur.present = true;
            }
            strncpy(cur.serial, serial, sizeof(cur.serial) - 1);
            cur.haveSerial = true;
            cur.sampleTsMs = now_ms();
            (void)snapshot_update_fem(ctx->snap, ctx->femUnit, &cur);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdEepromWriteSerial:
        rc = eeprom_write_serial(ctx->bus, ctx->femUnit, job->arg.serial.serial);
        break;

    case JobCmdShutdownLane:
        rc = STATUS_OK;
        break;

    default:
        rc = STATUS_NOK;
        break;
    }

    return rc;
}

static int exec_ctrl_job(LaneCtx *ctx, const Job *job) {

    int rc = STATUS_NOK;

    if (!ctx || !job) return STATUS_NOK;

    switch (job->cmd) {
    case JobCmdSampleCtrl: {
        CtrlSnapshot s;
        if (ctrl_sample(ctx, &s) == STATUS_OK) {
            (void)snapshot_update_ctrl(ctx->snap, &s);
            rc = STATUS_OK;
        }
        break;
    }

    case JobCmdShutdownLane:
        rc = STATUS_OK;
        break;

    default:
        rc = STATUS_NOK;
        break;
    }

    return rc;
}

static void *lane_main(void *arg) {

    LaneCtx *ctx = (LaneCtx *)arg;
    Job job;
    uint32_t t;
    int rc;

    if (!ctx) return NULL;

    for (;;) {
        t = now_ms();
        if (jobs_dequeue(ctx->jobs, ctx->lane, &job, t) != STATUS_OK) {
            break;
        }

        if (job.cmd == JobCmdShutdownLane) {
            (void)jobs_set_op_state(ctx->jobs, job.opId, OpStateDone, STATUS_OK, now_ms());
            break;
        }

        if (ctx->lane == LaneCtrl) {
            rc = exec_ctrl_job(ctx, &job);
        } else {
            rc = exec_fem_job(ctx, &job);
        }

        (void)jobs_set_op_state(ctx->jobs,
                                job.opId,
                                (rc == STATUS_OK) ? OpStateDone : OpStateFailed,
                                rc,
                                now_ms());
    }

    return NULL;
}

int lanes_init(Lanes *lanes,
               Jobs *jobs,
               SnapshotStore *snap,
               GpioController *gpio,
               int busCtrl,
               int busFem1,
               int busFem2) {

    if (!lanes || !jobs || !snap) return STATUS_NOK;

    memset(lanes, 0, sizeof(*lanes));

    lanes->jobs = jobs;
    lanes->snap = snap;
    lanes->gpio = gpio;

    if (i2c_bus_init(&lanes->busCtrl, busCtrl) != STATUS_OK) return STATUS_NOK;
    if (i2c_bus_init(&lanes->busFem1, busFem1) != STATUS_OK) return STATUS_NOK;
    if (i2c_bus_init(&lanes->busFem2, busFem2) != STATUS_OK) return STATUS_NOK;

    lanes->initialized = true;
    return STATUS_OK;
}

int lanes_start(Lanes *lanes) {

    static LaneCtx ctxCtrl;
    static LaneCtx ctxFem1;
    static LaneCtx ctxFem2;

    if (!lanes || !lanes->initialized) return STATUS_NOK;

    memset(&ctxCtrl, 0, sizeof(ctxCtrl));
    memset(&ctxFem1, 0, sizeof(ctxFem1));
    memset(&ctxFem2, 0, sizeof(ctxFem2));

    ctxCtrl.lane    = LaneCtrl;
    ctxCtrl.femUnit = 0;
    ctxCtrl.jobs    = lanes->jobs;
    ctxCtrl.snap    = lanes->snap;
    ctxCtrl.gpio    = lanes->gpio;
    ctxCtrl.bus     = &lanes->busCtrl;

    ctxFem1.lane    = LaneFem1;
    ctxFem1.femUnit = FEM_UNIT_1;
    ctxFem1.jobs    = lanes->jobs;
    ctxFem1.snap    = lanes->snap;
    ctxFem1.gpio    = lanes->gpio;
    ctxFem1.bus     = &lanes->busFem1;

    ctxFem2.lane    = LaneFem2;
    ctxFem2.femUnit = FEM_UNIT_2;
    ctxFem2.jobs    = lanes->jobs;
    ctxFem2.snap    = lanes->snap;
    ctxFem2.gpio    = lanes->gpio;
    ctxFem2.bus     = &lanes->busFem2;

    if (pthread_create(&lanes->threadCtrl, NULL, lane_main, &ctxCtrl) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem1, NULL, lane_main, &ctxFem1) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem2, NULL, lane_main, &ctxFem2) != 0) return STATUS_NOK;

    return STATUS_OK;
}

int lanes_stop(Lanes *lanes) {

    Job j;
    uint32_t t;

    if (!lanes || !lanes->initialized) return STATUS_NOK;

    t = now_ms();
    memset(&j, 0, sizeof(j));
    j.prio = JobPrioHi;
    j.cmd  = JobCmdShutdownLane;

    j.lane = LaneCtrl;
    (void)jobs_enqueue(lanes->jobs, &j, t);

    j.lane = LaneFem1;
    (void)jobs_enqueue(lanes->jobs, &j, t);

    j.lane = LaneFem2;
    (void)jobs_enqueue(lanes->jobs, &j, t);

    (void)pthread_join(lanes->threadCtrl, NULL);
    (void)pthread_join(lanes->threadFem1, NULL);
    (void)pthread_join(lanes->threadFem2, NULL);

    return STATUS_OK;
}

void lanes_cleanup(Lanes *lanes) {

    if (!lanes || !lanes->initialized) return;

    i2c_bus_cleanup(&lanes->busCtrl);
    i2c_bus_cleanup(&lanes->busFem1);
    i2c_bus_cleanup(&lanes->busFem2);

    memset(lanes, 0, sizeof(*lanes));
}
