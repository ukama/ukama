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
#include <math.h>
#include <time.h>
#include <stdlib.h>
#include <sys/ioctl.h>
#include <linux/i2c-dev.h>

#include "i2c_bus.h"

#include "usys_log.h"

#define I2C_TX_MAX        256  /* reg + payload */

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

static const char *femd_sysroot(void) {
    const char *v = getenv(ENV_FEMD_SYSROOT);
    return (v && v[0] != '\0') ? v : NULL;
}

static int build_i2c_dev_path(char *out, size_t outsz, int busNum) {
    const char *root = femd_sysroot();
    int n = 0;

    if (root) {
        n = snprintf(out, outsz, "%s/dev/i2c-%d", root, busNum);
    } else {
        n = snprintf(out, outsz, "/dev/i2c-%d", busNum);
    }

    if (n < 0 || n >= (int)outsz) return -1;
    return 0;
}


/* dev-laptop simulation switch
 * Default is real hardware.
 * Enable simulation when:
 *   - FEMD_SIM=1/true/yes/on
 *   - OR FEMD_SYSROOT == /tmp/sys
 */
static bool env_truthy(const char *v) {
    if (!v || v[0] == '\0') return false;
    return (!strcmp(v, "1") ||
            !strcasecmp(v, "true") ||
            !strcasecmp(v, "yes") || !strcasecmp(v, "on"));
}

static bool femd_sim_enabled(void) {
    const char *sim = getenv("FEMD_SIM");
    if (env_truthy(sim)) return true;

    const char *root = getenv(ENV_FEMD_SYSROOT);
    if (root && !strcmp(root, "/tmp/sys")) return true;

    return false;
}

static uint64_t now_ms_mono(void) {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (uint64_t)ts.tv_sec * 1000ULL + (uint64_t)ts.tv_nsec / 1000000ULL;
}

static float clampf(float v, float lo, float hi) {
    if (v < lo) return lo;
    if (v > hi) return hi;
    return v;
}

/* LM75A: generate temp (°C) and return register 0x00 raw bytes */
static int sim_lm75a_read(uint64_t tMs, int busNum, uint8_t reg, uint8_t *data, size_t len) {

    if (reg != 0x00 || len < 2) return STATUS_NOK;

    /* bus1/bus2: different phase, 35..55C sine, 120s period */
    double t = (double)tMs / 1000.0;
    double phase = (busNum == 1) ? 0.0 : 1.2;
    double tempC = 45.0 + 10.0 * sin((2.0 * M_PI * t / 120.0) + phase);
    tempC = clampf((float)tempC, 0.0f, 100.0f);

    /* LM75A raw: 9-bit in bits[15:7], 0.5C/LSB */
    int16_t temp9 = (int16_t)lround(tempC / 0.5);
    if (temp9 < 0) temp9 = (int16_t)(temp9 + 512);
    uint16_t raw = (uint16_t)((temp9 << 7) & 0xFF80);

    data[0] = (uint8_t)((raw >> 8) & 0xFF);
    data[1] = (uint8_t)(raw & 0xFF);

    return STATUS_OK;
}

/* TMP10x-style: reg 0x00 returns 12-bit temp, 0.0625C/LSB (as used by ctrl_temp_read_tmp10x) */
static int sim_tmp10x_read(uint64_t tMs, uint8_t reg, uint8_t *data, size_t len) {

    if (reg != 0x00 || len < 2) return STATUS_NOK;

    (void)tMs;
    float tempC = 42.0f; /* stable controller temp */
    int16_t code = (int16_t)lround(tempC / 0.0625f);
    int16_t raw  = (int16_t)(code << 4);

    data[0] = (uint8_t)((raw >> 8) & 0xFF);
    data[1] = (uint8_t)(raw & 0xFF);

    return STATUS_OK;
}

/* ADS1015 conversion register (0x00). We only care about which channel 
 * was selected by config writes.
 */
