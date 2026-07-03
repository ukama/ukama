/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <termios.h>
#include <unistd.h>
#include <sys/select.h>

#include "serial.h"

static speed_t baud_to_speed(int baud)
{
    switch (baud) {
    case 9600:   return B9600;
    case 19200:  return B19200;
    case 38400:  return B38400;
    case 57600:  return B57600;
    case 115200: return B115200;
    default:     return B9600;
    }
}

static bool serial_apply_config(SerialPort *port, int baud)
{
    struct termios tio;

    if (port == NULL || port->fd < 0) {
        return false;
    }

    memset(&tio, 0, sizeof(tio));
    if (tcgetattr(port->fd, &tio) != 0) {
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
    tio.c_lflag = 0;
    tio.c_oflag = 0;
    tio.c_cc[VMIN] = 0;
    tio.c_cc[VTIME] = 0;

    if (tcsetattr(port->fd, TCSANOW, &tio) != 0) {
        return false;
    }

    port->baud = baud;
    tcflush(port->fd, TCIOFLUSH);

    return true;
}

bool serial_open(SerialPort *port, const char *device, int baud)
{
    if (port == NULL || device == NULL) {
        return false;
    }

    memset(port, 0, sizeof(SerialPort));
    port->fd = -1;
    snprintf(port->device, sizeof(port->device), "%s", device);

    port->fd = open(device, O_RDWR | O_NOCTTY | O_SYNC);
    if (port->fd < 0) {
        return false;
    }

    if (!serial_apply_config(port, baud)) {
        serial_close(port);
        return false;
    }

    return true;
}

void serial_close(SerialPort *port)
{
    if (port == NULL) {
        return;
    }

    if (port->fd >= 0) {
        close(port->fd);
    }

    port->fd = -1;
}

bool serial_set_baud(SerialPort *port, int baud)
{
    if (port == NULL || port->fd < 0) {
        return false;
    }

    if (port->baud == baud) {
        return true;
    }

    return serial_apply_config(port, baud);
}

bool serial_flush_rx(SerialPort *port)
{
    if (port == NULL || port->fd < 0) {
        return false;
    }

    return tcflush(port->fd, TCIFLUSH) == 0;
}

bool serial_write_all(SerialPort *port, const uint8_t *data, size_t len)
{
    size_t off;
    ssize_t n;

    if (port == NULL || port->fd < 0 || data == NULL) {
        return false;
    }

    off = 0;
    while (off < len) {
        n = write(port->fd, data + off, len - off);
        if (n < 0 && errno == EINTR) {
            continue;
        }
        if (n <= 0) {
            return false;
        }
        off += (size_t)n;
    }

    return tcdrain(port->fd) == 0;
}

bool serial_read_frame(SerialPort *port,
                       uint8_t *buf,
                       size_t size,
                       size_t *len,
                       int timeoutMs)
{
    fd_set rfds;
    struct timeval tv;
    uint8_t byte;
    bool seenFlag;
    size_t off;
    int rc;
    ssize_t n;

    if (port == NULL || port->fd < 0 || buf == NULL || len == NULL || size == 0) {
        return false;
    }

    off = 0;
    seenFlag = false;

    while (off < size) {
        FD_ZERO(&rfds);
        FD_SET(port->fd, &rfds);
        tv.tv_sec = timeoutMs / 1000;
        tv.tv_usec = (timeoutMs % 1000) * 1000;

        rc = select(port->fd + 1, &rfds, NULL, NULL, &tv);
        if (rc < 0 && errno == EINTR) {
            continue;
        }
        if (rc <= 0) {
            return false;
        }

        n = read(port->fd, &byte, 1);
        if (n < 0 && errno == EINTR) {
            continue;
        }
        if (n != 1) {
            return false;
        }

        if (byte == 0x7E) {
            if (!seenFlag) {
                seenFlag = true;
                off = 0;
                buf[off++] = byte;
                continue;
            }

            buf[off++] = byte;
            *len = off;
            return true;
        }

        if (!seenFlag) {
            continue;
        }

        buf[off++] = byte;
    }

    return false;
}
