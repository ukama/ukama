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
#ifndef __AT91_PMC_V3_H__
#define __AT91_PMC_V3_H__

#define PMC_PLL_CTRL0		0x0C		/* PLL Control Register 0 */
#define AT91C_PLL_CTRL0_DIVPMC	(0xFFUL << 0)
#define		AT91C_PLL_CTRL0_DIVPMC_(x)	((x) << 0)
#define AT91C_PLL_CTRL0_DIVIO	(0xFFUL << 12)
#define		AT91C_PLL_CTRL0_DIVIO_(x)	((x) << 12)
#define AT91C_PLL_CTRL0_ENPLL	(0x1UL << 28)
#define AT91C_PLL_CTRL0_ENPLLCK	(0x1UL << 29)
#define AT91C_PLL_CTRL0_ENIOPLLCK (0x1UL << 30)
#define AT91C_PLL_CTRL0_ENLOCK	(0x1UL << 31)

#define PMC_PLL_CTRL1		0x10		/* PLL Control Register 1 */
#define AT91C_PLL_CTRL1_FRACR	(0x3FFFFFUL << 0)
#define		AT91C_PLL_CTRL1_FRACR_(x)	((x) << 0)
#define AT91C_PLL_CTRL1_MUL	(0xFFUL << 24)
#define		AT91C_PLL_CTRL1_MUL_(x)		((x) << 24)

#define PMC_PLL_SSR		0x14		/* PLL Spread Spectrum Register */
#define AT91C_PLL_SSR_STEP	(0xFFFFUL << 0)
#define		AT91C_PLL_SSR_STEP_(x)		((x) << 0)
#define AT91C_PLL_SSR_NSTEP	(0xFFUL << 16)
#define		AT91C_PLL_SSR_NSTEP_(x)		((x) << 16)
#define AT91C_PLL_SSR_ENSPREAD	(0x1UL << 28)

#define PMC_PLL_ACR		0x18		/* PLL Analog Control Register */
#define AT91C_PLL_ACR_CONTROL	(0xFFFUL << 0)
#define		AT91C_PLL_ACR_CONTROL_(x)	((x) << 0)
#define AT91C_PLL_ACR_UTMIVR	(0x1UL << 12)
#define AT91C_PLL_ACR_UTMIBG	(0x1UL << 13)
#define AT91C_PLL_ACR_LOCK_THR	(0x7UL << 16)
#define		AT91C_PLL_ACR_LOCK_THR_(x)	((x) << 16)
#define AT91C_PLL_ACR_LOOP_FILTER (0x3FUL << 24)
#define		AT91C_PLL_ACR_LOOP_FILTER_(x)	((x) << 24)
#define AT91C_PLL_ACR_DEFAULT_UTMI	0x12020010UL
#define AT91C_PLL_ACR_DEFAULT_PLLA	0x20010UL

#define PMC_PLL_UPDT		0x1c		/* PLL Update Register */
#define AT91C_PLL_UPDT_ID	(0xFUL << 0)
#define		AT91C_PLL_UPDT_ID_(x)	((x) << 0)
#define AT91C_PLL_UPDT_UPDATE	(0x1UL << 8)
#define AT91C_PLL_UPDT_STUPTIM	(0x3FUL << 16)
#define AT91C_PLL_UPDT_STUPTIM_(x) ((x) << 16)

#define PMC_MCKR		0x28		/* Master Clock Register */
#define AT91C_PMC_CSS		(0x3UL <<  0)
#define 	AT91C_PMC_CSS_SLOW_CLK		(0x0UL)
#define 	AT91C_PMC_CSS_MAIN_CLK		(0x1UL)
#define 	AT91C_PMC_CSS_CPUPLL_CLK	(0x2UL)
#define 	AT91C_PMC_CSS_SYSPLL_CLK	(0x3UL)
#define	AT91C_PMC_PRES		(0x7UL << 4)
#define		AT91C_PMC_PRES_CLK		(0x0UL << 4)
#define		AT91C_PMC_PRES_CLK_2		(0x1UL << 4)
#define		AT91C_PMC_PRES_CLK_4		(0x2UL << 4)
#define		AT91C_PMC_PRES_CLK_8		(0x3UL << 4)
#define		AT91C_PMC_PRES_CLK_16		(0x4UL << 4)
#define		AT91C_PMC_PRES_CLK_32		(0x5UL << 4)
#define		AT91C_PMC_PRES_CLK_64		(0x6UL << 4)
#define		AT91C_PMC_PRES_CLK_3		(0x7UL << 4)
#define AT91C_PMC_MDIV		(0x7UL <<  8)
#define 	AT91C_PMC_MDIV_1		(0x0UL <<  8)
#define 	AT91C_PMC_MDIV_2		(0x1UL <<  8)
#define 	AT91C_PMC_MDIV_4		(0x2UL <<  8)
#define 	AT91C_PMC_MDIV_3		(0x3UL <<  8)
#define 	AT91C_PMC_MDIV_5		(0x4UL <<  8)
#define AT91C_PMC_PLLADIV2	(0x0)
#define AT91C_PMC_H32MXDIV	(0x0)

