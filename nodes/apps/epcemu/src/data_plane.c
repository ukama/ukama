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
#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <sys/select.h>
#include <sys/socket.h>
#include <unistd.h>

#include "data_plane.h"
#include "tun.h"
#include "ue.h"

#define IPV4_MIN_HEADER 20

static uint16_t ipv4_total_len(const unsigned char *pkt, ssize_t len) {

    if (pkt == NULL || len < IPV4_MIN_HEADER) return 0;

    return ((uint16_t)pkt[2] << 8) | pkt[3];
}

static int ipv4_src_dst(const unsigned char *pkt,
                        ssize_t len,
                        char *src,
                        size_t srcLen,
                        char *dst,
                        size_t dstLen) {

    struct in_addr saddr;
    struct in_addr daddr;
    uint8_t version;

    if (pkt == NULL || len < IPV4_MIN_HEADER ||
        src == NULL || dst == NULL) {
        return USYS_FALSE;
    }

    version = pkt[0] >> 4;
    if (version != 4) {
        return USYS_FALSE;
    }

    memcpy(&saddr, pkt + 12, sizeof(saddr));
    memcpy(&daddr, pkt + 16, sizeof(daddr));

    if (inet_ntop(AF_INET, &saddr, src, srcLen) == NULL) {
        return USYS_FALSE;
    }

    if (inet_ntop(AF_INET, &daddr, dst, dstLen) == NULL) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

static void record_drop(DataPlane *dp, ssize_t bytes) {

    pthread_mutex_lock(&dp->lock);
    dp->droppedPackets++;
    if (bytes > 0) dp->droppedBytes += (uint64_t)bytes;
    pthread_mutex_unlock(&dp->lock);
}

static void record_uplink(DataPlane *dp, ssize_t bytes) {

    pthread_mutex_lock(&dp->lock);
    dp->uplinkPackets++;
    if (bytes > 0) dp->uplinkBytes += (uint64_t)bytes;
    pthread_mutex_unlock(&dp->lock);
}

static void record_downlink(DataPlane *dp, ssize_t bytes) {

    pthread_mutex_lock(&dp->lock);
    dp->downlinkPackets++;
    if (bytes > 0) dp->downlinkBytes += (uint64_t)bytes;
    pthread_mutex_unlock(&dp->lock);
}

static int udp_socket_create(int port) {

    struct sockaddr_in addr;
    int fd;
    int opt;

    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) {
        usys_log_error("failed to create UDP socket: %s", strerror(errno));
        return -1;
    }

    opt = 1;
    setsockopt(fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons((uint16_t)port);
    addr.sin_addr.s_addr = htonl(INADDR_ANY);

    if (bind(fd, (struct sockaddr *)&addr, sizeof(addr)) != 0) {
        usys_log_error("failed to bind UDP data port %d: %s",
                       port, strerror(errno));
        close(fd);
        return -1;
    }

    return fd;
}

static void handle_udp_to_tun(DataPlane *dp) {

    unsigned char pkt[EPCEMU_MAX_PACKET];
    struct sockaddr_storage peer;
    socklen_t peerLen;
    ssize_t n;
    char src[64];
    char dst[64];

    peerLen = sizeof(peer);
    memset(&peer, 0, sizeof(peer));
    memset(src, 0, sizeof(src));
    memset(dst, 0, sizeof(dst));

    n = recvfrom(dp->udpFd,
                 pkt,
                 sizeof(pkt),
                 0,
                 (struct sockaddr *)&peer,
                 &peerLen);
    if (n <= 0) return;

    if (!ipv4_src_dst(pkt, n, src, sizeof(src), dst, sizeof(dst))) {
        record_drop(dp, n);
        return;
    }

    if (!ue_find_by_src_ip(src, &(UeEntry){0})) {
        record_drop(dp, n);
        usys_log_debug("dropping uplink packet from unknown UE src=%s", src);
        return;
    }

    if (!ue_bind_peer(src, &peer, peerLen)) {
        record_drop(dp, n);
        return;
    }

    if (write(dp->tunFd, pkt, (size_t)n) != n) {
        record_drop(dp, n);
        usys_log_error("failed to write uplink packet to tun");
        return;
    }

    record_uplink(dp, n);
    ue_record_uplink(src, (uint64_t)n);
}

static void handle_tun_to_udp(DataPlane *dp) {

    unsigned char pkt[EPCEMU_MAX_PACKET];
    ssize_t n;
    char src[64];
    char dst[64];
    UeEntry ue;

    memset(src, 0, sizeof(src));
    memset(dst, 0, sizeof(dst));
    memset(&ue, 0, sizeof(ue));

    n = read(dp->tunFd, pkt, sizeof(pkt));
    if (n <= 0) return;

    if (!ipv4_src_dst(pkt, n, src, sizeof(src), dst, sizeof(dst))) {
        record_drop(dp, n);
        return;
    }

    if (!ue_find_by_dst_ip(dst, &ue)) {
        record_drop(dp, n);
        usys_log_debug("dropping downlink packet for unknown UE dst=%s", dst);
        return;
    }

    if (!ue.peerSet) {
        record_drop(dp, n);
        usys_log_debug("dropping downlink packet for UE without peer dst=%s", dst);
        return;
    }

    if (sendto(dp->udpFd,
               pkt,
               (size_t)n,
               0,
               (struct sockaddr *)&ue.peerAddr,
               ue.peerLen) != n) {
        record_drop(dp, n);
        usys_log_error("failed to send downlink packet to UE");
        return;
    }

    record_downlink(dp, n);
    ue_record_downlink(dst, (uint64_t)n);
}

static void *data_plane_loop(void *arg) {

    DataPlane *dp;
    fd_set rfds;
    int maxfd;
    int rc;

    dp = (DataPlane *)arg;
    if (dp == NULL) return NULL;

    while (dp->running) {
        FD_ZERO(&rfds);
        FD_SET(dp->udpFd, &rfds);
        FD_SET(dp->tunFd, &rfds);

        maxfd = dp->udpFd > dp->tunFd ? dp->udpFd : dp->tunFd;

        rc = select(maxfd + 1, &rfds, NULL, NULL, NULL);
        if (rc < 0) {
            if (errno == EINTR) continue;
            usys_log_error("data plane select failed: %s", strerror(errno));
            break;
        }

        if (FD_ISSET(dp->udpFd, &rfds)) {
            handle_udp_to_tun(dp);
        }

        if (FD_ISSET(dp->tunFd, &rfds)) {
            handle_tun_to_udp(dp);
        }
    }

    return NULL;
}

int data_plane_start(DataPlane *dp,
                     EpcemuConfig *config,
                     EpcemuStatus *status) {

    if (dp == NULL || config == NULL || status == NULL) {
        return USYS_FALSE;
    }

    memset(dp, 0, sizeof(DataPlane));
    pthread_mutex_init(&dp->lock, NULL);

    status_set(status, EpcemuStateStartingDataPlane,
               "starting data plane");

    if (!tun_configure(config->tunName, config->tunAddr)) {
        status_fail(status, "failed to configure EPC tun interface");
        return USYS_FALSE;
    }

    dp->tunFd = tun_create(config->tunName);
    if (dp->tunFd < 0) {
        status_fail(status, "failed to open EPC tun interface");
        return USYS_FALSE;
    }

    dp->udpFd = udp_socket_create(config->dataPort);
    if (dp->udpFd < 0) {
        tun_close(dp->tunFd);
        status_fail(status, "failed to open epcemu-data UDP socket");
        return USYS_FALSE;
    }

    dp->running = true;
    dp->ready = true;

    if (pthread_create(&dp->thread, NULL, data_plane_loop, dp) != 0) {
        dp->running = false;
        close(dp->udpFd);
        tun_close(dp->tunFd);
        status_fail(status, "failed to create data plane thread");
        return USYS_FALSE;
    }

    config->dataPlaneReady = true;

    usys_log_info("epcemu data plane ready udp=%d tun=%s addr=%s",
                  config->dataPort,
                  config->tunName,
                  config->tunAddr);

    return USYS_TRUE;
}

void data_plane_stop(DataPlane *dp) {

    if (dp == NULL) return;

    dp->running = false;

    if (dp->udpFd >= 0) close(dp->udpFd);
    if (dp->tunFd >= 0) tun_close(dp->tunFd);

    if (dp->ready) {
        pthread_join(dp->thread, NULL);
    }

    pthread_mutex_destroy(&dp->lock);
}

JsonObj *data_plane_json(DataPlane *dp, EpcemuConfig *config) {

    JsonObj *obj;
    uint64_t uplinkPackets;
    uint64_t uplinkBytes;
    uint64_t downlinkPackets;
    uint64_t downlinkBytes;
    uint64_t droppedPackets;
    uint64_t droppedBytes;
    bool ready;

    obj = json_object();
    if (obj == NULL) return NULL;

    if (dp == NULL) {
        json_object_set_new(obj, "enabled", json_false());
        return obj;
    }

    pthread_mutex_lock(&dp->lock);
    ready = dp->ready;
    uplinkPackets = dp->uplinkPackets;
    uplinkBytes = dp->uplinkBytes;
    downlinkPackets = dp->downlinkPackets;
    downlinkBytes = dp->downlinkBytes;
    droppedPackets = dp->droppedPackets;
    droppedBytes = dp->droppedBytes;
    pthread_mutex_unlock(&dp->lock);

    json_object_set_new(obj, "enabled", json_boolean(ready));
    json_object_set_new(obj, "mode", json_string("tun-udp"));
    json_object_set_new(obj, "service", json_string(EPCEMU_DATA_SERVICE_NAME));
    json_object_set_new(obj, "udpPort",
                        json_integer(config ? config->dataPort : 0));
    json_object_set_new(obj, "tun",
                        json_string(config ? config->tunName : ""));
    json_object_set_new(obj, "tunAddress",
                        json_string(config ? config->tunAddr : ""));

    json_object_set_new(obj, "uplinkPackets", json_integer(uplinkPackets));
    json_object_set_new(obj, "uplinkBytes", json_integer(uplinkBytes));
    json_object_set_new(obj, "downlinkPackets", json_integer(downlinkPackets));
    json_object_set_new(obj, "downlinkBytes", json_integer(downlinkBytes));
    json_object_set_new(obj, "droppedPackets", json_integer(droppedPackets));
    json_object_set_new(obj, "droppedBytes", json_integer(droppedBytes));

    return obj;
}
