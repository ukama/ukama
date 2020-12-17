/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "registry/alarm.h"

#include "headers/ubsp/property.h"
#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "inc/alarmhandler.h"
#include "inc/dbhandler.h"
#include "inc/reghelper.h"
#include "inc/regdb.h"
#include "headers/objects/alarm.h"
#include "headers/utils/log.h"

const DBFxnTable alarmdb_fxn_tbl = {
    .db_add_dev_to_reg =
        NULL, /* This function is called implicitly by registry driver */
    .db_read_data_from_dev = drdb_read_alarm_inst_data_from_dev,
    .db_write_data_from_dev = drdb_write_alarm_inst_data_to_dev,
    .db_search_inst_in_reg = NULL,
    .db_free_inst_data_from_reg = free_alarm_data,
    .db_update_inst_in_reg = NULL,
    .db_exec = drdb_exec_alarm_inst_rsrc
};

static int64_t get_time_stamp() {
    int64_t tsec = 0;
    tsec = time(NULL);
    return tsec;
}

static void populate_alarm_resource_id(void **pdata) {
    AlarmData *data = *pdata;
    data->eventtype.resourceId = RES_M_AL_EVENTTYPE;
    data->realtime.resourceId = RES_M_AL_REALTIME;
    data->state.resourceId = RES_M_AL_STATE;
    data->lowthreshold.resourceId = RES_M_AL_LOW_THRESHOLD;
    data->highthreshold.resourceId = RES_M_AL_HIGH_THRESHOLD;
    data->crithreshold.resourceId = RES_M_AL_CRIT_THRESHOLD;
    data->eventcount.resourceId = RES_M_AL_EVT_COUNT;
    data->time.resourceId = RES_M_AL_RECRD_TIME;
    data->clear.resourceId = RES_M_AL_CLEAR;
    data->sobjid.resourceId = RES_M_AL_OBJ_ID;
    data->sinstid.resourceId = RES_M_AL_INST_ID;
    data->srsrcid.resourceId = RES_M_AL_RSRC_ID;
    data->sensorvalue.resourceId = RES_M_SENSOR_VALUE;
    data->sensorunits.resourceId = RES_M_SENSOR_UNITS;
    data->applicationtype.resourceId = RES_O_APPLICATION_TYPE;
}

void free_alarm_data(void *pdata) {
    if (pdata) {
        AlarmData *data = pdata;
        free_sdata(&data->eventtype.prop);
        free_sdata(&data->realtime.prop);
        free_sdata(&data->state.prop);
        free_sdata(&data->lowthreshold.prop);
        free_sdata(&data->highthreshold.prop);
        free_sdata(&data->crithreshold.prop);
        free_sdata(&data->eventcount.prop);
        free_sdata(&data->time.prop);
        free_sdata(&data->clear.prop);
        free_sdata(&data->sobjid.prop);
        free_sdata(&data->sinstid.prop);
        free_sdata(&data->srsrcid.prop);
        free_sdata(&data->sensorvalue.prop);
        free_sdata(&data->sensorunits.prop);
        free_sdata(&data->applicationtype.prop);
    }
}

void *copy_alarm_data(void *pdata) {
    AlarmData *data = pdata;
    AlarmData *ndata = NULL;
    if (data) {
        ndata = malloc(sizeof(AlarmData));
        if (ndata) {
            memcpy(ndata, pdata, sizeof(AlarmData));
            /* Try deep  copy for properties of pdata now */
            ndata->eventtype.prop = reg_copy_pdata_prop(data->eventtype.prop);
            ndata->realtime.prop = reg_copy_pdata_prop(data->realtime.prop);
            ndata->state.prop = reg_copy_pdata_prop(data->state.prop);
            ndata->lowthreshold.prop =
                reg_copy_pdata_prop(data->lowthreshold.prop);
            ndata->highthreshold.prop =
                reg_copy_pdata_prop(data->highthreshold.prop);
            ndata->crithreshold.prop =
                reg_copy_pdata_prop(data->crithreshold.prop);
            ndata->eventcount.prop = reg_copy_pdata_prop(data->eventcount.prop);
            ndata->time.prop = reg_copy_pdata_prop(data->time.prop);
            ndata->clear.prop = reg_copy_pdata_prop(data->clear.prop);
            ndata->sobjid.prop = reg_copy_pdata_prop(data->sobjid.prop);
            ndata->sinstid.prop = reg_copy_pdata_prop(data->sinstid.prop);
            ndata->srsrcid.prop = reg_copy_pdata_prop(data->srsrcid.prop);
            ndata->sensorvalue.prop =
                reg_copy_pdata_prop(data->sensorvalue.prop);
            ndata->sensorunits.prop =
                reg_copy_pdata_prop(data->sensorunits.prop);
            ndata->applicationtype.prop =
                reg_copy_pdata_prop(data->applicationtype.prop);
        }
    }
    return ndata;
}

