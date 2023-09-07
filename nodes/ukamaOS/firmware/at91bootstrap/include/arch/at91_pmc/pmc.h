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
#ifndef __AT91_PMC_H__
#define __AT91_PMC_H__

#if defined(CONFIG_PMC_V1)
#include "pmc-v1.h"
#elif defined(CONFIG_PMC_V2)
#include "pmc-v2.h"
#elif defined(CONFIG_PMC_V3)
#include "pmc-v3.h"
#else
#error "Invalid PMC version"
#endif

/* Register offset in AT91C_PMC structure */
#define PMC_SCER	0x00	/* System Clock Enable Register */
#define PMC_SCDR	0x04	/* System Clock Disable Register */
#define PMC_SCSR	0x08	/* System Clock Status Register */
/* 0x0C Reversed */
#define PMC_MOR		0x20	/* Main Oscillator Register */
#define PMC_MCFR	0x24	/* Main Clock  Frequency Register */
/* 0x34 Reserved */
#define PMC_USB		0x38	/* USB clock register */
#define PMC_SMD		0x3C	/* SMD clock register */
#define PMC_PCKR	0x40	/* Programmable Clock 0 Register */
#define PMC_PCKR1	0x44	/* Programmable Clock 1 Register */
#define PMC_PCKR2	0x48	/* Programmable Clock 2 Register */
/* 0x48 ~ 0x5C */
#define PMC_IER		0x60	/* Interrupt Enable Register */
#define PMC_IDR		0x64	/* Interrupt Disable Register */
#define PMC_SR		0x68	/* Status Register */
#define PMC_IMR		0x6C	/* Interrupt Mask Register */
#define PMC_FSMR	0x70	/* Fast Startup Mode Register */
#define PMC_FSPR	0x74	/* Fast Startup Polarity Register */
#define PMC_FOCR	0x78	/* Fault Output Clear Register */
/* 0x7C Reserved */
#define PMC_WPMR	0xE4	/* Write Protect Mode Register */
#define PMC_WPSR	0xE8	/* Write Protect Status Register */
/* 0xEC ~ 0xF8 */
#define PMC_VERSION	0xFC	/* Version Register */
#define PMC_OCR		0x110	/* Oscillator Calibration Register */

/* -------- PMC_SCER : (PMC Offset: 0x0) System Clock Enable Register --------*/ 
#define AT91C_PMC_DDR		(0x1UL <<  2)
#define AT91C_PMC_LCDCK		(0x1UL <<  3)
#define AT91C_PMC_SMDCK		(0x1UL <<  4)
#define AT91C_PMC_UHP		(0x1UL <<  6)
#define AT91C_PMC_UDP		(0x1UL <<  7)
#define AT91C_PMC_PCK0		(0x1UL <<  8)
#define AT91C_PMC_PCK1		(0x1UL <<  9)
#define AT91C_PMC_PCK2		(0x1UL << 10)
#define AT91C_PMC_QSPICLK	(0x1UL << 19)

#define AVAILABLE_SYS_CLK (AT91C_PMC_DDR \
			 | AT91C_PMC_LCDCK \
			 | AT91C_PMC_SMDCK \
			 | AT91C_PMC_UHP \
			 | AT91C_PMC_UDP \
			 | AT91C_PMC_PCK0 \
			 | AT91C_PMC_PCK1 \
			 | AT91C_PMC_PCK2 \
			 | AT91C_PMC_QSPICLK)

/* -------- PMC_SCDR : (PMC Offset: 0x4) System Clock Disable Register --------*/ 
/* -------- PMC_SCSR : (PMC Offset: 0x8) System Clock Status Register ---------*/
#define AT91C_PMC_PCK		(0x1UL <<  0)

/* -------- CKGR_MOR : (PMC Offset: 0x20) Main Oscillator Register --------*/ 
#define AT91C_CKGR_MOSCXTEN	(0x1UL <<  0)
#define AT91C_CKGR_MOSCXTBY	(0x1UL <<  1)
#define AT91C_CKGR_WAITMODE	(0x1UL <<  2)
#define AT91C_CKGR_MOSCRCEN	(0x1UL <<  3)
#define AT91C_CKGR_MOSCXTST	(0xFFUL <<  8)
#define		AT91_CKGR_MOSCXTST_SET(value)	(value << 8)
#define AT91C_CKGR_KEY		(0xFFUL << 16)
#define		AT91C_CKGR_PASSWD		(0x37 << 16)
#define AT91C_CKGR_MOSCSEL	(0x1UL << 24)
#define AT91C_CKGR_CFDEN	(0x1UL << 25)

