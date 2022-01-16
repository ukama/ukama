/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <string.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/types.h>

#include "capp.h"
#include "log.h"

/*
 * is_valid_capp --
 *
 */
static int is_valid_capp(char *path) {

  char fileName[CAPP_MAX_BUFFER] = {0};

  if (path==NULL) return FALSE;

  sprintf(fileName, "%s/%s", path, DEF_CONFIG);

  if (access(fileName, F_OK)) {
    return TRUE;
  }

  return FALSE;
}

/*
 * capp_unpack -- unpack the capp to default location.
 *                Currently using tar eventually should be done via
 *                libarchieve (TODO)
 */
int capp_unpack(char *name, char *tag, char **dest) {

  char runMe[CAPP_MAX_BUFFER] = {0};
  char path[CAPP_MAX_BUFFER] = {0};
  struct stat stats;
  
  if (name == NULL || tag == NULL) return FALSE;

  /* Check if directory exist or not */
  stat(DEF_CAPP_UNPACK_PATH, &stats);
  if (!S_ISDIR(stats.st_mode)) {
    if(mkdir(DEF_CAPP_UNPACK_PATH, 0700) < 0) {
      log_error("Error creating default unpack dir: %s Error: %s",
		DEF_CAPP_UNPACK_PATH, strerror(errno));
      return FALSE;
    }
  }
  
  sprintf(runMe, "tar xfz %s_%s.tar.gz -C %s", name, tag,
	  DEF_CAPP_UNPACK_PATH);
  if (system(runMe) < 0) {
    log_error("Unable to unpack the capp: %s_%s.tar.gz to %s", name, tag,
	      DEF_CAPP_UNPACK_PATH);
    return FALSE;
  }

  /* check to see if config.json exists */
  sprintf(path, "%s/%s_%s/", DEF_CAPP_UNPACK_PATH, name, tag);
  if (is_valid_capp(path)) {
    *dest = strdup(path);
    return TRUE;
  }

  *dest = NULL;
  return FALSE;
}
