/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>
#include <unistd.h>

#include "gpio_controller.h"
#include "usys_log.h"
#include "usys_types.h"

typedef struct {
    const char *name;
    bool        isOutput;
} gpio_meta_t;

static const gpio_meta_t gpioMeta[GPIO_MAX] = {
    [GPIO_28V_VDS]   = { "pa_disable",    true  },
    [GPIO_TX_RF]     = { "tx_rf_enable",  true  },
    [GPIO_RX_RF]     = { "rx_rf_enable",  true  },
    [GPIO_PA_VDS]    = { "pa_vds_enable", true  },
    [GPIO_TX_RFPAL]  = { "rf_pal_enable", true  },
    [GPIO_PSU_PGOOD] = { "pg_reg_5v",     false }
};

static int femd_sim_write_default_gpio_file(const char *valuePath, const char *gpioName);

static inline bool valid_pin(GpioPin p) {
    return (p >= 0 && p < GPIO_MAX);
}

static int build_path(char *dst, size_t n, const char *base, FemUnit unit, GpioPin pin, const char *leaf) {

    if (!dst || n == 0 || !base || !leaf) return STATUS_NOK;
    if (!valid_pin(pin)) return STATUS_NOK;
    if (unit != FEM_UNIT_1 && unit != FEM_UNIT_2) return STATUS_NOK;

    if (snprintf(dst, n, "%s/fem%d/%s/%s", base, (unit == FEM_UNIT_1) ? 1 : 2, gpioMeta[pin].name, leaf) >= (int)n) {
        return STATUS_NOK;
    }

    return STATUS_OK;
}

static int write_bool_file(const char *path, bool v) {

    FILE *f;

    if (!path) return STATUS_NOK;

    f = fopen(path, "w");
    if (!f) return STATUS_NOK;

    if (fprintf(f, "%d\n", v ? 1 : 0) < 0) {
        fclose(f);
        return STATUS_NOK;
    }

    fclose(f);
    return STATUS_OK;
}

static int read_bool_file(const char *path, bool *out) {

    FILE *f;
    int v = 0;

    if (!path || !out) return STATUS_NOK;

    f = fopen(path, "r");
    if (!f) return STATUS_NOK;

    if (fscanf(f, "%d", &v) != 1) {
        fclose(f);
        return STATUS_NOK;
    }

    fclose(f);
    *out = (v != 0);
    return STATUS_OK;
}

/* Dev-laptop helper: create /tmp/sys GPIO files on-demand so femd can run
+ * without any pre-created mock tree. No-op on real sysfs.
+ */
static int ensure_dir(const char *path) {
    if (!path || path[0] == '\0') return STATUS_NOK;
    if (mkdir(path, 0755) == 0) return STATUS_OK;
    if (errno == EEXIST) return STATUS_OK;
    return STATUS_NOK;
}

static int ensure_file_with_default(const char *path, int def) {
    FILE *f;
    if (!path) return STATUS_NOK;
    if (access(path, F_OK) == 0) return STATUS_OK;
    f = fopen(path, "w");
    if (!f) return STATUS_NOK;
    fprintf(f, "%d\n", def);
    fclose(f);
    return STATUS_OK;
}

static void maybe_init_mock_tree(const char *basePath) {
    char p[256];

    if (!basePath) return;
    /* Only ever create under /tmp/sys (avoid touching real sysfs). */
    if (strncmp(basePath, "/tmp/sys", 8) != 0) return;

    /* Ensure parent dirs exist (non-recursive, but good enough for /tmp/sys/... path). */
    (void)ensure_dir("/tmp/sys");
    (void)ensure_dir("/tmp/sys/devices");
    (void)ensure_dir(basePath);

    /* base/fem{1,2}/{pin}/value */
    for (int u = 1; u <= 2; u++) {
        snprintf(p, sizeof(p), "%s/fem%d", basePath, u);
        (void)ensure_dir(p);
        for (int pin = 0; pin < GPIO_MAX; pin++) {
            snprintf(p, sizeof(p), "%s/fem%d/%s", basePath, u, gpioMeta[pin].name);
            (void)ensure_dir(p);
            snprintf(p, sizeof(p), "%s/fem%d/%s/value", basePath, u, gpioMeta[pin].name);
            (void)femd_sim_write_default_gpio_file(p, gpioMeta[pin].name);
        }
    }
}

