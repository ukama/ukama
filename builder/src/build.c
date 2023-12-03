/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>

#include "builder.h"
#include "config.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_string.h"
#include "usys_types.h"

#define SCRIPT      "./scripts/build.sh"
#define MAX_BUFFER  1024
#define DELIMINATOR ","
#define UKAMA_AUTH  "ukama-auth"

static bool build_system(char *name, char *path) {

	char runMe[MAX_BUFFER] = {0};

	if (name == NULL || path == NULL) return USYS_FALSE;

	sprintf(runMe, "%s system %s", SCRIPT, path);
	if (system(runMe) < 0) return USYS_FALSE;

    return USYS_TRUE;
}

bool build_all_systems(char *systemsList, char *ukamaRepo, char *authRepo) {

    char list[MAX_BUFFER] = {0};
    char systemPath[MAX_BUFFER] = {0};
    char *systemName = NULL;

    strncpy(list, systemsList, strlen(systemsList));
    list[strlen(systemsList) - 1] = '\0';

    systemName = strtok(list, DELIMINATOR);
    while (systemName != NULL) {

        if (strcasecmp(systemName, UKAMA_AUTH) == 0) {
            if (!build_system(systemName, authRepo)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, authRepo);
                return USYS_FALSE;
            }
        } else {
            sprintf(systemPath, "%s/systems/%s/", ukamaRepo, systemName);
            if (!build_system(systemName, systemPath)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, systemPath);
                return USYS_FALSE;
            }
        }

        systemName = strtok(NULL, DELIMINATOR);
    }

    return USYS_TRUE;
}
