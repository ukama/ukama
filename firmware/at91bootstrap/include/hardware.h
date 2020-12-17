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
#ifndef __HARDWARE_H__
#define __HARDWARE_H__

#ifdef CONFIG_AT91SAM9260
#include "arch/at91sam9260.h"
#endif

#ifdef CONFIG_AT91SAM9XE
#include "arch/at91sam9xe.h"
#endif

#ifdef CONFIG_AT91SAM9G10
#include "arch/at91sam9g10.h"
#endif

#ifdef CONFIG_AT91SAM9G20
#include "arch/at91sam9g20.h"
#endif

#ifdef CONFIG_AT91SAM9G45
#include "arch/at91sam9g45.h"
#endif

#ifdef CONFIG_AT91SAM9261
#include "arch/at91sam9261.h"
#endif

#ifdef CONFIG_AT91SAM9263
#include "arch/at91sam9263.h"
#endif

#ifdef CONFIG_AT91SAM9RL
#include "arch/at91sam9rl.h"
#endif

#ifdef CONFIG_AT91SAM9X5
#include "arch/at91sam9x5.h"
#endif

#ifdef CONFIG_AT91SAM9N12
#include "arch/at91sam9n12.h"
#endif

#ifdef CONFIG_SAM9X60
#include "arch/sam9x60.h"
#endif

#ifdef CONFIG_SAMA5D3X
#include "arch/sama5d3x.h"
#endif

#ifdef CONFIG_SAMA5D4
#include "arch/sama5d4.h"
#endif

#ifdef CONFIG_SAMA5D2
#include "arch/sama5d2.h"
#endif

#ifdef CONFIG_SAMA7G5
#include "arch/sama7g5.h"
#endif

/* I/O Function Macro */
#define writel(value, addr) \
	(*(volatile unsigned int *)(addr)) = (value)
#define readl(addr) \
	(*(volatile unsigned int *)(addr))

#define writew(value, addr) \
	 (*(volatile unsigned short *)(addr) = (value))
#define readw(addr) \
	(*(volatile unsigned short *)(addr))

#define writeb(value, addr) \
	 (*(volatile unsigned char *)(addr) = (value))
#define readb(addr) \
	(*(volatile unsigned char *)(addr))

#endif /* #ifndef __HARDWARE_H__ */
