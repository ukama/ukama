/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/tmp/se98.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "utils/irqdb.h"
#include "utils/irqhelper.h"
#include "headers/utils/log.h"
#include "devdb/sysfs/drvrsysfs.h"
#include "devdb/tmp/drvrse98.h"

static SensorCallbackFxn sensor_cb;

const DrvDBFxnTable drvr_s_e98_fxn_table = { .init = drvr_se98_init,
                                             .configure = drvr_se98_configure,
                                             .read = drvr_se98_read,
                                             .write = drvr_se98_write,
                                             .enable = drvr_se98_enable,
                                             .disable = drvr_se98_disable,
                                             .register_cb = drvr_se98_reg_cb,
                                             .dregister_cb = drvr_se98_dreg_cb,
                                             .enable_irq = drvr_se98_enable_irq,
                                             .disable_irq =
                                                 drvr_se98_disable_irq };

static Property *g_property = NULL;
static int g_property_count = 0;

static Property se98_property[MAXTEMPPROP] = {
    [T1TEMPVALUE] = { .name = "TEMPERATURE",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_STATUS,
                      .units = "milliCelsius",
                      .sysfname = "temp1_input",
                      .dep_prop = NULL },
    [T1MINLIMIT] = { .name = "LOW LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_min",
                     .dep_prop = NULL },
    [T1MAXLIMIT] = { .name = "HIGH LIMIT",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_max",
                     .dep_prop = NULL },
    [T1CRITLIMIT] = { .name = "CRITICAL LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "milliCelsius",
                      .sysfname = "temp1_crit",
                      .dep_prop = NULL },
    [T1MINALARM] = { .name = "LOW LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_min_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MINLIMIT,
                                                 .cond = LESSTHENEQUALTO } },
    [T1MAXALARM] = { .name = "HIGH LIMIT ALERT",
                     .data_type = TYPE_BOOL,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_ALERT,
                     .units = "NA",
                     .sysfname = "temp1_max_alarm",
                     .dep_prop = &(DepProperty){ .curr_idx = T1TEMPVALUE,
                                                 .lmt_idx = T1MAXLIMIT,
                                                 .cond = GREATERTHENEQUALTO } },
    [T1CRITALARM] = { .name = "CRITICAL LIMIT ALERT",
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
    [T1CRITHYST] = { .name = "CRITICAL HYSTERESIS",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD | PERM_WR,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_CONFIG,
                     .units = "milliCelsius",
                     .sysfname = "temp1_crit_hyst",
                     .dep_prop = NULL },
    [T1MAXHYST] = { .name = "MAX HYSTERESIS",
                    .data_type = TYPE_INT32,
                    .perm = PERM_RD | PERM_WR,
                    .available = PROP_AVAIL,
                    .prop_type = PROP_TYPE_CONFIG,
                    .units = "milliCelsius",
                    .sysfname = "temp1_max_hyst",
                    .dep_prop = NULL },
    [T1OFFSET] = { .name = "T1 OFFSET",
                   .data_type = TYPE_INT32,
                   .perm = PERM_RD | PERM_WR,
                   .available = PROP_NOTAVAIL,
                   .prop_type = PROP_TYPE_CONFIG,
                   .units = "NA",
                   .sysfname = "",
                   .dep_prop = NULL }
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_s_e98_fxn_table;
    }
}

int se98_get_irq_type(int pidx, uint8_t *alertstate) {
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

int se98_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int se98_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void se98_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
}

int se98_reg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int se98_dreg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int se98_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXTEMPPROP;
        g_property = se98_property;
        log_debug("SE98: Using static property table with %d property.",
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

int se98_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int se98_enable_irq(void *p_dev, void *prop, void *data) {
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

int se98_disable_irq(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for SE98 device */
int se98_confirm_irq(Device *p_dev, AlertCallBackData **acbdata, char *fpath,
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
