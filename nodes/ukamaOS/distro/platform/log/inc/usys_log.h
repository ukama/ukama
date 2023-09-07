/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef USYS_SYS_LOG_H
#define USYS_SYS_LOG_H

#ifdef __cplusplus
extern "C" {
#endif

#include "log.h"

/* Log levels */
#define USYS_LOG_TRACE 		LOG_TRACE
#define USYS_LOG_DEBUG		LOG_DEBUG
#define USYS_LOG_INFO		LOG_INFO
#define USYS_LOG_WARN		LOG_WARN
#define USYS_LOG_ERROR		LOG_ERROR
#define USYS_LOG_FATAL		LOG_FATAL

/* Logging */
#define usys_log_trace(...) log_log(USYS_LOG_TRACE, __FILE__, __LINE__, __VA_ARGS__)
#define usys_log_debug(...) log_log(USYS_LOG_DEBUG, __FILE__, __LINE__, __VA_ARGS__)
#define usys_log_info(...) log_log(USYS_LOG_INFO, __FILE__, __LINE__, __VA_ARGS__)
#define usys_log_warn(...) log_log(USYS_LOG_WARN, __FILE__, __LINE__, __VA_ARGS__)
#define usys_log_error(...) log_log(USYS_LOG_ERROR, __FILE__, __LINE__, __VA_ARGS__)
#define usys_log_fatal(...) log_log(USYS_LOG_FATAL, __FILE__, __LINE__, __VA_ARGS__)

/**
 * @brief Returns the name of the given log level as a string.
 *
 * @param level
 * @return const char*
 */
static inline const char *usys_log_level_string(int level) {
    return log_level_string(level);
}

/**
 * @brief If the log will be written to from multiple threads a lock function can be set.
 *        The function is passed the boolean true if the lock should be acquired or false
 *        if the lock should be released and the given udata value.
 *
 * @param fn
 * @param udata
 */
void static inline usys_log_set_lock(log_LockFn fn, void *udata) {
    log_set_lock(fn, udata);
}

/**
 * @brief The current logging level can be set by using the log_set_level() function.
 *        All logs below the given level will not be written to stderr.
 *        By default the level is LOG_TRACE, such that nothing is ignored.
 *
 * @param level
 */
static inline void usys_log_set_level(int level) {
    log_set_level(level);
}

/**
 * @brief Quiet-mode can be enabled by passing true to the log_set_quiet() function.
 *        While this mode is enabled the library will not output anything to stderr,
 *        but will continue to write to files and callbacks if any are set.
 *
 * @param enable
 */
static inline void usys_log_set_quiet(bool enable) {
    log_set_quiet(enable);
}

/**
 * @brief One or more callback functions which are called with the log data can be
 *        provided to the library by using the log_add_callback() function.
 *        A callback function is passed a log_Event structure containing the line number,
 *       filename, fmt string, va printf va_list, level and the given udata.
 *
 * @param fn
 * @param udata
 * @param level
 * @return int
 */
static inline int usys_log_add_callback(log_LogFn fn, void *udata, int level) {
    return log_add_callback(fn, udata, level);
}

/**
 * @brief One or more file pointers where the log will be written can be provided to
 *        the library by using the log_add_fp() function.
 *
 * @param fp
 * @param level
 * @return int
 */
static inline int usys_log_add_fp(FILE *fp, int level) {
    return log_add_fp(fp, level);
}

#ifdef __cplusplus
extern "C" {
#endif

#endif /* USYS_SYS_LOG_H */
