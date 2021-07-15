/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "inc/alarmhandler.h"
#include "dmt.h"

#include "headers/errorcode.h"
#include "headers/globalheader.h"
#include "ifmsg.h"
#include "msghandler.h"
#include "headers/utils/file.h"
#include "headers/utils/list.h"
#include "headers/utils/log.h"
#include "utils/timer.h"

#define RALARMFILE "/history/alarmdata.csv"

/* Periodic timer alarm reporting */
#define ALARM_PERIODIC_TIMER 15000

/* AlarmDB active and clear alarms which are not yet acknowledged by lwm2m client */
ListInfo alarmdb;

static pthread_t galarmthreadid = 0;
size_t galarmtimer;

/* Comparing Alarm node to the Alarms in AlarmDB based on URI and Device Object.*/
static int cmp_alarm(void *ip1, void *ip2) {
    int ret = 0;
    if (ip1 && ip2) {
        AlarmSchema *inp1 = (AlarmSchema *)ip1;
        AlarmSchema *inp2 = (AlarmSchema *)ip2;
        if (COMPARE_DEV_OBJ((*inp1).obj, (*inp2).obj)) {
            /* Compare uri and alarm state */
            if ((inp1->objid == inp2->objid) &&
                (inp1->instid == inp2->instid) &&
                (inp1->rsrcid == inp2->rsrcid)) {
                ret = 1;
            }
        }
    }
    return ret;
}

static void rmv_alarm(void *ip) {
    ListNode *node = (ListNode *)ip;
    if (node) {
        if (node->data) {
            dmt_free(node->data);
        }
        dmt_free(node);
    }
}

static void *cpy_alarm(void *ndata) {
    AlarmSchema *data = NULL;
    if (ndata) {
        data = dmt_malloc(sizeof(AlarmSchema));
        if (data) {
            memcpy(data, ndata, sizeof(AlarmSchema));
        }
    }
    return data;
}

void alarmdb_init() {
    list_new(&alarmdb, sizeof(AlarmSchema), rmv_alarm, cmp_alarm, cpy_alarm);
}

void alarmdb_exit() {
    log_debug("ALARMHANDLER: Clearing alarmdb.");
    list_destroy(&alarmdb);
}

void alarmdb_prepend(AlarmSchema *node) {
    list_prepend(&alarmdb, node);
}

int alarmdb_remove(AlarmSchema *node) {
    return list_remove(&alarmdb, node);
}

/* Updates the node in the Alarmdb */
int alarmdb_update(AlarmSchema *aschema) {
    return list_update(&alarmdb, aschema);
}

int alarmdb_node_exist(AlarmSchema *aschema) {
    return list_if_element_found(&alarmdb, aschema);
}

/* Search node based on object id, instance id, resource id and alert state */
AlarmSchema *alarmdb_search_node(AlarmSchema *snode) {
    AlarmSchema *fnode = NULL;
    if (snode) {
        fnode = list_search(&alarmdb, snode);
        if (fnode) {
            log_debug(
                "ALARMHANDLER:: Alarm instance %d for Device %d/%d/%d Name %s, Disc: %s Module UUID: %s found.",
                fnode->data.instanceId, fnode->objid, fnode->instid,
                fnode->rsrcid, fnode->obj.name, fnode->obj.disc,
                fnode->obj.mod_UUID);
        } else {
            log_error(
                "ALARMHANDLER:: Alarm instance %d for Device /%d/%d/%d  Name %s, Disc: %s Module UUID: %s not found.",
                snode->data.instanceId, snode->objid, snode->instid,
                snode->rsrcid, snode->obj.name, snode->obj.disc,
                snode->obj.mod_UUID);
        }
    }
    return fnode;
}

/* Create a node for the AlarmDb based on object id, instance id and resource id */
AlarmSchema *alarmdb_create_node(DevObj *obj, uint16_t objid, uint16_t instid,
                                 uint16_t rsrcid) {
    AlarmSchema *node = NULL;
    node = dmt_malloc(sizeof(AlarmSchema));
    if (node) {
        memset(node, '\0', sizeof(AlarmSchema));
        node->instid = instid;
        node->objid = objid;
        node->rsrcid = rsrcid;
        memcpy(&node->obj, obj, sizeof(DevObj));
    }
    return node;
}

