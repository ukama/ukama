/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>

#include "persistence.h"

int persistence_save(const EmuModel *model, const char *path) {
    FILE *fp = NULL;

    if (path == NULL || path[0] == '\0') {
        return STATUS_OK;
    }

    fp = fopen(path, "w");
    if (fp == NULL) {
        return STATUS_NOK;
    }

    fprintf(fp,
            "{\n"
            "  \"scenario\": \"%s\",\n"
            "  \"softwareVersion\": \"%s\"\n"
            "}\n",
            model->activeScenario, model->info.softwareVersion);
    fclose(fp);

    return STATUS_OK;
}

int persistence_load(EmuModel *model, const char *path) {
    (void)model;
    (void)path;
    return STATUS_OK;
}
