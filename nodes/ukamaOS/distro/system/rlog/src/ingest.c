/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "ingest.h"

#include <errno.h>
#include <jansson.h>
#include <limits.h>
#include <pthread.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/un.h>
#include <unistd.h>

#define INGEST_MAX_FRAME  (64U * 1024U)
#define INGEST_BACKLOG    4
#define INGEST_STATE_FILE "ingest.json"
#define INGEST_PRODUCER_MAX 128
#define INGEST_SOCKET_PATH_MAX \
    sizeof(((struct sockaddr_un *)0)->sun_path)

struct IngestServer {
    char socketPath[INGEST_SOCKET_PATH_MAX];
    char statePath[PATH_MAX];
    int listenFd;
    int clientFd;
    pthread_t thread;
    bool running;
    LogStore *store;

    char producerBootId[INGEST_PRODUCER_MAX];
    uint64_t acceptedCaptureSeq;
};

static int mkdir_one(const char *path) {
    if (mkdir(path, 0755) == 0 || errno == EEXIST) {
        return 0;
    }
    return -1;
}

static int ensure_parent_dir(const char *path) {
    char parent[PATH_MAX];
    char *slash;
    char *p;

    if (!path || strlen(path) >= sizeof(parent)) {
        return -1;
    }

    strcpy(parent, path);
    slash = strrchr(parent, '/');
    if (!slash) {
        return 0;
    }
    *slash = '\0';

    for (p = parent + 1; *p; p++) {
        if (*p != '/') {
            continue;
        }
        *p = '\0';
        if (mkdir_one(parent) != 0) {
            return -1;
        }
        *p = '/';
    }

    return mkdir_one(parent);
}

static int send_json(int fd, json_t *json) {
    char *data;
    size_t len;
    ssize_t sent;

    data = json_dumps(json, JSON_COMPACT | JSON_ENSURE_ASCII);
    if (!data) {
        return -1;
    }

    len = strlen(data);
    sent = send(fd, data, len, MSG_NOSIGNAL);
    free(data);

    return sent == (ssize_t)len ? 0 : -1;
}

static int send_error(int fd, const char *message) {
    json_t *reply;
    int rc;

    reply = json_pack("{s:s,s:s}",
                      "op", "error",
                      "message", message ? message : "invalid request");
    if (!reply) {
        return -1;
    }

    rc = send_json(fd, reply);
    json_decref(reply);
    return rc;
}

static int send_ack(int fd, const char *producerBootId,
                    uint64_t acceptedThrough) {
    json_t *reply;
    int rc;

    reply = json_pack("{s:s,s:s,s:I}",
                      "op", "ack",
                      "producer_boot_id",
                      producerBootId ? producerBootId : "",
                      "accepted_through",
                      (json_int_t)acceptedThrough);
    if (!reply) {
        return -1;
    }

    rc = send_json(fd, reply);
    json_decref(reply);
    return rc;
}

static int write_state(IngestServer *server) {
    char temporary[PATH_MAX];
    json_t *state;
    char *data;
    FILE *file;
    int rc;

    if (!server) {
        return -1;
    }

    if (snprintf(temporary, sizeof(temporary), "%s.tmp",
                 server->statePath) >= (int)sizeof(temporary)) {
        return -1;
    }

    state = json_pack("{s:s,s:s,s:I}",
                      "schema", "ukama.rlog.ingest.v1",
                      "producer_boot_id", server->producerBootId,
                      "accepted_capture_seq",
                      (json_int_t)server->acceptedCaptureSeq);
    if (!state) {
        return -1;
    }

    data = json_dumps(state, JSON_COMPACT | JSON_SORT_KEYS);
    json_decref(state);
    if (!data) {
        return -1;
    }

    file = fopen(temporary, "w");
    if (!file) {
        free(data);
        return -1;
    }

    rc = fprintf(file, "%s\n", data) < 0 ? -1 : 0;
    if (rc == 0 && fflush(file) != 0) {
        rc = -1;
    }
    if (rc == 0 && fsync(fileno(file)) != 0) {
        rc = -1;
    }
    fclose(file);
    free(data);

    if (rc == 0 && rename(temporary, server->statePath) != 0) {
        rc = -1;
    }
    if (rc != 0) {
        unlink(temporary);
    }

    return rc;
}

static void load_state(IngestServer *server) {
    json_error_t error;
    json_t *state;
    json_t *accepted;
    const char *producerBootId;
    size_t len;

    if (!server) {
        return;
    }

    state = json_load_file(server->statePath, 0, &error);
    if (!state || !json_is_object(state)) {
        if (state) {
            json_decref(state);
        }
        return;
    }

    producerBootId = json_string_value(
        json_object_get(state, "producer_boot_id"));
    accepted = json_object_get(state, "accepted_capture_seq");
    if (producerBootId && json_is_integer(accepted)) {
        len = strlen(producerBootId);
        if (len < sizeof(server->producerBootId)) {
            memcpy(server->producerBootId, producerBootId, len + 1);
            server->acceptedCaptureSeq =
                (uint64_t)json_integer_value(accepted);
        }
    }

    json_decref(state);
}

