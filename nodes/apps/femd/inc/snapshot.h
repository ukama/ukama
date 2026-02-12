/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#ifndef SNAPSHOT_H
#define SNAPSHOT_H

#include <stdint.h>
#include <stdbool.h>
#include <pthread.h>

#include "femd.h"
#include "gpio_controller.h"

#define SNAPSHOT_SERIAL_MAX_LEN  17

typedef struct {
    uint32_t   sampleTsMs;
    bool       present;

    bool       haveGpio;
    GpioStatus gpio;

    bool       haveTemp;
    float      tempC;

    bool       haveAdc;
    float      reversePowerDbm;
    float      forwardPowerDbm;
    float      paCurrentA;
    float      adcTempVolts;

    bool       haveDac;
    float      carrierVoltage;
    float      peakVoltage;

    bool       haveSerial;
    char       serial[SNAPSHOT_SERIAL_MAX_LEN];
} FemSnapshot;

typedef struct {
    uint32_t sampleTsMs;
    bool     present;

    bool     haveTemp;
    float    tempC;
} CtrlSnapshot;

typedef struct {
    pthread_rwlock_t lock;
    FemSnapshot      fem[3];
    CtrlSnapshot     ctrl;
    bool             initialized;
} SnapshotStore;

int  snapshot_init(SnapshotStore *store);
void snapshot_cleanup(SnapshotStore *store);

int  snapshot_set_fem_present(SnapshotStore *store, FemUnit unit, bool present, uint32_t tsMs);
int  snapshot_set_ctrl_present(SnapshotStore *store, bool present, uint32_t tsMs);

int  snapshot_update_fem(SnapshotStore *store, FemUnit unit, const FemSnapshot *in);
int  snapshot_update_ctrl(SnapshotStore *store, const CtrlSnapshot *in);

int  snapshot_get_fem(SnapshotStore *store, FemUnit unit, FemSnapshot *out);
int  snapshot_get_ctrl(SnapshotStore *store, CtrlSnapshot *out);

uint32_t snapshot_now_ms(void);

#endif /* SNAPSHOT_H */
