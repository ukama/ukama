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
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

#define SHUTDOWN_SCRIPT "./scripts/shutdown.sh"

bool shutdown_nodes(int count, char **nodesIDList) {

    int i;
    char runMe[MAX_BUFFER] = {0};

    for (i=0; i<count; i++) {
        sprintf(runMe, "%s node %s", SHUTDOWN_SCRIPT, nodesIDList[i]);
        if (system(runMe) < 0) {
            usys_log_error("Unable to shutdown node: %s", nodesIDList[i]);
            continue;
        }
    }

    return USYS_TRUE;
}

bool shutdown_all_systems(char *systems,
                          char *ukamaRepo,
                          char *authRepo) {

    char runMe[2*MAX_BUFFER+1] = {0};
    char list[MAX_BUFFER] = {0};
    char systemPath[MAX_BUFFER] = {0};
    char *systemName = NULL;

    strncpy(list, systems, strlen(systems));
    list[strlen(systems) - 1] = '\0';

    systemName = strtok(list, DELIMINATOR);
    while (systemName != NULL) {

        if (strcasecmp(systemName, UKAMA_AUTH) == 0) {
            sprintf(runMe, "%s system %s %s",
                    SHUTDOWN_SCRIPT,
                    authRepo,
                    UKAMA_AUTH);
        } else {
            sprintf(systemPath, "%s/systems/%s/", ukamaRepo, systemName);
            sprintf(runMe, "%s system %s %s",
                    SHUTDOWN_SCRIPT,
                    systemPath,
                    systemName);
        }

        if (system(runMe) < 0) {
            usys_log_error("Unable to execute: %s", runMe);
        }

        systemName = strtok(NULL, DELIMINATOR);
    }

    return USYS_TRUE;
}
