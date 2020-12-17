/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "headers/errorcode.h"
#include "headers/ubsp/ubsp.h"
#include "inc/dbhandler.h"
#include "inc/regdb.h"
#include "inc/reghelper.h"
#include "inc/registry.h"
#include "headers/objects/objects.h"
#include "registry/adc.h"
#include "registry/alarm.h"
#include "registry/atten.h"
#include "registry/gpio.h"
#include "registry/led.h"
#include "registry/module.h"
#include "registry/pwr.h"
#include "registry/tmp.h"
#include "registry/unit.h"
#include "headers/utils/log.h"

ListInfo ukamadevdb[OBJ_TYPE_MAX] = { 0 };

ListInfo *reg_getdb(ObjectType type) {
    ListInfo *drdb = NULL;
    switch (type) {
    case OBJ_TYPE_UNIT:
    case OBJ_TYPE_MOD:
    case OBJ_TYPE_TMP:
    case OBJ_TYPE_VOLT:
    case OBJ_TYPE_CURR:
    case OBJ_TYPE_PWR:
    case OBJ_TYPE_DIP:
    case OBJ_TYPE_DOP:
    case OBJ_TYPE_LED:
    case OBJ_TYPE_ADC:
    case OBJ_TYPE_ATT:
    case OBJ_TYPE_ALARM:
        drdb = &ukamadevdb[type];
        break;
    default: {
        drdb = NULL;
    }
    }
    return drdb;
}

int reg_validate_dev_type(ObjectType type) {
    int ret = 0;
    if ((type <= OBJ_TYPE_NULL) && (type >= OBJ_TYPE_MAX)) {
        ret = MSG_RESP_INVALID_DEVTYPE;
    }
    return ret;
}

Property *reg_copy_pdata_prop(Property *sdata) {
    Property *prop = NULL;
    if (sdata) {
        prop = malloc(sizeof(Property));
        if (prop) {
            memcpy(prop, sdata, sizeof(Property));
        }
    }
    return prop;
}

void free_sdata(Property **prop) {
    if (*prop) {
        free(*prop);
        *prop = NULL;
    }
}

size_t size_of_schema_data(ObjectType type) {
    size_t size = 0;
    switch (type) {
    case DEV_TYPE_TMP:
        size = sizeof(TempData);
        break;
    case OBJ_TYPE_VOLT:
    case OBJ_TYPE_CURR:
    case OBJ_TYPE_PWR:
        size = sizeof(GenPwrData);
        break;
    case OBJ_TYPE_DIP:
    case OBJ_TYPE_DOP:
        size = sizeof(DigitalData);
        break;
    case OBJ_TYPE_LED:
        size = sizeof(LedData);
        break;
    case OBJ_TYPE_ADC:
        size = sizeof(AdcData);
        break;
    case OBJ_TYPE_ATT:
        size = sizeof(AttData);
        break;
    case OBJ_TYPE_ALARM:
        size = sizeof(AlarmData);
        break;
    default: {
    }
    }
    return size;
}

void *copy_schema_data(ObjectType type, void *srcdata) {
    switch (type) {
    case OBJ_TYPE_TMP:
        return copy_tmp_data(srcdata);
    case OBJ_TYPE_VOLT:
    case OBJ_TYPE_CURR:
    case OBJ_TYPE_PWR:
        return copy_pwr_data(srcdata);
    case OBJ_TYPE_DIP:
    case OBJ_TYPE_DOP:
        return copy_gpio_data(srcdata);
    case OBJ_TYPE_LED:
        return copy_led_data(srcdata);
    case OBJ_TYPE_ADC:
        return copy_adc_data(srcdata);
    case OBJ_TYPE_ATT:
        return copy_att_data(srcdata);
    case OBJ_TYPE_ALARM:
        return copy_alarm_data(srcdata);
    default: {
        return NULL;
    }
    }
}

