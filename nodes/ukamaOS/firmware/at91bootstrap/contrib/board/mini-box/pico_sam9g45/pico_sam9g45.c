/* ----------------------------------------------------------------------------
 *          Support for pico-SAM9G45 Board   http://arm.mini-box.com/
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Jochen Friedrich <jochen@scram.de>
 *
 * Based on at91sam9m10g45ek.c 
 * Copyright (c) 2008, Atmel Corporation
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
#include "common.h"
#include "hardware.h"
#include "arch/at91_ccfg.h"
#include "arch/at91_rstc.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_smc.h"
#include "arch/at91_pio.h"
#include "arch/at91_ddrsdrc.h"
#include "gpio.h"
#include "pmc.h"
#include "usart.h"
#include "debug.h"
#include "ddramc.h"
#include "timer.h"
#include "watchdog.h"
#include "string.h"
#include "pico_sam9g45.h"

static void at91_dbgu_hw_init(void)
{
	/* Configure DBGU pin */
	const struct pio_desc dbgu_pins[] = {
		{"RXD", AT91C_PIN_PB(12), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TXD", AT91C_PIN_PB(13), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the dbgu pins */
	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(dbgu_pins);
}

static void initialize_dbgu(void)
{
	at91_dbgu_hw_init();
	usart_init(BAUDRATE(MASTER_CLOCK, 115200));
}

#ifdef CONFIG_DDR2
static void ddramc_reg_config(struct ddramc_register *ddramc_config)
{
	ddramc_config->mdr = (AT91C_DDRC2_DBW_16_BITS
			| AT91C_DDRC2_MD_DDR2_SDRAM);

	ddramc_config->cr = (AT91C_DDRC2_NC_DDR10_SDR9	/* 10 column bits */
			| AT91C_DDRC2_NR_14		/* 14 row bits(8K)*/
			| AT91C_DDRC2_CAS_3		/* CAS Latency 3 */
			| AT91C_DDRC2_DISABLE_RESET_DLL);

	ddramc_config->rtr = 0x24B;

	ddramc_config->t0pr = (AT91C_DDRC2_TRAS_(6)	/* 6 * 7.5 = 45 ns */
			| AT91C_DDRC2_TRCD_(2)		/* 2 * 7.5 = 22.5 ns */
			| AT91C_DDRC2_TWR_(2)		/* 2 * 7.5 = 15 ns */
			| AT91C_DDRC2_TRC_(8)		/* 8 * 7.5 = 75 ns */
			| AT91C_DDRC2_TRP_(2)		/* 2 * 7.5 = 22.5 ns */
			| AT91C_DDRC2_TRRD_(1)		/* 1 * 7.5 = 7.5 ns */
			| AT91C_DDRC2_TWTR_(1)		/* 1 clock cycle */
			| AT91C_DDRC2_TMRD_(2));	/* 2 clock cycles */

	ddramc_config->t1pr = (AT91C_DDRC2_TXP_(2)	/* 2 * 7.5 = 15 ns */
			| AT91C_DDRC2_TXSRD_(200)	/* 200 clock cycles */
			| AT91C_DDRC2_TXSNR_(16)	/* 16 * 7.5 = 120 ns */
			| AT91C_DDRC2_TRFC_(14));	/* 14 * 7.5 = 142 ns */

	ddramc_config->t2pr = (AT91C_DDRC2_TRTP_(1)	/* 1 * 7.5 = 7.5 ns */
			| AT91C_DDRC2_TRPA_(0)		/* 0 * 7.5 = 0 ns */
			| AT91C_DDRC2_TXARDS_(7)	/* 7 clock cycles */
			| AT91C_DDRC2_TXARD_(2));	/* 2 clock cycles */
}

static void ddramc_init(void)
{
	unsigned long csa;
	struct ddramc_register  ddramc_reg;

	ddramc_reg_config(&ddramc_reg);

	/* ENABLE DDR2 clock */ 
	pmc_enable_system_clock(AT91C_PMC_DDR);

	/* Chip select 1 is for DDR2/SDRAM */
	csa = readl(AT91C_BASE_CCFG + CCFG_EBICSA);
	csa |= AT91C_EBI_CS1A_SDRAMC;
	csa &= ~AT91C_VDDIOM_SEL_33V;
	writel(csa, AT91C_BASE_CCFG + CCFG_EBICSA);

	/* DDRAM2 Controller initialize */
	ddram_initialize(AT91C_BASE_DDRSDRC0, AT91C_BASE_CS6, &ddramc_reg);
	ddram_initialize(AT91C_BASE_DDRSDRC1, AT91C_BASE_CS1, &ddramc_reg);
}
#endif /* #ifdef CONFIG_DDR2 */

#ifdef CONFIG_HW_INIT
void hw_init(void)
{
	/* Disable watchdog */
	at91_disable_wdt();

	/*
	 * At this stage the main oscillator
	 * is supposed to be enabled PCK = MCK = MOSC
	 */
	pmc_init_pll(0);

	/* Configure PLLA = MOSC * (PLL_MULA + 1) / PLL_DIVA */
	pmc_cfg_plla(PLLA_SETTINGS);

	/* Switch PCK/MCK on Main clock output */
	pmc_mck_cfg_set(0, BOARD_PRESCALER_MAIN_CLOCK,
			AT91C_PMC_PLLADIV2 | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Switch PCK/MCK on PLLA output */
	pmc_mck_cfg_set(0, BOARD_PRESCALER_PLLA,
			AT91C_PMC_PLLADIV2 | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Enable External Reset */
	writel(AT91C_RSTC_KEY_UNLOCK | AT91C_RSTC_URSTEN, AT91C_BASE_RSTC + RSTC_RMR);

	/* Init timer */
	timer_init();

	/* Initialize dbgu */
	initialize_dbgu();

#ifdef CONFIG_DDR2
	/* Initialize DDRAM Controller */
	ddramc_init();
#endif

}
#endif /* #ifdef CONFIG_HW_INIT */

#ifdef CONFIG_SDCARD
#ifdef CONFIG_OF_LIBFDT
void at91_board_set_dtb_name(char *of_name)
{
	strcpy(of_name, "at91-picosam9g45.dtb");
}
#endif

void at91_mci0_hw_init(void)
{
	const struct pio_desc mci_pins[] = {
		{"MCCK", 	AT91C_PIN_PA(0), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCCDA",	AT91C_PIN_PA(1), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCDA0",	AT91C_PIN_PA(2), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCDA1",	AT91C_PIN_PA(3), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCDA2",	AT91C_PIN_PA(4), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCDA3",	AT91C_PIN_PA(5), 0, PIO_PULLUP, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the PIO controller */
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(mci_pins);

	/* Enable the clock */
	pmc_enable_periph_clock(AT91C_ID_MCI0, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_SDCARD */
