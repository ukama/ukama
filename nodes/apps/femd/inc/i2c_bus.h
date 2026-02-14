/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef I2C_BUS_H
#define I2C_BUS_H

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

#include "femd.h"

#include "usys_types.h"

typedef struct {
    int  busNum;
    int  fd;
    int  currentSlave;
    bool initialized;

    /* Userspace simulation mode (dev laptop / CI).
     * Enabled when FEMD_SIM=1/true/yes/on OR FEMD_SYSROOT==/tmp/sys.
     */
    bool     sim;
    uint64_t simStartMs;
    int      simAdcChannel; /* 0..3 (ADS1015 mux channel selection) */

    bool  dacKnown;
    float dacCarrierV;
    float dacPeakV;
} I2cBus;

/* lifecycle */
int  i2c_bus_init(I2cBus *bus, int busNum);
void i2c_bus_cleanup(I2cBus *bus);

/* slave selection (cached; only ioctl when addr changes) */
int  i2c_bus_set_slave(I2cBus *bus, uint8_t devAddr);

/* register IO helpers */
int  i2c_bus_write_reg(I2cBus *bus, uint8_t devAddr, uint8_t reg,
                       const uint8_t *data, size_t len);

int  i2c_bus_read_reg(I2cBus *bus, uint8_t devAddr, uint8_t reg,
                      uint8_t *data, size_t len);

/* simple presence probe */
int  i2c_bus_detect(I2cBus *bus, uint8_t devAddr);

#endif /* I2C_BUS_H */
