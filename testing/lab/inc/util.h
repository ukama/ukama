/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef ULAB_UTIL_H_
#define ULAB_UTIL_H_

#include "ulab.h"
#include <stdint.h>
#include <stddef.h>

char *ulab_trim(char *s);
int ulab_streq(const char *a, const char *b);
int ulab_starts(const char *s, const char *prefix);
int ulab_ends(const char *s, const char *suffix);
int ulab_parse_u32(const char *s, uint32_t *out);
int ulab_parse_u64(const char *s, uint64_t *out);
int ulab_parse_double(const char *s, double *out);
int ulab_copy(char *dst, size_t n, const char *src);
int ulab_mkdir_p(const char *path);
uint32_t ulab_hash32(const char *s, uint32_t seed);
int ulab_within_pct(uint64_t expected, uint64_t actual, uint32_t pct);
int ulab_run_cmd(const char *cmd, char *out, size_t out_len);
void ulab_json_escape(const char *in, char *out, size_t out_len);
const char *ulab_getenv_default(const char *name, const char *def);

#endif /* ULAB_UTIL_H_ */
