/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include "devdb/sysfs/drvrsysfs.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/driverfxn.h"
#include "devdb/sysfs/sysfs.h"

SensorCallbackFxn sensor_cb;

const DrvDBFxnTable drvr_sysfs_fxn_table = { .init = drvr_sysfs_init,
                                             .configure = drvr_sysfs_configure,
                                             .read = drvr_sysfs_read,
                                             .write = drvr_sysfs_write,
                                             .enable = drvr_sysfs_enable,
                                             .disable = drvr_sysfs_disable,
                                             .register_cb = drvr_sysfs_reg_cb,
                                             .dregister_cb = drvr_sysfs_dreg_cb,
                                             .enable_irq =
                                                 drvr_sysfs_enable_irq,
                                             .disable_irq =
                                                 drvr_sysfs_disable_irq };

const DrvDBFxnTable *drvr_sysfs_get_fxn_tbl() {
    return &drvr_sysfs_fxn_table;
}

int drvr_sysfs_init() {
    return 0;
}

int drvr_sysfs_registration(Device *dev) {
    return 0;
}

int drvr_sysfs_configure(void *hwattr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwattr)) {
        ret = sysfs_write(hwattr, data, property->data_type);
        if (ret <= 0) {
            ret = ERR_UBSP_SYSFS_WRITE_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_UBSP_SYSFS_FILE_MISSING;
    }
    return ret;
}

int drvr_sysfs_read(void *hwattr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwattr)) {
        ret = sysfs_read(hwattr, data, property->data_type);
        if (ret) {
            ret = ERR_UBSP_SYSFS_READ_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_UBSP_SYSFS_FILE_MISSING;
    }
    return ret;
}

int drvr_sysfs_write(void *hwattr, void *prop, void *data) {
    int ret = 0;
    Property *property = (Property *)prop;
    if (sysfs_exist(hwattr)) {
        ret = sysfs_write(hwattr, data, property->data_type);
        if (ret) {
            ret = ERR_UBSP_SYSFS_WRITE_FAILED;
        } else {
            ret = 0;
        }
    } else {
        ret = ERR_UBSP_SYSFS_FILE_MISSING;
    }
    return ret;
}

int drvr_sysfs_enable(void *hwattr, void *prop, void *data) {
    int ret = 0;

    return 0;
}

int drvr_sysfs_disable(void *hwattr, void *prop, void *data) {
    return 0;
}

int drvr_sysfs_reg_cb(void *hwattr, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int drvr_sysfs_dreg_cb(void *hwattr, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int drvr_sysfs_enable_irq(void *hwattr, void *prop, void *data) {
    return 0;
}

int drvr_sysfs_disable_irq(void *hwattr, void *prop, void *data) {
    return 0;
}
