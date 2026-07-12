/**
 * Copyright (c) 2020 rxi
 *
 * This library is free software; you can redistribute it and/or modify it
 * under the terms of the MIT license. See `log.c` for details.
 */

#ifndef LOG_H
#define LOG_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdarg.h>
#include <stdio.h>
#include <time.h>

#ifdef __cplusplus
extern "C" {
#endif

#define LOG_VERSION "0.2.0"
#define LOG_SCHEMA_VERSION "ukama.log.v1"
#define LOG_MAX_FIELDS 32
#define LOG_MAX_MESSAGE 2048
#define LOG_MAX_RECORD 8192

typedef struct {
    va_list ap;
    const char *fmt;
    const char *file;
    struct tm *time;
    void *udata;
    int line;
    int level;
    const char *function;
} log_Event;

typedef void (*log_LogFn)(log_Event *ev);
typedef void (*log_LockFn)(bool lock, void *udata);

typedef enum {
    LOG_FIELD_STRING = 0,
    LOG_FIELD_INTEGER,
    LOG_FIELD_UNSIGNED,
    LOG_FIELD_BOOLEAN,
    LOG_FIELD_DOUBLE
} log_FieldType;

typedef struct {
    const char *key;
    log_FieldType type;
    union {
        const char *stringValue;
        int64_t integerValue;
        uint64_t unsignedValue;
        bool booleanValue;
        double doubleValue;
    } value;
} log_Field;

enum {
    LOG_TRACE,
    LOG_DEBUG,
    LOG_INFO,
    LOG_WARN,
    LOG_ERROR,
    LOG_FATAL
};

#define log_trace(...)                                                     \
    log_log_ex(LOG_TRACE, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_debug(...)                                                     \
    log_log_ex(LOG_DEBUG, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_info(...)                                                      \
    log_log_ex(LOG_INFO, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_warn(...)                                                      \
    log_log_ex(LOG_WARN, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_error(...)                                                     \
    log_log_ex(LOG_ERROR, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_fatal(...)                                                     \
    log_log_ex(LOG_FATAL, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define log_critical(...)                                                  \
    log_log_ex(LOG_FATAL, __FILE__, __LINE__, __func__, __VA_ARGS__)

const char *log_level_string(int level);
void log_set_lock(log_LockFn fn, void *udata);
void log_set_level(int level);
void log_set_quiet(bool enable);
void log_set_service(char *service);
int log_add_callback(log_LogFn fn, void *udata, int level);
int log_add_fp(FILE *fp, int level);

void log_log(int level, const char *file, int line, const char *fmt, ...);
void log_log_ex(int level, const char *file, int line,
                const char *function, const char *fmt, ...);
void log_event(int level, const char *file, int line,
               const char *function, const char *component,
               const char *event, const char *message,
               const log_Field *fields, size_t fieldCount);

/*
 * Deprecated compatibility routines. Structured logs are written to stderr
 * and captured by starterd. These routines intentionally perform no remote
 * logging.
 */
void log_remote_init(char *serviceName);
int log_rlogd(char *message);
int is_connect_with_rlogd(void);
void log_enable_rlogd(int flag);

#ifdef __cplusplus
}
#endif

#endif
