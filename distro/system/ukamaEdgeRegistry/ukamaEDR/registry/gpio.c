/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/gpio.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/digital_input.h"
#include "headers/objects/digital_output.h"
#include "headers/utils/log.h"

const DBFxnTable gpiodb_fxn_tbl = {
    .db_add_dev_to_reg = drdb_add_gpio_dev_to_reg,
    .db_read_data_from_dev = drdb_read_gpio_inst_data_from_dev,
    .db_write_data_from_dev = drdb_write_gpio_inst_data_from_dev,
    .db_search_inst_in_reg = NULL,
    .db_free_inst_data_from_reg = free_gpio_data,
    .db_update_inst_in_reg = NULL,
    .db_exec = NULL
};

void free_gpio_data(void *pdata) {
    if (pdata) {
        DigitalData *data = pdata;
        free_sdata(&data->direction.prop);
        free_sdata(&data->state.prop);
        free_sdata(&data->counter.prop);
        free_sdata(&data->polarity.prop);
        free_sdata(&data->ontime.prop);
        free_sdata(&data->offtime.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

/* Used by List to copy nodes data */
void *copy_gpio_data(void *pdata) {
    DigitalData *data = pdata;
    DigitalData *ndata = NULL;
    if (data) {
        ndata = malloc(sizeof(DigitalData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(DigitalData));
            /* Try deep  copy for properties of pdata now */
            ndata->direction.prop = reg_copy_pdata_prop(data->direction.prop);
            ndata->state.prop = reg_copy_pdata_prop(data->state.prop);
            ndata->counter.prop = reg_copy_pdata_prop(data->counter.prop);
            ndata->polarity.prop = reg_copy_pdata_prop(data->polarity.prop);
            ndata->ontime.prop = reg_copy_pdata_prop(data->ontime.prop);
            ndata->offtime.prop = reg_copy_pdata_prop(data->offtime.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

static void populate_digital_input_resource_id(void **pdata) {
    DigitalData *data = *pdata;
    data->state.resourceId = RES_M_DIGITAL_INPUT_STATE;
    data->counter.resourceId = RES_O_DIGITAL_INPUT_COUNTER;
    data->polarity.resourceId = RES_O_DIGITAL_INPUT_POLARITY;
    data->debounce.resourceId = RES_O_DIGITAL_INPUT_DEBOUNCE;
    data->edge.resourceId = RES_O_DIGITIAL_INPUT_EDGE_SELECTION;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
    data->sensortype.resourceId = RES_O_SENSOR_TYPE;
}

static void populate_digital_output_resource_id(void **pdata) {
    DigitalData *data = *pdata;
    data->state.resourceId = RES_M_DIGITAL_OUTPUT_STATE;
    data->polarity.resourceId = RES_O_DIGITAL_OUTPUT_POLARITY;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void drdb_add_gpio_inst_to_reg(Device *dev, uint8_t inst, uint8_t subdev) {
    int ret = 0;
    int pcount = 0;
    int pidx = 0;
    DigitalData *data = NULL;
    Property *prop = NULL;
    DRDBSchema *reg = malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, &dev->obj, sizeof(DevObj));
        data = malloc(sizeof(DigitalData));
        if (data) {
            memset(data, '\0', sizeof(DigitalData));

            /* Read Property table of the sensor. */
            prop = db_read_dev_property(&dev->obj, &pcount);
            if (prop) {
                reg_data_add_property(VALUE, prop, &data->state);
                reg_data_add_property(POLARITY, prop, &data->polarity);
                reg_data_add_property(DIRECTION, prop, &data->direction);
            } else {
                ret = -1;
                goto cleanup;
            }
            /* Application type */
            memcpy(data->applicationtype.value.stringval, reg->obj.disc,
                   strlen(reg->obj.disc));
            data->counter.value.intval = 0;
            data->ontime.value.intval = 0;
            data->offtime.value.intval = 0;
        } else {
            ret = -1;
            log_debug("Err: GPIO:: Failed to allocate memory for Schema Data"
                      "Device Name: %s, Disc: %s, Type: %d, Module Id %s.",
                      reg->obj.name, reg->obj.disc, reg->obj.type,
                      reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &gpiodb_fxn_tbl;
        reg->instance = inst;
        reg->data = data;
        strcpy(reg->UUID, dev->obj.mod_UUID);
        if (!db_read_inst_data_from_dev(&reg->obj, &data->direction)) {
            if (!strcmp(data->direction.value.stringval, GPIO_TYPE_INPUT)) {
                reg->type = OBJ_TYPE_DIP;
                reg->instance = list_size(reg_getdb(reg->type));
                populate_digital_input_resource_id(&reg->data);
            } else {
                reg->type = OBJ_TYPE_DOP;
                reg->instance = list_size(reg_getdb(reg->type));
                populate_digital_output_resource_id(&reg->data);
            }
        } else {
            goto cleanup;
        }
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug(
            "DRDB:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d is added to DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
cleanup:
    if (ret) {
        free_gpio_data(data);
        UKAMA_FREE(data);
    }
    UKAMA_FREE(prop);
    free_reg(&reg);
}

void drdb_add_gpio_dev_to_reg(void *pdev) {
    if (pdev) {
        Device *dev = pdev;
        uint8_t inst = 0; /* Instance value is postponed till property read.*/
        drdb_add_gpio_inst_to_reg(dev, inst, 0);
    }
}

int drdb_read_gpio_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    DigObjInfo *msgdata = NULL;
    if (reg->data) {
        DigitalData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
			 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            UKAMA_FREE(rqmsg->data);
        }
        msgdata = malloc(sizeof(DigObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(DigObjInfo));

        int propid = 0;

        /* Direction */
        ret = db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->direction.resourceId),
            &msgdata->direction, &rdata->direction, &size);

        /* State */
        ret |= db_read_boolval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->state.resourceId),
            &msgdata->digital_state, &rdata->state, &size);

        /* Counter */
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->counter.resourceId),
            &msgdata->digital_counter, &rdata->counter, &size);

        /* Polarity */
        ret |= db_read_boolval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->polarity.resourceId),
            &msgdata->digital_polarity, &rdata->polarity, &size);

        /* Debounce */
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->debounce.resourceId),
            &msgdata->digital_debounce, &rdata->debounce, &size);

        /* Edge */
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->edge.resourceId),
            &msgdata->digitial_edge_selection, &rdata->edge, &size);

        /* Ontime  */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->ontime.resourceId),
                                   &msgdata->ontime, &rdata->ontime, &size);

        /* Offtime */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->offtime.resourceId),
                                   &msgdata->offtime, &rdata->offtime, &size);

        /* Application Type */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->applicationtype.resourceId),
            msgdata->application_type, &rdata->applicationtype, &size);

        /* Sensor Type */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->sensortype.resourceId),
            msgdata->sensor_type, &rdata->sensortype, &size);

        /* Instance Id */
        msgdata->instanceId = reg->instance;

        /* Check if was single prop read or full struct read */
        if (rqmsg->resourceId == ALL_RESOURCE_ID) {
            rqmsg->datasize = sizeof(DigObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(DigObjInfo);
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

int drdb_write_gpio_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = -1;
    int size = 0;
    DigitalData *rdata = reg->data;
    if (reg->data && rqmsg->data) {
        DigObjInfo *msgdata = rqmsg->data;

        /* State */
        ret = db_write_boolval_prop(&reg->obj,
                                    (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                        (rqmsg->resourceId ==
                                         rdata->state.resourceId),
                                    msgdata->digital_state, &rdata->state);

        /* Polarity */
        ret |= db_write_boolval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->polarity.resourceId),
            msgdata->digital_polarity, &rdata->polarity);

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
