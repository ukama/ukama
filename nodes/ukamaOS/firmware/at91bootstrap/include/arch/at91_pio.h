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
#ifndef __AT91_PIO_H__
#define __AT91_PIO_H__

/* Register offset in AT91S_PIO structure */
#ifdef CONFIG_CPU_HAS_PIO4
#define	PIO_MSKR	0x0000	/* PIO Mask Register */
#define	PIO_CFGR	0x0004	/* PIO Configuration Register */
#define	PIO_PDSR	0x0008	/* PIO Pin Data Status Register */
#define	PIO_LOCKSR	0x000c	/* PIO Lock Status Register */
#define	PIO_SODR	0x0010	/* PIO Set Output Data Register */
#define	PIO_CODR	0x0014	/* PIO Clear Output Data Register */
#define	PIO_ODSR	0x0018	/* PIO Output Data Status Register */
/* 0x001c Reserved */
#define	PIO_IER		0x0020	/* PIO Interrupt Enable Register */
#define	PIO_IDR		0x0024	/* PIO Interrupt Disable Register */
#define	PIO_IMR		0x0028	/* PIO Interrupt Mask Register */
#define	PIO_ISR		0x002c	/* PIO Interrupt Status Register */
#define	PIO_SIONR	0x0030	/* PIO Set I/O Non-Secure Register */
#define	PIO_SIOSR	0x0034	/* PIO Set I/O Secure Register */
#define	PIO_IOSSR	0x0038	/* PIO I/O Secure Status Register */
#define	PIO_IOFR	0x003c	/* PIO I/O Freeze Register */
#define	PIO_DELAYR	0x0040	/* PIO I/O Delay Register */

#define	S_PIO_SCDR	0x1500	/* Secure PIO Slow Clock Divider Debouncing Register */
#define	S_PIO_TPFR	0x15D0	/* Secure PIO Tamper Freeze Register */
/* 0x15D4 reserved */
#define	S_PIO_WPMR	0x15E0	/* Secure PIO Write Protection Mode Register */
#define	S_PIO_WPSR	0x15E4	/* Secure PIO Write Protection Status Register */

#else

#define PIO_PER		0x0000	/* PIO Enable Register */
#define PIO_PDR		0x0004	/* PIO Disable Register */
#define PIO_PSR		0x0008	/* PIO Status Register */
/* 0x000c */
#define PIO_OER		0x0010	/* Output Enable Register */
#define PIO_ODR		0x0014	/* Output Disable Registerr */
#define PIO_OSR		0x0018	/* Output Status Register */
/* 0x001c */
#define PIO_IFER	0x0020	/* Input Filter Enable Register */
#define PIO_IFDR	0x0024	/* Input Filter Disable Register */
#define PIO_IFSR	0x0028	/* Input Filter Status Register */
/* 0x002c */
#define PIO_SODR	0x0030	/* Set Output Data Register */
#define PIO_CODR	0x0034	/* Clear Output Data Register */
#define PIO_ODSR	0x0038	/* Output Data Status Register */
#define PIO_PDSR	0x003c	/* Pin Data Status Register */
#define PIO_IER		0x0040	/* Interrupt Enable Register */
#define PIO_IDR		0x0044	/* Interrupt Disable Register */
#define PIO_IMR		0x0048	/* Interrupt Mask Register */
#define PIO_ISR		0x004c	/* Interrupt Status Register */
#define PIO_MDER	0x0050	/* Multi-driver Enable Register */
#define PIO_MDDR	0x0054	/* Multi-driver Disable Register */
#define PIO_MDSR	0x0058	/* Multi-driver Status Register */
/* 0x005c */
#define PIO_PPUDR	0x0060	/* Pull-up Disable Register */
#define PIO_PPUER	0x0064	/* Pull-up Enable Register */
#define PIO_PPUSR	0x0068	/* Pull-up Status Register */
/* 0x006c */
#ifndef CONFIG_CPU_HAS_PIO3
#define PIO_ASR		0x0070	/* Peripheral Select Register 1 */
#define PIO_BSR		0x0074	/* Peripheral Select Register 2 */
#define PIO_ABSR	0x0078	/* AB Select Status Register */
/* 0x007c ~ 0x009c */
#else
#define PIO_SP1		0x0070	/* Select Peripheral 1 Register */
#define PIO_SP2		0x0074	/* Select Peripheral 2 Register */
#define PIO_ABSR	0x0078	/* AB Select Status Register */
/* 0x007c ~ 0x008c */
#define PIO_PPDDR	0x0090	/* Pull-down Disable Register */
#define PIO_PPDER	0x0094	/* Pull-down Enable Register */
#define PIO_PPDSR	0x0098	/* Pull-down Status Register */
/* 0x009c ~ 0x010c */
#define PIO_SLEWR	0x0110	/* I/O Slewrate Control Register */
/* 0x0114 */
#define PIO_DRIVER1	0x0118	/* I/O Drive Register 1 */
/* 0x011a */
#endif	/* !CONFIG_CPU_HAS_PIO3 */
#define PIO_OWER	0x00a0	/* Output Write Enable Register */
#define PIO_OWDR	0x00a4	/* Output Write Disable Register */
#define PIO_OWSR	0x00a8	/* Output Write Status Register */
/* 0x00ac */
#endif	/* CONFIG_CPU_HAS_PIO4 */

