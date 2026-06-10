/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SERIAL_H_
#define SERIAL_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

typedef struct {
    int fd;
    char device[256];
    int baud;
} SerialPort;

bool serial_open(SerialPort *port, const char *device, int baud);
void serial_close(SerialPort *port);
bool serial_write_all(SerialPort *port, const uint8_t *data, size_t len);
bool serial_read_frame(SerialPort *port,
                       uint8_t *buf,
                       size_t size,
                       size_t *len,
                       int timeoutMs);

#endif /* SERIAL_H_ */
