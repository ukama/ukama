/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef JOBS_H
#define JOBS_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "femd.h"
#include "gpio_controller.h"

typedef enum {
    LaneCtrl = 0,
    LaneFem1 = 1,
    LaneFem2 = 2,
    LaneMax  = 3
} LaneId;

typedef enum {
    JobPrioLo = 0,
    JobPrioHi = 1
} JobPrio;

typedef enum {
    JobCmdNone = 0,

    JobCmdSampleFem,
    JobCmdSampleCtrl,

    JobCmdGpioReadAll,
    JobCmdGpioApply,
    JobCmdGpioDisablePa,

    JobCmdDacInit,
    JobCmdDacSetCarrier,
    JobCmdDacSetPeak,
    JobCmdDacDisablePa,

    JobCmdTempInit,
    JobCmdTempRead,
    JobCmdTempSetThreshold,

    JobCmdAdcInit,
    JobCmdAdcReadAll,
    JobCmdAdcReadReversePower,
    JobCmdAdcReadPaCurrent,

    JobCmdEepromReadSerial,
    JobCmdEepromWriteSerial,

    JobCmdShutdownLane
} JobCmd;

typedef enum {
    OpStateUnknown = 0,
    OpStateQueued,
    OpStateRunning,
    OpStateDone,
    OpStateFailed,
    OpStateCanceled
} OpState;

typedef struct {
    float voltage;
} JobArgVoltage;

typedef struct {
    float thresholdC;
} JobArgTempThreshold;

typedef struct {
    GpioStatus gpio;
} JobArgGpioApply;

typedef struct {
    char serial[17];
} JobArgSerial;

typedef struct {
    LaneId   lane;
    FemUnit  femUnit;
    JobCmd   cmd;
    JobPrio  prio;
    uint64_t opId;

    union {
        JobArgVoltage       voltage;
        JobArgTempThreshold tempThr;
        JobArgGpioApply     gpioApply;
        JobArgSerial        serial;
    } arg;
} Job;

typedef struct {
    uint64_t opId;
    LaneId   lane;
    FemUnit  femUnit;
    JobCmd   cmd;

    OpState  state;
    int      result;

    uint32_t createdMs;
    uint32_t startedMs;
    uint32_t endedMs;
} OpStatus;

typedef struct {
    pthread_mutex_t mu;
    pthread_cond_t  cv;

    Job             hiQ[64];
    Job             loQ[256];

    uint32_t        hiHead;
    uint32_t        hiTail;
    uint32_t        hiCount;

    uint32_t        loHead;
    uint32_t        loTail;
    uint32_t        loCount;

    bool            stop;
} LaneQueue;

typedef struct {
    pthread_mutex_t mu;
    uint64_t        nextOpId;

    OpStatus        table[512];
    uint32_t        tablePos;

    LaneQueue       lane[LaneMax];
    bool            initialized;
} Jobs;

int  jobs_init(Jobs *jobs);
void jobs_cleanup(Jobs *jobs);

uint64_t jobs_enqueue(Jobs *jobs, const Job *in, uint32_t nowMs);

int  jobs_dequeue(Jobs *jobs, LaneId lane, Job *out, uint32_t nowMs);
int  jobs_set_op_state(Jobs *jobs, uint64_t opId, OpState state, int result, uint32_t nowMs);

int  jobs_get_op(Jobs *jobs, uint64_t opId, OpStatus *out);

int  jobs_request_lane_stop(Jobs *jobs, LaneId lane);

#endif /* JOBS_H */
