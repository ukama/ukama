/*
 * This file is part of the coreboot project.
 *
 * Copyright (c) 2018 Qualcomm Technologies
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2 and
 * only version 2 as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/mmio.h>
#include <stdlib.h>
#include <console/console.h>
#include <delay.h>
#include <soc/usb.h>
#include <soc/clock.h>
#include <soc/addressmap.h>
#include <soc/efuse.h>
#include <timer.h>

struct usb_qusb_phy_dig {
	u8 rsvd1[16];
	u32 pwr_ctrl1;
	u32 pwr_ctrl2;
	u8 rsvd2[8];
	u32 imp_ctrl1;
	u32 imp_ctrl2;
	u8 rsvd3[20];
	u32 chg_ctrl2;
	u32 tune1;
	u32 tune2;
	u32 tune3;
	u32 tune4;
	u32 tune5;
	u8 rsvd4[44];
	u32 debug_ctrl2;
	u8 rsvd5[28];
	u32 debug_stat5;
};
check_member(usb_qusb_phy_dig, tune5, 0x50);
check_member(usb_qusb_phy_dig, debug_ctrl2, 0x80);
check_member(usb_qusb_phy_dig, debug_stat5, 0xA0);

struct usb_qusb_phy_pll {
	u8 rsvd0[4];
	u32 analog_controls_two;
	u8 rsvd1[36];
	u32 cmode;
	u8 rsvd2[132];
	u32 dig_tim;
	u8 rsvd3[204];
	u32 lock_delay;
	u8 rsvd4[4];
	u32 clock_inverters;
	u8 rsvd5[4];
	u32 bias_ctrl_1;
	u32 bias_ctrl_2;
};
check_member(usb_qusb_phy_pll, cmode, 0x2C);
check_member(usb_qusb_phy_pll, bias_ctrl_2, 0x198);
check_member(usb_qusb_phy_pll, dig_tim, 0xB4);

/* Only for QMP V3 PHY - QSERDES COM registers */
struct usb3_phy_qserdes_com_reg_layout {
	u8 _reserved1[16];
	u32 com_ssc_en_center;
	u32 com_ssc_adj_per1;
	u32 com_ssc_adj_per2;
	u32 com_ssc_per1;
	u32 com_ssc_per2;
	u32 com_ssc_step_size1;
	u32 com_ssc_step_size2;
	u8 _reserved2[8];
	u32 com_bias_en_clkbuflr_en;
	u32 com_sys_clk_enable1;
	u32 com_sys_clk_ctrl;
	u32 com_sysclk_buf_enable;
	u32 com_pll_en;
	u32 com_pll_ivco;
	u8 _reserved3[20];
	u32 com_cp_ctrl_mode0;
	u8 _reserved4[4];
	u32 com_pll_rctrl_mode0;
	u8 _reserved5[4];
	u32 com_pll_cctrl_mode0;
	u8 _reserved6[12];
	u32 com_sysclk_en_sel;
	u8 _reserved7[8];
	u32 com_resetsm_ctrl2;
	u32 com_lock_cmp_en;
	u32 com_lock_cmp_cfg;
	u32 com_lock_cmp1_mode0;
	u32 com_lock_cmp2_mode0;
	u32 com_lock_cmp3_mode0;
	u8 _reserved8[12];
	u32 com_dec_start_mode0;
	u8 _reserved9[4];
	u32 com_div_frac_start1_mode0;
	u32 com_div_frac_start2_mode0;
	u32 com_div_frac_start3_mode0;
	u8 _reserved10[20];
	u32 com_integloop_gain0_mode0;
	u32 com_integloop_gain1_mode0;
	u8 _reserved11[16];
	u32 com_vco_tune_map;
	u32 com_vco_tune1_mode0;
	u32 com_vco_tune2_mode0;
	u8 _reserved12[60];
	u32 com_clk_select;
	u32 com_hsclk_sel;
	u8 _reserved13[8];
	u32 com_coreclk_div_mode0;
	u8 _reserved14[8];
	u32 com_core_clk_en;
	u32 com_c_ready_status;
	u32 com_cmn_config;
	u32 com_cmn_rate_override;
	u32 com_svs_mode_clk_sel;
};
check_member(usb3_phy_qserdes_com_reg_layout, com_ssc_en_center, 0x010);
check_member(usb3_phy_qserdes_com_reg_layout, com_ssc_adj_per1, 0x014);
check_member(usb3_phy_qserdes_com_reg_layout, com_ssc_adj_per2, 0x018);
check_member(usb3_phy_qserdes_com_reg_layout, com_ssc_per1, 0x01c);
check_member(usb3_phy_qserdes_com_reg_layout, com_ssc_per2, 0x020);
check_member(usb3_phy_qserdes_com_reg_layout, com_bias_en_clkbuflr_en, 0x034);
check_member(usb3_phy_qserdes_com_reg_layout, com_pll_ivco, 0x048);
check_member(usb3_phy_qserdes_com_reg_layout, com_cp_ctrl_mode0, 0x060);
check_member(usb3_phy_qserdes_com_reg_layout, com_sysclk_en_sel, 0x080);
check_member(usb3_phy_qserdes_com_reg_layout, com_resetsm_ctrl2, 0x08c);
check_member(usb3_phy_qserdes_com_reg_layout, com_dec_start_mode0, 0x0b0);
check_member(usb3_phy_qserdes_com_reg_layout, com_div_frac_start1_mode0, 0x0b8);
check_member(usb3_phy_qserdes_com_reg_layout, com_integloop_gain0_mode0, 0x0d8);
check_member(usb3_phy_qserdes_com_reg_layout, com_vco_tune_map, 0x0f0);
check_member(usb3_phy_qserdes_com_reg_layout, com_clk_select, 0x138);
check_member(usb3_phy_qserdes_com_reg_layout, com_coreclk_div_mode0, 0x148);
check_member(usb3_phy_qserdes_com_reg_layout, com_core_clk_en, 0x154);
check_member(usb3_phy_qserdes_com_reg_layout, com_svs_mode_clk_sel, 0x164);

