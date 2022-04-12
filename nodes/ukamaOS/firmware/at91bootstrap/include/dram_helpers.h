/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 * Copyright (c) 2020, Microchip Technology Inc. and its subsidiaries
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

#ifndef __DRAM_HELPERS_H__
#define __DRAM_HELPERS_H__

#define WAIT_WHILE_COND(x, time)						\
{ 									\
	unsigned int timeout = (time);					\
	while ((x) && (timeout--));					\
									\
	if (timeout == 0xFFFFFFFF) {					\
		dbg_printf("Timed out while waiting for  " #x "\n");	\
		return -1;						\
	}								\
}

#define DIV_ROUND_UP(x, y)	((x) + (y) - 1)/(y)
#define DIV_ROUND_DOWN(x, y)	((x) / (y))

/*
 * ns * 10^-9 (ns to seconds)  * CLOCK * 10^6 (beats/sec in Mhz, to seconds)
 * / 1000 = clock cycles
 * ps * 10^-12 (ps to seconds) * CLOCK * 10^6 (beats/sec in Mhz, to seconds)
 * /10^6 = clock cycles
 */
#define NS_TO_CYCLES_UP(ns)	DIV_ROUND_UP(((ns) * CONFIG_MEM_CLOCK), 1000UL)
#define NS_TO_CYCLES_DOWN(ns)	DIV_ROUND_DOWN(((ns) * CONFIG_MEM_CLOCK), 1000UL)
#define PS_TO_CYCLES_UP(ps)	DIV_ROUND_UP(((ps) * CONFIG_MEM_CLOCK), 1000000UL)
#define PS_TO_CYCLES_DOWN(ps)	DIV_ROUND_DOWN(((ps) * CONFIG_MEM_CLOCK), 1000000UL)
/*
 * 1/ frequency = period.
 * CLOCK is in Mhz
 * 1/ (CLOCK * 10^6) = seconds. To obtain picoseconds, multiply by 10^12
 * 1 * 10 ^12 / (CLOCK * 10^6) = 10^6/ CLOCK (picoseconds)
 */
#define CYCLE_IN_PS		(1000000UL / CONFIG_MEM_CLOCK)
#define CYCLE_IN_NS		DIV_ROUND_UP(CYCLE_IN_PS, 1000UL)

#define MAX(x, y)	((x) > (y) ? (x) : (y))
#define MIN(x, y)	((x) > (y) ? (y) : (x))

struct dram_timings
{
	unsigned long	tRFC;		/* Refresh to Refresh, ns */
	unsigned long	tREFI;		/* Refresh Intervals, ns */
	unsigned long	tWR;		/* Write Recovery, ns */
	unsigned long	tRP;		/* Row Precharge command, ns */
	unsigned long	tRP_ps;		/* Row Precharge command, ps */
	unsigned long	tRCD;		/* Row to Column delay, ns */
	unsigned long	tRCD_ps;	/* Row to Column delay, ps */
	unsigned long	tCCD;		/* Column to Column delay, ns */
	unsigned long	tRAS;		/* Row Active Strobe, ns */
	unsigned long	tRASMAX;	/* Row Active Strobe Max value, ns */
	unsigned long	tRC_ps;		/* Row Cycle, ps */
	unsigned long	tFAW;		/* Four Activation Window, ns */
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long	tPRECKE;
#endif
	unsigned long	tPOSTCKE;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long	tDQSCK_MIN;
	unsigned long	tDQSCK_MAX;
#endif
	unsigned long	CL;		/* CAS Latency, Clock cycles */
	unsigned long	CWL;		/* CAS Write Latency, Clock cycles */
	unsigned long	AL;		/* Additive Latency, Clock cycles */
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long	RL;		/* Read Latency, Clock cycles, LPDDR2 */
	unsigned long	WL;		/* Write Latency, Clock cycles, LPDDR2 */
#endif
	unsigned long	TZQOPER;	/* Calibration time, Clock cycles */
	unsigned long	TZQCS;		/* Calibration time short, Clock cycles */
	unsigned long	MRD;		/* Mode Register Delay, Clock cycles */

/* All timings below starting with a caps T should be in clock cycles */
#if defined(CONFIG_LPDDR3)
#define TRTP		MAX(4, PS_TO_CYCLES_UP(7500UL))		/* Read to Precharge, Clock Cycles */
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
#define TRTP		MAX(2, PS_TO_CYCLES_UP(7500UL))		/* Read to Precharge, Clock Cycles */
#endif
#define TWTR		MAX(4, PS_TO_CYCLES_UP(7500UL))		/* Write to Read, Clock Cycles */
#define TRP		PS_TO_CYCLES_UP(tRP_ps)			/* Row Precharge delay, Clock Cycles */
#define TRCD		PS_TO_CYCLES_UP(tRCD_ps)		/* Row to Column delay, Clock Cycles */
#define TRAS		NS_TO_CYCLES_UP(tRAS)			/* Row Active Strobe, Clock Cycles */
#define TRRD		MAX(4, NS_TO_CYCLES_UP(10UL))		/* Row to Row Delay, Clock Cycles */
#define TRC		PS_TO_CYCLES_UP(tRC_ps)			/* Row Cycle, Clock Cycles */
#define TFAW		NS_TO_CYCLES_UP(tFAW)			/* Four Activation Window, Clock Cycles */
#define TMRD		MRD					/* Mode Register Delay, Clock Cycles */
#if defined(CONFIG_LPDDR2)
#define TMRW		5UL					/* Mode Register Wait, Clock Cycles */
#endif
#if defined(CONFIG_LPDDR3)
#define TMRW		10UL					/* Mode Register Wait, Clock Cycles */
#endif
#define TCCD		NS_TO_CYCLES_UP(tCCD)			/* Column to Column Delay, Clock Cycles */
#define TMOD		MAX(12, NS_TO_CYCLES_UP(15UL))		/* Mode Update Delay, Clock Cycles */
#define TRFC		NS_TO_CYCLES_UP(tRFC)			/* Refresh Counter, Clock Cycles */
#define TXS		MAX(5, NS_TO_CYCLES_UP(tRFC + 10UL))	/* Self-refresh Exit Delay, Clock Cycles */
#define TXSDLL		512UL
#define TXSR		NS_TO_CYCLES_UP(tRFC + 10UL)
#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
#define TXP		MAX(3, PS_TO_CYCLES_UP(7500UL))		/* Power down exit Delay, Clock Cycles */
#define TXPDLL		MAX(10, NS_TO_CYCLES_UP(24UL))
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
#define TXP		MAX(3, PS_TO_CYCLES_UP(7500UL))		/* Power down exit Delay, Clock Cycles */
#endif
#define TXPR		MAX(5, NS_TO_CYCLES_UP(tRFC + 10))	/* Reset Clock exit time, Clock Cycles */
#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
#define TCKE		MAX(3, PS_TO_CYCLES_UP(7500UL))		/* CKE minimum width, Clock Cycles */
#define TCKESR		(TCKE + 1)
#define TCKSRX		MAX(5, NS_TO_CYCLES_UP(10UL))
#define TCKSRE		MAX(5, NS_TO_CYCLES_UP(10UL))
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
#define TCKE		MAX(3, PS_TO_CYCLES_UP(7500UL))		/* CKE minimum width, Clock Cycles */
#define TCKESR		NS_TO_CYCLES_UP(15UL)
#define TCKSRX		2
#define TCKSRE		2
#define TCKCSX		(2 + TXP)
#define TCKDPDX		2
#define TCKDPDE		2
#define TCKPDX		2
#define TCKPDE		2
#endif
#define TDLLK		512UL					/* DLL Lock time, Clock Cycles */
#define TAOND		2UL					/* ODT Turn-On Delay , Clock Cycles */
#define TDQSCK_MIN	PS_TO_CYCLES_DOWN(tDQSCK_MIN)		/* DQS Output access time */
#define TDQSCK_MAX	PS_TO_CYCLES_DOWN(tDQSCK_MAX)		/* DQS Output access time */
#define BL		8UL					/* Burst Length */
#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
#define RL		(CL + AL)				/* Read Latency */
#define WL		(CWL + AL)				/* Write Latency */
#endif
#if defined(CONFIG_LPDDR3)
#define TWR		MAX(4, NS_TO_CYCLES_UP(tWR))		/* Write Recovery */
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
#define TWR		NS_TO_CYCLES_UP(tWR)			/* Write Recovery */
#endif
#define TINIT1		NS_TO_CYCLES_UP(100UL)			/* Clock Cycles */
#define TINIT2		5UL					/* Clock Cycles */
#define TINIT4		NS_TO_CYCLES_UP(1000UL)			/* Clock Cycles */
#define TINIT5		NS_TO_CYCLES_UP(10000UL)		/* Clock Cycles */
#define TPOSTCKE	NS_TO_CYCLES_UP(tPOSTCKE)		/* Clock Cycles */
#define TZQINIT		PS_TO_CYCLES_UP(1000000UL)		/* Clock Cycles */
#define TZQRESET	MAX(3, NS_TO_CYCLES_UP(50UL))		/* Clock Cycles */
};

#define OFFSETOF(TYPE, ELEMENT) ((unsigned int)&(((TYPE *)0)->ELEMENT))

#endif
