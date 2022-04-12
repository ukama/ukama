/* ----------------------------------------------------------------------------
 *         MICROCHIP/ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2017, Microchip Technology Inc.
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
#include "debug.h"
#include "ddramc.h"
#include "gpio.h"
#include "hardware.h"
#include "l2cc.h"
#include "matrix.h"
#include "pmc.h"
#include "timer.h"
#include "usart.h"
#include "watchdog.h"

#include "arch/at91_ddrsdrc.h"
#include "arch/at91_pio.h"
#include "arch/at91_pmc/pmc.h"
#include "arch/at91_rstc.h"
#include "arch/at91_sfr.h"
#include "arch/sama5_smc.h"
#include "arch/tz_matrix.h"
#include "sama5d2_ptc_ek.h"

static void at91_dbgu_hw_init(void)
{
	const struct pio_desc dbgu_pins[] = {
		{"RXD0", CONFIG_SYS_DBGU_RXD_PIN, 0, PIO_DEFAULT, PIO_PERIPH_C},
		{"TXD0", CONFIG_SYS_DBGU_TXD_PIN, 0, PIO_DEFAULT, PIO_PERIPH_C},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(dbgu_pins);
	pmc_enable_periph_clock(CONFIG_SYS_DBGU_ID, PMC_PERIPH_CLK_DIVIDER_NA);
}

static void initialize_dbgu(void)
{
	unsigned int baudrate = 115200;

	at91_dbgu_hw_init();

	if (pmc_mck_check_h32mxdiv())
		usart_init(BAUDRATE(MASTER_CLOCK / 2, baudrate));
	else
		usart_init(BAUDRATE(MASTER_CLOCK, baudrate));
}

#ifdef CONFIG_MATRIX
static int matrix_configure_slave(void)
{
	unsigned int ddr_port;
	unsigned int ssr_setting, sasplit_setting, srtop_setting;

	/*
	 * Matrix 0 (H64MX)
	 */

	/*
	 * 0: Bridge from H64MX to AXIMX
	 * (Internal ROM, Crypto Library, PKCC RAM): Always Secured
	 */

	/* 1: H64MX Peripheral Bridge */

	/* 2 ~ 9 DDR2 Port0 ~ 7: Non-Secure */
	srtop_setting = MATRIX_SRTOP(0, MATRIX_SRTOP_VALUE_128M);
	sasplit_setting = (MATRIX_SASPLIT(0, MATRIX_SASPLIT_VALUE_128M) |
			   MATRIX_SASPLIT(1, MATRIX_SASPLIT_VALUE_128M) |
			   MATRIX_SASPLIT(2, MATRIX_SASPLIT_VALUE_128M) |
			   MATRIX_SASPLIT(3, MATRIX_SASPLIT_VALUE_128M));
	ssr_setting = (MATRIX_LANSECH_NS(0) |
		       MATRIX_LANSECH_NS(1) |
		       MATRIX_LANSECH_NS(2) |
		       MATRIX_LANSECH_NS(3) |
		       MATRIX_RDNSECH_NS(0) |
		       MATRIX_RDNSECH_NS(1) |
		       MATRIX_RDNSECH_NS(2) |
		       MATRIX_RDNSECH_NS(3) |
		       MATRIX_WRNSECH_NS(0) |
		       MATRIX_WRNSECH_NS(1) |
		       MATRIX_WRNSECH_NS(2) |
		       MATRIX_WRNSECH_NS(3));
	for (ddr_port = 0; ddr_port < 8; ddr_port++) {
		matrix_configure_slave_security(AT91C_BASE_MATRIX64,
					(H64MX_SLAVE_DDR2_PORT_0 + ddr_port),
					srtop_setting,
					sasplit_setting,
					ssr_setting);
	}

	/*
	 * 10: Internal SRAM 128K
	 * TOP0 is set to 128K
	 * SPLIT0 is set to 64K
	 * LANSECH0 is set to 0, the low area of region 0 is the Securable one
	 * RDNSECH0 is set to 0, region 0 Securable area is secured for reads.
	 * WRNSECH0 is set to 0, region 0 Securable area is secured for writes
	 */
	srtop_setting = MATRIX_SRTOP(0, MATRIX_SRTOP_VALUE_128K);
	sasplit_setting = MATRIX_SASPLIT(0, MATRIX_SASPLIT_VALUE_64K);
	ssr_setting = (MATRIX_LANSECH_S(0) |
		       MATRIX_RDNSECH_S(0) |
		       MATRIX_WRNSECH_S(0));
	matrix_configure_slave_security(AT91C_BASE_MATRIX64,
					H64MX_SLAVE_INTERNAL_SRAM,
					srtop_setting,
					sasplit_setting,
					ssr_setting);

	/* 11:  Internal SRAM 128K (Cache L2) */
	/* 12:  QSPI0 */
	/* 13:  QSPI1 */
	/* 14:  AESB */

	/*
	 * Matrix 1 (H32MX)
	 */

	/* 0: Bridge from H32MX to H64MX: Not Secured */

	/* 1: H32MX Peripheral Bridge 0: Not Secured */

	/* 2: H32MX Peripheral Bridge 1: Not Secured */

	/*
	 * 3: External Bus Interface
	 * EBI CS0 Memory(256M) ----> Slave Region 0, 1
	 * EBI CS1 Memory(256M) ----> Slave Region 2, 3
	 * EBI CS2 Memory(256M) ----> Slave Region 4, 5
	 * EBI CS3 Memory(128M) ----> Slave Region 6
	 * NFC Command Registers(128M) -->Slave Region 7
	 *
	 * NANDFlash(EBI CS3) --> Slave Region 6: Non-Secure
	 */
	srtop_setting =	MATRIX_SRTOP(6, MATRIX_SRTOP_VALUE_128M);
	srtop_setting |= MATRIX_SRTOP(7, MATRIX_SRTOP_VALUE_128M);
	sasplit_setting = MATRIX_SASPLIT(6, MATRIX_SASPLIT_VALUE_128M);
	sasplit_setting |= MATRIX_SASPLIT(7, MATRIX_SASPLIT_VALUE_128M);
	ssr_setting = (MATRIX_LANSECH_NS(6) |
		       MATRIX_RDNSECH_NS(6) |
		       MATRIX_WRNSECH_NS(6));
	ssr_setting |= (MATRIX_LANSECH_NS(7) |
			MATRIX_RDNSECH_NS(7) |
			MATRIX_WRNSECH_NS(7));
	matrix_configure_slave_security(AT91C_BASE_MATRIX32,
					H32MX_EXTERNAL_EBI,
					srtop_setting,
					sasplit_setting,
					ssr_setting);

	/* 4: NFC SRAM (4K): Non-Secure */
	srtop_setting = MATRIX_SRTOP(0, MATRIX_SRTOP_VALUE_8K);
	sasplit_setting = MATRIX_SASPLIT(0, MATRIX_SASPLIT_VALUE_8K);
	ssr_setting = (MATRIX_LANSECH_NS(0) |
		       MATRIX_RDNSECH_NS(0) |
		       MATRIX_WRNSECH_NS(0));
	matrix_configure_slave_security(AT91C_BASE_MATRIX32,
					H32MX_NFC_SRAM,
					srtop_setting,
					sasplit_setting,
					ssr_setting);

	/* 5:
	 * USB Device High Speed Dual Port RAM (DPR): 1M
	 * USB Host OHCI registers: 1M
	 * USB Host EHCI registers: 1M
	 */
	srtop_setting = (MATRIX_SRTOP(0, MATRIX_SRTOP_VALUE_1M) |
			 MATRIX_SRTOP(1, MATRIX_SRTOP_VALUE_1M) |
			 MATRIX_SRTOP(2, MATRIX_SRTOP_VALUE_1M));
	sasplit_setting = (MATRIX_SASPLIT(0, MATRIX_SASPLIT_VALUE_1M) |
			   MATRIX_SASPLIT(1, MATRIX_SASPLIT_VALUE_1M) |
			   MATRIX_SASPLIT(2, MATRIX_SASPLIT_VALUE_1M));
	ssr_setting = (MATRIX_LANSECH_NS(0) |
		       MATRIX_LANSECH_NS(1) |
		       MATRIX_LANSECH_NS(2) |
		       MATRIX_RDNSECH_NS(0) |
		       MATRIX_RDNSECH_NS(1) |
		       MATRIX_RDNSECH_NS(2) |
		       MATRIX_WRNSECH_NS(0) |
		       MATRIX_WRNSECH_NS(1) |
		       MATRIX_WRNSECH_NS(2));
	matrix_configure_slave_security(AT91C_BASE_MATRIX32,
					H32MX_USB,
					srtop_setting,
					sasplit_setting,
					ssr_setting);

	return 0;
}