int drdb_enable_alarm(void *data) {
    int ret = 0;
    //TODO:
    return ret;
}

int drdb_disable_alarm(void *data) {
    int ret = 0;
    //TODO:
    return ret;
}

int drdb_clear_alarm(void *data) {
    int ret = 0;
    //TODO:
    return ret;
}

void drdb_add_alarm_inst_to_reg(DevObj *obj, DevIDT *idt,
                                AlarmPropertyData *pdata) {
    int ret = 0;
    uint16_t instid = 0;
    AlarmData *data = NULL;
    DRDBSchema *reg = malloc(sizeof(DRDBSchema));
    if (reg) {
        memset(reg, '\0', sizeof(DRDBSchema));
        memcpy(&reg->obj, obj, sizeof(DevObj));

        data = malloc(sizeof(AlarmData));
        if (data) {
            memset(data, '\0', sizeof(AlarmData));

            /* Event Type */
            data->eventtype.value.intval = EVENT_TYPE_ALARM_CURR_STATE;

            /* Real Time */
            data->realtime.value.boolval = EVENT_REALTIME;

            /* URI identifiers for device */
            if (idt && pdata) {
                /* URI settings */
                data->sobjid.value.intval = idt->sobjid;
                data->sinstid.value.intval = idt->sinstid;
                data->srsrcid.value.intval = idt->srsrcid;

                /* Threshold settings */
                /* Critical threshold value */
                if (pdata->pcritthreshold) {
                    reg_data_copy_property(&data->crithreshold.prop,
                                           pdata->pcritthreshold);
                }

                /* High threshold value */
                if (pdata->phighthreshold) {
                    reg_data_copy_property(&data->highthreshold.prop,
                                           pdata->phighthreshold);
                }

                /* Low threshold value */
                if (pdata->plowthreshold) {
                    reg_data_copy_property(&data->lowthreshold.prop,
                                           pdata->plowthreshold);
                }

                /* Sensor Value */
                if (pdata->psensorvalue) {
                    reg_data_copy_property(&data->sensorvalue.prop,
                                           pdata->psensorvalue);

                    /* Units */
                    memcpy(data->sensorunits.value.stringval,
                           pdata->psensorvalue->units,
                           strlen(pdata->psensorvalue->units));

                    /* Application type */
                    memcpy(data->applicationtype.value.stringval,
                           pdata->psensorvalue->name,
                           strlen(pdata->psensorvalue->name));

                    /* Resource with execute permission. */

                    /* Clear Alarm */
                    data->clear.execFunc = &drdb_clear_alarm;

                    /* Enable Alarm */
                    data->enbalarm.execFunc = &drdb_enable_alarm;

                    /* Disable Alarm */
                    data->disalarm.execFunc = &drdb_disable_alarm;
                }
            } else {
                ret = -1;
                goto cleanup;
            }
        } else {
            ret = -1;
            log_debug("Err: ALARM:: Failed to allocate memory for Schema Data "
                      "%d/%d/%d Device Name: %s, Disc: %s, Type: %d, "
                      "Module Id %s.",
                      idt->sobjid, idt->sinstid, idt->srsrcid, reg->obj.name,
                      reg->obj.disc, reg->obj.type, reg->obj.mod_UUID);
            goto cleanup;
        }
        reg->dbfxntbl = &alarmdb_fxn_tbl;
        reg->instance = list_size(reg_getdb(OBJ_TYPE_ALARM));
        reg->data = data;
        reg->type = OBJ_TYPE_ALARM;
        strcpy(reg->UUID, obj->mod_UUID);
        populate_alarm_resource_id(&reg->data);
        reg_append_inst(reg_getdb(OBJ_TYPE_ALARM), reg);
        log_debug("DRDB:: Device Name: %s, Disc: %s, Type: %d, Module Id %s"
                  " Instance %d is added to DB.",
                  reg->obj.name, reg->obj.disc, reg->obj.type,
                  reg->obj.mod_UUID, reg->instance);
    }
cleanup:
    if (ret) {
        free_alarm_data(data);
        UKAMA_FREE(data);
    }
    free_reg(&reg);
}

