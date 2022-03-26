/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef SYSFS_H_
#define SYSFS_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "property.h"

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

/* 32 byte most of these are int values(4 bytes)
 *  so length wont go beyond it.*/
#define SYS_FILE_MAX_LENGTH			0x0040

/**
 * @fn      int sysfs_erase(char*, uint16_t)
 * @brief   Erase size number of bytes from the sysfs files by setting them
 *             to 0xff
 *
 * @param   name
 * @param   size
 * @return  On success, number of bytes overwritten with 0xff
 *          On failure, -1
 */
int sysfs_erase(char* name, uint16_t size);

/**
 * @fn      int sysfs_exist(char*)
 * @brief   Check if sysfs file exist.
 *
 * @param   name
 * @return  On success, 1
 *          On failure, 0
 */
int sysfs_exist(char* name);

/**
 * @fn      int sysfs_init(char*, void*)
 * @brief   Dummy function so far.
 *
 * @param   name
 * @param   data
 * @return  On success, 0
 *          On failure, non zero value
 */
int sysfs_init(char* name, void* data);

/**
 * @fn      int sysfs_open(char*, int)
 * @brief   Open the sysfs file with name.
 *
 * @param   name
 * @param   flags
 * @return  On success, positive value file descriptor value
 *          On failure, -1
 */
int sysfs_open(char* name, int flags);

/**
 * @fn      int sysfs_read(char*, void*, DataType)
 * @brief   Formatted read operation form the sysfs file.
 *          Read the data from sysfs file and typecast it to type mentioned
 *          by DataType.
 *
 * @param   name
 * @param   data
 * @param   type
 * @return  On success, 0
 *          On failure, -1
 */
int sysfs_read(char* name, void* data, DataType type);

/**
 * @fn      int sysfs_read_block(char*, void*, uint16_t)
 * @brief   raw read operation from the sysfs file. Read size number of bytes
 *          from the sysfs file.
 *
 * @param   name
 * @param   buff
 * @param   size
 * @return  On success, Number of bytes read
 *          On failure, -1
 */
int sysfs_read_block(char* name, void* buff, uint16_t size);

/**
 * @fn      int sysfs_write(char*, void*, DataType)
 * @brief   Formatted write operation to the sysfs file.
 *          Write the data to sysfs file.
 *
 * @param   name
 * @param   data
 * @param   type
 * @return  On success, 0
 *          On failure, -1
 */
int sysfs_write(char* name, void* data, DataType type);

/**
 * @fn      int sysfs_write_block(char*, void*, uint16_t)
 * @brief   Write size number of bytes to the sysfs file.
 *
 * @param   name
 * @param   buff
 * @param   size
 * @return   On success, Number of bytes written
 *          On failure, -1
 */
int sysfs_write_block(char* name, void* buff, uint16_t size);

/**
 * @fn      void sysfs_close(int)
 * @brief   closes the open file file referred by file descriptor.
 *
 * @param   fd
 */
void sysfs_close(int fd);

#ifdef __cplusplus
}
#endif

#endif /* SYSFS_H_ */
