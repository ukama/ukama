/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * cspace.h
 */

#ifndef LXCE_CSPACE_H
#define LXCE_CSPACE_H

#include <sys/types.h>
#include "manifest.h"
#include "capp.h"

#define CSPACE_DEFAULT_HOSTNAME "localhost"

/* For parsing the contained configuration file */

#define JSON_UID          "uid"
#define JSON_GID          "gid"
#define JSON_TYPE         "type"
#define JSON_HOSTNAME     "hostname"
#define JSON_NAMESPACES   "namespaces"
#define JSON_CAPABILITIES "capabilities"

#define STACK_SIZE      (1024*1024)
#define CONFIG_MAX_SIZE 1000000

#define LXCE_MAX_PATH  256
#define USER_NS_OFFSET 1000

#define CONTD_MAX_CAPS 20

#define CSPACE_MAX_BUFFER   1024

#define CSPACE_READ_ERROR   1
#define CSPACE_READ_TIMEOUT 2
#define CSPACE_MEMORY_ERROR 3

#define LXCE_SERIAL "serial"

/* Definition of Ukama's contained space as per config file */
typedef struct cSpace_t {

  char *version;      /* contained space version */
  
  char *serial;       /* serial of device, if applicable */
  int  target;        /* Target of this contained space (serial or general) */
  
  char *name;         /* name of the contained space */
  char *hostName;     /* host name associated with space */

  char *mountDir;     /* Mount directory */

  uid_t uid;          /* default uid of space */
  gid_t gid;          /* default gid of space */

  int nameSpaces;     /* linux namespaces enabled in this space */

  int capCount;       /* number of linux capabilities */
  int cap[CONTD_MAX_CAPS]; /* list of capabilities enabled in this space */

  int sockets[2];     /* socket pair */
  char *configFile;   /* Config file - defined in the config.toml */
  
  CApps *apps;        /* Apps associated with this space. */

  struct cSpace_t *next; /* pointer to next contained space */
} CSpace;

int create_cspace(CSpace *space, pid_t *pid);
int process_cspace_config(char *fileName, CSpace *space);

#endif /* LXCE_CSPACE_H */
