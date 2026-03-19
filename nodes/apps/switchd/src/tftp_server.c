/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>
#include <errno.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/select.h>
#include <sys/time.h>
#include <pthread.h>
#include "tftp_server.h"
#include "types.h"
#include "log.h"

#define TFTP_OPCODE_RRQ   1
#define TFTP_OPCODE_DATA  3
#define TFTP_OPCODE_ACK   4
#define TFTP_OPCODE_ERROR 5
#define TFTP_BLKSZ        512

static void send_error(int fd, struct sockaddr_in *peer, socklen_t plen, uint16_t code, const char *msg) {
    uint8_t pkt[516];
    size_t n = strlen(msg);
    pkt[0] = 0; pkt[1] = TFTP_OPCODE_ERROR;
    pkt[2] = (uint8_t)(code >> 8); pkt[3] = (uint8_t)(code & 0xFF);
    memcpy(pkt + 4, msg, n);
    pkt[4 + n] = '\0';
    sendto(fd, pkt, n + 5, 0, (struct sockaddr *)peer, plen);
}

static int serve_file(int fd, struct sockaddr_in *peer, socklen_t plen, const char *path) {
    FILE *fp;
    uint8_t pkt[4 + TFTP_BLKSZ];
    uint8_t ack[4];
    uint16_t blk = 1;
    size_t n;
    fd_set rfds;
    struct timeval tv;
    ssize_t r;

    fp = fopen(path, "rb");
    if (fp == NULL) {
        send_error(fd, peer, plen, 1, "File not found");
        return -1;
    }
    for (;;) {
        pkt[0] = 0; pkt[1] = TFTP_OPCODE_DATA;
        pkt[2] = (uint8_t)(blk >> 8); pkt[3] = (uint8_t)(blk & 0xFF);
        n = fread(pkt + 4, 1, TFTP_BLKSZ, fp);
        if (sendto(fd, pkt, n + 4, 0, (struct sockaddr *)peer, plen) < 0) {
            fclose(fp);
            return -1;
        }
        FD_ZERO(&rfds);
        FD_SET(fd, &rfds);
        tv.tv_sec = 2;
        tv.tv_usec = 0;
        r = select(fd + 1, &rfds, NULL, NULL, &tv);
        if (r <= 0) {
            fclose(fp);
            return -1;
        }
        r = recvfrom(fd, ack, sizeof(ack), 0, NULL, NULL);
        if (r < 4 || ack[1] != TFTP_OPCODE_ACK) {
            fclose(fp);
            return -1;
        }
        if (n < TFTP_BLKSZ) break;
        blk++;
    }
    fclose(fp);
    return 0;
}

static void *tftp_main(void *arg) {
    TftpServer *srv = (TftpServer *)arg;
    int fd;
    struct sockaddr_in addr;
    struct sockaddr_in peer;
    socklen_t plen;
    uint8_t req[600];
    ssize_t n;
    char path[512];

    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) {
        srv->running = false;
        return NULL;
    }
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons((uint16_t)srv->port);
    addr.sin_addr.s_addr = inet_addr(srv->bindIp);
    if (bind(fd, (struct sockaddr *)&addr, sizeof(addr)) != 0) {
        close(fd);
        srv->running = false;
        return NULL;
    }
    srv->running = true;
    while (!srv->terminate) {
        plen = sizeof(peer);
        n = recvfrom(fd, req, sizeof(req), 0, (struct sockaddr *)&peer, &plen);
        if (n <= 0) {
            if (errno == EINTR) continue;
            break;
        }
        if (n < 4 || req[1] != TFTP_OPCODE_RRQ) {
            send_error(fd, &peer, plen, 4, "Bad request");
            continue;
        }
        snprintf(path, sizeof(path), "%s/%s", srv->root, srv->filename);
        if (serve_file(fd, &peer, plen, path) != 0) {
            log_warn("tftp serve failed for %s", path);
        }
    }
    close(fd);
    srv->running = false;
    return NULL;
}

int tftp_server_start(TftpServer *srv, const char *bindIp, int port,
                      const char *root, const char *filename) {
    memset(srv, 0, sizeof(*srv));
    snprintf(srv->bindIp, sizeof(srv->bindIp), "%s", bindIp);
    snprintf(srv->root, sizeof(srv->root), "%s", root);
    snprintf(srv->filename, sizeof(srv->filename), "%s", filename);
    srv->port = port;
    srv->terminate = false;
    if (pthread_create(&srv->thread, NULL, tftp_main, srv) != 0) return SWITCHD_ERR_INTERNAL;
    { struct timespec ts; ts.tv_sec = 0; ts.tv_nsec = 150000000L; nanosleep(&ts, NULL); }
    return srv->running ? SWITCHD_OK : SWITCHD_ERR_IO;
}

void tftp_server_stop(TftpServer *srv) {
    if (!srv->running) return;
    srv->terminate = true;
    pthread_cancel(srv->thread);
    pthread_join(srv->thread, NULL);
    srv->running = false;
}
