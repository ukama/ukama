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
#ifndef __AT91_NAND_ECC_H__
#define __AT91_NAND_ECC_H__

/* PMECC Register Definitions */
#define PMECC_CFG	0x00	/* PMECC Configuration Register */
#define PMECC_SAREA	0x04	/* PMECC Spare Area Size Register */
#define PMECC_SADDR	0x08	/* PMECC Start Address Register */
#define PMECC_EADDR	0x0C	/* PMECC End Address Register */
#define PMECC_CLK	0x10	/* PMECC Clock Control Register */
#define PMECC_CTRL	0x14	/* PMECC Control Register */
#define PMECC_SR	0x18	/* PMECC Status Register */
#define PMECC_IER	0x1C	/* PMECC Interrupt Enable Register */
#define PMECC_IDR	0x20	/* PMECC Interrupt Disable Register */
#define PMECC_IMR	0x24	/* PMECC Interrupt Mask Register */
#define PMECC_ISR	0x28	/* PMECC Interrupt Status Register */
/* 0x2c ~ 0x3c */
#define PMECC_ECC	0x40	/* PMECC Redundancy structure */
#define PMECC_REM	0x240	/* PMECC Remainder structure */

/* -------- PMECC_CFG: (Offset: 0x0) PMECC Configuration Register --------*/ 
#define AT91_PMECC_BCH_ERR	(0x7UL << 0)	/* Error Correct Capability */
#define 	AT91C_PMECC_BCH_ERR2		(0x0UL)
#define 	AT91C_PMECC_BCH_ERR4		(0x1UL)
#define 	AT91C_PMECC_BCH_ERR8		(0x2UL)
#define 	AT91C_PMECC_BCH_ERR12		(0x3UL)
#define 	AT91C_PMECC_BCH_ERR24		(0x4UL)
#define 	AT91C_PMECC_BCH_ERR32		(0x5UL)
#define AT91C_PMECC_SECTORSZ	(0x1UL <<  4)	/* Sector Size */
#define 	AT91C_PMECC_SECTORSZ_512	(0x0UL << 4)
#define 	AT91C_PMECC_SECTORSZ_1024	(0x1UL << 4)
#define AT91C_PMECC_PAGESIZE	(0x3UL <<  8)	/* Number of Sectors in the Page */
#define 	AT91C_PMECC_PAGESIZE_1SEC	(0x0UL << 8)
#define 	AT91C_PMECC_PAGESIZE_2SEC	(0x1UL << 8)
#define 	AT91C_PMECC_PAGESIZE_4SEC	(0x2UL << 8)
#define 	AT91C_PMECC_PAGESIZE_8SEC	(0x3UL << 8)
#define AT91C_PMECC_NANDWR	(0x1UL << 12)	/* Nand Write Access */
#define 	AT91C_PMECC_NANDWR_0		(0x0UL << 12) 
#define 	AT91C_PMECC_NANDWR_1		(0x1UL << 12)
#define AT91C_PMECC_SPAREENA	(0x1UL << 16)	/* Spare Enable */
#define 	AT91C_PMECC_SPAREENA_DIS	(0x0UL << 16)
#define 	AT91C_PMECC_SPAREENA_ENA	(0x1UL << 16)
#define AT91C_PMECC_AUTO	(0x1UL << 20)	/* Automatic Mode Enable */
#define 	AT91C_PMECC_AUTO_DIS		(0x0UL << 20)
#define 	AT91C_PMECC_AUTO_ENA		(0x1UL << 20)

/* --------PMECC_SAREA: (Offset: 0x4) PMECC Spare Area Size Register -------- */
#define AT91C_PMECC_SPARESIZE		(0x1FFUL <<  0)

/* -------- PMECC_SADDR: (Offset: 0x8) PMECC Start Address Register -------- */
#define AT91C_PMECC_STARTADDR		(0x1FFUL <<  0)

/* -------- PMECC_EADDR: (Offset: 0xc) PMECC End Address Register -------- */ 
#define AT91C_PMECC_ENDADDR		(0x1FFUL <<  0)

