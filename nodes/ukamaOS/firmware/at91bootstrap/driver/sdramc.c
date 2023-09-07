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
#include "hardware.h"
#include "board.h"
#include "arch/at91_sdramc.h"
#include "sdramc.h"

static inline void sdramc_writel(unsigned int reg, const unsigned int value)
{
	writel(value, reg + AT91C_BASE_SDRAMC);
}

static inline unsigned int sdramc_readl(unsigned int reg)
{
	return readl(reg + AT91C_BASE_SDRAMC);
}

int sdramc_initialize(struct sdramc_register *sdramc_config,
			unsigned int sdram_address)
{
	unsigned int i;

	/* Step#1 SDRAM feature must be in the configuration register */
	sdramc_writel(SDRAMC_CR, sdramc_config->cr);

	/* Step#1 bis (SAM9X60 SDRAMC ) - initialize CFR1 register */

	if (sdramc_config->cfr1) /* avoid write for products where cfr1 is unavailable */
		sdramc_writel(SDRAMC_CFR1, sdramc_config->cfr1);

	/* Step#2 For mobile SDRAM, temperature-compensated self refresh(TCSR),... */

	/* Step#3 The SDRAM memory type must be set in the Memory Device Register */
	sdramc_writel(SDRAMC_MDR, sdramc_config->mdr);

	/* Step#4 The minimum pause of 200 us is provided to precede any single toggle */
	for (i = 0; i < 1000; i++) ;

	/* Step#5 A NOP command is issued to the SDRAM devices */
	sdramc_writel(SDRAMC_MR, AT91C_SDRAMC_MODE_NOP);
	writel(0x00000000, sdram_address);

	/* Step#6 An All Banks Precharge command is issued to the SDRAM devices  */
	sdramc_writel(SDRAMC_MR, AT91C_SDRAMC_MODE_PRECHARGE);
	writel(0x00000000, sdram_address);

	for (i = 0; i < 10000; i++) ;

	/* Step#7 Eight auto-refresh cycles are provided */
	for (i = 0; i < 8; i++) {
		sdramc_writel(SDRAMC_MR, AT91C_SDRAMC_MODE_AUTO_REFRESH);
		writel(0x00000001 + i, sdram_address + 4 + 4 * i);
	}

	/*  Pause cycles */
	for (i = 0; i < 1000; i++) ;

	/* Step#8 A Mode Register set (MRS) cyscle is issued to program the SDRAM parameters(TCSR, PASR, DS) */
	sdramc_writel(SDRAMC_MR, AT91C_SDRAMC_MODE_LOAD_MODE);
	writel(0xcafedede, sdram_address + 0x24);

	/*  Pause cycles */
	for (i = 0; i < 1000; i++) ;

	/* Step#9 For mobile SDRAM initialization, an Extended Mode Register set cycle is issued to ... */

	/* Step#10 The application must go into Normal Mode, setting Mode to 0 in the Mode Register
	   and perform a write access at any location in the SDRAM. */
	sdramc_writel(SDRAMC_MR, AT91C_SDRAMC_MODE_NORMAL);	// Set Normal mode
	writel(0x00000000, sdram_address);	// Perform Normal mode

	/* Step#11 Write the refresh rate into the count field in the SDRAMC Refresh Timer Rgister. */
	sdramc_writel(SDRAMC_TR, sdramc_config->tr);

	return 0;
}
