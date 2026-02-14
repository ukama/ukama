/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef LANES_H
#define LANES_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "jobs.h"
#include "snapshot.h"
#include "safety.h"
#include "i2c_bus.h"
#include "gpio_controller.h"

typedef struct {
    struct Lanes   *lanes;
    LaneId          lane;
    FemUnit         femUnit;
    I2cBus         *bus;
} LaneCtx;

typedef struct Lanes {
    Jobs          *jobs;
    SnapshotStore *snap;
    Safety        *safety;
    GpioController *gpio;

    I2cBus        *busCtrl;
    I2cBus        *busFem1;
    I2cBus        *busFem2;

    uint32_t       samplePeriodMs;
    uint32_t       safetyPeriodMs;

    pthread_t      threadCtrl;
    pthread_t      threadFem1;
    pthread_t      threadFem2;

    LaneCtx        ctxCtrl;
    LaneCtx        ctxFem1;
    LaneCtx        ctxFem2;

    bool           running;
    bool           initialized;
} Lanes;

int  lanes_init(Lanes *lanes,
                Jobs *jobs,
                SnapshotStore *snap,
                Safety *safety,
                I2cBus *busFem1,
                I2cBus *busFem2,
                I2cBus *busCtrl,
                GpioController *gpio,
                uint32_t samplePeriodMs,
                uint32_t safetyPeriodMs);

int  lanes_start(Lanes *lanes);
int  lanes_stop(Lanes *lanes);
void lanes_cleanup(Lanes *lanes);

#endif /* LANES_H */
