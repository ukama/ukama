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

#ifdef CONFIG_RSTC

#include "common.h"
#include "hardware.h"
#include "rstc.h"
#include "arch/at91_rstc.h"

static inline void rstc_write(unsigned int offset, unsigned int value)
{
	writel(value, AT91C_BASE_RSTC + offset);
}

static inline unsigned int rstc_read(unsigned int offset)
{
	return readl(AT91C_BASE_RSTC + offset);
}

/*
 * The external reset is asserted during a time of 2^(ERSTL+1) Slow Clock cycles
 * the external reset length is 15ms (1/32768 * (2 ^ 9).
 */
#define EXTERNAL_RESET_LENGTH	8

void rstc_external_reset(void)
{
	unsigned int timeout = 100000;

	/* Set the external reset length */
	rstc_write(RSTC_RMR, AT91C_RSTC_ERSTL(EXTERNAL_RESET_LENGTH)
					| AT91C_RSTC_KEY_UNLOCK);

	/* Assert the NRST pin */
	rstc_write(RSTC_RCR, (AT91C_RSTC_EXTRST | AT91C_RSTC_RCRKEY));

	/* Wait for the software reset completed */
	while(rstc_read(RSTC_RSR) & AT91C_RSTC_SRCMP);

	while (--timeout)
		;
}

void cpu_reset()
{
	rstc_write(RSTC_RCR, AT91C_RSTC_RCRKEY
		| AT91C_RSTC_PROCRST	/* Processor Reset */
		| AT91C_RSTC_PERRST	/* Peripheral Reset */
		| AT91C_RSTC_EXTRST);	/* External Reset (assert nRST pin) */
}

void rstc_ddr_rst_deassert(void)
{
	unsigned int grstr = rstc_read(RSTC_GRSTR);

	grstr |= AT91C_GRSTR_DDR_RST;
	rstc_write(RSTC_GRSTR, grstr);
}

void rstc_ddr_phy_rst_deassert(void)
{
	unsigned int grstr = rstc_read(RSTC_GRSTR);

	grstr |= AT91C_GRSTR_DDR_PHY_RST;
	rstc_write(RSTC_GRSTR, grstr);
}

#else
void rstc_external_reset(void)
{
}

void cpu_reset(void)
{
}
#endif
