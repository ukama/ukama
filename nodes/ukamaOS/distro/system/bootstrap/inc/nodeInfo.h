/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef NODE_INFO_H
#define NODE_INFO_H

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

#define MAX_URL_LEN 1024

/* overall node info */
typedef struct nodeInfo_ {

	char *uuid;
    char *name;
	char *partNumber;
	char *skew;
	char *mac;
	char *assemblyDate;
	char *oem;
} NodeInfo;

int get_nodeID_from_noded(char **nodeID, char *host, char *port);
void free_node_info(NodeInfo *nodeInfo);

#endif /* NODE_INFO_H */
