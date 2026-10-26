/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <linux/if_tun.h>
#include <net/if.h>
#include <stdio.h>
#include <string.h>
#include <sys/ioctl.h>
#include <unistd.h>

#include "epcemu.h"
#include "tun.h"

static int interface_flags(const char *name, struct ifreq *ifr) {

    int fd;
    int ret;
    int error;

    if (name == NULL || name[0] == '\0' || strlen(name) >= IFNAMSIZ) {
        errno = EINVAL;
        return -1;
    }

    memset(ifr, 0, sizeof(*ifr));
    snprintf(ifr->ifr_name, IFNAMSIZ, "%s", name);
    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd < 0) return -1;
    ret = ioctl(fd, SIOCGIFFLAGS, ifr);
    error = errno;
    close(fd);
    errno = error;
    return ret;
}

int tun_create(const char *name) {

    struct ifreq ifr;
    unsigned int index;
    int fd;
    int error;

    /* TUNSETIFF attaches to a prepared device; it must not create one. */
    if (interface_flags(name, &ifr) < 0) return -1;
    if (!(ifr.ifr_flags & IFF_UP)) {
        errno = ENETDOWN;
        return -1;
    }
    index = if_nametoindex(name);
    if (index == 0) return -1;

    fd = open("/dev/net/tun", O_RDWR | O_CLOEXEC);
    if (fd < 0) return -1;

    ifr.ifr_flags = IFF_TUN | IFF_NO_PI;
    if (ioctl(fd, TUNSETIFF, &ifr) < 0 ||
        ioctl(fd, TUNGETIFF, &ifr) < 0) {
        error = errno;
        close(fd);
        errno = error;
        return -1;
    }

    if (!(ifr.ifr_flags & IFF_PERSIST) ||
        if_nametoindex(name) != index) {
        close(fd);
        errno = ENODEV;
        return -1;
    }

    return fd;
}

void tun_close(int fd) {

    if (fd >= 0) {
        close(fd);
    }
}
