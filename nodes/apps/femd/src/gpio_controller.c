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

#include "usys_api.h"
#include "usys_file.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

#include "gpio_controller.h"
#include "femd.h"

static const char* femUnitNames[3] = {
    [0]           = "",
    [FEM_UNIT_1]  = "fema1-gpios",
    [FEM_UNIT_2]  = "fema2-gpios"
};

/* Per-pin metadata: path leaf, inverted logic, readable, writable */
typedef struct {
    const char *name;
    bool inverted;
    bool readable;
    bool writable;
} gpio_meta_t;

static const gpio_meta_t gpioMeta[GPIO_MAX] = {
    [GPIO_28V_VDS]   = { "pa_disable",   true,  true,  true  },
    [GPIO_TX_RF]     = { "tx_rf_enable", false, true,  true  },
    [GPIO_RX_RF]     = { "rx_rf_enable", false, true,  true  },
    [GPIO_PA_VDS]    = { "pa_vds_enable",false, true,  true  },
    [GPIO_TX_RFPAL]  = { "rf_pal_enable",false, true,  true  },
    [GPIO_PSU_PGOOD] = { "pg_reg_5v",    false, true,  false }
};

static inline bool valid_unit(FemUnit u) { return u == FEM_UNIT_1 || u == FEM_UNIT_2; }
static inline bool valid_pin(GpioPin p)  { return p >= 0 && p < GPIO_MAX; }

static int build_path(char *dst, size_t n, const char *base, FemUnit unit, GpioPin pin) {
    if (!valid_unit(unit) || !valid_pin(pin) || !gpioMeta[pin].name)
        return STATUS_NOK;

    if (snprintf(dst, n, "%s/%s/%s", base, femUnitNames[unit], gpioMeta[pin].name) >= (int)n)
        return STATUS_NOK;

    return STATUS_OK;
}

static int write_bool_file(const char *path, bool v) {
    FILE *file = fopen(path, "w");

    if (!file) {
        usys_log_error("open(w) %s failed", path);
        return STATUS_NOK;
    }

    if (fprintf(file, "%d", v ? 1 : 0) < 0) {
        fclose(file);
        usys_log_error("write %s failed", path);
        return STATUS_NOK;
    }

    fclose(file);
    return STATUS_OK;
}

static int read_bool_file(const char *path, bool *out) {

    char buf[16]={0};
    FILE *file = fopen(path, "r");
    
    if (!file) {
        usys_log_error("open(r) %s failed", path);
        return STATUS_NOK;
    }

    if (!fgets(buf, sizeof(buf), file)) {
        fclose(file);
        usys_log_error("read %s failed", path);
        return STATUS_NOK;
    }

    fclose(file);
    *out = atoi(buf) != 0;

    return STATUS_OK;
}

int gpio_controller_init(GpioController *ctl, const char *basePath) {

    if (ctl == NULL)
        return STATUS_NOK;

    memset(ctl, 0, sizeof(*ctl));
    if (!basePath) {
        basePath = GPIO_BASE_PATH;
    }
    ctl->basePath = strdup(basePath);
    if (!ctl->basePath) return STATUS_NOK;

    /* Soft sanity checks for both FEM units */
    for (FemUnit u = FEM_UNIT_1; u <= FEM_UNIT_2; ++u) {
        char p[GPIO_PATH_MAX_LEN];
        if (snprintf(p, sizeof(p), "%s/%s", basePath, femUnitNames[u]) < (int)sizeof(p)) {
            if (access(p, F_OK) != 0) usys_log_warn("GPIO path missing: %s", p);
        }
    }

    ctl->initialized = true;
    usys_log_info("GPIO controller initialized (base=%s)", ctl->basePath);

    return STATUS_OK;
}

void gpio_controller_cleanup(GpioController *ctl) {
    
    if (!ctl) return;
    
    usys_free(ctl->basePath);
    ctl->basePath    = NULL;
    ctl->initialized = false;

    usys_log_info("GPIO controller cleaned up");
}

