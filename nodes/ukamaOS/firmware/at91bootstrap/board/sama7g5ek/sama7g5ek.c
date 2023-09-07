/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Microchip Technology Inc. and its subsidiaries
 *
 * All rights reserved.
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
#include "flexcom.h"
#include "usart.h"
#include "hardware.h"
#include "arch/at91_pio.h"
#include "string.h"
#include "tzc400-interface.h"
#include "debug.h"
#include "umctl2.h"
#include "gpio.h"
#include "mcp16502.h"
#include "pmc.h"
#include "arch/at91_pmc/pmc.h"
#include "publ.h"
#include "umctl2.h"
#include "watchdog.h"
#include "timer.h"
#include "sdhc_cal.h"

#include "sama7g5ek.h"

static void ca7_enable_smp()
{
	/* asm code to enable ACTLR.SMP . this is needed for SCU for L2 cache */
	asm  volatile (
		"MRC p15, 0, R1, c1, c0, 1;"
		"ORR R1, R1, #64;"
		"MCR p15, 0, R1, c1, c0, 1;"
	);
}

static void initialize_serial(void)
{
	unsigned int baudrate = 115200;

	usart_init(BAUDRATE(MASTER_CLOCK, baudrate));
}

#if defined CONFIG_FLEXCOM
static struct at91_flexcom flexcoms[] = {
	{AT91C_ID_FLEXCOM0, FLEXCOM_USART, AT91C_BASE_FLEXCOM0}, /* BT serial */
	{AT91C_ID_FLEXCOM1, FLEXCOM_TWI, AT91C_BASE_FLEXCOM1}, /* PMIC I2C */
	{AT91C_ID_FLEXCOM2, FLEXCOM_USART, AT91C_BASE_FLEXCOM2}, /* Unused */
	{AT91C_ID_FLEXCOM3, FLEXCOM_USART, AT91C_BASE_FLEXCOM3}, /* DBGU */
};

void board_flexcoms_init()
{
	flexcoms_init(flexcoms);
}
#endif

#if defined(CONFIG_FLEXCOM)
unsigned int at91_flexcom3_init(void)
{
	const struct pio_desc flx_pins[] = {
		{"FLX_IO0", CONFIG_SYS_DBGU_TXD_PIN, 0, PIO_DEFAULT, PIO_PERIPH_F},
		{"FLX_IO1", CONFIG_SYS_DBGU_RXD_PIN, 0, PIO_PULLUP, PIO_PERIPH_F},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_F},
	};

	pio_configure(flx_pins);
	pmc_enable_periph_clock(AT91C_ID_FLEXCOM3, PMC_PERIPH_CLK_DIVIDER_NA);

	flexcom_init(3);

	return flexcom_get_regmap(3);
}

static unsigned int at91_flexcom1_init(void)
{
	const struct pio_desc flx_pins[] = {
		{"FLX_IO0", AT91C_PIN_PC(9), 0, PIO_DEFAULT, PIO_PERIPH_F},
		{"FLX_IO1", AT91C_PIN_PC(10), 0, PIO_DEFAULT, PIO_PERIPH_F},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_F},
	};

	pio_configure(flx_pins);
	pmc_enable_periph_clock(AT91C_ID_FLEXCOM1, PMC_PERIPH_CLK_DIVIDER_NA);

	flexcom_init(1);

	return flexcom_get_regmap(1);
}
#endif

static struct umctl2_config_state umctl2_config;