/* Only for QMP V3 PHY - TX registers */
struct usb3_phy_qserdes_tx_reg_layout {
	u8 _reserved1[68];
	u32 tx_res_code_lane_offset_tx;
	u32 tx_res_code_lane_offset_rx;
	u8 _reserved2[20];
	u32 tx_highz_drvr_en;
	u8 _reserved3[40];
	u32 tx_lane_mode_1;
	u8 _reserved4[20];
	u32 tx_rcv_detect_lvl_2;
};
check_member(usb3_phy_qserdes_tx_reg_layout, tx_res_code_lane_offset_tx, 0x044);
check_member(usb3_phy_qserdes_tx_reg_layout, tx_res_code_lane_offset_rx, 0x048);
check_member(usb3_phy_qserdes_tx_reg_layout, tx_highz_drvr_en, 0x060);
check_member(usb3_phy_qserdes_tx_reg_layout, tx_lane_mode_1, 0x08c);
check_member(usb3_phy_qserdes_tx_reg_layout, tx_rcv_detect_lvl_2, 0x0a4);

/* Only for QMP V3 PHY - RX registers */
struct usb3_phy_qserdes_rx_reg_layout {
	u8 _reserved1[8];
	u32 rx_ucdr_fo_gain;
	u32 rx_ucdr_so_gain_half;
	u8 _reserved2[32];
	u32 rx_ucdr_fastlock_fo_gain;
	u32 rx_ucdr_so_saturtn_and_en;
	u8 _reserved3[12];
	u32 rx_ucdr_pi_cntrls;
	u8 _reserved4[120];
	u32 rx_vga_cal_ctrl2;
	u8 _reserved5[16];
	u32 rx_rx_equ_adap_ctrl2;
	u32 rx_rx_equ_adap_ctrl3;
	u32 rx_rx_equ_adap_ctrl4;
	u8 _reserved6[24];
	u32 rx_rx_eq_offset_adap_ctrl1;
	u32 rx_rx_offset_adap_ctrl2;
	u32 rx_sigdet_enables;
	u32 rx_sigdet_ctrl;
	u8 _reserved7[4];
	u32 rx_sigdet_deglitch_ctrl;
	u32 rx_rx_band;
	u8 _reserved8[80];
	u32 rx_rx_mode_00;
};
check_member(usb3_phy_qserdes_rx_reg_layout, rx_ucdr_fo_gain, 0x008);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_ucdr_so_gain_half, 0x00c);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_ucdr_fastlock_fo_gain, 0x030);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_ucdr_so_saturtn_and_en, 0x034);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_ucdr_pi_cntrls, 0x044);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_vga_cal_ctrl2, 0x0c0);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_equ_adap_ctrl2, 0x0d4);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_equ_adap_ctrl3, 0x0d8);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_equ_adap_ctrl4, 0x0dc);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_eq_offset_adap_ctrl1, 0x0f8);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_offset_adap_ctrl2, 0x0fc);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_sigdet_enables, 0x100);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_sigdet_ctrl, 0x104);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_sigdet_deglitch_ctrl, 0x10c);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_band, 0x110);
check_member(usb3_phy_qserdes_rx_reg_layout, rx_rx_mode_00, 0x164);

/* Only for QMP V3 PHY - PCS registers */
struct usb3_phy_pcs_reg_layout {
	u32 pcs_sw_reset;
	u32 pcs_power_down_control;
	u32 pcs_start_control;
	u32 pcs_txmgn_v0;
	u32 pcs_txmgn_v1;
	u32 pcs_txmgn_v2;
	u32 pcs_txmgn_v3;
	u32 pcs_txmgn_v4;
	u32 pcs_txmgn_ls;
	u32 pcs_txdeemph_m6db_v0;
	u32 pcs_txdeemph_m3p5db_v0;
	u32 pcs_txdeemph_m6db_v1;
	u32 pcs_txdeemph_m3p5db_v1;
	u32 pcs_txdeemph_m6db_v2;
	u32 pcs_txdeemph_m3p5db_v2;
	u32 pcs_txdeemph_m6db_v3;
	u32 pcs_txdeemph_m3p5db_v3;
	u32 pcs_txdeemph_m6db_v4;
	u32 pcs_txdeemph_m3p5db_v4;
	u32 pcs_txdeemph_m6db_ls;
	u32 pcs_txdeemph_m3p5db_ls;
	u8 _reserved1[8];
	u32 pcs_rate_slew_cntrl;
	u8 _reserved2[4];
	u32 pcs_power_state_config2;
	u8 _reserved3[8];
	u32 pcs_rcvr_dtct_dly_p1u2_l;
	u32 pcs_rcvr_dtct_dly_p1u2_h;
	u32 pcs_rcvr_dtct_dly_u3_l;
	u32 pcs_rcvr_dtct_dly_u3_h;
	u32 pcs_lock_detect_config1;
	u32 pcs_lock_detect_config2;
	u32 pcs_lock_detect_config3;
	u32 pcs_tsync_rsync_time;
	u8 _reserved4[16];
	u32 pcs_pwrup_reset_dly_time_auxclk;
	u8 _reserved5[12];
	u32 pcs_lfps_ecstart_eqtlock;
	u8 _reserved6[4];
	u32 pcs_rxeqtraining_wait_time;
	u32 pcs_rxeqtraining_run_time;
	u8 _reserved7[4];
	u32 pcs_fll_ctrl1;
	u32 pcs_fll_ctrl2;
	u32 pcs_fll_cnt_val_l;
	u32 pcs_fll_cnt_val_h_tol;
	u32 pcs_fll_man_code;
	u32 pcs_autonomous_mode_ctrl;
	u8 _reserved8[152];
	u32 pcs_ready_status;
	u8 _reserved9[96];
	u32 pcs_rx_sigdet_lvl;
	u8 _reserved10[48];
	u32 pcs_refgen_req_config1;
	u32 pcs_refgen_req_config2;
};
check_member(usb3_phy_pcs_reg_layout, pcs_sw_reset, 0x000);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_v0, 0x00c);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_v1, 0x010);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_v2, 0x014);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_v3, 0x018);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_v4, 0x01c);
check_member(usb3_phy_pcs_reg_layout, pcs_txmgn_ls, 0x020);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_v0, 0x024);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_v0, 0x028);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_v1, 0x02c);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_v1, 0x030);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_v2, 0x034);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_v2, 0x038);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_v3, 0x03c);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_v3, 0x040);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_v4, 0x044);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_v4, 0x048);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m6db_ls, 0x04c);
check_member(usb3_phy_pcs_reg_layout, pcs_txdeemph_m3p5db_ls, 0x050);
check_member(usb3_phy_pcs_reg_layout, pcs_rate_slew_cntrl, 0x05c);
check_member(usb3_phy_pcs_reg_layout, pcs_power_state_config2, 0x064);
check_member(usb3_phy_pcs_reg_layout, pcs_rcvr_dtct_dly_p1u2_l, 0x070);
check_member(usb3_phy_pcs_reg_layout, pcs_rcvr_dtct_dly_p1u2_h, 0x074);
check_member(usb3_phy_pcs_reg_layout, pcs_rcvr_dtct_dly_u3_l, 0x078);
check_member(usb3_phy_pcs_reg_layout, pcs_rcvr_dtct_dly_u3_h, 0x07c);
check_member(usb3_phy_pcs_reg_layout, pcs_lock_detect_config1, 0x080);
check_member(usb3_phy_pcs_reg_layout, pcs_lock_detect_config2, 0x084);
check_member(usb3_phy_pcs_reg_layout, pcs_lock_detect_config3, 0x088);
check_member(usb3_phy_pcs_reg_layout, pcs_pwrup_reset_dly_time_auxclk, 0x0a0);
check_member(usb3_phy_pcs_reg_layout, pcs_rxeqtraining_wait_time, 0x0b8);
check_member(usb3_phy_pcs_reg_layout, pcs_fll_cnt_val_h_tol, 0x0d0);
check_member(usb3_phy_pcs_reg_layout, pcs_autonomous_mode_ctrl, 0x0d8);
check_member(usb3_phy_pcs_reg_layout, pcs_ready_status, 0x174);
check_member(usb3_phy_pcs_reg_layout, pcs_refgen_req_config2, 0x210);

