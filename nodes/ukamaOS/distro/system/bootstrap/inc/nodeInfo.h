/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
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

int get_nodeID_from_noded(char **nodeID, char *host, int port);
void free_node_info(NodeInfo *nodeInfo);

#endif /* NODE_INFO_H */
