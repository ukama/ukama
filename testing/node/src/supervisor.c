/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* supervisor.d related stuff for virual node */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

#include "log.h"
#include "config.h"
#include "supervisor.h"

#define SVISOR_FILENAME      "supervisor.conf"
#define SVISOR_CONFIg_HEADER "[supervisord] \n nodaemon=true \n\n"
#define SVISOR_

static void append_service_to_group(char* group, char* name, char* version) {
	if (name) {
		if(strlen(group)>0) {
			strcat(group, ",");
		}
		strcat(group, name);
		strcat(group, "_");
		strcat(group, version);
	}
}

/*
 * init_supervisor_config --
 *
 */
static FILE* init_supervisor_config(char *fileName) {

	FILE *fp=NULL;

	if ((fp = fopen(fileName, "w+")) == NULL) {
		log_error("Error opening file: %s Error: %s", fileName,
				  strerror(errno));
		return NULL;
	}

	/* Header */
	if (fwrite(SVISOR_HEADER, strlen(SVISOR_HEADER), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				  SVISOR_HEADER, strerror(errno));
		return NULL;
	}

	/* Supervisord */
	if (fwrite(SVISOR_SVISORD, strlen(SVISOR_SVISORD), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				SVISOR_SVISORD, strerror(errno));
		return NULL;
	}

	/* Supervisord_rpcinterface */
	if (fwrite(SVISOR_RPCINTERFACE, strlen(SVISOR_RPCINTERFACE), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				SVISOR_RPCINTERFACE, strerror(errno));
		return NULL;
	}

	/* Supervisorctl */
	if (fwrite(SVISOR_SVISOR_CTL, strlen(SVISOR_SVISOR_CTL), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				SVISOR_SVISOR_CTL, strerror(errno));
		return NULL;
	}

	/* Includes */
	if (fwrite(SVISOR_INCLUDE, strlen(SVISOR_INCLUDE), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				SVISOR_INCLUDE, strerror(errno));
		return NULL;
	}

	/* Kickstart */
	if (fwrite(SVISOR_KICKSTART, strlen(SVISOR_KICKSTART), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				SVISOR_KICKSTART, strerror(errno));
		return NULL;
	}

	return fp;
}

/*
 * Update Groups
 *
 * 	on-boot: Services which will be started by supervisorctl
 * 	 		  before bootstrap.
 * 	system-services: Services started by supervisorctl once
 * 			bootstrap is completed and meshd is started.
 *
 * Note: If Services has dependency that should be handled by events not
 * 		 by init.
 */
int create_supervisor_groups(FILE* fp ,Configs *configs, char* onBootGroup,
		char* sysGroup) {
	Configs    *ptr=NULL;
	CappConfig *capp=NULL;
	char buffer[SVISOR_MAX_SIZE] = {0};

	for (ptr = configs; ptr; ptr=ptr->next) {
	    if (!ptr->valid)         continue;
	    if (!ptr->config)        continue;
	    if (!ptr->config->capp)  continue;
	    if (!ptr->config->capp->group)  continue;

	    capp = ptr->config->capp;

	    if (strcmp(capp->group, SVISOR_GROUP_ON_BOOT)==0) {
	        append_service_to_group(onBootGroup, capp->name, capp->version);
	    } else if (strcmp(capp->group, SVISOR_GROUP_SYS_SVC)==0) {
	        append_service_to_group(sysGroup, capp->name, capp->version);
	    } else {
	        continue;
	    }
	}

	if (strlen(onBootGroup) !=0 ) {
	    /* On-boot group */
	    sprintf(buffer, SVISOR_GROUP_ONBOOT, onBootGroup);
	}

	if (strlen(sysGroup) !=0 ) {
	    /* sys-service group */
	    sprintf(buffer, SVISOR_GROUP_ONBOOT, sysGroup);
	}

	if (fwrite(buffer, strlen(buffer), 1, fp) <=0 ) {
	    log_error("Error writing Group config to %s. Str: %s. Error: %s", SVISOR_FILENAME,
	            buffer, strerror(errno));
	    return FALSE;
	}

	return TRUE;
}

/*
 * create_supervisor_config --
 *
 */
int create_supervisor_config(Configs *configs) {

	Configs    *ptr=NULL;
	CappConfig *capp=NULL;
	char buffer[SVISOR_MAX_SIZE] = {0};
	char cmd[SVISOR_MAX_SIZE] = {0};
	char onBootGroup[SVISOR_GROUP_LIST_MAX_SIZE] = {0};
	char sysGroup[SVISOR_GROUP_LIST_MAX_SIZE]={0};

	FILE *fp=NULL;

	if (configs == NULL) return FALSE;

	fp = init_supervisor_config(SVISOR_FILENAME);
	if (!fp) {
		log_error("Error initializing supervisor config file: %s",
				  SVISOR_FILENAME);
		return FALSE;
	}

	/* Do the service grouping */
	if (!create_supervisor_groups(fp , configs, onBootGroup, sysGroup)) {
	    log_error("Error grouping supervisor config file: %s",
	            SVISOR_FILENAME);
	    return FALSE;
	}

	for (ptr = configs; ptr; ptr=ptr->next) {

	    if (!ptr->valid)         continue;
	    if (!ptr->config)        continue;
	    if (!ptr->config->capp)  continue;
	    if (!ptr->config->build) continue;

	    capp = ptr->config->capp;

	    memset(buffer, 0, SVISOR_MAX_SIZE);

	    sprintf(buffer, SVISOR_PROGRAM, capp->name, capp->version);

	    if (capp->args) {
	        sprintf(buffer, SVISOR_COMMAND_WITH_ARGS, buffer,
	                capp->path, capp->bin, capp->args);
	    } else {
	        sprintf(buffer, SVISOR_COMMAND, buffer,
	                capp->path, capp->bin);
	    }

	    sprintf(buffer, SVISOR_AUTOSTART, buffer,
	            (capp->autostart ? "true": "false"));
	    sprintf(buffer, SVISOR_AUTORESTART, buffer,
	            (capp->autorestart ? "true": "false"));
	    sprintf(buffer, SVISOR_STARTRETRIES, buffer, capp->startretries);
	    sprintf(buffer, SVISOR_STDERR_LOGFILE, buffer, SVISOR_DEFAULT_STDERR);
	    sprintf(buffer, SVISOR_STDOUT_LOGFILE, buffer, SVISOR_DEFAULT_STDOUT);
	    sprintf(buffer, SVISOR_STDERR_LOGFILE_MAX_BYTES, buffer,
	            SVISOR_DEFAULT_STDERR_MAXBYTES);
	    sprintf(buffer, SVISOR_STDOUT_LOGFILE_MAX_BYTES, buffer,
	            SVISOR_DEFAULT_STDOUT_MAXBYTES);
	    sprintf(buffer, "%s \n", buffer);

	    if (fwrite(buffer, strlen(buffer), 1, fp) <=0 ) {
	        log_error("Error writing to %s. Str: %s. Error: %s",
	                SVISOR_FILENAME, buffer, strerror(errno));
	        fclose(fp);
	        return FALSE;
	    }
	}

	fclose(fp);
	return TRUE;
}

/*
 * purge_supervisor_config --
 *
 */
void purge_supervisor_config(char *fileName) {

	if (remove(fileName) == 0) {
		log_debug("supervisor config file removed: %s", fileName);
	} else {
		log_error("Unable to delete supervisor config file: %s", fileName);
	}

	return;
}
