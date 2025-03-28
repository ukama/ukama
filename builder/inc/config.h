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

typedef struct {
    char *key;
    char *value;
} KeyValuePair;

typedef struct setupConfig_ {

    char *networkInterface;
    char *buildOS;
    char *ukamaRepo;
    char *authRepo;
    int  statusInterval;
} SetupConfig;

typedef struct buildConfig_ {

    char *nodesIDFilename;
    char *kernelImage;
    char *initRAMImage;
    char *diskImage;
    char *systemsList;
    char *interfacesList;

    int  nodesCount;
    char **nodesIDList;
} BuildConfig;

typedef struct deployConfig_ {

    int          envCount;
    KeyValuePair *keyValuePair;

    char *systemsList;
    char *nodesIDFilename;

    int  nodesCount;
    char **nodesIDList;
} DeployConfig;

typedef struct config_ {

    char        *fileName;
    SetupConfig *setup;
    BuildConfig *build;
    DeployConfig *deploy;
} Config;

bool read_config_file(Config **config, char *fileName);
void free_config(Config *config);

#endif /* CONFIG_H_ */
