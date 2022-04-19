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
#ifndef __AT91SAM9RL_H__
#define __AT91SAM9RL_H__

/*
 * Peripheral ID Definitions
 */
#define AT91C_ID_FIQ		0	/* Advanced Interrupt Controller (FIQ) */
#define AT91C_ID_SYS		1	/* System Controller */
#define AT91C_ID_PIOA		2	/* Parallel I/O Controller A */
#define AT91C_ID_PIOB		3	/* Parallel I/O Controller B */
#define AT91C_ID_PIOC		4	/* Parallel I/O Controller C */
#define AT91C_ID_PIOD		5	/* Parallel I/O Controller D */
#define AT91C_ID_US0		6	/* USART 0 */
#define AT91C_ID_US1		7	/* USART 1 */
#define AT91C_ID_US2		8	/* USART 2 */
#define AT91C_ID_US3		9	/* USART 3 */
#define AT91C_ID_MCI		10	/* Multimedia Card Interface */
#define AT91C_ID_TWI0		11	/* Two-Wire Interface 0 */
#define AT91C_ID_TWI1		12	/* Two-Wire Interfaca 1 */
#define AT91C_ID_SPI		13	/* Serial Peripheral Interface */
#define AT91C_ID_SSC0		14	/* Serial Synchronous Controller 0 */
#define AT91C_ID_SSC1		15	/* Serial Synchronous Controller 1 */
#define AT91C_ID_TC0		16	/* Timer Counter 0 */
#define AT91C_ID_TC1		17	/* Timer Counter 1 */
#define AT91C_ID_TC2		18	/* Timer Counter 2 */
#define AT91C_ID_PWMC		19	/* Pulse Width Modulation Controller */
#define AT91C_ID_TSADCC		20	/* Touch Screen ADC Controller */
#define AT91C_ID_DMAC		21	/* DMA Controller */
#define AT91C_ID_UDPHS		22	/* USB Device High Speed */
#define AT91C_ID_LCDC		23	/* LCD Controller(AT91SAM9RL64 only) */
#define AT91C_ID_AC97		24	/* AC97 Controller */
#define AT91C_ID_IRQ		31	/* Advanced Interrupt Controller (IRQ0) */

/*
 * User Peripherals physical base addresses.
 */
#define AT91C_BASE_TC012	0xfffa0000
#define AT91C_BASE_MCI		0xfffa4000
#define AT91C_BASE_TWI0		0xfffa8000
#define AT91C_BASE_TWI1		0xfffac000
#define AT91C_BASE_USART0	0xfffb0000
#define AT91C_BASE_USART1	0xfffb4000
#define AT91C_BASE_USART2	0xfffb8000
#define AT91C_BASE_USART3	0xfffbc000
#define AT91C_BASE_SSC0		0xfffc0000
#define AT91C_BASE_SSC1		0xfffc4000
#define AT91C_BASE_PWMC		0xfffc8000
#define AT91C_BASE_SPI		0xfffcc000
#define AT91C_BASE_TSADC	0xfffd0000
#define AT91C_BASE_UDPHS	0xfffd4000
#define AT91C_BASE_AC97C	0xfffd8000

/*
 * Base Address Definitions
 */
#define AT91C_BASE_DMAC		0xffffe600
#define AT91C_BASE_ECC		0xffffe800
#define AT91C_BASE_SDRAMC	0xffffea00
#define AT91C_BASE_SMC		0xffffec00
#define AT91C_BASE_MATRIX	0xffffee00
#define AT91C_BASE_CCFG		0xffffef10
#define AT91C_BASE_AIC		0xfffff000
#define AT91C_BASE_DBGU		0xfffff200
#define AT91C_BASE_PIOA		0xfffff400
#define AT91C_BASE_PIOB		0xfffff600
#define AT91C_BASE_PIOC		0xfffff800
#define AT91C_BASE_PIOD		0xfffffa00
#define AT91C_BASE_PMC		0xfffffc00
#define AT91C_BASE_RSTC		0xfffffd00
#define AT91C_BASE_SHDWC	0xfffffd10
#define AT91C_BASE_RTTC		0xfffffd20
#define AT91C_BASE_PITC		0xfffffd30
#define AT91C_BASE_WDT		0xfffffd40
#define AT91C_BASE_SCKCR	0xfffffd50
#define AT91C_BASE_GPBR		0xfffffd60
#define AT91C_BASE_RTCC		0xfffffe00

/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_BOOT		0x00000000
#define AT91C_BASE_ITCM		0x00100000
#define AT91C_BASE_DTCM		0x00200000
#define AT91C_BASE_SRAM		0x00300000
#define AT91C_BASE_ROM		0x00400000
#define AT91C_BASE_LCDC		0x00500000
#define AT91C_BASE_UDPHS_RAM	0x00600000

/*
 * External memory
 */
#define AT91C_BASE_CS0		0x10000000
#define AT91C_BASE_CS1		0x20000000	/* SDRAMC */
#define AT91C_BASE_CS2		0x30000000
#define AT91C_BASE_CS3		0x40000000	/* NANDFlash */
#define AT91C_BASE_CS4		0x50000000	/* Compact Flash Slot 0 */
#define AT91C_BASE_CS5		0x60000000	/* Compact Flash Slot 1 */

#define AT91C_NUM_PIO		4

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME		"AT91SAM9RL"

#endif /* #ifndef __AT91SAM9RL_H__ */
