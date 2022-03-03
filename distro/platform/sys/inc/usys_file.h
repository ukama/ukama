/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_FILE_H_
#define USYS_FILE_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     char usys_file_read_sym_link*(char*)
 * @brief  Read symbolink link pathname.
 *
 * @param  fname
 * @return On success char*
 *         On error NULL
 */
char *usys_file_read_sym_link(char *fname);

/**
 * @fn     int usys_file_add_record(char*, char*, char*)
 * @brief  Add a record to a file if it exist if not, create a new file.
 *
 * @param  filename
 * @param  rowdesc
 * @param  data
 * @return On success bytes of the data written on success
 *         On error -1
 */
int usys_file_add_record(char *filename, char *rowdesc, char *data);

/**
 * @fn     int usys_file_append(void*, void*, off_t, uint16_t)
 * @brief  Append a data to end of file if exist otherwise error.
 *
 * @param  fname
 * @param  buff
 * @param  offset
 * @param  size
 * @return On success bytes of the data written on success
 *         On error -1
 */
int usys_file_append(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn     int usys_file_cleanup(void*)
 * @brief  remove file.
 *
 * @param  fname
 * @return On success 0
 *         On error -1
 */
int usys_file_cleanup(void *fname);

/**
 * @fn     int usys_file_erase(void*, off_t, uint16_t)
 * @brief  Erase is meant for eeprom devices where all the data in file
 *         is set to 0xFF
 *
 * @param  fname
 * @param  offset
 * @param  size
 * @return On success number of bytes erased.
 *         On error -1
 */
int usys_file_erase(void *fname, off_t offset, uint16_t size);

/**
 * @fn     int usys_file_exist(char*)
 * @brief  Check if file exsit and if it's a regular file
 *
 * @param  fname
 * @return On success 1
 *         Otherwise  0
 */
int usys_file_exist(char *fname);

/**
 * @fn     int usys_file_init(void*)
 * @brief  Check if file exist ot not. If not then create a new file.
 *
 * @param  data
 * @return On success 0
 *         On failure -1
 */
int usys_file_init(void *data);

/**
 * @fn     int usys_file_open(char*, int)
 * @brief  Opens a file with the flags provided.
 *
 * @param  fname
 * @param  flags
 * @return On success file descriptor
 *         On failure -1
 */
int usys_file_open(char *fname, int flags);

/**
 * @fn     int usys_file_path_exist(char*)
 * @brief  Check if file exist
 *
 * @param  fname
 * @return On success 1
 *         Otherwise  0
 */
int usys_file_path_exist(char *fname);

/**
 * @fn     int usys_file_protect(void*)
 * @brief  TBU
 *
 * @param  fname
 * @return
 */
int usys_file_protect(void *fname);

/**
 * @fn     int usys_file_raw_read(char*, void*, off_t, uint16_t)
 * @brief  This API is meant to read data from eeprom
 *
 * @param  fname
 * @param  buff
 * @param  offset
 * @param  size
 * @return On success number of bytes read
 *         On error -1
 */
int usys_file_raw_read(char *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn     int usys_file_read(void*, void*, off_t, uint16_t)
 * @brief  Read the data from file.
 *
 * @param  fname
 * @param  buff
 * @param  offset
 * @param  size
 * @return On success number of bytes read
 *         On error -1
 */
int usys_file_read(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn     int usys_file_read_number(void*, void*, off_t, uint16_t, uint8_t)
 * @brief  Read a number from file. Used to read data from sysfs files.
 *
 * @param  fname
 * @param  data
 * @param  offset
 * @param  count
 * @param  size
 * @return On success 0
 *         On error -1
 */
int usys_file_read_number(void *fname, void *data, off_t offset, uint16_t count,
                          uint8_t size);
/**
 * @fn     int usys_file_remove(void*)
 * @brief  delete a file
 *
 * @param  data
 * @return On success 0
 *         On error -1
 */
int usys_file_remove(void *data);

/**
 * @fn     int usys_file_rename(char*, char*)
 * @brief  Rename a file.
 *
 * @param  old_name
 * @param  new_name
 * @return On success 0
 *         On error -1
 */
int usys_file_rename(char *old_name, char *new_name);

/**
 * @fn     int usys_file_symlink_exists(const char*)
 * @brief  Check if the path provided is sunbolic link.
 *
 * @param  path
 * @return On success 1
 *         On error 0
 */
int usys_file_symlink_exists(const char *path);

/**
 * @fn     int usys_file_write(void*, void*, off_t, uint16_t)
 * @brief  Write data to file.
 *
 * @param  fname
 * @param  buff
 * @param  offset
 * @param  size
 * @return On success number of bytes written.
 *         On failure -1
 */
int usys_file_write(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn     int usys_file_write_number(void*, void*, off_t, uint16_t, uint8_t)
 * @brief  Write a numbers to a file. used for sysfs files.
 *
 * @param  fname
 * @param  data
 * @param  offset
 * @param  count
 * @param  size
 * @return On success 0
 *         On error -1
 */
int usys_file_write_number(void *fname, void *data, off_t offset,
                           uint16_t count, uint8_t size);
/**
 * @fn     void usys_file_close(int)
 * @brief  close a file
 *
 * @param  fd
 */
void usys_file_close(int fd);

#ifdef __cplusplus
}
#endif

#endif /* USYS_FILE_H_ */
