/**
 * Copyright (c) 2022-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef CONFIG_H
#define CONFIG_H

/* Various tables */
#define TABLE_BUILD         "build"
#define TABLE_CAPP          "capp"
#define TABLE_BUILD_FROM    "build-from"
#define TABLE_BUILD_COMPILE "build-compile"
#define TABLE_BUILD_ROOTFS  "build-rootfs"
#define TABLE_BUILD_CONF    "build-conf"
#define TABLE_CAPP_EXEC     "capp-exec"

/* Keys for various table */
#define KEY_BASE        "base"
#define KEY_VERSION     "version"
#define KEY_STATIC      "static"
#define KEY_SOURCE      "source"
#define KEY_CMD         "cmd"
#define KEY_BIN_FROM    "bin_from"
#define KEY_BIN_TO      "bin_to"
#define KEY_MKDIR       "mkdir"
#define KEY_FROM        "from"
#define KEY_TO          "to"
#define KEY_FORMAT      "format"
#define KEY_EXEC        "exec"
#define KEY_PATH        "path"
#define KEY_ARGS        "args"
#define KEY_ENVS        "envs"
#define KEY_NAME        "name"
#define KEY_BIN         "bin"
#define KEY_AUTOSTART   "autostart"
#define KEY_AUTORESTART "autorestart"
#define KEY_DEPENDS_ON  "depends_on"
#define KEY_WAIT_FOR    "wait_for"

#define VALUE_YES       "yes"
#define VALUE_NO        "no"

#define DATUM_BOOL      0x01
#define DATUM_STRING    0x02
#define DATUM_MANDATORY 0x04

#define BUILD_ONLY      0x08
#define CAPP_ONLY       0x16

#define MAX_BUFFER        1024
#define MAX_ERROR_BUFFER  1024

#define TRUE  1
#define FALSE 0

typedef struct build_config_t {

	/* from */
	char *baseImage;
	char *baseVersion;

	/* build */
	char *version;
	int  staticFlag;
	char *source;
	char *cmd;
	char *binFrom;
	char *binTo;

	/* rootfs */
	char *mkdir;

	/* conf */
	char *from;
	char *to;
} BuildConfig;

typedef struct capp_config_t {

	char *name;       /* Name of the capp */
	char *version;    /* Version of the capp */
	char *bin;        /* capp binary */
	char *path;       /* Absolute path to the bin */
	char *args;       /* Runtime arguments, if any */
	char *envs;       /* Environment variables, if any */
	char *dependsOn;  /* wait on program(s) to finish exec */
	char *waitFor;    /* time to wait before executing */
	int  autostart;   /* autostart for supervisor.d */
	int  autorestart; /* autorestart for supervisor.d */
} CappConfig;

typedef struct config_t {

	CappConfig  *capp;
	BuildConfig *build;
} Config;

typedef struct configs_t {

	char   *fileName;
	int    valid;
	char   *errorStr;
	Config *config;
  
	struct configs_t *next;
} Configs;

/* Function headers */
int read_config_files(Configs **configs, char *configDir);
void free_configs(Configs *configs);
void free_config(Config *config, int flag);
void log_config(Config *config);

#endif /* CONFIG_H */
