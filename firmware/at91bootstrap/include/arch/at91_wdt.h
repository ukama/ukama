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
#ifndef __AT91_WDT_H__
#define __AT91_WDT_H__

/**** Register offset in AT91S_WDTC structure ***/
#define WDTC_CR		0x00	/* Watchdog Control Register */
#define WDTC_MR		0x04	/* Watchdog Mode Register */
#define WDTC_SR		0x08	/* Watchdog Status Register */

/*-------- WDTC_CR : (WDTC Offset: 0x0) Periodic Interval Image Register --------*/
#define AT91C_WDTC_WDRSTT	(0x1UL << 0)
#define AT91C_WDTC_KEY		(0xA5 << 24)

/*-------- WDTC_WDMR : (WDTC Offset: 0x4) Watchdog Mode Register --------*/
#define AT91C_WDTC_WDV		(0xFFFUL << 0)
#define AT91C_WDTC_WDFIEN	(0x1UL << 12)
#define AT91C_WDTC_WDRSTEN	(0x1UL << 13)
#define AT91C_WDTC_WDRPROC	(0x1UL << 14)
#define AT91C_WDTC_WDDIS	(0x1UL << 15)
#define AT91C_WDTC_WDD		(0xFFFUL << 16)
#define AT91C_WDTC_WDDBGHLT	(0x1UL << 28)
#define AT91C_WDTC_WDIDLEHLT	(0x1UL << 29)

/*-------- WDTC_WDSR : (WDTC Offset: 0x8) Watchdog Status Register -------*/
#define AT91C_WDTC_WDUNF	(0x1UL << 0)
#define AT91C_WDTC_WDERR	(0x1UL << 1)

#endif /* #ifndef __AT91_WDT_H__ */
