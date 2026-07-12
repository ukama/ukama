/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "log_router.h"

#include <errno.h>
#include <fcntl.h>
#include <jansson.h>
#include <limits.h>
#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/epoll.h>
#include <sys/eventfd.h>
#include <time.h>
#include <unistd.h>

#include "rlog_client.h"

#define ROUTER_MAX_EVENTS 32
#define ROUTER_BUFFER_MAX (32U * 1024U)
#define BOOT_ID_PATH      "/proc/sys/kernel/random/boot_id"

#define STREAM_STDOUT "stdout"
#define STREAM_STDERR "stderr"

typedef struct LogStream {
    int fd;
    char *space;
    char *app;
    pid_t pid;
    uint32_t generation;
    char stream[8];
    char *buffer;
    size_t used;
    size_t capacity;
    struct LogStream *next;
} LogStream;

typedef struct {
    int epollFd;
    int wakeFd;
    int maxRecordBytes;
    bool running;
    pthread_t thread;
    pthread_mutex_t mutex;
    LogStream *streams;
    RlogClient *client;
    char producerBootId[128];
    uint64_t captureSeq;
} LogRouter;

static LogRouter *gRouter = NULL;

static void trim_newline(char *value) {
    size_t len;

    if (!value) return;

    len = strlen(value);
    while (len > 0 &&
           (value[len - 1] == '\n' || value[len - 1] == '\r')) {
        value[--len] = '\0';
    }
}

static void read_boot_id(char *buffer, size_t size) {
    FILE *fp;

    if (!buffer || size == 0) return;

    snprintf(buffer, size, "starterd-%ld", (long)getpid());
    fp = fopen(BOOT_ID_PATH, "r");
    if (!fp) return;

    if (fgets(buffer, (int)size, fp) != NULL) {
        char bootId[128];

        trim_newline(buffer);
        snprintf(bootId, sizeof(bootId), "%s", buffer);
        snprintf(buffer, size, "%s:%ld", bootId, (long)getpid());
    }
    fclose(fp);
}

static void write_fallback(json_t *record) {
    char *data;
    size_t len;
    const char newline = '\n';

    if (!record) return;

    data = json_dumps(record, JSON_COMPACT | JSON_ENSURE_ASCII);
    if (!data) return;

    len = strlen(data);
    (void)write(STDERR_FILENO, data, len);
    (void)write(STDERR_FILENO, &newline, 1);
    free(data);
}

static const char *fallback_level(const char *stream) {
    if (stream && strcmp(stream, STREAM_STDERR) == 0) {
        return "warn";
    }
    return "info";
}

static json_t *wrap_unstructured(const LogStream *stream,
                                 const char *line) {
    json_t *record;

    record = json_object();
    if (!record) return NULL;

    json_object_set_new(record, "schema", json_string("ukama.log.v1"));
    json_object_set_new(record, "level",
                        json_string(fallback_level(stream->stream)));
    json_object_set_new(record, "app", json_string(stream->app));
    json_object_set_new(record, "component", json_string("process"));
    json_object_set_new(record, "event", json_string("process_output"));
    json_object_set_new(record, "msg", json_string(line ? line : ""));
    json_object_set_new(record, "structured", json_false());

    return record;
}

static json_t *parse_record(const LogStream *stream,
                            const char *line,
                            size_t len) {
    json_error_t error;
    json_t *record;

    record = json_loadb(line, len, 0, &error);
    if (!record || !json_is_object(record)) {
        if (record) json_decref(record);
        return wrap_unstructured(stream, line);
    }

    return record;
}

static bool is_rlog_app(const char *name) {
    if (!name) return false;
    return strcmp(name, "rlog") == 0 ||
           strcmp(name, "rlog.d") == 0 ||
           strcmp(name, "rlogd") == 0;
}

static void route_line(LogRouter *router, LogStream *stream,
                       const char *line, size_t len) {
    json_t *record;
    char captureId[256];
    bool sent;

    if (!router || !stream || !line) return;

    record = parse_record(stream, line, len);
    if (!record) return;

    router->captureSeq++;
    snprintf(captureId, sizeof(captureId), "%s:%llu",
             router->producerBootId,
             (unsigned long long)router->captureSeq);

    json_object_set_new(record, "space",
                        json_string(stream->space ? stream->space : ""));
    json_object_set_new(record, "app", json_string(stream->app));
    json_object_set_new(record, "pid", json_integer(stream->pid));
    json_object_set_new(record, "app_generation",
                        json_integer(stream->generation));
    json_object_set_new(record, "stream", json_string(stream->stream));
    json_object_set_new(record, "producer_boot_id",
                        json_string(router->producerBootId));
    json_object_set_new(record, "capture_seq",
                        json_integer((json_int_t)router->captureSeq));
    json_object_set_new(record, "capture_id", json_string(captureId));

    sent = false;
    if (!is_rlog_app(stream->app)) {
        sent = rlog_client_send(router->client,
                                record,
                                stream->space,
                                stream->app,
                                stream->pid,
                                stream->generation,
                                stream->stream,
                                router->captureSeq,
                                captureId);
    }

    if (!sent) {
        write_fallback(record);
    }

    json_decref(record);
}

