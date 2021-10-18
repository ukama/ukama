/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * contained.h
 */

#ifndef LXCE_CONTAINED_H
#define LXCE_CONTAINED_H

#include "manifest.h"

#define POD_TYPE_BOOT     "boot"
#define POD_TYPE_SERVICE  "service"
#define POD_TYPE_SHUTDOWN "shutdown"

#define CONTD_DEFAULT_HOSTNAME "localhost"

/* For parsing the contained configuration file */

#define JSON_UID          "uid"
#define JSON_GID          "gid"
#define JSON_TYPE         "type"
#define JSON_HOSTNAME     "hostname"
#define JSON_NAMESPACES   "namespaces"
#define JSON_CAPABILITIES "capabilitles"

#define STACK_SIZE (1024*1024)
#define CONFIG_MAX_SIZE 1000000

#define LXCE_MAX_PATH  256
#define USER_NS_OFFSET 1000

#define CONTD_MAX_CAPS 20

#define LXCE_SERIAL "serial"

/* Definition of Ukama's contained space as per config file */

typedef struct _contdSpace {

  char *version;      /* contained space version */
  
  char *serial;       /* serial of device, if applicable */
  int  target;        /* Target of this contained space (serial or general) */
  
  char *name;         /* name of the contained space */
  char *hostName;     /* host name associated with space */

  uid_t uid;          /* default uid of space */
  gid_t gid;          /* default gid of space */

  int nameSpaces;  /* linux namespaces enabled in this space */

  int capCount;       /* number of linux capabilities */
  int cap[CONTD_MAX_CAPS]; /* list of capabilities enabled in this space */

  struct _contdSpace *next; /* pointer to next contained space */
} ContdSpace; 

typedef struct _pod {

  int  sockets[2]; /* socket pair between the cInit.d and lxce.d */
  char *type;      /* POD_TYPE_XXX */
  char *hostName;  /* hostname for the Pod */
  
  uid_t uid;       /* uid the cInit will run as */
  gid_t gid;       /* gid the cInit will run as */

  char *mountDir;
} Pod;

int create_ukama_pods(Pod *pods, Manifest *manifest, char *type);

#endif /* LXCE_CONTAINED_H */
