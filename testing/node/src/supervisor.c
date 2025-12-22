/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdarg.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

#include "log.h"
#include "config.h"
#include "supervisor.h"

#define SVISOR_FILENAME "supervisor.conf"

static int streq(const char *a, const char *b) {
    return (a && b && strcmp(a, b) == 0);
}

/* Append formatted string to dst safely. Returns 0 on success, -1 on overflow/error. */
static int appendf(char *dst, size_t dstsz, const char *fmt, ...) {
    size_t len = strnlen(dst, dstsz);
    if (len >= dstsz) return -1;

    va_list ap;
    va_start(ap, fmt);
    int n = vsnprintf(dst + len, dstsz - len, fmt, ap);
    va_end(ap);

    if (n < 0) return -1;
    if ((size_t)n >= (dstsz - len)) return -1; /* truncated */
    return 0;
}

static int is_bootstrap_program(const CappConfig *capp) {
    return (capp && capp->name && strcmp(capp->name, "bootstrap") == 0);
}

/* group string is a comma-separated list of program names */
static void append_service_to_group(char *group,
                                    size_t groupsz,
                                    const char *name,
                                    const char *version) {
    if (!group || groupsz == 0 || !name || !version) return;

    /* Compose "name_version" */
    char entry[256];
    int n = snprintf(entry, sizeof(entry), "%s_%s", name, version);
    if (n < 0 || (size_t)n >= sizeof(entry)) return;

    size_t glen = strnlen(group, groupsz);

    /* Add comma if needed */
    if (glen > 0) {
        if (glen + 1 >= groupsz) return;
        group[glen] = ',';
        group[glen + 1] = '\0';
        glen++;
    }

    size_t elen = strnlen(entry, sizeof(entry));
    if (glen + elen >= groupsz) return;

    strncat(group, entry, groupsz - glen - 1);
}

static FILE* init_supervisor_config(const char *fileName) {
    FILE *fp = NULL;

    if (!fileName) return NULL;

    fp = fopen(fileName, "w+");
    if (!fp) {
        log_error("Error opening file: %s Error: %s", fileName, strerror(errno));
        return NULL;
    }

    /* Header blocks */
    if (fwrite(SVISOR_HEADER, strlen(SVISOR_HEADER), 1, fp) <= 0 ||
        fwrite(SVISOR_SVISORD, strlen(SVISOR_SVISORD), 1, fp) <= 0 ||
        fwrite(SVISOR_RPCINTERFACE, strlen(SVISOR_RPCINTERFACE), 1, fp) <= 0 ||
        fwrite(SVISOR_SVISOR_CTL, strlen(SVISOR_SVISOR_CTL), 1, fp) <= 0 ||
        fwrite(SVISOR_KICKSTART, strlen(SVISOR_KICKSTART), 1, fp) <= 0) {

        log_error("Error writing supervisor header to %s. Error: %s", fileName, strerror(errno));
        fclose(fp);
        return NULL;
    }

    return fp;
}

/*
 * on-boot: started by kickstart before bootstrap
 * sys-service: started after meshd is running
 */
