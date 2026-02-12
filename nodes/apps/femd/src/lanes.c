/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <string.h>
#include <unistd.h>

#include "usys_log.h"

#include "lanes.h"
#include "i2c_controller.h"

static uint32_t now_ms(void) {
    return snapshot_now_ms();
}

static void lane_sample_ctrl(Lanes *lanes) {
    CtrlSnapshot s;
    float t = 0.0f;
    uint32_t ts = now_ms();

    memset(&s, 0, sizeof(s));
    s.sampleTsMs = ts;

    if (ctrl_temp_read_tmp10x(lanes->busCtrl, &t) == STATUS_OK) {
        s.present  = true;
        s.haveTemp = true;
        s.tempC    = t;
    } else {
        s.present = false;
    }

    (void)snapshot_update_ctrl(lanes->snap, &s);
}

static void lane_sample_fem(Lanes *lanes, FemUnit unit, I2cBus *bus) {
    FemSnapshot s;
    uint32_t ts = now_ms();

    memset(&s, 0, sizeof(s));
    s.sampleTsMs = ts;
    s.present = true;

    if (gpio_read_all(lanes->gpio, unit, &s.gpio) == STATUS_OK) {
        s.haveGpio = true;
    }

    if (temp_sensor_read(bus, &s.tempC) == STATUS_OK) {
        s.haveTemp = true;
    }

    if (adc_read_reverse_power(bus, &s.reversePowerDbm) == STATUS_OK) {
        s.haveAdc = true;
    }
    if (adc_read_forward_power(bus, &s.forwardPowerDbm) == STATUS_OK) {
        s.haveAdc = true;
    }
    if (adc_read_pa_current(bus, &s.paCurrentA) == STATUS_OK) {
        s.haveAdc = true;
    }

    if (dac_get_cached(bus, &s.carrierVoltage, &s.peakVoltage) == STATUS_OK) {
        s.haveDac = true;
    }

    (void)snapshot_update_fem(lanes->snap, unit, &s);
}

static void handle_job(LaneCtx *ctx, const Job *j) {
    Lanes *lanes = ctx->lanes;
    int rc = STATUS_OK;

    switch (j->cmd) {
    case JobCmdSampleCtrl:
        lane_sample_ctrl(lanes);
        break;

    case JobCmdSampleFem:
        lane_sample_fem(lanes, ctx->femUnit, ctx->bus);
        break;

    case JobCmdGpioReadAll: {
        FemSnapshot cur;
        memset(&cur, 0, sizeof(cur));
        if (snapshot_get_fem(lanes->snap, ctx->femUnit, &cur) == STATUS_OK) {
            if (gpio_read_all(lanes->gpio, ctx->femUnit, &cur.gpio) == STATUS_OK) {
                cur.haveGpio = true;
                cur.sampleTsMs = now_ms();
                cur.present = true;
                (void)snapshot_update_fem(lanes->snap, ctx->femUnit, &cur);
            }
        }
        break;
    }

    case JobCmdGpioApply:
        rc = gpio_apply(lanes->gpio, ctx->femUnit, &j->arg.gpioApply.gpio);
        break;

    case JobCmdGpioDisablePa:
        rc = gpio_disable_pa(lanes->gpio, ctx->femUnit);
        break;

    case JobCmdDacInit:
        rc = dac_init(ctx->bus);
        break;

    case JobCmdDacSetCarrier:
        rc = dac_set_carrier_voltage(ctx->bus, j->arg.voltage.voltage);
        break;

    case JobCmdDacSetPeak:
        rc = dac_set_peak_voltage(ctx->bus, j->arg.voltage.voltage);
        break;

    case JobCmdDacDisablePa:
        rc = dac_disable_pa(ctx->bus);
        break;

    case JobCmdTempInit:
        rc = temp_sensor_init(ctx->bus);
        break;

    case JobCmdTempRead: {
        float t = 0.0f;
        rc = temp_sensor_read(ctx->bus, &t);
        if (rc == STATUS_OK) {
            FemSnapshot cur;
            memset(&cur, 0, sizeof(cur));
            if (snapshot_get_fem(lanes->snap, ctx->femUnit, &cur) == STATUS_OK) {
                cur.tempC = t;
                cur.haveTemp = true;
                cur.sampleTsMs = now_ms();
                cur.present = true;
                (void)snapshot_update_fem(lanes->snap, ctx->femUnit, &cur);
            }
        }
        break;
    }

    case JobCmdTempSetThreshold:
        rc = temp_sensor_set_threshold(ctx->bus, j->arg.tempThr.thresholdC);
        break;

    case JobCmdAdcInit:
        rc = adc_init(ctx->bus);
        break;

    case JobCmdAdcReadReversePower:
    case JobCmdAdcReadPaCurrent:
    case JobCmdAdcReadAll:
        lane_sample_fem(lanes, ctx->femUnit, ctx->bus);
        rc = STATUS_OK;
        break;

    case JobCmdSafetyDisablePa: {
        Job jj;
        uint32_t t = now_ms();

        memset(&jj, 0, sizeof(jj));
        jj.lane   = ctx->lane;
        jj.femUnit= ctx->femUnit;
        jj.prio   = JobPrioHi;

        jj.cmd = JobCmdGpioDisablePa;
        (void)jobs_enqueue(lanes->jobs, &jj, t);

        jj.cmd = JobCmdDacDisablePa;
        (void)jobs_enqueue(lanes->jobs, &jj, t);
        rc = STATUS_OK;
        break;
    }

    case JobCmdSafetyRestorePa: {
        FemSnapshot cur;
        GpioStatus desired;
        float carrierV = 0.0f;
        float peakV = 0.0f;
        uint32_t t = now_ms();

        memset(&cur, 0, sizeof(cur));
        if (snapshot_get_fem(lanes->snap, ctx->femUnit, &cur) == STATUS_OK && cur.haveGpio) {
            desired = cur.gpio;
        } else {
            memset(&desired, 0, sizeof(desired));
        }

        desired.pa_vds_enable = true;
        desired.pa_disable    = true;

        Job jj;
        memset(&jj, 0, sizeof(jj));
        jj.lane    = ctx->lane;
        jj.femUnit = ctx->femUnit;
        jj.prio    = JobPrioHi;

        jj.cmd = JobCmdGpioApply;
        jj.arg.gpioApply.gpio = desired;
        (void)jobs_enqueue(lanes->jobs, &jj, t);

        rc = STATUS_OK;
        break;
    }

    case JobCmdShutdownLane:
        rc = STATUS_OK;
        break;

    default:
        rc = STATUS_NOK;
        break;
    }

    (void)jobs_set_op_state(lanes->jobs, j->opId, rc == STATUS_OK ? OpStateDone : OpStateFailed, rc, now_ms());
}

