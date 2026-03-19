/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef SWITCHD_UTILS_H
#define SWITCHD_UTILS_H

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

char *trim_ws(char *s);
int safe_snprintf(char *dst, size_t dstLen, const char *fmt, ...);
double parse_prefixed_double(const char *s);
uint64_t monotonic_msec(void);
int mkdir_p(const char *path, int mode);
int copy_file(const char *src, const char *dst);
const char *state_to_str(int state);
const char *op_type_to_str(int type);
const char *op_state_to_str(int state);
const char *fw_state_to_str(int state);
const char *alarm_severity_to_str(int severity);
const char *switch_error_to_str(int code);

#endif
