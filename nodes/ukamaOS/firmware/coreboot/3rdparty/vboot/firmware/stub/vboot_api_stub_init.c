/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Stub implementations of firmware-provided API functions.
 */


#include <stdarg.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>

#include "2common.h"
#include "vboot_api.h"

uint64_t VbExGetTimer(void)
{
	struct timeval tv;
	gettimeofday(&tv, NULL);
	return (uint64_t)tv.tv_sec * VB_USEC_PER_SEC + (uint64_t)tv.tv_usec;
}

vb2_error_t VbExNvStorageRead(uint8_t *buf)
{
	return VB2_SUCCESS;
}

vb2_error_t VbExNvStorageWrite(const uint8_t *buf)
{
	return VB2_SUCCESS;
}
