/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 *
 * This represents a driver for the Synopsys PHY Utility Block Lite (PUBL)
 * which is a external RAM memory PHY.
 *
 * It provides access to the register map and holds a configuration inside
 * a structure of parameters.
 *
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

#include "board.h"
#include "debug.h"
#include "hardware.h"
#include "timer.h"

#include "publ_regs.h"
#include "publ.h"

#include "dram_helpers.h"

static struct publ_regs		*PUBL;

void publ_init(void * config_data)
{
	struct dram_timings *timings = config_data;

#ifndef CONFIG_SYS_BASE_PUBL
#error "CONFIG_SYS_BASE_PUBL undefined"
#endif

	unsigned long tRFC = timings->tRFC;
	unsigned long tWR = timings->tWR;
	unsigned long tRP_ps = timings->tRP_ps;
	unsigned long tRCD_ps = timings->tRCD_ps;
	unsigned long tRAS = timings->tRAS;
	unsigned long tRC_ps = timings->tRC_ps;
	unsigned long tFAW = timings->tFAW;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long tDQSCK_MIN = timings->tDQSCK_MIN;
	unsigned long tDQSCK_MAX = timings->tDQSCK_MAX;
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long CL = timings->CL;
#endif
#ifdef CONFIG_DDR3
	unsigned long CWL = timings->CWL;
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long AL = timings->AL;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long RL = timings->RL;
	unsigned long WL = timings->WL;
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2) || defined(CONFIG_LPDDR3)
	unsigned long MRD = timings->MRD;
#endif

	PUBL = (struct publ_regs *) CONFIG_SYS_BASE_PUBL;
#ifdef CONFIG_DDR3
	PUBL->PUBL_DCR = PUBL_DCR_DDRMD_DDR3 | PUBL_DCR_DDRMD_DDR8BNK;
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_DCR = PUBL_DCR_DDRMD_DDR2 | PUBL_DCR_DDRMD_DDR8BNK;
#endif
#ifdef CONFIG_LPDDR2
	PUBL->PUBL_DCR = PUBL_DCR_DDRMD_LPDDR2 | PUBL_DCR_DDRMD_DDR8BNK |
			PUBL_DCR_DDRTYPE_S4;
#endif
#ifdef CONFIG_LPDDR3
	PUBL->PUBL_DCR = PUBL_DCR_DDRMD_LPDDR3 | PUBL_DCR_DDRMD_DDR8BNK;
#endif
	dbg_very_loud("PUBL_DCR %x\n", PUBL->PUBL_DCR);

	PUBL->PUBL_PGCR =
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
			PUBL_PGCR_DFTCMP
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			PUBL_PGCR_DQSCFG
#endif
			| PUBL_PGCR_CKEN(1) | PUBL_PGCR_CKDV(2) |
			PUBL_PGCR_RANKEN(1) | PUBL_PGCR_ZCKSEL(1) |
			PUBL_PGCR_PDDISDX;
	dbg_very_loud("PUBL_PGCR %x\n", PUBL->PUBL_PGCR);

	PUBL->PUBL_PTR0 = PUBL_PTR0_TDLLSRST(MAX(NS_TO_CYCLES_UP(50UL), 8)) |
			PUBL_PTR0_TDLLLOCK(NS_TO_CYCLES_UP(5120UL)) |
			PUBL_PTR0_TITMSRST(8);
	dbg_very_loud("PUBL_PTR0 %x\n", PUBL->PUBL_PTR0);

#ifdef CONFIG_DDR3
	PUBL->PUBL_PTR1 = PUBL_PTR1_TDINIT0(NS_TO_CYCLES_UP(520000UL)) |
			PUBL_PTR1_TDINIT1(MAX(NS_TO_CYCLES_UP(tRFC + 10), 5));
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_PTR1 = PUBL_PTR1_TDINIT0(NS_TO_CYCLES_UP(2000UL)) |
			PUBL_PTR1_TDINIT1(NS_TO_CYCLES_UP(400UL));
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_PTR1 = PUBL_PTR1_TDINIT0(NS_TO_CYCLES_UP(200000UL)) |
			PUBL_PTR1_TDINIT1(NS_TO_CYCLES_UP(100UL));
#endif
	dbg_very_loud("PUBL_PTR1 %x\n", PUBL->PUBL_PTR1);

#ifdef CONFIG_DDR3
	PUBL->PUBL_PTR2 = PUBL_PTR2_TDINIT2(NS_TO_CYCLES_UP(220000UL)) |
			PUBL_PTR2_TDINIT3(534);
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_PTR2 = PUBL_PTR2_TDINIT2(NS_TO_CYCLES_UP(400UL)) |
			PUBL_PTR2_TDINIT3(534);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_PTR2 = PUBL_PTR2_TDINIT2(MAX(5, NS_TO_CYCLES_UP(11000UL))) |
			PUBL_PTR2_TDINIT3(NS_TO_CYCLES_UP(200000UL));
