/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>

#include "config_app.h"
#include "log_app.h"

#define SCRIPT "./scripts/make-app.sh"

#define MAX_BUFFER   1024
#define DEF_HOSTNAME "localhost"
#define DEF_CONFIG   "config.json"

#define JSON "{ \n \t \"version\": \"%s\", \n \t \"target\": \"all\", \n \t \"process\": { \n \t \t \"exec\": \"%s/%s\", \n \t \t \"args\": \"%s\", \n \t \t \"env\": \"%s\" \n \t }, \n \t \"hostname\": \"%s\",\n \t \"namespaces\" : [ \n \t \t { \"type\" : \"pid\"}, \n \t {\t \"type\" : \"mount\"}, \n \t { \t \"type\" : \"user\"} \n \n \t ] \n } \n"

static int create_app_config(Config *config) {

  char str[2048] = {0};
  char *envs, *args;
  FILE *fp=NULL;
  int ret=TRUE;

  remove(DEF_CONFIG);

  fp = fopen(DEF_CONFIG, "a+");

  if (fp == NULL) {
    log_error("Error opening config file: config.json %s", strerror(errno));
    return FALSE;
  }

  if (config->capp->envs) {
    envs = config->capp->envs;
  } else {
    envs = "";
  }

  if (config->capp->args) {
    args = config->capp->args;
  } else {
    args = "";
  }

  /* version & target */
  sprintf(str, JSON, config->build->version,
	  config->capp->path, config->capp->bin,
	  args,
	  envs,
	  DEF_HOSTNAME);

  log_debug("config.json: \n %s", str);
  
  if (fwrite(str,1,strlen(str)+1,fp)<=0) {
    log_error("Error writing to json file: %s", strerror(errno));
    ret=FALSE;
    goto done;
  }
  fflush(fp);

 done:
  fclose(fp);
  if (!ret) remove(DEF_CONFIG);
  return ret;
}

int create_app(Config *config) {

  char runMe[MAX_BUFFER] = {0};

  /* Flow is:
   * 1. create config.json and mv to rootfs
   * 2. rename the rootfs to match app
   * 3. tar things up and clean up
   */
  if (!create_app_config(config)) {
      log_error("Error creating %s", DEF_CONFIG);
      return FALSE;
  }

  /* Copy to / in rootfs and delete it locally */
  sprintf(runMe, "%s cp %s %s_%s/", SCRIPT, DEF_CONFIG,
          config->capp->name, config->capp->version);
  if (system(runMe) < 0) {
      log_error("Error copying %s to rootfs", DEF_CONFIG);
      return FALSE;
  }
  remove(DEF_CONFIG);

  /* delete the directory afterwards */
  sprintf(runMe, "%s pack %s %s_%s.tar.gz %s_%s %d",
          SCRIPT,
          getenv("UKAMA_ROOT"),
          config->capp->name, config->capp->version,
          config->capp->name, config->capp->version, TRUE);
  if (system(runMe) < 0) {
      log_error("Error packing the capp to %s_%s", config->capp->name,
                config->capp->version);
      return FALSE;
  }

  /* clean up */
  sprintf(runMe, "%s clean %s_%s", SCRIPT,
          config->capp->name, config->capp->version);
  if (system(runMe) < 0) return FALSE;

  return TRUE;
}