void free_reg_data(DRDBSchema *reg) {
    switch (reg->type) {
    case OBJ_TYPE_TMP:
        free_tmp_data(reg->data);
        break;
    case OBJ_TYPE_VOLT:
    case OBJ_TYPE_CURR:
    case OBJ_TYPE_PWR:
        free_pwr_data(reg->data);
        break;
    case OBJ_TYPE_DIP:
    case OBJ_TYPE_DOP:
        free_gpio_data(reg->data);
        break;
    case OBJ_TYPE_LED:
        free_led_data(reg->data);
        break;
    case OBJ_TYPE_ADC:
        free_adc_data(reg->data);
        break;
    case OBJ_TYPE_ATT:
        free_att_data(reg->data);
        break;
    case OBJ_TYPE_ALARM:
        free_alarm_data(reg->data);
        break;
    default: {
    }
    }
    UKAMA_FREE(reg->data);
}

/* Remove the Data(DRDBSchema) from the Registry node */
void free_reg(DRDBSchema **reg) {
    if (*reg) {
        if ((*reg)->data) {
            /* Deletes the property for the DRDBSchema data.*/
            free_reg_data(*reg);
            (*reg)->data = NULL;
        }
        /*Delete DRDB SChema*/
        free(*reg);
        *reg = NULL;
    }
}

/* Removing node from the  Registry*/
static void remove_inst_from_reg(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            free_reg_data(node->data);
            free(node->data);
        }
        /* Delete Registry Node*/
        free(node);
    }
}

#if 0
/* Comparing device to the device in Device registry based on instance and device type.*/
static int cmp_dev_inst_in_reg(void *ip1, void *ip2) {
    DRDBSchema *inp1 = (DRDBSchema *)ip1;
    DRDBSchema *inp2 = (DRDBSchema *)ip2;
    int ret = 0;
    if ((inp1->instance == inp2->instance) &&
        (inp1->obj.type == inp2->obj.type)) {
        ret = 1;
    }
    return ret;
}

static int cmp_misc_inst_in_reg(void *ip1, void *ip2) {
    DRDBSchema *inp1 = (DRDBSchema *)ip1;
    DRDBSchema *inp2 = (DRDBSchema *)ip2;
    int ret = 0;
    if ((inp1->instance == inp2->instance) &&
        !(strcmp(inp1->UUID, inp2->UUID))) {
        ret = 1;
    }
    return ret;
}
#endif

static int cmp_inst_in_reg(void *ip1, void *ip2) {
    DRDBSchema *inp1 = (DRDBSchema *)ip1;
    DRDBSchema *inp2 = (DRDBSchema *)ip2;
    int ret = 0;
    if (inp1->instance == inp2->instance) {
        ret = 1;
    }
    return ret;
}

static void *cpy_inst_in_reg(void *srschema) {
    DRDBSchema *dschema = NULL;
    DRDBSchema *sschema = srschema;
    if (sschema) {
        dschema = malloc(sizeof(DRDBSchema));
        if (dschema) {
            memcpy(dschema, sschema, sizeof(DRDBSchema));
            dschema->data = NULL;
            if (sschema->data) {
                dschema->data = copy_schema_data(sschema->type, sschema->data);
            }
        }
    }
    return dschema;
}

/* Searching device in the Device registry*/
DRDBSchema *reg_search_inst(int instance, uint16_t misc, ObjectType type) {
    DRDBSchema *fdev = NULL;
    DRDBSchema *sdev = NULL;
    ListInfo *drdb = reg_getdb(type);
    if (!drdb) {
        goto cleanup;
    }
    sdev = malloc(sizeof(DRDBSchema));
    if (sdev) {
        memset(sdev, '\0', sizeof(DRDBSchema));
        sdev->instance = instance;
        sdev->type = type;
        fdev = list_search(drdb, sdev);
        if (fdev) {
            log_debug(
                "REGHELPER:: Device Name %s, Disc: %s Module UUID: %s Instance id %d found.",
                fdev->obj.name, fdev->obj.disc, fdev->obj.mod_UUID,
                fdev->instance);
        } else {
            log_debug(
                "REGHELPER:: Instance Id %d for device type 0x%x not found.",
                instance, type);
        }
    }
    if (sdev) {
        free(sdev);
        sdev = NULL;
    }
cleanup:
    return fdev;
}

