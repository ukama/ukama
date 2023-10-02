/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <time.h>

#include "jserdes.h"
#include "lookout.h"
#include "errorcode.h"
#include "json_types.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

/* define in resources.c */
extern int get_memory_usage(int pid);
extern int get_disk_usage(int pid);
extern double get_cpu_usage(int pid);
extern char *get_radio_status(void);

void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static void add_capp_to_list(CappList **list,
                             const char *name,
                             const char *tag,
                             const char *status,
                             int pid) {

    CappList *ptr=NULL;

    if (name == NULL || tag == NULL ||
        status == NULL) return;

    if (*list == NULL) { /* First entry */
        *list = (CappList *)calloc(1, sizeof(CappList));
        if (*list == NULL) return;
        ptr = *list;
    } else {
        (*list)->next = (CappList *)calloc(1, sizeof(CappList));
        if ((*list)->next == NULL) return;
        ptr = (*list)->next;
    }

    ptr->capp          = (Capp *)calloc(1, sizeof(Capp));
    ptr->capp->runtime = (CappRuntime *)calloc(1, sizeof(CappRuntime));

    ptr->capp->name            = strdup(name);
    ptr->capp->tag             = strdup(tag);
    ptr->capp->runtime->status = strdup(status);
    ptr->capp->runtime->pid    = pid;
    ptr->capp->runtime->memory = -1;
    ptr->capp->runtime->disk   = -1;
    ptr->capp->runtime->memory = -1;
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue,
                           double *doubleValue) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING): 
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = json_integer_value(jEntry);
        break;
    case (JSON_REAL):
        *doubleValue = json_real_value(jEntry);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

/*
 * deserialize_node_info --
 *
 * {
 * "nodeInfo": {
 *   "UUID": "ukma-7001-tnode-sa03-1100",
 *   "name": "tNode",
 *   "type": 2,
 *   "partNumber": "LTE-BAND-3-0XXXX",
 *   "skew": "UK_TNODE-LTE-0001",
 *   "mac": "10:20:30:20:50:60",
 *   "prodSwVersion": {
 *     "major": 1,
 *     "minor": 1
 *   },
 *   "swVersion": {
 *     "major": 0,
 *     "minor": 0
 *   },
 *   "assemblyDate": "30-07-2020",
 *   "oemName": "SANMINA",
 *   "moduleCount": 3
 * }
 *}
 *
 */
bool json_deserialize_node_id(char **nodeID, JsonObj *json) {

    JsonObj *jNodeInfo=NULL;

    if (json == NULL) return USYS_FALSE;

    jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }
    
    if (get_json_entry(jNodeInfo, JTAG_UUID, JSON_STRING,
                       nodeID, NULL, NULL) == USYS_FALSE) {
        log_error("Error deserializing node info");
        json_log(json);
        *nodeID = NULL;

        return USYS_FALSE;
    }
    
    return USYS_TRUE;
}

/*
 * "capps" : [
 *    {
 *       "name" : "example",
 *      "tag" : "0.0.1",
 *      "status" : "run",
 *      "pid" : "123"
 *    }
 * ]
 */
bool json_deserialize_capps(CappList **cappList, JsonObj *json) {

    int i=0, count=0;
    JsonObj *jCapp=NULL, *jArray=NULL;
    JsonObj *jName=NULL, *jTag=NULL;
    JsonObj *jStatus=NULL, *jPid=NULL;

    if (json == NULL) {
        return USYS_FALSE;
    }

    jArray = json_object_get(json, JTAG_CAPPS);
    if (!json_is_array(jArray)) {
        return USYS_FALSE;
    }

    count = json_array_size(jArray);
    if (count == 0) return USYS_FALSE;

    for (i=0; i<count; i++) {
        jCapp = json_array_get(jArray, i);

        if (jCapp == NULL) continue;

        jName   = json_object_get(jCapp, JTAG_NAME);
        jTag    = json_object_get(jCapp, JTAG_TAG);
        jStatus = json_object_get(jCapp, JTAG_STATUS);
        jPid    = json_object_get(jCapp, JTAG_PID);

        if (jName && jTag && jStatus && jPid) {
            add_capp_to_list(cappList,
                             json_string_value(jName),
                             json_string_value(jTag),
                             json_string_value(jStatus),
                             json_integer_value(jPid));
        }
    }

    usys_log_debug("Recevied %d capps from starter.d", count);

    return USYS_TRUE;
}

