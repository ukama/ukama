/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
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
