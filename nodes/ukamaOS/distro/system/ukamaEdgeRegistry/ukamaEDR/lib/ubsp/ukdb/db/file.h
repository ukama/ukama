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
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

#define MIN_UKDB_OFFSET 0
#define MAX_UKDB_OFFSET 65536

int file_init(void* fname);
int file_cleanup(void* fname);
int file_exist(char *fname);
int file_erase(void* fname, off_t offset, uint16_t size);
int file_open(char* fname, int flags);
int file_protect(void* fname);
char* file_read_sym_link(char* fname );
int file_raw_read(char* fname, void* buff, off_t offset, uint16_t size);
int file_read(void* fname, void* buff, off_t offset, uint16_t size);
int file_rename(char* old_name, char* new_name);
int file_write(void* fname, void* buff, off_t offset, uint16_t size) ;
int file_read_number(void* fname, void* value, off_t offset, uint16_t count, uint8_t size);
int file_write_number(void* fname, void* value, off_t offset, uint16_t count, uint8_t size);
int file_remove(void* fname);
void file_close();

#endif /*FILE_H_*/

