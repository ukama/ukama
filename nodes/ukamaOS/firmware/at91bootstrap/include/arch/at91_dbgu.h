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
#ifndef __AT91_DBGU_H__
#define __AT91_DBGU_H__

/* *** Register offset in AT91S_DBGU structure ***/
#define DBGU_CR		0x00	/* Control Register */
#define DBGU_MR		0x04	/* Mode Register */
#define DBGU_IER	0x08	/* Interrupt Enable Register */
#define DBGU_IDR	0x0C	/* Interrupt Disable Register */
#define DBGU_IMR	0x10	/* Interrupt Mask Register */
#define DBGU_CSR	0x14	/* Channel Status Register */
#define DBGU_RHR	0x18	/* Receiver Holding Register */
#define DBGU_THR	0x1C	/* Transmitter Holding Register */
#define DBGU_BRGR	0x20	/* Baud Rate Generator Register */
#define DBGU_CIDR	0x40	/* Chip ID Register */
#define DBGU_EXID	0x44	/* Chip ID Extension Register */
#define DBGU_FNTR	0x48	/* Force NTRST Register */
#define DBGU_ADDRSIZE	0xEC	/* DBGU ADDRSIZE REGISTER */
#define DBGU_IPNAME1	0xF0	/* DBGU IPNAME1 REGISTER */
#define DBGU_IPNAME2	0xF4	/* DBGU IPNAME2 REGISTER */
#define DBGU_FEATURES	0xF8	/* DBGU FEATURES REGISTER */
#define DBGU_VER	0xFC	/* DBGU VERSION REGISTER */

/* -------- DBGU_CR : (DBGU Offset: 0x0) Debug Unit Control Register --------*/ 
#define AT91C_DBGU_RSTRX	(0x1UL << 2)
#define AT91C_DBGU_RSTTX	(0x1UL << 3)
#define AT91C_DBGU_RXEN		(0x1UL << 4)
#define AT91C_DBGU_RXDIS	(0x1UL << 5)
#define AT91C_DBGU_TXEN		(0x1UL << 6)
#define AT91C_DBGU_TXDIS	(0x1UL << 7)
#define AT91C_DBGU_RSTSTA	(0x1UL << 8)

/* -------- DBGU_MR : (DBGU Offset: 0x4) Debug Unit Mode Register --------*/ 
#define AT91C_DBGU_NBSTOP	(0x3UL << 12)
#define		AT91C_DBGU_NBSTOP_1BIT			(0x0UL << 12)
#define		AT91C_DBGU_NBSTOP_1_5BIT		(0x1UL << 12)
#define		AT91C_DBGU_NBSTOP_2BIT			(0x2UL << 12)
#define AT91C_DBGU_CHRL		(0x3UL << 6)
#define		AT91C_DBGU_CHRL_5BIT			(0x0UL << 6)
#define		AT91C_DBGU_CHRL_6BIT			(0x1UL << 6)
#define		AT91C_DBGU_CHRL_7BIT			(0x2UL << 6)
#define		AT91C_DBGU_CHRL_8BIT			(0x3UL << 6)
#define AT91C_DBGU_PAR		(0x7UL << 9)
#define 	AT91C_DBGU_PAR_EVEN			(0x0UL << 9)
#define 	AT91C_DBGU_PAR_ODD			(0x1UL << 9)
#define 	AT91C_DBGU_PAR_SPACE			(0x2UL << 9)
#define 	AT91C_DBGU_PAR_MARK			(0x3UL << 9)
#define 	AT91C_DBGU_PAR_NONE			(0x4UL << 9)
#define AT91C_DBGU_CHMODE	(0x3UL << 14)
#define 	AT91C_DBGU_CHMODE_NORMAL		(0x0UL << 14)
#define 	AT91C_DBGU_CHMODE_AUTO			(0x1UL << 14)
#define 	AT91C_DBGU_CHMODE_LOCAL			(0x2UL << 14)
#define 	AT91C_DBGU_CHMODE_REMOTE		(0x3UL << 14)

/* -------- DBGU_IER : (DBGU Offset: 0x8) Debug Unit Interrupt Enable Register -------- */
#define AT91C_DBGU_RXRDY	(0x1UL <<  0)
#define AT91C_DBGU_TXRDY	(0x1UL <<  1)
#define AT91C_DBGU_ENDRX	(0x1UL <<  3)
#define AT91C_DBGU_ENDTX	(0x1UL <<  4)
#define AT91C_DBGU_OVRE		(0x1UL <<  5)
#define AT91C_DBGU_FRAME	(0x1UL <<  6)
#define AT91C_DBGU_PARE		(0x1UL <<  7)
#define AT91C_DBGU_TXEMPTY	(0x1UL <<  9)
#define AT91C_DBGU_TXBUFE	(0x1UL << 11)
#define AT91C_DBGU_RXBUFF	(0x1UL << 12)
#define AT91C_DBGU_COMM_TX	(0x1UL << 30)
#define AT91C_DBGU_COMM_RX	(0x1UL << 31)
/* -------- DBGU_IDR : (DBGU Offset: 0xc) Debug Unit Interrupt Disable Register -------- */
/* -------- DBGU_IMR : (DBGU Offset: 0x10) Debug Unit Interrupt Mask Register -------- */
/* -------- DBGU_CSR : (DBGU Offset: 0x14) Debug Unit Channel Status Register -------- */

/* -------- DBGU_FNTR : (DBGU Offset: 0x48) Debug Unit FORCE_NTRST Register --------*/ 
#define AT91C_DBGU_FORCE_NTRST	(0x1UL << 0)

#endif /* #ifndef __AT91_DBGU_H__ */
