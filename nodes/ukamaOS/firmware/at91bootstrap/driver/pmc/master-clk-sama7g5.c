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
#include "types.h"
#include "arch/at91_pmc/pmc.h"

static unsigned long pmc_mck_mck0_rate(unsigned int mckr)
{
	unsigned int css  = mckr & AT91C_PMC_CSS;
	unsigned int pres = mckr & AT91C_PMC_PRES;
	unsigned int mdiv = mckr & AT91C_PMC_MDIV;
	unsigned long rate = 0;

	switch (css) {
	case AT91C_PMC_CSS_SLOW_CLK:
		rate = 32768;
		break;
	case AT91C_PMC_CSS_MAIN_CLK:
		rate = pmc_mainck_get_rate();
		break;
	case AT91C_PMC_CSS_CPUPLL_CLK:
	case AT91C_PMC_CSS_SYSPLL_CLK:
		rate = pmc_get_pll_freq(css - AT91C_PMC_CSS_CPUPLL_CLK);
		break;
	default:
		return rate;
	}

	switch (pres) {
	case AT91C_PMC_PRES_CLK:
	case AT91C_PMC_PRES_CLK_2:
	case AT91C_PMC_PRES_CLK_4:
	case AT91C_PMC_PRES_CLK_8:
	case AT91C_PMC_PRES_CLK_16:
	case AT91C_PMC_PRES_CLK_32:
	case AT91C_PMC_PRES_CLK_64:
		rate = rate >> (pres >> 4);
		break;
	case AT91C_PMC_PRES_CLK_3:
		rate /= 3;
		break;
	default:
		return rate;
	}

	switch (mdiv) {
	case AT91C_PMC_MDIV_1:
	case AT91C_PMC_MDIV_2:
	case AT91C_PMC_MDIV_4:
		rate = rate >> (mdiv >> 8);
		break;
	case AT91C_PMC_MDIV_3:
		rate /= 3;
		break;
	case AT91C_PMC_MDIV_5:
		rate /= 5;
		break;
	default:
		return rate;
	}

	return rate;
}

static void pmc_mck_mck0_set(unsigned int updates, unsigned int bits,
			     bool css_first)
{
	unsigned int tmp;

	/* Update CSS/MDIV. */
	if (css_first) {
		if (updates & AT91C_PMC_CSS) {
			tmp = read_pmc(PMC_MCKR);
			tmp &= ~AT91C_PMC_CSS;
			tmp |= (bits & AT91C_PMC_CSS);
			write_pmc(PMC_MCKR, tmp);
			while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
		}
	} else {
		if (updates & AT91C_PMC_MDIV) {
			tmp = read_pmc(PMC_MCKR);
			tmp &= ~AT91C_PMC_MDIV;
			tmp |= (bits & AT91C_PMC_MDIV);
			write_pmc(PMC_MCKR, tmp);
			while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
		}
	}

	/* Update PRES. */
	if (updates & AT91C_PMC_PRES) {
		tmp = read_pmc(PMC_MCKR);
		tmp &= ~AT91C_PMC_PRES;
		tmp |= (bits & AT91C_PMC_PRES);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}

	/* Update MDIV/CSS. */
	if (css_first) {
		if (updates & AT91C_PMC_MDIV) {
			tmp = read_pmc(PMC_MCKR);
			tmp &= ~AT91C_PMC_MDIV;
			tmp |= (bits & AT91C_PMC_MDIV);
			write_pmc(PMC_MCKR, tmp);
			while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
		}
	} else {
		if (updates & AT91C_PMC_CSS) {
			tmp = read_pmc(PMC_MCKR);
			tmp &= ~AT91C_PMC_CSS;
			tmp |= (bits & AT91C_PMC_CSS);
			write_pmc(PMC_MCKR, tmp);
			while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
		}
	}
}

void pmc_mck_cfg_set(unsigned int mckid, unsigned int bits, unsigned int mask)
{
	unsigned int updates, tmp;
	unsigned long cur_rate = 0, new_rate = 0;

	/* MCK0 is programmed via PMC_MCKR. */
	if (!mckid) {
		tmp = read_pmc(PMC_MCKR);
		updates = (tmp & mask) ^ bits;

		if (updates & AT91C_PMC_CSS) {
			/*
			 * We have different progamming order b/w CSS, PRES,
			 * MDIV based on selected clock source.
			 */
			cur_rate = pmc_mck_mck0_rate(tmp);
			new_rate = pmc_mck_mck0_rate(bits & mask);

		}

		pmc_mck_mck0_set(updates, bits, cur_rate < new_rate);

		return;
	}

	/* Read current config. */
	write_pmc(PMC_MCR, mckid);
	tmp = read_pmc(PMC_MCR);
	updates = (tmp & mask) ^ bits;

	/* Prepare new config. */
	if (updates & AT91C_MCR_DIV) {
		tmp &= ~AT91C_MCR_DIV;
		tmp |= (bits & AT91C_MCR_DIV);
	}
	if (updates & AT91C_MCR_CSS) {
		tmp &= ~AT91C_MCR_CSS;
		tmp |= (bits & AT91C_MCR_CSS);
	}
	if (updates & AT91C_MCR_EN) {
		tmp &= ~AT91C_MCR_EN;
		tmp |= (bits & AT91C_MCR_EN);
	}

	/* Update it. */
	write_pmc(PMC_MCR, AT91C_MCR_CMD | mckid | tmp);

	while ((updates & AT91C_MCR_CSS) &&
	       !(read_pmc(PMC_SR) & AT91C_PMC_MCKXRDY)) ;
}

unsigned int pmc_mck_get_rate(unsigned int mckid)
{
	unsigned int clock_source, div, tmp, rate = 0;

	if (!mckid) {
		tmp = read_pmc(PMC_MCKR);
		return pmc_mck_mck0_rate(tmp);
	}

	write_pmc(PMC_MCR, mckid);
	tmp = read_pmc(PMC_MCR);
	clock_source = tmp & AT91C_MCR_CSS;
	div = (tmp & AT91C_MCR_DIV) >> 8;

	switch (clock_source) {
	case AT91C_MCR_CSS_MD_SLOW_CLK:
	case AT91C_MCR_CSS_TD_SLOW_CLK:
		rate = 32768;
		break;
	case AT91C_MCR_CSS_MAIN_CLK:
#ifdef BOARD_MAINOSC
		rate = BOARD_MAINOSC;
#endif
		break;
	case AT91C_MCR_CSS_MCK0_CLK:
		tmp = read_pmc(PMC_MCKR);
		rate = pmc_mck_mck0_rate(tmp);
		break;
	case AT91C_MCR_CSS_SYSPLL_CLK:
	case AT91C_MCR_CSS_DDRPLL_CLK:
	case AT91C_MCR_CSS_IMGPLL_CLK:
	case AT91C_MCR_CSS_BAUDPLL_CLK:
	case AT91C_MCR_CSS_AUDIOPLL_CLK:
	case AT91C_MCR_CSS_ETHPLL_CLK:
		tmp = ((clock_source - AT91C_MCR_CSS_SYSPLL_CLK) >> 16) + 1;
		rate = pmc_get_pll_freq(tmp);
		break;
	default:
		break;
	}

	return rate >> div;
}

int pmc_mck_check_h32mxdiv(void)
{
	return 0;
}
