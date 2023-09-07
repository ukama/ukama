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

#ifndef _UMCTL2_H_
#define _UMCTL2_H_

#define __I  volatile const /**< Defines 'read-only'  permissions */
#define __O  volatile       /**< Defines 'write-only' permissions */
#define __IO volatile       /**< Defines 'read/write' permissions */

/*
 * Synopsys UMCTL2 register map for multi port configuration.
 */

struct uddrc_mp {
	/* Unused */
	__I  unsigned int UDDRC_NOT_USED;
	/* Port status register */
	__I  unsigned int UDDRC_PSTAT;
	/* Port common configuration register */
	__IO unsigned int UDDRC_PCCFG;
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_0;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_0;
	__I  unsigned int Reserved1[33];
	/* Port n Control register */
	__IO unsigned int UDDRC_PCTRL_0;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_0;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_0;
	/* Port n Write QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_0;
	/* Port n Write QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_0;
	__I  unsigned int Reserved2[4];
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_1;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_1;
	__I  unsigned int UDDRC_DUMMY1;
	__I  unsigned int Reserved3[32];
	/* Port n Control Register */
	__IO unsigned int UDDRC_PCTRL_1;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_1;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_1;
	/* Port n Write QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_1;
	/* Port n Write QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_1;
	__I  unsigned int Reserved4[4];
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_2;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_2;
	__I  unsigned int Reserved5[33];
	/* Port n Control register */
	__IO unsigned int UDDRC_PCTRL_2;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_2;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_2;
	/* Port n Write QoS Confguration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_2;
	/* Port n Write QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_2;
	__I  unsigned int Reserved6[4];
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_3;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_3;
	__I  unsigned int Reserved7[33];
	/* Port n Control Register */
	__IO unsigned int UDDRC_PCTRL_3;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_3;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_3;
	/* Port n Write QoS Confguration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_3;
	/* Port n Write QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_3;
	__I  unsigned int Reserved8[4];
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_4;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_4;
	__I  unsigned int Reserved9[33];
	/* Port n Control Register */
	__IO unsigned int UDDRC_PCTRL_4;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_4;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_4;
	/* Port n Write QoS Confguration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_4;
	/* Port n Write QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_4;
	__I  unsigned int Reserved10[4];
	/* Port n Configuration Read register */
	__IO unsigned int UDDRC_PCFGR_5;
	/* Port n Configuration Write register */
	__IO unsigned int UDDRC_PCFGW_5;
	__I  unsigned int Reserved11[33];
	/* Port n Control Register */
	__IO unsigned int UDDRC_PCTRL_5;
	/* Port n Read QoS Configuration register 0 */
	__IO unsigned int UDDRC_PCFGQOS0_5;
	/* Port n Read QoS Configuration register 1 */
	__IO unsigned int UDDRC_PCFGQOS1_5;
	/* Port n Write QoS Confguration register 0 */
	__IO unsigned int UDDRC_PCFGWQOS0_5;
	/* Port n Write QoS Confguration register 1 */
	__IO unsigned int UDDRC_PCFGWQOS1_5;
	__I  unsigned int Reserved12[444];
	/* SAR Base Address Register 0 */
	__IO unsigned int UDDRC_SARBASE0;
	/* SAR Size Register 0 */
	__IO unsigned int UDDRC_SARSIZE0;
};

/* UMCTL2 MP register helpers */
/* { */
/* -------- UDDRC_PSTAT : (UDDRC_MP Offset: 0x4) Port Status Register -------- */
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 0. */
#define UDDRC_PSTAT_rd_port_busy_0 (0x1u << 0)
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 1. */
#define UDDRC_PSTAT_rd_port_busy_1 (0x1u << 1)
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 2. */
#define UDDRC_PSTAT_rd_port_busy_2 (0x1u << 2)
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 3. */
#define UDDRC_PSTAT_rd_port_busy_3 (0x1u << 3)
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 4. */
#define UDDRC_PSTAT_rd_port_busy_4 (0x1u << 4)
/* (UDDRC_PSTAT) Indicates if there are outstanding reads for AXI port 5. */
#define UDDRC_PSTAT_rd_port_busy_5 (0x1u << 5)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 0. */
#define UDDRC_PSTAT_wr_port_busy_0 (0x1u << 16)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 1. */
#define UDDRC_PSTAT_wr_port_busy_1 (0x1u << 17)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 2. */
#define UDDRC_PSTAT_wr_port_busy_2 (0x1u << 18)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 3. */
#define UDDRC_PSTAT_wr_port_busy_3 (0x1u << 19)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 4. */
#define UDDRC_PSTAT_wr_port_busy_4 (0x1u << 20)
/* (UDDRC_PSTAT) Indicates if there are outstanding writes for AXI port 5. */
#define UDDRC_PSTAT_wr_port_busy_5 (0x1u << 21)

/* -------- UDDRC_PCCFG : (UDDRC_MP Offset: 0x8)
 * Port Common Configuration Register --------
 */
/* (UDDRC_PCCFG) If set to 1 (enabled), sets co_gs_go2critical_wr and
 * co_gs_go2critical_lpr/co_gs_go2critical_hpr signals going to DDRC based on
 * urgent input (awurgent, arurgent) coming from AXI master.
 * If set to 0 (disabled), co_gs_go2critical_wr and
 * co_gs_go2critical_lpr/co_gs_go2critical_hpr signals at DDRC
 * are driven to 1b'0.
 */
#define UDDRC_PCCFG_go2critical_en (0x1u << 0)
/* (UDDRC_PCCFG) Page match four limit.
 * If set to 1, limits the number of consecutive same page DDRC transactions
 * that can be granted by the Port Arbiter to four when Page Match feature
 * is enabled.
 * If set to 0, there is no limit imposed on number of consecutive
 * same page DDRC transactions.
 */
#define UDDRC_PCCFG_pagematch_limit (0x1u << 4)
/* (UDDRC_PCCFG) Burst length expansion mode.
 * By default (i.e. bl_exp_mode==0) XPI expands every AXI burst
 * into multiple HIF commands, using the memory burst length as a unit.
 * If set to 1, then XPI will use half of the memory burst length as a unit.
 * This applies to both reads and writes.
 * When MSTR.data_bus_width==00, setting bl_exp_mode to 1 has no effect.
 * This can be used in cases where Partial Writes is enabled
 * (UMCTL2_PARTIAL_WR=1), in order to avoid or minimize t_ccd_l penalty
 * in DDR4 and t_ccd_mw penalty in LPDDR4.
 * Hence, bl_exp_mode=1 is only recommended if DDR4 or LPDDR4.
 *
 * Note that if DBICTL.dm_en=0, functionality is not supported in the
 * following cases:
 *   - UMCTL2_PARTIAL_WR=0
 *   - UMCTL2_PARTIAL_WR=1, MSTR.data_bus_width=01,
 * MEMC_BURST_LENGTH=8 and MSTR.burst_rdwr=1000 (LPDDR4 only)
 *   - UMCTL2_PARTIAL_WR=1, MSTR.data_bus_width=01,
 * MEMC_BURST_LENGTH=4 and MSTR.burst_rdwr=0100 (DDR4 only),
 * with either MSTR.burstchop=0 or CRCPARCTL1.crc_enable=1
 *
 * Functionality is also not supported if Data Channel Interleave is enabled
 */
#define UDDRC_PCCFG_bl_exp_mode (0x1u << 8)
/* -------- UDDRC_PCFGR_0 : (UDDRC_MP Offset: 0xC)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_0) Port priority :
 * Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port.
 * The aging counters down-count every clock cycle where the port
 * is requesting but not granted.
 * The higher significant 5-bits of the read aging counter sets the
 * priority of the read channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the
 * counter starts to decrease.
 * When the aging counter becomes 0, the corresponding port channel will have
 * the highest priority level (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority
 * (2 priority levels: low priority read - LPR, high priority read - HPR) on a
 * command by command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_0_rd_port_priority_Pos 0
#define UDDRC_PCFGR_0_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_0_rd_port_priority_Pos)
#define UDDRC_PCFGR_0_rd_port_priority(value) \
	(UDDRC_PCFGR_0_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_0_rd_port_priority_Pos))

/* (UDDRC_PCFGR_0) bypass read reordering */
#define UDDRC_PCFGR_0_read_reorder_bypass_en (0x1u << 11)

/* (UDDRC_PCFGR_0) If set to 1, enables aging function for the read channel
 * of the port.
 */
#define UDDRC_PCFGR_0_rd_port_aging_en (0x1u << 12)

/*(UDDRC_PCFGR_0) port urgent enable.
 * If set to 1, enables the AXI urgent sideband signal (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr
 * signal to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking (it is not associated with
 * any particular command).
 */
#define UDDRC_PCFGR_0_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_0) port match enable.
 * If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same memory page
 * (same bank and same row). See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_0_rd_port_pagematch_en (0x1u << 14)

/* (UDDRC_PCFGR_0) read write ordered enable */

#define UDDRC_PCFGR_0_rdwr_ordered_en (0x1u << 16)

/* -------- UDDRC_PCFGW_0 : (UDDRC_MP Offset: 0x10)
 Port n Configuration Write Register --------
 */
#define UDDRC_PCFGW_0_wr_port_priority_Pos 0

/* (UDDRC_PCFGW_0) Write port priority
 * Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port.
 * The aging counters down-count every clock cycle where the port is requesting
 * but not granted.
 * The higher significant 5-bits of the write aging counter sets the initial
 * priority of the write channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the
 * counter starts to decrease. When the aging counter becomes 0, the
 * corresponding port channel will have the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set
 * port priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when
 * they timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_0_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_0_wr_port_priority_Pos)
#define UDDRC_PCFGW_0_wr_port_priority(value) \
	((UDDRC_PCFGW_0_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_0_wr_port_priority_Pos)))

/* (UDDRC_PCFGW_0) Write port aging enable.
 * If set to 1, enables aging function for the
 * write channel of the port.
 */
#define UDDRC_PCFGW_0_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_0) Write port pagematch enable.
 * If set to 1, enables the AXI urgent sideband signal (awurgent).
 * When enabled and awurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_wr signal to DDRC is asserted
 * if enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_0_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_0) Write port pagematch enable.
 * If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_0_wr_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCTRL_0 : (UDDRC_MP Offset: 0x98)
 * Port n Control Register --------
 */
/* (UDDRC_PCTRL_0) Enables AXI port n. */
#define UDDRC_PCTRL_0_port_en (0x1u << 0)

/* -------- UDDRC_PCFGQOS0_0 : (UDDRC_MP Offset: 0x9C)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_0) map level 1
 * Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0.
 * Possible values for level1 are 0 to 13 (for dual RAQ) or 0 to 14
 * (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_0_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_0_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_0_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_0_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_0_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_0_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_0) map level 2
 * Separation level2 indicating the end of region1 mapping.
 */
#define UDDRC_PCFGQOS0_0_rqos_map_level2_Pos 8
#define UDDRC_PCFGQOS0_0_rqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGQOS0_0_rqos_map_level2_Pos)
#define UDDRC_PCFGQOS0_0_rqos_map_level2(value) \
	((UDDRC_PCFGQOS0_0_rqos_map_level2_Msk & \
	((value) << UDDRC_PCFGQOS0_0_rqos_map_level2_Pos)))

/* (UDDRC_PCFGQOS0_0) Region 0 qos value.
 * This bitfield indicates the traffic class of region 0.
 * Valid values are:
 * 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are:
 * 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class
 * of region0 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_0_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_0_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_0_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_0_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_0_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_0_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_0) Region 1 qos value.
 * This bitfield indicates the traffic class of region 1.
 * Valid values are:
 * 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are
 * 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class
 * of region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_0_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_0_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_0_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_0_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_0_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_0_rqos_map_region1_Pos)))

/* (UDDRC_PCFGQOS0_0) Region 2 qos value.
 * This bitfield indicates the traffic class of region 2.
 * Valid values are:
 * 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are
 * 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class
 * of region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_0_rqos_map_region2_Pos 24
#define UDDRC_PCFGQOS0_0_rqos_map_region2_Msk \
	(0x3u << UDDRC_PCFGQOS0_0_rqos_map_region2_Pos)
#define UDDRC_PCFGQOS0_0_rqos_map_region2(value) \
	((UDDRC_PCFGQOS0_0_rqos_map_region2_Msk & \
	((value) << UDDRC_PCFGQOS0_0_rqos_map_region2_Pos)))

/* -------- UDDRC_PCFGQOS1_0 : (UDDRC_MP Offset: 0xA0)
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_0) rqos map timeout blue.
 * Specifies the timeout value for transactions mapped to the blue address queue.
 */
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_0_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_0_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_0_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_0) rqos map timeout red.
 * Specifies the timeout value for transactions mapped to the red address queue.
 */
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_0_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_0_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_0_rqos_map_timeoutr_Msk & \
	((value) << UDDRC_PCFGQOS1_0_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_0 : (UDDRC_MP Offset: 0xA4)
 * Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_0) wqos map level 1.
 * Separation level indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 14
 * which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_0_wqos_map_level1_Pos 0
#define UDDRC_PCFGWQOS0_0_wqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGWQOS0_0_wqos_map_level1_Pos)
#define UDDRC_PCFGWQOS0_0_wqos_map_level1(value) \
	((UDDRC_PCFGWQOS0_0_wqos_map_level1_Msk & \
	((value) << UDDRC_PCFGWQOS0_0_wqos_map_level1_Pos)))

/* (UDDRC_PCFGWQOS0_0) wqos map level 2.
 * Separation level indicating the end of region1 mapping;
 * start of region1 is level1 + 1. Possible values for level2 are (leve1+1)
 * to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_0_wqos_map_level2_Pos 8
#define UDDRC_PCFGWQOS0_0_wqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGWQOS0_0_wqos_map_level2_Pos)
#define UDDRC_PCFGWQOS0_0_wqos_map_level2(value) \
	((UDDRC_PCFGWQOS0_0_wqos_map_level2_Msk & \
	((value) << UDDRC_PCFGWQOS0_0_wqos_map_level2_Pos)))

/* (UDDRC_PCFGWQOS0_0) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0)
 * and traffic class of region0 is set to 1 (VPW),
 * VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_0_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_0_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_0_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_0_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_0_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_0_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_0) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0)
 * and traffic class of region 1 is set to 1 (VPW),
 * VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_0_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_0_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_0_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_0_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_0_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_0_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_0 : (UDDRC_MP Offset: 0xA8)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_0) Specifies the timeout value for write transactions.
 */
#define UDDRC_PCFGWQOS1_0_wqos_map_timeout_Pos 0
#define UDDRC_PCFGWQOS1_0_wqos_map_timeout_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_0_wqos_map_timeout_Pos)
#define UDDRC_PCFGWQOS1_0_wqos_map_timeout(value) \
	((UDDRC_PCFGWQOS1_0_wqos_map_timeout_Msk & \
	((value) << UDDRC_PCFGWQOS1_0_wqos_map_timeout_Pos)))

/* -------- UDDRC_PCFGR_1 : (UDDRC_MP Offset: 0xBC)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_1) Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset,
 * or after each grant to the corresponding port.
 * The aging counters down-count every clock cycle where the port
 * is requesting but not granted.
 * The higher significant 5-bits of the read aging counter sets the priority
 * of the read channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the counter
 * starts to decrease. When the aging counter becomes 0, the corresponding port
 * channel will have the highest priority level (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when
 * they timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority
 * (2 priority levels: low priority read - LPR, high priority read - HPR)
 * on a command by command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_1_rd_port_priority_Pos 0
#define UDDRC_PCFGR_1_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_1_rd_port_priority_Pos)
#define UDDRC_PCFGR_1_rd_port_priority(value) \
	((UDDRC_PCFGR_1_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_1_rd_port_priority_Pos)))

/* (UDDRC_PCFGR_1) bypass read reordering */
#define UDDRC_PCFGR_1_read_reorder_bypass_en (0x1u << 11)

/* (UDDRC_PCFGR_1) If set to 1, enables aging function for the read channel of the port. */
#define UDDRC_PCFGR_1_rd_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGR_1) If set to 1, enables the AXI urgent sideband signal
 * (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr signal
 * to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 *(it is not associated with any particular command).
 */
#define UDDRC_PCFGR_1_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_1) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_1_rd_port_pagematch_en (0x1u << 14)

/* (UDDRC_PCFGR_0) read write ordered enable */
#define UDDRC_PCFGR_1_rdwr_ordered_en (0x1u << 16)

/* -------- UDDRC_PCFGW_1 : (UDDRC_MP Offset: 0xC0)
 * Port n Configuration Write Register --------
 */
/* (UDDRC_PCFGW_1) Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset,
 * or after each grant to the corresponding port.
 * The aging counters down-count every clock cycle where the port is requesting
 * but not granted. The higher significant 5-bits of the write aging counter
 * sets the initial priority of the write channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the counter
 * starts to decrease. When the aging counter becomes 0, the corresponding port
 * channel will have the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_1_wr_port_priority_Pos 0
#define UDDRC_PCFGW_1_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_1_wr_port_priority_Pos)
#define UDDRC_PCFGW_1_wr_port_priority(value) \
	((UDDRC_PCFGW_1_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_1_wr_port_priority_Pos)))

/* (UDDRC_PCFGW_1) If set to 1, enables aging function for the
 * write channel of the port.
 */
#define UDDRC_PCFGW_1_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_1) If set to 1, enables the AXI urgent sideband signal
 * (awurgent). When enabled and awurgent is asserted by the master,
 * that port becomes the highest priority and co_gs_go2critical_wr signal to
 * DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be  asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_1_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_1) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_1_wr_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCTRL_1 : (UDDRC_MP Offset: 0x148)
 * Port n Control Register --------
 */
/* (UDDRC_PCTRL_1) Enables AXI port n.
 */
#define UDDRC_PCTRL_1_port_en (0x1u << 0)

/* -------- UDDRC_PCFGQOS0_1 : (UDDRC_MP Offset: 0x14C)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_1) Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 13 (for dual RAQ)
 * or 0 to 14 (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_1_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_1_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_1_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_1_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_1_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_1_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_1) Separation level2 indicating the end of region1 mapping;
 */
#define UDDRC_PCFGQOS0_1_rqos_map_level2_Pos 8
#define UDDRC_PCFGQOS0_1_rqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGQOS0_1_rqos_map_level2_Pos)
#define UDDRC_PCFGQOS0_1_rqos_map_level2(value) \
	((UDDRC_PCFGQOS0_1_rqos_map_level2_Msk & \
	((value) << UDDRC_PCFGQOS0_1_rqos_map_level2_Pos)))

/* (UDDRC_PCFGQOS0_1) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are: 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of region0
 * is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_1_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_1_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_1_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_1_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_1_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_1_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_1) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class
 * of region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_1_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_1_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_1_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_1_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_1_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_1_rqos_map_region1_Pos)))

/* (UDDRC_PCFGQOS0_1) This bitfield indicates the traffic class of region 2.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region2 maps to the red? address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class
 * of region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_1_rqos_map_region2_Pos 24
#define UDDRC_PCFGQOS0_1_rqos_map_region2_Msk \
	(0x3u << UDDRC_PCFGQOS0_1_rqos_map_region2_Pos)
#define UDDRC_PCFGQOS0_1_rqos_map_region2(value) \
	((UDDRC_PCFGQOS0_1_rqos_map_region2_Msk & \
	((value) << UDDRC_PCFGQOS0_1_rqos_map_region2_Pos)))

/* -------- UDDRC_PCFGQOS1_1 : (UDDRC_MP Offset: 0x150)
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_1) Specifies the timeout value for transactions
 * mapped to the blue address queue.
 */
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_1_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_1_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_1_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_1) Specifies the timeout value for transactions
 * mapped to the red address queue.
 */
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_1_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_1_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_1_rqos_map_timeoutr_Msk & \
	((value) << UDDRC_PCFGQOS1_1_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_1 : (UDDRC_MP Offset: 0x154)
 * Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_1) Separation level indicating the end of region0 mapping;
 * start of region0 is 0.
 * Possible values for level1 are 0 to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_1_wqos_map_level1_Pos 0
#define UDDRC_PCFGWQOS0_1_wqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGWQOS0_1_wqos_map_level1_Pos)
#define UDDRC_PCFGWQOS0_1_wqos_map_level1(value) \
	((UDDRC_PCFGWQOS0_1_wqos_map_level1_Msk & \
	((value) << UDDRC_PCFGWQOS0_1_wqos_map_level1_Pos)))

/* (UDDRC_PCFGWQOS0_1) Separation level indicating the end of region1 mapping;
 * start of region1 is level1 + 1.
 * Possible values for level2 are (leve1+1) to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_1_wqos_map_level2_Pos 8
#define UDDRC_PCFGWQOS0_1_wqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGWQOS0_1_wqos_map_level2_Pos)
#define UDDRC_PCFGWQOS0_1_wqos_map_level2(value) \
	((UDDRC_PCFGWQOS0_1_wqos_map_level2_Msk & \
	((value) << UDDRC_PCFGWQOS0_1_wqos_map_level2_Pos)))

/* (UDDRC_PCFGWQOS0_1) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region0
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_1_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_1_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_1_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_1_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_1_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_1_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_1) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region 1
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_1_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_1_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_1_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_1_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_1_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_1_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_1 : (UDDRC_MP Offset: 0x158)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_1) Specifies the timeout value for write transactions.
 */
#define UDDRC_PCFGWQOS1_1_wqos_map_timeout_Pos 0
#define UDDRC_PCFGWQOS1_1_wqos_map_timeout_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_1_wqos_map_timeout_Pos)
#define UDDRC_PCFGWQOS1_1_wqos_map_timeout(value) \
	((UDDRC_PCFGWQOS1_1_wqos_map_timeout_Msk & \
	((value) << UDDRC_PCFGWQOS1_1_wqos_map_timeout_Pos)))

/* -------- UDDRC_PCFGR_2 : (UDDRC_MP Offset: 0x16C)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_2) Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset,
 * or after each grant to the corresponding port.
 * The aging counters down-count every clock cycle where the port is requesting
 * but not granted. The higher significant 5-bits of the read aging counter
 * sets the priority of the read channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the counter
 * starts to decrease. When the aging counter becomes 0, the corresponding port
 * channel will have the highest priority level (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority
 * (2 priority levels: low priority read - LPR, high priority read - HPR)
 * on a command by command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_2_rd_port_priority_Pos 0
#define UDDRC_PCFGR_2_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_2_rd_port_priority_Pos)
#define UDDRC_PCFGR_2_rd_port_priority(value) \
	((UDDRC_PCFGR_2_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_2_rd_port_priority_Pos)))

/* (UDDRC_PCFGR_2) bypass read reordering */
#define UDDRC_PCFGR_2_read_reorder_bypass_en (0x1u << 11)

/* (UDDRC_PCFGR_2) If set to 1, enables aging function for the read channel
 * of the port.
 */
#define UDDRC_PCFGR_2_rd_port_aging_en (0x1u << 12)

/*(UDDRC_PCFGR_2) If set to 1, enables the AXI urgent sideband signal (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr
 * signal to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGR_2_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_2) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_2_rd_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCFGW_2 : (UDDRC_MP Offset: 0x170)
 * Port n Configuration Write Register --------
 */
/* (UDDRC_PCFGR_2) read write ordered enable */
#define UDDRC_PCFGR_2_rdwr_ordered_en (0x1u << 16)

/*(UDDRC_PCFGW_2) Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not granted.
 * The higher significant 5-bits of the write aging counter sets the initial
 * priority of the write channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the counter
 * starts to decrease. When the aging counter becomes 0, the corresponding port
 * channel will have the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_2_wr_port_priority_Pos 0
#define UDDRC_PCFGW_2_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_2_wr_port_priority_Pos)
#define UDDRC_PCFGW_2_wr_port_priority(value) \
	((UDDRC_PCFGW_2_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_2_wr_port_priority_Pos)))

/*(UDDRC_PCFGW_2) If set to 1, enables aging function for the write channel
 * of the port.
 */
#define UDDRC_PCFGW_2_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_2) If set to 1, enables the AXI urgent sideband signal (awurgent).
 * When enabled and awurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_wr signal to DDRC is asserted
 * if enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_2_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_2) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same memory page
 * (same bank and same row). See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_2_wr_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCTRL_2 : (UDDRC_MP Offset: 0x1F8)
 * Port n Control Register --------
 */
/* (UDDRC_PCTRL_2) Enables AXI port n. */
#define UDDRC_PCTRL_2_port_en (0x1u << 0)

/* -------- UDDRC_PCFGQOS0_2 : (UDDRC_MP Offset: 0x1FC)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_2) Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0.
 * Possible values for level1 are 0 to 13 (for dual RAQ)
 * or 0 to 14 (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_2_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_2_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_2_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_2_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_2_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_2_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_2) Separation level2 indicating the end of region1 mapping;
 * start of region1 is (level1 + 1).
 * Possible values for level2 are (level1 + 1) to 14 which corresponds to arqos.
 * Region2 starts from (level2 + 1) up to 15.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_2_rqos_map_level2_Pos 8
#define UDDRC_PCFGQOS0_2_rqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGQOS0_2_rqos_map_level2_Pos)
#define UDDRC_PCFGQOS0_2_rqos_map_level2(value) \
	((UDDRC_PCFGQOS0_2_rqos_map_level2_Msk & \
	((value) << UDDRC_PCFGQOS0_2_rqos_map_level2_Pos)))

/* (UDDRC_PCFGQOS0_2) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are: 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of region0
 * is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_2_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_2_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_2_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_2_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_2_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_2_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_2) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of region 1
 * is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_2_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_2_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_2_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_2_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_2_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_2_rqos_map_region1_Pos)))

/* (UDDRC_PCFGQOS0_2) This bitfield indicates the traffic class of region2.
 * For dual address queue configurations, region2 maps to the red address queue.
 * Valid values are 1: VPR and 2: HPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of region2
 * is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_2_rqos_map_region2_Pos 24
#define UDDRC_PCFGQOS0_2_rqos_map_region2_Msk \
	(0x3u << UDDRC_PCFGQOS0_2_rqos_map_region2_Pos)
#define UDDRC_PCFGQOS0_2_rqos_map_region2(value) \
	((UDDRC_PCFGQOS0_2_rqos_map_region2_Msk & \
	((value) << UDDRC_PCFGQOS0_2_rqos_map_region2_Pos)))
/* -------- UDDRC_PCFGQOS1_2 : (UDDRC_MP Offset: 0x200)
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_2) Specifies the timeout value for transactions
 * mapped to the blue address queue.
*/
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_2_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_2_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_2_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_2) Specifies the timeout value for transactions
 * mapped to the red address queue.
 */
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_2_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_2_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_2_rqos_map_timeoutr_Msk & \
	((value) << UDDRC_PCFGQOS1_2_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_2 : (UDDRC_MP Offset: 0x204)
 * Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_2) Separation level indicating the end of region0 mapping;
 * start of region0 is 0.
 * Possible values for level1 are 0 to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_2_wqos_map_level1_Pos 0
#define UDDRC_PCFGWQOS0_2_wqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGWQOS0_2_wqos_map_level1_Pos)
#define UDDRC_PCFGWQOS0_2_wqos_map_level1(value) \
	((UDDRC_PCFGWQOS0_2_wqos_map_level1_Msk & \
	((value) << UDDRC_PCFGWQOS0_2_wqos_map_level1_Pos)))

/* (UDDRC_PCFGWQOS0_1) Separation level indicating the end of region1 mapping;
 * start of region1 is level1 + 1.
 * Possible values for level2 are (leve1+1) to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_2_wqos_map_level2_Pos 8
#define UDDRC_PCFGWQOS0_2_wqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGWQOS0_2_wqos_map_level2_Pos)
#define UDDRC_PCFGWQOS0_2_wqos_map_level2(value) \
	((UDDRC_PCFGWQOS0_2_wqos_map_level2_Msk & \
	((value) << UDDRC_PCFGWQOS0_2_wqos_map_level2_Pos)))

/* (UDDRC_PCFGWQOS0_2) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region0
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_2_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_2_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_2_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_2_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_2_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_2_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_2) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region 1
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_2_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_2_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_2_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_2_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_2_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_2_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_2 : (UDDRC_MP Offset: 0x208)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_2) Specifies the timeout value for write transactions. */
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout1_Pos 0
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout1_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_2_wqos_map_timeout1_Pos)
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout1(value) \
	((UDDRC_PCFGWQOS1_2_wqos_map_timeout1_Msk & \
	((value) << UDDRC_PCFGWQOS1_2_wqos_map_timeout1_Pos)))

/* (UDDRC_PCFGWQOS1_2) Specifies the timeout value for write transactions. */
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout2_Pos 16
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout2_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_2_wqos_map_timeout2_Pos)
#define UDDRC_PCFGWQOS1_2_wqos_map_timeout2(value) \
	((UDDRC_PCFGWQOS1_2_wqos_map_timeout2_Msk & \
	((value) << UDDRC_PCFGWQOS1_2_wqos_map_timeout2_Pos)))

/* -------- UDDRC_PCFGR_3 : (UDDRC_MP Offset: 0x21C)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_3) Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port.
 * The aging counters down-count every clock cycle where the port is requesting
 * but not granted. The higher significant 5-bits of the read aging counter
 * sets the priority of the read channel of a given port.
 * Port's priority will increase as the higher significant 5-bits of the counter
 * starts to decrease. When the aging counter becomes 0, the corresponding port
 * channel will have the highest priority level (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority (2 priority levels:
 * low priority read - LPR, high priority read - HPR)
 * on a command by command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_3_rd_port_priority_Pos 0
#define UDDRC_PCFGR_3_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_3_rd_port_priority_Pos)
#define UDDRC_PCFGR_3_rd_port_priority(value) \
	((UDDRC_PCFGR_3_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_3_rd_port_priority_Pos)))

/* (UDDRC_PCFGR_3) bypass read reordering */
#define UDDRC_PCFGR_3_read_reorder_bypass_en (0x1u << 11)

/* (UDDRC_PCFGR_3) If set to 1, enables aging function for the read
 * channel of the port.
 */
#define UDDRC_PCFGR_3_rd_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGR_3) If set to 1, enables the AXI urgent sideband signal (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr
 * signal to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGR_3_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_3) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same memory page
 * (same bank and same row). See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_3_rd_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCFGW_3 : (UDDRC_MP Offset: 0x220)
 * Port n Configuration Write Register --------
 */
/* (UDDRC_PCFGR_3) read write ordered enable */
#define UDDRC_PCFGR_3_rdwr_ordered_en (0x1u << 16)

