/*
 * Copyright (c) 2020 rxi
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */

#include "log.h"

#include <errno.h>
#include <math.h>
#include <pthread.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/syscall.h>
#include <sys/types.h>
#include <unistd.h>

#define MAX_CALLBACKS 32
#define MAX_SERVICE_NAME 64
#define FALLBACK_MESSAGE "log record exceeded size limit"

typedef struct {
    log_LogFn fn;
    void *udata;
    int level;
} Callback;

typedef struct {
    char *data;
    size_t size;
    size_t used;
    bool failed;
} JsonBuffer;

static struct {
    void *udata;
    log_LockFn lock;
    int level;
    bool quiet;
    char service[MAX_SERVICE_NAME];
    Callback callbacks[MAX_CALLBACKS];
} l = {
    .level = LOG_TRACE,
    .service = "unknown"
};

static pthread_mutex_t logMutex = PTHREAD_MUTEX_INITIALIZER;

static const char *levelStrings[] = {
    "TRACE",
    "DEBUG",
    "INFO",
    "WARN",
    "ERROR",
    "FATAL"
};

static const char *jsonLevelStrings[] = {
    "trace",
    "debug",
    "info",
    "warn",
    "error",
    "critical"
};

static bool level_is_valid(int level) {
    return level >= LOG_TRACE && level <= LOG_FATAL;
}

static void buffer_append(JsonBuffer *buffer, const char *value,
                          size_t length) {
    if (buffer->failed) {
        return;
    }

    if (length >= buffer->size - buffer->used) {
        buffer->failed = true;
        return;
    }

    memcpy(buffer->data + buffer->used, value, length);
    buffer->used += length;
    buffer->data[buffer->used] = '\0';
}

static void buffer_append_char(JsonBuffer *buffer, char value) {
    buffer_append(buffer, &value, 1);
}

static void buffer_append_format(JsonBuffer *buffer, const char *format, ...) {
    int written;
    va_list args;

    if (buffer->failed) {
        return;
    }

    va_start(args, format);
    written = vsnprintf(buffer->data + buffer->used,
                        buffer->size - buffer->used,
                        format,
                        args);
    va_end(args);

    if (written < 0 || (size_t)written >= buffer->size - buffer->used) {
        buffer->failed = true;
        return;
    }

    buffer->used += (size_t)written;
}

static void buffer_append_json_string(JsonBuffer *buffer,
                                      const char *value) {
    const unsigned char *cursor;

    if (value == NULL) {
        buffer_append(buffer, "null", 4);
        return;
    }

    buffer_append_char(buffer, '"');
    for (cursor = (const unsigned char *)value; *cursor != '\0'; cursor++) {
        switch (*cursor) {
        case '"':
            buffer_append(buffer, "\\\"", 2);
            break;
        case '\\':
            buffer_append(buffer, "\\\\", 2);
            break;
        case '\b':
            buffer_append(buffer, "\\b", 2);
            break;
        case '\f':
            buffer_append(buffer, "\\f", 2);
            break;
        case '\n':
            buffer_append(buffer, "\\n", 2);
            break;
        case '\r':
            buffer_append(buffer, "\\r", 2);
            break;
        case '\t':
            buffer_append(buffer, "\\t", 2);
            break;
        default:
            if (*cursor < 0x20) {
                buffer_append_format(buffer, "\\u%04x", *cursor);
            } else {
                buffer_append_char(buffer, (char)*cursor);
            }
            break;
        }

        if (buffer->failed) {
            return;
        }
    }
    buffer_append_char(buffer, '"');
}

static void buffer_append_key(JsonBuffer *buffer, bool *first,
                              const char *key) {
    if (!*first) {
        buffer_append_char(buffer, ',');
    }
    *first = false;
    buffer_append_json_string(buffer, key);
    buffer_append_char(buffer, ':');
}

static void buffer_append_string_field(JsonBuffer *buffer, bool *first,
                                       const char *key,
                                       const char *value) {
    buffer_append_key(buffer, first, key);
    buffer_append_json_string(buffer, value);
}

static void buffer_append_integer_field(JsonBuffer *buffer, bool *first,
                                        const char *key, int64_t value) {
    buffer_append_key(buffer, first, key);
    buffer_append_format(buffer, "%lld", (long long)value);
}

