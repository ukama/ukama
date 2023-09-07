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
#ifndef __AT91_PIT_H__
#define __AT91_PIT_H__

/**** Register offset in AT91S_PITC structure ***/
#define PIT_MR		0x00	/* Period Interval Mode Register */
#define PIT_SR		0x04	/* Period Interval Status Register */
#define PIT_PIVR	0x08	/* Period Interval Value Register */
#define PIT_PIIR	0x0c	/* Period Interval Image Register */

/*-------- PITC_PIMR : (PITC Offset: 0x0) Periodic Interval Mode Register -------*/
#define AT91C_PIT_PIV		(0xFFFFF << 0)	/* Periodic Interval Value */
#define AT91C_PIT_PITEN	(0x01 << 24)	/* Periodic Interval Timer Enabled */
#define AT91C_PIT_PITIEN	(0x01 << 25)	/* Periodic Interval Timer Interrupt Enable */
/*-------- PITC_PISR : (PITC Offset: 0x4) Periodic Interval Status Register --------*/
#define AT91C_PIT_PITS		(0x01 << 0)	/* (PITC) Periodic Interval Timer Status */
/*-------- PITC_PIVR : (PITC Offset: 0x8) Periodic Interval Value Register --------*/
/*-------- PITC_PIIR : (PITC Offset: 0xc) Periodic Interval Image Register --------*/
#define AT91C_PIT_CPIV		(0xFFFFF << 0)	/* (PITC) Current Periodic Interval Value */
#define AT91C_PIT_PICNT	(0xFFF << 20)	/* (PITC) Periodic Interval Counter */

#endif /* #ifndef __AT91_PIT_H__ */
