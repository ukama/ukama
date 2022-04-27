/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2013, Atmel Corporation

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
#ifndef __SAMA5D4_H__
#define __SAMA5D4_H__

/*
 * Peripheral identifiers/interrupts.
 */
#define AT91C_ID_FIQ		0	/* FIQ Interrupt ID */
#define AT91C_ID_SYS		1	/* System Controller Interrupt */
#define AT91C_ID_ARM		2	/* Performance Monitor Unit */
#define AT91C_ID_PIT		3	/* Periodic Interval Timer Interrupt */
#define AT91C_ID_WDT		4	/* Watchdog Timer Interrupt */
#define AT91C_ID_PIOD		5	/* Parallel I/O Controller D*/
#define AT91C_ID_USART0		6	/* USART0 */
#define AT91C_ID_USART1		7	/* USART1 */
#define AT91C_ID_XDMAC0		8	/* DMA Controller 0 */
#define AT91C_ID_ICM		9	/* Integrity Check Monitor */
#define AT91C_ID_PKCC		10	/* Public Key Crypto Controller */
#define AT91C_ID_SCI		11	/* Smart Card Interface */
#define AT91C_ID_AES		12	/* Advanced Encryption Standard */
#define AT91C_ID_AESB		13	/* AES bridge */
#define AT91C_ID_TDES		14	/* Triple Data Encrytion Standard */
#define AT91C_ID_SHA		15	/* SHA Signature */
#define AT91C_ID_MPDDRC		16	/* MPDDR controller */
#define AT91C_ID_MATRIX1	17	/* H32MX, 32-bit AHB Matrix */
#define AT91C_ID_MATRIX0	18	/* H64MX, 64-bit AHB Matrix */
#define AT91C_ID_VDEC		19	/* Video Decoder */
#define AT91C_ID_SECUMOD	20	/* Security Module */
#define AT91C_ID_MSADCC		21	/* Magstrioe ADC */
#define AT91C_ID_HSMC		22	/* Multi-bit ECC interrupt */
#define AT91C_ID_PIOA		23	/* Parallel I/O Controller A */
#define AT91C_ID_PIOB		24	/* Parallel I/O Controller B */
#define AT91C_ID_PIOC		25	/* Parallel I/O Controller C */
#define AT91C_ID_PIOE		26	/* Parallel I/O Controller E */
#define AT91C_ID_UART0		27	/* UART0 */
#define AT91C_ID_UART1		28	/* UART1 */
#define AT91C_ID_USART2		29	/* USART2 */
#define AT91C_ID_USART3		30	/* USART3 */
#define AT91C_ID_USART4		31	/* USART4 */
#define AT91C_ID_TWI0		32	/* Two-wire Interface 0 */
#define AT91C_ID_TWI1		33	/* Two-wire Interface 1 */
#define AT91C_ID_TWI2		34	/* Two-wire Interface 2 */
#define AT91C_ID_HSMCI0		35	/* HSMCI0 */
#define AT91C_ID_HSMCI1		36	/* HSMCI1 */
#define AT91C_ID_SPI0		37	/* Serial Peripheral Interface 0 */
#define AT91C_ID_SPI1		38	/* Serial Peripheral Interface 1 */
#define AT91C_ID_SPI2		39	/* Serial Peripheral Interface 2 */
#define AT91C_ID_TC0		40	/* Timer Counter 0 (ch.0,1,2) */
#define AT91C_ID_TC1		41	/* Timer Counter 1 (ch.3,4,5) */
#define AT91C_ID_TC2		42	/* Timer Counter 2 (ch.6,7,8) */
#define AT91C_ID_PWM		43	/* Pulse Width Modulation Controller */
#define AT91C_ID_ADC		44	/* Touch Screen ADC Controller */
#define AT91C_ID_DBGU		45	/* Debug Unit Interrupt */
#define AT91C_ID_UHPHS		46	/* USB Host High Speed */
#define AT91C_ID_UDPHS		47	/* USB Device High Speed */
#define AT91C_ID_SSC0		48	/* Serial Synchronous Controller 0 */
#define AT91C_ID_SSC1		49	/* Serial Synchronous Controller 1 */
#define AT91C_ID_XDMAC1		50	/* DMA Controller 1 */
#define AT91C_ID_LCDC		51	/* LCD Controller */
#define AT91C_ID_ISI		52	/* Image Sensor Interface */
#define AT91C_ID_TRNG		53	/* True Random Number Generator */
#define AT91C_ID_GMAC		54	/* Gigabit Ethernet MAC */
#define AT91C_ID_GMAC1		55	/* Gigabit Ethernet MAC */
#define AT91C_ID_IRQ		56	/* IRQ Interrupt ID */
#define AT91C_ID_SFC		57	/* Fuse Controller */
/* 58 reserved */
#define AT91C_ID_SECURAM	59	/* Secure RAM */
#define AT91C_ID_CTB		60	/* Capacitive Touch Button */
#define AT91C_ID_SMD		61	/* SMD Soft Modem */
#define AT91C_ID_TWI3		62	/* Two-wire Interface 3 */
#define AT91C_ID_CATB		63	/* Capacitive Touch Controller */
#define AT91C_ID_SFR		64	/* Special Function Register */
#define AT91C_ID_AIC		65	/* Advanced Interrupt Controller */
#define AT91C_ID_SAIC		66	/* SAIC */
#define AT91C_ID_L2CC		67	/* L2 Cache Controller */
#define AT91C_ID_PMC		68	/* Power Management Controller */

