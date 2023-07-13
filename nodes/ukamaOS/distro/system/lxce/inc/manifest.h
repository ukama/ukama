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

#define JTAG_VERSION  "version"
#define JTAG_TARGET   "target"
#define JTAG_BOOT     "boot"
#define JTAG_SERVICES "services"
#define JTAG_REBOOT   "reboot"

#define JTAG_NAME     "name"
#define JTAG_TAG       "tag"
#define JTAG_RESTART   "restart"

typedef struct _arrayElem {
    char *name;      /* Name of the cApp */
    char *tag;       /* cApp tag/version */
    char *rootfs;    /* Location where the rootfs is at */
    int  restart;    /* 1: yes, always restart. 0: No */
    int  contained;

    struct _arrayElem *next; /* Next in the list */
} ArrayElem;

typedef struct {

    char *version; /* version of manifest file */
    char *serial;  /* Serial number this config applies (optional) */
    char *target;  /* serial or anyone */

    ArrayElem *boot;
    ArrayElem *services;
    ArrayElem *reboot;
} Manifest;

/* Function headers. */
int process_manifest(Manifest **manifest, char *fileName, void *space);
void get_containers_local_path(Manifest *manifest, Config *config);
void free_manifest(Manifest *manifest);
void copy_capps_to_rootfs(Manifest *manifest);

#endif /* LXCE_MANIFEST_H */
