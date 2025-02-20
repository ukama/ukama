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

#define SCRIPT        "./build.sh"
#define NODE_SCRIPT   "./build-virtual-node.sh"
#define BASE_IMAGE_ID "uk-ma0000-tnode-a1-1234"

#define BUILD_ANODE_SCRIPT "./build-amplifier-node.sh"

/* board_config.c */
extern char *getAppsFromBoardConfigs(const char *commonFile, const char *boardFile);

static bool build_system(char *name, char *path) {

	char runMe[MAX_BUFFER] = {0};

	if (name == NULL || path == NULL) return USYS_FALSE;

	sprintf(runMe, "cd scripts; %s system %s; cd -", SCRIPT, path);
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

bool build_nodes(char *repo, int count, char **list) {

    int i;
    char runMe[MAX_BUFFER] = {0};

    if (getenv(ENV_DOCKER_BUILD)) {
        sprintf(runMe, "cd scripts; %s base-image %s %s; cd -",
                NODE_SCRIPT, repo, BASE_IMAGE_ID);
    } else {
        sprintf(runMe, "cd scripts; sudo %s base-image %s %s; cd -",
                NODE_SCRIPT, repo, BASE_IMAGE_ID);
    }
    if (system(runMe) < 0) {
        usys_log_error("Unable to create base image via repo: %s", repo);
        return USYS_FALSE;
    }

    for (i=0; i<count; i++) {

        if (getenv(ENV_DOCKER_BUILD)) {
            sprintf(runMe, "cd scripts; %s create-node %s %s %s; cd -",
                    NODE_SCRIPT, repo, list[i], BASE_IMAGE_ID);
        } else {
            sprintf(runMe, "cd scripts; sudo %s create-node %s %s %s; cd -",
                    NODE_SCRIPT, repo, list[i], BASE_IMAGE_ID);
        }
        if (system(runMe) < 0) {
            usys_log_error("Unable to create node with ID: %s", list[i]);
            continue;
        }
    }

    return USYS_TRUE;
}

bool build_ukamaos_image(char *repo) {

    char runMe[MAX_BUFFER] = {0};

    if (getenv(ENV_DOCKER_BUILD)) {
        sprintf(runMe, "cd scripts; %s base-image %s %s; cd -",
                SCRIPT, repo, BASE_IMAGE_ID);
    } else {
        sprintf(runMe, "cd scripts; sudo %s base-image %s %s; cd -",
                SCRIPT, repo, BASE_IMAGE_ID);
    }

    if (system(runMe) < 0) {
        usys_log_error("Unable to create base image via repo: %s", repo);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool build_amplifier_node(char *repo, char *nodeID) {

    char runMe[MAX_BUFFER] = {0};
    char *appsList = NULL;

    appsList = getAppsFromBoardConfigs(BOARD_COMMON_CONFIG,
                                       BOARD_CONTROLLER_CONFIG);

    if (appsList) {
        sprintf(runMe, "cd scripts; sudo %s %s 0.0.1 %s %s; cd -",
                BUILD_ANODE_SCRIPT,
                repo,
                nodeID,
                appsList);
    } else {
        sprintf(runMe, "cd scripts; %s %s 0.0.1 %s %s; cd -",
                BUILD_ANODE_SCRIPT,
                repo,
                nodeID,
                "");
    }

    if (system(runMe) < 0) {
        usys_log_error("Unable to build amplifier image", repo);
        usys_log_error(" common config: %s controller config: %s",
                       BOARD_COMMON_CONFIG, BOARD_CONTROLLER_CONFIG);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}