/* (UDDRC_PCFGW_3) Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not granted.
 * The higher significant 5-bits of the write aging counter sets the initial
 * priority of the write channel of a given port. Port's priority will increase
 * as the higher significant 5-bits of the counter starts to decrease.
 * When the aging counter becomes 0, the corresponding port channel will have
 * the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_3_wr_port_priority_Pos 0
#define UDDRC_PCFGW_3_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_3_wr_port_priority_Pos)
#define UDDRC_PCFGW_3_wr_port_priority(value) \
	((UDDRC_PCFGW_3_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_3_wr_port_priority_Pos)))

/* (UDDRC_PCFGW_3) If set to 1, enables aging function for the write
 * channel of the port.
 */
#define UDDRC_PCFGW_3_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_3) If set to 1, enables the AXI urgent sideband signal (awurgent).
 * When enabled and awurgent is asserted by the master, that port becomes the
 * highest priority and co_gs_go2critical_wr signal to DDRC is asserted if
 * enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_3_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_3) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_3_wr_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCTRL_3 : (UDDRC_MP Offset: 0x2A8)
 * Port n Control Register --------
 */

/* (UDDRC_PCTRL_3) Enables AXI port n. */
#define UDDRC_PCTRL_3_port_en (0x1u << 0)
/* -------- UDDRC_PCFGQOS0_3 : (UDDRC_MP Offset: 0x2AC)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_3) Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 13 (for dual RAQ)
 * or 0 to 14 (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_3_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_3_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_3_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_3_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_3_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_3_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_3) Separation level2 indicating the end of region1 mapping;
 * start of region1 is (level1 + 1). Possible values for level2 are (level1 + 1)
 * to 14 which corresponds to arqos.
 * Region2 starts from (level2 + 1) up to 15.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_3_rqos_map_level2_Pos 8
#define UDDRC_PCFGQOS0_3_rqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGQOS0_3_rqos_map_level2_Pos)
#define UDDRC_PCFGQOS0_3_rqos_map_level2(value) \
	((UDDRC_PCFGQOS0_3_rqos_map_level2_Msk & \
((value) << UDDRC_PCFGQOS0_3_rqos_map_level2_Pos)))

/* (UDDRC_PCFGQOS0_3) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are: 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region0 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_3_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_3_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_3_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_3_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_3_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_3_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_3) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_3_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_3_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_3_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_3_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_3_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_3_rqos_map_region1_Pos)))

/* (UDDRC_PCFGQOS0_3) This bitfield indicates the traffic class of region2.
 * For dual address queue configurations, region2 maps to the red address queue.
 * Valid values are 1: VPR and 2: HPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region2 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_3_rqos_map_region2_Pos 24
#define UDDRC_PCFGQOS0_3_rqos_map_region2_Msk \
	(0x3u << UDDRC_PCFGQOS0_3_rqos_map_region2_Pos)
#define UDDRC_PCFGQOS0_3_rqos_map_region2(value) \
	((UDDRC_PCFGQOS0_3_rqos_map_region2_Msk & \
	((value) << UDDRC_PCFGQOS0_3_rqos_map_region2_Pos)))

/* -------- UDDRC_PCFGQOS1_3 : (UDDRC_MP Offset: 0x2B0)
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_3) Specifies the timeout value for transactions mapped
 * to the blue address queue.
 */
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_3_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_3_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_3_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_3) Specifies the timeout value for transactions mapped
 * to the red address queue.
 */
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_3_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_3_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_3_rqos_map_timeoutr_Msk & \
	((value) << UDDRC_PCFGQOS1_3_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_3 : (UDDRC_MP Offset: 0x2B4)
 Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_3) Separation level indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 14 which
 * corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_3_wqos_map_level1_Pos 0
#define UDDRC_PCFGWQOS0_3_wqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGWQOS0_3_wqos_map_level1_Pos)
#define UDDRC_PCFGWQOS0_3_wqos_map_level1(value) \
	((UDDRC_PCFGWQOS0_3_wqos_map_level1_Msk & \
	((value) << UDDRC_PCFGWQOS0_3_wqos_map_level1_Pos)))

/* (UDDRC_PCFGWQOS0_3) Separation level indicating the end of region1 mapping;
 * start of region1 is level1 + 1. Possible values for level2 are (leve1+1)
 * to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_3_wqos_map_level2_Pos 8
#define UDDRC_PCFGWQOS0_3_wqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGWQOS0_3_wqos_map_level2_Pos)
#define UDDRC_PCFGWQOS0_3_wqos_map_level2(value) \
	((UDDRC_PCFGWQOS0_3_wqos_map_level2_Msk & \
	((value) << UDDRC_PCFGWQOS0_3_wqos_map_level2_Pos)))

/* (UDDRC_PCFGWQOS0_3) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of
 * region0 is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_3_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_3_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_3_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_3_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_3_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_3_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_3) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of
 * region 1 is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_3_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_3_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_3_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_3_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_3_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_3_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_3 : (UDDRC_MP Offset: 0x2B8)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_3) Specifies the timeout value for write transactions.
 */
#define UDDRC_PCFGWQOS1_3_wqos_map_timeout_Pos 0
#define UDDRC_PCFGWQOS1_3_wqos_map_timeout_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_3_wqos_map_timeout_Pos)
#define UDDRC_PCFGWQOS1_3_wqos_map_timeout(value) \
	((UDDRC_PCFGWQOS1_3_wqos_map_timeout_Msk & \
	((value) << UDDRC_PCFGWQOS1_3_wqos_map_timeout_Pos)))

/* -------- UDDRC_PCFGR_4 : (UDDRC_MP Offset: 0x2CC)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_4) Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not	granted.
 * The higher significant 5-bits of the read aging counter sets the priority
 * of the read channel of a given port.	 Port's priority will increase as the
 * higher significant 5-bits of the counter starts to decrease.
 * When the aging counter becomes 0, the corresponding port channel will have
 * the highest priority level (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority  (2 priority levels:
 * low priority read - LPR, high priority read - HPR) on a command by
 * command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_4_rd_port_priority_Pos 0
#define UDDRC_PCFGR_4_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_4_rd_port_priority_Pos)
#define UDDRC_PCFGR_4_rd_port_priority(value) \
	((UDDRC_PCFGR_4_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_4_rd_port_priority_Pos)))


/* (UDDRC_PCFGR_4) bypass read reordering */
#define UDDRC_PCFGR_4_read_reorder_bypass_en (0x1u << 11)

/* (UDDRC_PCFGR_4) If set to 1, enables aging function for the read channel
 * of the port.
 */
#define UDDRC_PCFGR_4_rd_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGR_4) If set to 1, enables the AXI urgent sideband signal (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr signal
 * to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGR_4_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_4) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_4_rd_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCFGW_4 : (UDDRC_MP Offset: 0x2D0)
 * Port n Configuration Write Register --------
 */
/* (UDDRC_PCFGR_4) read write ordered enable */
#define UDDRC_PCFGR_4_rdwr_ordered_en (0x1u << 16)

/* (UDDRC_PCFGW_4) Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not granted. The higher significant 5-bits
 * of the write aging counter sets the initial priority of the write channel
 * of a given port. Port's priority will increase as the higher significant
 * 5-bits of the counter starts to decrease. When the aging counter becomes 0,
 * the corresponding port channel will have the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_4_wr_port_priority_Pos 0
#define UDDRC_PCFGW_4_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_4_wr_port_priority_Pos)
#define UDDRC_PCFGW_4_wr_port_priority(value) \
	((UDDRC_PCFGW_4_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_4_wr_port_priority_Pos)))

/* (UDDRC_PCFGW_4) If set to 1, enables aging function for the write
 * channel of the port.
 */
#define UDDRC_PCFGW_4_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_4) If set to 1, enables the AXI urgent sideband signal (awurgent).
 * When enabled and awurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_wr signal to DDRC is asserted
 * if enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_4_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_4) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued
 * to be granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_4_wr_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCTRL_4 : (UDDRC_MP Offset: 0x358)
 * Port n Control Register --------
 */
/* (UDDRC_PCTRL_4) Enables AXI port n. */
#define UDDRC_PCTRL_4_port_en (0x1u << 0)

/* -------- UDDRC_PCFGQOS0_4 : (UDDRC_MP Offset: 0x35C)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_4) Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 13 (for dual RAQ)
 * or 0 to 14 (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_4_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_4_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_4_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_4_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_4_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_4_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_4) Separation level1 indicating the end of region1 mapping;
 */
#define UDDRC_PCFGQOS0_4_rqos_map_level2_Pos 8
#define UDDRC_PCFGQOS0_4_rqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGQOS0_4_rqos_map_level2_Pos)
	#define UDDRC_PCFGQOS0_4_rqos_map_level2(value) \
	((UDDRC_PCFGQOS0_4_rqos_map_level2_Msk & \
	((value) << UDDRC_PCFGQOS0_4_rqos_map_level2_Pos)))

/* (UDDRC_PCFGQOS0_4) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are: 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of region0
 * is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_4_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_4_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_4_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_4_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_4_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_4_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_4) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_4_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_4_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_4_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_4_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_4_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_4_rqos_map_region1_Pos)))

/* (UDDRC_PCFGQOS0_4) This bitfield indicates the traffic class of region 2.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region2 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_4_rqos_map_region2_Pos 24
#define UDDRC_PCFGQOS0_4_rqos_map_region2_Msk \
	(0x3u << UDDRC_PCFGQOS0_4_rqos_map_region2_Pos)
#define UDDRC_PCFGQOS0_4_rqos_map_region2(value) \
	((UDDRC_PCFGQOS0_4_rqos_map_region2_Msk & \
	((value) << UDDRC_PCFGQOS0_4_rqos_map_region2_Pos)))

/* -------- UDDRC_PCFGQOS1_4 : (UDDRC_MP Offset: 0x360) \
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_4) Specifies the timeout value for transactions mapped
 * to the blue address queue.
 */
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_4_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_4_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_4_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_4) Specifies the timeout value for transactions mapped
 * to the red address queue.
 */
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_4_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_4_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_4_rqos_map_timeoutr_Msk & \
((value) << UDDRC_PCFGQOS1_4_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_4 : (UDDRC_MP Offset: 0x364)
 * Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_4) Separation level indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 14
 * which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_4_wqos_map_level1_Pos 0
#define UDDRC_PCFGWQOS0_4_wqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGWQOS0_4_wqos_map_level1_Pos)
#define UDDRC_PCFGWQOS0_4_wqos_map_level1(value) \
	((UDDRC_PCFGWQOS0_4_wqos_map_level1_Msk & \
((value) << UDDRC_PCFGWQOS0_4_wqos_map_level1_Pos)))

/* (UDDRC_PCFGWQOS0_4) Separation level indicating the end of region1 mapping;
 * start of region1 is level1 + 1. Possible values for level2 are (leve1+1)
 * to 14 which corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_4_wqos_map_level2_Pos 8
#define UDDRC_PCFGWQOS0_4_wqos_map_level2_Msk \
	(0xfu << UDDRC_PCFGWQOS0_4_wqos_map_level2_Pos)
#define UDDRC_PCFGWQOS0_4_wqos_map_level2(value) \
	((UDDRC_PCFGWQOS0_4_wqos_map_level2_Msk & \
	((value) << UDDRC_PCFGWQOS0_4_wqos_map_level2_Pos)))

/* (UDDRC_PCFGWQOS0_4) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region0
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_4_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_4_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_4_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_4_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_4_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_4_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_4) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of
 * region 1 is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_4_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_4_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_4_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_4_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_4_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_4_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_4 : (UDDRC_MP Offset: 0x368)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_4) Specifies the timeout value for write transactions.
 */
#define UDDRC_PCFGWQOS1_4_wqos_map_timeout_Pos 0
#define UDDRC_PCFGWQOS1_4_wqos_map_timeout_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_4_wqos_map_timeout_Pos)
#define UDDRC_PCFGWQOS1_4_wqos_map_timeout(value) \
	((UDDRC_PCFGWQOS1_4_wqos_map_timeout_Msk & \
	((value) << UDDRC_PCFGWQOS1_4_wqos_map_timeout_Pos)))

/* -------- UDDRC_PCFGR_5 : (UDDRC_MP Offset: 0x37C)
 * Port n Configuration Read Register --------
 */
/* (UDDRC_PCFGR_5) Determines the initial load value of read aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not granted. The higher significant 5-bits
 * of the read aging counter sets the priority of the read channel of
 * a given port. Port's priority will increase as the higher significant 5-bits
 * of the counter starts to decrease. When the aging counter becomes 0,
 * the corresponding port channel will have the highest priority level
 * (timeout condition - Priority0).
 * For multi-port configurations, the aging counters cannot be used to set
 * port priorities when external dynamic priority inputs (arqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * In this case, external dynamic priority input, arqos (for reads only)
 * can still be used to set the DDRC read priority (2 priority levels:
 * low priority read - LPR, high priority read - HPR)
 * on a command by command basis.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGR_5_rd_port_priority_Pos 0
#define UDDRC_PCFGR_5_rd_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGR_5_rd_port_priority_Pos)
#define UDDRC_PCFGR_5_rd_port_priority(value) \
	((UDDRC_PCFGR_5_rd_port_priority_Msk & \
	((value) << UDDRC_PCFGR_5_rd_port_priority_Pos)))

/* (UDDRC_PCFGR_5) If set to 1, enables aging function for the read channel
 * of the port.
 */
#define UDDRC_PCFGR_5_rd_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGR_5) If set to 1, enables the AXI urgent sideband signal (arurgent).
 * When enabled and arurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_lpr/co_gs_go2critical_hpr signal
 * to DDRC is asserted if enabled in PCCFG.go2critical_en register.
 * Note that arurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGR_5_rd_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGR_5) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same
 * memory page (same bank and same row).
 * See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGR_5_rd_port_pagematch_en (0x1u << 14)

/* -------- UDDRC_PCFGW_5 : (UDDRC_MP Offset: 0x380)
 * Port n Configuration Write Register --------
 */
/* (UDDRC_PCFGW_5) Determines the initial load value of write aging counters.
 * These counters will be parallel loaded after reset, or after each grant
 * to the corresponding port. The aging counters down-count every clock cycle
 * where the port is requesting but not granted. The higher significant 5-bits
 * of the write aging counter sets the initial priority of the write channel
 * of a given port. Port's priority will increase as the higher significant
 * 5-bits of the counter starts to decrease. When the aging counter becomes 0,
 * the corresponding port channel will have the highest priority level.
 * For multi-port configurations, the aging counters cannot be used to set port
 * priorities when external dynamic priority inputs (awqos) are enabled
 * (timeout is still applicable).
 * For single port configurations, the aging counters are only used when they
 * timeout (become 0) to force read-write direction switching.
 * Note: The two LSBs of this register field are tied internally to 2'b00.
 */
#define UDDRC_PCFGW_5_wr_port_priority_Pos 0
#define UDDRC_PCFGW_5_wr_port_priority_Msk \
	(0x3ffu << UDDRC_PCFGW_5_wr_port_priority_Pos)
#define UDDRC_PCFGW_5_wr_port_priority(value) \
	((UDDRC_PCFGW_5_wr_port_priority_Msk & \
	((value) << UDDRC_PCFGW_5_wr_port_priority_Pos)))

/* (UDDRC_PCFGW_5) If set to 1, enables aging function for the write channel
 * of the port.
 */
#define UDDRC_PCFGW_5_wr_port_aging_en (0x1u << 12)

/* (UDDRC_PCFGW_5) If set to 1, enables the AXI urgent sideband signal (awurgent).
 * When enabled and awurgent is asserted by the master, that port becomes
 * the highest priority and co_gs_go2critical_wr signal to DDRC is asserted
 * if enabled in PCCFG.go2critical_en register.
 * Note that awurgent signal can be asserted anytime and as long as required
 * which is independent of address handshaking
 * (it is not associated with any particular command).
 */
#define UDDRC_PCFGW_5_wr_port_urgent_en (0x1u << 13)

/* (UDDRC_PCFGW_5) If set to 1, enables the Page Match feature.
 * If enabled, once a requesting port is granted, the port is continued to be
 * granted if the following immediate commands are to the same memory page
 * (same bank and same row). See also related PCCFG.pagematch_limit register.
 */
#define UDDRC_PCFGW_5_wr_port_pagematch_en (0x1u << 14)
/* -------- UDDRC_PCTRL_5 : (UDDRC_MP Offset: 0x408)
 Port n Control Register --------
 */
/* (UDDRC_PCTRL_5) Enables AXI port n. */
#define UDDRC_PCTRL_5_port_en (0x1u << 0)

/* -------- UDDRC_PCFGQOS0_5 : (UDDRC_MP Offset: 0x40C)
 * Port n Read QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGQOS0_5) Separation level1 indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 13 (for dual RAQ)
 * or 0 to 14 (for single RAQ) which corresponds to arqos.
 * Note that for PA, arqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 * All of the map_level* registers must be set to distinct values.
 */
#define UDDRC_PCFGQOS0_5_rqos_map_level1_Pos 0
#define UDDRC_PCFGQOS0_5_rqos_map_level1_Msk \
	(0xfu << UDDRC_PCFGQOS0_5_rqos_map_level1_Pos)
#define UDDRC_PCFGQOS0_5_rqos_map_level1(value) \
	((UDDRC_PCFGQOS0_5_rqos_map_level1_Msk & \
	((value) << UDDRC_PCFGQOS0_5_rqos_map_level1_Pos)))

/* (UDDRC_PCFGQOS0_5) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region 0 maps to the blue address queue.
 * In this case, valid values are: 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region0 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_5_rqos_map_region0_Pos 16
#define UDDRC_PCFGQOS0_5_rqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGQOS0_5_rqos_map_region0_Pos)
#define UDDRC_PCFGQOS0_5_rqos_map_region0(value) \
	((UDDRC_PCFGQOS0_5_rqos_map_region0_Msk & \
	((value) << UDDRC_PCFGQOS0_5_rqos_map_region0_Pos)))

/* (UDDRC_PCFGQOS0_5) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0 : LPR, 1: VPR, 2: HPR.
 * For dual address queue configurations, region1 maps to the blue address queue.
 * In this case, valid values are 0: LPR and 1: VPR only.
 * When VPR support is disabled (UMCTL2_VPR_EN = 0) and traffic class of
 * region 1 is set to 1 (VPR), VPR traffic is aliased to LPR traffic.
 */
#define UDDRC_PCFGQOS0_5_rqos_map_region1_Pos 20
#define UDDRC_PCFGQOS0_5_rqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGQOS0_5_rqos_map_region1_Pos)
#define UDDRC_PCFGQOS0_5_rqos_map_region1(value) \
	((UDDRC_PCFGQOS0_5_rqos_map_region1_Msk & \
	((value) << UDDRC_PCFGQOS0_5_rqos_map_region1_Pos)))

/* -------- UDDRC_PCFGQOS1_5 : (UDDRC_MP Offset: 0x410)
 * Port n Read QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGQOS1_5) Specifies the timeout value for transactions mapped
 * to the blue address queue.
 */
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutb_Pos 0
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutb_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_5_rqos_map_timeoutb_Pos)
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutb(value) \
	((UDDRC_PCFGQOS1_5_rqos_map_timeoutb_Msk & \
	((value) << UDDRC_PCFGQOS1_5_rqos_map_timeoutb_Pos)))

/* (UDDRC_PCFGQOS1_5) Specifies the timeout value for transactions mapped
 * to the red address queue.
 */
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutr_Pos 16
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutr_Msk \
	(0x7ffu << UDDRC_PCFGQOS1_5_rqos_map_timeoutr_Pos)
#define UDDRC_PCFGQOS1_5_rqos_map_timeoutr(value) \
	((UDDRC_PCFGQOS1_5_rqos_map_timeoutr_Msk & \
	((value) << UDDRC_PCFGQOS1_5_rqos_map_timeoutr_Pos)))

/* -------- UDDRC_PCFGWQOS0_5 : (UDDRC_MP Offset: 0x414)
 * Port n Write QoS Configuration Register 0 --------
 */
/* (UDDRC_PCFGWQOS0_5) Separation level indicating the end of region0 mapping;
 * start of region0 is 0. Possible values for level1 are 0 to 14 which
 * corresponds to awqos.
 * Note that for PA, awqos values are used directly as port priorities,
 * where the higher the value corresponds to higher port priority.
 */
#define UDDRC_PCFGWQOS0_5_wqos_map_level_Pos 0
#define UDDRC_PCFGWQOS0_5_wqos_map_level_Msk \
	(0xfu << UDDRC_PCFGWQOS0_5_wqos_map_level_Pos)
#define UDDRC_PCFGWQOS0_5_wqos_map_level(value) \
	((UDDRC_PCFGWQOS0_5_wqos_map_level_Msk & \
	((value) << UDDRC_PCFGWQOS0_5_wqos_map_level_Pos)))

/* (UDDRC_PCFGWQOS0_5) This bitfield indicates the traffic class of region 0.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of region0
 * is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_5_wqos_map_region0_Pos 16
#define UDDRC_PCFGWQOS0_5_wqos_map_region0_Msk \
	(0x3u << UDDRC_PCFGWQOS0_5_wqos_map_region0_Pos)
#define UDDRC_PCFGWQOS0_5_wqos_map_region0(value) \
	((UDDRC_PCFGWQOS0_5_wqos_map_region0_Msk & \
	((value) << UDDRC_PCFGWQOS0_5_wqos_map_region0_Pos)))

/* (UDDRC_PCFGWQOS0_5) This bitfield indicates the traffic class of region 1.
 * Valid values are: 0: NPW, 1: VPW.
 * When VPW support is disabled (UMCTL2_VPW_EN = 0) and traffic class of
 * region 1 is set to 1 (VPW), VPW traffic is aliased to NPW traffic.
 */
#define UDDRC_PCFGWQOS0_5_wqos_map_region1_Pos 20
#define UDDRC_PCFGWQOS0_5_wqos_map_region1_Msk \
	(0x3u << UDDRC_PCFGWQOS0_5_wqos_map_region1_Pos)
#define UDDRC_PCFGWQOS0_5_wqos_map_region1(value) \
	((UDDRC_PCFGWQOS0_5_wqos_map_region1_Msk & \
	((value) << UDDRC_PCFGWQOS0_5_wqos_map_region1_Pos)))

/* -------- UDDRC_PCFGWQOS1_5 : (UDDRC_MP Offset: 0x418)
 * Port n Write QoS Configuration Register 1 --------
 */
/* (UDDRC_PCFGWQOS1_5) Specifies the timeout value for write transactions.*/
#define UDDRC_PCFGWQOS1_5_wqos_map_timeout_Pos 0
#define UDDRC_PCFGWQOS1_5_wqos_map_timeout_Msk \
	(0x7ffu << UDDRC_PCFGWQOS1_5_wqos_map_timeout_Pos)
#define UDDRC_PCFGWQOS1_5_wqos_map_timeout(value) \
	((UDDRC_PCFGWQOS1_5_wqos_map_timeout_Msk & \
	((value) << UDDRC_PCFGWQOS1_5_wqos_map_timeout_Pos)))

/* } */
/* End of UMCTL2 MP register helpers */

/*
 * Synopsys UMCTL2 register map for main configuration.
 */

struct uddrc_regs {
	/* (Uddrc_regs Offset: 0x0) Master Register0 */
	__IO unsigned int UDDRC_MSTR;
	/* (Uddrc_regs Offset: 0x4) Operating Mode Status Register */
	__I  unsigned int UDDRC_STAT;
	__I  unsigned int Reserved1[2];
	/* (Uddrc_regs Offset: 0x10) Mode Register Read/Write Control Register 0. */
	__IO unsigned int UDDRC_MRCTRL0;
	/* (Uddrc_regs Offset: 0x14) Mode Register Read/Write Control Register 1 */
	__IO unsigned int UDDRC_MRCTRL1;
	/* (Uddrc_regs Offset: 0x18) Mode Register Read/Write Status Register */
	__I  unsigned int UDDRC_MRSTAT;
	/* (Uddrc_regs Offset: 0x1C Mode Register Read/Write Control Register 2 */
	__IO unsigned int UDDRC_MRCTRL2;
	/* (Uddrc_regs Offset: 0x20) Temperature Derate Enable Register */
	__IO unsigned int UDDRC_DERATEEN;
	/* (Uddrc_regs Offset: 0x24) Temperature Derate Interval Register */
	__IO unsigned int UDDRC_DERATEINT;
	__I  unsigned int Reserved3[2];
	/* (Uddrc_regs Offset: 0x30) Low Power Control Register */
	__IO unsigned int UDDRC_PWRCTL;
	/* (Uddrc_regs Offset: 0x34) Low Power Timing Register */
	__IO unsigned int UDDRC_PWRTMG;
	/* (Uddrc_regs Offset: 0x38) Hardware Low Power Control Register */
	__IO unsigned int UDDRC_HWLPCTL;
	__I  unsigned int Reserved4[5];
	/* (Uddrc_regs Offset: 0x50) Refresh Control Register 0 */
	__IO unsigned int UDDRC_RFSHCTL0;
	__I  unsigned int Reserved5[3];
	/* (Uddrc_regs Offset: 0x60) Refresh Control Register 3 */
	__IO unsigned int UDDRC_RFSHCTL3;
	/* (Uddrc_regs Offset: 0x64) Refresh Timing Register */
	__IO unsigned int UDDRC_RFSHTMG;
	__I  unsigned int Reserved6[22];
	/* (Uddrc_regs Offset: 0xC0) CRC Parity Control Register0 */
	__IO unsigned int UDDRC_CRCPARCTL0;
	/* (Uddrc_regs Offset: 0xC4) CRC Parity Control Register1 */
	__IO unsigned int UDDRC_CRCPARCTL1;
	/* (Uddrc_regs Offset: 0xC8) CRC Parity Control Register2 */
	__IO unsigned int UDDRC_CRCPARCTL2;
	/* (Uddrc_regs Offset: 0xCC) CRC Parity Status Register */
	__I  unsigned int UDDRC_CRCPARSTAT;
	/* (Uddrc_regs Offset: 0xD0) SDRAM Initialization Register 0 */
	__IO unsigned int UDDRC_INIT0;
	/* (Uddrc_regs Offset: 0xD4) SDRAM Initialization Register 1 */
	__IO unsigned int UDDRC_INIT1;
	/* (Uddrc_regs Offset: 0xD8) SDRAM Initialization Register 2 */
	__IO unsigned int UDDRC_INIT2;
	/* (Uddrc_regs Offset: 0xDC) SDRAM Initialization Register 3 */
	__IO unsigned int UDDRC_INIT3;
	/* (Uddrc_regs Offset: 0xE0) SDRAM Initialization Register 4 */
	__IO unsigned int UDDRC_INIT4;
	/* (Uddrc_regs Offset: 0xE4) SDRAM Initialization Register 5 */
	__IO unsigned int UDDRC_INIT5;
	__I  unsigned int Reserved8[2];
	/* (Uddrc_regs Offset: 0xF0) DIMM Control Register */
	__IO unsigned int UDDRC_DIMMCTL;
	__I  unsigned int Reserved9[3];
	/* (Uddrc_regs Offset: 0x100) SDRAM Timing Register 0 */
	__IO unsigned int UDDRC_DRAMTMG0;
	/* (Uddrc_regs Offset: 0x104) SDRAM Timing Register 1 */
	__IO unsigned int UDDRC_DRAMTMG1;
	/* (Uddrc_regs Offset: 0x108) SDRAM Timing Register 2 */
	__IO unsigned int UDDRC_DRAMTMG2;
	/* (Uddrc_regs Offset: 0x10C) SDRAM Timing Register 3 */
	__IO unsigned int UDDRC_DRAMTMG3;
	/* (Uddrc_regs Offset: 0x110) SDRAM Timing Register 4 */
	__IO unsigned int UDDRC_DRAMTMG4;
	/* (Uddrc_regs Offset: 0x114) SDRAM Timing Register 5 */
	__IO unsigned int UDDRC_DRAMTMG5;
	/* (Uddrc_regs Offset: 0x118) SDRAM Timing Register 6 */
	__IO unsigned int UDDRC_DRAMTMG6;
	/* (Uddrc_regs Offset: 0x11C) SDRAM Timing Register 7 */
	__IO unsigned int UDDRC_DRAMTMG7;
	/* (Uddrc_regs Offset: 0x120) SDRAM Timing Register 8 */
	__IO unsigned int UDDRC_DRAMTMG8;
	__I  unsigned int Reserved10[5];
	/* (Uddrc_regs Offset: 0x138) SDRAM Timing Register 14 */
	__IO unsigned int UDDRC_DRAMTMG14;
	/* (Uddrc_regs Offset: 0x13C) SDRAM Timing Register 15 */
	__IO unsigned int UDDRC_DRAMTMG15;
	__I  unsigned int Reserved11[16];
	/* (Uddrc_regs Offset: 0x180) ZQ Control Register 0 */
	__IO unsigned int UDDRC_ZQCTL0;
	/* (Uddrc_regs Offset: 0x184) ZQ Control Register 1 */
	__IO unsigned int UDDRC_ZQCTL1;
	/* (Uddrc_regs Offset: 0x188) ZQ Control Register 2 */
	__IO unsigned int UDDRC_ZQCTL2;
	/* (Uddrc_regs Offset: 0x18C) ZQ Status Register */
	__I  unsigned int UDDRC_ZQSTAT;
	/* (Uddrc_regs Offset: 0x190) DFI Timing Register 0 */
	__IO unsigned int UDDRC_DFITMG0;
	/* (Uddrc_regs Offset: 0x194) DFI Timing Register 1 */
	__IO unsigned int UDDRC_DFITMG1;
	/* (Uddrc_regs Offset: 0x198) DFI Low Power Configuration Register 0 */
	__IO unsigned int UDDRC_DFILPCFG0;
	/* (Uddrc_regs Offset: 0x19C) DFI Low Power Configuration Register 1 */
	__IO unsigned int UDDRC_DFILPCFG1;
	/* (Uddrc_regs Offset: 0x1A0) DFI Update Register 0 */
	__IO unsigned int UDDRC_DFIUPD0;
	/* (Uddrc_regs Offset: 0x1A4) DFI Update Register 1 */
	__IO unsigned int UDDRC_DFIUPD1;
	/* (Uddrc_regs Offset: 0x1A8) DFI Update Register 2 */
	__IO unsigned int UDDRC_DFIUPD2;
	__I  unsigned int Reserved13[1];
	/* (Uddrc_regs Offset: 0x1B0) DFI Miscellaneous Control Register */
	__IO unsigned int UDDRC_DFIMISC;
	/* (Uddrc_regs Offset: 0x1B4) DFI Timing Register 1 */
	__IO unsigned int UDDRC_DFITMG2;
	__I  unsigned int Reserved14[1];
	/* (Uddrc_regs Offset: 0x1BC) DFI Status Register */
	__I  unsigned int UDDRC_DFISTAT;
	__I  unsigned int Reserved15[1];
	/* (Uddrc_regs Offset: 0x1C4) DFI PHY Master */
	__IO unsigned int UDDRC_DFIPHYMSTR;
	__I  unsigned int Reserved16[15];
	/* (Uddrc_regs Offset: 0x204) Address Map Register 1 */
	__IO unsigned int UDDRC_ADDRMAP1;
	/* (Uddrc_regs Offset: 0x208) Address Map Register 2 */
	__IO unsigned int UDDRC_ADDRMAP2;
	/* (Uddrc_regs Offset: 0x20C) Address Map Register 3 */
	__IO unsigned int UDDRC_ADDRMAP3;
	/* (Uddrc_regs Offset: 0x210) Address Map Register 4 */
	__IO unsigned int UDDRC_ADDRMAP4;
	/* (Uddrc_regs Offset: 0x214) Address Map Register 5 */
	__IO unsigned int UDDRC_ADDRMAP5;
	/* (Uddrc_regs Offset: 0x218) Address Map Register 6 */
	__IO unsigned int UDDRC_ADDRMAP6;
	__I  unsigned int Reserved17[2];
	/* (Uddrc_regs Offset: 0x224) Address Map Register 9 */
	__IO unsigned int UDDRC_ADDRMAP9;
	/* (Uddrc_regs Offset: 0x228) Address Map Register 10 */
	__IO unsigned int UDDRC_ADDRMAP10;
	/* (Uddrc_regs Offset: 0x22C) Address Map Register 11 */
	__IO unsigned int UDDRC_ADDRMAP11;
	__I  unsigned int Reserved18[4];
	/* (Uddrc_regs Offset: 0x240) ODT Configuration Register */
	__IO unsigned int UDDRC_ODTCFG;
	/* (Uddrc_regs Offset: 0x244) ODT/Rank Map Register */
	__IO unsigned int UDDRC_ODTMAP;
	__I  unsigned int Reserved19[2];
	/* (Uddrc_regs Offset: 0x250) Scheduler Control Register */
	__IO unsigned int UDDRC_SCHED;
	/* (Uddrc_regs Offset: 0x254) Scheduler Control Register 1 */
	__IO unsigned int UDDRC_SCHED1;
	__I  unsigned int Reserved20[1];
	/* (Uddrc_regs Offset: 0x25C) High Priority Read CAM Register 1 */
	__IO unsigned int UDDRC_PERFHPR1;
	__I  unsigned int Reserved21[1];
	/* (Uddrc_regs Offset: 0x264) Low Priority Read CAM Register 1 */
	__IO unsigned int UDDRC_PERFLPR1;
	__I  unsigned int Reserved22[1];
	/* (Uddrc_regs Offset: 0x26C) Write CAM Register 1 */
	__IO unsigned int UDDRC_PERFWR1;
	__I  unsigned int Reserved23[36];
	/* (Uddrc_regs Offset: 0x300) Debug Register 0 */
	__IO unsigned int UDDRC_DBG0;
	/* (Uddrc_regs Offset: 0x304) Debug Register 1 */
	__IO unsigned int UDDRC_DBG1;
	/* (Uddrc_regs Offset: 0x308) CAM Debug Register */
	__I  unsigned int UDDRC_DBGCAM;
	/* (Uddrc_regs Offset: 0x30C) Command Debug Register */
	__IO unsigned int UDDRC_DBGCMD;
	/* (Uddrc_regs Offset: 0x310) Status Debug Register */
	__I  unsigned int UDDRC_DBGSTAT;
	__I  unsigned int Reserved24[3];
	/* (Uddrc_regs Offset: 0x320) Software Register Programming Control Enable */
	__IO unsigned int UDDRC_SWCTL;
	/* (Uddrc_regs Offset: 0x324) Software Register Programming Control Status */
	__I  unsigned int UDDRC_SWSTAT;
	__I  unsigned int Reserved25[17];
	/* (Uddrc_regs Offset: 0x36C) AXI Poison Configuration Register.
	 * Common for all AXI ports */
	__IO unsigned int UDDRC_POISONCFG;
	/* (Uddrc_regs Offset: 0x370) AXI Poison Status Register */
	__I  unsigned int UDDRC_POISONSTAT;
};


