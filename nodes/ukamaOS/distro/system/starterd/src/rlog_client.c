/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include "rlog_client.h"

#include <errno.h>
#include <limits.h>
#include <poll.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/time.h>
#include <sys/un.h>
#include <time.h>
#include <unistd.h>

#define RLOG_REPLY_MAX 4096
#define RLOG_ACK_TIMEOUT_MS 500

struct RlogClient {
    char socketPath[PATH_MAX];
    char producerBootId[128];
    int reconnectMs;
    int fd;
    int64_t nextRetryMs;
};

static int64_t monotonic_ms(void) {
    struct timespec now;

    clock_gettime(CLOCK_MONOTONIC, &now);
    return ((int64_t)now.tv_sec * 1000LL) +
           ((int64_t)now.tv_nsec / 1000000LL);
}

static void client_close(RlogClient *client) {
    if (!client) return;
    if (client->fd >= 0) close(client->fd);
    client->fd = -1;
    client->nextRetryMs = monotonic_ms() + client->reconnectMs;
}

static bool send_json(int fd, json_t *json) {
    char *data;
    size_t len;
    ssize_t sent;

    data = json_dumps(json, JSON_COMPACT | JSON_ENSURE_ASCII);
    if (!data) return false;

    len = strlen(data);
    sent = send(fd, data, len, MSG_NOSIGNAL);
    free(data);

    return sent == (ssize_t)len;
}

static json_t *recv_json(int fd) {
    struct pollfd pfd;
    char buffer[RLOG_REPLY_MAX + 1];
    ssize_t received;
    json_error_t error;

    memset(&pfd, 0, sizeof(pfd));
    pfd.fd = fd;
    pfd.events = POLLIN;

    if (poll(&pfd, 1, RLOG_ACK_TIMEOUT_MS) <= 0) {
        return NULL;
    }

    received = recv(fd, buffer, RLOG_REPLY_MAX, 0);
    if (received <= 0) return NULL;

    buffer[received] = '\0';
    return json_loadb(buffer, (size_t)received, 0, &error);
}

static bool client_connect(RlogClient *client) {
    struct sockaddr_un address;
    json_t *hello;
    json_t *reply;
    const char *op;

    if (!client) return false;
    if (client->fd >= 0) return true;
    if (monotonic_ms() < client->nextRetryMs) return false;

    client->fd = socket(AF_UNIX, SOCK_SEQPACKET, 0);
    if (client->fd < 0) {
        client_close(client);
        return false;
    }

    memset(&address, 0, sizeof(address));
    address.sun_family = AF_UNIX;
    snprintf(address.sun_path, sizeof(address.sun_path), "%s",
             client->socketPath);

    if (connect(client->fd, (struct sockaddr *)&address,
                sizeof(address)) != 0) {
        client_close(client);
        return false;
    }

    hello = json_pack("{s:s,s:s,s:s}",
                      "op", "hello",
                      "producer", "starterd",
                      "producer_boot_id", client->producerBootId);
    if (!hello || !send_json(client->fd, hello)) {
        if (hello) json_decref(hello);
        client_close(client);
        return false;
    }
    json_decref(hello);

    reply = recv_json(client->fd);
    if (!reply) {
        client_close(client);
        return false;
    }

    op = json_string_value(json_object_get(reply, "op"));
    if (!op || strcmp(op, "hello") != 0 ||
        !json_is_true(json_object_get(reply, "ready"))) {
        json_decref(reply);
        client_close(client);
        return false;
    }

    json_decref(reply);
    return true;
}

RlogClient *rlog_client_create(const char *socketPath,
                               const char *producerBootId,
                               int reconnectMs) {
    RlogClient *client;

    if (!socketPath || !*socketPath ||
        !producerBootId || !*producerBootId) {
        return NULL;
    }

    client = calloc(1, sizeof(*client));
    if (!client) return NULL;

    snprintf(client->socketPath, sizeof(client->socketPath), "%s",
             socketPath);
    snprintf(client->producerBootId,
             sizeof(client->producerBootId), "%s", producerBootId);
    client->reconnectMs = reconnectMs > 0 ? reconnectMs : 1000;
    client->fd = -1;
    client->nextRetryMs = 0;

    return client;
}

void rlog_client_destroy(RlogClient *client) {
    if (!client) return;
    if (client->fd >= 0) close(client->fd);
    free(client);
}

bool rlog_client_send(RlogClient *client,
                      json_t *record,
                      const char *space,
                      const char *app,
                      int pid,
                      uint32_t generation,
                      const char *stream,
                      uint64_t captureSeq,
                      const char *captureId) {
    json_t *frame;
    json_t *reply;
    const char *op;
    json_t *accepted;
    bool ok;

    if (!client || !record || !app || !stream || !captureId) {
        return false;
    }

    if (!client_connect(client)) return false;

    frame = json_pack("{s:s,s:s,s:I,s:s,s:s,s:s,s:i,s:i,s:s,s:o}",
                      "op", "event",
                      "producer_boot_id", client->producerBootId,
                      "capture_seq", (json_int_t)captureSeq,
                      "capture_id", captureId,
                      "space", space ? space : "",
                      "app", app,
                      "pid", pid,
                      "app_generation", generation,
                      "stream", stream,
                      "record", json_incref(record));
    if (!frame) return false;

    ok = send_json(client->fd, frame);
    json_decref(frame);
    if (!ok) {
        client_close(client);
        return false;
    }

    reply = recv_json(client->fd);
    if (!reply) {
        client_close(client);
        return false;
    }

    op = json_string_value(json_object_get(reply, "op"));
    accepted = json_object_get(reply, "accepted_through");
    ok = op && strcmp(op, "ack") == 0 &&
         json_is_integer(accepted) &&
         (uint64_t)json_integer_value(accepted) >= captureSeq;

    json_decref(reply);
    if (!ok) client_close(client);
    return ok;
}
