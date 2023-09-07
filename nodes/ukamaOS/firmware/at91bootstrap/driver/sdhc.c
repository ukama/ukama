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
#include "hardware.h"
#include "board.h"
#include "mci_media.h"
#include "div.h"
#include "timer.h"
#include "debug.h"
#include "pmc.h"
#include "sdhc.h"
#include "sdhc_cal.h"

/*
 * Registers Definitions
 */
#define	SDMMC_SSAR	0x00	/* SDMA System Address / Argument 2 Register */
#define	SDMMC_BSR	0x04	/* Block Size Register */
#define	SDMMC_BCR	0x06	/* Block Count Register */
#define	SDMMC_ARG1R	0x08	/* Argument 1 Register */
#define	SDMMC_TMR	0x0c	/* Transfer Mode Register */
#define	SDMMC_CR	0x0e	/* Command Register */
#define	SDMMC_RR0	0x10	/* Response Register 0 */
#define	SDMMC_RR1	0x14	/* Response Register 1 */
#define	SDMMC_RR2	0x18	/* Response Register 2 */
#define	SDMMC_RR3	0x1c	/* Response Register 3 */
#define	SDMMC_BDPR	0x20	/* Buffer Data Port Register */
#define	SDMMC_PSR	0x24	/* Present State Register */
#define	SDMMC_HC1R	0x28	/* Host Control 1 Register */
#define	SDMMC_PCR	0x29	/* Power Control Register */
#define	SDMMC_BGCR	0x2a	/* Block Gap Control Register */
#define	SDMMC_WCR	0x2b	/* Wakeup Control Register */
#define	SDMMC_CCR	0x2c	/* Clock Control Register */
#define	SDMMC_TCR	0x2e	/* Timeout Control Register */
#define	SDMMC_SRR	0x2f	/* Software Reset Register */
#define	SDMMC_NISTR	0x30	/* Normal Interrupt Status Register */
#define	SDMMC_EISTR	0x32	/* Error Interrupt Status Register */
#define	SDMMC_NISTER	0x34	/* Normal Interrupt Status Enable Register */
#define	SDMMC_EISTER	0x36	/* Error Interrupt Status Enable Register */
#define	SDMMC_NISIER	0x38	/* Normal Interrupt Signal Enable Register */
#define	SDMMC_EISIER	0x3a	/* Error Interrupt Signal Enable Register */
#define	SDMMC_ACESR	0x3c	/* Auto CMD Error Status Register */
#define	SDMMC_HC2R	0x3e	/* Host Control 2 Register */
#define	SDMMC_CA0R	0x40	/* Capabilities 0 Register */
#define	SDMMC_CA1R	0x44	/* Capabilities 1 Register */
#define	SDMMC_MCARR	0x48	/* Maximum Current Capabilities Register */
/* 0x4c: Reserved */
#define	SDMMC_FERACES	0x50	/* Force Event Register for Auto CMD Error Status */
#define	SDMMC_FEREIS	0x52	/* Force Event Register for Error Interrupt Status */
#define	SDMMC_AESR	0x54	/* ADMA Error Status Register */
#define	SDMMC_ASAR0	0x58	/* ADMA System Address Register 0 */
/* 0x5c: Reserved */
#define	SDMMC_PVR0	0x60	/* Preset Value Register 0 (for initialization) */
#define	SDMMC_PVR1	0x62	/* Preset Value Register 1 (for Default Speed) */
#define	SDMMC_PVR2	0x64	/* Preset Value Register 2 (for High Speed) */
#define	SDMMC_PVR3	0x66	/* Preset Value Register 3 (for SDR12) */
#define	SDMMC_PVR4	0x68	/* Preset Value Register 4 (for SDR25) */
#define	SDMMC_PVR5	0x6a	/* Preset Value Register 5 (for SDR50) */
#define	SDMMC_PVR6	0x6c	/* Preset Value Register 6 (for SDR104) */
#define	SDMMC_PVR7	0x6e	/* Preset Value Register 7 (for DDR50) */
/* 0x70 ~ 0xf8: Reserved */
#define	SDMMC_SISR	0xfc	/* Slot Interrupt Status Register */
#define	SDMMC_HCVR	0xfe	/* Host Control Version Register */
/* 0x100 ~ 0x1fc: reserved */
#define	SDMMC_APSR	0x200	/* Additional Present State Register */
#define	SDMMC_MC1R	0x204	/* MMC Control 1 Register */
#define	SDMMC_MC2R	0x205	/* MMC Control 2 Register */
#define	SDMMC_ACR	0x208	/* AHB Control Register */
#define	SDMMC_CC2R	0x20c	/* Clock Control 2 Register */
#define	SDMMC_RTC1R	0x210	/* Re-Tuning Timer Control 1 Register */
#define	SDMMC_RTC2R	0x211	/* Re-Tuning Timer Control 2 Register */
#define	SDMMC_RTCVR	0x214	/* Re-Tuning Timer Control Value Register */
#define	SDMMC_RTISTER	0x218	/* Re-Tuning Timer Interrupt Status Enable Register */
#define	SDMMC_RTISIER	0x219	/* Re-Tuning Timer Interrupt Signal Enable Register */
#define	SDMMC_RTISTR	0x21c	/* Re-Tuning Timer Interrupt Signal Enable Register */
#define	SDMMC_RTSSTR	0x21d	/* Re-Tuning Timer Statue Slots Register */
#define	SDMMC_TUNCR	0x220	/* Tuning Control Register */
/* 0x224 ~ 0x22c: Reserved */
#define	SDMMC_CACR	0x230	/* Capabilities Control Register */
/* 0x234 ~ 0x23c: Reserved */
/* SDMMC_CACR register definition moved to sdhc_cal.h file */
/* 0x224 ~ 0x2fc: Reserved */

