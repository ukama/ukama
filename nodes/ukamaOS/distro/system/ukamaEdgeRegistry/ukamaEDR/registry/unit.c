/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/unit.h"

#include "dmt.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/unit.h"
#include "headers/utils/log.h"

const DBFxnTable unitdb_fxn_tbl = {
    .db_add_dev_to_reg = drdb_add_unit_to_reg,
    .db_read_data_from_dev = drdb_read_unit_inst_data_from_dev,
    .db_write_data_from_dev = drdb_write_unit_inst_data_from_dev,
    .db_search_inst_in_reg = NULL,
    .db_free_inst_data_from_reg = free_unit_data,
    .db_update_inst_in_reg = NULL,
    .db_exec = drdb_exec_unit_inst_rsrc
};

void unitdb_init(DeviceType type) {
    if (type == DEV_TYPE_NULL) {
        reg_init_type(type);
    }
}

void free_unit_data(void *pdata) {
    if (pdata) {
        UnitData *data = pdata;
        free_sdata(&data->UUID.prop);
        free_sdata(&data->name.prop);
        free_sdata(&data->unit.prop);
        free_sdata(&data->asmdate.prop);
        free_sdata(&data->oemname.prop);
        free_sdata(&data->skew.prop);
        free_sdata(&data->mac.prop);
        free_sdata(&data->swversion.prop);
        free_sdata(&data->pswversion.prop);
        free_sdata(&data->modcount.prop);
    }
}

void *copy_unit_data(void *pdata) {
    UnitData *data = pdata;
    UnitData *ndata = NULL;
    if (data) {
        ndata = dmt_malloc(sizeof(UnitData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(UnitData));
            /* Try deep  copy for properties of pdata now */
            ndata->UUID.prop = reg_copy_pdata_prop(data->UUID.prop);
            ndata->name.prop = reg_copy_pdata_prop(data->name.prop);
            ndata->unit.prop = reg_copy_pdata_prop(data->unit.prop);
            ndata->asmdate.prop = reg_copy_pdata_prop(data->asmdate.prop);
            ndata->oemname.prop = reg_copy_pdata_prop(data->oemname.prop);
            ndata->skew.prop = reg_copy_pdata_prop(data->skew.prop);
            ndata->mac.prop = reg_copy_pdata_prop(data->mac.prop);
            ndata->swversion.prop = reg_copy_pdata_prop(data->swversion.prop);
            ndata->pswversion.prop = reg_copy_pdata_prop(data->pswversion.prop);
            ndata->modcount.prop = reg_copy_pdata_prop(data->modcount.prop);
        }
    }
    return ndata;
}

int drdb_reset_unit_counters(void *data) {
    int ret = 0;

    return ret;
}

static void populate_unit_resource_id(void **pdata) {
    UnitData *data = *pdata;
    data->UUID.resourceId = RES_M_UNIT_UUID;
    data->name.resourceId = RES_M_UNIT_NAME;
    data->unit.resourceId = RES_M_UNIT_CLASS;
    data->skew.resourceId = RES_M_SKEW;
    data->asmdate.resourceId = RES_M_UNIT_ASMDATE;
    data->oemname.resourceId = RES_M_UNIT_OEMNAME;
    data->mac.resourceId = RES_M_UNIT_MAC;
    data->swversion.resourceId = RES_M_UNIT_SW_VERSION;
    data->pswversion.resourceId = RES_M_UNIT_PSW_VERSION;
    data->modcount.resourceId = RES_M_UNIT_MOD_COUNT;
}

