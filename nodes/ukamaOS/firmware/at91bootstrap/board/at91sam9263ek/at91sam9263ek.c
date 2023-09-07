/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2009, Atmel Corporation
 *
 * All rights reserved.
 * PSRAM,16 bit SDRAM support courtesy of Epsilon Group
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
#include "psram.h"
#include "timer.h"
#include "watchdog.h"
#include "string.h"
#include "at91sam9263ek.h"

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
	/* Bus Matrix Master Configuration Register */
	matrix_writel(AT91C_MATRIX_ULBT_16_BEAT, MATRIX_MCFG0);		/* OHCI */
	matrix_writel(AT91C_MATRIX_ULBT_8_BEAT, MATRIX_MCFG1);		/* ISI */
	matrix_writel(AT91C_MATRIX_ULBT_8_BEAT, MATRIX_MCFG2);		/* 2D */
	matrix_writel(AT91C_MATRIX_ULBT_8_BEAT, MATRIX_MCFG3);		/* DMAC */
	matrix_writel(AT91C_MATRIX_ULBT_4_BEAT, MATRIX_MCFG4);		/* MACB */
	matrix_writel(AT91C_MATRIX_ULBT_16_BEAT, MATRIX_MCFG5);		/* LCDC */
	matrix_writel(AT91C_MATRIX_ULBT_SINGLE_ACCESS, MATRIX_MCFG6);	/* PDC */
	matrix_writel(AT91C_MATRIX_ULBT_8_BEAT, MATRIX_MCFG7);		/* DBUS */
	matrix_writel(AT91C_MATRIX_ULBT_4_BEAT, MATRIX_MCFG8);		/* IBUS */

	/* Bus Matrix Slave Configuration Registers */
	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_ARM926I
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(32)),
			MATRIX_SCFG0);		/* ROM */

	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_EMAC
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(32)),
			MATRIX_SCFG1);		/* RAM80K */

	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_USB
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(16)),
			MATRIX_SCFG2);		/* RAM16K */

	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_PDC
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(4)),
			MATRIX_SCFG3);		/* PERIPHERALS */

	matrix_writel((AT91C_MATRIX_ARBT_ROUND_ROBIN
			| AT91C_MATRIX_FIXED_DEFMSTR_ARM926I
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(32)),
			MATRIX_SCFG4);		/* EBI0 */

	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_LCDC
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(64)),
			MATRIX_SCFG5);		/* EBI1 */

	matrix_writel((AT91C_MATRIX_ARBT_FIXED_PRIORITY
			| AT91C_MATRIX_FIXED_DEFMSTR_ARM926D
			| AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR
			| AT91C_MATRIX_SLOT_CYCLE_(4)),
			MATRIX_SCFG6);		/* APB */

	/* ROM */
	matrix_writel((AT91C_MATRIX_M0PR_(1)
			| AT91C_MATRIX_M1PR_(0)
			| AT91C_MATRIX_M2PR_(2)
			| AT91C_MATRIX_M3PR_(1)
			| AT91C_MATRIX_M4PR_(0)
			| AT91C_MATRIX_M5PR_(3)
			| AT91C_MATRIX_M6PR_(2)
			| AT91C_MATRIX_M7PR_(3)),
			MATRIX_PRAS0);

	matrix_writel(AT91C_MATRIX_M8PR_(0),
			MATRIX_PRBS0);

	/* RAM80K */
	matrix_writel((AT91C_MATRIX_M0PR_(1)
			| AT91C_MATRIX_M1PR_(2)
			| AT91C_MATRIX_M2PR_(1)
			| AT91C_MATRIX_M3PR_(3)
			| AT91C_MATRIX_M4PR_(0)
			| AT91C_MATRIX_M5PR_(0)
			| AT91C_MATRIX_M6PR_(3)
			| AT91C_MATRIX_M7PR_(0)),
			MATRIX_PRAS1);

	matrix_writel(AT91C_MATRIX_M8PR_(2),
			MATRIX_PRBS1);

	/* RAM16K */
	matrix_writel((AT91C_MATRIX_M0PR_(1)
			| AT91C_MATRIX_M1PR_(0)
			| AT91C_MATRIX_M2PR_(2)
			| AT91C_MATRIX_M3PR_(1)
			| AT91C_MATRIX_M4PR_(0)
			| AT91C_MATRIX_M5PR_(3)
			| AT91C_MATRIX_M6PR_(3)
			| AT91C_MATRIX_M7PR_(2)),
			MATRIX_PRAS2);

	matrix_writel(AT91C_MATRIX_M8PR_(0),
			MATRIX_PRBS2);

	/* PERIPHERALS */
	matrix_writel((AT91C_MATRIX_M0PR_(0)
			| AT91C_MATRIX_M1PR_(1)
			| AT91C_MATRIX_M2PR_(0)
			| AT91C_MATRIX_M3PR_(2)
			| AT91C_MATRIX_M4PR_(1)
			| AT91C_MATRIX_M5PR_(0)
			| AT91C_MATRIX_M6PR_(3)
			| AT91C_MATRIX_M7PR_(2)),
			MATRIX_PRAS3);

	matrix_writel(AT91C_MATRIX_M8PR_(3),
			MATRIX_PRBS3);

