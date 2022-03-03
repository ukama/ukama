/**
 * Copyright (c) 2020-present, Ukama.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_PROCESS_H_
#define USYS_PROCESS_H_

#ifdef __cplusplus
extern "C" {
#endif

#include "usys_types.h"

/**
 * @fn     USysPid usys_fork(void)
 * @brief  creating a new process
 *
 * @return Negative Value: creation of a child process was unsuccessful.
           Zero: Returned to the newly created child process.
           Positive value: Returned to parent or caller.
           The value contains process ID of newly created child process.
 */
static inline USysPid usys_fork(void) {
	return fork();
}

/**
 * @fn     int usys_execv(const char*, char* const[])
 * @brief  Replaces the current process image with a new process image.
 *         An array of pointers to null-terminated strings that represent
 *         the argument list available to the new program. The first argument,
 *         by convention, should point to the filename associated with the file
 *         being executed. The array of pointers must be terminated by a NULL
 *         pointer.
 *
 * @param  path
 * @param  argv
 * @return The return value is -1 on error and errno is set to indicate the error
 *         otherwise doesn't return.
 */
static inline int usys_execv(const char *path, char *const argv[]){
 return execv(path, argv);
}

/**
 * @fn     USysPid usys_wait(int*)
 * @brief  used to wait for state changes in a
 *         child of the calling process, and obtain information about the
 *         child whose state has changed
 *
 * @param  wstatus
 * @return On success, returns the process ID of the terminated child
 *         On failure, -1 is returned.
 */
static inline USysPid usys_wait(int *wstatus) {
    return wait(wstatus);
}

/**
 * @fn     USysPid usys_waitpid(USysPid, int*, int)
 * @brief  function allows the calling thread to obtain status information
 *         for specified child processes.
 *
 * @param  pid
 * @param  status
 * @param  options
 * @return On Success the value returned indicates the process ID of the child process
 *         whose status information was recorded in the storage pointed
 *         to by stat_loc.
 *         0 if no child process was immediately available.
 *         On error -1 The errno value is set to indicate the error.
 */
static inline USysPid usys_waitpid(USysPid pid, int *status, int options){
	return waitpid(pid, status, options);
}

/**
 * @fn     USysPid usys_getpid()
 * @brief  returns the process ID of the calling process.
 *
 * @return USysPid
 */
static inline USysPid usys_getpid(){
	return getpid();
}

/**
 * @fn     USysPid usys_getppid()
 * @brief  returns the parent process ID of the calling process.
 *
 * @return USysPid
 */
static inline USysPid usys_getppid(){
	return getppid();
}

/**
 * @fn     USysPid usys_getpgrp()
 * @brief  returns the process group ID of the calling process.
 *
 * @return USysPid
 */
static inline USysPid usys_getpgrp(){
	return getpgrp();
}

/**
 * @fn     int usys_setpgid(USysPid, USysPid)
 * @brief  join an existing process group or create a new process group
 *         within the session of the calling process
 *
 * @param  pid
 * @param  pgid
 * @return On Success 0
 * 		   On error -1. The errno variable is set to indicate the error.
 */
static inline int usys_setpgid(USysPid pid, USysPid pgid){
	return setpgid(pid, pgid);
}

/**
 * @fn     long int usys_ulimit(int, ...)
 * @brief  function provides a way to get and set process resource limits.
 *
 * @param  cmd
 * @return On success 0
 * 		   On error -1. The errno variable is set to indicate the error.
 */
static inline long int usys_ulimit(int cmd, ...){
	va_list args;
	return ulimit(cmd, args);
}

/**
 * @fn     int usys_getrlimit(int, struct rlimit*)
 * @brief  function returns the resource limit for the specified resource.
 *
 * @param  resource
 * @param  rlp
 * @return On success 0
 * 		   On error -1. The errno variable is set to indicate the error.
 */
static inline int usys_getrlimit(int resource, struct rlimit *rlp){
	return getrlimit(resource, rlp);
}

/**
 * @fn     int usys_setrlimit(int, struct rlimit*)
 * @brief  function set the resource limit for the specified resource.
 *
 * @param  resource
 * @param  rlp
 * @return On success 0
 * 		   On error -1. The errno variable is set to indicate the error.
 */
static inline int usys_setrlimit(int resource, struct rlimit *rlp){
	return setrlimit(resource, rlp);
}

/**
 * @fn     int usys_getopt(int, char* const[], const char*)
 * @brief  function returns the next flag letter in the argv
 *         list that matches a letter in optionstring.
 *
 * @param  argc
 * @param  argv
 * @param  optionstring
 * @return EOF on successful processing of all flags.
 * 		   ? on unkown flag.
 */
static inline int usys_getopt(int argc, char * const argv[],
     const char *optionstring) {
	return getopt(argc, argv, optionstring);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_PROCESS_H_ */
