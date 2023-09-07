/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/led/ledtricol.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "headers/utils/log.h"
#include "devdb/led/drvrledtricol.h"
#include "devdb/sysfs/drvrsysfs.h"

const DrvDBFxnTable drvr_led_tricol_fxn_table = {
    .init = drvr_led_tricol_init,
    .configure = drvr_led_tricol_configure,
    .read = drvr_led_tricol_read,
    .write = drvr_led_tricol_write,
    .enable = drvr_led_tricol_enable,
    .disable = drvr_led_tricol_disable,
    .register_cb = NULL,
    .dregister_cb = NULL,
    .enable_irq = NULL,
    .disable_irq = NULL
};

static Property *g_property = NULL;
static int g_property_count = 0;

static Property led_tricol_property[MAXLEDTRICOLPROP] = {
    [RBRIGHTNESS] = { .name = "RED LED BRIGHTNESS",
                      .data_type = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysfname = "red/brightness",
                      .dep_prop = NULL },
    [RMAX_BRIGHTNESS] = { .name = "RED LED MAX BRIGHTNESS",
                          .data_type = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .prop_type = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysfname = "red/max_brightness",
                          .dep_prop = NULL },
    [RTRIGGER] = { .name = "RED LED TRIGGER",
                   .data_type = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "red/trigger",
                   .dep_prop = NULL },
    [GBRIGHTNESS] = { .name = "GREEN LED BRIGHTNESS",
                      .data_type = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysfname = "green/brightness",
                      .dep_prop = NULL },
    [GMAX_BRIGHTNESS] = { .name = "GREEN LED MAX BRIGHTNESS",
                          .data_type = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .prop_type = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysfname = "green/max_brightness",
                          .dep_prop = NULL },
    [GTRIGGER] = { .name = "GREEN LED TRIGGER",
                   .data_type = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "green/trigger",
                   .dep_prop = NULL },
    [BBRIGHTNESS] = { .name = "BLUE LED BRIGHTNESS",
                      .data_type = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysfname = "blue/brightness",
                      .dep_prop = NULL },
    [BMAX_BRIGHTNESS] = { .name = "BLUE LED MAX BRIGHTNESS",
                          .data_type = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .prop_type = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysfname = "blue/max_brightness",
                          .dep_prop = NULL },
    [BTRIGGER] = { .name = "BLUE LED TRIGGER",
                   .data_type = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "blue/trigger",
                   .dep_prop = NULL }
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_led_tricol_fxn_table;
    }
}

int led_tricol_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int led_tricol_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void led_tricol_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //led_tricol_read(void* p_dev, void* prop, void* data );
}
int led_tricol_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXLEDTRICOLPROP;
        g_property = led_tricol_property;
        log_debug("LEDTRICOL: Using static property table with %d property.",
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

int led_tricol_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int led_tricol_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int led_tricol_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int led_tricol_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int led_tricol_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int led_tricol_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
