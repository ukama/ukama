/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2022-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <errno.h>
#include <stdio.h>
#include <unistd.h>
#include <dirent.h>
#include <sys/types.h>
#include <sys/stat.h>

#include "config.h"
#include "toml.h"
#include "log.h"

static int read_entry(toml_table_t *table, char *key, char **destStr,
                      int *destInt, int flag);
static int read_capp_table(toml_table_t *table, Config *config, char *type);
static int read_build_table(toml_table_t *table, Config *config, char *type);
static int get_table(toml_table_t *src, char *key, toml_table_t **table);
static int read_build_config(Config *config, char *fileName,
                             toml_table_t *buildFrom,
                             toml_table_t *buildCompile,
                             toml_table_t *buildRootfs,
                             toml_table_t *buildConf);
static int read_capp_config(Config *config, char *fileName,
                            toml_table_t *cappExec);
static int add_to_configs(Configs **configs, Config *config, char *fileName,
                          int status, char *errorStr);
static int read_config_file(Config *config, char *fileName, char **error);
static int is_valid_file(char *fileName);

static int append_env_var(EnvVar **envList, const char *key, const char *value);
static int read_env_table(toml_table_t *table, EnvVar **envList);
static void free_env_vars(EnvVar *envList);
static void log_env_vars(EnvVar *envList);

static int read_entry(toml_table_t *table,
                      char *key,
                      char **destStr,
                      int *destInt,
                      int flag) {

    toml_datum_t datum;
    int ret = TRUE;

    if (table == NULL || key == NULL) return FALSE;

    if (flag & DATUM_BOOL) {
        if (destInt == NULL) return FALSE;

        datum = toml_bool_in(table, key);
        if (datum.ok) {
            *destInt = datum.u.b ? TRUE : FALSE;
        } else {
            toml_datum_t s = toml_string_in(table, key);
            if (s.ok) {
                if (strcasecmp(s.u.s, "TRUE") == 0) {
                    *destInt = TRUE;
                } else if (strcasecmp(s.u.s, "FALSE") == 0) {
                    *destInt = FALSE;
                } else {
                    log_error("[%s] is invalid, expect true/false", key);
                    *destInt = -1;
                    ret = FALSE;
                }
                free(s.u.s);
            } else if (flag & DATUM_MANDATORY) {
                log_error("[%s] is missing but is required", key);
                return FALSE;
            }
        }
    }

    if (flag & DATUM_STRING) {
        if (destStr == NULL) return FALSE;

        datum = toml_string_in(table, key);
        if (datum.ok) {
            *destStr = strdup(datum.u.s);
            free(datum.u.s);
        } else if (flag & DATUM_MANDATORY) {
            log_error("[%s] is missing but is required", key);
            return FALSE;
        }
    }

    if (flag & DATUM_INT) {
        if (destInt == NULL) return FALSE;

        datum = toml_int_in(table, key);
        if (datum.ok) {
            *destInt = (int)datum.u.i;
        } else if (flag & DATUM_MANDATORY) {
            log_error("[%s] is missing but is required", key);
            return FALSE;
        }
    }

    return ret;
}

static int append_env_var(EnvVar **envList, const char *key, const char *value) {

    EnvVar *node = NULL;
    EnvVar *tail = NULL;

    if (!envList || !key || !*key || !value) {
        return FALSE;
    }

    node = (EnvVar *)calloc(1, sizeof(EnvVar));
    if (!node) {
        log_error("Error allocating memory of size: %lu", sizeof(EnvVar));
        return FALSE;
    }

    node->key = strdup(key);
    node->value = strdup(value);
    if (!node->key || !node->value) {
        if (node->key) free(node->key);
        if (node->value) free(node->value);
        free(node);
        return FALSE;
    }

    if (*envList == NULL) {
        *envList = node;
        return TRUE;
    }

    for (tail = *envList; tail->next; tail = tail->next);
    tail->next = node;

    return TRUE;
}