/*
 * Register Field Definitions
 */
/* SDMMC_TMR */
#define	SDMMC_TMR_DMAEN		(0x1 << 0)	/* DMA Enable */
#define	SDMMC_TMR_BCEN		(0x1 << 1)	/* Block Count Enable */
#define	SDMMC_TMR_ACMDEN	(0x3 << 2)	/* Auto Command Enable */
#define		SDMMC_TMR_ACMDEN_DISABLED	(0x0 << 2)
#define		SDMMC_TMR_ACMDEN_CMD12		(0x1 << 2)
#define		SDMMC_TMR_ACMDEN_CMD23		(0x2 << 2)
#define	SDMMC_TMR_DTDSEL	(0x1 << 4)	/* Data Transfer Direction Selection */
#define		SDMMC_TMR_DTDSEL_WRITE		(0x0 << 4)
#define		SDMMC_TMR_DTDSEL_READ		(0x1 << 4)
#define	SDMMC_TMR_MSBSEL	(0x1 << 5)	/* Multi/Single Block Selection */

/* SDMMC_CR */
#define	SDMMC_CR_RESPTYP	(0x3 << 0)	/* Respone Type */
#define		SDMMC_CR_RESPTYP_(x)		((x) << 0)
#define		SDMMC_CR_RESPTYP_OFFSET		 0
#define		SDMMC_CR_RESPTYP_NORESP		0x0
#define		SDMMC_CR_RESPTYP_RL136		0x1
#define		SDMMC_CR_RESPTYP_RL48		0x2
#define		SDMMC_CR_RESPTYP_RL48BUSY	0x3
#define	SDMMC_CR_CMDCCEN	(0x1 << 3)	/* Command CRC Check Enable */
#define	SDMMC_CR_CMDICEN	(0x1 << 4)	/* Command Index Check Enable */
#define	SDMMC_CR_DPSEL		(0x1 << 5)	/* Data Present Select */
#define	SDMMC_CR_CMDTYP		(0x3 << 6)	/* Command Type */
#define		SDMMC_CR_CMDTYP_(x)		((x) << 6)
#define		SDMMC_CR_CMDTYP_OFFSET		6
#define		SDMMC_CR_CMDTYP_NORMAL		0x0
#define		SDMMC_CR_CMDTYP_SUSPEND		0x1
#define		SDMMC_CR_CMDTYP_RESUME		0x2
#define		SDMMC_CR_CMDTYP_ABORT		0x3
#define	SDMMC_CR_CMDIDX		(0x3f << 8)	/* Command Index */
#define		SDMMC_CR_CMDIDX_(x)		((x) << 8)

/* SDMMC_PSR */
#define	SDMMC_PSR_CMDINHC	(0x1 << 0)	/* Command Inhibit(CMD) */
#define	SDMMC_PSR_CMDINHD	(0x1 << 1)	/* Command Inhibit(DAT) */
#define	SDMMC_PSR_DLACT		(0x1 << 2)	/* DAT Line Active */
#define	SDMMC_PSR_WTACT		(0x1 << 8)	/* Write Transfer Active */
#define	SDMMC_PSR_RTACT		(0x1 << 9)	/* Read Transfer Active */
#define	SDMMC_PSR_BUFWREN	(0x1 << 10)	/* Buffer Write Enable */
#define	SDMMC_PSR_BUFRDEN	(0x1 << 11)	/* Buffer Read Enable */
#define	SDMMC_PSR_CARDINS	(0x1 << 16)	/* Card Inserted*/
#define	SDMMC_PSR_CARDSS	(0x1 << 17)	/* Card State Stable */
#define	SDMMC_PSR_CARDDPL	(0x1 << 18)	/* Card Detect Pin Level */
#define	SDMMC_PSR_WRPPL		(0x1 << 19)	/* Write Protect Pin Level */
#define	SDMMC_PSR_DATLL		(0xf << 20)	/* DAT[3:0] Line Level */
#define	SDMMC_PSR_CMDLL		(0xf << 24)	/* CMD Line Level */

/* SDMMC_NISTR/SDMMC_NISTER/SDMMC_NISIER(e.MMC) */
#define	SDMMC_NISTR_CMDC	(0x1 << 0)	/* Command Complete */
#define	SDMMC_NISTR_TRFC	(0x1 << 1)	/* Transfer Complete */
#define	SDMMC_NISTR_BLKGE	(0x1 << 2)	/* Block Gap Event */
#define	SDMMC_NISTR_DMAINT	(0x1 << 3)	/* DMA Interrupt */
#define	SDMMC_NISTR_BWRRDY	(0x1 << 4)	/* Buffer Write Ready */
#define	SDMMC_NISTR_BRDRDY	(0x1 << 5)	/* Buffer Read Ready */
#define	SDMMC_NISTR_CINS	(0x1 << 6)	/* Card Insertion */
#define	SDMMC_NISTR_CREM	(0x1 << 7)	/* Card Removal */
#define	SDMMC_NISTR_CINT	(0x1 << 8)	/* Card Interrupt */
#define	SDMMC_NISTR_BOOTAR	(0x1 << 14)	/* Boot Acknowledge Received */
#define	SDMMC_NISTR_ERRINT	(0x1 << 15)	/* Error Interrupt */