static struct usb3_phy_qserdes_com_reg_layout *const qserdes_com_reg_layout =
	(void *)QMP_PHY_QSERDES_COM_REG_BASE;
static struct usb3_phy_qserdes_tx_reg_layout *const qserdes_tx_reg_layout =
	(void *)QMP_PHY_QSERDES_TX_REG_BASE;
static struct usb3_phy_qserdes_rx_reg_layout *const qserdes_rx_reg_layout =
	(void *)QMP_PHY_QSERDES_RX_REG_BASE;
static struct usb3_phy_pcs_reg_layout *const pcs_reg_layout =
	(void *)QMP_PHY_PCS_REG_BASE;

static struct usb3_phy_qserdes_com_reg_layout *const
	uniphy_qserdes_com_reg_layout =
	(void *)QMP_UNIPHY_QSERDES_COM_REG_BASE;
static struct usb3_phy_qserdes_tx_reg_layout
	*const uniphy_qserdes_tx_reg_layout =
	(void *)QMP_UNIPHY_QSERDES_TX_REG_BASE;
static struct usb3_phy_qserdes_rx_reg_layout
	*const uniphy_qserdes_rx_reg_layout =
	(void *)QMP_UNIPHY_QSERDES_RX_REG_BASE;
static struct usb3_phy_pcs_reg_layout *const uniphy_pcs_reg_layout =
	(void *)QMP_UNIPHY_PCS_REG_BASE;

struct usb_dwc3 {
	u32 sbuscfg0;
	u32 sbuscfg1;
	u32 txthrcfg;
	u32 rxthrcfg;
	u32 ctl;
	u32 pmsts;
	u32 sts;
	u32 uctl1;
	u32 snpsid;
	u32 gpio;
	u32 uid;
	u32 uctl;
	u64 buserraddr;
	u64 prtbimap;
	u8 reserved1[32];
	u32 dbgfifospace;
	u32 dbgltssm;
	u32 dbglnmcc;
	u32 dbgbmu;
	u32 dbglspmux;
	u32 dbglsp;
	u32 dbgepinfo0;
	u32 dbgepinfo1;
	u64 prtbimap_hs;
	u64 prtbimap_fs;
	u8 reserved2[112];
	u32 usb2phycfg;
	u8 reserved3[60];
	u32 usb2i2cctl;
	u8 reserved4[60];
	u32 usb2phyacc;
	u8 reserved5[60];
	u32 usb3pipectl;
	u8 reserved6[60];
};
check_member(usb_dwc3, usb3pipectl, 0x1c0);