/* Add event to a file for persistent storage */
int alramdb_record_data(AlarmSchema *node) {
    int ret = 0;
    char rowdesc[512] = { '\0' };
    char rowdata[1024] = { '\0' };
    if (node) {
        /* Row description*/
        sprintf(rowdesc,
                "%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,"
                "%s,%s,%s,%s,%s,%s,%s,\n",
                "RecordTime", "Token", "Reported", "TxCount", "Acknowledged",
                "Device Name", "Device Description", "Module", "Device Type",
                "Alarm Time", "Alarm ObjectId", "Alarm Instance Id",
                "Alarm resource Id", "Event Type", "Real Time", "State",
                "Alarm Description", "Low threshold", "High Threshold",
                "Critical Threshold", "Event Count", "Sensor Object Id",
                "Sensor InstanceId", "Sensor Resource ID", "Sensor Value",
                "Sensor Units", "Application Type");

        /* Alert Data */
        sprintf(rowdata,
                "%ld,%d,%d,%d,%d,%s,%s,%s,%d,%ld,%d,%d,%d,%d,%d,%d,%s,%lf,%lf,"
                "%lf,%d,%d,%d,%d,%lf,%s,%s,\n",
                time(NULL), node->token, node->reported, node->txcount,
                node->response, node->obj.name, node->obj.disc,
                node->obj.mod_UUID, node->obj.type, node->data.time,
                node->objid, node->instid, node->rsrcid, node->data.eventtype,
                node->data.realtime, node->data.state, node->data.disc,
                node->data.lowthreshold, node->data.highthreshold,
                node->data.crithreshold, node->data.eventcount,
                node->data.sobjid, node->data.sinstid, node->data.srsrcid,
                node->data.sensorvalue, node->data.sensorunits,
                node->data.applicationtype);

        ret = file_add_record(RALARMFILE, rowdesc, rowdata);
    }
    return ret;
}

AlarmObjInfo *alarm_drdbschema_to_objectdb(AlarmData *ndata) {
    AlarmObjInfo *obj = dmt_malloc(sizeof(AlarmObjInfo));
    if (obj) {
        memset(obj, '\0', sizeof(AlarmObjInfo));
        /* All values are already read from sensors during alert trigger so just need to copy these*/
        obj->eventtype = ndata->eventtype.value.intval;
        obj->realtime = ndata->realtime.value.boolval;
        obj->state = ndata->state.value.sintval;
        strcpy(obj->disc, ndata->disc.value.stringval);
        obj->lowthreshold = ndata->lowthreshold.value.doubleval;
        obj->highthreshold = ndata->highthreshold.value.doubleval;
        obj->crithreshold = ndata->crithreshold.value.doubleval;
        obj->eventcount = ndata->eventcount.value.intval;
        obj->time = ndata->time.value.lintval;
        obj->sobjid = ndata->sobjid.value.sintval;
        obj->sinstid = ndata->sinstid.value.sintval;
        obj->srsrcid = ndata->srsrcid.value.sintval;
        obj->sensorvalue = ndata->sensorvalue.value.doubleval;
        strcpy(obj->sensorunits, ndata->sensorunits.value.stringval);
        strcpy(obj->applicationtype, ndata->applicationtype.value.stringval);
    }
    return obj;
}

