/* Copyright 2016 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef __CROS_EC_ASSERT_H__
#define __CROS_EC_ASSERT_H__

/* Include CONFIG definitions for EC sources. */
#ifndef THIRD_PARTY
#include "common.h"
#endif

#ifdef __cplusplus
extern "C" {
#endif

#ifdef CONFIG_DEBUG_ASSERT
#ifdef CONFIG_DEBUG_ASSERT_REBOOTS

#ifdef CONFIG_DEBUG_ASSERT_BRIEF
extern void panic_assert_fail(const char *fname, int linenum)
	__attribute__((noreturn));
#define ASSERT(cond) do {					\
		if (!(cond))					\
			panic_assert_fail(__FILE__, __LINE__);	\
	} while (0)
#else
extern void panic_assert_fail(const char *msg, const char *func,
		const char *fname, int linenum) __attribute__((noreturn));
#define ASSERT(cond) do {					     \
		if (!(cond))					     \
			panic_assert_fail(#cond, __func__, __FILE__, \
					__LINE__);		     \
	} while (0)
#endif
#else
#define ASSERT(cond) do {			\
		if (!(cond))			\
			__asm("bkpt");		\
			__builtin_unreachable();\
	} while (0)
#endif
#else
#define ASSERT(cond)
#endif

/* This collides with cstdlib, so exclude it where cstdlib is supported. */
#ifndef assert
#define assert(x...) ASSERT(x)
#endif

#ifdef __cplusplus
}
#endif

#endif /* __CROS_EC_ASSERT_H__ */
