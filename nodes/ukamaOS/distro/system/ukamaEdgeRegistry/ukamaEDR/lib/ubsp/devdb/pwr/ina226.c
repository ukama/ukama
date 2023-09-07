/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devdb/pwr/ina226.h"

#include "headers/errorcode.h"
#include "headers/ubsp/property.h"
#include "inc/devhelper.h"
#include "utils/irqdb.h"
#include "utils/irqhelper.h"
#include "headers/utils/log.h"
#include "devdb/pwr/drvrina226.h"
#include "devdb/sysfs/drvrsysfs.h"

static SensorCallbackFxn sensor_cb;

const DrvDBFxnTable drvr_ina226_fxn_table = {
    .init = drvr_ina226_init,
    .configure = drvr_ina226_configure,
    .read = drvr_ina226_read,
    .write = drvr_ina226_write,
    .enable = drvr_ina226_enable,
    .disable = drvr_ina226_disable,
    .register_cb = drvr_ina226_reg_cb,
    .dregister_cb = drvr_ina226_dreg_cb,
    .enable_irq = drvr_ina226_enable_irq,
    .disable_irq = drvr_ina226_disable_irq
};

static Property *g_property = NULL;
static int g_property_count = 0;

static Property ina226_property[MAXINAPROP] = {
    [SHUNTVOLTAGE] = { .name = "SHUNT VOLTAGE",
                       .data_type = TYPE_INT32,
                       .perm = PERM_RD,
                       .available = PROP_AVAIL,
                       .prop_type = PROP_TYPE_STATUS,
                       .units = "milliVolts",
                       .sysfname = "in0_input",
                       .dep_prop = NULL },
    [BUSVOLTAGE] = { .name = "BUS VOLTAGE",
                     .data_type = TYPE_INT32,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .prop_type = PROP_TYPE_STATUS,
                     .units = "milliVolts",
                     .sysfname = "in1_input",
                     .dep_prop = NULL },
    [CURRENT] = { .name = "CURRENT",
                  .data_type = TYPE_INT32,
                  .perm = PERM_RD,
                  .available = PROP_AVAIL,
                  .prop_type = PROP_TYPE_STATUS,
                  .units = "milliAmp",
                  .sysfname = "curr1_input",
                  .dep_prop = NULL },
    [POWER] = { .name = "POWER",
                .data_type = TYPE_INT32,
                .perm = PERM_RD,
                .available = PROP_AVAIL,
                .prop_type = PROP_TYPE_STATUS,
                .units = "microWatt",
                .sysfname = "power1_input",
                .dep_prop = NULL },
    [SHUNTRESISTOR] = { .name = "SHUNT RESISTANCE",
                        .data_type = TYPE_INT32,
                        .perm = PERM_RD | PERM_WR,
                        .available = PROP_AVAIL,
                        .prop_type = PROP_TYPE_CONFIG,
                        .units = "microOhm",
                        .sysfname = "shunt_resistor",
                        .dep_prop = NULL },
    [CALIBRATION] = { .name = "CALIBRATION",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "",
                      .sysfname = "calibration",
                      .dep_prop = NULL },
    [CRITLOWSHUNTVOLTAGE] = { .name = "CRIT LOW SHUNT VOLTAGE",
                              .data_type = TYPE_INT32,
                              .perm = PERM_RD | PERM_WR,
                              .available = PROP_AVAIL,
                              .prop_type = PROP_TYPE_CONFIG,
                              .units = "milliVolts",
                              .sysfname = "in0_lcrit",
                              .dep_prop = NULL },
    [CRITHIGHSHUNTVOLTAGE] = { .name = "CRIT HIGH SHUNT VOLTAGE",
                               .data_type = TYPE_INT32,
                               .perm = PERM_RD | PERM_WR,
                               .available = PROP_AVAIL,
                               .prop_type = PROP_TYPE_CONFIG,
                               .units = "milliVolts",
                               .sysfname = "in0_crit",
                               .dep_prop = NULL },
    [SHUNTVOLTAGECRITLOWALARM] = { .name = "SHUNT VOLTAGE CRIT LOW ALARM",
                                   .data_type = TYPE_BOOL,
                                   .perm = PERM_RD,
                                   .available = PROP_AVAIL,
                                   .prop_type = PROP_TYPE_ALERT,
                                   .units = "NA",
                                   .sysfname = "in0_lcrit_alarm",
                                   .dep_prop =
                                       &(DepProperty){
                                           .curr_idx = SHUNTVOLTAGE,
                                           .lmt_idx = CRITLOWSHUNTVOLTAGE,
                                           .cond = LESSTHENEQUALTO } },
    [SHUNTVOLTAGECRITHIGHALARM] = { .name = "SHUNT VOLTAGE CRIT HIGH ALARM",
                                    .data_type = TYPE_BOOL,
                                    .perm = PERM_RD,
                                    .available = PROP_AVAIL,
                                    .prop_type = PROP_TYPE_ALERT,
                                    .units = "NA",
                                    .sysfname = "in0_crit_alarm",
                                    .dep_prop =
                                        &(DepProperty){
                                            .curr_idx = SHUNTVOLTAGE,
                                            .lmt_idx = CRITHIGHSHUNTVOLTAGE,
                                            .cond = GREATERTHENEQUALTO } },
    [CRITLOWBUSVOLTAGE] = { .name = "LOW VOLTAGE LIMIT",
                            .data_type = TYPE_INT32,
                            .perm = PERM_RD | PERM_WR,
                            .available = PROP_AVAIL,
                            .prop_type = PROP_TYPE_CONFIG,
                            .units = "milliVolts",
                            .sysfname = "in1_lcrit",
                            .dep_prop = NULL },
    [CRITHIGHBUSVOLTAGE] = { .name = "HIGH VOLTAGE LIMIT",
                             .data_type = TYPE_INT32,
                             .perm = PERM_RD,
                             .available = PROP_AVAIL,
                             .prop_type = PROP_TYPE_CONFIG,
                             .units = "milliVolts",
                             .sysfname = "in1_crit",
                             .dep_prop = NULL },
    [BUSVOLTAGECRITLOWALARM] = { .name = "BUS VOLTAGE CRIT LOW ALARM",
                                 .data_type = TYPE_BOOL,
                                 .perm = PERM_RD,
                                 .available = PROP_AVAIL,
                                 .prop_type = PROP_TYPE_ALERT,
                                 .units = "NA",
                                 .sysfname = "in1_lcrit_alarm",
                                 .dep_prop = &(
                                     DepProperty){ .curr_idx = BUSVOLTAGE,
                                                   .lmt_idx = CRITLOWBUSVOLTAGE,
                                                   .cond = LESSTHENEQUALTO } },
    [BUSVOLTAGECRITHIGHALARM] = { .name = "BUS VOLTAGE CRIT HIGH ALARM",
                                  .data_type = TYPE_BOOL,
                                  .perm = PERM_RD,
                                  .available = PROP_AVAIL,
                                  .prop_type = PROP_TYPE_ALERT,
                                  .units = "NA",
                                  .sysfname = "in1_crit_alarm",
                                  .dep_prop =
                                      &(DepProperty){
                                          .curr_idx = BUSVOLTAGE,
                                          .lmt_idx = CRITHIGHBUSVOLTAGE,
                                          .cond = GREATERTHENEQUALTO } },
    [CRITHIGHPWR] = { .name = "CRITICAL HIGH POWER LIMIT",
                      .data_type = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .prop_type = PROP_TYPE_CONFIG,
                      .units = "microWatt",
                      .sysfname = "power1_crit",
                      .dep_prop = NULL },
    [CRITHIGHPWRALARM] = { .name = "CRITICAL HIGH POWER",
                           .data_type = TYPE_BOOL,
                           .perm = PERM_RD,
                           .available = PROP_AVAIL,
                           .prop_type = PROP_TYPE_ALERT,
                           .units = "NA",
                           .sysfname = "power1_crit_alarm",
                           .dep_prop =
                               &(DepProperty){ .curr_idx = POWER,
                                               .lmt_idx = CRITHIGHPWR,
                                               .cond = GREATERTHENEQUALTO } },
    [UPDATEINTERVAL] = { .name = "DATA CONVERSION TIME",
                         .data_type = TYPE_INT32,
                         .perm = PERM_RD | PERM_WR,
                         .available = PROP_AVAIL,
                         .prop_type = PROP_TYPE_CONFIG,
                         .units = "NA",
                         .sysfname = "update_interval",
                         .dep_prop = NULL }
};