int drdb_read_alarm_inst_data_from_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    int size = 0;
    AlarmObjInfo *msgdata = NULL;
    if (reg->data) {
        AlarmData *rdata = reg->data;
        /* TODO:  For now whatever be the case we are sending whole struct
    	 * Property requested will be updated and rest all will be zero. */
        /* free any memory if allocated and re-assign*/
        if (rqmsg->data) {
            UKAMA_FREE(rqmsg->data);
        }
        msgdata = malloc(sizeof(AlarmObjInfo));
        if (!msgdata) {
            return -1;
        }
        memset(msgdata, '\0', sizeof(AlarmObjInfo));

        /* Event Type*/
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->eventtype.resourceId),
            &msgdata->eventtype, &rdata->eventtype, &size);

        /* Real Time */
        ret |= db_read_boolval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->realtime.resourceId),
            &msgdata->realtime, &rdata->realtime, &size);

        /* Alarm State */
        ret |= db_read_shortintval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->state.resourceId),
            &msgdata->state, &rdata->state, &size);

        /* Description Type */
        ret |= db_read_strval_prop(&reg->obj,
                                   (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                       (rqmsg->resourceId ==
                                        rdata->disc.resourceId),
                                   msgdata->disc, &rdata->disc, &size);

        /* Low Threshold */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->lowthreshold.resourceId),
            &msgdata->lowthreshold, &rdata->lowthreshold, &size);

        /* High Threshold */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->highthreshold.resourceId),
            &msgdata->highthreshold, &rdata->highthreshold, &size);

        /* Critical Threshold */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->crithreshold.resourceId),
            &msgdata->crithreshold, &rdata->crithreshold, &size);

        /* Event Count */
        ret |= db_read_intval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->eventcount.resourceId),
            &msgdata->eventcount, &rdata->eventcount, &size);

        /* Time*/
        ret |= db_read_longintval_prop(&reg->obj,
                                       (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                                           (rqmsg->resourceId ==
                                            rdata->time.resourceId),
                                       &msgdata->time, &rdata->time, &size);

        /* Object Id */
        ret |= db_read_shortintval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->sobjid.resourceId),
            &msgdata->sobjid, &rdata->sobjid, &size);

        /* Instance Id*/
        ret |= db_read_shortintval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->sinstid.resourceId),
            &msgdata->sinstid, &rdata->sinstid, &size);

        /* Resource id */
        ret |= db_read_shortintval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->srsrcid.resourceId),
            &msgdata->srsrcid, &rdata->srsrcid, &size);

        /* Sensor Value */
        ret |= db_read_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->sensorvalue.resourceId),
            &msgdata->sensorvalue, &rdata->sensorvalue, &size);

        /* Units */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->sensorunits.resourceId),
            msgdata->sensorunits, &rdata->sensorunits, &size);

        /* Application Type */
        ret |= db_read_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->applicationtype.resourceId),
            msgdata->applicationtype, &rdata->applicationtype, &size);

        /* Instance Id */
        msgdata->instanceId = reg->instance;

        /* Check if was single prop read or full struct read */
        if (rqmsg->resourceId == ALL_RESOURCE_ID) {
            rqmsg->datasize = sizeof(AlarmObjInfo);
        } else {
            /* Because of above TODO */
            //rqmsg->datasize = size;
            rqmsg->datasize = sizeof(AlarmObjInfo);
        }
        /* Request message updated with requested data */
        rqmsg->data = msgdata;
    }
    log_trace("ALARM:: Reading request for resource %d for Device %s Disc: %s"
              "Module ID %s",
              rqmsg->resourceId, reg->obj.name, reg->obj.disc,
              reg->obj.mod_UUID);
    return ret;
}

