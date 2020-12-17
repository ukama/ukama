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
#ifndef __SAMA5D3X_H__
#define __SAMA5D3X_H__

/*
 * Peripheral identifiers/interrupts.
 */
#define AT91C_ID_AIC	0	/* Advanced Interrupt Controller (AIC) */
#define AT91C_ID_SYS	1	/* System Controller Interrupt */
#define AT91C_ID_DBGU	2	/* Debug Unit Interrupt */
#define AT91C_ID_PIT	3	/* Periodic Interval Timer Interrupt */
#define AT91C_ID_WDT	4	/* Watchdog Timer Interrupt */
#define AT91C_ID_SMC	5	/* Multi-bit ECC Interrupt */
#define AT91C_ID_PIOA	6	/* Parallel I/O Controller A */
#define AT91C_ID_PIOB	7	/* Parallel I/O Controller B */
#define AT91C_ID_PIOC	8	/* Parallel I/O Controller C */
#define AT91C_ID_PIOD	9	/* Parallel I/O Controller D */
#define AT91C_ID_PIOE	10	/* Parallel I/O Controller E */
#define AT91C_ID_SMD	11	/* SMD Soft Modem */
#define AT91C_ID_USART0	12	/* USART 0 */
#define AT91C_ID_USART1	13	/* USART 1 */
#define AT91C_ID_USART2	14	/* USART 2 */
#define AT91C_ID_USART3	15	/* USART 3 */
#define AT91C_ID_UART0	16	/* UART 0 */
#define AT91C_ID_UART1	17	/* UART 1 */
#define AT91C_ID_TWI0	18	/* Two-Wire Interface 0 */
#define AT91C_ID_TWI1	19	/* Two-Wire Interface 1 */
#define AT91C_ID_TWI2	20	/* Two-Wire Interface 2 */
#define AT91C_ID_HSMCI0	21	/* High Speed Multimedia Card Interface 0 */
#define AT91C_ID_HSMCI1	22	/* High Speed Multimedia Card Interface 1 */
#define AT91C_ID_HSMCI2	23	/* High Speed Multimedia Card Interface 2 */
#define AT91C_ID_SPI0	24	/* Serial Peripheral Interface 0 */
#define AT91C_ID_SPI1	25	/* Serial Peripheral Interface 1 */
#define AT91C_ID_TC0	26	/* Timer Counter 0 (ch.0,1,3) */
#define AT91C_ID_TC1	27	/* Timer Counter 1 (ch.3,4,5) */
#define AT91C_ID_PWM	28	/* Pulse Width Modulation */
#define AT91C_ID_ADC	29	/* Touch Screen ADC Controller */
#define AT91C_ID_DMAC0	30	/* DMA Controller 0 */
#define AT91C_ID_DMAC1	31	/* DMA Controller 1 */
#define AT91C_ID_UHPHS	32	/* USB Host High Speed */
#define AT91C_ID_UDPHS	33	/* USB Device High Speed */
#define AT91C_ID_GMAC	34	/* Gigabit Ethernet MAC */
#define AT91C_ID_EMAC	35	/* Ethernet MAC */
#define AT91C_ID_LCDC	36	/* LCD Controller */
#define AT91C_ID_ISI	37	/* Image Sensor Interface */
#define AT91C_ID_SSC0	38	/* Serial Synchronous Controller 0 */
#define AT91C_ID_SSC1	39	/* Serial Synchronous Controller 1 */
#define AT91C_ID_CAN0	40	/* CAN Controller 0 */
#define AT91C_ID_CAN1	41	/* CAN Controller 1 */
#define AT91C_ID_SHA	42	/* Secure Hash Algorithm */
#define AT91C_ID_AES	43	/* Anvanced Encryption Standard */
#define AT91C_ID_TDES	44	/* Triple Data Encryption Standard */
#define AT91C_ID_TRNG	45	/* True Random Number Generator */
#define AT91C_ID_ARM	46	/* Performance Monitor Unit */
#define AT91C_ID_IRQ	47	/* Advanced Interrupt Controller (IRQ) */
#define AT91C_ID_FUSE	48	/* Fuse Controller */
#define AT91C_ID_MPDDRC	49	/* MPDDR Controller */

/*
 * User Peripherals physical base addresses.
 */