static int femd_sim_default_gpio_value(const char *gpioName) {
    if (!gpioName) return 0;

    /* Keep board "healthy" in SIM */
    if (!strcmp(gpioName, "psu_pgood")) return 1;
    if (!strcmp(gpioName, "pgood")) return 1;

    /* RF chain defaults in SIM so PA metrics make sense */
    if (!strcmp(gpioName, "pa_enable")) return 1;
    if (!strcmp(gpioName, "tx_enable")) return 1;

    /* Optional: if you have these signals and want them ON by default */
    if (!strcmp(gpioName, "rf_enable")) return 1;
    if (!strcmp(gpioName, "trx_enable")) return 1;

    return 0;
}

static int femd_sim_write_default_gpio_file(const char *valuePath, const char *gpioName) {
    int def = femd_sim_default_gpio_value(gpioName);
    return ensure_file_with_default(valuePath, def);
}


int gpio_controller_init(GpioController *ctl, const char *basePath) {

    if (!ctl) return STATUS_NOK;

    memset(ctl, 0, sizeof(*ctl));
    pthread_mutex_init(&ctl->mu, NULL);

    if (!basePath) basePath = GPIO_BASE_PATH;

    ctl->basePath = strdup(basePath);
    if (!ctl->basePath) {
        pthread_mutex_destroy(&ctl->mu);
        return STATUS_NOK;
    }

    ctl->initialized = 1;

    /* If we're pointing at /tmp/sys..., create the minimal mock tree. */
    maybe_init_mock_tree(ctl->basePath);

    return STATUS_OK;
}

void gpio_controller_cleanup(GpioController *ctl) {

    if (!ctl) return;

    pthread_mutex_lock(&ctl->mu);

    if (ctl->basePath) free(ctl->basePath);
    ctl->basePath = NULL;
    ctl->initialized = 0;

    pthread_mutex_unlock(&ctl->mu);
    pthread_mutex_destroy(&ctl->mu);

    memset(ctl, 0, sizeof(*ctl));
}

int gpio_set(GpioController *ctl, FemUnit unit, GpioPin pin, bool value) {

    char path[256];

    if (!ctl || !ctl->initialized) return STATUS_NOK;
    if (!valid_pin(pin)) return STATUS_NOK;
    if (!gpioMeta[pin].isOutput) return STATUS_NOK;

    pthread_mutex_lock(&ctl->mu);

    if (build_path(path, sizeof(path), ctl->basePath, unit, pin, "value") != STATUS_OK) {
        pthread_mutex_unlock(&ctl->mu);
        return STATUS_NOK;
    }

    if (write_bool_file(path, value) != STATUS_OK) {
        pthread_mutex_unlock(&ctl->mu);
        return STATUS_NOK;
    }

    pthread_mutex_unlock(&ctl->mu);
    return STATUS_OK;
}

int gpio_get(GpioController *ctl, FemUnit unit, GpioPin pin, bool *out) {

    char path[256];

    if (!ctl || !ctl->initialized || !out) return STATUS_NOK;
    if (!valid_pin(pin)) return STATUS_NOK;

    pthread_mutex_lock(&ctl->mu);

    if (build_path(path, sizeof(path), ctl->basePath, unit, pin, "value") != STATUS_OK) {
        pthread_mutex_unlock(&ctl->mu);
        return STATUS_NOK;
    }

    if (read_bool_file(path, out) != STATUS_OK) {
        pthread_mutex_unlock(&ctl->mu);
        return STATUS_NOK;
    }

    pthread_mutex_unlock(&ctl->mu);
    return STATUS_OK;
}

int gpio_read_all(GpioController *ctl, FemUnit unit, GpioStatus *out) {

    bool v;

    if (!ctl || !out) return STATUS_NOK;

    memset(out, 0, sizeof(*out));

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

    usys_log_debug("GPIO FEM%d: pa_disable=%d tx=%d rx=%d pa_vds=%d rf_pal=%d pgood=%d",
                   (unit == FEM_UNIT_1) ? 1 : 2,
                   (int)out->pa_disable,
                   (int)out->tx_rf_enable,
                   (int)out->rx_rf_enable,
                   (int)out->pa_vds_enable,
                   (int)out->rf_pal_enable,
                   (int)out->pg_reg_5v);

    return STATUS_OK;
}

int gpio_apply(GpioController *ctl, FemUnit unit, const GpioStatus *desired) {

    if (!ctl || !desired) return STATUS_NOK;

    if (gpio_set(ctl, unit, GPIO_28V_VDS,  desired->pa_disable)    != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RF,    desired->tx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_RX_RF,    desired->rx_rf_enable)  != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_PA_VDS,   desired->pa_vds_enable) != STATUS_OK) return STATUS_NOK;
    if (gpio_set(ctl, unit, GPIO_TX_RFPAL, desired->rf_pal_enable) != STATUS_OK) return STATUS_NOK;

    return STATUS_OK;
}
