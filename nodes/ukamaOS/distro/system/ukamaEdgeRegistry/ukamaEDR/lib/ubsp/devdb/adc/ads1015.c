/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/adc/ads1015.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "headers/utils/log.h"
#include "devdb/adc/drvrads1015.h"
#include "devdb/sysfs/drvrsysfs.h"

const DrvDBFxnTable drvr_ads1015_fxn_table = { .init = drvr_ads1015_init,
                                               .configure =
                                                   drvr_ads1015_configure,
                                               .read = drvr_ads1015_read,
                                               .write = drvr_ads1015_write,
                                               .enable = drvr_ads1015_enable,
                                               .disable = drvr_ads1015_disable,
                                               .register_cb = NULL,
                                               .dregister_cb = NULL,
                                               .enable_irq = NULL,
                                               .disable_irq = NULL };

static Property *g_property = NULL;
static int g_property_count = 0;

static Property ads1015_property[MAXADCPROP] = {
    [VAIN0AIN1] = { .name = "VOLT OVER AIN0 AND AIN1",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysfname = "in0_input",
                    .dep_prop = NULL },
    [VAIN0AIN3] = { .name = "VOLT OVER AIN0 AND AIN3",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysfname = "in1_input",
                    .dep_prop = NULL },
    [VAIN1AIN3] = { .name = "VOLT OVER AIN1 AND AIN3",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysfname = "in2_input",
                    .dep_prop = NULL },
    [VAIN2AIN3] = { .name = "VOLT OVER AIN2 AND AIN3",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysfname = "in3_input",
                    .dep_prop = NULL },
    [VAIN0GND] = { .name = "VOLT OVER AIN0 AND GND",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysfname = "in4_input",
                   .dep_prop = NULL },
    [VAIN1GND] = { .name = "VOLT OVER AIN1 AND GND",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysfname = "in5_input",
                   .dep_prop = NULL },
    [VAIN2GND] = { .name = "VOLT OVER AIN2 AND GND",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysfname = "in6_input",
                   .dep_prop = NULL },
    [VAIN3GND] = { .name = "VOLT OVER AIN3 AND GND",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysfname = "in7_input",
                   .dep_prop = NULL }

};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_ads1015_fxn_table;
    }
}

int ads1015_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int ads1015_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void ads1015_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //ads1015_read(void* p_dev, void* prop, void* data );
}

int ads1015_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXADCPROP;
        g_property = ads1015_property;
        log_debug("ADS1015: Using static property table with %d property.",
                  g_property_count);
    }

    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_init_driver(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ads1015_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