static const struct qmp_phy_init_tbl qmp_v3_usb3_serdes_tbl[] = {
	{&qserdes_com_reg_layout->com_pll_ivco, 0x07},
	{&qserdes_com_reg_layout->com_sysclk_en_sel, 0x14},
	{&qserdes_com_reg_layout->com_bias_en_clkbuflr_en, 0x08},
	{&qserdes_com_reg_layout->com_clk_select, 0x30},
	{&qserdes_com_reg_layout->com_sys_clk_ctrl, 0x02},
	{&qserdes_com_reg_layout->com_resetsm_ctrl2, 0x08},
	{&qserdes_com_reg_layout->com_cmn_config, 0x16},
	{&qserdes_com_reg_layout->com_svs_mode_clk_sel, 0x01},
	{&qserdes_com_reg_layout->com_hsclk_sel, 0x80},
	{&qserdes_com_reg_layout->com_dec_start_mode0, 0x82},
	{&qserdes_com_reg_layout->com_div_frac_start1_mode0, 0xab},
	{&qserdes_com_reg_layout->com_div_frac_start2_mode0, 0xea},
	{&qserdes_com_reg_layout->com_div_frac_start3_mode0, 0x02},
	{&qserdes_com_reg_layout->com_cp_ctrl_mode0, 0x06},
	{&qserdes_com_reg_layout->com_pll_rctrl_mode0, 0x16},
	{&qserdes_com_reg_layout->com_pll_cctrl_mode0, 0x36},
	{&qserdes_com_reg_layout->com_integloop_gain1_mode0, 0x00},
	{&qserdes_com_reg_layout->com_integloop_gain0_mode0, 0x3f},
	{&qserdes_com_reg_layout->com_vco_tune2_mode0, 0x01},
	{&qserdes_com_reg_layout->com_vco_tune1_mode0, 0xc9},
	{&qserdes_com_reg_layout->com_coreclk_div_mode0, 0x0a},
	{&qserdes_com_reg_layout->com_lock_cmp3_mode0, 0x00},
	{&qserdes_com_reg_layout->com_lock_cmp2_mode0, 0x34},
	{&qserdes_com_reg_layout->com_lock_cmp1_mode0, 0x15},
	{&qserdes_com_reg_layout->com_lock_cmp_en, 0x04},
	{&qserdes_com_reg_layout->com_core_clk_en, 0x00},
	{&qserdes_com_reg_layout->com_lock_cmp_cfg, 0x00},
	{&qserdes_com_reg_layout->com_vco_tune_map, 0x00},
	{&qserdes_com_reg_layout->com_sysclk_buf_enable, 0x0a},
	{&qserdes_com_reg_layout->com_ssc_en_center, 0x01},
	{&qserdes_com_reg_layout->com_ssc_per1, 0x31},
	{&qserdes_com_reg_layout->com_ssc_per2, 0x01},
	{&qserdes_com_reg_layout->com_ssc_adj_per1, 0x00},
	{&qserdes_com_reg_layout->com_ssc_adj_per2, 0x00},
	{&qserdes_com_reg_layout->com_ssc_step_size1, 0x85},
	{&qserdes_com_reg_layout->com_ssc_step_size2, 0x07},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_tx_tbl[] = {
	{&qserdes_tx_reg_layout->tx_highz_drvr_en, 0x10},
	{&qserdes_tx_reg_layout->tx_rcv_detect_lvl_2, 0x12},
	{&qserdes_tx_reg_layout->tx_lane_mode_1, 0x16},
	{&qserdes_tx_reg_layout->tx_res_code_lane_offset_rx, 0x09},
	{&qserdes_tx_reg_layout->tx_res_code_lane_offset_tx, 0x06},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_rx_tbl[] = {
	{&qserdes_rx_reg_layout->rx_ucdr_fastlock_fo_gain, 0x0b},
	{&qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl2, 0x0f},
	{&qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl3, 0x4e},
	{&qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl4, 0x18},
	{&qserdes_rx_reg_layout->rx_rx_eq_offset_adap_ctrl1, 0x77},
	{&qserdes_rx_reg_layout->rx_rx_offset_adap_ctrl2, 0x80},
	{&qserdes_rx_reg_layout->rx_sigdet_ctrl, 0x03},
	{&qserdes_rx_reg_layout->rx_sigdet_deglitch_ctrl, 0x16},
	{&qserdes_rx_reg_layout->rx_ucdr_so_saturtn_and_en, 0x75},
	{&qserdes_rx_reg_layout->rx_ucdr_pi_cntrls, 0x80},
	{&qserdes_rx_reg_layout->rx_ucdr_fo_gain, 0x0a},
	{&qserdes_rx_reg_layout->rx_ucdr_so_gain_half, 0x06},
	{&qserdes_rx_reg_layout->rx_sigdet_enables, 0x00},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_pcs_tbl[] = {
	/* FLL settings */
	{&pcs_reg_layout->pcs_fll_ctrl2, 0x83},
	{&pcs_reg_layout->pcs_fll_cnt_val_l, 0x09},
	{&pcs_reg_layout->pcs_fll_cnt_val_h_tol, 0xa2},
	{&pcs_reg_layout->pcs_fll_man_code, 0x40},
	{&pcs_reg_layout->pcs_fll_ctrl1, 0x02},

	/* Lock Det settings */
	{&pcs_reg_layout->pcs_lock_detect_config1, 0xd1},
	{&pcs_reg_layout->pcs_lock_detect_config2, 0x1f},
	{&pcs_reg_layout->pcs_lock_detect_config3, 0x47},
	{&pcs_reg_layout->pcs_power_state_config2, 0x1b},

	{&pcs_reg_layout->pcs_rx_sigdet_lvl, 0xba},
	{&pcs_reg_layout->pcs_txmgn_v0, 0x9f},
	{&pcs_reg_layout->pcs_txmgn_v1, 0x9f},
	{&pcs_reg_layout->pcs_txmgn_v2, 0xb7},
	{&pcs_reg_layout->pcs_txmgn_v3, 0x4e},
	{&pcs_reg_layout->pcs_txmgn_v4, 0x65},
	{&pcs_reg_layout->pcs_txmgn_ls, 0x6b},
	{&pcs_reg_layout->pcs_txdeemph_m6db_v0, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_v0, 0x0d},
	{&pcs_reg_layout->pcs_txdeemph_m6db_v1, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_v1, 0x0d},
	{&pcs_reg_layout->pcs_txdeemph_m6db_v2, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_v2, 0x0d},
	{&pcs_reg_layout->pcs_txdeemph_m6db_v3, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_v3, 0x1d},
	{&pcs_reg_layout->pcs_txdeemph_m6db_v4, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_v4, 0x0d},
	{&pcs_reg_layout->pcs_txdeemph_m6db_ls, 0x15},
	{&pcs_reg_layout->pcs_txdeemph_m3p5db_ls, 0x0d},
	{&pcs_reg_layout->pcs_rate_slew_cntrl, 0x02},
	{&pcs_reg_layout->pcs_pwrup_reset_dly_time_auxclk, 0x04},
	{&pcs_reg_layout->pcs_tsync_rsync_time, 0x44},
	{&pcs_reg_layout->pcs_rcvr_dtct_dly_p1u2_l, 0xe7},
	{&pcs_reg_layout->pcs_rcvr_dtct_dly_p1u2_h, 0x03},
	{&pcs_reg_layout->pcs_rcvr_dtct_dly_u3_l, 0x40},
	{&pcs_reg_layout->pcs_rcvr_dtct_dly_u3_h, 0x00},
	{&pcs_reg_layout->pcs_rxeqtraining_wait_time, 0x75},
	{&pcs_reg_layout->pcs_lfps_ecstart_eqtlock, 0x86},
	{&pcs_reg_layout->pcs_rxeqtraining_run_time, 0x13},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_uniphy_serdes_tbl[] = {
	{&uniphy_qserdes_com_reg_layout->com_pll_ivco, 0x07},
	{&uniphy_qserdes_com_reg_layout->com_sysclk_en_sel, 0x14},
	{&uniphy_qserdes_com_reg_layout->com_bias_en_clkbuflr_en, 0x04},
	{&uniphy_qserdes_com_reg_layout->com_clk_select, 0x30},
	{&uniphy_qserdes_com_reg_layout->com_sys_clk_ctrl, 0x02},
	{&uniphy_qserdes_com_reg_layout->com_resetsm_ctrl2, 0x08},
	{&uniphy_qserdes_com_reg_layout->com_cmn_config, 0x06},
	{&uniphy_qserdes_com_reg_layout->com_svs_mode_clk_sel, 0x01},
	{&uniphy_qserdes_com_reg_layout->com_hsclk_sel, 0x80},
	{&uniphy_qserdes_com_reg_layout->com_dec_start_mode0, 0x82},
	{&uniphy_qserdes_com_reg_layout->com_div_frac_start1_mode0, 0xab},
	{&uniphy_qserdes_com_reg_layout->com_div_frac_start2_mode0, 0xea},
	{&uniphy_qserdes_com_reg_layout->com_div_frac_start3_mode0, 0x02},
	{&uniphy_qserdes_com_reg_layout->com_cp_ctrl_mode0, 0x06},
	{&uniphy_qserdes_com_reg_layout->com_pll_rctrl_mode0, 0x16},
	{&uniphy_qserdes_com_reg_layout->com_pll_cctrl_mode0, 0x36},
	{&uniphy_qserdes_com_reg_layout->com_integloop_gain1_mode0, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_integloop_gain0_mode0, 0x3f},
	{&uniphy_qserdes_com_reg_layout->com_vco_tune2_mode0, 0x01},
	{&uniphy_qserdes_com_reg_layout->com_vco_tune1_mode0, 0xc9},
	{&uniphy_qserdes_com_reg_layout->com_coreclk_div_mode0, 0x0a},
	{&uniphy_qserdes_com_reg_layout->com_lock_cmp3_mode0, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_lock_cmp2_mode0, 0x34},
	{&uniphy_qserdes_com_reg_layout->com_lock_cmp1_mode0, 0x15},
	{&uniphy_qserdes_com_reg_layout->com_lock_cmp_en, 0x04},
	{&uniphy_qserdes_com_reg_layout->com_core_clk_en, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_lock_cmp_cfg, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_vco_tune_map, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_sysclk_buf_enable, 0x0a},
	{&uniphy_qserdes_com_reg_layout->com_ssc_en_center, 0x01},
	{&uniphy_qserdes_com_reg_layout->com_ssc_per1, 0x31},
	{&uniphy_qserdes_com_reg_layout->com_ssc_per2, 0x01},
	{&uniphy_qserdes_com_reg_layout->com_ssc_adj_per1, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_ssc_adj_per2, 0x00},
	{&uniphy_qserdes_com_reg_layout->com_ssc_step_size1, 0x85},
	{&uniphy_qserdes_com_reg_layout->com_ssc_step_size2, 0x07},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_uniphy_tx_tbl[] = {
	{&uniphy_qserdes_tx_reg_layout->tx_highz_drvr_en, 0x10},
	{&uniphy_qserdes_tx_reg_layout->tx_rcv_detect_lvl_2, 0x12},
	{&uniphy_qserdes_tx_reg_layout->tx_lane_mode_1, 0xc6},
	{&uniphy_qserdes_tx_reg_layout->tx_res_code_lane_offset_rx, 0x06},
	{&uniphy_qserdes_tx_reg_layout->tx_res_code_lane_offset_tx, 0x06},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_uniphy_rx_tbl[] = {
	{&uniphy_qserdes_rx_reg_layout->rx_vga_cal_ctrl2, 0x0c},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_mode_00, 0x50},
	{&uniphy_qserdes_rx_reg_layout->rx_ucdr_fastlock_fo_gain, 0x0b},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl2, 0x0e},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl3, 0x4e},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_equ_adap_ctrl4, 0x18},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_eq_offset_adap_ctrl1, 0x77},
	{&uniphy_qserdes_rx_reg_layout->rx_rx_offset_adap_ctrl2, 0x80},
	{&uniphy_qserdes_rx_reg_layout->rx_sigdet_ctrl, 0x03},
	{&uniphy_qserdes_rx_reg_layout->rx_sigdet_deglitch_ctrl, 0x1c},
	{&uniphy_qserdes_rx_reg_layout->rx_ucdr_so_saturtn_and_en, 0x75},
	{&uniphy_qserdes_rx_reg_layout->rx_ucdr_pi_cntrls, 0x80},
	{&uniphy_qserdes_rx_reg_layout->rx_ucdr_fo_gain, 0x0a},
	{&uniphy_qserdes_rx_reg_layout->rx_ucdr_so_gain_half, 0x06},
	{&uniphy_qserdes_rx_reg_layout->rx_sigdet_enables, 0x00},
};

static const struct qmp_phy_init_tbl qmp_v3_usb3_uniphy_pcs_tbl[] = {
	/* FLL settings */
	{&uniphy_pcs_reg_layout->pcs_fll_ctrl2, 0x83},
	{&uniphy_pcs_reg_layout->pcs_fll_cnt_val_l, 0x09},
	{&uniphy_pcs_reg_layout->pcs_fll_cnt_val_h_tol, 0xa2},
	{&uniphy_pcs_reg_layout->pcs_fll_man_code, 0x40},
	{&uniphy_pcs_reg_layout->pcs_fll_ctrl1, 0x02},

	/* Lock Det settings */
	{&uniphy_pcs_reg_layout->pcs_lock_detect_config1, 0xd1},
	{&uniphy_pcs_reg_layout->pcs_lock_detect_config2, 0x1f},
	{&uniphy_pcs_reg_layout->pcs_lock_detect_config3, 0x47},
	{&uniphy_pcs_reg_layout->pcs_power_state_config2, 0x1b},

	{&uniphy_pcs_reg_layout->pcs_rx_sigdet_lvl, 0xba},
	{&uniphy_pcs_reg_layout->pcs_txmgn_v0, 0x9f},
	{&uniphy_pcs_reg_layout->pcs_txmgn_v1, 0x9f},
	{&uniphy_pcs_reg_layout->pcs_txmgn_v2, 0xb5},
	{&uniphy_pcs_reg_layout->pcs_txmgn_v3, 0x4c},
	{&uniphy_pcs_reg_layout->pcs_txmgn_v4, 0x64},
	{&uniphy_pcs_reg_layout->pcs_txmgn_ls, 0x6a},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_v0, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_v0, 0x0d},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_v1, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_v1, 0x0d},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_v2, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_v2, 0x0d},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_v3, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_v3, 0x1d},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_v4, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_v4, 0x0d},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m6db_ls, 0x15},
	{&uniphy_pcs_reg_layout->pcs_txdeemph_m3p5db_ls, 0x0d},
	{&uniphy_pcs_reg_layout->pcs_rate_slew_cntrl, 0x02},
	{&uniphy_pcs_reg_layout->pcs_pwrup_reset_dly_time_auxclk, 0x04},
	{&uniphy_pcs_reg_layout->pcs_tsync_rsync_time, 0x44},
	{&uniphy_pcs_reg_layout->pcs_rcvr_dtct_dly_p1u2_l, 0xe7},
	{&uniphy_pcs_reg_layout->pcs_rcvr_dtct_dly_p1u2_h, 0x03},
	{&uniphy_pcs_reg_layout->pcs_rcvr_dtct_dly_u3_l, 0x40},
	{&uniphy_pcs_reg_layout->pcs_rcvr_dtct_dly_u3_h, 0x00},
	{&uniphy_pcs_reg_layout->pcs_rxeqtraining_wait_time, 0x75},
	{&uniphy_pcs_reg_layout->pcs_lfps_ecstart_eqtlock, 0x86},
	{&uniphy_pcs_reg_layout->pcs_rxeqtraining_run_time, 0x13},
	{&uniphy_pcs_reg_layout->pcs_refgen_req_config1, 0x21},
	{&uniphy_pcs_reg_layout->pcs_refgen_req_config2, 0x60},
};

