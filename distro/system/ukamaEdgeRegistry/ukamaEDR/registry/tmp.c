/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/tmp.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/temperature.h"
#include "registry/alarm.h"
#include "headers/utils/log.h"

const AlarmSensorData stmp_alarm_data[3] = {
    { .plowthresholdidx = T1MINLIMIT,
      .phighthresholdidx = T1MAXLIMIT,
      .pcrithresholdidx = T1CRITLIMIT,
      .psensorvalueidx = T1TEMPVALUE,
      .plowlimitalarmidx = T1MINALARM,
      .phighlimitalarmidx = T1MAXALARM,
      .pcrilimitalarmidx = T1CRITALARM },
    { .plowthresholdidx = T2MINLIMIT,
      .phighthresholdidx = T2MAXLIMIT,
      .pcrithresholdidx = T2CRITLIMIT,
      .psensorvalueidx = T2TEMPVALUE,
      .plowlimitalarmidx = T2MINALARM,
      .phighlimitalarmidx = T2MAXALARM,
      .pcrilimitalarmidx = T2CRITALARM },
    { .plowthresholdidx = T3MINLIMIT,
      .phighthresholdidx = T3MAXLIMIT,
      .pcrithresholdidx = T3CRITLIMIT,
      .psensorvalueidx = T3TEMPVALUE,
      .plowlimitalarmidx = T3MINALARM,
      .phighlimitalarmidx = T3MAXALARM,
      .pcrilimitalarmidx = T3CRITALARM },
};

const DBFxnTable tmpdb_fxn_tbl = {
    .db_add_dev_to_reg = drdb_add_tmp_dev_to_reg,
    .db_read_data_from_dev = drdb_read_tmp_inst_data_from_dev,
    .db_write_data_from_dev = drdb_write_tmp_inst_data_from_dev,
    .db_search_inst_in_reg = NULL,
    .db_free_inst_data_from_reg = free_tmp_data,
    .db_update_inst_in_reg = drdb_update_tmp_inst_data,
    .db_exec = drdb_exec_tmp_inst_rsrc
};

void tmpdb_init(DeviceType type) {
    if (type == DEV_TYPE_TMP) {
        reg_init_type(type);
    }
}