/* Create a new Alarm node for reporting to client from the DRDBSchema of the Alarm instance. */
void alarmdb_push(DRDBSchema *node) {
    if (!node->data) {
        return;
    }
    AlarmData *adata = node->data;

    /* Create a Alarm Schema node from Registry schema data provided.*/
    AlarmSchema *cnode =
        alarmdb_create_node(&node->obj, adata->sobjid.value.sintval,
                            adata->sinstid.value.sintval,
                            adata->srsrcid.value.sintval);
    if (cnode) {
        AlarmSchema *snode = alarmdb_search_node(cnode);
        if (snode) {
            /* Update the entry in database i.e remove the old and add new one to front.*/
            alarmdb_remove(snode);
            /* Free snode */
            dmt_free(snode);
            alarmdb_push(node);
        } else {
            /* Translate DRDBschema to ObjectInfo */
            AlarmObjInfo *obj = alarm_drdbschema_to_objectdb(adata);
            if (obj) {
                /* Add new entry in the database*/
                AlarmSchema *aschema = dmt_malloc(sizeof(AlarmSchema));
                if (aschema) {
                    memset(aschema, '\0', sizeof(AlarmSchema));
                    aschema->instid = adata->sinstid.value.sintval;
                    aschema->objid = adata->sobjid.value.sintval;
                    aschema->rsrcid = adata->srsrcid.value.sintval;
                    aschema->reported = ALARM_NOT_REPORTED;
                    aschema->response = ALARM_RESPONSE_PENDING;
                    aschema->txcount = 0;
                    memcpy(&aschema->obj, &node->obj, sizeof(DevObj));
                    memcpy(&aschema->data, obj, sizeof(AlarmObjInfo));
                }

                /* Add to history data*/
                alramdb_record_data(aschema);

                /* Add to list in FIFO manner*/
                alarmdb_prepend(aschema);
                log_debug(
                    "ALARMHANDLER:: Alarm instance %d for Device /%d/%d/%d Name  %s, Disc: %s Module UUID: %s added.",
                    obj->instanceId, aschema->objid, aschema->instid,
                    aschema->rsrcid, aschema->obj.name, aschema->obj.disc,
                    aschema->obj.mod_UUID);
                dmt_free(aschema);
                dmt_free(obj);
            }
        }
    }
    dmt_free(cnode);
}

/* Remove node from the AlarmDB */
void alarmdb_pop(AlarmSchema *aschema) {
    if (aschema) {
        log_trace(
            "ALARMHANDLER:: Removing Alarm instance for %d for Device /%d/%d/%d Name %s, Disc: %s Module UUID: %s removed form Alarm Queue.",
            aschema->data.instanceId, aschema->objid, aschema->instid,
            aschema->rsrcid, aschema->obj.name, aschema->obj.disc,
            aschema->obj.mod_UUID);

        list_remove(&alarmdb, aschema);
    }
}

MsgFrame *alarmhandler_create_tx_frame(AlarmSchema *anode, bool newtoken) {
    MsgFrame *smsg = NULL;
    if (anode) {
        /* Create a if-message */
        smsg = create_msgframe(MSG_TYPE_ALERT_REP, anode->instid,
                               ALL_RESOURCE_ID, OBJ_TYPE_ALARM,
                               sizeof(AlarmObjInfo), newtoken, &anode->data);
    }
    return smsg;
}

int alarmhandler_verify_resp(MsgFrame *rmsg, MsgFrame *smsg, void *data,
                             int sflag) {
    int ret = 0;
    if (rmsg && smsg) {
        if (sflag) {
            /* Failure in processing request */
            log_error(
                "Err(%d): ALARMHANDLER:: Failure while handling response for "
                "Inst: %d of Type 0x%x RId: %d",
                ret, rmsg->instance, rmsg->objecttype, rmsg->resourceId);
        } else {
            if (rmsg) {
                if (!(rmsg->response)) {
                    /* No error */
                    /* Copy response */
                    if (rmsg->data) {
                        if (data) {
                            memcpy(data, rmsg->data, rmsg->datasize);
                        }
                    }
                } else {
                    /* Response failure from the Client */
                    ret = rmsg->response;
                    log_error(
                        "Err(%d): ALARMHANDLER:: Failure response from the client for %d of Type 0x%x RId: %d",
                        ret, rmsg->instance, rmsg->objecttype,
                        rmsg->resourceId);
                }
            }
        }
    }
    return ret;
}

/* Process Alarm response received from the client.*/
int alarmhandler_proc_alarm_resp(AlarmSchema *anode, MsgFrame *rmsg,
                                 MsgFrame *smsg, int flag) {
    int ret = 0;
    if (rmsg && (rmsg->msgtype == MSG_TYPE_ALERT_RESP)) {
        /* Verify response */
        ret = alarmhandler_verify_resp(rmsg, smsg, &anode->data, flag);
        if (!ret) {
            anode->response = true;
        } else {
            anode->response = false;
        }

        /* Add response to history */
        alramdb_record_data(anode);

    } else {
        ret = ERR_UNEXPECTED_RESP_MSG; /* ERR_UNEXPECTED_RESP_MSG */
    }
    return ret;
}

