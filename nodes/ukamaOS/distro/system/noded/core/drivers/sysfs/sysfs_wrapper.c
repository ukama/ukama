/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "drivers/sysfs_wrapper.h"

#include "errorcode.h"
#include "property.h"
#include "driver_ops.h"
#include "drivers/sysfs.h"

SensorCallbackFxn sensorCb;

const DrvrOps drvrSysFsOps = { .init = sysfs_wrapper_init,
                               .configure = sysfs_wrapper_configure,
                               .read = sysfs_wrapper_read,
                               .write = sysfs_wrapper_write,
                               .enable = sysfs_wrapper_enable,
                               .disable = sysfs_wrapper_disable,
                               .registerCb = sysfs_wrapper_reg_cb,
                               .dregisterCb = sysfs_wrapper_dreg_cb,
                               .enableIrq = sysfs_wrapper_enable_irq,
                               .disableIrq = sysfs_wrapper_disable_irq };

const DrvrOps *sysfs_wrapper_get_ops() {
    return &drvrSysFsOps;
}

int sysfs_wrapper_init() {
    return 0;
}

int sysfs_wrapper_registration(Device *dev) {
    return 0;
}

int sysfs_wrapper_configure(void *hwAttr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwAttr)) {
        ret = sysfs_write(hwAttr, data, property->dataType);
        if (ret <= 0) {
            ret = ERR_NODED_SYSFS_WRITE_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_NODED_SYSFS_FILE_MISSING;
    }
    return ret;
}

int sysfs_wrapper_read(void *hwAttr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwAttr)) {
        ret = sysfs_read(hwAttr, data, property->dataType);
        if (ret) {
            ret = ERR_NODED_SYSFS_READ_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_NODED_SYSFS_FILE_MISSING;
    }
    return ret;
}

int sysfs_wrapper_write(void *hwAttr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwAttr)) {
        ret = sysfs_write(hwAttr, data, property->dataType);
        if (ret) {
            ret = ERR_NODED_SYSFS_WRITE_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_NODED_SYSFS_FILE_MISSING;
    }
    return ret;
}

int sysfs_wrapper_enable(void *hwAttr, void *prop, void *data) {
    int ret = 0;

    return 0;
}

int sysfs_wrapper_disable(void *hwAttr, void *prop, void *data) {
    return 0;
}

int sysfs_wrapper_reg_cb(void *hwAttr, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensorCb = fun;
    }
    return ret;
}

int sysfs_wrapper_dreg_cb(void *hwAttr, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensorCb = NULL;
    }
    return ret;
}

int sysfs_wrapper_enable_irq(void *hwAttr, void *prop, void *data) {
    return 0;
}

int sysfs_wrapper_disable_irq(void *hwAttr, void *prop, void *data) {
    return 0;
}
