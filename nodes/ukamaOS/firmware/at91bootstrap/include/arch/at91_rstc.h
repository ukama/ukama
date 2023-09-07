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
#ifndef __AT91_RSTC_H__
#define __AT91_RSTC_H__

/* *** Register offset in AT91S_RSTC structure ****/
#define RSTC_RCR	0x00	/* Reset Control Register */
#define RSTC_RSR	0x04	/* Reset Status Register */
#define RSTC_RMR	0x08	/* Reset Mode Register */
#define RSTC_GRSTR	0xE4	/* Generic Reset Register */
#define RSTC_VER	0xFC	/* Version Register */

/* -------- RSTC_RCR : (RSTC Offset: 0x0) Reset Control Register --------*/ 
#define AT91C_RSTC_PROCRST	(0x1UL << 0)
#define AT91C_RSTC_ICERST	(0x1UL << 1)
#define AT91C_RSTC_PERRST	(0x1UL << 2)
#define AT91C_RSTC_EXTRST	(0x1UL << 3)
#define AT91C_RSTC_RCRKEY_MASK	(0xFFUL << 24)
#define 	AT91C_RSTC_RCRKEY	(0xA5UL << 24)

/* -------- RSTC_RSR : (RSTC Offset: 0x4) Reset Status Register --------*/ 
#define AT91C_RSTC_URSTS	(0x1UL << 0)
#define AT91C_RSTC_RSTTYP	(0x7UL << 8)
#define 	AT91C_RSTC_RSTTYP_GENERAL	(0x0UL << 8)
#define 	AT91C_RSTC_RSTTYP_WAKEUP	(0x1UL << 8)
#define 	AT91C_RSTC_RSTTYP_WATCHDOG	(0x2UL << 8)
#define 	AT91C_RSTC_RSTTYP_SOFTWARE	(0x3UL << 8)
#define 	AT91C_RSTC_RSTTYP_USER		(0x4UL << 8)
#define AT91C_RSTC_NRSTL	(0x1UL << 16)
#define AT91C_RSTC_SRCMP	(0x1UL << 17)

/* -------- RSTC_RMR : (RSTC Offset: 0x8) Reset Mode Register --------*/ 
#define AT91C_RSTC_URSTEN	(0x1UL << 0)
#define AT91C_RSTC_URSTIEN	(0x1UL << 4)
#define AT91C_RSTC_ERSTL_MASK	(0xFUL << 8)
#define AT91C_RSTC_ERSTL(x)	((x) & AT91C_RSTC_ERSTL_MASK << 8)
#define AT91C_RSTC_RMRKEY_MASK	(0xFFUL << 24)
#define		AT91C_RSTC_KEY_UNLOCK		(0xA5UL << 24)

/* -------- RSTC_GRSTR : (RSTC Offset: 0xE4) Generic Reset Register --------*/
#define AT91C_GRSTR_DDR_RST	(0x1UL << 0)
#define AT91C_GRSTR_DDR_PHY_RST	(0x1UL << 2)

#endif /* #ifndef __AT91_RSTC_H__ */