#if defined(CONFIG_PSRAM)
	/* EBI0 */
	matrix_writel((AT91C_MATRIX_M0PR_(2)
			| AT91C_MATRIX_M1PR_(1)
			| AT91C_MATRIX_M2PR_(1)
			| AT91C_MATRIX_M3PR_(3)
			| AT91C_MATRIX_M4PR_(0)
			| AT91C_MATRIX_M5PR_(3)
			| AT91C_MATRIX_M6PR_(0)
			| AT91C_MATRIX_M7PR_(0)),
			MATRIX_PRAS4);

	matrix_writel(AT91C_MATRIX_M8PR_(2),
			MATRIX_PRBS4);
#else
	/* EBI0 */
	matrix_writel((AT91C_MATRIX_M0PR_(1)
			| AT91C_MATRIX_M1PR_(3)
			| AT91C_MATRIX_M2PR_(0)
			| AT91C_MATRIX_M3PR_(2)
			| AT91C_MATRIX_M4PR_(3)
			| AT91C_MATRIX_M5PR_(0)
			| AT91C_MATRIX_M6PR_(0)
			| AT91C_MATRIX_M7PR_(1)),
			MATRIX_PRAS4);

	matrix_writel(AT91C_MATRIX_M8PR_(2),
			MATRIX_PRBS4);
#endif /* #if defined(CONFIG_PSRAM) */
	/* EBI1 */
	matrix_writel((AT91C_MATRIX_M0PR_(0)
			| AT91C_MATRIX_M1PR_(1)
			| AT91C_MATRIX_M2PR_(0)
			| AT91C_MATRIX_M3PR_(0)
			| AT91C_MATRIX_M4PR_(3)
			| AT91C_MATRIX_M5PR_(2)
			| AT91C_MATRIX_M6PR_(3)
			| AT91C_MATRIX_M7PR_(2)),
			MATRIX_PRAS5);

	matrix_writel(AT91C_MATRIX_M8PR_(1),
			MATRIX_PRBS5);

	/* APB */
	matrix_writel((AT91C_MATRIX_M0PR_(1)
			| AT91C_MATRIX_M1PR_(0)
			| AT91C_MATRIX_M2PR_(2)
			| AT91C_MATRIX_M3PR_(1)
			| AT91C_MATRIX_M4PR_(0)
			| AT91C_MATRIX_M5PR_(0)
			| AT91C_MATRIX_M6PR_(3)
			| AT91C_MATRIX_M7PR_(3)),
			MATRIX_PRAS4);

	matrix_writel(AT91C_MATRIX_M8PR_(2),
			MATRIX_PRBS4);
}

