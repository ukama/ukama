/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * pod.h
 */

#ifndef LXCE_POD_H
#define LXCE_POD_H

#include "manifest.h"

#define POD_TYPE_BOOT     "boot"
#define POD_TYPE_SERVICE  "service"
#define POD_TYPE_SHUTDOWN "shutdown"

#define POD_DEFAULT_HOSTNAME "localhost"

#define STACK_SIZE (1024*1024)

#define LXCE_MAX_PATH  256
#define USER_NS_OFFSET 1000

typedef struct _pod {

  int  sockets[2]; /* socket pair between the cInit.d and lxce.d */
  char *type;      /* POD_TYPE_XXX */
  char *hostName;  /* hostname for the Pod */
  
  uid_t uid;       /* uid the cInit will run as */
  gid_t gid;       /* gid the cInit will run as */

  char *mountDir;
} Pod;

int create_ukama_pods(Pod *pods, Manifest *manifest, char *type);

#endif /* LXCE_POD_H */