static void umctl2_config_state_init()
{
	/* enable 5 AXI ports */
	umctl2_config.axi_port_bitmap = MP_AXI_PORT_ENABLE(0) |
					MP_AXI_PORT_ENABLE(1) |
					MP_AXI_PORT_ENABLE(2) |
					MP_AXI_PORT_ENABLE(3) |
					MP_AXI_PORT_ENABLE(4);

	/* port 0 configuration : CPU port */
		/* region 0 is ArQOS 0x0 to 0x8; */
		/* region 1 is ArQOS 0x9 to 0xF; */
	umctl2_config.port_x_read_region0_last[0] = 0x3;
		/* region 2 doesn't exist (single queue port) */
	umctl2_config.port_x_read_region1_last[0] = 0xE;
		/* region 0 is AwQOS 0x0 to 0xF; */
	umctl2_config.port_x_write_region0_last[0] = 0x3;
	umctl2_config.port_x_write_region1_last[0] = 0xE;

	/* port 1 configuration : AHB port */
		/* region 0 is ArQOS 0x0 to 0x8; */
		/* region 1 is ArQOS 0x9 to 0xF; */
	umctl2_config.port_x_read_region0_last[1] = 0x3;
		/* region 2 doesn't exist (single queue port) */
	umctl2_config.port_x_read_region1_last[1] = 0xE;
		/* region 0 is AwQOS 0x0 to 0xF; */
	umctl2_config.port_x_write_region0_last[1] = 0x3;
	umctl2_config.port_x_write_region1_last[1] = 0xE;

	/* port 2 configuration : PSS port XDMAC0,1+... */
		/* region 0 is ArQOS 0x0 to 0x3; */
		/* region 1 is ArQOS 0x3 to 0x8; */
	umctl2_config.port_x_read_region0_last[2] = 0x3;
		/* region 2 is ArQOS 0x8 to 0xF */
	umctl2_config.port_x_read_region1_last[2] = 0x8;
		/* region 0 is AwQOS 0x0 to 0xF; */
	umctl2_config.port_x_write_region0_last[2] = 0x3;
	umctl2_config.port_x_write_region1_last[2] = 0xE;

	/* port 3 configuration : ISS port XISC */
		/* region 0 is ArQOS 0x0 to 0x8; */
		/* region 1 is ArQOS 0x9 to 0xF; */
	umctl2_config.port_x_read_region0_last[3] = 0x3;
		/* region 2 doesn't exist (single queue port) */
	umctl2_config.port_x_read_region1_last[3] = 0xE;
		/* region 0 is AwQOS 0x0 to 0xF; */
	umctl2_config.port_x_write_region0_last[3] = 0x3;
	umctl2_config.port_x_write_region1_last[3] = 0xE;

	umctl2_config.port_x_rdwr_ordered_en = 0x0;

	/* port 4 configuration : MSS port SD+GMAC+... */
		/* region 0 is ArQOS 0x0 to 0x8; */
		/* region 1 is ArQOS 0x9 to 0xF; */
	umctl2_config.port_x_read_region0_last[4] = 0x3;
		/* region 2 doesn't exist (single queue port) */
	umctl2_config.port_x_read_region1_last[4] = 0xE;
		/* region 0 is AwQOS 0x0 to 0xF; */
	umctl2_config.port_x_write_region0_last[4] = 0x3;
	umctl2_config.port_x_write_region1_last[4] = 0xE;

	/* configure 16 entries from 32 as low priority. other 16 will be high */
	/* MEMC_NO_OF_ENTRY is total number. synthesised as 32 */
	umctl2_config.lpr_num_entries = 16;

	/* configure time queues can be unserved until going critical; DFI clocks */
	umctl2_config.lpr_max_starve = 0x1FF;
	umctl2_config.hpr_max_starve = 0x1F;
	umctl2_config.w_max_starve = 0x7F;

	/* configure how long queues should run after going critical; number of transactions */
	umctl2_config.lpr_xact_run_length = 0xF;
	umctl2_config.hpr_xact_run_length = 0xF;
	umctl2_config.w_xact_run_length = 0xF;

	/* wait 3 cycles and then go to alternate store if it's non empty */
	/* for now go immediately to alternate store */
	umctl2_config.rdwr_idle_gap = 0;

	/* prefer the read transaction store. */
	umctl2_config.prefer_write = 0;

	/* enable pageclose mechanism */
	umctl2_config.pageclose = 0;
	umctl2_config.pageclose_timer = 0x27;

	/* disable auto-precharge for flushed command in collision case */
	umctl2_config.dis_collision_page_opt = 0;

	/* do not disable write combine */
	umctl2_config.dis_wc = 0;

	/* Use PUBL as PHY */
	umctl2_config.phy_init = &publ_init;
	umctl2_config.phy_idone = &publ_idone;
	umctl2_config.phy_start = &publ_start;
	umctl2_config.phy_train = &publ_train;
}

#if defined(CONFIG_SDCARD)
#if defined(CONFIG_OF_LIBFDT)
void at91_board_set_dtb_name(char *of_name)
{
	strcpy(of_name, "at91-sama7g5ek.dtb");
}
#endif