/* UMCTL2 main register helpers */
/* { */

/* -------- UDDRC_MSTR : (UDDRC_REGS Offset: 0x0) Master Register0 -------- */
/* (UDDRC_MSTR) Select DDR3 SDRAM
 * - 1 - DDR3 SDRAM device in use
 * - 0 - non-DDR3 SDRAM device in use
 * Only present in designs that support DDR3.
 */
#define UDDRC_MSTR_ddr3 (0x1u << 0)

/* (UDDRC_MSTR) Select mDDR SDRAM
 * - 1 - Mobile/LPDDR SDRAM device in use
 * - 0 - non-mobile SDRAM device in use
 * This is only present for designs supporting mDDR.
 */
#define UDDRC_MSTR_mobile (0x1u << 1)

/* (UDDRC_MSTR) Select LPDDR2 SDRAM
 * - 1 - LPDDR2 SDRAM device in use.
 * - 0 - non-LPDDR2 device in use
 * Present only in designs configured to support LPDDR2.
 */
#define UDDRC_MSTR_lpddr2 (0x1u << 2)

/* (UDDRC_MSTR) Select LPDDR3 SDRAM
 * - 1 - LPDDR3 SDRAM device in use.
 * - 0 - non-LPDDR3 device in use
 * Present only in designs configured to support LPDDR3.
 */
#define UDDRC_MSTR_lpddr3 (0x1u << 3)

/* (UDDRC_MSTR) When set, enable interleaved burst mode */
#define UDDRC_MSTR_burst_mode (0x1u << 8)

/* (UDDRC_MSTR) When set, enable burst-chop (BC4 or 8 on-the-fly) in DDR3/DDR4.
 * Burst Chop for Reads is exercised only in HIF configurations
 * (UMCTL2_INCL_ARB not set) and if in full bus width mode
 * (MSTR.data_bus_width = 00) and if MEMC_BURST_LENGTH=8 or 16.
 * Burst Chop for Writes is exercised only if Partial Writes enabled
 * (UMCTL2_PARTIAL_WR=1) and if CRC is disabled (CRCPARCTL1.crc_enable = 0).
 * If DDR4 CRC/parity retry is enabled (CRCPARCTL1.crc_parity_retry_enable = 1),
 * burst chop is not supported, and this bit must be set to '0'.
 * BC4 (fixed) mode is not supported.
 */
#define UDDRC_MSTR_burstchop (0x1u << 9)

/* (UDDRC_MSTR) If 1, then uMCTL2 uses 2T timing. Otherwise, uses 1T timing.
 * In 2T timing, all command signals (except chip select) are held for
 * 2 clocks on the SDRAM bus.
 * Chip select is asserted on the second cycle of the command
 * Note: 2T timing is not supported in LPDDR2/LPDDR3/LPDDR4 mode
 * Note: 2T timing is not supported if the configuration parameter
 * MEMC_CMD_RTN2IDLE is set
 * Note: 2T timing is not supported in DDR4 geardown mode.
 * Note: 2T timing is not supported in Shared-AC dual channel mode and
 * the register value is don't care.
 */
#define UDDRC_MSTR_en_2t_timing_mode (0x1u << 10)

/* (UDDRC_MSTR) geardown_mode */
#define UDDRC_MSTR_geardown_mode (0x1u << 11)

/* (UDDRC_MSTR) Selects proportion of DQ bus width that is used by the SDRAM
 * - 00 - Full DQ bus width to SDRAM
 * - 01 - Half DQ bus width to SDRAM
 * - 10 - Quarter DQ bus width to SDRAM
 * - 11 - Reserved.
 * Note that half bus width mode is only supported when the SDRAM bus width
 * is a multiple of 16, and quarter bus width mode is only supported when
 * the SDRAM bus width is a multiple of 32 and the configuration parameter
 * MEMC_QBUS_SUPPORT is set.
 * Bus width refers to DQ bus width (excluding any ECC width).
 */
#define UDDRC_MSTR_data_bus_width_Pos 12
#define UDDRC_MSTR_data_bus_width_Msk (0x3u << UDDRC_MSTR_data_bus_width_Pos)
#define UDDRC_MSTR_data_bus_width(value) \
	((UDDRC_MSTR_data_bus_width_Msk & \
	((value) << UDDRC_MSTR_data_bus_width_Pos)))

/* (UDDRC_MSTR) Set to 1 when the uMCTL2 and DRAM has to be put in
 * DLL-off mode for low frequency operation.
 * Set to 0 to put uMCTL2 and DRAM in DLL-on mode for normal frequency operation.
 * If DDR4 CRC/parity retry is enabled (CRCPARCTL1.crc_parity_retry_enable = 1),
 * dll_off_mode is not supported, and this bit must be set to '0'.
 */
#define UDDRC_MSTR_dll_off_mode (0x1u << 15)

/* (UDDRC_MSTR) SDRAM burst length used:
 * - 0001 - Burst length of 2 (only supported for mDDR)
 * - 0010 - Burst length of 4
 * - 0100 - Burst length of 8
 * - 1000 - Burst length of 16 (only supported for mDDR, LPDDR2, and LPDDR4)
 * All other values are reserved.
 * This controls the burst size used to access the SDRAM.
 * This must match the burst length mode register setting in the SDRAM.
 * (For BC4/8 on-the-fly mode of DDR3 and DDR4, set this field to 0x0100)
 * Burst length of 2 is not supported with AXI ports when MEMC_BURST_LENGTH is 8.
 * Burst length of 2 is only supported when the controller is operating in 1:1
 * frequency mode.
 * For DDR3, DDR4 and LPDDR3, this must be set to 0x0100 (BL8).
 * For LPDDR4, this must be set to 0x1000 (BL16).
 */
#define UDDRC_MSTR_burst_rdwr_BL2 (1) /* 0001 */
#define UDDRC_MSTR_burst_rdwr_BL4 (2) /* 0010 */
#define UDDRC_MSTR_burst_rdwr_BL8 (4) /* 0100 */
#define UDDRC_MSTR_burst_rdwr_BL16 (8) /* 1000 */
#define UDDRC_MSTR_burst_rdwr_Pos 16
#define UDDRC_MSTR_burst_rdwr_Msk (0xfu << UDDRC_MSTR_burst_rdwr_Pos)
#define UDDRC_MSTR_burst_rdwr(value) \
	((UDDRC_MSTR_burst_rdwr_Msk & ((value) << UDDRC_MSTR_burst_rdwr_Pos)))

/* (UDDRC_MSTR) Choose which registers are used.
 * - 0 - Original registers
 * - 1 - When UMCTL2_FREQUENCY_NUM=2: FREQ1 registers
 *      When UMCTL2_FREQUENCY_NUM>2: Choosen by MSTR2.target_frequency register
 */
#define UDDRC_MSTR_frequency_ratio (0x1u << 22)
#define UDDRC_MSTR_frequency_mode (0x1u << 29)

/* -------- UDDRC_STAT : (UDDRC_REGS Offset: 0x4)
 * Operating Mode Status Register --------
 */
/* (UDDRC_STAT) Operating mode.
 * This is 3-bits wide in configurations with mDDR/LPDDR2/LPDDR3/LPDDR4/DDR4
 * support and 2-bits in all other configurations.
non-mDDR/LPDDR2/LPDDR3/LPDDR4 and non-DDR4 designs:
 - 00 - Init
 - 01 - Normal
 - 10 - Power-down
 - 11 - Self refresh
mDDR/LPDDR2/LPDDR3 or DDR4 designs:
 - 000 - Init
 - 001 - Normal
 - 010 - Power-down
 - 011 - Self refresh
 - 1XX - Deep power-down / Maximum Power Saving Mode
LPDDR4 designs:
 - 000 - Init
 - 001 - Normal
 - 010 - Power-down
 - 011 - Self refresh / Self refresh power-down
 */
#define UDDRC_STAT_operating_mode_Pos 0
#define UDDRC_STAT_operating_mode_Msk (0x3u << UDDRC_STAT_operating_mode_Pos)
#define UDDRC_STAT_operating_mode_Normal (0x1u)

/* (UDDRC_STAT) Flags if Self Refresh (except LPDDR4) or SR-Powerdown (LPDDR4)
 * is entered and if it was under Automatic Self Refresh control only or not.
 - 00 - SDRAM is not in Self Refresh (except LPDDR4) or SR-Powerdown (LPDDR4).
If retry is enabled by CRCPARCTRL1.crc_parity_retry_enable, this also indicates
SRE command is still in parity error window or retry is in-progress.
 - 11 - SDRAM is in Self Refresh (except LPDDR4) or SR-Powerdown (LPDDR4),
which was caused by Automatic Self Refresh only. If retry is enabled, this
guarantees SRE command is executed correctly without parity error.
 - 10 - SDRAM is in Self Refresh (except LPDDR4) or SR-Powerdown (LPDDR4),
which was not caused solely under Automatic Self Refresh control.
It could have been caused by Hardware Low Power Interface and/or Software
(PWRCTL.selfref_sw). If retry is enabled, this guarantees SRE command is
executed correctly without parity error.
 - 01 - SDRAM is in Self Refresh, which was caused by PHY Master Request.
 */
#define UDDRC_STAT_selfref_type_Pos 4
#define UDDRC_STAT_selfref_type_Msk (0x3u << UDDRC_STAT_selfref_type_Pos)

/* (UDDRC_STAT) Self refresh with CAMs not empty. Set to 1 when Self Refresh
 * is entered but CAMs are not drained. Cleared after exiting Self Refresh.
 */
#define UDDRC_STAT_selfref_cam_not_empty (0x1u << 12)

/* -------- UDDRC_MRCTRL0 : (UDDRC_REGS Offset: 0x10)
 * Mode Register Read/Write Control Register 0. --------
 */
/* (UDDRC_MRCTRL0) Indicates whether the mode register operation is
 *read or write. Only used for LPDDR2/LPDDR3/LPDDR4/DDR4.
 - 0 - Write
 - 1 - Read
 */
#define UDDRC_MRCTRL0_mr_type (0x1u << 0)

/* (UDDRC_MRCTRL0) Controls which rank is accessed by MRCTRL0.mr_wr.
 * Normally, it is desired to access all ranks, so all bits should be set to 1.
 * However, for multi-rank UDIMMs/RDIMMs/LRDIMMs which implement address
 * mirroring, it may be necessary to access ranks individually.
 * Examples (assume uMCTL2 is configured for 4 ranks):
    - 0x1 - select rank 0 only
    - 0x2 - select rank 1 only
    - 0x5 - select ranks 0 and 2
    - 0xA - select ranks 1 and 3
    - 0xF - select ranks 0, 1, 2 and 3
 */
#define UDDRC_MRCTRL0_mr_rank_Pos 4
#define UDDRC_MRCTRL0_mr_rank_Msk (0xfu << UDDRC_MRCTRL0_mr_rank_Pos)
#define UDDRC_MRCTRL0_mr_rank(value) \
	((UDDRC_MRCTRL0_mr_rank_Msk & \
	((value) << UDDRC_MRCTRL0_mr_rank_Pos)))

/* (UDDRC_MRCTRL0) Address of the mode register that is to be written to.
 - 0000 - MR0
 - 0001 - MR1
 - 0010 - MR2
 - 0011 - MR3
 - 0100 - MR4
 - 0101 - MR5
 - 0110 - MR6
 - 0111 - MR7
 * Don't Care for LPDDR2/LPDDR3/LPDDR4 (see MRCTRL1.mr_data for mode register
 * addressing in LPDDR2/LPDDR3/LPDDR4)
 * This signal is also used for writing to control words of the register chip
 * on RDIMMs/LRDIMMs. In that case, it corresponds to the bank address bits
 * sent to the RDIMM/LRDIMM
 * In case of DDR4, the bit[3:2] corresponds to the bank group bits.
 * Therefore, the bit[3] as well as the bit[2:0] must be set to an appropriate
 * value which is considered both the Address Mirroring of UDIMMs/RDIMMs/LRDIMMs
 * and the Output Inversion of RDIMMs/LRDIMMs.
 */
#define UDDRC_MRCTRL0_mr_addr_Pos 12
#define UDDRC_MRCTRL0_mr_addr_Msk (0xfu << UDDRC_MRCTRL0_mr_addr_Pos)
#define UDDRC_MRCTRL0_mr_addr(value) \
	((UDDRC_MRCTRL0_mr_addr_Msk & \
	((value) << UDDRC_MRCTRL0_mr_addr_Pos)))

/* (UDDRC_MRCTRL0) Setting this register bit to 1 triggers a mode register
 * read or write operation. When the MR operation is complete, the uMCTL2
 * automatically clears this bit. The other register fields of this register
 * must be written in a separate APB transaction, before setting this mr_wr bit.
 * It is recommended NOT to set this signal if in Init, Deep power-down
 * or MPSM operating modes.
 */
#define UDDRC_MRCTRL0_mr_wr (0x1u << 31)

/* -------- UDDRC_MRCTRL1 : (UDDRC_REGS Offset: 0x14)
 * Mode Register Read/Write Control Register 1 --------
 */
/* (UDDRC_MRCTRL1) Mode register write data for all
 * non-LPDDR2/non-LPDDR3/non-LPDDR4 modes.
 * For LPDDR2/LPDDR3/LPDDR4, MRCTRL1[15:0] are interpreted as [15:8] MR Address
 * [7:0] MR data for writes, don't care for reads.
 * This is 18-bits wide in configurations with DDR4 support and
 * 16-bits in all other configurations.
 */
#define UDDRC_MRCTRL1_mr_data_Pos 0
#define UDDRC_MRCTRL1_mr_data_Msk (0x3ffffu << UDDRC_MRCTRL1_mr_data_Pos)
#define UDDRC_MRCTRL1_mr_data(value) \
	((UDDRC_MRCTRL1_mr_data_Msk & \
	((value) << UDDRC_MRCTRL1_mr_data_Pos)))

/* -------- UDDRC_MRSTAT : (UDDRC_REGS Offset: 0x18)
 * Mode Register Read/Write Status Register --------
 */
/* (UDDRC_MRSTAT) The SoC core may initiate a MR write operation only if
 * this signal is low. This signal goes high in the clock after the uMCTL2
 * accepts the MRW/MRR request. It goes low when the MRW/MRR command is issued
 * to the SDRAM. It is recommended not to perform MRW/MRR commands when
 * 'MRSTAT.mr_wr_busy' is high.
 - 0 - Indicates that the SoC core can initiate a mode register write operation
 - 1 - Indicates that mode register write operation is in progress
 */

/* -------- UDDRC_DERATEEN : (UDDRC_REGS Offset: 0x20)
 * Temperature Derate Enable Register --------
 */
#define UDDRC_MRSTAT_mr_wr_busy (0x1u << 0)

/* (UDDRC_DERATEEN) Enables derating
 - 0 - Timing parameter derating is disabled
 - 1 - Timing parameter derating is enabled using MR4 read value.
 * Present only in designs configured to support LPDDR2/LPDDR3/LPDDR4
 * This field must be set to '0' for non-LPDDR2/LPDDR3/LPDDR4 mode.
 */
#define UDDRC_DERATEEN_derate_enable (0x1u << 0)

/* (UDDRC_DERATEEN) Derate value
 - 0 - Derating uses +1.
 - 1 - Derating uses +2.
 * Present only in designs configured to support LPDDR2/LPDDR3/LPDDR4
 * Set to 0 for all LPDDR2 speed grades as derating value of +1.875 ns is less
 * than a core_ddrc_core_clk period.
 * For LPDDR3/4, if the period of core_ddrc_core_clk is less than 1.875ns,
 * this register field should be set to 1; otherwise it should be set to 0.
 */
#define UDDRC_DERATEEN_derate_value_Pos 1
#define UDDRC_DERATEEN_derate_value_Msk \
	(0x3u << UDDRC_DERATEEN_derate_value_Pos)
#define UDDRC_DERATEEN_derate_value(value) \
	((UDDRC_DERATEEN_derate_value_Msk & \
	((value) << UDDRC_DERATEEN_derate_value_Pos)))

/* (UDDRC_DERATEEN) Derate byte
 * Present only in designs configured to support LPDDR2/LPDDR3/LPDDR4
 * Indicates which byte of the MRR data is used for derating.
 * The maximum valid value depends on MEMC_DRAM_TOTAL_DATA_WIDTH.
 */
#define UDDRC_DERATEEN_derate_byte_Pos 4
#define UDDRC_DERATEEN_derate_byte_Msk \
	(0xfu << UDDRC_DERATEEN_derate_byte_Pos)
#define UDDRC_DERATEEN_derate_byte(value) \
	((UDDRC_DERATEEN_derate_byte_Msk & \
	((value) << UDDRC_DERATEEN_derate_byte_Pos)))

/* -------- UDDRC_DERATEINT : (UDDRC_REGS Offset: 0x24)
 * Temperature Derate Interval Register --------
 */
/* (UDDRC_DERATEINT) Interval between two MR4 reads, used to derate
 * the timing parameters.
 * Present only in designs configured to support LPDDR2/LPDDR3/LPDDR4.
 * This register must not be set to zero.
 * Unit: DFI clock cycle.
 */
#define UDDRC_DERATEINT_mr4_read_interval_Pos 0
#define UDDRC_DERATEINT_mr4_read_interval_Msk \
	(0xffffffffu << UDDRC_DERATEINT_mr4_read_interval_Pos)
#define UDDRC_DERATEINT_mr4_read_interval(value) \
	((UDDRC_DERATEINT_mr4_read_interval_Msk & \
	((value) << UDDRC_DERATEINT_mr4_read_interval_Pos)))

/* -------- UDDRC_PWRCTL : (UDDRC_REGS Offset: 0x30)
 * Low Power Control Register --------
 */
/* (UDDRC_PWRCTL) If true then the uMCTL2 puts the SDRAM into Self Refresh
 * after a programmable number of cycles
 * "maximum idle clocks before Self Refresh (PWRTMG.selfref_to_x32)".
 * This register bit may be re-programmed during the course of normal operation.
 */
#define UDDRC_PWRCTL_selfref_en (0x1u << 0)

/* (UDDRC_PWRCTL) If true then the uMCTL2 goes into power-down after
 * a programmable number of cycles "maximum idle clocks before power down"
 * (PWRTMG.powerdown_to_x32).
 * This register bit may be re-programmed during the course of normal operation.
 */
#define UDDRC_PWRCTL_powerdown_en (0x1u << 1)

/* (UDDRC_PWRCTL) When this is 1, uMCTL2 puts the SDRAM into deep power-down
 * mode when the transaction store is empty.
 * This register must be reset to '0' to bring uMCTL2 out of deep power-down
 * mode. Controller performs automatic SDRAM initialization on deep power-down
 * exit.
 * Present only in designs configured to support mDDR or LPDDR2 or LPDDR3.
 * For non-mDDR/non-LPDDR2/non-LPDDR3, this register should not be set to 1.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PWRCTL_deeppowerdown_en (0x1u << 2)

/* (UDDRC_PWRCTL) Enable the assertion of dfi_dram_clk_disable whenever
 * a clock is not required by the SDRAM.
 * If set to 0, dfi_dram_clk_disable is never asserted.
 * Assertion of dfi_dram_clk_disable is as follows:
 * In DDR2/DDR3, can only be asserted in Self Refresh.
 * In DDR4, can be asserted in following:
  - in Self Refresh.
  - in Maximum Power Saving Mode
 * In mDDR/LPDDR2/LPDDR3, can be asserted in following:
  - in Self Refresh
  - in Power Down
  - in Deep Power Down
  - during Normal operation (Clock Stop)
 * In LPDDR4, can be asserted in following:
  - in Self Refresh Power Down
  - in Power Down
  - during Normal operation (Clock Stop)
 */
#define UDDRC_PWRCTL_en_dfi_dram_clk_disable (0x1u << 3)

/* (UDDRC_PWRCTL) A value of 1 to this register causes system to move to
 * Self Refresh state immediately, as long as it is not in INIT or DPD/MPSM
 * operating_mode. This is referred to as Software Entry/Exit to Self Refresh.
  - 1 - Software Entry to Self Refresh
  - 0 - Software Exit from Self Refresh
 */
#define UDDRC_PWRCTL_selfref_sw (0x1u << 5)

/* (UDDRC_PWRCTL) Indicates whether skipping CAM draining is allowed when
 * entering Self-Refresh.
 * This register field cannot be modified while PWRCTL.selfref_sw==1.
  - 0 - CAMs must be empty before entering SR
  - 1 - CAMs are not emptied before entering SR
 */
#define UDDRC_PWRCTL_dis_cam_drain_selfref (0x1u << 7)

/* -------- UDDRC_PWRTMG : (UDDRC_REGS Offset: 0x34)
 * Low Power Timing Register --------
 */
/* (UDDRC_PWRTMG) After this many clocks of the DDRC command channel being idle
 * the uMCTL2 automatically puts the SDRAM into power-down.
 * The DDRC command channel is considered idle when there are no HIF commands
 * outstanding. This must be enabled in the PWRCTL.powerdown_en.
 * Unit: Multiples of 32 DFI clocks
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PWRTMG_powerdown_to_x32_Pos 0
#define UDDRC_PWRTMG_powerdown_to_x32_Msk \
	(0x1fu << UDDRC_PWRTMG_powerdown_to_x32_Pos)
#define UDDRC_PWRTMG_powerdown_to_x32(value) \
	((UDDRC_PWRTMG_powerdown_to_x32_Msk & \
	((value) << UDDRC_PWRTMG_powerdown_to_x32_Pos)))

/* (UDDRC_PWRTMG) Minimum deep power-down time.
 * For mDDR, value from the JEDEC specification is 0 as mDDR exits from
 * deep power-down mode immediately after PWRCTL.deeppowerdown_en is de-asserted.
 * For LPDDR2/LPDDR3, value from the JEDEC specification is 500us.
 * Unit: Multiples of 4096 DFI clocks.
 * Present only in designs configured to support mDDR, LPDDR2 or LPDDR3.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PWRTMG_t_dpd_x4096_Pos 8
#define UDDRC_PWRTMG_t_dpd_x4096_Msk (0xffu << UDDRC_PWRTMG_t_dpd_x4096_Pos)
#define UDDRC_PWRTMG_t_dpd_x4096(value) \
	((UDDRC_PWRTMG_t_dpd_x4096_Msk & \
	((value) << UDDRC_PWRTMG_t_dpd_x4096_Pos)))

/* (UDDRC_PWRTMG) After this many clocks of the DDRC command channel being idle
 * the uMCTL2 automatically puts the SDRAM into Self Refresh.
 * The DDRC command channel is considered idle when there are no HIF
 * commands outstanding. This must be enabled in the PWRCTL.selfref_en.
 * Unit: Multiples of 32 DFI clocks.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PWRTMG_selfref_to_x32_Pos 16
#define UDDRC_PWRTMG_selfref_to_x32_Msk \
	(0xffu << UDDRC_PWRTMG_selfref_to_x32_Pos)

#define UDDRC_PWRTMG_selfref_to_x32(value) \
	((UDDRC_PWRTMG_selfref_to_x32_Msk & \
	((value) << UDDRC_PWRTMG_selfref_to_x32_Pos)))

/* -------- UDDRC_HWLPCTL : (UDDRC_REGS Offset: 0x38)
 * Hardware Low Power Control Register --------
 */
/* (UDDRC_HWLPCTL) Enable for Hardware Low Power Interface. */
#define UDDRC_HWLPCTL_hw_lp_en (0x1u << 0)

/* (UDDRC_HWLPCTL) When this bit is programmed to 1 the cactive_in_ddrc pin
 * of the DDRC can be used to exit from the automatic clock stop,
 * automatic power down or automatic self-refresh modes.
 * Note, it will not cause exit of Self-Refresh that was caused by
 * Hardware Low Power Interface and/or Software (PWRCTL.selfref_sw).
 */
#define UDDRC_HWLPCTL_hw_lp_exit_idle_en (0x1u << 1)

/* (UDDRC_HWLPCTL) Hardware idle period. The cactive_ddrc output is driven low
 * if the DDRC command channel is idle for hw_lp_idle * 32 cycles
 * if not in INIT or DPD/MPSM operating_mode.
 * The DDRC command channel is considered idle when there are no HIF
 * commands outstanding.
 * The hardware idle function is disabled when hw_lp_idle_x32=0.
 * Unit: Multiples of 32 DFI clocks.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_HWLPCTL_hw_lp_idle_x32_Pos 16
#define UDDRC_HWLPCTL_hw_lp_idle_x32_Msk \
	(0xfffu << UDDRC_HWLPCTL_hw_lp_idle_x32_Pos)
#define UDDRC_HWLPCTL_hw_lp_idle_x32(value) \
	((UDDRC_HWLPCTL_hw_lp_idle_x32_Msk & \
	((value) << UDDRC_HWLPCTL_hw_lp_idle_x32_Pos)))

/* -------- UDDRC_RFSHCTL0 : (UDDRC_REGS Offset: 0x50)
 * Refresh Control Register 0 --------
 */
/* (UDDRC_RFSHCTL0)
 - 1 - Per bank refresh;
 - 0 - All bank refresh.
 * Per bank refresh allows traffic to flow to other banks.
 * Per bank refresh is not supported by all LPDDR2 devices but should be
 * supported by all LPDDR3/LPDDR4 devices.
 * Present only in designs configured to support LPDDR2/LPDDR3/LPDDR4
 */
#define UDDRC_RFSHCTL0_per_bank_refresh (0x1u << 2)

