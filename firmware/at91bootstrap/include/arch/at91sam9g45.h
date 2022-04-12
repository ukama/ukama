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
#ifndef __AT91SAM9G45_H__
#define __AT91SAM9G45_H__

/*
 * Peripheral identifiers/interrupts.
 */
#define AT91C_ID_FIQ		0	/* Advanced Interrupt Controller (FIQ) */
#define AT91C_ID_SYS		1	/* System Controller Interrupt */
#define AT91C_ID_PIOA		2	/* Parallel I/O Controller A */
#define AT91C_ID_PIOB		3	/* Parallel I/O Controller B */
#define AT91C_ID_PIOC		4	/* Parallel I/O Controller C */
#define AT91C_ID_PIOD_E		5	/* Parallel I/O Controller D/E */
#define AT91C_ID_TRNG		6	/* True Random Number Generator */
#define AT91C_ID_US0		7	/* USART 0 */
#define AT91C_ID_US1		8	/* USART 1 */
#define AT91C_ID_US2		9	/* USART 2 */
#define AT91C_ID_US3		10	/* USART 3 */
#define AT91C_ID_MCI0		11	/* High Speed Multimedia Card Interface 0 */
#define AT91C_ID_TWI0		12	/* Two-Wire Interface 0 */
#define AT91C_ID_TWI1		13	/* Two-Wire Interface 1 */
#define AT91C_ID_SPI0		14	/* Serial Peripheral Interface 0 */
#define AT91C_ID_SPI1		15	/* Serial Peripheral Interface 1 */
#define AT91C_ID_SSC0		16	/* Serial Synchronous Controller 0 */
#define AT91C_ID_SSC1		17	/* Serial Synchronous Controller 1 */
#define AT91C_ID_TCB		18	/* Timer Counter */
#define AT91C_ID_PWM		19	/* Pulse Width Modulation Controller */
#define AT91C_ID_TSADC		20	/* Touch Screen ADC Controller */
#define AT91C_ID_DMAC		21	/* DMA Controller */
#define AT91C_ID_UHPHS		22	/* USB Host High Speed */
#define AT91C_ID_LCDC		23	/* LCD Controller */
#define AT91C_ID_AC97C		24	/* AC97 Controller */
#define AT91C_ID_EMAC		25	/* Ethernet MAC */
#define AT91C_ID_ISI		26	/* Image Sensor Interface */
#define AT91C_ID_UDPHS		27	/* USB Device High Speed */
/* Reserved		28 */
#define AT91C_ID_MCI1		29	/* High Speed Multimedia Card Interface 1 */
#define AT91C_ID_VDEC		30	/* Video Decoder */
#define AT91C_ID_IRQ		31	/* Advanced Interrupt Controller (IRQ) */

/*
 * User Peripherals physical base addresses.
 */
#define AT91C_BASE_UDPHS	0xfff78000
#define AT91C_BASE_TCB0		0xfff7c000
#define AT91C_BASE_MCI0		0xfff80000
#define AT91C_BASE_TWI0		0xfff84000
#define AT91C_BASE_TWI1		0xfff88000
#define AT91C_BASE_USART0	0xfff8c000
#define AT91C_BASE_USART1	0xfff90000
#define AT91C_BASE_USART2	0xfff94000
#define AT91C_BASE_USART3	0xfff98000
#define AT91C_BASE_SSC0		0xfff9c000
#define AT91C_BASE_SSC1		0xfffa0000
#define AT91C_BASE_SPI0		0xfffa4000
#define AT91C_BASE_SPI1		0xfffa8000
#define AT91C_BASE_AC97		0xfffac000
#define AT91C_BASE_TSADC	0xfffb0000
#define AT91C_BASE_ISI		0xfffb4000
#define AT91C_BASE_PWMC		0xfffb8000
#define AT91C_BASE_EMAC		0xfffbc000
#define AT91C_BASE_TRNG		0xfffcc000
#define AT91C_BASE_MCI1		0xfffd0000
#define AT91C_BASE_TCB2		0xfffd4000

/*
 * System Peripherals physical base addresses.
 */
#define AT91C_BASE_ECC		0xffffe200
#define AT91C_BASE_DDRSDRC1	0xffffe400
#define AT91C_BASE_DDRSDRC0	0xffffe600
#define AT91C_BASE_SMC		0xffffe800
#define AT91C_BASE_MATRIX	0xffffea00
#define AT91C_BASE_DMAC		0xffffec00
#define AT91C_BASE_DBGU		0xffffee00
#define AT91C_BASE_AIC		0xfffff000
#define AT91C_BASE_PIOA		0xfffff200
#define AT91C_BASE_PIOB		0xfffff400
#define AT91C_BASE_PIOC		0xfffff600
#define AT91C_BASE_PIOD		0xfffff800
#define AT91C_BASE_PIOE		0xfffffa00
#define AT91C_BASE_PMC		0xfffffc00
#define AT91C_BASE_RSTC		0xfffffd00
#define AT91C_BASE_SHDC		0xfffffd10
#define AT91C_BASE_RTTC		0xfffffd20
#define AT91C_BASE_PITC		0xfffffd30
#define AT91C_BASE_WDT		0xfffffd40
#define AT91C_BASE_SCKCR	0xfffffd50
#define AT91C_BASE_GPBR		0xfffffd60
#define AT91C_BASE_RTCC		0xfffffdb0

/* (CCFG) Base Address */
#define AT91C_BASE_CCFG		0xffffeb20

/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_BOOT		0x00000000	/* Boot mapped area */
#define AT91C_BASE_ITCM		0x00100000	/* Boot mapped area */
#define AT91C_BASE_DTCM		0x00200000	/* Boot mapped area */
#define AT91C_BASE_SRAM		0x00300000	/* Internal SRAM */
#define AT91C_BASE_ROM		0x00400000	/* Internal ROM base address */
#define AT91C_BASE_LCD		0x00500000	/* LCD User Interface */
#define AT91C_BASE_UDPHS_RAM	0x00600000	/* USB Device controller */
#define AT91C_BASE_UHP_OHCI	0x00700000      /* UHPHS_OHCI Base address */
#define AT91C_BASE_UHP_EHCI	0x00800000      /* UHPHS_EHCI Base address */
#define AT91C_BASE_VDEC		0x00900000      /* Video Decoder */

/*
 * External memory
 */
#define AT91C_BASE_CS0		0x10000000	/* typically NOR */
#define AT91C_BASE_CS1		0x20000000	/* DRAM1 */
#define AT91C_BASE_CS2		0x30000000
#define AT91C_BASE_CS3		0x40000000	/* typically NAND */
#define AT91C_BASE_CS4		0x50000000
#define AT91C_BASE_CS5		0x60000000
#define AT91C_BASE_CS6		0x70000000	/* DRAM0 */

#define AT91C_NUM_PIO		5

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME		"AT91SAM9G45"

#endif	/* #ifndef __AT91SAM9G45_H__ */
