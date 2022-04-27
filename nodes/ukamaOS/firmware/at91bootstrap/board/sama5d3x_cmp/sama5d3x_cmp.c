/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2015, Atmel Corporation
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
#include "pmc.h"
#include "usart.h"
#include "debug.h"
#include "ddramc.h"
#include "spi.h"
#include "gpio.h"
#include "timer.h"
#include "watchdog.h"
#include "string.h"

#include "arch/at91_pmc/pmc.h"
#include "arch/at91_rstc.h"
#include "arch/sama5_smc.h"
#include "arch/at91_pio.h"
#include "arch/at91_ddrsdrc.h"
#include "sama5d3x_cmp.h"
#include "act8865.h"
#include "twi.h"
#include "board_hw_info.h"

static void at91_dbgu_hw_init(void)
{
	const struct pio_desc dbgu_pins[] = {
		{"RXD", AT91C_PIN_PB(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TXD", AT91C_PIN_PB(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(dbgu_pins);

	pmc_enable_periph_clock(AT91C_ID_DBGU, PMC_PERIPH_CLK_DIVIDER_NA);
}

static void initialize_dbgu(void)
{
	at91_dbgu_hw_init();
	usart_init(BAUDRATE(MASTER_CLOCK, 115200));
}

#if defined(CONFIG_LPDDR2)
static void lpddr2_reg_config(struct ddramc_register *ddramc_config)
{
	ddramc_config->mdr = (AT91C_DDRC2_DBW_32_BITS
				| AT91C_DDRC2_MD_LPDDR2_SDRAM);

	ddramc_config->cr = (AT91C_DDRC2_NC_DDR10_SDR9
				| AT91C_DDRC2_NR_14
				| AT91C_DDRC2_CAS_3
				| AT91C_DDRC2_ZQ_SHORT
				| AT91C_DDRC2_NB_BANKS_8
				| AT91C_DDRC2_UNAL_SUPPORTED);

	ddramc_config->lpddr2_lpr = AT91C_LPDDRC2_DS(0x03);
	/* 90n short calibration: ZQCS */
	ddramc_config->tim_calr = AT91C_DDRC2_ZQCS(12);

	/*
	 * The MT42128M32 refresh window: 32ms
	 * Required number of REFRESH commands(MIN): 8192
	 * (32ms / 8192) * 132MHz = 514 i.e. 0x202
	 */
	ddramc_config->rtr = 0x202;

	ddramc_config->cal_mr4r = AT91C_DDRC2_COUNT_CAL(0xC852);

	ddramc_config->t0pr = (AT91C_DDRC2_TRAS_(6)
			| AT91C_DDRC2_TRCD_(2)
			| AT91C_DDRC2_TWR_(3)
			| AT91C_DDRC2_TRC_(8)
			| AT91C_DDRC2_TRP_(2)
			| AT91C_DDRC2_TRRD_(2)
			| AT91C_DDRC2_TWTR_(2)
			| AT91C_DDRC2_TMRD_(3));

	ddramc_config->t1pr = (AT91C_DDRC2_TXP_(2)
			| AT91C_DDRC2_TXSNR_(18)
			| AT91C_DDRC2_TRFC_(17));

	ddramc_config->t2pr = (AT91C_DDRC2_TFAW_(8)
			| AT91C_DDRC2_TRTP_(2)
			| AT91C_DDRC2_TRPA_(3)
			| AT91C_DDRC2_TXARDS_(1)
			| AT91C_DDRC2_TXARD_(1));
}

static void lpddr2_init(void)
{
	struct ddramc_register ddramc_reg;
	unsigned int reg;

	lpddr2_reg_config(&ddramc_reg);

	pmc_enable_periph_clock(AT91C_ID_MPDDRC, PMC_PERIPH_CLK_DIVIDER_NA);
	pmc_enable_system_clock(AT91C_PMC_DDR);

	/*
	 * Initialize the special register for the SAMA5D3X_CMP.
	 * MPDDRC DLL Slave Offset Register: DDR2 configuration
	 */
	reg = AT91C_MPDDRC_S0OFF(0x04)
		| AT91C_MPDDRC_S1OFF(0x03)
		| AT91C_MPDDRC_S2OFF(0x04)
		| AT91C_MPDDRC_S3OFF(0x04);
	writel(reg, (AT91C_BASE_MPDDRC + MPDDRC_DLL_SOR));

	/*
	 * MPDDRC DLL Master Offset Register
	 * write master + clk90 offset
	 */
	reg = AT91C_MPDDRC_MOFF(7)
		| AT91C_MPDDRC_CLK90OFF(0x1F)
		| AT91C_MPDDRC_SELOFF_ENABLED | AT91C_MPDDRC_KEY;
	writel(reg, (AT91C_BASE_MPDDRC + MPDDRC_DLL_MOR));

	/*
	 * MPDDRC I/O Calibration Register
	 * DDR2 RZQ = 50 Ohm
	 * TZQIO = 4
	 */
	reg = readl(AT91C_BASE_MPDDRC + MPDDRC_IO_CALIBR);
	reg &= ~AT91C_MPDDRC_RDIV;
	reg &= ~AT91C_MPDDRC_TZQIO;
	reg |= AT91C_MPDDRC_RDIV_DDR2_RZQ_50;
	reg |= AT91C_MPDDRC_TZQIO_3;
	writel(reg, (AT91C_BASE_MPDDRC + MPDDRC_IO_CALIBR));

	lpddr2_sdram_initialize(AT91C_BASE_MPDDRC,
				AT91C_BASE_DDRCS, &ddramc_reg);
}
#else
#error "No right DDR-SDRAM device type provided"
#endif /* #ifdef CONFIG_LPDDR2 */

static void one_wire_hw_init(void)
{
	const struct pio_desc one_wire_pio[] = {
		{"1-Wire", AT91C_PIN_PE(25), 1, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOE, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(one_wire_pio);
}

static void SiI9022_hw_reset(void)
{
	pio_set_gpio_output(CONFIG_SYS_HDMI_RESET_PIN, 1);
	pio_set_gpio_output(CONFIG_SYS_HDMI_RESET_PIN, 0);
	udelay(500);
	pio_set_gpio_output(CONFIG_SYS_HDMI_RESET_PIN, 1);
}

#ifdef CONFIG_HW_INIT
void hw_init(void)
{
	/* Disable watchdog */
	at91_disable_wdt();

	/* Configure PLLA = MOSC * (PLL_MULA + 1) / PLL_DIVA */
	pmc_cfg_plla(PLLA_SETTINGS);

	/* Initialize PLLA charge pump */
	pmc_init_pll(AT91C_PMC_IPLLA_3);

	/* Switch PCK/MCK on Main clock output */
	pmc_mck_cfg_set(0, BOARD_PRESCALER_MAIN_CLOCK,
			AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Switch PCK/MCK on PLLA output */
	pmc_mck_cfg_set(0, BOARD_PRESCALER_PLLA,
			AT91C_PMC_MDIV | AT91C_PMC_CSS);

	/* Initialize timer */
	timer_init();

	/* Initialize the DBGU */
	initialize_dbgu();

#if defined(CONFIG_LPDDR2)
	lpddr2_init();
#endif

	/* Initialize 1-Wire */
	one_wire_hw_init();

	/* Reset HDMI SiI9022 */
	SiI9022_hw_reset();

#if defined(CONFIG_TWI)
	twi_init();
#endif
}
#endif /* #ifdef CONFIG_HW_INIT */

#ifdef CONFIG_DATAFLASH
void at91_spi0_hw_init(void)
{
	const struct pio_desc spi0_pins[] = {
		{"MISO",	AT91C_PIN_PD(10), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MOSI",	AT91C_PIN_PD(11), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"SPCK",	AT91C_PIN_PD(12), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"NPCS",	CONFIG_SYS_SPI_PCS, 1, PIO_DEFAULT, PIO_OUTPUT},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOD, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(spi0_pins);

	pmc_enable_periph_clock(AT91C_ID_SPI0, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_DATAFLASH */

#ifdef CONFIG_SDCARD
#ifdef CONFIG_OF_LIBFDT
void at91_board_set_dtb_name(char *of_name)
{
	strcat(of_name, "sama5d3x_cmp.dtb");
}
#endif

void at91_mci0_hw_init(void)
{
	const struct pio_desc mci_pins[] = {
		{"MCCK", AT91C_PIN_PD(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCCDA", AT91C_PIN_PD(0), 0, PIO_DEFAULT, PIO_PERIPH_A},

		{"MCDA0", AT91C_PIN_PD(1), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA1", AT91C_PIN_PD(2), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA2", AT91C_PIN_PD(3), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA3", AT91C_PIN_PD(4), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA4", AT91C_PIN_PD(5), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA5", AT91C_PIN_PD(6), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA6", AT91C_PIN_PD(7), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"MCDA7", AT91C_PIN_PD(8), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_HSMCI0, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(mci_pins);

	pmc_enable_periph_clock(AT91C_ID_HSMCI0, PMC_PERIPH_CLK_DIVIDER_NA);
}
#endif /* #ifdef CONFIG_SDCARD */

#ifdef CONFIG_NANDFLASH
void nandflash_hw_init(void)
{
	const struct pio_desc nand_pins[] = {
		{"NANDALE", AT91C_PIN_PE(21), 0, PIO_PULLUP, PIO_PERIPH_A},
		{"NANDCLE", AT91C_PIN_PE(22), 0, PIO_PULLUP, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pmc_enable_periph_clock(AT91C_ID_PIOE, PMC_PERIPH_CLK_DIVIDER_NA);
	pio_configure(nand_pins);

	pmc_enable_periph_clock(AT91C_ID_SMC, PMC_PERIPH_CLK_DIVIDER_NA);

	writel(AT91C_SMC_SETUP_NWE(1)
		| AT91C_SMC_SETUP_NCS_WR(1)
		| AT91C_SMC_SETUP_NRD(2)
		| AT91C_SMC_SETUP_NCS_RD(1),
		(ATMEL_BASE_SMC + SMC_SETUP3));

	writel(AT91C_SMC_PULSE_NWE(5)
		| AT91C_SMC_PULSE_NCS_WR(7)
		| AT91C_SMC_PULSE_NRD(5)
		| AT91C_SMC_PULSE_NCS_RD(7),
		(ATMEL_BASE_SMC + SMC_PULSE3));

	writel(AT91C_SMC_CYCLE_NWE(8)
		| AT91C_SMC_CYCLE_NRD(9),
		(ATMEL_BASE_SMC + SMC_CYCLE3));

	writel(AT91C_SMC_TIMINGS_TCLR(3)
		| AT91C_SMC_TIMINGS_TADL(10)
		| AT91C_SMC_TIMINGS_TAR(3)
		| AT91C_SMC_TIMINGS_TRR(4)
		| AT91C_SMC_TIMINGS_TWB(5)
		| AT91C_SMC_TIMINGS_RBNSEL(3)
		| AT91C_SMC_TIMINGS_NFSEL,
		(ATMEL_BASE_SMC + SMC_TIMINGS3));

	writel(AT91C_SMC_MODE_READMODE_NRD_CTRL
		| AT91C_SMC_MODE_WRITEMODE_NWE_CTRL
		| AT91C_SMC_MODE_EXNWMODE_DISABLED
		| AT91C_SMC_MODE_DBW_8
		| AT91C_SMC_MODE_TDF_CYCLES(1),
		(ATMEL_BASE_SMC + SMC_MODE3));
}
#endif /* #ifdef CONFIG_NANDFLASH */

/*--------------------- PM Support ----------------------------*/

#if defined(CONFIG_MAC0_PHY)
unsigned int at91_eth0_hw_init(void)
{
	unsigned int base_addr = AT91C_BASE_GMAC;

	const struct pio_desc macb_pins[] = {
		{"GMDC",	AT91C_PIN_PB(16), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"GMDIO",	AT91C_PIN_PB(17), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(macb_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOB, PMC_PERIPH_CLK_DIVIDER_NA);

	pmc_enable_periph_clock(AT91C_ID_GMAC, PMC_PERIPH_CLK_DIVIDER_NA);

	return base_addr;
}
#endif

#if defined(CONFIG_MAC1_PHY)
unsigned int at91_eth1_hw_init(void)
{
	unsigned int base_addr = AT91C_BASE_EMAC;

	const struct pio_desc macb_pins[] = {
		{"EMDC",	AT91C_PIN_PC(8), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"EMDIO",	AT91C_PIN_PC(9), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(macb_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOC, PMC_PERIPH_CLK_DIVIDER_NA);

	pmc_enable_periph_clock(AT91C_ID_EMAC, PMC_PERIPH_CLK_DIVIDER_NA);

	return base_addr;
}
#endif

#if defined(CONFIG_MACB)
void at91_disable_mac_clock(void)
{
#if defined(CONFIG_MAC0_PHY)
	pmc_disable_periph_clock(AT91C_ID_GMAC);
#endif
#if defined(CONFIG_MAC1_PHY)
	pmc_disable_periph_clock(AT91C_ID_EMAC);
#endif
}
#endif

#if defined(CONFIG_TWI)
#if defined(CONFIG_AUTOCONFIG_TWI_BUS)
void at91_board_config_twi_bus(void)
{
	unsigned int version = 'D';

	act8865_twi_bus = 1;

	/*
	 * ISI_VSYNC/ISI_HSYNC reused by TWI0, so move Audio / HDMI
	 * TWI interface from TWI0 to TWI1 from the ver.D
	 */
#if defined(CONFIG_LOAD_HW_INFO)
	version = get_ek_rev();
#endif
	if (version < 'D') {
		hdmi_twi_bus = 0;
		wm8904_twi_bus = 0;
	} else {
		hdmi_twi_bus = 1;
		wm8904_twi_bus = 1;
	}
}
#endif

unsigned int at91_twi0_hw_init(void)
{
	unsigned int base_addr = AT91C_BASE_TWI0;

	const struct pio_desc twi_pins[] = {
		{"TWD", AT91C_PIN_PA(30), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{"TWCK", AT91C_PIN_PA(31), 0, PIO_DEFAULT, PIO_PERIPH_A},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(twi_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOA, PMC_PERIPH_CLK_DIVIDER_NA);

	pmc_enable_periph_clock(AT91C_ID_TWI0, PMC_PERIPH_CLK_DIVIDER_NA);

	return base_addr;
}

unsigned int at91_twi1_hw_init(void)
{
	unsigned int base_addr = AT91C_BASE_TWI1;

	const struct pio_desc twi_pins[] = {
		{"TWD", AT91C_PIN_PC(26), 0, PIO_DEFAULT, PIO_PERIPH_B},
		{"TWCK", AT91C_PIN_PC(27), 0, PIO_DEFAULT, PIO_PERIPH_B},
		{(char *)0, 0, 0, PIO_DEFAULT, PIO_PERIPH_A},
	};

	pio_configure(twi_pins);
	pmc_enable_periph_clock(AT91C_ID_PIOC, PMC_PERIPH_CLK_DIVIDER_NA);

	pmc_enable_periph_clock(AT91C_ID_TWI1, PMC_PERIPH_CLK_DIVIDER_NA);

	return base_addr;
}

unsigned int at91_twi2_hw_init(void)
{
	return 0;
}
#endif

#if defined(CONFIG_ACT8865_SET_VOLTAGE)
int at91_board_act8865_set_reg_voltage(void)
{
	unsigned char reg, value;
	int ret;

	/* Check ACT8865 I2C interface */
	if (act8865_check_i2c_disabled())
		return 0;

	/* Enable REG4(VDDANA) output 3.3V */
	reg = REG4_0;
	value = ACT8865_3V3;
	ret = act8865_set_reg_voltage(reg, value);
	if (ret)
		console_printf("ACT8865: Failed to make REG4 output 3300mV\n");

	/* Enable REG2 output 1.25V */
	reg = REG2_0;
	value = ACT8865_1V25;
	ret = act8865_set_reg_voltage(reg, value);
	if (ret)
		console_printf("ACT8865: Failed to make REG2 output 1250mV\n");

	return 0;
}
#endif

#ifdef CONFIG_PM
void at91_disable_smd_clock(void)
{
	/*
	 * Set pin DIBP to pull-up and DIBN to pull-down
	 * to save power on VDDIOP0
	 */
	pmc_enable_system_clock(AT91C_PMC_SMDCK);
	pmc_set_smd_clock_divider(AT91C_PMC_SMDDIV);
	pmc_enable_periph_clock(AT91C_ID_SMD, PMC_PERIPH_CLK_DIVIDER_NA);
	writel(0xF, (0x0C + AT91C_BASE_SMD));
	pmc_disable_periph_clock(AT91C_ID_SMD);
	pmc_disable_system_clock(AT91C_PMC_SMDCK);
}
#endif


#if defined(CONFIG_TWI)
void twi_init()
{
	twi_bus_init(at91_twi0_hw_init);
	twi_bus_init(at91_twi1_hw_init);
	twi_bus_init(at91_twi2_hw_init);
#if defined(CONFIG_AUTOCONFIG_TWI_BUS)
	dbg_loud("Auto-Config the TWI Bus by the board\n");
	at91_board_config_twi_bus();
#endif
}
#endif
