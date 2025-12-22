/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

#include "log.h"
#include "config.h"
#include "supervisor.h"

#define SVISOR_FILENAME "supervisor.conf"

/* small helper */
static int streq(const char *a, const char *b) {
    return (a && b && strcmp(a, b) == 0);
}

static void append_service_to_group(char* group, char* name, char* version) {
    if (name && version) {
        if (strlen(group) > 0) {
            strcat(group, ",");
        }
        strcat(group, name);
        strcat(group, "_");
        strcat(group, version);
    }
}

static FILE* init_supervisor_config(char *fileName) {
    FILE *fp = NULL;

    if ((fp = fopen(fileName, "w+")) == NULL) {
        log_error("Error opening file: %s Error: %s", fileName, strerror(errno));
        return NULL;
    }

    if (fwrite(SVISOR_HEADER, strlen(SVISOR_HEADER), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_HEADER, strerror(errno));
        fclose(fp);
        return NULL;
    }

    if (fwrite(SVISOR_SVISORD, strlen(SVISOR_SVISORD), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_SVISORD, strerror(errno));
        fclose(fp);
        return NULL;
    }

    if (fwrite(SVISOR_RPCINTERFACE, strlen(SVISOR_RPCINTERFACE), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_RPCINTERFACE, strerror(errno));
        fclose(fp);
        return NULL;
    }

    if (fwrite(SVISOR_SVISOR_CTL, strlen(SVISOR_SVISOR_CTL), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_SVISOR_CTL, strerror(errno));
        fclose(fp);
        return NULL;
    }

    if (fwrite(SVISOR_INCLUDE, strlen(SVISOR_INCLUDE), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_INCLUDE, strerror(errno));
        fclose(fp);
        return NULL;
    }

    if (fwrite(SVISOR_KICKSTART, strlen(SVISOR_KICKSTART), 1, fp) <= 0) {
        log_error("Error writing to %s. Str: %s. Error: %s",
                  fileName, SVISOR_KICKSTART, strerror(errno));
        fclose(fp);
        return NULL;
    }

    return fp;
}

int create_supervisor_groups(FILE* fp, Configs *configs, char* onBootGroup, char* sysGroup) {
    Configs    *ptr = NULL;
    CappConfig *capp = NULL;
    char buffer[SVISOR_MAX_SIZE];

    if (!fp || !configs) return FALSE;

    for (ptr = configs; ptr; ptr = ptr->next) {
        if (!ptr->valid) continue;
        if (!ptr->config || !ptr->config->capp) continue;
        if (!ptr->config->capp->group) continue;

        capp = ptr->config->capp;

        if (streq(capp->group, SVISOR_GROUP_ON_BOOT)) {
            append_service_to_group(onBootGroup, capp->name, capp->version);
        } else if (streq(capp->group, SVISOR_GROUP_SYS_SVC)) {
            append_service_to_group(sysGroup, capp->name, capp->version);
        }
    }

    if (strlen(onBootGroup) != 0) {
        memset(buffer, 0, sizeof(buffer));
        snprintf(buffer, sizeof(buffer), SVISOR_GROUP_ONBOOT, onBootGroup);
        if (fwrite(buffer, strlen(buffer), 1, fp) <= 0) {
            log_error("Error writing On-boot Group to %s. Str: %s. Error: %s",
                      SVISOR_FILENAME, buffer, strerror(errno));
            return FALSE;
        }
    }

    if (strlen(sysGroup) != 0) {
        memset(buffer, 0, sizeof(buffer));
        snprintf(buffer, sizeof(buffer), SVISOR_GROUP_SYSSVC, sysGroup);
        if (fwrite(buffer, strlen(buffer), 1, fp) <= 0) {
            log_error("Error writing Sys-service Group to %s. Str: %s. Error: %s",
                      SVISOR_FILENAME, buffer, strerror(errno));
            return FALSE;
        }
    }

    return TRUE;
}

static int is_bootstrap_program(const CappConfig *capp) {
    if (!capp || !capp->name) return FALSE;
    return (strcmp(capp->name, "bootstrap") == 0);
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
int create_supervisor_config(Configs *configs) {

    Configs    *ptr = NULL;
    CappConfig *capp = NULL;
    char buffer[SVISOR_MAX_SIZE];
    char onBootGroup[SVISOR_GROUP_LIST_MAX_SIZE] = {0};
    char sysGroup[SVISOR_GROUP_LIST_MAX_SIZE] = {0};

    FILE *fp = NULL;

    if (configs == NULL) return FALSE;

    fp = init_supervisor_config(SVISOR_FILENAME);
    if (!fp) {
        log_error("Error initializing supervisor config file: %s", SVISOR_FILENAME);
        return FALSE;
    }

    if (!create_supervisor_groups(fp, configs, onBootGroup, sysGroup)) {
        log_error("Error grouping supervisor config file: %s", SVISOR_FILENAME);
        fclose(fp);
        return FALSE;
    }

    for (ptr = configs; ptr; ptr = ptr->next) {

        if (!ptr->valid) continue;
        if (!ptr->config || !ptr->config->capp || !ptr->config->build) continue;

        capp = ptr->config->capp;

        memset(buffer, 0, sizeof(buffer));

        /* program header */
        snprintf(buffer, sizeof(buffer), SVISOR_PROGRAM, capp->name, capp->version);

        /* command */
        if (capp->args) {
            snprintf(buffer, sizeof(buffer), SVISOR_COMMAND_WITH_ARGS, buffer,
                     capp->path, capp->bin, capp->args);
        } else {
            snprintf(buffer, sizeof(buffer), SVISOR_COMMAND, buffer,
                     capp->path, capp->bin);
        }

        snprintf(buffer, sizeof(buffer), SVISOR_AUTOSTART, buffer,
                 (capp->autostart ? "true" : "false"));

        if (is_bootstrap_program(capp)) {
            snprintf(buffer, sizeof(buffer), SVISOR_AUTORESTART, buffer, "false");
            snprintf(buffer, sizeof(buffer), SVISOR_STARTRETRIES, buffer, 0);
            snprintf(buffer, sizeof(buffer), SVISOR_STARTSECS, buffer, 0);
            snprintf(buffer, sizeof(buffer), SVISOR_EXITCODES, buffer);
        } else {
            snprintf(buffer, sizeof(buffer), SVISOR_AUTORESTART, buffer,
                     (capp->autorestart ? "true" : "false"));
            snprintf(buffer, sizeof(buffer), SVISOR_STARTRETRIES, buffer, capp->startretries);

            snprintf(buffer, sizeof(buffer), SVISOR_STARTSECS, buffer, 2);
        }

        snprintf(buffer, sizeof(buffer), SVISOR_STDERR_LOGFILE, buffer);
        snprintf(buffer, sizeof(buffer), SVISOR_STDOUT_LOGFILE, buffer);

        strncat(buffer, "\n", sizeof(buffer) - strlen(buffer) - 1);

        if (fwrite(buffer, strlen(buffer), 1, fp) <= 0) {
            log_error("Error writing to %s. Str: %s. Error: %s",
                      SVISOR_FILENAME, buffer, strerror(errno));
            fclose(fp);
            return FALSE;
        }
    }

    fclose(fp);
    return TRUE;
}

void purge_supervisor_config(char *fileName) {

    if (remove(fileName) == 0) {
        log_debug("supervisor config file removed: %s", fileName);
    } else {
        log_error("Unable to delete supervisor config file: %s", fileName);
    }
}
