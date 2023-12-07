/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>
#include <jansson.h>

#include "config.h"
#include "builder.h"
#include "json_types.h"

#include "usys_error.h"
#include "usys_log.h"
#include "usys_mem.h"
#include "usys_string.h"
#include "usys_types.h"

void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue, int *intValue) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        usys_log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_INTEGER):
        *intValue = json_integer_value(jEntry);
        break;
    default:
        usys_log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static bool deserialize_config_file(Config **config, json_t *json) {

    int ret = USYS_TRUE;

    json_t *jSetup  = NULL;
    json_t *jBuild  = NULL;
    json_t *jDeploy = NULL;
    json_t *jNodes  = NULL;
    json_t *jSystems     = NULL;
    json_t *jInterfaces  = NULL;
    json_t *jDeployEnv   = NULL;

    int count = 0;
    const char *key = NULL;
    json_t *value = NULL;

    jSetup = json_object_get(json, JTAG_SETUP);
    jBuild = json_object_get(json, JTAG_BUILD);
    jDeploy = json_object_get(json, JTAG_DEPLOY);

    if (jSetup == NULL || jBuild == NULL || jDeploy == NULL) {
        usys_log_error("Missing setup, build and/or deploy in config");
        return USYS_FALSE;
    }

    jNodes      = json_object_get(jBuild, JTAG_NODES);
    jSystems    = json_object_get(jBuild, JTAG_SYSTEMS);
    jInterfaces = json_object_get(jBuild, JTAG_INTERFACES);
    if (jNodes == NULL || jSystems == NULL || jInterfaces == NULL) {
        usys_log_error("Missing entries in build");
        return USYS_FALSE;
    }
    
    jDeployEnv = json_object_get(jDeploy, JTAG_ENV);
    if (jDeployEnv == NULL) {
        usys_log_debug("No env variable define in <deploy>");
    }

    *config = (Config *)calloc(1, sizeof(Config));
    if (*config == NULL) {
        usys_log_error("Error allocating memory of size: %ld", sizeof(Config));
        json_log(json);
        return USYS_FALSE;
    }

    (*config)->setup  = (SetupConfig *)calloc(1, sizeof(SetupConfig));
    (*config)->build  = (BuildConfig *)calloc(1, sizeof(BuildConfig));
    (*config)->deploy = (DeployConfig *)calloc(1, sizeof(DeployConfig));

    if ((*config)->setup == NULL ||
        (*config)->build == NULL ||
        (*config)->deploy == NULL) {
        usys_log_error("Error allocating memory");
        usys_free(*config);

        return USYS_FALSE;
    } 

    jDeployEnv = json_object_get(jDeploy, JTAG_ENV);
    if (jDeployEnv == NULL) {
        usys_log_debug("No env variable define in <deploy>");
    } else {
        (*config)->deploy->envCount     = json_object_size(jDeployEnv);
        (*config)->deploy->keyValuePair =
            calloc((*config)->deploy->envCount, sizeof(KeyValuePair));
    }

    /* setup */
    ret |= get_json_entry(jSetup, JTAG_NETWORK_INTERFACE, JSON_STRING,
                          &(*config)->setup->networkInterface, NULL);
    ret |= get_json_entry(jSetup, JTAG_BUILD_OS, JSON_STRING,
                          &(*config)->setup->buildOS, NULL);
    ret |= get_json_entry(jSetup, JTAG_UKAMA_REPO, JSON_STRING,
                          &(*config)->setup->ukamaRepo, NULL);
    ret |= get_json_entry(jSetup, JTAG_AUTH_REPO, JSON_STRING,
                          &(*config)->setup->authRepo, NULL);
    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing <setup>");
        free_config(*config);
        json_log(json);
        return USYS_FALSE;
    }

    /* build */
    ret |= get_json_entry(jNodes, JTAG_COUNT, JSON_INTEGER,
                          NULL, &(*config)->build->nodeCount);
    ret |= get_json_entry(jNodes, JTAG_IDS, JSON_STRING,
                          &(*config)->build->nodeIDsList, NULL);
    ret |= get_json_entry(jSystems, JTAG_LIST, JSON_STRING,
                          &(*config)->build->systemsList, NULL);
    ret |= get_json_entry(jInterfaces, JTAG_LIST, JSON_STRING,
                          &(*config)->build->interfacesList, NULL);
    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing <build>");
        free_config(*config);
        json_log(json);
        return USYS_FALSE;
    }

    /* deploy */
    json_object_foreach(jDeployEnv, key, value) {
        if (json_is_string(value)) {
            (*config)->deploy->keyValuePair[count].key   = strdup(key);
            (*config)->deploy->keyValuePair[count].value = strdup(json_string_value(value));
            count++;
        }
    }
    ret |= get_json_entry(jDeploy, JTAG_SYSTEMS, JSON_STRING,
                          &(*config)->deploy->systemsList, NULL);
    ret |= get_json_entry(jDeploy, JTAG_NODES, JSON_STRING,
                          &(*config)->deploy->nodeIDsList, NULL);

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing <deploy>");
        free_config(*config);
        json_log(json);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

bool read_config_file(Config **config, char *fileName) {

    bool ret;
    FILE *fp = NULL;
    char *buffer = NULL;
    long size = 0;
    json_t *json = NULL;
    json_error_t jerror;

    if ((fp = fopen(fileName, "r")) == NULL) {
        usys_log_error("Error opening config file: %s Error: %s",
                       fileName, strerror(errno));
        return USYS_FALSE;
    }

    /* Read everything into buffer */
    fseek(fp, 0, SEEK_END);
    size = ftell(fp);
    fseek(fp, 0, SEEK_SET);

    if (size > MAX_CONFIG_FILE_SIZE) {
        usys_log_error("Error opening config file: %s "
                       "Error: File size too big: %ld",
                       fileName, size);
        fclose(fp);
        return USYS_FALSE;
    }

    buffer = (char *)malloc(size+1);
    if (buffer == NULL) {
        usys_log_error("Error allocating memory of size: %ld", size+1);
        fclose(fp);
        return USYS_FALSE;
    }
    memset(buffer, 0, size+1);
    fread(buffer, 1, size, fp);

    /* load as JSON */
    json = json_loads(buffer, 0, &jerror);
    if (json == NULL) {
        usys_log_error("Error loading manifest into JSON format."
                       "File: %s Size: %ld", fileName, size);
        usys_log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
    } else {
        ret = deserialize_config_file(config, json);
    }

    if (buffer) free(buffer);

    fclose(fp);
    json_decref(json);

    return ret;
}

void free_config(Config *config) {

    int count;

    if (config == NULL) return;

    if (config->setup) {
        usys_free(config->setup->networkInterface);
        usys_free(config->setup->buildOS);
        usys_free(config->setup->ukamaRepo);
        usys_free(config->setup->authRepo);
    }

    if (config->build) {
        usys_free(config->build->nodeIDsList);
        usys_free(config->build->systemsList);
        usys_free(config->build->interfacesList);
    }

    if (config->deploy) {

        for (count = 0; count < config->deploy->envCount; count++) {
            usys_free(config->deploy->keyValuePair[count].key);
            usys_free(config->deploy->keyValuePair[count].value);
        }
        usys_free(config->deploy->keyValuePair);
        usys_free(config->deploy->systemsList);
        usys_free(config->deploy->nodeIDsList);
    }

    usys_free(config->setup);
    usys_free(config->build);
    usys_free(config->deploy);

    usys_free(config);
}
