/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_adt7481.h"

#include "devhelper.h"
#include "errorcode.h"
#include "irqdb.h"
#include "irqhelper.h"
#include "property.h"
#include "drivers/adt7481_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"

static SensorCallbackFxn sensor_cb;

const DrvrOps adt7841WrapperOps = {
    .init = adt7841_wrapper_init,
    .configure = adt7841_wrapper_configure,
    .read = adt7841_wrapper_read,
    .write = adt7841_wrapper_write,
    .enable = adt7841_wrapper_enable,
    .disable = adt7841_wrapper_disable,
    .registerCb = adt7841_wrapper_reg_cb,
    .dregisterCb = adt7841_wrapper_dreg_cb,
    .enableIrq = adt7841_wrapper_enable_irq,
    .disable_Irq = adt7841_wrapper_disable_irq
};

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property adt7481_property[MAXTEMPPROP] = {
    [T1TEMPVALUE] = { .name = "T1 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp1_input",
                      .depProp = NULL },
    [T1MINLIMIT] = { .name = "T1 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_min",
                     .depProp = NULL },
    [T1MAXLIMIT] = { .name = "T1 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_max",
                     .depProp = NULL },
    [T1CRITLIMIT] = { .name = "T1 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp1_crit",
                      .depProp = NULL },
    [T1MINALARM] = { .name = "T1 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_min_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T1MAXALARM] = { .name = "T1 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_max_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T1CRITALARM] = { .name = "T1 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp1_crit_alarm",
                      .depProp =
                          &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                          .lmt_idx = T1CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T1CRITHYST] = { .name = "T1 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_crit_hyst",
                     .depProp = NULL },
    [T1MAXHYST] = { .name = "T1 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp1_max_hyst",
                    .depProp = NULL },
    [T1OFFSET] = { .name = "T1 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_NOTAVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "temp1_offset",
                   .depProp = NULL },
    [T2TEMPVALUE] = { .name = "T2 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp2_input",
                      .depProp = NULL },
    [T2MINLIMIT] = { .name = "T2 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_min",
                     .depProp = NULL },
    [T2MAXLIMIT] = { .name = "T2 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_max",
                     .depProp = NULL },
    [T2CRITLIMIT] = { .name = "T2 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp2_crit",
                      .depProp = NULL },
    [T2MINALARM] = { .name = "T2 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp2_min_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                                 .lmt_idx = T2MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T2MAXALARM] = { .name = "T2 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp2_max_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                                 .lmt_idx = T2MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T2CRITALARM] = { .name = "T2 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp2_crit_alarm",
                      .depProp =
                          &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                          .lmt_idx = T2CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T2CRITHYST] = { .name = "T2 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_crit_hyst",
                     .depProp = NULL },
    [T2MAXHYST] = { .name = "T2 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp2_max_hyst",
                    .depProp = NULL },
    [T2OFFSET] = { .name = "T2 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_NOTAVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "",
                   .depProp = NULL },
    [T3TEMPVALUE] = { .name = "T3 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp3_input",
                      .depProp = NULL },
    [T3MINLIMIT] = { .name = "T3 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_min",
                     .depProp = NULL },
    [T3MAXLIMIT] = { .name = "T3 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_max",
                     .depProp = NULL },
    [T3CRITLIMIT] = { .name = "T3 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp3_crit",
                      .depProp = NULL },
    [T3MINALARM] = { .name = "T3 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp3_min_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                                 .lmt_idx = T3MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T3MAXALARM] = { .name = "T3 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp3_max_alarm",
                     .depProp = &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                                 .lmt_idx = T3MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T3CRITALARM] = { .name = "T3 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp3_crit_alarm",
                      .depProp =
                          &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                          .lmt_idx = T3CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T3CRITHYST] = { .name = "T3 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_crit_hyst",
                     .depProp = NULL },
    [T3MAXHYST] = { .name = "T3 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp3_max_hyst",
                    .depProp = NULL },
    [T3OFFSET] = { .name = "T3 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_NOTAVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "temp3_offset",
                   .depProp = NULL }
};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &adt7841WrapperOps;
    }
}

int bsp_adt7481_get_irq_type(int pidx, uint8_t *alertstate) {
    int ret = 0;
    if ((pidx == T1MINALARM) || (pidx == T2MINALARM) || (pidx == T3MINALARM)) {
        *alertstate = ALARM_STATE_LOW_ALARM_ACTIVE;
    } else if ((pidx == T1MAXALARM) || (pidx == T2MAXALARM) ||
               (pidx == T3MAXALARM)) {
        *alertstate = ALARM_STATE_HIGH_ALARM_ACTIVE;
    } else if ((pidx == T1CRITALARM) || (pidx == T2CRITALARM) ||
               (pidx == T3CRITALARM)) {
        *alertstate = ALARM_STATE_CRIT_ALARM_ACTIVE;
    } else {
        ret = -1;
    }
    return ret;
}

int bsp_adt7481_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_adt7481_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void adt7481_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    return;
}

int bsp_adt7481_reg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int bsp_adt7481_dreg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int bsp_adt7481_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty,
                                            &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXTEMPPROP;
        gProperty = adt7481_property;
        log_debug("ADT7481: Using static property table with %d property.",
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

int bsp_adt7481_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_adt7481_enable_irq(void *pDev, void *prop, void *data) {
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

int bsp_adt7481_disable_irq(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for ADT7481 device */
int bsp_adt7481_confirm_irq(Device *pDev, AlertCallBackData **acbdata, char *fpath,
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
