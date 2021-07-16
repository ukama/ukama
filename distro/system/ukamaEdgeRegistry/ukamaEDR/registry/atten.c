/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/atten.h"
#include "dmt.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/atten.h"
#include "headers/utils/log.h"

const DBFxnTable attdb_fxn_tbl = { .db_add_dev_to_reg = drdb_add_att_dev_to_reg,
                                   .db_read_data_from_dev =
                                       drdb_read_att_inst_data_from_dev,
                                   .db_write_data_from_dev =
                                       drdb_write_att_inst_data_from_dev,
                                   .db_search_inst_in_reg = NULL,
                                   .db_free_inst_data_from_reg = free_att_data,
                                   .db_update_inst_in_reg = NULL,
                                   .db_exec = NULL };

void free_att_data(void *pdata) {
    if (pdata) {
        AttData *data = pdata;
        free_sdata(&data->attvalue.prop);
        free_sdata(&data->latchenable.prop);
        free_sdata(&data->minrange.prop);
        free_sdata(&data->maxrange.prop);
        free_sdata(&data->units.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

/* Used by List to copy nodes data */
void *copy_att_data(void *pdata) {
    AttData *data = pdata;
    AttData *ndata = NULL;
    if (data) {
        ndata = dmt_malloc(sizeof(AttData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(AttData));
            /* Try deep  copy for properties of pdata now */
            ndata->attvalue.prop = reg_copy_pdata_prop(data->attvalue.prop);
            ndata->latchenable.prop =
                reg_copy_pdata_prop(data->latchenable.prop);
            ndata->minrange.prop = reg_copy_pdata_prop(data->minrange.prop);
            ndata->maxrange.prop = reg_copy_pdata_prop(data->maxrange.prop);
            ndata->units.prop = reg_copy_pdata_prop(data->units.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

static void populate_atten_resource_id(void **pdata) {
    AttData *data = *pdata;
    data->attvalue.resourceId = RES_M_ATTVALUE;
    data->latchenable.resourceId = RES_M_LATCH;
    data->minrange.resourceId = RES_M_MINRANGE;
    data->maxrange.resourceId = RES_M_MAXRANGE;
    data->units.resourceId = RES_M_SENSOR_UNITS;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void drdb_add_att_inst_to_reg(Device *dev, uint8_t inst, uint8_t subdev) {
    int ret = 0;
    int pcount = 0;
    int pidx = 0;
    AttData *data = NULL;
    Property *prop = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = dmt_malloc(sizeof(AttData));
        if (data) {
            memset(data, '\0', sizeof(AttData));
            /* Read Property table of the sensor. */
            prop = db_read_dev_property(&dev->obj, &pcount);
            if (prop) {
                /* Attenuation Value */
                reg_data_add_property(ATTVALUE, prop, &data->attvalue);

                /* Units */
                memcpy(data->units.value.stringval, prop[pidx].units,
                       strlen(prop[pidx].units));

                /* Latch enable */
                reg_data_add_property(LATCHENABLE, prop, &data->latchenable);

            } else {
                ret = -1;
                goto cleanup;
            }

            /* Application type */
            memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                   strlen(reg->obj.disc));

            /* Min and Max settings */
            data->minrange.value.intval = MINRANGE_ATTVALUE; //0dB
            data->maxrange.value.intval = MAXRANGE_ATTVALUE; //2*63dB
        } else {
            ret = -1;
            log_debug("Err: ATTEN:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &attdb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_ATT;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_atten_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug(
            "ATTEN:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d is added to DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
cleanup:
    if (ret) {
        free_att_data(data);
        dmt_free(data);
    }
    dmt_free(prop);
    free_reg(&reg);
}

void drdb_add_att_dev_to_reg(void *pdev) {
    if (pdev) {
        Device *dev = pdev;
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_ATT));
        drdb_add_att_inst_to_reg(dev, inst, 0);
    }
}

int drdb_read_att_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    AttObjInfo *msgdata = NULL;
    if (reg->data) {
        AttData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            dmt_free(rqmsg->data);
        }
        msgdata = dmt_malloc(sizeof(AttObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(AttObjInfo));

        int propid = 0;

        /* Value */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->attvalue.resourceId),
                                   &msgdata->attvalue, &rdata->attvalue, &size);

        /* Min Value */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->minrange.resourceId),
                                   &msgdata->minrange, &rdata->minrange, &size);

        /* Max Value */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->maxrange.resourceId),
                                   &msgdata->maxrange, &rdata->maxrange, &size);

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
            rqmsg->datasize = sizeof(AttObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(AttObjInfo);
        }

        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("ATT:: Reading request for resource %d for Device %s Disc: %s"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    return ret;
}

int drdb_write_att_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = -1;
    int size = 0;
    AttData *rdata = reg->data;
    if (reg->data && rqmsg->data) {
        AttObjInfo *msgdata = rqmsg->data;
        /* Attenuation value */
        ret = db_write_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->attvalue.resourceId),
                                   msgdata->attvalue, &rdata->attvalue);

        /* Latch enable */
        ret |= db_write_intval_prop(&reg->obj,
                                    (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                        (rqmsg->resourceId ==
                                         rdata->latchenable.resourceId),
                                    msgdata->latchenable, &rdata->latchenable);

        if (ret) {
            log_trace(
                "Err(%d): ATT:: Write filed for resource %d for Device %s Disc: %s"
                "Module ID %s is not writable",
                ret, rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                reg->obj.mod_UUID);
        } else {
            log_trace("ATT:: Write done for resource %d for Device %s Disc: %s"
                      "Module ID %s is not writable",
                      rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                      reg->obj.mod_UUID);
        }
    }
    return ret;
}
