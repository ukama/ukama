/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "devices/bsp_ina226.h"

#include "devhelper.h"
#include "errorcode.h"
#include "irqdb.h"
#include "irqhelper.h"
#include "property.h"
#include "drivers/ina226_wrapper.h"
#include "drivers/sysfs_wrapper.h"

#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"

static SensorCallbackFxn sensorCb;

const DrvrOps ina226WrapperOps = { .init = ina226_wrapper_init,
                                   .configure = ina226_wrapper_configure,
                                   .read = ina226_wrapper_read,
                                   .write = ina226_wrapper_write,
                                   .enable = ina226_wrapper_enable,
                                   .disable = ina226_wrapper_disable,
                                   .registerCb = ina226_wrapper_reg_cb,
                                   .dregisterCb = ina226_wrapper_dreg_cb,
                                   .enableIrq = ina226_wrapper_enable_irq,
                                   .disableIrq = ina226_wrapper_disable_irq };

static Property *gProperty = NULL;
static int gPropertyCount = 0;

static Property ina226Property[MAXINAPROP] = {
    [SHUNTVOLTAGE] = { .name = "SHUNT VOLTAGE",
                       .dataType = TYPE_INT32,
                       .perm = PERM_RD,
                       .available = PROP_AVAIL,
                       .propType = PROP_TYPE_STATUS,
                       .units = "milliVolts",
                       .sysFname = "in0_input",
                       .depProp = NULL },
    [BUSVOLTAGE] = { .name = "BUS VOLTAGE",
                     .dataType = TYPE_INT32,
                     .perm = PERM_RD,
                     .available = PROP_AVAIL,
                     .propType = PROP_TYPE_STATUS,
                     .units = "milliVolts",
                     .sysFname = "in1_input",
                     .depProp = NULL },
    [CURRENT] = { .name = "CURRENT",
                  .dataType = TYPE_INT32,
                  .perm = PERM_RD,
                  .available = PROP_AVAIL,
                  .propType = PROP_TYPE_STATUS,
                  .units = "milliAmp",
                  .sysFname = "curr1_input",
                  .depProp = NULL },
    [POWER] = { .name = "POWER",
                .dataType = TYPE_INT32,
                .perm = PERM_RD,
                .available = PROP_AVAIL,
                .propType = PROP_TYPE_STATUS,
                .units = "microWatt",
                .sysFname = "power1_input",
                .depProp = NULL },
    [SHUNTRESISTOR] = { .name = "SHUNT RESISTANCE",
                        .dataType = TYPE_INT32,
                        .perm = PERM_RD | PERM_WR,
                        .available = PROP_AVAIL,
                        .propType = PROP_TYPE_CONFIG,
                        .units = "microOhm",
                        .sysFname = "shunt_resistor",
                        .depProp = NULL },
    [CALIBRATION] = { .name = "CALIBRATION",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "",
                      .sysFname = "calibration",
                      .depProp = NULL },
    [CRITLOWSHUNTVOLTAGE] = { .name = "CRIT LOW SHUNT VOLTAGE",
                              .dataType = TYPE_INT32,
                              .perm = PERM_RD | PERM_WR,
                              .available = PROP_AVAIL,
                              .propType = PROP_TYPE_CONFIG,
                              .units = "milliVolts",
                              .sysFname = "in0_lcrit",
                              .depProp = NULL },
    [CRITHIGHSHUNTVOLTAGE] = { .name = "CRIT HIGH SHUNT VOLTAGE",
                               .dataType = TYPE_INT32,
                               .perm = PERM_RD | PERM_WR,
                               .available = PROP_AVAIL,
                               .propType = PROP_TYPE_CONFIG,
                               .units = "milliVolts",
                               .sysFname = "in0_crit",
                               .depProp = NULL },
    [SHUNTVOLTAGECRITLOWALARM] = { .name = "SHUNT VOLTAGE CRIT LOW ALARM",
                                   .dataType = TYPE_BOOL,
                                   .perm = PERM_RD,
                                   .available = PROP_AVAIL,
                                   .propType = PROP_TYPE_ALERT,
                                   .units = "NA",
                                   .sysFname = "in0_lcrit_alarm",
                                   .depProp =
                                       &(DepProperty){
                                           .currIdx = SHUNTVOLTAGE,
                                           .lmtIdx = CRITLOWSHUNTVOLTAGE,
                                           .cond = LESSTHENEQUALTO } },
    [SHUNTVOLTAGECRITHIGHALARM] = { .name = "SHUNT VOLTAGE CRIT HIGH ALARM",
                                    .dataType = TYPE_BOOL,
                                    .perm = PERM_RD,
                                    .available = PROP_AVAIL,
                                    .propType = PROP_TYPE_ALERT,
                                    .units = "NA",
                                    .sysFname = "in0_crit_alarm",
                                    .depProp =
                                        &(DepProperty){
                                            .currIdx = SHUNTVOLTAGE,
                                            .lmtIdx = CRITHIGHSHUNTVOLTAGE,
                                            .cond = GREATERTHENEQUALTO } },
    [CRITLOWBUSVOLTAGE] = { .name = "LOW VOLTAGE LIMIT",
                            .dataType = TYPE_INT32,
                            .perm = PERM_RD | PERM_WR,
                            .available = PROP_AVAIL,
                            .propType = PROP_TYPE_CONFIG,
                            .units = "milliVolts",
                            .sysFname = "in1_lcrit",
                            .depProp = NULL },
    [CRITHIGHBUSVOLTAGE] = { .name = "HIGH VOLTAGE LIMIT",
                             .dataType = TYPE_INT32,
                             .perm = PERM_RD,
                             .available = PROP_AVAIL,
                             .propType = PROP_TYPE_CONFIG,
                             .units = "milliVolts",
                             .sysFname = "in1_crit",
                             .depProp = NULL },
    [BUSVOLTAGECRITLOWALARM] = { .name = "BUS VOLTAGE CRIT LOW ALARM",
                                 .dataType = TYPE_BOOL,
                                 .perm = PERM_RD,
                                 .available = PROP_AVAIL,
                                 .propType = PROP_TYPE_ALERT,
                                 .units = "NA",
                                 .sysFname = "in1_lcrit_alarm",
                                 .depProp = &(
                                     DepProperty){ .currIdx = BUSVOLTAGE,
                                                   .lmtIdx = CRITLOWBUSVOLTAGE,
                                                   .cond = LESSTHENEQUALTO } },
    [BUSVOLTAGECRITHIGHALARM] = { .name = "BUS VOLTAGE CRIT HIGH ALARM",
                                  .dataType = TYPE_BOOL,
                                  .perm = PERM_RD,
                                  .available = PROP_AVAIL,
                                  .propType = PROP_TYPE_ALERT,
                                  .units = "NA",
                                  .sysFname = "in1_crit_alarm",
                                  .depProp =
                                      &(DepProperty){
                                          .currIdx = BUSVOLTAGE,
                                          .lmtIdx = CRITHIGHBUSVOLTAGE,
                                          .cond = GREATERTHENEQUALTO } },
    [CRITHIGHPWR] = { .name = "CRITICAL HIGH POWER LIMIT",
                      .dataType = TYPE_INT32,
                      .perm = PERM_RD | PERM_WR,
                      .available = PROP_AVAIL,
                      .propType = PROP_TYPE_CONFIG,
                      .units = "microWatt",
                      .sysFname = "power1_crit",
                      .depProp = NULL },
    [CRITHIGHPWRALARM] = { .name = "CRITICAL HIGH POWER",
                           .dataType = TYPE_BOOL,
                           .perm = PERM_RD,
                           .available = PROP_AVAIL,
                           .propType = PROP_TYPE_ALERT,
                           .units = "NA",
                           .sysFname = "power1_crit_alarm",
                           .depProp =
                               &(DepProperty){ .currIdx = POWER,
                                               .lmtIdx = CRITHIGHPWR,
                                               .cond = GREATERTHENEQUALTO } },
    [UPDATEINTERVAL] = { .name = "DATA CONVERSION TIME",
                         .dataType = TYPE_INT32,
                         .perm = PERM_RD | PERM_WR,
                         .available = PROP_AVAIL,
                         .propType = PROP_TYPE_CONFIG,
                         .units = "NA",
                         .sysFname = "update_interval",
                         .depProp = NULL }
};