#endif
	dbg_very_loud("PUBL_PTR2 %x\n", PUBL->PUBL_PTR2);

#ifdef CONFIG_DDR3
	PUBL->PUBL_MR0 = PUBL_MR0_CL((CL - 4) << 2 ) |
			PUBL_MR0_WR(NS_TO_CYCLES_UP(tWR) - 4) | 1;
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_MR0 = PUBL_MR0_CL(CL) | PUBL_MR0_BL_8 |
			PUBL_MR0_WR(NS_TO_CYCLES_UP(tWR) - 1) | 1;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_MR0 = 0;
#endif
	dbg_very_loud("PUBL_MR0 %x\n", PUBL->PUBL_MR0);

#ifdef CONFIG_DDR3
	PUBL->PUBL_MR1 = PUBL_MR1_RTT1 | PUBL_MR1_AL(AL);
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_MR1 = PUBL_MR1_RTT1 | PUBL_MR1_AL(AL) | PUBL_MR1_OCD(0);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_MR1 = PUBL_MR1_BL(0x3) | PUBL_MR1_NWR(TWR - 2);
#endif
	dbg_very_loud("PUBL_MR1 %x\n", PUBL->PUBL_MR1);

#ifdef CONFIG_DDR3
	PUBL->PUBL_MR2 = PUBL_MR2_CWL(CWL - 5);
#endif
#ifdef CONFIG_DDR2
	PUBL->PUBL_MR2 = 0;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	if (RL == 6 && WL == 3)
		PUBL->PUBL_MR2 = PUBL_MR2_RLWL_6_3;
	if (RL == 8 && WL == 4)
		PUBL->PUBL_MR2 = PUBL_MR2_RLWL_8_4;
#endif
	dbg_very_loud("PUBL_MR2 %x\n", PUBL->PUBL_MR2);

#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	PUBL->PUBL_MR3 = 0;
#endif
#ifdef CONFIG_LPDDR2
	PUBL->PUBL_MR3 = PUBL_MR3_DS_48OHM;
#endif
#ifdef CONFIG_LPDDR3
	PUBL->PUBL_MR3 = 0;
#endif
	dbg_very_loud("PUBL_MR3 %x\n", PUBL->PUBL_MR3);

#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	PUBL->PUBL_ODTCR = PUBL_ODTCR_WRODT0(1);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_ODTCR = 0;
#endif
	dbg_very_loud("PUBL_ODTCR %x\n", PUBL->PUBL_ODTCR);

	PUBL->PUBL_DTPR0 =
#ifdef CONFIG_DDR3
			PUBL_DTPR0_TMRD(TMRD - 4) |
#endif
#if defined(CONFIG_DDR2)
			PUBL_DTPR0_TMRD(TMRD) |
#endif
#ifdef CONFIG_LPDDR3
			PUBL_DTPR0_TMRD(TMRD - 8) |
#endif
			PUBL_DTPR0_TRTP(MAX(TRTP, 2)) |
			PUBL_DTPR0_TWTR(MAX(TWTR, 1)) |
#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
			PUBL_DTPR0_TRP(MAX(TRP, 2)) |
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			PUBL_DTPR0_TRP(MAX(TRP - 1, 2)) | /* PHY adds + 1 back */
			PUBL_DTPR0_TCCD |
#endif
			PUBL_DTPR0_TRCD(MAX(TRCD, 2)) |
			PUBL_DTPR0_TRAS(MAX(TRAS, 2)) |
			PUBL_DTPR0_TRRD(TRRD) |
			PUBL_DTPR0_TRC(MAX(TRC, 2));
	dbg_very_loud("PUBL_DTPR0 %x\n", PUBL->PUBL_DTPR0);

	PUBL->PUBL_DTPR1 = PUBL_DTPR1_TFAW(MAX(TFAW, 2)) |
#ifdef CONFIG_DDR3
			PUBL_DTPR1_TMOD(TMOD - 12) |
#endif
#ifdef CONFIG_DDR2
			PUBL_DTPR1_TAOND(TAOND - 2) |
#endif
			PUBL_DTPR1_TRFC(TRFC) |
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
			PUBL_DTPR1_TDQSCKMIN(1) | PUBL_DTPR1_TDQSCKMAX(1)
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			PUBL_DTPR1_TDQSCKMIN(TDQSCK_MIN)
#if defined(CONFIG_LPDDR2)
			| PUBL_DTPR1_TDQSCKMAX(TDQSCK_MAX)
