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
#ifndef __AT91_MATRIX_H__
#define __AT91_MATRIX_H__

#define MATRIX_MCFG0	0x00	/*  Master Configuration Register 0 */
#define MATRIX_MCFG1	0x04	/*  Master Configuration Register 1 */
#define MATRIX_MCFG2	0x08	/*  Master Configuration Register 2 */
#define MATRIX_MCFG3	0x0c	/*  Master Configuration Register 3 */
#define MATRIX_MCFG4	0x10	/*  Master Configuration Register 4 */
#define MATRIX_MCFG5	0x14	/*  Master Configuration Register 5 */
#define MATRIX_MCFG6	0x18	/*  Master Configuration Register 6 */
#define MATRIX_MCFG7	0x2c	/*  Master Configuration Register 7 */
#define MATRIX_MCFG8	0x20	/*  Master Configuration Register 8 */
/* 0x18 ~ 0x3c reserved */
#define MATRIX_SCFG0	0x40	/*  Slave Configuration Register 0 */
#define MATRIX_SCFG1	0x44	/*  Slave Configuration Register 1 */
#define MATRIX_SCFG2	0x48	/*  Slave Configuration Register 2 */
#define MATRIX_SCFG3	0x4c	/*  Slave Configuration Register 3 */
#define MATRIX_SCFG4	0x50	/*  Slave Configuration Register 4 */
#define MATRIX_SCFG5	0x54	/*  Slave Configuration Register 5 */
#define MATRIX_SCFG6	0x58	/*  Slave Configuration Register 6 */
/* 0x54 ~ 0x7c reserved */
#define MATRIX_PRAS0	0x80	/*  Priority Register A for Slave 0 */
#define MATRIX_PRBS0	0x84	/*  Priority Register B for Slave 0 */
#define MATRIX_PRAS1	0x88	/*  Priority Register A for Slave 1 */
#define MATRIX_PRBS1	0x8c	/*  Priority Register B for Slave 1 */
#define MATRIX_PRAS2	0x90	/*  Priority Register A for Slave 2 */
#define MATRIX_PRBS2	0x94	/*  Priority Register B for Slave 2 */
#define MATRIX_PRAS3	0x98	/*  Priority Register A for Slave 3 */
#define MATRIX_PRBS3	0x9c	/*  Priority Register B for Slave 3 */
#define MATRIX_PRAS4	0xa0	/*  Priority Register A for Slave 4 */
#define MATRIX_PRBS4	0xa4	/*  Priority Register B for Slave 4 */
#define MATRIX_PRAS5	0xa8	/*  Priority Register A for Slave 5 */
#define MATRIX_PRBS5	0xac	/*  Priority Register B for Slave 5 */
#define MATRIX_PRAS6	0xb0	/*  Priority Register A for Slave 6 */
#define MATRIX_PRBS6	0xb4	/*  Priority Register B for Slave 6 */

/*-------- MATRIX_MCFG0:(MATRIX Offset: 0x0) Master Configuration Register --------*/
#define AT91C_MATRIX_ULBT	(0x7 << 0)	/* (MATRIX) Undefined Length Burst Type */
#define		AT91C_MATRIX_ULBT_INFINITE		0
#define		AT91C_MATRIX_ULBT_SINGLE_ACCESS		1
#define		AT91C_MATRIX_ULBT_4_BEAT		2
#define		AT91C_MATRIX_ULBT_8_BEAT		3
#define		AT91C_MATRIX_ULBT_16_BEAT		4

