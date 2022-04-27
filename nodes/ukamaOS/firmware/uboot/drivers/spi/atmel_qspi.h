/*
 * Copyright (C) 2016
 *
 * SPDX-License-Identifier:	GPL-2.0+
 */

#ifndef __ATMEL_QSPI_H__
#define __ATMEL_QSPI_H__

/*
 * Register Definitions
 */
#define	QSPI_CR		0x00	/* Control Register */
#define	QSPI_MR		0x04	/* Mode Register */
#define	QSPI_RDR	0x08	/* Receive Data Register */
#define	QSPI_TDR	0x0c	/* Transmit Data Register */
#define	QSPI_SR		0x10	/* Status Register */
#define	QSPI_IER	0x14	/* Interrupt Enable Register */
#define	QSPI_IDR	0x18	/* Interrupt Disable Register */
#define	QSPI_IMR	0x1c	/* Interrupt Mask Register */
#define	QSPI_SCR	0x20	/* Serial Clock Register */
#define	QSPI_IAR	0x30	/* Instruction Address Register */
#define	QSPI_ICR	0x34	/* Instruction Code Register */
#define	QSPI_IFR	0x38	/* Instruction Frame Register */
/* 0x3c Reserved */
#define	QSPI_SMR	0x40	/* Scrambling Mode Register */
#define	QSPI_SKR	0x44	/* Scrambling Key Register */
/* 0x48 ~ 0xe0 */
#define	QSPI_WPMR	0xe4	/* Write Protection Mode Register */
#define	QSPI_WPSR	0xe8	/* Write Protection Status Register */
/* 0xec ~ 0xf8 Reserved */
/* 0xfc Reserved */

/*
 * Register Field Definitions
 */
/* QSPI_CR */
#define	QSPI_CR_QSPIEN		BIT(0)	/* QSPI Enable */
#define	QSPI_CR_QSPIDIS		BIT(1)	/* QSPI Disable */
#define	QSPI_CR_SWRST		BIT(7)	/* QSPI Software Reset */
#define	QSPI_CR_LASTXFER	BIT(24)	/* Last Transfer */

/* QSPI_MR */
#define	QSPI_MR_SMM		BIT(0)	/* Serial Memort Mode */
#define		QSPI_MR_SMM_SPI		0
#define		QSPI_MR_SMM_MEMORY	QSPI_MR_SMM
#define	QSPI_MR_LLB		BIT(1)	/* Local Localback Enable */
#define		QSPI_MR_LLB_DISABLED	0
#define		QSPI_MR_LLB_ENABLED	QSPI_MR_LLB
#define	QSPI_MR_WDRBT		BIT(2)	/* Wait Data Read Before Transfer */
#define		QSPI_MR_WDRBT_DISABLED	0
#define		QSPI_MR_WDRBT_ENABLED	QSPI_MR_WDRBT
#define	QSPI_MR_SMRM		BIT(3)	/* Serial Memory Register Mode */
#define		QSPI_MR_SMRM_AHB	0
#define		QSPI_MR_SMRM_APB	QSPI_MR_SMRM
#define	QSPI_MR_CSMODE		GENMASK(5, 4)	/* Chip Select Mode */
#define		QSPI_MR_CSMODE_NOT_RELOADED	(0x0u << 4)
#define		QSPI_MR_CSMODE_LASTXFER		(0x1u << 4)
#define		QSPI_MR_CSMODE_SYSTEMATICALLY	(0x2u << 4)
#define	QSPI_MR_NBBITS		GENMASK(11, 8)	/*
						 * Number of Bits Per
						 * Transfer
						 */
#define		QSPI_MR_NBBITS_8_BIT		(0x0u << 8)
#define		QSPI_MR_NBBITS_16_BIT		(0x8u << 8)
#define	QSPI_MR_DLYBCT		GENMASK(23, 16)	/*
						 * Delay Between Consecutive
						 * Transfers
						 */
#define	QSPI_MR_DLYCS		GENMASK(31, 24)	/* Minimum Inactive QCS Delay */

/* QSPI_SR */
#define	QSPI_SR_RDRF		BIT(0)	/* Receive Data Register Full */
#define	QSPI_SR_TDRE		BIT(1)	/* Transmit Data Register Empty */
#define	QSPI_SR_TXEMPTY		BIT(2)	/* Transmission Registers Empty */
#define	QSPI_SR_OVRES		BIT(3)	/* Overrun Error Status */
#define	QSPI_SR_CSR		BIT(8)	/* Chip Select Rise */
#define	QSPI_SR_CSS		BIT(9)	/* Chip Select Status */
#define	QSPI_SR_INSTRE		BIT(10)	/* Instruction End Status */
#define	QSPI_SR_QSPIENS		BIT(24)	/* QSPI Enable Status */

