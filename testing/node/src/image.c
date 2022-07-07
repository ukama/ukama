/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* create virtual node */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <string.h>

#include "log.h"
#include "node.h"
#include "image.h"
#include "config.h"
#include "supervisor.h"

static char *module_schema_file(char* nodeType, char *type);
static void set_schema_args(Node *node, char **buffer);
static FILE* init_container_file(char *fileName);
static int write_to_container_file(char *buffer, char *fileName, FILE *fp);
static int create_container_file(char *target, Configs *config, Node *node);

/*
 * module_schema_file --
 *
 */
static char *module_schema_file(char* nodeType, char *type) {

	if (strcmp(type, MODULE_TYPE_COM) == 0)  return SCHEMA_FILE_COM;

	/* TRX module schema is selected based on node type */
	if (strcmp(nodeType, NODE_TYPE_TNODE) == 0) {
		if (strcmp(type, MODULE_TYPE_TRX) == 0)  return SCHEMA_FILE_TRX;
	} else if ((strcmp(nodeType, NODE_TYPE_HNODE) == 0)) {
		if (strcmp(type, MODULE_TYPE_TRX) == 0)  return SCHEMA_FILE_HNODE_TRX;
	}

	if (strcmp(type, MODULE_TYPE_MASK) == 0) return SCHEMA_FILE_MASK;
	if (strcmp(type, MODULE_TYPE_CTRL) == 0) return SCHEMA_FILE_RFCTRL;
	if (strcmp(type, MODULE_TYPE_FE) == 0)   return SCHEMA_FILE_RFFE;

	return SCHEMA_FILE_NONE;
}

/*
 * set_schema_args
 *
 */
static void set_schema_args(Node *node, char **buffer) {

	char temp[MAX_BUFFER]={0}, temp1[MAX_BUFFER]={0};
	NodeConfig *ptr=NULL;

	if (node == NULL || *buffer == NULL) return;
	if (strcmp(node->nodeInfo->type, NODE_TYPE_NONE) == 0) return;

	ptr = node->nodeConfig;
	while (ptr) {
		sprintf(temp, " --n %s --m %s --f ./schemas/%s",
				ptr->type, ptr->moduleID, module_schema_file(node->nodeInfo->type, ptr->type));
		strcat(temp1, temp);
		ptr = ptr->next;
	}

	sprintf(*buffer, "%s=%s", ENV_VNODE_SCHEMA_ARGS, temp1);
	//	sprintf(*buffer, "%s", temp1);
}

/*
 * init_container_file --
 *
 */
static FILE* init_container_file(char *fileName) {

	FILE *fp=NULL;

	if ((fp = fopen(fileName, "w+")) == NULL) {
		log_error("Error opening file: %s Error: %s", fileName,
				  strerror(errno));
		return NULL;
	}

	if (fwrite(CF_HEADER, strlen(CF_HEADER), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				  CF_HEADER, strerror(errno));
		return NULL;
	}

	return fp;
}

/*
 * write_to_container_file --
 *
 */
static int write_to_container_file(char *buffer, char *fileName, FILE *fp) {

	if (buffer == NULL || fp == NULL) return FALSE;

	if (fwrite(buffer, strlen(buffer), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s",
				  SVISOR_FILENAME, buffer, strerror(errno));
		fclose(fp);
		return FALSE;
	}

	return TRUE;
}

/*
 * create_container_file --
 *
 */
