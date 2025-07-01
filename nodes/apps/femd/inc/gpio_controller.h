/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#ifndef GPIO_CONTROLLER_H
#define GPIO_CONTROLLER_H

#include <stdint.h>
#include <stdbool.h>

#define STATUS_OK                 0
#define STATUS_NOK               -1

#define GPIO_PATH_MAX_LEN         256
#define GPIO_BASE_PATH           "/sys/devices/platform"

typedef enum {
    FEM_UNIT_1 = 1,
    FEM_UNIT_2 = 2
} FemUnit;

typedef enum {
    GPIO_28V_VDS,
    GPIO_TX_RF,
    GPIO_RX_RF,
    GPIO_PA_VDS,
    GPIO_TX_RFPAL,
    GPIO_PSU_PGOOD,
    GPIO_MAX
} GpioPin;

typedef struct {
    bool pa_disable;        // 28V_VDS Enable (inverted logic)
    bool tx_rf_enable;      // TX_RF Enable
    bool rx_rf_enable;      // RX_RF Enable
    bool pa_vds_enable;     // PA_VDS Enable
    bool rf_pal_enable;     // TX_RFPAL Enable
    bool pg_reg_5v;         // PSU_PGOOD (read-only)
} GpioStatus;

typedef struct {
    char *basePath;
    bool initialized;
} GpioController;

int gpio_controller_init(GpioController *controller, const char *basePath);
void gpio_controller_cleanup(GpioController *controller);

int gpio_set_28v_vds(GpioController *controller, FemUnit unit, bool enable);
int gpio_set_tx_rf(GpioController *controller, FemUnit unit, bool enable);
int gpio_set_rx_rf(GpioController *controller, FemUnit unit, bool enable);
int gpio_set_pa_vds(GpioController *controller, FemUnit unit, bool enable);
int gpio_set_tx_rfpal(GpioController *controller, FemUnit unit, bool enable);
int gpio_get_psu_pgood(GpioController *controller, FemUnit unit, bool *status);

int gpio_get_all_status(GpioController *controller, FemUnit unit, GpioStatus *status);
int gpio_disable_pa(GpioController *controller, FemUnit unit);

#endif /* GPIO_CONTROLLER_H */