void drdb_add_unit_inst_to_reg(UnitInfo *uinfo, uint8_t instance,
                               uint8_t subdev) {
    UnitData *data = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(reg->UUID, uinfo->uuid, strlen(uinfo->uuid));
        data = dmt_malloc(sizeof(UnitData));
        if (data) {
            memset(data, '\0', sizeof(UnitData));

            /* DB read module info */
            if (uinfo) {
                /* UIID */
                memcpy(data->UUID.value.stringval, uinfo->uuid,
                       strlen(uinfo->uuid));

                /* Name */
                memcpy(data->name.value.stringval, uinfo->name,
                       strlen(uinfo->name));

                /* OEM Name */
                memcpy(data->oemname.value.stringval, uinfo->oem_name,
                       strlen(uinfo->oem_name));

                /* Assembly Date */
                memcpy(data->asmdate.value.stringval, uinfo->assm_date,
                       strlen(uinfo->assm_date));

                /* Unit class */
                data->unit.value.intval = uinfo->unit;

                /* Skew  */
                memcpy(data->skew.value.stringval, uinfo->skew,
                       strlen(uinfo->skew));

                /* MAC */
                memcpy(data->mac.value.stringval, uinfo->mac,
                       strlen(uinfo->mac));

                /* SW Version */
                //memcpy(data->swversion.value.stringval, uinfo->pswver, strlen(uinfo->pswver));
                db_versiontostr(uinfo->swver, data->swversion.value.stringval);

                /* Production SW Version */
                //memcpy(data->pswversion.value.stringval, uinfo->swver, strlen(uinfo->swver));
                db_versiontostr(uinfo->pswver,
                                data->pswversion.value.stringval);
                /* MAC */
                memcpy(data->mac.value.stringval, uinfo->mac,
                       strlen(uinfo->mac));

                /* Module Count */
                data->modcount.value.intval = uinfo->mod_count;

                /* Registry UUID would be unit Id*/
                strcpy(reg->UUID, uinfo->uuid);
            }
        } else {
            log_debug("Err: UNIT:: Failed to allocate memory for Schema Data.");
            goto cleanup;
        }
        reg->dbfxntbl = &unitdb_fxn_tbl;
        reg->instance = instance;
        reg->data = data;
        reg->type = OBJ_TYPE_UNIT;
        populate_unit_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug("UNIT:: Module Id %s Instance %d is added to DB.",
                  uinfo->uuid, reg->instance);
    }
cleanup:
    free_reg(&reg);
}

void drdb_add_unit_to_reg(void *pinfo) {
    if (pinfo) {
        /* This should always have one instance only */
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_UNIT));
        if (inst > 0) {
            reg_exit_type(OBJ_TYPE_UNIT);
        }
        drdb_add_unit_inst_to_reg(pinfo, 0, 0);
    }
}

int drdb_read_unit_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    UnitObjInfo *msgdata = NULL;
    if (reg->data) {
        UnitData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            dmt_free(rqmsg->data);
        }
        msgdata = dmt_malloc(sizeof(UnitObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(UnitObjInfo));

        int propid = 0;
        /* UUID */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->UUID.resourceId),
                                   msgdata->uuid, &rdata->UUID, &size);

        /* Name */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->name.resourceId),
                                   msgdata->name, &rdata->name, &size);

        /* Unit Type*/
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->unit.resourceId),
                                   &msgdata->class, &rdata->unit, &size);

        /* Module Count */
        ret = db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->UUID.resourceId),
            &msgdata->module_count, &rdata->modcount, &size);

        /* OEM Name */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->oemname.resourceId),
                                   msgdata->oemname, &rdata->oemname, &size);

        /* Assembly Date */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->asmdate.resourceId),
                                   msgdata->asmdate, &rdata->asmdate, &size);

        /* Skew */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->skew.resourceId),
                                   msgdata->skew, &rdata->skew, &size);

        /* MAC */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->mac.resourceId),
                                   msgdata->mac, &rdata->mac, &size);

        /* SW Version */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->swversion.resourceId),
            msgdata->sw_version, &rdata->swversion, &size);

        /* PSW Version */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->pswversion.resourceId),
            msgdata->psw_version, &rdata->pswversion, &size);

        /* Instance Id */
        msgdata->instanceId = reg->instance;

        /* Check if was single prop read or full struct read */
        if (rqmsg->resourceId == ALL_RESOURCE_ID) {
            rqmsg->datasize = sizeof(UnitObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(UnitObjInfo);
        }

        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("UNIT:: Reading request for resource %d for"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.mod_UUID);
    return ret;
}

int drdb_write_unit_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    log_trace("UNIT:: Reading request for resource %d for"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.mod_UUID);
    ret = ERR_EDGEREG_PROP_PERMDENIED;
    return ret;
}

int drdb_exec_unit_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    UnitObjInfo *msgdata = rqmsg->data;
    if (reg->data) {
        UnitData *rdata = reg->data;
        int propid = 0;
        /* No function added yet. */
        log_trace("UNIT:: Reading request for resource %d for"
                  "Module ID %s",
                  rqmsg->resourceId, reg->obj.mod_UUID);
    }
    return ret;
}
