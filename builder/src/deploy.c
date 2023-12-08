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

#define SCRIPT "./scripts/deploy-system.sh"

void toLowerCase(char *str) {
    while (*str) {
        *str = tolower((unsigned char) *str);
        str++;
    }
}

bool isDuplicate(char **variables, int count, const char *variable) {

    for (int i = 0; i < count; i++) {
        if (strcasecmp(variables[i], variable) == 0) {
            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

bool isVariablePresent(Config *config, char *varName) {

    for (int count = 0; count < config->deploy->envCount; count++) {
        if (strcasecmp(config->deploy->keyValuePair[count].key,
                       varName) == 0) {
            return USYS_TRUE;
        }
    }

    return USYS_FALSE;
}

static bool deploy_system(DeployConfig *deployConfig, char *name, char *path) {

	char runMe[MAX_BUFFER]     = {0};
    char fileName[MAX_BUFFER]  = {0};
    char line[MAX_LINE_LENGTH] = {0};

    FILE *file = NULL;
    int count = 0, i, j;
    char **envs = NULL;

    sprintf(fileName, "%s/docker-compose.yml", path);

    file = fopen(fileName, "r");
    if (file == NULL) {
        usys_log_error("Error opening file: %s", fileName);
        return USYS_FALSE;
    }

    envs = (char **)calloc(MAX_VARIABLES, MAX_LINE_LENGTH);
    if (envs == NULL) {
        usys_log_error("Unable to allocate memory of size: %d",
                       MAX_VARIABLES * MAX_LINE_LENGTH);
        return USYS_FALSE;
    }

    while (fgets(line, MAX_LINE_LENGTH, file)) {
        char *start, *end;
        start = line;
        while ((start = strchr(start, '$')) && *(start + 1) == '{') {
            end = strchr(start, '}');
            if (end && (end - start > 2)) {
                char varName[MAX_LINE_LENGTH] = {0};
                *end = '\0';
                strncpy(varName, start + 2, end - start - 2);
                varName[end - start - 2] = '\0';
                toLowerCase(varName);

                if (!isDuplicate(envs, count, varName)) {
                    strcpy(envs[count++], varName);
                }

                *end = '}';
                start = end + 1;
            } else {
                break;
            }
        }
    }
    fclose(file);

    if (count > deployConfig->envCount) {
        usys_log_error("No enough var defined in config file");
        usys_free(envs);

        return USYS_FALSE;
    }

    /* setup env variables */
    for (i = 0; i < count; i++) {
        for (j = 0; j < deployConfig->envCount; j++) {
            if (strcasecmp(deployConfig->keyValuePair[i].key,
                           envs[i]) == 0 ||
                ((strcasecmp(name, "init") == 0 ||
                  strcasecmp(name, "ukama-auth") == 0)) ) {
                if (setenv(deployConfig->keyValuePair[i].key,
                           deployConfig->keyValuePair[i].value, 1) == -1) {
                    usys_log_error("Unable to set env variable");
                    usys_free(envs);

                    return USYS_FALSE;
                }
            }
        }
    }

    sprintf(runMe, "%s system %s %s", SCRIPT, name, path);
    if (system(runMe) < 0) return USYS_FALSE;

    usys_free(envs);
    return USYS_TRUE;
}

bool deploy_all_systems(DeployConfig *deployConfig, char *ukamaRepo, char *authRepo) {

    char list[MAX_BUFFER] = {0};
    char systemPath[MAX_BUFFER] = {0};
    char *systemName = NULL;

    strncpy(list, deployConfig->systemsList, strlen(deployConfig->systemsList));
    list[strlen(deployConfig->systemsList) - 1] = '\0';

    systemName = strtok(list, DELIMINATOR);
    while (systemName != NULL) {

        if (strcasecmp(systemName, UKAMA_AUTH) == 0) {
            if (!deploy_system(deployConfig, systemName, authRepo)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, authRepo);
                return USYS_FALSE;
            }
        } else {
            sprintf(systemPath, "%s/systems/%s/", ukamaRepo, systemName);
            if (!deploy_system(deployConfig, systemName, systemPath)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, systemPath);
                return USYS_FALSE;
            }
        }

        systemName = strtok(NULL, DELIMINATOR);
    }

    return USYS_TRUE;
}

bool deploy_node(char *id) {

    char *nodeID = NULL;
    char runMe[MAX_BUFFER] = {0};

    if (strcmp(id, "random") == 0) {
        nodeID = DEF_NODE_ID;
    } else {
        nodeID = id;
    }

    sprintf(runMe, "%s node %s", SCRIPT, nodeID);
    if (system(runMe) < 0) return USYS_FALSE;

    return USYS_TRUE;
}
