/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>

#include "config_app.h"
#include "log_app.h"

#define SCRIPT           "builder/scripts/make-app.sh"
#define MAX_BUFFER       4096
#define DEF_VERSION_FILE "VERSION"

static int run_command(const char *command) {

    int status;

    if (command == NULL || command[0] == '\0') {
        return FALSE;
    }

    status = system(command);
    if (status == -1) {
        log_error("Unable to execute command: %s", command);
        return FALSE;
    }

    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        if (WIFEXITED(status)) {
            log_error("Command exited with status %d: %s",
                      WEXITSTATUS(status), command);
        } else {
            log_error("Command terminated abnormally: %s", command);
        }
        return FALSE;
    }

    return TRUE;
}

static int create_version_file(Config *config) {

    char path[MAX_BUFFER] = {0};
    FILE *fp              = NULL;
    int written;

    if (config == NULL || config->capp == NULL ||
        config->capp->name == NULL || config->capp->version == NULL) {
        return FALSE;
    }

    written = snprintf(path, sizeof(path), "%s_%s/%s",
                       config->capp->name,
                       config->capp->version,
                       DEF_VERSION_FILE);
    if (written < 0 || (size_t)written >= sizeof(path)) {
        log_error("VERSION path is too long");
        return FALSE;
    }

    fp = fopen(path, "w");
    if (fp == NULL) {
        log_error("Error opening VERSION file %s: %s",
                  path, strerror(errno));
        return FALSE;
    }

    if (fprintf(fp, "%s\n", config->capp->version) < 0) {
        log_error("Error writing VERSION file %s: %s",
                  path, strerror(errno));
        fclose(fp);
        remove(path);
        return FALSE;
    }

    if (fclose(fp) != 0) {
        log_error("Error closing VERSION file %s: %s",
                  path, strerror(errno));
        remove(path);
        return FALSE;
    }

    return TRUE;
}

int create_app(Config *config) {

    char command[MAX_BUFFER] = {0};
    char *ukamaRoot          = NULL;
    int written;

    if (config == NULL || config->capp == NULL) {
        return FALSE;
    }

    ukamaRoot = getenv("UKAMA_ROOT");
    if (ukamaRoot == NULL || ukamaRoot[0] == '\0') {
        log_error("UKAMA_ROOT is not set");
        return FALSE;
    }

    if (!create_version_file(config)) {
        log_error("Error creating VERSION file");
        return FALSE;
    }

    written = snprintf(command, sizeof(command),
                       "\"%s/%s\" pack \"%s\" "
                       "\"%s_%s.tar.gz\" \"%s_%s\" 0",
                       ukamaRoot, SCRIPT, ukamaRoot,
                       config->capp->name,
                       config->capp->version,
                       config->capp->name,
                       config->capp->version);
    if (written < 0 || (size_t)written >= sizeof(command) ||
        !run_command(command)) {
        log_error("Error packing the capp %s_%s",
                  config->capp->name,
                  config->capp->version);
        return FALSE;
    }

    written = snprintf(command, sizeof(command),
                       "\"%s/%s\" clean \"%s_%s\"",
                       ukamaRoot, SCRIPT,
                       config->capp->name,
                       config->capp->version);
    if (written < 0 || (size_t)written >= sizeof(command) ||
        !run_command(command)) {
        return FALSE;
    }

    return TRUE;
}
