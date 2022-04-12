/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 * Copyright (c) 2018, Microchip Technology Inc. and its subsidiaries
 *
 * All rights reserved.
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

#ifndef __SAMA7G5EK_H__
#define __SAMA7G5EK_H__

#ifdef CONFIG_FLEXCOM
#include "flexcom.h"
#endif

#define BOARD_PRESCALER_CPUPLL	(AT91C_PMC_PRES_CLK | \
				AT91C_PMC_CSS_CPUPLL_CLK)

#define BOARD_PRESCALER_MCK1	(AT91C_MCR_MASTER_DIV2 | \
				AT91C_MCR_CSS_SYSPLL_CLK | \
				AT91C_MCR_EN)

#define BOARD_PRESCALER_MCK1_CLEAN	(AT91C_MCR_MASTER_DIV4 | \
					AT91C_MCR_CSS_MCK0_CLK | \
					AT91C_MCR_EN)

#define BOARD_PRESCALER_MCK2	(AT91C_MCR_MASTER_DIV1 | \
				AT91C_MCR_CSS_DDRPLL_CLK | \
				AT91C_MCR_EN)

#define BOARD_PRESCALER_MCK3	(AT91C_MCR_MASTER_DIV1 | \
				AT91C_MCR_CSS_IMGPLL_CLK | \
				AT91C_MCR_EN)

#define BOARD_PRESCALER_MCK4	(AT91C_MCR_MASTER_DIV1 | \
				AT91C_MCR_CSS_SYSPLL_CLK | \
				AT91C_MCR_EN)

#define BOARD_PRESCALER_MCK4_CLEAN	(AT91C_MCR_MASTER_DIV4 | \
					AT91C_MCR_CSS_MCK0_CLK | \
					AT91C_MCR_EN)

#define MASTER_CLOCK			200000000
#define BOARD_MAINOSC			24000000

#define CONFIG_SYS_BASE_TZC400		AT91C_BASE_TZC400
#define TZC400_BASE			CONFIG_SYS_BASE_TZC400

#ifdef CONFIG_FLEXCOM
#define USART_BASE			(AT91C_BASE_FLEXCOM3 + AT91C_OFFSET_FLEXCOM_USART)
#endif

#define CONFIG_SYS_DBGU_RXD_PIN		AT91C_PIN_PD(17)
#define CONFIG_SYS_DBGU_TXD_PIN		AT91C_PIN_PD(16)

#define CONFIG_SYS_BASE_PUBL		AT91C_BASE_PUBL
#define CONFIG_SYS_BASE_UMCTL2		AT91C_BASE_UMCTL2
#define CONFIG_SYS_BASE_UMCTL2_MP	AT91C_BASE_UMCTL2_MP

#if defined(CONFIG_SDHC0)
#define CONFIG_SYS_BASE_SDHC		AT91C_BASE_SDMMC0
#define	CONFIG_SYS_ID_SDHC		AT91C_ID_SDMMC0
#elif defined(CONFIG_SDHC1)
#define CONFIG_SYS_BASE_SDHC		AT91C_BASE_SDMMC1
#define	CONFIG_SYS_ID_SDHC		AT91C_ID_SDMMC1
#endif

#ifdef CONFIG_MT41K256M16TW_107
#define _tRFC		350UL		/* Refresh to Refresh, ns */
#define _tREFI		7800UL		/* Refresh Intervals, ns */
#define _tWR		15UL		/* Write Recovery, ns */
#define _tRP		14UL		/* Row Precharge command, ns */
#define _tRP_ps		13125UL		/* Row Precharge command, ps */
#define _tRCD		14UL		/* Row to Column delay, ns */
#define _tRCD_ps	13125UL		/* Row to Column delay, ps */
#define _tCCD		4UL		/* Column to Column delay, ns */
#define _tRAS		38UL		/* Row Active Strobe, ns */
#define _tRASMAX	70000UL		/* RAS max, ns */
#define _tRC_ps		50625UL		/* Row Cycle, ps */
#define _tFAW		50UL		/* Four Activation Window, ns */
#define _tPRECKE	500000UL
#define _tPOSTCKE	0UL		/* Unused with DDR3 */
#define _CL		7UL		/* CAS Latency, Clock cycles */
#define _CWL		6UL		/* CAS Write Latency, Clock cycles */
#define _AL		0UL		/* Additive Latency, Clock cycles */
#define _TZQOPER	256UL		/* TZQ Long Calibration, Clock cycles */
#define _TZQCS		64UL		/* TZQ Short Calibration, Clock cycles */
#define _MRD		4UL		/* Mode Register Delay, Clock cycles */
#endif

