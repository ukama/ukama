/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2012, Atmel Corporation
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
#include "hardware.h"
#include "arch/at91_wdt.h"
#include "debug.h"

static inline void wdt_write(unsigned int offset, unsigned int value)
{
	writel(value, (AT91C_BASE_WDT + offset));
}

static inline int wdt_read(unsigned int offset)
{
	return(readl(AT91C_BASE_WDT + offset));
}

#ifdef CONFIG_DISABLE_WATCHDOG
void at91_disable_wdt(void)
{
	unsigned int reg;

	reg = wdt_read(WDTC_MR);
	reg |= AT91C_WDTC_WDDIS;
	wdt_write(WDTC_MR, reg);
}
#else
void at91_disable_wdt(void) {}
#endif

#if defined(CONFIG_ENTER_NWD)
unsigned int at91_wdt_set_counter(unsigned int count)
{
	unsigned int reg;

	reg = wdt_read(WDTC_MR);
	if (reg & AT91C_WDTC_WDDIS) {
		dbg_info("at91_wdt: watchdog is disabled\n");
		return -1;
	}

	reg = AT91C_WDTC_WDRSTEN
		| AT91C_WDTC_WDDBGHLT
		| AT91C_WDTC_WDD
		| (count & AT91C_WDTC_WDV);
	wdt_write(WDTC_MR, reg);

	return 0;
}

unsigned int at91_wdt_reload_counter(void)
{
	wdt_write(WDTC_CR, (AT91C_WDTC_WDRSTT | AT91C_WDTC_KEY));

	return 0;
}
#endif
