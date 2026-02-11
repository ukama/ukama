/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2025-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_types.h"

#include "gpio_controller.h"

static const char *femUnitNames[3] = {
    [0]          = "",
    [FEM_UNIT_1] = "fema1-gpios",
    [FEM_UNIT_2] = "fema2-gpios"
};

typedef struct {
    const char *name;
    bool inverted;
    bool readable;
    bool writable;
} gpio_meta_t;

static const gpio_meta_t gpioMeta[GPIO_MAX] = {
    [GPIO_28V_VDS]   = { "pa_disable",    true,  true,  true  },
    [GPIO_TX_RF]     = { "tx_rf_enable",  false, true,  true  },
    [GPIO_RX_RF]     = { "rx_rf_enable",  false, true,  true  },
    [GPIO_PA_VDS]    = { "pa_vds_enable", false, true,  true  },
    [GPIO_TX_RFPAL]  = { "rf_pal_enable", false, true,  true  },
    [GPIO_PSU_PGOOD] = { "pg_reg_5v",     false, true,  false }
};

static inline bool valid_unit(FemUnit u) { return u == FEM_UNIT_1 || u == FEM_UNIT_2; }
static inline bool valid_pin(GpioPin p)  { return p >= 0 && p < GPIO_MAX; }

static int build_path(char *dst, size_t n, const char *base, FemUnit unit, GpioPin pin) {

    if (!dst || !base) return STATUS_NOK;
    if (!valid_unit(unit) || !valid_pin(pin) || !gpioMeta[pin].name) return STATUS_NOK;

    if (snprintf(dst, n, "%s/%s/%s", base, femUnitNames[unit], gpioMeta[pin].name) >= (int)n) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

static int write_bool_file(const char *path, bool v) {

    FILE *f;

    if (!path) return STATUS_NOK;

    f = fopen(path, "w");
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

    FILE *f;
    char buf[16];

    if (!path || !out) return STATUS_NOK;

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

int gpio_controller_init(GpioController *ctl, const char *basePath) {

    if (!ctl) return STATUS_NOK;

    memset(ctl, 0, sizeof(*ctl));

    if (!basePath) basePath = GPIO_BASE_PATH;

    ctl->basePath = strdup(basePath);
    if (!ctl->basePath) return STATUS_NOK;

    for (FemUnit u = FEM_UNIT_1; u <= FEM_UNIT_2; u++) {
        char p[GPIO_PATH_MAX_LEN];
        if (snprintf(p, sizeof(p), "%s/%s", ctl->basePath, femUnitNames[u]) < (int)sizeof(p)) {
            if (access(p, F_OK) != 0) usys_log_warn("GPIO path missing: %s", p);
        }
    }

    ctl->initialized = true;
    usys_log_info("GPIO controller initialized (base=%s)", ctl->basePath);

    return STATUS_OK;
}

void gpio_controller_cleanup(GpioController *ctl) {

    if (!ctl) return;

    if (ctl->basePath) {
        usys_free(ctl->basePath);
        ctl->basePath = NULL;
    }

    ctl->initialized = false;
    usys_log_info("GPIO controller cleaned up");
}

int gpio_set(GpioController *ctl, FemUnit unit, GpioPin pin, bool value) {

    char path[GPIO_PATH_MAX_LEN];
    bool fileVal;

    if (!ctl || !ctl->initialized) return STATUS_NOK;
    if (!valid_unit(unit) || !valid_pin(pin)) return STATUS_NOK;
    if (!gpioMeta[pin].writable) return STATUS_NOK;

    if (build_path(path, sizeof(path), ctl->basePath, unit, pin) != STATUS_OK) return STATUS_NOK;

    fileVal = gpioMeta[pin].inverted ? !value : value;

    if (write_bool_file(path, fileVal) != STATUS_OK) return STATUS_NOK;

    usys_log_debug("gpio set fem=%d pin=%d path=%s file=%d logical=%d",
                   unit, pin, path, (int)fileVal, (int)value);

    return STATUS_OK;
}

int gpio_get(GpioController *ctl, FemUnit unit, GpioPin pin, bool *out) {

    char path[GPIO_PATH_MAX_LEN];
    bool fileVal;

    if (!ctl || !ctl->initialized || !out) return STATUS_NOK;
    if (!valid_unit(unit) || !valid_pin(pin)) return STATUS_NOK;
    if (!gpioMeta[pin].readable) return STATUS_NOK;

    if (build_path(path, sizeof(path), ctl->basePath, unit, pin) != STATUS_OK) return STATUS_NOK;

    if (read_bool_file(path, &fileVal) != STATUS_OK) return STATUS_NOK;

    *out = gpioMeta[pin].inverted ? !fileVal : fileVal;

    usys_log_debug("gpio get fem=%d pin=%d path=%s file=%d logical=%d",
                   unit, pin, path, (int)fileVal, (int)(*out));

    return STATUS_OK;
}

int gpio_read_all(GpioController *ctl, FemUnit unit, GpioStatus *out) {

    bool v;

    if (!out) return STATUS_NOK;

    if (gpio_get(ctl, unit, GPIO_28V_VDS, &v) != STATUS_OK) return STATUS_NOK;
    out->pa_disable = v;

    if (gpio_get(ctl, unit, GPIO_TX_RF, &v) != STATUS_OK) return STATUS_NOK;
    out->tx_rf_enable = v;

    if (gpio_get(ctl, unit, GPIO_RX_RF, &v) != STATUS_OK) return STATUS_NOK;
    out->rx_rf_enable = v;

    if (gpio_get(ctl, unit, GPIO_PA_VDS, &v) != STATUS_OK) return STATUS_NOK;
    out->pa_vds_enable = v;

    if (gpio_get(ctl, unit, GPIO_TX_RFPAL, &v) != STATUS_OK) return STATUS_NOK;
    out->rf_pal_enable = v;

    if (gpio_get(ctl, unit, GPIO_PSU_PGOOD, &v) != STATUS_OK) return STATUS_NOK;
    out->pg_reg_5v = v;

    usys_log_debug("gpio all fem=%d 28v_en=%d tx=%d rx=%d pa_vds=%d pal=%d pgood=%d",
                   unit,
                   (int)gpio_vds_28v_enabled(out),
                   (int)out->tx_rf_enable,
                   (int)out->rx_rf_enable,
                   (int)out->pa_vds_enable,
                   (int)out->rf_pal_enable,
                   (int)out->pg_reg_5v);

    return STATUS_OK;
}

int gpio_apply(GpioController *ctl, FemUnit unit, const GpioStatus *desired) {

    if (!desired) return STATUS_NOK;

    if (gpio_set(ctl, unit, GPIO_28V_VDS,  desired->pa_disable)    != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RF,    desired->tx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_RX_RF,    desired->rx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_PA_VDS,   desired->pa_vds_enable) != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RFPAL, desired->rf_pal_enable) != STATUS_OK) return STATUS_NOK;

    return STATUS_OK;
}

int gpio_disable_pa(GpioController *ctl, FemUnit unit) {

    if (!ctl || !ctl->initialized) return STATUS_NOK;

    usys_log_warn("gpio emergency pa disable fem=%d", unit);

    if (gpio_set(ctl, unit, GPIO_PA_VDS,  false) != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_28V_VDS, false) != STATUS_OK) return STATUS_NOK;

    return STATUS_OK;
}
