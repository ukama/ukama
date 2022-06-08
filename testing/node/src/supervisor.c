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

	if (fwrite(SVISOR_HEADER, strlen(SVISOR_HEADER), 1, fp) <=0 ) {
		log_error("Error writing to %s. Str: %s. Error: %s", fileName,
				  SVISOR_HEADER, strerror(errno));
		return NULL;
	}

	return fp;
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
	FILE *fp=NULL;

	if (configs == NULL) return FALSE;

	fp = init_supervisor_config(SVISOR_FILENAME);
	if (!fp) {
		log_error("Error initializing supervisor config file: %s",
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

		/* formulate the command */
		/* wait for progs to finish exec */
		if (capp->dependsOn) {
			sprintf(cmd, "%s %s;", SVISOR_WAITFOR_SH, capp->dependsOn);
		} else {
			sprintf(cmd, " ");
		}
		/* sleep for sometime */
		if (capp->waitFor) {
			sprintf(cmd, "%s sleep %s;", cmd, capp->waitFor);
		} else {
			sprintf(cmd, "%s ", cmd);
		}

		if (capp->args) {
			sprintf(buffer, SVISOR_COMMAND_WITH_ARGS, buffer, cmd,
					SVISOR_RUNME_SH, capp->name,
					capp->path, capp->bin, capp->args);
		} else {
			sprintf(buffer, SVISOR_COMMAND, buffer, cmd,
					SVISOR_RUNME_SH, capp->name, capp->path, capp->bin);
		}

		sprintf(buffer, SVISOR_AUTOSTART, buffer,
				(capp->autostart ? "true": "false"));
		sprintf(buffer, SVISOR_AUTORESTART, buffer,
				(capp->autorestart ? "true": "false"));
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
