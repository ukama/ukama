/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include "gpio_controller.h"
#include "femd.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

static const char* gpio_pin_names[GPIO_MAX] = {
    [GPIO_28V_VDS]   = "pa_disable",
    [GPIO_TX_RF]     = "tx_rf_enable", 
    [GPIO_RX_RF]     = "rx_rf_enable",
    [GPIO_PA_VDS]    = "pa_vds_enable",
    [GPIO_TX_RFPAL]  = "rf_pal_enable",
    [GPIO_PSU_PGOOD] = "pg_reg_5v"
};

static const char* fem_unit_names[3] = {
    [0] = "",
    [FEM_UNIT_1] = "fema1-gpios",
    [FEM_UNIT_2] = "fema2-gpios"
};

static int gpio_write_value(const char *basePath, FemUnit unit, GpioPin pin, bool value) {
    char path[GPIO_PATH_MAX_LEN];
    char valueStr[8];
    FILE *file;
    
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        usys_log_error("Invalid FEM unit: %d", unit);
        return STATUS_NOK;
    }
    
    if (pin >= GPIO_MAX) {
        usys_log_error("Invalid GPIO pin: %d", pin);
        return STATUS_NOK;
    }
    
    snprintf(path, sizeof(path), "%s/%s/%s", 
             basePath, fem_unit_names[unit], gpio_pin_names[pin]);
    
    snprintf(valueStr, sizeof(valueStr), "%d", value ? 1 : 0);
    
    file = fopen(path, "w");
    if (!file) {
        usys_log_error("Failed to open GPIO %s for writing", path);
        return STATUS_NOK;
    }
    
    if (fprintf(file, "%s", valueStr) < 0) {
        usys_log_error("Failed to write value %s to GPIO %s", valueStr, path);
        fclose(file);
        return STATUS_NOK;
    }
    
    fclose(file);
    usys_log_debug("GPIO %s set to %s", path, valueStr);
    return STATUS_OK;
}

static int gpio_read_value(const char *basePath, FemUnit unit, GpioPin pin, bool *value) {
    char path[GPIO_PATH_MAX_LEN];
    char buffer[16];
    FILE *file;
    
    if (unit < FEM_UNIT_1 || unit > FEM_UNIT_2) {
        usys_log_error("Invalid FEM unit: %d", unit);
        return STATUS_NOK;
    }
    
    if (pin >= GPIO_MAX) {
        usys_log_error("Invalid GPIO pin: %d", pin);
        return STATUS_NOK;
    }
    
    if (!value) {
        usys_log_error("Null value pointer");
        return STATUS_NOK;
    }
    
    snprintf(path, sizeof(path), "%s/%s/%s", 
             basePath, fem_unit_names[unit], gpio_pin_names[pin]);
    
    file = fopen(path, "r");
    if (!file) {
        usys_log_error("Failed to open GPIO %s for reading", path);
        return STATUS_NOK;
    }
    
    if (!fgets(buffer, sizeof(buffer), file)) {
        usys_log_error("Failed to read from GPIO %s", path);
        fclose(file);
        return STATUS_NOK;
    }
    
    fclose(file);
    
    int intValue = atoi(buffer);
    *value = (intValue != 0);
    
    usys_log_debug("GPIO %s read value: %d", path, intValue);
    return STATUS_OK;
}

int gpio_controller_init(GpioController *controller, const char *basePath) {
    if (!controller) {
        usys_log_error("Null controller pointer");
        return STATUS_NOK;
    }
    
    if (!basePath) {
        basePath = GPIO_BASE_PATH;
    }
    
    memset(controller, 0, sizeof(GpioController));
    controller->basePath = strdup(basePath);
    if (!controller->basePath) {
        usys_log_error("Failed to allocate memory for GPIO base path");
        return STATUS_NOK;
    }
    
    char testPath[GPIO_PATH_MAX_LEN];
    snprintf(testPath, sizeof(testPath), "%s/%s", basePath, fem_unit_names[FEM_UNIT_1]);
    if (access(testPath, F_OK) != 0) {
        usys_log_warn("GPIO path %s does not exist", testPath);
    }
    
    snprintf(testPath, sizeof(testPath), "%s/%s", basePath, fem_unit_names[FEM_UNIT_2]);
    if (access(testPath, F_OK) != 0) {
        usys_log_warn("GPIO path %s does not exist", testPath);
    }
    
    controller->initialized = true;
    usys_log_info("GPIO controller initialized with base path: %s", basePath);
    
    return STATUS_OK;
}