struct usb_dwc3_cfg {
	struct usb_dwc3 *usb_host_dwc3;
	struct usb_qusb_phy_pll *qusb_phy_pll;
	struct usb_qusb_phy_dig *qusb_phy_dig;
	/* Init sequence for QMP PHY blocks - serdes, tx, rx, pcs */
	const struct qmp_phy_init_tbl *serdes_tbl;
	int serdes_tbl_num;
	const struct qmp_phy_init_tbl *tx_tbl;
	int tx_tbl_num;
	const struct qmp_phy_init_tbl *rx_tbl;
	int rx_tbl_num;
	const struct qmp_phy_init_tbl *pcs_tbl;
	int pcs_tbl_num;
	struct usb3_phy_pcs_reg_layout *qmp_pcs_reg;

	u32 *usb3_bcr;
	u32 *qusb2phy_bcr;
	u32 *gcc_usb3phy_bcr_reg;
	u32 *gcc_qmpphy_bcr_reg;
	struct usb_board_data *board_data;
	u32 efuse_offset;
};

static struct usb_dwc3_cfg usb_port0 = {
	.usb_host_dwc3 =	(void *)USB_HOST0_DWC3_BASE,
	.qusb_phy_pll =		(void *)QUSB_PRIM_PHY_BASE,
	.qusb_phy_dig =		(void *)QUSB_PRIM_PHY_DIG_BASE,
	.serdes_tbl =		qmp_v3_usb3_serdes_tbl,
	.serdes_tbl_num	=	ARRAY_SIZE(qmp_v3_usb3_serdes_tbl),
	.tx_tbl =		qmp_v3_usb3_tx_tbl,
	.tx_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_tx_tbl),
	.rx_tbl =		qmp_v3_usb3_rx_tbl,
	.rx_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_rx_tbl),
	.pcs_tbl =		qmp_v3_usb3_pcs_tbl,
	.pcs_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_pcs_tbl),
	.qmp_pcs_reg =		(void *)QMP_PHY_PCS_REG_BASE,
	.usb3_bcr =		&gcc->usb30_prim_bcr,
	.qusb2phy_bcr =		&gcc->qusb2phy_prim_bcr,
	.gcc_usb3phy_bcr_reg =	&gcc->usb3_dp_phy_prim_bcr,
	.gcc_qmpphy_bcr_reg =	&gcc->usb3_phy_prim_bcr,
	.efuse_offset =		25,
};
static struct usb_dwc3_cfg usb_port1 = {
	.usb_host_dwc3 =	(void *)USB_HOST1_DWC3_BASE,
	.qusb_phy_pll =		(void *)QUSB_SEC_PHY_BASE,
	.qusb_phy_dig =		(void *)QUSB_SEC_PHY_DIG_BASE,
	.serdes_tbl =		qmp_v3_usb3_uniphy_serdes_tbl,
	.serdes_tbl_num	=	ARRAY_SIZE(qmp_v3_usb3_uniphy_serdes_tbl),
	.tx_tbl =		qmp_v3_usb3_uniphy_tx_tbl,
	.tx_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_uniphy_tx_tbl),
	.rx_tbl =		qmp_v3_usb3_uniphy_rx_tbl,
	.rx_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_uniphy_rx_tbl),
	.pcs_tbl =		qmp_v3_usb3_uniphy_pcs_tbl,
	.pcs_tbl_num =		ARRAY_SIZE(qmp_v3_usb3_uniphy_pcs_tbl),
	.qmp_pcs_reg =		(void *)QMP_UNIPHY_PCS_REG_BASE,
	.usb3_bcr =		&gcc->usb30_sec_bcr,
	.qusb2phy_bcr =		&gcc->qusb2phy_sec_bcr,
	.gcc_usb3phy_bcr_reg =	&gcc->usb3phy_phy_sec_bcr,
	.gcc_qmpphy_bcr_reg =	&gcc->usb3_phy_sec_bcr,
	.efuse_offset =		30,
};