/* (UDDRC_RFSHCTL0) The programmed value + 1 is the number of refresh timeouts
 * that is allowed to accumulate before traffic is blocked and the refreshes
 * are forced to execute. Closing pages to perform a refresh is a one-time
 * penalty that must be paid for each group of refreshes.
 * Therefore, performing refreshes in a burst reduces the per-refresh penalty
 * of these page closings. Higher numbers for RFSHCTL.refresh_burst slightly
 * increases utilization; lower numbers decreases the worst-case
 * latency associated with refreshes.
 - 0 - single refresh
 - 1 - burst-of-2 refresh
 - 7 - burst-of-8 refresh
 * For information on burst refresh feature refer to section 3.9
 * of DDR2 JEDEC specification - JESD79-2F.pdf.
 * For DDR2/3, the refresh is always per-rank and not per-bank. The rank refresh
 * can be accumulated over 8*tREFI cycles using the burst refresh feature.
 * In DDR4 mode, according to Fine Granularity feature, 8 refreshes can be
 * postponed in 1X mode, 16 refreshes in 2X mode and 32 refreshes in 4X mode.
 * If using PHY-initiated updates, care must be taken in the setting of
 * RFSHCTL0.refresh_burst, to ensure that tRFCmax is not violated due to a
 * PHY-initiated update occurring shortly before a refresh burst was due.
 * In this situation, the refresh burst will be delayed until the PHY-initiated
 * update is complete.
 */

#define UDDRC_RFSHCTL0_refresh_burst_Pos 4
#define UDDRC_RFSHCTL0_refresh_burst_Msk \
	(0x1fu << UDDRC_RFSHCTL0_refresh_burst_Pos)
#define UDDRC_RFSHCTL0_refresh_burst(value) \
	((UDDRC_RFSHCTL0_refresh_burst_Msk & \
	((value) << UDDRC_RFSHCTL0_refresh_burst_Pos)))

/* (UDDRC_RFSHCTL0) If the refresh timer (tRFCnom, also known as tREFI)
 * has expired at least once, but it has not expired (RFSHCTL0.refresh_burst+1)
 * times yet, then a speculative refresh may be performed.
 * A speculative refresh is a refresh performed at a time when refresh would be
 * useful, but before it is absolutely required. When the SDRAM bus is idle for
 * a period of time determined by this RFSHCTL0.refresh_to_x32 and the refresh
 * timer has expired at least once since the last refresh, then a speculative
 * refresh is performed. Speculative refreshes continues successively until
 * there are no refreshes pending or until new reads or writes are issued
 * to the uMCTL2.
 * FOR PERFORMANCE ONLY.
 * Unit: Multiples of 32 DFI clocks.
 */
#define UDDRC_RFSHCTL0_refresh_to_x32_Pos 12
#define UDDRC_RFSHCTL0_refresh_to_x32_Msk \
	(0x1fu << UDDRC_RFSHCTL0_refresh_to_x32_Pos)
#define UDDRC_RFSHCTL0_refresh_to_x32(value) \
	((UDDRC_RFSHCTL0_refresh_to_x32_Msk & \
	((value) << UDDRC_RFSHCTL0_refresh_to_x32_Pos)))

/* (UDDRC_RFSHCTL0) Threshold value in number of DFI clock cycles before the
 * critical refresh or page timer expires. A critical refresh is to be issued
 * before this threshold is reached. It is recommended that this not be changed
 * from the default value, currently shown as 0x2. It must always be less than
 * internally used t_rfc_nom_x32. Note that, in LPDDR2/LPDDR3/LPDDR4,
 * internally used t_rfc_nom_x32 may be equal to RFSHTMG.t_rfc_nom_x32>>2
 * if derating is enabled (DERATEEN.derate_enable=1).
 * Otherwise, internally used t_rfc_nom_x32 will be equal to RFSHTMG.t_rfc_nom_x32.
 * Unit: Multiples of 32 DFI clocks.
 */
#define UDDRC_RFSHCTL0_refresh_margin_Pos 20
#define UDDRC_RFSHCTL0_refresh_margin_Msk \
	(0xfu << UDDRC_RFSHCTL0_refresh_margin_Pos)
#define UDDRC_RFSHCTL0_refresh_margin(value) \
	((UDDRC_RFSHCTL0_refresh_margin_Msk & \
	((value) << UDDRC_RFSHCTL0_refresh_margin_Pos)))

/* -------- UDDRC_RFSHCTL3 : (UDDRC_REGS Offset: 0x60)
 * Refresh Control Register 3 --------
 */
/* (UDDRC_RFSHCTL3) When '1', disable auto-refresh generated by the uMCTL2.
 * When auto-refresh is disabled, the SoC core must generate refreshes using
 * the registers DBGCMD.rankn_refresh.
 * When dis_auto_refresh transitions from 0 to 1, any pending refreshes
 * are immediately scheduled by the uMCTL2.
 * If DDR4 CRC/parity retry is enabled (CRCPARCTL1.crc_parity_retry_enable = 1),
 * disable auto-refresh is not supported, and this bit must be set to '0'.
 * (DDR4 only) If FGR mode is enabled (RFSHCTL3.refresh_mode > 0),
 * disable auto-refresh is not supported, and this bit must be set to '0'.
 * This register field is changeable on the fly.
 */
#define UDDRC_RFSHCTL3_dis_auto_refresh (0x1u << 0)

/* (UDDRC_RFSHCTL3) Toggle this signal (either from 0 to 1 or from 1 to 0)
 * to indicate that the refresh register(s) have been updated.
 * refresh_update_level must not be toggled when the DDRC is in reset
 * (core_ddrc_rstn = 0).
 * The refresh register(s) are automatically updated when exiting reset.
 */
#define UDDRC_RFSHCTL3_refresh_update_level (0x1u << 1)

/* -------- UDDRC_RFSHTMG : (UDDRC_REGS Offset: 0x64)
 * Refresh Timing Register --------
 */
/* (UDDRC_RFSHTMG) tRFC (min): Minimum time from refresh to refresh or activate.
 * When the controller is operating in 1:1 mode, t_rfc_min should be set to
 * RoundUp(tRFCmin/tCK).
 * When the controller is operating in 1:2 mode, t_rfc_min should be set to
 * RoundUp(RoundUp(tRFCmin/tCK)/2).
 * In LPDDR2/LPDDR3/LPDDR4 mode:
  - if using all-bank refreshes, the tRFCmin value in the above equations is
	equal to tRFCab
  - if using per-bank refreshes, the tRFCmin value in the above equations is
	equal to tRFCpb
 * In DDR4 mode, the tRFCmin value in the above equations is different depending
 * on the refresh mode (fixed 1X,2X,4X) and the device density.
 * The user should program the appropriate value from the spec based on the
 * 'refresh_mode' and the device density that is used.
 * Unit: Clocks.
 */
#define UDDRC_RFSHTMG_t_rfc_min_Pos 0
#define UDDRC_RFSHTMG_t_rfc_min_Msk (0x3ffu << UDDRC_RFSHTMG_t_rfc_min_Pos)
#define UDDRC_RFSHTMG_t_rfc_min(value) \
	((UDDRC_RFSHTMG_t_rfc_min_Msk & \
	((value) << UDDRC_RFSHTMG_t_rfc_min_Pos)))

/* (UDDRC_RFSHTMG) Used only when LPDDR3 memory type is connected.
 * Should only be changed when uMCTL2 is in reset. Specifies whether to use the
 * tREFBW parameter (required by some LPDDR3 devices which comply with earlier
 * versions of the LPDDR3 JEDEC specification) or not:
   - 0 - tREFBW parameter not used
   - 1 - tREFBW parameter used
 */
#define UDDRC_RFSHTMG_lpddr3_trefbw_en (0x1u << 15)

/* (UDDRC_RFSHTMG) tREFI: Average time interval between refreshes per rank
 * (Specification: 7.8us for DDR2, DDR3 and DDR4. See JEDEC specification for
 * mDDR, LPDDR2, LPDDR3 and LPDDR4).
 * For LPDDR2/LPDDR3/LPDDR4:
  - if using all-bank refreshes (RFSHCTL0.per_bank_refresh = 0),
	this register should be set to tREFIab
  - if using per-bank refreshes (RFSHCTL0.per_bank_refresh = 1),
	this register should be set to tREFIpb
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to (tREFI/2), no rounding up.
 *
 * In DDR4 mode, tREFI value is different depending on the refresh mode.
 * The user should program the appropriate value from the spec based on the
 * value programmed in the refresh mode register.
 * Note that RFSHTMG.t_rfc_nom_x32 * 32 must be greater than RFSHTMG.t_rfc_min,
 * and RFSHTMG.t_rfc_nom_x32 must be greater than 0x1.
  - Non-DDR4 or DDR4 Fixed 1x mode: RFSHTMG.t_rfc_nom_x32 must be less than
 	or equal to 0xFFE.
  - DDR4 Fixed 2x mode: RFSHTMG.t_rfc_nom_x32 must be less than or equal to 0x7FF.
  - DDR4 Fixed 4x mode: RFSHTMG.t_rfc_nom_x32 must be less than or equal to 0x3FF.
 * Unit: Multiples of 32 clocks.
 */
#define UDDRC_RFSHTMG_t_rfc_nom_x32_Pos 16
#define UDDRC_RFSHTMG_t_rfc_nom_x32_Msk \
	(0xfffu << UDDRC_RFSHTMG_t_rfc_nom_x32_Pos)
#define UDDRC_RFSHTMG_t_rfc_nom_x32(value) \
	((UDDRC_RFSHTMG_t_rfc_nom_x32_Msk & \
((value) << UDDRC_RFSHTMG_t_rfc_nom_x32_Pos)))

/* -------- UDDRC_CRCPARCTL0 : (UDDRC_REGS Offset: 0xC0)
 * CRC Parity Control Register0 --------
 */
/* (UDDRC_CRCPARCTL0) Interrupt enable bit for DFI alert error.
 * If this bit is set, any parity/CRC error detected on the dfi_alert_n input
 * will result in an interrupt being set on CRCPARSTAT.dfi_alert_err_int.
 */
#define UDDRC_CRCPARCTL0_dfi_alert_err_int_en (0x1u << 0)

/* (UDDRC_CRCPARCTL0) Interrupt clear bit for DFI alert error.
 * If this bit is set, the alert error interrupt on CRCPARSTAT.dfi_alert_err_int
 * will be cleared. When the clear operation is complete,
 * the uMCTL2 automatically clears this bit.
 */
#define UDDRC_CRCPARCTL0_dfi_alert_err_int_clr (0x1u << 1)

/* (UDDRC_CRCPARCTL0) DFI alert error count clear.
 * Clear bit for DFI alert error counter. Asserting this bit will clear the DFI
 * alert error counter, CRCPARSTAT.dfi_alert_err_cnt.
 * When the clear operation is complete,
 * the uMCTL2 automatically clears this bit.
 */
#define UDDRC_CRCPARCTL0_dfi_alert_err_cnt_clr (0x1u << 2)

/* -------- UDDRC_CRCPARSTAT : (UDDRC_REGS Offset: 0xCC)
 * CRC Parity Status Register --------
 */
/* (UDDRC_CRCPARSTAT) DFI alert error count.
 * If a parity/CRC error is detected on dfi_alert_n, this counter be incremented.
 * This is independent of the setting of CRCPARCTL0.dfi_alert_err_int_en.
 * It will saturate at 0xFFFF, and can be cleared by asserting
 * CRCPARCTL0.dfi_alert_err_cnt_clr.
 */
#define UDDRC_CRCPARSTAT_dfi_alert_err_cnt_Pos 0
#define UDDRC_CRCPARSTAT_dfi_alert_err_cnt_Msk \
	(0xffffu << UDDRC_CRCPARSTAT_dfi_alert_err_cnt_Pos)

/* (UDDRC_CRCPARSTAT) DFI alert error interrupt.
 * If a parity/CRC error is detected on dfi_alert_n, and the interrupt
 * is enabled by CRCPARCTL0.dfi_alert_err_int_en, this interrupt bit will be set.
 * It will remain set until cleared by CRCPARCTL0.dfi_alert_err_int_clr
 */
#define UDDRC_CRCPARSTAT_dfi_alert_err_int (0x1u << 16)

/* -------- UDDRC_INIT0 : (UDDRC_REGS Offset: 0xD0)
 * SDRAM Initialization Register 0 --------
 */
/* (UDDRC_INIT0) Cycles to wait after reset before driving CKE high to start
 * the SDRAM initialization sequence.
 * Unit: 1024 DFI clock cycles.
 * DDR2 specifications typically require this to be programmed for
 * a delay of >= 200 us.
 * LPDDR2/LPDDR3: tINIT1 of 100 ns (min)
 * LPDDR4: tINIT3 of 2 ms (min)
 * When the controller is operating in 1:2 frequency ratio mode, program this
 * to JEDEC spec value divided by 2, and round it up to the next integer value.
 * For DDR3/DDR4 RDIMMs, this should include the time needed to satisfy tSTAB
 */
#define UDDRC_INIT0_pre_cke_x1024_Pos 0
#define UDDRC_INIT0_pre_cke_x1024_Msk \
	(0xfffu << UDDRC_INIT0_pre_cke_x1024_Pos)
#define UDDRC_INIT0_pre_cke_x1024(value) \
	((UDDRC_INIT0_pre_cke_x1024_Msk & \
	((value) << UDDRC_INIT0_pre_cke_x1024_Pos)))

/* (UDDRC_INIT0) Cycles to wait after driving CKE high to start
 * the SDRAM initialization sequence.
 * Unit: 1024 DFI clock cycles.
 * DDR2 typically requires a 400 ns delay, requiring this value to be
 * programmed to 2 at all clock speeds.
 * LPDDR2/LPDDR3 typically requires this to be programmed for a delay of 200 us.
 * LPDDR4 typically requires this to be programmed for a delay of 2 us.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * JEDEC spec value divided by 2, and round it up to the next integer value.
 */
#define UDDRC_INIT0_post_cke_x1024_Pos 16
#define UDDRC_INIT0_post_cke_x1024_Msk \
	(0x3ffu << UDDRC_INIT0_post_cke_x1024_Pos)
#define UDDRC_INIT0_post_cke_x1024(value) \
	((UDDRC_INIT0_post_cke_x1024_Msk & \
	((value) << UDDRC_INIT0_post_cke_x1024_Pos)))

/* (UDDRC_INIT0) If lower bit is enabled the SDRAM initialization routine
 * is skipped. The upper bit decides what state the controller starts up
 * in when reset is removed
 - 00 - SDRAM Intialization routine is run after power-up
 - 01 - SDRAM Intialization routine is skipped after power-up.
 	Controller starts up in Normal Mode
 - 11 - SDRAM Intialization routine is skipped after power-up.
	Controller starts up in Self-refresh Mode
 - 10 - SDRAM Intialization routine is run after power-up.
 */
#define UDDRC_INIT0_skip_dram_init_Pos 30
#define UDDRC_INIT0_skip_dram_init_Msk \
	(0x3u << UDDRC_INIT0_skip_dram_init_Pos)
#define UDDRC_INIT0_skip_dram_init(value) \
	((UDDRC_INIT0_skip_dram_init_Msk & \
	((value) << UDDRC_INIT0_skip_dram_init_Pos)))

/* -------- UDDRC_INIT1 : (UDDRC_REGS Offset: 0xD4)
 * SDRAM Initialization Register 1 --------
 */
/* (UDDRC_INIT1) Wait period before driving the OCD complete command to SDRAM.
 * Unit: Counts of a global timer that pulses every 32 DFI clock cycles.
 * There is no known specific requirement for this; it may be set to zero.
 */
#define UDDRC_INIT1_pre_ocd_x32_Pos 0
#define UDDRC_INIT1_pre_ocd_x32_Msk \
	(0xfu << UDDRC_INIT1_pre_ocd_x32_Pos)
#define UDDRC_INIT1_pre_ocd_x32(value) \
	((UDDRC_INIT1_pre_ocd_x32_Msk & \
	((value) << UDDRC_INIT1_pre_ocd_x32_Pos)))

/* (UDDRC_INIT1) Number of cycles to assert SDRAM reset signal during
 * init sequence.
 * This is only present for designs supporting DDR3, DDR4 or LPDDR4 devices.
 * For use with a Synopsys DDR PHY, this should be set to a minimum of 1.
 * When the controller is operating in 1:2 frequency ratio mode, program this
 * to JEDEC spec value divided by 2, and round it up to the next integer value.
 * Unit: 1024 DFI clock cycles.
 */
#define UDDRC_INIT1_dram_rstn_x1024_Pos 16
#define UDDRC_INIT1_dram_rstn_x1024_Msk \
	(0x1ffu << UDDRC_INIT1_dram_rstn_x1024_Pos)
#define UDDRC_INIT1_dram_rstn_x1024(value) \
	((UDDRC_INIT1_dram_rstn_x1024_Msk & \
	((value) << UDDRC_INIT1_dram_rstn_x1024_Pos)))

/* -------- UDDRC_INIT2 : (UDDRC_REGS Offset: 0xD8)
 * SDRAM Initialization Register 2 --------
 */
/* (UDDRC_INIT2) Time to wait after the first CKE high, tINIT2.
 * Present only in designs configured to support LPDDR2/LPDDR3.
 * LPDDR2/LPDDR3 typically requires 5 x tCK delay.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to JEDEC spec value divided by 2, and round it up to the
 * next integer value.
 * Unit: DFI clock cycles.
 */
#define UDDRC_INIT2_min_stable_clock_x1_Pos 0
#define UDDRC_INIT2_min_stable_clock_x1_Msk \
	(0xfu << UDDRC_INIT2_min_stable_clock_x1_Pos)
#define UDDRC_INIT2_min_stable_clock_x1(value) \
	((UDDRC_INIT2_min_stable_clock_x1_Msk & \
	((value) << UDDRC_INIT2_min_stable_clock_x1_Pos)))

/* (UDDRC_INIT2) Idle time after the reset command, tINIT4.
 * Present only in designs configured to support LPDDR2.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to JEDEC spec value divided by 2, and round it up to
 * the next integer value.
 * Unit: 32 DFI clock cycles.
 */
#define UDDRC_INIT2_idle_after_reset_x32_Pos 8
#define UDDRC_INIT2_idle_after_reset_x32_Msk \
	(0xffu << UDDRC_INIT2_idle_after_reset_x32_Pos)
#define UDDRC_INIT2_idle_after_reset_x32(value) \
	((UDDRC_INIT2_idle_after_reset_x32_Msk & \
	((value) << UDDRC_INIT2_idle_after_reset_x32_Pos)))

/* -------- UDDRC_INIT3 : (UDDRC_REGS Offset: 0xDC)
 * SDRAM Initialization Register 3 --------
 */
/* (UDDRC_INIT3) DDR2: Value to write to EMR register.
 * Bits 9:7 are for OCD and the setting in this register is ignored.
 * The uMCTL2 sets those bits appropriately.
 * DDR3/DDR4: Value to write to MR1 register  Set bit 7 to 0.
 * If PHY-evaluation mode training is enabled, this bit is set appropriately by
 * the uMCTL2 during write leveling.
 * mDDR: Value to write to EMR register.
 * LPDDR2/LPDDR3/LPDDR4 - Value to write to MR2 register
 */
#define UDDRC_INIT3_emr_Pos 0
#define UDDRC_INIT3_emr_Msk (0xffffu << UDDRC_INIT3_emr_Pos)
#define UDDRC_INIT3_emr(value) \
	((UDDRC_INIT3_emr_Msk & ((value) << UDDRC_INIT3_emr_Pos)))

/* (UDDRC_INIT3) DDR2: Value to write to MR register.
 * Bit 8 is for DLL and the setting here is ignored.
 * The uMCTL2 sets this bit appropriately.
 * DDR3/DDR4: Value loaded into MR0 register.
 * mDDR: Value to write to MR register.
 * LPDDR2/LPDDR3/LPDDR4 - Value to write to MR1 register
 */
#define UDDRC_INIT3_mr_Pos 16
#define UDDRC_INIT3_mr_Msk (0xffffu << UDDRC_INIT3_mr_Pos)
#define UDDRC_INIT3_mr(value) \
	((UDDRC_INIT3_mr_Msk & ((value) << UDDRC_INIT3_mr_Pos)))

/* -------- UDDRC_INIT4 : (UDDRC_REGS Offset: 0xE0) SDRAM Initialization Register 4 -------- */
/* (UDDRC_INIT4) DDR2: Value to write to EMR3 register.
 * DDR3/DDR4: Value to write to MR3 register
 * mDDR/LPDDR2/LPDDR3: Unused
 * LPDDR4: Value to write to MR13 register
 */
#define UDDRC_INIT4_emr3_Pos 0
#define UDDRC_INIT4_emr3_Msk (0xffffu << UDDRC_INIT4_emr3_Pos)
#define UDDRC_INIT4_emr3(value) \
	((UDDRC_INIT4_emr3_Msk & ((value) << UDDRC_INIT4_emr3_Pos)))

/* (UDDRC_INIT4) DDR2: Value to write to EMR2 register.
 * DDR3/DDR4: Value to write to MR2 register
 * LPDDR2/LPDDR3/LPDDR4: Value to write to MR3 register
 * mDDR: Unused
 */
#define UDDRC_INIT4_emr2_Pos 16
#define UDDRC_INIT4_emr2_Msk (0xffffu << UDDRC_INIT4_emr2_Pos)
#define UDDRC_INIT4_emr2(value) \
	((UDDRC_INIT4_emr2_Msk & ((value) << UDDRC_INIT4_emr2_Pos)))

/* -------- UDDRC_INIT5 : (UDDRC_REGS Offset: 0xE4)
 * SDRAM Initialization Register 5 --------
 */
/* (UDDRC_INIT5) Maximum duration of the auto initialization, tINIT5.
 * Present only in designs configured to support LPDDR2/LPDDR3.
 * LPDDR2/LPDDR3 typically requires 10 us.
 * Unit: 1024 DFI clock cycles.
 */
#define UDDRC_INIT5_max_auto_init_x1024_Pos 0
#define UDDRC_INIT5_max_auto_init_x1024_Msk \
	(0x3ffu << UDDRC_INIT5_max_auto_init_x1024_Pos)
#define UDDRC_INIT5_max_auto_init_x1024(value) \
	((UDDRC_INIT5_max_auto_init_x1024_Msk & \
	((value) << UDDRC_INIT5_max_auto_init_x1024_Pos)))

/* (UDDRC_INIT5) ZQ initial calibration, tZQINIT.
 * Present only in designs configured to support DDR3 or DDR4 or LPDDR2/LPDDR3.
 * DDR3 typically requires 512 SDRAM clock cycles.
 * DDR4 requires 1024 SDRAM clock cycles.
 * LPDDR2/LPDDR3 requires 1 us.
 * When the controller is operating in 1:2 frequency ratio mode, program this
 * to JEDEC spec value divided by 2, and round it up to the next integer value.
 * Unit: 32 DFI clock cycles.
 */
#define UDDRC_INIT5_dev_zqinit_x32_Pos 16
#define UDDRC_INIT5_dev_zqinit_x32_Msk (0xffu << UDDRC_INIT5_dev_zqinit_x32_Pos)
#define UDDRC_INIT5_dev_zqinit_x32(value) \
	((UDDRC_INIT5_dev_zqinit_x32_Msk & \
	((value) << UDDRC_INIT5_dev_zqinit_x32_Pos)))

/* -------- UDDRC_DIMMCTL : (UDDRC_REGS Offset: 0xF0)
 * DIMM Control Register --------
 */
/* (UDDRC_DIMMCTL) Staggering enable for multi-rank accesses
 * (for multi-rank UDIMM, RDIMM and LRDIMM implementations only).
 * This is not supported for mDDR, LPDDR2, LPDDR3 or LPDDR4 SDRAMs.
 * Note: Even if this bit is set it does not take care of software driven
 * MR commands (via MRCTRL0/MRCTRL1), where software is responsible to
 * send them to seperate ranks as appropriate.
    - 1 - (DDR4) Send MRS commands to each ranks seperately
    - 1 - (non-DDR4) Send all commands to even and odd ranks seperately
    - 0 - Do not stagger accesses
 */
#define UDDRC_DIMMCTL_dimm_stagger_cs_en (0x1u << 0)

/* (UDDRC_DIMMCTL) Address Mirroring Enable
 * (for multi-rank UDIMM implementations and multi-rank DDR4
 * RDIMM/LRDIMM implementations).
 * Some UDIMMs and DDR4 RDIMMs/LRDIMMs implement address mirroring for
 * odd ranks, which means that the following address, bank address and
 * bank group bits are swapped: (A3, A4), (A5, A6), (A7, A8), (BA0, BA1)
 * and also (A11, A13), (BG0, BG1) for the DDR4. Setting this bit ensures that,
 * for mode register accesses during the automatic initialization routine,
 * these bits are swapped within the uMCTL2 to compensate for this
 * UDIMM/RDIMM/LRDIMM swapping. In addition to the automatic initialization
 * routine, in case of DDR4 UDIMM/RDIMM/LRDIMM, they are swapped during
 * the automatic MRS access to enable/disable of a particular DDR4 feature.
 * Note: This has no effect on the address of any other memory accesses,
 * or of software-driven mode register accesses.
 * This is not supported for mDDR, LPDDR2, LPDDR3 or LPDDR4 SDRAMs.
 * Note: In case of x16 DDR4 DIMMs, BG1 output of MRS for the odd ranks is same
 * as BG0 because BG1 is invalid, hence dimm_dis_bg_mirroring register
 * must be set to 1.
    - 1 - For odd ranks, implement address mirroring for MRS commands to
	during initialization and for any automatic DDR4 MRS commands
	(to be used if UDIMM/RDIMM/LRDIMM implements address mirroring)
    - 0 - Do not implement address mirroring
 */
#define UDDRC_DIMMCTL_dimm_addr_mirr_en (0x1u << 1)

/* -------- UDDRC_DRAMTMG0 : (UDDRC_REGS Offset: 0x100)
 * SDRAM Timing Register 0 --------
 */
/* (UDDRC_DRAMTMG0) tRAS(min):  Minimum time between activate and
 * precharge to the same bank.
 * When the controller is operating in 1:2 frequency mode, 1T mode,
 * program this to tRAS(min)/2. No rounding up.
 * When the controller is operating in 1:2 frequency ratio mode, 2T mode
 * or LPDDR4 mode, program this to (tRAS(min)/2) and round it up to the next
 * integer value.
 * Unit: Clocks
 */
#define UDDRC_DRAMTMG0_t_ras_min_Pos 0
#define UDDRC_DRAMTMG0_t_ras_min_Msk (0x3fu << UDDRC_DRAMTMG0_t_ras_min_Pos)
#define UDDRC_DRAMTMG0_t_ras_min(value) \
	((UDDRC_DRAMTMG0_t_ras_min_Msk & \
	((value) << UDDRC_DRAMTMG0_t_ras_min_Pos)))

/* (UDDRC_DRAMTMG0) tRAS(max):  Maximum time between activate and precharge
 * to same bank. This is the maximum time that a page can be kept open
 * Minimum value of this register is 1. Zero is invalid.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to (tRAS(max)-1)/2. No rounding up.
 * Unit: Multiples of 1024 clocks.
 */
#define UDDRC_DRAMTMG0_t_ras_max_Pos 8
#define UDDRC_DRAMTMG0_t_ras_max_Msk (0x7fu << UDDRC_DRAMTMG0_t_ras_max_Pos)
#define UDDRC_DRAMTMG0_t_ras_max(value) \
	((UDDRC_DRAMTMG0_t_ras_max_Msk & \
	((value) << UDDRC_DRAMTMG0_t_ras_max_Pos)))

/* (UDDRC_DRAMTMG0) tFAW Valid only when 8 or more banks(or banks x bank groups)
 * are present.
 * In 8-bank design, at most 4 banks must be activated in a rolling window
 * of tFAW cycles.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to (tFAW/2) and round up to next integer value.
 * In a 4-bank design, set this register to 0x1 independent of
 * the 1:1/1:2 frequency mode.
Unit: Clocks
 */
#define UDDRC_DRAMTMG0_t_faw_Pos 16
#define UDDRC_DRAMTMG0_t_faw_Msk (0x3fu << UDDRC_DRAMTMG0_t_faw_Pos)
#define UDDRC_DRAMTMG0_t_faw(value) \
	((UDDRC_DRAMTMG0_t_faw_Msk & \
	((value) << UDDRC_DRAMTMG0_t_faw_Pos)))

/* (UDDRC_DRAMTMG0) Minimum time between write and precharge to same bank.
 * Unit: Clocks
 * Specifications: WL + BL/2 + tWR = approximately 8 cycles + 15 ns =
 * 14 clocks @400MHz and less for lower frequencies
 * where:
 - WL = write latency
 - BL = burst length. This must match the value programmed in the BL bit of the
	mode register to the SDRAM.
	BST (burst terminate) is not supported at present.
 - tWR = Write recovery time. This comes directly from the SDRAM specification.
 * Add one extra cycle for LPDDR2/LPDDR3/LPDDR4 for this parameter.
 *
 * When the controller is operating in 1:2 frequency ratio mode, 1T mode,
 * divide the above value by 2. No rounding up.
 * When the controller is operating in 1:2 frequency ratio mode, 2T mode
 * or LPDDR4 mode, divide the above value by 2 and round it up to the
 * next integer value.
 *
 * Note that, depending on the PHY, if using LRDIMM, it may be necessary
 * to adjust the value of this parameter to compensate for the extra cycle
 * of latency through the LRDIMM.
 */
#define UDDRC_DRAMTMG0_wr2pre_Pos 24
#define UDDRC_DRAMTMG0_wr2pre_Msk (0x7fu << UDDRC_DRAMTMG0_wr2pre_Pos)
#define UDDRC_DRAMTMG0_wr2pre(value) \
	((UDDRC_DRAMTMG0_wr2pre_Msk & \
	((value) << UDDRC_DRAMTMG0_wr2pre_Pos)))

/* -------- UDDRC_DRAMTMG1 : (UDDRC_REGS Offset: 0x104)
 * SDRAM Timing Register 1 --------
 */
/* (UDDRC_DRAMTMG1) tRC:  Minimum time between activates to same bank.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to (tRC/2) and round up to next integer value.
 * Unit: Clocks.
 */
#define UDDRC_DRAMTMG1_t_rc_Pos 0
#define UDDRC_DRAMTMG1_t_rc_Msk (0x7fu << UDDRC_DRAMTMG1_t_rc_Pos)
#define UDDRC_DRAMTMG1_t_rc(value) \
	((UDDRC_DRAMTMG1_t_rc_Msk & ((value) << UDDRC_DRAMTMG1_t_rc_Pos)))

