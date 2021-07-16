/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/module.h"

#include "dmt.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/module.h"
#include "headers/utils/log.h"

const DBFxnTable moddb_fxn_tbl = { .db_add_dev_to_reg = drdb_add_mod_to_reg,
                                   .db_read_data_from_dev =
                                       drdb_read_mod_inst_data_from_dev,
                                   .db_write_data_from_dev =
                                       drdb_write_mod_inst_data_from_dev,
                                   .db_search_inst_in_reg = NULL,
                                   .db_free_inst_data_from_reg = free_mod_data,
                                   .db_update_inst_in_reg = NULL,
                                   .db_exec = drdb_exec_mod_inst_rsrc };

void moddb_init(DeviceType type) {
    if (type == DEV_TYPE_NULL) {
        reg_init_type(type);
    }
}

void free_mod_data(void *pdata) {
    if (pdata) {
        ModuleData *data = pdata;
        free_sdata(&data->UUID.prop);
        free_sdata(&data->name.prop);
        free_sdata(&data->moduleclass.prop);
        free_sdata(&data->partnumber.prop);
        free_sdata(&data->mfgdate.prop);
        free_sdata(&data->mfgname.prop);
        free_sdata(&data->hwversion.prop);
        free_sdata(&data->mac.prop);
        free_sdata(&data->swversion.prop);
        free_sdata(&data->pswversion.prop);
        free_sdata(&data->devcount.prop);
    }
}

void *copy_module_data(void *pdata) {
    ModuleData *data = pdata;
    ModuleData *ndata = NULL;
    if (data) {
        ndata = dmt_malloc(sizeof(ModuleData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(ModuleData));
            /* Try deep  copy for properties of pdata now */
            ndata->UUID.prop = reg_copy_pdata_prop(data->UUID.prop);
            ndata->name.prop = reg_copy_pdata_prop(data->name.prop);
            ndata->moduleclass.prop =
                reg_copy_pdata_prop(data->moduleclass.prop);
            ndata->partnumber.prop = reg_copy_pdata_prop(data->partnumber.prop);
            ndata->mfgdate.prop = reg_copy_pdata_prop(data->mfgdate.prop);
            ndata->mfgname.prop = reg_copy_pdata_prop(data->mfgname.prop);
            ndata->hwversion.prop = reg_copy_pdata_prop(data->hwversion.prop);
            ndata->mac.prop = reg_copy_pdata_prop(data->mac.prop);
            ndata->swversion.prop = reg_copy_pdata_prop(data->swversion.prop);
            ndata->pswversion.prop = reg_copy_pdata_prop(data->pswversion.prop);
            ndata->devcount.prop = reg_copy_pdata_prop(data->devcount.prop);
        }
    }
    return ndata;
}

static void populate_mod_resource_id(void **pdata) {
    ModuleData *data = *pdata;
    data->UUID.resourceId = RES_M_MOD_UUID;
    data->name.resourceId = RES_M_MOD_NAME;
    data->moduleclass.resourceId = RES_M_MOD_CLASS;
    data->partnumber.resourceId = RES_M_PART_NUMBER;
    data->mfgdate.resourceId = RES_M_MOD_MFGDATE;
    data->mfgname.resourceId = RES_M_MOD_MFGNAME;
    data->mac.resourceId = RES_M_MOD_MAC;
    data->hwversion.resourceId = RES_M_MOD_HW_VERSION;
    data->swversion.resourceId = RES_M_MOD_SW_VERSION;
    data->pswversion.resourceId = RES_M_MOD_PSW_VERSION;
    data->devcount.resourceId = RES_M_MOD_DEV_COUNT;
}

int drdb_reset_mod_counters(void *data) {
    int ret = 0;

    return ret;
}

