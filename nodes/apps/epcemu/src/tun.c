/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <linux/if.h>
#include <linux/if_tun.h>
#include <stdio.h>
#include <string.h>
#include <sys/ioctl.h>
#include <sys/wait.h>
#include <unistd.h>

#include "epcemu.h"
#include "tun.h"

static int exec_cmd(const char *cmd,
                    const char *a1,
                    const char *a2,
                    const char *a3,
                    const char *a4,
                    const char *a5,
                    const char *a6) {

    pid_t pid;
    int status;

    pid = fork();
    if (pid < 0) {
        usys_log_error("fork failed for %s: %s", cmd, strerror(errno));
        return USYS_FALSE;
    }

    if (pid == 0) {
        execlp(cmd, cmd, a1, a2, a3, a4, a5, a6, (char *)NULL);
        _exit(127);
    }

    while (waitpid(pid, &status, 0) < 0) {
        if (errno == EINTR) continue;
        usys_log_error("waitpid failed for %s: %s", cmd, strerror(errno));
        return USYS_FALSE;
    }

    if (!WIFEXITED(status) || WEXITSTATUS(status) != 0) {
        usys_log_error("command failed cmd=%s rc=%d",
                       cmd,
                       WIFEXITED(status) ? WEXITSTATUS(status) : -1);
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

int tun_create(const char *name) {

    struct ifreq ifr;
    int fd;

    if (name == NULL || name[0] == '\0') return -1;

    fd = open("/dev/net/tun", O_RDWR);
    if (fd < 0) {
        usys_log_error("failed to open /dev/net/tun: %s", strerror(errno));
        return -1;
    }

    memset(&ifr, 0, sizeof(ifr));
    ifr.ifr_flags = IFF_TUN | IFF_NO_PI;
    snprintf(ifr.ifr_name, IFNAMSIZ, "%s", name);

    if (ioctl(fd, TUNSETIFF, (void *)&ifr) < 0) {
        usys_log_error("failed to create tun %s: %s", name, strerror(errno));
        close(fd);
        return -1;
    }

    return fd;
}

int tun_configure(const char *name, const char *addr) {

    if (name == NULL || addr == NULL) return USYS_FALSE;

    exec_cmd("ip", "link", "delete", name, NULL, NULL, NULL, NULL);

    if (!exec_cmd("ip", "tuntap", "add", "dev", name, "mode", "tun")) {
        return USYS_FALSE;
    }

    if (!exec_cmd("ip", "addr", "replace", addr, "dev", name, NULL, NULL)) {
        return USYS_FALSE;
    }

    if (!exec_cmd("ip", "link", "set", name, "up", NULL, NULL, NULL)) {
        return USYS_FALSE;
    }

    return USYS_TRUE;
}

void tun_close(int fd) {

    if (fd >= 0) {
        close(fd);
    }
}
