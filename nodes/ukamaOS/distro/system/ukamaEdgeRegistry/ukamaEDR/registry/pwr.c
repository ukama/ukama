/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/pwr.h"

#include "dmt.h"
#include "headers/errorcode.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/power.h"
#include "headers/objects/voltage.h"
#include "headers/utils/log.h"

const AlarmSensorData svolt_alarm_data = {
    .plowthresholdidx = CRITLOWBUSVOLTAGE,
    .phighthresholdidx = PROPERY_IDX_NOT_APPLICABLE,
    .pcrithresholdidx = CRITHIGHBUSVOLTAGE,
    .psensorvalueidx = BUSVOLTAGE,
    .plowlimitalarmidx = BUSVOLTAGECRITLOWALARM,
    .phighlimitalarmidx = PROPERY_IDX_NOT_APPLICABLE,
    .pcrilimitalarmidx = BUSVOLTAGECRITHIGHALARM
};

const AlarmSensorData spwr_alarm_data = {
    .plowthresholdidx = PROPERY_IDX_NOT_APPLICABLE,
    .phighthresholdidx = PROPERY_IDX_NOT_APPLICABLE,
    .pcrithresholdidx = CRITHIGHPWR,
    .psensorvalueidx = POWER,
    .plowlimitalarmidx = PROPERY_IDX_NOT_APPLICABLE,
    .phighlimitalarmidx = PROPERY_IDX_NOT_APPLICABLE,
    .pcrilimitalarmidx = CRITHIGHPWRALARM
};

const DBFxnTable pwrdb_fxn_tbl = {
    .db_add_dev_to_reg = drdb_add_pwr_dev_to_reg,
    .db_read_data_from_dev = drdb_read_pwr_inst_data_from_dev,
    .db_write_data_from_dev = drdb_write_pwr_inst_data_from_dev,
    .db_search_inst_in_reg = NULL,
    .db_free_inst_data_from_reg = free_pwr_data,
    .db_update_inst_in_reg = drdb_update_pwr_inst_data,
    .db_exec = drdb_exec_pwr_inst_rsrc
};

