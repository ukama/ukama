/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/adc.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/analog_output.h"
#include "headers/utils/log.h"

const DBFxnTable adcdb_fxn_tbl = { .db_add_dev_to_reg = drdb_add_adc_dev_to_reg,
                                   .db_read_data_from_dev =
                                       drdb_read_adc_inst_data_from_dev,
                                   .db_write_data_from_dev =
                                       drdb_write_adc_inst_data_from_dev,
                                   .db_search_inst_in_reg = NULL,
                                   .db_free_inst_data_from_reg = free_adc_data,
                                   .db_update_inst_in_reg = NULL,
                                   .db_exec = NULL };

void free_adc_data(void *pdata) {
    if (pdata) {
        AdcData *data = pdata;
        free_sdata(&data->outputcurr.prop);
        free_sdata(&data->minrange.prop);
        free_sdata(&data->maxrange.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

void *copy_adc_data(void *pdata) {
    AdcData *data = pdata;
    AdcData *ndata = NULL;
    if (data) {
        ndata = malloc(sizeof(AdcData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(AdcData));
            /* Try deep  copy for properties of pdata now */
            ndata->outputcurr.prop = reg_copy_pdata_prop(data->outputcurr.prop);
            ndata->maxrange.prop = reg_copy_pdata_prop(data->maxrange.prop);
            ndata->minrange.prop = reg_copy_pdata_prop(data->minrange.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

static void populate_adc_resource_id(void **pdata) {
    AdcData *data = *pdata;
    data->outputcurr.resourceId = RES_M_OUT_CURR_VALUE;
    data->minrange.resourceId = RES_O_MIN_RANGE_VALUE;
    data->maxrange.resourceId = RES_O_MAX_RANGE_VALUE;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

int drdb_add_ads1015_inst_to_reg(Device *dev, Property *prop, uint8_t inst,
                                 uint8_t subdev) {
    int ret = 0;
    int pcount = 0;
    int pidx = 0;
    AdcData *data = NULL;
    DRDBSchema *reg = malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = malloc(sizeof(AdcData));
        if (data) {
            memset(data, '\0', sizeof(AdcData));
            /* If Property table of the sensor exist. */
            if (prop) {
                /* Sensor instance for ADC is translated to Property table */
                pidx = subdev;
            } else {
                ret = -1;
                goto cleanup;
            }

            /* Sensor value */
            reg_data_add_property(pidx, prop, &data->outputcurr);

            /* Application type */
            memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                   strlen(reg->obj.disc));

            //TODO: Check values an add it as a part of property json.*/
            data->minrange.value.doubleval = MINRANGE_ADCVALUE;
            data->maxrange.value.doubleval = MAXRANGE_ADCVALUE;

        } else {
            ret = -1;
            log_debug("Err: ADC:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &adcdb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_ADC;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_adc_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug(
            "DRDB:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d is added to DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
cleanup:
    if (ret) {
        free_adc_data(data);
        UKAMA_FREE(data);
    }
    free_reg(&reg);
    return ret;
}

void drdb_add_adc_dev_to_reg(void *pdev) {
    int ret = 0;
    int pcount = 0;
    Property *prop = NULL;
    if (pdev) {
        Device *dev = pdev;
        prop = db_read_dev_property(&dev->obj, &pcount);
        for (uint8_t iter = 0; iter < 8; iter++) {
            uint8_t inst = list_size(reg_getdb(OBJ_TYPE_ADC));
            ret |= drdb_add_ads1015_inst_to_reg(dev, prop, inst, iter);
        }
        /* Register alert callback if alert property exist */
        if (!ret) {
        }
    }
    UKAMA_FREE(prop);
}

int drdb_read_adc_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    AdcObjInfo *msgdata = NULL;
    if (reg->data) {
        AdcData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            UKAMA_FREE(rqmsg->data);
        }

        msgdata = malloc(sizeof(AdcObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(AdcObjInfo));

        int propid = 0;

        /* Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->outputcurr.resourceId),
            &msgdata->outputcurr, &rdata->outputcurr, &size);

        /* Min Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->minrange.resourceId),
            &msgdata->minrange, &rdata->minrange, &size);

        /* Max Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->maxrange.resourceId),
            &msgdata->maxrange, &rdata->maxrange, &size);

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
            rqmsg->datasize = sizeof(AdcObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(AdcObjInfo);
        }

        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("ADC:: Reading request for resource %d for Device %s Disc: %s"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    return ret;
}

int drdb_write_adc_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    log_trace("ADC::Resource %d for Device %s Disc: %s"
              "Module ID %s is not writable",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    ret = ERR_EDGEREG_PROP_PERMDENIED;
    return ret;
}
