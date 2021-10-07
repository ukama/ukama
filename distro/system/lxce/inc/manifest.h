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

#define MANIFEST_MAX_SIZE 1000000 /* 1MB */

#define JSON_VERSION  "version"
#define JSON_SERIAL   "serial"
#define JSON_TARGET   "target"

#define JSON_BOOT     "boot"
#define JSON_SERVICE  "service"
#define JSON_SHUTDOWN "shutdown"

#define JSON_ORDER    "order"
#define JSON_NAME     "name"
#define JSON_TAG      "tag"
#define JSON_ACTIVE   "active"
#define JSON_RESTART  "restart"

#define MANIFEST_ALL    "all"
#define MANIFEST_SERIAL "serial"

#define TRUE  1
#define FALSE 0

typedef struct _container {

  int  order;    /* Start order */
  char *name;    /* Name of the container */
  char *tag;     /* container tag */
  int  active;   /* 1: yes, start it. 0: skip it */
  int  restart;  /* 1: yes, always restart. 0: No */

  char *path;    /* local path as per WIMC. */

  struct _container *next; /* Next in the list. */
} Container;

typedef struct {

  char *version; /* version of manifest file */
  char *serial;  /* Serial number this config applies (optional) */
  char *target;  /* serial or anyone. */

  Container *boot;     /* Container cfg to start upon booting */
  Container *service;  /* Container cfg post boot */
  Container *shutdown; /* Contaienrs when unit is being shutdown */
} Manifest;

/* Function headers. */
int process_manifest(char *fileName, Manifest *manifest);
void get_containers_local_path(Manifest *manifest, Config *config);
void clear_manifest(Manifest *manifest);

#endif /* LXCE_MANIFEST_H */