static int create_supervisor_groups(FILE *fp,
                                    Configs *configs,
                                    char *onBootGroup,
                                    size_t onBootSz,
                                    char *sysGroup,
                                    size_t sysSz) {
    Configs *ptr = NULL;
    CappConfig *capp = NULL;
    char block[SVISOR_MAX_SIZE];

    if (!fp || !configs || !onBootGroup || !sysGroup) return FALSE;

    onBootGroup[0] = '\0';
    sysGroup[0] = '\0';

    for (ptr = configs; ptr; ptr = ptr->next) {
        if (!ptr->valid) continue;
        if (!ptr->config || !ptr->config->capp) continue;
        if (!ptr->config->capp->group) continue;

        capp = ptr->config->capp;

        if (streq(capp->group, SVISOR_GROUP_ON_BOOT)) {
            append_service_to_group(onBootGroup, onBootSz, capp->name, capp->version);
        } else if (streq(capp->group, SVISOR_GROUP_SYS_SVC)) {
            append_service_to_group(sysGroup, sysSz, capp->name, capp->version);
        }
    }

    if (strlen(onBootGroup) > 0) {
        memset(block, 0, sizeof(block));
        if (snprintf(block, sizeof(block), SVISOR_GROUP_ONBOOT, onBootGroup) < 0) return FALSE;
        if (fwrite(block, strlen(block), 1, fp) <= 0) {
            log_error("Error writing on-boot group to %s. Error: %s",
                      SVISOR_FILENAME, strerror(errno));
            return FALSE;
        }
    }

    if (strlen(sysGroup) > 0) {
        memset(block, 0, sizeof(block));
        if (snprintf(block, sizeof(block), SVISOR_GROUP_SYSSVC, sysGroup) < 0) return FALSE;
        if (fwrite(block, strlen(block), 1, fp) <= 0) {
            log_error("Error writing sys-service group to %s. Error: %s",
                      SVISOR_FILENAME,
                      strerror(errno));
            return FALSE;
        }
    }

    return TRUE;
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
    Configs *ptr = NULL;
    CappConfig *capp = NULL;

    char buffer[SVISOR_MAX_SIZE];
    char onBootGroup[SVISOR_GROUP_LIST_MAX_SIZE];
    char sysGroup[SVISOR_GROUP_LIST_MAX_SIZE];

    FILE *fp = NULL;

    if (!configs) return FALSE;

    fp = init_supervisor_config(SVISOR_FILENAME);
    if (!fp) {
        log_error("Error initializing supervisor config file: %s", SVISOR_FILENAME);
        return FALSE;
    }

    if (!create_supervisor_groups(fp, configs,
                                  onBootGroup, sizeof(onBootGroup),
                                  sysGroup, sizeof(sysGroup))) {
        log_error("Error creating supervisor groups in: %s", SVISOR_FILENAME);
        fclose(fp);
        return FALSE;
    }

    for (ptr = configs; ptr; ptr = ptr->next) {
        if (!ptr->valid) continue;
        if (!ptr->config || !ptr->config->capp || !ptr->config->build) continue;

        capp = ptr->config->capp;

        /* minimal sanity checks */
        if (!capp->name || !capp->version || !capp->path || !capp->bin) {
            log_error("Skipping invalid capp (missing fields): name=%s version=%s path=%s bin=%s",
                      capp->name ? capp->name : "(null)",
                      capp->version ? capp->version : "(null)",
                      capp->path ? capp->path : "(null)",
                      capp->bin ? capp->bin : "(null)");
            continue;
        }

        memset(buffer, 0, sizeof(buffer));

        /* [program:name_version] */
        if (appendf(buffer, sizeof(buffer), SVISOR_PROGRAM, capp->name, capp->version) < 0) {
            log_error("Supervisor config overflow building program header for %s_%s",
                      capp->name, capp->version);
            fclose(fp);
            return FALSE;
        }

        /* command=... */
        if (capp->args && strlen(capp->args) > 0) {
            if (appendf(buffer, sizeof(buffer), "command=%s/%s %s\n",
                        capp->path, capp->bin, capp->args) < 0) {
                log_error("Supervisor config overflow building command (with args) for %s_%s",
                          capp->name, capp->version);
                fclose(fp);
                return FALSE;
            }
        } else {
            if (appendf(buffer, sizeof(buffer), "command=%s/%s\n", capp->path, capp->bin) < 0) {
                log_error("Supervisor config overflow building command for %s_%s",
                          capp->name, capp->version);
                fclose(fp);
                return FALSE;
            }
        }

        /* autostart */
        if (appendf(buffer, sizeof(buffer), "autostart=%s\n",
                    capp->autostart ? "true" : "false") < 0) {
            log_error("Supervisor config overflow building autostart for %s_%s",
                      capp->name, capp->version);
            fclose(fp);
            return FALSE;
        }

        /*
         * bootstrap is one-shot:
         * - do not restart
         * - must exit 0
         * - startretries=0, startsecs=0
         */
        if (is_bootstrap_program(capp)) {
            if (appendf(buffer, sizeof(buffer), "autorestart=false\n") < 0 ||
                appendf(buffer, sizeof(buffer), "startretries=0\n") < 0 ||
                appendf(buffer, sizeof(buffer), "startsecs=0\n") < 0 ||
                appendf(buffer, sizeof(buffer), "exitcodes=0\n") < 0) {
                log_error("Supervisor config overflow building bootstrap policy for %s_%s",
                          capp->name, capp->version);
                fclose(fp);
                return FALSE;
            }
        } else {
            /* daemons follow config, but use startsecs=2 for stability */
            if (appendf(buffer, sizeof(buffer), "autorestart=%s\n",
                        capp->autorestart ? "true" : "false") < 0 ||
                appendf(buffer, sizeof(buffer), "startretries=%d\n", capp->startretries) < 0 ||
                appendf(buffer, sizeof(buffer), "startsecs=2\n") < 0) {
                log_error("Supervisor config overflow building daemon policy for %s_%s",
                          capp->name, capp->version);
                fclose(fp);
                return FALSE;
            }
        }

        /* container-friendly logs */
        if (appendf(buffer, sizeof(buffer), "stdout_logfile=/dev/stdout\n") < 0 ||
            appendf(buffer, sizeof(buffer), "stderr_logfile=/dev/stderr\n\n") < 0) {
            log_error("Supervisor config overflow building logs for %s_%s",
                      capp->name, capp->version);
            fclose(fp);
            return FALSE;
        }

        if (fwrite(buffer, strlen(buffer), 1, fp) <= 0) {
            log_error("Error writing to %s. Error: %s", SVISOR_FILENAME, strerror(errno));
            fclose(fp);
            return FALSE;
        }
    }

    fclose(fp);
    return TRUE;
}

void purge_supervisor_config(char *fileName) {
    if (!fileName) return;

    if (remove(fileName) == 0) {
        log_debug("supervisor config file removed: %s", fileName);
    } else {
        log_error("Unable to delete supervisor config file: %s", fileName);
    }
}