void gpio_controller_cleanup(GpioController *controller) {
    if (!controller) {
        return;
    }
    
    if (controller->basePath) {
        free(controller->basePath);
        controller->basePath = NULL;
    }
    
    controller->initialized = false;
    usys_log_info("GPIO controller cleanup completed");
}

int gpio_set_28v_vds(GpioController *controller, FemUnit unit, bool enable) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_write_value(controller->basePath, unit, GPIO_28V_VDS, !enable);
}

int gpio_set_tx_rf(GpioController *controller, FemUnit unit, bool enable) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_write_value(controller->basePath, unit, GPIO_TX_RF, enable);
}

int gpio_set_rx_rf(GpioController *controller, FemUnit unit, bool enable) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_write_value(controller->basePath, unit, GPIO_RX_RF, enable);
}

int gpio_set_pa_vds(GpioController *controller, FemUnit unit, bool enable) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_write_value(controller->basePath, unit, GPIO_PA_VDS, enable);
}

int gpio_set_tx_rfpal(GpioController *controller, FemUnit unit, bool enable) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_write_value(controller->basePath, unit, GPIO_TX_RFPAL, enable);
}

int gpio_get_psu_pgood(GpioController *controller, FemUnit unit, bool *status) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    return gpio_read_value(controller->basePath, unit, GPIO_PSU_PGOOD, status);
}

int gpio_get_all_status(GpioController *controller, FemUnit unit, GpioStatus *status) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    if (!status) {
        usys_log_error("Null status pointer");
        return STATUS_NOK;
    }
    
    bool pa_disable_raw;
    if (gpio_read_value(controller->basePath, unit, GPIO_28V_VDS, &pa_disable_raw) != STATUS_OK) {
        usys_log_error("Failed to read 28V_VDS status");
        return STATUS_NOK;
    }
    status->pa_disable = !pa_disable_raw;
    
    if (gpio_read_value(controller->basePath, unit, GPIO_TX_RF, &status->tx_rf_enable) != STATUS_OK) {
        usys_log_error("Failed to read TX_RF status");
        return STATUS_NOK;
    }
    
    if (gpio_read_value(controller->basePath, unit, GPIO_RX_RF, &status->rx_rf_enable) != STATUS_OK) {
        usys_log_error("Failed to read RX_RF status");
        return STATUS_NOK;
    }
    
    if (gpio_read_value(controller->basePath, unit, GPIO_PA_VDS, &status->pa_vds_enable) != STATUS_OK) {
        usys_log_error("Failed to read PA_VDS status");
        return STATUS_NOK;
    }
    
    if (gpio_read_value(controller->basePath, unit, GPIO_TX_RFPAL, &status->rf_pal_enable) != STATUS_OK) {
        usys_log_error("Failed to read TX_RFPAL status");
        return STATUS_NOK;
    }
    
    if (gpio_read_value(controller->basePath, unit, GPIO_PSU_PGOOD, &status->pg_reg_5v) != STATUS_OK) {
        usys_log_error("Failed to read PSU_PGOOD status");
        return STATUS_NOK;
    }
    
    usys_log_debug("GPIO status for FEM%d: 28V_VDS=%d, TX_RF=%d, RX_RF=%d, PA_VDS=%d, TX_RFPAL=%d, PSU_PGOOD=%d",
           unit, status->pa_disable, status->tx_rf_enable, status->rx_rf_enable,
           status->pa_vds_enable, status->rf_pal_enable, status->pg_reg_5v);
    
    return STATUS_OK;
}

int gpio_disable_pa(GpioController *controller, FemUnit unit) {
    if (!controller || !controller->initialized) {
        usys_log_error("GPIO controller not initialized");
        return STATUS_NOK;
    }
    
    usys_log_warn("Emergency PA disable for FEM%d", unit);
    
    if (gpio_set_pa_vds(controller, unit, false) != STATUS_OK) {
        usys_log_error("Failed to disable PA_VDS for FEM%d", unit);
        return STATUS_NOK;
    }
    
    if (gpio_set_28v_vds(controller, unit, false) != STATUS_OK) {
        usys_log_error("Failed to disable 28V_VDS for FEM%d", unit);
        return STATUS_NOK;
    }
    
    usys_log_info("PA disabled successfully for FEM%d", unit);
    return STATUS_OK;
}