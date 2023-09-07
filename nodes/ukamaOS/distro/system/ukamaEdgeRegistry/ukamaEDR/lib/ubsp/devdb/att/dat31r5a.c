/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/att/dat31r5a.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "headers/utils/log.h"
#include "devdb/att/drvrdat31r5a.h"
#include "devdb/sysfs/drvrsysfs.h"

const DrvDBFxnTable drvr_dat31r5a_fxn_table = { .init = drvr_dat31r5a_init,
                                                .configure =
                                                    drvr_dat31r5a_configure,
                                                .read = drvr_dat31r5a_read,
                                                .write = drvr_dat31r5a_write,
                                                .enable = drvr_dat31r5a_enable,
                                                .disable =
                                                    drvr_dat31r5a_disable,
                                                .register_cb = NULL,
                                                .dregister_cb = NULL,
                                                .enable_irq = NULL,
                                                .disable_irq = NULL };

static Property *g_property = NULL;
static int g_property_count = 0;

static Property dat31r5a_property[MAXATTPROP] = {
    [ATTVALUE] = { .name = "ATTENUATION VALUE",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "dB",
                   .sysfname = "in0_attvalue",
                   .dep_prop = NULL },
    [LATCHENABLE] = { .name = "LATCH FOR ATTENUATION",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysfname = "in0_latch",
                      .dep_prop = NULL },
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_dat31r5a_fxn_table;
    }
}

int dat31r5a_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int dat31r5a_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void dat31r5a_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //dat31r5a_read(void* p_dev, void* prop, void* data );
}

int dat31r5a_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXATTPROP;
        g_property = dat31r5a_property;
        log_debug("DAT31R5A: Using static property table with %d property.",
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

int dat31r5a_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int dat31r5a_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int dat31r5a_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int dat31r5a_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int dat31r5a_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int dat31r5a_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
