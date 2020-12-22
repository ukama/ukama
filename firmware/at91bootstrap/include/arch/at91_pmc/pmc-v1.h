/* ----------------------------------------------------------------------------
 *         Microchip Microprocessor (MPU) Software Team
 * ----------------------------------------------------------------------------
 * Copyright (C) 2019 Microchip Technology Inc. and its subsidiaries
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Microchip's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY MICROCHIP "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL MICROCHIP BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#ifndef __AT91_PMC_V1_H__
#define __AT91_PMC_V1_H__

#define PMC_UCKR		0x1C		/* UTMI Clock Configuration Register */
#define AT91C_CKGR_UOSCEN			(0x1UL <<  0)
#define AT91C_CKGR_OSCBYPASS			(0x1UL <<  1)
#define AT91C_CKGR_OSCOUNT			(0xFFUL << 8)
#define AT91C_CKGR_UPLLEN			(0x1UL << 16)
#define		AT91C_CKGR_UPLLEN_DISABLED 	(0x0UL << 16)
#define		AT91C_CKGR_UPLLEN_ENABLED  	(0x1UL << 16)
#define AT91C_CKGR_UPLLCOUNT			(0xFUL << 20)
#define		AT91C_CKGR_UPLLCOUNT_DEFAULT	(0x1UL << 20)
#define AT91C_CKGR_BIASEN			(0x1UL << 24)
#define		AT91C_CKGR_BIASEN_DISABLED	(0x0UL << 24)
#define		AT91C_CKGR_BIASEN_ENABLED	(0x1UL << 24)
#define AT91C_CKGR_BIASCOUNT			(0xFUL << 28)

#define PMC_PLLAR		0x28		/* PLL A Register */
#define AT91C_CKGR_DIVA		(0xff << 0)
#define		AT91C_CKGR_DIVA_MSK		0xff
#define		AT91C_CKGR_DIVA_OFFSET		0
#define		AT91C_CKGR_DIVA_0		0x0
#define		AT91C_CKGR_DIVA_BYPASS		0x1
#define AT91C_CKGR_PLLACOUNT	(0x3FUL <<  8)
#define AT91C_CKGR_OUTA		(0x3UL << 14)
#define 	AT91C_CKGR_OUTA_0		(0x0UL << 14)
#define 	AT91C_CKGR_OUTA_1		(0x1UL << 14)
#define 	AT91C_CKGR_OUTA_2		(0x2UL << 14)
#define 	AT91C_CKGR_OUTA_3		(0x3UL << 14)
#define AT91C_CKGR_MULA		(0xff << 16)
#define	AT91C_CKGR_ALT_MULA	(0x7f << 18)
#define		AT91C_CKGR_ALT_MULA_MSK		0x7f
#define		AT91C_CKGR_ALT_MULA_OFFSET	18
#define AT91C_CKGR_SRCA		(0x1UL << 29)

