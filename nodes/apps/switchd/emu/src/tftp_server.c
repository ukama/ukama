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
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#include "model.h"
#include "tftp_server.h"

static void *tftp_main(void *arg) {
    EmuModel *model = (EmuModel *)arg;
    char buf[EMU_TFTP_BUF] = {0};

    while (model->running) {
        struct sockaddr_in cli;
        socklen_t slen = sizeof(cli);
        ssize_t bytes = recvfrom(model->tftpFd, buf, sizeof(buf) - 1, 0,
                                 (struct sockaddr *)&cli, &slen);
        if (bytes < 0) {
            if (errno == EINTR) {
                continue;
            }
            if (!model->running) {
                break;
            }
            continue;
        }

        buf[bytes] = '\0';
        if (model->faults.tftpFail) {
            continue;
        }

        if (strncmp(buf, "WRQ ", 4) == 0) {
            char fileName[128] = {0};
            (void)sscanf(buf + 4, "%127s", fileName);

            pthread_mutex_lock(&model->lock);
            model_stage_firmware(model, "/tmp", fileName);
            pthread_mutex_unlock(&model->lock);

            sendto(model->tftpFd, "ACK\n", 4, 0,
                   (struct sockaddr *)&cli, slen);
        } else if (strncmp(buf, "DATA ", 5) == 0) {
            sendto(model->tftpFd, "ACK\n", 4, 0,
                   (struct sockaddr *)&cli, slen);
        }
    }

    return NULL;
}

int tftp_server_start(EmuModel *model) {
    struct sockaddr_in addr;

    model->tftpFd = socket(AF_INET, SOCK_DGRAM, 0);
    if (model->tftpFd < 0) {
        return STATUS_NOK;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port   = htons((uint16_t)model->cfg.tftpPort);
    addr.sin_addr.s_addr = inet_addr(model->cfg.bindAddr);

    if (bind(model->tftpFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        return STATUS_NOK;
    }

    return pthread_create(&model->tftpThread, NULL, tftp_main, model);
}

void tftp_server_stop(EmuModel *model) {
    if (model->tftpFd >= 0) {
        close(model->tftpFd);
        model->tftpFd = -1;
    }

    if (model->tftpThread != 0U) {
        pthread_join(model->tftpThread, NULL);
        model->tftpThread = 0;
    }
}
