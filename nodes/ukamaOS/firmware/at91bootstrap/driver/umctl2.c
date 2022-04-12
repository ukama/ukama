/* ----------------------------------------------------------------------------
 *         Microchip Technology AT91Bootstrap project
 * ----------------------------------------------------------------------------
 *
 * This represents a driver for the Synopsys Universal Memory Controller 2
 * which is a external RAM memory controller.
 *
 * It provides access to the register map and holds a configuration inside
 * a structure of parameters.
 *
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

#include "board.h"
#include "debug.h"
#include "hardware.h"
#include "timer.h"

#ifdef CONFIG_RSTC
#include "rstc.h"
#endif

#include "umctl2_regs.h"
#include "umctl2.h"

#ifdef CONFIG_PUBL
#include "publ.h"
#endif

#include "dram_helpers.h"

#define MP_AXI_PORT_ENABLE(x) (1 << (x))

static struct umctl2_config_state *umctl2_config;

static struct uddrc_regs	*UDDRC_REGS;
static struct uddrc_mp		*UDDRC_MP;

unsigned long tRFC;
unsigned long tREFI;
unsigned long tWR;
unsigned long tRP;
unsigned long tRP_ps;
unsigned long tRCD;
unsigned long tRCD_ps;
unsigned long tCCD;
unsigned long tRAS;
unsigned long tRASMAX;
unsigned long tRC_ps;
unsigned long tFAW;
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
unsigned long tPRECKE;
#endif
unsigned long tPOSTCKE;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
unsigned long tDQSCK_MIN;
unsigned long tDQSCK_MAX;
#endif
unsigned long CL;
unsigned long CWL;
unsigned long AL;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
unsigned long RL;
unsigned long WL;
#endif
unsigned long TZQOPER;
unsigned long TZQCS;
unsigned long MRD;

struct dram_timings timings = {
	_tRFC, _tREFI, _tWR, _tRP, _tRP_ps, _tRCD, _tRCD_ps, _tCCD,_tRAS, _tRASMAX,
	_tRC_ps, _tFAW,
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	_tPRECKE,
#endif
	_tPOSTCKE,
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	_tDQSCK_MIN, _tDQSCK_MAX,
#endif
	_CL, _CWL, _AL,
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	_RL, _WL,
#endif
	_TZQOPER, _TZQCS, _MRD,
};

inline static void uddrc_mp_setup()
{
	/* how many clocks until changing to alternate transaction store (rd or wr depending on config) */
	UDDRC_REGS->UDDRC_SCHED	= UDDRC_SCHED_rdwr_idle_gap (umctl2_config->rdwr_idle_gap) |
	/* number of entries in low prio store minus one */
			UDDRC_SCHED_lpr_num_entries(umctl2_config->lpr_num_entries - 1) |
	/* midway policy between Open and Close page policies. meaning keep open or close after use */
			(umctl2_config->pageclose ? UDDRC_SCHED_pageclose : 0) |
	/* choose which transaction store is preffered */
			(umctl2_config->prefer_write ? UDDRC_SCHED_prefer_write : 0) |
	/* this is active low, so disabling force low priority */
			UDDRC_SCHED_force_low_pri_n;

	/* timer for pageclose policy of precharging */
	UDDRC_REGS->UDDRC_SCHED1 = UDDRC_SCHED1_pageclose_timer(umctl2_config->pageclose_timer);

	/* add starvation counters setup for the queues
	 * performance starvation settings on different queues (for 5ns clock)
	 */
	/* how many transactions to service once high prio read queue goes critical */
	UDDRC_REGS->UDDRC_PERFHPR1 = UDDRC_PERFHPR1_hpr_xact_run_length(umctl2_config->hpr_xact_run_length) |
	/* how many DFI clocks that high prio read queue can be starved until going critical */
			UDDRC_PERFHPR1_hpr_max_starve(umctl2_config->hpr_max_starve);
	/* how many transactions to service once low prio read queue goes critical */
	UDDRC_REGS->UDDRC_PERFLPR1 = UDDRC_PERFLPR1_lpr_xact_run_length(umctl2_config->hpr_xact_run_length) |
	/* how many DFI clocks that low prio read queue can be starved until going critical */
			UDDRC_PERFLPR1_lpr_max_starve(umctl2_config->lpr_max_starve);

	/* how many transactions to service once write queue goes critical (low + high prio) */
	UDDRC_REGS->UDDRC_PERFWR1 = UDDRC_PERFWR1_w_xact_run_length(umctl2_config->w_xact_run_length) |
	/* how many DFI clocks that write queue can be starved until going critical (low+high prio)*/
			UDDRC_PERFWR1_w_max_starve(umctl2_config->w_max_starve);

	/* auto precharge for collision cases*/
	UDDRC_REGS->UDDRC_DBG0 = (umctl2_config->dis_collision_page_opt ? UDDRC_DBG0_dis_collision_page_opt : 0) |
	/* disable write combined */
			(umctl2_config->dis_wc ? UDDRC_DBG0_dis_wc : 0 );

	/* reenable dequeue transactions from CAM and issues to DRAM*/
	UDDRC_REGS->UDDRC_DBG1 = 0;

	/* FPGA specific */
	UDDRC_REGS->UDDRC_DBGCMD = 0;

	/* enable interrupt for read//write/slverr AXI transaction poisoning
	   signal to AXI  */
	/* deprecated by TZC ? */

#if defined(AXI_TRANSACTION_POISONING_IRQ_EN)
	UDDRC_REGS->UDDRC_POISONCFG = UDDRC_POISONCFG_rd_poison_intr_en |
			UDDRC_POISONCFG_wr_poison_intr_en |
			UDDRC_POISONCFG_wr_poison_slverr_en;