static const DrvDBFxnTable *get_fxn_tbl(Device *p_dev) {
    if (IF_SYSFS_SUPPORT(p_dev->sysfile)) {
        return drvr_sysfs_get_fxn_tbl();
    } else {
        return &drvr_ina226_fxn_table;
    }
}

int ina226_get_irq_type(int pidx, uint8_t *alertstate) {
    int ret = 0;
    if ((pidx == SHUNTVOLTAGECRITLOWALARM) ||
        (pidx == BUSVOLTAGECRITLOWALARM) || (pidx == T3MINALARM)) {
        *alertstate = ALARM_STATE_LOW_ALARM_ACTIVE;
    } else if ((pidx == SHUNTVOLTAGECRITHIGHALARM) ||
               (pidx == BUSVOLTAGECRITHIGHALARM) ||
               (pidx == CRITHIGHPWRALARM)) {
        *alertstate = ALARM_STATE_CRIT_ALARM_ACTIVE;
    } else {
        ret = -1;
    }
    return ret;
}

int ina226_read_prop_count(Device *p_dev, uint16_t *count) {
    int ret = 0;
    *count = g_property_count;
    return 0;
}

int ina226_read_properties(Device *p_dev, void *prop) {
    int ret = 0;
    if (prop) {
        memset(prop, '\0', sizeof(Property) * g_property_count);
        memcpy(prop, g_property, sizeof(Property) * g_property_count);
    } else {
        ret = -1;
    }
    return ret;
}

