/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_dat31r5a.h"

#include "devhelper.h"
#include "errorcode.h"
#include "property.h"
#include "drivers/dat31r5a_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

const DrvrOps dat31r5aWrapperOps = { .init = dat31r5a_wrapper_init,
                                     .configure = dat31r5a_wrapper_configure,
                                     .read = dat31r5a_wrapper_read,
                                     .write = dat31r5a_wrapper_write,
                                     .enable = dat31r5a_wrapper_enable,
                                     .disable = dat31r5a_wrapper_disable,
                                     .registerCb = NULL,
                                     .dregisterCb = NULL,
                                     .enableIrq = NULL,
                                     .disableIrq = NULL };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property dat31r5aProperty[MAXATTPROP] = {
    [ATTVALUE] = { .name = "ATTENUATION VALUE",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "dB",
                   .sysFname = "in0_attvalue",
                   .depProp = NULL },
    [LATCHENABLE] = { .name = "LATCH FOR ATTENUATION",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysFname = "in0_latch",
                      .depProp = NULL },
};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &dat31r5aWrapperOps;
    }
}

int bsp_dat31r5a_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_dat31r5a_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void dat31r5a_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //dat31r5a_read(void* pDev, void* prop, void* data );
}

int bsp_dat31r5a_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXATTPROP;
        gProperty = dat31r5aProperty;
        log_debug("DAT31R5A: Using static property table with %d property.",
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

int bsp_dat31r5a_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_dat31r5a_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_dat31r5a_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_dat31r5a_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_dat31r5a_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_dat31r5a_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
