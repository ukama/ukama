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
#include "common.h"
#include "hardware.h"
#include "arch/at91_ccfg.h"
#include "arch/at91_matrix.h"
#include "arch/at91_rstc.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_smc.h"
#include "arch/at91_pio.h"
#include "arch/at91_sdramc.h"
#include "spi.h"
#include "gpio.h"
#include "pmc.h"
#include "usart.h"
#include "debug.h"
#include "sdramc.h"
#include "timer.h"
#include "watchdog.h"
#include "at91sam9rlek.h"

static void at91_dbgu_hw_init(void)
{
	/* Configure DBGU pin */
	const struct pio_desc dbgu_pins[] = {
		{"RXD", AT91C_PIN_PA(21), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TXD", AT91C_PIN_PA(22), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the dbgu pins */
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(dbgu_pins);
}

static void initialize_dbgu(void)
{
	at91_dbgu_hw_init();

	usart_init(BAUDRATE(MASTER_CLOCK, 115200));
}

#ifdef CONFIG_SDRAM
static void sdramc_hw_init(void)
{
	/* Configure sdramc pins */
	const struct pio_desc sdramc_pins[] = {
		{"D16", AT91C_PIN_PB(16), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D17", AT91C_PIN_PB(17), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D18", AT91C_PIN_PB(18), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D19", AT91C_PIN_PB(19), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D20", AT91C_PIN_PB(20), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D21", AT91C_PIN_PB(21), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D22", AT91C_PIN_PB(22), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D23", AT91C_PIN_PB(23), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D24", AT91C_PIN_PB(24), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D25", AT91C_PIN_PB(25), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D26", AT91C_PIN_PB(26), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D27", AT91C_PIN_PB(27), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D28", AT91C_PIN_PB(28), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D29", AT91C_PIN_PB(29), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D30", AT91C_PIN_PB(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D31", AT91C_PIN_PB(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *) 0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the sdramc pins */
	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(sdramc_pins);
}

static void sdramc_init(void)
{
	struct sdramc_register sdramc_config;
	unsigned int reg;

#ifdef MCK_100
	sdramc_config.cr = (AT91C_SDRAMC_NC_9
				| AT91C_SDRAMC_NR_13
				| AT91C_SDRAMC_CAS_2
				| AT91C_SDRAMC_NB_4_BANKS
				| AT91C_SDRAMC_DBW_32_BITS
				| AT91C_SDRAMC_TWR_2
				| AT91C_SDRAMC_TRC_7
				| AT91C_SDRAMC_TRP_2
				| AT91C_SDRAMC_TRCD_2
				| AT91C_SDRAMC_TRAS_5
				| AT91C_SDRAMC_TXSR_8);
#else
	sdramc_config.cr = (AT91C_SDRAMC_NC_9
				| AT91C_SDRAMC_NR_13
				| AT91C_SDRAMC_CAS_3
				| AT91C_SDRAMC_NB_4_BANKS
				| AT91C_SDRAMC_DBW_32_BITS
				| AT91C_SDRAMC_TWR_2
				| AT91C_SDRAMC_TRC_9
				| AT91C_SDRAMC_TRP_3
				| AT91C_SDRAMC_TRCD_3
				| AT91C_SDRAMC_TRAS_6
				| AT91C_SDRAMC_TXSR_10);
#endif /* #ifdef MCK_100 */

	sdramc_config.tr = (MASTER_CLOCK * 7) / 1000000;
	sdramc_config.mdr = AT91C_SDRAMC_MD_SDRAM;

	sdramc_hw_init();

	/* Initialize the matrix (memory voltage = 3.3) */
	reg = readl(AT91C_BASE_CCFG + CCFG_EBI0CSA);
	reg |= AT91C_EBI_CS1A_SDRAMC;
	writel(reg, AT91C_BASE_CCFG + CCFG_EBI0CSA);

	sdramc_initialize(&sdramc_config, AT91C_BASE_CS1);
}
#endif /* #ifdef CONFIG_SDRAM */

static inline void matrix_writel(const unsigned int value, unsigned int reg)
{
	writel(value, reg + AT91C_BASE_MATRIX);
}

static inline unsigned int matrix_readl(unsigned int reg)
{
	return readl(reg + AT91C_BASE_MATRIX);
}

static void at91_matrix_hw_init(void)
{
	unsigned int reg;

	reg = matrix_readl(MATRIX_SCFG0);
	reg |= AT91C_MATRIX_ARBT_FIXED_PRIORITY;
	reg |= AT91C_MATRIX_FIXED_DEFMSTR_ARM926D;
	matrix_writel(reg, MATRIX_SCFG0);	/* ROM */


	reg = matrix_readl(MATRIX_SCFG4);
	reg |= AT91C_MATRIX_ARBT_FIXED_PRIORITY;
	reg |= AT91C_MATRIX_FIXED_DEFMSTR_ARM926D;
	reg |= AT91C_MATRIX_SLOT_CYCLE_(0x40);
	matrix_writel(reg, MATRIX_SCFG4);	/* EBI0 */

}

#if defined(CONFIG_NANDFLASH_RECOVERY) || defined(CONFIG_DATAFLASH_RECOVERY)
static void recovery_buttons_hw_init(void)
{
	/* Configure recovery button PINs */
	const struct pio_desc recovery_button_pins[] = {
		{"RECOVERY_BUTTON", CONFIG_SYS_RECOVERY_BUTTON_PIN, 0, PIO_PULLUP, PIO_INPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(recovery_button_pins);
}
#endif /* #if defined(CONFIG_NANDFLASH_RECOVERY) || defined(CONFIG_DATAFLASH_RECOVERY) */

#ifdef CONFIG_HW_INIT
void hw_init(void)
{
	/* Disable watchdog */
	at91_disable_wdt();

	/*
	 * At this stage the main oscillator is supposed to be enabled
	 *  PCK = MCK = MOSC
	 */
	/* Configure PLLA = MOSC * (PLL_MULA + 1) / PLL_DIVA */
	pmc_cfg_plla(PLL_SETTINGS);

	/* PCK = PLL = 2 * MCK */
	pmc_mck_cfg_set(0, MCKR_SETTINGS, AT91C_PMC_PRES | AT91C_PMC_MDIV);

	/* Switch MCK on PLLA output */
	pmc_mck_cfg_set(0, MCKR_CSS_SETTINGS,
			AT91C_PMC_PRES | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Enable External Reset */
	writel(AT91C_RSTC_KEY_UNLOCK | AT91C_RSTC_URSTEN, AT91C_BASE_RSTC + RSTC_RMR);

	/* Init timer */
	timer_init();

	/* Initialize dbgu */
	initialize_dbgu();

	/* Initialize matrix configuration */
	at91_matrix_hw_init();

#ifdef CONFIG_SDRAM
	/* Configure SDRAM Controller */
	sdramc_init();
#endif

#if defined(CONFIG_NANDFLASH_RECOVERY) || defined(CONFIG_DATAFLASH_RECOVERY)
	/* Init the recovery buttons pins */
	recovery_buttons_hw_init();
#endif
}
#endif /* #ifdef CONFIG_HW_INIT */

#ifdef CONFIG_DATAFLASH
void at91_spi0_hw_init(void)
{
	/* Configure spi0 PINs */
	const struct pio_desc spi0_pins[] = {
		{"MISO",	AT91C_PIN_PA(25),	0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MOSI",	AT91C_PIN_PA(26),	0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SPCK",	AT91C_PIN_PA(27),	0, PIO_DEFAULT, PIO_PERIPH_A},
		{"NPCS",	CONFIG_SYS_SPI_PCS,	1, PIO_PULLUP, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the spi0 pins */
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(spi0_pins);

	/* Enable the spi0 clock */
	pmc_enable_periph_clock(AT91C_ID_SPI, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_DATAFLASH */

#ifdef CONFIG_SDCARD
void at91_mci0_hw_init(void)
{
	const struct pio_desc mci_pins[] = {
		{"MC_CK",	AT91C_PIN_PA(2), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MC_CDA",	AT91C_PIN_PA(1), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MC_DA0",	AT91C_PIN_PA(0), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MC_DA1",	AT91C_PIN_PA(3), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MC_DA2",	AT91C_PIN_PA(4), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MC_DA3",	AT91C_PIN_PA(5), 0, PIO_PULLUP, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},

	};

	/* Configure the PIO controller */
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(mci_pins);

	/* Enable the clock */
	pmc_enable_periph_clock(AT91C_ID_MCI, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_SDCARD */

#ifdef CONFIG_NANDFLASH
void nandflash_hw_init(void)
{
	unsigned int reg;

	/* Configure NANDFlash pins*/
	const struct pio_desc nand_pins[] = {
		{"NANDALE",	AT91C_PIN_PB(2),		0, PIO_PULLUP, PIO_PERIPH_A},
		{"NANDCLE",	AT91C_PIN_PB(3),		0, PIO_PULLUP, PIO_PERIPH_A},
		{"NANDOE",	AT91C_PIN_PB(4),		0, PIO_PULLUP, PIO_PERIPH_A},
		{"NANDWE",	AT91C_PIN_PB(5),		0, PIO_PULLUP, PIO_PERIPH_A},
		{"NANDCS",	CONFIG_SYS_NAND_ENABLE_PIN,	1, PIO_PULLUP, PIO_OUTPUT},
		{(char *)0, 	0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Setup Smart Media, first enable the address range of CS3 in HMATRIX user interface  */
	reg = readl(AT91C_BASE_CCFG + CCFG_EBI0CSA);
	reg |= AT91C_EBI_CS3A_SM;
	writel(reg, AT91C_BASE_CCFG + CCFG_EBI0CSA);

	/* Configure SMC CS3 */
	writel((AT91C_SMC_NWESETUP_(1)
		| AT91C_SMC_NCS_WRSETUP_(0)
		| AT91C_SMC_NRDSETUP_(1)
		| AT91C_SMC_NCS_RDSETUP_(0)),
		AT91C_BASE_SMC + SMC_SETUP3);

	writel((AT91C_SMC_NWEPULSE_(3)
		| AT91C_SMC_NCS_WRPULSE_(3)
		| AT91C_SMC_NRDPULSE_(3)
		| AT91C_SMC_NCS_RDPULSE_(3)),
		AT91C_BASE_SMC + SMC_PULSE3);

	writel((AT91C_SMC_NWECYCLE_(5)
		|  AT91C_SMC_NRDCYCLE_(5)),
		AT91C_BASE_SMC + SMC_CYCLE3);

	writel((AT91C_SMC_READMODE
		| AT91C_SMC_WRITEMODE
		| AT91C_SMC_NWAITM_NWAIT_DISABLE
		| AT91C_SMC_DBW_WIDTH_BITS_8
		| AT91_SMC_TDF_(2)),
		AT91C_BASE_SMC + SMC_CTRL3);

	/* Configure the NANDFlash pins */
	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(nand_pins);
}
#endif /* #ifdef CONFIG_NANDFLASH */