#define PMC_MCR		0x30		/* Master Clock Register */
#define AT91C_MCR_ID		(0xFUL << 0)
#define		AT91C_MCR_ID_(x) 		((x) << 0)
#define AT91C_MCR_CMD		(0x1UL << 7)
#define AT91C_MCR_DIV		(0x7UL << 8)
#define		AT91C_MCR_MASTER_DIV1		(0x0UL << 8)
#define		AT91C_MCR_MASTER_DIV2		(0x1UL << 8)
#define		AT91C_MCR_MASTER_DIV4		(0x2UL << 8)
#define		AT91C_MCR_MASTER_DIV8		(0x3UL << 8)
#define		AT91C_MCR_MASTER_DIV16		(0x4UL << 8)
#define		AT91C_MCR_MASTER_DIV32		(0x5UL << 8)
#define		AT91C_MCR_MASTER_DIV64		(0x6UL << 8)
#define		AT91C_MCR_MASTER_DIV3		(0x7UL << 8)
#define AT91C_MCR_CSS		(0x1FUL << 16)
#define		AT91C_MCR_CSS_MD_SLOW_CLK	(0UL << 16)
#define		AT91C_MCR_CSS_TD_SLOW_CLK	(1UL << 16)
#define		AT91C_MCR_CSS_MAIN_CLK		(2UL << 16)
#define		AT91C_MCR_CSS_MCK0_CLK		(3UL << 16)
#define		AT91C_MCR_CSS_SYSPLL_CLK	(5UL << 16)
#define		AT91C_MCR_CSS_DDRPLL_CLK	(6UL << 16)
#define		AT91C_MCR_CSS_IMGPLL_CLK	(7UL << 16)
#define		AT91C_MCR_CSS_BAUDPLL_CLK	(8UL << 16)
#define		AT91C_MCR_CSS_AUDIOPLL_CLK	(9UL << 16)
#define		AT91C_MCR_CSS_ETHPLL_CLK	(10UL << 16)
#define AT91C_MCR_EN		(0x1UL << 28)

#define PMC_PCR			0x88		/* Peripheral Control Register */
#define	AT91C_PMC_PID		(0x7FUL << 0)
#define	AT91C_PMC_GCKCSS	(0x1FUL << 8)
#define		AT91C_PMC_GCKCSS_MD_SLOW_CLK	(0UL << 8)
#define		AT91C_PMC_GCKCSS_SLOW_CLK	(1UL << 8)
#define		AT91C_PMC_GCKCSS_MAIN_CLK	(2UL << 8)
#define		AT91C_PMC_GCKCSS_MCK_CLK	(3UL << 8)
#define		AT91C_PMC_GCKCSS_SYSPLL_CLK	(5UL << 8)
#define		AT91C_PMC_GCKCSS_DDRPLL_CLK	(6UL << 8)
#define		AT91C_PMC_GCKCSS_IMGPLL_CLK	(7UL << 8)
#define		AT91C_PMC_GCKCSS_BAUDPLL_CLK	(8UL << 8)
#define		AT91C_PMC_GCKCSS_AUDIOPLL_CLK	(9UL << 8)
#define		AT91C_PMC_GCKCSS_ETHPLL_CLK	(10UL << 8)
#define AT91C_PMC_DIV		(0x3UL << 14)
#define		AT91C_PMC_DIV_(x)		((x) << 14)
#define	AT91C_PMC_MCKID		(0xFUL << 16)
#define	AT91C_PMC_GCKDIV	(0xFFUL << 20)
#define		AT91C_PMC_GCKDIV_MSK		0xFF
#define		AT91C_PMC_GCKDIV_OFFSET		20
#define		AT91C_PMC_GCKDIV_(x)		(((x) & AT91C_PMC_GCKDIV_MSK) << AT91C_PMC_GCKDIV_OFFSET)
#define	AT91C_PMC_EN		(0x1UL << 28)
#define	AT91C_PMC_GCKEN		(0x1UL << 29)
#define	AT91C_PMC_CMD		(0x1UL << 31)

#define PMC_GCSR0		0xC0		/* Generic Clock Status Register 0 */
#define PMC_GCSR1		0xC4		/* Generic Clock Status Register 1 */
#define PMC_GCSR2		0xC8		/* Generic Clock Status Register 2 */
#define PMC_GCSR3		0xCC		/* Generic Clock Status Register 3 */

#define PMC_PLL_ISR0		0xEC		/* PLL Interrupt Status Register 0 */
#define AT91C_PLL_ISR0_LOCKA	(0x1UL << 0)

#define PMC_PLL_ISR1		0xF0		/* PLL Interrupt Status Register 1 */

#endif

