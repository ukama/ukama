/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <unistd.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>

#include "i2c_bus.h"

#include "usys_log.h"

#define I2C_DEV_PATH_MAX  32
#define I2C_TX_MAX        256  /* reg + payload */

int i2c_bus_init(I2cBus *bus, int busNum) {

    char devPath[I2C_DEV_PATH_MAX] = {0};

    if (!bus) {
        usys_log_error("i2c_bus_init: bus is NULL");
        return STATUS_NOK;
    }

    memset(bus, 0, sizeof(*bus));
    bus->busNum       = busNum;
    bus->fd           = -1;
    bus->currentSlave = -1;
    bus->initialized  = false;

    if (snprintf(devPath, sizeof(devPath), "/dev/i2c-%d", busNum) >= (int)sizeof(devPath)) {
        usys_log_error("i2c_bus_init: dev path truncated for bus %d", busNum);
        return STATUS_NOK;
    }

    bus->fd = open(devPath, O_RDWR | O_CLOEXEC);
    if (bus->fd < 0) {
        usys_log_error("i2c_bus_init: open %s failed: %s", devPath, strerror(errno));
        return STATUS_NOK;
    }

    bus->initialized = true;
    usys_log_info("I2C bus initialized: %s (fd=%d)", devPath, bus->fd);

    return STATUS_OK;
}

void i2c_bus_cleanup(I2cBus *bus) {

    if (!bus) return;

    if (bus->fd >= 0) {
        close(bus->fd);
        bus->fd = -1;
    }

    bus->currentSlave = -1;
    bus->initialized  = false;

    usys_log_info("I2C bus cleaned up (bus=%d)", bus->busNum);
}

int i2c_bus_set_slave(I2cBus *bus, uint8_t devAddr) {

    if (!bus || !bus->initialized || bus->fd < 0) {
        return STATUS_NOK;
    }

    /* Avoid ioctl if already selected */
    if (bus->currentSlave == (int)devAddr) {
        return STATUS_OK;
    }

    if (ioctl(bus->fd, I2C_SLAVE, devAddr) < 0) {
        usys_log_error("i2c_bus_set_slave: bus=%d addr=0x%02X ioctl failed: %s",
                       bus->busNum, devAddr, strerror(errno));
        return STATUS_NOK;
    }

    bus->currentSlave = (int)devAddr;
    return STATUS_OK;
}

int i2c_bus_write_reg(I2cBus *bus, uint8_t devAddr, uint8_t reg,
                      const uint8_t *data, size_t len) {

    uint8_t tx[I2C_TX_MAX] = {0};
    ssize_t wrote = 0;

    if (!bus || !bus->initialized || bus->fd < 0) {
        return STATUS_NOK;
    }

    if (len > (I2C_TX_MAX - 1)) {
        usys_log_error("i2c_bus_write_reg: len too large (%zu)", len);
        return STATUS_NOK;
    }

    if (i2c_bus_set_slave(bus, devAddr) != STATUS_OK) {
        return STATUS_NOK;
    }

    tx[0] = reg;
    if (len > 0 && data) {
        memcpy(&tx[1], data, len);
    }

    wrote = write(bus->fd, tx, (size_t)(len + 1));
    if (wrote != (ssize_t)(len + 1)) {
        usys_log_error("i2c_bus_write_reg: bus=%d addr=0x%02X reg=0x%02X len=%zu failed: %s",
                       bus->busNum, devAddr, reg, len, strerror(errno));
        return STATUS_NOK;
    }

    usys_log_debug("i2c_bus_write_reg: bus=%d addr=0x%02X reg=0x%02X len=%zu OK",
                   bus->busNum, devAddr, reg, len);

    return STATUS_OK;
}

int i2c_bus_read_reg(I2cBus *bus, uint8_t devAddr, uint8_t reg,
                     uint8_t *data, size_t len) {

    ssize_t wrote = 0;
    ssize_t readn = 0;

    if (!bus || !bus->initialized || bus->fd < 0 || !data || len == 0) {
        return STATUS_NOK;
    }

    if (i2c_bus_set_slave(bus, devAddr) != STATUS_OK) {
        return STATUS_NOK;
    }

    /* Write register pointer */
    wrote = write(bus->fd, &reg, 1);
    if (wrote != 1) {
        usys_log_error("i2c_bus_read_reg: bus=%d addr=0x%02X reg=0x%02X write(reg) failed: %s",
                       bus->busNum, devAddr, reg, strerror(errno));
        return STATUS_NOK;
    }

    /* Read payload */
    readn = read(bus->fd, data, len);
    if (readn != (ssize_t)len) {
        usys_log_error("i2c_bus_read_reg: bus=%d addr=0x%02X reg=0x%02X len=%zu read failed: %s",
                       bus->busNum, devAddr, reg, len, strerror(errno));
        return STATUS_NOK;
    }

    usys_log_debug("i2c_bus_read_reg: bus=%d addr=0x%02X reg=0x%02X len=%zu OK",
                   bus->busNum, devAddr, reg, len);

    return STATUS_OK;
}

int i2c_bus_detect(I2cBus *bus, uint8_t devAddr) {

    uint8_t dummy = 0x00;
    ssize_t wrote = 0;

    if (!bus || !bus->initialized || bus->fd < 0) {
        return STATUS_NOK;
    }

    if (i2c_bus_set_slave(bus, devAddr) != STATUS_OK) {
        return STATUS_NOK;
    }

    /*
     * Best-effort probe: attempt a 1-byte write (no reg).
     * Some devices may NACK depending on state; this is only a soft check.
     */
    wrote = write(bus->fd, &dummy, 0);
    if (wrote < 0) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}
