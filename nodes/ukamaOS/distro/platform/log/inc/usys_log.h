/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef USYS_SYS_LOG_H
#define USYS_SYS_LOG_H

#include "log.h"

#ifdef __cplusplus
extern "C" {
#endif

/* Log levels */
#define USYS_LOG_TRACE LOG_TRACE
#define USYS_LOG_DEBUG LOG_DEBUG
#define USYS_LOG_INFO LOG_INFO
#define USYS_LOG_WARN LOG_WARN
#define USYS_LOG_ERROR LOG_ERROR
#define USYS_LOG_FATAL LOG_FATAL
#define USYS_LOG_CRITICAL LOG_FATAL

/* Existing formatted logging API. */
#define usys_log_trace(...)                                                \
    log_log_ex(USYS_LOG_TRACE, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_debug(...)                                                \
    log_log_ex(USYS_LOG_DEBUG, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_info(...)                                                 \
    log_log_ex(USYS_LOG_INFO, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_warn(...)                                                 \
    log_log_ex(USYS_LOG_WARN, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_error(...)                                                \
    log_log_ex(USYS_LOG_ERROR, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_fatal(...)                                                \
    log_log_ex(USYS_LOG_FATAL, __FILE__, __LINE__, __func__, __VA_ARGS__)
#define usys_log_critical(...)                                             \
    log_log_ex(USYS_LOG_CRITICAL, __FILE__, __LINE__, __func__,           \
               __VA_ARGS__)

typedef log_Field USysLogField;
typedef log_FieldType USysLogFieldType;

#define USYS_LOG_FIELD_COUNT(fields)                                       \
    (sizeof(fields) / sizeof((fields)[0]))

static inline USysLogField usys_log_field_string(const char *key,
                                                  const char *value) {
    USysLogField field = {NULL, LOG_FIELD_STRING, {NULL}};

    field.key  = key;
    field.type = LOG_FIELD_STRING;
    field.value.stringValue = value;
    return field;
}

static inline USysLogField usys_log_field_integer(const char *key,
                                                   int64_t value) {
    USysLogField field = {NULL, LOG_FIELD_INTEGER, {NULL}};

    field.key  = key;
    field.type = LOG_FIELD_INTEGER;
    field.value.integerValue = value;
    return field;
}

static inline USysLogField usys_log_field_unsigned(const char *key,
                                                    uint64_t value) {
    USysLogField field = {NULL, LOG_FIELD_UNSIGNED, {NULL}};

    field.key  = key;
    field.type = LOG_FIELD_UNSIGNED;
    field.value.unsignedValue = value;
    return field;
}

static inline USysLogField usys_log_field_boolean(const char *key,
                                                   bool value) {
    USysLogField field = {NULL, LOG_FIELD_BOOLEAN, {NULL}};

    field.key  = key;
    field.type = LOG_FIELD_BOOLEAN;
    field.value.booleanValue = value;
    return field;
}

static inline USysLogField usys_log_field_double(const char *key,
                                                  double value) {
    USysLogField field = {NULL, LOG_FIELD_DOUBLE, {NULL}};

    field.key = key;
    field.type = LOG_FIELD_DOUBLE;
    field.value.doubleValue = value;
    return field;
}

#define USYS_LOG_STR(key, value) usys_log_field_string((key), (value))
#define USYS_LOG_INT(key, value) usys_log_field_integer((key), (value))
#define USYS_LOG_UINT(key, value) usys_log_field_unsigned((key), (value))
#define USYS_LOG_BOOL(key, value) usys_log_field_boolean((key), (value))
#define USYS_LOG_DOUBLE(key, value) usys_log_field_double((key), (value))

/* Structured event API. Message text is intentionally not printf-style. */
#define usys_log_event_trace(component, event, message, fields, count)     \
    log_event(USYS_LOG_TRACE, __FILE__, __LINE__, __func__, (component),  \
              (event), (message), (fields), (count))
#define usys_log_event_debug(component, event, message, fields, count)     \
    log_event(USYS_LOG_DEBUG, __FILE__, __LINE__, __func__, (component),  \
              (event), (message), (fields), (count))
#define usys_log_event_info(component, event, message, fields, count)      \
    log_event(USYS_LOG_INFO, __FILE__, __LINE__, __func__, (component),   \
              (event), (message), (fields), (count))
#define usys_log_event_warn(component, event, message, fields, count)      \
    log_event(USYS_LOG_WARN, __FILE__, __LINE__, __func__, (component),   \
              (event), (message), (fields), (count))
#define usys_log_event_error(component, event, message, fields, count)     \
    log_event(USYS_LOG_ERROR, __FILE__, __LINE__, __func__, (component),  \
              (event), (message), (fields), (count))
#define usys_log_event_fatal(component, event, message, fields, count)     \
    log_event(USYS_LOG_FATAL, __FILE__, __LINE__, __func__, (component),  \
              (event), (message), (fields), (count))
#define usys_log_event_critical(component, event, message, fields, count)  \
    log_event(USYS_LOG_CRITICAL, __FILE__, __LINE__, __func__,            \
              (component), (event), (message), (fields), (count))

static inline const char *usys_log_level_string(int level) {
    return log_level_string(level);
}

static inline void usys_log_set_lock(log_LockFn fn, void *udata) {
    log_set_lock(fn, udata);
}

static inline void usys_log_set_level(int level) {
    log_set_level(level);
}

static inline void usys_log_set_quiet(bool enable) {
    log_set_quiet(enable);
}

static inline int usys_log_add_callback(log_LogFn fn, void *udata,
                                        int level) {
    return log_add_callback(fn, udata, level);
}

static inline int usys_log_add_fp(FILE *fp, int level) {
    return log_add_fp(fp, level);
}

/* Deprecated compatibility function; intentionally a no-op. */
static inline void usys_log_remote_init(char *name) {
    log_remote_init(name);
}

static inline void usys_log_set_service(char *name) {
    log_set_service(name);
}

#ifdef __cplusplus
}
#endif

#endif /* USYS_SYS_LOG_H */
