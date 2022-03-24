/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_ads1015.h"

#include "devhelper.h"
#include "errorcode.h"
#include "property.h"
#include "drivers/ads1015_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

const DrvrOps ads1015WrapperOps = { .init = ads1015_wrapper_init,
                                    .configure = ads1015_wrapper_configure,
                                    .read = ads1015_wrapper_read,
                                    .write = ads1015_wrapper_write,
                                    .enable = ads1015_wrapper_enable,
                                    .disable = ads1015_wrapper_disable,
                                    .registerCb = NULL,
                                    .dregisterCb = NULL,
                                    .enableIrq = NULL,
                                    .disableIrq = NULL };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property ads1015Property[MAXADCPROP] = {
    [VAIN0AIN1] = { .name = "VOLT OVER AIN0 AND AIN1",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysFname = "in0_input",
                    .depProp = NULL },
    [VAIN0AIN3] = { .name = "VOLT OVER AIN0 AND AIN3",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysFname = "in1_input",
                    .depProp = NULL },
    [VAIN1AIN3] = { .name = "VOLT OVER AIN1 AND AIN3",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysFname = "in2_input",
                    .depProp = NULL },
    [VAIN2AIN3] = { .name = "VOLT OVER AIN2 AND AIN3",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_STATUS,
                    .units = "NA",
                    .sysFname = "in3_input",
                    .depProp = NULL },
    [VAIN0GND] = { .name = "VOLT OVER AIN0 AND GND",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysFname = "in4_input",
                   .depProp = NULL },
    [VAIN1GND] = { .name = "VOLT OVER AIN1 AND GND",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysFname = "in5_input",
                   .depProp = NULL },
    [VAIN2GND] = { .name = "VOLT OVER AIN2 AND GND",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysFname = "in6_input",
                   .depProp = NULL },
    [VAIN3GND] = { .name = "VOLT OVER AIN3 AND GND",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_STATUS,
                   .units = "NA",
                   .sysFname = "in7_input",
                   .depProp = NULL }

};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &ads1015WrapperOps;
    }
}

int bsp_ads1015_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_ads1015_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void bsp_drvr_ads1015_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //ads1015_read(void* pDev, void* prop, void* data );
}

int bsp_ads1015_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXADCPROP;
        gProperty = ads1015Property;
        usys_log_debug("ADS1015: Using static property table with %d property.",
                       gPropertyCount);
    }

    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_init_driver(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ads1015_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