#endif

	/* multi-port register settings (urgent bit are not connected in the DESIGN
	 *	because they cant be used in the matrix (no master capable))
 	 */

	/* burst length expansion mode = 0 -> each axi burst into multiple HIF commands.
			(this is recommended only for DDR4)
	 * pagematch limit =0 -> no limit on consecutive same page DDRC transactions
	 * go2critical_en = 0 -> signals to go to critical are not in design .
	 */

	UDDRC_MP->UDDRC_PCCFG = 0;

	/* if port is enabled */
	/* port 0 should always be enabled otherwise CPU can't access DDR */
	if (umctl2_config->axi_port_bitmap & MP_AXI_PORT_ENABLE(0)) {
	/* keep read/writes in order ? */
	UDDRC_MP->UDDRC_PCFGR_0 = (umctl2_config->port_x_rdwr_ordered_en ?
				UDDRC_PCFGR_0_rdwr_ordered_en : 0) |
	/* enable pagematch feature: keep port granted if same memory page */
			UDDRC_PCFGR_0_rd_port_pagematch_en |
	/* enable aging function for the read channel */
			UDDRC_PCFGR_0_rd_port_aging_en |
	/* reorder bypass ? */
			UDDRC_PCFGR_0_read_reorder_bypass_en |
	/* initial aging counter for port */
			UDDRC_PCFGR_0_rd_port_priority(0x3FF);

	/* enable pagematch feature: keep port granted if same memory page */
	UDDRC_MP->UDDRC_PCFGW_0	= UDDRC_PCFGW_0_wr_port_pagematch_en |
	/* enable aging for write on the port */
			UDDRC_PCFGW_0_wr_port_aging_en |
	/* initial load value of aging counter for write */
			UDDRC_PCFGW_0_wr_port_priority(0x3FF);

	/* arqos region mapping */
	UDDRC_MP->UDDRC_PCFGQOS0_0 =
			/* end of region 0 mapping */
			UDDRC_PCFGQOS0_0_rqos_map_level1(umctl2_config->port_x_read_region0_last[0]) |
			/* region 2 not used as this is single queue */
			UDDRC_PCFGQOS0_0_rqos_map_level2(0xe) |
			/* region 0 is LPR */
			UDDRC_PCFGQOS0_0_rqos_map_region0(0) |
			/* region 1 is VPR */
			UDDRC_PCFGQOS0_0_rqos_map_region1(1) |
			/* region 2 is unused */
			UDDRC_PCFGQOS0_0_rqos_map_region2(2);

	/* timeout for transactions on blue and red queue */
	UDDRC_MP->UDDRC_PCFGQOS1_0 = UDDRC_PCFGQOS1_0_rqos_map_timeoutb(0x0100) |
			UDDRC_PCFGQOS1_0_rqos_map_timeoutr(0x0100);
	/* awqos region mapping */
	UDDRC_MP->UDDRC_PCFGWQOS0_0 =
			UDDRC_PCFGWQOS0_0_wqos_map_level1(umctl2_config->port_x_write_region0_last[0]) |
			UDDRC_PCFGWQOS0_0_wqos_map_level2(umctl2_config->port_x_write_region1_last[0]) |
			UDDRC_PCFGWQOS0_0_wqos_map_region0(0) |
			UDDRC_PCFGWQOS0_0_wqos_map_region1(1);

	/* timeout for write transactions in the regions */
	UDDRC_MP->UDDRC_PCFGWQOS1_0 = UDDRC_PCFGWQOS1_0_wqos_map_timeout(0x0100);
	}

	/* if port is enabled */
	/* port 1 should always be enabled otherwise we can't access regmap */
	if (umctl2_config->axi_port_bitmap & MP_AXI_PORT_ENABLE(1)) {

	/* keep read/writes in order ? */
	UDDRC_MP->UDDRC_PCFGR_1 = (umctl2_config->port_x_rdwr_ordered_en ?
				UDDRC_PCFGR_1_rdwr_ordered_en : 0) |
	/* enable pagematch feature: keep port granted if same memory page */
			UDDRC_PCFGR_1_rd_port_pagematch_en |
	/* enable aging function for the read channel */
			UDDRC_PCFGR_1_rd_port_aging_en |
	/* reorder bypass ? */
			UDDRC_PCFGR_1_read_reorder_bypass_en |
	/* initial aging counter for port */
			UDDRC_PCFGR_1_rd_port_priority(0x3FF);

	/* enable pagematch feature: keep port granted if same memory page */
	UDDRC_MP->UDDRC_PCFGW_1	= UDDRC_PCFGW_1_wr_port_pagematch_en |
	/* enable aging for write on the port */
			UDDRC_PCFGW_1_wr_port_aging_en |
	/* initial load value of aging counter for write */
			UDDRC_PCFGW_1_wr_port_priority(0x3FF);

	/* arqos region mapping */
	UDDRC_MP->UDDRC_PCFGQOS0_1 =
			/* end of region 0 mapping */
			UDDRC_PCFGQOS0_1_rqos_map_level1(umctl2_config->port_x_read_region0_last[1]) |
			/* region 2 not used as this is single queue */
			UDDRC_PCFGQOS0_1_rqos_map_level2(umctl2_config->port_x_read_region1_last[1]) |
			/* region 0 is VPR */
			UDDRC_PCFGQOS0_1_rqos_map_region0(1) |
			/* region 1 is VPR */
			UDDRC_PCFGQOS0_1_rqos_map_region1(1) |
			/* region 2 is unused */
			UDDRC_PCFGQOS0_1_rqos_map_region2(2);

	/* timeout for transactions on blue and red queue */
	UDDRC_MP->UDDRC_PCFGQOS1_1 = UDDRC_PCFGQOS1_1_rqos_map_timeoutb(0x040) |
			UDDRC_PCFGQOS1_1_rqos_map_timeoutr(0x040);
	/* awqos region mapping */
	UDDRC_MP->UDDRC_PCFGWQOS0_1 =
			UDDRC_PCFGWQOS0_1_wqos_map_level1(umctl2_config->port_x_write_region0_last[1]) |
			UDDRC_PCFGWQOS0_1_wqos_map_level2(umctl2_config->port_x_write_region1_last[1]) |
			UDDRC_PCFGWQOS0_1_wqos_map_region0(1) |
			UDDRC_PCFGWQOS0_1_wqos_map_region1(1);

	/* timeout for write transactions in the regions */
	UDDRC_MP->UDDRC_PCFGWQOS1_1 = UDDRC_PCFGWQOS1_1_wqos_map_timeout(0x040);
	}

	/* if port is enabled */
	/* port 2 is dual queue, for PSS */
	if (umctl2_config->axi_port_bitmap & MP_AXI_PORT_ENABLE(2)) {

	/* keep read/writes in order ? */
	UDDRC_MP->UDDRC_PCFGR_2 = (umctl2_config->port_x_rdwr_ordered_en ?
				UDDRC_PCFGR_2_rdwr_ordered_en : 0) |
	/* enable pagematch feature: keep port granted if same memory page */
			UDDRC_PCFGR_2_rd_port_pagematch_en |
	/* enable aging function for the read channel */
			UDDRC_PCFGR_2_rd_port_aging_en |
	/* reorder bypass ? */
			UDDRC_PCFGR_2_read_reorder_bypass_en |
	/* initial aging counter for port */
			UDDRC_PCFGR_2_rd_port_priority(0x3FF);

	/* enable pagematch feature: keep port granted if same memory page */
	UDDRC_MP->UDDRC_PCFGW_2	= UDDRC_PCFGW_2_wr_port_pagematch_en |
	/* enable aging for write on the port */
			UDDRC_PCFGW_2_wr_port_aging_en |
	/* initial load value of aging counter for write */
			UDDRC_PCFGW_2_wr_port_priority(0x3FF);

	/* arqos region mapping */
	UDDRC_MP->UDDRC_PCFGQOS0_2 =
			/* end of region 0 mapping */
			UDDRC_PCFGQOS0_2_rqos_map_level1(umctl2_config->port_x_read_region0_last[2]) |
			/* end of region 1 mapping */
			UDDRC_PCFGQOS0_2_rqos_map_level2(umctl2_config->port_x_read_region1_last[2]) |
			/* region 0 is LPR */
			UDDRC_PCFGQOS0_2_rqos_map_region0(0) |
			/* region 1 is vPR */
			UDDRC_PCFGQOS0_2_rqos_map_region1(1) |
			/* region 2 is HPR */
			UDDRC_PCFGQOS0_2_rqos_map_region2(2);

	/* timeout for transactions on blue and red queue */
	UDDRC_MP->UDDRC_PCFGQOS1_2 = UDDRC_PCFGQOS1_2_rqos_map_timeoutb(0x0200) |
			UDDRC_PCFGQOS1_2_rqos_map_timeoutr(0x0200);
	/* awqos region mapping */
	UDDRC_MP->UDDRC_PCFGWQOS0_2 =
			UDDRC_PCFGWQOS0_2_wqos_map_level1(umctl2_config->port_x_write_region0_last[2]) |
			UDDRC_PCFGWQOS0_2_wqos_map_level2(umctl2_config->port_x_write_region1_last[2]) |
			UDDRC_PCFGWQOS0_2_wqos_map_region0(0) |
			UDDRC_PCFGWQOS0_2_wqos_map_region1(1);

	/* timeout for write transactions in the regions */
	UDDRC_MP->UDDRC_PCFGWQOS1_2 = UDDRC_PCFGWQOS1_2_wqos_map_timeout1(0x0200) |
				UDDRC_PCFGWQOS1_2_wqos_map_timeout2(0x0200);
	}

	if (umctl2_config->axi_port_bitmap & MP_AXI_PORT_ENABLE(3)) {

	/* keep read/writes in order ? */
	UDDRC_MP->UDDRC_PCFGR_3 = (umctl2_config->port_x_rdwr_ordered_en ?
				UDDRC_PCFGR_3_rdwr_ordered_en : 0) |
	/* enable pagematch feature: keep port granted if same memory page */
			UDDRC_PCFGR_3_rd_port_pagematch_en |
	/* enable aging function for the read channel */
			UDDRC_PCFGR_3_rd_port_aging_en |
	/* reorder bypass ? */
			UDDRC_PCFGR_3_read_reorder_bypass_en |
	/* initial aging counter for port */
			UDDRC_PCFGR_3_rd_port_priority(0x3FF);

	/* enable pagematch feature: keep port granted if same memory page */
	UDDRC_MP->UDDRC_PCFGW_3	= UDDRC_PCFGW_3_wr_port_pagematch_en |
	/* enable aging for write on the port */
			UDDRC_PCFGW_3_wr_port_aging_en |
	/* initial load value of aging counter for write */
			UDDRC_PCFGW_3_wr_port_priority(0x3FF);

	/* arqos region mapping */
	UDDRC_MP->UDDRC_PCFGQOS0_3 =
			/* end of region 0 mapping */
			UDDRC_PCFGQOS0_3_rqos_map_level1(umctl2_config->port_x_read_region0_last[3]) |
			UDDRC_PCFGQOS0_3_rqos_map_level2(umctl2_config->port_x_read_region1_last[3]) |
			/* region 0 is LPR */
			UDDRC_PCFGQOS0_3_rqos_map_region0(0) |
			/* region 1 is VPR */
			UDDRC_PCFGQOS0_3_rqos_map_region1(1) |
			UDDRC_PCFGQOS0_3_rqos_map_region2(2);

	/* timeout for transactions on blue and red queue */
	UDDRC_MP->UDDRC_PCFGQOS1_3 = UDDRC_PCFGQOS1_3_rqos_map_timeoutb(0x0100) |
			UDDRC_PCFGQOS1_3_rqos_map_timeoutr(0x0100);
	/* awqos region mapping */
	UDDRC_MP->UDDRC_PCFGWQOS0_3 =
			UDDRC_PCFGWQOS0_3_wqos_map_level1(umctl2_config->port_x_write_region0_last[3]) |
			UDDRC_PCFGWQOS0_3_wqos_map_level2(umctl2_config->port_x_write_region1_last[3]) |
			UDDRC_PCFGWQOS0_3_wqos_map_region0(0) |
			UDDRC_PCFGWQOS0_3_wqos_map_region1(1);

	/* timeout for write transactions in the regions */
	UDDRC_MP->UDDRC_PCFGWQOS1_3 = UDDRC_PCFGWQOS1_3_wqos_map_timeout(0x0100);
	}


	if (umctl2_config->axi_port_bitmap & MP_AXI_PORT_ENABLE(4)) {

	/* keep read/writes in order ? */
	UDDRC_MP->UDDRC_PCFGR_4 = (umctl2_config->port_x_rdwr_ordered_en ?
				UDDRC_PCFGR_4_rdwr_ordered_en : 0) |
	/* enable pagematch feature: keep port granted if same memory page */
			UDDRC_PCFGR_4_rd_port_pagematch_en |
	/* enable aging function for the read channel */
			UDDRC_PCFGR_4_rd_port_aging_en |
	/* reorder bypass ? */
			UDDRC_PCFGR_4_read_reorder_bypass_en |
	/* initial aging counter for port */
			UDDRC_PCFGR_4_rd_port_priority(0x3FF);

	/* enable pagematch feature: keep port granted if same memory page */
	UDDRC_MP->UDDRC_PCFGW_4	= UDDRC_PCFGW_4_wr_port_pagematch_en |
	/* enable aging for write on the port */
			UDDRC_PCFGW_4_wr_port_aging_en |
	/* initial load value of aging counter for write */
			UDDRC_PCFGW_4_wr_port_priority(0x3FF);

	/* arqos region mapping */
	UDDRC_MP->UDDRC_PCFGQOS0_4 =
			/* end of region 0 mapping */
			UDDRC_PCFGQOS0_4_rqos_map_level1(umctl2_config->port_x_read_region0_last[4]) |
			UDDRC_PCFGQOS0_4_rqos_map_level2(umctl2_config->port_x_read_region1_last[4]) |
			/* region 0 is LPR */
			UDDRC_PCFGQOS0_4_rqos_map_region0(0) |
			/* region 1 is VPR */
			UDDRC_PCFGQOS0_4_rqos_map_region1(1) |
			UDDRC_PCFGQOS0_4_rqos_map_region2(2);

	/* timeout for transactions on blue and red queue */
	UDDRC_MP->UDDRC_PCFGQOS1_4 = UDDRC_PCFGQOS1_4_rqos_map_timeoutb(0x0100) |
			UDDRC_PCFGQOS1_4_rqos_map_timeoutr(0x0100);
	/* awqos region mapping */
	UDDRC_MP->UDDRC_PCFGWQOS0_4 =
			UDDRC_PCFGWQOS0_4_wqos_map_level1(umctl2_config->port_x_write_region0_last[4]) |
			UDDRC_PCFGWQOS0_4_wqos_map_level2(umctl2_config->port_x_write_region1_last[4]) |
			UDDRC_PCFGWQOS0_4_wqos_map_region0(0) |
			UDDRC_PCFGWQOS0_4_wqos_map_region1(1);

	/* timeout for write transactions in the regions */
	UDDRC_MP->UDDRC_PCFGWQOS1_4 = UDDRC_PCFGWQOS1_4_wqos_map_timeout(0x0100);
	}

	/* SARBASE/SIZE must be chunks of 512MB */
	UDDRC_MP->UDDRC_SARBASE0 = ((unsigned int) MEM_BANK) /
				((unsigned int) (512 * 1024 * 1024));
	UDDRC_MP->UDDRC_SARSIZE0 = (2 * 1024) / (unsigned int) MEM_SIZE - 1;
	dbg_very_loud("sarbase %x, offset = %x\n", UDDRC_MP->UDDRC_SARBASE0,
			OFFSETOF(struct uddrc_mp, UDDRC_SARBASE0));
}

