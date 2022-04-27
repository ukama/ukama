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
#include "clk-common.h"
#include "common.h"
#include "debug.h"
#include "hardware.h"
#include "pmc.h"
#include "timer.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_sfr.h"
#include "board.h"

#if defined(AT91C_BASE_SFR) && !defined(CONFIG_SAMA5D4)
static unsigned long pmc_get_main_clock(void)
{
	unsigned long main_clock;

#ifdef BOARD_MAINOSC
	main_clock = BOARD_MAINOSC;
#else
	unsigned int tmp;

	do {
		tmp = read_pmc(PMC_MCFR);
	} while (!(tmp & AT91C_CKGR_MAINRDY));

	tmp &= AT91C_CKGR_MAINF;
	main_clock = tmp * (CONFIG_SYS_AT91_SLOW_CLOCK / 16);
#endif
	return main_clock;
}

static int clock_freq_in_range(unsigned long freq,
			     unsigned long norm, unsigned int delta)
{
	if (delta > norm)
		return 0;

	if ((freq >= (norm - delta)) && (freq <= (norm + delta)))
		return 1;
	else
		return 0;
}

static int pmc_configure_utmi_ref_clk(void)
{
	unsigned long main_clock;
	unsigned int utmi_ref_clk_freq, tmp;
	unsigned long delta = 50000;

	/*
	 * If mainck rate is different from 12 MHz, we have to configure
	 * the FREQ field of the SFR_UTMICKTRIM register to generate properly
	 * the utmi clock.
	 */
	main_clock = pmc_get_main_clock();
	if (clock_freq_in_range(main_clock, 12000000, delta)) {
		utmi_ref_clk_freq = 0;
	} else if (clock_freq_in_range(main_clock, 16000000, delta)) {
		utmi_ref_clk_freq = 1;
	} else if (clock_freq_in_range(main_clock, 24000000, delta)) {
		utmi_ref_clk_freq = 2;
	/*
	 * Not supported on SAMA5D2 but it's not an issue since MAINCK
	 * maximum value is 24 MHz.
	 */
	} else if (clock_freq_in_range(main_clock, 48000000, delta)) {
		utmi_ref_clk_freq = 3;
	} else {
		dbg_info("The MAINCK frequency: %d\n", main_clock);
		return -1;
	}

	tmp = readl(SFR_UTMICKTRIM + AT91C_BASE_SFR);
	tmp &= ~AT91C_UTMICKTRIM_FREQ;
	tmp |= utmi_ref_clk_freq;
	writel(tmp, SFR_UTMICKTRIM + AT91C_BASE_SFR);

	return 0;
}
#else
static int pmc_configure_utmi_ref_clk(void)
{
	return 0;
}
#endif /* AT91C_BASE_SFR && !CONFIG_SAMA5D4 */

int pmc_uckr_clk(unsigned int is_on)
{
	unsigned int uckr = read_pmc(PMC_UCKR);
	unsigned int sr;
	int ret;

	if (is_on) {
		ret = pmc_configure_utmi_ref_clk();
		if (ret)
			return ret;

		uckr |= (AT91C_CKGR_UPLLCOUNT_DEFAULT |
			 AT91C_CKGR_UPLLEN | AT91C_CKGR_BIASEN);
		sr = AT91C_PMC_LOCKU;
	} else {
		uckr &= ~(AT91C_CKGR_UPLLEN | AT91C_CKGR_BIASEN);
		sr = 0;
	}

	write_pmc(PMC_UCKR, uckr);

	do {
		udelay(1);
	} while ((read_pmc(PMC_SR) & AT91C_PMC_LOCKU) != sr);

	return 0;
}

