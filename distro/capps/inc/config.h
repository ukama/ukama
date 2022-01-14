/*
 * config.h
 */

#ifndef CONFIG_H
#define CONFIG_H

/* Various table defs */
#define TABLE_BUILD         "build"
#define TABLE_CAPP          "capp"
#define TABLE_BUILD_FROM    "build-from"
#define TABLE_BUILD_COMPILE "build-compile"
#define TABLE_BUILD_ROOTFS  "build-rootfs"
#define TABLE_BUILD_CONF    "build-conf"
#define TABLE_CAPP_EXEC     "capp-exec"
#define TABLE_CAPP_OUTPUT   "capp-output"

/* Keys for various table */
#define KEY_ROOTFS    "rootfs"
#define KEY_CONTAINED "contained"
#define KEY_VERSION   "version"
#define KEY_STATIC    "static"
#define KEY_SOURCE    "source"
#define KEY_CMD       "cmd"
#define KEY_BIN_FROM  "bin_from"
#define KEY_BIN_TO    "bin_to"
#define KEY_MKDIR     "mkdir"
#define KEY_FROM      "from"
#define KEY_TO        "to"
#define KEY_FORMAT    "format"
#define KEY_EXEC      "exec"
#define KEY_PATH      "path"
#define KEY_ARGS      "args"
#define KEY_ENVS      "envs"
#define KEY_NAME      "name"
#define KEY_BIN       "bin"

#define DATUM_BOOL      0x01
#define DATUM_STRING    0x02
#define DATUM_MANDATORY 0x04

#define BUILD_ONLY 0x01
#define CAPP_ONLY  0x02

#define MAX_ERR_BUFFER 256

#define TRUE  1
#define FALSE 0

typedef struct build_config_t {
  
  /* from */
  char *rootfs;
  char *contained;

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
}BuildConfig;

typedef struct capp_config_t {

  char *name;   /* Name of the capp */
  char *version;/* Version of the capp */
  char *bin;    /* capp binary */
  char *path;   /* Absolute path to the bin */
  char *args;   /* Runtime arguments, if any */
  char *envs;   /* Environment variables, if any */
  char *format; /* Output format */
}CappConfig;

typedef struct config_t {

  CappConfig  *capp;
  BuildConfig *build;
} Config;

/* Function headers */
int read_config_file(Config *config, char *fileName);
void clear_config(Config *config, int flag);
void log_config(Config *config);

#endif /* CONFIG_H */
