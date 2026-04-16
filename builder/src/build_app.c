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

#define SCRIPT       "builder/scripts/make-app.sh"
#define LIB_USYS     "nodes/ukamaOS/distro/platform/build/libusys.so"
#define MAX_BUFFER   1024
#define MAX_LINE     512

static int read_version_from_header(const char *srcDir, char **versionOut) {

    char path[MAX_BUFFER] = {0};
    char line[MAX_LINE]   = {0};
    char *start           = NULL;
    char *end             = NULL;
    char *version         = NULL;
    FILE *fp              = NULL;

    if (versionOut == NULL || srcDir == NULL) {
        return FALSE;
    }

    *versionOut = NULL;

    snprintf(path, sizeof(path), "%s/version.h", srcDir);

    fp = fopen(path, "r");
    if (fp == NULL) {
        log_error("Unable to open version header: %s", path);
        return FALSE;
    }

    while (fgets(line, sizeof(line), fp) != NULL) {
        if (strstr(line, "#define VERSION") == NULL) {
            continue;
        }

        start = strchr(line, '"');
        if (start == NULL) {
            continue;
        }

        start++;

        end = strchr(start, '"');
        if (end == NULL || end <= start) {
            continue;
        }

        *end = '\0';

        version = strdup(start);
        if (version == NULL) {
            fclose(fp);
            return FALSE;
        }

        *versionOut = version;
        fclose(fp);
        return TRUE;
    }

    fclose(fp);
    log_error("VERSION not found in: %s", path);
    return FALSE;
}

int build_app(Config *config) {

    char *ukamaRoot    = NULL;
    char *builtVersion = NULL;
    char runMe[MAX_BUFFER] = {0};
    BuildConfig *build;

    if (config == NULL)        return FALSE;
    if (config->build == NULL) return FALSE;
    if (config->capp == NULL)  return FALSE;

    ukamaRoot = getenv("UKAMA_ROOT");
    if (ukamaRoot == NULL) return FALSE;

    build = config->build;

    /* Build first so version.h is generated/updated by the app build. */
    snprintf(runMe, sizeof(runMe), "%s/%s build app %s \"%s\"",
             ukamaRoot, SCRIPT, build->source, build->cmd);
    if (system(runMe) < 0) return FALSE;

    /* Source of truth for app package version is version.h. */
    if (!read_version_from_header(build->source, &builtVersion)) {
        log_error("Unable to read app version from version.h");
        return FALSE;
    }

    free(config->capp->version);
    config->capp->version = builtVersion;

    /* initialize package staging area using the real VERSION. */
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