static int sim_ads1015_read(uint64_t tMs, int channel, uint8_t reg, uint8_t *data, size_t len) {
    if (reg != 0x00 || len < 2) return STATUS_NOK;

    double t = (double)tMs / 1000.0;
    float v  = 0.0f;

    switch (channel) {
        case 0: /* reverse power */
            v = 2.05f + 0.10f * (float)sin(2.0 * M_PI * t / 30.0);
            break;
        case 1: /* forward power */
            v = 2.25f + 0.12f * (float)sin(2.0 * M_PI * t / 25.0);
            break;
        case 2: /* PA current */
            v = 1.00f + 0.15f * (float)sin(2.0 * M_PI * t / 20.0);
            break;
        case 3: /* temperature channel (if wired) */
        default:
            v = 1.50f;
            break;
    }

    /* adc_raw_to_voltage(): ((int16_t)raw >> 4) * 4.096 / 2048
     * => code = v * 2048 / 4.096, raw = code << 4
     */
    float codeF  = (v * 2048.0f) / 4.096f;
    int16_t code = (int16_t)lround(clampf(codeF, -2048.0f, 2047.0f));
    uint16_t raw = (uint16_t)((uint16_t)code << 4);

    data[0] = (uint8_t)((raw >> 8) & 0xFF);
    data[1] = (uint8_t)(raw & 0xFF);

    return STATUS_OK;
}

int i2c_bus_init(I2cBus *bus, int busNum) {

    char devPath[I2C_DEV_PATH_MAX] = {0};

    if (!bus) {
        usys_log_error("i2c_bus_init: bus is NULL");
        return STATUS_NOK;
    }

    memset(bus, 0, sizeof(*bus));
    bus->busNum        = busNum;
    bus->fd            = -1;
    bus->currentSlave  = -1;
    bus->initialized   = false;
    bus->sim           = femd_sim_enabled();
    bus->simStartMs    = now_ms_mono();
    bus->simAdcChannel = 0;
    bus->dacKnown      = false;
    bus->dacCarrierV   = 0.0f;
    bus->dacPeakV      = 0.0f;
    bus->dacKnown      = false;
    bus->dacCarrierV   = 0.0f;
    bus->dacPeakV      = 0.0f;

    if (bus->sim) {
        bus->initialized = true;
        usys_log_info("I2C bus initialized: SIM (bus=%d)", bus->busNum);
        return STATUS_OK;
    }

    if (build_i2c_dev_path(devPath, sizeof(devPath), busNum) != 0) {
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
    bus->dacKnown     = false;
    bus->dacCarrierV  = 0.0f;
    bus->dacPeakV     = 0.0f;

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

int i2c_bus_write_reg(I2cBus *bus,
                      uint8_t devAddr,
                      uint8_t reg,
                      const uint8_t *data,
                      size_t len) {

    uint8_t tx[I2C_TX_MAX] = {0};
    ssize_t wrote = 0;

    if (!bus || !bus->initialized) return STATUS_NOK;

    /* SIM path: interpret just enough to support existing drivers. */
    if (bus->sim) {
        /* ADS1015 config write (reg 0x01): track selected channel via mux bits. */
        if (devAddr == 0x48 && reg == 0x01 && data && len == 2) {
            uint16_t cfg = (uint16_t)((data[0] << 8) | data[1]);
            uint16_t mux = (uint16_t)(cfg & 0x7000);

            if      (mux == 0x4000) bus->simAdcChannel = 0;
            else if (mux == 0x5000) bus->simAdcChannel = 1;
            else if (mux == 0x6000) bus->simAdcChannel = 2;
            else if (mux == 0x7000) bus->simAdcChannel = 3;
        }
        return STATUS_OK;
    }

    if (bus->fd < 0) return STATUS_NOK;

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

    if (!bus || !bus->initialized || !data || len == 0) return STATUS_NOK;

    /* SIM path: synthesize reads for the devices femd currently uses. */
    if (bus->sim) {
        uint64_t t = now_ms_mono() - bus->simStartMs;

        if (devAddr == 0x49) {
            return sim_lm75a_read(t, bus->busNum, reg, data, len);
        }

        if (devAddr == 0x48) {
            if (bus->busNum == 0) return sim_tmp10x_read(t, reg, data, len);
            return sim_ads1015_read(t, bus->simAdcChannel, reg, data, len);
        }

        memset(data, 0, len);
        return STATUS_OK;
    }

    if (bus->fd < 0) return STATUS_NOK;

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

    if (!bus || !bus->initialized) return STATUS_NOK;

    if (bus->sim) {
        /* In SIM we assume the fixed topology exists. */
        (void)devAddr;
        return STATUS_OK;
    }

    if (bus->fd < 0) return STATUS_NOK;

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