#define AT91C_BASE_HSMCI0	0xf0000000
#define AT91C_BASE_SPI0		0xf0004000
#define AT91C_BASE_SSC0		0xf0008000
#define AT91C_BASE_CAN0		0xf000c000
#define AT91C_BASE_TC012	0xf0010000
#define AT91C_BASE_TWI0		0xf0014000
#define AT91C_BASE_TWI1		0xf0018000
#define AT91C_BASE_USART0	0xf001c000
#define AT91C_BASE_USART1	0xf0020000
#define AT91C_BASE_UART0	0xf0024000
#define AT91C_BASE_GMAC		0xf0028000
#define AT91C_BASE_PWMC		0xf002c000
#define AT91C_BASE_LCDC		0xf0030000
#define AT91C_BASE_ISI		0xf0034000
#define AT91C_BASE_SFR		0xf0038000
/* Reserved */
#define AT91C_BASE_HSMCI1	0xf8000000
#define AT91C_BASE_HSMCI2	0xf8004000
#define AT91C_BASE_SPI1		0xf8008000
#define AT91C_BASE_SSC1		0xf800c000
#define AT91C_BASE_CAN1		0xf8010000
#define AT91C_BASE_TC345	0xf8014000
#define AT91C_BASE_TSADC	0xf8018000
#define AT91C_BASE_TWI2		0xf801c000
#define AT91C_BASE_USART2	0xf8020000
#define AT91C_BASE_USART3	0xf8024000
#define AT91C_BASE_UART1	0xf8028000
#define AT91C_BASE_EMAC		0xf802c000
#define AT91C_BASE_UDPHS	0xf8030000
#define AT91C_BASE_SHA		0xf8034000
#define AT91C_BASE_AES		0xf8038000
#define AT91C_BASE_TDES		0xf803c000
#define AT91C_BASE_TRNG		0xf8044000

/*
 * System Peripherals physical base addresses.
 */
#define AT91C_BASE_SMC		0xffffc000
/* Reserved */
#define AT91C_BASE_FUSE		0xffffe400
#define AT91C_BASE_DMAC0	0xffffe600
#define AT91C_BASE_DMAC1	0xffffe800
#define AT91C_BASE_MPDDRC	0xffffea00
#define AT91C_BASE_MATRIX	0xffffec00
#define AT91C_BASE_DBGU		0xffffee00
#define AT91C_BASE_AIC		0xfffff000
#define AT91C_BASE_PIOA		0xfffff200
#define AT91C_BASE_PIOB		0xfffff400
#define AT91C_BASE_PIOC		0xfffff600
#define AT91C_BASE_PIOD		0xfffff800
#define AT91C_BASE_PIOE		0xfffffa00
#define AT91C_BASE_PMC		0xfffffc00
#define AT91C_BASE_RSTC		0xfffffe00
#define AT91C_BASE_SHDC		0xfffffe10
/* Reserved 0xffff fe20 */
#define AT91C_BASE_PITC		0xfffffe30
#define AT91C_BASE_WDT		0xfffffe40
#define AT91C_BASE_SCKCR	0xfffffe50
#define AT91C_BASE_BSC		0xfffffe54
#define AT91C_BASE_GPBR		0xfffffe60
/* Reserved 0xffff fe70 */
#define AT91C_BASE_RTCC		0xfffffeb0
/* Reserved 0xffff fee0 */
/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_BOOT		0x00000000	/* Boot mapped area */
#define AT91C_BASE_ROM		0x00100000	/* Internal ROM base address */
#define AT91C_BASE_NFC_SRAM	0x00200000	/* Internal NFC SRAM */
#define AT91C_BASE_SRAM0	0x00300000	/* Internal SRAM0 */
#define AT91C_BASE_SRAM1	0x00310000	/* Internal SRAM1 */
#define AT91C_BASE_SMD		0x00400000	/* SMD */
#define AT91C_BASE_UDPH_SRAM	0x00500000	/* UDPH SRAM */
#define AT91C_BASE_UHP_OHCI	0x00600000	/* UHP OHCI */
#define AT91C_BASE_UHP_EHCI	0x00700000	/* UHP EHCI */
#define AT91C_BASE_AXI_MATRIX	0x00800000	/* AXI Maxtrix */
#define AT91C_BASE_DAP		0x00900000	/* DAP */

/*
 * External memory
 */
#define AT91C_BASE_CS0		0x10000000
#define AT91C_BASE_DDRCS	0x20000000
#define AT91C_BASE_CS1		0x40000000
#define AT91C_BASE_CS2		0x50000000
#define AT91C_BASE_CS3		0x60000000
#define AT91C_BASE_NFC_CMD	0x70000000

/*
 * Other misc defines
 */
#define AT91C_BASE_PMECC	(AT91C_BASE_SMC + 0x70)
#define AT91C_BASE_PMERRLOC	(AT91C_BASE_SMC + 0x500)

#define ATMEL_BASE_SMC		(AT91C_BASE_SMC + 0x600)
#define AT91C_SYS_SCKCR		AT91C_BASE_SCKCR

#define AT91C_NUM_PIO		5
#define	AT91C_NUM_TWI		3

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME		"SAMA5D3X"

#define BACKUP_REGISTER_BOOT_MODE_R4 (AT91C_BASE_GPBR + 0x4)

#endif /* #ifndef __SAMA5D3X_H__ */