static void json_add_resources_to_capp_report(JsonObj **json,
                                              CappRuntime *runtime) {

    JsonObj *jArray = NULL;
    JsonObj *jMemory = NULL, *jDisk = NULL, *jCpu = NULL;

    char buffer[MAX_BUFFER] = {0};

    json_object_set_new(*json, JTAG_RESOURCES, json_array());
    jArray = json_object_get(*json, JTAG_RESOURCES);
    if (jArray == NULL) return;

    jMemory = json_object();
    jDisk   = json_object();
    jCpu    = json_object();

    /* memory */
    sprintf(buffer, "%d", runtime->memory);
    json_object_set_new(jMemory,
                        JTAG_NAME,
                        json_string("memory"));
    json_object_set_new(jMemory,
                        JTAG_VALUE,
                        json_string(buffer));

    /* disk */
    sprintf(buffer, "%d", runtime->disk);
    json_object_set_new(jDisk,
                        JTAG_NAME,
                        json_string("disk"));
    json_object_set_new(jDisk,
                        JTAG_VALUE,
                        json_string(buffer));

    /* cpu */
    sprintf(buffer, "%f", runtime->cpu);
    json_object_set_new(jCpu,
                        JTAG_NAME,
                        json_string("cpu"));
    json_object_set_new(jCpu,
                        JTAG_VALUE,
                        json_string(buffer));

    json_array_append_new(jArray, jMemory);
    json_array_append_new(jArray, jDisk);
    json_array_append_new(jArray, jCpu);
}

static void json_add_system_info_to_report(JsonObj **json) {

    JsonObj *jArray = NULL;
    JsonObj *jRadio = NULL;

    jRadio = json_object();

    /* system */
    json_object_set_new(*json, JTAG_SYSTEM, json_array());
    jArray = json_object_get(*json, JTAG_SYSTEM);
    if (jArray == NULL) return;

    /* radio */
    json_object_set_new(jRadio,
                        JTAG_NAME,
                        json_string("radio"));
    json_object_set_new(jRadio,
                        JTAG_VALUE,
                        json_string(get_radio_status()));

    json_array_append_new(jArray, jRadio);
}

/*
http://localhost:8080/v1/health/{nodeID}
{
  "nodeID" : "ukma-xx-xxx-xxxx-xxx"
  "timestamp": "12345678",
  "system": [
    {
      "name": "radio",
      "value": "off"
    }
  ],
  "capps": [
    {
      "name": "bootstrap",
      "tag": "0.0.1",
      "status": "run",
      "resources": [
        {
          "name": "disk",
          "value": "3456"
        },
        {
          "name": "memory",
          "value": "12345"
        }
      ]
    }
  ]
}
*/
bool json_serialize_health_report(JsonObj **json,
                                  char *nodeID,
                                  CappList *list) {

    JsonObj *jArray     = NULL;
    JsonObj *jCapp      = NULL;
    JsonObj *jResources = NULL;
    CappList *ptr       = NULL;
    Capp     *capp      = NULL;
    char     buffer[MAX_BUFFER] = {0};
    time_t   currTime;

    *json = json_object();
    if (*json == NULL) return USYS_FALSE;

    /* nodeID */
    json_object_set_new(*json,
                        JTAG_NODE_ID,
                        json_string(nodeID));

    /* time-stamp */
    time(&currTime);
    sprintf(buffer, "%ld", (long)currTime);
    json_object_set_new(*json,
                        JTAG_TIMESTAMP,
                        json_string(buffer));

    /* capps */
    json_object_set_new(*json, JTAG_CAPPS, json_array());
    jArray = json_object_get(*json, JTAG_CAPPS);
    if (jArray == NULL) return USYS_FALSE;

    for (ptr = list; ptr; ptr=ptr->next) {
        capp = ptr->capp;

        jCapp      = json_object();
        jResources = json_object();
        if (jCapp == NULL || jResources == NULL) return USYS_FALSE;

        json_object_set_new(jCapp,
                            JTAG_NAME,
                            json_string(capp->name));
        json_object_set_new(jCapp,
                            JTAG_TAG,
                            json_string(capp->tag));
        json_object_set_new(jCapp,
                            JTAG_STATUS,
                            json_string(capp->runtime->status));

        json_add_resources_to_capp_report(&jCapp, capp->runtime);
        json_array_append_new(jArray, jCapp);
    }

    /* system */
    json_add_system_info_to_report(json);

    return USYS_TRUE;
}

void json_free(JsonObj** json) {
    if (*json){
        json_decref(*json);
        *json = NULL;
    }
}