static const DrvrOps* get_fxn_tbl(Device *pDev) {
    if (IF_SYSFS_SUPPORT(pDev->sysFile)) {
        return sysfs_wrapper_get_ops();
    } else {
        return &ina226WrapperOps;
    }
}

int bsp_ina226_get_irq_type(int pidx, uint8_t *alertstate) {
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

int bsp_ina226_read_prop_count(Device *pDev, uint16_t *count) {
    int ret = 0;
    *count = gPropertyCount;
    return 0;
}

int bsp_ina226_read_properties(Device *pDev, void *prop) {
    int ret = 0;
    if (prop) {
        usys_memset(prop, '\0', sizeof(Property) * gPropertyCount);
        usys_memcpy(prop, gProperty, sizeof(Property) * gPropertyCount);
    } else {
        ret = -1;
    }
    return ret;
}

void ina226_irq_callback(DevObj *obj, void *prop, void *data) {
    /* Read and confirm the IRQ's for device.*/
    //ina226_read(void* pDev, void* prop, void* data );
}

int bsp_ina226_reg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensorCb = fun;
    }
    return ret;
}

int bsp_ina226_dreg_cb(void *pDev, SensorCallbackFxn fun) {
    int ret = 0;
    if (fun) {
        sensorCb = NULL;
    }
    return ret;
}

int bsp_ina226_init(Device *pDev) {
    int ret = 0;
    ret = dhelper_init_property_from_parser(pDev, &gProperty, &gPropertyCount);
    if (ret) {
        gPropertyCount = MAXINAPROP;
        gProperty = ina226Property;
        log_debug("INA226: Using static property table with %d property.",
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

int bsp_ina226_registration(Device *pDev) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_registration(drvr, pDev);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_configure(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_configure(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_read(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_read(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_write(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_write(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_enable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_disable(void *pDev, void *prop, void *data) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_disable(drvr, pDev, gProperty, *(int *)prop, data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_enable_irq(void *pDev, void *prop, void *data) {
    //TODO: check if IRQ has to enable and disabled here or in driver layer below it.
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_enable_irq(drvr, sensorCb, pDev, gProperty, *(int *)prop,
                                 data);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}

int bsp_ina226_disable_irq(void *pDev, void *prop, void *data) {
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
int bsp_ina226_confirm_irq(Device *pDev, AlertCallBackData **acbdata,
                           char *fpath, int *evt) {
    int ret = 0;
    const DrvrOps *drvr = get_fxn_tbl(pDev);
    if (drvr) {
        ret = dhelper_confirm_irq(drvr, pDev, gProperty, acbdata, fpath,
                                  MAXINAPROP, evt);
    } else {
        ret = ERR_NODED_DEV_DRVR_MISSING;
    }
    return ret;
}
