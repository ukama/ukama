/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>

#include "config_app.h"
#include "toml.h"
#include "log_app.h"

static int read_entry(toml_table_t *table, char *key, char **destStr,
                      int *destInt, int flag);
static int read_capp_table(toml_table_t *table, Config *config, char *type);
static int read_build_table(toml_table_t *table, Config *config, char *type);
static int read_build_table(toml_table_t *table, Config *config, char *type);
static int get_table(toml_table_t *src, char *key, toml_table_t **table);
static int read_build_config(Config *config, char *fileName,
                             toml_table_t *buildFrom,
                             toml_table_t *buildCompile,
                             toml_table_t *buildRootfs,
                             toml_table_t *buildConf,
                             toml_table_t *buildMisc);
static int read_capp_config(Config *config, char *fileName,
                            toml_table_t *cappExec,
                            toml_table_t *cappOutput);

static int read_entry(toml_table_t *table, char *key, char **destStr,
                      int *destInt, int flag) {

    char *rootPath = NULL;
    toml_datum_t datum;
    int ret=TRUE, size=0;

    if (table == NULL || key == NULL) return FALSE;

    datum = toml_string_in(table, key);

    if (datum.ok) {
        if (flag & DATUM_BOOL) {
            if (strcasecmp(datum.u.s, "TRUE")==0) {
                *destInt = TRUE;
            } else if (strcasecmp(datum.u.s, "FALSE")==0) {
                *destInt = FALSE;
            } else {
                log_error("[%s] is invalid, except 'true' or 'false'", key);
                *destInt = -1;
                ret      = FALSE;
            }
        } else if (flag & DATUM_STRING) {

            rootPath = getenv("UKAMA_ROOT");

            if (strcmp(key, KEY_FROM) != 0 &&
                strcmp(key, KEY_BIN_FROM) != 0 &&
                strcmp(key, KEY_SOURCE) != 0) {
                *destStr = strdup(datum.u.s);
            } else {
                if (rootPath == NULL) {
                    *destStr = strdup(datum.u.s);
                } else {

                    size = strlen(rootPath) + strlen(datum.u.s) + 2;
                    *destStr = (char *)malloc(size);
                    if ( *destStr == NULL) {
                        ret = FALSE;
                    }

                    sprintf(*destStr, "%s/%s", rootPath, datum.u.s);
                }
            }
        } else {
            ret = FALSE;
        }
    } else {
        if (flag & DATUM_MANDATORY) {
            log_error("[%s] is missing but is required", key);
            return FALSE;
        }
    }

    if (datum.ok) free(datum.u.s);

    return ret;
}