/* SDMMC_EISTR/SDMMC_EISTER/SDMMC_EISIER(e.MMC) */
#define	SDMMC_EISTR_CMDTEO	(0x1 << 0)	/* Command Timeout Error */
#define	SDMMC_EISTR_CMDCRC	(0x1 << 1)	/* Command CRC Error */
#define	SDMMC_EISTR_CMDEND	(0x1 << 2)	/* Command End Bit Error */
#define	SDMMC_EISTR_CMDIDX	(0x1 << 3)	/* Command Index Error */
#define	SDMMC_EISTR_DATTEO	(0x1 << 4)	/* Data Timeout Error */
#define	SDMMC_EISTR_DATCRC	(0x1 << 5)	/* Data CRC Error */
#define	SDMMC_EISTR_DATEND	(0x1 << 6)	/* Data End Bit Error */
#define	SDMMC_EISTR_CURLIM	(0x1 << 7)	/* Current Limit Error */
#define	SDMMC_EISTR_ACMD	(0x1 << 8)	/* Auto CMD Error */
#define	SDMMC_EISTR_ADMA	(0x1 << 9)	/* ADMA Error */
#define	SDMMC_EISTR_BOOTAE	(0x1 << 12)	/* Boot Acknowledge Error */

/* SDMMC_MC1R */
#define	SDMMC_MC1R_CMDTYP	(0x3 << 0)	/* e.MMC Command Type */
#define		SDMMC_MC1R_CMDTYP_NORMAL	(0x0 << 0)
#define		SDMMC_MC1R_CMDTYP_WAITIRQ	(0x1 << 0)
#define		SDMMC_MC1R_CMDTYP_STREAM	(0x2 << 0)
#define		SDMMC_MC1R_CMDTYP_BOOT		(0x3 << 0)
#define	SDMMC_MC1R_DDR		(0x01 << 3)	/* e.MMC HSDDR Mode */
#define	SDMMC_MC1R_OPD		(0x01 << 4)	/* e.MMC Open Drain Mode */
#define	SDMMC_MC1R_BOOTA	(0x01 << 5)	/* e.MMC Boot Acknowledge Enable */
#define	SDMMC_MC1R_RSRN		(0x01 << 6)	/* e.MMC Reset Signal */
#define	SDMMC_MC1R_FCD		(0x01 << 7)	/* e.MMC Force Card Detect */

/* SDMMC_MC2R(SDMMC e.MMC Control 2 Register) */
#define	SDMMC_MC2R_SRESP	(0x1 >> 0)	/* e.MMC Abort Wait IRQ */
#define	SDMMC_MC2R_ABOOT	(0x1 >> 1)	/* e.MMC Abort Boot */

/* SDMMC_PCR */
#define	SDMMC_PCR_SDBPWR	(0x1 << 0)	/* SD Bud Power */
#define	SDMMC_PCR_SDBVSEL	(0x7 << 1)	/* SD Bus Voltage Select */

/* SDMMC_CCR */
#define	SDMMC_CCR_INTCLKEN	(0x1 << 0)	/* Internal Clock Enable */
#define	SDMMC_CCR_INTCLKS	(0x1 << 1)	/* Internal Clock Stable */
#define	SDMMC_CCR_SDCLKEN	(0x1 << 2)	/* SD Clock Enable */
#define	SDMMC_CCR_CLKGSEL	(0x1 << 5)	/* Clock Generator Select */
#define	SDMMC_CCR_USDCLKFSEL	(0x3 << 6)	/* Upper Bits of SDCLK Frequency Select */
#define		SDMMC_CCR_USDCLKFSEL_MSK	3
#define		SDMMC_CCR_USDCLKFSEL_OFFSET	6
#define	SDMMC_CCR_SDCLKFSEL	(0xff << 8)	/* SDCLK Frequency Select */
#define		SDMMC_CCR_SDCLKFSEL_MSK		0xff
#define		SDMMC_CCR_SDCLKFSEL_OFFSET	8

/* SDMMC_SRR */
#define	SDMMC_SRR_SWRSTALL	(0x1 << 0)	/* Software reset for All */
#define	SDMMC_SRR_SWRSTCMD	(0x1 << 1)	/* Software reset for CMD line */
#define	SDMMC_SRR_SWRSTDAT	(0x1 << 2)	/* Software reset for DAT line */

