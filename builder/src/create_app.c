/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>

#include "config_app.h"
#include "log_app.h"

#define SCRIPT           "./scripts/make-app.sh"
#define MAX_BUFFER       1024
#define DEF_VERSION_FILE "VERSION"

static int create_version_file(Config *config) {

    char path[MAX_BUFFER] = {0};
    FILE *fp              = NULL;

    if (config == NULL ||
        config->capp == NULL ||
        config->capp->name == NULL ||
        config->capp->version == NULL) {
        return FALSE;
    }

    snprintf(path, sizeof(path), "%s_%s/%s",
             config->capp->name,
             config->capp->version,
             DEF_VERSION_FILE);

    fp = fopen(path, "w");
    if (fp == NULL) {
        log_error("Error opening VERSION file %s: %s", path, strerror(errno));
        return FALSE;
    }

    if (fprintf(fp, "%s\n", config->capp->version) < 0) {
        log_error("Error writing VERSION file %s: %s", path, strerror(errno));
        fclose(fp);
        remove(path);
        return FALSE;
    }

    fclose(fp);
    return TRUE;
}

int create_app(Config *config) {

    char runMe[MAX_BUFFER] = {0};

    if (config == NULL || config->capp == NULL) {
        return FALSE;
    }

    if (!create_version_file(config)) {
        log_error("Error creating VERSION file");
        return FALSE;
    }

    snprintf(runMe, sizeof(runMe), "%s pack %s %s_%s.tar.gz %s_%s %d",
             SCRIPT,
             getenv("UKAMA_ROOT"),
             config->capp->name, config->capp->version,
             config->capp->name, config->capp->version, TRUE);
    if (system(runMe) < 0) {
        log_error("Error packing the capp to %s_%s",
                  config->capp->name, config->capp->version);
        return FALSE;
    }

    snprintf(runMe, sizeof(runMe), "%s clean %s_%s",
             SCRIPT,
             config->capp->name, config->capp->version);
    if (system(runMe) < 0) return FALSE;

    return TRUE;
}
