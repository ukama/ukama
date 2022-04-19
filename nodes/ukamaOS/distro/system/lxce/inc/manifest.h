/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef LXCE_MANIFEST_H
#define LXCE_MANIFEST_H

#include <jansson.h>

#include "lxce_config.h"

#define MANIFEST_MAX_SIZE 1000000 /* 1MB max. file size */

#define JSON_VERSION  "version"
#define JSON_SERIAL   "serial"
#define JSON_TARGET   "target"

#define JSON_CAPP     "ukama-cApp"

/* defines for Ukama Contained App cApp */
#define JSON_NAME      "name"
#define JSON_TAG       "tag"
#define JSON_RESTART   "restart"
#define JSON_CONTAINED "contained"

#define MANIFEST_ALL    "all"
#define MANIFEST_SERIAL "serial"

#define TRUE  1
#define FALSE 0

typedef struct _arrayElem {

  char *name;      /* Name of the cApp */
  char *tag;       /* cApp tag */
  char *contained; /* where this app is contained (boot, service, shutdown) */
  char *rootfs;
  int  restart;    /* 1: yes, always restart. 0: No */

  struct _arrayElem *next; /* Next in the list */
} ArrayElem;

typedef struct {

  char *version; /* version of manifest file */
  char *serial;  /* Serial number this config applies (optional) */
  char *target;  /* serial or anyone */

  ArrayElem *arrayElem;  /* cApps array elements */
} Manifest;

/* Function headers. */
int process_manifest(char *fileName, Manifest *manifest, void *space);
void get_containers_local_path(Manifest *manifest, Config *config);
void clear_manifest(Manifest *manifest);
void copy_capps_to_cspace_rootfs(Manifest *manifest, char *sPath, char *dPath);

#endif /* LXCE_MANIFEST_H */
