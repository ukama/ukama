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
#ifndef __AT91SAM9G10_MATRIX_H__
#define __AT91SAM9G10_MATRIX_H__

#define MATRIX_MCFG	0x00	/* Master Configuration Register */
#define MATRIX_SCFG0	0x04	/* Slave Configuration Register 0 */
#define MATRIX_SCFG1	0x08	/* Slave Configuration Register 1 */
#define MATRIX_SCFG2	0x0c	/* Slave Configuration Register 2 */
#define MATRIX_SCFG3	0x10	/* Slave Configuration Register 3 */
#define MATRIX_SCFG4	0x14	/* Slave Configuration Register 4 */

/*-------- MATRIX_MCFG : (MATRIX Offset: 0x0) Master Configuration Register --------*/
#define AT91C_MATRIX_RCA926I	(0x1 << 0)	/* Remap Command for ARM926EJ-S Instruction Master */
#define AT91C_MATRIX_RCA926D	(0x1 << 1)	/* Remap Command for ARM926EJ-S Data Master */

/*-------- MATRIX_SCFG0 : (MATRIX Offset: 0x04) Slave Configuration Register 0 -------*/
/*-------- MATRIX_SCFG1 : (MATRIX Offset: 0x08) Slave Configuration Register 1 -------*/
/*-------- MATRIX_SCFG2 : (MATRIX Offset: 0x0c) Slave Configuration Register 2 -------*/
/*-------- MATRIX_SCFG3 : (MATRIX Offset: 0x10) Slave Configuration Register 3 -------*/
/*-------- MATRIX_SCFG4 : (MATRIX Offset: 0x14) Slave Configuration Register 4 -------*/
#define AT91C_MATRIX_SLOT_CYCLE		(0xFF << 0)	/* Maximum Number of Allowed Cycles for a Bursti */
#define 	AT91C_MATRIX_SLOT_CYCLE_(x)		(x << 0)
#define AT91C_MATRIX_DEFMSTR_TYPE	(0x3 << 16)	/* Default Master Type */
#define 	AT91C_MATRIX_DEFMSTR_TYPE_NO_DEFMSTR		(0x0 << 16)
#define 	AT91C_MATRIX_DEFMSTR_TYPE_LAST_DEFMSTR		(0x1 << 16)
#define 	AT91C_MATRIX_DEFMSTR_TYPE_FIXED_DEFMSTR		(0x2 << 16)
#define AT91C_MATRIX_FIXED_DEFMSTR	(0x7 << 18)	/* Fixed Index of Default Master */
#define 	AT91C_MATRIX_FIXED_DEFMSTR_ARM926I		(0x0 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_ARM926D		(0x1 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_HPDC3		(0x2 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_LCDC			(0x3 << 18)
#define 	AT91C_MATRIX_FIXED_DEFMSTR_UHP			(0x4 << 18)

#endif /* #ifndef __AT91SAM9G10_MATRIX_H__ */
