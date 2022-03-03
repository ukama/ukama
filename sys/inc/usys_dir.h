/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_DIR_H_
#define USYS_DIR_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     int usys_mkdir(const char*, mode_t)
 * @brief  Create the directory named pathname.
 *
 * @param  pathname
 * @param  mode
 * @return On success, 0
 *         On error, -1
 */
static inline int usys_mkdir(const char *pathname, mode_t mode) {
	return mkdir(pathname, mode);
}

/**
 * @fn     DIR opendir*(const char*)
 * @brief  function opens a directory stream corresponding to the directory
 *         name, and returns a pointer to the directory stream.
 *
 * @param  name
 * @return On success, pointer to a directory stream.
 * 		  On error, NULL is returned, and errno is set to indicate the error.
 */
static inline DIR *usys_opendir(const char *name) {
	return opendir(name);
}

/**
 * @fn     struct dirent readdir*(DIR*)
 * @brief  returns a pointer to a dirent structure representing the next
 *         directory entry in the directory stream pointed to by dirp.
 *
 * @param  dirp
 * @return On success, readdir() returns a pointer to a dirent structure.
 *         If an error occurs, NULL is returned and errno is set to
 *         indicate the error.
 */
static inline struct dirent *usys_readdir(DIR *dirp) {
	return readdir(dirp);
}

/**
 * @fn     int usys_rmdir(const char*)
 * @brief  remove empty directories
 *
 * @param  pathname
 * @return On success, 0.
 *         On error, -1
 */
static inline int usys_rmdir(const char *pathname) {
	return rmdir(pathname);
}

/**
 * @fn     int closedir(DIR*)
 * @brief  closes the directory stream associated with dirp
 *
 * @param  dirp
 * @return On success, 0
 * 		   On error, -1
 */
static inline int usys_closedir(DIR *dirp) {
	return closedir(dirp);
}

/**
 * @fn     char usys_getcwd*(char*, size_t)
 * @brief  These functions return a null-terminated string containing an
 *         absolute pathname that is the current working directory of the
 *         calling process.
 *
 * @param  buf
 * @param  size
 * @return On success, these functions return a pointer to a string
 *         containing the pathname of the current working directory.
 *         On failure, these functions return NULL, and errno is set to
 *         indicate the error.
 */
static inline char *usys_getcwd(char *buf, size_t size) {
	return getcwd(buf, size);
}

/**
 * @fn    void usys_rewinddir(DIR*)
 * @brief resets the position of the directory stream dirp to the beginning
 *        of the directory
 *
 * @param dirp
 */
static inline void usys_rewinddir(DIR *dirp) {
	rewinddir(dirp);
}

/**
 * @fn     void usys_seekdir(DIR*, long)
 * @brief  set the position of the next readdir() call in the directory stream.
 *
 * @param  dirp
 * @param  loc
 */
static inline void usys_seekdir(DIR *dirp, long loc) {
	seekdir(dirp, loc);
}

/**
 * @fn     long usys_telldir(DIR*)
 * @brief  return current location in directory stream
 *
 * @param  dirp
 * @return On success, returns the current location of directory stream
 *         On error, -1 is returned, and errno is set to indicate the error.
 */
static inline long int usys_telldir(DIR *dirp){
	return telldir(dirp);
}

/**
 * @fn     int usys_chdir(const char*)
 * @brief  changes the current working directory of the calling
 *         process to the directory specified in path.
 *
 * @param  path
 * @return On success, 0.
 *         On error, -1
 */
static inline int usys_chdir(const char *path) {
	return chdir(path);
}

/**
 * @fn     int usy_chroot(const char*)
 * @brief  changes the root directory of the calling process to
 *         that specified in path.
 *
 * @param  path
 * @return On success, 0.
 *         On error, -1, errno is set to indicate the error.
 */
static inline  int usy_chroot(const char *path) {
	return chroot(path);
}

#ifdef __cplusplus
}
#endif

#endif /* SYS_INC_USYS_DIR_H_ */
