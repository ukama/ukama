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
#include "board.h"
#include "clk-common.h"
#include "pmc.h"
#include "arch/at91_pmc/pmc.h"

int pmc_enable_periph_clock(unsigned int periph_id, int divider)
{
	unsigned int val;

	if (divider == PMC_PERIPH_CLK_DIVIDER_NA)
		divider = 0;

	write_pmc(PMC_PCR, ((periph_id & AT91C_PMC_PID)));
	val = read_pmc(PMC_PCR);
	val &= ~AT91C_PMC_DIV;
	val |= AT91C_PMC_DIV_(divider) | AT91C_PMC_CMD | AT91C_PMC_EN;
	write_pmc(PMC_PCR, val);

	return 0;
}

int pmc_disable_periph_clock(unsigned int periph_id)
{
	unsigned int val;

	write_pmc(PMC_PCR, (periph_id & AT91C_PMC_PID));
	val = read_pmc(PMC_PCR);
	val &= ~AT91C_PMC_EN;
	val |= AT91C_PMC_CMD;
	write_pmc(PMC_PCR, val);

	return 0;
}

unsigned int pmc_periph_clock_get_rate(unsigned int periph_id)
{
#ifdef CONFIG_SAMA7G5
	unsigned int val, mck_rate, mckid, div;

	write_pmc(PMC_PCR, (periph_id & AT91C_PMC_PID));
	val = read_pmc(PMC_PCR);
	mckid = (val & AT91C_PMC_MCKID) >> 16;
	div = (val & AT91C_PMC_DIV) >> 14;

	mck_rate = pmc_mck_get_rate(mckid);

	return mck_rate / (1 << div);
#else
	return at91_get_ahb_clock();
#endif
}
