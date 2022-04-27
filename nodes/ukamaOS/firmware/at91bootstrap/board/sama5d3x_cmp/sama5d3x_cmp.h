/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2015, Atmel Corporation
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
#ifndef __SAMA5D3X_CMP_H__
#define __SAMA5D3X_CMP_H__

/*
 * PMC Setting
 *
 * The main oscillator is enabled as soon as possible in the lowlevel_clock_init
 * and MCK is switched on the main oscillator.
 */
#define BOARD_MAINOSC		12000000

#if defined(CONFIG_BUS_SPEED_133MHZ)

#define MASTER_CLOCK		132000000

#if defined(CONFIG_CPU_CLK_528MHZ)

/* PCK = 528MHz, MCK = 132MHz */
#define PLLA_MULA		43
#define BOARD_PCK		((unsigned long)(BOARD_MAINOSC * \
							(PLLA_MULA + 1)))
#define BOARD_MCK		((unsigned long)((BOARD_MAINOSC * \
							(PLLA_MULA + 1)) / 4))

#define BOARD_CKGR_PLLA		(AT91C_CKGR_SRCA | AT91C_CKGR_OUTA_0)
#define BOARD_PLLACOUNT		(0x3F << 8)
#define BOARD_MULA		((AT91C_CKGR_MULA << 2) & (PLLA_MULA << 18))
#define BOARD_DIVA		(AT91C_CKGR_DIVA & 1)

#define BOARD_PRESCALER_MAIN_CLOCK	(AT91C_PMC_MDIV_4 \
					| AT91C_PMC_CSS_MAIN_CLK)

#define BOARD_PRESCALER_PLLA		(AT91C_PMC_MDIV_4 \
					| AT91C_PMC_CSS_PLLA_CLK)

#else
#error "No cpu clock provided!"
#endif /* #if defined(CONFIG_CPU_CLK_528MHZ) */

#else
#error "No main clock provided!"
#endif /* #if defined(CONFIG_BUS_SPEED_133MHZ) */

#define PLLA_SETTINGS		(BOARD_CKGR_PLLA \
				| BOARD_PLLACOUNT \
				| BOARD_MULA \
				| BOARD_DIVA)

/*
 * DataFlash Settings
 */
#define CONFIG_SYS_SPI_CLOCK	AT91C_SPI_CLK
#define CONFIG_SYS_SPI_MODE	SPI_MODE3

#if defined(CONFIG_SPI_BUS0)
#define CONFIG_SYS_BASE_SPI	AT91C_BASE_SPI0
#elif defined(CONFIG_SPI_BUS1)
#define CONFIG_SYS_BASE_SPI	AT91C_BASE_SPI1
#endif

#if (AT91C_SPI_PCS_DATAFLASH == AT91C_SPI_PCS0_DATAFLASH)
#define CONFIG_SYS_SPI_PCS	AT91C_PIN_PD(13)
#elif (AT91C_SPI_PCS_DATAFLASH == AT91C_SPI_PCS1_DATAFLASH)
#define CONFIG_SYS_SPI_PCS	AT91C_PIN_PD(14)
#endif

/*
 * NandFlash Settings
 */
#define CONFIG_SYS_NAND_BASE		AT91C_BASE_CS3
#define CONFIG_SYS_NAND_MASK_ALE	(1 << 21)
#define CONFIG_SYS_NAND_MASK_CLE	(1 << 22)

#undef CONFIG_SYS_NAND_ENABLE_PIN

#define CONFIG_LOOKUP_TABLE_ALPHA_OFFSET	0x14000
#define CONFIG_LOOKUP_TABLE_INDEX_OFFSET	0x10000

#define CONFIG_LOOKUP_TABLE_ALPHA_OFFSET_1024	0x20000
#define CONFIG_LOOKUP_TABLE_INDEX_OFFSET_1024	0x18000

/*
 * MCI Settings
 */
#define CONFIG_SYS_BASE_MCI	AT91C_BASE_HSMCI0

/*
 * 1-Wire Pin
 */
#define CONFIG_SYS_ONE_WIRE_PIN		AT91C_PIN_PE(25)

/*
 * Others
 */
#define CONFIG_SYS_HDMI_RESET_PIN	AT91C_PIN_PC(31)

#endif /* #ifndef __SAMA5D3X_CMP_H__ */
