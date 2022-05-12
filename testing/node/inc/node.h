/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef NODE_H
#define NODE_H

#include <jansson.h>

#define TRUE 1
#define FALSE 0

#define NODE_TYPE_HNODE "hnode"
#define NODE_TYPE_ANODE "anode"
#define NODE_TYPE_TNODE "tnode"
#define NODE_TYPE_NONE  ""

#define MODULE_TYPE_COM  "com"
#define MODULE_TYPE_TRX  "trx"
#define MODULE_TYPE_CTRL "ctrl"
#define MODULE_TYPE_FE   "fe"
#define MODULE_TYPE_MASK "mask"
#define MODULE_TYPE_NONE ""

#define SCHEMA_FILE_COM  "com.json"
#define SCHEMA_FILE_TRX  "trx.json"
#define SCHEMA_FILE_MASK "mask.json"
#define SCHEMA_FILE_CTRL "ctrl.json"
#define SCHEMA_FILE_FE   "fe.json"
#define SCHEMA_FILE_NONE ""

/* overall node info */
typedef struct nodeInfo_ {

	char *uuid;
	char *type;
	char *partNumber;
	char *skew;
	char *mac;
	char *swVersion;
	char *mfgSWVersion;
	char *assemblyDate;
	char *oem;
    char *mfgTestStatus;
    char *status;
	int  moduleCount;
} NodeInfo;

/* For module(s) */
typedef struct nodeConfig_ {

	char *moduleID;
	char *type;
	char *partNumber;
	char *hwVersion;
	char *mac;
	char *swVersion;
	char *mfgSWVersion;
	char *mfgDate;
	char *oem;
	char *status;

	struct nodeConfig_ *next;
} NodeConfig;

/* Node meta data via VNODE_METADATA */
typedef struct node_ {

	NodeInfo   *nodeInfo;
	NodeConfig *nodeConfig;
} Node;

void free_node(Node *node);
void free_node_config(NodeConfig *nodeConfig);
int add_node_config_entry(Node **node, NodeConfig *nodeConfig);

#endif /* NODE_H */
