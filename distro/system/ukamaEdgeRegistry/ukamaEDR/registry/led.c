/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/led.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/led.h"
#include "headers/utils/log.h"

const DBFxnTable leddb_fxn_tbl = { .db_add_dev_to_reg = drdb_add_led_dev_to_reg,
                                   .db_read_data_from_dev =
                                       drdb_read_led_inst_data_from_dev,
                                   .db_write_data_from_dev =
                                       drdb_write_led_inst_data_from_dev,
                                   .db_search_inst_in_reg = NULL,
                                   .db_free_inst_data_from_reg = free_led_data,
                                   .db_update_inst_in_reg = NULL,
                                   .db_exec = NULL };

void free_led_data(void *pdata) {
    if (pdata) {
        LedData *data = pdata;
        free_sdata(&data->state.prop);
        free_sdata(&data->dimmer.prop);
        free_sdata(&data->ontime.prop);
        free_sdata(&data->cumactivepower.prop);
        free_sdata(&data->colour.prop);
        free_sdata(&data->units.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

/* Used by List to copy nodes data */
void *copy_led_data(void *pdata) {
    LedData *data = pdata;
    LedData *ndata = NULL;
    if (data) {
        ndata = malloc(sizeof(LedData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(LedData));
            /* Try deep  copy for properties of pdata now */
            ndata->state.prop = reg_copy_pdata_prop(data->state.prop);
            ndata->dimmer.prop = reg_copy_pdata_prop(data->dimmer.prop);
            ndata->ontime.prop = reg_copy_pdata_prop(data->ontime.prop);
            ndata->cumactivepower.prop =
                reg_copy_pdata_prop(data->cumactivepower.prop);
            ndata->colour.prop = reg_copy_pdata_prop(data->colour.prop);
            ndata->units.prop = reg_copy_pdata_prop(data->units.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

static void populate_led_resource_id(void **pdata) {
    LedData *data = *pdata;
    data->state.resourceId = RES_M_ONOFF_VALUE;
    data->ontime.resourceId = RES_M_ONTIME_VALUE;
    data->dimmer.resourceId = RES_M_DIMMER_VALUE;
    data->cumactivepower.resourceId = RES_O_CUMM_ACTIVE_PWR_VALUE;
    data->power.resourceId = RES_O_CUMM_ACTIVE_PWR_VALUE;
    data->colour.resourceId = RES_M_COLOUR_VALUE;
    data->units.resourceId = RES_M_SENSOR_UNITS;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void drdb_add_led_inst_to_reg(Device *dev, uint8_t inst, uint8_t subdev) {
    int ret = 0;
    int pcount = 0;
    int pidx = 0;
    LedData *data = NULL;
    Property *prop = NULL;
    DRDBSchema *reg = malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        data = malloc(sizeof(LedData));
        if (data) {
            char str[32] = { '\0' };
            memset(data, '\0', sizeof(LedData));
            memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
            reg->instance = inst;
            if (subdev == 0) {
                pidx = RBRIGHTNESS;
                strcpy(str, "RED");
            } else if (subdev == 1) {
                pidx = GBRIGHTNESS;
                strcpy(str, "GREEN");
            } else if (subdev == 2) {
                pidx = BBRIGHTNESS;
                strcpy(str, "BLUE");
            } else {
                ret = -1;
                goto cleanup;
            }
            prop = db_read_dev_property(&dev->obj, &pcount);
            if (prop) {
                /* State */
                reg_data_add_property(pidx, prop, &data->state);
            } else {
                goto cleanup;
            }

            /* Color */
            memcpy(data->colour.value.stringval, str, sizeof(str));

            /* Application type */
            memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                   strlen(reg->obj.disc));

            data->cumactivepower.value.intval = 0;
            data->power.value.intval = 0;
            data->ontime.value.intval = 0;
        } else {
            ret = -1;
            log_debug("Err: LED:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &leddb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        reg->type = OBJ_TYPE_LED;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        populate_led_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug(
            "LED:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d is added to DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
cleanup:
    if (ret) {
        free_led_data(data);
        UKAMA_FREE(data);
    }
    UKAMA_FREE(prop);
    free_reg(&reg);
}

void drdb_add_led_dev_to_reg(void *pdev) {
    if (pdev) {
        Device *dev = pdev;
        for (uint8_t iter = 0; iter < 3; iter++) {
            uint8_t inst = list_size(reg_getdb(OBJ_TYPE_LED));
            drdb_add_led_inst_to_reg(dev, inst, iter);
        }
    }
}

int drdb_read_led_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    LedObjInfo *msgdata = NULL;
    if (reg->data) {
        LedData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            UKAMA_FREE(rqmsg->data);
        }
        msgdata = malloc(sizeof(LedObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(LedObjInfo));

        int propid = 0;

        /* State */
        ret = db_read_boolval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->state.resourceId),
                                   &msgdata->onoff, &rdata->state, &size);

        /* Dimmer */
        ret = db_read_intval_prop(&reg->obj,
                                  (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                      (rqmsg->resourceId ==
                                       rdata->dimmer.resourceId),
                                  &msgdata->dimmer, &rdata->dimmer, &size);

        /* Ontime */
        ret = db_read_intval_prop(&reg->obj,
                                  (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                      (rqmsg->resourceId ==
                                       rdata->ontime.resourceId),
                                  &msgdata->ontime, &rdata->ontime, &size);

        /* Cummulative power */
        ret = db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->ontime.resourceId),
            &msgdata->cumm_active_pwr, &rdata->cumactivepower, &size);

        /* Power */
        ret = db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->ontime.resourceId),
            &msgdata->pwr_factor, &rdata->power, &size);

        /* Colour */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->colour.resourceId),
                                   msgdata->colour, &rdata->colour, &size);

        /* sensor units */
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
            rqmsg->datasize = sizeof(LedObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(LedObjInfo);
        }

        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("LED:: Reading request for resource %d for Device %s Disc: %s"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    return ret;
}

int drdb_write_led_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = -1;
    int size = 0;
    LedData *rdata = reg->data;
    LedObjInfo *msgdata = NULL;
    if (reg->data && rqmsg->data) {
        LedObjInfo *msgdata = rqmsg->data;

        /* State */
        ret = db_write_boolval_prop(&reg->obj,
                                    (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                        (rqmsg->resourceId ==
                                         rdata->state.resourceId),
                                    msgdata->onoff, &rdata->state);

        if (ret) {
            log_trace(
                "Err(%d): GPIO:: Write filed for resource %d for Device %s Disc: %s"
                "Module ID %s is not writable",
                ret, rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                reg->obj.mod_UUID);
        } else {
            log_trace("GPIO:: Write done for resource %d for Device %s Disc: %s"
                      "Module ID %s is not writable",
                      rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                      reg->obj.mod_UUID);
        }
    }
    return ret;
}