void at91_leds_init(void)
{
	const struct pio_desc leds_pins[] = {
		{"RED", AT91C_PIN_PB(8), 0, PIO_DEFAULT, PIO_OUTPUT},
		{"GREEN", AT91C_PIN_PA(13), 1, PIO_DEFAULT, PIO_OUTPUT},
		{"BLUE", AT91C_PIN_PD(20), 0, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_B},
	};

	pio_configure(leds_pins);
}

#define ATMEL_SDHC_GCKDIV_VALUE		1

void at91_sdhc_hw_init(void)
{
	u32 reg;
#if defined(CONFIG_SDHC1)
	const struct pio_desc sdmmc_pins[] = {
		{"SDMMC1_CK",   AT91C_PIN_PB(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_CMD",  AT91C_PIN_PB(29), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_DAT0", AT91C_PIN_PB(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_DAT1", AT91C_PIN_PC(0), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_DAT2", AT91C_PIN_PC(1), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_DAT3", AT91C_PIN_PC(2), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_VDDSEL", AT91C_PIN_PC(5), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_WP",   AT91C_PIN_PC(3), 1, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_CD",   AT91C_PIN_PC(4), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC1_RST",  AT91C_PIN_PB(22), 0, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
	const struct pio_desc sdmmc_pins_reset[] = {
		{"SDMMC1_RST",  AT91C_PIN_PB(22), 1, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
#endif
#if defined(CONFIG_SDHC0)
	const struct pio_desc sdmmc_pins[] = {
		{"SDMMC0_CK",   AT91C_PIN_PA(0), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_CMD",  AT91C_PIN_PA(1), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_RST",  AT91C_PIN_PA(2), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT0", AT91C_PIN_PA(3), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT1", AT91C_PIN_PA(4), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT2", AT91C_PIN_PA(5), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT3", AT91C_PIN_PA(6), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT4", AT91C_PIN_PA(7), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT5", AT91C_PIN_PA(8), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT6", AT91C_PIN_PA(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_DAT7", AT91C_PIN_PA(10), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_WP",   AT91C_PIN_PA(12), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SDMMC0_CD",  AT91C_PIN_PA(14), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};
#endif

	pio_configure(sdmmc_pins);
#if defined(CONFIG_SDHC1)
	mdelay(500);
	pio_configure(sdmmc_pins_reset);
#endif

	pmc_enable_periph_clock(CONFIG_SYS_ID_SDHC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_generic_clock(CONFIG_SYS_ID_SDHC,
				 GCK_CSS_SYSPLL_CLK,
				 ATMEL_SDHC_GCKDIV_VALUE);

	/* Launch calibration and wait till it's completed */
	reg = readl(CONFIG_SYS_BASE_SDHC + SDMMC_CALCR);
	reg |= SDMMC_CALCR_CLKDIV(7) | SDMMC_CALCR_ALWYSON;
	writel(reg, CONFIG_SYS_BASE_SDHC + SDMMC_CALCR);
	reg |= SDMMC_CALCR_EN;
	writel(reg, CONFIG_SYS_BASE_SDHC + SDMMC_CALCR);
	while (readl(CONFIG_SYS_BASE_SDHC + SDMMC_CALCR) & SDMMC_CALCR_EN)
		;
}
#endif

#ifdef CONFIG_TWI
static int twi1_bus_id = -1;

void twi_init()
{
#ifdef CONFIG_FLEXCOM
	twi1_bus_id = twi_bus_init(at91_flexcom1_init);
#endif
}
#endif /* CONFIG_TWI */

void cpu_voltage_select(void)
{
#ifdef CONFIG_MCP16502
	const struct mcp16502_cfg regulators_cfg[] = {
		{ .regulator = MCP16502_BUCK4,
#ifdef CONFIG_CPU_CLK_600MHZ
		 .uV = 1100000,
#endif
#ifdef CONFIG_CPU_CLK_800MHZ
		 .uV = 1200000,
#endif
#ifdef CONFIG_CPU_CLK_1000MHZ
		 .uV = 1250000,
#endif
		 .enable = 1, },
	};

	/*
	 * SAMA7G5's SHDWC keeps LPM pin low by default, so there is no need
	 * to pass pin descriptor to mcp16502_init() for switching to active
	 * state.
	 */
	if (mcp16502_init(twi1_bus_id, 0x5b, NULL, regulators_cfg,
				ARRAY_SIZE(regulators_cfg)))
		dbg_printf("MCP16502: init failure");
	else
		dbg_printf("MCP16502: CPU VDD @ %umV\n", regulators_cfg[0].uV / 1000);
#endif /* CONFIG_MCP16502 */
}

void hw_init(void)
{
	struct pmc_pll_cfg plla_config;
	struct pmc_pll_cfg ddrpll_config;
	struct pmc_pll_cfg syspll_config;
	struct pmc_pll_cfg imgpll_config;

#if BOOTSTRAP_DEBUG_LEVEL==DEBUG_VERY_LOUD
	volatile unsigned int *EXTRAM_CS;
	int i;
#endif

	unsigned int mck0_prescaler;

	/* Watchdog might be enabled out of reset. Let's make sure it's off */
	at91_disable_wdt();
#ifdef CONFIG_WDTS
	at91_disable_wdts();
#endif
	/* SMP is needed for L2 cache in cortex A7 */
	ca7_enable_smp();

	/* We need timers in the following steps */
	timer_init();

	/* Configure & Enable CPU PLL at a safe speed of 600 Mhz*/
	plla_config.mul = 24; /* 25 * 24 = 600 */
	plla_config.div = 0;
	plla_config.count = 0x3f;
	plla_config.fracr = 0;
	plla_config.acr = 0x1b040010;

	mck0_prescaler = BOARD_PRESCALER_CPUPLL | AT91C_PMC_MDIV_3;

	pmc_sam9x60_cfg_pll(PLL_ID_CPUPLL, &plla_config);

	pmc_mck_cfg_set(0, mck0_prescaler,
			AT91C_PMC_PRES | AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/*
	 * Out of Romcode, MCK1 & MCK4 are already configured with SYSPLL
	 * at unwanted frequencies.
	 * We cannot reconfigure SYSPLL without first reconfiguring
	 * MCK1 & MCK4 to different clock source.
	 * After we complete configuration for SYSPLL,
	 * we can configure MCK1 & MCK4 accordingly
	 */
	pmc_mck_cfg_set(4, BOARD_PRESCALER_MCK4_CLEAN,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	pmc_mck_cfg_set(1, BOARD_PRESCALER_MCK1_CLEAN,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	/*
	 * We are done here cleaning up MCK1 & MCK4, so we can configure
	 * SYSPLL at cruise speed
	 */

	/* Configure & Enable SYS PLL */
	syspll_config.mul = 49; /* (49 + 1) * 24 = 1200 */
	syspll_config.div = 2; /* Feed to PMC 1200/3 = 400 Mhz */
	syspll_config.count = 0x3f;
	syspll_config.fracr = 0;
	syspll_config.acr = 0x1b040010;
	/* SYSPLL @ 1200 MHz */
	pmc_sam9x60_cfg_pll(PLL_ID_SYSPLL, &syspll_config);

	/* MCK4 @ 400 Mhz (== SYSPLL) */
	pmc_mck_cfg_set(4, BOARD_PRESCALER_MCK4,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	/* MCK1 @ 200 Mhz (== SYSPLL/2) */
	pmc_mck_cfg_set(1, BOARD_PRESCALER_MCK1,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	board_flexcoms_init();
	at91_flexcom3_init();
	at91_leds_init();

	initialize_serial();

#ifdef CONFIG_TWI
	twi_init();
#endif

	dbg_very_loud("CA7 early uart\n");

	/* Configure & Enable DDR PLL */
#if CONFIG_MEM_CLOCK == 533
	ddrpll_config.mul = 43; /* (43 + 1) * 24 = 1056 */
	ddrpll_config.div = 1;
	ddrpll_config.divio = 100;
	ddrpll_config.count = 0x3f;
	ddrpll_config.fracr = 0x1aaaab; /* (10/24) * 2^22 to get extra 10 MHz */
	ddrpll_config.acr = 0x00020033;
	/* DDRPLL @ 1066 MHz */
#endif
#if CONFIG_MEM_CLOCK == 400
	ddrpll_config.mul = 32; /* (33 + 1) * 24 =  792*/
	ddrpll_config.div = 1;
	ddrpll_config.divio = 100;
	ddrpll_config.count = 0x3f;
	ddrpll_config.fracr = 0x155556; /* 2^22 / 3 */
	ddrpll_config.acr = 0x00020033;
	/* DDRPLL @ 800 MHz */
#endif
	pmc_sam9x60_cfg_pll(PLL_ID_DDRPLL, &ddrpll_config);

	/* MCK2 @ DDRPLL/2 MHz */
	pmc_mck_cfg_set(2, BOARD_PRESCALER_MCK2,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	/* Configure & Enable IMG PLL */
	imgpll_config.mul = 43; /* (43 + 1) * 24 = 1056 */
	imgpll_config.div = 3;
	imgpll_config.divio = 3;
	imgpll_config.count = 0x3f;
	imgpll_config.fracr = 0x155555; /* (8/24) * 2^22 to get extra 8 MHz */
	imgpll_config.acr = 0x1b040010;
	/* IMGPLL @ 1064 MHz */
	pmc_sam9x60_cfg_pll(PLL_ID_IMGPLL, &imgpll_config);

	/* MCK3 @ 266 MHz */
	pmc_mck_cfg_set(3, BOARD_PRESCALER_MCK3,
			AT91C_MCR_DIV | AT91C_MCR_CSS | AT91C_MCR_EN);

	dbg_printf("MCK: mck domains initialization complete.\n");

	tzc400_init();

	/* All MCK MUST be started before UMCTL2. Otherwise UMCTL2 will
	 * have the AXI ports blocked.
	 */
	umctl2_config_state_init();
	if (umctl2_init(&umctl2_config)) {
		console_printf("UMCTL2: Error initializing.\n");
	} else {
		console_printf("UMCTL2: Initialization complete.\n");
	}

#ifdef CONFIG_DDR3
	dbg_printf("MEMORY: %dMB DDR3 @ %d Mhz\n", MEM_SIZE / 1024 / 1024,
		   CONFIG_MEM_CLOCK);
#endif
#ifdef CONFIG_DDR2
	dbg_printf("MEMORY: %dMB DDR2 @ %d Mhz\n", MEM_SIZE / 1024 / 1024,
		   CONFIG_MEM_CLOCK);
#endif
#ifdef CONFIG_LPDDR2
	dbg_printf("MEMORY: %dMB LPDDR2 @ %d Mhz\n", MEM_SIZE / 1024 / 1024,
		   CONFIG_MEM_CLOCK);
#endif
#ifdef CONFIG_LPDDR3
	dbg_printf("MEMORY: %dMB LPDDR3 @ %d Mhz\n", MEM_SIZE / 1024 / 1024,
		   CONFIG_MEM_CLOCK);
#endif

#if BOOTSTRAP_DEBUG_LEVEL==DEBUG_VERY_LOUD
	EXTRAM_CS = (unsigned int *) MEM_BANK;
	for (i = 0; i < MEM_SIZE / 4; i++) {
		*EXTRAM_CS = 0xabcdef04;
		EXTRAM_CS++;
		if (!(i % 10000000)) dbg_very_loud("WRITING DDR\n");
	}

	EXTRAM_CS = (unsigned int *) MEM_BANK;
	for (i = 0; i < MEM_SIZE / 4; i++) {
		if (*EXTRAM_CS != 0xabcdef04)
			dbg_printf("ERRORS IN DDR %x:%x\n", EXTRAM_CS, *EXTRAM_CS);
		EXTRAM_CS++;
		if (!(i % 10000000)) dbg_very_loud("TESTING DDR\n");
	}
#endif

	cpu_voltage_select();

#ifdef CONFIG_CPU_CLK_800MHZ
	plla_config.mul = 32; /* 33 * 24 = 792 */
	plla_config.div = 0;
	plla_config.count = 0x3f;
	plla_config.fracr = 0x155556; /* 2^22 / 3 */
	plla_config.acr = 0x1b040010;

	mck0_prescaler = BOARD_PRESCALER_CPUPLL | AT91C_PMC_MDIV_4;
#endif
#ifdef CONFIG_CPU_CLK_1000MHZ
	plla_config.mul = 40; /* 41 * 24 = 984 */
	plla_config.div = 0;
	plla_config.count = 0x3f;
	plla_config.fracr = 0x2AAAAB; /* 2^22  * 2 / 3 */
	plla_config.acr = 0x1b040010;

	mck0_prescaler = BOARD_PRESCALER_CPUPLL | AT91C_PMC_MDIV_5;
#endif

#ifndef CONFIG_CPU_CLK_600MHZ
	pmc_sam9x60_cfg_pll(PLL_ID_CPUPLL, &plla_config);

	pmc_mck_cfg_set(0, mck0_prescaler,
			AT91C_PMC_PRES | AT91C_PMC_MDIV | AT91C_PMC_CSS);
#endif
}
