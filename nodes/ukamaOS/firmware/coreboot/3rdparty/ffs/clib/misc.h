/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: clib/misc.h $                                                 */
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
 *   File: misc.h
 * Author: Shaun Wetzstein <shaun@us.ibm.com?
 *  Descr: Misc. library helper functions
 *   Note:
 *   Date: 10/22/10
 */

#ifndef __MISC_H__
#define __MISC_H__

#include <stdint.h>
#include <stdio.h>

#include <sys/syscall.h>

#include "attribute.h"

#include "align.h"
#include "min.h"
#include "max.h"

#define INT32_BIT	(CHAR_BIT*sizeof(int32_t))

#define ARRAY_SIZE(x)	(sizeof(x) / sizeof((x)[0]))

#ifndef be16toh
#include <byteswap.h>
#if __BYTE_ORDER == __LITTLE_ENDIAN
#define be16toh(x) __bswap_16(x)
#define htobe16(x) __bswap_16(x)
#define le16toh(x) (x)
#define htole16(x) (x)
#else
#define be16toh(x) (x)
#define htobe16(x) (x)
#define le16toh(x) __bswap_16(x)
#define htole16(x) __bswap_16(x)
#endif
#endif

#ifndef be32toh
#include <byteswap.h>
#if __BYTE_ORDER == __LITTLE_ENDIAN
#define be32toh(x) __bswap_32(x)
#define htobe32(x) __bswap_32(x)
#define le32toh(x) (x)
#define htole32(x) (x)
#else
#define be32toh(x) (x)
#define htobe32(x) (x)
#define le32toh(x) __bswap_32(x)
#define htole32(x) __bswap_32(x)
#endif
#endif

#ifndef be64toh
#include <byteswap.h>
#if __BYTE_ORDER == __LITTLE_ENDIAN
#define be64toh(x) __bswap_64(x)
#define htobe64(x) __bswap_64(x)
#define le64toh(x) (x)
#define htole64(x) (x)
#else
#define be64toh(x) (x)
#define htobe64(x) (x)
#define le64toh(x) __bswap_64(x)
#define htole64(x) __bswap_64(x)
#endif
#endif

#ifndef gettid
#define gettid()	({			\
    pid_t __tid = (pid_t)syscall(SYS_gettid);	\
    __tid;					\
			 })
#endif

#define is_pow2(x)				\
({						\
    typeof(x) _x = (x);				\
    (((_x) != 0) && !((_x) & ((_x) - 1)));	\
})

#define __align(x,y)				\
({						\
    assert(is_pow2((y)));			\
    typeof((y)) _y = (y) - 1;			\
    typeof((x)) _x = ((x)+_y) & ~_y;		\
    _x;						\
})

extern void prefetch(void *addr, size_t len, ...) __nonnull((1));
extern void print_binary(FILE *, void *, size_t);
extern void dump_memory(FILE *, uint32_t, const void *__restrict,
			size_t) __nonnull((1, 3));
extern unsigned long align(unsigned long size, unsigned long offset);
extern int __round_pow2(int size);

#endif				/* __MISC_H__ */
