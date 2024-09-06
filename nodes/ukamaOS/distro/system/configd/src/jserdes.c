/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include "jserdes.h"
#include "configd.h"
#include "errorcode.h"
#include "json_types.h"
#include "web_service.h"

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
        if (jEntry->type == JSON_NULL) {
        	*strValue = NULL;
        } else {
        	*strValue = strdup(json_string_value(jEntry));
        }
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
"filename":"version.json",
"app":"configd",
"version":"abcdef",
"timestamp":189295,
"reason": 1,
"file_count":4,
"data": ""}
*/
bool json_deserialize_session_data(JsonObj *json,
                                   SessionData **sd) {

    int ret = USYS_TRUE;
    
    if (json == NULL) return USYS_FALSE;

    *sd = (SessionData *)usys_calloc(1, sizeof(SessionData));
    if (*sd == NULL) {
        usys_log_error("Error allocating memory of size: %d",
                       sizeof(SessionData));
        return USYS_FALSE;
    }

    ret &= get_json_entry(json, JTAG_FILE_NAME, JSON_STRING,
                          &(*sd)->fileName, NULL, NULL);
    ret &= get_json_entry(json, JTAG_APP_NAME, JSON_STRING,
                          &(*sd)->app, NULL, NULL);
    ret &= get_json_entry(json, JTAG_TIME_STAMP, JSON_INTEGER,
                          NULL, &(*sd)->timestamp, NULL);
    ret &= get_json_entry(json, JTAG_REASON, JSON_INTEGER,
                          NULL, &(*sd)->reason, NULL);
    ret &= get_json_entry(json, JTAG_DATA, JSON_STRING,
                          &(*sd)->data, NULL, NULL);
    ret &= get_json_entry(json, JTAG_VERSION, JSON_STRING,
                          &(*sd)->version, NULL, NULL);
    ret &= get_json_entry(json, JTAG_FILE_COUNT, JSON_INTEGER,
                          NULL, &(*sd)->fileCount, NULL);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the config JSON");
        json_log(json);
        free_session_data(*sd);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}