static void buffer_append_unsigned_field(JsonBuffer *buffer, bool *first,
                                         const char *key,
                                         uint64_t value) {
    buffer_append_key(buffer, first, key);
    buffer_append_format(buffer, "%llu", (unsigned long long)value);
}

static void buffer_append_boolean_field(JsonBuffer *buffer, bool *first,
                                        const char *key, bool value) {
    buffer_append_key(buffer, first, key);
    if (value) {
        buffer_append(buffer, "true", 4);
    } else {
        buffer_append(buffer, "false", 5);
    }
}

static void buffer_append_double_field(JsonBuffer *buffer, bool *first,
                                       const char *key, double value) {
    buffer_append_key(buffer, first, key);
    if (isfinite(value)) {
        buffer_append_format(buffer, "%.17g", value);
    } else {
        buffer_append(buffer, "null", 4);
    }
}

static bool key_is_reserved(const char *key) {
    static const char *reserved[] = {
        "schema",
        "ts",
        "mono_ms",
        "level",
        "app",
        "component",
        "event",
        "msg",
        "pid",
        "tid",
        "source_file",
        "source_line",
        "source_function",
        "truncated"
    };
    size_t idx;

    if (key == NULL || key[0] == '\0') {
        return true;
    }

    for (idx = 0; idx < sizeof(reserved) / sizeof(reserved[0]); idx++) {
        if (strcmp(key, reserved[idx]) == 0) {
            return true;
        }
    }

    return false;
}

static pid_t current_thread_id(void) {
#ifdef SYS_gettid
    return (pid_t)syscall(SYS_gettid);
#else
    return getpid();
#endif
}

static void make_timestamps(char *timestamp, size_t timestampSize,
                            uint64_t *monoMs) {
    struct timespec realtime;
    struct timespec monotonic;
    struct tm utc;
    size_t length;

    clock_gettime(CLOCK_REALTIME, &realtime);
    clock_gettime(CLOCK_MONOTONIC, &monotonic);
    gmtime_r(&realtime.tv_sec, &utc);

    length = strftime(timestamp, timestampSize,
                      "%Y-%m-%dT%H:%M:%S", &utc);
    if (length > 0 && length < timestampSize) {
        snprintf(timestamp + length, timestampSize - length,
                 ".%03ldZ", realtime.tv_nsec / 1000000L);
    }

    *monoMs = ((uint64_t)monotonic.tv_sec * 1000ULL) +
              ((uint64_t)monotonic.tv_nsec / 1000000ULL);
}

static bool append_custom_fields(JsonBuffer *buffer, bool *first,
                                 const log_Field *fields,
                                 size_t fieldCount) {
    size_t idx;

    if (fields == NULL) {
        return true;
    }

    if (fieldCount > LOG_MAX_FIELDS) {
        fieldCount = LOG_MAX_FIELDS;
    }

    for (idx = 0; idx < fieldCount; idx++) {
        const log_Field *field = &fields[idx];

        if (key_is_reserved(field->key)) {
            continue;
        }

        switch (field->type) {
        case LOG_FIELD_STRING:
            buffer_append_string_field(buffer, first, field->key,
                                       field->value.stringValue);
            break;
        case LOG_FIELD_INTEGER:
            buffer_append_integer_field(buffer, first, field->key,
                                        field->value.integerValue);
            break;
        case LOG_FIELD_UNSIGNED:
            buffer_append_unsigned_field(buffer, first, field->key,
                                         field->value.unsignedValue);
            break;
        case LOG_FIELD_BOOLEAN:
            buffer_append_boolean_field(buffer, first, field->key,
                                        field->value.booleanValue);
            break;
        case LOG_FIELD_DOUBLE:
            buffer_append_double_field(buffer, first, field->key,
                                       field->value.doubleValue);
            break;
        default:
            continue;
        }

        if (buffer->failed) {
            return false;
        }
    }

    return true;
}

