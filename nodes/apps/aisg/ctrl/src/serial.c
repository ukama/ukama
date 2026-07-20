/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <fcntl.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <termios.h>
#include <time.h>
#include <unistd.h>
#include <sys/select.h>

#include "serial.h"

#define HDLC_FLAG                          0x7E

static int64_t monotonic_ms(void) {
    struct timespec ts;

    if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) return 0;

    return ((int64_t)ts.tv_sec * 1000) + ((int64_t)ts.tv_nsec / 1000000);
}

static speed_t baud_to_speed(int baud) {
    switch (baud) {
    case 9600:   return B9600;
    case 19200:  return B19200;
    case 38400:  return B38400;
    case 57600:  return B57600;
    case 115200: return B115200;
    default:     return B9600;
    }
}

bool serial_open(SerialPort *port, const char *device, int baud) {
    struct termios tio;

    if (port == NULL || device == NULL) return false;

    memset(port, 0, sizeof(SerialPort));
    port->fd = -1;
    snprintf(port->device, sizeof(port->device), "%s", device);
    port->baud = baud;

    port->fd = open(device, O_RDWR | O_NOCTTY | O_SYNC);
    if (port->fd < 0) return false;

    memset(&tio, 0, sizeof(tio));
    if (tcgetattr(port->fd, &tio) != 0) {
        serial_close(port);
        return false;
    }

    cfmakeraw(&tio);
    cfsetispeed(&tio, baud_to_speed(baud));
    cfsetospeed(&tio, baud_to_speed(baud));

    tio.c_cflag |= CLOCAL | CREAD;
    tio.c_cflag &= ~PARENB;
    tio.c_cflag &= ~CSTOPB;
    tio.c_cflag &= ~CSIZE;
    tio.c_cflag |= CS8;
#ifdef CRTSCTS
    tio.c_cflag &= ~CRTSCTS;
#endif
    tio.c_iflag &= ~(IXON | IXOFF | IXANY);
    tio.c_cc[VMIN] = 0;
    tio.c_cc[VTIME] = 0;

    if (tcsetattr(port->fd, TCSANOW, &tio) != 0) {
        serial_close(port);
        return false;
    }

    tcflush(port->fd, TCIOFLUSH);
    return true;
}

void serial_close(SerialPort *port) {
    if (port == NULL) return;
    if (port->fd >= 0) close(port->fd);
    port->fd = -1;
}

bool serial_write_all(SerialPort *port, const uint8_t *data, size_t len) {
    size_t off;
    ssize_t n;

    if (port == NULL || port->fd < 0 || data == NULL) return false;

    off = 0;
    while (off < len) {
        n = write(port->fd, data + off, len - off);
        if (n < 0 && errno == EINTR) continue;
        if (n <= 0) return false;
        off += (size_t)n;
    }

    /* Start the AISG reply timer only after the final flag left the tty. */
    while (tcdrain(port->fd) != 0) {
        if (errno != EINTR) return false;
    }

    return true;
}

bool serial_discard_input(SerialPort *port) {
    if (port == NULL || port->fd < 0) return false;

    port->openingFlagPending = false;
    while (tcflush(port->fd, TCIFLUSH) != 0) {
        if (errno != EINTR) return false;
    }

    return true;
}

bool serial_read_frame(SerialPort *port,
                       uint8_t *buf,
                       size_t size,
                       size_t *len,
                       int timeoutMs) {
    fd_set rfds;
    struct timeval tv;
    uint8_t byte;
    bool seenFlag;
    size_t off;
    int64_t deadlineMs;
    int64_t nowMs;
    int waitMs;
    int ready;
    ssize_t n;

    if (port == NULL || port->fd < 0 || buf == NULL || len == NULL ||
        size < 2) {
        return false;
    }

    *len = 0;
    off = 0;
    seenFlag = port->openingFlagPending;
    port->openingFlagPending = false;

    if (seenFlag) {
        buf[off++] = HDLC_FLAG;
    }

    if (timeoutMs <= 0) timeoutMs = 1;
    deadlineMs = monotonic_ms() + timeoutMs;

    for (;;) {
        nowMs = monotonic_ms();
        waitMs = (int)(deadlineMs - nowMs);
        if (waitMs <= 0) {
            *len = off;
            port->openingFlagPending = seenFlag && off == 1;
            return false;
        }

        FD_ZERO(&rfds);
        FD_SET(port->fd, &rfds);
        tv.tv_sec = waitMs / 1000;
        tv.tv_usec = (waitMs % 1000) * 1000;

        ready = select(port->fd + 1, &rfds, NULL, NULL, &tv);
        if (ready < 0 && errno == EINTR) continue;
        if (ready <= 0) {
            *len = off;
            port->openingFlagPending = seenFlag && off == 1;
            return false;
        }

        n = read(port->fd, &byte, 1);
        if (n < 0 && errno == EINTR) continue;
        if (n != 1) {
            *len = off;
            port->openingFlagPending = seenFlag && off == 1;
            return false;
        }

        if (byte == HDLC_FLAG) {
            if (!seenFlag) {
                seenFlag = true;
                off = 0;
                buf[off++] = byte;
                continue;
            }

            /* Fill flags and a shared closing/opening flag are not frames. */
            if (off == 1) continue;

            if (off >= size) {
                *len = off;
                return false;
            }

            buf[off++] = byte;
            if (off > 2) {
                *len = off;
                port->openingFlagPending = true;
                return true;
            }

            continue;
        }

        if (!seenFlag) continue;

        if (off >= size) {
            *len = off;
            return false;
        }

        buf[off++] = byte;
    }
}
