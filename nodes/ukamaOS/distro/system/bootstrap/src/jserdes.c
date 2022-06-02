/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <stdio.h>
#include <string.h>
#include <jansson.h>

#include "jserdes.h"
#include "nodeInfo.h"
#include "log.h"

static void log_json(json_t *json);
static int get_json_entry(json_t *json, char *key, int type, char **strValue,
						  int *intValue);
/*
 * log_json --
 *
 */
static void log_json(json_t *json) {

	char *str = NULL;

	str = json_dumps(json, 0);
	if (str) {
		log_debug("json str: %s", str);
		free(str);
	}
}

/*
 * get_json_entry --
 *
 */
static int get_json_entry(json_t *json, char *key, int type, char **strValue,
						  int *intValue) {

	json_t *jEntry=NULL;

	if (json == NULL || key == NULL) return FALSE;

	jEntry = json_object_get(json, key);
	if (jEntry == NULL) {
		log_error("Missing %s key in json", key);
		return FALSE;
	}

	if (type == JSON_STRING) {
		*strValue = strdup(json_string_value(jEntry));
	} else if (type == JSON_INTEGER) {
		*intValue = json_integer_value(jEntry);
	} else {
		log_error("Invalid type for json key-value: %d", type);
		return FALSE;
	}

	return TRUE;
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
int deserialize_node_info(NodeInfo **nodeInfo, json_t *json) {

	int ret=TRUE;
	json_t *jNodeInfo=NULL;

	if (json == NULL) return FALSE;

	log_json(json);
	
	jNodeInfo = json_object_get(json, JSON_NODE_INFO);

	if (jNodeInfo == NULL) {
		log_error("Missing mandatory %s from JSON", JSON_NODE_INFO);
		return FALSE;
	}

	*nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
	if (*nodeInfo == NULL) {
		log_error("Error allocating memory of size: %lu", sizeof(NodeInfo));
		return FALSE;
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

	if (ret == FALSE) {
		log_error("Error deserializing node info");
		log_json(json);
		free_node_info(*nodeInfo);
		*nodeInfo = NULL;
	}

	return ret;
}

/*
 * deserialize_server_info --
 *
 * {
 *   "node": "uk-sa2220-hnode-v0-dcf4",
 *   "org": "test",
 *   "ip": "192.168.0.1",
 *   "certificate": "aGVscG1lCg=="
 * }
 *
 */
int deserialize_server_info(ServerInfo *serverInfo, json_t *json) {

	int ret=TRUE;

	if (serverInfo == NULL || json == NULL) return FALSE;

	log_json(json);

	ret |= get_json_entry(json, JSON_IP,  JSON_STRING, &serverInfo->IP,  NULL);
	ret |= get_json_entry(json, JSON_ORG, JSON_STRING, &serverInfo->org, NULL);
	ret |= get_json_entry(json, JSON_CERTIFICATE, JSON_STRING,
						  &serverInfo->cert, NULL);

	if (ret == FALSE) {
		log_error("Error deserializing server info");
		log_json(json);
		free_server_info(serverInfo);
	}

	return ret;
}