/* -------- CKGR_MCFR : (PMC Offset: 0x24) Main Clock Frequency Register --------*/ 
#define AT91C_CKGR_MAINF	(0xFFFFUL <<  0)
#define AT91C_CKGR_MAINRDY	(0x1UL << 16)
#define AT91C_CKGR_RCMEAS	(0x1UL << 20)
#define AT91C_CKGR_CCSS         (0x1UL << 24)
#define		AT91C_CKGR_CCSS_RC_OSC		(0x0UL << 24)
#define		AT91C_CKGR_CCSS_XTAL_OSC	(0x1UL << 24)

/* -------- CKGR_PLLBR : (PMC Offset: 0x2c) PLL B Register --------*/
#define AT91C_CKGR_DIVB		(0xFFUL <<  0)
#define 	AT91C_CKGR_DIVB_0		(0x0UL)
#define 	AT91C_CKGR_DIVB_BYPASS		(0x1UL)
#define AT91C_CKGR_PLLBCOUNT	(0x3FUL <<  8)
#define AT91C_CKGR_OUTB		(0x3UL << 14)
#define 	AT91C_CKGR_OUTB_0		(0x0UL << 14)
#define 	AT91C_CKGR_OUTB_1		(0x1UL << 14)
#define 	AT91C_CKGR_OUTB_2		(0x2UL << 14)
#define 	AT91C_CKGR_OUTB_3		(0x3UL << 14)
#define AT91C_CKGR_MULB		(0x7FFUL << 16)

/* -------- PMC_USB : (PMC Offset: 0x38) USB Clock Register --------*/ 
#define AT91C_PMC_USBS		(0x1UL <<  0)
#define 	AT91C_PMC_USBS_USB_PLLA		(0x0UL)
#define 	AT91C_PMC_USBS_USB_UPLL		(0x1UL)
#define AT91C_PMC_USBDIV	(0xFUL <<  8)

/* -------- PMC_SMD : (PMC Offset: 0x3c) SMD Clock Register --------*/ 
#define AT91C_PMC_SMDS		(0x1UL <<  0)
#define 	AT91C_PMC_SMDS_SMD_PLLA		(0x0UL)
#define 	AT91C_PMC_SMDS_SMD_PLLB		(0x1UL)
#define AT91C_PMC_SMDDIV	(0x1FUL <<  8)
#define AT91C_PMC_SMDDIV_(x)	((x & 0x1F) << 8)

/* -------- PMC_PCKR : (PMC Offset: 0x40) Programmable Clock Register --------*/ 
#define AT91C_PMC_SLCKMCK	(0x1UL <<  8)
#define 	AT91C_PMC_SLCKMCK_SLCK                 (0x0UL <<  8)
#define 	AT91C_PMC_SLCKMCK_MCK                  (0x1UL <<  8)

/* -------- PMC_IER : (PMC Offset: 0x60) PMC Interrupt Enable Register --------*/ 
#define AT91C_PMC_MOSCXTS	(0x1UL << 0)
#define AT91C_PMC_LOCKA		(0x1UL << 1)
#define AT91C_PMC_LOCKB		(0x1UL << 2)
#define AT91C_PMC_MCKRDY	(0x1UL << 3)
#define AT91C_PMC_LOCKU		(0x1UL << 6)
#define AT91C_PMC_PCKRDY0	(0x1UL << 8)
#define AT91C_PMC_PCKRDY1	(0x1UL << 9)
#define AT91C_PMC_MOSCSELS	(0x1UL << 16)
#define AT91C_PMC_MOSCRCS	(0x1UL << 17)
#define AT91C_PMC_CFDEV		(0x1UL << 18)
#define AT91C_PMC_GCKRDY	(0x1UL << 24)
#define AT91C_PMC_MCKXRDY	(0x1UL << 26)
/* -------- PMC_IDR : (PMC Offset: 0x64) PMC Interrupt Disable Register --------*/ 
/* -------- PMC_SR : (PMC Offset: 0x68) PMC Status Register --------*/ 
#define AT91C_PMC_CFDS		(0x1UL << 19)
#define AT91C_PMC_FOS		(0x1UL << 20)
/* -------- PMC_IMR : (PMC Offset: 0x6c) PMC Interrupt Mask Register --------*/ 

/* -------- PMC_FSMR : (PMC Offset: 0x70) Fast Startup Mode Register --------*/ 
#define AT91C_PMC_FSTT		(0xFFFFUL << 0)
#define AT91C_PMC_RTTAL		(0x1UL << 16)
#define AT91C_PMC_RTCAL		(0x1UL << 17)
#define AT91C_PMC_USBAL		(0x1UL << 18)
#define AT91C_PMC_LPM		(0x1UL << 20)

/* -------- PMC_FSPR : (PMC Offset: 0x74) Fast Startup Polarity Register --------*/ 
#define AT91C_PMC_FSTP		(0xFFFFUL <<  0)

/* -------- PMC_FOCR : (PMC Offset: 0x78) Fault Output Clear Register --------*/ 
#define AT91C_PMC_FOCLR		(0x1UL << 0)

