/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <sys/time.h>

#include "controller_client.h"
#include "usys_log.h"

static bool write_all(int fd, const char *data, size_t len) {
    size_t off;
    ssize_t n;

    off = 0;
    while (off < len) {
        n = write(fd, data + off, len - off);
        if (n <= 0) return false;
        off += (size_t)n;
    }

    return true;
}

static bool read_line(int fd, char *buf, size_t size) {
    size_t off;
    ssize_t n;
    char c;

    off = 0;
    while (off < size - 1) {
        n = read(fd, &c, 1);
        if (n <= 0) return false;
        if (c == '\n') break;
        buf[off++] = c;
    }

    buf[off] = '\0';
    return true;
}

static int open_socket(ControllerClient *client) {
    int fd;
    struct sockaddr_un addr;
    struct timeval tv;

    fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (fd < 0) return -1;

    tv.tv_sec = client->timeoutMs / 1000;
    tv.tv_usec = (client->timeoutMs % 1000) * 1000;
    setsockopt(fd, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));
    setsockopt(fd, SOL_SOCKET, SO_SNDTIMEO, &tv, sizeof(tv));

    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    snprintf(addr.sun_path, sizeof(addr.sun_path), "%s", client->path);

    if (connect(fd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        close(fd);
        return -1;
    }

    return fd;
}

void controller_client_init(ControllerClient *client, Config *config) {
    if (client == NULL || config == NULL) return;

    memset(client, 0, sizeof(ControllerClient));
    snprintf(client->path, sizeof(client->path), "%s", config->controllerPath);
    client->timeoutMs = config->controllerTimeoutMs;
}

void ctrl_response_free(CtrlResponse *response) {
    if (response == NULL) return;
    json_decref(response->payload);
    memset(response, 0, sizeof(CtrlResponse));
}

JsonObj *ctrl_response_steal_payload(CtrlResponse *response) {
    JsonObj *payload;

    if (response == NULL) return NULL;

    payload = response->payload;
    response->payload = NULL;
    return payload;
}

static bool parse_response(JsonObj *root, CtrlResponse *response) {
    JsonObj *value;

    if (root == NULL || response == NULL) return false;

    memset(response, 0, sizeof(CtrlResponse));
    value = json_object_get(root, "ok");
    response->ok = json_is_true(value);

    value = json_object_get(root, "code");
    snprintf(response->code,
             sizeof(response->code),
             "%s",
             json_is_string(value) ? json_string_value(value) : "");

    value = json_object_get(root, "reason");
    snprintf(response->reason,
             sizeof(response->reason),
             "%s",
             json_is_string(value) ? json_string_value(value) : "");

    value = json_object_get(root, "payload");
    response->payload = value ? json_deep_copy(value) : json_object();

    return response->payload != NULL;
}

bool controller_client_call(ControllerClient *client,
                            const char *type,
                            JsonObj *payload,
                            CtrlResponse *response) {
    int fd;
    char *wire;
    char line[8192];
    JsonObj *req;
    JsonObj *reply;
    json_error_t error;
    bool ok;

    if (client == NULL || type == NULL || response == NULL) {
        json_decref(payload);
        return false;
    }

    req = json_object();
    json_object_set_new(req, "id", json_string("req-1"));
    json_object_set_new(req, "type", json_string(type));
    json_object_set_new(req, "payload", payload ? payload : json_object());

    wire = json_dumps(req, JSON_COMPACT);
    json_decref(req);
    if (wire == NULL) return false;

    fd = open_socket(client);
    if (fd < 0) {
        usys_log_error("failed to connect controller socket %s", client->path);
        free(wire);
        return false;
    }

    ok = write_all(fd, wire, strlen(wire));
    ok = ok && write_all(fd, "\n", 1);
    free(wire);

    if (!ok || !read_line(fd, line, sizeof(line))) {
        close(fd);
        return false;
    }
    close(fd);

    reply = json_loads(line, 0, &error);
    if (reply == NULL) {
        usys_log_error("invalid controller response: %s", error.text);
        return false;
    }

    ok = parse_response(reply, response);
    json_decref(reply);

    return ok;
}