/* Transmit and receive response from the client for the Alarm */
int alramhandler_prepare_to_tx(AlarmSchema *anode, bool newalarm) {
    int ret = 0;
    int flag = 0;
    size_t size = sizeof(AlarmObjInfo);
    MsgFrame *rmsg = NULL;
    MsgFrame *smsg = alarmhandler_create_tx_frame(anode, newalarm);
    if (smsg) {
        anode->reported = true;
        if (newalarm) {
            anode->token = smsg->reqtoken;
        } else {
            smsg->reqtoken = anode->token;
        }
        anode->txcount++;

        /* Add response to history */
        alramdb_record_data(anode);

        /* Send msg and Receive */
        rmsg = msghandler_client_send(smsg, &size, &flag);
        if (rmsg) {
            ret = alarmhandler_proc_alarm_resp(anode, rmsg, smsg, flag);
        }
    } else {
        ret = ERR_UBSP_INVALID_POINTER;
    }
    free_msgframe(&smsg);
    free_msgframe(&rmsg);
    return ret;
}

/* Check if the alarm needs to be reported.
 * Return 1 on successful processing
 * Return 0 on nay failure
 */
int alarmhandler_proc_node(void *node) {
    int ret = 1;
    if (!node) {
        return ret;
    }
    /* node is pointer to data in the AlarmDB list */
    AlarmSchema *aschema = node;
    log_trace(
        "ALARMHANDLER:: Alarm instance %d for Device  /%d/%d/%d Name %s, Disc: %s Module UUID: %s is getting reported.",
        aschema->data.instanceId, aschema->objid, aschema->instid,
        aschema->rsrcid, aschema->obj.name, aschema->obj.disc,
        aschema->obj.mod_UUID);
    if (!aschema->reported) {
        /* New alarm */
        ret = alramhandler_prepare_to_tx(aschema, true);
    } else if (!aschema->response) {
        /* Not responded by client yet */
        ret = alramhandler_prepare_to_tx(aschema, false);
    }

    if (!ret) {
        /* Check if alarm is successfully transmitted and responded */
        if (aschema->reported && aschema->response) {
            /* Success */
            /* Reported and responded by client */
            log_debug(
                "ALARMHANDLER:: Alarm instance %d for Device  /%d/%d/%d Name %s, Disc: %s Module UUID: %s reported successfully.",
                aschema->data.instanceId, aschema->objid, aschema->instid,
                aschema->rsrcid, aschema->obj.name, aschema->obj.disc,
                aschema->obj.mod_UUID);
            alarmdb_pop(aschema);
        } else {
            ret = -1;
        }
    }

    if (ret) {
        log_debug(
            "ALARMHANDLER:: Alarm instance %d for Device  /%d/%d/%d Name %s, Disc: %s Module UUID: %s couldn't be reported..",
            aschema->data.instanceId, aschema->objid, aschema->instid,
            aschema->rsrcid, aschema->obj.name, aschema->obj.disc,
            aschema->obj.mod_UUID);
        alarmdb_update(aschema);
    }
    return ret;
}

void alarmhandler_service(size_t timer_id, void *data) {
    int ret = 0;
    /* Read data from every device type */
    log_trace("ALARMHANDLER:: Reporting %d active/clear alarms.\n",
              list_size(&alarmdb));
    if (list_size(&alarmdb))
        list_for_each(&alarmdb, &alarmhandler_proc_node);
}

void alarmhandler_start() {
    size_t ptimer = 0;
    galarmtimer = start_timer(ALARM_PERIODIC_TIMER, &alarmhandler_service,
                              TIMER_PERIODIC, NULL);
    if (galarmtimer) {
        log_debug(
            "ALARMHANDLER:: Periodic timer started with period of %d millisec.\n",
            (ALARM_PERIODIC_TIMER));
    }
}

void alarmhandler_stop(size_t timer) {
    stop_timer(timer);
}

void alarmhandler_init() {
    alarmdb_init();
    initialize(&galarmthreadid);
}

void alarmhandler_exit() {
    alarmhandler_stop(galarmtimer);
    finalize(galarmthreadid);
    alarmdb_exit();
}