static int read_capp_table(toml_table_t *table,
                           Config *config,
                           char *type) {

    CappConfig *capp = NULL;

    if (table == NULL || config == NULL) return FALSE;

    if (config->capp == NULL) {
        config->capp = (CappConfig *)calloc(1, sizeof(CappConfig));
        if (config->capp == NULL) {
            log_error("[%s]: Error allocating memory of size: %d", TABLE_CAPP,
                      sizeof(CappConfig));
            return FALSE;
        }
    }

    capp = config->capp;

    if (strcmp(type, TABLE_CAPP_EXEC)==0) { /* [capp.exec] */
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

        if (!read_entry(table, KEY_ENVS, &capp->envs, NULL, DATUM_STRING)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_CAPP_OUTPUT)==0) { /* [capp.format] */

        if (!read_entry(table, KEY_FORMAT, &capp->format, NULL, DATUM_STRING)) {
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
                      sizeof(BuildConfig));
            return FALSE;
        }
    }

    build = config->build;

    if (strcmp(type, TABLE_BUILD_FROM)==0) { /* [build.from] */
        if (!read_entry(table, KEY_ROOTFS, &build->rootfs, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_CONTAINED, &build->contained, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_BUILD_COMPILE)==0) { /* [build.compile] */
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
    } else if (strcmp(type, TABLE_BUILD_ROOTFS)==0) { /* [build.rootfs] */
        if (!read_entry(table, KEY_MKDIR, &build->mkdir, NULL, DATUM_STRING)) {
            return FALSE;
        }
    } else if (strcmp(type, TABLE_BUILD_CONF)==0) { /* [build.conf] */
        if (!read_entry(table, KEY_FROM, &build->from, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_TO, &build->to, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

    } else if (strcmp(type, TABLE_BUILD_MISC)==0) { /* [build-misc] */
        if (!read_entry(table, KEY_FROM, &build->miscFrom, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }

        if (!read_entry(table, KEY_TO, &build->miscTo, NULL,
                        DATUM_STRING | DATUM_MANDATORY)) {
            return FALSE;
        }
    } else {
        return FALSE;
    }

    return TRUE;
}

static int get_table(toml_table_t *src, char *key, toml_table_t **table) {

    *table = toml_table_in(src, key);
    if (*table == NULL) {
        log_error("[%s] section parsing error in config file", key);
        return FALSE;
    }

    return TRUE;
}

static int read_build_config(Config *config, char *fileName,
                             toml_table_t *buildFrom,
                             toml_table_t *buildCompile,
                             toml_table_t *buildRootfs,
                             toml_table_t *buildConf,
                             toml_table_t *buildMisc) {

    int ret=TRUE;

    ret = read_build_table(buildFrom, config, TABLE_BUILD_FROM);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_BUILD_FROM, fileName);
        clear_config(config, BUILD_ONLY);
        goto done;
    }

    ret = read_build_table(buildCompile, config, TABLE_BUILD_COMPILE);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_BUILD_COMPILE, fileName);
        clear_config(config, BUILD_ONLY);
        goto done;
    }

    if (buildRootfs) {
        if (read_build_table(buildRootfs, config, TABLE_BUILD_ROOTFS) == FALSE) {
            log_error("[%s] section parsing error in config file: %s\n",
                      TABLE_BUILD_ROOTFS, fileName);
            clear_config(config, BUILD_ONLY);
            goto done;
        }
    }

    if (buildConf) {
        if (read_build_table(buildConf, config, TABLE_BUILD_CONF) == FALSE) {
            log_error("[%s] section parsing error in config file: %s\n",
                      TABLE_BUILD_CONF, fileName);
            clear_config(config, BUILD_ONLY);
            goto done;
        }
    }

    if (buildMisc) {
        if (read_build_table(buildMisc, config, TABLE_BUILD_MISC) == FALSE) {
            log_error("[%s] section parsing error in config file: %s Ignoring",
                      TABLE_BUILD_MISC, fileName);
            goto done;
        }
    }

done:
    return ret;
}

static int read_capp_config(Config *config, char *fileName,
                            toml_table_t *cappExec,
                            toml_table_t *cappOutput) {

    int ret=TRUE;

    ret = read_capp_table(cappExec, config, TABLE_CAPP_EXEC);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_CAPP_EXEC, fileName);
        clear_config(config, CAPP_ONLY);
        goto done;
    }

    ret = read_capp_table(cappOutput, config, TABLE_CAPP_OUTPUT);
    if (ret == FALSE) {
        log_error("[%s] section parsing error in config file: %s\n",
                  TABLE_CAPP_OUTPUT, fileName);
        clear_config(config, CAPP_ONLY);
        goto done;
    }

done:
    return ret;
}

