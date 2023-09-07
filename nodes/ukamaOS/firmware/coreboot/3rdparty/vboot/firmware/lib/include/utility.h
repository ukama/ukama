/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Helper functions/wrappers for memory allocations, manipulation and
 * comparison.
 */

#ifndef VBOOT_REFERENCE_UTILITY_H_
#define VBOOT_REFERENCE_UTILITY_H_

#include "2common.h"
#include "2sysincludes.h"
#include "vboot_api.h"

/*
 * Buffer size required to hold the longest possible output of Uint64ToString()
 * - that is, Uint64ToString(~0, 2).
 */
#define UINT64_TO_STRING_MAX 65

/**
 * Convert a value to a string in the specified radix (2=binary, 10=decimal,
 * 16=hex) and store it in <buf>, which is <bufsize> chars long.  If
 * <zero_pad_width>, left-pads the string to at least that width with '0'.
 * Returns the length of the stored string, not counting the terminating null.
 */
uint32_t Uint64ToString(char *buf, uint32_t bufsize, uint64_t value,
			uint32_t radix, uint32_t zero_pad_width);

/**
 * Concatenate <src> onto <dest>, which has space for <destlen> characters
 * including the terminating null.  Note that <dest> will always be
 * null-terminated if <destlen> > 0.  Returns the number of characters used in
 * <dest>, not counting the terminating null.
 */
uint32_t StrnAppend(char *dest, const char *src, uint32_t destlen);

#endif  /* VBOOT_REFERENCE_UTILITY_H_ */
