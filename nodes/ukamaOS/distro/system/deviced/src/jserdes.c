/**
 * Copyright (c) 2023-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "web_service.h"
#include "deviced.h"
#include "json_types.h"
#include "config.h"

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

bool json_serialize_alarm_notification(JsonObj **json,
                                       Config *config) {

    *json = json_object();
    if (!*json) return USYS_FALSE;

    json_object_set_new(*json, JTAG_SERVICE_NAME,
                        json_string(config->serviceName));
    json_object_set_new(*json, JTAG_SEVERITY, json_string(ALARM_HIGH));
    json_object_set_new(*json, JTAG_TIME,     json_integer(time(NULL)));
    json_object_set_new(*json, JTAG_MODULE,   json_string(MODULE_NONE));
    json_object_set_new(*json, JTAG_NAME,     json_string(ALARM_NODE));
    json_object_set_new(*json, JTAG_VALUE,    json_string(ALARM_REBOOT));
    json_object_set_new(*json, JTAG_UNITS,    json_string(EMPTY_STRING));
    json_object_set_new(*json, JTAG_DETAILS,  json_string(ALARM_REBOOT_DESCRP));

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
bool json_deserialize_node_info(char **data, char *tag, json_t *json) {

    json_t *jNodeInfo=NULL;

    if (json == NULL) return USYS_FALSE;

    jNodeInfo = json_object_get(json, JTAG_NODE_INFO);
    if (jNodeInfo == NULL) {
        log_error("Missing mandatory %s from JSON", JTAG_NODE_INFO);
        return USYS_FALSE;
    }

    if (get_json_entry(jNodeInfo, tag, JSON_STRING,
                       data, NULL, NULL) == USYS_FALSE) {
        log_error("Error deserializing node info. tag: %s", tag);
        json_log(json);
        *data = NULL;
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
