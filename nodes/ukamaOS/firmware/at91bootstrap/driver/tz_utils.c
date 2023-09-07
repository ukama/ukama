/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation
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
#include "mon_macros.h"
#include "mon.h"
#include "debug.h"

void dacr_swd_init(void)
{
	asm volatile (
		"mcr	p15, 0, %0, c3, c0, 0\n\t"
		"isb"
		:
		: "r" (0x55555555));
}

/* Enable access to CP 10 & 11 */
void cpacr_init(void)
{
	asm volatile (
		"mrc	p15, 0, r0, c1, c0, 2\n\t"
		"orr	r0, r0, %0\n\t"
		"mcr	p15, 0, r0, c1, c0, 2\n\t"
		"isb"
		:
		: "r" (0xf << 20)
		: "r0");
}

void enter_normal_world(void)
{
	asm volatile ("smc #0");
}

void switch_normal_world(void)
{
	monitor_init();

	dbg_info("Enter Normal World, Run at %x\n", NWD_BOOT_ADDR);

	enter_normal_world();
}
