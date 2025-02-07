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

#define DEPLOY_SCRIPT "./scripts/deploy-system.sh"
#define STATUS_SCRIPT "./scripts/status.sh"

#define FREE_ENVS(envs, max_variables) \
    do { \
        for (int i = 0; i < (max_variables); i++) { \
            usys_free((envs)[i]); \
        } \
        usys_free(envs); \
    } while(0)

void toLowerCase(char *str) {
    while (*str) {
        *str = tolower((unsigned char) *str);
        str++;
    }
}

bool isDuplicate(char **variables, int count, const char *variable) {

    for (int i = 0; i < count; i++) {
        if (variables[i] && strcasecmp(variables[i], variable) == 0) {
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

static bool deploy_system(char *configFile,
                          DeployConfig *deployConfig,
                          char *name, char *path) {

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

    envs = (char **)calloc(MAX_VARIABLES, sizeof (char *));
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

                    envs[count] = (char *)calloc(MAX_LINE_LENGTH, sizeof(char));
                    if (envs[count] == NULL) {
                        usys_log_error("Unable to allocate memory of size: %d",
                                       MAX_LINE_LENGTH * sizeof(char));
                        FREE_ENVS(envs, MAX_VARIABLES);
                        return USYS_FALSE;
                    }
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
        FREE_ENVS(envs, MAX_VARIABLES);
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
                    FREE_ENVS(envs, MAX_VARIABLES);
                    return USYS_FALSE;
                }
            }
        }
    }

    sprintf(runMe, "%s system %s %s %s", DEPLOY_SCRIPT, name, path, configFile);
    if (system(runMe) < 0) return USYS_FALSE;

    FREE_ENVS(envs, MAX_VARIABLES);
    return USYS_TRUE;
}

bool deploy_all_systems(char *configFilename,
                        DeployConfig *deployConfig,
                        char *ukamaRepo,
                        char *authRepo) {

    int index;
    char list[MAX_BUFFER]       = {0};
    char systemPath[MAX_BUFFER] = {0};
    char runMe[MAX_BUFFER]      = {0};
    char *systemName = NULL;

    strncpy(list, deployConfig->systemsList, strlen(deployConfig->systemsList));
    list[strlen(deployConfig->systemsList) - 1] = '\0';

    systemName = strtok(list, DELIMINATOR);
    while (systemName != NULL) {

        if (strcasecmp(systemName, UKAMA_AUTH) == 0) {
            if (!deploy_system(configFilename, deployConfig, systemName, authRepo)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, authRepo);
                return USYS_FALSE;
            }
        } else {
            sprintf(systemPath, "%s/systems/%s/", ukamaRepo, systemName);
            if (!deploy_system(configFilename, deployConfig, systemName, systemPath)) {
                usys_log_error("Build failed: %s path: %s",
                               systemName, systemPath);
                return USYS_FALSE;
            }
        }

        systemName = strtok(NULL, DELIMINATOR);
    }

    /* add org and nodes to init system */
    sprintf(runMe, "%s add-org-to-init-system %s", DEPLOY_SCRIPT, ukamaRepo);
    if (system(runMe) < 0) {
        usys_log_error("Unable to add default org to init system");
        return USYS_FALSE;
    }

    for (index=0; index < deployConfig->nodesCount; index++) {
        sprintf(runMe, "%s add-node-to-init-system %s %s",
                DEPLOY_SCRIPT,
                ukamaRepo,
                deployConfig->nodesIDList[index]);
        if  (system(runMe) < 0) {
            usys_log_error("Unable to add node id %s to init system",
                           deployConfig->nodesIDList[index]);
            continue;
        }
    }

    return USYS_TRUE;
}

bool deploy_nodes(int count, char **nodesIDList) {

    int i;
    char runMe[MAX_BUFFER] = {0};

    for (i=0; i<count; i++) {
        sprintf(runMe, "%s node %s", DEPLOY_SCRIPT, nodesIDList[i]);
        if (system(runMe) < 0) {
            usys_log_error("Unable to deploy node: %s", nodesIDList[i]);
            continue;
        }
    }

    return USYS_TRUE;
}

bool display_all_systems_status(char *systems, int interval) {

    char runMe[MAX_BUFFER] = {0};

    sprintf(runMe, "%s %s %d", STATUS_SCRIPT, systems, interval);
    if (system(runMe) < 0) return USYS_FALSE;

    return USYS_TRUE;
}