void ina226_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //ina226_read(void* p_dev, void* prop, void* data );
}

int ina226_reg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = fun;
    }
    return ret;
}

int ina226_dreg_cb(void *p_dev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensor_cb = NULL;
    }
    return ret;
}

int ina226_init(Device *p_dev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(p_dev, &g_property,
                                            &g_property_count);
    if (ret) {
        g_property_count = MAXINAPROP;
        g_property = ina226_property;
        log_debug("INA226: Using static property table with %d property.",
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

int ina226_registration(Device *p_dev) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_registration(drvr, p_dev);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_configure(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_configure(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_read(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_read(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_write(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_write(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_enable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_enable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_disable(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

int ina226_enable_irq(void *p_dev, void *prop, void *data) {
    //TODO: check if IRQ has to enable and disabled here or in driver layer below it.
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

int ina226_disable_irq(void *p_dev, void *prop, void *data) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_disable_irq(drvr, p_dev, g_property, *(int *)prop, data);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}

/* Reading and confirming interrupts for ADT7481 device */
int ina226_confirm_irq(Device *p_dev, AlertCallBackData **acbdata, char *fpath,
                       int *evt) {
    int ret = 0;
    const DrvDBFxnTable *drvr = get_fxn_tbl(p_dev);
    if (drvr) {
        ret = dhelper_confirm_irq(drvr, p_dev, g_property, acbdata, fpath,
                                  MAXINAPROP, evt);
    } else {
        ret = ERR_UBSP_DEV_DRVR_MISSING;
    }
    return ret;
}