#define AT91C_ID_COUNTS		(AT91C_ID_PMC + 1)

/*
 * User Peripherals physical base addresses.
 */
#define AT91C_BASE_LCDC		0xf0000000
#define AT91C_BASE_DMAC1	0xf0004000
#define AT91C_BASE_ISI		0xf0008000

#define AT91C_BASE_HSMCI0	0xf8000000
#define AT91C_BASE_UART0	0xf8004000
#define AT91C_BASE_SSC0		0xf8008000
#define AT91C_BASE_PWMC		0xf800c000
#define AT91C_BASE_SPI0		0xf8010000
#define AT91C_BASE_TWI0		0xf8014000
#define AT91C_BASE_TWI1		0xf8018000
#define AT91C_BASE_TC012	0xf801c000
#define AT91C_BASE_GMAC		0xf8020000
#define AT91C_BASE_TWI2		0xf8024000

#define AT91C_BASE_HSMCI1	0xfc000000
#define AT91C_BASE_UART1	0xfc004000
#define AT91C_BASE_USART2	0xfc008000
#define AT91C_BASE_USART3	0xfc00c000
#define AT91C_BASE_USART4	0xfc010000
#define AT91C_BASE_SSC1		0xfc014000
#define AT91C_BASE_SPI1		0xfc018000
#define AT91C_BASE_SPI2		0xfc01c000
#define AT91C_BASE_TC345	0xf8020000
#define AT91C_BASE_TC678	0xf8024000
#define AT91C_BASE_GMAC1	0xfc028000

#define AT91C_BASE_UDPHS	0xfc02c000
#define AT91C_BASE_TRNG		0xfc030000
#define AT91C_BASE_ADC		0xfc034000
#define AT91C_BASE_TWI3		0xfc038000

#define AT91C_BASE_DBGU		0xfc069000
#define AT91C_BASE_PIOA		0xfc06a000
#define AT91C_BASE_PIOB		0xfc06b000
#define AT91C_BASE_PIOC		0xfc06c000
#define AT91C_BASE_PIOE		0xfc06d000
#define AT91C_BASE_AIC		0xfc06e000

/* Always Secure Mapping */
#define AT91C_BASE_PKCC		0xf000c000
#define AT91C_BASE_MPDDRC	0xf0010000
#define AT91C_BASE_DMAC0	0xf0014000
#define AT91C_BASE_PMC		0xf0018000
#define AT91C_BASE_MATRIX64	0xf001c000
#define AT91C_BASE_AESB		0xf0020000

#define AT91C_BASE_SFR		0xf8028000
#define AT91C_BASE_USART0	0xf802c000
#define AT91C_BASE_USART1	0xf8030000
#define AT91C_BASE_SCI		0xf8034000
#define AT91C_BASE_CTB		0xf8038000

#define	AT91C_BASE_MAGSADC	0xfc03c000
#define	AT91C_BASE_ICM		0xfc040000
#define	AT91C_BASE_AES		0xfc044000

#define	AT91C_BASE_TDES		0xfc04c000
#define	AT91C_BASE_SHA		0xfc050000
#define	AT91C_BASE_MATRIX32	0xfc054000
#define	AT91C_BASE_SECURAM	0xfc058000
#define	AT91C_BASE_SMC		0xfc05c000
#define	AT91C_BASE_SFC		0xfc060000

#define	AT91C_BASE_PIOD		0xfc068000
#define	AT91C_BASE_SECUMOD	0xfc068200
#define	AT91C_BASE_SAIC		0xfc068400
#define	AT91C_BASE_RSTC		0xfc068600
#define	AT91C_BASE_SHDC		0xfc068610
#define	AT91C_BASE_PITC		0xfc068630
#define	AT91C_BASE_WDT		0xfc068640
#define	AT91C_BASE_SCKCR	0xfc068650
#define	AT91C_BASE_RTCC		0xfc0686b0

/*
 * System Peripherals physical base addresses.
 */
#define AT91C_BASE_SYS		0xffffc000
/* Reserved */
/*
 * Internal Memory common on all these SoCs
 */
