/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_TFTP_SERVER_H
#define SWITCHD_TFTP_SERVER_H

#include <stdbool.h>
#include <pthread.h>

typedef struct {
    char bindIp[64];
    int port;
    char root[256];
    char filename[128];
    bool running;
    bool terminate;
    pthread_t thread;
} TftpServer;

int tftp_server_start(TftpServer *srv, const char *bindIp, int port,
                      const char *root, const char *filename);
void tftp_server_stop(TftpServer *srv);

#endif