int reg_update_dev(DRDBSchema *reg) {
    int ret = 0;
    if (reg) {
        ListInfo *drdb = reg_getdb(reg->type);
        if (!drdb) {
            return MSG_RESP_INVALID_DEVTYPE;
        }
        /*Return 0 for update success.*/
        //ret = list_update(drdb, reg);
        if (ret) {
            log_error(
                "REGHELPER:: Device Name %s, Disc: %s Module UUID: %s Instance %d update failed.",
                reg->obj.name, reg->obj.disc, reg->obj.mod_UUID, reg->instance);
        }
    }
    return ret;
}

void reg_exit_type(ObjectType type) {
    ListInfo *drdb = reg_getdb(type);
    if (drdb) {
        list_destroy(drdb);
    }
}

void reg_init_type(ObjectType type) {
    ListInfo *drdb = reg_getdb(type);
    if (drdb) {
        list_new(drdb, sizeof(DRDBSchema), remove_inst_from_reg,
                 cmp_inst_in_reg, cpy_inst_in_reg);
    }
}

void reg_init() {
    for (int iter = OBJ_TYPE_UNIT; iter < OBJ_TYPE_MAX; iter++) {
        list_new(reg_getdb(iter), sizeof(DRDBSchema), remove_inst_from_reg,
                 cmp_inst_in_reg, cpy_inst_in_reg);
    }
}

void reg_exit() {
    for (int iter = OBJ_TYPE_UNIT; iter < OBJ_TYPE_MAX; iter++) {
        log_debug("REGHELPER:: Cleaning registry for Object type 0x%2x.", iter);
        list_destroy(reg_getdb(iter));
    }
}

void reg_append_inst(ListInfo *regdb, void *node) {
    if (regdb && node) {
        list_append(regdb, node);
    }
}

void reg_prepend_inst(ListInfo *regdb, void *node) {
    if (regdb && node) {
        list_prepend(regdb, node);
    }
}

void reg_update_inst(ListInfo *regdb, void *node) {
    if (regdb && node) {
        list_update(regdb, node);
    }
}

int reg_read_max_instance(ObjectType type) {
    ListNode *node = NULL;
    int instance = 0;
    ListInfo *drdb = reg_getdb(type);
    if (drdb) {
        while (TRUE) {
            list_next(drdb, &node);
            if (node) {
                DRDBSchema *dev_reg = node->data;
                if (type == dev_reg->type) {
                    instance = MAX(dev_reg->instance, instance);
                }
            } else {
                break;
            }
        }
    }
    return instance;
}

