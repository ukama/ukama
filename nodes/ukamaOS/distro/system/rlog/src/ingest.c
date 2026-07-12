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
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <sys/un.h>
#include <unistd.h>

#define INGEST_MAX_FRAME (32U * 1024U)
#define INGEST_BACKLOG   4

struct IngestServer {
    char socketPath[PATH_MAX];
    int listenFd;
    int clientFd;
    pthread_t thread;
    bool running;
    LogStore *store;

    char producerBootId[128];
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

static int handle_hello(IngestServer *server, int fd, json_t *request) {
    const char *producerBootId;
    json_t *reply;
    int rc;

    producerBootId = json_string_value(
        json_object_get(request, "producer_boot_id"));
    if (!producerBootId || !*producerBootId) {
        return send_error(fd, "producer_boot_id is required");
    }

    if (strcmp(server->producerBootId, producerBootId) != 0) {
        snprintf(server->producerBootId,
                 sizeof(server->producerBootId), "%s", producerBootId);
        server->acceptedCaptureSeq = 0;
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

    if (strcmp(server->producerBootId, producerBootId) != 0) {
        snprintf(server->producerBootId,
                 sizeof(server->producerBootId), "%s", producerBootId);
        server->acceptedCaptureSeq = 0;
    }

    if (captureSeq <= server->acceptedCaptureSeq) {
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

    if (!socketPath || !*socketPath || !store ||
        strlen(socketPath) >= sizeof(address.sun_path)) {
        return NULL;
    }

    if (ensure_parent_dir(socketPath) != 0) {
        return NULL;
    }

    server = calloc(1, sizeof(*server));
    if (!server) {
        return NULL;
    }

    snprintf(server->socketPath, sizeof(server->socketPath), "%s",
             socketPath);
    server->listenFd = socket(AF_UNIX, SOCK_SEQPACKET, 0);
    server->clientFd = -1;
    if (server->listenFd < 0) {
        free(server);
        return NULL;
    }

    memset(&address, 0, sizeof(address));
    address.sun_family = AF_UNIX;
    snprintf(address.sun_path, sizeof(address.sun_path), "%s",
             socketPath);

    unlink(socketPath);
    if (bind(server->listenFd, (struct sockaddr *)&address,
             sizeof(address)) != 0 ||
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
    if (server->clientFd >= 0) {
        shutdown(server->clientFd, SHUT_RDWR);
    }
    shutdown(server->listenFd, SHUT_RDWR);
    close(server->listenFd);
    pthread_join(server->thread, NULL);
    unlink(server->socketPath);
    free(server);
}

const char *ingest_socket_path(const IngestServer *server) {
    return server ? server->socketPath : NULL;
}