#define PMC_MCKR		0x30		/* Master Clock Register */
#define AT91C_PMC_CSS		(0x7UL <<  0)
#define 	AT91C_PMC_CSS_SLOW_CLK		(0x0UL)
#define 	AT91C_PMC_CSS_MAIN_CLK		(0x1UL)
#define 	AT91C_PMC_CSS_PLLA_CLK		(0x2UL)
#define 	AT91C_PMC_CSS_UPLL_CLK		(0x3UL)
#define 	AT91C_PMC_CSS_SYS_CLK		(0x4UL)
#define	AT91C_PMC_PRES		(0xFUL << 2)
#define		AT91C_PMC_PRES_CLK		(0x0UL << 2)
#define		AT91C_PMC_PRES_CLK_2		(0x1UL << 2)
#define		AT91C_PMC_PRES_CLK_4		(0x2UL << 2)
#define		AT91C_PMC_PRES_CLK_8		(0x3UL << 2)
#define		AT91C_PMC_PRES_CLK_16		(0x4UL << 2)
#define		AT91C_PMC_PRES_CLK_32		(0x5UL << 2)
#define		AT91C_PMC_PRES_CLK_64		(0x6UL << 2)
#define		AT91C_PMC_PRES_CLK_3		(0x7UL << 2)
#define		AT91C_PMC_PRES_CLK_1_5		(0x8UL << 2)
/* Master/Processor Clock Prescaler [alternate location ] */
#define	AT91C_PMC_ALT_PRES	(0xFUL << 4)
#define		AT91C_PMC_PRES_ALT_CLK		(0x0UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_2	(0x1UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_4	(0x2UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_8	(0x3UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_16	(0x4UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_32	(0x5UL << 4)
#define		AT91C_PMC_PRES_ALT_CLK_64	(0x6UL << 4)
#define AT91C_PMC_MDIV		(0x3UL <<  8)
#define 	AT91C_PMC_MDIV_1		(0x0UL <<  8)
#define 	AT91C_PMC_MDIV_2		(0x1UL <<  8)
#define 	AT91C_PMC_MDIV_4		(0x2UL <<  8)
#define 	AT91C_PMC_MDIV_3		(0x3UL <<  8)
#define AT91C_PMC_PLLADIV2	(0x1UL << 12)
#define		AT91C_PMC_PLLADIV2_1		(0x0UL << 12)
#define		AT91C_PMC_PLLADIV2_2		(0x1UL << 12)
#define AT91C_PMC_H32MXDIV	(0x1UL << 24)
#define		AT91C_PMC_H32MXDIV_H32MXDIV1	(0x0UL << 24)
#define		AT91C_PMC_H32MXDIV_H32MXDIV2	(0x1UL << 24)

#define PMC_PLLICPR		0x80		/* PLL Charge Pump Current Register */
#define AT91C_PMC_ICPPLLA	(0xFUL <<  0)
#define 	AT91C_PMC_ICPPLLA_0		(0x0UL << 0)
#define 	AT91C_PMC_ICPPLLA_1		(0x1UL << 0)
#define AT91C_PMC_REALLOCK	(0x1UL <<  7)
#define AT91C_PMC_IPLLA		(0xFUL <<  8)
#define 	AT91C_PMC_IPLLA_0		(0x0UL <<  8)
#define 	AT91C_PMC_IPLLA_1		(0x1UL <<  8)
#define 	AT91C_PMC_IPLLA_2		(0x2UL <<  8)
#define 	AT91C_PMC_IPLLA_3		(0x3UL <<  8)

#define PMC_PCR			0x10C		/* Peripheral Control Register */
#define	AT91C_PMC_PID		(0x3F << 0)
#define	AT91C_PMC_GCKCSS	(0x7 << 8)
#define		AT91C_PMC_GCKCSS_SLOW_CLK	(0x0 << 8)
#define		AT91C_PMC_GCKCSS_MAIN_CLK	(0x1 << 8)
#define		AT91C_PMC_GCKCSS_PLLA_CLK	(0x2 << 8)
#define		AT91C_PMC_GCKCSS_UPLL_CLK	(0x3 << 8)
#define		AT91C_PMC_GCKCSS_MCK_CLK	(0x4 << 8)
#define		AT91C_PMC_GCKCSS_AUDIO_CLK	(0x5 << 8)
#define	AT91C_PMC_CMD		(1 << 12)
#define AT91C_PMC_DIV		(0x3UL << 16)
#define		AT91C_PMC_DIV_(x)		((x) << 16)
#define	AT91C_PMC_GCKDIV	(0xff << 20)
#define		AT91C_PMC_GCKDIV_MSK		0xff
#define		AT91C_PMC_GCKDIV_OFFSET		20
#define		AT91C_PMC_GCKDIV_(x)		(((x) & AT91C_PMC_GCKDIV_MSK) << AT91C_PMC_GCKDIV_OFFSET)
#define	AT91C_PMC_EN		(0x1 << 28)
#define	AT91C_PMC_GCKEN		(0x1 << 29)

#endif
