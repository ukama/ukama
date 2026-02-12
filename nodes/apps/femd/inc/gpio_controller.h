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

#include "femd.h"

#define GPIO_PATH_MAX_LEN 128

typedef struct {
    bool tx_rf_enable;
    bool rx_rf_enable;
    bool pa_vds_enable;
    bool rf_pal_enable;
    bool pa_disable;     /* inverted logic in sysfs */
    bool psu_pgood;

} GpioStatus;

typedef struct {
    char txRfEnable[GPIO_PATH_MAX_LEN];
    char rxRfEnable[GPIO_PATH_MAX_LEN];
    char paVdsEnable[GPIO_PATH_MAX_LEN];
    char rfPalEnable[GPIO_PATH_MAX_LEN];
    char vds28Enable[GPIO_PATH_MAX_LEN];
    char psuPgood[GPIO_PATH_MAX_LEN];
} GpioPaths;

typedef struct {
    char     basePath[GPIO_PATH_MAX_LEN];
    GpioPaths fem[3];
    bool     initialized;
} GpioController;

int  gpio_controller_init(GpioController *ctrl, const char *gpioBasePath);
void gpio_controller_cleanup(GpioController *ctrl);

int  gpio_read_all(GpioController *ctrl, FemUnit unit, GpioStatus *out);
int  gpio_apply(GpioController *ctrl, FemUnit unit, const GpioStatus *desired);

int  gpio_disable_pa(GpioController *ctrl, FemUnit unit);

#endif /* GPIO_CONTROLLER_H */
