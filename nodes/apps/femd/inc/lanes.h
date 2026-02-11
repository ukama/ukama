/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef LANES_H
#define LANES_H

#include <pthread.h>
#include <stdbool.h>

#include "femd.h"
#include "jobs.h"
#include "snapshot.h"
#include "i2c_bus.h"
#include "gpio_controller.h"

typedef struct {
    Jobs          *jobs;
    SnapshotStore *snap;
    GpioController *gpio;

    I2cBus         busCtrl;
    I2cBus         busFem1;
    I2cBus         busFem2;

    pthread_t      threadCtrl;
    pthread_t      threadFem1;
    pthread_t      threadFem2;

    bool           initialized;
} Lanes;

int  lanes_init(Lanes *lanes,
                Jobs *jobs,
                SnapshotStore *snap,
                GpioController *gpio,
                int busCtrl,
                int busFem1,
                int busFem2);

int  lanes_start(Lanes *lanes);
int  lanes_stop(Lanes *lanes);
void lanes_cleanup(Lanes *lanes);

#endif /* LANES_H */
