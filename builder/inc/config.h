/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#ifndef CONFIG_H_
#define CONFIG_H_

#include <stdbool.h>

typedef struct setupConfig_ {

    char *networkInterface;
    char *buildOS;
    char *ukamaRepo;
    char *authRepo;
} SetupConfig;

typedef struct buildConfig_ {

    int  nodeCount;
    char *nodeIDsList;
    char *systemsList;
    char *interfacesList;
} BuildConfig;

typedef struct deployConfig_ {

    char *email;
    char *name;
    char *orgName;
    char *orgID;

    char *systemsList;
    char *nodeIDsList;
} DeployConfig;

typedef struct config_ {

    SetupConfig *setup;
    BuildConfig *build;
    DeployConfig *deploy;
} Config;

bool read_config_file(Config **config, char *fileName);
void free_config(Config *config);

#endif /* CONFIG_H_ */