/* SDMMC_CA0R */
#define	SDMMC_CA0R_TEOCLKF	(0x3f << 0)	/* Timeout Clock Frequency */
#define	SDMMC_CA0R_TEOCLKU	(0x1 << 7)	/* Timeout Clock Unit */
#define	SDMMC_CA0R_BASECLKF		(0xff << 8)	/* Base Clock Frequency */
#define		SDMMC_CA0R_BASECLKF_MSK		0xff
#define		SDMMC_CA0R_BASECLKF_OFFSET	8
#define	SDMMC_CA0R_MAXBLKL	(0x3 << 16)	/* Max Block Length */
#define		SDMMC_CA0R_MAXBLKL_512		(0x0 << 16)
#define		SDMMC_CA0R_MAXBLKL_1024		(0x1 << 16)
#define		SDMMC_CA0R_MAXBLKL_2048		(0x2 << 16)
#define	SDMMC_CA0R_ED8SUP	(0x1 << 18)	/* 8-Bit Support for Embedded Device */
#define	SDMMC_CA0R_ADMA2SUP	(0x1 << 19)	/* ADMA2 Support */
#define	SDMMC_CA0R_HSSUP	(0x1 << 21)	/* High Speed Support */
#define	SDMMC_CA0R_SDMASUP	(0x1 << 22)	/* SDMA Support */
#define	SDMMC_CA0R_SRSUP	(0x1 << 23)	/* Suspend/Resume Support */
#define	SDMMC_CA0R_V33VSUP	(0x1 << 24)	/* Voltage Support 3.3V */
#define	SDMMC_CA0R_V30VSUP	(0x1 << 25)	/* Voltage Support 3.0V */
#define	SDMMC_CA0R_V18VSUP	(0x1 << 26)	/* Voltage Support 1.8V */
#define	SDMMC_CA0R_SB64SUP	(0x1 << 28)	/* 64-Bit System Bus Support */
#define	SDMMC_CA0R_ASINTSUP	(0x1 << 29)	/* Asynchronous Interrupt Support */
#define	SDMMC_CA0R_SLTYPE	(0x3 << 30)	/* Slot Type */

/* SDMMC_CA1R */
#define	SDMMC_CA1R_SDR50SUP	(0x1 << 0)	/* SDR50 Support */
#define	SDMMC_CA1R_SDR104SUP	(0x1 << 1)	/* SDR104 Support */
#define	SDMMC_CA1R_DDR50SUP	(0x1 << 2)	/* DDR50 Support */
#define	SDMMC_CA1R_DRVASUP	(0x1 << 4)	/* Driver Type A Support */
#define	SDMMC_CA1R_DRVCSUP	(0x1 << 5)	/* Driver Type C Support */
#define	SDMMC_CA1R_DRVDSUP	(0x1 << 6)	/* Driver Type D Support */
#define	SDMMC_CA1R_DRVDSUP	(0x1 << 6)	/* Driver Type D Support */
#define	SDMMC_CA1R_TCNTRT	(0xf << 8)	/* Timer Count For Re-Tuning */
#define	SDMMC_CA1R_TSDR50	(0x1 << 13)	/* Use Tuning for SDR50 */
#define	SDMMC_CA1R_RTMOD	(0x3 << 14)	/* Re-Tuning Modes */
#define		SDMMC_CA1R_RTMOD_MODE1		(0x0 << 14)
#define		SDMMC_CA1R_RTMOD_MODE2		(0x1 << 14)
#define		SDMMC_CA1R_RTMOD_MODE3		(0x2 << 14)
#define	SDMMC_CA1R_CLKMULT	(0xff << 16)	/* Clock Multiplier */
#define		SDMMC_CA1R_CLKMULT_MSK		0xff
#define		SDMMC_CA1R_CLKMULT_OFFSET	16

/* SDMMC_HC1R */
#define	SDMMC_HC1R_LEDCTRL	(0x1 << 0)	/* LED Control */
#define	SDMMC_HC1R_DW		(0x1 << 1)	/* Data Width */
#define		SDMMC_HC1R_DW_1_BIT	(0x0 << 1)
#define		SDMMC_HC1R_DW_4_BIT	(0x1 << 1)
#define	SDMMC_HC1R_HSEN		(0x1 << 2)	/* High Speed Enable */
#define	SDMMC_HC1R_DMASEL	(0x3 << 3)	/* DMA Select */
#define		SDMMC_HC1R_DMASEL_SDMA		(0x0 << 3)
#define		SDMMC_HC1R_DMASEL_ADMA32	(0x2 << 3)
#define	SDMMC_HC1R_EXTDW	(0x1 << 5)	/* Extended Data Width */
#define	SDMMC_HC1R_CARDDTL	(0x1 << 6)	/* Card Detect Test Level */
#define	SDMMC_HC1R_CARDDSEL	(0x1 << 7)	/* Card Detect Signal Selection */

/*---------------------------------------------------------------*/

static struct sd_host sdhc_host;

static unsigned int sdhc_get_base(void)
{
	return CONFIG_SYS_BASE_SDHC;
}

static unsigned int sdhc_readl(unsigned int reg)
{
	return readl(sdhc_get_base() + reg);
}

static void sdhc_writel(unsigned int reg, unsigned int value)
{
	writel(value, sdhc_get_base() + reg);
}

static unsigned short sdhc_readw(unsigned int reg)
{
	return readw(sdhc_get_base() + reg);
}

static void sdhc_writew(unsigned int reg, unsigned short value)
{
	writew(value, sdhc_get_base() + reg);
}

