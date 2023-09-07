/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2010, Intel Corp. Huang Ying <ying.huang@intel.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef _UUID_H_
#define _UUID_H_

#include <string.h>

#define UUID_LEN	16
#define UUID_STRLEN	36

/*
 * Parses a canonical UUID string into the common byte representation
 * where the first three words are interpreted as little endian:
 *
 * The UUID
 *   "00112233-4455-6677-8899-aabbccddeeff"
 * is stored as
 *   33 22 11 00 55 44 77 66 88 99 aa bb cc dd ee ff
 *
 * Returns negative value on error, 0 on success.
 */
int parse_uuid(uint8_t *uuid, const char *uuid_str);

typedef struct {
	uint8_t b[16];
} __packed guid_t;

#define GUID_INIT(a, b, c, d0, d1, d2, d3, d4, d5, d6, d7)		\
((guid_t)								\
{{ (a) & 0xff, ((a) >> 8) & 0xff, ((a) >> 16) & 0xff, ((a) >> 24) & 0xff, \
   (b) & 0xff, ((b) >> 8) & 0xff,					\
   (c) & 0xff, ((c) >> 8) & 0xff,					\
   (d0), (d1), (d2), (d3), (d4), (d5), (d6), (d7) } })

static inline int guidcmp(const guid_t *guid1, const guid_t *guid2)
{
	return memcmp(guid1, guid2, sizeof(guid_t));
}

static inline guid_t *guidcpy(guid_t *dest, const guid_t *src)
{
	return (guid_t *)memcpy(dest, src, sizeof(guid_t));
}

#endif /* _UUID_H_ */