/*
 * Register Field definition
 */
/*-------- PIO_CFGR: (Offset: 0x4) -------*/
#define	AT91C_PIO_CFGR_FUNC	(0x07 << 0)	/* I/O Line Function */
#define		AT91C_PIO_CFGR_FUNC_GPIO	0x00
#define		AT91C_PIO_CFGR_FUNC_PERIPH_A	0x01
#define		AT91C_PIO_CFGR_FUNC_PERIPH_B	0x02
#define		AT91C_PIO_CFGR_FUNC_PERIPH_C	0x03
#define		AT91C_PIO_CFGR_FUNC_PERIPH_D	0x04
#define		AT91C_PIO_CFGR_FUNC_PERIPH_E	0x05
#define		AT91C_PIO_CFGR_FUNC_PERIPH_F	0x06
#define		AT91C_PIO_CFGR_FUNC_PERIPH_G	0x07
#define	AT91C_PIO_CFGR_DIR	(0x01 << 8)	/* Direction */
#define	AT91C_PIO_CFGR_PUEN	(0x01 << 9)	/* Pull-Up Enable */
#define	AT91C_PIO_CFGR_PDEN	(0x01 << 10)	/* Pull-Down Enable */
#define	AT91C_PIO_CFGR_IFEN	(0x01 << 12)	/* Input Filter Enable */
#define	AT91C_PIO_CFGR_IFSCEN	(0x01 << 13)	/* Input Filter Slow Clock Enable */
#define	AT91C_PIO_CFGR_OPD	(0x01 << 14)	/* Open-Drain */
#define	AT91C_PIO_CFGR_SCHMITT	(0x01 << 15)	/* Schmitt Trigger*/
#define	AT91C_PIO_CFGR_DRVSTR	(0x03 << 16)	/* Drive Strength*/
#define		AT91C_PIO_CFGR_DRVSTR_HIGH	(0x03 << 16)
#define		AT91C_PIO_CFGR_DRVSTR_MEDIUM	(0x02 << 16)
#define		AT91C_PIO_CFGR_DRVSTR_LOW	(0x01 << 16)
#define	AT91C_PIO_CFGR_EVTSEL	(0x07 << 24)	/* Event Selection */
#define		AT91C_PIO_CFGR_EVTSEL_FALLING	(0x00 << 24)
#define		AT91C_PIO_CFGR_EVTSEL_RISING	(0x01 << 24)
#define		AT91C_PIO_CFGR_EVTSEL_BOTH	(0x02 << 24)
#define		AT91C_PIO_CFGR_EVTSEL_LOW	(0x03 << 24)
#define		AT91C_PIO_CFGR_EVTSEL_HIGH	(0x04 << 24)
#define	AT91C_PIO_CFGR_PCFS	(0x01 << 29)	/* Physical Configuration Freeze Status */
#define	AT91C_PIO_CFGR_ICFS	(0x01 << 30)	/* Interrupt Configuration Freeze Status */
#define	AT91C_PIO_CFGR_TAMPEN	(0x01 << 31)	/* Tamper Enable */

/* Number of IO handled by one PIO controller */
#define PIO_NUM_IO		32

#define	AT91C_PIN_PA(io)	(0 * PIO_NUM_IO + io)
#define	AT91C_PIN_PB(io)	(1 * PIO_NUM_IO + io)
#define	AT91C_PIN_PC(io)	(2 * PIO_NUM_IO + io)
#define	AT91C_PIN_PD(io)	(3 * PIO_NUM_IO + io)
#define	AT91C_PIN_PE(io)	(4 * PIO_NUM_IO + io)

/* define the register */
#define PIO_REG_PER		0x00
#define PIO_REG_OER		0x10
#define PIO_REG_CODR		0x34
#define PIO_REG_PPUDR		0x60
#define	PIO_REG_PPDDR		0x90

#endif	/* !__AT91_PIO_H__ */