int read_config_file(Config *config, char *fileName) {

    int ret=FALSE;
    FILE *fp=NULL;

    toml_table_t *fileData=NULL;
    toml_table_t *buildFrom=NULL, *buildCompile=NULL, *buildRootfs=NULL;
    toml_table_t *buildConf=NULL, *buildMisc=NULL;
    toml_table_t *cappExec=NULL, *cappOutput=NULL;

    char errBuf[MAX_ERR_BUFFER] ={0};

    /* Sanity check. */
    if (fileName == NULL || config == NULL) {
        return FALSE;
    }

    if ((fp = fopen(fileName, "r")) == NULL) {
        log_error("Error opening config file: %s: %s\n", fileName,
                  strerror(errno));
        return FALSE;
    }

    /* Parse the TOML file entries. */
    fileData = toml_parse_file(fp, errBuf, sizeof(errBuf));

    fclose(fp);

    if (!fileData) {
        log_error("Error parsing the config file %s: %s\n", fileName, &errBuf[0]);
        return FALSE;
    }

    /* get all mandatory tables for build and capp */
    if (!get_table(fileData, TABLE_BUILD_FROM,    &buildFrom))    goto done;
    if (!get_table(fileData, TABLE_BUILD_COMPILE, &buildCompile)) goto done;
    if (!get_table(fileData, TABLE_CAPP_EXEC,     &cappExec))     goto done;
    if (!get_table(fileData, TABLE_CAPP_OUTPUT,   &cappOutput))   goto done;

    /* non-mandatory. */
    get_table(fileData, TABLE_BUILD_MISC, &buildMisc);
    get_table(fileData, TABLE_BUILD_ROOTFS,  &buildRootfs);
    get_table(fileData, TABLE_BUILD_CONF,    &buildConf);

    ret = read_build_config(config, fileName, buildFrom, buildCompile,
                            buildRootfs, buildConf, buildMisc);
    if (ret == FALSE) goto done;

    ret = read_capp_config(config, fileName, cappExec, cappOutput);
    if (ret == FALSE) goto done;

    ret=TRUE;

done:
    toml_free(fileData);
    return ret;
}

void clear_config(Config *config, int flag) {

    BuildConfig *build;
    CappConfig *capp;

    if (!config) return;

    if ((flag & BUILD_ONLY) && config->build) {

        build = config->build;

        if (build->rootfs)    free(build->rootfs);
        if (build->contained) free(build->contained);

        if (build->version) free(build->version);
        if (build->source)  free(build->source);
        if (build->cmd)     free(build->cmd);
        if (build->binFrom) free(build->binFrom);
        if (build->binTo)   free(build->binTo);

        if (build->mkdir) free(build->mkdir);

        if (build->from)     free(build->from);
        if (build->to)       free(build->to);

        if (build->miscFrom) free(build->miscFrom);
        if (build->miscTo)   free(build->miscTo);

        free(config->build);
    }

    if ((flag & CAPP_ONLY) && config->capp) {

        capp = config->capp;

        if (capp->name)    free(capp->name);
        if (capp->version) free(capp->version);
        if (capp->bin)     free(capp->bin);
        if (capp->path)    free(capp->path);
        if (capp->args)    free(capp->args);
        if (capp->envs)    free(capp->envs);
        if (capp->format)  free(capp->format);

        free(config->capp);
    }
}

void log_config(Config *config) {

    BuildConfig *build=NULL;
    CappConfig *capp=NULL;

    if (config == NULL) return;

    if (config->build) {
        build = config->build;
        log_debug("--- Build Configuration ---");

        log_debug("[FROM:]");
        log_debug("\t rootfs:    %s", build->rootfs);
        log_debug("\t contained: %s", build->contained);

        log_debug("[BUILD:]");
        log_debug("\t version: %s", build->version);
        log_debug("\t static:  %s", (build->staticFlag ? "true": "false"));
        log_debug("\t source:  %s", build->source);
        log_debug("\t command: %s", build->cmd);
        log_debug("\t from:    %s", build->binFrom);
        log_debug("\t to:      %s", build->binTo);

        log_debug("[ROOTFS:]");
        log_debug("\t mkdir: %s", build->mkdir);

        log_debug("[CONF:]");
        log_debug("\t from: %s", build->from);
        log_debug("\t to:   %s", build->to);

        log_debug("[MISC:]");
        log_debug("\t from: %s", build->miscFrom);
        log_debug("\t to:   %s", build->miscTo);
    } else {
        log_debug("No Build configuration found");
    }

    if (config->capp) {
        capp = config->capp;
        log_debug("--- CAPP Configuration ---");

        log_debug("[EXEC:]");
        log_debug("\t name: %s", capp->name);
        log_debug("\t version: %s", capp->version);
        log_debug("\t binary:  %s", capp->bin);
        log_debug("\t path: %s", capp->path);
        if (capp->args) {
            log_debug("\t args: %s", capp->args);
        } else {
            log_debug("\t No args");
        }
        if (capp->envs) {
            log_debug("\t envs: %s", capp->envs);
        } else {
            log_debug("\t No env");
        }

        log_debug("[OUTPUT:]");
        log_debug("\t format: %s", capp->format);
    }
}
