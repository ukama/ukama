/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_se98.h"

#include "devhelper.h"
#include "errorcode.h"
#include "irqdb.h"
#include "irqhelper.h"
#include "property.h"
#include "drivers/se98_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static SensorCallbackFxn sensor_cb;

const DrvrOps se98WrapperOps = { .init = se98_wrapper_init,
                                             .configure = se98_wrapper_configure,
                                             .read = se98_wrapper_read,
                                             .write = se98_wrapper_write,
                                             .enable = se98_wrapper_enable,
                                             .disable = se98_wrapper_disable,
                                             .registerCb = se98_wrapper_reg_cb,
                                             .dregisterCb = se98_wrapper_dreg_cb,
                                             .enableIrq = se98_wrapper_enable_irq,
                                             .disableIrq =
                                                 se98_wrapper_disable_irq };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property se98_property[MAXTEMPPROP] = {
    [T1TEMPVALUE] = { .name = "TEMPERATURE",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysFname = "temp1_input",
                      .depProp = NULL },
    [T1MINLIMIT] = { .name = "LOW LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_min",
                     .depProp = NULL },
    [T1MAXLIMIT] = { .name = "HIGH LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_max",
                     .depProp = NULL },
    [T1CRITLIMIT] = { .name = "CRITICAL LIMIT",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysFname = "temp1_crit",
                      .depProp = NULL },
    [T1MINALARM] = { .name = "LOW LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp1_min_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T1MAXALARM] = { .name = "HIGH LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp1_max_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T1CRITALARM] = { .name = "CRITICAL LIMIT ALERT",
                      .dataType = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysFname = "temp1_crit_alarm",
                      .depProp =
                          &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                          .lmt_idx = T1CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T1CRITHYST] = { .name = "CRITICAL HYSTERESIS",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_crit_hyst",
                     .depProp = NULL },
    [T1MAXHYST] = { .name = "MAX HYSTERESIS",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_AVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysFname = "temp1_max_hyst",
                    .depProp = NULL },
    [T1OFFSET] = { .name = "T1 OFFSET",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_NOTAVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "",
                   .depProp = NULL }
};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &se98WrapperOps;
    }
}

int bsp_se98_get_irq_type(int pidx, uint8_t *alertstate) {
    int ret = 0;
    if (pidx == T1MINALARM) {
        *alertstate = ALARM_STATE_LOW_ALARM_ACTIVE;
    } else if (pidx == T1MAXALARM) {
        *alertstate = ALARM_STATE_HIGH_ALARM_ACTIVE;
    } else if (pidx == T1CRITALARM) {
        *alertstate = ALARM_STATE_CRIT_ALARM_ACTIVE;
    } else {
        ret = -1;
    }
    return ret;
}

int bsp_se98_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_se98_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void se98_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
}

int bsp_se98_reg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int bsp_se98_dreg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int bsp_se98_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty,
                                            &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXTEMPPROP;
        gProperty = se98_property;
        log_debug("SE98: Using static property table with %d property.",
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

int bsp_se98_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_enable_irq(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable_irq(drvr, sensor_cb, pDev, gProperty,
                                 *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_se98_disable_irq(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for SE98 device */
int bsp_se98_confirm_irq(Device *pDev, AlertCallBackData **acbdata, char *fpath,
                     int *evt) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_confirm_irq(drvr, pDev, gProperty, acbdata, fpath,
                                  MAXTEMPPROP, evt);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