static unsigned int security_ps_peri_id[] = {
	0,
};

static int matrix_config_periheral(void)
{
	unsigned int *peri_id = security_ps_peri_id;
	unsigned int array_size = sizeof(security_ps_peri_id) / sizeof(unsigned int);
	int ret;

	ret = matrix_configure_peri_security(peri_id, array_size);
	if (ret)
		return -1;

	return 0;
}

static int matrix_init(void)
{
	int ret;

	matrix_write_protect_disable(AT91C_BASE_MATRIX64);
	matrix_write_protect_disable(AT91C_BASE_MATRIX32);

	ret = matrix_configure_slave();
	if (ret)
		return -1;

	ret = matrix_config_periheral();
	if (ret)
		return -1;

	return 0;
}
#endif	/* CONFIG_MATRIX */

#ifdef CONFIG_DDR2
static void ddramc_reg_config(struct ddramc_register *ddramc_config)
{
	ddramc_config->mdr = (AT91C_DDRC2_DBW_32_BITS |
			      AT91C_DDRC2_MD_DDR2_SDRAM);

	ddramc_config->cr = (AT91C_DDRC2_NC_DDR10_SDR9 |
			     AT91C_DDRC2_NR_14 |
			     AT91C_DDRC2_CAS_3 |
			     AT91C_DDRC2_DISABLE_RESET_DLL |
			     AT91C_DDRC2_DISABLE_DLL |
			     AT91C_DDRC2_NB_BANKS_8 |
			     AT91C_DDRC2_DECOD_INTERLEAVED |
			     AT91C_DDRC2_UNAL_SUPPORTED);

	ddramc_config->rtr = 0x500;

	/* Assume timings for 6ns min clock period */
	ddramc_config->t0pr = (AT91C_DDRC2_TRAS_(8) |
			       AT91C_DDRC2_TRCD_(3) |
			       AT91C_DDRC2_TWR_(3) |
			       AT91C_DDRC2_TRC_(10) |
			       AT91C_DDRC2_TRP_(3) |
			       AT91C_DDRC2_TRRD_(4) |
			       AT91C_DDRC2_TWTR_(4) |
			       AT91C_DDRC2_TMRD_(2));

	ddramc_config->t1pr = (AT91C_DDRC2_TRFC_(33) |
			       AT91C_DDRC2_TXSNR_(35) |
			       AT91C_DDRC2_TXSRD_(200) |
			       AT91C_DDRC2_TXP_(3));

	ddramc_config->t2pr = (AT91C_DDRC2_TXARD_(2) |
			       AT91C_DDRC2_TXARDS_(8) |
			       AT91C_DDRC2_TRPA_(4) |
			       AT91C_DDRC2_TRTP_(2) |
			       AT91C_DDRC2_TFAW_(8));
}

