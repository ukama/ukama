/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

#include <stdbool.h>
#include <stdint.h>
#include <sys/types.h>

bool log_router_start(const char *socketPath,
                      int maxRecordBytes,
                      int reconnectMs);
void log_router_stop(void);
bool log_router_register(const char *space,
                         const char *app,
                         pid_t pid,
                         uint32_t generation,
                         int stdoutFd,
                         int stderrFd);