static struct qfprom_corr * const qfprom_corr_efuse = (void *)QFPROM_BASE;

static void reset_usb(struct usb_dwc3_cfg *dwc3)
{
	/* Assert Core reset */
	clock_reset_bcr(dwc3->usb3_bcr, 1);

	/* Assert QUSB PHY reset */
	clock_reset_bcr(dwc3->qusb2phy_bcr, 1);

	/* Assert QMP PHY reset */
	clock_reset_bcr(dwc3->gcc_usb3phy_bcr_reg, 1);
	clock_reset_bcr(dwc3->gcc_qmpphy_bcr_reg, 1);
}

void reset_usb0(void)
{
	/* Before Resetting PHY, put Core in Reset */
	printk(BIOS_INFO, "Starting DWC3 and PHY resets for USB(0)\n");

	reset_usb(&usb_port0);
}

void reset_usb1(void)
{
	/* Before Resetting PHY, put Core in Reset */
	printk(BIOS_INFO, "Starting DWC3 and PHY resets for USB(1)\n");

	reset_usb(&usb_port1);
}
/*
 * Update board specific PHY tuning override values that specified from
 * board file.
 */
static void qusb2_phy_override_phy_params(struct usb_dwc3_cfg *dwc3)
{
	/* Override preemphasis value */
	write32(&dwc3->qusb_phy_dig->tune1,
		dwc3->board_data->port_tune1);

	/* Override BIAS_CTRL_2 to reduce the TX swing overshooting. */
	write32(&dwc3->qusb_phy_pll->bias_ctrl_2,
			dwc3->board_data->pll_bias_control_2);

	/* Override IMP_RES_OFFSET value */
	write32(&dwc3->qusb_phy_dig->imp_ctrl1,
		dwc3->board_data->imp_ctrl1);
}