/* (UDDRC_DRAMTMG1) tRTP:  Minimum time from read to precharge of same bank.
 - DDR2: tAL + BL/2 + max(tRTP, 2) - 2
 - DDR3: tAL + max (tRTP, 4)
 - DDR4: Max of following two equations:
           tAL + max (tRTP, 4) or,
           RL + BL/2 - tRP (*).
 - mDDR: BL/2
 - LPDDR2: Depends on if it's LPDDR2-S2 or LPDDR2-S4:
           LPDDR2-S2: BL/2 + tRTP - 1.
           LPDDR2-S4: BL/2 + max(tRTP,2) - 2.
 - LPDDR3: BL/2 +  max(tRTP,4) - 4
 - LPDDR4: BL/2 + max(tRTP,8) - 8
(*) When both DDR4 SDRAM and ST-MRAM are used simultaneously,
 * use SDRAM's tRP value for calculation.
 *
 * When the controller is operating in 1:2 mode, 1T mode, divide the above
 * value by 2. No rounding up.
 * When the controller is operating in 1:2 mode, 2T mode or LPDDR4 mode,
 * divide the above value by 2 and round it up to the next integer value.
 * Unit: Clocks.
 */
#define UDDRC_DRAMTMG1_rd2pre_Pos 8
#define UDDRC_DRAMTMG1_rd2pre_Msk (0x3fu << UDDRC_DRAMTMG1_rd2pre_Pos)
#define UDDRC_DRAMTMG1_rd2pre(value) \
	((UDDRC_DRAMTMG1_rd2pre_Msk & ((value) << UDDRC_DRAMTMG1_rd2pre_Pos)))

/* (UDDRC_DRAMTMG1) tXP: Minimum time after power-down exit to any operation.
 * For DDR3, this should be programmed to tXPDLL if slow powerdown exit
 * is selected in MR0[12].
 * If C/A parity for DDR4 is used, set to (tXP+PL) instead.
 * If LPDDR4 is selected and its spec has tCKELPD parameter,
 * set to the larger of tXP and tCKELPD instead.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to (tXP/2) and round it up to the next integer value.
 * Units: Clocks
 */
#define UDDRC_DRAMTMG1_t_xp_Pos 16
#define UDDRC_DRAMTMG1_t_xp_Msk (0x1fu << UDDRC_DRAMTMG1_t_xp_Pos)
#define UDDRC_DRAMTMG1_t_xp(value) \
	((UDDRC_DRAMTMG1_t_xp_Msk & ((value) << UDDRC_DRAMTMG1_t_xp_Pos)))

/* -------- UDDRC_DRAMTMG2 : (UDDRC_REGS Offset: 0x108)
 * SDRAM Timing Register 2 --------
 */

/* (UDDRC_DRAMTMG2)
 * DDR4: CWL + PL + BL/2 + tWTR_L
 * LPDDR2/3/4: WL + BL/2 + tWTR + 1
 * Others: CWL + BL/2 + tWTR
 * In DDR4, minimum time from write command to read command for same bank group.
 * In others, minimum time from write command to read command. Includes time for
 * bus turnaround, recovery times, and all per-bank, per-rank,
 * and global constraints.
 * Unit: Clocks.
 *Where:
 - CWL = CAS write latency
 - WL = Write latency
 - PL = Parity latency
 - BL = burst length. This must match the value programmed in the BL bit of the
	mode register to the SDRAM
 - tWTR_L = internal write to read command delay for same bank group.
	This comes directly from the SDRAM specification.
 - tWTR = internal write to read command delay.
	This comes directly from the SDRAM specification.
 * Add one extra cycle for LPDDR2/LPDDR3/LPDDR4 operation.
 * When the controller is operating in 1:2 mode, divide the value calculated
 * using the above equation by 2, and round it up to next integer.
 */
#define UDDRC_DRAMTMG2_wr2rd_Pos 0
#define UDDRC_DRAMTMG2_wr2rd_Msk (0x3fu << UDDRC_DRAMTMG2_wr2rd_Pos)
#define UDDRC_DRAMTMG2_wr2rd(value) \
	((UDDRC_DRAMTMG2_wr2rd_Msk & ((value) << UDDRC_DRAMTMG2_wr2rd_Pos)))

/* (UDDRC_DRAMTMG2)
 * DDR2/3/mDDR: RL + BL/2 + 2 - WL
 * DDR4: RL + BL/2 + 1 + WR_PREAMBLE - WL
 * LPDDR2/LPDDR3: RL + BL/2 + RU(tDQSCKmax/tCK) + 1 - WL
 * LPDDR4(DQ ODT is Disabled): RL + BL/2 + RU(tDQSCKmax/tCK) +
		WR_PREAMBLE + RD_POSTAMBLE - WL
 * LPDDR4(DQ ODT is Enabled) : RL + BL/2 + RU(tDQSCKmax/tCK) +
		RD_POSTAMBLE - ODTLon - RU(tODTon(min)/tCK)
 *
 * Minimum time from read command to write command. Include time for bus
 * turnaround and all per-bank, per-rank, and global constraints.
 * Please see the relevant PHY databook for details of what should be
 * included here.
 * Unit: Clocks.
 * Where:
 - WL = write latency
 - BL = burst length. This must match the value programmed in the BL bit of
	the mode register to the SDRAM
 - RL = read latency = CAS latency
 - WR_PREAMBLE = write preamble. This is unique to DDR4 and LPDDR4.
 - RD_POSTAMBLE = read postamble. This is unique to LPDDR4.
 * For LPDDR2/LPDDR3/LPDDR4, if derating is enabled (DERATEEN.derate_enable=1),
 * derated tDQSCKmax should be used.
 *
 * When the controller is operating in 1:2 frequency ratio mode,
 * divide the value calculated using the above equation by 2,
 * and round it up to next integer.
 *
 * Note that, depending on the PHY, if using LRDIMM, it may be necessary to
 * adjust the value of this parameter to compensate for the extra cycle of
 * latency through the LRDIMM.
 */
#define UDDRC_DRAMTMG2_rd2wr_Pos 8
#define UDDRC_DRAMTMG2_rd2wr_Msk (0x3fu << UDDRC_DRAMTMG2_rd2wr_Pos)
#define UDDRC_DRAMTMG2_rd2wr(value) \
	((UDDRC_DRAMTMG2_rd2wr_Msk & ((value) << UDDRC_DRAMTMG2_rd2wr_Pos)))

/* (UDDRC_DRAMTMG2) Set to RL
 * Time from read command to read data on SDRAM interface.
 * This must be set to RL.
 * Note that, depending on the PHY, if using RDIMM/LRDIMM, it may be necessary
 * to adjust the value of RL to compensate for the extra cycle of latency
 * through the RDIMM/LRDIMM.
 * When the controller is operating in 1:2 frequency ratio mode, divide the
 * value calculated using the above equation by 2, and round it up
 * to next integer.
 * This register field is not required for DDR2 and DDR3 (except if
 * MEMC_TRAINING is set), as the DFI read and write latencies defined
 * in DFITMG0 and DFITMG1 are sufficient for those protocols
 * Unit: clocks
 */
#define UDDRC_DRAMTMG2_read_latency_Pos 16
#define UDDRC_DRAMTMG2_read_latency_Msk \
	(0x3fu << UDDRC_DRAMTMG2_read_latency_Pos)
#define UDDRC_DRAMTMG2_read_latency(value) \
	((UDDRC_DRAMTMG2_read_latency_Msk & \
	((value) << UDDRC_DRAMTMG2_read_latency_Pos)))

/* (UDDRC_DRAMTMG2) Set to WL
 * Time from write command to write data on SDRAM interface.
 * This must be set to WL.
 * For mDDR, it should normally be set to 1.
 * Note that, depending on the PHY, if using RDIMM/LRDIMM, it may be necessary
 * to adjust the value of WL to compensate for the extra cycle of latency
 * through the RDIMM/LRDIMM.
 * When the controller is operating in 1:2 frequency ratio mode, divide the
 * value calculated using the above equation by 2, and round it up to
 * next integer.
 * This register field is not required for DDR2 and DDR3 (except if
 * MEMC_TRAINING is set), as the DFI read and write latencies defined in
 * DFITMG0 and DFITMG1 are sufficient for those protocols
 * Unit: clocks
 */
#define UDDRC_DRAMTMG2_write_latency_Pos 24
#define UDDRC_DRAMTMG2_write_latency_Msk \
	(0x3fu << UDDRC_DRAMTMG2_write_latency_Pos)
#define UDDRC_DRAMTMG2_write_latency(value) \
	((UDDRC_DRAMTMG2_write_latency_Msk & \
	((value) << UDDRC_DRAMTMG2_write_latency_Pos)))

/* -------- UDDRC_DRAMTMG3 : (UDDRC_REGS Offset: 0x10C)
 * SDRAM Timing Register 3 --------
 */
/* (UDDRC_DRAMTMG3) tMOD: Parameter used only in DDR3 and DDR4.
 * Cycles between load mode command and following non-load mode command.
 * If C/A parity for DDR4 is used, set to tMOD_PAR(tMOD+PL) instead.
 * If MPR writes for DDR4 are used, set to tMOD + AL
 * (or tMPD_PAR + AL if C/A parity is also used).
 * Set to tMOD if controller is operating in 1:1 frequency ratio mode,
 * or tMOD/2 (rounded up to next integer) if controller is operating in
 * 1:2 frequency ratio mode.
 * Note that if using RDIMM/LRDIMM, depending on the PHY, it may be necessary
 * to adjust the value of this parameter to compensate for the extra cycle of
 * latency applied to mode register writes by the RDIMM/LRDIMM chip.
 * Also note that if using LRDIMM, the minimum value of this register is
 * tMRD_L2 if controller is operating in 1:1 frequency ratio mode,
 * or tMRD_L2/2 (rounded up to next integer) if controller is operating in
 * 1:2 frequency ratio mode.
 */
#define UDDRC_DRAMTMG3_t_mod_Pos 0
#define UDDRC_DRAMTMG3_t_mod_Msk (0x3ffu << UDDRC_DRAMTMG3_t_mod_Pos)
#define UDDRC_DRAMTMG3_t_mod(value) \
	((UDDRC_DRAMTMG3_t_mod_Msk & ((value) << UDDRC_DRAMTMG3_t_mod_Pos)))

/* (UDDRC_DRAMTMG3) tMRD: Cycles to wait after a mode register write or read.
 * Depending on the connected SDRAM, tMRD represents:
 * DDR2/mDDR: Time from MRS to any command
 * DDR3/4: Time from MRS to MRS command
 * LPDDR2: not used
 * LPDDR3/4: Time from MRS to non-MRS command.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * (tMRD/2) and round it up to the next integer value.
 * If C/A parity for DDR4 is used, set to tMRD_PAR(tMOD+PL) instead.
 */
#define UDDRC_DRAMTMG3_t_mrd_Pos 12
#define UDDRC_DRAMTMG3_t_mrd_Msk (0x3fu << UDDRC_DRAMTMG3_t_mrd_Pos)
#define UDDRC_DRAMTMG3_t_mrd(value) \
	((UDDRC_DRAMTMG3_t_mrd_Msk & ((value) << UDDRC_DRAMTMG3_t_mrd_Pos)))

/* (UDDRC_DRAMTMG3) Time to wait after a mode register write or read (MRW or MRR).
 * Present only in designs configured to support LPDDR2, LPDDR3 or LPDDR4.
 * LPDDR2 typically requires value of 5.
 * LPDDR3 typically requires value of 10.
 * LPDDR4: Set this to the larger of tMRW and tMRWCKEL.
 * For LPDDR2, this register is used for the time from a MRW/MRR
 * to all other commands.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to the above values divided by 2 and round it up
 * to the next integer value.
 * For LDPDR3, this register is used for the time from a MRW/MRR to a MRW/MRR.
 */
#define UDDRC_DRAMTMG3_t_mrw_Pos 20
#define UDDRC_DRAMTMG3_t_mrw_Msk (0x3ffu << UDDRC_DRAMTMG3_t_mrw_Pos)
#define UDDRC_DRAMTMG3_t_mrw(value) \
	((UDDRC_DRAMTMG3_t_mrw_Msk & ((value) << UDDRC_DRAMTMG3_t_mrw_Pos)))


/* -------- UDDRC_DRAMTMG4 : (UDDRC_REGS Offset: 0x110)
 * SDRAM Timing Register 4 --------
 */
/* (UDDRC_DRAMTMG4) tRP:  Minimum time from precharge to activate of same bank.
 * When the controller is operating in 1:1 frequency ratio mode, t_rp should be
 * set to RoundUp(tRP/tCK).
 * When the controller is operating in 1:2 frequency ratio mode, t_rp should be
 * set to RoundDown(RoundUp(tRP/tCK)/2) + 1.
 * When the controller is operating in 1:2 frequency ratio mode in LPDDR4, t_rp
 * should be set to RoundUp(RoundUp(tRP/tCK)/2).
 * Unit: Clocks.
 */
#define UDDRC_DRAMTMG4_t_rp_Pos 0
#define UDDRC_DRAMTMG4_t_rp_Msk (0x1fu << UDDRC_DRAMTMG4_t_rp_Pos)
#define UDDRC_DRAMTMG4_t_rp(value) \
	((UDDRC_DRAMTMG4_t_rp_Msk & ((value) << UDDRC_DRAMTMG4_t_rp_Pos)))

/* (UDDRC_DRAMTMG4) DDR4: tRRD_L: Minimum time between activates from
 * bank "a" to bank "b" for same bank group.
 * Others: tRRD: Minimum time between activates from bank "a" to bank "b"
 * When the controller is operating in 1:2 frequency ratio mode, program this
 * to (tRRD_L/2 or tRRD/2) and round it up to the next integer value.
 * Unit: Clocks.
 */
#define UDDRC_DRAMTMG4_t_rrd_Pos 8
#define UDDRC_DRAMTMG4_t_rrd_Msk (0xfu << UDDRC_DRAMTMG4_t_rrd_Pos)
#define UDDRC_DRAMTMG4_t_rrd(value) \
	((UDDRC_DRAMTMG4_t_rrd_Msk & ((value) << UDDRC_DRAMTMG4_t_rrd_Pos)))

/* (UDDRC_DRAMTMG4) DDR4: tCCD_L: This is the minimum time between
 * two reads or two writes for same bank group.
 * Others: tCCD: This is the minimum time between two reads or two writes.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * (tCCD_L/2 or tCCD/2) and round it up to the next integer value.
 * Unit: clocks.
 */
#define UDDRC_DRAMTMG4_t_ccd_Pos 16
#define UDDRC_DRAMTMG4_t_ccd_Msk (0xfu << UDDRC_DRAMTMG4_t_ccd_Pos)
#define UDDRC_DRAMTMG4_t_ccd(value) \
	((UDDRC_DRAMTMG4_t_ccd_Msk & ((value) << UDDRC_DRAMTMG4_t_ccd_Pos)))

/* (UDDRC_DRAMTMG4) tRCD - tAL: Minimum time from activate to read or write
 * command to same bank.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * ((tRCD - tAL)/2) and round it up to the next integer value.
 * Minimum value allowed for this register is 1, which implies minimum
 * (tRCD - tAL) value to be 2 when the controller is operating in 1:2 frequency
 * ratio mode.
 *Unit: Clocks.
 */
#define UDDRC_DRAMTMG4_t_rcd_Pos 24
#define UDDRC_DRAMTMG4_t_rcd_Msk (0x1fu << UDDRC_DRAMTMG4_t_rcd_Pos)
#define UDDRC_DRAMTMG4_t_rcd(value) \
	((UDDRC_DRAMTMG4_t_rcd_Msk & ((value) << UDDRC_DRAMTMG4_t_rcd_Pos)))

/* -------- UDDRC_DRAMTMG5 : (UDDRC_REGS Offset: 0x114)
 SDRAM Timing Register 5 --------
 */
/* (UDDRC_DRAMTMG5) Minimum number of cycles of CKE HIGH/LOW during
 * power-down and self refresh.
 - LPDDR2/LPDDR3 mode: Set this to the larger of tCKE or tCKESR
 - LPDDR4 mode: Set this to the larger of tCKE or tSR.
 - Non-LPDDR2/non-LPDDR3/non-LPDDR4 designs: Set this to tCKE value.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * (value described above)/2 and round it up to the next integer value.
 * Unit: Clocks.
 */
#define UDDRC_DRAMTMG5_t_cke_Pos 0
#define UDDRC_DRAMTMG5_t_cke_Msk (0x1fu << UDDRC_DRAMTMG5_t_cke_Pos)
#define UDDRC_DRAMTMG5_t_cke(value) \
	((UDDRC_DRAMTMG5_t_cke_Msk & ((value) << UDDRC_DRAMTMG5_t_cke_Pos)))

/* (UDDRC_DRAMTMG5) Minimum CKE low width for Self refresh or Self refresh
 * power down entry to exit timing in memory clock cycles.
 * Recommended settings:
    - mDDR: tRFC
    - LPDDR2: tCKESR
    - LPDDR3: tCKESR
    - LPDDR4: max(tCKE, tSR)
    - DDR2: tCKE
    - DDR3: tCKE + 1
    - DDR4: tCKE + 1 (+ PL(parity latency)(*))
 * (*)Only if CRCPARCTL1.caparity_disable_before_sr=0,
 * this register should be increased by PL.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 */
#define UDDRC_DRAMTMG5_t_ckesr_Pos 8
#define UDDRC_DRAMTMG5_t_ckesr_Msk (0x3fu << UDDRC_DRAMTMG5_t_ckesr_Pos)
#define UDDRC_DRAMTMG5_t_ckesr(value) \
	((UDDRC_DRAMTMG5_t_ckesr_Msk & ((value) << UDDRC_DRAMTMG5_t_ckesr_Pos)))

/* (UDDRC_DRAMTMG5) This is the time after Self Refresh Down Entry that CK is
 * maintained as a valid clock. Specifies the clock disable delay after SRE.
 * Recommended settings:
    - mDDR: 0
    - LPDDR2: 2
    - LPDDR3: 2
    - LPDDR4: tCKELCK
    - DDR2: 1
    - DDR3: max (10 ns, 5 tCK)
    - DDR4: max (10 ns, 5 tCK) (+ PL(parity latency)(*))
 *
 * (*)Only if CRCPARCTL1.caparity_disable_before_sr=0,
 * this register should be increased by PL.
 * When the controller is operating in 1:2 frequency ratio mode,
 * program this to recommended value divided by two and round it up to
 * next integer.
 */
#define UDDRC_DRAMTMG5_t_cksre_Pos 16
#define UDDRC_DRAMTMG5_t_cksre_Msk (0xfu << UDDRC_DRAMTMG5_t_cksre_Pos)
#define UDDRC_DRAMTMG5_t_cksre(value) \
	((UDDRC_DRAMTMG5_t_cksre_Msk & ((value) << UDDRC_DRAMTMG5_t_cksre_Pos)))

/* (UDDRC_DRAMTMG5) This is the time before Self Refresh Exit that CK
 * is maintained as a valid clock before issuing SRX.
 * Specifies the clock stable time before SRX.
 * Recommended settings:
    - mDDR: 1
    - LPDDR2: 2
    - LPDDR3: 2
    - LPDDR4: tCKCKEH
    - DDR2: 1
    - DDR3: tCKSRX
    - DDR4: tCKSRX
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 */
#define UDDRC_DRAMTMG5_t_cksrx_Pos 24
#define UDDRC_DRAMTMG5_t_cksrx_Msk (0xfu << UDDRC_DRAMTMG5_t_cksrx_Pos)
#define UDDRC_DRAMTMG5_t_cksrx(value) \
	((UDDRC_DRAMTMG5_t_cksrx_Msk & ((value) << UDDRC_DRAMTMG5_t_cksrx_Pos)))

/* -------- UDDRC_DRAMTMG6 : (UDDRC_REGS Offset: 0x118)
 * SDRAM Timing Register 6 --------
 */
/* (UDDRC_DRAMTMG6) This is the time before Clock Stop Exit that CK is
 * maintained as a valid clock before issuing Clock Stop Exit.
 * Specifies the clock stable time before next command after Clock Stop Exit.
 * Recommended settings:
    - mDDR: 1
    - LPDDR2: tXP + 2
    - LPDDR3: tXP + 2
    - LPDDR4: tXP + 2
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3/LPDDR4
 * devices.
 */
#define UDDRC_DRAMTMG6_t_ckcsx_Pos 0
#define UDDRC_DRAMTMG6_t_ckcsx_Msk (0xfu << UDDRC_DRAMTMG6_t_ckcsx_Pos)
#define UDDRC_DRAMTMG6_t_ckcsx(value) \
	((UDDRC_DRAMTMG6_t_ckcsx_Msk & ((value) << UDDRC_DRAMTMG6_t_ckcsx_Pos)))

/* (UDDRC_DRAMTMG6) This is the time before Deep Power Down Exit that CK is
 * maintained as a valid clock before issuing DPDX.
 * Specifies the clock stable time before DPDX.
 * Recommended settings:
    - mDDR: 1
    - LPDDR2: 2
    - LPDDR3: 2
 *
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * This is only present for designs supporting mDDR or LPDDR2 devices.
 */
#define UDDRC_DRAMTMG6_t_ckdpdx_Pos 16
#define UDDRC_DRAMTMG6_t_ckdpdx_Msk (0xfu << UDDRC_DRAMTMG6_t_ckdpdx_Pos)
#define UDDRC_DRAMTMG6_t_ckdpdx(value) \
	((UDDRC_DRAMTMG6_t_ckdpdx_Msk & ((value) << UDDRC_DRAMTMG6_t_ckdpdx_Pos)))

/* (UDDRC_DRAMTMG6) This is the time after Deep Power Down Entry that CK
 * is maintained as a valid clock. Specifies the clock disable delay after DPDE. 
 * Recommended settings:
    - mDDR: 0
    - LPDDR2: 2
    - LPDDR3: 2
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3 devices.
 */
#define UDDRC_DRAMTMG6_t_ckdpde_Pos 24
#define UDDRC_DRAMTMG6_t_ckdpde_Msk (0xfu << UDDRC_DRAMTMG6_t_ckdpde_Pos)
#define UDDRC_DRAMTMG6_t_ckdpde(value) \
	((UDDRC_DRAMTMG6_t_ckdpde_Msk & ((value) << UDDRC_DRAMTMG6_t_ckdpde_Pos)))

/* -------- UDDRC_DRAMTMG7 : (UDDRC_REGS Offset: 0x11C)
 * SDRAM Timing Register 7 --------
 */
/* (UDDRC_DRAMTMG7) This is the time before Power Down Exit that CK
 * is maintained as a valid clock before issuing PDX.
 * Specifies the clock stable time before PDX.
 * Recommended settings:
    - mDDR: 0
    - LPDDR2: 2
    - LPDDR3: 2
    - LPDDR4: 2
 * When using DDR2/3/4 SDRAM, this register should be set to the same value
 * as DRAMTMG5.t_cksrx.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_DRAMTMG7_t_ckpdx_Pos 0
#define UDDRC_DRAMTMG7_t_ckpdx_Msk (0xfu << UDDRC_DRAMTMG7_t_ckpdx_Pos)
#define UDDRC_DRAMTMG7_t_ckpdx(value) \
	((UDDRC_DRAMTMG7_t_ckpdx_Msk & ((value) << UDDRC_DRAMTMG7_t_ckpdx_Pos)))

/* (UDDRC_DRAMTMG7) This is the time after Power Down Entry that CK is
 * maintained as a valid clock. Specifies the clock disable delay after PDE.
 * Recommended settings:
    - mDDR: 0
    - LPDDR2: 2
    - LPDDR3: 2
    - LPDDR4: tCKELCK
 * When using DDR2/3/4 SDRAM, this register should be set to the same value
 * as DRAMTMG5.t_cksre.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3/LPDDR4
 * devices.
 */
#define UDDRC_DRAMTMG7_t_ckpde_Pos 8
#define UDDRC_DRAMTMG7_t_ckpde_Msk (0xfu << UDDRC_DRAMTMG7_t_ckpde_Pos)
#define UDDRC_DRAMTMG7_t_ckpde(value) \
	((UDDRC_DRAMTMG7_t_ckpde_Msk & ((value) << UDDRC_DRAMTMG7_t_ckpde_Pos)))

/* -------- UDDRC_DRAMTMG8 : (UDDRC_REGS Offset: 0x120)
 * SDRAM Timing Register 8 --------
 */
/* (UDDRC_DRAMTMG8) tXS: Exit Self Refresh to commands not requiring a
 * locked DLL.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * the above value divided by 2 and round up to next integer value.
 * Unit: Multiples of 32 clocks.
 * Note: Used only for DDR2, DDR3 and DDR4 SDRAMs.
 */
#define UDDRC_DRAMTMG8_t_xs_x32_Pos 0
#define UDDRC_DRAMTMG8_t_xs_x32_Msk (0x7fu << UDDRC_DRAMTMG8_t_xs_x32_Pos)
#define UDDRC_DRAMTMG8_t_xs_x32(value) \
	((UDDRC_DRAMTMG8_t_xs_x32_Msk & ((value) << UDDRC_DRAMTMG8_t_xs_x32_Pos)))

/* (UDDRC_DRAMTMG8) tXSDLL: Exit Self Refresh to commands requiring a locked DLL.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * the above value divided by 2 and round up to next integer value.
 * Unit: Multiples of 32 clocks.
 * Note: Used only for DDR2, DDR3 and DDR4 SDRAMs.
 */
#define UDDRC_DRAMTMG8_t_xs_dll_x32_Pos 8
#define UDDRC_DRAMTMG8_t_xs_dll_x32_Msk \
	(0x7fu << UDDRC_DRAMTMG8_t_xs_dll_x32_Pos)
#define UDDRC_DRAMTMG8_t_xs_dll_x32(value) \
	((UDDRC_DRAMTMG8_t_xs_dll_x32_Msk & \
	((value) << UDDRC_DRAMTMG8_t_xs_dll_x32_Pos)))

/* -------- UDDRC_DRAMTMG14 : (UDDRC_REGS Offset: 0x138)
 * SDRAM Timing Register 14 --------
 */
/* (UDDRC_DRAMTMG14) tXSR: Exit Self Refresh to any command.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * the above value divided by 2 and round up to next integer value.
 * Note: Used only for mDDR/LPDDR2/LPDDR3/LPDDR4 mode.
 */
#define UDDRC_DRAMTMG14_t_xsr_Pos 0
#define UDDRC_DRAMTMG14_t_xsr_Msk (0xfffu << UDDRC_DRAMTMG14_t_xsr_Pos)
#define UDDRC_DRAMTMG14_t_xsr(value) \
	((UDDRC_DRAMTMG14_t_xsr_Msk & ((value) << UDDRC_DRAMTMG14_t_xsr_Pos)))


/* -------- UDDRC_DRAMTMG15 : (UDDRC_REGS Offset: 0x13C)
 * SDRAM Timing Register 15 --------
 */
/* (UDDRC_DRAMTMG15) tSTAB: Stabilization time.
 * It is required in the following two cases for DDR3/DDR4 RDIMM : 
   - when exiting power saving mode, if the clock was stopped, after re-enabling
	 it the clock must be stable for a time specified by tSTAB
   - in the case of input clock frequency change (DDR4)
   - after issuing control words that refers to clock timing
 * (Specification: 6us for DDR3, 5us for DDR4)
 *
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * recommended value divided by two and round it up to next integer.
 * Unit: Multiples of 32 clock cycles.
 */
#define UDDRC_DRAMTMG15_t_stab_x32_Pos 0
#define UDDRC_DRAMTMG15_t_stab_x32_Msk (0xffu << UDDRC_DRAMTMG15_t_stab_x32_Pos)
#define UDDRC_DRAMTMG15_t_stab_x32(value) \
	((UDDRC_DRAMTMG15_t_stab_x32_Msk & \
	((value) << UDDRC_DRAMTMG15_t_stab_x32_Pos)))

/* (UDDRC_DRAMTMG15)
  - 1 - Enable using tSTAB when exiting DFI LP. Needs to be set when the PHY is
	stopping the clock during DFI LP to save maximum power.
  - 0 - Disable using tSTAB when exiting DFI LP
 */
#define UDDRC_DRAMTMG15_en_dfi_lp_t_stab (0x1u << 31)

/* -------- UDDRC_ZQCTL0 : (UDDRC_REGS Offset: 0x180)
 * ZQ Control Register 0 --------
 */
/* (UDDRC_ZQCTL0) tZQCS for DDR3/DD4/LPDDR2/LPDDR3, tZQLAT for LPDDR4:
 * Number of DFI clock cycles of NOP required after a ZQCS
 * (ZQ calibration short)/MPC(ZQ Latch) command is issued to SDRAM.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * tZQCS/2 and round it up to the next integer value.
 * This is only present for designs supporting DDR3/DDR4 or LPDDR2/LPDDR3/LPDDR4
 * devices.
 */
#define UDDRC_ZQCTL0_t_zq_short_nop_Pos 0
#define UDDRC_ZQCTL0_t_zq_short_nop_Msk \
	(0x3ffu << UDDRC_ZQCTL0_t_zq_short_nop_Pos)
#define UDDRC_ZQCTL0_t_zq_short_nop(value) \
	((UDDRC_ZQCTL0_t_zq_short_nop_Msk & \
	((value) << UDDRC_ZQCTL0_t_zq_short_nop_Pos)))

/* (UDDRC_ZQCTL0) tZQoper for DDR3/DDR4, tZQCL for LPDDR2/LPDDR3,
 * tZQCAL for LPDDR4: Number of DFI clock cycles of NOP required after a
 * ZQCL (ZQ calibration long)/MPC(ZQ Start) command is issued to SDRAM.
 * When the controller is operating in 1:2 frequency ratio mode:
 * DDR3/DDR4: program this to tZQoper/2 and round it up to the next
 * integer value.
 * LPDDR2/LPDDR3: program this to tZQCL/2 and round it up to the next
 * integer value.
 * LPDDR4: program this to tZQCAL/2 and round it up to the next integer value.
 * This is only present for designs supporting DDR3/DDR4 or
 * LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL0_t_zq_long_nop_Pos 16
#define UDDRC_ZQCTL0_t_zq_long_nop_Msk \
	(0x7ffu << UDDRC_ZQCTL0_t_zq_long_nop_Pos)
#define UDDRC_ZQCTL0_t_zq_long_nop(value) \
	((UDDRC_ZQCTL0_t_zq_long_nop_Msk & \
	((value) << UDDRC_ZQCTL0_t_zq_long_nop_Pos)))

/* (UDDRC_ZQCTL0)
 * - 1 - Denotes that ZQ resistor is shared between ranks.
 * Means ZQinit/ZQCL/ZQCS/MPC(ZQ calibration) commands are sent to one rank
 * at a time with tZQinit/tZQCL/tZQCS/tZQCAL/tZQLAT timing met between commands
 * so that commands to different ranks do not overlap.
 * - 0 - ZQ resistor is not shared.
 * This is only present for designs supporting DDR3/DDR4 or
 * LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL0_zq_resistor_shared (0x1u << 29)

/* (UDDRC_ZQCTL0)
 * - 1 - Disable issuing of ZQCL/MPC(ZQ calibration) command
 * at Self-Refresh/SR-Powerdown exit. Only applicable when run in DDR3 or DDR4
 * or LPDDR2 or LPDDR3 or LPDDR4 mode.
 * - 0 - Enable issuing of ZQCL/MPC(ZQ calibration) command at
 * Self-Refresh/SR-Powerdown exit. Only applicable when run in DDR3 or DDR4
 * or LPDDR2 or LPDDR3 or LPDDR4 mode.
 * This is only present for designs supporting DDR3/DDR4 or
 * LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL0_dis_srx_zqcl (0x1u << 30)

/* (UDDRC_ZQCTL0)
 * - 1 - Disable uMCTL2 generation of ZQCS/MPC(ZQ calibration) command.
 * Register DBGCMD.zq_calib_short can be used instead to issue ZQ calibration
 * request from APB module.
 * - 0 - Internally generate ZQCS/MPC(ZQ calibration) commands based on
 * ZQCTL1.t_zq_short_interval_x1024.
 * This is only present for designs supporting DDR3/DDR4 or
 * LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL0_dis_auto_zq (0x1u << 31)

/* -------- UDDRC_ZQCTL1 : (UDDRC_REGS Offset: 0x184)
 * ZQ Control Register 1 --------
 */
