/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#define _XOPEN_SOURCE 600
#define _DEFAULT_SOURCE

#include <errno.h>
#include <fcntl.h>
#include <poll.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>
#include <termios.h>
#include <unistd.h>
#include <stdlib.h>

#include "hdlc.h"
#include "ret_pty.h"
#include "usys_log.h"

#define RET_PTY_POLL_MS 250

static void set_raw_termios(int fd)
{
    struct termios tio;

    if (tcgetattr(fd, &tio) != 0) {
        return;
    }

    cfmakeraw(&tio);
    tio.c_cflag |= CLOCAL | CREAD;
    tio.c_cc[VMIN] = 1;
    tio.c_cc[VTIME] = 0;

    tcsetattr(fd, TCSANOW, &tio);
}

bool ret_pty_open(const char *linkPath,
                  int *masterFd,
                  char *slaveName,
                  size_t slaveNameSize)
{
    int master;
    int slave;
    char *name;

    if (linkPath == NULL || masterFd == NULL ||
        slaveName == NULL || slaveNameSize == 0) {
        return false;
    }

    master = posix_openpt(O_RDWR | O_NOCTTY);
    if (master < 0) {
        usys_log_error("ret-emu: posix_openpt failed: %s", strerror(errno));
        return false;
    }

    if (grantpt(master) != 0 || unlockpt(master) != 0) {
        usys_log_error("ret-emu: grantpt/unlockpt failed: %s", strerror(errno));
        close(master);
        return false;
    }

    name = ptsname(master);
    if (name == NULL) {
        usys_log_error("ret-emu: ptsname failed: %s", strerror(errno));
        close(master);
        return false;
    }

    snprintf(slaveName, slaveNameSize, "%s", name);

    set_raw_termios(master);

    slave = open(slaveName, O_RDWR | O_NOCTTY);
    if (slave >= 0) {
        set_raw_termios(slave);
        close(slave);
    }

    unlink(linkPath);
    if (symlink(slaveName, linkPath) != 0) {
        usys_log_error("ret-emu: failed to symlink %s -> %s: %s",
                       linkPath,
                       slaveName,
                       strerror(errno));
        close(master);
        return false;
    }

    *masterFd = master;

    usys_log_info("ret-emu: PTY ready %s -> %s", linkPath, slaveName);

    return true;
}

bool ret_pty_read_hdlc_frame(int fd,
                             volatile sig_atomic_t *running,
                             uint8_t *buf,
                             size_t size,
                             size_t *len)
{
    struct pollfd pfd;
    size_t off = 0;
    bool started = false;
    uint8_t byte;
    ssize_t n;
    int rc;

    if (fd < 0 || running == NULL || buf == NULL || len == NULL || size < 2) {
        return false;
    }

    *len = 0;

    while (*running) {
        memset(&pfd, 0, sizeof(pfd));
        pfd.fd = fd;
        pfd.events = POLLIN;

        rc = poll(&pfd, 1, RET_PTY_POLL_MS);
        if (rc == 0) {
            continue;
        }

        if (rc < 0) {
            if (errno == EINTR) {
                continue;
            }
            usys_log_error("ret-emu: PTY poll failed: %s", strerror(errno));
            return false;
        }

        n = read(fd, &byte, 1);
        if (n <= 0) {
            if (n < 0 && errno == EINTR) {
                continue;
            }
            return false;
        }

        if (byte == HDLC_FLAG) {
            if (!started) {
                started = true;
                off = 0;
                buf[off++] = byte;
                continue;
            }

            if (off >= size) {
                return false;
            }
            buf[off++] = byte;
            *len = off;
            return true;
        }

        if (!started) {
            continue;
        }

        if (off >= size) {
            usys_log_warn("ret-emu: HDLC frame too large");
            started = false;
            off = 0;
            continue;
        }

        buf[off++] = byte;
    }

    return false;
}

bool ret_pty_write_all(int fd, const uint8_t *buf, size_t len)
{
    size_t off = 0;
    ssize_t n;

    if (fd < 0 || (buf == NULL && len != 0)) {
        return false;
    }

    while (off < len) {
        n = write(fd, buf + off, len - off);
        if (n <= 0) {
            if (errno == EINTR) {
                continue;
            }
            usys_log_error("ret-emu: PTY write failed: %s", strerror(errno));
            return false;
        }
        off += (size_t)n;
    }

    return true;
}

void ret_pty_close(int fd, const char *linkPath)
{
    if (fd >= 0) {
        close(fd);
    }

    if (linkPath != NULL) {
        unlink(linkPath);
    }
}