/* -------- PMECC_CLK: (Offset: 0x10) PMECC Clock Control Register -------- */
#define AT91C_PMECC_CLKCTRL		(0x7UL <<  0)

/* -------- PMECC_CTRL: (Offset: 0x14) PMECC Control Register --------*/ 
#define AT91C_PMECC_RST		(0x1UL << 0)	/* (BCH) Reset the PMECC module */
#define AT91C_PMECC_DATA	(0x1UL << 1)	/* (BCH) Start a data phase */
#define AT91C_PMECC_USER	(0x1UL << 2)	/* (BCH) Start a user mode phase */
#define AT91C_PMECC_ENABLE	(0x1UL << 4)	/* (BCH) Enable BCH function */
#define AT91C_PMECC_DISABLE	(0x1UL << 5)	/* (BCH) Disable BCH function */

/* -------- PMECC_SR: (Offset: 0x18) PMECC Status Register --------*/ 
#define AT91C_PMECC_BUSY	(0x1UL << 0)	/* (BCH) The Kernel of the PMECC is busy */

/* -------- PMECC_IER: (Offset: 0x1c) PMECC Interrupt Enable Register --------*/ 
#define AT91C_PMECC_ERRIE	(0x1UL << 0)	/* (BCH) Error interrupt enable */

/* -------- PMECC_IDR: (Offset: 0x20) PMECC Interrupt Disable Register --------*/ 
#define AT91C_PMECC_ERRID	(0x1UL << 0)	/* (BCH) Error interrupt disable */

/* -------- PMECC_IMR: (Offset: 0x24) PMECC Interrupt Mask Register --------*/ 
#define AT91C_PMECC_ERRIM	(0x1UL << 0)	/* (BCH) Error interrupt Mask */

/* -------- PMECC_ISR: (Offset: 0x28) PMECC Interrupt Status Register -------- */
#define AT91C_PMECC_ERRIS	(0xFFUL << 0)	/* (BCH) Error interrupt Status */

/**** Register offset in AT91C_BCHEL structure ****/
/* PMERRLOC Register Definitions */
#define PMERRLOC_ELCFG		0x000	/* Error Location Configuration Register */
#define PMERRLOC_ELPRIM		0x004	/* Error Location Primitive Register */
#define PMERRLOC_ELEN		0x008	/* Error Location Enable Register */
#define PMERRLOC_ELDIS		0x00C	/* Error Location Disable Register */
#define PMERRLOC_ELSR		0x010	/* Error Location Status Register */
#define PMERRLOC_ELIER		0x014	/* Error Location Interrupt Enable Register */
#define PMERRLOC_ELIDR		0x018	/* Error Location Interrupt Disable Register */
#define PMERRLOC_ELIMR		0x01C	/* Error Location Interrupt Mask Register */

#define PMERRLOC_ELISR		0x020	/* Error Location Interrupt Status Register */
/* -------- PMERRLOC_ELISR: (Offset: 0x20) Error Location Interrupt Status Register --------*/
#define PMERRLOC_ELISR_DONE		(0x1)
#ifdef CONFIG_SAMA5D2
#define PMERRLOC_ELISR_ERR_CNT		(0x3f << 8)
#else
#define PMERRLOC_ELISR_ERR_CNT		(0x1f << 8)
#endif

/* 0x24 reserved */
#define PMERRLOC_SIGMA0		0x028	/* PMECC Error Location SIGMA0 Register */

#ifdef CONFIG_SAMA5D2
#define PMERRLOC_EL0		0x0AC	/* PMECC Error Location 0 Register */
#else
#define PMERRLOC_EL0		0x08C	/* PMECC Error Location 0 Register */
#endif

#define	PMERRLOC_VERSION	0x1FC	/* PMECC Version Register */
#define		AT91C_PMECC_VERSION_SAMA5D4	0x113
#define		AT91C_PMECC_VERSION_SAMA5D3	0x112
#define		AT91C_PMECC_VERSION_AT91SAM9N12	0x102
#define		AT91C_PMECC_VERSION_AT91SAM9X5	0x101

#endif /* #ifndef __AT91_NAND_ECC_H__  */