/* QSPI_SCR */
#define	QSPI_SCR_CPOL		BIT(0)	/* Clock Polarity */
#define	QSPI_SCR_CPOL_(x)	((x) << 0)
#define	QSPI_SCR_CPHA		BIT(1)	/* Clock Phase */
#define	QSPI_SCR_CPHA_(x)	((x) << 1)
#define	QSPI_SCR_SCBR		GENMASK(15, 8)	/* Serial Clock Baud Rate */
#define	QSPI_SCR_SCBR_(x)	(((x) << 8) & QSPI_SCR_SCBR)
#define QSPI_SCR_DLYBS		GENMASK(23, 16)
#define	QSPI_SCR_DLYBS_(x)	(((x) << 16) & QSPI_SCR_DLYBS)	/*
								 * Delay Before
								 * QSCK
								 */

/* QSPI_ICR */
#define QSPI_ICR_INST		GENMASK(7, 0)
#define	QSPI_ICR_INST_(x)	(((x) << 0) & QSPI_ICR_INST)	/*
								 * Instruction
								 * Code
								 */
#define QSPI_ICR_OPT		GENMASK(23, 16)
#define	QSPI_ICR_OPT_(x)	(((x) << 16) & QSPI_ICR_OPT)	/*
								 * Option
								 * Code
								 */

/* QSPI_IFR */
#define	QSPI_IFR_WIDTH		GENMASK(2, 0)	/*
						 * Width of Instruction Code,
						 * Address, Option Code and Data
						 */
#define		QSPI_IFR_WIDTH_SINGLE_BIT_SPI	(0x0u << 0)
#define		QSPI_IFR_WIDTH_DUAL_OUTPUT	(0x1u << 0)
#define		QSPI_IFR_WIDTH_QUAD_OUTPUT	(0x2u << 0)
#define		QSPI_IFR_WIDTH_DUAL_IO		(0x3u << 0)
#define		QSPI_IFR_WIDTH_QUAD_IO		(0x4u << 0)
#define		QSPI_IFR_WIDTH_DUAL_CMD		(0x5u << 0)
#define		QSPI_IFR_WIDTH_QUAD_CMD		(0x6u << 0)
#define QSPI_IFR_WIDTH_(x)	(((x) << 0) & QSPI_IFR_WIDTH)
#define	QSPI_IFR_INSTEN		BIT(4)	/* Instruction Enable*/
#define	QSPI_IFR_ADDREN		BIT(5)	/* Address Enable*/
#define	QSPI_IFR_OPTEN		BIT(6)	/* Option Enable*/
#define	QSPI_IFR_DATAEN		BIT(7)	/* Data Enable*/
#define	QSPI_IFR_OPTL		GENMASK(9, 8)	/* Option Code Length */
#define		QSPI_IFR_OPTL_1BIT		(0x0u << 8)
#define		QSPI_IFR_OPTL_2BIT		(0x1u << 8)
#define		QSPI_IFR_OPTL_4BIT		(0x2u << 8)
#define		QSPI_IFR_OPTL_8BIT		(0x3u << 8)
#define	QSPI_IFR_ADDRL		BIT(10)	/* Address Length */
#define		QSPI_IFR_ADDRL_24_BIT		0
#define		QSPI_IFR_ADDRL_32_BIT		QSPI_IFR_ADDRL
#define	QSPI_IFR_TFRTYPE	GENMASK(13, 12)	/* Data Transfer Type */
#define		QSPI_IFR_TFRTYPE_READ		(0x0u << 12)
#define		QSPI_IFR_TFRTYPE_READ_MEMORY	(0x1u << 12)
#define		QSPI_IFR_TFRTYPE_WRITE		(0x2u << 12)
#define		QSPI_IFR_TFRTYPE_WRITE_MEMORY	(0x3u << 12)
#define QSPI_IFR_TFRTYPE_(x)	(((x) << 12) & QSPI_IFR_TFRTYPE)
#define	QSPI_IFR_CRM		BIT(14)	/* Continuous Read Mode */
#define QSPI_IFR_NBDUM		GENMASK(20, 16)
#define	QSPI_IFR_NBDUM_(x)	(((x) << 16) & QSPI_IFR_NBDUM)	/*
								 * Number Of
								 * Dummy Cycles
								 */


struct atmel_qspi_platdata {
	void		*regbase;
	void		*membase;
};

struct atmel_qspi_priv {
	ulong		bus_clk_rate;
	void		*regbase;
	void		*membase;
};

#include <asm/io.h>

static inline u32 qspi_readl(struct atmel_qspi_priv *aq, u32 reg)
{
	return readl(aq->regbase + reg);
}

static inline void qspi_writel(struct atmel_qspi_priv *aq, u32 reg, u32 value)
{
	writel(value, aq->regbase + reg);
}

#endif /* __ATMEL_QSPI_H__ */
