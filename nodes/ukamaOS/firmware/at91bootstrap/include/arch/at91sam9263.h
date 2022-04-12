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
#ifndef __AT91SAM9263__
#define __AT91SAM9263__

/*
 * Peripheral identifiers/interrupts.
 */
#define AT91C_ID_FIQ		0	/* Advanced Interrupt Controller (FIQ) */
#define AT91C_ID_SYS		1	/* System Controller */
#define AT91C_ID_PIOA		2	/* Parallel IO Controller A */
#define AT91C_ID_PIOB		3	/* Parallel IO Controller B */
#define AT91C_ID_PIOCDE		4	/* Parallel IO Controller C, D, and E */
/* 5 reserved */
/* 6 reserved */
#define AT91C_ID_US0		7	/* USART 0 */
#define AT91C_ID_US1		8	/* USART 1 */
#define AT91C_ID_US2		9	/* USART 2 */
#define AT91C_ID_MCI0		10	/* Multimedia Card Interface 0 */
#define AT91C_ID_MCI1		11	/* Multimedia Card Interface 1 */
#define AT91C_ID_CAN		12	/* CAN Controller */
#define AT91C_ID_TWI		13	/* Two-Wire Interface */
#define AT91C_ID_SPI0		14	/* Serial Peripheral Interface 0 */
#define AT91C_ID_SPI1		15	/* Serial Peripheral Interface 1 */
#define AT91C_ID_SSC0		16	/* Serial Synchronous Controller 0 */
#define AT91C_ID_SSC1		17	/* Serial Synchronous Controller 1 */
#define AT91C_ID_AC97C		18	/* AC97 Controller */
#define AT91C_ID_TC012		19	/* Timer/Counter 0, 1, and 2 */
#define AT91C_ID_PWMC		20	/* Pulse Width Modulation Controller */
#define AT91C_ID_EMAC		21	/* Ethernet Mac */
/* 22 reserved */
#define AT91C_ID_2DGE		23	/* 2D Graphic Engine */
#define AT91C_ID_UDP		24	/* USB Device Port */
#define AT91C_ID_ISI		25	/* Image Sensor Interface */
#define AT91C_ID_LCDC		26	/* LCD Controller */
#define AT91C_ID_DMA		27	/* DMA Controller */
/* 28 reserved */
#define AT91C_ID_UHP		29	/* USB Host Port */
#define AT91C_ID_IRQ0		30	/* Advanced Interrupt Controller (IRQ0) */
#define AT91C_ID_IRQ1		31	/* Advanced Interrupt Controller (IRQ1) */

/*
 * User Peripherals physical base addresses.
 */
#define AT91C_BASE_UDP		0xfff78000
#define AT91C_BASE_TC0		0xfff7c000
#define AT91C_BASE_TC1		0xfff7c040
#define AT91C_BASE_TC2		0xfff7c080
#define AT91C_BASE_MCI0		0xfff80000
#define AT91C_BASE_MCI1		0xfff84000
#define AT91C_BASE_TWI		0xfff88000
#define AT91C_BASE_USART0	0xfff8c000
#define AT91C_BASE_USART1	0xfff90000
#define AT91C_BASE_USART2	0xfff94000
#define AT91C_BASE_SSC0		0xfff98000
#define AT91C_BASE_SSC1		0xfff9c000
#define AT91C_BASE_AC97C	0xfffa0000
#define AT91C_BASE_SPI0		0xfffa4000
#define AT91C_BASE_SPI1		0xfffa8000
#define AT91C_BASE_CAN0		0xfffac000
#define AT91C_BASE_PWMC		0xfffb8000
#define AT91C_BASE_EMAC		0xfffbc000
#define AT91C_BASE_ISI		0xfffc4000
#define AT91C_BASE_2DGE		0xfffc8000

/*
 * System Peripherals physical base addresses.
 */
#define AT91C_BASE_ECC0		0xffffe200
#define AT91C_BASE_SDRAMC0	0xffffe200
#define AT91C_BASE_SMC0		0xffffe400
#define AT91C_BASE_ECC1		0xffffe600
#define AT91C_BASE_SDRAMC1	0xffffe800
#define AT91C_BASE_SMC1		0xffffea00
#define AT91C_BASE_MATRIX	0xffffec00
#define AT91C_BASE_CCFG		0xffffed10
#define AT91C_BASE_DBGU		0xffffee00
#define AT91C_BASE_AIC		0xfffff000
#define AT91C_BASE_PIOA		0xfffff200
#define AT91C_BASE_PIOB		0xfffff400
#define AT91C_BASE_PIOC		0xfffff600
#define AT91C_BASE_PIOD		0xfffff800
#define AT91C_BASE_PIOE		0xfffffa00
#define AT91C_BASE_PMC		0xfffffc00
#define AT91C_BASE_RSTC		0xfffffd00
#define AT91C_BASE_SHDWC	0xfffffd10
#define AT91C_BASE_RTT0		0xfffffd20
#define AT91C_BASE_PITC		0xfffffd30
#define AT91C_BASE_WDT		0xfffffd40
#define AT91C_BASE_RTT1		0xfffffd50
#define AT91C_BASE_GPBR		0xfffffd60

/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_BOOT		0x00000000	/* Boot Memory */
#define AT91C_BASE_ITCM		0x00100000
#define AT91C_BASE_DTCM		0x00200000
#define AT91C_BASE_SRAM		0x00300000
#define AT91C_BASE_ROM		0x00400000
#define AT91C_BASE_SRAM0	0x00500000

/*
 * External memory
 */
#define AT91C_BASE_EBI0_CS0	0x10000000
#define AT91C_BASE_EBI0_CS1	0x20000000	/* EBI0 SDRAMC */
#define AT91C_BASE_EBI0_CS2	0x30000000
#define AT91C_BASE_EBI0_CS3	0x40000000	/* EBI0 NANDFlash */
#define AT91C_BASE_EBI0_CS4	0x50000000	/* Compact Flash Slot 0 */
#define AT91C_BASE_EBI0_CS5	0x60000000	/* Compact Flash Slot 1 */
#define AT91C_BASE_EBI1_CS0	0x70000000
#define AT91C_BASE_EBI1_CS1	0x80000000	/* EBI1 SDRAMC */
#define AT91C_BASE_EBI1_CS2	0x90000000	/* EBI1 NANDFlash */

#define AT91C_NUM_PIO		5

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME		"AT91SAM9263"

#endif /* #ifndef __AT91SAM9263__ */