static int read_env_table(toml_table_t *table, EnvVar **envList) {

    int i;
    const char *key;
    toml_datum_t datum;

    if (!envList) {
        return FALSE;
    }

    if (!table) {
        return TRUE;
    }

    i = 0;
    while ((key = toml_key_in(table, i++)) != NULL) {
        datum = toml_string_in(table, key);
        if (!datum.ok) {
            log_error("[capp-exec.env.%s] must be a string", key);
            return FALSE;
        }

        if (!append_env_var(envList, key, datum.u.s)) {
            free(datum.u.s);
            return FALSE;
        }

        free(datum.u.s);
    }

    return TRUE;
}

static void free_env_vars(EnvVar *envList) {

    EnvVar *curr = NULL;
    EnvVar *next = NULL;

    curr = envList;
    while (curr) {
        next = curr->next;
        if (curr->key) free(curr->key);
        if (curr->value) free(curr->value);
        free(curr);
        curr = next;
    }
}

static void log_env_vars(EnvVar *envList) {

    EnvVar *curr = NULL;

    if (!envList) {
        log_debug("\t No env");
        return;
    }

    for (curr = envList; curr; curr = curr->next) {
        log_debug("\t env[%s]=%s", curr->key, curr->value);
    }
}

static int read_capp_table(toml_table_t *table, Config *config,
                           char *type) {

    CappConfig *capp = NULL;
    toml_table_t *envTable = NULL;

    if (table == NULL || config == NULL) return FALSE;

    if (config->capp == NULL) {
        config->capp = (CappConfig *)calloc(1, sizeof(CappConfig));
        if (config->capp == NULL) {
            log_error("[%s]: Error allocating memory of size: %d", TABLE_CAPP,
                      (int)sizeof(CappConfig));
            return FALSE;
        }
    }

    capp = config->capp;

    if (strcmp(type, TABLE_CAPP_EXEC)==0) {
        if (!read_entry(table, KEY_NAME, &capp->name, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_VERSION, &capp->version, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_BIN, &capp->bin, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_PATH, &capp->path, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_ARGS, &capp->args, NULL, DATUM_STRING)) {
            return FALSE;
        }

        envTable = toml_table_in(table, KEY_ENV);
        if (!read_env_table(envTable, &capp->env)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_AUTOSTART, NULL, &capp->autostart,
                        DATUM_BOOL | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_AUTORESTART, NULL, &capp->autorestart,
                        DATUM_BOOL | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_DEPENDS_ON, &capp->dependsOn, NULL,
                        DATUM_STRING)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_WAIT_FOR, &capp->waitFor, NULL,
                        DATUM_STRING)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_GROUP, &capp->group, NULL,
                        DATUM_STRING)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_RETRY, NULL, &capp->startretries,
                        DATUM_INT | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else {
        return FALSE;
    }

    return TRUE;
}