#endif
#endif
			;
	dbg_very_loud("PUBL_DTPR1 %x\n", PUBL->PUBL_DTPR1);

	PUBL->PUBL_DTPR2 =
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
			PUBL_DTPR2_TXS(MAX(TXS, TXSDLL)) |
			PUBL_DTPR2_TXP(MAX(2, MAX(TXP, TXPDLL))) |
			PUBL_DTPR2_TCKE(MAX(TCKESR, 2)) |
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			PUBL_DTPR2_TXS(MAX(2, TXSR)) |
			PUBL_DTPR2_TXP(MAX(2, TXP)) |
			PUBL_DTPR2_TCKE(MAX(TCKE, 2)) |
#endif
			PUBL_DTPR2_TDLLK(MAX(TDLLK, 2));
	dbg_very_loud("PUBL_DTPR2 %x\n", PUBL->PUBL_DTPR2);

	PUBL->PUBL_DSGCR = PUBL_DSGCR_PUREN | PUBL_DSGCR_BDISEN |
			PUBL_DSGCR_ZUEN | PUBL_DSGCR_LPIOPD |
			PUBL_DSGCR_LPDLLPD
#if defined(CONFIG_LPDDR2)
			| PUBL_DSGCR_DQSGX(MAX(1, NS_TO_CYCLES_UP(tDQSCK_MAX - tDQSCK_MIN)))
			| PUBL_DSGCR_DQSGE(MAX(1, NS_TO_CYCLES_UP(tDQSCK_MAX - tDQSCK_MIN)))
#endif
#if defined(CONFIG_LPDDR3)
			| PUBL_DSGCR_DQSGX(MAX(1, NS_TO_CYCLES_UP(tDQSCK_MAX - tDQSCK_MIN)) + 1)
			| PUBL_DSGCR_DQSGE(MAX(1, NS_TO_CYCLES_UP(tDQSCK_MAX - tDQSCK_MIN)) + 1)
#endif
			| PUBL_DSGCR_NOBUB |
			PUBL_DSGCR_FXDLAT | PUBL_DSGCR_NL2OE |
			PUBL_DSGCR_TPDOE | PUBL_DSGCR_CKOE |
			PUBL_DSGCR_ODTOE | PUBL_DSGCR_RSTOE |
			PUBL_DSGCR_CKEOE;
	dbg_very_loud("PUBL_DSCGR %x\n", PUBL->PUBL_DSGCR);

#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	PUBL->PUBL_DXCCR = 0;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	PUBL->PUBL_DXCCR = PUBL_DXCCR_DQSRES_688OHM | PUBL_DXCCR_DQSNRES_688OHM
			| PUBL_DXCCR_DQSNRES_PU;
#endif
	dbg_very_loud("PUBL_DXCCR %x\n", PUBL->PUBL_DXCCR);

#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
	PUBL->PUBL_ZQ0CR1 = PUBL_ZQ0CR1_ZPROG_OID(11) |
			PUBL_ZQ0CR1_ZPROG_ODT(1);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	/* 48 OHM */
	PUBL->PUBL_ZQ0CR1 = PUBL_ZQ0CR1_ZPROG_OID(9) |
			PUBL_ZQ0CR1_ZPROG_ODT(1);
#endif
	dbg_very_loud("PUBL_ZQ0CR1 %x\n", PUBL->PUBL_ZQ0CR1);

	dbg_very_loud("PUBL_ACIOCR %x\n", PUBL->PUBL_ACIOCR);
	dbg_very_loud("PUBL_DLLGCR %x\n", PUBL->PUBL_DLLGCR);
}

int publ_idone()
{
	udelay(100);

	WAIT_WHILE_COND(!(PUBL->PUBL_PGSR & PUBL_PGSR_IDONE), 50000);

	return 0;
}

int publ_start()
{
	PUBL->PUBL_PIR = PUBL_PIR_INIT | PUBL_PIR_CTLDINIT;

	publ_idone();

	return publ_idone();
}

int publ_train()
{
	PUBL->PUBL_PGCR = (PUBL->PUBL_PGCR & ~PUBL_PGCR_RANKEN_MASK) |
			 PUBL_PGCR_RANKEN(1);

	PUBL->PUBL_DTAR = (0x7F0UL << 0) | (0x1FFFUL << 12) | (0x1UL << 28);

	PUBL->PUBL_PIR = PUBL_PIR_INIT | PUBL_PIR_QSTRN | PUBL_PIR_RVTRN;

	if (publ_idone())
		return -1;

	dbg_very_loud("PUBL_DCUSR0 %x\n", PUBL->PUBL_DCUSR0);
	dbg_very_loud("PUBL_DCUSR1 %x\n", PUBL->PUBL_DCUSR1);
	if (PUBL->PUBL_PGSR & (PUBL_PGSR_DTERR | PUBL_PGSR_DTIERR |
		PUBL_PGSR_RVERR | PUBL_PGSR_RVEIRR)) {
		dbg_info("PUBL: Error Training PHY : PGSR = %x\n", PUBL->PUBL_PGSR);
		return -1;
	} else {
		dbg_info("PUBL: Training complete.\n");
	}

	return 0;
}
