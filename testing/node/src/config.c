/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Config.c
 *
 */

#include <stdlib.h>
#include <string.h>
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

/*
 * read_entry --
 *
 */
static int read_entry(toml_table_t *table, char *key, char **destStr,
					  int *destInt, int flag) {

	toml_datum_t datum;
	int ret=TRUE;

	/* sanity check */
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
				ret = FALSE;
			}
		} else if (flag & DATUM_STRING) {
			*destStr = strdup(datum.u.s);
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

/*
 * read_capp_table --
 *
 */
static int read_capp_table(toml_table_t *table, Config *config,
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

	if (strcmp(type, TABLE_CAPP_EXEC)==0) { /* [capp-exec] */
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
	} else {
		return FALSE;
	}

	return TRUE;
}

/*
 * read_build_table --
 *
 */
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

	if (strcmp(type, TABLE_BUILD_FROM)==0) { /* [build-from] */
		if (!read_entry(table, KEY_BASE, &build->baseImage, NULL,
						DATUM_STRING | DATUM_MANDATORY)) {
			return FALSE;
		}

		if (!read_entry(table, KEY_VERSION, &build->baseVersion, NULL,
						DATUM_STRING | DATUM_MANDATORY)) {
			return FALSE;
		}
	} else if (strcmp(type, TABLE_BUILD_COMPILE)==0) { /* [build-compile] */
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
	} else if (strcmp(type, TABLE_BUILD_ROOTFS)==0) { /* [build-rootfs] */
		if (!read_entry(table, KEY_MKDIR, &build->mkdir, NULL, DATUM_STRING)) {
			return FALSE;
		}
	} else if (strcmp(type, TABLE_BUILD_CONF)==0) { /* [build-conf] */
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

/*
 * get_table --
 *
 */
static int get_table(toml_table_t *src, char *key, toml_table_t **table) {

	*table = toml_table_in(src, key);
	if (*table == NULL) {
		log_error("[%s] section parsing error in config file", key);
		return FALSE;
	}

	return TRUE;
}

/*
 * read_build_config --
 *
 */
static int read_build_config(Config *config, char *fileName,
							 toml_table_t *buildFrom,
							 toml_table_t *buildCompile,
							 toml_table_t *buildRootfs,
							 toml_table_t *buildConf) {

	int ret=TRUE;

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

	if (read_build_table(buildRootfs, config, TABLE_BUILD_ROOTFS) == FALSE) {
	    log_debug("[%s] section parsing error in config file: %s\n",
				  TABLE_BUILD_ROOTFS, fileName);
	}

	if (read_build_table(buildConf, config, TABLE_BUILD_CONF) == FALSE) {
		log_debug("[%s] section parsing error in config file: %s\n",
				  TABLE_BUILD_CONF, fileName);
	}

 done:
	return ret;
}

/*
 * read_capp_config --
 *
 */
static int read_capp_config(Config *config, char *fileName,
							toml_table_t *cappExec) {

	int ret=TRUE;

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

/*
 * add_to_configs --
 *
 */
static int add_to_configs(Configs **configs, Config *config, char *fileName,
						  int status, char *errorStr) {

	Configs *ptr=NULL;

	if ((*configs) == NULL) {
		*configs = (Configs *)calloc(1, sizeof(Configs));
		if (*configs == NULL) {
			log_error("Error allocating memory of size: %lu", sizeof(Configs));
			return FALSE;
		}
		ptr = *configs;
	} else {
		for (ptr=(*configs); ptr->next; ptr=ptr->next);
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

/*
 * is_valid_file --
 *
 */
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
	
/*
 * read_config_files -- read all the config files within the dir
 *
 */
int read_config_files(Configs **configs, char *configDir) {

    int ret=TRUE, configStatus=FALSE;
	struct stat statBuf;
	struct dirent *dp=NULL;
	DIR *dir=NULL;
	char *configFile=NULL, *errorStr=NULL;
	Config *config=NULL;
	char buffer[MAX_BUFFER] = {0};

	if (configs == NULL || configDir == NULL) return FALSE;

	/* Check to see if the configDir is a valid one. */
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
	    sprintf(buffer , "%s/%s", configDir, dp->d_name);
	    configFile = realpath(buffer, NULL);
		if (!is_valid_file(configFile)) {
		  free(configFile);
		  continue;
		}

		config = (Config *)calloc(1, sizeof(Config));
		if (config == NULL) {
			log_error("Error allocating memory of size: %lu", sizeof(Config));
			goto failure;
		}

		if (!read_config_file(config, configFile, &errorStr)) {
			log_error("Parsing error for: %s", configFile);
			free_config(config, BUILD_ONLY | CAPP_ONLY);
			free(config);
			ret = FALSE;
			configStatus = FALSE;
		} else {
			configStatus = TRUE;
			log_config(config);
		}

		add_to_configs(configs, config, configFile, configStatus, errorStr);

		free(errorStr);
		free(configFile);
	}

	closedir(dir);
	return ret;
    
 failure:
	free_configs(*configs);
	if (config) free_config(config, BUILD_ONLY | CAPP_ONLY);
	closedir(dir);

	return FALSE;
}

/*
 * read_config_file -- read and parse the config file
 *
 *
 */
static int read_config_file(Config *config, char *fileName, char **error) {

	int ret=FALSE;
	FILE *fp=NULL;

	toml_table_t *fileData=NULL;
	toml_table_t *buildFrom=NULL, *buildCompile=NULL, *buildRootfs=NULL;
	toml_table_t *buildConf=NULL;
	toml_table_t *cappExec=NULL;

	/* Sanity check. */
	if (fileName == NULL || config == NULL) {
		return FALSE;
	}

	if ((fp = fopen(fileName, "r")) == NULL) {
		log_error("Error opening config file: %s: %s\n", fileName,
				  strerror(errno));
		return FALSE;
	}

	/* Allocate memory for error buffer. This will be freed if no error 
	 * was found while parsing the toml file 
	 */
	*error = (char *)calloc(1, MAX_ERROR_BUFFER);
	if (*error == NULL) {
		log_error("Error allocating memory of size: %lu", MAX_ERROR_BUFFER);
		return FALSE;
	}

	/* Parse the TOML file entries. */
	fileData = toml_parse_file(fp, *error, MAX_ERROR_BUFFER);
	fclose(fp);

	if (!fileData) {
		log_error("Error parsing the config file %s: %s\n", fileName, *error);
		return FALSE;
	}

	/* get all mandatory tables for build and capp */
	if (!get_table(fileData, TABLE_BUILD_FROM,    &buildFrom))    goto done;
	if (!get_table(fileData, TABLE_BUILD_COMPILE, &buildCompile)) goto done;
	if (!get_table(fileData, TABLE_CAPP_EXEC,     &cappExec))     goto done;

	/* optional table */
	get_table(fileData, TABLE_BUILD_ROOTFS,  &buildRootfs);
	get_table(fileData, TABLE_BUILD_CONF,    &buildConf);

	ret = read_build_config(config, fileName, buildFrom, buildCompile,
							buildRootfs, buildConf);
	if (ret == FALSE) goto done;

	ret = read_capp_config(config, fileName, cappExec);
	if (ret == FALSE) goto done;

	ret=TRUE;

 done:
	toml_free(fileData);
	return ret;
}

/*
 * free_configs --
 *
 */
void free_configs(Configs *configs) {

	Configs *ptr=NULL, *next=NULL;

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

/*
 * free_config --
 */
void free_config(Config *config, int flag) {

	BuildConfig *build;
	CappConfig *capp;

	if (!config) return;

	if ((flag & BUILD_ONLY) && config->build) {

		build = config->build;

		if (build->baseImage)    free(build->baseImage);
		if (build->baseVersion)  free(build->baseVersion);

		if (build->version) free(build->version);
		if (build->source)  free(build->source);
		if (build->cmd)     free(build->cmd);
		if (build->binFrom) free(build->binFrom);
		if (build->binTo)   free(build->binTo);
		if (build->mkdir)   free(build->mkdir);
		if (build->from)    free(build->from);
		if (build->to)      free(build->to);

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
		if (capp->waitFor) free(capp->waitFor);
		if (capp->dependsOn) free(capp->dependsOn);

		free(config->capp);
	}

	if ((flag & CAPP_ONLY) && (flag & BUILD_ONLY)) {
	    free(config);
	}
}

/*
 * log_config --
 *
 */
void log_config(Config *config) {

	BuildConfig *build=NULL;
	CappConfig *capp=NULL;

	if (config == NULL) return;

	if (config->build) {
		build = config->build;
		log_debug("--- CAPP Build Configuration ---");

		log_debug("[FROM:]");
		log_debug("\t base:    %s", build->baseImage);
		log_debug("\t version: %s", build->baseVersion);

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

	} else {
		log_debug("No Build configuration found");
	}

	if (config->capp) {
		capp = config->capp;
		log_debug("--- CAPP Run Configuration ---");

		log_debug("[EXEC:]");
		log_debug("\t name: %s", capp->name);
		log_debug("\t version: %s", capp->version);
		log_debug("\t autostart:  %s", (capp->autostart ? "true": "false"));
		log_debug("\t autorestart:  %s", (capp->autorestart ? "true": "false"));
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
		if (capp->dependsOn) {
			log_debug("\t dependsOn: %s", capp->dependsOn);
		}
		if (capp->waitFor) {
			log_debug("\t waitFor: %s", capp->waitFor);
		}
	}
}