DRDBSchema *reg_read_instance(int instance, DeviceType type) {
    int ret = 0;
    DRDBSchema *reg = NULL;
    reg = reg_search_inst(instance, 0, type);
    if (reg) {
        log_debug(
            "REGHELPER:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d found in DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    } else {
        UKAMA_FREE(reg->data);
        UKAMA_FREE(reg);
        log_debug(
            "REGHELPER:: Device Name: %s, Disc: %s, Type: %d, Module Id %s Instance %d not found in DB.",
            reg->obj.name, reg->obj.disc, reg->obj.type, reg->obj.mod_UUID,
            reg->instance);
    }
    return reg;
}

void *reg_data_value(PData *sd) {
    void *data = NULL;
    if (sd && sd->prop) {
        //Property* prop = sd->prop;
        switch (sd->prop->data_type) {
        case TYPE_NULL:
        case TYPE_BOOL:
        case TYPE_UINT8:
        case TYPE_INT8:
        case TYPE_UINT16:
        case TYPE_INT16:
        case TYPE_UINT32:
        case TYPE_INT32:
        case TYPE_ENUM: {
            data = &sd->value.intval;
            break;
        }
        case TYPE_FLOAT:
        case TYPE_DOUBLE: {
            data = &sd->value.doubleval;
            break;
        }
        case TYPE_STRING: {
            data = sd->value.stringval;
            break;
        }
        default: {
            data = &sd->value.intval;
        }
        }
    }
    return data;
}

int reg_add_dev(Device *dev) {
    int ret = 0;
    switch (dev->obj.type) {
    case DEV_TYPE_TMP:
        drdb_add_tmp_dev_to_reg(dev);
        break;
    case DEV_TYPE_PWR:
        drdb_add_pwr_dev_to_reg(dev);
        break;
    case DEV_TYPE_GPIO:
        drdb_add_gpio_dev_to_reg(dev);
        break;
    case DEV_TYPE_LED:
        drdb_add_led_dev_to_reg(dev);
        break;
    case DEV_TYPE_ADC:
        drdb_add_adc_dev_to_reg(dev);
        break;
    case DEV_TYPE_ATT:
        drdb_add_att_dev_to_reg(dev);
        break;
    default: {
        ret = -1;
    }
    }
    return ret;
}

int reg_add_misc(uint16_t misc, void *data) {
    int ret = 0;
    switch (misc) {
    case MISC_TYPE_UNIT:
        drdb_add_unit_to_reg(data);
        break;
    case MISC_TYPE_MODULE:
        drdb_add_mod_to_reg(data);
        break;
    default: {
        ret = -1;
    }
    }
    return ret;
}

int reg_read_dev(MsgFrame *req) {
    int ret = 0;
    /* Search the instance id */
    DRDBSchema *reg = NULL;
    reg = reg_search_inst(req->instance, req->misc, req->objecttype);
    if (reg) {
        const DBFxnTable *fxntbl = reg->dbfxntbl;
        if (fxntbl->db_read_data_from_dev) {
            ret = fxntbl->db_read_data_from_dev(reg, req);
        } else {
            ret = ERR_EDGEREG_INST_INAVLIDOP;
        }
    } else {
        ret = ERR_EDGEREG_INST_FXNTBL_MSNG;
    }
    return ret;
}

int reg_read_inst_count(MsgFrame *req) {
    int ret = ERR_EDGEREG_NOTAVAIL;
    int size = 0;
    ListInfo *db = reg_getdb(req->objecttype);
    if (db) {
        size = list_size(db);
        if (!(req->data)) {
            req->data = malloc(sizeof(int));
        }
        if (req->data) {
            memcpy(req->data, &size, sizeof(int));
            ret = 0;
        } else {
            ret = ERR_UBSP_MEMORY_EXHAUSTED;
        }
    }
    return ret;
}

int reg_write_dev(MsgFrame *req) {
    int ret = 0;
    /* Search the instance id */
    DRDBSchema *reg = NULL;
    reg = reg_search_inst(req->instance, req->misc, req->objecttype);
    if (reg) {
        if (reg->dbfxntbl) {
            const DBFxnTable *fxntbl = reg->dbfxntbl;
            if (fxntbl->db_write_data_from_dev) {
                ret = fxntbl->db_write_data_from_dev(reg, req);
            } else {
                ret = ERR_EDGEREG_INST_INAVLIDOP;
            }
        } else {
            ret = ERR_EDGEREG_INST_FXNTBL_MSNG;
        }
    }
    return ret;
}

int reg_exec_dev(MsgFrame *req) {
    int ret = 0;
    /* Search the instance id */
    DRDBSchema *reg = NULL;
    reg = reg_search_inst(req->instance, req->misc, req->objecttype);
    if (reg) {
        if (reg->dbfxntbl) {
            const DBFxnTable *fxntbl = reg->dbfxntbl;
            if (fxntbl->db_exec) {
                ret = fxntbl->db_exec(reg, req);
            } else {
                ret = ERR_EDGEREG_INST_INAVLIDOP;
            }
        } else {
            ret = ERR_EDGEREG_INST_FXNTBL_MSNG;
        }
    }
    return ret;
}

int reg_register_modules(char *puuid) {
    int ret = 0;
    ModuleInfo *minfo = NULL;
    minfo = db_read_module_info(puuid);
    if (minfo) {
        drdb_add_mod_to_reg(minfo);
        UKAMA_FREE(minfo);
    } else {
        ret = ERR_EDGEREG_REGCREFAILED;
    }
    return ret;
}

int reg_register_devices() {
    int ret = 0;
    for (DeviceType typeiter = DEV_TYPE_TMP; typeiter <= DEV_TYPE_ATT;
         typeiter++) {
        uint16_t count = 0;
        ret = ubsp_read_registered_dev_count(typeiter, &count);
        if (ret) {
            return ret;
        }
        log_debug("REGHELPER::Read %d registered devices for device type 0x%x.",
                  count, typeiter);
        if (count > 0) {
            Device *dev = malloc(sizeof(Device) * count);
            if (dev) {
                ret = ubsp_read_registered_dev(typeiter, dev);
                if (!ret) {
                    /* Add to list */
                    for (int iter = 0; iter < count; iter++) {
                        reg_add_dev(&dev[iter]);
                    }
                } else {
                }
                UKAMA_FREE(dev);
            } else {
                ret = ERR_UBSP_MEMORY_EXHAUSTED;
            }
        }
    }
    return ret;
}

int reg_register_misc() {
    int ret = 0;
    UnitInfo *uinfo = NULL;
    UnitCfg *ucfg = NULL;
    /* Read Unit Info UUID has to be "" for master*/
    uinfo = db_read_unit_info("");
    if (uinfo) {
        drdb_add_unit_to_reg(uinfo);
        /* Read Unit CFG */
        ucfg = db_read_unit_cfg("", uinfo->mod_count);
        if (ucfg) {
            for (int iter = 0; iter < uinfo->mod_count; iter++) {
                ret = reg_register_modules(ucfg[iter].mod_uuid);
                if (ret) {
                    ret = ERR_EDGEREG_REGCREFAILED;
                    goto cleanucfg;
                }
            }
        }

    } else {
        ret = ERR_EDGEREG_REGCREFAILED;
        goto cleanunit;
    }
cleanucfg:
    db_free_unit_cfg(&ucfg, uinfo->mod_count);
cleanunit:
    UKAMA_FREE(uinfo);
    return ret;
}

/* Assign property values for sensor data .*/
Property *assign_property(int propid, Property *prop) {
    Property *data = NULL;
    if (prop) {
        if (prop[propid].available == PROP_AVAIL) {
            data = malloc(sizeof(Property));
            if (data) {
                memcpy(data, &prop[propid], sizeof(Property));
            }
        }
    }
    return data;
}

/* Add propersty to registry schema from senor properties*/
void reg_data_add_property(int pidx, Property *prop, PData *pdata) {
    if (prop[pidx].available == PROP_AVAIL) {
        pdata->prop = malloc(sizeof(Property));
        if (pdata->prop) {
            memcpy(pdata->prop, &prop[pidx], sizeof(Property));
        }
    }
}

void reg_data_copy_property(Property **destp, Property *srcp) {
    if (srcp) {
        *destp = malloc(sizeof(Property));
        if (*destp) {
            memcpy(*destp, srcp, sizeof(Property));
        }
    }
}

void reg_free_alarm_prop(AlarmPropertyData **pdata) {
    AlarmPropertyData *data = *pdata;
    if (data) {
        if (data->pcritthreshold) {
            free(data->pcritthreshold);
            (data)->pcritthreshold = NULL;
        }
        if (data->phighthreshold) {
            free(data->phighthreshold);
            data->phighthreshold = NULL;
        }
        if (data->plowthreshold) {
            free(data->plowthreshold);
            data->plowthreshold = NULL;
        }
        if (data->plowlimitalarm) {
            free(data->plowlimitalarm);
            data->plowlimitalarm = NULL;
        }
        if (data->phighlimitalarm) {
            free(data->phighlimitalarm);
            data->phighlimitalarm = NULL;
        }
        if (data->pcrilimitalarm) {
            free(data->pcrilimitalarm);
            data->pcrilimitalarm = NULL;
        }
        if (data->psensorvalue) {
            free(data->psensorvalue);
            data->psensorvalue = NULL;
        }
        free(data);
        data = NULL;
    }
}

void *reg_initialize_alarm_prop(Property *prop, const AlarmSensorData *sdata) {
    AlarmPropertyData *data = malloc(sizeof(AlarmPropertyData));
    if (data) {
        memset(data, '\0', sizeof(AlarmPropertyData));

        if (sdata->pcrithresholdidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->pcritthreshold =
                assign_property(sdata->pcrithresholdidx, prop);
        }

        if (sdata->phighthresholdidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->phighthreshold =
                assign_property(sdata->phighthresholdidx, prop);
        }

        if (sdata->plowthresholdidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->plowthreshold =
                assign_property(sdata->plowthresholdidx, prop);
        }

        if (sdata->pcrilimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->pcrilimitalarm =
                assign_property(sdata->pcrilimitalarmidx, prop);
        }

        if (sdata->phighlimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->phighlimitalarm =
                assign_property(sdata->phighlimitalarmidx, prop);
        }

        if (sdata->plowlimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->plowlimitalarm =
                assign_property(sdata->plowlimitalarmidx, prop);
        }

        if (sdata->psensorvalueidx != PROPERY_IDX_NOT_APPLICABLE) {
            data->psensorvalue = assign_property(sdata->psensorvalueidx, prop);
        }
    }
    return data;
}

void reg_initialize_dev_idt(void *pidt, uint16_t inst, uint16_t oid,
                            uint16_t rid) {
    DevIDT *idt = pidt;
    if (idt) {
        idt->sinstid = inst;
        idt->sobjid = oid;
        idt->srsrcid = rid;
    }
}

void reg_register_sensor_alarms(DevObj *obj, Property *prop,
                                const AlarmSensorData *sdata, uint16_t inst,
                                uint16_t objid, uint16_t rsrcid) {
    DevIDT devidt = { 0 };
    AlarmPropertyData *pdata = NULL;
    reg_initialize_dev_idt(&devidt, inst, objid, rsrcid);
    pdata = reg_initialize_alarm_prop(prop, sdata);
    drdb_add_alarm_inst_to_reg(obj, &devidt, pdata);
    if (pdata) {
        reg_free_alarm_prop(&pdata);
    }
}

int reg_check_if_alarm_property_exist(Property *prop, int pcount) {
    int ret = 0;
    for (int iter = 0; iter < pcount; iter++) {
        if (prop[iter].prop_type == PROP_TYPE_ALERT) {
            /* Alert property found for device */
            ret = 1;
            break;
        }
    }
    return ret;
}

int reg_enable_alarms(DevObj *obj, const AlarmSensorData *sdata) {
    int ret = 0;
    if (sdata->pcrilimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
        ret = db_enable_alarm(obj, sdata->pcrilimitalarmidx);
    }

    if (sdata->phighlimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
        ret = db_enable_alarm(obj, sdata->phighlimitalarmidx);
    }

    if (sdata->plowlimitalarmidx != PROPERY_IDX_NOT_APPLICABLE) {
        ret = db_enable_alarm(obj, sdata->plowlimitalarmidx);
    }
    return ret;
}

int reg_print_node(void *pnode) {
    int ret = 0;
    if (pnode) {
        DRDBSchema *node = pnode;
        log_trace("REG:: Type 0x%2x, Instance %2d, UUID %24s", node->type,
                  node->instance, node->UUID);
        /* For sensors. */
        if (node->type > OBJ_TYPE_MOD) {
            log_trace("REG:: Name %24s, Disc %24s", node->obj.name,
                      node->obj.disc);
        }
        ret = 1;
    }
    return ret;
}

void reg_list_reg_devices() {
    for (ObjectType typeiter = OBJ_TYPE_UNIT; typeiter < OBJ_TYPE_MAX;
         typeiter++) {
        uint16_t count = list_size(reg_getdb(typeiter));
        list_for_each(reg_getdb(typeiter), &reg_print_node);
    }
}