inline static void uddrc_spin_cam_pipeline_empty()
{
	/* spin until both rd and wr pipelines are empty */
	while (!(UDDRC_REGS->UDDRC_DBGCAM & UDDRC_DBGCAM_wr_data_pipeline_empty));
	while (!(UDDRC_REGS->UDDRC_DBGCAM & UDDRC_DBGCAM_rd_data_pipeline_empty));
	/* must read them twice according to datasheet */
	while (!(UDDRC_REGS->UDDRC_DBGCAM & UDDRC_DBGCAM_wr_data_pipeline_empty));
	while (!(UDDRC_REGS->UDDRC_DBGCAM & UDDRC_DBGCAM_rd_data_pipeline_empty));
}

inline static void uddrc_dis_dq()
{
	/* disable dequeuing of rd/wr commands from CAM to DRAM */
	UDDRC_REGS->UDDRC_DBG1 |= UDDRC_DBG1_dis_dq;
#if defined SPIN_CAM_PIPELINE_EMPTY
	uddrc_spin_cam_pipeline_empty();
#endif
}

inline static void uddrc_selfref_en()
{
	/* put DRAM into self refresh after a number of clocks PWRTMG.selfref_to_x32 */
	UDDRC_REGS->UDDRC_PWRCTL |= UDDRC_PWRCTL_selfref_en;
}

inline static void uddrc_mstr()
{
	/* initialize MSTR - master register */
	/* DRAM burst length of 8 for DDR3/DDR2 */
	UDDRC_REGS->UDDRC_MSTR = UDDRC_MSTR_burst_rdwr(UDDRC_MSTR_burst_rdwr_BL8)
	/* use burstchop only for DDR3*/
#ifdef CONFIG_DDR3
							| UDDRC_MSTR_burstchop
#endif
#ifdef CONFIG_DDR3
							| UDDRC_MSTR_ddr3
#endif
#ifdef CONFIG_LPDDR2
							| UDDRC_MSTR_lpddr2
#endif
#ifdef CONFIG_LPDDR3
							| UDDRC_MSTR_lpddr3
#endif
	;

	dbg_very_loud("UMCTL2 MSTR %x\n", UDDRC_REGS->UDDRC_MSTR);
}