/* (UDDRC_ZQCTL1) Average interval to wait between automatically issuing
 * ZQCS (ZQ calibration short)/MPC(ZQ calibration) commands to
 * DDR3/DDR4/LPDDR2/LPDDR3/LPDDR4 devices.
 * Meaningless, if ZQCTL0.dis_auto_zq=1.
 * Unit: 1024 DFI clock cycles.
 * This is only present for designs supporting DDR3/DDR4 or
 * LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL1_t_zq_short_interval_x1024_Pos 0
#define UDDRC_ZQCTL1_t_zq_short_interval_x1024_Msk \
	(0xfffffu << UDDRC_ZQCTL1_t_zq_short_interval_x1024_Pos)


#define UDDRC_ZQCTL1_t_zq_short_interval_x1024(value) \
	((UDDRC_ZQCTL1_t_zq_short_interval_x1024_Msk & \
	((value) << UDDRC_ZQCTL1_t_zq_short_interval_x1024_Pos)))

/* (UDDRC_ZQCTL1) tZQReset: Number of DFI clock cycles of NOP required after
 * a ZQReset (ZQ calibration Reset) command is issued to SDRAM.
 * When the controller is operating in 1:2 frequency ratio mode, program this to
 * tZQReset/2 and round it up to the next integer value.
 * This is only present for designs supporting LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL1_t_zq_reset_nop_Pos 20
#define UDDRC_ZQCTL1_t_zq_reset_nop_Msk \
	(0x3ffu << UDDRC_ZQCTL1_t_zq_reset_nop_Pos)
#define UDDRC_ZQCTL1_t_zq_reset_nop(value) \
	((UDDRC_ZQCTL1_t_zq_reset_nop_Msk & \
	((value) << UDDRC_ZQCTL1_t_zq_reset_nop_Pos)))

/* -------- UDDRC_ZQCTL2 : (UDDRC_REGS Offset: 0x188)
 * ZQ Control Register 2 --------
 */
/* (UDDRC_ZQCTL2)
 * Setting this register bit to 1 triggers a ZQ Reset operation.
 * When the ZQ Reset operation is complete, the uMCTL2 automatically clears
 * this bit. It is recommended NOT to set this signal if in Init,
 * Self-Refresh(except LPDDR4) or SR-Powerdown(LPDDR4) or Deep power-down
 * operating modes.
 * This is only present for designs supporting LPDDR2/LPDDR3/LPDDR4 devices.
 */
#define UDDRC_ZQCTL2_zq_reset (0x1u << 0)

/* -------- UDDRC_ZQSTAT : (UDDRC_REGS Offset: 0x18C)
 * ZQ Status Register --------
 */

/* (UDDRC_ZQSTAT) SoC core may initiate a ZQ Reset operation only if this signal
 * is low. This signal goes high in the clock after the uMCTL2 accepts the
 * ZQ Reset request. It goes low when the ZQ Reset command is issued to the
 * SDRAM and the associated NOP period is over. It is recommended not to
 * perform ZQ Reset commands when this signal is high.
 * - 0 - Indicates that the SoC core can initiate a ZQ Reset operation
 * - 1 - Indicates that ZQ Reset operation is in progress
 */
#define UDDRC_ZQSTAT_zq_reset_busy (0x1u << 0)
/* -------- UDDRC_DFITMG0 : (UDDRC_REGS Offset: 0x190)
 * DFI Timing Register 0 --------
 */

/* (UDDRC_DFITMG0) Write latency
 * Number of clocks from the write command to write data enable (dfi_wrdata_en).
 * This corresponds to the DFI timing parameter tphy_wrlat.
 * Refer to PHY specification for correct value.Note that, depending on the PHY,
 * if using RDIMM/LRDIMM, it may be necessary to use the adjusted value of CL
 * in the calculation of tphy_wrlat. This is to compensate for the extra
 * cycle(s) of latency through the RDIMM/LRDIMM.
 * For LPDDR4, dfi_tphy_wrlat>60 is not supported.
 * Unit: DFI clock cycles or DFI PHY clock cycles, depending on
 * DFITMG0.dfi_wrdata_use_dfi_phy_clk.
 */
#define UDDRC_DFITMG0_dfi_tphy_wrlat_Pos 0
#define UDDRC_DFITMG0_dfi_tphy_wrlat_Msk \
	(0x3fu << UDDRC_DFITMG0_dfi_tphy_wrlat_Pos)
#define UDDRC_DFITMG0_dfi_tphy_wrlat(value) \
	((UDDRC_DFITMG0_dfi_tphy_wrlat_Msk & \
	((value) << UDDRC_DFITMG0_dfi_tphy_wrlat_Pos)))

/* (UDDRC_DFITMG0) Specifies the number of clock cycles between when
 * dfi_wrdata_en is asserted to when the associated write data is driven
 * on the dfi_wrdata signal.  This corresponds to the DFI timing parameter
 * tphy_wrdata.  Refer to PHY specification for correct value.
 * Note, max supported value is 8.
 * Unit: DFI clock cycles or DFI PHY clock cycles, depending on
 * DFITMG0.dfi_wrdata_use_dfi_phy_clk.
 */
#define UDDRC_DFITMG0_dfi_tphy_wrdata_Pos 8
#define UDDRC_DFITMG0_dfi_tphy_wrdata_Msk \
	(0x3fu << UDDRC_DFITMG0_dfi_tphy_wrdata_Pos)
#define UDDRC_DFITMG0_dfi_tphy_wrdata(value) \
	((UDDRC_DFITMG0_dfi_tphy_wrdata_Msk & \
	((value) << UDDRC_DFITMG0_dfi_tphy_wrdata_Pos)))

/* (UDDRC_DFITMG0) Defines whether dfi_wrdata_en/dfi_wrdata/dfi_wrdata_mask is
 * generated using HDR (DFI clock) or SDR (DFI PHY clock) values
 * Selects whether value in DFITMG0.dfi_tphy_wrlat  is in terms of
 * HDR (DFI clock) or SDR (DFI PHY clock) cycles
 * Selects whether value in DFITMG0.dfi_tphy_wrdata is in terms of
 * HDR (DFI clock) or SDR (DFI PHY clock) cycles
 * - 0 in terms of HDR (DFI clock) cycles
 * - 1 in terms of SDR (DFI PHY clock) cycles
 * Refer to PHY specification for correct value.
 * If using a Synopsys DWC DDR3/2 PHY, DWC DDR2/3-Lite/mDDR PHY,
 * DWC DDR multiPHY or DWC Gen2 DDR multiPHY, this field must be set to 0;
 * otherwise:
 * - If MEMC_PROG_FREQ_RATIO=1 and MSTR.frequency_ratio=1,
	 this field should be set to 0
 * - Else, it must be set to 1
 */
#define UDDRC_DFITMG0_dfi_wrdata_use_dfi_phy_clk (0x1u << 15)

/* (UDDRC_DFITMG0) Time from the assertion of a read command on the
 * DFI interface to the assertion of the dfi_rddata_en signal.
 * Refer to PHY specification for correct value.
 * This corresponds to the DFI parameter trddata_en.
 * Note that, depending on the PHY, if using RDIMM/LRDIMM, it may be necessary
 * to use the adjusted value of CL in the calculation of trddata_en.
 * This is to compensate for the extra cycle(s) of latency through the
 * RDIMM/LRDIMM.
 * Unit: DFI clock cycles or DFI PHY clock cycles, depending on
 * DFITMG0.dfi_rddata_use_dfi_phy_clk.
 */
#define UDDRC_DFITMG0_dfi_t_rddata_en_Pos 16
#define UDDRC_DFITMG0_dfi_t_rddata_en_Msk \
	(0x7fu << UDDRC_DFITMG0_dfi_t_rddata_en_Pos)
#define UDDRC_DFITMG0_dfi_t_rddata_en(value) \
	((UDDRC_DFITMG0_dfi_t_rddata_en_Msk & \
	((value) << UDDRC_DFITMG0_dfi_t_rddata_en_Pos)))

/* (UDDRC_DFITMG0) Defines whether dfi_rddata_en/dfi_rddata/dfi_rddata_valid is
 * generated using HDR (DFI clock) or SDR (DFI PHY clock) values.
 * Selects whether value in DFITMG0.dfi_t_rddata_en is in terms of
 * HDR (DFI clock) or SDR (DFI PHY clock) cycles:
 * - 0 in terms of HDR (DFI clock) cycles
 * - 1 in terms of SDR (DFI PHY clock) cycles
 * Refer to PHY specification for correct value.
 * If using a Synopsys DWC DDR3/2 PHY, DWC DDR2/3-Lite/mDDR PHY,
 * DWC DDR multiPHY or DWC Gen2 DDR multiPHY, this field must be set to 0;
 * otherwise:
 * - If MEMC_PROG_FREQ_RATIO=1 and MSTR.frequency_ratio=1, this field should be
	set to 0
 * - Else, it must be set to 1
 */
#define UDDRC_DFITMG0_dfi_rddata_use_dfi_phy_clk (0x1u << 23)

/* (UDDRC_DFITMG0) Specifies the number of DFI clock cycles after an assertion
 * or de-assertion of the DFI control signals that the control signals at the
 * PHY-DRAM interface reflect the assertion or de-assertion.
 * If the DFI clock and the memory clock are not phase-aligned, this timing
 * parameter should be rounded up to the next integer value.
 * Note that if using RDIMM/LRDIMM, it is necessary to increment this parameter
 * by RDIMM's/LRDIMM's extra cycle of latency in terms of DFI clock.
 */
#define UDDRC_DFITMG0_dfi_t_ctrl_delay_Pos 24
#define UDDRC_DFITMG0_dfi_t_ctrl_delay_Msk \
	(0x1fu << UDDRC_DFITMG0_dfi_t_ctrl_delay_Pos)
#define UDDRC_DFITMG0_dfi_t_ctrl_delay(value) \
	((UDDRC_DFITMG0_dfi_t_ctrl_delay_Msk & \
	((value) << UDDRC_DFITMG0_dfi_t_ctrl_delay_Pos)))

/* -------- UDDRC_DFITMG1 : (UDDRC_REGS Offset: 0x194)
 * DFI Timing Register 1 --------
 */
/* (UDDRC_DFITMG1) Specifies the number of DFI clock cycles from the
 * de-assertion of the dfi_dram_clk_disable signal on the DFI until the first
 * valid rising edge of the clock to the DRAM memory devices,
 * at the PHY-DRAM boundary. If the DFI clock and the memory clock are not
 * phase aligned, this timing parameter should be rounded up to the next
 * integer value.
 */
#define UDDRC_DFITMG1_dfi_t_dram_clk_enable_Pos 0
#define UDDRC_DFITMG1_dfi_t_dram_clk_enable_Msk \
	(0x1fu << UDDRC_DFITMG1_dfi_t_dram_clk_enable_Pos)
#define UDDRC_DFITMG1_dfi_t_dram_clk_enable(value) \
	((UDDRC_DFITMG1_dfi_t_dram_clk_enable_Msk & \
	((value) << UDDRC_DFITMG1_dfi_t_dram_clk_enable_Pos)))

/* (UDDRC_DFITMG1) Specifies the number of DFI clock cycles from the assertion
 * of the dfi_dram_clk_disable signal on the DFI until the clock to the DRAM
 * memory devices, at the PHY-DRAM boundary, maintains a low value.
 * If the DFI clock and the memory clock are not phase aligned, this timing
 * parameter should be rounded up to the next integer value.
 */
#define UDDRC_DFITMG1_dfi_t_dram_clk_disable_Pos 8
#define UDDRC_DFITMG1_dfi_t_dram_clk_disable_Msk \
	(0x1fu << UDDRC_DFITMG1_dfi_t_dram_clk_disable_Pos)
#define UDDRC_DFITMG1_dfi_t_dram_clk_disable(value) \
	((UDDRC_DFITMG1_dfi_t_dram_clk_disable_Msk & \
	((value) << UDDRC_DFITMG1_dfi_t_dram_clk_disable_Pos)))

/* (UDDRC_DFITMG1) Specifies the number of DFI clock cycles between when the
 * dfi_wrdata_en signal is asserted and when the corresponding write data
 * transfer is completed on the DRAM bus.
 * This corresponds to the DFI timing parameter twrdata_delay.
 * Refer to PHY specification for correct value.
 * For DFI 3.0 PHY, set to twrdata_delay, a new timing parameter introduced
 * in DFI 3.0.
 * For DFI 2.1 PHY, set to tphy_wrdata + (delay of DFI write data to the DRAM).
 * Value to be programmed is in terms of DFI clocks, not PHY clocks.
 * In FREQ_RATIO=2, divide PHY's value by 2 and round up to next integer.
 * If using DFITMG0.dfi_wrdata_use_dfi_phy_clk=1, add 1 to the value.
 * Unit: Clocks
 */
#define UDDRC_DFITMG1_dfi_t_wrdata_delay_Pos 16
#define UDDRC_DFITMG1_dfi_t_wrdata_delay_Msk \
	(0x1fu << UDDRC_DFITMG1_dfi_t_wrdata_delay_Pos)
#define UDDRC_DFITMG1_dfi_t_wrdata_delay(value) \
	((UDDRC_DFITMG1_dfi_t_wrdata_delay_Msk & \
	((value) << UDDRC_DFITMG1_dfi_t_wrdata_delay_Pos)))

/* (UDDRC_DFITMG1) */
#define UDDRC_DFITMG1_dfi_t_parrin_lat_Pos 24
#define UDDRC_DFITMG1_dfi_t_parrin_lat_Msk \
	(0x3u << UDDRC_DFITMG1_dfi_t_parrin_lat_Pos)
#define UDDRC_DFITMG1_dfi_t_parrin_lat(value) \
	((UDDRC_DFITMG1_dfi_t_parrin_lat_Msk & \
	((value) << UDDRC_DFITMG1_dfi_t_parrin_lat_Pos)))

/* (UDDRC_DFITMG1) */
#define UDDRC_DFITMG1_dfi_t_cmd_lat_Pos 28
#define UDDRC_DFITMG1_dfi_t_cmd_lat_Msk \
	(0xfu << UDDRC_DFITMG1_dfi_t_cmd_lat_Pos)


#define UDDRC_DFITMG1_dfi_t_cmd_lat(value) \
	((UDDRC_DFITMG1_dfi_t_cmd_lat_Msk & \
	((value) << UDDRC_DFITMG1_dfi_t_cmd_lat_Pos)))

/* -------- UDDRC_DFILPCFG0 : (UDDRC_REGS Offset: 0x198)
 * DFI Low Power Configuration Register 0 --------
 */
/* (UDDRC_DFILPCFG0) Enables DFI Low Power interface handshaking
 * during Power Down Entry/Exit.
 * - 0 - Disabled
 * - 1 - Enabled
 */
#define UDDRC_DFILPCFG0_dfi_lp_en_pd (0x1u << 0)

/* (UDDRC_DFILPCFG0) Value in DFI clock cycles to drive on dfi_lp_wakeup signal
 * when Power Down mode is entered.
 * Determines the DFI's tlp_wakeup time:
  - 0x0 - 16 cycles
  - 0x1 - 32 cycles
  - 0x2 - 64 cycles
  - 0x3 - 128 cycles
  - 0x4 - 256 cycles
  - 0x5 - 512 cycles
  - 0x6 - 1024 cycles
  - 0x7 - 2048 cycles
  - 0x8 - 4096 cycles
  - 0x9 - 8192 cycles
  - 0xA - 16384 cycles
  - 0xB - 32768 cycles
  - 0xC - 65536 cycles
  - 0xD - 131072 cycles
  - 0xE - 262144 cycles
  - 0xF - Unlimited
*/
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_pd_Pos 4
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_pd_Msk \
	(0xfu << UDDRC_DFILPCFG0_dfi_lp_wakeup_pd_Pos)
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_pd(value) \
	((UDDRC_DFILPCFG0_dfi_lp_wakeup_pd_Msk & \
	((value) << UDDRC_DFILPCFG0_dfi_lp_wakeup_pd_Pos)))

/* (UDDRC_DFILPCFG0) Enables DFI Low Power interface handshaking during
 * Self Refresh Entry/Exit.
 * - 0 - Disabled
 *  - 1 - Enabled
 */

#define UDDRC_DFILPCFG0_dfi_lp_en_sr (0x1u << 8)
/* (UDDRC_DFILPCFG0) Value in DFI clpck cycles to drive on dfi_lp_wakeup signal
 * when Self Refresh mode is entered.
 * Determines the DFI's tlp_wakeup time:
  - 0x0 - 16 cycles
  - 0x1 - 32 cycles
  - 0x2 - 64 cycles
  - 0x3 - 128 cycles
  - 0x4 - 256 cycles
  - 0x5 - 512 cycles
  - 0x6 - 1024 cycles
  - 0x7 - 2048 cycles
  - 0x8 - 4096 cycles
  - 0x9 - 8192 cycles
  - 0xA - 16384 cycles
  - 0xB - 32768 cycles
  - 0xC - 65536 cycles
  - 0xD - 131072 cycles
  - 0xE - 262144 cycles
  - 0xF - Unlimited
 */
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_sr_Pos 12
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_sr_Msk \
	(0xfu << UDDRC_DFILPCFG0_dfi_lp_wakeup_sr_Pos)
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_sr(value) \
	((UDDRC_DFILPCFG0_dfi_lp_wakeup_sr_Msk & \
	((value) << UDDRC_DFILPCFG0_dfi_lp_wakeup_sr_Pos)))

/* (UDDRC_DFILPCFG0) Enables DFI Low Power interface handshaking during
 * Deep Power Down Entry/Exit.
  - 0 - Disabled
  - 1 - Enabled
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3 devices.
 */
#define UDDRC_DFILPCFG0_dfi_lp_en_dpd (0x1u << 16)

/* (UDDRC_DFILPCFG0) Value in DFI clock cycles to drive on dfi_lp_wakeup signal
 * when Deep Power Down mode is entered.
 * Determines the DFI's tlp_wakeup time:
  - 0x0 - 16 cycles
  - 0x1 - 32 cycles
  - 0x2 - 64 cycles
  - 0x3 - 128 cycles
  - 0x4 - 256 cycles
  - 0x5 - 512 cycles
  - 0x6 - 1024 cycles
  - 0x7 - 2048 cycles
  - 0x8 - 4096 cycles
  - 0x9 - 8192 cycles
  - 0xA - 16384 cycles
  - 0xB - 32768 cycles
  - 0xC - 65536 cycles
  - 0xD - 131072 cycles
  - 0xE - 262144 cycles
  - 0xF - Unlimited
 * This is only present for designs supporting mDDR or LPDDR2/LPDDR3 devices.
 */
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd_Pos 20
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd_Msk \
	(0xfu << UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd_Pos)
#define UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd(value) \
	((UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd_Msk & \
	((value) << UDDRC_DFILPCFG0_dfi_lp_wakeup_dpd_Pos)))

/* (UDDRC_DFILPCFG0) Setting in DFI clock cycles for DFI's tlp_resp time.
 * Same value is used for both Power Down, Self Refresh, Deep Power Down and
 * Maximum Power Saving modes.
 * DFI 2.1 specification onwards, recommends using a fixed value of 7 always.
 */
#define UDDRC_DFILPCFG0_dfi_tlp_resp_Pos 24
#define UDDRC_DFILPCFG0_dfi_tlp_resp_Msk \
	(0x1fu << UDDRC_DFILPCFG0_dfi_tlp_resp_Pos)
#define UDDRC_DFILPCFG0_dfi_tlp_resp(value) \
	((UDDRC_DFILPCFG0_dfi_tlp_resp_Msk & \
	((value) << UDDRC_DFILPCFG0_dfi_tlp_resp_Pos)))

/* -------- UDDRC_DFILPCFG1 : (UDDRC_REGS Offset: 0x19C)
 * DFI Low Power Configuration Register 1 --------
 */
#define UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm_Pos 4
#define UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm_Msk \
	(0xfu << UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm_Pos)
#define UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm(value) \
	((UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm_Msk & \
	((value) << UDDRC_DFILPCFG1_dfi_lp_wakeup_mpsm_Pos)))

#define UDDRC_DFILPCFG1_dfi_lp_en_mpsm (0x1u << 0)
/* -------- UDDRC_DFIUPD0 : (UDDRC_REGS Offset: 0x1A0)
 * DFI Update Register 0 --------
 */
/* (UDDRC_DFIUPD0) Specifies the minimum number of DFI clock cycles that the
 * dfi_ctrlupd_req signal must be asserted. The uMCTL2 expects the PHY to
 * respond within this time.  If the PHY does not respond, the uMCTL2 will
 * de-assert dfi_ctrlupd_req after dfi_t_ctrlup_min + 2 cycles.
 * Lowest value to assign to this variable is 0x3.
 */
#define UDDRC_DFIUPD0_dfi_t_ctrlup_min_Pos 0
#define UDDRC_DFIUPD0_dfi_t_ctrlup_min_Msk \
	(0x3ffu << UDDRC_DFIUPD0_dfi_t_ctrlup_min_Pos)
#define UDDRC_DFIUPD0_dfi_t_ctrlup_min(value) \
	((UDDRC_DFIUPD0_dfi_t_ctrlup_min_Msk & \
	((value) << UDDRC_DFIUPD0_dfi_t_ctrlup_min_Pos)))

/* (UDDRC_DFIUPD0) Specifies the maximum number of DFI clock cycles that the
 * dfi_ctrlupd_req signal can assert.
 * Lowest value to assign to this variable is 0x40.
 */
#define UDDRC_DFIUPD0_dfi_t_ctrlup_max_Pos 16
#define UDDRC_DFIUPD0_dfi_t_ctrlup_max_Msk \
	(0x3ffu << UDDRC_DFIUPD0_dfi_t_ctrlup_max_Pos)
#define UDDRC_DFIUPD0_dfi_t_ctrlup_max(value) \
	((UDDRC_DFIUPD0_dfi_t_ctrlup_max_Msk & \
	((value) << UDDRC_DFIUPD0_dfi_t_ctrlup_max_Pos)))

/* (UDDRC_DFIUPD0) Selects dfi_ctrlupd_req requirements at SRX:
 *  - 0 : send ctrlupd after SRX
 * - 1 : send ctrlupd before SRX
 * If DFIUPD0.dis_auto_ctrlupd_srx=1, this register has no impact,
 * because no dfi_ctrlupd_req will be issued when SRX.
 */
#define UDDRC_DFIUPD0_ctrlupd_pre_srx (0x1u << 29)

/* (UDDRC_DFIUPD0) When '1', disable the automatic dfi_ctrlupd_req generation
 * by the uMCTL2 at self-refresh exit.
 * When '0', uMCTL2 issues a dfi_ctrlupd_req before or after exiting
 * self-refresh,  depending on DFIUPD0.ctrlupd_pre_srx.
 */
#define UDDRC_DFIUPD0_dis_auto_ctrlupd_srx (0x1u << 30)

/* (UDDRC_DFIUPD0) When '1', disable the automatic dfi_ctrlupd_req generation
 * by the uMCTL2. The core must issue the dfi_ctrlupd_req signal using register
 * DBGCMD.ctrlupd.
 * When '0', uMCTL2 issues dfi_ctrlupd_req periodically.
 */
#define UDDRC_DFIUPD0_dis_auto_ctrlupd (0x1u << 31)

/* -------- UDDRC_DFIUPD1 : (UDDRC_REGS Offset: 0x1A4)
 * DFI Update Register 1 --------
 */
/* (UDDRC_DFIUPD1) This is the maximum amount of time between uMCTL2 initiated
 * DFI update requests. This timer resets with each update request;
 * when the timer expires dfi_ctrlupd_req is sent and traffic is blocked
 * until the dfi_ctrlupd_ackx is received. PHY can use this idle time to
 * recalibrate the delay lines to the DLLs. The DFI controller update is also
 * used to reset PHY FIFO pointers in case of data capture errors.
 * Updates are required to maintain calibration over PVT, but frequent updates
 * may impact performance. Minimum allowed value for this field is 1.
 * Note: Value programmed for DFIUPD1.dfi_t_ctrlupd_interval_max_x1024 must be
 * greater than DFIUPD1.dfi_t_ctrlupd_interval_min_x1024.
 * Unit: 1024 DFI clock cycles
 */
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024_Pos 0
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024_Msk \
	(0xffu << UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024_Pos)
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024(value) \
	((UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024_Msk & \
	((value) << UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_max_x1024_Pos)))

/* (UDDRC_DFIUPD1) This is the minimum amount of time between uMCTL2 initiated
 * DFI update requests (which is executed whenever the uMCTL2 is idle).
 * Set this number higher to reduce the frequency of update requests,
 * which can have a small impact on the latency of the first read request
 * when the uMCTL2 is idle. Minimum allowed value for this field is 1.
 * Unit: 1024 DFI clock cycles
 */
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024_Pos 16
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024_Msk \
	(0xffu << UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024_Pos)
#define UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024(value) \
	((UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024_Msk & \
	((value) << UDDRC_DFIUPD1_dfi_t_ctrlupd_interval_min_x1024_Pos)))

/* -------- UDDRC_DFIUPD2 : (UDDRC_REGS Offset: 0x1A8)
 * DFI Update Register 2 --------
 */
/* (UDDRC_DFIUPD2) Enables the support for acknowledging PHY-initiated updates:
 *    - 0 - Disabled
 *   - 1 - Enabled
 */
#define UDDRC_DFIUPD2_dfi_phyupd_en (0x1u << 31)
/* -------- UDDRC_DFIMISC : (UDDRC_REGS Offset: 0x1B0)
 * DFI Miscellaneous Control Register --------
 */
/* (UDDRC_DFIMISC) PHY initialization complete enable signal.
 * When asserted the dfi_init_complete signal can be used to trigger
 * SDRAM initialisation
 */
#define UDDRC_DFIMISC_dfi_init_complete_en (0x1u << 0)

/* (UDDRC_DFIMISC) Enables support of ctl_idle signal, which is non-DFI related
 * pin specific to certain Synopsys PHYs. See signal description of ctl_idle
 * signal for further details of ctl_idle functionality.
 */
#define UDDRC_DFIMISC_ctl_idle_en (0x1u << 4)

/* (UDDRC_DFIMISC) PHY init start request signal.When asserted it triggers
 * the PHY init start request
 */
#define UDDRC_DFIMISC_dfi_init_start (0x1u << 5)

/* (UDDRC_DFIMISC) Indicates the operating frequency of the system.
 * The number of supported frequencies and the mapping of signal values to clock
 * frequencies are defined by the PHY.
 */
#define UDDRC_DFIMISC_dfi_frequency_Pos 8
#define UDDRC_DFIMISC_dfi_frequency_Msk \
	(0x1fu << UDDRC_DFIMISC_dfi_frequency_Pos)
#define UDDRC_DFIMISC_dfi_frequency(value) \
	((UDDRC_DFIMISC_dfi_frequency_Msk & \
	((value) << UDDRC_DFIMISC_dfi_frequency_Pos)))

/* -------- UDDRC_DFITMG2 : (UDDRC_REGS Offset: 0x1B4)
 * DFI Timing Register 2 --------
 */
/* (UDDRC_DFITMG2) Number of DFI PHY clock cycles between when a write command
 * is sent on the DFI control interface and when the associated dfi_wrdata_cs
 * signal is asserted.
 * This corresponds to the DFI timing parameter tphy_wrcslat.
 */
#define UDDRC_DFITMG2_dfi_tphy_wrcslat_Pos 0
#define UDDRC_DFITMG2_dfi_tphy_wrcslat_Msk \
	(0x1fu << UDDRC_DFITMG2_dfi_tphy_wrcslat_Pos)
#define UDDRC_DFITMG2_dfi_tphy_wrcslat(value) \
	((UDDRC_DFITMG2_dfi_tphy_wrcslat_Msk & \
	((value) << UDDRC_DFITMG2_dfi_tphy_wrcslat_Pos)))

/* (UDDRC_DFITMG2) Number of DFI PHY clock cycles between when a read command
 * is sent on the DFI control interface and when the associated
 * dfi_rddata_cs signal is asserted.
 * This corresponds to the DFI timing parameter tphy_rdcslat.
 * Refer to PHY specification for correct value.
 */
#define UDDRC_DFITMG2_dfi_tphy_rdcslat_Pos 8
#define UDDRC_DFITMG2_dfi_tphy_rdcslat_Msk \
	(0x7fu << UDDRC_DFITMG2_dfi_tphy_rdcslat_Pos)
#define UDDRC_DFITMG2_dfi_tphy_rdcslat(value) \
	((UDDRC_DFITMG2_dfi_tphy_rdcslat_Msk & \
	((value) << UDDRC_DFITMG2_dfi_tphy_rdcslat_Pos)))

/* -------- UDDRC_DFISTAT : (UDDRC_REGS Offset: 0x1BC)
 * DFI Status Register --------
 */
/* (UDDRC_DFISTAT) The status flag register which announces when the DFI
 * initialization has been completed. The DFI INIT triggered by dfi_init_start
 * signal and then the dfi_init_complete flag is polled to know when the
 * initialization is done.
 */