static unsigned char sdhc_readb(unsigned int reg)
{
	return readb(sdhc_get_base() + reg);
}

static void sdhc_writeb(unsigned int reg, unsigned char value)
{
	writeb(value, sdhc_get_base() + reg);
}

static void sdhc_software_reset(void)
{
	sdhc_writeb(SDMMC_SRR, SDMMC_SRR_SWRSTALL);

	while (sdhc_readb(SDMMC_SRR) & SDMMC_SRR_SWRSTALL)
		;
}

static void sdhc_software_reset_cmd(void)
{
	sdhc_writeb(SDMMC_SRR, SDMMC_SRR_SWRSTCMD);

	while (sdhc_readb(SDMMC_SRR) & SDMMC_SRR_SWRSTCMD)
		;
}

static void sdhc_software_reset_dat(void)
{
	sdhc_writeb(SDMMC_SRR, SDMMC_SRR_SWRSTDAT);

	while (sdhc_readb(SDMMC_SRR) & SDMMC_SRR_SWRSTDAT)
		;
}

static void sdhc_set_power(void)
{
	unsigned char value = sdhc_readb(SDMMC_PCR);

	sdhc_writeb(SDMMC_PCR, value | SDMMC_PCR_SDBPWR);
}

static int sdhc_set_clock(struct sd_card *sdcard, unsigned int clock)
{
	struct sd_host *host = sdcard->host;
	unsigned int clk_gen_sel = 0;
	unsigned int clk_div;
	unsigned int reg;
	unsigned int timeout;

	timeout = 100000;
	while ((--timeout) &&
	       (sdhc_readl(SDMMC_PSR) & (SDMMC_PSR_CMDINHC | SDMMC_PSR_CMDINHD)))
		;

	if (!timeout)
		dbg_info("SDHC: Timeout waiting for CMD and DAT Inhibit bits\n");

	reg = sdhc_readw(SDMMC_CCR);
	reg &= ~SDMMC_CCR_SDCLKEN;
	sdhc_writew(SDMMC_CCR, reg);

	if (clock > host->caps_max_clock)
		clock = host->caps_max_clock;

	if (clock < host->caps_min_clock)
		clock = host->caps_min_clock;

	if (host->caps_clk_mult) {
		for (clk_div = 1; clk_div <= 1024; clk_div++) {
			if (div(host->caps_max_clock, clk_div) <= clock)
				break;
		}
		clk_gen_sel = SDMMC_CCR_CLKGSEL;
		clk_div -= 1;
	} else {
		if (host->caps_max_clock <= clock) {
			clk_div = 1;
		} else {
			for (clk_div = 2; clk_div < 2048; clk_div += 2) {
				if (div(host->caps_max_clock, clk_div) <= clock)
					break;
			}
		}
		clk_div >>= 1;
	}

	sdhc_writew(SDMMC_CCR, sdhc_readw(SDMMC_CCR) & ~SDMMC_CCR_SDCLKEN);

	sdhc_writew(SDMMC_CCR, SDMMC_CCR_INTCLKEN | clk_gen_sel
			| ((clk_div & 0xff) << SDMMC_CCR_SDCLKFSEL_OFFSET)
			| (((clk_div >> 8) & SDMMC_CCR_USDCLKFSEL_MSK)
					<< SDMMC_CCR_USDCLKFSEL_OFFSET));

	timeout = 1000000;
	while ((--timeout) && (!(sdhc_readw(SDMMC_CCR) & SDMMC_CCR_INTCLKS)))
		;

	if (!timeout)
		dbg_info("SDHC: Timeout waiting for internal clock ready\n");

	sdhc_writew(SDMMC_CCR, sdhc_readw(SDMMC_CCR) | SDMMC_CCR_SDCLKEN);

	reg = sdhc_readb(SDMMC_HC1R);
	if (clock > 26000000)
		reg |= SDMMC_HC1R_HSEN;
	else
		reg &= ~SDMMC_HC1R_HSEN;

	sdhc_writeb(SDMMC_HC1R, reg);

	return 0;
}

static int sdhc_set_bus_width(struct sd_card *sdcard, unsigned int width)
{
	struct sd_host *host = sdcard->host;
	unsigned int reg;

	reg = sdhc_readb(SDMMC_HC1R);
	if (width == 8) {
		if (!(host->caps_bus_width & BUS_WIDTH_8_BIT))
			return -1;

		reg |= SDMMC_HC1R_EXTDW;
	} else if (width == 4) {
		if (!(host->caps_bus_width & BUS_WIDTH_4_BIT))
			return -1;

		reg &= ~SDMMC_HC1R_EXTDW;
		reg |= SDMMC_HC1R_DW_4_BIT;
	} else {
		reg &= ~SDMMC_HC1R_EXTDW;
		reg &= ~SDMMC_HC1R_DW;
	}
	sdhc_writeb(SDMMC_HC1R, reg);
	sdcard->configured_bus_w = width;

	return 0;
}

static int sdhc_set_ddr(struct sd_card *sdcard)
{
	if (!sdcard->host->caps_ddr)
		return -1;

	sdhc_writeb(SDMMC_MC1R, sdhc_readb(SDMMC_MC1R) | SDMMC_MC1R_DDR);
	sdcard->ddr = 1;
	return 0;
}

