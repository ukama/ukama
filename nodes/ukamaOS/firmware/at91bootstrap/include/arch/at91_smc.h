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
#ifndef __AT91_SMC_H__
#define __AT91_SMC_H__

/**** Register offset in AT91C SMC structure ***/
#define SMC_SETUP0	0x00	/* Setup Register for CS 0 */
#define SMC_PULSE0	0x04	/* Pulse Register for CS 0 */
#define SMC_CYCLE0	0x08	/* Cycle Register for CS 0 */
#define SMC_CTRL0	0x0C	/* Control Register for CS 0 */
#define SMC_SETUP1	0x10	/* Setup Register for CS 1 */
#define SMC_PULSE1	0x14	/* Pulse Register for CS 1 */
#define SMC_CYCLE1	0x18	/* Cycle Register for CS 1 */
#define SMC_CTRL1	0x1C	/* Control Register for CS 1 */
#define SMC_SETUP2	0x20	/* Setup Register for CS 2 */
#define SMC_PULSE2	0x24	/* Pulse Register for CS 2 */
#define SMC_CYCLE2	0x28	/* Cycle Register for CS 2 */
#define SMC_CTRL2	0x2C	/* Control Register for CS 2 */
#define SMC_SETUP3	0x30	/* Setup Register for CS 3 */
#define SMC_PULSE3	0x34	/* Pulse Register for CS 3 */
#define SMC_CYCLE3	0x38	/* Cycle Register for CS 3 */
#define SMC_CTRL3	0x3C	/* Control Register for CS 3 */
#define SMC_SETUP4	0x40	/* Setup Register for CS 4 */
#define SMC_PULSE4	0x44	/* Pulse Register for CS 4 */
#define SMC_CYCLE4	0x48	/* Cycle Register for CS 4 */
#define SMC_CTRL4	0x4C	/* Control Register for CS 4 */
#define SMC_SETUP5	0x50	/* Setup Register for CS 5 */
#define SMC_PULSE5	0x54	/* Pulse Register for CS 5 */
#define SMC_CYCLE5	0x58	/* Cycle Register for CS 5 */
#define SMC_CTRL5	0x5C	/* Control Register for CS 5 */
#define SMC_SETUP6	0x60	/* Setup Register for CS 6 */
#define SMC_PULSE6	0x64	/* Pulse Register for CS 6 */
#define SMC_CYCLE6	0x68	/* Cycle Register for CS 6 */
#define SMC_CTRL6	0x6C	/* Control Register for CS 6 */
#define SMC_SETUP7	0x70	/* Setup Register for CS 7 */
#define SMC_PULSE7	0x74	/* Pulse Register for CS 7 */
#define SMC_CYCLE7	0x78	/* Cycle Register for CS 7 */
#define SMC_CTRL7	0x7C	/* Control Register for CS 7 */
#define SMC_DELAY1	0xC0	/* SMC Delay Control Register */
#define SMC_DELAY2	0xC4	/* SMC Delay Control Register */
#define SMC_DELAY3	0xC8	/* SMC Delay Control Register */
#define SMC_DELAY4	0xCC	/* SMC Delay Control Register */
#define SMC_DELAY5	0xD0	/* SMC Delay Control Register */
#define SMC_DELAY6	0xD4	/* SMC Delay Control Register */
#define SMC_DELAY7	0xD8	/* SMC Delay Control Register */
#define SMC_DELAY8	0xDC	/* SMC Delay Control Register */

/* -------- SMC_SETUP : (SMC Offset: 0x0) Setup Register for CS x --------*/ 
#define AT91C_SMC_NWESETUP	(0x3FUL << 0)
#define		AT91C_SMC_NWESETUP_(x)		((x) << 0)
#define AT91C_SMC_NCSSETUPWR	(0x3FUL << 8)
#define		AT91C_SMC_NCS_WRSETUP_(x)	((x) << 8)
#define AT91C_SMC_NRDSETUP	(0x3FUL << 16)
#define		AT91C_SMC_NRDSETUP_(x)		((x) << 16)
#define AT91C_SMC_NCSSETUPRD	(0x3FUL << 24)
#define		AT91C_SMC_NCS_RDSETUP_(x)	((x) << 24)

/* -------- SMC_PULSE : (SMC Offset: 0x4) Pulse Register for CS x --------*/ 
#define AT91C_SMC_NWEPULSE	(0x7FUL << 0)
#define		AT91C_SMC_NWEPULSE_(x)		((x) << 0)
#define AT91C_SMC_NCSPULSEWR	(0x7FUL << 8)
#define		AT91C_SMC_NCS_WRPULSE_(x)	((x) << 8)
#define AT91C_SMC_NRDPULSE	(0x7FUL << 16)
#define		AT91C_SMC_NRDPULSE_(x)		((x) << 16)
#define AT91C_SMC_NCSPULSERD	(0x7FUL << 24)
#define		AT91C_SMC_NCS_RDPULSE_(x)	((x) << 24)

/* -------- SMC_CYC : (SMC Offset: 0x8) Cycle Register for CS x --------*/ 
#define AT91C_SMC_NWECYCLE	(0x1FFUL << 0)
#define		AT91C_SMC_NWECYCLE_(x)		((x) << 0)
#define AT91C_SMC_NRDCYCLE	(0x1FFUL << 16)
#define		AT91C_SMC_NRDCYCLE_(x)		((x) << 16)

/* -------- SMC_CTRL : (SMC Offset: 0xc) Control Register for CS x --------*/ 
#define AT91C_SMC_READMODE	(0x1UL << 0)
#define AT91C_SMC_WRITEMODE	(0x1UL << 1)
#define AT91C_SMC_NWAITM	(0x3UL << 4)
#define 	AT91C_SMC_NWAITM_NWAIT_DISABLE		(0x0UL << 4)
#define 	AT91C_SMC_NWAITM_NWAIT_ENABLE_FROZEN	(0x2UL << 4)
#define 	AT91C_SMC_NWAITM_NWAIT_ENABLE_READY	(0x3UL << 4)
#define AT91C_SMC_BAT		(0x1UL << 8)
#define 	AT91C_SMC_BAT_BYTE_SELECT	(0x0UL << 8)
#define 	AT91C_SMC_BAT_BYTE_WRITE	(0x1UL << 8)
#define AT91C_SMC_DBW		(0x3UL << 12)
#define 	AT91C_SMC_DBW_WIDTH_BITS_8	(0x0UL << 12)
#define 	AT91C_SMC_DBW_WIDTH_BITS_16	(0x1UL << 12)
#define 	AT91C_SMC_DBW_WIDTH_BITS_32	(0x2UL << 12)
#define AT91C_SMC_TDF		(0xFUL << 16)
#define		AT91_SMC_TDF_(x)		((x) << 16)
#define AT91C_SMC_TDFEN		(0x1UL << 20)
#define AT91C_SMC_PMEN		(0x1UL << 24)
#define AT91C_SMC_PS		(0x3UL << 28)
#define 	AT91C_SMC_PS_SIZE_FOUR_BYTES		(0x0UL << 28)
#define 	AT91C_SMC_PS_SIZE_EIGHT_BYTES		(0x1UL << 28)
#define 	AT91C_SMC_PS_SIZE_SIXTEEN_BYTES		(0x2UL << 28)
#define 	AT91C_SMC_PS_SIZE_THIRTY_TWO_BYTES	(0x3UL << 28)

#endif /* #ifndef __AT91_SMC_H__ */
