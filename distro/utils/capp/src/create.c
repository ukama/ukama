/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/* to build capp */

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <errno.h>
#include <string.h>

#include "config.h"
#include "log.h"

#define SCRIPT "./mk_capp_rootfs.sh"

#define MAX_BUFFER   1024
#define DEF_HOSTNAME "localhost"
#define DEF_CONFIG   "config.json"

#define JSON "{ \n \t \"version\": \"%s\", \n \t \"target\": \"all\", \n \t \"process\": { \n \t \t \"exec\": \" %s/%s \", \n \t \t \"args\": [ \"%s\" ], \n \t \t \"env\": [ \"%s\" ], \n \t } \n \t \"hostname\": \"%s\",\n \t \"namespace\" : [ \n \t \t { \"type\" : \"pid\", \n \t \t \"type\" : \"ipc\", \n \t \t \"type\" : \"mount\", \n \t \t \"type\" : \"user\", \n \t } \n \t ] \n ] \n"

/*
 * create_capp_config --
 *
 */
static int create_capp_config(Config *config) {

  char str[2048] = {0};
  char runMe[MAX_BUFFER] = {0};
  int fd, ret=TRUE;

  if ((fd = open(DEF_CONFIG, O_CREAT | O_TRUNC,S_IRWXU)) < 0) {
    log_error("Error opening config file: config.json %s", strerror(errno));
    return FALSE;
  }

  /* version & target */
  sprintf(str, JSON, config->build->version,
	  config->capp->path, config->capp->bin,
	  config->capp->args,
	  config->capp->envs,
	  DEF_HOSTNAME);

  log_debug("config.json: \n %s", str);
  
  write(fd, str, strlen(str)+1);
  close(fd);

  /* Copy to / in rootfs and delete it locally */
  sprintf(runMe, "%s cp %s /", SCRIPT, DEF_CONFIG);
  if (system(runMe) < 0) {
    log_error("Error copying %s to rootfs", DEF_CONFIG);
    ret = FALSE;
    goto cleanup;
  }

 cleanup:
  remove(DEF_CONFIG);
  return ret;
}

/*
 * create_capp --
 *
 */
int create_capp(Config *config) {

  char runMe[MAX_BUFFER] = {0};

  /* Flow is:
   * 1. create config.json and mv to rootfs
   * 2. rename the rootfs to match capp
   * 3. tar things up and clean up
   */
  if (!create_capp_config(config)) {
    log_error("Error creating %s", DEF_CONFIG);
    return FALSE;
  }

  sprintf(runMe, "%s rename %s_%s", SCRIPT, config->capp->name,
	  config->capp->version);
  if (system(runMe) < 0) {
    log_error("Error renaming the dir to %s_%s", config->capp->name,
	      config->capp->version);
    return FALSE;
  }

  /* delete the directory afterwards */
  sprintf(runMe, "%s pack %s_%s.tar.gz %d", SCRIPT, config->capp->name,
	  config->capp->version, TRUE);

  return TRUE;
}