static void consume_lines(LogRouter *router, LogStream *stream) {
    size_t start;
    size_t idx;

    if (!router || !stream) return;

    start = 0;
    for (idx = 0; idx < stream->used; idx++) {
        if (stream->buffer[idx] != '\n') continue;

        stream->buffer[idx] = '\0';
        if (idx > start && stream->buffer[idx - 1] == '\r') {
            stream->buffer[idx - 1] = '\0';
        }
        route_line(router, stream,
                   stream->buffer + start,
                   strlen(stream->buffer + start));
        start = idx + 1;
    }

    if (start > 0) {
        memmove(stream->buffer, stream->buffer + start,
                stream->used - start);
        stream->used -= start;
    }
}

static void stream_remove(LogRouter *router, LogStream *stream) {
    LogStream **cursor;

    if (!router || !stream) return;

    epoll_ctl(router->epollFd, EPOLL_CTL_DEL, stream->fd, NULL);
    close(stream->fd);

    pthread_mutex_lock(&router->mutex);
    cursor = &router->streams;
    while (*cursor) {
        if (*cursor == stream) {
            *cursor = stream->next;
            break;
        }
        cursor = &(*cursor)->next;
    }
    pthread_mutex_unlock(&router->mutex);

    free(stream->space);
    free(stream->app);
    free(stream->buffer);
    free(stream);
}

static void stream_eof(LogRouter *router, LogStream *stream) {
    if (stream->used > 0) {
        stream->buffer[stream->used] = '\0';
        route_line(router, stream, stream->buffer, stream->used);
        stream->used = 0;
    }
    stream_remove(router, stream);
}

static void stream_read(LogRouter *router, LogStream *stream) {
    while (1) {
        ssize_t n;
        size_t room;

        if (stream->used >= stream->capacity - 1) {
            stream->buffer[stream->used] = '\0';
            route_line(router, stream, stream->buffer, stream->used);
            stream->used = 0;
        }

        room = stream->capacity - stream->used - 1;
        n = read(stream->fd, stream->buffer + stream->used, room);
        if (n > 0) {
            stream->used += (size_t)n;
            stream->buffer[stream->used] = '\0';
            consume_lines(router, stream);
            continue;
        }

        if (n == 0) {
            stream_eof(router, stream);
            return;
        }

        if (errno == EINTR) continue;
        if (errno == EAGAIN || errno == EWOULDBLOCK) return;

        stream_eof(router, stream);
        return;
    }
}

static void *router_thread(void *data) {
    LogRouter *router;
    struct epoll_event events[ROUTER_MAX_EVENTS];

    router = (LogRouter *)data;
    while (router->running) {
        int count;
        int idx;

        count = epoll_wait(router->epollFd, events,
                           ROUTER_MAX_EVENTS, -1);
        if (count < 0) {
            if (errno == EINTR) continue;
            break;
        }

        for (idx = 0; idx < count; idx++) {
            if (events[idx].data.ptr == NULL) {
                uint64_t value;
                (void)read(router->wakeFd, &value, sizeof(value));
                continue;
            }

            stream_read(router,
                        (LogStream *)events[idx].data.ptr);
        }
    }

    return NULL;
}