int drdb_write_alarm_inst_data_to_dev(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = -1;
    AlarmData *rdata = reg->data;
    if (reg->data && rqmsg->data) {
        AlarmObjInfo *msgdata = rqmsg->data;

        /* Low Threshold */
        ret |= db_write_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->lowthreshold.resourceId),
            msgdata->lowthreshold, &rdata->lowthreshold);

        /* High Threshold */
        ret |= db_write_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->highthreshold.resourceId),
            msgdata->highthreshold, &rdata->highthreshold);

        /* Critical Threshold */
        ret |= db_write_doubleval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->crithreshold.resourceId),
            msgdata->crithreshold, &rdata->crithreshold);

        /* Application Type */
        ret |= db_write_strval_prop(
            &reg->obj,
            (rqmsg->resourceId == ALL_RESOURCE_ID) ||
                (rqmsg->resourceId == rdata->applicationtype.resourceId),
            msgdata->applicationtype, &rdata->applicationtype);

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

/* Compare Alarm instance based on the device object and dependent property*/
int drdb_compare_alarm_node(DRDBSchema *ndata, DevObj *obj, int pidx) {
    int ret = 0;
    /* Compare Object first */
    if (!(strcmp(obj->name, ndata->obj.name)) &&
        !(strcmp(obj->disc, ndata->obj.disc)) &&
        !(strcmp(obj->mod_UUID, ndata->obj.mod_UUID))) {
        /* Compare sensorvalue property id if that matches reported alarm*/
        if (ndata->data) {
            AlarmData *adata = ndata->data;
            PData *pdata = &adata->sensorvalue;
            if (pdata->prop) {
                if (pdata->prop->id == pidx) {
                    ret = 1;
                }
            }
        }
    }
    return ret;
}

/* Iterating through the  registry. */
DRDBSchema *drdb_search_alarm_inst(DevObj *obj, int pidx) {
    int ret = 0;
    DRDBSchema *nalarmreg = NULL;
    ListInfo *alarmdb = reg_getdb(OBJ_TYPE_ALARM);
    if (alarmdb) {
        if (alarmdb->logicalLength > 0) {
            ListNode *head = alarmdb->head;
            while (head) {
                if (head->data) {
                    DRDBSchema *ndata = head->data;
                    ret = drdb_compare_alarm_node(ndata, obj, pidx);
                    if (ret) {
                        nalarmreg = malloc(sizeof(DRDBSchema));
                        if (nalarmreg) {
                            memset(nalarmreg, '\0', sizeof(DRDBSchema));
                            memcpy(nalarmreg, ndata, sizeof(DRDBSchema));
                            if (ndata->data) {
                                /* Copy Alarm Data */
                                nalarmreg->data = copy_alarm_data(ndata->data);
                            }
                        }
                        break;
                    }
                }
                head = head->next;
            }
        }
    }
    return nalarmreg;
}

int drdb_exec_alarm_inst_rsrc(DRDBSchema *reg, MsgFrame *rqmsg) {
    int ret = 0;
    AlarmObjInfo *msgdata = rqmsg->data;
    if (reg->data) {
        AlarmData *rdata = reg->data;
        int propid = 0;
        if (rqmsg->resourceId == rdata->clear.resourceId) {
            if (rdata->clear.execFunc) {
                ret = rdata->clear.execFunc(reg);
            } else {
                ret = ERR_EDGEREG_RESR_NOTIMPL;
            }
        }
        if (rqmsg->resourceId == rdata->enbalarm.resourceId) {
            if (rdata->enbalarm.execFunc) {
                ret = rdata->enbalarm.execFunc(reg);
            } else {
                ret = ERR_EDGEREG_RESR_NOTIMPL;
            }
        }
        if (rqmsg->resourceId == rdata->disalarm.resourceId) {
            if (rdata->disalarm.execFunc) {
                ret = rdata->disalarm.execFunc(reg);
            } else {
                ret = ERR_EDGEREG_RESR_NOTIMPL;
            }
        }
        log_trace("ALARM:: Executed resource %d for Device %s Disc: %s"
                  "Module ID %s",
                  rqmsg->resourceId, reg->obj.name, reg->obj.disc,
                  reg->obj.mod_UUID);
    }
    return ret;
}