/* -------- PMC_FEATURES : (PMC Offset: 0xf8)   --------*/ 
#define AT91C_PMC_CFGAHBCLK	(0x1UL <<  0)
#define 	AT91C_PMC_CFGAHBCLK_0		(0x0UL << 0)
#define 	AT91C_PMC_CFGAHBCLK_1		(0x1UL << 0)
#define AT91C_PMC_HCLKEN	(0x1UL <<  1)
#define 	AT91C_PMC_HCLKEN_0		(0x0UL <<  1)
#define 	AT91C_PMC_HCLKEN_1		(0x1UL <<  1)
#define AT91C_PMC_PERMCLK	(0x1UL <<  2)
#define 	AT91C_PMC_PERMCLK_0		(0x0UL <<  2)
#define 	AT91C_PMC_PERMCLK_1		(0x1UL <<  2)
#define AT91C_PMC_CORE2		(0x1UL <<  3)
#define 	AT91C_PMC_CORE2_0		(0x0UL <<  3)
#define 	AT91C_PMC_CORE2_1		(0x1UL <<  3)
#define AT91C_PMC_USBDEVCK	(0x1UL <<  4)
#define 	AT91C_PMC_USBDEVCK_0		(0x0UL <<  4)
#define 	AT91C_PMC_USBDEVCK_1		(0x1UL <<  4)
#define AT91C_PMC_USBHOSTCK	(0x1UL <<  5)
#define 	AT91C_PMC_USBHOSTCK_0		(0x0UL <<  5)
#define 	AT91C_PMC_USBHOSTCK_1		(0x1UL <<  5)
#define AT91C_PMC_USBOTGCK	(0x1UL <<  6)
#define 	AT91C_PMC_USBOTGCK_0		(0x0UL <<  6)
#define 	AT91C_PMC_USBOTGCK_1		(0x1UL <<  6)
#define AT91C_PMC_UHSYNRST	(0x1UL <<  7)
#define 	AT91C_PMC_UHSYNRST_0		(0x0UL <<  7)
#define 	AT91C_PMC_UHSYNRST_1		(0x1UL <<  7)
#define AT91C_PMC_UOSYNRST	(0x1UL <<  8)
#define 	AT91C_PMC_UOSYNRST_0		(0x0UL <<  8)
#define 	AT91C_PMC_UOSYNRST_1		(0x1UL <<  8)
#define AT91C_PMC_PLLENPOL	(0x1UL <<  9)
#define 	AT91C_PMC_PLLENPOL_0		(0x0UL <<  9)
#define 	AT91C_PMC_PLLENPOL_1		(0x1UL <<  9)
#define AT91C_PMC_BIASREG	(0x1UL << 10)
#define 	AT91C_PMC_BIASREG_0		(0x0UL << 10)
#define 	AT91C_PMC_BIASREG_1		(0x1UL << 10)
#define AT91C_PMC_OUTPLL	(0x1UL << 11)
#define 	AT91C_PMC_OUTPLL_0		(0x0UL << 11)
#define 	AT91C_PMC_OUTPLL_1		(0x1UL << 11)
#define AT91C_PMC_OUTCURR	(0x1UL << 12)
#define 	AT91C_PMC_OUTCURR_0		(0x0UL << 12)
#define 	AT91C_PMC_OUTCURR_1		(0x1UL << 12)
#define AT91C_PMC_FWUP		(0x1UL << 13)
#define 	AT91C_PMC_FWUP_0		(0x0UL << 13)
#define 	AT91C_PMC_FWUP_1		(0x1UL << 13)
#define AT91C_PMC_SELMAINCLK	(0x1UL << 14)
#define 	AT91C_PMC_SELMAINCLK_0		(0x0UL << 14)
#define 	AT91C_PMC_SELMAINCLK_1		(0x1UL << 14)
#define AT91C_PMC_RSTCLKM	(0x1UL << 15)
#define 	AT91C_PMC_RSTCLKM_0		(0x0UL << 15)
#define 	AT91C_PMC_RSTCLKM_1		(0x1UL << 15)
#define AT91C_PMC_NB_PERIPH_CLOCK	(0xFFUL << 16)

/* -------- PMC_VERSION : (PMC Offset: 0xfc)   --------*/ 
#define AT91C_PMC_Version	(0x0FFFUL << 0)
#define 	AT91C_PMC_Version_0		(0x0UL)
#define 	AT91C_PMC_Version_1		(0x1UL)
#define AT91C_PMC_MFN		(0x7UL << 16)
#define 	AT91C_PMC_MFN_0		(0x0UL << 16)
#define 	AT91C_PMC_MFN_1		(0x1UL << 16)

#endif /* #ifndef __AT91_PMC_H__ */
