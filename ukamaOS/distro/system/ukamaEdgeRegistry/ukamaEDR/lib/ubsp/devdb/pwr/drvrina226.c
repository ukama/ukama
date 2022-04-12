/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
