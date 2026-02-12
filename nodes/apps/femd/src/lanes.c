/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <pthread.h>
#include <string.h>
#include <sys/time.h>
#include <unistd.h>

#include "usys_log.h"
#include "usys_types.h"

#include "lanes.h"
#include "jobs.h"
#include "snapshot.h"
#include "gpio_controller.h"

typedef struct {
    Lanes   *lanes;
    LaneId   lane;
    FemUnit  femUnit;
    void    *bus;
} LaneCtx;

static uint32_t now_ms(void) {

    struct timeval tv;

    gettimeofday(&tv, NULL);
    return (uint32_t)(tv.tv_sec * 1000UL + tv.tv_usec / 1000UL);
}

static void sample_ctrl(LaneCtx *ctx, uint32_t tsMs) {

    CtrlSnapshot s;

    if (!ctx || !ctx->lanes || !ctx->lanes->snap) return;

    memset(&s, 0, sizeof(s));
    s.sampleTsMs = tsMs;
    s.present    = true;

    (void)snapshot_update_ctrl(ctx->lanes->snap, &s);
}

static void sample_fem_gpio(LaneCtx *ctx, uint32_t tsMs) {

    FemSnapshot s;

    if (!ctx || !ctx->lanes || !ctx->lanes->snap) return;

    memset(&s, 0, sizeof(s));
    s.sampleTsMs = tsMs;
    s.present    = true;

    if (ctx->lanes->gpio && gpio_read_all(ctx->lanes->gpio, ctx->femUnit, &s.gpio) == STATUS_OK) {
        s.haveGpio = true;
    } else {
        s.haveGpio = false;
    }

    (void)snapshot_update_fem(ctx->lanes->snap, ctx->femUnit, &s);
}

static void lane_sample(LaneCtx *ctx) {

    uint32_t tsMs;

    if (!ctx) return;

    tsMs = now_ms();

    if (ctx->lane == LaneCtrl) {
        sample_ctrl(ctx, tsMs);
        return;
    }

    sample_fem_gpio(ctx, tsMs);
}

static void handle_job(LaneCtx *ctx, const Job *j, int *stop) {

    if (!ctx || !j) return;

    switch (j->cmd) {

    case JobCmdShutdownLane:
        if (stop) *stop = 1;
        break;

    case JobCmdGpioReadAll:
        if (ctx->lane == LaneFem1 || ctx->lane == LaneFem2) {
            sample_fem_gpio(ctx, now_ms());
        }
        break;

    case JobCmdGpioApply:
        if (ctx->lane == LaneFem1 || ctx->lane == LaneFem2) {
            (void)gpio_apply(ctx->lanes->gpio, ctx->femUnit, &j->arg.gpioApply.gpio);
        }
        break;

    case JobCmdGpioDisablePa:
        if (ctx->lane == LaneFem1 || ctx->lane == LaneFem2) {
            (void)gpio_disable_pa(ctx->lanes->gpio, ctx->femUnit);
        }
        break;

    case JobCmdDacDisablePa:
        if (ctx->lane == LaneFem1 || ctx->lane == LaneFem2) {
            (void)gpio_disable_pa(ctx->lanes->gpio, ctx->femUnit);
        }
        break;

    default:
        break;
    }
}

static void *lane_main(void *arg) {

    LaneCtx *ctx;

    uint32_t nextSampleMs;
    int stop = 0;

    ctx = (LaneCtx *)arg;
    if (!ctx || !ctx->lanes || !ctx->lanes->initialized) return NULL;

    nextSampleMs = now_ms();

    while (!stop) {

        uint32_t nowMs = now_ms();

        if ((int32_t)(nowMs - nextSampleMs) >= 0) {
            lane_sample(ctx);
            nextSampleMs = nowMs + ctx->lanes->samplePeriodMs;
        }

        {
            Job j;

            memset(&j, 0, sizeof(j));
            if (jobs_dequeue(ctx->lanes->jobs, ctx->lane, &j, nowMs) == STATUS_OK) {
                handle_job(ctx, &j, &stop);
            } else {
                (void)usleep(2000);
            }
        }
    }

    return NULL;
}

int lanes_init(Lanes          *lanes,
               Jobs           *jobs,
               SnapshotStore  *snap,
               Safety         *safety,
               void           *busFem1,
               void           *busFem2,
               void           *busCtrl,
               GpioController *gpio,
               uint32_t        samplePeriodMs,
               uint32_t        safetyPeriodMs) {

    if (!lanes || !jobs || !snap || !gpio) return STATUS_NOK;

    memset(lanes, 0, sizeof(*lanes));

    lanes->jobs           = jobs;
    lanes->snap           = snap;
    lanes->safety         = safety;
    lanes->busFem1        = busFem1;
    lanes->busFem2        = busFem2;
    lanes->busCtrl        = busCtrl;
    lanes->gpio           = gpio;
    lanes->samplePeriodMs = samplePeriodMs;
    lanes->safetyPeriodMs = safetyPeriodMs;
    lanes->initialized    = true;

    memset(&lanes->ctxCtrl, 0, sizeof(lanes->ctxCtrl));
    lanes->ctxCtrl.lanes   = lanes;
    lanes->ctxCtrl.lane    = LaneCtrl;
    lanes->ctxCtrl.femUnit = 0;
    lanes->ctxCtrl.bus     = lanes->busCtrl;

    memset(&lanes->ctxFem1, 0, sizeof(lanes->ctxFem1));
    lanes->ctxFem1.lanes   = lanes;
    lanes->ctxFem1.lane    = LaneFem1;
    lanes->ctxFem1.femUnit = FEM_UNIT_1;
    lanes->ctxFem1.bus     = lanes->busFem1;

    memset(&lanes->ctxFem2, 0, sizeof(lanes->ctxFem2));
    lanes->ctxFem2.lanes   = lanes;
    lanes->ctxFem2.lane    = LaneFem2;
    lanes->ctxFem2.femUnit = FEM_UNIT_2;
    lanes->ctxFem2.bus     = lanes->busFem2;

    if (pthread_create(&lanes->threadCtrl, NULL, lane_main, &lanes->ctxCtrl) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem1, NULL, lane_main, &lanes->ctxFem1) != 0) return STATUS_NOK;
    if (pthread_create(&lanes->threadFem2, NULL, lane_main, &lanes->ctxFem2) != 0) return STATUS_NOK;

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

    if (!lanes) return;
    memset(lanes, 0, sizeof(*lanes));
}