#define UDDRC_DFISTAT_dfi_init_complete (0x1u << 0)

/* (UDDRC_DFISTAT) Stores the value of the dfi_lp_ack input to the controller.
 */
#define UDDRC_DFISTAT_dfi_lp_ack (0x1u << 1)

/* -------- UDDRC_DFIPHYMSTR : (UDDRC_REGS Offset: 0x1C4)
 * DFI PHY Master --------
 */
/* (UDDRC_DFIPHYMSTR) Enables the PHY Master Interface:
 *   - 0 - Disabled
 *   - 1 - Enabled
 */
#define UDDRC_DFIPHYMSTR_dfi_phymstr_en (0x1u << 0)

/* -------- UDDRC_ADDRMAP1 : (UDDRC_REGS Offset: 0x204)
 * Address Map Register 1 --------
 */
/* (UDDRC_ADDRMAP1) Selects the HIF address bits used as bank address bit 0.
 * Valid Range: 0 to 32 and 63
 * Internal Base: 2
 * The selected HIF address bit for each of the bank address bits is determined
 * by adding the internal base to the value of this field.
 * If unused, set to 63 and then bank address bit 0 is set to 0.
 */
#define UDDRC_ADDRMAP1_addrmap_bank_b0_Pos 0
#define UDDRC_ADDRMAP1_addrmap_bank_b0_Msk \
	(0x3fu << UDDRC_ADDRMAP1_addrmap_bank_b0_Pos)
#define UDDRC_ADDRMAP1_addrmap_bank_b0(value) \
	((UDDRC_ADDRMAP1_addrmap_bank_b0_Msk & \
	((value) << UDDRC_ADDRMAP1_addrmap_bank_b0_Pos)))

/* (UDDRC_ADDRMAP1) Selects the HIF address bits used as bank address bit 1.
 * Valid Range: 0 to 32 and 63
 * Internal Base: 3
 * The selected HIF address bit for each of the bank address bits is determined
 * by adding the internal base to the value of this field.
 * If unused, set to 63 and then bank address bit 1 is set to 0.
 */
#define UDDRC_ADDRMAP1_addrmap_bank_b1_Pos 8
#define UDDRC_ADDRMAP1_addrmap_bank_b1_Msk \
	(0x3fu << UDDRC_ADDRMAP1_addrmap_bank_b1_Pos)
#define UDDRC_ADDRMAP1_addrmap_bank_b1(value) \
	((UDDRC_ADDRMAP1_addrmap_bank_b1_Msk & \
	((value) << UDDRC_ADDRMAP1_addrmap_bank_b1_Pos)))

/* (UDDRC_ADDRMAP1) Selects the HIF address bit used as bank address bit 2.
 * Valid Range: 0 to 31 and 63
 * Internal Base: 4
 * The selected HIF address bit is determined by adding the internal base to
 * the value of this field.
 * If unused, set to 63 and then bank address bit 2 is set to 0.
 */
#define UDDRC_ADDRMAP1_addrmap_bank_b2_Pos 16
#define UDDRC_ADDRMAP1_addrmap_bank_b2_Msk \
	(0x3fu << UDDRC_ADDRMAP1_addrmap_bank_b2_Pos)
#define UDDRC_ADDRMAP1_addrmap_bank_b2(value) \
	((UDDRC_ADDRMAP1_addrmap_bank_b2_Msk & \
	((value) << UDDRC_ADDRMAP1_addrmap_bank_b2_Pos)))

/* -------- UDDRC_ADDRMAP2 : (UDDRC_REGS Offset: 0x208)
 * Address Map Register 2 --------
 */
/* (UDDRC_ADDRMAP2)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	 bit 2.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	 bit 3.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	 bit 4.
 * Valid Range: 0 to 7
 * Internal Base: 2
 * The selected HIF address bit is determined by adding the internal base to
 * the value of this field.
 * Note, if UMCTL2_INCL_ARB=1 and MEMC_BURST_LENGTH=8, it is required to program
 * this to 0 unless:
  - in Half or Quarter bus width (MSTR.data_bus_width!=00) and
  - PCCFG.bl_exp_mode==1 and either
  - In DDR4   and ADDRMAP8.addrmap_bg_b0==0 or
  - In LPDDR4 and ADDRMAP1.addrmap_bank_b0==0
 * If UMCTL2_INCL_ARB=1 and MEMC_BURST_LENGTH=16, it is required to program this
 * to 0 unless:
  - in Half or Quarter bus width (MSTR.data_bus_width!=00) and
  - PCCFG.bl_exp_mode==1 and
  - In DDR4 and ADDRMAP8.addrmap_bg_b0==0
 * Otherwise, if MEMC_BURST_LENGTH=8 and Full Bus Width
 * (MSTR.data_bus_width==00), it is recommended to program this to 0 so that
 * HIF[2] maps to column address bit 2.
 *
 * If MEMC_BURST_LENGTH=16 and Full Bus Width (MSTR.data_bus_width==00), it is
 * recommended to program this to 0 so that HIF[2] maps to column address bit 2.
 *
 * If MEMC_BURST_LENGTH=16 and Half Bus Width (MSTR.data_bus_width==01), it is
 * recommended to program this to 0 so that HIF[2] maps to column address bit 3.
 */
#define UDDRC_ADDRMAP2_addrmap_col_b2_Pos 0
#define UDDRC_ADDRMAP2_addrmap_col_b2_Msk \
	(0xfu << UDDRC_ADDRMAP2_addrmap_col_b2_Pos)
#define UDDRC_ADDRMAP2_addrmap_col_b2(value) \
	((UDDRC_ADDRMAP2_addrmap_col_b2_Msk & \
	((value) << UDDRC_ADDRMAP2_addrmap_col_b2_Pos)))


#define UDDRC_ADDRMAP2_addrmap_col_b3_Pos 8
#define UDDRC_ADDRMAP2_addrmap_col_b3_Msk \
	(0xfu << UDDRC_ADDRMAP2_addrmap_col_b3_Pos)
/* (UDDRC_ADDRMAP2)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 3.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 4.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 5.
 * Valid Range: 0 to 7
 * Internal Base: 3
 * The selected HIF address bit is determined by adding the internal base to
 * the value of this field.
 * Note, if UMCTL2_INCL_ARB=1, MEMC_BURST_LENGTH=16, Full bus width
 * (MSTR.data_bus_width=00) and BL16 (MSTR.burst_rdwr=1000),
 * it is recommended to program this to 0.
 */
#define UDDRC_ADDRMAP2_addrmap_col_b3(value) \
	((UDDRC_ADDRMAP2_addrmap_col_b3_Msk & \
	((value) << UDDRC_ADDRMAP2_addrmap_col_b3_Pos)))

/**< \brief (UDDRC_ADDRMAP2)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 4.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 5.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 6.
 * Valid Range: 0 to 7, and 15
 * Internal Base: 4
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field. If unused, set to 15 and then this column address bit is
 * set to 0.
 */
#define UDDRC_ADDRMAP2_addrmap_col_b4_Pos 16
#define UDDRC_ADDRMAP2_addrmap_col_b4_Msk \
	(0xfu << UDDRC_ADDRMAP2_addrmap_col_b4_Pos)
#define UDDRC_ADDRMAP2_addrmap_col_b4(value) \
	((UDDRC_ADDRMAP2_addrmap_col_b4_Msk & \
	((value) << UDDRC_ADDRMAP2_addrmap_col_b4_Pos)))

/* (UDDRC_ADDRMAP2)
 * - Full bus width mode: Selects the HIF address bit used as column address
	bit 5.
 * - Half bus width mode: Selects the HIF address bit used as column address
	bit 6.
 * - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 7 .
 * Valid Range: 0 to 7, and 15
 * Internal Base: 5
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field. If unused, set to 15 and then this column address bit is
 * set to 0.
 */
#define UDDRC_ADDRMAP2_addrmap_col_b5_Pos 24
#define UDDRC_ADDRMAP2_addrmap_col_b5_Msk \
	(0xfu << UDDRC_ADDRMAP2_addrmap_col_b5_Pos)
#define UDDRC_ADDRMAP2_addrmap_col_b5(value) \
	((UDDRC_ADDRMAP2_addrmap_col_b5_Msk & \
	((value) << UDDRC_ADDRMAP2_addrmap_col_b5_Pos)))

/* -------- UDDRC_ADDRMAP3 : (UDDRC_REGS Offset: 0x20C)
 * Address Map Register 3 --------
 */
/* (UDDRC_ADDRMAP3)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 6.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 7.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 8.
 * Valid Range: 0 to 7, and 15
 * Internal Base: 6
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field. If unused, set to 15 and then this column address bit is
 * set to 0.
 */
#define UDDRC_ADDRMAP3_addrmap_col_b6_Pos 0
#define UDDRC_ADDRMAP3_addrmap_col_b6_Msk \
	(0xfu << UDDRC_ADDRMAP3_addrmap_col_b6_Pos)
#define UDDRC_ADDRMAP3_addrmap_col_b6(value) \
	((UDDRC_ADDRMAP3_addrmap_col_b6_Msk & \
	((value) << UDDRC_ADDRMAP3_addrmap_col_b6_Pos)))

/* (UDDRC_ADDRMAP3)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 7.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 8.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 9.
 * Valid Range: 0 to 7, x, and 31. x indicate a valid value in inline
 * ECC configuration.
 * Internal Base: 7
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field. If unused, set to 31 and then this column address bit is
 * set to 0.
 *
 * In Inline ECC configuration (MEMC_INLINE_ECC=1) and ECC is enabled
 * (ECCCFG0.ecc_mode>0), the highest 3 column address bits must map to the
 * highest 3 valid HIF address bits.
 *
 * If column bit 7 is the third highest column address bit, it must map to the
 * third highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 2 - internal base)
 *
 * if it is unused, set to 31.
 */
#define UDDRC_ADDRMAP3_addrmap_col_b7_Pos 8
#define UDDRC_ADDRMAP3_addrmap_col_b7_Msk \
	(0x1fu << UDDRC_ADDRMAP3_addrmap_col_b7_Pos)
#define UDDRC_ADDRMAP3_addrmap_col_b7(value) \
	((UDDRC_ADDRMAP3_addrmap_col_b7_Msk & \
	((value) << UDDRC_ADDRMAP3_addrmap_col_b7_Pos)))

/* (UDDRC_ADDRMAP3)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 8.
 *  - Half bus width mode: Selects the HIF address bit used as column address
 	bit 9.
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 11 (10 in LPDDR2/LPDDR3 mode).
 * Valid Range: 0 to 7, x, and 31. x indicate a valid value in inline ECC
 * configuration.
 * Internal Base: 8
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 31 and then this column address bit is set to 0.
 *
 * Note: Per JEDEC DDR2/3/mDDR specification, column address bit 10 is reserved
 * for indicating auto-precharge, and hence no source address bit can be mapped
 * to column address bit 10.
 *
 * In LPDDR2/LPDDR3, there is a dedicated bit for auto-precharge in the CA bus
 * and hence column bit 10 is used.
 *
 * In Inline ECC configuration (MEMC_INLINE_ECC=1) and ECC is enabled
 * (ECCCFG0.ecc_mode>0), the highest 3 column address bits must map to the
 * highest 3 valid HIF address bits.
 *
 * If column bit 8 is the second highest column address bit, it must map to the
 * second highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 1 - internal base)
 *
 * If column bit 8 is the third highest column address bit, it must map to the
 * third highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 2 - internal base)
 *
 * if it is unused, set to 31.
 */
#define UDDRC_ADDRMAP3_addrmap_col_b8_Pos 16
#define UDDRC_ADDRMAP3_addrmap_col_b8_Msk \
	(0x1fu << UDDRC_ADDRMAP3_addrmap_col_b8_Pos)
#define UDDRC_ADDRMAP3_addrmap_col_b8(value) \
	((UDDRC_ADDRMAP3_addrmap_col_b8_Msk & \
	((value) << UDDRC_ADDRMAP3_addrmap_col_b8_Pos)))
/* (UDDRC_ADDRMAP3)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 9.
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 11 (10 in LPDDR2/LPDDR3 mode).
 *  - Quarter bus width mode: Selects the HIF address bit used as column address
	bit 13 (11 in LPDDR2/LPDDR3 mode).
 * Valid Range: 0 to 7, x, and 31. x indicate a valid value in inline ECC
 * configuration.
 * Internal Base: 9
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 31 and then this column address bit is set to 0.
 *
 * Note: Per JEDEC DDR2/3/mDDR specification, column address bit 10 is reserved
 * for indicating auto-precharge, and hence no source address bit can be mapped
 * to column address bit 10.
 *
 * In LPDDR2/LPDDR3, there is a dedicated bit for auto-precharge in the CA bus
 * and hence column bit 10 is used.
 *
 * In Inline ECC configuration (MEMC_INLINE_ECC=1) and ECC is enabled
 * (ECCCFG0.ecc_mode>0), the highest 3 column address bits must map to the
 * highest 3 valid HIF address bits.
 *
 * If column bit 9 is the highest column address bit, it must map to the highest
 * valid HIF address bit. (x = the highest valid HIF address bit - internal base)

 * If column bit 9 is the second highest column address bit, it must map to the
 * second highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 1 - internal base)
 *
 * If column bit 9 is the third highest column address bit, it must map to the
 * third highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 2 - internal base)
 *
 * if it is unused, set to 31.
 */
#define UDDRC_ADDRMAP3_addrmap_col_b9_Pos 24
#define UDDRC_ADDRMAP3_addrmap_col_b9_Msk \
	(0x1fu << UDDRC_ADDRMAP3_addrmap_col_b9_Pos)
#define UDDRC_ADDRMAP3_addrmap_col_b9(value) \
	((UDDRC_ADDRMAP3_addrmap_col_b9_Msk & \
	((value) << UDDRC_ADDRMAP3_addrmap_col_b9_Pos)))

/* -------- UDDRC_ADDRMAP4 : (UDDRC_REGS Offset: 0x210)
 * Address Map Register 4 --------
 */
/* (UDDRC_ADDRMAP4)
 *  - Full bus width mode: Selects the HIF address bit used as column address
	bit 11 (10 in LPDDR2/LPDDR3 mode).
 *  - Half bus width mode: Selects the HIF address bit used as column address
	bit 13 (11 in LPDDR2/LPDDR3 mode).
 *  - Quarter bus width mode: UNUSED. To make it unused, this must be tied to
	4'hF.
 * Valid Range: 0 to 7, x, and 31. x indicate a valid value in inline ECC
	configuration.
 * Internal Base: 10
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 31 and then this column address bit is set to 0.
 *
 * Note: Per JEDEC DDR2/3/mDDR specification, column address bit 10 is reserved
 * for indicating auto-precharge, and hence no source address bit can be mapped
 * to column address bit 10.
 * In LPDDR2/LPDDR3, there is a dedicated bit for auto-precharge in the CA bus
 * and hence column bit 10 is used.
 *
 * In Inline ECC configuration (MEMC_INLINE_ECC=1) and ECC is enabled
 * (ECCCFG0.ecc_mode>0), the highest 3 column address bits must map to the
 * highest 3 valid HIF address bits.
 *
 * If column bit 10 is the highest column address bit, it must map to the
 * highest valid HIF address bit.
 * (x = the highest valid HIF address bit - internal base)
 *
 * If column bit 10 is the second highest column address bit, it must map to the
 * second highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 1 - internal base)
 *
 * If column bit 10 is the third highest column address bit, it must map to the
 * third highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 2 - internal base)
 *
 * if it is unused, set to 31.
 */
#define UDDRC_ADDRMAP4_addrmap_col_b10_Pos 0
#define UDDRC_ADDRMAP4_addrmap_col_b10_Msk \
	(0x1fu << UDDRC_ADDRMAP4_addrmap_col_b10_Pos)
#define UDDRC_ADDRMAP4_addrmap_col_b10(value) \
	((UDDRC_ADDRMAP4_addrmap_col_b10_Msk & \
	((value) << UDDRC_ADDRMAP4_addrmap_col_b10_Pos)))

/* (UDDRC_ADDRMAP4)
 * - Full bus width mode: Selects the HIF address bit used as column address
	bit 13 (11 in LPDDR2/LPDDR3 mode).
 * - Half bus width mode: Unused. To make it unused, this should be tied
	to 4'hF.
 * - Quarter bus width mode: Unused. To make it unused, this must be tied
	to 4'hF.
 * Valid Range: 0 to 7, x, and 31. x indicate a valid value in inline ECC
 * configuration.
 * Internal Base: 11
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 31 and then this column address bit is set to 0.
 *
 * Note: Per JEDEC DDR2/3/mDDR specification, column address bit 10 is reserved
 * for indicating auto-precharge, and hence no source address bit can be mapped
 * to column address bit 10.
 * In LPDDR2/LPDDR3, there is a dedicated bit for auto-precharge in the CA bus
 * and hence column bit 10 is used.
 *
 * In Inline ECC configuration (MEMC_INLINE_ECC=1) and ECC is enabled
 * (ECCCFG0.ecc_mode>0), the highest 3 column address bits must map to the
 * highest 3 valid HIF address bits.
 *
 * If column bit 11 is the highest column address bit, it must map to the
 * highest valid HIF address bit.
 * (x = the highest valid HIF address bit - internal base)
 *
 * If column bit 11 is the second highest column address bit, it must map to the
 * second highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 1 - internal base)
 *
 * If column bit 11 is the third highest column address bit, it must map to the
 * third highest valid HIF address bit.
 * (x = the highest valid HIF address bit - 2 - internal base)
 *
 * if it is unused, set to 31.
 */
#define UDDRC_ADDRMAP4_addrmap_col_b11_Pos 8
#define UDDRC_ADDRMAP4_addrmap_col_b11_Msk \
	(0x1fu << UDDRC_ADDRMAP4_addrmap_col_b11_Pos)
#define UDDRC_ADDRMAP4_addrmap_col_b11(value) \
	((UDDRC_ADDRMAP4_addrmap_col_b11_Msk & \
	((value) << UDDRC_ADDRMAP4_addrmap_col_b11_Pos)))

/* -------- UDDRC_ADDRMAP5 : (UDDRC_REGS Offset: 0x214)
 * Address Map Register 5 --------
 */
/* (UDDRC_ADDRMAP5) Selects the HIF address bits used as row address bit 0.
 * Valid Range: 0 to 11
 * Internal Base: 6
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field.
 */
#define UDDRC_ADDRMAP5_addrmap_row_b0_Pos 0
#define UDDRC_ADDRMAP5_addrmap_row_b0_Msk \
	(0xfu << UDDRC_ADDRMAP5_addrmap_row_b0_Pos)
#define UDDRC_ADDRMAP5_addrmap_row_b0(value) \
	((UDDRC_ADDRMAP5_addrmap_row_b0_Msk & \
	((value) << UDDRC_ADDRMAP5_addrmap_row_b0_Pos)))

/* (UDDRC_ADDRMAP5) Selects the HIF address bits used as row address bit 1.
 * Valid Range: 0 to 11
 * Internal Base: 7
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field.
 */
#define UDDRC_ADDRMAP5_addrmap_row_b1_Pos 8
#define UDDRC_ADDRMAP5_addrmap_row_b1_Msk \
	(0xfu << UDDRC_ADDRMAP5_addrmap_row_b1_Pos)
#define UDDRC_ADDRMAP5_addrmap_row_b1(value) \
	((UDDRC_ADDRMAP5_addrmap_row_b1_Msk & \
	((value) << UDDRC_ADDRMAP5_addrmap_row_b1_Pos)))

/* (UDDRC_ADDRMAP5) Selects the HIF address bits used as row address
 * bits 2 to 10.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 8 (for row address bit 2), 9 (for row address bit 3),
 * 10 (for row address bit 4) etc increasing to 16 (for row address bit 10)
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. When value 15 is used
 * the values of row address bits 2 to 10 are defined by registers
 * ADDRMAP9, ADDRMAP10, ADDRMAP11.
 */
#define UDDRC_ADDRMAP5_addrmap_row_b2_10_Pos 16
#define UDDRC_ADDRMAP5_addrmap_row_b2_10_Msk \
	(0xfu << UDDRC_ADDRMAP5_addrmap_row_b2_10_Pos)
#define UDDRC_ADDRMAP5_addrmap_row_b2_10(value) \
	((UDDRC_ADDRMAP5_addrmap_row_b2_10_Msk & \
	((value) << UDDRC_ADDRMAP5_addrmap_row_b2_10_Pos)))

/* (UDDRC_ADDRMAP5) Selects the HIF address bit used as row address bit 11.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 17
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 15 and then row address bit 11 is set to 0.
 */
#define UDDRC_ADDRMAP5_addrmap_row_b11_Pos 24
#define UDDRC_ADDRMAP5_addrmap_row_b11_Msk \
	(0xfu << UDDRC_ADDRMAP5_addrmap_row_b11_Pos)
#define UDDRC_ADDRMAP5_addrmap_row_b11(value) \
	((UDDRC_ADDRMAP5_addrmap_row_b11_Msk & \
	((value) << UDDRC_ADDRMAP5_addrmap_row_b11_Pos)))

/* -------- UDDRC_ADDRMAP6 : (UDDRC_REGS Offset: 0x218)
 Address Map Register 6 --------
 */
/* (UDDRC_ADDRMAP6) Selects the HIF address bit used as row address bit 12.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 18
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 15 and then row address bit 12 is set to 0.
 */
#define UDDRC_ADDRMAP6_addrmap_row_b12_Pos 0
#define UDDRC_ADDRMAP6_addrmap_row_b12_Msk \
	(0xfu << UDDRC_ADDRMAP6_addrmap_row_b12_Pos)
#define UDDRC_ADDRMAP6_addrmap_row_b12(value) \
	((UDDRC_ADDRMAP6_addrmap_row_b12_Msk & \
	((value) << UDDRC_ADDRMAP6_addrmap_row_b12_Pos)))

/* (UDDRC_ADDRMAP6) Selects the HIF address bit used as row address bit 13.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 19
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 15 and then row address bit 13 is set to 0.
 */
#define UDDRC_ADDRMAP6_addrmap_row_b13_Pos 8
#define UDDRC_ADDRMAP6_addrmap_row_b13_Msk \
	(0xfu << UDDRC_ADDRMAP6_addrmap_row_b13_Pos)
#define UDDRC_ADDRMAP6_addrmap_row_b13(value) \
	((UDDRC_ADDRMAP6_addrmap_row_b13_Msk & \
	((value) << UDDRC_ADDRMAP6_addrmap_row_b13_Pos)))

/* (UDDRC_ADDRMAP6) Selects the HIF address bit used as row address bit 14.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 20
 * The selected HIF address bit is determined by adding the internal base to the
 * value of this field.
 * If unused, set to 15 and then row address bit 14 is set to 0.
 */
#define UDDRC_ADDRMAP6_addrmap_row_b14_Pos 16
#define UDDRC_ADDRMAP6_addrmap_row_b14_Msk \
	(0xfu << UDDRC_ADDRMAP6_addrmap_row_b14_Pos)
#define UDDRC_ADDRMAP6_addrmap_row_b14(value) \
	((UDDRC_ADDRMAP6_addrmap_row_b14_Msk & \
	((value) << UDDRC_ADDRMAP6_addrmap_row_b14_Pos)))
/* (UDDRC_ADDRMAP6) Selects the HIF address bit used as row address bit 15.
 * Valid Range: 0 to 11, and 15
 * Internal Base: 21
 * The selected HIF address bit is determined by adding the internal base
 * to the value of this field.
 * If unused, set to 15 and then row address bit 15 is set to 0.
 */
#define UDDRC_ADDRMAP6_addrmap_row_b15_Pos 24
#define UDDRC_ADDRMAP6_addrmap_row_b15_Msk \
	(0xfu << UDDRC_ADDRMAP6_addrmap_row_b15_Pos)
#define UDDRC_ADDRMAP6_addrmap_row_b15(value) \
	((UDDRC_ADDRMAP6_addrmap_row_b15_Msk & \
	((value) << UDDRC_ADDRMAP6_addrmap_row_b15_Pos)))

/* (UDDRC_ADDRMAP6)
 * Set this to 1 if there is an LPDDR3 SDRAM 6Gb or 12Gb device in use.
 * - 1 - LPDDR3 SDRAM 6Gb/12Gb device in use. Every address having
 * row[14:13]==2'b11 is considered as invalid
 * - 0 - non-LPDDR3 6Gb/12Gb device in use. All addresses are valid
 * Present only in designs configured to support LPDDR3.
 */
#define UDDRC_ADDRMAP6_lpddr3_6gb_12gb (0x1u << 31)

/* -------- UDDRC_ADDRMAP9 : (UDDRC_REGS Offset: 0x224)
 * Address Map Register 9 --------
 */
/* (UDDRC_ADDRMAP9) Selects the HIF address bits used as row address bit 2.
 * Valid Range: 0 to 11
 * Internal Base: 8
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP9_addrmap_row_b2_Pos 0
#define UDDRC_ADDRMAP9_addrmap_row_b2_Msk \
	(0xfu << UDDRC_ADDRMAP9_addrmap_row_b2_Pos)
#define UDDRC_ADDRMAP9_addrmap_row_b2(value) \
	((UDDRC_ADDRMAP9_addrmap_row_b2_Msk & \
	((value) << UDDRC_ADDRMAP9_addrmap_row_b2_Pos)))

/* (UDDRC_ADDRMAP9) Selects the HIF address bits used as row address bit 3.
 * Valid Range: 0 to 11
 * Internal Base: 9
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP9_addrmap_row_b3_Pos 8
#define UDDRC_ADDRMAP9_addrmap_row_b3_Msk \
	(0xfu << UDDRC_ADDRMAP9_addrmap_row_b3_Pos)
#define UDDRC_ADDRMAP9_addrmap_row_b3(value) \
	((UDDRC_ADDRMAP9_addrmap_row_b3_Msk & \
	((value) << UDDRC_ADDRMAP9_addrmap_row_b3_Pos)))

/* (UDDRC_ADDRMAP9) Selects the HIF address bits used as row address bit 4.
 * Valid Range: 0 to 11
 * Internal Base: 10
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP9_addrmap_row_b4_Pos 16
#define UDDRC_ADDRMAP9_addrmap_row_b4_Msk \
	(0xfu << UDDRC_ADDRMAP9_addrmap_row_b4_Pos)
#define UDDRC_ADDRMAP9_addrmap_row_b4(value) \
	((UDDRC_ADDRMAP9_addrmap_row_b4_Msk & \
	((value) << UDDRC_ADDRMAP9_addrmap_row_b4_Pos)))

/* (UDDRC_ADDRMAP9) Selects the HIF address bits used as row address bit 5.
 * Valid Range: 0 to 11
 * Internal Base: 11
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP9_addrmap_row_b5_Pos 24
#define UDDRC_ADDRMAP9_addrmap_row_b5_Msk \
	(0xfu << UDDRC_ADDRMAP9_addrmap_row_b5_Pos)
#define UDDRC_ADDRMAP9_addrmap_row_b5(value) \
	((UDDRC_ADDRMAP9_addrmap_row_b5_Msk & \
	((value) << UDDRC_ADDRMAP9_addrmap_row_b5_Pos)))

/* -------- UDDRC_ADDRMAP10 : (UDDRC_REGS Offset: 0x228)
 * Address Map Register 10 --------
 */
/* (UDDRC_ADDRMAP10) Selects the HIF address bits used as row address bit 6.
 * Valid Range: 0 to 11
 * Internal Base: 12
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP10_addrmap_row_b6_Pos 0
#define UDDRC_ADDRMAP10_addrmap_row_b6_Msk \
	(0xfu << UDDRC_ADDRMAP10_addrmap_row_b6_Pos)
#define UDDRC_ADDRMAP10_addrmap_row_b6(value) \
	((UDDRC_ADDRMAP10_addrmap_row_b6_Msk & \
	((value) << UDDRC_ADDRMAP10_addrmap_row_b6_Pos)))

/* (UDDRC_ADDRMAP10) Selects the HIF address bits used as row address bit 7.
 * Valid Range: 0 to 11
 * Internal Base: 13
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP10_addrmap_row_b7_Pos 8
#define UDDRC_ADDRMAP10_addrmap_row_b7_Msk \
	(0xfu << UDDRC_ADDRMAP10_addrmap_row_b7_Pos)


#define UDDRC_ADDRMAP10_addrmap_row_b7(value) \
	((UDDRC_ADDRMAP10_addrmap_row_b7_Msk & \
	((value) << UDDRC_ADDRMAP10_addrmap_row_b7_Pos)))

/* (UDDRC_ADDRMAP10) Selects the HIF address bits used as row address bit 8.
 * Valid Range: 0 to 11
 * Internal Base: 14
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP10_addrmap_row_b8_Pos 16
#define UDDRC_ADDRMAP10_addrmap_row_b8_Msk \
	(0xfu << UDDRC_ADDRMAP10_addrmap_row_b8_Pos)
#define UDDRC_ADDRMAP10_addrmap_row_b8(value) \
	((UDDRC_ADDRMAP10_addrmap_row_b8_Msk & \
	((value) << UDDRC_ADDRMAP10_addrmap_row_b8_Pos)))

/* (UDDRC_ADDRMAP10) Selects the HIF address bits used as row address bit 9.
 * Valid Range: 0 to 11
 * Internal Base: 15
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP10_addrmap_row_b9_Pos 24
#define UDDRC_ADDRMAP10_addrmap_row_b9_Msk \
	(0xfu << UDDRC_ADDRMAP10_addrmap_row_b9_Pos)

#define UDDRC_ADDRMAP10_addrmap_row_b9(value) \
	((UDDRC_ADDRMAP10_addrmap_row_b9_Msk & \
	((value) << UDDRC_ADDRMAP10_addrmap_row_b9_Pos)))

/* -------- UDDRC_ADDRMAP11 : (UDDRC_REGS Offset: 0x22C)
 * Address Map Register 11 --------
 */
/* (UDDRC_ADDRMAP11) Selects the HIF address bits used as row address bit 10.
 * Valid Range: 0 to 11
 * Internal Base: 16
 * The selected HIF address bit for each of the row address bits is determined
 * by adding the internal base to the value of this field. This register field
 * is used only when ADDRMAP5.addrmap_row_b2_10 is set to value 15.
 */
#define UDDRC_ADDRMAP11_addrmap_row_b10_Pos 0
#define UDDRC_ADDRMAP11_addrmap_row_b10_Msk \
	(0xfu << UDDRC_ADDRMAP11_addrmap_row_b10_Pos)
