/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * space.h
 */

#ifndef LXCE_SPACE_H
#define LXCE_SPACE_H

#define SPACE_STACK_SIZE (1024*1024)

#define AREA_TYPE_CSPACE "cSpace"
#define AREA_TYPE_CAPP   "cApp"

int setup_mounts(char *areaType, char *rootfs, char *name);
int create_space(char *areaType, int sockets[2], int namespaces,
		 char *name, pid_t *pid, int (*func)(void *), void *arg);

#endif /* LXCE_SPACE_H */