static void ddramc_init(void)
{
	struct ddramc_register ddramc_reg;
	unsigned int reg;

	ddramc_reg_config(&ddramc_reg);

	pmc_enable_periph_clock(AT91C_ID_MPDDRC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_system_clock(AT91C_PMC_DDR);

	/* Configure Shift Sampling Point of Data. */
	reg = AT91C_MPDDRC_RD_DATA_PATH_ONE_CYCLES;
	writel(reg, (AT91C_BASE_MPDDRC + MPDDRC_RD_DATA_PATH));

	/* MPDDRC I/O Calibration Register */
	reg = readl(AT91C_BASE_MPDDRC + MPDDRC_IO_CALIBR);
	reg &= ~AT91C_MPDDRC_RDIV;
	reg |= AT91C_MPDDRC_RDIV_DDR2_RZQ_50;
	reg &= ~AT91C_MPDDRC_TZQIO;
	reg |= AT91C_MPDDRC_TZQIO_(100);
	writel(reg, (AT91C_BASE_MPDDRC + MPDDRC_IO_CALIBR));

	ddram_initialize(AT91C_BASE_MPDDRC, AT91C_BASE_DDRCS, &ddramc_reg);

	ddramc_dump_regs(AT91C_BASE_MPDDRC);
}
#endif /* CONFIG DDR2 */

#ifdef CONFIG_HW_INIT
void hw_init(void)
{
	/* Disable watchdog */
	at91_disable_wdt();

	/* Configure PLLA */
	pmc_cfg_plla(PLLA_SETTINGS);

	/* Initialize PLLA charge pump */
	/*
	 * The field named ICP_PLLA[1:0] must be written to 0.
	 * Even if its default value is 0, it is wrongly re-written to 0x3
	 * by the ROMCode.
	 */
	pmc_init_pll(AT91C_PMC_ICPPLLA_0);

	/* Switch MCK clock source to PLLA */
	pmc_mck_cfg_set(0, BOARD_PRESCALER_PLLA,
			AT91C_PMC_H32MXDIV | AT91C_PMC_PLLADIV2 |
			AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Enable external reset */
	writel(AT91C_RSTC_KEY_UNLOCK | AT91C_RSTC_URSTEN,
	       AT91C_BASE_RSTC + RSTC_RMR);

#if defined(CONFIG_MATRIX)
	/* Initialize the matrix */
	matrix_init();
#endif
	/* initialize the dbgu */
	initialize_dbgu();

	/* Init timer */
	timer_init();

#if defined(CONFIG_DDR2)
	/* Initialize MPDDR Controller */
	ddramc_init();
#endif
	/* Prepare L2 cache setup */
	l2cache_prepare();
}
#endif /* #ifdef CONFIG_HW_INIT */

#ifdef CONFIG_DATAFLASH
void at91_spi0_hw_init(void)
{
	const struct pio_desc spi_pins[] = {
		{"SPI0_SPCK",	AT91C_PIN_PA(14), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SPI0_MOSI",	AT91C_PIN_PA(15), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SPI0_MISO",	AT91C_PIN_PA(16), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SPI0_NPCS",	CONFIG_SYS_SPI_PCS, 1, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(spi_pins);
	pmc_enable_periph_clock(CONFIG_SYS_ID_SPI, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif

#ifdef CONFIG_NANDFLASH
void nandflash_hw_init(void)
{
	const struct pio_desc nand_pins[] = {
		{"NANDOE", CONFIG_SYS_NAND_OE_PIN, 0, PIO_PULLUP, PIO_PERIPH_B},
		{"NANDWE", CONFIG_SYS_NAND_WE_PIN, 0, PIO_PULLUP, PIO_PERIPH_B},
		{"NANDALE", CONFIG_SYS_NAND_ALE_PIN, 0, PIO_PULLUP, PIO_PERIPH_B},
		{"NANDCLE", CONFIG_SYS_NAND_CLE_PIN, 0, PIO_PULLUP, PIO_PERIPH_B},
		{"NANDCS", CONFIG_SYS_NAND_ENABLE_PIN, 1, PIO_DEFAULT, PIO_OUTPUT},
		{"D0", AT91C_PIN_PA(22), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D1", AT91C_PIN_PA(23), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D2", AT91C_PIN_PA(24), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D3", AT91C_PIN_PA(25), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D4", AT91C_PIN_PA(26), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D5", AT91C_PIN_PA(27), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D6", AT91C_PIN_PA(28), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{"D7", AT91C_PIN_PA(29), 0, PIO_DRVSTR_ME, PIO_PERIPH_B},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(nand_pins);
	pmc_enable_periph_clock(AT91C_ID_HSMC, PMC_PERIPH_CLK_DIVIDER_NA);

	/* EBI Configuration Register */
	writel((AT91C_EBICFG_DRIVE0_HIGH |
		AT91C_EBICFG_PULL0_NONE |
		AT91C_EBICFG_DRIVE1_HIGH |
		AT91C_EBICFG_PULL1_NONE), SFR_EBICFG + AT91C_BASE_SFR);

	/* Configure SMC CS3 for NAND/SmartMedia */
	writel(AT91C_SMC_SETUP_NWE(1) |
	       AT91C_SMC_SETUP_NCS_WR(1) |
	       AT91C_SMC_SETUP_NRD(1) |
	       AT91C_SMC_SETUP_NCS_RD(1), (ATMEL_BASE_SMC + SMC_SETUP3));

	writel(AT91C_SMC_PULSE_NWE(2) |
	       AT91C_SMC_PULSE_NCS_WR(3) |
	       AT91C_SMC_PULSE_NRD(2) |
	       AT91C_SMC_PULSE_NCS_RD(3), (ATMEL_BASE_SMC + SMC_PULSE3));

	writel(AT91C_SMC_CYCLE_NWE(5) |
	       AT91C_SMC_CYCLE_NRD(5), (ATMEL_BASE_SMC + SMC_CYCLE3));

	writel(AT91C_SMC_TIMINGS_TCLR(2) |
	       AT91C_SMC_TIMINGS_TADL(7) |
	       AT91C_SMC_TIMINGS_TAR(2) |
	       AT91C_SMC_TIMINGS_TRR(3) |
	       AT91C_SMC_TIMINGS_TWB(7) |
	       AT91C_SMC_TIMINGS_RBNSEL(2) |
	       AT91C_SMC_TIMINGS_NFSEL, (ATMEL_BASE_SMC + SMC_TIMINGS3));

	writel(AT91C_SMC_MODE_READMODE_NRD_CTRL |
	       AT91C_SMC_MODE_WRITEMODE_NWE_CTRL |
	       AT91C_SMC_MODE_DBW_8 |
	       AT91C_SMC_MODE_TDF_CYCLES(1), (ATMEL_BASE_SMC + SMC_MODE3));
}
#endif /* CONFIG_NANDFLASH */

#ifdef CONFIG_SDCARD
#define ATMEL_SDHC_GCKDIV_VALUE	1

void at91_sdhc_hw_init(void)
{
#ifdef CONFIG_SDHC0
	const struct pio_desc sdmmc_pins[] = {
		{"SDMMC0_CK",   AT91C_PIN_PA(0), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_CMD",  AT91C_PIN_PA(1), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT0", AT91C_PIN_PA(2), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT1", AT91C_PIN_PA(3), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT2", AT91C_PIN_PA(4), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT3", AT91C_PIN_PA(5), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT4", AT91C_PIN_PA(6), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT5", AT91C_PIN_PA(7), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT6", AT91C_PIN_PA(8), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT7", AT91C_PIN_PA(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_RSTN", AT91C_PIN_PA(10), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_VDDSEL", AT91C_PIN_PA(11), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_CD",   AT91C_PIN_PA(13), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
#endif

#ifdef CONFIG_SDHC1
	const struct pio_desc sdmmc_pins[] = {
		{"SDMMC1_CD",   AT91C_PIN_PA(30), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_CMD",  AT91C_PIN_PA(28), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_CK",   AT91C_PIN_PA(22), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_DAT0", AT91C_PIN_PA(18), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_DAT1", AT91C_PIN_PA(19), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_DAT2", AT91C_PIN_PA(20), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{"SDMMC1_DAT3", AT91C_PIN_PA(21), 0, PIO_DEFAULT, PIO_PERIPH_E},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
#endif

	pio_configure(sdmmc_pins);

	pmc_enable_periph_clock(CONFIG_SYS_ID_SDHC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_generic_clock(CONFIG_SYS_ID_SDHC,
				 GCK_CSS_UPLL_CLK,
				 ATMEL_SDHC_GCKDIV_VALUE);
}
#endif /* CONFIG_SDCARD */
