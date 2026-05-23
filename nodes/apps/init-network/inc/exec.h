/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef EXEC_H_
#define EXEC_H_

#include <stdbool.h>

#define EXEC_MAX_ARGS 64
#define EXEC_MAX_PATH 512

bool exec_tool_exists(const char *cmd);
int exec_cmd(int timeoutSec, const char *cmd, ...);
int exec_cmd_argv(int timeoutSec, char *const argv[]);

#endif /* EXEC_H_ */
