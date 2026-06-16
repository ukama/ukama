/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_LOG_H_
#define ULAB_LOG_H_

#include <stdarg.h>

void ulab_log_set_verbose(int verbose);
void ulab_log_set_quiet(int quiet);
void ulab_log_debug(const char *fmt, ...);
void ulab_log_info(const char *fmt, ...);
void ulab_log_warn(const char *fmt, ...);
void ulab_log_error(const char *fmt, ...);
void ulab_status(const char *state, const char *fmt, ...);

#endif /* ULAB_LOG_H_ */
