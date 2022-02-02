/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
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
#include "cspace.h"
#include "log.h"
#include "capp.h"
#include "utils.h"

/* Local functions */
static int deserialize_capp_config_file(CAppConfig *config, json_t *json);

/*
 * log_config --
 */

static void log_config(CAppConfig *config) {

  int i;

  log_debug("Version: %s", (config->version ? config->version: "NULL"));
  log_debug("Serial:  %s", (config->serial ? config->serial: "NULL"));
  log_debug("hostname: %s", (config->hostName ? config->hostName : "NULL"));

  if (config->process) {
    log_debug("process:", config->process->exec);
    log_debug("\t exec: %s", config->process->exec);

    log_debug("\t argc: %d", config->process->argc);
    for (i=0; i<config->process->argc; i++) {
      log_debug("\t argv[%d]: %s", i, config->process->argv[i]);
    }

    log_debug("\t envc: %d", config->process->envc);
    for (i=0; i<config->process->envc; i++) {
      log_debug("\t env[%d]: %s", i, config->process->env[i]);
    }
  } else {
    log_debug("process: NULL");
  }
}

/*
 * valid_path -- path is valid if there is readable json file present
 *
 */
int valid_path(char *path) {

  char fileName[MAX_BUFFER] = {0};
  struct stat stats;

  if (!path) return FALSE;

  sprintf(fileName, "%s/config.json", path);

  stat(fileName, &stats);

  if (S_ISREG(stats.st_mode)) {
    log_debug("Valid config file at: %s", fileName);
  } else {
    log_error("Default config file not found: %s", fileName);
    return FALSE;
  }

  return TRUE;
}

/*
 * process_capp_config_file --
 *
 */
int process_capp_config_file(CAppConfig *config, char *fileName) {

  int ret;
  FILE *fp=NULL;
  char *buffer=NULL;
  long size=0;
  json_t *json=NULL;
  json_error_t jerror;

  if (!fileName) return FALSE;

  if ((fp = fopen(fileName, "r")) == NULL) {
    log_error("Error opening file: %s Error %s", fileName, strerror(errno));
    return FALSE;
  }

  log_debug("Reading the capp's config.json file: %s", fileName);

  /* Read everything into buffer */
  fseek(fp, 0, SEEK_END);
  size = ftell(fp);
  fseek(fp, 0, SEEK_SET);

  if (size > CONFIG_MAX_SIZE) {
    log_error("Error opening file: %s Error: File size too big: %ld",
              fileName, size);
    fclose(fp);
    return FALSE;
  }

  buffer = (char *)malloc(size+1);
  if (buffer==NULL) {
    log_error("Error allocating memory of size: %ld", size+1);
    fclose(fp);
    return FALSE;
  }
  memset(buffer, 0, size+1);
  fread(buffer, 1, size, fp); /* Read everything into buffer */
  fclose(fp);

  /* Trying loading it as JSON */
  json = json_loads(buffer, 0, &jerror);
  if (json==NULL) {
    log_error("Error loading contd config into JSON format. File: %s Size: %ld",
              fileName, size);
    log_error("JSON error on line: %d: %s", jerror.line, jerror.text);
    goto done;
  } else {
    log_debug("JSON successfully loaded. %s", fileName);
  }

  /* convert into internal structure */
  ret=deserialize_capp_config_file(config, json);
  if (ret) {
    log_config(config);
  } else {
    log_error("Error deser capp config file");
  }

  return ret;

 done:
  if (buffer) free(buffer);
  json_decref(json);

  return FALSE;
}

/*
 * deserialize_capp_config_file -- convert the json into internal struct
 *
 */
static int deserialize_capp_config_file(CAppConfig *config, json_t *json) {

  int j=0, size=0;
  json_t *obj, *jProc;
  json_t *jArray, *jElem;
  char *str=NULL;

  if (config == NULL || json == NULL) return FALSE;

  if (!set_str_object_value(json, &(config->version), JSON_VERSION,
			      TRUE, NULL)) {
    return FALSE;
  }

  if (!set_str_object_value(json, &(config->target), JSON_TARGET, TRUE, NULL)) {
    return FALSE;
  }

  if (config->target == LXCE_SERIAL) {
    if (!set_str_object_value(json, &(config->serial), JSON_SERIAL,
				TRUE, NULL)) {
      return FALSE;
    }
  } else {
    set_str_object_value(json, &(config->serial), JSON_SERIAL, FALSE, NULL);
  }

  set_str_object_value(json, &(config->hostName), JSON_HOSTNAME, FALSE,
		       CAPP_DEFAULT_HOSTNAME);

  /* Look for process info */
  jProc = json_object_get(json, JSON_PROCESS);
  if (jProc != NULL) {
    config->process = (CAppProc *)malloc(sizeof(CAppProc));
    if (!set_str_object_value(jProc, &(config->process->exec), JSON_EXEC,
			      TRUE, NULL)) {
      return FALSE;
    }

    /* Get arguments */
    get_json_arrary_elems(jProc, &(config->process->argc),
			  &(config->process->argv), JSON_ARGS);

    /* Get env variables */
    get_json_arrary_elems(jProc, &(config->process->envc),
			  &(config->process->env), JSON_ENV);

  } else {
    log_error("No valid process info found.");
  }

  /* Look for namespaces. */
  config->nameSpaces = 0;
  jArray = json_object_get(json, JSON_NAMESPACES);
  if (jArray != NULL) {
    size = json_array_size(jArray);

    for (j=0; j<size; j++) {
      jElem = json_array_get(jArray, j);
      if (jElem) {
	obj = json_object_get(jElem, JSON_TYPE);
	if (obj) {
	  str = json_string_value(obj);
	  /* For capp, currently only allow mount, user and pid. */
	  if (strcmp(str, "mount") == 0 ||
	      strcmp(str, "user")  == 0 ||
	      strcmp(str, "pid")   == 0 ) {
	    config->nameSpaces |= namespaces_flag(json_string_value(obj));
	  } else {
	    log_error("Invalid namespace for capp: %s", str);
	  }
	}
      }
    }
  } else {
    log_debug("No valid namespaces found.");
  }

  log_debug("config.json successfully parsed");

  return TRUE;
}