#define AT91C_BASE_ROM		0x00000000	/* ROM */
#define AT91C_BASE_NFC_SRAM	0x00100000	/* NFC SRAM */
#define AT91C_BASE_SRAM		0x00200000	/* SRAM0 */
#define AT91C_BASE_VDEC		0x00300000	/* VDEC */
#define AT91C_BASE_UDPHS_SRAM	0x00400000	/* UDPHS RAM */
#define AT91C_BASE_UHP_OHCI	0x00500000	/* UHP OHCI */
#define AT91C_BASE_UHP_EHCI	0x00600000	/* UHP EHCI */
#define AT91C_BASE_AXI_MATRIX	0x00700000	/* AXI Maxtrix */
#define AT91C_BASE_DAP		0x00800000	/* DAP */
#define AT91C_BASE_SMD		0x00900000	/* SMD */
#define AT91C_BASE_L2CC		0x00A00000	/* L2CC */

/*
 * External memory
 */
#define AT91C_BASE_CS0		0x10000000
#define AT91C_BASE_DDRCS	0x20000000
#define AT91C_BASE_DDRCS_AES	0x40000000
#define AT91C_BASE_CS1		0x60000000
#define AT91C_BASE_CS2		0x70000000
#define AT91C_BASE_CS3		0x80000000
#define AT91C_BASE_NFC_CMDR	0x90000000

/*
 * Other misc defines
 */
#define AT91C_BASE_PMECC	(AT91C_BASE_SMC + 0x70)
#define AT91C_BASE_PMERRLOC	(AT91C_BASE_SMC + 0x500)

#define AT91_PMECC		(AT91C_BASE_PMECC - AT91C_BASE_SYS)
#define AT91_PMERRLOC		(AT91C_BASE_PMERRLOC - AT91C_BASE_SYS)

#define ATMEL_BASE_SMC		(AT91C_BASE_SMC + 0x600)
#define AT91C_SYS_SCKCR		AT91C_BASE_SCKCR

#define AT91C_NUM_PIO		5
#define	AT91C_NUM_TWI		4

/* AICREDIR Unlock Key */
#define	AICREDIR_KEY		0x5F67B102

/*
 * Matrix Slaves ID
 */

/* MATRIX0(H64MX) Matrix Slaves */
/* Bridge from H64MX to AXIMX (Internal ROM, Cryto Library, PKCC RAM) */
#define	H64MX_SLAVE_BRIDGE_TO_AXIMX	0
#define	H64MX_SLAVE_PERI_BRIDGE		1	/* H64MX Peripheral Bridge */
#define	H64MX_SLAVE_VIDEO_DECODER	2	/* Video Decoder */
#define	H64MX_SLAVE_DDR2_PORT_0		3	/* DDR2 Port0-AESOTF */
#define	H64MX_SLAVE_DDR2_PORT_1		4	/* DDR2 Port1 */
#define	H64MX_SLAVE_DDR2_PORT_2		5	/* DDR2 Port2 */
#define	H64MX_SLAVE_DDR2_PORT_3		6	/* DDR2 Port3 */
#define	H64MX_SLAVE_DDR2_PORT_4		7	/* DDR2 Port4 */
#define	H64MX_SLAVE_DDR2_PORT_5		8	/* DDR2 Port5 */
#define	H64MX_SLAVE_DDR2_PORT_6		9	/* DDR2 Port6 */
#define	H64MX_SLAVE_DDR2_PORT_7		10	/* DDR2 Port7 */
#define	H64MX_SLAVE_INTERNAL_SRAM	11	/* Internal SRAM 128K */
#define	H64MX_SLAVE_BRIDGE_TO_H32MX	12	/* Bridge from H64MX to H32MX */

/* MATRIX1(H32MX) Matrix Slaves */
#define	H32MX_BRIDGE_TO_H64MX		0	/* Bridge from H32MX to H64MX */
#define	H32MX_PERI_BRIDGE_0		1	/* H32MX Peripheral Bridge 0 */
#define	H32MX_PERI_BRIDGE_1		2	/* H32MX Peripheral Bridge 1 */
#define	H32MX_EXTERNAL_EBI		3	/* External Bus Interface */
#define	H32MX_NFC_CMD_REG		3	/* NFC command Register */
#define	H32MX_NFC_SRAM			4	/* NFC SRAM */
/* USB Device High Speed Dual Port RAM (DPR)
 * USB Host OHCI registers
 * USB Host EHCI registers
 */
#define	H32MX_USB			5
#define	H32MX_SMD			6	/* Soft Modem(SMD) */

/*
 * SoC specific defines
 */
#define AT91C_CPU_NAME		"SAMA5D4"

#endif /* #ifndef __SAMA5D4_H__ */