#ifdef CONFIG_MT47H128M16RT_25E_C
#define _tRFC		195UL		/* Refresh to Refresh, ns */
#define _tREFI		7800UL		/* Refresh Intervals, ns */
#define _tWR		15UL		/* Write Recovery, ns */
#define _tRP		13UL		/* Row Precharge command, ns */
#define _tRP_ps		12500UL		/* Row Precharge command, ps */
#define _tRCD		13UL		/* Row to Column delay, ns */
#define _tRCD_ps	12500UL		/* Row to Column delay, ps */
#define _tCCD		2UL		/* Column to Column delay, ns */
#define _tRAS		40UL		/* Row Active Strobe, ns */
#define _tRASMAX	70000UL		/* RAS max, ns */
#define _tRC_ps		55000UL		/* Row Cycle, ps */
#define _tFAW		50UL		/* Four Activation Window, ns */
#define _tPRECKE	200000UL
#define _tPOSTCKE	400UL
#define _CL		5UL		/* CAS Latency, Clock cycles */
#define _CWL		4UL		/* CAS Write Latency, Clock cycles */
#define _AL		0UL		/* Additive Latency, Clock cycles */
#define _TZQOPER	256UL		/* TZQ Long Calibration, Clock cycles */
#define _TZQCS		64UL		/* TZQ Short Calibration, Clock cycles */
#define _MRD		2UL		/* Mode Register Delay, Clock cycles */
#endif

#ifdef CONFIG_IS43LD16128B_25BLI
#define _tRFC		130UL		/* Refresh to Refresh, ns */
#define _tREFI		3900UL		/* Refresh Intervals, ns */
#define _tWR		15UL		/* Write Recovery, ns */
#define _tRP		21UL		/* Row Precharge command, ns */
#define _tRP_ps		21000UL		/* Row Precharge command, ps */
#define _tRCD		18UL		/* Row to Column delay, ns */
#define _tRCD_ps	18000UL		/* Row to Column delay, ps */
#define _tCCD		2UL		/* Column to Column delay, ns */
#define _tRAS		42UL		/* Row Active Strobe, ns */
#define _tRASMAX	70000UL		/* RAS max, ns */
#define _tRC_ps		((_tRAS + _tRP)*1000)		/* Row Cycle, ps */
#define _tFAW		50UL		/* Four Activation Window, ns */
#define _tPOSTCKE	200000UL	/* Post CKE, ns */
#define _CL		0UL		/* CAS Latency, Clock cycles, unused in LPDDR2 */
#define _CWL		0UL		/* CAS Write Latency, Clock cycles, unused in LPDDR2 */
#define _AL		0UL		/* Additive Latency, Clock cycles */
#define _RL		6UL		/* For LPDDR2, Read Latency, Clock cycles */
#define _WL		3UL		/* For LPDDR2, Write Latency, Clock cycles */
#define _TZQOPER	NS_TO_CYCLES_UP(360UL)		/* TZQ Long Calibration, Clock cycles */
#define _TZQCS		NS_TO_CYCLES_UP(90UL)		/* TZQ Short Calibration, Clock cycles */
#define _MRD		0UL		/* Mode Register Delay, Clock cycles , unused in LPDDR2*/
#define _tDQSCK_MIN	2500UL
#define _tDQSCK_MAX	5500UL
#endif

#ifdef CONFIG_MT52L256M32D1PF_107_WT_B
#define _tRFC		210UL		/* Refresh to Refresh, ns */
#define _tREFI		3900UL		/* Refresh Intervals, ns */
#define _tWR		MAX(15UL, 3 * CYCLE_IN_NS)		/* Write Recovery, ns */
#define _tRP		MAX(21UL, 3 * CYCLE_IN_NS)		/* Row Precharge command, ns */
#define _tRP_ps		MAX(21000UL, 3 * CYCLE_IN_PS)		/* Row Precharge command, ps */
#define _tRCD		MAX(18UL, 3 * CYCLE_IN_NS)		/* Row to Column delay, ns */
#define _tRCD_ps	MAX(18000UL, 3 * CYCLE_IN_PS)		/* Row to Column delay, ps */
#define _tCCD		4UL		/* Column to Column delay, ns */
#define _tRAS		MAX(42UL, 3 * CYCLE_IN_NS)		/* Row Active Strobe, ns */
#define _tRASMAX	70200UL		/* RAS max, ns */
#define _tRC_ps		((_tRAS + _tRP) * 1000)		/* Row Cycle, ps */
#define _tFAW		MAX(50UL, 8 * CYCLE_IN_NS)		/* Four Activation Window, ns */
#define _tPOSTCKE	200000UL	/* Post CKE, ns */
#define _CL		0UL		/* CAS Latency, Clock cycles, unused in LPDDR2 */
#define _CWL		0UL		/* CAS Write Latency, Clock cycles, unused in LPDDR2 */
#define _AL		0UL		/* Additive Latency, Clock cycles */
#define _RL		8UL		/* For LPDDR3, Read Latency, Clock cycles */
#define _WL		4UL		/* For LPDDR3, Write Latency, Clock cycles */
#define _TZQOPER	NS_TO_CYCLES_UP(360UL)		/* TZQ Long Calibration, Clock cycles */
#define _TZQCS		NS_TO_CYCLES_UP(90UL)		/* TZQ Short Calibration, Clock cycles */
#define _MRD		MAX(NS_TO_CYCLES_UP(14UL), 10UL)		/* Mode Register Delay, Clock cycles */
#define _tDQSCK_MIN	2500UL
#define _tDQSCK_MAX	5500UL
#endif
#endif
