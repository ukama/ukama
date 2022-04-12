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
#ifndef __AT91_EEFC_H__
#define __AT91_EEFC_H__

#define EEFC_FMR		0x00	/* EEFC Flash Mode Register */
#define EEFC_FCR		0x04	/* EEFC Flash Command Register */
#define EEFC_FSR		0x08	/* EEFC Flash Status Register */
#define EEFC_FRR		0x0c	/* EEFC Flash Result Register */

/*-------- EEFC_FMR : (EEFC Offset: 0x0) EEFC Flash Mode Register --------*/ 
#define AT91C_EEFC_FRDY		(0x1 << 0)	/* Ready Interrupt Enable */
#define AT91C_EEFC_FWS		(0xF << 8)	/* Flash Wait State */
#define 	AT91C_EEFC_FWS_0WS		(0x0 << 8)
#define 	AT91C_EEFC_FWS_1WS		(0x1 << 8)
#define 	AT91C_EEFC_FWS_2WS		(0x2 << 8)
#define 	AT91C_EEFC_FWS_3WS		(0x3 << 8)
#define 	AT91C_EEFC_FWS_4WS		(0x4 << 8)
#define 	AT91C_EEFC_FWS_5WS		(0x5 << 8)
#define 	AT91C_EEFC_FWS_6WS		(0x6 << 8)

/*-------- EEFC_FCR : (EEFC Offset: 0x4) EEFC Flash Command Register --------*/ 
#define AT91C_EEFC_FCMD		(0xFF << 0)	/* Flash Command */
#define 	AT91C_EEFC_FCMD_GETD		(0x0)
#define 	AT91C_EEFC_FCMD_WP		(0x1)
#define 	AT91C_EEFC_FCMD_WPL		(0x2)
#define 	AT91C_EEFC_FCMD_EWP		(0x3)
#define 	AT91C_EEFC_FCMD_EWPL		(0x4)
#define 	AT91C_EEFC_FCMD_EA		(0x5)
#define 	AT91C_EEFC_FCMD_EPL		(0x6)
#define 	AT91C_EEFC_FCMD_EPA		(0x7)
#define 	AT91C_EEFC_FCMD_SLB		(0x8)
#define 	AT91C_EEFC_FCMD_CLB		(0x9)
#define 	AT91C_EEFC_FCMD_GLB		(0xA)
#define 	AT91C_EEFC_FCMD_SFB		(0xB)
#define 	AT91C_EEFC_FCMD_CFB		(0xC)
#define 	AT91C_EEFC_FCMD_GFB		(0xD)
#define AT91C_EEFC_FARG		(0xFFFF << 8)	/* Flash Command Argument */
#define AT91C_EEFC_FKEY		(0xFF << 24)	/* Flash Writing Protection Key */

/*-------- EEFC_FSR : (EEFC Offset: 0x8) EEFC Flash Status Register --------*/ 
#define AT91C_EEFC_FRDYS	(0x1 << 0)	/* Flash Ready Status */
#define AT91C_EEFC_FCMDE	(0x1 << 1)	/* Flash Command Error Status */
#define AT91C_EEFC_LOCKE	(0x1 << 2)	/* Flash Lock Error Status */

/*-------- EEFC_FRR : (EEFC Offset: 0xc) EEFC Flash Result Register --------*/ 
#define AT91C_EEFC_FVALUE	(0x0 << 0)	/* Flash Result Value */

#endif /* #ifndef __AT91_EEEFC_H__ */
