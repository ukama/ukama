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
#ifndef __AT91SAM9X5_H__
#define __AT91SAM9X5_H__

/*
 * Peripheral identifiers/interrupts.
 */
#define AT91C_ID_FIQ		0	/* Advanced Interrupt Controller (FIQ) */
#define AT91C_ID_SYS		1	/* System Controller Interupt */
#define AT91C_ID_PIOA_B		2	/* Parallel I/O Controller A and B */
#define AT91C_ID_PIOC_D		3	/* Parallel I/O Controller C and D */
#define AT91C_ID_SMD		4	/* SMD Software Modem */
#define AT91C_ID_USART0		5	/* USART 0 */
#define AT91C_ID_USART1		6	/* USART 1 */
#define AT91C_ID_USART2		7	/* USART 2 */
#define AT91C_ID_USART3		8	/* USART 3 */
#define AT91C_ID_TWI0		9	/* Two-Wire Interface 0 */
#define AT91C_ID_TWI1		10	/* Two-Wire Interface 1 */
#define AT91C_ID_TWI2		11	/* Two-Wire Interface 2 */
#define AT91C_ID_HSMCI0		12	/* High Speed Multimedia Card Interface 0 */
#define AT91C_ID_SPI0		13	/* Serial Peripheral Interface 0 */
#define AT91C_ID_SPI1		14	/* Serial Peripheral Interface 1 */
#define AT91C_ID_UART0		15	/* UART 0 */
#define AT91C_ID_UART1		16	/* UART 1 */
#define AT91C_ID_TCB		17	/* Timer Counter 0, 1, 2, 3, 4, 5 */
#define AT91C_ID_PWM		18	/* Pulse Width Modulation Controller */
#define AT91C_ID_ADC		19	/* ADC Controller */
#define AT91C_ID_DMAC0		20	/* DMA Controller 0 */
#define AT91C_ID_DMAC1		21	/* DMA Controller 1 */
#define AT91C_ID_UHPHS		22	/* USB Host High Speed */
#define AT91C_ID_UDPHS		23	/* USB Device High Speed */
#define AT91C_ID_EMAC		24	/* Ethernet MAC 0 */
#define AT91C_ID_LCDC_ISI	25	/* LCD or Image Sensor Interface */
#define AT91C_ID_HSMCI1		26	/* High Speed Multimedia Card Interface 1 */
#define AT91C_ID_EMAC1		27	/* Ethernet MAC 1 */
#define AT91C_ID_SSC		28	/* Serial Synchronous Controller */
#define AT91C_ID_CAN0		29	/* CAN Controller 0 */
#define AT91C_ID_CAN1		30	/* CAN Controller 1 */
#define AT91C_ID_IRQ		31	/* Advanced Interrupt Controller (IRQ) */

/*
 * User Peripherals physical base addresses
 */
#define AT91C_BASE_SPI0		0xf0000000
#define AT91C_BASE_SPI1		0xf0004000
#define AT91C_BASE_HSMCI0	0xf0008000
#define AT91C_BASE_HSMCI1	0xf000c000
#define AT91C_BASE_SSC		0xf0010000
#define AT91C_BASE_CAN0		0xf8000000
#define AT91C_BASE_CAN1		0xf8004000
#define AT91C_BASE_TC0		0xf8008000
#define AT91C_BASE_TC3		0xf800c000
#define AT91C_BASE_TWI0		0xf8010000
#define AT91C_BASE_TWI1		0xf8014000
#define AT91C_BASE_TWI2		0xf8018000
#define AT91C_BASE_USART0	0xf801c000
#define AT91C_BASE_USART1	0xf8020000
#define AT91C_BASE_USART2	0xf8024000
#define AT91C_BASE_EMAC		0xf802c000
#define AT91C_BASE_PWMC		0xf8034000
#define AT91C_BASE_LCDC		0xf8038000
#define AT91C_BASE_UDPHS	0xf803c000
#define AT91C_BASE_UART0	0xf8040000
#define AT91C_BASE_UART1	0xf8044000
#define AT91C_BASE_TSADC	0xf804c000

/*
 * System Peripherals physical base addresses.
 */
#define AT91C_BASE_MATRIX	0xffffde00
#define AT91C_BASE_PMECC	0xffffe000
#define AT91C_BASE_PMERRLOC	0xffffe600
#define AT91C_BASE_DDRSDRC	0xffffe800
#define AT91C_BASE_SMC		0xffffea00
#define AT91C_BASE_DMAC0	0xffffec00
#define AT91C_BASE_DMAC1	0xffffee00
#define AT91C_BASE_AIC		0xfffff000
#define AT91C_BASE_DBGU		0xfffff200
#define AT91C_BASE_PIOA		0xfffff400
#define AT91C_BASE_PIOB		0xfffff600
#define AT91C_BASE_PIOC		0xfffff800
#define AT91C_BASE_PIOD		0xfffffa00
#define AT91C_BASE_PMC		0xfffffc00
#define AT91C_BASE_RSTC		0xfffffe00
#define AT91C_BASE_SHDC		0xfffffe10
#define AT91C_BASE_PITC		0xfffffe30
#define AT91C_BASE_WDT		0xfffffe40
#define AT91C_BASE_SCKCR	0xfffffe50
#define AT91C_BASE_BSCR		0xfffffe54
#define AT91C_BASE_GPBR		0xfffffe60
#define AT91C_BASE_RTC		0xfffffeb0

/* (CCFG) Base Address */
#define AT91C_BASE_CCFG		(0xffffdf18)
/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_BOOT		0x00000000      /* Boot mapped area */
#define AT91C_BASE_ROM		0x00100000      /* Internal ROM base address */
#define AT91C_BASE_SRAM		0x00300000      /* Internal SRAM */
#define AT91C_BASE_SMD		0x00400000      /* Soft Modem */
#define AT91C_BASE_UDPHS_RAM	0x00500000      /* USB Device controller */
#define AT91C_BASE_UHPHS_OHCI	0x00600000      /* UHPHS_OHCI Base address*/
#define AT91C_BASE_UHPHS_EHCI	0x00700000      /* UHPHS_EHCI Base address*/

/*
 * External memory
 */
#define AT91C_BASE_CS0          0x10000000      /* typically NOR */
#define AT91C_BASE_CS1          0x20000000      /* DRAM */
#define AT91C_BASE_CS2          0x30000000
#define AT91C_BASE_CS3          0x40000000      /* typically NAND */
#define AT91C_BASE_CS4          0x50000000
#define AT91C_BASE_CS5          0x60000000

#define AT91C_NUM_PIO		4

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME         "AT91SAM9X5"

#endif	/* #ifndef __AT91SAM9X5_H__ */