#define UDDRC_ADDRMAP11_addrmap_row_b10(value) \
	((UDDRC_ADDRMAP11_addrmap_row_b10_Msk & \
	((value) << UDDRC_ADDRMAP11_addrmap_row_b10_Pos)))

/* -------- UDDRC_ODTCFG : (UDDRC_REGS Offset: 0x240)
 * ODT Configuration Register --------
 */
/* (UDDRC_ODTCFG)
 * The delay, in DFI PHY clock cycles, from issuing a read command to setting
 * ODT values associated with that command. ODT setting must remain constant
 * for the entire time that DQS is driven by the uMCTL2.
 *
 * Recommended values:
 * DDR2:
 * - CL + AL - 4 (not DDR2-1066),  CL + AL - 5 (DDR2-1066)
 * If (CL + AL - 4 < 0),  uMCTL2 does not support ODT for read operation.
 *
 * DDR3:
 * - CL - CWL
 * DDR4:
 * - CL - CWL - RD_PREAMBLE + WR_PREAMBLE + DFITMG1.dfi_t_cmd_lat
 * (to adjust for CAL mode)
 *
 *  WR_PREAMBLE = 1 (1tCK write preamble),  2 (2tCK write preamble)
 *  RD_PREAMBLE = 1 (1tCK write preamble),  2 (2tCK write preamble)
 *  If (CL - CWL - RD_PREAMBLE + WR_PREAMBLE) < 0,  uMCTL2 does not support
 * ODT for read operation.
 *
 * LPDDR3:
 * - RL + RD(tDQSCK(min)/tCK) - 1 - RU(tODTon(max)/tCK)
 */

#define UDDRC_ODTCFG_rd_odt_delay_Pos 2
#define UDDRC_ODTCFG_rd_odt_delay_Msk \
	(0x1fu << UDDRC_ODTCFG_rd_odt_delay_Pos)
#define UDDRC_ODTCFG_rd_odt_delay(value) \
	((UDDRC_ODTCFG_rd_odt_delay_Msk & \
	((value) << UDDRC_ODTCFG_rd_odt_delay_Pos)))

/* (UDDRC_ODTCFG)
 * DFI PHY clock cycles to hold ODT for a read command.
 * The minimum supported value is 2.
 * Recommended values:
 *
 * DDR2:
 * - BL8: 0x6 (not DDR2-1066),  0x7 (DDR2-1066)
 * - BL4: 0x4 (not DDR2-1066),  0x5 (DDR2-1066)
 * DDR3:
 * - BL8 - 0x6
 * DDR4:
 * - BL8: 5 + RD_PREAMBLE
 * RD_PREAMBLE = 1 (1tCK write preamble),  2 (2tCK write preamble)
 * LPDDR3:
 * - BL8:  5 + RU(tDQSCK(max)/tCK) - RD(tDQSCK(min)/tCK) + RU(tODTon(max)/tCK)
 */
#define UDDRC_ODTCFG_rd_odt_hold_Pos 8
#define UDDRC_ODTCFG_rd_odt_hold_Msk \
	(0xfu << UDDRC_ODTCFG_rd_odt_hold_Pos)
#define UDDRC_ODTCFG_rd_odt_hold(value) \
	((UDDRC_ODTCFG_rd_odt_hold_Msk & \
	((value) << UDDRC_ODTCFG_rd_odt_hold_Pos)))

/* (UDDRC_ODTCFG)
 * The delay, in DFI PHY clock cycles, from issuing a write command to setting
 * ODT values associated with that command. ODT setting must remain constant
 * for the entire time that DQS is driven by the uMCTL2.
 * Recommended values:
 *
 * DDR2:
 * - CWL + AL - 3 (DDR2-400/533/667),  CWL + AL - 4 (DDR2-800),
 *  CWL + AL - 5 (DDR2-1066)
 * If (CWL + AL - 3  < 0),  uMCTL2 does not support ODT for write operation.
 * DDR3:
 * - 0x0
 * DDR4:
 * - DFITMG1.dfi_t_cmd_lat (to adjust for CAL mode)
 * LPDDR3:
 * - WL - 1 - RU(tODTon(max)/tCK))
 */
#define UDDRC_ODTCFG_wr_odt_delay_Pos 16
#define UDDRC_ODTCFG_wr_odt_delay_Msk (0x1fu << UDDRC_ODTCFG_wr_odt_delay_Pos)
#define UDDRC_ODTCFG_wr_odt_delay(value) \
	((UDDRC_ODTCFG_wr_odt_delay_Msk & \
	((value) << UDDRC_ODTCFG_wr_odt_delay_Pos)))

/* (UDDRC_ODTCFG)
 * DFI PHY clock cycles to hold ODT for a write command.
 * The minimum supported value is 2.
 * Recommended values:
 * DDR2:
 * - BL8:  0x5 (DDR2-400/533/667),  0x6 (DDR2-800),  0x7 (DDR2-1066)
 * - BL4:  0x3 (DDR2-400/533/667),  0x4 (DDR2-800),  0x5 (DDR2-1066)
 * DDR3:
 * - BL8: 0x6
 * DDR4:
 * - BL8: 5 + WR_PREAMBLE + CRC_MODE
 * WR_PREAMBLE = 1 (1tCK write preamble),  2 (2tCK write preamble)
 * CRC_MODE = 0 (not CRC mode),  1 (CRC mode)
 * LPDDR3:
 * - BL8: 7 + RU(tODTon(max)/tCK)
 */
#define UDDRC_ODTCFG_wr_odt_hold_Pos 24
#define UDDRC_ODTCFG_wr_odt_hold_Msk (0xfu << UDDRC_ODTCFG_wr_odt_hold_Pos)
#define UDDRC_ODTCFG_wr_odt_hold(value) \
	((UDDRC_ODTCFG_wr_odt_hold_Msk & \
	((value) << UDDRC_ODTCFG_wr_odt_hold_Pos)))

/* -------- UDDRC_ODTMAP : (UDDRC_REGS Offset: 0x244)
 * ODT/Rank Map Register --------
 */
/* (UDDRC_ODTMAP) Indicates which remote ODTs must be turned on during a write
 * to rank 0.
 * Each rank has a remote ODT (in the SDRAM) which can be turned on by setting
 * the appropriate bit here.
 * Rank 0 is controlled by the LSB; rank 1 is controlled by bit next to the
 * LSB, etc.
 * For each rank, set its bit to 1 to enable its ODT.
 */
#define UDDRC_ODTMAP_rank0_wr_odt (0x1u << 0)

/* (UDDRC_ODTMAP) Indicates which remote ODTs must be turned on during a read
 * from rank 0.
 * Each rank has a remote ODT (in the SDRAM) which can be turned on by setting
 * the appropriate bit here.
 * Rank 0 is controlled by the LSB; rank 1 is controlled by bit next to the LSB,
 * etc.
 * For each rank, set its bit to 1 to enable its ODT.
 */
#define UDDRC_ODTMAP_rank0_rd_odt (0x1u << 4)

/* -------- UDDRC_SCHED : (UDDRC_REGS Offset: 0x250)
 * Scheduler Control Register --------
 */
/* (UDDRC_SCHED) Active low signal. When asserted ('0'), all incoming
 * transactions are forced to low priority. This implies that all
 * High Priority Read (HPR) and Variable Priority Read commands (VPR) will be
 * treated as Low Priority Read (LPR) commands. On the write side, all
 * Variable Priority Write (VPW) commands will be treated as
 * Normal Priority Write (NPW) commands. Forcing the incoming transactions to
 * low priority implicitly turns off Bypass path for read commands.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_SCHED_force_low_pri_n (0x1u << 0)

/* (UDDRC_SCHED) If set then the bank selector prefers writes over reads.
 * FOR DEBUG ONLY.
 */
#define UDDRC_SCHED_prefer_write (0x1u << 1)

/* (UDDRC_SCHED) If true, bank is kept open only while there are page hit
 * transactions available in the CAM to that bank. The last read or write
 * command in the CAM with a bank and page hit will be executed with
 * auto-precharge if SCHED1.pageclose_timer=0. Even if this register set to 1
 * and SCHED1.pageclose_timer is set to 0, explicit precharge
 * (and not auto-precharge) may be issued in some cases where there is a mode
 * switch between Write and Read or between LPR and HPR. The Read and Write
 * commands that are executed as part of the ECC scrub requests are also
 * executed without auto-precharge.
 *
 * If false, the bank remains open until there is a need to close it
 * (to open a different page, or for page timeout or refresh timeout)
 * - also known as open page policy. The open page policy can be overridden
 * by setting the per-command-autopre bit on the HIF interface (hif_cmd_autopre).
 * The pageclose feature provids a midway between Open and Close page policies.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_SCHED_pageclose (0x1u << 2)

/* (UDDRC_SCHED) Number of entries in the low priority transaction store
 * is this value + 1.
 * (MEMC_NO_OF_ENTRY - (SCHED.lpr_num_entries + 1)) is the number of entries
 * available for the high priority transaction store.
 * Setting this to maximum value allocates all entries to low priority
 * transaction store.
 * Setting this to 0 allocates 1 entry to low priority transaction store and
 * the rest to high priority transaction store.
 * Note: In ECC configurations, the numbers of write and low priority read
 * credits issued is one less than in the non-ECC case.  One entry each is
 * reserved in the write and low-priority read CAMs for storing the RMW requests
 * arising out of single bit error correction RMW operation.
 */
#define UDDRC_SCHED_lpr_num_entries_Pos 8
#define UDDRC_SCHED_lpr_num_entries_Msk \
	(0x1fu << UDDRC_SCHED_lpr_num_entries_Pos)
#define UDDRC_SCHED_lpr_num_entries(value) \
	((UDDRC_SCHED_lpr_num_entries_Msk & \
	((value) << UDDRC_SCHED_lpr_num_entries_Pos)))

/* (UDDRC_SCHED) UNUSED */
#define UDDRC_SCHED_go2critical_hysteresis_Pos 16
#define UDDRC_SCHED_go2critical_hysteresis_Msk \
	(0xffu << UDDRC_SCHED_go2critical_hysteresis_Pos)
#define UDDRC_SCHED_go2critical_hysteresis(value) \
	((UDDRC_SCHED_go2critical_hysteresis_Msk & \
	((value) << UDDRC_SCHED_go2critical_hysteresis_Pos)))

/* (UDDRC_SCHED) When the preferred transaction store is empty for these many
 * clock cycles, switch to the alternate transaction store if it is non-empty.
 * The read transaction store (both high and low priority) is the default
 * preferred transaction store and the write transaction store is the
 * alternative store.
 * When prefer write over read is set this is reversed.
 * 0x0 is a legal value for this register. When set to 0x0, the transaction
 * store switching will happen immediately when the switching conditions
 * become true.
 * FOR PERFORMANCE ONLY
 */
#define UDDRC_SCHED_rdwr_idle_gap_Pos 24
#define UDDRC_SCHED_rdwr_idle_gap_Msk (0x7fu << UDDRC_SCHED_rdwr_idle_gap_Pos)
#define UDDRC_SCHED_rdwr_idle_gap(value) \
	((UDDRC_SCHED_rdwr_idle_gap_Msk & \
	((value) << UDDRC_SCHED_rdwr_idle_gap_Pos)))

/* -------- UDDRC_SCHED1 : (UDDRC_REGS Offset: 0x254)
 * Scheduler Control Register 1 --------
 */
/* (UDDRC_SCHED1) This field works in conjunction with SCHED.pageclose.
 * It only has meaning if SCHED.pageclose==1.
 *
 * If SCHED.pageclose==1 and pageclose_timer==0, then an auto-precharge
 * may be scheduled for last read or write command in the CAM with a bank
 * and page hit.
 * Note, sometimes an explicit precharge is scheduled instead of the
 * auto-precharge. See SCHED.pageclose for details of when this may happen.
 *
 * If SCHED.pageclose==1 and pageclose_timer>0, then an auto-precharge is not
 * scheduled for last read or write command in the CAM with a bank and page hit.
 * Instead, a timer is started, with pageclose_timer as the initial value.
 * There is a timer on a per bank basis.
 * The timer decrements unless the next read or write in the CAM to a bank is a
 * page hit. It gets reset to pageclose_timer value if the next read or write
 * in the CAM to a bank is a page hit. Once the timer has reached zero,
 * an explcit precharge will be attempted to be scheduled.
 */
#define UDDRC_SCHED1_pageclose_timer_Pos 0
#define UDDRC_SCHED1_pageclose_timer_Msk \
	(0xffu << UDDRC_SCHED1_pageclose_timer_Pos)
#define UDDRC_SCHED1_pageclose_timer(value) \
	((UDDRC_SCHED1_pageclose_timer_Msk & \
	((value) << UDDRC_SCHED1_pageclose_timer_Pos)))

/* -------- UDDRC_PERFHPR1 : (UDDRC_REGS Offset: 0x25C)
 * High Priority Read CAM Register 1 --------
 */
/* (UDDRC_PERFHPR1) Number of DFI clocks that the HPR queue can be starved
 * before it goes critical. The minimum valid functional value for this register
 * is 0x1. Programming it to 0x0 will disable the starvation functionality;
 * during normal operation, this function should not be disabled as it will
 * cause excessive latencies.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFHPR1_hpr_max_starve_Pos 0
#define UDDRC_PERFHPR1_hpr_max_starve_Msk \
	(0xffffu << UDDRC_PERFHPR1_hpr_max_starve_Pos)
#define UDDRC_PERFHPR1_hpr_max_starve(value) \
	((UDDRC_PERFHPR1_hpr_max_starve_Msk & \
	((value) << UDDRC_PERFHPR1_hpr_max_starve_Pos)))

/* (UDDRC_PERFHPR1)
 * Number of transactions that are serviced once the HPR queue goes critical
 * is the smaller of:
 * - (a) This number
 * - (b) Number of transactions available.
 * Unit: Transaction.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFHPR1_hpr_xact_run_length_Pos 24
#define UDDRC_PERFHPR1_hpr_xact_run_length_Msk \
	(0xffu << UDDRC_PERFHPR1_hpr_xact_run_length_Pos)
#define UDDRC_PERFHPR1_hpr_xact_run_length(value) \
	((UDDRC_PERFHPR1_hpr_xact_run_length_Msk & \
	((value) << UDDRC_PERFHPR1_hpr_xact_run_length_Pos)))

/* -------- UDDRC_PERFLPR1 : (UDDRC_REGS Offset: 0x264)
 * Low Priority Read CAM Register 1 --------
 */
/* (UDDRC_PERFLPR1) Number of DFI clocks that the LPR queue can be starved
 * before it goes critical. The minimum valid functional value for this register
 * is 0x1. Programming it to 0x0 will disable the starvation functionality;
 * during normal operation, this function should not be disabled as it will
 * cause excessive latencies.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFLPR1_lpr_max_starve_Pos 0
#define UDDRC_PERFLPR1_lpr_max_starve_Msk \
	(0xffffu << UDDRC_PERFLPR1_lpr_max_starve_Pos)
#define UDDRC_PERFLPR1_lpr_max_starve(value) \
	((UDDRC_PERFLPR1_lpr_max_starve_Msk & \
	((value) << UDDRC_PERFLPR1_lpr_max_starve_Pos)))

/* (UDDRC_PERFLPR1)
 * Number of transactions that are serviced once the LPR queue goes critical
 * is the smaller of:
 * - (a) This number
 * - (b) Number of transactions available.
 * Unit: Transaction.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFLPR1_lpr_xact_run_length_Pos 24
#define UDDRC_PERFLPR1_lpr_xact_run_length_Msk \
	(0xffu << UDDRC_PERFLPR1_lpr_xact_run_length_Pos)
#define UDDRC_PERFLPR1_lpr_xact_run_length(value) \
	((UDDRC_PERFLPR1_lpr_xact_run_length_Msk & \
	((value) << UDDRC_PERFLPR1_lpr_xact_run_length_Pos)))
/* -------- UDDRC_PERFWR1 : (UDDRC_REGS Offset: 0x26C)
 * Write CAM Register 1 --------
 */
/* (UDDRC_PERFWR1) Number of DFI clocks that the WR queue can be starved before
 * it goes critical. The minimum valid functional value for this register is 0x1.
 * Programming it to 0x0 will disable the starvation functionality;
 * during normal operation, this function should not be disabled as it will
 * cause excessive latencies.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFWR1_w_max_starve_Pos 0
#define UDDRC_PERFWR1_w_max_starve_Msk \
	(0xffffu << UDDRC_PERFWR1_w_max_starve_Pos)
#define UDDRC_PERFWR1_w_max_starve(value) \
	((UDDRC_PERFWR1_w_max_starve_Msk & \
	((value) << UDDRC_PERFWR1_w_max_starve_Pos)))

/* (UDDRC_PERFWR1)
 * Number of transactions that are serviced once the WR queue goes critical
 * is the smaller of:
 * - (a) This number
 * - (b) Number of transactions available.
 * Unit: Transaction.
 * FOR PERFORMANCE ONLY.
 */
#define UDDRC_PERFWR1_w_xact_run_length_Pos 24
#define UDDRC_PERFWR1_w_xact_run_length_Msk \
	(0xffu << UDDRC_PERFWR1_w_xact_run_length_Pos)
#define UDDRC_PERFWR1_w_xact_run_length(value) \
	((UDDRC_PERFWR1_w_xact_run_length_Msk & \
	((value) << UDDRC_PERFWR1_w_xact_run_length_Pos)))

/* -------- UDDRC_DBG0 : (UDDRC_REGS Offset: 0x300)
 * Debug Register 0 --------
 */
/* (UDDRC_DBG0) When 1, disable write combine.
 * FOR DEBUG ONLY
 */
/* (UDDRC_DBG0) Only present in designs supporting read bypass.
 * When 1, disable bypass path for high priority read page hits
 * FOR DEBUG ONLY.
 */
#define UDDRC_DBG0_dis_wc (0x1u << 0)

/* (UDDRC_DBG0) Only present in designs supporting activate bypass.
 * When 1, disable bypass path for high priority read activates
 * FOR DEBUG ONLY.
 */
#define UDDRC_DBG0_dis_rd_bypass (0x1u << 1)
#define UDDRC_DBG0_dis_act_bypass (0x1u << 2)

/* (UDDRC_DBG0) When this is set to '0', auto-precharge is disabled for the
 * flushed command in a collision case. Collision cases are write followed by
 * read to same address, read followed by write to same address,
 * or write followed by write to same address with DBG0.dis_wc bit = 1
 * (where same address comparisons exclude the two address bits representing
 * critical word).
 * FOR DEBUG ONLY.
 */
#define UDDRC_DBG0_dis_collision_page_opt (0x1u << 4)
/* -------- UDDRC_DBG1 : (UDDRC_REGS Offset: 0x304)
 * Debug Register 1 --------
 */
/* (UDDRC_DBG1) When 1, uMCTL2 will not de-queue any transactions from the CAM.
 * Bypass is also disabled. All transactions are queued in the CAM. No reads or
 * writes are issued to SDRAM as long as this is asserted.
 * This bit may be used to prevent reads or writes being issued by the uMCTL2,
 * which makes it safe to modify certain register fields associated with
 * reads and writes (see User Guide for details).  After setting this bit,
 * it is strongly recommended to poll DBGCAM.wr_data_pipeline_empty and
 * DBGCAM.rd_data_pipeline_empty, before making changes to any registers which
 * affect reads and writes.
 * This will ensure that the relevant logic in the DDRC is idle.
 * This bit is intended to be switched on-the-fly.
 */
#define UDDRC_DBG1_dis_dq (0x1u << 0)

/* (UDDRC_DBG1) When 1, uMCTL2 asserts the HIF command signal hif_cmd_stall.
 * uMCTL2 will ignore the hif_cmd_valid and all other associated request signals.
 * This bit is intended to be switched on-the-fly.
 */
#define UDDRC_DBG1_dis_hif (0x1u << 1)

/* -------- UDDRC_DBGCAM : (UDDRC_REGS Offset: 0x308)
 * CAM Debug Register --------
 */
/* (UDDRC_DBGCAM) High priority read queue depth
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_hpr_q_depth_Pos 0
#define UDDRC_DBGCAM_dbg_hpr_q_depth_Msk \
	(0x3fu << UDDRC_DBGCAM_dbg_hpr_q_depth_Pos)
/* (UDDRC_DBGCAM) Low priority read queue depth
 * The last entry of Lpr queue is reserved for ECC SCRUB operation.
 * This entry is not included in the calculation of the queue depth.
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_lpr_q_depth_Pos 8
#define UDDRC_DBGCAM_dbg_lpr_q_depth_Msk \
	(0x3fu << UDDRC_DBGCAM_dbg_lpr_q_depth_Pos)

/* (UDDRC_DBGCAM) Write queue depth
 * The last entry of WR queue is reserved for ECC SCRUB operation.
 * This entry is not included in the calculation of the queue depth.
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_w_q_depth_Pos 16
#define UDDRC_DBGCAM_dbg_w_q_depth_Msk \
	(0x3fu << UDDRC_DBGCAM_dbg_w_q_depth_Pos)

/* (UDDRC_DBGCAM) Stall
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_stall (0x1u << 24)

/* (UDDRC_DBGCAM) When 1, all the Read command queues and Read data buffers
 * inside DDRC are empty. This register is to be used for debug purpose.
 * An example use-case scenario: When Controller enters Self-Refresh using the
 * Low-Power entry sequence, Controller is expected to have executed all the
 * commands in its queues and the write and read data drained.
 * Hence this register should be 1 at that time.
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_rd_q_empty (0x1u << 25)

/* (UDDRC_DBGCAM) When 1, all the Write command queues and Write data buffers
 * inside DDRC are empty. This register is to be used for debug purpose.
 * An example use-case scenario: When Controller enters Self-Refresh using the
 * Low-Power entry sequence, Controller is expected to have executed all the
 * commands in its queues and the write and read data drained.
 * Hence this register should be 1 at that time.
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_wr_q_empty (0x1u << 26)

/* (UDDRC_DBGCAM) This bit indicates that the read data pipeline on the
 * DFI interface is empty.  This register is intended to be polled at least
 * twice after setting DBG1.dis_dq, to ensure that all remaining commands/data
 * have completed.
 */
#define UDDRC_DBGCAM_rd_data_pipeline_empty (0x1u << 28)

/* (UDDRC_DBGCAM) This bit indicates that the write data pipeline on the
 * DFI interface is empty.  This register is intended to be polled at least
 * twice after setting DBG1.dis_dq, to ensure that all remaining commands/data
 * have completed.
 */
#define UDDRC_DBGCAM_wr_data_pipeline_empty (0x1u << 29)

/* (UDDRC_DBGCAM) Stall for Write channel
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_stall_wr (0x1u << 30)

/* (UDDRC_DBGCAM) Stall for Read channel
 * FOR DEBUG ONLY
 */
#define UDDRC_DBGCAM_dbg_stall_rd (0x1u << 31)

/* -------- UDDRC_DBGCMD : (UDDRC_REGS Offset: 0x30C)
 * Command Debug Register --------
 */
/* (UDDRC_DBGCMD)
 * Setting this register bit to 1 indicates to the uMCTL2 to issue a refresh to
 * rank 0. Writing to this bit causes DBGSTAT.rank0_refresh_busy to be set.
 * When DBGSTAT.rank0_refresh_busy is cleared, the command has been stored in
 * uMCTL2.
 * For 3DS configuration, refresh is sent to rank index 0.
 * This operation can be performed only when RFSHCTL3.dis_auto_refresh=1.
 * It is recommended NOT to set this register bit if in Init or Deep power-down
 * operating modes or Maximum Power Saving Mode.
 */
#define UDDRC_DBGCMD_rank0_refresh (0x1u << 0)

/* (UDDRC_DBGCMD)
 * Setting this register bit to 1 indicates to the uMCTL2 to issue a
 * ZQCS (ZQ calibration short)/MPC(ZQ calibration) command to the SDRAM.
 * When this request is stored in the uMCTL2, the bit is automatically cleared.
 * This operation can be performed only when ZQCTL0.dis_auto_zq=1. It is
 * recommended NOT to set this register bit if in Init operating mode.
 * This register bit is ignored when in Self-Refresh(except LPDDR4) and
 * SR-Powerdown(LPDDR4) and Deep power-down operating modes and
 * Maximum Power Saving Mode.
 */
#define UDDRC_DBGCMD_zq_calib_short (0x1u << 4)

/* (UDDRC_DBGCMD)
 * Setting this register bit to 1 indicates to the uMCTL2 to issue a
 * dfi_ctrlupd_req to the PHY.  When this request is stored in the uMCTL2,
 * the bit is automatically cleared. This operation must only be performed when
 * DFIUPD0.dis_auto_ctrlupd=1.
 */
#define UDDRC_DBGCMD_ctrlupd (0x1u << 5)

/* -------- UDDRC_DBGSTAT : (UDDRC_REGS Offset: 0x310)
 * Status Debug Register --------
 */
/* (UDDRC_DBGSTAT) SoC core may initiate a rank0_refresh operation
 * (refresh operation to rank 0) only if this signal is low. This signal
 * goes high in the clock after DBGCMD.rank0_refresh is set to one.
 * It goes low when the rank0_refresh operation is stored in the uMCTL2.
 * It is recommended not to perform rank0_refresh operations when this signal
 * is high.
 * - 0 - Indicates that the SoC core can initiate a rank0_refresh operation
 * - 1 - Indicates that rank0_refresh operation has not been stored yet in the uMCTL2
 */
#define UDDRC_DBGSTAT_rank0_refresh_busy (0x1u << 0)

/* (UDDRC_DBGSTAT) SoC core may initiate a ZQCS (ZQ calibration short) operation
 * only if this signal is low. This signal goes high in the clock after the
 * uMCTL2 accepts the ZQCS request. It goes low when the ZQCS operation is
 * initiated in the uMCTL2. It is recommended not to perform ZQCS operations
 * when this signal is high.
 * - 0 - Indicates that the SoC core can initiate a ZQCS operation
 * - 1 - Indicates that ZQCS operation has not been initiated yet in the uMCTL2
 */
#define UDDRC_DBGSTAT_zq_calib_short_busy (0x1u << 4)

/* (UDDRC_DBGSTAT) SoC core may initiate a ctrlupd operation only if this
 * signal is low. This signal goes high in the clock after the uMCTL2 accepts
 * the ctrlupd request. It goes low when the ctrlupd operation is initiated
 * in the uMCTL2. It is recommended not to perform ctrlupd operations when this
 * signal is high.
 * - 0 - Indicates that the SoC core can initiate a ctrlupd operation
 * - 1 - Indicates that ctrlupd operation has not been initiated yet in the uMCTL2
 */
#define UDDRC_DBGSTAT_ctrlupd_busy (0x1u << 5)

/* -------- UDDRC_SWCTL : (UDDRC_REGS Offset: 0x320)
 * Software Register Programming Control Enable --------
 */
/* (UDDRC_SWCTL) Enable quasi-dynamic register programming outside reset.
 * Program register to 0 to enable quasi-dynamic programming.
 * Set back register to 1 once programming is done.
 */
#define UDDRC_SWCTL_sw_done (0x1u << 0)

/* -------- UDDRC_SWSTAT : (UDDRC_REGS Offset: 0x324)
 * Software Register Programming Control Status --------
 */
/* (UDDRC_SWSTAT) Register programming done. This register is the echo of
 * SWCTL.sw_done. Wait for sw_done value 1 to propagate to sw_done_ack
 * at the end of the programming sequence to ensure that the correct registers
 * values are propagated to the destination clock domains.
 */
#define UDDRC_SWSTAT_sw_done_ack (0x1u << 0)

/* -------- UDDRC_POISONCFG : (UDDRC_REGS Offset: 0x36C)
 * AXI Poison Configuration Register. Common for all AXI ports --------
 */
/* (UDDRC_POISONCFG) If set to 1, enables SLVERR response for write
 * transaction poisoning
 */
#define UDDRC_POISONCFG_wr_poison_slverr_en (0x1u << 0)

/* (UDDRC_POISONCFG) If set to 1, enables interrupts for write
 * transaction poisoning
 */
#define UDDRC_POISONCFG_wr_poison_intr_en (0x1u << 4)

/* (UDDRC_POISONCFG) Interrupt clear for write transaction poisoning.
 * Allow 2/3 clock cycles for correct value to propagate to core logic and
 * clear the interrupts.
 */
#define UDDRC_POISONCFG_wr_poison_intr_clr (0x1u << 8)

/* (UDDRC_POISONCFG) If set to 1, enables SLVERR response for read
 * transaction poisoning
 */
#define UDDRC_POISONCFG_rd_poison_slverr_en (0x1u << 16)

/* (UDDRC_POISONCFG) If set to 1, enables interrupts for read
 * transaction poisoning
 */
#define UDDRC_POISONCFG_rd_poison_intr_en (0x1u << 20)

/* (UDDRC_POISONCFG) Interrupt clear for read transaction poisoning.
 * Allow 2/3 clock cycles for correct value to propagate to core logic and
 * clear the interrupts.
 */
#define UDDRC_POISONCFG_rd_poison_intr_clr (0x1u << 24)

/* -------- UDDRC_POISONSTAT : (UDDRC_REGS Offset: 0x370)
 * AXI Poison Status Register --------
 */
/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 0.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_0 (0x1u << 0)

/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 1.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_1 (0x1u << 1)

/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 2.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_2 (0x1u << 2)

/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 3.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_3 (0x1u << 3)

/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 4.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_4 (0x1u << 4)

/* (UDDRC_POISONSTAT) Write transaction poisoning error interrupt for port 5.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's write address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register wr_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_wr_poison_intr_5 (0x1u << 5)

/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 0.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */

#define UDDRC_POISONSTAT_rd_poison_intr_0 (0x1u << 16)
/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 1.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_rd_poison_intr_1 (0x1u << 17)

/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 2.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_rd_poison_intr_2 (0x1u << 18)

/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 3.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_rd_poison_intr_3 (0x1u << 19)

/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 4.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_rd_poison_intr_4 (0x1u << 20)

/* (UDDRC_POISONSTAT) Read transaction poisoning error interrupt for port 5.
 * This register is a APB clock copy (double register synchronizer) of the
 * interrupt asserted when a transaction is poisoned on the corresponding
 * AXI port's read address channel. Bit 0 corresponds to Port 0, and so on.
 * Interrupt is cleared by register rd_poison_intr_clr, then value propagated
 * to APB clock.
 */
#define UDDRC_POISONSTAT_rd_poison_intr_5 (0x1u << 21)

/* } */

/* UMCTL2 main register helpers end */
#endif