/* Update the alarm time and other required parameter for the node */
void drdb_alarm_update_data(DRDBSchema *node, Property *prop, int pidx,
                            void *value, uint8_t alertstate) {
    int ret = 0;
    double dummy = 0;
    int size = 0;
    if (node->data) {
        AlarmData *adata = node->data;

        /* Set time */
        adata->time.value.lintval = get_time_stamp();

        /* Set state*/
        adata->state.value.sintval = alertstate;

        /* Set description*/
        strcpy(adata->disc.value.stringval, prop[pidx].name);

        /* Set sensor value */
        db_set_prop_val(&adata->sensorvalue,
                        prop[prop[pidx].dep_prop->curr_idx].data_type, value);

        /* Low Threshold */
        ret |= db_read_doubleval_prop(&node->obj, true, &dummy,
                                      &adata->lowthreshold, &size);

        /* High Threshold */
        ret |= db_read_doubleval_prop(&node->obj, true, &dummy,
                                      &adata->highthreshold, &size);

        /* Critical Threshold */
        ret |= db_read_doubleval_prop(&node->obj, true, &dummy,
                                      &adata->crithreshold, &size);

        /* Event count */
        adata->eventcount.value.sintval++;

        /* Remove from list */
        //list_remove(reg_getdb(OBJ_TYPE_ALARM), node);

        /* Add to list */
        //list_prepend(reg_getdb(OBJ_TYPE_ALARM), node);

        /*Update List */
        reg_update_inst(reg_getdb(OBJ_TYPE_ALARM), node);
    }
}

/* This is called from the ISR context should be released as soon as possible to start monitor again.
 *  We could have a thread waiting in the app. Once we get this callback release a flag/semaphore and
 *  pass the required info for processing. Here I have just used ISR context to print info which is for demo only.
 */
void drdb_alarm_cb(DevObj *obj, AlertCallBackData **acbdata, int *count) {
    int pcount = 0;
    uint8_t alertstate = 0;
    if (*acbdata && *count) {
        Property *prop = db_read_dev_property(obj, &pcount);
        if (prop) {
            AlertCallBackData *adata = *acbdata;
            log_trace(
                "ALARM:: App Callback function for Device name %s Disc: %s Module UUID %s called "
                "from ubsp with %d alerts.",
                obj->name, obj->disc, obj->mod_UUID, *count);
            int pidx = adata->pidx;
            int didx = prop[pidx].dep_prop->curr_idx;
            alertstate = adata->alertstate;
            int size = get_sizeof(prop[didx].data_type);
            void *value = malloc(sizeof(size));
            if (value) {
                memcpy(value, adata->svalue, size);
            }

            log_trace(
                "ALARM:: Alert %d received for Property[%d], Name: %s , Value %lf %s.",
                alertstate, pidx, prop[pidx].name, *(double *)value,
                prop[didx].units);
            UKAMA_FREE(adata->svalue);
            UKAMA_FREE(adata);

            /* Search the alarm instance in registry  based on Object and
             *  alarm dependent property. For each sensor value we can have
             *  a alarm instance. For Temp we have alert only on current
             *  Temp value so one instance in registry where as in INA226
             *  we can alarm on current value or on power value and on voltage
             *  values so for INA we would have three alarm instances registered.  */
            DRDBSchema *regnode = drdb_search_alarm_inst(obj, didx);
            if (regnode) {
                /* Update the alarm Info to Alarm instance list */
                drdb_alarm_update_data(regnode, prop, pidx, value, alertstate);

                /* Push node to Alarmdb  */
                alarmdb_push(regnode);
                /* Free reg node */
                if (regnode) {
                    /* Free registery data */
                    if (regnode->data) {
                        free_reg_data(regnode);
                        regnode->data = NULL;
                    }
                    /* Free node */
                    free(regnode);
                    regnode = NULL;
                }
            } else {
                log_error(
                    "ALARM:: Alert received for Device %s Disc %s Module %s for property %d but not registered in alarm registry.",
                    obj->name, obj->disc, obj->mod_UUID, pidx);
            }
        } else {
            log_error(
                "ALARM:: Alert received for Device %s Disc %s Module %s but failed to read props.",
                obj->name, obj->disc, obj->mod_UUID);
        }
    } else {
        log_error("ALARM:: Alert received but callback data is corrupted.");
    }
}
