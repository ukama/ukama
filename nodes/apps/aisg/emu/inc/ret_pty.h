/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef AISG_EMU_RET_PTY_H_
#define AISG_EMU_RET_PTY_H_

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <signal.h>

bool ret_pty_open(const char *linkPath,
                  int *masterFd,
                  char *slaveName,
                  size_t slaveNameSize);
bool ret_pty_read_hdlc_frame(int fd,
                             volatile sig_atomic_t *running,
                             uint8_t *buf,
                             size_t size,
                             size_t *len);
bool ret_pty_write_all(int fd, const uint8_t *buf, size_t len);
void ret_pty_close(int fd, const char *linkPath);

#endif /* AISG_EMU_RET_PTY_H_ */