static int select_producer(IngestServer *server,
                           const char *producerBootId) {
    size_t len;

    if (!server || !producerBootId || !*producerBootId) {
        return -1;
    }

    len = strlen(producerBootId);
    if (len >= sizeof(server->producerBootId)) {
        return -1;
    }

    if (strcmp(server->producerBootId, producerBootId) == 0) {
        return 0;
    }

    memcpy(server->producerBootId, producerBootId, len + 1);
    server->acceptedCaptureSeq = 0;
    return write_state(server);
}

static int handle_hello(IngestServer *server, int fd, json_t *request) {
    const char *producerBootId;
    json_t *reply;
    int rc;

    producerBootId = json_string_value(
        json_object_get(request, "producer_boot_id"));
    if (!producerBootId || !*producerBootId) {
        return send_error(fd, "producer_boot_id is required");
    }

    if (select_producer(server, producerBootId) != 0) {
        return send_error(fd, "unable to persist producer state");
    }

    {
        uint64_t recovered;

        recovered = log_store_last_capture_seq(server->store,
                                               producerBootId);
        if (recovered > server->acceptedCaptureSeq) {
            server->acceptedCaptureSeq = recovered;
            if (write_state(server) != 0) {
                return send_error(fd,
                                  "unable to recover ingest cursor");
            }
        }
    }

    reply = json_pack("{s:s,s:b,s:s,s:I,s:s}",
                      "op", "hello",
                      "ready", 1,
                      "schema", "ukama.log.v1",
                      "accepted_capture_seq",
                      (json_int_t)server->acceptedCaptureSeq,
                      "boot_id", log_store_boot_id(server->store));
    if (!reply) {
        return -1;
    }

    rc = send_json(fd, reply);
    json_decref(reply);
    return rc;
}

static int set_string_if_present(json_t *record, json_t *request,
                                 const char *key) {
    json_t *value;

    value = json_object_get(request, key);
    if (!json_is_string(value)) {
        return 0;
    }

    return json_object_set(record, key, value);
}

static int set_integer_if_present(json_t *record, json_t *request,
                                  const char *key) {
    json_t *value;

    value = json_object_get(request, key);
    if (!json_is_integer(value)) {
        return 0;
    }

    return json_object_set(record, key, value);
}

static int handle_event(IngestServer *server, int fd, json_t *request) {
    const char *producerBootId;
    json_t *captureSeqValue;
    json_t *recordValue;
    json_t *record;
    uint64_t captureSeq;
    uint64_t canonicalSeq;

    producerBootId = json_string_value(
        json_object_get(request, "producer_boot_id"));
    captureSeqValue = json_object_get(request, "capture_seq");
    recordValue = json_object_get(request, "record");

    if (!producerBootId || !*producerBootId ||
        !json_is_integer(captureSeqValue) ||
        !json_is_object(recordValue)) {
        return send_error(fd, "invalid event frame");
    }

    captureSeq = (uint64_t)json_integer_value(captureSeqValue);
    if (captureSeq == 0) {
        return send_error(fd, "capture_seq must be greater than zero");
    }

    if (select_producer(server, producerBootId) != 0) {
        return send_error(fd, "unable to persist producer state");
    }

    if (captureSeq <= server->acceptedCaptureSeq) {
        if (write_state(server) != 0) {
            return send_error(fd, "unable to persist ingest cursor");
        }
        return send_ack(fd, producerBootId,
                        server->acceptedCaptureSeq);
    }

    record = json_deep_copy(recordValue);
    if (!record) {
        return send_error(fd, "unable to copy event record");
    }

    json_object_set_new(record, "producer_boot_id",
                        json_string(producerBootId));
    json_object_set_new(record, "capture_seq",
                        json_integer((json_int_t)captureSeq));

    set_string_if_present(record, request, "capture_id");
    set_string_if_present(record, request, "space");
    set_string_if_present(record, request, "app");
    set_string_if_present(record, request, "stream");
    set_integer_if_present(record, request, "pid");
    set_integer_if_present(record, request, "app_generation");

    canonicalSeq = 0;
    if (log_store_append(server->store, record, &canonicalSeq) != 0) {
        json_decref(record);
        return send_error(fd, "canonical store append failed");
    }

    json_decref(record);
    server->acceptedCaptureSeq = captureSeq;
    (void)canonicalSeq;

    if (write_state(server) != 0) {
        return send_error(fd, "unable to persist ingest cursor");
    }

    return send_ack(fd, producerBootId, captureSeq);
}

