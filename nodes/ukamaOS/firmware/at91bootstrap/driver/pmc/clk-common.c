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
#include "clk-common.h"
#include "timer.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_sfr.h"
#include "rstc.h"
#include "debug.h"
#include "pmc.h"

void lowlevel_clock_init()
{
	unsigned long tmp;

	/*
	 * Switch the master clock to the slow clock without modifying other
	 * parameters. It is assumed that ROM code set H32MXDIV, PLLADIV2,
	 * PCK_DIV3.
	 */
	pmc_mck_cfg_set(0, AT91C_PMC_CSS_SLOW_CLK, AT91C_PMC_CSS);

#ifdef CONFIG_SAMA5D3X_CMP
	/*
	 * On the sama5d3x_cmp board, a phy is not in the proper reset state
	 * after power-up, additional reset
	 */
	rstc_external_reset();
#endif

#if defined(CONFIG_AT91SAM9X5) || defined(CONFIG_AT91SAM9N12) \
	|| defined(CONFIG_SAMA5D3X) || defined(CONFIG_SAMA5D4) \
	|| defined(CONFIG_SAMA5D2) || defined(CONFIG_SAM9X60) \
	|| defined(CONFIG_SAMA7G5)
	/*
	 * Enable the Main Crystal Oscillator
	 * tST_max = 2ms
	 * Startup Time: 32768 * 2ms / 8 = 8
	 */
	tmp = read_pmc(PMC_MOR);
	tmp &= (~AT91C_CKGR_MOSCXTST);
	tmp &= (~AT91C_CKGR_KEY);
	tmp |= AT91C_CKGR_MOSCXTEN;
	tmp |= AT91_CKGR_MOSCXTST_SET(8);
	tmp |= AT91C_CKGR_PASSWD;
	write_pmc(PMC_MOR, tmp);

	while (!(read_pmc(PMC_SR) & AT91C_PMC_MOSCXTS))
		;

#if defined(CONFIG_SAMA5D2) || defined(CONFIG_SAMA7G5)
	/* Enable a measurement of the Main Cristal Oscillator */
	tmp = read_pmc(PMC_MCFR);
	tmp |= AT91C_CKGR_CCSS_XTAL_OSC;
	tmp |= AT91C_CKGR_RCMEAS;
	write_pmc(PMC_MCFR, tmp);

	while (!(tmp = read_pmc(PMC_MCFR) & AT91C_CKGR_MAINRDY))
		;

#if defined(CONFIG_SAMA7G5) && 0
	tmp = (tmp & AT91C_CKGR_MAINF) * 32768 / 16;
	if (tmp != pmc_mck_get_rate(0))
		return;
#endif
#endif

	/* Switch from internal 12MHz RC to the Main Cristal Oscillator */
	tmp = read_pmc(PMC_MOR);
#ifdef CONFIG_MCK_BYPASS
	tmp |= (AT91C_CKGR_MOSCXTBY);
#endif
	tmp &= (~AT91C_CKGR_KEY);
	tmp |= AT91C_CKGR_PASSWD;
	write_pmc(PMC_MOR, tmp);

	tmp = read_pmc(PMC_MOR);
	tmp |= AT91C_CKGR_MOSCSEL;
	tmp &= (~AT91C_CKGR_KEY);
	tmp |= AT91C_CKGR_PASSWD;
	write_pmc(PMC_MOR, tmp);

	while (!(read_pmc(PMC_SR) & AT91C_PMC_MOSCSELS))
		;

#if !defined(CONFIG_SAMA5D4) && !defined(CONFIG_SAMA5D2) \
	&& !defined(CONFIG_SAMA7G5) && !defined(CONFIG_AT91SAM9N12)
	/* Disable the 12MHz RC Oscillator */
	tmp = read_pmc(PMC_MOR);
	tmp &= (~AT91C_CKGR_MOSCRCEN);
	tmp &= (~AT91C_CKGR_KEY);
	tmp |= AT91C_CKGR_PASSWD;
	write_pmc(PMC_MOR, tmp);
#endif

#else
	/*
	 * Enable the Main Crystal Oscillator
	 * tST_max = 2ms
	 * Startup Time: 32768 * 2ms / 8 = 8
	 */
	tmp = read_pmc(PMC_MOR);
	tmp &= (~AT91C_CKGR_MOSCXTST);
	tmp |= AT91C_CKGR_MOSCXTEN;
	tmp |= AT91_CKGR_MOSCXTST_SET(8);
	write_pmc(PMC_MOR, tmp);

	while (!(read_pmc(PMC_SR) & AT91C_PMC_MOSCXTS))
		;
#endif /* CONFIG_AT91SAM9X5 || CONFIG_AT91SAM9N12 || CONFIG_SAMA5D3X ||
        * CONFIG_SAMA5D4 || CONFIG_SAMA5D2 || CONFIG_SAM9X60 */

	/* After stablization, switch to Main Clock */
#if defined(CONFIG_AT91SAM9X5) || defined(CONFIG_AT91SAM9N12) \
	|| defined(CONFIG_SAMA5D3X) || defined(CONFIG_SAMA5D4) \
	|| defined(CONFIG_SAMA5D2)
	pmc_mck_cfg_set(0, AT91C_PMC_CSS_MAIN_CLK | AT91C_PMC_PRES_CLK,
			AT91C_PMC_CSS | AT91C_PMC_ALT_PRES);
#else
	pmc_mck_cfg_set(0, AT91C_PMC_CSS_MAIN_CLK | AT91C_PMC_PRES_CLK,
			AT91C_PMC_CSS | AT91C_PMC_PRES);
#endif

	return;
}

