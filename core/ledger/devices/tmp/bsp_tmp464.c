/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_tmp464.h"

#include "devhelper.h"
#include "errorcode.h"
#include "irqdb.h"
#include "irqhelper.h"
#include "property.h"
#include "drivers/tmp464_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static SensorCallbackFxn sensor_cb;

const DrvrOps drvr_tmp464_fxn_table = { .init = tmp464_wrapper_init,
                                        .configure = tmp464_wrapper_configure,
                                        .read = tmp464_wrapper_read,
                                        .write = tmp464_wrapper_write,
                                        .enable = tmp464_wrapper_enable,
                                        .disable = tmp464_wrapper_disable,
                                        .registerCb = tmp464_wrapper_reg_cb,
                                        .dregisterCb = tmp464_wrapper_dreg_cb,
                                        .enableIrq = tmp464_wrapper_enable_irq,
                                        .disableIrq =
                                            tmp464_wrapper_disable_irq };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property tmp464_property[MAXTEMPPROP] = {
    [T1TEMPVALUE] = { .name = "T1 TEMPERATURE",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysFname = "temp1_input",
                      .depProp = NULL },
    [T1MINLIMIT] = { .name = "T1 LOW LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_min",
                     .depProp = NULL },
    [T1MAXLIMIT] = { .name = "T1 HIGH LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_max",
                     .depProp = NULL },
    [T1CRITLIMIT] = { .name = "T1 CRITICAL LIMIT",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysFname = "temp1_crit",
                      .depProp = NULL },
    [T1MINALARM] = { .name = "T1 LOW LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp1_min_alarm",
                     .depProp = &(DepProperty){ .currIdx = T1TEMPVALUE,
                                                .lmtIdx = T1MINLIMIT,
                                                .cond = LESSTHENEQUALTO } },
    [T1MAXALARM] = { .name = "T1 HIGH LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp1_max_alarm",
                     .depProp = &(DepProperty){ .currIdx = T1TEMPVALUE,
                                                .lmtIdx = T1MAXLIMIT,
                                                .cond = GREATERTHENEQUALTO } },
    [T1CRITALARM] = { .name = "T1 CRITICAL LIMIT ALERT",
                      .dataType = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysFname = "temp1_crit_alarm",
                      .depProp = &(DepProperty){ .currIdx = T1TEMPVALUE,
                                                 .lmtIdx = T1CRITLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T1CRITHYST] = { .name = "T1 CRITICAL HYSTERESIS",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp1_crit_hyst",
                     .depProp = NULL },
    [T1MAXHYST] = { .name = "T1 MAX HYSTERESIS",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysFname = "temp1_max_hyst",
                    .depProp = NULL },
    [T1OFFSET] = { .name = "T1 OFFSET",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "temp1_offset",
                   .depProp = NULL },
    [T2TEMPVALUE] = { .name = "T2 TEMPERATURE",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysFname = "temp2_input",
                      .depProp = NULL },
    [T2MINLIMIT] = { .name = "T2 LOW LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp2_min",
                     .depProp = NULL },
    [T2MAXLIMIT] = { .name = "T2 HIGH LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp2_max",
                     .depProp = NULL },
    [T2CRITLIMIT] = { .name = "T2 CRITICAL LIMIT",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysFname = "temp2_crit",
                      .depProp = NULL },
    [T2MINALARM] = { .name = "T2 LOW LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp2_min_alarm",
                     .depProp = &(DepProperty){ .currIdx = T2TEMPVALUE,
                                                .lmtIdx = T2MINLIMIT,
                                                .cond = LESSTHENEQUALTO } },
    [T2MAXALARM] = { .name = "T2 HIGH LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp2_max_alarm",
                     .depProp = &(DepProperty){ .currIdx = T2TEMPVALUE,
                                                .lmtIdx = T2MAXLIMIT,
                                                .cond = GREATERTHENEQUALTO } },
    [T2CRITALARM] = { .name = "T2 CRITICAL LIMIT ALERT",
                      .dataType = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysFname = "temp2_crit_alarm",
                      .depProp = &(DepProperty){ .currIdx = T2TEMPVALUE,
                                                 .lmtIdx = T2CRITLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T2CRITHYST] = { .name = "T2 CRITICAL HYSTERESIS",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp2_crit_hyst",
                     .depProp = NULL },
    [T2MAXHYST] = { .name = "T2 MAX HYSTERESIS",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysFname = "temp2_max_hyst",
                    .depProp = NULL },
    [T2OFFSET] = { .name = "T2 OFFSET",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "",
                   .depProp = NULL },
    [T3TEMPVALUE] = { .name = "T3 TEMPERATURE",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysFname = "temp3_input",
                      .depProp = NULL },
    [T3MINLIMIT] = { .name = "T3 LOW LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp3_min",
                     .depProp = NULL },
    [T3MAXLIMIT] = { .name = "T3 HIGH LIMIT",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp3_max",
                     .depProp = NULL },
    [T3CRITLIMIT] = { .name = "T3 CRITICAL LIMIT",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysFname = "temp3_crit",
                      .depProp = NULL },
    [T3MINALARM] = { .name = "T3 LOW LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp3_min_alarm",
                     .depProp = &(DepProperty){ .currIdx = T3TEMPVALUE,
                                                .lmtIdx = T3MINLIMIT,
                                                .cond = LESSTHENEQUALTO } },
    [T3MAXALARM] = { .name = "T3 HIGH LIMIT ALERT",
                     .dataType = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysFname = "temp3_max_alarm",
                     .depProp = &(DepProperty){ .currIdx = T3TEMPVALUE,
                                                .lmtIdx = T3MAXLIMIT,
                                                .cond = GREATERTHENEQUALTO } },
    [T3CRITALARM] = { .name = "T3 CRITICAL LIMIT ALERT",
                      .dataType = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysFname = "temp3_crit_alarm",
                      .depProp = &(DepProperty){ .currIdx = T3TEMPVALUE,
                                                 .lmtIdx = T3CRITLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T3CRITHYST] = { .name = "T3 CRITICAL HYSTERESIS",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysFname = "temp3_crit_hyst",
                     .depProp = NULL },
    [T3MAXHYST] = { .name = "T3 MAX HYSTERESIS",
                    .dataType = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .propType = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysFname = "temp3_max_hyst",
                    .depProp = NULL },
    [T3OFFSET] = { .name = "T3 OFFSET",
                   .dataType = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .propType = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysFname = "temp3_offset",
                   .depProp = NULL }
};

static const DrvrOps *get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &drvr_tmp464_fxn_table;
    }
}

int bsp_tmp464_get_irq_type(int pidx, uint8_t *alertstate) {
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

int bsp_tmp464_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_tmp464_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void tmp464_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //tmp464_read(void* pDev, void* prop, void* data );
}

int bsp_tmp464_reg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int bsp_tmp464_dreg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int bsp_tmp464_init(Device *pDev) {
    int ret = 0;

    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXTEMPPROP;
        gProperty = tmp464_property;
        log_debug("TMP464: Using static property table with %d property.",
                  gPropertyCount);
    }

    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_init_driver(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }

    /* Register end points to https server */

    return ret;
}

int bsp_tmp464_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_enable_irq(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable_irq(drvr, sensor_cb, pDev, gProperty, *(int *)prop,
                                 data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_tmp464_disable_irq(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for TMP464 device */
int bsp_tmp464_confirm_irq(Device *pDev, AlertCallBackData **acbdata,
                           char *fpath, int *evt) {
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
