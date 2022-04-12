/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2006, Atmel Corporation

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
#ifndef __SAMA5_SMC_H__
#define __SAMA5_SMC_H__

/**** Register offset in SAMA5 SMC structure ***/
#define SMC_SETUP0	0x00	/* Setup Register for CS 0 */
#define SMC_PULSE0	0x04	/* Pulse Register for CS 0 */
#define SMC_CYCLE0	0x08	/* Cycle Register for CS 0 */
#define SMC_TIMINGS0	0x0C	/* Timings Register for CS 0 */
#define SMC_MODE0	0x10	/* Mode Register for CS 0 */
#define SMC_SETUP1	0x14	/* Setup Register for CS 1 */
#define SMC_PULSE1	0x18	/* Pulse Register for CS 1 */
#define SMC_CYCLE1	0x1C	/* Cycle Register for CS 1 */
#define SMC_TIMINGS1	0x20	/* Timings Register for CS 1 */
#define SMC_MODE1	0x24	/* Mode Register for CS 1 */
#define SMC_SETUP2	0x28	/* Setup Register for CS 2 */
#define SMC_PULSE2	0x2C	/* Pulse Register for CS 2 */
#define SMC_CYCLE2	0x30	/* Cycle Register for CS 2 */
#define SMC_TIMINGS2	0x34	/* Timings Register for CS 2 */
#define SMC_MODE2	0x38	/* Mode Register for CS 2 */
#define SMC_SETUP3	0x3C	/* Setup Register for CS 3 */
#define SMC_PULSE3	0x40	/* Pulse Register for CS 3 */
#define SMC_CYCLE3	0x44	/* Cycle Register for CS 3 */
#define SMC_TIMINGS3	0x48	/* Timings Register for CS 3 */
#define SMC_MODE3	0x4C	/* Mode Register for CS 3 */

#define AT91C_SMC_SETUP_NWE(x)		(x & 0x3f)
#define AT91C_SMC_SETUP_NCS_WR(x)	((x & 0x3f) << 8)
#define AT91C_SMC_SETUP_NRD(x)		((x & 0x3f) << 16)
#define AT91C_SMC_SETUP_NCS_RD(x)	((x & 0x3f) << 24)

#define AT91C_SMC_PULSE_NWE(x)		(x & 0x3f)
#define AT91C_SMC_PULSE_NCS_WR(x)	((x & 0x3f) << 8)
#define AT91C_SMC_PULSE_NRD(x)		((x & 0x3f) << 16)
#define AT91C_SMC_PULSE_NCS_RD(x)	((x & 0x3f) << 24)

#define AT91C_SMC_CYCLE_NWE(x)		(x & 0x1ff)
#define AT91C_SMC_CYCLE_NRD(x)		((x & 0x1ff) << 16)

#define AT91C_SMC_TIMINGS_TCLR(x)	(x & 0xf)
#define AT91C_SMC_TIMINGS_TADL(x)	((x & 0xf) << 4)
#define AT91C_SMC_TIMINGS_TAR(x)		((x & 0xf) << 8)
#define AT91C_SMC_TIMINGS_OCMS		(0x01 << 12)
#define AT91C_SMC_TIMINGS_TRR(x)		((x & 0xf) << 16)
#define AT91C_SMC_TIMINGS_TWB(x)		((x & 0xf) << 24)
#define AT91C_SMC_TIMINGS_RBNSEL(x)	((x & 0x7) << 28)
#define AT91C_SMC_TIMINGS_NFSEL		(0x01 << 31)

#define AT91C_SMC_MODE_READMODE		(0x01 << 0)
#define 	AT91C_SMC_MODE_READMODE_NCS_CTRL	(0x00 << 0)
#define 	AT91C_SMC_MODE_READMODE_NRD_CTRL	(0x01 << 0)
#define AT91C_SMC_MODE_WRITEMODE		(0x01 << 1)
#define 	AT91C_SMC_MODE_WRITEMODE_NCS_CTRL	(0x00 << 1)
#define 	AT91C_SMC_MODE_WRITEMODE_NWE_CTRL	(0x01 << 1)
#define AT91C_SMC_MODE_EXNWMODE		(0x03 << 4)
#define 	AT91C_SMC_MODE_EXNWMODE_DISABLED	(0x00 << 4)
#define 	AT91C_SMC_MODE_EXNWMODE_FROZEN		(0x02 << 4)
#define 	AT91C_SMC_MODE_EXNWMODE_READY		(0x03 << 4)
#define AT91C_SMC_MODE_BAT		(0x01 << 8)
#define 	AT91C_SMC_MODE_BAT_SELECT	(0x00 << 8)
#define 	AT91C_SMC_MODE_BAT_WRITE		(0x01 << 8)
#define AT91C_SMC_MODE_DBW		(0x01 << 12)
#define 	AT91C_SMC_MODE_DBW_8		(0x00 << 12)
#define 	AT91C_SMC_MODE_DBW_16	(0x01 << 12)
#define AT91C_SMC_MODE_TDF_CYCLES(x)		((x & 0xf) << 16)
#define AT91C_SMC_MODE_TDF_MODE		(0x01 << 20)
#define 	AT91C_SMC_MODE_TDF_MODE_DISABLED		(0x00 << 20)
#define 	AT91C_SMC_MODE_TDF_MODE_ENABLED		(0x01 << 20)

#endif	/* #ifndef __SAMA5_SMC_H__ */
