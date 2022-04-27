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
#include "arch/at91_pmc/pmc.h"

void pmc_mck_cfg_set(unsigned int mckid, unsigned int bits, unsigned int mask)
{
	unsigned int tmp = read_pmc(PMC_MCKR) & mask;
	unsigned int updates = tmp ^ bits;

	/*
	 * Program the PRES field in the PMC_MCKR register
	 */
#if defined(CONFIG_AT91SAM9X5) || defined(CONFIG_AT91SAM9N12) \
	|| defined(CONFIG_SAMA5D3X) || defined(CONFIG_SAMA5D4) \
	|| defined(CONFIG_SAMA5D2)
	if (updates & AT91C_PMC_ALT_PRES) {
		tmp &= (~(0x1 << 13));
		tmp &= (~AT91C_PMC_ALT_PRES);
		tmp |= (bits & AT91C_PMC_ALT_PRES);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}
#else
	if (updates & AT91C_PMC_PRES) {
		tmp &= (~(0x1 << 13));
		tmp &= (~AT91C_PMC_PRES);
		tmp |= (bits & AT91C_PMC_PRES);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}
#endif

	/*
	 * Program the MDIV field in the PMC_MCKR register
	 */
	if (updates & AT91C_PMC_MDIV) {
		tmp = read_pmc(PMC_MCKR);
		tmp &= (~AT91C_PMC_MDIV);
		tmp |= (bits & AT91C_PMC_MDIV);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}

	/*
	 * Program the PLLADIV2 field in the PMC_MCKR register
	 */
	if (updates & AT91C_PMC_PLLADIV2) {
		tmp = read_pmc(PMC_MCKR);
		tmp &= (~AT91C_PMC_PLLADIV2);
		tmp |= (bits & AT91C_PMC_PLLADIV2);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}

	/*
	 * Program the H32MXDIV field in the PMC_MCKR register
	 */
	if (updates & AT91C_PMC_H32MXDIV) {
		tmp = read_pmc(PMC_MCKR);
		tmp &= (~AT91C_PMC_H32MXDIV);
		tmp |= (bits & AT91C_PMC_H32MXDIV);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}

	/*
	 * Program the CSS field in the PMC_MCKR register,
	 * wait for MCKRDY bit to be set in the PMC_SR register
	 */
	if (updates & AT91C_PMC_CSS) {
		tmp = read_pmc(PMC_MCKR);
		tmp &= (~AT91C_PMC_CSS);
		tmp |= (bits & AT91C_PMC_CSS);
		write_pmc(PMC_MCKR, tmp);
		while (!(read_pmc(PMC_SR) & AT91C_PMC_MCKRDY)) ;
	}
}

int pmc_mck_check_h32mxdiv(void)
{
#ifdef CONFIG_CPU_HAS_H32MXDIV
	return read_pmc(PMC_MCKR) & AT91C_PMC_H32MXDIV;
#else
	return 0;
#endif
}

unsigned int pmc_mck_get_rate(unsigned int mckid)
{
	return MASTER_CLOCK;
}