static int read_build_table(toml_table_t *table, Config *config,
                            char *type) {

    BuildConfig *build = NULL;

    if (table == NULL || config == NULL) return FALSE;

    if (config->build == NULL) {
        config->build = (BuildConfig *)calloc(1, sizeof(BuildConfig));
        if (config->build == NULL) {
            log_error("[%s]: Error allocating memory of size: %d", TABLE_BUILD,
                      (int)sizeof(BuildConfig));
            return FALSE;
        }
    }

    build = config->build;

    if (strcmp(type, TABLE_BUILD_FROM)==0) {
        if (!read_entry(table, KEY_BASE, &build->baseImage, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_VERSION, &build->baseVersion, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_BUILD_COMPILE)==0) {
        if (!read_entry(table, KEY_VERSION, &build->version, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_STATIC, NULL, &build->staticFlag,
                        DATUM_BOOL | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_SOURCE, &build->source, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_CMD, &build->cmd, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_BIN_FROM, &build->binFrom, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_BIN_TO, &build->binTo, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_BUILD_ROOTFS)==0) {
        if (!read_entry(table, KEY_MKDIR, &build->mkdir, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_BUILD_CONF)==0) {
        if (!read_entry(table, KEY_FROM, &build->from, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_TO, &build->to, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else {
        return FALSE;
    }

    return TRUE;
}

static int get_table(toml_table_t *src, char *key, toml_table_t **table) {

    if (src == NULL || key == NULL || table == NULL) return FALSE;

    *table = toml_table_in(src, key);
    if ((*table) == NULL) {
        log_error("Missing required section: [%s]", key);
        return FALSE;
    }

    return TRUE;
}

static int read_build_config(Config *config, char *fileName,
                             toml_table_t *buildFrom,
                             toml_table_t *buildCompile,
                             toml_table_t *buildRootfs,
                             toml_table_t *buildConf) {

    int ret = TRUE;

    ret = read_build_table(buildFrom, config, TABLE_BUILD_FROM);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_BUILD_FROM, fileName);
        free_config(config, BUILD_ONLY);
        goto done;
    }

    ret = read_build_table(buildCompile, config, TABLE_BUILD_COMPILE);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_BUILD_COMPILE, fileName);
        free_config(config, BUILD_ONLY);
        goto done;
    }

    if (buildRootfs) {
        ret = read_build_table(buildRootfs, config, TABLE_BUILD_ROOTFS);
        if (ret == FALSE) {
            log_error("[%s] section parsing error in config file: %s\n",
                      TABLE_BUILD_ROOTFS, fileName);
            free_config(config, BUILD_ONLY);
            goto done;
        }
    }

    if (buildConf) {
        ret = read_build_table(buildConf, config, TABLE_BUILD_CONF);
        if (ret == FALSE) {
            log_error("[%s] section parsing error in config file: %s\n",
                      TABLE_BUILD_CONF, fileName);
            free_config(config, BUILD_ONLY);
            goto done;
        }
    }

done:
    return ret;
}

static int read_capp_config(Config *config, char *fileName,
                            toml_table_t *cappExec) {

    int ret = TRUE;

    ret = read_capp_table(cappExec, config, TABLE_CAPP_EXEC);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_CAPP_EXEC, fileName);
        free_config(config, CAPP_ONLY);
        goto done;
    }

done:
    return ret;
}

static int add_to_configs(Configs **configs, Config *config, char *fileName,
                          int status, char *errorStr) {

    Configs *ptr = NULL;

    if ((*configs) == NULL) {
        *configs = (Configs *)calloc(1, sizeof(Configs));
        if (*configs == NULL) {
            log_error("Error allocating memory of size: %lu", sizeof(Configs));
            return FALSE;
        }
        ptr = *configs;
    } else {
        for (ptr = (*configs); ptr->next; ptr = ptr->next);
        ptr->next = (Configs *)calloc(1, sizeof(Configs));
        if (ptr->next == NULL) {
            log_error("Error allocating memory of size: %lu", sizeof(Configs));
            return FALSE;
        }
        ptr = ptr->next;
    }

    ptr->fileName = strdup(fileName);
    ptr->valid    = status;
    ptr->config   = config;
    ptr->next     = NULL;

    if (!ptr->valid && errorStr) {
        ptr->errorStr = strdup(errorStr);
    }

    return TRUE;
}

static int is_valid_file(char *fileName) {

    struct stat statBuf;

    if (fileName == NULL) return FALSE;

    if (stat(fileName, &statBuf) != 0) {
        return FALSE;
    }

    if (!S_ISREG(statBuf.st_mode)) {
        return FALSE;
    }

    return TRUE;
}