static bool build_record(char *record, size_t recordSize, int level,
                         const char *file, int line,
                         const char *function,
                         const char *component,
                         const char *event,
                         const char *message,
                         const log_Field *fields,
                         size_t fieldCount,
                         bool truncated) {
    JsonBuffer buffer = {
        .data = record,
        .size = recordSize,
        .used = 0,
        .failed = false
    };
    char timestamp[40] = {0};
    uint64_t monoMs = 0;
    bool first = true;

    make_timestamps(timestamp, sizeof(timestamp), &monoMs);

    buffer_append_char(&buffer, '{');
    buffer_append_string_field(&buffer, &first, "schema",
                               LOG_SCHEMA_VERSION);
    buffer_append_string_field(&buffer, &first, "ts", timestamp);
    buffer_append_unsigned_field(&buffer, &first, "mono_ms", monoMs);
    buffer_append_string_field(&buffer, &first, "level",
                               jsonLevelStrings[level]);
    buffer_append_string_field(&buffer, &first, "app", l.service);
    buffer_append_string_field(&buffer, &first, "component", component);
    buffer_append_string_field(&buffer, &first, "event", event);
    buffer_append_string_field(&buffer, &first, "msg", message);
    buffer_append_integer_field(&buffer, &first, "pid", getpid());
    buffer_append_integer_field(&buffer, &first, "tid",
                                current_thread_id());
    buffer_append_string_field(&buffer, &first, "source_file", file);
    buffer_append_integer_field(&buffer, &first, "source_line", line);

    if (function != NULL && function[0] != '\0') {
        buffer_append_string_field(&buffer, &first, "source_function",
                                   function);
    }

    if (!append_custom_fields(&buffer, &first, fields, fieldCount)) {
        return false;
    }

    if (truncated || fieldCount > LOG_MAX_FIELDS) {
        buffer_append_boolean_field(&buffer, &first, "truncated", true);
    }

    buffer_append_char(&buffer, '}');
    buffer_append_char(&buffer, '\n');

    return !buffer.failed;
}

static void write_record(const char *record) {
    size_t remaining = strlen(record);
    const char *cursor = record;

    while (remaining > 0) {
        ssize_t written = write(STDERR_FILENO, cursor, remaining);

        if (written < 0) {
            if (errno == EINTR) {
                continue;
            }
            return;
        }

        cursor += written;
        remaining -= (size_t)written;
    }
}

static void file_callback(log_Event *event) {
    char timestamp[64];

    timestamp[strftime(timestamp, sizeof(timestamp),
                       "%Y-%m-%d %H:%M:%S", event->time)] = '\0';
    fprintf(event->udata, "%s %s %-5s %s:%d: ", l.service, timestamp,
            levelStrings[event->level], event->file, event->line);
    vfprintf(event->udata, event->fmt, event->ap);
    fprintf(event->udata, "\n");
    fflush(event->udata);
}

static void dispatch_callback(Callback *callback, int level,
                              const char *file, int line,
                              const char *function,
                              const char *format, ...) {
    log_Event event = {
        .fmt = format,
        .file = file,
        .function = function,
        .line = line,
        .level = level,
        .udata = callback->udata
    };
    time_t now = time(NULL);
    struct tm localTime;

    localtime_r(&now, &localTime);
    event.time = &localTime;

    va_start(event.ap, format);
    callback->fn(&event);
    va_end(event.ap);
}

static void dispatch_callbacks(int level, const char *file, int line,
                               const char *function,
                               const char *message) {
    int idx;

    for (idx = 0; idx < MAX_CALLBACKS && l.callbacks[idx].fn; idx++) {
        Callback *callback = &l.callbacks[idx];

        if (level >= callback->level) {
            dispatch_callback(callback, level, file, line, function,
                              "%s", message);
        }
    }
}

static void output_event(int level, const char *file, int line,
                         const char *function,
                         const char *component,
                         const char *event,
                         const char *message,
                         const log_Field *fields,
                         size_t fieldCount,
                         bool messageTruncated) {
    char record[LOG_MAX_RECORD] = {0};
    bool built;

    pthread_mutex_lock(&logMutex);
    if (l.lock != NULL) {
        l.lock(true, l.udata);
    }

    if (!l.quiet && level >= l.level) {
        built = build_record(record, sizeof(record), level, file, line,
                             function, component, event, message, fields,
                             fieldCount, messageTruncated);
        if (!built) {
            memset(record, 0, sizeof(record));
            built = build_record(record, sizeof(record), level, file, line,
                                 function, component, event,
                                 FALLBACK_MESSAGE, NULL, 0, true);
        }

        if (built) {
            write_record(record);
        }
    }

    dispatch_callbacks(level, file, line, function, message);

    if (l.lock != NULL) {
        l.lock(false, l.udata);
    }
    pthread_mutex_unlock(&logMutex);
}

