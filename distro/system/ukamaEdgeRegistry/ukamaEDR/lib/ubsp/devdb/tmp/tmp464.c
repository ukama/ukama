/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/tmp/tmp464.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "utils/irqdb.h"
#include "utils/irqhelper.h"
#include "headers/utils/log.h"
#include "devdb/tmp/drvrtmp464.h"
#include "devdb/sysfs/drvrsysfs.h"

static SensorCallbackFxn sensor_cb;

const DrvDBFxnTable drvr_tmp464_fxn_table = {
    .init = drvr_tmp464_init,
    .configure = drvr_tmp464_configure,
    .read = drvr_tmp464_read,
    .write = drvr_tmp464_write,
    .enable = drvr_tmp464_enable,
    .disable = drvr_tmp464_disable,
    .register_cb = drvr_tmp464_reg_cb,
    .dregister_cb = drvr_tmp464_dreg_cb,
    .enable_irq = drvr_tmp464_enable_irq,
    .disable_irq = drvr_tmp464_disable_irq
};

static Property *g_property = NULL;
static int g_property_count = 0;

static Property tmp464_property[MAXTEMPPROP] = {
    [T1TEMPVALUE] = { .name = "T1 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp1_input",
                      .dep_prop = NULL },
    [T1MINLIMIT] = { .name = "T1 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_min",
                     .dep_prop = NULL },
    [T1MAXLIMIT] = { .name = "T1 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_max",
                     .dep_prop = NULL },
    [T1CRITLIMIT] = { .name = "T1 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp1_crit",
                      .dep_prop = NULL },
    [T1MINALARM] = { .name = "T1 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_min_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T1MAXALARM] = { .name = "T1 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_max_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T1CRITALARM] = { .name = "T1 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp1_crit_alarm",
                      .dep_prop =
                          &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                          .lmt_idx = T1CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T1CRITHYST] = { .name = "T1 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_crit_hyst",
                     .dep_prop = NULL },
    [T1MAXHYST] = { .name = "T1 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .prop_type = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp1_max_hyst",
                    .dep_prop = NULL },
    [T1OFFSET] = { .name = "T1 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "temp1_offset",
                   .dep_prop = NULL },
    [T2TEMPVALUE] = { .name = "T2 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp2_input",
                      .dep_prop = NULL },
    [T2MINLIMIT] = { .name = "T2 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_min",
                     .dep_prop = NULL },
    [T2MAXLIMIT] = { .name = "T2 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_max",
                     .dep_prop = NULL },
    [T2CRITLIMIT] = { .name = "T2 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp2_crit",
                      .dep_prop = NULL },
    [T2MINALARM] = { .name = "T2 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp2_min_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                                 .lmt_idx = T2MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T2MAXALARM] = { .name = "T2 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp2_max_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                                 .lmt_idx = T2MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T2CRITALARM] = { .name = "T2 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp2_crit_alarm",
                      .dep_prop =
                          &(DepProperty){ .curr_idx = T2TEMPVALUE,
                                          .lmt_idx = T2CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T2CRITHYST] = { .name = "T2 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp2_crit_hyst",
                     .dep_prop = NULL },
    [T2MAXHYST] = { .name = "T2 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .prop_type = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp2_max_hyst",
                    .dep_prop = NULL },
    [T2OFFSET] = { .name = "T2 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "",
                   .dep_prop = NULL },
    [T3TEMPVALUE] = { .name = "T3 TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp3_input",
                      .dep_prop = NULL },
    [T3MINLIMIT] = { .name = "T3 LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_min",
                     .dep_prop = NULL },
    [T3MAXLIMIT] = { .name = "T3 HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_max",
                     .dep_prop = NULL },
    [T3CRITLIMIT] = { .name = "T3 CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp3_crit",
                      .dep_prop = NULL },
    [T3MINALARM] = { .name = "T3 LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp3_min_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                                 .lmt_idx = T3MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T3MAXALARM] = { .name = "T3 HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp3_max_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                                 .lmt_idx = T3MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T3CRITALARM] = { .name = "T3 CRITICAL LIMIT ALERT",
                      .data_type = TYPE_BOOL,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_ALERT,
                      .units = "NA",
                      .sysfname = "temp3_crit_alarm",
                      .dep_prop =
                          &(DepProperty){ .curr_idx = T3TEMPVALUE,
                                          .lmt_idx = T3CRITLIMIT,
                                          .cond = GREATERTHENEQUALTO } },
    [T3CRITHYST] = { .name = "T3 CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp3_crit_hyst",
                     .dep_prop = NULL },
    [T3MAXHYST] = { .name = "T3 MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_NOTAVAIL,
                    .prop_type = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp3_max_hyst",
                    .dep_prop = NULL },
    [T3OFFSET] = { .name = "T3 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_AVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "temp3_offset",
                   .dep_prop = NULL }
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_tmp464_fxn_table;
    }
}

int tmp464_get_irq_type(int pidx, uint8_t *alertstate) {
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

int tmp464_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int tmp464_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void tmp464_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //tmp464_read(void* p_dev, void* prop, void* data );
}

int tmp464_reg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int tmp464_dreg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int tmp464_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXTEMPPROP;
        g_property = tmp464_property;
        log_debug("TMP464: Using static property table with %d property.",
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

int tmp464_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_enable_irq(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable_irq(drvr, sensor_cb, p_dev, g_property,
                                 *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int tmp464_disable_irq(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for TMP464 device */
int tmp464_confirm_irq(Device *p_dev, AlertCallBackData **acbdata, char *fpath,
                       int *evt) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_confirm_irq(drvr, p_dev, g_property, acbdata, fpath,
                                  MAXTEMPPROP, evt);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