int read_config_files(Configs **configs, char *configDir, BoardConfig *boardCfg) {

    int ret = TRUE;
    int configStatus = FALSE;
    struct stat statBuf;
    struct dirent *dp = NULL;
    DIR *dir = NULL;
    char *configFile = NULL;
    char *errorStr = NULL;
    Config *config = NULL;
    char buffer[MAX_BUFFER] = {0};

    if (configs == NULL || configDir == NULL) return FALSE;

    if (stat(configDir, &statBuf) != 0) {
        log_error("Error reading dir at: %s. Error: %s", configDir,
                  strerror(errno));
        return FALSE;
    }

    if (!S_ISDIR(statBuf.st_mode)) {
        log_error("Invalid capp config dir: %s", configDir);
        return FALSE;
    }

    dir = opendir(configDir);
    if (!dir) {
        log_error("Unable to open capp config dir: %s", configDir);
        return FALSE;
    }

    while ((dp = readdir(dir)) != NULL) {
        memset(buffer, 0, MAX_BUFFER);
        sprintf(buffer, "%s/%s", configDir, dp->d_name);
        configFile = realpath(buffer, NULL);
        if (!configFile || !is_valid_file(configFile)) {
            if (configFile) free(configFile);
            continue;
        }

        config = (Config *)calloc(1, sizeof(Config));
        if (config == NULL) {
            log_error("Error allocating memory of size: %lu", sizeof(Config));
            goto failure;
        }

        if (!read_config_file(config, configFile, &errorStr)) {
            log_error("Parsing error for: %s", configFile);
            ret = FALSE;
            configStatus = FALSE;
        } else {
            configStatus = TRUE;
            log_config(config);
        }

        if (boardCfg && config->capp && config->capp->name) {
            if (!board_is_app_enabled(boardCfg, config->capp->name)) {
                log_debug("Skipping %s (disabled for this board)",
                          config->capp->name);

                if (errorStr) {
                    free(errorStr);
                    errorStr = NULL;
                }

                free(configFile);
                configFile = NULL;

                free_config(config, BUILD_ONLY | CAPP_ONLY);
                config = NULL;

                continue;
            }
        }

        add_to_configs(configs, config, configFile, configStatus, errorStr);

        if (errorStr) {
            free(errorStr);
            errorStr = NULL;
        }

        free(configFile);
        configFile = NULL;
        config = NULL;
    }

    closedir(dir);
    return ret;

failure:
    free_configs(*configs);
    if (config) free_config(config, BUILD_ONLY | CAPP_ONLY);
    if (configFile) free(configFile);
    if (errorStr) free(errorStr);
    closedir(dir);

    return FALSE;
}

static int read_config_file(Config *config, char *fileName, char **error) {

    int ret = FALSE;
    FILE *fp = NULL;

    toml_table_t *fileData = NULL;
    toml_table_t *buildFrom = NULL;
    toml_table_t *buildCompile = NULL;
    toml_table_t *buildRootfs = NULL;
    toml_table_t *buildConf = NULL;
    toml_table_t *cappExec = NULL;

    if (fileName == NULL || config == NULL) {
        return FALSE;
    }

    if ((fp = fopen(fileName, "r")) == NULL) {
        log_error("Error opening config file: %s: %s\n", fileName,
                  strerror(errno));
        return FALSE;
    }

    *error = (char *)calloc(1, MAX_ERROR_BUFFER);
    if (*error == NULL) {
        fclose(fp);
        log_error("Error allocating memory of size: %lu", MAX_ERROR_BUFFER);
        return FALSE;
    }

    fileData = toml_parse_file(fp, *error, MAX_ERROR_BUFFER);
    fclose(fp);

    if (!fileData) {
        log_error("Error parsing the config file %s: %s\n", fileName, *error);
        return FALSE;
    }

    if (!get_table(fileData, TABLE_BUILD_FROM,    &buildFrom))    goto done;
    if (!get_table(fileData, TABLE_BUILD_COMPILE, &buildCompile)) goto done;
    if (!get_table(fileData, TABLE_CAPP_EXEC,     &cappExec))     goto done;

    get_table(fileData, TABLE_BUILD_ROOTFS, &buildRootfs);
    get_table(fileData, TABLE_BUILD_CONF,   &buildConf);

    ret = read_build_config(config, fileName, buildFrom, buildCompile,
                            buildRootfs, buildConf);
    if (ret == FALSE) goto done;

    ret = read_capp_config(config, fileName, cappExec);
    if (ret == FALSE) goto done;

    ret = TRUE;

done:
    toml_free(fileData);
    return ret;
}