int pmc_cfg_pck(unsigned char x, unsigned int clk_sel, unsigned int prescaler)
{
	write_pmc(PMC_PCKR + x * 4, clk_sel | prescaler);
	write_pmc(PMC_SCER, 1 << (x + 8));
	while (!(read_pmc(PMC_SR) & (1 << (x + 8)))) ;

	return 0;
}

void pmc_enable_system_clock(unsigned int clock_id)
{
	 write_pmc(PMC_SCER, clock_id);
}

void pmc_disable_system_clock(unsigned int clock_id)
{
	 write_pmc(PMC_SCDR, clock_id);
};

void pmc_set_smd_clock_divider(unsigned int divider)
{
	unsigned int tmp = read_pmc(PMC_SMD);

	tmp &= ~AT91C_PMC_SMDDIV;
	tmp |= AT91C_PMC_SMDDIV_(divider);

	write_pmc(PMC_SMD, tmp);
}

unsigned int at91_get_ahb_clock(void)
{
	if (pmc_mck_check_h32mxdiv())
		return MASTER_CLOCK / 2;

	return MASTER_CLOCK;
}

unsigned int pmc_mainck_get_rate(void)
{
#if defined(CONFIG_AT91SAM9X5) || defined(CONFIG_AT91SAM9N12) || \
	defined(CONFIG_SAMA5D3X) || defined(CONFIG_SAMA5D4) || \
	defined(CONFIG_SAMA5D2) || defined(CONFIG_SAM9X60) || \
	defined(CONFIG_SAMA7G5)
	unsigned int val = read_pmc(PMC_MOR);

	if (val & AT91C_CKGR_MOSCSEL)
		/* Crystal oscillator. */
#if defined(BOARD_MAINOSC)
		return BOARD_MAINOSC;
#else
		return 0;
#endif
	else
		/* Embedded RC oscillator. */
		return 12000000;
#elif defined(BOARD_MAINOSC)
	return BOARD_MAINOSC;
#endif

	return 0;
}

#if defined(CONFIG_ENTER_NWD)
unsigned int pmc_read_reg(unsigned int reg_offset)
{
	if (reg_offset <= PMC_OCR)
		return read_pmc(reg_offset);
	else
		return 0xffffffff;
}

void pmc_enable_clock(unsigned int per_id)
{
	unsigned long regval;

	regval = AT91C_PMC_CMD | AT91C_PMC_EN;
	regval |= (per_id & AT91C_PMC_PID);

	/* No need to use Peripheral Divisor value
	 * peripherals are supposed to work @ 90MHz
	 */
	write_pmc(PMC_PCR, regval);
}

void pmc_disable_clock(unsigned int per_id)
{
	unsigned long regval;

	regval = AT91C_PMC_CMD;
	regval |= (per_id & AT91C_PMC_PID);
	write_pmc(PMC_PCR, regval);
}

int pmc_periph_clk(unsigned int periph_id, unsigned int is_on)
{
	if ((periph_id > AT91C_ID_FIQ) && (periph_id < AT91C_ID_COUNTS)) {
		if (is_on)
			pmc_enable_clock(periph_id);
		else
			pmc_disable_clock(periph_id);
		return 0;
	} else {
		return -1;
	}
}

unsigned int sys_mask_to_per_id(unsigned int sys_mask)
{
	switch (sys_mask) {
	case AT91C_PMC_LCDCK:
		return AT91C_ID_LCDC;
	case AT91C_PMC_SMDCK:
		return AT91C_ID_SMD;
	case AT91C_PMC_UHP:
		return AT91C_ID_UHPHS;
	case AT91C_PMC_UDP:
		return AT91C_ID_UDPHS;
	}

	return 0xffffffff;
}

int pmc_sys_clk(unsigned int sys_clock_mask, unsigned int is_on)
{
	if (sys_clock_mask & AVAILABLE_SYS_CLK) {
		if (is_on)
			write_pmc(PMC_SCER, sys_clock_mask);
		else
			write_pmc(PMC_SCDR, sys_clock_mask);
		return 0;
	} else {
		return -1;
	}
}

void pmc_pck_setup(unsigned int reg_offset, unsigned int reg_value)
{
	write_pmc(reg_offset, reg_value);
}

unsigned int pmc_usb_setup(void)
{
	unsigned int usbr = AT91C_PMC_USBS_USB_UPLL;

	/* Setup divider by 10 to reach 48 MHz */
	usbr |= ((10 - 1) << 8) & AT91C_PMC_USBDIV;

	write_pmc(PMC_USB, usbr);

	return usbr;
}

void pmc_smd_setup(unsigned int val)
{
	val &= (AT91C_PMC_SMDS | AT91C_PMC_SMDDIV);

	write_pmc(PMC_SMD, val);
}
#endif /* CONFIG_ENTER_NWD */