/*-------- MATRIX_SCFG0 : (MATRIX Offset: 0x40) Slave Configuration Register 0 ----*/
#define AT91C_MATRIX_SLOT_CYCLE   	(0xff << 0)	/* Maximum Number of Allowed Cycles for a Burst */
#define 	AT91C_MATRIX_SLOT_CYCLE_(x)			(x << 0)
#define AT91C_MATRIX_DEFMSTR_TYPE 	(0x3 << 16)	/* Default Master Type */
#define 	AT91C_MATRIX_DEFMSTR_TYPE_NO_DEFMSTR		(0x0 << 16)
#define 	AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR		(0x1 << 16)
#define 	AT91C_MATRIX_DEFMSTR_TYPE_FIXED_DEFMSTR		(0x2 << 16)
#define AT91C_MATRIX_FIXED_DEFMSTR 	(0xf << 18)	/* Fixed Default Master */
#define 	AT91C_MATRIX_FIXED_DEFMSTR_ARM926I		(0x0 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_ARM926D		(0x1 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_PDC			(0x2 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_LCDC			(0x3 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_2DGC			(0x4 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_ISI			(0x5 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_DMA			(0x6 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_EMAC			(0x7 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_USB			(0x8 << 18)
#define AT91C_MATRIX_ARBT 		(0x3 << 24)	/* Arbitration Type */
#define 	AT91C_MATRIX_ARBT_ROUND_ROBIN	 		(0x0 << 24)
#define 	AT91C_MATRIX_ARBT_FIXED_PRIORITY		(0x1 << 24)

/*-------- MATRIX_PRAS0 : (MATRIX Offset: 0x80) PRAS0 Register --------*/
#define AT91C_MATRIX_M0PR	(0x3 << 0)	/* ARM926EJ-S Instruction priority */
#define 	AT91C_MATRIX_M0PR_(x)		(x << 0)
#define AT91C_MATRIX_M1PR	(0x3 << 4)	/* ARM926EJ-S Data priority */
#define 	AT91C_MATRIX_M1PR_(x)		(x << 4)
#define AT91C_MATRIX_M2PR	(0x3 << 8)	/* PDC priority */
#define 	AT91C_MATRIX_M2PR_(x)		(x << 8)
#define AT91C_MATRIX_M3PR	(0x3 << 12)	/* LCDC priority */
#define 	AT91C_MATRIX_M3PR_(x)		(x << 12)
#define AT91C_MATRIX_M4PR	(0x3 << 16)	/* 2DGC priority */
#define 	AT91C_MATRIX_M4PR_(x)		(x << 16)
#define AT91C_MATRIX_M5PR	(0x3 << 20)	/* ISI priority */
#define 	AT91C_MATRIX_M5PR_(x)		(x << 20)
#define AT91C_MATRIX_M6PR	(0x3 << 24)	/* DMA priority */
#define 	AT91C_MATRIX_M6PR_(x)		(x << 24)
#define AT91C_MATRIX_M7PR	(0x3 << 28)	/* EMAC priority */
#define 	AT91C_MATRIX_M7PR_(x)		(x << 28)

/*-------- MATRIX_PRBS0 : (MATRIX Offset: 0x84) PRBS0 Register --------*/
#define AT91C_MATRIX_M8PR	(0x3 << 0)	/* (MATRIX) USB priority */
#define 	AT91C_MATRIX_M8PR_(x)		(x << 0)

/*-------- MATRIX_MRCR : (MATRIX Offset: 0x100) MRCR Register --------*/
#define AT91C_MATRIX_RCA926I	(0x1 << 0)	/* Remap Command for ARM926EJ-S Instruction Master */
#define AT91C_MATRIX_RCA926D	(0x1 << 1)	/* Remap Command for ARM926EJ-S Data Master */
#define AT91C_MATRIX_RCB2	(0x1 << 2)	/* Remap Command Bit for PDC */
#define AT91C_MATRIX_RCB3	(0x1 << 3)	/* Remap Command Bit for LCD */
#define AT91C_MATRIX_RCB4	(0x1 << 4)	/* Remap Command Bit for 2DGC */
#define AT91C_MATRIX_RCB5	(0x1 << 5)	/* Remap Command Bit for ISI */
#define AT91C_MATRIX_RCB6	(0x1 << 6)	/* Remap Command Bit for DMA */
#define AT91C_MATRIX_RCB7	(0x1 << 7)	/* Remap Command Bit for EMAC */
#define AT91C_MATRIX_RCB8	(0x1 << 8)	/* Remap Command Bit for USB */

#endif	/* #ifndef __AT91_MATRIX_H__ */