const char *log_level_string(int level) {
    if (!level_is_valid(level)) {
        return "UNKNOWN";
    }

    return levelStrings[level];
}

void log_set_lock(log_LockFn fn, void *udata) {
    pthread_mutex_lock(&logMutex);
    l.lock = fn;
    l.udata = udata;
    pthread_mutex_unlock(&logMutex);
}

void log_set_level(int level) {
    if (!level_is_valid(level)) {
        return;
    }

    pthread_mutex_lock(&logMutex);
    l.level = level;
    pthread_mutex_unlock(&logMutex);
}

void log_set_quiet(bool enable) {
    pthread_mutex_lock(&logMutex);
    l.quiet = enable;
    pthread_mutex_unlock(&logMutex);
}

void log_set_service(char *service) {
    pthread_mutex_lock(&logMutex);
    if (service == NULL || service[0] == '\0') {
        snprintf(l.service, sizeof(l.service), "%s", "unknown");
    } else {
        snprintf(l.service, sizeof(l.service), "%s", service);
    }
    pthread_mutex_unlock(&logMutex);
}

int log_add_callback(log_LogFn fn, void *udata, int level) {
    int idx;

    if (fn == NULL || !level_is_valid(level)) {
        return -1;
    }

    pthread_mutex_lock(&logMutex);
    for (idx = 0; idx < MAX_CALLBACKS; idx++) {
        if (l.callbacks[idx].fn == NULL) {
            l.callbacks[idx] = (Callback){fn, udata, level};
            pthread_mutex_unlock(&logMutex);
            return 0;
        }
    }
    pthread_mutex_unlock(&logMutex);

    return -1;
}

int log_add_fp(FILE *fp, int level) {
    if (fp == NULL) {
        return -1;
    }

    return log_add_callback(file_callback, fp, level);
}

static void log_vmessage(int level, const char *file, int line,
                         const char *function, const char *format,
                         va_list args) {
    char message[LOG_MAX_MESSAGE + 1] = {0};
    int written;
    bool truncated = false;

    if (!level_is_valid(level) || format == NULL) {
        return;
    }

    written = vsnprintf(message, sizeof(message), format, args);
    if (written < 0) {
        snprintf(message, sizeof(message), "%s", "unable to format log");
        truncated = true;
    } else if ((size_t)written >= sizeof(message)) {
        truncated = true;
    }

    output_event(level, file, line, function, "general", "log_message",
                 message, NULL, 0, truncated);
}

void log_log(int level, const char *file, int line, const char *format, ...) {
    va_list args;

    va_start(args, format);
    log_vmessage(level, file, line, NULL, format, args);
    va_end(args);
}

void log_log_ex(int level, const char *file, int line,
                const char *function, const char *format, ...) {
    va_list args;

    va_start(args, format);
    log_vmessage(level, file, line, function, format, args);
    va_end(args);
}

void log_event(int level, const char *file, int line,
               const char *function, const char *component,
               const char *event, const char *message,
               const log_Field *fields, size_t fieldCount) {
    char boundedMessage[LOG_MAX_MESSAGE + 1] = {0};
    bool truncated = false;
    size_t messageLength;

    if (!level_is_valid(level) || event == NULL || event[0] == '\0') {
        return;
    }

    if (component == NULL || component[0] == '\0') {
        component = "general";
    }
    if (message == NULL) {
        message = "";
    }

    messageLength = strlen(message);
    if (messageLength > LOG_MAX_MESSAGE) {
        messageLength = LOG_MAX_MESSAGE;
        truncated = true;
    }
    memcpy(boundedMessage, message, messageLength);
    boundedMessage[messageLength] = '\0';

    output_event(level, file, line, function, component, event,
                 boundedMessage, fields, fieldCount, truncated);
}
