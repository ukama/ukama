/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * capp creation related functions
 */

#include <stdio.h>
#include <jansson.h>
#include <errno.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <sys/types.h>
#include <dirent.h>

#include "capp_config.h"
#include "starter.h"

#include "usys_mem.h"

static void log_capp_config(CappConfig *config) {

    usys_log_debug("Version: %s",
                   (config->version ? config->version: "NULL"));
    usys_log_debug("hostname: %s",
                   (config->hostName ? config->hostName : "NULL"));

    if (config->process) {
        usys_log_debug("process:", config->process->exec);
        usys_log_debug("\t exec: %s", config->process->exec);
        if (config->process->argv) {
            usys_log_debug("\t argv: %s", config->process->argv);
        }
        if (config->process->env) {
            usys_log_debug("\t env: %s", config->process->env);
        }
    } else {
        usys_log_debug("process: NULL");
    }
}

static void json_log(json_t *json) {

    char *str = NULL;

    str = json_dumps(json, 0);
    if (str) {
        log_debug("json str: %s", str);
        usys_free(str);
    }
}

static bool get_json_entry(json_t *json, char *key, json_type type,
                           char **strValue,
                           json_t **jsonObj) {

    json_t *jEntry=NULL;

    if (json == NULL || key == NULL) return USYS_FALSE;

    jEntry = json_object_get(json, key);
    if (jEntry == NULL) {
        log_error("Missing %s key in json", key);
        return USYS_FALSE;
    }

    switch(type) {
    case (JSON_STRING):
        *strValue = strdup(json_string_value(jEntry));
        break;
    case (JSON_OBJECT):
        *jsonObj = json_object_get(json, key);
        break;
    default:
        log_error("Invalid type for json key-value: %d", type);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void free_capp_config(CappConfig *config) {

    if (config == NULL) return;

    usys_free(config->version);
    usys_free(config->hostName);

    if (config->process) {
        usys_free(config->process->exec);
        usys_free(config->process->argv);
        usys_free(config->process->env);
        usys_free(config->process);
    }

    usys_free(config);
}

static int deserialize_capp_config_file(CappConfig **config, json_t *json) {

    int ret = USYS_TRUE;
    json_t *jProc = NULL;

    if (json == NULL) return USYS_FALSE;

    *config = (CappConfig *)calloc(1, sizeof(CappConfig));

    ret &= get_json_entry(json, JTAG_VERSION, JSON_STRING,
                          &(*config)->version, NULL);
    ret &= get_json_entry(json, JTAG_PROCESS, JSON_OBJECT,
                          NULL, &jProc);
    if (get_json_entry(json, JTAG_HOSTNAME, JSON_STRING,
                       &(*config)->hostName, NULL) == USYS_FALSE) {
        (*config)->hostName = strdup(CAPP_DEFAULT_HOSTNAME);
    }

    if (ret == USYS_FALSE) {
        usys_log_error("Error deserializing the capp config.json");
        json_log(json);
        free_capp_config(*config);

        return USYS_FALSE;
    }

    (*config)->process = (CappProc *)calloc(1, sizeof(CappProc));

    ret &= get_json_entry(jProc, JTAG_EXEC, JSON_STRING,
                          &(*config)->process->exec, NULL);
    ret &= get_json_entry(jProc, JTAG_ARGS, JSON_STRING,
                          &(*config)->process->argv, NULL);
    ret &= get_json_entry(jProc, JTAG_ENV, JSON_STRING,
                          &(*config)->process->env, NULL);

    if (ret == USYS_FALSE) {
        usys_log_error("No valid process info found.");
        json_log(jProc);
        free_capp_config(*config);

        return USYS_FALSE;
    }
    
    usys_log_debug("capp config.json successfully parsed");
    
    return USYS_TRUE;
}

bool process_capp_config_file(CappConfig **config, char *fileName) {

    bool ret;
    FILE *fp=NULL;
    char *buffer=NULL;
    long size=0;
    json_t *json=NULL;
    json_error_t jerror;

    if (fileName == NULL) return USYS_FALSE;

    if ((fp = fopen(fileName, "r")) == NULL) {
        usys_log_error("Error opening file: %s Error %s",
                       fileName, strerror(errno));
        return USYS_FALSE;
    }

    usys_log_debug("Reading the capp's config.json file: %s", fileName);

    /* Read everything into buffer */
    fseek(fp, 0, SEEK_END);
    size = ftell(fp);
    fseek(fp, 0, SEEK_SET);
 
    if (size > CAPP_CONFIG_MAX_SIZE) {
        log_error("Error opening file: %s Error: File size too big: %ld",
                  fileName, size);
        fclose(fp);
        return USYS_FALSE;
    }

    buffer = (char *)calloc(1, size+1);
    if (buffer == NULL) {
        usys_log_error("Error allocating memory of size: %ld", size+1);
        fclose(fp);
        return USYS_FALSE;
    }
    memset(buffer, 0, size+1);
    fread(buffer, 1, size, fp);
    fclose(fp);

    json = json_loads(buffer, 0, &jerror);
    if (json == NULL) {
        log_error("Error loading capp config into JSON format."
                  "File: %s Size: %ld", fileName, size);
        log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
        usys_free(buffer);
        return USYS_FALSE;
    } else {
        usys_free(buffer);
        log_debug("JSON successfully loaded. %s", fileName);
    }

    ret = deserialize_capp_config_file(config, json);

    if (ret == USYS_TRUE) {
        log_capp_config(*config);
    } else {
        usys_log_error("Error with capp config file: %s", fileName);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}