static void* lane_main(void *p) {
    LaneCtx *ctx = (LaneCtx *)p;
    Lanes *lanes = ctx->lanes;

    uint32_t nextSample = now_ms();
    uint32_t nextSafety = now_ms();

    for (;;) {
        Job j;
        uint32_t t = now_ms();

        if (jobs_try_dequeue(lanes->jobs, ctx->lane, &j, t) == STATUS_OK) {
            (void)jobs_set_op_state(lanes->jobs, j.opId, OpStateRunning, STATUS_OK, t);
            if (j.cmd == JobCmdShutdownLane) {
                (void)jobs_set_op_state(lanes->jobs, j.opId, OpStateDone, STATUS_OK, now_ms());
                break;
            }
            handle_job(ctx, &j);
        }

        t = now_ms();

        if (ctx->lane == LaneCtrl && lanes->samplePeriodMs > 0 && t >= nextSample) {
            lane_sample_ctrl(lanes);
            nextSample = t + lanes->samplePeriodMs;
        }

        if ((ctx->lane == LaneFem1 || ctx->lane == LaneFem2) && lanes->samplePeriodMs > 0 && t >= nextSample) {
            lane_sample_fem(lanes, ctx->femUnit, ctx->bus);
            nextSample = t + lanes->samplePeriodMs;
        }

        if ((ctx->lane == LaneFem1 || ctx->lane == LaneFem2) && lanes->safetyPeriodMs > 0 && lanes->safety && t >= nextSafety) {
            (void)safety_tick(lanes->safety, ctx->femUnit);
            nextSafety = t + lanes->safetyPeriodMs;
        }

        usleep(20000);
    }

    return NULL;
}

int lanes_init(Lanes *lanes,
               Jobs *jobs,
               SnapshotStore *snap,
               Safety *safety,
               I2cBus *busFem1,
               I2cBus *busFem2,
               I2cBus *busCtrl,
               GpioController *gpio,
               uint32_t samplePeriodMs,
               uint32_t safetyPeriodMs) {

    if (!lanes || !jobs || !snap || !busFem1 || !busFem2 || !busCtrl || !gpio) return STATUS_NOK;

    memset(lanes, 0, sizeof(*lanes));

    lanes->jobs = jobs;
    lanes->snap = snap;
    lanes->safety = safety;
    lanes->gpio = gpio;

    lanes->busFem1 = busFem1;
    lanes->busFem2 = busFem2;
    lanes->busCtrl = busCtrl;

    lanes->samplePeriodMs = samplePeriodMs;
    lanes->safetyPeriodMs = safetyPeriodMs;

    memset(&lanes->ctxCtrl, 0, sizeof(lanes->ctxCtrl));
    lanes->ctxCtrl.lanes = lanes;
    lanes->ctxCtrl.lane  = LaneCtrl;
    lanes->ctxCtrl.femUnit = 0;
    lanes->ctxCtrl.bus   = lanes->busCtrl;

    memset(&lanes->ctxFem1, 0, sizeof(lanes->ctxFem1));
    lanes->ctxFem1.lanes = lanes;
    lanes->ctxFem1.lane  = LaneFem1;
    lanes->ctxFem1.femUnit = FEM_UNIT_1;
    lanes->ctxFem1.bus   = lanes->busFem1;

    memset(&lanes->ctxFem2, 0, sizeof(lanes->ctxFem2));
    lanes->ctxFem2.lanes = lanes;
    lanes->ctxFem2.lane  = LaneFem2;
    lanes->ctxFem2.femUnit = FEM_UNIT_2;
    lanes->ctxFem2.bus   = lanes->busFem2;

    lanes->initialized = true;
    return STATUS_OK;
}

int lanes_start(Lanes *lanes) {
    if (!lanes || !lanes->initialized || lanes->running) return STATUS_NOK;

    if (pthread_create(&lanes->threadCtrl, NULL, lane_main, &lanes->ctxCtrl) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem1, NULL, lane_main, &lanes->ctxFem1) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem2, NULL, lane_main, &lanes->ctxFem2) != 0) return STATUS_NOK;

    lanes->running = true;
    return STATUS_OK;
}

int lanes_stop(Lanes *lanes) {
    Job j;
    uint32_t t;

    if (!lanes || !lanes->initialized || !lanes->running) return STATUS_NOK;

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

    lanes->running = false;
    return STATUS_OK;
}

void lanes_cleanup(Lanes *lanes) {
    if (!lanes || !lanes->initialized) return;
    memset(lanes, 0, sizeof(*lanes));
}
