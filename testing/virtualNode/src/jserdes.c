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
#include "node.h"
#include "log.h"

static void log_json(json_t *json);
static int deserialize_node_info(Node **node, json_t *json);
static int deserialize_node_config(Node **node, json_t *json, int count);

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
 * deserialize_node --
 *
 * {
 * "nodeInfo": {
 *  "UUID": "UK-9001-HNODE-SA03-1103",
 *  "name": "tNode",
 *  "type": 2,
 *  "partNumber": "LTE-BAND-3-0XXXX",
 *  "skew": "UK_TNODE-LTE-0001",
 *  "mac": "10:20:30:20:50:60",
 *  "swVersion": {
 *    "major": 0,
 *    "minor": 1
 *  },
 *  "prodSwVersion": {
 *    "major": 1,
 *    "minor": 1
 *  },
 *  "assemblyDate": "30-07-2020",
 *  "oemName": "SANMINA",
 *  "moduleCount": 3
 *},
 *"nodeConfig": [
 *  {
 *    "UUID": "UK-9001-COM-1103",
 *    "name": "COM",
 *    "type": 0,
 *    "partNumber": "COMv1-X86-0XXXX",
 *    "hwVersion": "REV-A",
 *    "mac": "10:20:30:20:50:60",
 *    "swVersion": {
 *      "major": 0,
 *      "minor": 1
 *    },
 *    "prodSwVersion": {
 *      "major": 1,
 *        "minor": 1
 *    },
 *    "manufacturingDate": "29-07-2020",
 *    "manufacturerName": "FOXCON"
 *  },
 * ]
 *}
 *
 */
int deserialize_node(Node **node, json_t *json) {

	json_t *jNodeInfo=NULL;
	json_t *jNodeConfig=NULL;
	int count;

	if (json == NULL) return FALSE;

	jNodeInfo   = json_object_get(json, JSON_NODE_INFO);
	jNodeConfig = json_object_get(json, JSON_NODE_CONFIG);

	if (jNodeInfo == NULL || jNodeConfig == NULL) {
		log_error("Missing mandatory %s or %s from the env variable",
				  JSON_NODE_INFO, JSON_NODE_CONFIG);
		return FALSE;
	}

	if (!json_is_array(jNodeConfig)) {
		log_error("Expecting %s array but is missing", JSON_NODE_CONFIG);
		log_json(json);
		return FALSE;
	}

	count = json_array_size(jNodeConfig);
	if (count == 0) {
		log_error("%s array with no element.", JSON_NODE_CONFIG);
		log_json(jNodeConfig);
		return FALSE;
	}

	*node = (Node *)calloc(1, sizeof(Node));
	if (*node == NULL) {
		log_error("Error allocating memory of size: %lu", sizeof(Node));
		return FALSE;
	}

	(*node)->nodeInfo = (NodeInfo *)calloc(1, sizeof(NodeInfo));
	if ((*node)->nodeInfo == NULL) {
		log_error("Error allocating memory of size: %lu", sizeof(NodeInfo));
		goto failure;
	}
	if (!deserialize_node_info(node, jNodeInfo)) {
		log_error("Error deserializing node info");
		goto failure;
	}

	/* check to see if we have config for all modules */
	if (count != (*node)->nodeInfo->moduleCount) {
	  log_error("Module count from nodeInfo: %d mismatch with array: %d",
				(*node)->nodeInfo->moduleCount, count);
	  goto failure;
	}

	if (!deserialize_node_config(node, jNodeConfig, count)) {
		log_error("Error deserializing node config");
		goto failure;
	}

	return TRUE;

 failure:
	log_error("Error deserializing node info");
	log_json(json);
	free_node(*node);
	*node = NULL;

	return FALSE;
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
 */
static int deserialize_node_info(Node **node, json_t *json) {

	int ret=TRUE;
	NodeInfo *nodeInfo=NULL;

	nodeInfo = (*node)->nodeInfo;
	if (nodeInfo == NULL) return FALSE;

	ret |= get_json_entry(json, JSON_UUID, JSON_STRING,
						  &nodeInfo->uuid, NULL);
	ret |= get_json_entry(json, JSON_NAME, JSON_STRING,
						  &nodeInfo->name, NULL);
	ret |= get_json_entry(json, JSON_TYPE, JSON_INTEGER,
						  NULL, &nodeInfo->type);
	ret |= get_json_entry(json, JSON_PART_NUMBER, JSON_STRING,
						  &nodeInfo->partNumber, NULL);
	ret |= get_json_entry(json, JSON_SKEW, JSON_STRING,
						  &nodeInfo->skew, NULL);
	ret |= get_json_entry(json, JSON_MAC, JSON_STRING,
						  &nodeInfo->mac, NULL);
	ret |= get_json_entry(json, JSON_ASSEMBLY_DATE, JSON_STRING,
						  &nodeInfo->assemblyDate, NULL);
	ret |= get_json_entry(json, JSON_OEM_NAME, JSON_STRING,
						  &nodeInfo->oemName, NULL);
	ret |= get_json_entry(json, JSON_MODULE_COUNT, JSON_INTEGER, NULL,
						  &nodeInfo->moduleCount);

	return ret;
}

/*
 * deserialize_node_config_elem --
 *
 */
static int deserialize_node_config_elem(NodeConfig **nodeConfig, json_t *json) {

	int ret=TRUE;

	if (json == NULL) return FALSE;

	*nodeConfig = (NodeConfig *)calloc(1, sizeof(NodeConfig));
	if (nodeConfig == NULL) {
	  log_error("Error allocating Memory of size: %lu", sizeof(NodeConfig));
	  return FALSE;
	}
	
	ret |= get_json_entry(json, JSON_UUID, JSON_STRING,
						  &(*nodeConfig)->uuid, NULL);
	ret |= get_json_entry(json, JSON_NAME, JSON_STRING,
						  &(*nodeConfig)->name, NULL);
	ret |= get_json_entry(json, JSON_TYPE, JSON_INTEGER,
						  NULL, &(*nodeConfig)->type);
	ret |= get_json_entry(json, JSON_PART_NUMBER, JSON_STRING,
						  &(*nodeConfig)->partNumber, NULL);
	ret |= get_json_entry(json, JSON_HW_VERSION, JSON_STRING,
						  &(*nodeConfig)->hwVersion, NULL);
	ret |= get_json_entry(json, JSON_MAC, JSON_STRING,
						  &(*nodeConfig)->mac, NULL);
	ret |= get_json_entry(json, JSON_MANUFACTURING_DATE, JSON_STRING,
						  &(*nodeConfig)->manufacturingDate, NULL);
	ret |= get_json_entry(json, JSON_MANUFACTURER_NAME, JSON_STRING,
						  &(*nodeConfig)->manufacturerName, NULL);

	return ret;
}

/*
 * deserialize_node_config --
 *
 */
static int deserialize_node_config(Node **node, json_t *json, int count) {

	int i;
	json_t *jElem=NULL;
	NodeConfig *nodeConfig = NULL;

	if (node==NULL || json==NULL) return FALSE;

	for (i=0; i<count; i++) {

		jElem = json_array_get(json, i);
		if (!deserialize_node_config_elem(&nodeConfig, jElem)) {
			free_node_config(nodeConfig);
			return FALSE;
		}

		if (!add_node_config_entry(node, nodeConfig)) {
			free_node_config(nodeConfig);
			return FALSE;
		}
	}

	return TRUE;
}
