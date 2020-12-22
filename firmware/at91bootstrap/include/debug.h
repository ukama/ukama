/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2006, Atmel Corporation
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Atmel's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY ATMEL "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL ATMEL BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#ifndef __DEBUG_H__
#define __DEBUG_H__

#define DEBUG_INFO        1
#define DEBUG_LOUD        2
#define DEBUG_VERY_LOUD   4

#define DUMP_WIDTH_BIT_8	1
#define DUMP_WIDTH_BIT_32	2

#ifdef CONFIG_DEBUG
extern int dbg_printf(const char *fmt_str, ...);
extern void dbg_hexdump(const unsigned char *buf,
			unsigned int size, unsigned int width);
#else
#define BOOTSTRAP_DEBUG_LEVEL 0

#include "usart.h"

static inline int dbg_printf(const char *fmt_str, ...)
{
	usart_puts(fmt_str);
	return 0;
}
static inline void dbg_hexdump(const unsigned char *buf,
			       unsigned int size, unsigned int width) { }
#endif

#define console_printf(fmt_str, args...) \
	dbg_printf(fmt_str , ## args)

#define dbg_log(level, fmt_str, args...) \
	({ \
		(level) <= BOOTSTRAP_DEBUG_LEVEL ? dbg_printf((fmt_str), ##args) : 0; \
	})

#define dbg_info(fmt_str, arg...)		\
	dbg_log(DEBUG_INFO, fmt_str , ## arg)
#define dbg_loud(fmt_str, arg...)		\
	dbg_log(DEBUG_LOUD, fmt_str , ## arg)
#define dbg_very_loud(fmt_str, arg...)		\
	dbg_log(DEBUG_VERY_LOUD, fmt_str , ## arg)

#endif /* #ifndef __DEBUG_H__ */
