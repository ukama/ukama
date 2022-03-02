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
 * @fn char usys_file_read_sym_link*(char*)
 * @brief
 *
 * @param fname
 * @return
 */
char *usys_file_read_sym_link(char *fname);

/**
 * @fn int usys_file_add_record(char*, char*, char*)
 * @brief
 *
 * @param filename
 * @param rowdesc
 * @param data
 * @return
 */
int usys_file_add_record(char *filename, char *rowdesc, char *data);

/**
 * @fn int usys_file_append(void*, void*, off_t, uint16_t)
 * @brief
 *
 * @param fname
 * @param buff
 * @param offset
 * @param size
 * @return
 */
int usys_file_append(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn int usys_file_cleanup(void*)
 * @brief
 *
 * @param fname
 * @return
 */
int usys_file_cleanup(void *fname);

/**
 * @fn int usys_file_erase(void*, off_t, uint16_t)
 * @brief
 *
 * @param fname
 * @param offset
 * @param size
 * @return
 */
int usys_file_erase(void *fname, off_t offset, uint16_t size);

/**
 * @fn int usys_file_exist(char*)
 * @brief
 *
 * @param fname
 * @return
 */
int usys_file_exist(char *fname);

/**
 * @fn int usys_file_init(void*)
 * @brief
 *
 * @param data
 * @return
 */
int usys_file_init(void *data);

/**
 * @fn int usys_file_open(char*, int)
 * @brief
 *
 * @param fname
 * @param flags
 * @return
 */
int usys_file_open(char *fname, int flags);

/**
 * @fn int usys_file_path_exist(char*)
 * @brief
 *
 * @param fname
 * @return
 */
int usys_file_path_exist(char *fname);

/**
 * @fn int usys_file_protect(void*)
 * @brief
 *
 * @param fname
 * @return
 */
int usys_file_protect(void *fname);

/**
 * @fn int usys_file_raw_read(char*, void*, off_t, uint16_t)
 * @brief
 *
 * @param fname
 * @param buff
 * @param offset
 * @param size
 * @return
 */
int usys_file_raw_read(char *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn int usys_file_read(void*, void*, off_t, uint16_t)
 * @brief
 *
 * @param fname
 * @param buff
 * @param offset
 * @param size
 * @return
 */
int usys_file_read(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn int usys_file_read_number(void*, void*, off_t, uint16_t, uint8_t)
 * @brief
 *
 * @param fname
 * @param data
 * @param offset
 * @param count
 * @param size
 * @return
 */
int usys_file_read_number(void *fname, void *data, off_t offset, uint16_t count,
                     uint8_t size);
/**
 * @fn int usys_file_remove(void*)
 * @brief
 *
 * @param data
 * @return
 */
int usys_file_remove(void *data);

/**
 * @fn int usys_file_rename(char*, char*)
 * @brief
 *
 * @param old_name
 * @param new_name
 * @return
 */
int usys_file_rename(char *old_name, char *new_name);

/**
 * @fn int usys_file_symlink_exists(const char*)
 * @brief
 *
 * @param path
 * @return
 */
int usys_file_symlink_exists(const char *path);

/**
 * @fn int usys_file_write(void*, void*, off_t, uint16_t)
 * @brief
 *
 * @param fname
 * @param buff
 * @param offset
 * @param size
 * @return
 */
int usys_file_write(void *fname, void *buff, off_t offset, uint16_t size);

/**
 * @fn int usys_file_write_number(void*, void*, off_t, uint16_t, uint8_t)
 * @brief
 *
 * @param fname
 * @param data
 * @param offset
 * @param count
 * @param size
 * @return
 */
int usys_file_write_number(void *fname, void *data, off_t offset, uint16_t count,
                      uint8_t size);
/**
 * @fn void usys_file_close(int)
 * @brief
 *
 * @param fd
 */
void usys_file_close(int fd);

#ifdef __cplusplus
}
#endif

#endif /* USYS_FILE_H_ */
