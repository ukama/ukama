/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef FILE_H_
#define FILE_H_

#include <errno.h>
#include <fcntl.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>

#define MIN_UKDB_OFFSET 0
#define MAX_UKDB_OFFSET 65536

int file_exist(char *fname);
int file_open(char *fname, int flags);
int file_remove(void *fname);
void file_close(int fd);

#endif /* FILE_H_ */