static int sdhc_host_capability(struct sd_card *sdcard)
{
	struct sd_host *host = sdcard->host;
	unsigned int caps;

	host->caps_max_clock = pmc_get_generic_clock(CONFIG_SYS_ID_SDHC);
	host->caps_min_clock = host->caps_max_clock / 2048;

	caps = sdhc_readl(SDMMC_CA0R);

	host->caps_bus_width = BUS_WIDTH_1_BIT | BUS_WIDTH_4_BIT;
	if (caps & SDMMC_CA0R_ED8SUP)
		host->caps_bus_width |= BUS_WIDTH_8_BIT;

	host->caps_high_speed = 0;
	host->caps_ddr = 0;
	host->caps_adma2 = 0;
	if (caps & SDMMC_CA0R_HSSUP)
		host->caps_high_speed = 1;
#if !defined(CONFIG_SDHC_NODMA)
	if (caps & SDMMC_CA0R_ADMA2SUP) {
		dbg_printf("MMC: ADMA supported\n");
		host->caps_adma2 = 1;
	}
#endif

	host->caps_voltages = 0;

	if (caps & SDMMC_CA0R_V33VSUP)
		host->caps_voltages |= SD_OCR_VDD_32_33 | SD_OCR_VDD_33_34;
	if (caps & SDMMC_CA0R_V30VSUP)
		host->caps_voltages |= SD_OCR_VDD_29_30 | SD_OCR_VDD_30_31;
	if (caps & SDMMC_CA0R_V18VSUP)
		host->caps_voltages |= SD_OCR_VDD_165_195;

	caps = sdhc_readl(SDMMC_CA1R);

	host->caps_clk_mult = (caps >> SDMMC_CA1R_CLKMULT_OFFSET)
						& SDMMC_CA1R_CLKMULT_MSK;
	if (caps & SDMMC_CA1R_DDR50SUP)
		host->caps_ddr = 1;

	return 0;
}

static int sdhc_is_card_inserted(struct sd_card *sdcard)
{
	/*
	 * Debouncing of the card detect pin is up to 13ms on sama5d2 rev B
	 * and later.
	 * Try to be safe and wait for up to 50ms (50000µs). Let assume
	 * the PCK (processor clock) frequency is 500MHz, hence 500 cycles/µs.
	 * 500 * 50000 = 25000000 cycles.
	 */
	unsigned int timeout = 25000000;
	int is_inserted = 0;

	/*
	 * Enable (unmask) the 'card inserted' bit in the Normal Interrupt
	 * Status Register.
	 */
	sdhc_writew(SDMMC_NISTER, SDMMC_NISTR_CINS);

	/*
	 * Check whether the 'card inserted' bit is already set in
	 * the Present State Register.
	 */
	if (sdhc_readl(SDMMC_PSR) & SDMMC_PSR_CARDINS) {
		is_inserted = 1;
		goto exit;
	}

	/* Poll the Normal Interrupt Status Register for bit 'card inserted'. */
	while (!(sdhc_readw(SDMMC_NISTR) & SDMMC_NISTR_CINS) &&
	       timeout--);

	is_inserted = !!(sdhc_readw(SDMMC_NISTR) & SDMMC_NISTR_CINS);

exit:
	/*
	 * Disable (mask) the 'card inserted' bit in the Normal Interrupt
	 * Status Register.
	 */
	sdhc_writew(SDMMC_NISTER, 0);

	/* Clear the pending 'card inserted' interrupt. */
	sdhc_writew(SDMMC_NISTR, SDMMC_NISTR_CINS);

	return is_inserted;
}

static int sdhc_init(struct sd_card *sdcard)
{
	unsigned int normal_status_mask, error_status_mask;

	sdhc_software_reset();

	sdhc_set_power();

	sdhc_host_capability(sdcard);
	
	if (sdhc_is_card_inserted(sdcard) <= 0) {
		dbg_info("SDHC: Error: No Card Inserted\n");
	}
	normal_status_mask = SDMMC_NISTR_CMDC
				| SDMMC_NISTR_TRFC
				| SDMMC_NISTR_BWRRDY
				| SDMMC_NISTR_BRDRDY
				| SDMMC_NISTR_DMAINT;
	error_status_mask = SDMMC_EISTR_ADMA
				| SDMMC_EISTR_CMDTEO
				| SDMMC_EISTR_CMDCRC
				| SDMMC_EISTR_CMDEND
				| SDMMC_EISTR_CMDIDX
				| SDMMC_EISTR_DATTEO
				| SDMMC_EISTR_DATCRC
				| SDMMC_EISTR_DATEND;

	sdhc_writew(SDMMC_NISTER, normal_status_mask);
	sdhc_writew(SDMMC_EISTER, error_status_mask);

	sdhc_writew(SDMMC_NISIER, 0);
	sdhc_writew(SDMMC_EISIER, 0);

	sdhc_set_clock(sdcard, 400000);
	sdhc_set_bus_width(sdcard, 1);

	return 0;
}

