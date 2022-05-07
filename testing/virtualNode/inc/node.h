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

/* For module(s) */
typedef struct nodeConfig_ {

	char *uuid;
	char *name;
	int  type;
	char *partNumber;
	char *hwVersion;
	char *mac;
	int  swVersionMinor;
	int  swVersionMajor;
	int  prodSWVersionMinor;
	int  prodSWVersionMajor;
	char *manufacturingDate;
	char *manufacturerName;

	struct nodeConfig_ *next;
} NodeConfig;

/* overall node info */
typedef struct nodeInfo_ {

	char *uuid;
	char *name;
	int  type;
	char *partNumber;
	char *skew;
	char *mac;
	int  swVersionMinor;
	int  swVersionMajor;
	int  prodSWVersionMinor;
	int  prodSWVersionMajor;
	char *assemblyDate;
	char *oemName;
	int  moduleCount;
} NodeInfo;

/* Node meta data via VNODE_METADATA */
typedef struct node_ {

	NodeInfo   *nodeInfo;
	NodeConfig *nodeConfig;
} Node;

void free_node(Node *node);
void free_node_config(NodeConfig *nodeConfig);
int add_node_config_entry(Node **node, NodeConfig *nodeConfig);

#endif /* NODE_H */
