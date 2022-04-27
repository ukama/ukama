/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SYSFS_H_
#define SYSFS_H_

#include "headers/ubsp/property.h"

#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <errno.h>
#include <stdbool.h>
#include <string.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdint.h>

#define SYS_DEF_OFFSET				0x0000
#define SYS_FILE_MAX_LENGTH			0x0040 /* 32 byte Remember most of these are int values so length wont go beyond it.*/

int sysfs_init(char* name, void* data);
int sysfs_open(char* name, int flags);
void sysfs_close(int fd);
int sysfs_exist(char* name);
int sysfs_erase(char* name, uint16_t size);
int sysfs_read_block(char* name, void* buff, uint16_t size);
int sysfs_write_block(char* name, void* buff, uint16_t size);
int sysfs_read(char* name, void* data, DataType type);
int sysfs_write(char* name, void* data, DataType type);

#endif /* SYSFS_H_ */