static int sdhc_read_data(struct sd_data *data)
{
	unsigned short normal_status, error_status;
	unsigned int psr;
	unsigned int timeout;
	unsigned int block = 0;
	unsigned short i;
	unsigned int *tmp;
	int done = 0;

	timeout = 100000;
	do {
		normal_status = sdhc_readw(SDMMC_NISTR);

		sdhc_writew(SDMMC_NISTR, normal_status);

		if ((normal_status & SDMMC_NISTR_ERRINT) == SDMMC_NISTR_ERRINT) {
			error_status = sdhc_readw(SDMMC_EISTR);

			sdhc_writew(SDMMC_EISTR, error_status);

			sdhc_software_reset_dat();

			dbg_info("SDHC: Error detected in status: %x, %x\n",
				normal_status, error_status);

			return -1;
		}
		if (((normal_status & SDMMC_NISTR_BRDRDY) ||
			(normal_status & SDMMC_NISTR_BWRRDY)) && !done ) {
			psr = sdhc_readl(SDMMC_PSR);
			if (data->direction == SD_DATA_DIR_RD &&
				!(psr & SDMMC_PSR_BUFRDEN))
			continue;
			if (data->direction == SD_DATA_DIR_WR &&
				!(psr & SDMMC_PSR_BUFWREN))
			continue;

			tmp = (unsigned int *)data->buff;
			for (i = 0; i < data->blocksize; i += 4) {
				if (data->direction == SD_DATA_DIR_RD)
					*tmp++ = sdhc_readl(SDMMC_BDPR);
				else
					sdhc_writel(SDMMC_BDPR, *tmp++);
			}

			data->buff += data->blocksize;
			if (++block >= data->blocks) {
				done = 1;
				goto sdhc_read_data_reset;
			}
		}

		if (timeout-- > 0) {
			udelay(10);
		} else {
			dbg_loud("SDHC: Transfer data timeout\n");
			goto sdhc_read_data_reset;
		}
	} while (!(normal_status & SDMMC_NISTR_TRFC) && !done);

	return 0;

sdhc_read_data_reset:
	/* There is an issue when writing data,
	 * even after the whole data has been written, the SDMMC will not
	 * trigger TRFC in some cases. Check for WTACT for that.
	 * In this case we use a sw workaround to reset the
	 * CMD and DAT lines from SRR and clear the interrupts
	 * Otherwise we are never getting inhibits bits back on
	 */
	if (data->direction == SD_DATA_DIR_WR &&
		(sdhc_readl(SDMMC_PSR) & SDMMC_PSR_WTACT)) {
		/* wait for BWRRDY to be cleared */
		unsigned int timeout2 = 10000;
		do {
			normal_status = sdhc_readw(SDMMC_NISTR);
		} while (timeout2-- && (normal_status & SDMMC_NISTR_BWRRDY));
	}

	error_status = sdhc_readw(SDMMC_EISTR);
	sdhc_software_reset_dat();
	sdhc_software_reset_cmd();
	sdhc_writew(SDMMC_EISTR, error_status);
	sdhc_writew(SDMMC_NISTR, normal_status);

	return -1;
}

