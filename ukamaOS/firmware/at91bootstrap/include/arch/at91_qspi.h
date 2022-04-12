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
#ifndef	__AT91_QSPI_H__
#define	__AT91_QSPI_H__

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
#define	QSPI_WICR	0x34	/* Write Instruction Code Register */
#define	QSPI_IFR	0x38	/* Instruction Frame Register */
#define	QSPI_RICR	0x3c	/* Read Instruction Code Register */
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
#define	QSPI_CR_QSPIEN		(0x1 << 0)	/* QSPI Enable */
#define	QSPI_CR_QSPIDIS		(0x1 << 1)	/* QSPI Disable */
#define	QSPI_CR_SWRST		(0x1 << 7)	/* QSPI Software Reset */
#define	QSPI_CR_LASTXFER	(0x1 << 24)	/* Last Transfer */

/* QSPI_MR */
#define	QSPI_MR_SMM		(0x1 << 0)	/* Serial Memory Mode */
#define		QSPI_MR_SMM_SPI		(0x0 << 0)
#define		QSPI_MR_SMM_MEMORY	(0x1 << 0)
#define	QSPI_MR_LLB		(0x1 << 1)	/* Local Localback Enable */
#define		QSPI_MR_LLB_DISABLED	(0x0 << 1)
#define		QSPI_MR_LLB_ENABLED	(0x1 << 1)
#define	QSPI_MR_WDRBT		(0x1 << 2)	/* Wait Data Read Before Transfer */
#define		QSPI_MR_WDRBT_DISABLED		(0x0 << 2)
#define		QSPI_MR_WDRBT_ENABLED		(0x1 << 2)
#define	QSPI_MR_SMRM		(0x1 << 3)	/* Serial Memory Register Mode */
#define		QSPI_MR_SMRM_AHB	(0x0 << 3)
#define		QSPI_MR_SMRM_APB	(0x1 << 3)
#define	QSPI_MR_CSMODE		(0x3 << 4)	/* Chip Select Mode */
#define		QSPI_MR_CSMODE_NOT_RELOADED	(0x0 << 4)
#define		QSPI_MR_CSMODE_LASTXFER		(0x1 << 4)
#define		QSPI_MR_CSMODE_SYSTEMATICALLY	(0x2 << 4)
#define	QSPI_MR_NBBITS		(0xf << 8)	/* Number of Bits Per Transfer */
#define		QSPI_MR_NBBITS_8_BIT		(0x0 << 8)
#define		QSPI_MR_NBBITS_16_BIT		(0x8 << 8)
#define	QSPI_MR_DLYBCT		(0xff << 16)	/* Delay Between Consecutive Transfers */
#define	QSPI_MR_DLYCS		(0xff << 24)	/* Minimum Inactive QCS Delay */

/* QSPI_SR */
#define	QSPI_SR_RDRF		(0x1 << 0)	/* Receive Data Register Full */
#define	QSPI_SR_TDRE		(0x1 << 1)	/* Transmit Data Register Empty */
#define	QSPI_SR_TXEMPTY		(0x1 << 2)	/* Transmission Registers Empty */
#define	QSPI_SR_OVRES		(0x1 << 3)	/* Overrun Error Status */
#define	QSPI_SR_CSR		(0x1 << 8)	/* Chip Select Rise */
#define	QSPI_SR_CSS		(0x1 << 9)	/* Chip Select Status */
#define	QSPI_SR_INSTRE		(0x1 << 10)	/* Instruction End Status */
#define	QSPI_SR_QSPIENS		(0x1 << 24)	/* QSPI Enable Status */

/* QSPI_SCR */
#define	QSPI_SCR_CPOL		(0x1 << 0)	/* Clock Polarity */
#define	QSPI_SCR_CPHA		(0x1 << 1)	/* Clock Phase */
#define	QSPI_SCR_SCBR		(0xff << 8)
#define	QSPI_SCR_SCBR_(x)	(((x) << 8) & QSPI_SCR_SCBR)	/* Serial Clock Baud Rate */
#define	QSPI_SCR_DLYBS_(x)	((x) << 16)	/* Delay Before QSCK */

/* QSPI_ICR */
#define QSPI_ICR_INST_MASK     0xff
#define QSPI_ICR_INST(inst)    ((inst) & QSPI_ICR_INST_MASK)
#define QSPI_ICR_OPT_MASK      (0xff << 16)
#define QSPI_ICR_OPT(opt)      (((opt) << 16) & QSPI_ICR_OPT_MASK)

/* QSPI_IFR */
#define	QSPI_IFR_WIDTH		(0x7 << 0)	/* Width of Instruction Code, Address, Option Code and Data */
#define		QSPI_IFR_WIDTH_SINGLE_BIT_SPI	(0x0 << 0)
#define		QSPI_IFR_WIDTH_DUAL_OUTPUT	(0x1 << 0)
#define		QSPI_IFR_WIDTH_QUAD_OUTPUT	(0x2 << 0)
#define		QSPI_IFR_WIDTH_DUAL_IO		(0x3 << 0)
#define		QSPI_IFR_WIDTH_QUAD_IO		(0x4 << 0)
#define		QSPI_IFR_WIDTH_DUAL_CMD		(0x5 << 0)
#define		QSPI_IFR_WIDTH_QUAD_CMD		(0x6 << 0)
#define QSPI_IFR_WIDTH_(x)	(((x) << 0) & QSPI_IFR_WIDTH)
#define	QSPI_IFR_INSTEN		(0x1 << 4)	/* Instruction Enable*/
#define	QSPI_IFR_ADDREN		(0x1 << 5)	/* Address Enable*/
#define	QSPI_IFR_OPTEN		(0x1 << 6)	/* Option Enable*/
#define	QSPI_IFR_DATAEN		(0x1 << 7)	/* Data Enable*/
#define	QSPI_IFR_OPTL		(0x3 << 8)	/* Option Code Length */
#define		QSPI_IFR_OPTL_1BIT		(0x0 << 8)
#define		QSPI_IFR_OPTL_2BIT		(0x1 << 8)
#define		QSPI_IFR_OPTL_4BIT		(0x2 << 8)
#define		QSPI_IFR_OPTL_8BIT		(0x3 << 8)
#define	QSPI_IFR_ADDRL		(0x1 << 10)	/* Address Length */
#define		QSPI_IFR_ADDRL_24_BIT		(0x0 << 10)
#define		QSPI_IFR_ADDRL_32_BIT		(0x1 << 10)
#define	QSPI_IFR_TFRTYPE	(0x3 << 12)	/* Data Transfer Type */
#define		QSPI_IFR_TFRTYPE_READ		(0x0 << 12)
#define		QSPI_IFR_TFRTYPE_READ_MEMORY	(0x1 << 12)
#define		QSPI_IFR_TFRTYPE_WRITE		(0x2 << 12)
#define		QSPI_IFR_TFRTYPE_WRITE_MEMORY	(0x3 << 12)
#define QSPI_IFR_TFRTYPE_(x)	(((x) << 12) & QSPI_IFR_TFRTYPE)
#define	QSPI_IFR_CRM		(0x1 << 14)	/* Continuous Read Mode */
#define	QSPI_IFR_NBDUM_(x)	((x) << 16)	/* Number Of Dummy Cycles */

#endif
