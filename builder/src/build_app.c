/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#include "config_app.h"
#include "log_app.h"

#define SCRIPT            "builder/scripts/make-app.sh"
#define VERSION_SCRIPT    "nodes/utils/scripts/generate_version.sh"
#define LIB_USYS          "nodes/ukamaOS/distro/platform/build/libusys.so"
#define MAX_BUFFER        1024
#define MAX_LINE          512

static int git_mark_safe(const char *repoRoot) {

    char cmd[MAX_BUFFER] = {0};

    if (repoRoot == NULL) {
        return FALSE;
    }

    snprintf(cmd, sizeof(cmd),
             "git config --global --add safe.directory \"%s\" >/dev/null 2>&1",
             repoRoot);

    if (system(cmd) < 0) {
        log_error("Unable to mark repo as git safe.directory: %s", repoRoot);
        return FALSE;
    }

    return TRUE;
}

static int get_app_version(const char *ukamaRoot, char **versionOut) {

    char cmd[MAX_BUFFER]  = {0};
    char line[MAX_LINE]   = {0};
    char *nl              = NULL;
    FILE *fp              = NULL;

    if (ukamaRoot == NULL || versionOut == NULL) {
        return FALSE;
    }

    *versionOut = NULL;

    snprintf(cmd, sizeof(cmd),
             "git config --global --add safe.directory \"%s\" >/dev/null 2>&1; "
             "cd \"%s\" && ./%s --print",
             ukamaRoot, ukamaRoot, VERSION_SCRIPT);

    fp = popen(cmd, "r");
    if (fp == NULL) {
        log_error("Unable to run version script");
        return FALSE;
    }

    if (fgets(line, sizeof(line), fp) == NULL) {
        pclose(fp);
        log_error("Unable to read version from script");
        return FALSE;
    }

    if (pclose(fp) != 0) {
        log_error("Version script failed");
        return FALSE;
    }

    nl = strchr(line, '\n');
    if (nl != NULL) {
        *nl = '\0';
    }

    if (line[0] == '\0' || strcmp(line, "-") == 0) {
        log_error("Invalid app version from script: '%s'", line);
        return FALSE;
    }

    *versionOut = strdup(line);
    if (*versionOut == NULL) {
        return FALSE;
    }

    return TRUE;
}

int build_app(Config *config) {

    char *ukamaRoot        = NULL;
    char *builtVersion     = NULL;
    char runMe[MAX_BUFFER] = {0};
    BuildConfig *build;

    if (config == NULL)        return FALSE;
    if (config->build == NULL) return FALSE;
    if (config->capp == NULL)  return FALSE;

    ukamaRoot = getenv("UKAMA_ROOT");
    if (ukamaRoot == NULL) return FALSE;

    build = config->build;

    if (!git_mark_safe(ukamaRoot)) {
        return FALSE;
    }

    /* Build first so app artifacts and version state are up to date. */
    snprintf(runMe, sizeof(runMe),
             "git config --global --add safe.directory \"%s\" >/dev/null 2>&1; "
             "%s/%s build app %s \"%s\"",
             ukamaRoot, ukamaRoot, SCRIPT, build->source, build->cmd);
    if (system(runMe) < 0) return FALSE;

    /* Source of truth for app package version comes from generate_version.sh. */
    if (!get_app_version(ukamaRoot, &builtVersion)) {
        log_error("Unable to determine app version");
        return FALSE;
    }

    free(config->capp->version);
    config->capp->version = builtVersion;

    /* Initialize package staging area using the real VERSION. */
    snprintf(runMe, sizeof(runMe), "%s/%s init %s_%s",
             ukamaRoot,
             SCRIPT,
             config->capp->name,
             config->capp->version);
    if (system(runMe) < 0) return FALSE;

    snprintf(runMe, sizeof(runMe), "%s/%s cp %s %s_%s%s",
             ukamaRoot, SCRIPT,
             build->binFrom, config->capp->name,
             config->capp->version, build->binTo);
    if (system(runMe) < 0) return FALSE;

    if (build->mkdir) {
        snprintf(runMe, sizeof(runMe), "%s/%s mkdir %s_%s%s",
                 ukamaRoot, SCRIPT,
                 config->capp->name, config->capp->version, build->mkdir);
        if (system(runMe) < 0) return FALSE;
    }

    if (build->from && build->to) {
        snprintf(runMe, sizeof(runMe), "%s/%s cp %s %s_%s%s",
                 ukamaRoot, SCRIPT,
                 build->from, config->capp->name,
                 config->capp->version, build->to);
        if (system(runMe) < 0) return FALSE;
    }

    if (build->miscFrom && build->miscTo) {
        snprintf(runMe, sizeof(runMe), "%s/%s cp %s %s_%s%s",
                 ukamaRoot, SCRIPT,
                 build->miscFrom, config->capp->name,
                 config->capp->version, build->miscTo);
        if (system(runMe) < 0) return FALSE;
    }

    if (!build->staticFlag) {
        snprintf(runMe, sizeof(runMe), "%s/%s libs %s %s_%s",
                 ukamaRoot, SCRIPT,
                 build->binFrom, config->capp->name, config->capp->version);
        if (system(runMe) < 0) return FALSE;
    }

#if 0
    // Currently, we are focusing on using alpine - commented it out.
    if (!build->staticFlag) {
        snprintf(runMe, sizeof(runMe), "%s/%s cp %s/%s %s_%s/lib",
                 ukamaRoot, SCRIPT, ukamaRoot, LIB_USYS,
                 config->capp->name, config->capp->version);
        if (system(runMe) < 0) return FALSE;
    }
#endif

    return TRUE;
}