static int create_container_file(char *target, Configs *config, Node *node) {

	FILE *fp=NULL;
	char buffer[MAX_BUFFER] = {0};

	if (config == NULL) return FALSE;

	fp = init_container_file(CONTAINER_FILE);
	if (!fp) {
		log_error("Error initalizing container file: %s", CONTAINER_FILE);
		return FALSE;
	}

	sprintf(buffer, CF_FROM, target);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	if (strstr(target, TARGET_ALPINE) != NULL) {
		sprintf(buffer, CF_RUN_APK, UPDATE_PKGS);
	} else {
		sprintf(buffer, CF_RUN_APT, UPDATE_PKGS);
	}
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf (buffer, CF_MKDIR, NODE_LIBS);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf (buffer, CF_MKDIR, SYSFS_DIR);
		if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/sbin", "/sbin");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/conf", "/conf");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/lib", NODE_LIBS);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/mfgdata", "/mfgdata");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/schemas", "/schemas");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/sys", SYSFS_DIR);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/capps", "/capps");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_COPY, "./build/bin", "/bin");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_ADD, "supervisor.conf", "/etc/supervisor.conf");
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_CMD, SUPERVISOR_CMD);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	sprintf(buffer, CF_ENV, LIB_PATH, NODE_LIBS);
	if (!write_to_container_file(buffer, CONTAINER_FILE, fp)) return FALSE;

	fclose(fp);
	return TRUE;
}

/*
 * create_vnode_image --
 *
 */
int create_vnode_image(char *target, Configs *config, Node *node) {

	char runMe[MAX_BUFFER]={0};
	char *buffer=NULL;
	NodeInfo *nodeInfo=NULL;

	if (node == NULL)             return FALSE;
	if (node->nodeInfo   == NULL) return FALSE;
	if (node->nodeConfig == NULL) return FALSE;

	nodeInfo   = node->nodeInfo;

	if (nodeInfo->moduleCount == 0){
		log_error("Node has no module. Node uuid: %s type: %s",
				  nodeInfo->uuid, nodeInfo->type);
		return FALSE;
	}

	buffer = (char *)calloc(1, MAX_BUFFER);
	if (buffer == NULL) {
		log_error("Memory allocation error of size: %lu", MAX_BUFFER);
		return FALSE;
	}

	/* steps are:
	 * 0. clean and build the needed tools
	 * 1. create sysfs
	 * 2. create ContainerFile
	 * 3. create container image
	 * 4. upload the image to registry
	 */

	/* Step:0 clean and build the needed tools */
	sprintf(runMe, "%s init", SCRIPT);
	if (system(runMe) < 0) goto failure;

	/* Step:1 create sys using prepare_env.sh */
	/* 'sysfs type uuid module_metadata' */
	set_schema_args(node, &buffer);

	if (putenv(buffer) != 0) {
		log_error("Unable to set environment variable: %s Error: %s",
				  ENV_VNODE_SCHEMA_ARGS, strerror(errno));
		goto failure;
	}

	sprintf(runMe, "%s sysfs %s %s", SCRIPT, nodeInfo->type, nodeInfo->uuid);
	if (system(runMe) < 0) goto failure;

	/* Step:2 create the container file */
	if (!create_container_file(target, config, node)) {
		log_error("Unable to create container file: %s", CONTAINER_FILE);
		goto failure;
	}

	/* Step:3 run buildah */
	/* build container_file uuid */
	sprintf(runMe, "%s build %s %s", SCRIPT, CONTAINER_FILE, nodeInfo->uuid);
	if (system(runMe) != 0) goto failure;

	/* Step:4 push image to the registry */
	sprintf(runMe, "%s push %s", SCRIPT, nodeInfo->uuid);
	if (system(runMe) != 0) goto failure;

	if (buffer) free(buffer);
	return TRUE;

 failure:
	if (buffer) free(buffer);
	return FALSE;
}

/*
 * purge_container_file --
 *
 */
void purge_container_file(char *fileName) {

	if (remove(fileName) == 0) {
		log_debug("ContainerFile removed: %s", fileName);
	} else {
		log_error("Unable to delete cotainer file: %s", fileName);
	}

	return;
}

/*
 * purge_vnode_image --
 *
 */
void purge_vnode_image(Node *node) {

	char runMe[MAX_BUFFER]  = {0};

	sprintf(runMe, "%s clean %s-%s", SCRIPT, node->nodeInfo->type,
			node->nodeInfo->uuid);
	system(runMe);
}
