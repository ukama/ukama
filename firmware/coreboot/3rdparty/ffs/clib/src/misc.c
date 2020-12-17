/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/src/misc.c $                                             */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */

/*
 *   File: misc.c
 * Author: Shaun Wetzstein <shaun@us.ibm.com>
 *  Descr:
 *   Note:
 *   Date: 07/26/09
 */

#include <stdarg.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <ctype.h>

#include "misc.h"

inline void prefetch(void *addr, size_t len, ...)
{
	va_list v;
	va_start(v, len);
	int w = va_arg(v, int);
	int l = va_arg(v, int);
	va_end(v);

	while (0 < len) {
		if (w) {
			switch (l) {
			case 0:
				__builtin_prefetch(addr, 1, 0);
				break;
			case 1:
				__builtin_prefetch(addr, 1, 1);
				break;
			case 2:
				__builtin_prefetch(addr, 1, 2);
				break;
			}
		} else {
			switch (l) {
			case 0:
				__builtin_prefetch(addr, 0, 0);
				break;
			case 1:
				__builtin_prefetch(addr, 0, 1);
				break;
			case 2:
				__builtin_prefetch(addr, 0, 2);
				break;
			}
		}
		len -= sizeof addr;
	}
}

unsigned long align(unsigned long size, unsigned long offset)
{
	--offset;
	return (size + offset) & ~offset;
}

void print_binary(FILE * __out, void *__data, size_t __len)
{
	static const char *__ascii[] = {
		"0000", "0001", "0010", "0011", "0100", "0101", "0110", "0111",
		"1000", "1001", "1010", "1011", "1100", "1101", "1110", "1111",
	};

	if (__data == NULL)
		return;
	if (__out == NULL)
		__out = stdout;

	size_t i;
	for (i = 0; i < __len; i++) {
		unsigned char c = *(unsigned char *)(__data + i);

		fprintf(__out, "%4.4s%4.4s",
			__ascii[(c & 0xF0) >> 4], __ascii[c & 0x0F]);
	}
}

#if 0
void dump_memory(FILE * file, unsigned long addr, const void *buf, size_t size)
{
	if (size <= 0 || buf == NULL)
		return;
	if (file == NULL)
		file = stdout;

	unsigned long *ul = (unsigned long *)buf;
	char hex[64] = { 0, }, ascii[32] = {
	0,}, c;
	int hl = 0, al = 0;

	int cnt = size / sizeof(unsigned long);

	int i;
	for (i = 0; i < cnt; i++) {
		hl += snprintf(hex + hl, sizeof hex, "%08lx", ul[i]);

		c = (ul[i] & 0xFF000000) >> 24;
		al +=
		    snprintf(ascii + al, sizeof ascii, "%c",
			     (isprint(c) ? c : '.'));

		c = (ul[i] & 0x00FF0000) >> 16;
		al +=
		    snprintf(ascii + al, sizeof ascii, "%c",
			     (isprint(c) ? c : '.'));

		c = (ul[i] & 0x0000FF00) >> 8;
		al +=
		    snprintf(ascii + al, sizeof ascii, "%c",
			     (isprint(c) ? c : '.'));

		c = (ul[i] & 0x000000FF) >> 0;
		al +=
		    snprintf(ascii + al, sizeof ascii, "%c",
			     (isprint(c) ? c : '.'));

		if ((i + 1) % 4 == 0) {
			fprintf(file, "  %08lx [%s] [%s]\n", addr, hex, ascii);

			addr += 4 * sizeof(unsigned long);
			hl = al = 0;

			memset(hex, 0, sizeof hex);
			memset(ascii, 0, sizeof ascii);
		} else {
			hl += snprintf(hex + hl, sizeof hex, " ");
		}
	}

	int j;
	for (j = 0; j < (i % 4); j++) {
		hl += snprintf(hex + hl, sizeof hex, "%08lx ", 0ul);
		al += snprintf(ascii + al, sizeof ascii, "%4.4s", "....");
	}

	if (i % 4 != 0) {
		fprintf(file, "  %08lx [%s] [%s]\n", addr, hex, ascii);
	}

	return;
}
#endif

/*
 * "aaaaaaaa [wwwwwwww xxxxxxxx yyyyyyyy zzzzzzzz] [........ ........]"
 */
void dump_memory(FILE * __out, uint32_t __addr, const void *__restrict __buf,
		 size_t __buf_sz)
{
	if (__buf_sz <= 0 || __buf == NULL)
		return;
	if (__out == NULL)
		__out = stdout;

	size_t hex_size = 16 * 2;	// 16 bytes per doubleword
	size_t ascii_size = 8 * 2;	// 8 bytes per doubleword

	uint8_t hex[hex_size + 1], ascii[hex_size + 1];
	memset(hex, '.', hex_size), memset(ascii, '.', ascii_size);

	void print_line(void) {
		fprintf(__out, "%08x ", __addr);

		fprintf(__out, "[%.8s %.8s %.8s %.8s] ",
			hex + 0, hex + 8, hex + 16, hex + 24);

		fprintf(__out, "[%.8s %.8s]\n", ascii + 0, ascii + 8);
	}

	size_t s = __addr % 16, i = 0;
	__addr &= ~0xF;

	for (i = s; i < __buf_sz + s; i++) {
		unsigned char c = ((unsigned char *)__buf)[i - s];

		hex[((i << 1) + 0) % hex_size] = "0123456789abcdef"[c >> 4];
		hex[((i << 1) + 1) % hex_size] = "0123456789abcdef"[c & 0xf];

		ascii[i % ascii_size] = isprint(c) ? c : '.';

		if (i == 0)
			continue;

		if ((i + 1) % ascii_size == 0) {
			print_line();
			memset(hex, '.', hex_size), memset(ascii, '.',
							   ascii_size);

			__addr += ascii_size;
		}
	}

	if (i % ascii_size)
		print_line();

	return;
}

int __round_pow2(int size)
{
	size--;

	size |= size >> 1;
	size |= size >> 2;
	size |= size >> 4;
	size |= size >> 8;
	size |= size >> 16;

	return ++size;
}