void free_configs(Configs *configs) {

    Configs *ptr = NULL;
    Configs *next = NULL;

    if (!configs) return;

    ptr = configs;
    while (ptr) {
        next = ptr->next;

        if (ptr->fileName) free(ptr->fileName);
        if (ptr->errorStr) free(ptr->errorStr);

        free_config(ptr->config, BUILD_ONLY | CAPP_ONLY);
        free(ptr);

        ptr = next;
    }
}

void free_config(Config *config, int flag) {

    BuildConfig *build;
    CappConfig *capp;

    if (!config) return;

    if ((flag & BUILD_ONLY) && config->build) {
        build = config->build;

        if (build->baseImage)   free(build->baseImage);
        if (build->baseVersion) free(build->baseVersion);
        if (build->version)     free(build->version);
        if (build->source)      free(build->source);
        if (build->cmd)         free(build->cmd);
        if (build->binFrom)     free(build->binFrom);
        if (build->binTo)       free(build->binTo);
        if (build->mkdir)       free(build->mkdir);
        if (build->from)        free(build->from);
        if (build->to)          free(build->to);

        free(config->build);
    }

    if ((flag & CAPP_ONLY) && config->capp) {
        capp = config->capp;

        if (capp->name)      free(capp->name);
        if (capp->version)   free(capp->version);
        if (capp->bin)       free(capp->bin);
        if (capp->path)      free(capp->path);
        if (capp->args)      free(capp->args);
        if (capp->env)       free_env_vars(capp->env);
        if (capp->waitFor)   free(capp->waitFor);
        if (capp->dependsOn) free(capp->dependsOn);
        if (capp->group)     free(capp->group);

        free(config->capp);
    }

    if ((flag & CAPP_ONLY) && (flag & BUILD_ONLY)) {
        free(config);
    }
}

void log_config(Config *config) {

    BuildConfig *build = NULL;
    CappConfig *capp = NULL;

    if (config == NULL) return;

    if (config->build) {
        build = config->build;
        log_debug("--- APP Build Configuration ---");

        log_debug("[FROM:]");
        log_debug("\t base:    %s", build->baseImage);
        log_debug("\t version: %s", build->baseVersion);

        log_debug("[BUILD:]");
        log_debug("\t version: %s", build->version);
        log_debug("\t static:  %s", (build->staticFlag ? "true" : "false"));
        log_debug("\t source:  %s", build->source);
        log_debug("\t command: %s", build->cmd);
        log_debug("\t from:    %s", build->binFrom);
        log_debug("\t to:      %s", build->binTo);

        log_debug("[ROOTFS:]");
        log_debug("\t mkdir: %s", build->mkdir);

        log_debug("[CONF:]");
        log_debug("\t from: %s", build->from);
        log_debug("\t to:   %s", build->to);
    } else {
        log_debug("No Build configuration found");
    }

    if (config->capp) {
        capp = config->capp;
        log_debug("--- APP Run Configuration ---");

        log_debug("[EXEC:]");
        log_debug("\t name: %s", capp->name);
        log_debug("\t version: %s", capp->version);
        log_debug("\t autostart:  %s", (capp->autostart ? "true" : "false"));
        log_debug("\t autorestart:  %s", (capp->autorestart ? "true" : "false"));
        log_debug("\t binary:  %s", capp->bin);
        log_debug("\t path: %s", capp->path);

        if (capp->args) {
            log_debug("\t args: %s", capp->args);
        } else {
            log_debug("\t No args");
        }

        log_env_vars(capp->env);

        if (capp->dependsOn) {
            log_debug("\t dependsOn: %s", capp->dependsOn);
        }

        if (capp->waitFor) {
            log_debug("\t waitFor: %s", capp->waitFor);
        }
    }
}
