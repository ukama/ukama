/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "devdb/pwr/drvrina226.h"

#include "headers/utils/log.h"

int drvr_ina226_init() {
    return 0;
}

int drvr_ina226_registration(Device *p_dev) {
    return 0;
}

int drvr_ina226_read_properties(DevObj *obj, void *prop, uint16_t *count) {
    return 0;
}

int drvr_ina226_configure(void *p_dev, void *prop, void *data) {
    return 0;
}

int drvr_ina226_read(void *p_dev, void *prop, void *data) {
    return 0;
}

int drvr_ina226_write(void *p_dev, void *prop, void *data) {
    return 0;
}
int drvr_ina226_enable(void *p_dev, void *prop, void *data) {
    return 0;
}

int drvr_ina226_disable(void *p_dev, void *prop, void *data) {
    return 0;
}

int drvr_ina226_reg_cb(void *p_dev, SensorCallbackFxn fun) {
    return 0;
}

int drvr_ina226_dreg_cb(void *p_dev, SensorCallbackFxn fun) {
    return 0;
}
int drvr_ina226_enable_irq(void *p_dev, void *prop, void *data) {
    return 0;
}

int drvr_ina226_disable_irq(void *p_dev, void *prop, void *data) {
    return 0;
}
