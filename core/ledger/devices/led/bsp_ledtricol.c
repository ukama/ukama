/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_ledtricol.h"

#include "devhelper.h"
#include "errorcode.h"
#include "property.h"
#include "drivers/ledtricol_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

const DrvrOps ledTricolWrapperOps = { .init = led_tricol_wrapper_init,
                                      .configure = led_tricol_wrapper_configure,
                                      .read = led_tricol_wrapper_read,
                                      .write = led_tricol_wrapper_write,
                                      .enable = led_tricol_wrapper_enable,
                                      .disable = led_tricol_wrapper_disable,
                                      .registerCb = NULL,
                                      .dregisterCb = NULL,
                                      .enableIrq = NULL,
                                      .disableIrq = NULL };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property led_tricol_property[MAXLEDTRICOLPROP] = {
    [RBRIGHTNESS] = { .name = "RED LED BRIGHTNESS",
                      .dataType = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysFname = "red/brightness",
                      .depProp = NULL },
    [RMAX_BRIGHTNESS] = { .name = "RED LED MAX BRIGHTNESS",
                          .dataType = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .propType = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysFname = "red/max_brightness",
                          .depProp = NULL },
    [RTRIGGER] = { .name = "RED LED TRIGGER",
                   .dataType = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "red/trigger",
                   .depProp = NULL },
    [GBRIGHTNESS] = { .name = "GREEN LED BRIGHTNESS",
                      .dataType = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysFname = "green/brightness",
                      .depProp = NULL },
    [GMAX_BRIGHTNESS] = { .name = "GREEN LED MAX BRIGHTNESS",
                          .dataType = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .propType = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysFname = "green/max_brightness",
                          .depProp = NULL },
    [GTRIGGER] = { .name = "GREEN LED TRIGGER",
                   .dataType = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "green/trigger",
                   .depProp = NULL },
    [BBRIGHTNESS] = { .name = "BLUE LED BRIGHTNESS",
                      .dataType = TYPE_UINT8,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "NA",
                      .sysFname = "blue/brightness",
                      .depProp = NULL },
    [BMAX_BRIGHTNESS] = { .name = "BLUE LED MAX BRIGHTNESS",
                          .dataType = TYPE_UINT8,
                          .perm = PERM_RD | PERM_WR,
                          .available = PROP_AVAIL,
                          .propType = PROP_TYPE_CONFIG,
                          .units = "NA",
                          .sysFname = "blue/max_brightness",
                          .depProp = NULL },
    [BTRIGGER] = { .name = "BLUE LED TRIGGER",
                   .dataType = TYPE_UINT8,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "blue/trigger",
                   .depProp = NULL }
};

static const DrvrOps *get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &ledTricolWrapperOps;
    }
}

int bsp_led_tricol_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_led_tricol_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void led_tricol_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //led_tricol_read(void* pDev, void* prop, void* data );
}
int bsp_led_tricol_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXLEDTRICOLPROP;
        gProperty = led_tricol_property;
        log_debug("LEDTRICOL: Using static property table with %d property.",
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

int bsp_led_tricol_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_led_tricol_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_led_tricol_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_led_tricol_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_led_tricol_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_led_tricol_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