inline static void uddrc_init0()
{
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long pre_cke_x1024 = 1 +
		DIV_ROUND_UP(NS_TO_CYCLES_UP(tPRECKE), 1024 * 2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long pre_cke_x1024 = 1 + DIV_ROUND_UP(TINIT1, 1024 * 2);
#endif
#ifdef CONFIG_DDR3
	unsigned long post_cke_x1024 = 1 +
		DIV_ROUND_UP(NS_TO_CYCLES_UP(TXPR), 1024 * 2);
#endif
#ifdef CONFIG_DDR2
	unsigned long post_cke_x1024 = 2;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long post_cke_x1024 = 1 + DIV_ROUND_UP(TPOSTCKE, 1024 * 2);
#endif

	/* DRAM init register 0 */
	UDDRC_REGS->UDDRC_INIT0	=
	/* wait val * 1024 DFI cycles after CKE is high to start the DRAM init */
	/* used in step 4) from DDR3 state machine init seq */
				UDDRC_INIT0_post_cke_x1024(post_cke_x1024) |
	/* wait val * 1024 DFI cycles before CKE is high , to start the DRAM init */
	/* used in step 3) from DDR3 state machine init seq */
				UDDRC_INIT0_pre_cke_x1024(pre_cke_x1024)
	;
	dbg_very_loud("UMCTL2 INIT 0 %x\n", UDDRC_REGS->UDDRC_INIT0);
}

inline static void uddrc_init1()
{
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long dram_rstn_x1024 = 1 +
			DIV_ROUND_UP(NS_TO_CYCLES_UP(200000), 1024 * 2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long dram_rstn_x1024 = 0;
#endif
	/* DRAM init register 1 */
	/* number of cycles to assert DRAM reset during init sequence - dfi_reset_n */
	/* used in step 2) from DDR3 state machine init seq */
	UDDRC_REGS->UDDRC_INIT1	= UDDRC_INIT1_dram_rstn_x1024(dram_rstn_x1024) |
	/* number of 32 DFI clocks to wait before driving OCD (driver calibration?) complete to DRAM */
				UDDRC_INIT1_pre_ocd_x32(0)
	;
	dbg_very_loud("UMCTL2 INIT 1 %x\n", UDDRC_REGS->UDDRC_INIT1);

}

inline static void uddrc_init2()
{
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long min_stable_clock_x1 = DIV_ROUND_UP(TINIT2, 2);
#endif
#if defined(CONFIG_LPDDR2)
	unsigned long idle_after_reset_x32 = DIV_ROUND_UP(TINIT4, 2 * 32);
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	UDDRC_REGS->UDDRC_INIT2 = 0;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	UDDRC_REGS->UDDRC_INIT2 = UDDRC_INIT2_min_stable_clock_x1(min_stable_clock_x1)
#if defined(CONFIG_LPDDR2)
				| UDDRC_INIT2_idle_after_reset_x32(idle_after_reset_x32 + 1)
#endif
				;
#endif
	dbg_very_loud("UMCTL2 INIT 2 %x\n", UDDRC_REGS->UDDRC_INIT2);
}

inline static void uddrc_init3()
{
	/* DRAM init register 3 */
	/* INIT3 : MR1 + MR0 (registers from module)
	 MR0[11:9].wr		=3'b010 (6)
	 MR0[8].dr    	=1'b1   (DLL-Reset)
	 MR0[7].tm    	=1'b0
	 MR0[6:4].cl.2:1	=3'b010 (6)
	 MR0[3].bt    	=1'b0   (sequential burst)
	 MR0[2].cl.0  	=1'b0   (7)
	 MR0[1:0].bl   	=2'b00  (8 fixed)
	 */
#ifdef CONFIG_DDR3
	UDDRC_REGS->UDDRC_INIT3	= UDDRC_INIT3_mr(
				(1 << 0) | ((CL - 4) << 4) | (1 << 8) |
				((TWR - 4) << 9)) |
				UDDRC_INIT3_emr(1 << 6);
#endif
#ifdef CONFIG_DDR2
	UDDRC_REGS->UDDRC_INIT3 = UDDRC_INIT3_mr(
				((BL == 8) ? (3 << 0) : (2 << 0))
				 | (CL << 4) | (1 << 8) | ((TWR - 1) << 9));
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	UDDRC_REGS->UDDRC_INIT3 = UDDRC_INIT3_mr(3 | ((TWR - 2) << 5))
				| UDDRC_INIT3_emr(RL - 2);
#endif

	dbg_very_loud("UMCTL2 INIT 3 %x\n", UDDRC_REGS->UDDRC_INIT3);
}

inline static void uddrc_init4()
{
	/* DRAM init register 4 */
	UDDRC_REGS->UDDRC_INIT4	=
#ifdef CONFIG_DDR3
				 UDDRC_INIT4_emr2(((CWL - 5) << 3))
#endif
#ifdef CONFIG_DDR2
				UDDRC_INIT4_emr2(0)
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
				UDDRC_INIT4_emr2(3)
#endif
				 | UDDRC_INIT4_emr3(0);
	dbg_very_loud("UMCTL2 INIT 4 %x\n", UDDRC_REGS->UDDRC_INIT4);
}

inline static void uddrc_init5()
{
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long dev_zqinit_x32 = 1 + DIV_ROUND_UP(512, 32*2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long dev_zqinit_x32 = 1 + DIV_ROUND_UP(TZQINIT, 32*2);
	unsigned long max_auto_init_x1024 = 1 + DIV_ROUND_UP(TINIT5, 1024);
#endif
	/* DRAM init register 5 */
	/* initial calibration time */
	UDDRC_REGS->UDDRC_INIT5	= UDDRC_INIT5_dev_zqinit_x32(dev_zqinit_x32)
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
				| UDDRC_INIT5_max_auto_init_x1024(max_auto_init_x1024)
#endif
				;
	dbg_very_loud("UMCTL2 INIT 5 %x\n", UDDRC_REGS->UDDRC_INIT5);
}

inline static void uddrc_dramtmg()
{
	unsigned long wr2pre = DIV_ROUND_UP(WL + (BL / 2) + TWR
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
		+ 1
#endif
		, 2);
	unsigned long tfaw = DIV_ROUND_UP(TFAW, 2);
	unsigned long t_ras_min = DIV_ROUND_UP(TRAS, 2);
	unsigned long t_ras_max = DIV_ROUND_DOWN(
			DIV_ROUND_DOWN(NS_TO_CYCLES_DOWN(tRASMAX), 1024) - 1, 2);
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long t_xp = DIV_ROUND_UP(TXPDLL, 2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long t_xp = DIV_ROUND_UP(TXP, 2);
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long rd2pre = AL + DIV_ROUND_UP(TRTP, 2);
#endif
#if defined(CONFIG_LPDDR2)
	unsigned long rd2pre = DIV_ROUND_UP((BL / 2) + MAX(TRTP, 2) - 2, 2);
#endif
#if defined(CONFIG_LPDDR3)
	unsigned long rd2pre = DIV_ROUND_UP((BL / 2) + MAX(TRTP, 4) - 4, 2);
#endif
	unsigned long t_rc = DIV_ROUND_UP(TRC, 2);
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long rd2wr = DIV_ROUND_UP(RL + (BL / 2) + 2 - WL, 2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long rd2wr = DIV_ROUND_UP(RL + (BL / 2) + TDQSCK_MAX + 1 - WL, 2);
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long wr2rd = DIV_ROUND_UP(CWL + (BL / 2) + TWTR, 2);
#if defined(CONFIG_DDR3)
	unsigned long t_mod = DIV_ROUND_UP(TMOD, 2);
#endif
#endif
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2) || defined(CONFIG_LPDDR3)
	unsigned long t_mrd = DIV_ROUND_UP(TMRD, 2);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long wr2rd = DIV_ROUND_UP(WL + (BL / 2) + TWTR, 2);
#endif
	unsigned long t_rcd = DIV_ROUND_UP(TRCD - AL, 2);
	unsigned long t_ccd = DIV_ROUND_UP(TCCD, 2);
	unsigned long t_rrd = DIV_ROUND_UP(TRRD, 2);
	unsigned long t_rp = DIV_ROUND_DOWN(TRP, 2) + 1;
	unsigned long t_cksrx = DIV_ROUND_UP(TCKSRX, 2);
	unsigned long t_cksre = DIV_ROUND_UP(TCKSRE, 2);
	unsigned long t_ckesr = DIV_ROUND_UP(TCKESR, 2);
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long t_cke = DIV_ROUND_UP(TCKE, 2);
	unsigned long t_xs_dll_x32 = DIV_ROUND_UP(TXSDLL, 32*2) + 1;
	unsigned long t_xs_x32 = DIV_ROUND_UP(TXS, 32*2) + 1;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long t_cke = DIV_ROUND_UP(MAX(TCKE, TCKESR), 2);

	unsigned long t_ckcsx = DIV_ROUND_UP(TCKCSX, 2);
	unsigned long t_ckdpdx = DIV_ROUND_UP(TCKDPDX, 2);
	unsigned long t_ckdpde = DIV_ROUND_UP(TCKDPDE, 2);
	unsigned long t_ckpdx = DIV_ROUND_UP(TCKPDX, 2);
	unsigned long t_ckpde = DIV_ROUND_UP(TCKPDE, 2);
#endif

	/* configure DRAMTMG* registers */
	/* write to precharge on same bank :
	 * (WriteLatency + Burstlength/2 + tWR write recovery)/2 - freq ratio 1:2
	 */
	UDDRC_REGS->UDDRC_DRAMTMG0 = UDDRC_DRAMTMG0_wr2pre(wr2pre) |
	/* tFAW/2 (en clock) : valid if 8 or more banks */
		UDDRC_DRAMTMG0_t_faw(tfaw) |
	/* tRAS max, max time activate to precharge, in amounts of 1024 clocks */
		UDDRC_DRAMTMG0_t_ras_max(t_ras_max) |
	/* tRAS min */
		UDDRC_DRAMTMG0_t_ras_min(t_ras_min);

	/* time from powerdown exit until any operation */
	UDDRC_REGS->UDDRC_DRAMTMG1 = UDDRC_DRAMTMG1_t_xp(t_xp) |
	/* read to precharge of same bank, (tAL + max(tRTP,4))/2 */
			UDDRC_DRAMTMG1_rd2pre(rd2pre) |
	/* tRC - row cycle ? time between activates to same bank , = tRCD + tWR + tRP = tRAS + tRP */
			UDDRC_DRAMTMG1_t_rc(t_rc);

	/* tWL - write latency, not required in DDR3. Only for LPDDR2*/
	UDDRC_REGS->UDDRC_DRAMTMG2 =
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			UDDRC_DRAMTMG2_write_latency(DIV_ROUND_UP(WL, 2)) |
	/* tRL - read latency, not required in DDR3. Only for LPDDR2 */
			UDDRC_DRAMTMG2_read_latency(DIV_ROUND_UP(RL, 2)) |
#endif
	/* read to write latency = (tRL + BL/2 +2 - tWL )/2 */
			UDDRC_DRAMTMG2_rd2wr(rd2wr) |
	/* write to read latency = (Cas write latency CWL + BL/2 + tWTR (internal write to read ))/2 */
			UDDRC_DRAMTMG2_wr2rd(wr2rd);

	/* tMRD/2 (cycles to wait after a mode reg write or read) */
	UDDRC_REGS->UDDRC_DRAMTMG3 =
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2) || defined(CONFIG_LPDDR3)
			UDDRC_DRAMTMG3_t_mrd(t_mrd)
#endif
#ifdef CONFIG_LPDDR2
			UDDRC_DRAMTMG3_t_mrw(DIV_ROUND_UP(TMRW, 2))
#endif
#ifdef CONFIG_DDR3
	/* tMOD/2 : cycles between load mode command and following non-load mode command */
			| UDDRC_DRAMTMG3_t_mod(t_mod)
#endif
			;

	/* (tRCD - tAL)/2 : row to column delay - activate to read or write to same bank */
	UDDRC_REGS->UDDRC_DRAMTMG4 = UDDRC_DRAMTMG4_t_rcd(t_rcd) |
	/* tCCD/2 :  minimum time between two reads or two writes */
			UDDRC_DRAMTMG4_t_ccd(t_ccd) |
	/* tRRD/2 : minimum time between activates from bank a to bank b */
			UDDRC_DRAMTMG4_t_rrd(t_rrd) |
	/* tRP : minimum time from precharge to activate of the same bank , precharge delay */
			UDDRC_DRAMTMG4_t_rp(t_rp);

	/* TCKSRX/2 (time before self refresh exit that CK is maintained as a valid clock before issuing SRX) */
	UDDRC_REGS->UDDRC_DRAMTMG5 = UDDRC_DRAMTMG5_t_cksrx(t_cksrx) |
	/* TCKSRE/2 (time after self refresh down entry that CK is maintained as a valid clock) */
			UDDRC_DRAMTMG5_t_cksre(t_cksre) |
	/* (tCKE+1)/2 (minimum CKE Low width for self refresh to exit) */
			UDDRC_DRAMTMG5_t_ckesr(t_ckesr) |
	/* RoundUp(tCKE/2) */
			UDDRC_DRAMTMG5_t_cke(t_cke);

#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	/* LPDDR only */
	UDDRC_REGS->UDDRC_DRAMTMG6 = UDDRC_DRAMTMG6_t_ckcsx(t_ckcsx) |
				 UDDRC_DRAMTMG6_t_ckdpdx(t_ckdpdx) |
				 UDDRC_DRAMTMG6_t_ckdpde(t_ckdpde);

	/* LPDDR only */
	UDDRC_REGS->UDDRC_DRAMTMG7 = UDDRC_DRAMTMG7_t_ckpdx(t_ckpdx) |
				 UDDRC_DRAMTMG7_t_ckpde(t_ckpde);
#endif

#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	/* tXSDLL/2 : exit self refresh to commands requiring a locked DLL */
	UDDRC_REGS->UDDRC_DRAMTMG8 = UDDRC_DRAMTMG8_t_xs_dll_x32(t_xs_dll_x32) |
	/* tXS/2 : exit self refresh to commands not requiring a locked DLL */
			UDDRC_DRAMTMG8_t_xs_x32(t_xs_x32);
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	UDDRC_REGS->UDDRC_DRAMTMG8 = 0;
#endif

	UDDRC_REGS->UDDRC_DRAMTMG14 = UDDRC_DRAMTMG14_t_xsr(0xA0);

	dbg_very_loud("DRAMTMG0 %x\n", UDDRC_REGS->UDDRC_DRAMTMG0);
	dbg_very_loud("DRAMTMG1 %x\n", UDDRC_REGS->UDDRC_DRAMTMG1);
	dbg_very_loud("DRAMTMG2 %x\n", UDDRC_REGS->UDDRC_DRAMTMG2);
	dbg_very_loud("DRAMTMG3 %x\n", UDDRC_REGS->UDDRC_DRAMTMG3);
	dbg_very_loud("DRAMTMG4 %x\n", UDDRC_REGS->UDDRC_DRAMTMG4);
	dbg_very_loud("DRAMTMG5 %x\n", UDDRC_REGS->UDDRC_DRAMTMG5);
	dbg_very_loud("DRAMTMG6 %x\n", UDDRC_REGS->UDDRC_DRAMTMG6);
	dbg_very_loud("DRAMTMG7 %x\n", UDDRC_REGS->UDDRC_DRAMTMG7);
	dbg_very_loud("DRAMTMG8 %x\n", UDDRC_REGS->UDDRC_DRAMTMG8);
	dbg_very_loud("DRAMTMG14 %x\n", UDDRC_REGS->UDDRC_DRAMTMG14);
}

inline static void uddrc_addrmap_init()
{
	/* HIF address to bank address:
	This is explained mapping, need a width of 200 columns for a proper display.


							   bank address:			bit2 bit1 bit0
							   col address:						bit11bit10bit9 bit8 bit7 bit6 bit5 bit4 bit3 bit2
HIF bit:											4+8  3+8  2+8	-     -   9+0  8+0  7+0  6+0  5+0  4+0  3+0  2+0
=												12   11   10	          9    8    7    6    5    4    3    2

row address: bit15bit14bit13bit12bit11bit10bit9 bit8 bit7 bit6 bit5 bit4 bit3 bit2 bit1 bit0
HIF bit:     -    7+20 7+19 7+18 7+17 7+16 7+15 7+14 7+13 7+12 7+11 7+10 7+9  7+8  7+7  7+6
=	     -    27   26   25   24   23   22   21   20   19   18   17   16   15   14   13

	 */
	/* internal base = 4 3 2 for these bits */
	UDDRC_REGS->UDDRC_ADDRMAP1 = UDDRC_ADDRMAP1_addrmap_bank_b2(8) |
				UDDRC_ADDRMAP1_addrmap_bank_b1(8) |
				UDDRC_ADDRMAP1_addrmap_bank_b0(8);

	/* internal base = 5 4 3 2 for these bits */
	UDDRC_REGS->UDDRC_ADDRMAP2 = UDDRC_ADDRMAP2_addrmap_col_b5(0) |
				UDDRC_ADDRMAP2_addrmap_col_b4(0) |
				UDDRC_ADDRMAP2_addrmap_col_b3(0) |
				UDDRC_ADDRMAP2_addrmap_col_b2(0);

	/* internal base = 9 8 7 6 for these bits */
	UDDRC_REGS->UDDRC_ADDRMAP3 = UDDRC_ADDRMAP3_addrmap_col_b9(0) |
				UDDRC_ADDRMAP3_addrmap_col_b8(0) |
				UDDRC_ADDRMAP3_addrmap_col_b7(0) |
				UDDRC_ADDRMAP3_addrmap_col_b6(0);

	/* 31 means unused */
	UDDRC_REGS->UDDRC_ADDRMAP4 = UDDRC_ADDRMAP4_addrmap_col_b11(31)	|
				UDDRC_ADDRMAP4_addrmap_col_b10(31);

	/* internal base = 17, 7, 6 for these bits, b2_10=15 means we use reg MAP9 */
	UDDRC_REGS->UDDRC_ADDRMAP5 = UDDRC_ADDRMAP5_addrmap_row_b11(7) |
				UDDRC_ADDRMAP5_addrmap_row_b2_10(7) |
				UDDRC_ADDRMAP5_addrmap_row_b1(7) |
				UDDRC_ADDRMAP5_addrmap_row_b0(7);

	/* 15 = unused, internal base for these bits 19 and 18 */
	UDDRC_REGS->UDDRC_ADDRMAP6 = UDDRC_ADDRMAP6_addrmap_row_b15(15)	|
#if defined(CONFIG_RAM_512MB)
				UDDRC_ADDRMAP6_addrmap_row_b14(7) |
#endif
#if defined(CONFIG_RAM_256MB)
				UDDRC_ADDRMAP6_addrmap_row_b14(15) |
#endif
				UDDRC_ADDRMAP6_addrmap_row_b13(7) |
				UDDRC_ADDRMAP6_addrmap_row_b12(7);

	/* internal base for these bits 11 10 9 8 */
	UDDRC_REGS->UDDRC_ADDRMAP9 = 0;

	/* internal base for these bits 15 14 13 12 */
	UDDRC_REGS->UDDRC_ADDRMAP10 = 0;

	/* internal base for this bit 16 */
	UDDRC_REGS->UDDRC_ADDRMAP11 = 0;//UDDRC_ADDRMAP11_addrmap_row_b10(7);
}

inline static void uddrc_configure_refresh()
{
	/* after this many clocks of DDRC idle go to selfrefresh, multiples of 32 DFI clocks */
	UDDRC_REGS->UDDRC_PWRTMG = UDDRC_PWRTMG_selfref_to_x32(2) |
	/* minimum deep power down time */
			UDDRC_PWRTMG_t_dpd_x4096(0xaf) |
	/* after this many clocks of DDRC idle go to powerdown, multiples of 32 DFI clocks */
			UDDRC_PWRTMG_powerdown_to_x32(3);

	/* number of 32 DFI clocks before critical refresh/page timer */
	UDDRC_REGS->UDDRC_RFSHCTL0 = UDDRC_RFSHCTL0_refresh_margin(0x2) |
	/* speculative refresh, number of 32 DFI clocks until can run speculative refresh */
			UDDRC_RFSHCTL0_refresh_to_x32(0x2) |
	/* how many refresh timeouts can handle before refresh is forced to execute (+1)*/
			UDDRC_RFSHCTL0_refresh_burst(0);

	/* not disable refresh on ranks */
	UDDRC_REGS->UDDRC_RFSHCTL3 = 0;

	dbg_very_loud("PWRTMG %x\n", UDDRC_REGS->UDDRC_PWRTMG);
	dbg_very_loud("RFSHCTL0 %x\n", UDDRC_REGS->UDDRC_RFSHCTL0);
	dbg_very_loud("RFSHCTL3 %x\n", UDDRC_REGS->UDDRC_RFSHCTL3);
}

inline static void uddrc_mr1()
{
	UDDRC_REGS->UDDRC_MRCTRL0 = 0
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	/* select mode register 1 MR1 */
			| UDDRC_MRCTRL0_mr_addr(0x1)
#endif
	/* select rank 0 only */
			| UDDRC_MRCTRL0_mr_rank(0x1);
	dbg_very_loud("MRCTRL0 %x\n", UDDRC_REGS->UDDRC_MRCTRL0);
	dbg_very_loud("MRCTRL1 %x\n", UDDRC_REGS->UDDRC_MRCTRL1);
}

inline static void uddrc_ena_hw_lowpwr()
{
	/* Enable low power control */
	UDDRC_REGS->UDDRC_HWLPCTL = UDDRC_HWLPCTL_hw_lp_exit_idle_en |
			UDDRC_HWLPCTL_hw_lp_idle_x32(0x4d);
}

inline static void uddrc_rank_refresh()
{
	unsigned long t_rfc_nom_x32 =
		DIV_ROUND_DOWN(NS_TO_CYCLES_DOWN(tREFI), 32*2);
	unsigned long t_rfc_min =
		DIV_ROUND_UP(NS_TO_CYCLES_UP(tRFC), 2);

	/* tREFI - average time interface between refreshes per rank. tREFI/2 */
	UDDRC_REGS->UDDRC_RFSHTMG =
			UDDRC_RFSHTMG_t_rfc_nom_x32(t_rfc_nom_x32) |
	/* tRFCmin - min time from refresh to refresh or activate (tRFCmin/tCK)/2*/
			UDDRC_RFSHTMG_t_rfc_min(t_rfc_min);

	dbg_very_loud("RFSHTMG %x\n", UDDRC_REGS->UDDRC_RFSHTMG);
}

inline static void uddrc_config_zq_calib()
{
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	unsigned long t_zq_short_interval_x1024 =
		DIV_ROUND_UP(NS_TO_CYCLES_UP(128000000UL), 1024*2) + 1;
#endif
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	unsigned long t_zq_short_interval_x1024 =
		DIV_ROUND_UP(NS_TO_CYCLES_UP(4000000UL), 1024*2) + 1;
#endif

	/* enable generation of auto ZQ calibration by controller */
	UDDRC_REGS->UDDRC_ZQCTL0 =
	/* ZQ calibration resistor is not shared between ranks */
	/* number of DFI nops to run after ZQ calibration long */
			UDDRC_ZQCTL0_t_zq_long_nop(DIV_ROUND_UP(TZQOPER, 2)) |
	/* number of DFI nops to run after ZQ calibration short */
			UDDRC_ZQCTL0_t_zq_short_nop(DIV_ROUND_UP(TZQCS, 2));

	UDDRC_REGS->UDDRC_ZQCTL1 =
			UDDRC_ZQCTL1_t_zq_short_interval_x1024(t_zq_short_interval_x1024)
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
			| UDDRC_ZQCTL1_t_zq_reset_nop(DIV_ROUND_UP(TZQRESET, 2))
#endif
			;
	dbg_very_loud("ZQCTL0 %x\n", UDDRC_REGS->UDDRC_ZQCTL0);
	dbg_very_loud("ZQCTL1 %x\n", UDDRC_REGS->UDDRC_ZQCTL1);
}

inline static void uddrc_dis_hw_lowpwr()
{
	/* disable hardware low power features */
	UDDRC_REGS->UDDRC_HWLPCTL = 0;
}

inline static void uddrc_config_odt_timings()
{
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
	/* DFI clocks to hold ODT for a write, 6 is recommended */
	UDDRC_REGS->UDDRC_ODTCFG = UDDRC_ODTCFG_wr_odt_hold(0x6) |
	/* DFI clocks to delay between write and setting ODT values , 0 recommended */
			UDDRC_ODTCFG_wr_odt_delay(0x0) |
	/* DFI clocks to hold ODT for a read, 6 is recommended */
			UDDRC_ODTCFG_rd_odt_hold(0x6) |
	/* DFI clocks to delay between read and setting ODT values , recommended CL - CWL */
			UDDRC_ODTCFG_rd_odt_delay(CL - CWL);

	/* which remote ODTs must be turned on during a read on rank X (we have just rank 0) */
	UDDRC_REGS->UDDRC_ODTMAP = UDDRC_ODTMAP_rank0_wr_odt | UDDRC_ODTMAP_rank0_rd_odt;
#endif
/* LPDDR2 does not use ODT */
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	UDDRC_REGS->UDDRC_ODTCFG = 0;
	UDDRC_REGS->UDDRC_ODTMAP = 0;
#endif
	dbg_very_loud("ODTCFG %x\n", UDDRC_REGS->UDDRC_ODTCFG);
	dbg_very_loud("ODTMAP %x\n", UDDRC_REGS->UDDRC_ODTMAP);
}

inline static void uddrc_dis_crc_parity()
{
	/* do not enable any CRC parity control */
	UDDRC_REGS->UDDRC_CRCPARCTL0 = 0;
}

inline static void uddrc_dimmctl()
{
	/* do not have address mirroring, do not have multi rank */
	UDDRC_REGS->UDDRC_DIMMCTL = 0;
}

inline static void uddrc_config_dfi()
{
	unsigned long rddata_en = DIV_ROUND_DOWN(RL - 2, 2);

	/* DFI clocks after assert/deassert of signals until PHY-DRAM reflects them */
	UDDRC_REGS->UDDRC_DFITMG0 = UDDRC_DFITMG0_dfi_t_ctrl_delay(0x3) |
	/* select to use DFI or DFI PHY clock for dfi signals */
	/* clocks from assertion of DFI signals until dfi_rddata_en -phy specific*/
			UDDRC_DFITMG0_dfi_t_rddata_en(rddata_en) |
	/* clocks from assertion of wrdata_en to write data is driven */
			UDDRC_DFITMG0_dfi_tphy_wrdata(0x0) |
	/* clocks from write command to wrdata_en tWRlat*/
			UDDRC_DFITMG0_dfi_tphy_wrlat(DIV_ROUND_DOWN(WL
#if defined(CONFIG_DDR3) || defined(CONFIG_DDR2)
			 - 1
#endif
			, 2));

	/* clocks between dfi_cs asserted and dfi_parity_in is driven */
	UDDRC_REGS->UDDRC_DFITMG1 =
	/* clocks between dfi_wrdata_en asserted and data transfered completed on bus */
			UDDRC_DFITMG1_dfi_t_wrdata_delay(0x0)	|
	/* clocks between dfi_dram_clk_disable asserted until clock to DRAM maintains low value */
			UDDRC_DFITMG1_dfi_t_dram_clk_disable(0x2) |
	/* clocks between dfi_dram-clk_disable asserted until first rising edge to DRAM clock */
			UDDRC_DFITMG1_dfi_t_dram_clk_enable(0x2);

	/* clocks for tlp_resp time low power response time - assumed*/
	UDDRC_REGS->UDDRC_DFILPCFG0 = UDDRC_DFILPCFG0_dfi_tlp_resp(0x7) |
	/* clocks to drive on dfi_lp_wakeup when enter self refresh - 4096 cycles should be 0x8 or 0x80 ? */
			UDDRC_DFILPCFG0_dfi_lp_wakeup_sr(0x7) 	|
	/* Enables DFI Low Power interface handshaking during Self Refresh Entry/Exit */
			UDDRC_DFILPCFG0_dfi_lp_en_sr |
	/* clocks to drive on dfi_lp_wakeup when enter power down - 16 cycles should be 0x3 or 0x0 ? */
			UDDRC_DFILPCFG0_dfi_lp_wakeup_pd(0x7) |
	/* Enables DFI Low Power interface handshaking during Power Down Entry/Exit */
			UDDRC_DFILPCFG0_dfi_lp_en_pd;

	/* enable automatic generation of dfi_ctrlupd_req */
	UDDRC_REGS->UDDRC_DFIUPD0 =
	/* max clocks that dfi_ctrlupd_req can assert */
			UDDRC_DFIUPD0_dfi_t_ctrlup_max(0x40) |
	/* min clocks that dfi_ctrlupd_req can assert */
			UDDRC_DFIUPD0_dfi_t_ctrlup_min(0x3);

	/* min number of 1024 clocks between DFI update requests initiated by umctl2 */
	UDDRC_REGS->UDDRC_DFIUPD1 = UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024(0x3) |
	/* max number of 1024 clocks between DFI update requests initiated by umctl2 */
			UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024(0x8d);

	/* enable support for acking phy initiated updates */
	UDDRC_REGS->UDDRC_DFIUPD2 = UDDRC_DFIUPD2_dfi_phyupd_en;

	/* disable PHY master interface */
	UDDRC_REGS->UDDRC_DFIPHYMSTR = 0;
	dbg_very_loud("DFITMG0 %x\n", UDDRC_REGS->UDDRC_DFITMG0);
	dbg_very_loud("DFITMG1 %x\n", UDDRC_REGS->UDDRC_DFITMG1);
	dbg_very_loud("DFILPCFG0 %x\n", UDDRC_REGS->UDDRC_DFILPCFG0);
	dbg_very_loud("DFIUPD0 %x\n", UDDRC_REGS->UDDRC_DFIUPD0);
	dbg_very_loud("DFIUPD1 %x\n", UDDRC_REGS->UDDRC_DFIUPD1);
}

/* Main entry point of the UMCTL2 DRAM driver.
 * state is a preconfigured umctl2_config.
 * The driver will initialize the Controller and then turn back control.
 * state is useless after this call returns.
 */
int umctl2_init (struct umctl2_config_state *state)
{
	unsigned int ret = 0;

#ifndef CONFIG_SYS_BASE_UMCTL2
#error "CONFIG_SYS_BASE_UMCTL2 undefined"
#endif

#ifndef CONFIG_SYS_BASE_UMCTL2_MP
#error "CONFIG_SYS_BASE_UMCTL2_MP undefined"
#endif

	UDDRC_REGS = (struct uddrc_regs *) CONFIG_SYS_BASE_UMCTL2;
	UDDRC_MP = (struct uddrc_mp *) CONFIG_SYS_BASE_UMCTL2_MP;

	/* hold a global pointer to the state */
	umctl2_config = state;

	tRFC = timings.tRFC;
	tREFI = timings.tREFI;
	tWR = timings.tWR;
	tRP = timings.tRP;
	tRP_ps = timings.tRP_ps;
	tRCD = timings.tRCD;
	tRCD_ps = timings.tRCD_ps;
	tCCD = timings.tCCD;
	tRAS = timings.tRAS;
	tRASMAX = timings.tRASMAX;
	tRC_ps = timings.tRC_ps;
	tFAW = timings.tFAW;
#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3)
	tPRECKE = timings.tPRECKE;
#endif
	tPOSTCKE = timings.tPOSTCKE;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	tDQSCK_MIN = timings.tDQSCK_MIN;
	tDQSCK_MAX = timings.tDQSCK_MAX;
#endif
	CL = timings.CL;
	CWL = timings.CWL;
	AL = timings.AL;
#if defined(CONFIG_LPDDR2) || defined(CONFIG_LPDDR3)
	RL = timings.RL;
	WL = timings.WL;
#endif
	TZQOPER = timings.TZQOPER;
	TZQCS = timings.TZQCS;
	MRD = timings.MRD;

#ifdef CONFIG_RSTC
	/*
	 * STEP 0: Deassert reset signal DDR PHY. This ensures that the
	 * PHY starts.
	 */
	rstc_ddr_phy_rst_deassert();
#else
#error "UMCTL2 requires RSTC access for ddrc_phy_rstn signal"
#endif
	/* wait for clock synchronization */
	udelay(100);

	/* STEP 1
	 * Program the DWC_ddr_umctl2 registers
	 * (Ensure that DFIMISC.dfi_init_complete_en is set to 0)
	 */

	UDDRC_REGS->UDDRC_DFIMISC = 0;

	/* disable dequeuing of rd/wr commands from CAM to DRAM */
	uddrc_dis_dq();
	/* configure master register MSTR - type of DDR, burst length*/
	uddrc_mstr();

	/* enable hardware low power features */
	uddrc_ena_hw_lowpwr();
	/* configure refresh-related parameters */
	uddrc_configure_refresh();
	/* disable crc parity control */
	uddrc_dis_crc_parity();

	/* configure INIT0 - clocks to wait before and after CKE is asserted,
	 * mode of operation after reset is removed
	 */
	uddrc_init0();
	/* configure INIT1 - clocks to assert DRAM reset during init,
	 * pre-ocd complete delay
	 */
	uddrc_init1();
	/* configure INIT2 - Min stable clock, idle after reset */
	uddrc_init2();
	/* configure INIT3 - DRAM registers MR0 and MR1 */
	uddrc_init3();
	/* configure INIT4 - DRAM registers MR2 and MR3 */
	uddrc_init4();
	/* configure INIT5 - initial calibration time */
	uddrc_init5();
	/* configure DRAM timings */
	uddrc_dramtmg();
	/* configure mapping between HIF addr to bank, column, row
	 * internal address
	 */
	uddrc_addrmap_init();
	/* configure refresh per rank */
	uddrc_rank_refresh();
	/* configure zq calibration timing */
	uddrc_config_zq_calib();
	/* configure DFI clocks to wait on ODT rd/wr commands */
	uddrc_config_odt_timings();
	/* configure DIMM control */
	uddrc_dimmctl();
	/* configure DFI interface to PHY */
	uddrc_config_dfi();

	/* multi-port register settings (urgent bit are not connected in the DESIGN) */
	uddrc_mp_setup();

#ifdef CONFIG_RSTC
	/* STEP 2
	 * Deassert reset signal core_ddrc_rstn ! (mandatory)
	 */
	rstc_ddr_rst_deassert();
#else
#error "UMCTL2 requires RSTC access for core_ddrc_rstn signal"
#endif

#if defined POWER_CYCLE_ON_PHY_INIT
	UDDRC_REGS->UDDRC_DBG1 = 0;
	/* power down and power up ? */
	UDDRC_REGS->UDDRC_PWRCTL &= UDDRC_PWRCTL_powerdown_en;
	UDDRC_REGS->UDDRC_PWRCTL = 0;
#endif

#if defined(CONFIG_DDR2) || defined(CONFIG_DDR3) || defined(CONFIG_LPDDR2)
	UDDRC_REGS->UDDRC_PWRCTL = UDDRC_PWRCTL_en_dfi_dram_clk_disable;
#endif

	/* STEP 3
	 * Init PHY
	 */
#if defined CONFIG_PUBL
	state->phy_init(&timings);
#else
	state->phy_init(NULL);
#endif
	/* STEP 4
	 * Wait PHY init
	 */
	if (state->phy_idone())
		return -1;

	/* STEP 5
	 * Enable uMCTL2 dfi init start signal
	 * in case of DRAM init is done via PHY: start PHY
	 */
	if (state->phy_start())
		return -1;

	/* STEP 6
	 * Set SWCTL.sw_done to 0 : to enable programming of
	 * quasi-static registers outside of reset
	 */
	/* DO NOT Check SWSTAT.sw_done_ack = 1 */
	UDDRC_REGS->UDDRC_SWCTL = 0;

	/* STEP 7
	 * set DFI init complete to 1
	 */
	UDDRC_REGS->UDDRC_DFIMISC = UDDRC_DFIMISC_dfi_init_complete_en;

	/* STEP 8
	 * Set software done to 1; also,
	 * Only quasi-dinamic registers can be modified from now on
	 */
	UDDRC_REGS->UDDRC_SWCTL	= UDDRC_SWCTL_sw_done;

	/* STEP 9
	 * Wait for SWSTAT.sw_done_ack to become 1 :
	 * ONLY AFTER UDDRC_SWCTL_sw_done
	 */
	WAIT_WHILE_COND((UDDRC_REGS->UDDRC_SWSTAT != UDDRC_SWSTAT_sw_done_ack), 0xA6);

	/* STEP 10
	 * Check for uMCTL2 in normal mode
	 */
	WAIT_WHILE_COND (((UDDRC_REGS->UDDRC_STAT & UDDRC_STAT_operating_mode_Msk) !=
		UDDRC_STAT_operating_mode_Normal), 10000);

	/* STEP 11
	 * Disable auto-refreshes
	 */
	UDDRC_REGS->UDDRC_RFSHCTL3 |= UDDRC_RFSHCTL3_dis_auto_refresh;
	UDDRC_REGS->UDDRC_RFSHCTL3 ^= UDDRC_RFSHCTL3_refresh_update_level;

	/* STEP 11a
	 * Mask transitions in phy_dfi_init_complete during PHY training
	 */
	UDDRC_REGS->UDDRC_SWCTL = 0;
	UDDRC_REGS->UDDRC_DFIMISC &= ~UDDRC_DFIMISC_dfi_init_complete_en;
	UDDRC_REGS->UDDRC_SWCTL = UDDRC_SWCTL_sw_done;

	WAIT_WHILE_COND((UDDRC_REGS->UDDRC_SWSTAT != UDDRC_SWSTAT_sw_done_ack), 0xA6);

	/* STEP 12
	 * Train PHY
	 */
	ret = state->phy_train();
	if (ret)
		return ret;

	/* STEP 13
	 * Revert steps 11 and 11a
	 */
	UDDRC_REGS->UDDRC_RFSHCTL3 &= ~UDDRC_RFSHCTL3_dis_auto_refresh;
	UDDRC_REGS->UDDRC_RFSHCTL3 ^= UDDRC_RFSHCTL3_refresh_update_level;

	UDDRC_REGS->UDDRC_SWCTL = 0;
	UDDRC_REGS->UDDRC_DFIMISC |= UDDRC_DFIMISC_dfi_init_complete_en;
	UDDRC_REGS->UDDRC_SWCTL = UDDRC_SWCTL_sw_done;

	WAIT_WHILE_COND((UDDRC_REGS->UDDRC_SWSTAT != UDDRC_SWSTAT_sw_done_ack), 0xA6);

	/* STEP 14
	 * AXI ports can now take transactions
	 */
	if (umctl2_config->axi_port_bitmap &  MP_AXI_PORT_ENABLE(0))
		UDDRC_MP->UDDRC_PCTRL_0	= UDDRC_PCTRL_0_port_en;
	if (umctl2_config->axi_port_bitmap &  MP_AXI_PORT_ENABLE(1))
		UDDRC_MP->UDDRC_PCTRL_1	= UDDRC_PCTRL_1_port_en;
	if (umctl2_config->axi_port_bitmap &  MP_AXI_PORT_ENABLE(2))
		UDDRC_MP->UDDRC_PCTRL_2	= UDDRC_PCTRL_2_port_en;
	if (umctl2_config->axi_port_bitmap &  MP_AXI_PORT_ENABLE(3))
		UDDRC_MP->UDDRC_PCTRL_3	= UDDRC_PCTRL_3_port_en;
	if (umctl2_config->axi_port_bitmap &  MP_AXI_PORT_ENABLE(4))
		UDDRC_MP->UDDRC_PCTRL_4	= UDDRC_PCTRL_4_port_en;

	return ret;
}

