/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef GPIO_CONTROLLER_H
#define GPIO_CONTROLLER_H

#include <stdbool.h>
#include <stdint.h>
#include <pthread.h>

#include "femd.h"

#define GPIO_BASE_PATH "/sys/class/gpio"

typedef enum {
    GPIO_28V_VDS = 0,
    GPIO_TX_RF,
    GPIO_RX_RF,
    GPIO_PA_VDS,
    GPIO_TX_RFPAL,
    GPIO_PSU_PGOOD,
    GPIO_MAX
} GpioPin;

typedef struct {
    bool pa_disable;
    bool tx_rf_enable;
    bool rx_rf_enable;
    bool pa_vds_enable;
    bool rf_pal_enable;
    bool pg_reg_5v;
    bool psu_pgood;
} GpioStatus;

static inline int gpio_vds_28v_enabled(const GpioStatus *s) {
    if (!s) return 0;
    return s->pa_disable ? 0 : 1;
}

typedef struct {
    pthread_mutex_t mu;
    char           *basePath;
    int             initialized;
} GpioController;

int  gpio_controller_init(GpioController *ctl, const char *basePath);
void gpio_controller_cleanup(GpioController *ctl);

int gpio_get(GpioController *ctl, FemUnit unit, GpioPin pin, bool *out);
int gpio_set(GpioController *ctl, FemUnit unit, GpioPin pin, bool value);

int gpio_read_all(GpioController *ctl, FemUnit unit, GpioStatus *out);
int gpio_apply(GpioController *ctl, FemUnit unit, const GpioStatus *desired);

#endif /* GPIO_CONTROLLER_H */
