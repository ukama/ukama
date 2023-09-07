/* ----------------------------------------------------------------------------
 *         Microchip Microprocessor (MPU) Software Team
 * ----------------------------------------------------------------------------
 * Copyright (C) 2019 Microchip Technology Inc. and its subsidiaries
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL MICROCHIP BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#include "hardware.h"
#include "arch/at91_wdt2.h"
#include "debug.h"

static inline void wdt_write(unsigned int offset, unsigned int value)
{
	writel(value, (AT91C_BASE_WDT + offset));
}

static inline int wdt_read(unsigned int offset)
{
	return(readl(AT91C_BASE_WDT + offset));
}

#ifdef CONFIG_WDTS
static inline void wdts_write(unsigned int offset, unsigned int value)
{
	writel(value, (AT91C_BASE_WDTS + offset));
}

static inline int wdts_read(unsigned int offset)
{
	return(readl(AT91C_BASE_WDTS + offset));
}
#endif

#ifdef CONFIG_DISABLE_WATCHDOG
void at91_disable_wdt(void)
{
	unsigned int reg;

	reg = wdt_read(WDTC_MR);
	reg |= AT91C_WDTC_WDDIS;
	wdt_write(WDTC_MR, reg);
}

#ifdef CONFIG_WDTS
void at91_disable_wdts(void)
{
	unsigned int reg;

	reg = wdts_read(WDTC_MR);
	reg |= AT91C_WDTC_WDDIS;
	wdts_write(WDTC_MR, reg);
}
#endif
#else
void at91_disable_wdt(void) {}
void at91_disable_wdts(void) {}
#endif
