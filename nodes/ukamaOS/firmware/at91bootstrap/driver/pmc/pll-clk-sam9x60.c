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
#include "common.h"
#include "div.h"
#include "pmc.h"
#include "timer.h"
#include "arch/at91_pmc/pmc.h"

static struct pmc_pll_cfg config[PLL_ID_MAX] = { 0 };

void pmc_sam9x60_cfg_pll(unsigned int pll_id, struct pmc_pll_cfg *cfg)
{

	unsigned int reg;

	if (pll_id < 0 || pll_id >= PLL_ID_MAX)
		return;

#ifdef CONFIG_PMC_V2
	if (pll_id == PLL_ID_UPLL) {
		if (cfg->div != 1)
			return;
	}
#endif

	reg = read_pmc(PMC_PLL_UPDT);
	reg &= ~(AT91C_PLL_UPDT_STUPTIM
			| AT91C_PLL_UPDT_UPDATE
			| AT91C_PLL_UPDT_ID);
	reg |= (AT91C_PLL_UPDT_STUPTIM_(cfg->count)
		| AT91C_PLL_UPDT_ID_(pll_id));
	write_pmc(PMC_PLL_UPDT, reg);

	write_pmc(PMC_PLL_ACR, cfg->acr);
	write_pmc(PMC_PLL_CTRL1, AT91C_PLL_CTRL1_MUL_(cfg->mul)
			| AT91C_PLL_CTRL1_FRACR_(cfg->fracr));

#ifdef CONFIG_PMC_V2
	if (pll_id == PLL_ID_UPLL) {
		reg = read_pmc(PMC_PLL_ACR);
		reg |= AT91C_PLL_ACR_UTMIBG;
		write_pmc(PMC_PLL_ACR, reg);
		udelay(10);
		reg = read_pmc(PMC_PLL_ACR);
		reg |= AT91C_PLL_ACR_UTMIVR;
		write_pmc(PMC_PLL_ACR, reg);
		udelay(10);
	}
#endif

	reg = read_pmc(PMC_PLL_UPDT);
	reg |= AT91C_PLL_UPDT_UPDATE;
	write_pmc(PMC_PLL_UPDT, reg);
	reg = read_pmc(PMC_PLL_CTRL0);
	reg &= (~AT91C_PLL_CTRL0_DIVPMC);
	reg &= (~AT91C_PLL_CTRL0_DIVIO);
	reg |= (AT91C_PLL_CTRL0_ENLOCK
		| AT91C_PLL_CTRL0_ENPLL
		| AT91C_PLL_CTRL0_DIVPMC_(cfg->div)
		| AT91C_PLL_CTRL0_DIVIO_(cfg->divio)
		| AT91C_PLL_CTRL0_ENIOPLLCK
		| AT91C_PLL_CTRL0_ENPLLCK);
	write_pmc(PMC_PLL_CTRL0, reg);

	reg = read_pmc(PMC_PLL_UPDT);
	reg |= AT91C_PLL_UPDT_UPDATE;
	write_pmc(PMC_PLL_UPDT, reg);

	while ((read_pmc(PMC_PLL_ISR0) & (AT91C_PLL_ISR0_LOCKA << pll_id))
			!= (AT91C_PLL_ISR0_LOCKA << pll_id))
		;

	config[pll_id] = *(struct pmc_pll_cfg *) cfg;
}

#ifdef CONFIG_PMC_V2
static int pmc_pll_get_parent_props(unsigned int pll_id,
				    unsigned int *parent_rate,
				    unsigned int *divider)
{
	switch (pll_id) {
	case PLL_ID_PLLA:
		*parent_rate = pmc_mainck_get_rate();
		*divider = config[pll_id].div + 1;
		break;
#ifdef BOARD_MAINOSC
	case PLL_ID_UPLL:
		*parent_rate = BOARD_MAINOSC;
		*divider = 2;
		break;
#endif
	}

	return 0;
}
#else
static int pmc_pll_get_parent_props(unsigned int pll_id,
				    unsigned int *parent_rate,
				    unsigned int *divider)
{
	*divider = config[pll_id].div + 1;

	switch (pll_id) {
	case PLL_ID_CPUPLL:
	case PLL_ID_SYSPLL:
	case PLL_ID_DDRPLL:
	case PLL_ID_IMGPLL:
	case PLL_ID_BAUDPLL:
		*parent_rate = pmc_mainck_get_rate();
		break;

#ifdef BOARD_MAINOSC
	case PLL_ID_AUDIOPLL:
	case PLL_ID_ETHPLL:
		*parent_rate = BOARD_MAINOSC;
		break;
#endif
	default:
		return -1;
	}

	return 0;
}
#endif

unsigned int pmc_get_pll_freq(unsigned int pll_id)
{
	unsigned int freq, core_freq, divider = 0, parent_rate = 0;
	unsigned long long frac_freq;
	int ret;

	ret = pmc_pll_get_parent_props(pll_id, &parent_rate, &divider);
	if (ret)
		return ret;

	core_freq = parent_rate * (config[pll_id].mul + 1);
	frac_freq = ((parent_rate >> 20) * config[pll_id].fracr) >> 2;
	core_freq += frac_freq;

	freq = div(core_freq, divider);

	return freq;
}