static void at91_dbgu_hw_init(void)
{
	/* Configure DBGU pin */
	const struct pio_desc dbgu_pins[] = {
		{"RXD", AT91C_PIN_PC(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TXD", AT91C_PIN_PC(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/*  Configure the dbgu pins */
	pmc_enable_periph_clock(AT91C_ID_PIOCDE, PMC_PERIPH_CLK_DIVIDER_NA);
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
		{"D16", AT91C_PIN_PD(16), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D17", AT91C_PIN_PD(17), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D18", AT91C_PIN_PD(18), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D19", AT91C_PIN_PD(19), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D20", AT91C_PIN_PD(20), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D21", AT91C_PIN_PD(21), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D22", AT91C_PIN_PD(22), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D23", AT91C_PIN_PD(23), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D24", AT91C_PIN_PD(24), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D25", AT91C_PIN_PD(25), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D26", AT91C_PIN_PD(26), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D27", AT91C_PIN_PD(27), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D28", AT91C_PIN_PD(28), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D29", AT91C_PIN_PD(29), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D30", AT91C_PIN_PD(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D31", AT91C_PIN_PD(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *) 0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the SDRAMC PINs */
	pmc_enable_periph_clock(AT91C_ID_PIOCDE, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(sdramc_pins);
}

static void sdramc0_init(void)
{
	unsigned int reg;
	struct sdramc_register sdramc_config;

#ifdef CONFIG_SDRAM_16BIT
	sdramc_config.cr = AT91C_SDRAMC_NC_10 | AT91C_SDRAMC_NR_13 | AT91C_SDRAMC_CAS_2
				| AT91C_SDRAMC_NB_4_BANKS | AT91C_SDRAMC_DBW_16_BITS
				| AT91C_SDRAMC_TWR_2 | AT91C_SDRAMC_TRC_7
				| AT91C_SDRAMC_TRP_2 | AT91C_SDRAMC_TRCD_2
				| AT91C_SDRAMC_TRAS_5 | AT91C_SDRAMC_TXSR_8;
#else
	sdramc_config.cr = AT91C_SDRAMC_NC_9 | AT91C_SDRAMC_NR_13 | AT91C_SDRAMC_CAS_2
				| AT91C_SDRAMC_NB_4_BANKS | AT91C_SDRAMC_DBW_32_BITS
				| AT91C_SDRAMC_TWR_2 | AT91C_SDRAMC_TRC_7
				| AT91C_SDRAMC_TRP_2 | AT91C_SDRAMC_TRCD_2
				| AT91C_SDRAMC_TRAS_5 | AT91C_SDRAMC_TXSR_8;
#endif	/* #ifdef CONFIG_SDRAM_16BIT */

	sdramc_config.tr = (MASTER_CLOCK * 7) / 1000000;
	sdramc_config.mdr = AT91C_SDRAMC_MD_SDRAM;

	sdramc_hw_init();

	reg =  readl(AT91C_BASE_CCFG + CCFG_EBI0CSA);
	reg |= AT91C_VDDIOM_SEL_33V;
	reg |= AT91C_EBI_CS1A_SDRAMC;
	writel(reg, AT91C_BASE_CCFG + CCFG_EBI0CSA);

	sdramc_initialize(&sdramc_config, AT91C_BASE_EBI0_CS1);
}
#endif /* #ifdef CONFIG_SDRAM */

#ifdef CONFIG_PSRAM

#define CONFIG_SYS_PSRAM_DATA_ACCESS_PIN	AT91C_PIN_PE(16)
static void data_access_enable(void)
{
	pio_set_value(CONFIG_SYS_PSRAM_DATA_ACCESS_PIN, 0);
}

static void psram_hw_init(void)
{
	unsigned short *addressMax = (unsigned short *)MICRON_8MB_ADDRESS_MAX;

	const struct pio_desc psram_pins[] = {
		{"CRE", CONFIG_SYS_PSRAM_DATA_ACCESS_PIN, 1, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure SMC1 CS0 */
	writel((AT91C_SMC_NWESETUP_(0)
		| AT91C_SMC_NCS_WRSETUP_(0)
		| AT91C_SMC_NRDSETUP_(0)
		| AT91C_SMC_NCS_RDSETUP_(0)),
		AT91C_BASE_SMC1 + SMC_SETUP0);

	writel((AT91C_SMC_NWEPULSE_(4)
		| AT91C_SMC_NCS_WRPULSE_(5)
		| AT91C_SMC_NRDPULSE_(2)
		| AT91C_SMC_NCS_RDPULSE_(5)),
		AT91C_BASE_SMC1 + SMC_PULSE0);

	writel((AT91C_SMC_NWECYCLE_(5)
		|  AT91C_SMC_NRDCYCLE_(7)),
		AT91C_BASE_SMC1 + SMC_CYCLE0);

	writel((AT91C_SMC_READMODE
		| AT91C_SMC_WRITEMODE
		| AT91C_SMC_BAT_BYTE_SELECT
		| AT91C_SMC_DBW_WIDTH_BITS_16
		| AT91C_SMC_PMEN
		| AT91C_SMC_PS_SIZE_THIRTY_TWO_BYTES),
		AT91C_BASE_SMC1 + SMC_CTRL0);

	/* Configure psram pins */
	pio_configure(psram_pins);

	data_access_enable();

	/* Enable page mode */
	readl(addressMax);
	readl(addressMax);
	writel(MICRON_RCR, addressMax);
	writel(MICRON_PAGE_MODE_ENABLE, addressMax);
}
#endif /* #ifdef CONFIG_PSRAM */

#if defined(CONFIG_NANDFLASH_RECOVERY) || defined(CONFIG_DATAFLASH_RECOVERY)
static void recovery_buttons_hw_init(void)
{
	/* Configure recovery button PINs */
	const struct pio_desc recovery_button_pins[] = {
		{"RECOVERY_BUTTON", CONFIG_SYS_RECOVERY_BUTTON_PIN, 0, PIO_PULLUP, PIO_INPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOCDE, PMC_PERIPH_CLK_DIVIDER_NA);
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
	 * PCK = MCK = MOSC
	 */
	/* Configure PLLA = MOSC * (PLL_MULA + 1) / PLL_DIVA */
	pmc_cfg_plla(PLLA_SETTINGS);

	/* PCK = PLLA = 2 * MCK */
	pmc_mck_cfg_set(0, MCKR_SETTINGS, AT91C_PMC_PRES | AT91C_PMC_MDIV);

	/* Switch MCK on PLLA output */
	pmc_mck_cfg_set(0, MCKR_CSS_SETTINGS,
			AT91C_PMC_PRES | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Enable External Reset */
	writel(AT91C_RSTC_KEY_UNLOCK | AT91C_RSTC_URSTEN, AT91C_BASE_RSTC + RSTC_RMR);

	at91_matrix_hw_init();

	/* Init timer */
	timer_init();

	/* Initialize dbgu */
	initialize_dbgu();

#ifdef CONFIG_SDRAM
	/* Initialize SDRAMC0 */
	sdramc0_init();
#endif

#if defined(CONFIG_PSRAM)
	psram_hw_init();
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
		{"MISO",	AT91C_PIN_PA(0),	0, PIO_DEFAULT, PIO_PERIPH_B},
		{"MOSI",	AT91C_PIN_PA(1),	0, PIO_DEFAULT, PIO_PERIPH_B},
		{"SPCK",	AT91C_PIN_PA(2),	0, PIO_DEFAULT, PIO_PERIPH_B},
		{"NPCS",	CONFIG_SYS_SPI_PCS,	1, PIO_PULLUP, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the spi0 pins */
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(spi0_pins);

	/* Enable the spi0 clock */
	pmc_enable_periph_clock(AT91C_ID_SPI0, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* CONFIG_DATAFLASH */

#ifdef CONFIG_SDCARD
#ifdef CONFIG_OF_LIBFDT
void at91_board_set_dtb_name(char *of_name)
{
	strcpy(of_name, "at91sam9263ek.dtb");
}
#endif

void at91_mci0_hw_init(void)
{
	/* Note: For at91sam9263ek, using MCI1 instead of MCI0 as SD Card booting */
	/* Configure the MCI1 pins */
	const struct pio_desc mci_pins[] = {
		{"MCI1_CK",	AT91C_PIN_PA(6), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCI1_CDA",	AT91C_PIN_PA(7), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCI1_DA0",	AT91C_PIN_PA(8), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCI1_DA1",	AT91C_PIN_PA(9), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCI1_DA2",	AT91C_PIN_PA(10), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"MCI1_DA3",	AT91C_PIN_PA(11), 0, PIO_PULLUP, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},

	};

	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(mci_pins);

	/* Enable the clock */
	pmc_enable_periph_clock(AT91C_ID_MCI1, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_SDCARD */

#ifdef CONFIG_NANDFLASH
void nandflash_hw_init(void)
{
	unsigned int reg;

	/* Configure NANDFlash pins*/
	const struct pio_desc nand_pins[] = {
		{"NANDCS",	CONFIG_SYS_NAND_ENABLE_PIN,	1, PIO_PULLUP, PIO_OUTPUT},
		{(char *)0, 	0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	/* Configure the NANDFlash pins */
	pio_configure(nand_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOCDE, PMC_PERIPH_CLK_DIVIDER_NA);

	/* Setup Smart Media, first enable the address range of CS3 in HMATRIX user interface  */
	reg = readl(AT91C_BASE_CCFG + CCFG_EBI0CSA);
	reg |= AT91C_EBI_CS3A_SM;
	writel(reg, AT91C_BASE_CCFG + CCFG_EBI0CSA);

	/* Configure SMC CS3 */
	writel((AT91C_SMC_NWESETUP_(1)
		| AT91C_SMC_NCS_WRSETUP_(1)
		| AT91C_SMC_NRDSETUP_(1)
		| AT91C_SMC_NCS_RDSETUP_(1)),
		AT91C_BASE_SMC0 + SMC_SETUP3);

	writel((AT91C_SMC_NWEPULSE_(3)
		| AT91C_SMC_NCS_WRPULSE_(3)
		| AT91C_SMC_NRDPULSE_(3)
		| AT91C_SMC_NCS_RDPULSE_(3)),
		AT91C_BASE_SMC0 + SMC_PULSE3);

	writel((AT91C_SMC_NWECYCLE_(5)
		|  AT91C_SMC_NRDCYCLE_(5)),
		AT91C_BASE_SMC0 + SMC_CYCLE3);

	writel((AT91C_SMC_READMODE
		| AT91C_SMC_WRITEMODE
		| AT91C_SMC_NWAITM_NWAIT_DISABLE
		| AT91C_SMC_DBW_WIDTH_BITS_8
		| AT91_SMC_TDF_(2)),
		AT91C_BASE_SMC0 + SMC_CTRL3);
}
#endif /* #ifdef CONFIG_NANDFLASH */