static int sdhc_send_command(struct sd_command *sd_cmd, struct sd_data *data)
{
	unsigned int normal_status, error_status, normal_status_mask;
	unsigned int cmd_reg, mode;
	unsigned int i;
	int ret;
	unsigned int timeout;
	struct adma_desc dma_desc[16] = {0};

	timeout = 100000;
	while ((--timeout) &&
	       (sdhc_readl(SDMMC_PSR) & (SDMMC_PSR_CMDINHC | SDMMC_PSR_CMDINHD)))
		;

	if (!timeout)
		dbg_info("SDHC: Timeout waiting for CMD and DAT Inhibit bits\n");

	normal_status_mask =  SDMMC_NISTR_CMDC;

	cmd_reg = SDMMC_CR_CMDIDX_(sd_cmd->cmd);

	if ((sd_cmd->resp_type == SD_RESP_TYPE_R1)
		|| (sd_cmd->resp_type == SD_RESP_TYPE_R5)
		|| (sd_cmd->resp_type == SD_RESP_TYPE_R6)
		|| (sd_cmd->resp_type == SD_RESP_TYPE_R7)) {
		cmd_reg |= (SDMMC_CR_RESPTYP_RL48
				| SDMMC_CR_CMDCCEN
				| SDMMC_CR_CMDICEN);
	} else if (sd_cmd->resp_type == SD_RESP_TYPE_R1B) {
		cmd_reg |= (SDMMC_CR_RESPTYP_RL48BUSY
				| SDMMC_CR_CMDCCEN
				| SDMMC_CR_CMDICEN);
		normal_status_mask |= SDMMC_NISTR_TRFC;
	} else if (sd_cmd->resp_type == SD_RESP_TYPE_R2) {
		cmd_reg |= (SDMMC_CR_RESPTYP_RL136
				| SDMMC_CR_CMDCCEN);
	} else if ((sd_cmd->resp_type == SD_RESP_TYPE_R3)
		|| (sd_cmd->resp_type == SD_RESP_TYPE_R4)) {
		cmd_reg |= SDMMC_CR_RESPTYP_RL48;
	} else {
		cmd_reg |= SDMMC_CR_RESPTYP_NORESP;
	}

	if (data) {
		if (sdhc_host.caps_adma2 && (sd_cmd->cmd == SD_CMD_READ_SINGLE_BLOCK ||
		    sd_cmd->cmd == SD_CMD_READ_MULTIPLE_BLOCK)) {
			/* for CMD17 and CMD18 we use ADMA2 */
			sdhc_writeb(SDMMC_HC1R, sdhc_readb(SDMMC_HC1R) |
						SDMMC_HC1R_DMASEL_ADMA32);
		}

		cmd_reg |= SDMMC_CR_DPSEL;

		mode = SDMMC_TMR_BCEN;
		mode |= (data->blocks > 1) ? SDMMC_TMR_MSBSEL : 0;
		mode |= (data->direction == SD_DATA_DIR_RD) ? SDMMC_TMR_DTDSEL_READ : 0;

		if (sdhc_host.caps_adma2 && (sd_cmd->cmd == SD_CMD_READ_SINGLE_BLOCK ||
		    sd_cmd->cmd == SD_CMD_READ_MULTIPLE_BLOCK))
			mode |= SDMMC_TMR_DMAEN;

		sdhc_writeb(SDMMC_TCR, 0xe);
		sdhc_writew(SDMMC_BSR, data->blocksize);
		if (data->blocks > 1)
			sdhc_writew(SDMMC_BCR, data->blocks);

		sdhc_writew(SDMMC_TMR, mode);

		/* for CMD17 and CMD18 we use ADMA to transfer faster */
		if (sdhc_host.caps_adma2 && (sd_cmd->cmd == SD_CMD_READ_SINGLE_BLOCK ||
		    sd_cmd->cmd == SD_CMD_READ_MULTIPLE_BLOCK)) {
			/* prepare descriptor table */
			if (data->blocks > 16 )
				dbg_printf("too many blocks requested at once, error\n");

			for (i = 0; i < data->blocks; i++) {
				/* last descriptor must have the end bit */
				if (i == data->blocks - 1)
					dma_desc[i].cmd = 0x23;
				else
					dma_desc[i].cmd = 0x21;

				dma_desc[i].len = data->blocksize;
				dma_desc[i].addr = (unsigned int)(data->buff + data->blocksize * i);
			}
		/* address of the first descriptor goes here */
		sdhc_writel(SDMMC_ASAR0, (unsigned int) &dma_desc[0]);
		}
	}

	sdhc_writel(SDMMC_ARG1R, sd_cmd->argu);

	sdhc_writew(SDMMC_CR, cmd_reg);

	timeout = 100000;
	do {
		normal_status = sdhc_readw(SDMMC_NISTR);
	} while ((--timeout) &&
		 ((normal_status & normal_status_mask) != normal_status_mask));

	if (!timeout)
		dbg_very_loud("SDHC: Timeout waiting for command complete\n");

	/* clear the status, except for read and write ready.
	 * those will be cleared by the read/write data routine, which
	 * bases itself on the fact that the hardware is ready to receive data
	 * or has data ready to be read
	 */
	sdhc_writew(SDMMC_NISTR, normal_status & ~ (SDMMC_NISTR_BWRRDY | SDMMC_NISTR_BRDRDY));

	if ((normal_status & normal_status_mask) == normal_status_mask) {
		if (sd_cmd->resp_type == SD_RESP_TYPE_R2) {
			for (i = 0; i < 4; i++)
				sd_cmd->resp[i] = sdhc_readl(SDMMC_RR0 + 0x04 * i);
		} else  {
			*sd_cmd->resp = sdhc_readl(SDMMC_RR0);
		}

		/* if we have data but not using block transfer, we use PIO mode */
		if (data && (!sdhc_host.caps_adma2 || (sd_cmd->cmd != SD_CMD_READ_SINGLE_BLOCK &&
		    sd_cmd->cmd != SD_CMD_READ_MULTIPLE_BLOCK))) {
			sdhc_read_data(data);
		} else if (data && sdhc_host.caps_adma2) {
			/* otherwise, ADMA will carry the data for us */
			/* Let's wait for ADMA to finish transferring */
			timeout = 1000000;
			do {
				normal_status = sdhc_readw(SDMMC_NISTR);
				udelay(1);
			} while (--timeout && !(normal_status & SDMMC_NISTR_TRFC));
			if (!timeout)
				console_printf("SDHC: Timeout waiting for ADMA\n");

			sdhc_writew(SDMMC_NISTR, SDMMC_NISTR_TRFC);
			error_status = sdhc_readw(SDMMC_EISTR);
		}

		ret = 0;
	} else {
		error_status = sdhc_readw(SDMMC_EISTR);

		sdhc_software_reset_cmd();

		sdhc_writew(SDMMC_EISTR, error_status);

		if (error_status & SDMMC_EISTR_CMDTEO)
			return ERROR_TIMEOUT;

		ret = -1;
	}

	return ret;
}

static struct host_ops sdhc_ops = {
	.init = sdhc_init,
	.send_command = sdhc_send_command,
	.set_clock = sdhc_set_clock,
	.set_bus_width = sdhc_set_bus_width,
	.set_ddr = sdhc_set_ddr,
};

int sdcard_register_sdhc(struct sd_card *sdcard)
{
	sdcard->host = &sdhc_host;
	sdcard->host->ops = &sdhc_ops;

	return 0;
}
