/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "jserdes.h"
#include "notification.h"
#include "errorcode.h"
#include "json_types.h"
#include "node.h"
#include "web_service.h"
#include "notify/notify.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"


void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
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
bool json_deserialize_node_id(char **nodeID, json_t *json) {

    json_t *jNodeInfo=NULL;

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
{
    "serviceName" : "noded",
    "time"        : 1234567,
    "status"      : xxx,
    "type"        : "alert"
    "nodeID"      : "ukma-aaa-bbb-cccc",
    "details": {
        "module"   : "trx"
        "property" : "metric-name",
        "value"    : "xxx",
        "units"    : "milli-seconds",
        "description" : "User are too many"
    }
}
*/

bool json_serialize_notification(JsonObj **json, Notification* notification,
                                 char *type, char *nodeID, int statusCode) {
                                
    JsonObj *jDetails=NULL;

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_SERVICE_NAME,
                        json_string(notification->serviceName));

    json_object_set_new(*json, JTAG_TIME,
                        json_integer(notification->epochTime));

    json_object_set_new(*json, JTAG_STATUS, json_integer(statusCode));
    json_object_set_new(*json, JTAG_TYPE, json_string(type));
    json_object_set_new(*json, JTAG_NODE_ID, json_string(nodeID));

    /* Add details about the event/alarm */
    json_object_set_new(*json, JTAG_DETAILS, json_object());
    jDetails = json_object_get(*json, JTAG_DETAILS);
    
    json_object_set_new(jDetails, JTAG_MODULE,
                        json_string(notification->module));
    json_object_set_new(jDetails, JTAG_NAME,
                        json_string(notification->propertyName));
    json_object_set_new(jDetails, JTAG_VALUE,
                        json_string(notification->propertyValue));
    json_object_set_new(jDetails, JTAG_UNITS,
                        json_string(notification->propertyUnit));
    json_object_set_new(jDetails, JTAG_DESCRIPTION,
                        json_string(notification->details));

    return USYS_TRUE;
}

/* Deserialize generic notification received from local services */
bool json_deserialize_notification(JsonObj *json,
                                   Notification **notification) {

    bool ret=USYS_TRUE;

    if (json == NULL) {
        usys_log_error("No data to deserialize");
        return USYS_FALSE;
    }

    *notification = (Notification *)calloc(1, sizeof(Notification));
    if (*notification == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(Notification));
        return USYS_FALSE;
    }
    
    ret |= get_json_entry(json, JTAG_SERVICE_NAME, JSON_STRING,
                          &(*notification)->serviceName, NULL, NULL);
    ret |= get_json_entry(json, JTAG_SEVERITY, JSON_STRING,
                          &(*notification)->severity, NULL, NULL);
    ret |= get_json_entry(json, JTAG_TIME, JSON_INTEGER,
                          NULL, &(*notification)->epochTime, NULL);
    ret |= get_json_entry(json, JTAG_NAME, JSON_STRING,
                          &(*notification)->propertyName, NULL, NULL);
    ret |= get_json_entry(json, JTAG_VALUE, JSON_STRING,
                          &(*notification)->propertyValue, NULL, NULL);
    ret |= get_json_entry(json, JTAG_UNITS, JSON_STRING,
                          &(*notification)->propertyUnit, NULL, NULL);
    ret |= get_json_entry(json, JTAG_DETAILS, JSON_STRING,
                          &(*notification)->details, NULL, NULL);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the notifiction JSON");
        json_log(json);
        free_notification(*notification);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

/* Decrement json references */
void json_free(JsonObj** json) {
    if (*json){
        json_decref(*json);
        *json = NULL;
    }
}