static bool add_stream(LogRouter *router,
                       const char *space,
                       const char *app,
                       pid_t pid,
                       uint32_t generation,
                       int fd,
                       const char *streamName) {
    LogStream *stream;
    struct epoll_event event;
    int flags;

    if (!router || fd < 0 || !app || !streamName) return false;

    stream = calloc(1, sizeof(*stream));
    if (!stream) return false;

    stream->space = strdup(space ? space : "");
    stream->app = strdup(app);
    stream->buffer = calloc(1, (size_t)router->maxRecordBytes + 1U);
    if (!stream->space || !stream->app || !stream->buffer) {
        free(stream->space);
        free(stream->app);
        free(stream->buffer);
        free(stream);
        return false;
    }

    stream->fd = fd;
    stream->pid = pid;
    stream->generation = generation;
    stream->capacity = (size_t)router->maxRecordBytes + 1U;
    snprintf(stream->stream, sizeof(stream->stream), "%s", streamName);

    flags = fcntl(fd, F_GETFL, 0);
    if (flags < 0 || fcntl(fd, F_SETFL, flags | O_NONBLOCK) != 0) {
        free(stream->space);
        free(stream->app);
        free(stream->buffer);
        free(stream);
        return false;
    }

    memset(&event, 0, sizeof(event));
    event.events = EPOLLIN | EPOLLRDHUP | EPOLLHUP;
    event.data.ptr = stream;

    pthread_mutex_lock(&router->mutex);
    stream->next = router->streams;
    router->streams = stream;
    pthread_mutex_unlock(&router->mutex);

    if (epoll_ctl(router->epollFd, EPOLL_CTL_ADD, fd, &event) != 0) {
        LogStream **cursor;

        pthread_mutex_lock(&router->mutex);
        cursor = &router->streams;
        while (*cursor) {
            if (*cursor == stream) {
                *cursor = stream->next;
                break;
            }
            cursor = &(*cursor)->next;
        }
        pthread_mutex_unlock(&router->mutex);

        free(stream->space);
        free(stream->app);
        free(stream->buffer);
        free(stream);
        return false;
    }

    return true;
}

bool log_router_start(const char *socketPath,
                      int maxRecordBytes,
                      int reconnectMs) {
    LogRouter *router;
    struct epoll_event event;

    if (gRouter || !socketPath || !*socketPath) return false;

    router = calloc(1, sizeof(*router));
    if (!router) return false;

    router->maxRecordBytes = maxRecordBytes > 0 ?
                             maxRecordBytes : (int)ROUTER_BUFFER_MAX;
    router->epollFd = epoll_create1(EPOLL_CLOEXEC);
    router->wakeFd = eventfd(0, EFD_NONBLOCK | EFD_CLOEXEC);
    if (router->epollFd < 0 || router->wakeFd < 0) {
        if (router->epollFd >= 0) close(router->epollFd);
        if (router->wakeFd >= 0) close(router->wakeFd);
        free(router);
        return false;
    }

    pthread_mutex_init(&router->mutex, NULL);
    read_boot_id(router->producerBootId,
                 sizeof(router->producerBootId));
    router->client = rlog_client_create(socketPath,
                                         router->producerBootId,
                                         reconnectMs);
    if (!router->client) {
        close(router->wakeFd);
        close(router->epollFd);
        pthread_mutex_destroy(&router->mutex);
        free(router);
        return false;
    }

    memset(&event, 0, sizeof(event));
    event.events = EPOLLIN;
    event.data.ptr = NULL;
    if (epoll_ctl(router->epollFd, EPOLL_CTL_ADD,
                  router->wakeFd, &event) != 0) {
        rlog_client_destroy(router->client);
        close(router->wakeFd);
        close(router->epollFd);
        pthread_mutex_destroy(&router->mutex);
        free(router);
        return false;
    }

    router->running = true;
    if (pthread_create(&router->thread, NULL,
                       router_thread, router) != 0) {
        rlog_client_destroy(router->client);
        close(router->wakeFd);
        close(router->epollFd);
        pthread_mutex_destroy(&router->mutex);
        free(router);
        return false;
    }

    gRouter = router;
    return true;
}

void log_router_stop(void) {
    LogRouter *router;
    LogStream *stream;
    uint64_t wake;

    router = gRouter;
    if (!router) return;

    router->running = false;
    wake = 1;
    (void)write(router->wakeFd, &wake, sizeof(wake));
    pthread_join(router->thread, NULL);

    pthread_mutex_lock(&router->mutex);
    stream = router->streams;
    router->streams = NULL;
    pthread_mutex_unlock(&router->mutex);

    while (stream) {
        LogStream *next = stream->next;
        close(stream->fd);
        free(stream->space);
        free(stream->app);
        free(stream->buffer);
        free(stream);
        stream = next;
    }

    rlog_client_destroy(router->client);
    close(router->wakeFd);
    close(router->epollFd);
    pthread_mutex_destroy(&router->mutex);
    free(router);
    gRouter = NULL;
}

bool log_router_register(const char *space,
                         const char *app,
                         pid_t pid,
                         uint32_t generation,
                         int stdoutFd,
                         int stderrFd) {
    bool stdoutOk;
    bool stderrOk;

    if (!gRouter || stdoutFd < 0 || stderrFd < 0) return false;

    stdoutOk = add_stream(gRouter, space, app, pid, generation,
                          stdoutFd, STREAM_STDOUT);
    stderrOk = add_stream(gRouter, space, app, pid, generation,
                          stderrFd, STREAM_STDERR);

    if (!stdoutOk || !stderrOk) {
        if (!stdoutOk) close(stdoutFd);
        if (!stderrOk) close(stderrFd);
        return false;
    }

    return true;
}
