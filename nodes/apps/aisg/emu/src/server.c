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
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

#include "ops.h"
#include "request.h"
#include "response.h"
#include "server.h"
#include "usys_log.h"

#define EMU_SERVER_READ_BUF 8192
#define EMU_SERVER_BACKLOG  8

static bool read_line(int fd, char *buf, size_t size)
{
    size_t off = 0;
    ssize_t n;
    char c;

    if (buf == NULL || size == 0) {
        return false;
    }

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

static bool write_all(int fd, const char *data, size_t len)
{
    size_t off = 0;
    ssize_t n;

    if (data == NULL) {
        return false;
    }

    while (off < len) {
        n = write(fd, data + off, len - off);
        if (n <= 0) {
            return false;
        }

        off += (size_t)n;
    }

    return true;
}

static bool send_response(int fd, EmuResponse *response)
{
    JsonObj *json = NULL;
    char *body = NULL;
    bool ok;

    json = emu_response_to_json(response);
    if (json == NULL) {
        return false;
    }

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

static bool parse_request_line(const char *line, EmuRequest *request)
{
    JsonObj *json = NULL;
    json_error_t error;
    bool ok;

    if (line == NULL || request == NULL) {
        return false;
    }

    json = json_loads(line, 0, &error);
    if (json == NULL) {
        usys_log_error("invalid emu request: %s", error.text);
        return false;
    }

    ok = emu_request_from_json(json, request);
    json_decref(json);

    return ok;
}

static void handle_invalid_request(int fd, const char *reason)
{
    EmuResponse response;

    emu_response_init(&response, "req");
    emu_response_set_error(&response,
                           "InvalidRequest",
                           reason ? reason : "invalid request");
    send_response(fd, &response);
    emu_response_free(&response);
}

static void handle_client(int fd, EmuModel *model)
{
    char line[EMU_SERVER_READ_BUF];
    EmuRequest request;
    EmuResponse response;

    emu_request_init(&request);
    emu_response_init(&response, "req");

    if (!read_line(fd, line, sizeof(line))) {
        handle_invalid_request(fd, "failed to read request");
        return;
    }

    if (!parse_request_line(line, &request)) {
        handle_invalid_request(fd, "invalid request envelope");
        return;
    }

    emu_response_free(&response);
    emu_response_init(&response, request.id);

    if (!emu_ops_handle(model, &request, &response)) {
        emu_response_set_error(&response,
                               "TransportError",
                               "emulator execution failed");
    }

    send_response(fd, &response);

    emu_request_free(&request);
    emu_response_free(&response);
}

static bool bind_server_socket(int serverFd, const char *socketPath)
{
    struct sockaddr_un addr;

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    snprintf(addr.sun_path, sizeof(addr.sun_path), "%s", socketPath);

    if (bind(serverFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        usys_log_error("failed to bind %s: %s", socketPath, strerror(errno));
        return false;
    }

    return true;
}

bool emu_server_run(EmuConfig *config, EmuModel *model, volatile bool *running)
{
    int serverFd;
    int clientFd;

    if (config == NULL || model == NULL || running == NULL) {
        return false;
    }

    unlink(config->socketPath);

    serverFd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (serverFd < 0) {
        usys_log_error("failed to create unix socket: %s", strerror(errno));
        return false;
    }

    if (!bind_server_socket(serverFd, config->socketPath)) {
        close(serverFd);
        return false;
    }

    if (listen(serverFd, EMU_SERVER_BACKLOG) < 0) {
        usys_log_error("failed to listen on %s: %s",
                       config->socketPath,
                       strerror(errno));
        close(serverFd);
        return false;
    }

    usys_log_info("aisg-emu listening on %s", config->socketPath);

    while (*running) {
        clientFd = accept(serverFd, NULL, NULL);
        if (clientFd < 0) {
            if (errno == EINTR) {
                continue;
            }

            usys_log_error("accept failed: %s", strerror(errno));
            break;
        }

        handle_client(clientFd, model);
        close(clientFd);
    }

    close(serverFd);
    unlink(config->socketPath);

    return true;
}
