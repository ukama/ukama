/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2024-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "usys_types.h"
#include "usys_log.h"
#include "usys_mem.h"

#include "jserdes.h"
#include "nodeInfo.h"

static void log_json(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        usys_log_debug("json str: %s", str);
        usys_free(str);
    }
}

static int get_json_entry(json_t *json, char *key, json_type type,
						  char **strValue, int *intValue) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        usys_log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    if (type == JSON_STRING) {
        *strValue = strdup(json_string_value(jEntry));
    } else if (type == JSON_INTEGER) {
        *intValue = json_integer_value(jEntry);
    } else {
        usys_log_error("Invalid type for json key-value: %d", type);
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
 *   "oemName": "abc",
 *   "moduleCount": 3
 * }
 *}
 *
 */
int deserialize_node_info(NodeInfo **nodeInfo, json_t *json) {

    int ret=USYS_TRUE;
    json_t *jNodeInfo=NULL;

    if (json == NULL) return USYS_FALSE;

    log_json(json);
	
    jNodeInfo = json_object_get(json, JSON_NODE_INFO);

    if (jNodeInfo == NULL) {
        usys_log_error("Missing mandatory %s from JSON", JSON_NODE_INFO);
        return USYS_FALSE;
    }

    *nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
    if (*nodeInfo == NULL) {
        usys_log_error("Error allocating memory of size: %lu", sizeof(NodeInfo));
        return USYS_FALSE;
    }

    ret |= get_json_entry(jNodeInfo, JSON_UUID, JSON_STRING,
                          &(*nodeInfo)->uuid, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_NAME, JSON_STRING,
                          &(*nodeInfo)->name, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_PART_NUMBER, JSON_STRING,
                          &(*nodeInfo)->partNumber, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_SKEW, JSON_STRING,
                          &(*nodeInfo)->skew, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_MAC, JSON_STRING,
                          &(*nodeInfo)->mac, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_ASSEMBLY_DATE, JSON_STRING,
                          &(*nodeInfo)->assemblyDate, NULL);
    ret |= get_json_entry(jNodeInfo, JSON_OEM, JSON_STRING,
                          &(*nodeInfo)->oem, NULL);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing node info");
        log_json(json);
        free_node_info(*nodeInfo);
        *nodeInfo = NULL;
    }

    return ret;
}