/*
 * Fetches HS Tx tuning value from efuse register and sets the
 * QUSB2PHY_PORT_TUNE1/2 register.
 * For error case, skip setting the value and use the default value.
 */
static void qusb2_phy_set_tune_param(struct usb_dwc3_cfg *dwc3)
{
	/*
	 * Efuse registers 4 bit value specifies tuning for HSTX
	 * output current in TUNE1 Register. Hence Extract 4 bits from
	 * EFUSE at correct position.
	 */

	const int efuse_bits = 4;
	int bit_pos = dwc3->efuse_offset;

	u32 bit_mask = (1 << efuse_bits) - 1;
	u32 tune_val =
		(read32(&qfprom_corr_efuse->qusb_hstx_trim_lsb) >> bit_pos)
		& bit_mask;

	if (bit_pos + efuse_bits > 32) {
		/*
		 * Value split between two EFUSE registers,
		 * get the second part.
		 */
		int done_bits = 32 - bit_pos;

		bit_mask = (1 << (efuse_bits - done_bits)) - 1;
		tune_val |=
			(read32(&qfprom_corr_efuse->qusb_hstx_trim_msb) &
			bit_mask) << done_bits;
	}

	/*
	 * if efuse reg is updated (i.e non-zero) then use it to program
	 * tune parameters.
	 */
	if (tune_val)
		clrsetbits_le32(&dwc3->qusb_phy_dig->tune1,
				PORT_TUNE1_MASK, tune_val << 4);
}

static void tune_phy(struct usb_dwc3_cfg *dwc3, struct usb_qusb_phy_dig *phy)
{
	write32(&phy->pwr_ctrl2, QUSB2PHY_PWR_CTRL2);
	/* IMP_CTRL1: Control the impedance reduction */
	write32(&phy->imp_ctrl1, QUSB2PHY_IMP_CTRL1);
	/* IMP_CTRL2: Impedance offset/mapping slope */
	write32(&phy->imp_ctrl2, QUSB2PHY_IMP_CTRL1);
	write32(&phy->chg_ctrl2, QUSB2PHY_IMP_CTRL2);
	/*
	 * TUNE1: Sets HS Impedance to approx 45 ohms
	 * then override with efuse value.
	 */
	write32(&phy->tune1, QUSB2PHY_PORT_TUNE1);
	/* TUNE2: Tuning for HS Disconnect Level */
	write32(&phy->tune2, QUSB2PHY_PORT_TUNE2);
	/* TUNE3: Tune squelch range */
	write32(&phy->tune3, QUSB2PHY_PORT_TUNE3);
	/* TUNE4: Sets EOP_DLY(Squelch rising edge to linestate falling edge) */
	write32(&phy->tune4, QUSB2PHY_PORT_TUNE4);
	write32(&phy->tune5, QUSB2PHY_PORT_TUNE5);

	if (dwc3->board_data) {
		/* Override board specific PHY tuning values */
		qusb2_phy_override_phy_params(dwc3);

		/* Set efuse value for tuning the PHY */
		qusb2_phy_set_tune_param(dwc3);
	}
}

static void hs_qusb_phy_init(struct usb_dwc3_cfg *dwc3)
{
	/* PWR_CTRL: set the power down bit to disable the PHY */
	setbits_le32(&dwc3->qusb_phy_dig->pwr_ctrl1, POWER_DOWN);

	write32(&dwc3->qusb_phy_pll->analog_controls_two,
			QUSB2PHY_PLL_ANALOG_CONTROLS_TWO);
	write32(&dwc3->qusb_phy_pll->clock_inverters,
			QUSB2PHY_PLL_CLOCK_INVERTERS);
	write32(&dwc3->qusb_phy_pll->cmode,
			QUSB2PHY_PLL_CMODE);
	write32(&dwc3->qusb_phy_pll->lock_delay,
			QUSB2PHY_PLL_LOCK_DELAY);
	write32(&dwc3->qusb_phy_pll->dig_tim,
			QUSB2PHY_PLL_DIGITAL_TIMERS_TWO);
	write32(&dwc3->qusb_phy_pll->bias_ctrl_1,
			QUSB2PHY_PLL_BIAS_CONTROL_1);
	write32(&dwc3->qusb_phy_pll->bias_ctrl_2,
			QUSB2PHY_PLL_BIAS_CONTROL_2);

	tune_phy(dwc3, dwc3->qusb_phy_dig);

	/* PWR_CTRL1: Clear the power down bit to enable the PHY */
	clrbits_le32(&dwc3->qusb_phy_dig->pwr_ctrl1, POWER_DOWN);

	write32(&dwc3->qusb_phy_dig->debug_ctrl2,
				DEBUG_CTRL2_MUX_PLL_LOCK_STATUS);

	/*
	 * DEBUG_STAT5: wait for 160uS for PLL lock;
	 * vstatus[0] changes from 0 to 1.
	 */
	long lock_us = wait_us(160, read32(&dwc3->qusb_phy_dig->debug_stat5) &
						VSTATUS_PLL_LOCK_STATUS_MASK);
	if (!lock_us)
		printk(BIOS_ERR, "ERROR: QUSB PHY PLL LOCK fails\n");
	else
		printk(BIOS_DEBUG, "QUSB PHY initialized and locked in %ldus\n",
				lock_us);
}

