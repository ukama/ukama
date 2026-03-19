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
#include <time.h>
#include <unistd.h>

#include "oid_map.h"
#include "snmp_agent.h"

static void *snmp_main(void *arg) {
    EmuModel *model = (EmuModel *)arg;
    char buf[EMU_SNMP_BUF] = {0};

    while (model->running) {
        struct sockaddr_in cli;
        socklen_t slen = sizeof(cli);
        ssize_t bytes = recvfrom(model->snmpFd, buf, sizeof(buf) - 1, 0,
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
        if (model->faults.unreachable || !model->info.reachable) {
            continue;
        }

        if (model->faults.snmpDelayMs > 0) {
            struct timespec ts;
            ts.tv_sec = model->faults.snmpDelayMs / 1000;
            ts.tv_nsec = (long)(model->faults.snmpDelayMs % 1000) * 1000000L;
            nanosleep(&ts, NULL);
        }

        {
            char cmd[8] = {0};
            char oid[256] = {0};
            char out[512] = {0};
            int iv = 0;

            if (sscanf(buf, "%7s %255s %d", cmd, oid, &iv) >= 2) {
                if (strcmp(cmd, "GET") == 0) {
                    int v = 0;
                    char sval[256] = {0};

                    pthread_mutex_lock(&model->lock);
                    if (oid_get_int(model, oid, &v) == STATUS_OK) {
                        (void)snprintf(out, sizeof(out), "INT %.255s %d\n", oid, v);
                    } else if (oid_get_string(model, oid, sval,
                                              sizeof(sval)) == STATUS_OK) {
                        (void)snprintf(out, sizeof(out), "STR %.200s %.200s\n", oid, sval);
                    } else {
                        snprintf(out, sizeof(out), "ERROR noSuchName\n");
                    }
                    pthread_mutex_unlock(&model->lock);
                } else if (strcmp(cmd, "SET") == 0) {
                    pthread_mutex_lock(&model->lock);
                    if (oid_set_int(model, oid, iv) == STATUS_OK) {
                        (void)snprintf(out, sizeof(out), "OK %.255s %d\n", oid, iv);
                    } else {
                        snprintf(out, sizeof(out), "ERROR setFailed\n");
                    }
                    pthread_mutex_unlock(&model->lock);
                } else {
                    snprintf(out, sizeof(out), "ERROR badCmd\n");
                }

                sendto(model->snmpFd, out, strlen(out), 0,
                       (struct sockaddr *)&cli, slen);
            }
        }
    }

    return NULL;
}

int snmp_agent_start(EmuModel *model) {
    struct sockaddr_in addr;

    model->snmpFd = socket(AF_INET, SOCK_DGRAM, 0);
    if (model->snmpFd < 0) {
        return STATUS_NOK;
    }

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port   = htons((uint16_t)model->cfg.snmpPort);
    addr.sin_addr.s_addr = inet_addr(model->cfg.bindAddr);

    if (bind(model->snmpFd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
        return STATUS_NOK;
    }

    return pthread_create(&model->snmpThread, NULL, snmp_main, model);
}

void snmp_agent_stop(EmuModel *model) {
    if (model->snmpFd >= 0) {
        close(model->snmpFd);
        model->snmpFd = -1;
    }

    if (model->snmpThread != 0U) {
        pthread_join(model->snmpThread, NULL);
        model->snmpThread = 0;
    }
}