void free_tmp_data(void *pdata) {
    if (pdata) {
        TempData *data = pdata;
        free_sdata(&data->value.prop);
        free_sdata(&data->min.prop);
        free_sdata(&data->max.prop);
        free_sdata(&data->minrange.prop);
        free_sdata(&data->maxrange.prop);
        free_sdata(&data->avg.prop);
        free_sdata(&data->cumm.prop);
        free_sdata(&data->counter.prop);
        free_sdata(&data->units.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

void *copy_tmp_data(void *pdata) {
    TempData *data = pdata;
    TempData *ndata = NULL;
    if (data) {
        ndata = malloc(sizeof(TempData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(TempData));
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
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

static void populate_tmp_resource_id(void **pdata) {
    TempData *data = *pdata;
    data->value.resourceId = RES_M_SENSOR_VALUE;
    data->min.resourceId = RES_O_MIN_MEASURED_VALUE;
    data->max.resourceId = RES_O_MAX_MEASURED_VALUE;
    data->minrange.resourceId = RES_O_MIN_RANGE_VALUE;
    data->maxrange.resourceId = RES_O_MAX_RANGE_VALUE;
    data->units.resourceId = RES_M_SENSOR_UNITS;
    data->resetcounter.resourceId = RES_O_RESET_MIN_AND_MAX_MEASURED_VALUE;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void get_tmp_range(char *name, double *max, double *min) {
    if (!strcmp("TMP646", name)) {
        *min = MIN_TMP_RANGE;
        *max = MAX_TMP_RANGE;
    } else if (!strcmp("ADT7481", name)) {
        *min = MIN_TMP_RANGE;
        *max = MAX_TMP_RANGE;
    } else if (!strcmp("SE98", name)) {
        *min = MIN_TMP_RANGE;
        *max = MAX_TMP_RANGE;
    }
}

void drdb_update_tmp_inst_data(double temp, PData *min, PData *max, PData *avg,
                               PData *cumm, PData *count) {
    min->value.doubleval = MIN(min->value.doubleval, temp);
    max->value.doubleval = MAX(max->value.doubleval, temp);
    cumm->value.doubleval = temp + cumm->value.doubleval;
    count->value.intval++;
    avg->value.doubleval = cumm->value.doubleval / count->value.intval;
}

int drdb_reset_tmp_counters(void *data) {
    int ret = 0;
    DRDBSchema *reg = data;
    TempData *rdata = reg->data;
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

int drdb_add_tmp_inst_to_reg(Device *dev, Property *prop, uint8_t instance,
                             uint8_t subdev) {
    int ret = 0;
    int pcount = 0;
    int pidx = 0;
    TempData *data = NULL;
    DRDBSchema *reg = malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        reg->data = malloc(sizeof(TempData));
        if (reg->data) {
            memset(reg->data, '\0', sizeof(TempData));
            /* If Property table of the sensor exist. */
            if (prop) {
                if (subdev == 0) {
                    pidx = T1TEMPVALUE;
                } else if (subdev == 1) {
                    pidx = T2TEMPVALUE;
                } else if (subdev == 2) {
                    pidx = T3TEMPVALUE;
                }
            } else {
                ret = -1;
                goto cleanup;
            }
            data = reg->data;

            /* Sensor Value */
            reg_data_add_property(pidx, prop, &data->value);

            /* Units */
            memcpy(data->units.value.stringval, prop[pidx].units,
                   strlen(prop[pidx].units));

            /* Application type */
            memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                   strlen(reg->obj.disc));

            /* Read min max*/
            get_tmp_range(reg->obj.name, &data->maxrange.value.doubleval,
                          &data->minrange.value.doubleval);

            /* Resource with execute permission. */
            data->resetcounter.execFunc = &drdb_reset_tmp_counters;
        } else {
            ret = -1;
            log_debug("Err: TMP:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &tmpdb_fxn_tbl;
        reg->instance = instance;
        reg->data = data;
        reg->type = OBJ_TYPE_TMP;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_tmp_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);

        /* Check if we have any alert property to be registered in registry*/
        reg_register_sensor_alarms(&dev->obj, prop, &stmp_alarm_data[subdev],
                                   instance, OBJECT_ID_TMP, RES_M_SENSOR_VALUE);

        /* Enable alarms */
        ret = reg_enable_alarms(&dev->obj, &stmp_alarm_data[subdev]);

        log_debug(
            "TMP:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d is added to DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
cleanup:
    if (ret) {
        free_tmp_data(data);
        UKAMA_FREE(data);
    }
    free_reg(&reg);
    return ret;
}

int drdb_add_tmp464_dev_to_reg(Device *dev, Property *prop) {
    int ret = 0;
    for (uint8_t iter = 0; iter < 3; iter++) {
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_TMP));
        ret |= drdb_add_tmp_inst_to_reg(dev, prop, inst, iter);
    }
    return ret;
}

int drdb_add_adt7481_dev_to_reg(Device *dev, Property *prop) {
    int ret = 0;
    for (uint8_t iter = 0; iter < 3; iter++) {
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_TMP));
        ret |= drdb_add_tmp_inst_to_reg(dev, prop, inst, iter);
    }
    return ret;
}

int drdb_add_se98_dev_to_reg(Device *dev, Property *prop) {
    int ret = 0;
    for (uint8_t iter = 0; iter < 1; iter++) {
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_TMP));
        ret |= drdb_add_tmp_inst_to_reg(dev, prop, inst, iter);
    }
    return ret;
}

void drdb_add_tmp_dev_to_reg(void *pdev) {
    int ret = 0;
    Property *prop = NULL;
    Device *dev = NULL;
    int pcount = 0;
    if (pdev) {
        dev = pdev;
        prop = db_read_dev_property(&dev->obj, &pcount);
        if (!strcmp("TMP464", dev->obj.name)) {
            ret = drdb_add_tmp464_dev_to_reg(dev, prop);
        } else if (!strcmp("ADT7481", dev->obj.name)) {
            ret = drdb_add_adt7481_dev_to_reg(dev, prop);
        } else if (!strcmp("SE98", dev->obj.name)) {
            ret = drdb_add_se98_dev_to_reg(dev, prop);
        }

        /* Check if alarm exists*/
        if (!ret) {
            if (reg_check_if_alarm_property_exist(prop, pcount)) {
                /* Register alarm callback if alarm property exist */
                ret = db_register_alarm_callback(&dev->obj, 0, &drdb_alarm_cb);
                log_debug("TMP:: Registered Callback function for "
                          "Device name %s Disc: %s Module UUID %s.",
                          dev->obj.name, dev->obj.disc, dev->obj.mod_UUID);
            }
        }
    }
    UKAMA_FREE(prop);
}

int drdb_read_tmp_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    TempObjInfo *msgdata = NULL;
    if (reg->data) {
        TempData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
    	 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            UKAMA_FREE(rqmsg->data);
        }

        msgdata = malloc(sizeof(TempObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(TempObjInfo));

        int propid = 0;
        /* TEMP Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->value.resourceId),
            &msgdata->sensor_value, &rdata->value, &size);

        /* Min TEMP Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->min.resourceId),
            &msgdata->min_measured_value, &rdata->min, &size);

        /* Max TEMP Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->max.resourceId),
            &msgdata->max_measured_value, &rdata->max, &size);

        /* Avg TEMP value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->avg.resourceId),
            &msgdata->avg_measured_value, &rdata->avg, &size);

        /* Min TEMP Range value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->minrange.resourceId),
            &msgdata->min_range_value, &rdata->minrange, &size);

        /* Max TEMP Range value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->maxrange.resourceId),
            &msgdata->max_range_value, &rdata->maxrange, &size);

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
            rqmsg->datasize = sizeof(TempObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(TempObjInfo);
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

int drdb_write_tmp_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    log_trace("TMP::Resource %d for Device %s Disc: %s"
              "Module ID %s is not writable",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    ret = ERR_EDGEREG_PROP_PERMDENIED;
    return ret;
}

int drdb_exec_tmp_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    TempObjInfo *msgdata = rqmsg->data;
    if (reg->data) {
        TempData *rdata = reg->data;
        int propid = 0;
        if (rqmsg->resourceId == rdata->resetcounter.resourceId) {
            if (rdata->resetcounter.execFunc) {
                ret = rdata->resetcounter.execFunc(reg);
            } else {
                ret = ERR_EDGEREG_RESR_NOTIMPL;
            }
        }
        log_trace("TMP::Executed resource %d for Device %s Disc: %s"
                  "Module ID %s with status %d",
                  rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                  reg->obj.mod_UUID, ret);
    }
    return ret;
}