static void qcom_qmp_phy_configure(const struct qmp_phy_init_tbl tbl[],
				int num)
{
	int i;
	const struct qmp_phy_init_tbl *t = tbl;

	if (!t)
		return;

	for (i = 0; i < num; i++, t++)
		write32(t->address, t->val);
}

static void ss_qmp_phy_init(struct usb_dwc3_cfg *dwc3)
{
	/* power up USB3 PHY */
	write32(&dwc3->qmp_pcs_reg->pcs_power_down_control, 0x01);

	 /* Serdes configuration */
	qcom_qmp_phy_configure(dwc3->serdes_tbl, dwc3->serdes_tbl_num);
	/* Tx, Rx, and PCS configurations */
	qcom_qmp_phy_configure(dwc3->tx_tbl, dwc3->tx_tbl_num);
	qcom_qmp_phy_configure(dwc3->rx_tbl, dwc3->rx_tbl_num);
	qcom_qmp_phy_configure(dwc3->pcs_tbl, dwc3->pcs_tbl_num);

	/* perform software reset of PCS/Serdes */
	write32(&dwc3->qmp_pcs_reg->pcs_sw_reset, 0x00);
	/* start PCS/Serdes to operation mode */
	write32(&dwc3->qmp_pcs_reg->pcs_start_control, 0x03);

	/*
	 * Wait for PHY initialization to be done
	 * PCS_STATUS: wait for 1ms for PHY STATUS;
	 * SW can continuously check for PHYSTATUS = 1.b0.
	 */
	long lock_us = wait_us(1000,
			!(read32(&dwc3->qmp_pcs_reg->pcs_ready_status) &
			USB3_PCS_PHYSTATUS));
	if (!lock_us)
		printk(BIOS_ERR, "ERROR: QMP PHY PLL LOCK fails:\n");
	else
		printk(BIOS_DEBUG, "QMP PHY initialized and locked in %ldus\n",
				lock_us);
}

static void setup_dwc3(struct usb_dwc3 *dwc3)
{
	/* core exits U1/U2/U3 only in PHY power state P1/P2/P3 respectively */
	clrsetbits_le32(&dwc3->usb3pipectl,
		DWC3_GUSB3PIPECTL_DELAYP1TRANS,
		DWC3_GUSB3PIPECTL_UX_EXIT_IN_PX);

	/*
	 * Configure USB phy interface of DWC3 core.
	 * 1. Select UTMI+ PHY with 16-bit interface.
	 * 2. Set USBTRDTIM to the corresponding value
	 * according to the UTMI+ PHY interface.
	 */
	clrsetbits_le32(&dwc3->usb2phycfg,
			(DWC3_GUSB2PHYCFG_USB2TRDTIM_MASK |
			DWC3_GUSB2PHYCFG_PHYIF_MASK),
			(DWC3_GUSB2PHYCFG_PHYIF(UTMI_PHYIF_8_BIT) |
			DWC3_GUSB2PHYCFG_USBTRDTIM(USBTRDTIM_UTMI_8_BIT)));

	clrsetbits_le32(&dwc3->ctl, (DWC3_GCTL_SCALEDOWN_MASK |
			DWC3_GCTL_DISSCRAMBLE),
			DWC3_GCTL_U2EXIT_LFPS | DWC3_GCTL_DSBLCLKGTNG);

	/* configure controller in Host mode */
	clrsetbits_le32(&dwc3->ctl, (DWC3_GCTL_PRTCAPDIR(DWC3_GCTL_PRTCAP_OTG)),
			DWC3_GCTL_PRTCAPDIR(DWC3_GCTL_PRTCAP_HOST));
	printk(BIOS_SPEW, "Configure USB in Host mode\n");
}

/* Initialization of DWC3 Core and PHY */
static void setup_usb_host(struct usb_dwc3_cfg *dwc3,
				struct usb_board_data *board_data)
{
	dwc3->board_data = board_data;

	 /* Clear core reset. */
	clock_reset_bcr(dwc3->usb3_bcr, 0);

	/* Clear QUSB PHY reset. */
	clock_reset_bcr(dwc3->qusb2phy_bcr, 0);

	/* Initialize QUSB PHY */
	hs_qusb_phy_init(dwc3);

	/* Clear QMP PHY resets. */
	clock_reset_bcr(dwc3->gcc_usb3phy_bcr_reg, 0);
	clock_reset_bcr(dwc3->gcc_qmpphy_bcr_reg, 0);

	/* Initialize QMP PHY */
	ss_qmp_phy_init(dwc3);

	setup_dwc3(dwc3->usb_host_dwc3);

	printk(BIOS_INFO, "DWC3 and PHY setup finished\n");
}

void setup_usb_host0(struct usb_board_data *board_data)
{
	printk(BIOS_INFO, "Setting up USB HOST0 controller.\n");
	setup_usb_host(&usb_port0, board_data);
}

void setup_usb_host1(struct usb_board_data *board_data)
{
	printk(BIOS_INFO, "Setting up USB HOST1 controller.\n");
	setup_usb_host(&usb_port1, board_data);
}