static int handle_request(IngestServer *server, int fd,
                          const char *data, size_t size) {
    json_error_t error;
    json_t *request;
    const char *op;
    int rc;

    request = json_loadb(data, size, 0, &error);
    if (!request || !json_is_object(request)) {
        if (request) {
            json_decref(request);
        }
        return send_error(fd, "request must be a JSON object");
    }

    op = json_string_value(json_object_get(request, "op"));
    if (!op) {
        rc = send_error(fd, "op is required");
    } else if (strcmp(op, "hello") == 0) {
        rc = handle_hello(server, fd, request);
    } else if (strcmp(op, "event") == 0) {
        rc = handle_event(server, fd, request);
    } else {
        rc = send_error(fd, "unsupported operation");
    }

    json_decref(request);
    return rc;
}

static void handle_client(IngestServer *server, int fd) {
    char *buffer;

    buffer = malloc(INGEST_MAX_FRAME + 1U);
    if (!buffer) {
        return;
    }

    while (server->running) {
        ssize_t received;

        received = recv(fd, buffer, INGEST_MAX_FRAME, 0);
        if (received == 0) {
            break;
        }
        if (received < 0) {
            if (errno == EINTR) {
                continue;
            }
            break;
        }

        buffer[received] = '\0';
        if (handle_request(server, fd, buffer,
                           (size_t)received) != 0) {
            break;
        }
    }

    free(buffer);
}

static void wake_listener(const char *socketPath) {
    struct sockaddr_un address;
    size_t pathLen;
    socklen_t addressLen;
    int fd;

    if (!socketPath || !*socketPath) {
        return;
    }

    pathLen = strlen(socketPath);
    if (pathLen >= sizeof(address.sun_path)) {
        return;
    }

    fd = socket(AF_UNIX, SOCK_SEQPACKET, 0);
    if (fd < 0) {
        return;
    }

    memset(&address, 0, sizeof(address));
    address.sun_family = AF_UNIX;
    memcpy(address.sun_path, socketPath, pathLen + 1);
    addressLen = (socklen_t)(offsetof(struct sockaddr_un, sun_path) +
                              pathLen + 1);
    (void)connect(fd, (struct sockaddr *)&address, addressLen);
    close(fd);
}

static void *ingest_thread(void *data) {
    IngestServer *server;

    server = (IngestServer *)data;
    while (server->running) {
        int clientFd;

        clientFd = accept(server->listenFd, NULL, NULL);
        if (clientFd < 0) {
            if (errno == EINTR) {
                continue;
            }
            if (!server->running) {
                break;
            }
            continue;
        }

        server->clientFd = clientFd;
        handle_client(server, clientFd);
        close(clientFd);
        server->clientFd = -1;
    }

    return NULL;
}

IngestServer *ingest_start(const char *socketPath, LogStore *store) {
    IngestServer *server;
    struct sockaddr_un address;
    const char *stateDir;
    size_t pathLen;
    socklen_t addressLen;

    if (!socketPath || !*socketPath || !store) {
        return NULL;
    }

    pathLen = strlen(socketPath);
    if (pathLen >= sizeof(address.sun_path)) {
        errno = ENAMETOOLONG;
        return NULL;
    }

    stateDir = log_store_state_dir(store);
    if (!stateDir || !*stateDir || ensure_parent_dir(socketPath) != 0) {
        return NULL;
    }

    server = calloc(1, sizeof(*server));
    if (!server) {
        return NULL;
    }

    memcpy(server->socketPath, socketPath, pathLen + 1);
    if (snprintf(server->statePath, sizeof(server->statePath), "%s/%s",
                 stateDir, INGEST_STATE_FILE) >=
        (int)sizeof(server->statePath)) {
        free(server);
        return NULL;
    }

    load_state(server);
    server->listenFd = socket(AF_UNIX, SOCK_SEQPACKET, 0);
    server->clientFd = -1;
    if (server->listenFd < 0) {
        free(server);
        return NULL;
    }

    memset(&address, 0, sizeof(address));
    address.sun_family = AF_UNIX;
    memcpy(address.sun_path, socketPath, pathLen + 1);
    addressLen = (socklen_t)(offsetof(struct sockaddr_un, sun_path) +
                              pathLen + 1);

    unlink(socketPath);
    if (bind(server->listenFd, (struct sockaddr *)&address,
             addressLen) != 0 ||
        chmod(socketPath, 0660) != 0 ||
        listen(server->listenFd, INGEST_BACKLOG) != 0) {
        close(server->listenFd);
        unlink(socketPath);
        free(server);
        return NULL;
    }

    server->store = store;
    server->running = true;

    if (pthread_create(&server->thread, NULL, ingest_thread, server) != 0) {
        close(server->listenFd);
        unlink(socketPath);
        free(server);
        return NULL;
    }

    return server;
}

void ingest_stop(IngestServer *server) {
    if (!server) {
        return;
    }

    server->running = false;
    (void)write_state(server);
    if (server->clientFd >= 0) {
        shutdown(server->clientFd, SHUT_RDWR);
    }
    wake_listener(server->socketPath);
    shutdown(server->listenFd, SHUT_RDWR);
    pthread_join(server->thread, NULL);
    close(server->listenFd);
    unlink(server->socketPath);
    free(server);
}

const char *ingest_socket_path(const IngestServer *server) {
    return server ? server->socketPath : NULL;
}
