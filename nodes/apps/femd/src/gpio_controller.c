/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <unistd.h>

#include "usys_log.h"
#include "usys_types.h"

#include "gpio_controller.h"

#define GPIO_DEFAULT_BASE "/sys/devices/platform"

static inline int valid_unit(FemUnit unit) {
    return unit == FEM_UNIT_1 || unit == FEM_UNIT_2;
}

static int write_bool_file(const char *path, bool v) {

    FILE *f = fopen(path, "w");
    if (!f) {
        usys_log_error("open(w) %s failed", path);
        return STATUS_NOK;
    }

    if (fprintf(f, "%d", v ? 1 : 0) < 0) {
        fclose(f);
        usys_log_error("write %s failed", path);
        return STATUS_NOK;
    }

    fclose(f);
    return STATUS_OK;
}

static int read_bool_file(const char *path, bool *out) {

    char buf[16];
    FILE *f;

    if (!out) return STATUS_NOK;

    memset(buf, 0, sizeof(buf));

    f = fopen(path, "r");
    if (!f) {
        usys_log_error("open(r) %s failed", path);
        return STATUS_NOK;
    }

    if (!fgets(buf, sizeof(buf), f)) {
        fclose(f);
        usys_log_error("read %s failed", path);
        return STATUS_NOK;
    }

    fclose(f);
    *out = (atoi(buf) != 0);
    return STATUS_OK;
}

static void build_paths_for_unit(GpioController *ctrl, FemUnit unit, const char *unitDir) {

    GpioPaths *p;

    p = &ctrl->fem[unit];

    (void)snprintf(p->txRfEnable, sizeof(p->txRfEnable), "%s/%s/tx_rf_enable", ctrl->basePath, unitDir);
    (void)snprintf(p->rxRfEnable, sizeof(p->rxRfEnable), "%s/%s/rx_rf_enable", ctrl->basePath, unitDir);
    (void)snprintf(p->paVdsEnable, sizeof(p->paVdsEnable), "%s/%s/pa_vds_enable", ctrl->basePath, unitDir);
    (void)snprintf(p->rfPalEnable, sizeof(p->rfPalEnable), "%s/%s/rf_pal_enable", ctrl->basePath, unitDir);

    (void)snprintf(p->vds28Enable, sizeof(p->vds28Enable), "%s/%s/pa_disable", ctrl->basePath, unitDir);

    (void)snprintf(p->psuPgood, sizeof(p->psuPgood), "%s/%s/pg_reg_5v", ctrl->basePath, unitDir);
}

int gpio_controller_init(GpioController *ctrl, const char *gpioBasePath) {

    const char *base;

    if (!ctrl) return STATUS_NOK;

    memset(ctrl, 0, sizeof(*ctrl));

    base = gpioBasePath ? gpioBasePath : GPIO_DEFAULT_BASE;
    (void)snprintf(ctrl->basePath, sizeof(ctrl->basePath), "%s", base);

    build_paths_for_unit(ctrl, FEM_UNIT_1, "fema1-gpios");
    build_paths_for_unit(ctrl, FEM_UNIT_2, "fema2-gpios");

    if (access(ctrl->fem[FEM_UNIT_1].txRfEnable, F_OK) != 0) {
        usys_log_warn("GPIO path missing (FEM1): %s", ctrl->fem[FEM_UNIT_1].txRfEnable);
    }
    if (access(ctrl->fem[FEM_UNIT_2].txRfEnable, F_OK) != 0) {
        usys_log_warn("GPIO path missing (FEM2): %s", ctrl->fem[FEM_UNIT_2].txRfEnable);
    }

    ctrl->initialized = true;
    usys_log_info("GPIO controller initialized (base=%s)", ctrl->basePath);

    return STATUS_OK;
}

void gpio_controller_cleanup(GpioController *ctrl) {

    if (!ctrl) return;

    memset(ctrl, 0, sizeof(*ctrl));
}

int gpio_read_all(GpioController *ctrl, FemUnit unit, GpioStatus *out) {

    bool fileVal;

    if (!ctrl || !ctrl->initialized || !out) return STATUS_NOK;
    if (!valid_unit(unit)) return STATUS_NOK;

    memset(out, 0, sizeof(*out));

    if (read_bool_file(ctrl->fem[unit].vds28Enable, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->pa_disable = !fileVal;

    if (read_bool_file(ctrl->fem[unit].txRfEnable, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->tx_rf_enable = fileVal;

    if (read_bool_file(ctrl->fem[unit].rxRfEnable, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->rx_rf_enable = fileVal;

    if (read_bool_file(ctrl->fem[unit].paVdsEnable, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->pa_vds_enable = fileVal;

    if (read_bool_file(ctrl->fem[unit].rfPalEnable, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->rf_pal_enable = fileVal;

    if (read_bool_file(ctrl->fem[unit].psuPgood, &fileVal) != STATUS_OK) return STATUS_NOK;
    out->psu_pgood = fileVal;

    usys_log_debug("FEM%d gpio: 28V_EN=%d TX=%d RX=%d PA_VDS=%d PAL=%d PGOOD=%d",
                   unit,
                   (int)out->pa_disable,
                   (int)out->tx_rf_enable,
                   (int)out->rx_rf_enable,
                   (int)out->pa_vds_enable,
                   (int)out->rf_pal_enable,
                   (int)out->psu_pgood);

    return STATUS_OK;
}

int gpio_apply(GpioController *ctrl, FemUnit unit, const GpioStatus *desired) {

    bool fileVal;

    if (!ctrl || !ctrl->initialized || !desired) return STATUS_NOK;
    if (!valid_unit(unit)) return STATUS_NOK;

    fileVal = !desired->pa_disable;
    if (write_bool_file(ctrl->fem[unit].vds28Enable, fileVal) != STATUS_OK) return STATUS_NOK;

    if (write_bool_file(ctrl->fem[unit].txRfEnable, desired->tx_rf_enable) != STATUS_OK) return STATUS_NOK;
    if (write_bool_file(ctrl->fem[unit].rxRfEnable, desired->rx_rf_enable) != STATUS_OK) return STATUS_NOK;
    if (write_bool_file(ctrl->fem[unit].paVdsEnable, desired->pa_vds_enable) != STATUS_OK) return STATUS_NOK;
    if (write_bool_file(ctrl->fem[unit].rfPalEnable, desired->rf_pal_enable) != STATUS_OK) return STATUS_NOK;

    return STATUS_OK;
}

int gpio_disable_pa(GpioController *ctrl, FemUnit unit) {

    GpioStatus s;

    if (!ctrl || !ctrl->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    if (!valid_unit(unit)) return STATUS_NOK;

    memset(&s, 0, sizeof(s));

    s.pa_disable = false;
    s.pa_vds_enable = false;
    s.tx_rf_enable = false;
    s.rx_rf_enable = false;
    s.rf_pal_enable = false;

    usys_log_warn("Emergency PA disable for FEM%d", unit);

    if (gpio_apply(ctrl, unit, &s) != STATUS_OK) {
        return STATUS_NOK;
    }

    usys_log_info("PA disabled for FEM%d", unit);

    return STATUS_OK;
}
