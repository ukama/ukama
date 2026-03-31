/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "http_status.h"
#include "network.h"
#include "utils.h"
#include "web_service.h"

static void *http_main(void *arg) {
    EmuModel *model = (EmuModel *)arg;

    while (model->running) {
        int clientFd = -1;
        struct sockaddr_in clientAddr;
        socklen_t clientLen = sizeof(clientAddr);
        char req[EMU_HTTP_REQ_BUF] = {0};
        char method[8] = {0};
        char path[256] = {0};
        char body[2048] = {0};
        char resp[EMU_HTTP_RESP_BUF] = {0};
        char header[256] = {0};
        char *bodyStart = NULL;
        int status = HttpStatus_OK;
        ssize_t bytes = 0;

        clientFd = accept(model->httpFd,
                          (struct sockaddr *)&clientAddr,
                          &clientLen);
        if (clientFd < 0) {
            if (errno == EINTR) {
                continue;
            }
            if (!model->running) {
                break;
            }
            continue;
        }

        bytes = read(clientFd, req, sizeof(req) - 1);
        if (bytes <= 0) {
            close(clientFd);
            continue;
        }

        req[bytes] = '\0';
        (void)sscanf(req, "%7s %255s", method, path);

        bodyStart = strstr(req, "\r\n\r\n");
        if (bodyStart != NULL) {
            snprintf(body, sizeof(body), "%s", bodyStart + 4);
        }

        web_service_handle(model, method, path, body,
                           resp, sizeof(resp), &status);

        snprintf(header, sizeof(header),
                 "HTTP/1.1 %d %s\r\n"
                 "Content-Type: application/json\r\n"
                 "Content-Length: %zu\r\n"
                 "Connection: close\r\n\r\n",
                 status, HttpStatusStr(status), strlen(resp));

        write_all(clientFd, header, strlen(header));
        write_all(clientFd, resp, strlen(resp));
        close(clientFd);
    }

    return NULL;
}

int network_start(EmuModel *model) {
    struct sockaddr_in addr;
    int yes = 1;

    model->httpFd = socket(AF_INET, SOCK_STREAM, 0);
    if (model->httpFd < 0) {
        return STATUS_NOK;
    }

    setsockopt(model->httpFd, SOL_SOCKET, SO_REUSEADDR, &yes, sizeof(yes));

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port   = htons((uint16_t)model->cfg.httpPort);
    addr.sin_addr.s_addr = inet_addr(model->cfg.bindAddr);

    if (bind(model->httpFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        return STATUS_NOK;
    }

    if (listen(model->httpFd, 16) < 0) {
        return STATUS_NOK;
    }

    return pthread_create(&model->httpThread, NULL, http_main, model);
}

void network_stop(EmuModel *model) {
    if (model->httpFd >= 0) {
        close(model->httpFd);
        model->httpFd = -1;
    }

    if (model->httpThread != 0U) {
        pthread_join(model->httpThread, NULL);
        model->httpThread = 0;
    }
}