int gpio_set(GpioController *ctl, FemUnit unit, GpioPin pin, bool value) {

    int  rc;
    bool fileVal;
    char path[GPIO_PATH_MAX_LEN]={0};
    
    if (!ctl || !ctl->initialized) {
        usys_log_error("controller not initialized");
        return STATUS_NOK;
    }
    
    if (!valid_unit(unit) || !valid_pin(pin)) {
        usys_log_error("bad unit/pin");
        return STATUS_NOK;
    }
    
    if (!gpioMeta[pin].writable) {
        usys_log_error("pin %d not writable", pin);
        return STATUS_NOK;
    }

    fileVal = gpioMeta[pin].inverted ? !value : value;
    if (build_path(path, sizeof(path), ctl->basePath, unit, pin) != STATUS_OK) {
        return STATUS_NOK;
    }

    rc = write_bool_file(path, fileVal);
    if (rc == STATUS_OK) {
        usys_log_debug("set %s := %d (logical=%d)", path, (int)fileVal, (int)value);
    }

    return rc;
}

int gpio_get(GpioController *ctl, FemUnit unit, GpioPin pin, bool *out) {

    int  rc;
    bool fileVal;
    char path[GPIO_PATH_MAX_LEN]={0};
    
    if (!ctl || !ctl->initialized || !out) {
        return STATUS_NOK;
    }
    
    if (!valid_unit(unit) || !valid_pin(pin)) {
        usys_log_error("bad unit/pin");
        return STATUS_NOK;
    }

    if (!gpioMeta[pin].readable) {
        usys_log_error("pin %d not readable", pin);
        return STATUS_NOK;
    }

    if (build_path(path, sizeof(path), ctl->basePath, unit, pin) != STATUS_OK) {
        return STATUS_NOK;
    }

    rc = read_bool_file(path, &fileVal);
    if (rc != STATUS_OK) {
        return rc;
    }

    *out = gpioMeta[pin].inverted ? !fileVal : fileVal;
    usys_log_debug("get %s -> %d (logical=%d)", path, (int)fileVal, (int)*out);

    return STATUS_OK;
}

int gpio_read_all(GpioController *ctl, FemUnit unit, GpioStatus *out) {

    bool val;

    if (!out) {
        return STATUS_NOK;
    }

    if (gpio_get(ctl, unit, GPIO_28V_VDS, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->pa_disable = val;

    if (gpio_get(ctl, unit, GPIO_TX_RF, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->tx_rf_enable = val;

    if (gpio_get(ctl, unit, GPIO_RX_RF, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->rx_rf_enable = val;

    if (gpio_get(ctl, unit, GPIO_PA_VDS, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->pa_vds_enable= val;

    if (gpio_get(ctl, unit, GPIO_TX_RFPAL, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->rf_pal_enable= val;

    if (gpio_get(ctl, unit, GPIO_PSU_PGOOD, &val) != STATUS_OK) {
        return STATUS_NOK;
    }
    out->pg_reg_5v = val;

    usys_log_debug("FEM%d status: 28V_EN=%d TX=%d RX=%d PA_VDS=%d PAL=%d PGOOD=%d",
                   unit, out->pa_disable, out->tx_rf_enable, out->rx_rf_enable,
                   out->pa_vds_enable, out->rf_pal_enable, out->pg_reg_5v);

    return STATUS_OK;
}

int gpio_apply(GpioController *ctl, FemUnit unit, const GpioStatus *desired) {

    if (!desired) {
        return STATUS_NOK;
    }

    if (gpio_set(ctl, unit, GPIO_28V_VDS,  desired->pa_disable)    != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RF,    desired->tx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_RX_RF,    desired->rx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_PA_VDS,   desired->pa_vds_enable) != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RFPAL, desired->rf_pal_enable) != STATUS_OK) return STATUS_NOK;

    return STATUS_OK;
}

int gpio_disable_pa(GpioController *ctl, FemUnit unit) {

    if (!ctl || !ctl->initialized) {
        usys_log_error("controller not initialized");
        return STATUS_NOK;
    }

    usys_log_warn("Emergency PA disable for FEM%d", unit);
    if (gpio_set(ctl, unit, GPIO_PA_VDS,  false) != STATUS_OK) {
        return STATUS_NOK;
    }

    if (gpio_set(ctl, unit, GPIO_28V_VDS, false) != STATUS_OK) {
        return STATUS_NOK;
    }

    usys_log_info("PA disabled for FEM%d", unit);
    return STATUS_OK;
}
