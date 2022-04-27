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
#include "cspace.h"

/*
 * is_valid_capp --
 *
 */
static int is_valid_capp(char *path) {

  char fileName[CAPP_MAX_BUFFER] = {0};
  struct stat stats;

  if (path == NULL) return FALSE;

  sprintf(fileName, "%s/%s", path, DEF_CONFIG);

  stat(fileName, &stats);
  if (S_ISREG(stats.st_mode)) {
    log_debug("Valid config file at: %s", fileName);
    return TRUE;
  } else {
    log_error("Default config file not found: %s", fileName);
  }

  return FALSE;
}

/*
 * capp_unpack -- unpack the capp
 *                Currently using tar eventually should be done via
 *                libarchive (TODO)
 */
static int capp_unpack(char *name, char *tag, char *unpackPath) {

  char runMe[CAPP_MAX_BUFFER] = {0};
  char path[CAPP_MAX_BUFFER]  = {0};
  struct stat stats;

  if (name == NULL || tag == NULL || unpackPath == NULL) return FALSE;

  stat(unpackPath, &stats);
  if (!S_ISDIR(stats.st_mode)) {
    if(mkdir(unpackPath, 0700) < 0) {
      log_error("Error creating default unpack dir: %s Error: %s",
		unpackPath, strerror(errno));
      return FALSE;
    }
  }

  log_debug("Unpacking the capp: %s/%s_%s.tar.gz",  DEF_CAPP_PATH, name, tag);
  sprintf(runMe, "/bin/tar xfz %s/%s_%s.tar.gz -C %s",  DEF_CAPP_PATH,  name,
	  tag, unpackPath);
  if (system(runMe) < 0) {
    log_error("Unable to unpack the capp: %s_%s.tar.gz to %s", name, tag,
	      unpackPath);
    return FALSE;
  }

  /* check to see if config.json exists */
  sprintf(path, "%s/%s_%s/", unpackPath, name, tag);
  if (is_valid_capp(path)) {
    return TRUE;
  }

  return FALSE;
}

/*
 * unpack_all_capps_to_cspace_rootfs --
 *
 */
int unpack_all_capps_to_cspace_rootfs(Manifest *manifest, char *rootfsPath,
				      char *cappPath) {
  int ret=0;
  ArrayElem *ptr=NULL;
  char unpackPath[CSPACE_MAX_BUFFER] = {0};
  char runMe[CSPACE_MAX_BUFFER]      = {0};

if (manifest == NULL) return FALSE;

  /* cleanup existing capps */
  for (ptr = manifest->arrayElem; ptr; ptr=ptr->next) {
    if (!ptr->name || !ptr->tag || !ptr->contained) continue;

    /* existing capps will be at: /capps/rootfs/[contained]/capps/pkgs/unpack
     * delete the whole directory for the cspace
     */
    sprintf(unpackPath, "%s/%s/%s/unpack", rootfsPath, ptr->contained, cappPath);
    sprintf(runMe, "/bin/rm -rf %s", unpackPath);
    log_debug("Running command: %s", runMe);
    if ((ret = system(runMe)) < 0) {
      log_error("Unable to execute cmd %s for space: %s Code: %d", runMe,
		ptr->contained, ret);
      continue;
    }
  }

  /* unpack the capps to: /capps/rootfs/[contained]/capps/pkgs/unpack/
   */
  for (ptr = manifest->arrayElem; ptr; ptr=ptr->next) {

    if (!ptr->name || !ptr->tag || !ptr->contained) continue;

    /* create unpack directory */
    sprintf(unpackPath, "%s/%s/%s/unpack", rootfsPath, ptr->contained, cappPath);
    sprintf(runMe, "/bin/mkdir -p %s", unpackPath);
    log_debug("Running command: %s", runMe);
    if ((ret = system(runMe)) < 0) {
      log_error("Unable to execute cmd %s for space: %s Code: %d", runMe,
		ptr->contained, ret);
      continue;
    }

    /* unpack and verify */
    if (!capp_unpack(ptr->name, ptr->tag, unpackPath)) {
      return FALSE;
    } else {
      /* 4: 2 for /, 1 for _ and 1 for NULL */
      ptr->rootfs = (char *)malloc(strlen(unpackPath) + strlen(ptr->name) +
				   strlen(ptr->tag) + 4);
      /* TODO, hard coded path. FixME */
      sprintf(ptr->rootfs, "/capps/pkgs/unpack/%s_%s", ptr->name, ptr->tag);
    }
  }

  return TRUE;
}
