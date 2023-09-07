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

#include "node.h"
#include "log.h"

static void free_node_info(NodeInfo *nodeInfo);

/*
 * add_node_config_entry --
 *
 */
int add_node_config_entry(Node **node, NodeConfig *nodeConfig) {

	NodeConfig *ptr=NULL;

	if ((*node)->nodeConfig == NULL) {
		(*node)->nodeConfig = nodeConfig;
	} else {
		for (ptr=(*node)->nodeConfig; ptr->next; ptr=ptr->next);
		ptr->next = nodeConfig;
	}

	return TRUE;
}

/*
 * free_node --
 *
 */
void free_node(Node *node) {

	NodeConfig *ptr, *next;

	if (node == NULL) return;

	free_node_info(node->nodeInfo);

	ptr = node->nodeConfig;

	while (ptr) {
	
		next = ptr->next;
		free_node_config(ptr);
		ptr = next;
	}

	free(node);
}

/*
 * free_node_info --
 *
 */
static void free_node_info(NodeInfo *nodeInfo) {

	NodeInfo *ptr;

	if (nodeInfo == NULL) return;

	ptr = nodeInfo;

	if (ptr->uuid)          free(ptr->uuid);
	if (ptr->type)          free(ptr->type);
	if (ptr->partNumber)    free(ptr->partNumber);
	if (ptr->skew)          free(ptr->skew);
	if (ptr->mac)           free(ptr->mac);
	if (ptr->swVersion)     free(ptr->swVersion);
	if (ptr->mfgSWVersion)  free(ptr->mfgSWVersion);
	if (ptr->assemblyDate)  free(ptr->assemblyDate);
	if (ptr->oem)           free(ptr->oem);
	if (ptr->mfgTestStatus) free(ptr->mfgTestStatus);
	if (ptr->status)        free(ptr->status);

	free(ptr); ptr=NULL;
}

/*
 * free_node_config --
 *
 */
void free_node_config(NodeConfig *nodeConfig) {

	NodeConfig *ptr;

	if (nodeConfig == NULL) return;

	ptr = nodeConfig;

	if (ptr->moduleID)     free(ptr->moduleID);
	if (ptr->type)         free(ptr->type);
	if (ptr->partNumber)   free(ptr->partNumber);
	if (ptr->hwVersion)    free(ptr->hwVersion);
	if (ptr->mac)          free(ptr->mac);
	if (ptr->swVersion)    free(ptr->swVersion);
	if (ptr->mfgSWVersion) free(ptr->mfgSWVersion);
	if (ptr->mfgDate)      free(ptr->mfgDate);
	if (ptr->oem)          free(ptr->oem);
	if (ptr->status)       free(ptr->status);

	free(ptr); ptr=NULL;
}
