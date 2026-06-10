/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */


#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/un.h>

#include "server.h"
#include "usys_log.h"

#define CTRL_SERVER_READ_BUF 8192

static bool read_line(int fd, char *buf, size_t size) {
    size_t off;
    ssize_t n;
    char c;

    off = 0;
    while (off < size - 1) {
        n = read(fd, &c, 1);
        if (n <= 0) {
            return false;
        }

        if (c == '\n') {
            break;
        }

        buf[off++] = c;
    }

    buf[off] = '\0';
    return true;
}

static bool write_all(int fd, const char *data, size_t len) {
    size_t off;
    ssize_t n;

    off = 0;
    while (off < len) {
        n = write(fd, data + off, len - off);
        if (n <= 0) {
            return false;
        }
        off += (size_t)n;
    }

    return true;
}

static bool send_response(int fd, CtrlResponse *response) {
    JsonObj *json;
    char *body;
    bool ok;

    json = ctrl_response_to_json(response);
    body = json_dumps(json, JSON_COMPACT);
    json_decref(json);

    if (body == NULL) {
        return false;
    }

    ok = write_all(fd, body, strlen(body));
    ok = ok && write_all(fd, "\n", 1);
    free(body);

    return ok;
}

static void handle_client(int fd, Backend *backend) {
    char buf[CTRL_SERVER_READ_BUF];
    JsonObj *json;
    json_error_t error;
    CtrlRequest request;
    CtrlResponse response;

    memset(&request, 0, sizeof(request));
    ctrl_response_init(&response, "req");

    if (!read_line(fd, buf, sizeof(buf))) {
        ctrl_response_set_error(&response,
                                CtrlCodeInvalidRequest,
                                "failed to read request");
        send_response(fd, &response);
        ctrl_response_free(&response);
        return;
    }

    json = json_loads(buf, 0, &error);
    if (json == NULL) {
        ctrl_response_set_error(&response,
                                CtrlCodeInvalidRequest,
                                "invalid json request");
        send_response(fd, &response);
        ctrl_response_free(&response);
        return;
    }

    if (!ctrl_request_from_json(json, &request)) {
        ctrl_response_set_error(&response,
                                CtrlCodeInvalidRequest,
                                "invalid request envelope");
        json_decref(json);
        send_response(fd, &response);
        ctrl_response_free(&response);
        return;
    }

    ctrl_response_free(&response);
    ctrl_response_init(&response, request.id);

    if (!backend_execute(backend, &request, &response)) {
        if (response.payload == NULL) {
            ctrl_response_set_error(&response,
                                    CtrlCodeTransportError,
                                    "backend execution failed");
        }
    }

    send_response(fd, &response);
    ctrl_request_free(&request);
    ctrl_response_free(&response);
    json_decref(json);
}

bool ctrl_server_run(Config *config, Backend *backend, volatile bool *running) {
    int serverFd;
    int clientFd;
    struct sockaddr_un addr;

    if (config == NULL || backend == NULL || running == NULL) {
        return false;
    }

    unlink(config->socketPath);
    serverFd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (serverFd < 0) {
        usys_log_error("failed to create unix socket: %s", strerror(errno));
        return false;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    snprintf(addr.sun_path, sizeof(addr.sun_path), "%s", config->socketPath);

    if (bind(serverFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        usys_log_error("failed to bind %s: %s",
                       config->socketPath,
                       strerror(errno));
        close(serverFd);
        return false;
    }

    if (listen(serverFd, 8) < 0) {
        usys_log_error("failed to listen on %s: %s",
                       config->socketPath,
                       strerror(errno));
        close(serverFd);
        return false;
    }

    usys_log_info("aisg-ctrl listening on %s", config->socketPath);
    while (*running) {
        clientFd = accept(serverFd, NULL, NULL);
        if (clientFd < 0) {
            if (errno == EINTR) {
                continue;
            }
            usys_log_error("accept failed: %s", strerror(errno));
            break;
        }

        handle_client(clientFd, backend);
        close(clientFd);
    }

    close(serverFd);
    unlink(config->socketPath);

    return true;
}
