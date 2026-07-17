/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>

#include "config_app.h"
#include "log_app.h"

#define SCRIPT         "builder/scripts/make-app.sh"
#define VERSION_SCRIPT "nodes/utils/scripts/generate_version.sh"
#define MAX_BUFFER     4096
#define MAX_LINE       512

static int make_command(char *buffer, size_t size, const char *format, ...) {

    int written;
    va_list args;

    va_start(args, format);
    written = vsnprintf(buffer, size, format, args);
    va_end(args);

    if (written < 0 || (size_t)written >= size) {
        log_error("Command is too long");
        return FALSE;
    }

    return TRUE;
}

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

static int get_app_version(const char *ukamaRoot, char **versionOut) {

    char command[MAX_BUFFER] = {0};
    char line[MAX_LINE]      = {0};
    char *newline            = NULL;
    FILE *fp                 = NULL;
    int status;

    if (ukamaRoot == NULL || versionOut == NULL) {
        return FALSE;
    }

    *versionOut = NULL;

    if (!make_command(command, sizeof(command),
                      "\"%s/%s\" --print",
                      ukamaRoot, VERSION_SCRIPT)) {
        return FALSE;
    }

    fp = popen(command, "r");
    if (fp == NULL) {
        log_error("Unable to run version script");
        return FALSE;
    }

    if (fgets(line, sizeof(line), fp) == NULL) {
        (void)pclose(fp);
        log_error("Unable to read version from script");
        return FALSE;
    }

    status = pclose(fp);
    if (status == -1 || !WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        log_error("Version script failed");
        return FALSE;
    }

    newline = strchr(line, '\n');
    if (newline != NULL) {
        *newline = '\0';
    }

    if (line[0] == '\0' || strcmp(line, "-") == 0) {
        log_error("Invalid app version from script: '%s'", line);
        return FALSE;
    }

    *versionOut = strdup(line);
    if (*versionOut == NULL) {
        log_error("Unable to allocate app version");
        return FALSE;
    }

    return TRUE;
}

int build_app(Config *config) {

    char *ukamaRoot         = NULL;
    char *builtVersion      = NULL;
    char command[MAX_BUFFER] = {0};
    BuildConfig *build      = NULL;

    if (config == NULL || config->build == NULL || config->capp == NULL) {
        return FALSE;
    }

    ukamaRoot = getenv("UKAMA_ROOT");
    if (ukamaRoot == NULL || ukamaRoot[0] == '\0') {
        log_error("UKAMA_ROOT is not set");
        return FALSE;
    }

    build = config->build;

    if (!make_command(command, sizeof(command),
                      "\"%s/%s\" build app \"%s\" \"%s\"",
                      ukamaRoot, SCRIPT, build->source, build->cmd) ||
        !run_command(command)) {
        return FALSE;
    }

    if (!get_app_version(ukamaRoot, &builtVersion)) {
        log_error("Unable to determine app version");
        return FALSE;
    }

    free(config->capp->version);
    config->capp->version = builtVersion;

    if (!make_command(command, sizeof(command),
                      "\"%s/%s\" init \"%s_%s\"",
                      ukamaRoot, SCRIPT,
                      config->capp->name,
                      config->capp->version) ||
        !run_command(command)) {
        return FALSE;
    }

    if (!make_command(command, sizeof(command),
                      "\"%s/%s\" cp \"%s\" \"%s_%s%s\"",
                      ukamaRoot, SCRIPT,
                      build->binFrom,
                      config->capp->name,
                      config->capp->version,
                      build->binTo) ||
        !run_command(command)) {
        return FALSE;
    }

    if (build->mkdir != NULL) {
        if (!make_command(command, sizeof(command),
                          "\"%s/%s\" mkdir \"%s_%s%s\"",
                          ukamaRoot, SCRIPT,
                          config->capp->name,
                          config->capp->version,
                          build->mkdir) ||
            !run_command(command)) {
            return FALSE;
        }
    }

    if (build->from != NULL && build->to != NULL) {
        if (!make_command(command, sizeof(command),
                          "\"%s/%s\" cp \"%s\" \"%s_%s%s\"",
                          ukamaRoot, SCRIPT,
                          build->from,
                          config->capp->name,
                          config->capp->version,
                          build->to) ||
            !run_command(command)) {
            return FALSE;
        }
    }

    if (build->miscFrom != NULL && build->miscTo != NULL) {
        if (!make_command(command, sizeof(command),
                          "\"%s/%s\" cp \"%s\" \"%s_%s%s\"",
                          ukamaRoot, SCRIPT,
                          build->miscFrom,
                          config->capp->name,
                          config->capp->version,
                          build->miscTo) ||
            !run_command(command)) {
            return FALSE;
        }
    }

    if (!build->staticFlag) {
        if (!make_command(command, sizeof(command),
                          "\"%s/%s\" libs \"%s\" \"%s_%s\"",
                          ukamaRoot, SCRIPT,
                          build->binFrom,
                          config->capp->name,
                          config->capp->version) ||
            !run_command(command)) {
            return FALSE;
        }
    }

    return TRUE;
}