void drdb_add_mod_inst_to_reg(ModuleInfo *minfo, uint8_t instance,
                              uint8_t subdev) {
    ModuleData *data = NULL;
    DRDBSchema *reg = dmt_malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(reg->UUID, minfo->uuid, strlen(minfo->uuid));
        data = dmt_malloc(sizeof(ModuleData));
        if (data) {
            memset(data, '\0', sizeof(ModuleData));
            /* DB read module info */
            if (minfo) {
                /* UIID */
                memcpy(data->UUID.value.stringval, minfo->uuid,
                       strlen(minfo->uuid));

                /* Name */
                memcpy(data->name.value.stringval, minfo->name,
                       strlen(minfo->name));

                /* Manufacturing Name */
                memcpy(data->mfgname.value.stringval, minfo->mfg_name,
                       strlen(minfo->mfg_name));

                /* Manufacturing Date */
                memcpy(data->mfgdate.value.stringval, minfo->mfg_date,
                       strlen(minfo->mfg_date));

                /* Module class */
                data->moduleclass.value.intval = minfo->module;

                /* Part Number */
                memcpy(data->partnumber.value.stringval, minfo->partno,
                       strlen(minfo->partno));

                /* HW Version */
                memcpy(data->hwversion.value.stringval, minfo->hwver,
                       strlen(minfo->hwver));

                /* SW Version */
                //memcpy(data->swversion.value.stringval, minfo->pswver, strlen(minfo->pswver));
                db_versiontostr(minfo->swver, data->swversion.value.stringval);
                /* Production SW Version */
                //memcpy(data->pswversion.value.stringval, minfo->swver, strlen(minfo->swver));
                db_versiontostr(minfo->pswver,
                                data->pswversion.value.stringval);
                /* MAC */
                memcpy(data->mac.value.stringval, minfo->mac,
                       strlen(minfo->mac));

                /* Dev Count */
                data->devcount.value.intval = minfo->dev_count;

                /* Module UUID*/
                strcpy(reg->UUID, minfo->uuid);
            }
        } else {
            log_debug(
                "Err: MODULE:: Failed to allocate memory for Schema Data.");
            goto cleanup;
        }
        reg->dbfxntbl = &moddb_fxn_tbl;
        reg->instance = instance;
        reg->data = data;
        reg->type = OBJ_TYPE_MOD;
        populate_mod_resource_id(&reg->data);
        reg_append_inst(reg_getdb(reg->type), reg);
        log_debug("DRDB:: Module Id %s Instance %d is added to DB.",
                  minfo->uuid, reg->instance);
    }
cleanup:
    free_reg(&reg);
}

void drdb_add_mod_to_reg(void *pinfo) {
    if (pinfo) {
        uint8_t inst = list_size(reg_getdb(OBJ_TYPE_MOD));
        drdb_add_mod_inst_to_reg(pinfo, inst, 0);
    }
}

int drdb_read_mod_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    ModuleObjInfo *msgdata = NULL;
    if (reg->data) {
        ModuleData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
		 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            dmt_free(rqmsg->data);
        }
        msgdata = dmt_malloc(sizeof(ModuleObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(ModuleObjInfo));

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

        /* Module Class */
        ret |= db_read_intval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->moduleclass.resourceId),
                                   &msgdata->class, &rdata->moduleclass, &size);

        /* Device Count */
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->devcount.resourceId),
            &msgdata->device_count, &rdata->devcount, &size);

        /* Manufacturer Name */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->mfgname.resourceId),
                                   msgdata->mfgname, &rdata->mfgname, &size);

        /* Manufacturing Date */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->mfgdate.resourceId),
                                   msgdata->mfgdate, &rdata->mfgdate, &size);

        /* Part Number */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->partnumber.resourceId),
            msgdata->partnumber, &rdata->partnumber, &size);

        /* HW Version */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->hwversion.resourceId),
            msgdata->hw_version, &rdata->hwversion, &size);

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

        /* MAC */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->mac.resourceId),
                                   msgdata->mac, &rdata->mac, &size);

        /* Instance Id */
        msgdata->instanceId = reg->instance;

        /* Check if was single prop read or full struct read */
        if (rqmsg->resourceId == ALL_RESOURCE_ID) {
            rqmsg->datasize = sizeof(ModuleObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(ModuleObjInfo);
        }

        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("UNIT:: Reading request for resource %d for"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.mod_UUID);
    return ret;
}

int drdb_write_mod_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    log_trace("UNIT:: Reading request for resource %d for"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.mod_UUID);
    ret = ERR_EDGEREG_PROP_PERMDENIED;
    return ret;
}

int drdb_exec_mod_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    ModuleObjInfo *msgdata = rqmsg->data;
    if (reg->data) {
        ModuleData *rdata = reg->data;
        int propid = 0;
        /* No function added yet. */
        log_trace("UNIT:: Reading request for resource %d for"
                  "Module ID %s",
                  rqmsg->resourceId, reg->obj.mod_UUID);
    }
    return ret;
}