static void populate_ina_resource_id(void **pdata) {
    GenPwrData *data = *pdata;
    data->value.resourceId = RES_M_SENSOR_VALUE;
    data->min.resourceId = RES_O_MIN_MEASURED_VALUE;
    data->max.resourceId = RES_O_MAX_MEASURED_VALUE;
    data->minrange.resourceId = RES_O_MIN_RANGE_VALUE;
    data->maxrange.resourceId = RES_O_MAX_RANGE_VALUE;
    data->units.resourceId = RES_M_SENSOR_UNITS;
    data->resetcounter.resourceId = RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE;
    data->calibration.resourceId = RES_O_CURR_CALIBRATION_VALUE;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void free_pwr_data(void *pdata) {
    if (pdata) {
        GenPwrData *data = pdata;
        free_sdata(&data->value.prop);
        free_sdata(&data->min.prop);
        free_sdata(&data->max.prop);
        free_sdata(&data->minrange.prop);
        free_sdata(&data->maxrange.prop);
        free_sdata(&data->avg.prop);
        free_sdata(&data->cumm.prop);
        free_sdata(&data->counter.prop);
        free_sdata(&data->units.prop);
        free_sdata(&data->calibration.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

void *copy_pwr_data(void *pdata) {
    GenPwrData *data = pdata;
    GenPwrData *ndata = NULL;
    if (data) {
        ndata = dmt_malloc(sizeof(GenPwrData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(GenPwrData));
            /* Try deep  copy for properties of pdata now */
            ndata->value.prop = reg_copy_pdata_prop(data->value.prop);
            ndata->min.prop = reg_copy_pdata_prop(data->min.prop);
            ndata->max.prop = reg_copy_pdata_prop(data->max.prop);
            ndata->minrange.prop = reg_copy_pdata_prop(data->minrange.prop);
            ndata->maxrange.prop = reg_copy_pdata_prop(data->maxrange.prop);
            ndata->avg.prop = reg_copy_pdata_prop(data->avg.prop);
            ndata->cumm.prop = reg_copy_pdata_prop(data->cumm.prop);
            ndata->counter.prop = reg_copy_pdata_prop(data->counter.prop);
            ndata->units.prop = reg_copy_pdata_prop(data->units.prop);
            ndata->calibration.prop =
                reg_copy_pdata_prop(data->calibration.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

#if 0
static void drdb_register_volt_alarms(DevObj* obj, Property* prop, int pcount, uint16_t inst) {
	int ret = 0;
	DevIDT devidt = {0};
	AlarmPropertyData *pdata = NULL;
	int critid = CRITHIGHBUSVOLTAGE;
	int lowid = CRITLOWBUSVOLTAGE;
	int svalid = BUSVOLTAGE;
	reg_initialize_dev_idt(&devidt, inst, OBJECT_ID_VOLT, RES_M_SENSOR_VALUE);
	pdata = reg_initialize_alarm_prop(prop, &critid, NULL, &lowid, &svalid);
	drdb_add_alarm_inst_to_reg(obj, &devidt, prop, pdata);

}

static void drdb_register_pwr_alarms(DevObj* obj, Property* prop, int pcount, uint16_t inst) {
	DevIDT devidt = {0};
	AlarmPropertyData *pdata = NULL;
	int critid = CRITHIGHPWR;
	int svalid = POWER;
	reg_initialize_dev_idt(&devidt, inst, OBJECT_ID_PWR, RES_M_SENSOR_VALUE);
	pdata = reg_initialize_alarm_prop(prop, &critid, NULL, NULL, &svalid);
	drdb_add_alarm_inst_to_reg(obj, &devidt, prop, pdata);
}

#endif

/* To do get this from the property.json */
void get_volt_range(double *max, double *min) {
    *min = 1900;
    *max = 12000;
}

/* To do get this from the property.json */
void get_curr_range(double *max, double *min) {
    *min = 200;
    *max = 7000;
}

/* To do get this from the property.json */
void get_pwr_range(double *max, double *min) {
    *min = 30000;
    *max = 600000;
}

void drdb_update_pwr_inst_data(double val, PData *min, PData *max, PData *avg,
                               PData *cumm, PData *count) {
    min->value.doubleval = MIN(min->value.doubleval, val);
    max->value.doubleval = MAX(max->value.doubleval, val);
    cumm->value.doubleval = val + cumm->value.doubleval;
    count->value.intval++;
    avg->value.doubleval = cumm->value.doubleval / count->value.intval;
}

int drdb_reset_pwr_counters(void *data) {
    int ret = 0;
    DRDBSchema *reg = data;
    GenPwrData *rdata = reg->data;
    if (rdata) {
        rdata->min.value.doubleval = 0;
        rdata->max.value.doubleval = 0;
        rdata->cumm.value.doubleval = 0;
        rdata->counter.value.intval = 0;
        rdata->avg.value.doubleval = 0;

        /* Update the instance in registry.*/
        ret = reg_update_dev(reg);
    } else {
        ret = -1;
    }
    return ret;
}

void drdb_add_volt_inst_to_reg(Device *dev, Property *prop, uint8_t inst,
                               int pidx) {
    int ret = 0;
    VoltData *data = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = dmt_malloc(sizeof(VoltData));
        if (data) {
            memset(data, '\0', sizeof(VoltData));

            if (prop && prop[pidx].available) {
                /* Sensor value */
                data->value.prop = dmt_malloc(sizeof(Property));
                if (data->value.prop) {
                    memcpy(data->value.prop, &prop[pidx], sizeof(Property));
                }

                /* Units */
                memcpy(data->units.value.stringval, prop[pidx].units,
                       strlen(prop[pidx].units));

                /* Application type */
                memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                       strlen(reg->obj.disc));

                /* Calibration or offset */
                /* TODO: Add property to property table and add sysfs file if required.*/
                reg_data_add_property(CALIBRATION, prop, &data->calibration);
            } else {
                ret = -1;
                goto cleanup;
            }

            /* Read min max*/
            get_volt_range(&data->maxrange.value.doubleval,
                           &data->minrange.value.doubleval);

            /* Resource with execute permission. */
            data->resetcounter.execFunc = &drdb_reset_pwr_counters;

        } else {
            ret = -1;
            log_debug("Err: VOLT:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &pwrdb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_VOLT;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_ina_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);

        /* Check if we have any alert property to be registered in registry*/
        reg_register_sensor_alarms(&dev->obj, prop, &svolt_alarm_data, inst,
                                   OBJECT_ID_VOLT, RES_M_SENSOR_VALUE);

        /* Enable alarms */
        reg_enable_alarms(&dev->obj, &svolt_alarm_data);

        log_debug("VOLT:: Device Name: %s, Disc: %s, Type: %d, Module Id %s"
                  " Instance %d is added to DB.",
                  reg->obj.name, reg->obj.disc, reg->obj.type,
                  reg->obj.mod_UUID, reg->instance);
    }
cleanup:
    if (ret) {
        free_pwr_data(data);
        dmt_free(data);
    }
    free_reg(&reg);
}

void drdb_add_curr_inst_to_reg(Device *dev, Property *prop, uint8_t inst,
                               int pidx) {
    int ret = 0;
    CurrData *data = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = dmt_malloc(sizeof(CurrData));
        if (data) {
            memset(data, '\0', sizeof(CurrData));

            if (prop && prop[pidx].available) {
                /* Sensor value */
                data->value.prop = dmt_malloc(sizeof(Property));
                if (data->value.prop) {
                    memcpy(data->value.prop, &prop[pidx], sizeof(Property));
                }

                /* Units */
                memcpy(data->units.value.stringval, prop[pidx].units,
                       strlen(prop[pidx].units));

                /* Application type */
                memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                       strlen(reg->obj.disc));

                /* Calibration or offset */
                /* TODO: Add property to property table and add sysfs file if required.*/
                reg_data_add_property(CALIBRATION, prop, &data->calibration);
            } else {
                ret = -1;
                goto cleanup;
            }

            /* Read min max*/
            get_curr_range(&data->maxrange.value.doubleval,
                           &data->minrange.value.doubleval);

            /* Resource with execute permission. */
            data->resetcounter.execFunc = &drdb_reset_pwr_counters;

        } else {
            ret = -1;
            log_debug("Err: CURR:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &pwrdb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_CURR;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_ina_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);

        log_debug("CURR:: Device Name: %s, Disc: %s, Type: %d, Module Id %s"
                  " Instance %d is added to DB.",
                  reg->obj.name, reg->obj.disc, reg->obj.type,
                  reg->obj.mod_UUID, reg->instance);
    }
cleanup:
    if (ret) {
        free_pwr_data(data);
        dmt_free(data);
    }
    free_reg(&reg);
}

void drdb_add_pwr_inst_to_reg(Device *dev, Property *prop, uint8_t inst,
                              int pidx) {
    int ret = 0;
    PwrData *data = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = dmt_malloc(sizeof(PwrData));
        if (data) {
            memset(data, '\0', sizeof(PwrData));

            if (prop && prop[pidx].available) {
                /* Sensor value */
                data->value.prop = dmt_malloc(sizeof(Property));
                if (data->value.prop) {
                    memcpy(data->value.prop, &prop[pidx], sizeof(Property));
                }

                /* Units */
                memcpy(data->units.value.stringval, prop[pidx].units,
                       strlen(prop[pidx].units));

                /* Application type */
                memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                       strlen(reg->obj.disc));

                /* Calibration or offset */
                /* TODO: Add property to property table and add sysfs file if required.*/
                reg_data_add_property(CALIBRATION, prop, &data->calibration);
            } else {
                ret = -1;
                goto cleanup;
            }

            /* Read min max*/
            get_pwr_range(&data->maxrange.value.doubleval,
                          &data->minrange.value.doubleval);

            /* Resource with execute permission. */
            data->resetcounter.execFunc = &drdb_reset_pwr_counters;

        } else {
            ret = -1;
            log_debug("Err: PWR:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &pwrdb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_PWR;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_ina_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);

        /* Check if we have any alert property to be registered in registry*/
        reg_register_sensor_alarms(&dev->obj, prop, &spwr_alarm_data, inst,
                                   OBJECT_ID_PWR, RES_M_SENSOR_VALUE);

        /* Enable alarms */
        reg_enable_alarms(&dev->obj, &spwr_alarm_data);

        log_debug("PWR:: Device Name: %s, Disc: %s, Type: %d, Module Id %s"
                  " Instance %d is added to DB.",
                  reg->obj.name, reg->obj.disc, reg->obj.type,
                  reg->obj.mod_UUID, reg->instance);
    }
cleanup:
    if (ret) {
        free_pwr_data(data);
        dmt_free(data);
    }
    free_reg(&reg);
}

void drdb_add_pwr_dev_to_reg(void *pdev) {
    int ret = 0;
    Property *prop = NULL;
    Device *dev = NULL;
    int pcount = 0;
    if (pdev) {
        dev = pdev;
        prop = db_read_dev_property(&dev->obj, &pcount);
        for (uint8_t iter = 0; iter < 3; iter++) {
            uint8_t inst = 0;
            switch (iter) {
            case 0: {
                inst = list_size(reg_getdb(OBJ_TYPE_VOLT));
                drdb_add_volt_inst_to_reg(dev, prop, inst, BUSVOLTAGE);
                break;
            }
            case 1: {
                inst = list_size(reg_getdb(OBJ_TYPE_CURR));
                drdb_add_curr_inst_to_reg(dev, prop, inst, CURRENT);
                break;
            }
            case 2:
                inst = list_size(reg_getdb(OBJ_TYPE_PWR));
                drdb_add_pwr_inst_to_reg(dev, prop, inst, POWER);
                break;
            default:
                return;
            }
        }

        /* Check if alarm exists*/
        if (!ret) {
            if (reg_check_if_alarm_property_exist(prop, pcount)) {
                /* Register alarm callback if alarm property exist */
                ret = db_register_alarm_callback(&dev->obj, 0, &drdb_alarm_cb);
                log_debug("INA:: Registered Callback function for "
                          "Device name %s Disc: %s Module UUID %s.",
                          dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
            }
        }
    }
    dmt_free(prop);
}

int drdb_read_pwr_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    GenPwrObjInfo *msgdata = NULL;
    if (reg->data) {
        GenPwrData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            dmt_free(rqmsg->data);
        }
        msgdata = dmt_malloc(sizeof(GenPwrObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(GenPwrObjInfo));

        int propid = 0;
        /* Sensor Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->value.resourceId),
            &msgdata->sensor_value, &rdata->value, &size);

        /* Min Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->min.resourceId),
            &msgdata->min_measured_value, &rdata->min, &size);

        /* Max Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->max.resourceId),
            &msgdata->max_measured_value, &rdata->max, &size);

        /* Avg value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->avg.resourceId),
            &msgdata->avg_measured_value, &rdata->avg, &size);

        /* Min Range value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->minrange.resourceId),
            &msgdata->min_range_value, &rdata->minrange, &size);

        /* Max Range value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->maxrange.resourceId),
            &msgdata->max_range_value, &rdata->maxrange, &size);

        /* Calibration value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->calibration.resourceId),
            &msgdata->calibration_value, &rdata->calibration, &size);

        /* Units */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->units.resourceId),
                                   msgdata->sensor_units, &rdata->units, &size);

        /* Application Type */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->applicationtype.resourceId),
            msgdata->application_type, &rdata->applicationtype, &size);

        /* Instance Id */
        msgdata->instanceId = reg->instance;

        /* Check if was single prop read or full struct read */
        if (rqmsg->resourceId == ALL_RESOURCE_ID) {
            rqmsg->datasize = sizeof(GenPwrObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(GenPwrObjInfo);
        }
        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("TMP:: Reading request for resource %d for Device %s Disc: %s"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    return ret;
}

int drdb_write_pwr_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = -1;
    GenPwrData *rdata = reg->data;
    if (reg->data && rqmsg->data) {
        GenPwrObjInfo *msgdata = rqmsg->data;

        /* Calibration Value */
        ret = db_write_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->calibration.resourceId),
            msgdata->calibration_value, &rdata->calibration);

        if (ret) {
            log_trace("Err(%d): PWR:: Write filed for resource %d for "
                      "Device %s Disc: %s Module ID %s is not writable",
                      ret, rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                      reg->obj.mod_UUID);
        } else {
            log_trace("PWR:: Write done for resource %d for Device %s Disc: %s"
                      "Module ID %s is not writable",
                      rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                      reg->obj.mod_UUID);
        }
    }
    return ret;
}

int drdb_exec_pwr_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    GenPwrObjInfo *msgdata = rqmsg->data;
    if (reg->data) {
        GenPwrData *rdata = reg->data;
        int propid = 0;
        if (rqmsg->resourceId == rdata->value.resourceId) {
            if (rdata->resetcounter.execFunc) {
                ret = rdata->resetcounter.execFunc(reg);
            } else {
                ret = ERR_EDGEREG_RESR_NOTIMPL;
            }
        }
        log_trace("PWR::Executed resource %d for Device %s Disc: %s"
                  "Module ID %s",
                  rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                  reg->obj.mod_UUID);
    }
    return ret;
}
