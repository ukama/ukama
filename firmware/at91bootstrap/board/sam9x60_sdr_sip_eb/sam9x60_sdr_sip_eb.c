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

#include "common.h"
#include "hardware.h"
#include "arch/at91_sfr.h"
#include "arch/at91_rstc.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_smc.h"
#include "arch/at91_pio.h"
#include "arch/at91_sdramc.h"
#include "gpio.h"
#include "pmc.h"
#include "usart.h"
#include "debug.h"
#include "sdramc.h"
#include "timer.h"
#include "watchdog.h"
#include "string.h"
#include "sam9x60_sdr_sip_eb.h"
#include "board_hw_info.h"
#include "twi.h"
#include "flexcom.h"
#include "board.h"

#define PLLA_DIV 1
#define PLLA_COUNT 0x3f
#define PLLA_CLOCK 200000000
#define PLLA_FRACR(_p, _q) \
	((unsigned int)((((unsigned long)(_p)) << 22) / (_q)))

static void at91_dbgu_hw_init(void)
{
	/* Configure DBGU pins */
	const struct pio_desc dbgu_pins[] = {
		{"RXD", AT91C_PIN_PA(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TXD", AT91C_PIN_PA(10), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(dbgu_pins);

	pmc_enable_periph_clock(AT91C_ID_DBGU, PMC_PERIPH_CLK_DIVIDER_NA);
}

static void initialize_dbgu(void)
{
	at91_dbgu_hw_init();
	usart_init(BAUDRATE(MASTER_CLOCK, BAUD_RATE));
}

#ifdef CONFIG_SDRAM

static void sdramc_reg_config(struct sdramc_register *sdramc_config)
{
	sdramc_config->cr = AT91C_SDRAMC_NC_8
	| AT91C_SDRAMC_NR_12
	| AT91C_SDRAMC_CAS_3
	| AT91C_SDRAMC_NB_4_BANKS
	| AT91C_SDRAMC_DBW_16_BITS
	| AT91C_SDRAMC_TWR_3
	| AT91C_SDRAMC_TRC_11
	| AT91C_SDRAMC_TRP_3
	| AT91C_SDRAMC_TRCD_3
	| AT91C_SDRAMC_TRAS_8
	| AT91C_SDRAMC_TXSR_12;

        sdramc_config->tr = (MASTER_CLOCK * 7) / 1000000;
        sdramc_config->mdr = AT91C_SDRAMC_MD_SDRAM
			 | AT91C_SDRAMC_MD_SHIFT_SAMPLING_2_CYCLE;

	sdramc_config->cfr1 = AT91C_SDRAMC_CFR1_TMRD_2 | AT91C_SDRAMC_CFR1_UNAL
				| AT91C_SDRAMC_CFR1_ADD_DATA_MUX_UNSUPPORTED
				| AT91C_SDRAMC_CFR1_CMD_MUX_UNSUPPORTED;
}

static void sdramc_init(void)
{
	struct sdramc_register sdramc_reg;
	unsigned int reg = 0;

	reg = readl(AT91C_BASE_SFR + SFR_DDRCFG);
	/*
	 * We need to also enable AT91C_EBI_NFD0_ON_D16 . Otherwise the DDR will
	 * not work if NAND lines have been previously used by RomCode
	 */
	reg |= (AT91C_EBI_CS1A | AT91C_EBI_NFD0_ON_D16);
	writel(reg, (AT91C_BASE_SFR + SFR_DDRCFG));

	pmc_enable_periph_clock(AT91C_ID_SDRAMC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_system_clock(AT91C_PMC_DDR);

	sdramc_reg_config(&sdramc_reg);

	sdramc_initialize(&sdramc_reg, AT91C_BASE_CS1);
}
#endif /* CONFIG SDRAM */

static void at91_green_led_on(void)
{
	pio_set_gpio_output(AT91C_PIN_PB(13), 0);
	pio_set_gpio_output(AT91C_PIN_PB(11), 0);
	pio_set_gpio_output(AT91C_PIN_PB(12), 1);
}

#if defined(CONFIG_FLEXCOM)
unsigned int at91_flexcom0_init(void)
{
	const struct pio_desc flx_pins[] = {
		{"FLX_IO0", AT91C_PIN_PA(0), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"FLX_IO1", AT91C_PIN_PA(1), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(flx_pins);
	pmc_enable_periph_clock(AT91C_ID_FLEXCOM0, PMC_PERIPH_CLK_DIVIDER_NA);

	flexcom_init(0);

	return flexcom_get_regmap(0);
}
#endif

#if defined CONFIG_TWI
void twi_init()
{
#if defined(CONFIG_FLEXCOM)
       twi_bus_init(at91_flexcom0_init);
#endif
#if defined(CONFIG_AUTOCONFIG_TWI_BUS)
	dbg_loud("Auto-Config the TWI Bus by the board\n");
	at91_board_config_twi_bus();
#endif
}
#endif

#if defined(CONFIG_AUTOCONFIG_TWI_BUS)
void at91_board_config_twi_bus(void)
{
	at24xx_twi_bus = 0;
}
#endif
#if defined CONFIG_FLEXCOM
static struct at91_flexcom flexcoms[] = {
#if defined CONFIG_TWI
	{AT91C_ID_FLEXCOM0, FLEXCOM_TWI, AT91C_BASE_FLEXCOM0, 0},
#endif
};

void board_flexcoms_init()
{
	flexcoms_init(flexcoms);
}
#endif

#ifdef CONFIG_HW_INIT
void hw_init(void)
{
	struct pmc_pll_cfg plla_config;

	/* Disable watchdog */
	at91_disable_wdt();

	/* Green Led ON */
	at91_green_led_on();

	/* Configure & Enable PLLA */
	plla_config.mul = 49;
	plla_config.div = PLLA_DIV;
	plla_config.count = PLLA_COUNT;
	plla_config.fracr = 0;
	plla_config.acr = AT91C_PLL_ACR_DEFAULT_PLLA;
	pmc_sam9x60_cfg_pll(PLL_ID_PLLA, &plla_config);

	pmc_mck_cfg_set(0, BOARD_PRESCALER_PLLA,
			AT91C_PMC_PRES | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Initialize dbgu */
	initialize_dbgu();

	/* Enable External Reset */
	writel(AT91C_RSTC_KEY_UNLOCK | AT91C_RSTC_URSTEN, AT91C_BASE_RSTC + RSTC_RMR);

	/* Init timer */
	timer_init();

#ifdef CONFIG_FLEXCOM
	board_flexcoms_init();
#endif

#ifdef CONFIG_TWI
	twi_init();
#endif

#ifdef CONFIG_SDRAM
	/* Initialize SDRAM Controller */
	sdramc_init();
#endif
}
#endif /* #ifdef CONFIG_HW_INIT */

#ifdef CONFIG_DATAFLASH

#if defined(CONFIG_QSPI)
void at91_qspi_hw_init(void)
{
	const struct pio_desc qspi_pins[] = {
		{"QSPI0_SCK", AT91C_PIN_PB(19), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"QSPI0_CS", AT91C_PIN_PB(20), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"QSPI0_IO0", AT91C_PIN_PB(21), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"QSPI0_IO1", AT91C_PIN_PB(22), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"QSPI0_IO2", AT91C_PIN_PB(23), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"QSPI0_IO3", AT91C_PIN_PB(24), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(qspi_pins);

	pmc_enable_system_clock(AT91C_PMC_QSPICLK);
	pmc_enable_periph_clock(CONFIG_SYS_ID_QSPI, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif  /* #ifdef CONFIG_QSPI */

#endif	/* #ifdef CONFIG_DATAFLASH */

#ifdef CONFIG_SDCARD
#ifdef CONFIG_OF_LIBFDT
void at91_board_set_dtb_name(char *of_name)
{
	strcpy(of_name, "at91-sam9x60ek.dtb");
}
#endif

#define ATMEL_SDHC_GCKDIV_VALUE     5

void at91_sdhc_hw_init(void)
{
#ifdef CONFIG_SDHC0
	const struct pio_desc sdmmc_pins[] = {
		{"SDMMC0_CMD",	AT91C_PIN_PA(16), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_CK",	AT91C_PIN_PA(17), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT0",	AT91C_PIN_PA(15), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT1",	AT91C_PIN_PA(18), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT2",	AT91C_PIN_PA(19), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT3",	AT91C_PIN_PA(20), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
	pio_configure(sdmmc_pins);

	pmc_enable_periph_clock(CONFIG_SYS_ID_SDHC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_generic_clock(CONFIG_SYS_ID_SDHC,
				 GCK_CSS_PLLA_CLK,
				 ATMEL_SDHC_GCKDIV_VALUE);
#endif
}
#endif /* #ifdef CONFIG_SDCARD */

#ifdef CONFIG_NANDFLASH
void nandflash_hw_init(void)
{
	unsigned int reg;
	const struct pio_desc nand_pins[] = {
		{"NANDOE", CONFIG_SYS_NAND_OE_PIN, 0, PIO_DEFAULT,
		 PIO_PERIPH_A},
		{"NANDWE", CONFIG_SYS_NAND_WE_PIN, 0, PIO_DEFAULT,
		 PIO_PERIPH_A},
		{"NANDALE", CONFIG_SYS_NAND_ALE_PIN, 0, PIO_DEFAULT,
		 PIO_PERIPH_A},
		{"NANDCLE", CONFIG_SYS_NAND_CLE_PIN, 0, PIO_DEFAULT,
		 PIO_PERIPH_A},
		{"NANDCS", CONFIG_SYS_NAND_ENABLE_PIN, 1, PIO_DEFAULT,
		 PIO_OUTPUT},
		{"D0", AT91C_PIN_PD(6), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D1", AT91C_PIN_PD(7), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D2", AT91C_PIN_PD(8), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D3", AT91C_PIN_PD(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D4", AT91C_PIN_PD(10), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D5", AT91C_PIN_PD(11), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D6", AT91C_PIN_PD(12), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"D7", AT91C_PIN_PD(13), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(nand_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOD, PMC_PERIPH_CLK_DIVIDER_NA);

	reg = readl(AT91C_BASE_SFR + SFR_CCFG_EBICSA);
	reg |= AT91C_EBI_CS3A_SM | AT91C_EBI_NFD0_ON_D16;
	reg &= ~AT91C_EBI_DRV;
	writel(reg, AT91C_BASE_SFR + SFR_CCFG_EBICSA);

	/* Configure SMC CS3 for NAND */
	writel(AT91C_SMC_NWESETUP_(4), AT91C_BASE_SMC + SMC_SETUP3);

	writel(AT91C_SMC_NWEPULSE_(10) | AT91C_SMC_NCS_WRPULSE_(20) |
	       AT91C_SMC_NRDPULSE_(10) | AT91C_SMC_NCS_RDPULSE_(20),
	       AT91C_BASE_SMC + SMC_PULSE3);

	writel(AT91C_SMC_NWECYCLE_(20) | AT91C_SMC_NRDCYCLE_(20),
	       AT91C_BASE_SMC + SMC_CYCLE3);

	writel(AT91C_SMC_READMODE | AT91C_SMC_WRITEMODE | AT91C_SMC_TDFEN |
	       AT91_SMC_TDF_(15), AT91C_BASE_SMC + SMC_CTRL3);
}
#endif /* #ifdef CONFIG_NANDFLASH */
