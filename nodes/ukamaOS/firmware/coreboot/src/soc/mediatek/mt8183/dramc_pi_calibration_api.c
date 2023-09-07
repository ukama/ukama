/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 MediaTek Inc.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <assert.h>
#include <console/console.h>
#include <delay.h>
#include <device/mmio.h>
#include <soc/emi.h>
#include <soc/dramc_register.h>
#include <soc/dramc_pi_api.h>
#include <timer.h>

enum {
	RX_VREF_BEGIN = 0,
	RX_VREF_END = 31,
	RX_VREF_STEP = 1,
	TX_VREF_BEGIN = 0,
	TX_VREF_END = 50,
	TX_VREF_STEP = 2,
};

enum {
	FIRST_DQ_DELAY = 0,
	FIRST_DQS_DELAY = -48,
	MAX_DQDLY_TAPS = 16,
	MAX_RX_DQDLY_TAPS = 63,
};

#define WRITE_LEVELING_MOVD_DQS 1
#define TEST2_1_CAL 0x55000000
#define TEST2_2_CAL 0xaa000400

enum CAL_TYPE {
	RX_WIN_RD_DQC = 0,
	RX_WIN_TEST_ENG,
	TX_WIN_DQ_ONLY,
	TX_WIN_DQ_DQM,
};

enum RX_TYPE {
	RX_DQ = 0,
	RX_DQM,
	RX_DQS,
};

struct win_perbit_dly {
	s16 first_pass;
	s16 last_pass;
	s16 best_first;
	s16 best_last;
	s16 best_dqdly;
	s16 win_center;
};

struct vref_perbit_dly {
	u8 best_vref;
	u16 max_win_sum;
	struct win_perbit_dly perbit_dly[DQ_DATA_WIDTH];
};

struct tx_dly_tune {
	u8 fine_tune;
	u8 coarse_tune_large;
	u8 coarse_tune_small;
	u8 coarse_tune_large_oen;
	u8 coarse_tune_small_oen;
};

struct per_byte_dly {
	u16 max_center;
	u16 min_center;
	u16 final_dly;
};

static void dramc_auto_refresh_switch(u8 chn, bool option)
{
	SET32_BITFIELDS(&ch[chn].ao.refctrl0, REFCTRL0_REFDIS, option ? 0 : 1);

	if (!option) {
		/*
		 * Because HW will actually disable autorefresh after
		 * refresh_queue empty, we need to wait until queue empty.
		 */
		udelay(READ32_BITFIELD(&ch[chn].nao.misc_statusa,
				       MISC_STATUSA_REFRESH_QUEUE_CNT) * 4);
	}
}

void dramc_cke_fix_onoff(u8 chn, bool cke_on, bool cke_off)
{
	SET32_BITFIELDS(&ch[chn].ao.ckectrl,
			CKECTRL_CKEFIXON, cke_on,
			CKECTRL_CKEFIXOFF, cke_off);
}

void dramc_mode_reg_write(u8 chn, u8 mr_idx, u8 value)
{
	u32 ckectrl_bak = read32(&ch[chn].ao.ckectrl);

	dramc_cke_fix_onoff(chn, true, false);
	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSMA, mr_idx);
	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSOP, value);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_MRWEN, 1);

	/* Wait MRW command fired */
	while (READ32_BITFIELD(&ch[chn].nao.spcmdresp, SPCMDRESP_MRW_RESPONSE)
	       == 0)
		;

	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_MRWEN, 0);
	write32(&ch[chn].ao.ckectrl, ckectrl_bak);
	dramc_dbg("Write MR%d =0x%x\n", mr_idx, value);
}

static void dramc_mode_reg_write_by_rank(u8 chn, u8 rank,
		u8 mr_idx, u8 value)
{
	u32 mrs_bak = READ32_BITFIELD(&ch[chn].ao.mrs, MRS_MRSRK);
	dramc_dbg("Mode reg write rank%d MR%d = 0x%x\n", rank, mr_idx, value);

	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSRK, rank);
	dramc_mode_reg_write(chn, mr_idx, value);
	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSRK, mrs_bak);
}

static void move_dramc_delay(u32 *reg_0, u32 *reg_1, u8 shift, s8 shift_coarse_tune)
{
	s32 sum;
	u32 tmp_0p5t, tmp_2t;

	tmp_0p5t = ((read32(reg_0) >> shift) & DQ_DIV_MASK) &
		~(1 << DQ_DIV_SHIFT);
	tmp_2t = (read32(reg_1) >> shift) & DQ_DIV_MASK;

	sum = (tmp_2t << DQ_DIV_SHIFT) + tmp_0p5t + shift_coarse_tune;

	if (sum < 0) {
		tmp_0p5t = 0;
		tmp_2t = 0;
	} else {
		tmp_2t = sum >> DQ_DIV_SHIFT;
		tmp_0p5t = sum - (tmp_2t << DQ_DIV_SHIFT);
	}

	clrsetbits_le32(reg_0, DQ_DIV_MASK << shift, tmp_0p5t << shift);
	clrsetbits_le32(reg_1, DQ_DIV_MASK << shift, tmp_2t << shift);
}

static void move_dramc_tx_dqs(u8 chn, u8 byte, s8 shift_coarse_tune)
{
	move_dramc_delay(&ch[chn].ao.shu[0].selph_dqs1,
		&ch[chn].ao.shu[0].selph_dqs0, byte * 4, shift_coarse_tune);
}

static void move_dramc_tx_dqs_oen(u8 chn, u8 byte,
		s8 shift_coarse_tune)
{
	move_dramc_delay(&ch[chn].ao.shu[0].selph_dqs1,
		&ch[chn].ao.shu[0].selph_dqs0, byte * 4 + OEN_SHIFT, shift_coarse_tune);
}

static void move_dramc_tx_dq(u8 chn, u8 rank, u8 byte, s8 shift_coarse_tune)
{
	/* DQM0 */
	move_dramc_delay(&ch[chn].ao.shu[0].rk[rank].selph_dq[3],
		&ch[chn].ao.shu[0].rk[rank].selph_dq[1], byte * 4, shift_coarse_tune);

	/* DQ0 */
	move_dramc_delay(&ch[chn].ao.shu[0].rk[rank].selph_dq[2],
		&ch[chn].ao.shu[0].rk[rank].selph_dq[0], byte * 4, shift_coarse_tune);
}

static void move_dramc_tx_dq_oen(u8 chn, u8 rank,
	u8 byte, s8 shift_coarse_tune)
{
	/* DQM_OEN_0 */
	move_dramc_delay(&ch[chn].ao.shu[0].rk[rank].selph_dq[3],
		&ch[chn].ao.shu[0].rk[rank].selph_dq[1],
		byte * 4 + OEN_SHIFT, shift_coarse_tune);

	/* DQ_OEN_0 */
	move_dramc_delay(&ch[chn].ao.shu[0].rk[rank].selph_dq[2],
		&ch[chn].ao.shu[0].rk[rank].selph_dq[0],
		byte * 4 + OEN_SHIFT, shift_coarse_tune);
}

static void write_leveling_move_dqs_instead_of_clk(u8 chn)
{
	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		move_dramc_tx_dqs(chn, byte, -WRITE_LEVELING_MOVD_DQS);
		move_dramc_tx_dqs_oen(chn, byte, -WRITE_LEVELING_MOVD_DQS);

		for (u8 rk = RANK_0; rk < RANK_MAX; rk++) {
			move_dramc_tx_dq(chn, rk, byte, -WRITE_LEVELING_MOVD_DQS);
			move_dramc_tx_dq_oen(chn, rk, byte, -WRITE_LEVELING_MOVD_DQS);
		}
	}
}

static void dramc_write_leveling(u8 chn, u8 rank, u8 freq_group,
		const u8 wr_level[CHANNEL_MAX][RANK_MAX][DQS_NUMBER])
{
	dramc_auto_refresh_switch(chn, false);

	if (rank == RANK_0 && (freq_group == LP4X_DDR3600 ||
			       freq_group == LP4X_DDR1600 ||
			       freq_group == LP4X_DDR2400))
		write_leveling_move_dqs_instead_of_clk(chn);

	SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].ca_cmd[9],
			SHU1_R0_CA_CMD9_RG_RK0_ARPI_CLK, 0);

	for (size_t byte = 0; byte < DQS_NUMBER; byte++) {
		u32 wrlevel_dq_delay = wr_level[chn][rank][byte] + 0x10;
		SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[7],
				FINE_TUNE_PBYTE, wr_level[chn][rank][byte]);
		if (wrlevel_dq_delay >= 0x40) {
			wrlevel_dq_delay -= 0x40;
			move_dramc_tx_dq(chn, rank, byte, 2);
			move_dramc_tx_dq_oen(chn, rank, byte, 2);
		}
		SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[7],
				FINE_TUNE_DQM, wrlevel_dq_delay,
				FINE_TUNE_DQ, wrlevel_dq_delay);
	}
}

static void dramc_cmd_bus_training(u8 chn, u8 rank, u8 freq_group,
		const struct sdram_params *params, const bool fast_calib)
{
	u32 final_vref, clk_dly, cmd_dly, cs_dly;

	clk_dly = params->cbt_clk_dly[chn][rank];
	cmd_dly = params->cbt_cmd_dly[chn][rank];
	cs_dly = params->cbt_cs_dly[chn][rank];
	final_vref = params->cbt_final_vref[chn][rank];

	if (fast_calib) {
		/* Set CLK and CA delay */
		SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].ca_cmd[9],
				SHU1_R0_CA_CMD9_RG_RK0_ARPI_CMD, cmd_dly,
				SHU1_R0_CA_CMD9_RG_RK0_ARPI_CLK, clk_dly);
		udelay(1);
	}

	/* Set CLK and CS delay */
	SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].ca_cmd[9],
			SHU1_R0_CA_CMD9_RG_RK0_ARPI_CS, cs_dly);

	/* CBT set vref */
	dramc_mode_reg_write_by_rank(chn, rank, 12, final_vref);
}

static void dramc_read_dbi_onoff(bool on)
{
	for (size_t chn = 0; chn < CHANNEL_MAX; chn++)
		for (size_t b = 0; b < 2; b++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].b[b].dq[7],
					SHU1_B0_DQ7_R_DMDQMDBI_SHU_B0, on);
}

static void dramc_write_dbi_onoff(bool onoff)
{
	for (size_t chn = 0; chn < CHANNEL_MAX; chn++)
		SET32_BITFIELDS(&ch[chn].ao.shu[0].wodt, SHU1_WODT_DBIWR, onoff);
}

static void dramc_phy_dcm_2_channel(u8 chn, bool en)
{
	clrsetbits_le32(&ch[chn].phy.misc_cg_ctrl0, (0x3 << 19) | (0x3ff << 8),
		((en ? 0 : 0x1) << 19) | ((en ? 0 : 0x1ff) << 9) | (1 << 8));

	for (size_t i = 0; i < DRAM_DFS_SHUFFLE_MAX; i++) {
		struct ddrphy_ao_shu *shu = &ch[chn].phy.shu[i];
		for (size_t b = 0; b < 2; b++)
			clrsetbits_le32(&shu->b[b].dq[8], 0x1fff << 19,
				((en ? 0 : 0x7ff) << 22) | (0x1 << 21) |
				((en ? 0 : 0x3) << 19));
		clrbits_le32(&shu->ca_cmd[8], 0x1fff << 19);
	}
	clrsetbits_le32(&ch[chn].phy.misc_cg_ctrl5, (0x7 << 16) | (0x7 << 20),
		((en ? 0x7 : 0) << 16) | ((en ? 0x7 : 0) << 20));
}

void dramc_enable_phy_dcm(bool en)
{
	for (size_t chn = 0; chn < CHANNEL_MAX ; chn++) {
		clrbits_le32(&ch[chn].phy.b[0].dll_fine_tune[1], 0x1 << 20);
		clrbits_le32(&ch[chn].phy.b[1].dll_fine_tune[1], 0x1 << 20);
		clrbits_le32(&ch[chn].phy.ca_dll_fine_tune[1], 0x1 << 20);

		for (size_t i = 0; i < DRAM_DFS_SHUFFLE_MAX; i++) {
			struct ddrphy_ao_shu *shu = &ch[chn].phy.shu[i];
			setbits_le32(&shu->b[0].dll[0], 0x1);
			setbits_le32(&shu->b[1].dll[0], 0x1);
			setbits_le32(&shu->ca_dll[0], 0x1);
		}

		clrsetbits_le32(&ch[chn].ao.dramc_pd_ctrl,
			(0x1 << 0) | (0x1 << 1) | (0x1 << 2) |
			(0x1 << 5) | (0x1 << 26) | (0x1 << 30) | (0x1 << 31),
			((en ? 0x1 : 0) << 0) | ((en ? 0x1 : 0) << 1) |
			((en ? 0x1 : 0) << 2) | ((en ? 0 : 0x1) << 5) |
			((en ? 0 : 0x1) << 26) | ((en ? 0x1 : 0) << 30) |
			((en ? 0x1 : 0) << 31));

		/* DCM on: CHANNEL_EMI free run; DCM off: mem_dcm */
		write32(&ch[chn].phy.misc_cg_ctrl2,
			0x8060033e | (0x40 << (en ? 0x1 : 0)));
		write32(&ch[chn].phy.misc_cg_ctrl2,
			0x8060033f | (0x40 << (en ? 0x1 : 0)));
		write32(&ch[chn].phy.misc_cg_ctrl2,
			0x8060033e | (0x40 << (en ? 0x1 : 0)));

		clrsetbits_le32(&ch[chn].phy.misc_ctrl3, 0x3 << 26,
			(en ? 0 : 0x3) << 26);
		for (size_t i = 0; i < DRAM_DFS_SHUFFLE_MAX; i++) {
			u32 mask = 0x7 << 17;
			u32 value = (en ? 0x7 : 0) << 17;
			struct ddrphy_ao_shu *shu = &ch[chn].phy.shu[i];

			clrsetbits_le32(&shu->b[0].dq[7], mask, value);
			clrsetbits_le32(&shu->b[1].dq[7], mask, value);
			clrsetbits_le32(&shu->ca_cmd[7], mask, value);
		}

		dramc_phy_dcm_2_channel(chn, en);
	}
}

static void dramc_reset_delay_chain_before_calibration(void)
{
	for (size_t chn = 0; chn < CHANNEL_MAX; chn++)
		for (size_t rank = 0; rank < RANK_MAX; rank++) {
			struct dramc_ddrphy_regs_shu_rk *rk;
			rk = &ch[chn].phy.shu[0].rk[rank];
			clrbits_le32(&rk->ca_cmd[0], 0xffffff << 0);
			clrbits_le32(&rk->b[0].dq[0], 0xfffffff << 0);
			clrbits_le32(&rk->b[1].dq[0], 0xfffffff << 0);
			clrbits_le32(&rk->b[0].dq[1], 0xf << 0);
			clrbits_le32(&rk->b[1].dq[1], 0xf << 0);
		}
}

void dramc_hw_gating_onoff(u8 chn, bool on)
{
	clrsetbits_le32(&ch[chn].ao.shuctrl2, 0x3 << 14,
		(on ? 0x3 : 0) << 14);
	clrsetbits_le32(&ch[chn].ao.stbcal2, 0x1 << 28, (on ? 0x1 : 0) << 28);
	clrsetbits_le32(&ch[chn].ao.stbcal, 0x1 << 24, (on ? 0x1 : 0) << 24);
	clrsetbits_le32(&ch[chn].ao.stbcal, 0x1 << 22, (on ? 0x1 : 0) << 22);
}

static void dramc_rx_input_delay_tracking_init_by_freq(u8 chn)
{
	struct ddrphy_ao_shu *shu = &ch[chn].phy.shu[0];

	clrsetbits_le32(&shu->b[0].dq[5], 0x7 << 20, 0x3 << 20);
	clrsetbits_le32(&shu->b[1].dq[5], 0x7 << 20, 0x3 << 20);
	clrbits_le32(&shu->b[0].dq[7], (0x1 << 12) | (0x1 << 13));
	clrbits_le32(&shu->b[1].dq[7], (0x1 << 12) | (0x1 << 13));
}

void dramc_apply_config_before_calibration(u8 freq_group)
{
	dramc_enable_phy_dcm(false);
	dramc_reset_delay_chain_before_calibration();

	setbits_le32(&ch[0].ao.shu[0].conf[3], 0x1ff << 16);
	setbits_le32(&ch[0].ao.spcmdctrl, 0x1 << 24);
	clrsetbits_le32(&ch[0].ao.shu[0].scintv, 0x1f << 1, 0x1b << 1);

	for (size_t shu = DRAM_DFS_SHUFFLE_1; shu < DRAM_DFS_SHUFFLE_MAX; shu++)
		setbits_le32(&ch[0].ao.shu[shu].conf[3], 0x1ff << 0);

	clrbits_le32(&ch[0].ao.dramctrl, 0x1 << 18);
	clrbits_le32(&ch[0].ao.spcmdctrl, 0x1 << 31);
	clrbits_le32(&ch[0].ao.spcmdctrl, 0x1 << 30);
	clrbits_le32(&ch[0].ao.dqsoscr, 0x1 << 26);
	clrbits_le32(&ch[0].ao.dqsoscr, 0x1 << 25);

	dramc_write_dbi_onoff(false);
	dramc_read_dbi_onoff(false);

	for (size_t chn = 0; chn < CHANNEL_MAX; chn++) {
		setbits_le32(&ch[chn].ao.spcmdctrl, 0x1 << 29);
		setbits_le32(&ch[chn].ao.dqsoscr, 0x1 << 24);
		for (size_t shu = DRAM_DFS_SHUFFLE_1; shu < DRAM_DFS_SHUFFLE_MAX; shu++)
			setbits_le32(&ch[chn].ao.shu[shu].scintv, 0x1 << 30);

		clrbits_le32(&ch[chn].ao.dummy_rd, (0x1 << 7) | (0x7 << 20));
		dramc_hw_gating_onoff(chn, false);
		clrbits_le32(&ch[chn].ao.stbcal2, 0x1 << 28);

		setbits_le32(&ch[chn].phy.misc_ctrl1, (0x1 << 7) | (0x1 << 11));
		clrbits_le32(&ch[chn].ao.refctrl0, 0x1 << 18);
		clrbits_le32(&ch[chn].ao.mrs, 0x3 << 24);
		setbits_le32(&ch[chn].ao.mpc_option, 0x1 << 17);
		clrsetbits_le32(&ch[chn].phy.b[0].dq[6], 0x3 << 0, 0x1 << 0);
		clrsetbits_le32(&ch[chn].phy.b[1].dq[6], 0x3 << 0, 0x1 << 0);
		clrsetbits_le32(&ch[chn].phy.ca_cmd[6], 0x3 << 0, 0x1 << 0);

		dramc_rx_input_delay_tracking_init_by_freq(chn);

		setbits_le32(&ch[chn].ao.dummy_rd, 0x1 << 25);
		setbits_le32(&ch[chn].ao.drsctrl, 0x1 << 0);
		if (freq_group == LP4X_DDR3200 || freq_group == LP4X_DDR3600)
			clrbits_le32(&ch[chn].ao.shu[1].drving[1], 0x1 << 31);
		else
			setbits_le32(&ch[chn].ao.shu[1].drving[1], 0x1 << 31);
	}

	for (size_t r = 0; r < 2; r++) {
		for (size_t b = 0; b < 2; b++)
			clrbits_le32(&ch[0].phy.r[r].b[b].rxdvs[2],
				(0x1 << 28) | (0x1 << 23) | (0x3 << 30));
		clrbits_le32(&ch[0].phy.r0_ca_rxdvs[2], 0x3 << 30);
	}
}

static void dramc_set_mr13_vrcg_to_Normal(u8 chn, const struct mr_value *mr)
{
	for (u8 rank = 0; rank < RANK_MAX; rank++)
		dramc_mode_reg_write_by_rank(chn, rank, 13,
					     mr->MR13Value & ~(0x1 << 3));

	for (u8 shu = 0; shu < DRAM_DFS_SHUFFLE_MAX; shu++)
		clrbits_le32(&ch[chn].ao.shu[shu].hwset_vrcg, 0x1 << 19);
}

void dramc_apply_config_after_calibration(const struct mr_value *mr)
{
	for (size_t chn = 0; chn < CHANNEL_MAX; chn++) {
		write32(&ch[chn].phy.misc_cg_ctrl4, 0x11400000);
		clrbits_le32(&ch[chn].ao.refctrl1, 0x1 << 7);
		clrbits_le32(&ch[chn].ao.shuctrl, 0x1 << 2);
		clrbits_le32(&ch[chn].phy.ca_cmd[6], 0x1 << 6);
		dramc_set_mr13_vrcg_to_Normal(chn, mr);

		clrbits_le32(&ch[chn].phy.b[0].dq[6], 0x3);
		clrbits_le32(&ch[chn].phy.b[1].dq[6], 0x3);
		clrbits_le32(&ch[chn].phy.ca_cmd[6], 0x3);
		setbits_le32(&ch[chn].phy.b[0].dq[6], 0x1 << 5);
		setbits_le32(&ch[chn].phy.b[1].dq[6], 0x1 << 5);
		setbits_le32(&ch[chn].phy.ca_cmd[6], 0x1 << 5);

		clrbits_le32(&ch[chn].ao.impcal, 0x3 << 24);
		clrbits_le32(&ch[chn].phy.misc_imp_ctrl0, 0x4);
		clrbits_le32(&ch[chn].phy.misc_cg_ctrl0, 0xf);

		clrbits_le32(&ch[chn].phy.misc_ctrl0, 0x1 << 31);
		clrbits_le32(&ch[chn].phy.misc_ctrl1, 0x1 << 25);

		setbits_le32(&ch[chn].ao.spcmdctrl, 1 << 29);
		setbits_le32(&ch[chn].ao.dqsoscr, 1 << 24);

		for (u8 shu = 0; shu < DRAM_DFS_SHUFFLE_MAX; shu++)
			clrbits_le32(&ch[chn].ao.shu[shu].scintv, 0x1 << 30);

		clrbits_le32(&ch[chn].ao.dummy_rd, (0x7 << 20) | (0x1 << 7));
		dramc_cke_fix_onoff(chn, false, false);
		clrbits_le32(&ch[chn].ao.dramc_pd_ctrl, 0x1 << 26);

		clrbits_le32(&ch[chn].ao.eyescan, 0x7 << 8);
		clrsetbits_le32(&ch[chn].ao.test2_4, 0x7 << 28, 0x4 << 28);
	}
}

static void dramc_rx_dqs_isi_pulse_cg_switch(u8 chn, bool flag)
{
	for (size_t b = 0; b < 2; b++)
		clrsetbits_le32(&ch[chn].phy.b[b].dq[6], 1 << 5,
				(flag ? 1 : 0) << 5);
}

static void dramc_set_rank_engine2(u8 chn, u8 rank)
{
	/*  Select CTO_AGENT1_RANK */
	SET32_BITFIELDS(&ch[chn].ao.dramctrl, DRAMCTRL_ADRDECEN_TARKMODE, 1);
	SET32_BITFIELDS(&ch[chn].ao.test2_4, TEST2_4_TESTAGENTRKSEL, 0);
	SET32_BITFIELDS(&ch[chn].ao.test2_4, TEST2_4_TESTAGENTRK, rank);
}

static void dramc_engine2_setpat(u8 chn, bool test_pat)
{
	SET32_BITFIELDS(&ch[chn].ao.test2_4,
			TEST2_4_TEST_REQ_LEN1, 0,
			TEST2_4_TESTXTALKPAT, 0,
			TEST2_4_TESTAUDMODE, 0,
			TEST2_4_TESTAUDBITINV, 0);

	if (!test_pat) {
		SET32_BITFIELDS(&ch[chn].ao.perfctl0, PERFCTL0_RWOFOEN, 1);
		SET32_BITFIELDS(&ch[chn].ao.test2_4,
				TEST2_4_TESTSSOPAT, 0,
				TEST2_4_TESTSSOXTALKPAT, 0,
				TEST2_4_TESTXTALKPAT, 1);
	} else {
		SET32_BITFIELDS(&ch[chn].ao.test2_4,
				TEST2_4_TESTAUDINIT, 0x11,
				TEST2_4_TESTAUDINC, 0xd,
				TEST2_4_TESTAUDBITINV, 1);
	}
	SET32_BITFIELDS(&ch[chn].ao.test2_3,
			TEST2_3_TESTAUDPAT, test_pat, TEST2_3_TESTCNT, 0);
}

static void dramc_engine2_init(u8 chn, u8 rank, u32 t2_1, u32 t2_2, bool test_pat)
{
	dramc_set_rank_engine2(chn, rank);

	SET32_BITFIELDS(&ch[chn].ao.dummy_rd,
			DUMMY_RD_DQSG_DMYRD_EN, 0,
			DUMMY_RD_DQSG_DMYWR_EN, 0,
			DUMMY_RD_DUMMY_RD_EN, 0,
			DUMMY_RD_SREF_DMYRD_EN, 0,
			DUMMY_RD_DMY_RD_DBG, 0,
			DUMMY_RD_DMY_WR_DBG, 0);
	SET32_BITFIELDS(&ch[chn].nao.testchip_dma1,
			TESTCHIP_DMA1_DMA_LP4MATAB_OPT, 0);
	SET32_BITFIELDS(&ch[chn].ao.test2_3,
			TEST2_3_TEST2W, 0,
			TEST2_3_TEST2R, 0,
			TEST2_3_TEST1, 0);
	SET32_BITFIELDS(&ch[chn].ao.test2_0, TEST2_0_PAT0, t2_1 >> 24,
			TEST2_0_PAT1, t2_2 >> 24);
	SET32_BITFIELDS(&ch[chn].ao.test2_1, TEST2_1_TEST2_BASE,
			t2_1 & 0xffffff);
	SET32_BITFIELDS(&ch[chn].ao.test2_2, TEST2_2_TEST2_OFF,
			t2_2 & 0xffffff);

	dramc_engine2_setpat(chn, test_pat);
}

static void dramc_engine2_check_complete(u8 chn, u8 status)
{
	u32 loop = 0;
	/* In some case test engine finished but the complete signal late come,
	 * system will wait very long time. Hence, we set a timeout here.
	 * After system receive complete signal or wait until time out
	 * it will return, the caller will check compare result to verify
	 * whether engine success.
	 */
	while (wait_us(100, read32(&ch[chn].nao.testrpt) & status) != status) {
		if (loop++ > 100)
			dramc_dbg("MEASURE_A timeout\n");
	}
}

static void dramc_engine2_compare(u8 chn, enum dram_te_op wr)
{
	u8 rank_status = ((read32(&ch[chn].ao.test2_3) & 0xf) == 1) ? 3 : 1;

	if (wr == TE_OP_WRITE_READ_CHECK) {
		dramc_engine2_check_complete(chn, rank_status);

		SET32_BITFIELDS(&ch[chn].ao.test2_3, TEST2_3_TEST2W, 0,
				TEST2_3_TEST2R, 0, TEST2_3_TEST1, 0);
		udelay(1);
		SET32_BITFIELDS(&ch[chn].ao.test2_3, TEST2_3_TEST2W, 1);
	}

	dramc_engine2_check_complete(chn, rank_status);
}

static u32 dramc_engine2_run(u8 chn, enum dram_te_op wr)
{
	u32 result;

	if (wr == TE_OP_READ_CHECK) {
		SET32_BITFIELDS(&ch[chn].ao.test2_4, TEST2_4_TESTAUDMODE, 0);
		SET32_BITFIELDS(&ch[chn].ao.test2_3,
				TEST2_3_TEST2W, 0, TEST2_3_TEST2R, 1, TEST2_3_TEST1, 0);
	} else if (wr == TE_OP_WRITE_READ_CHECK) {
		SET32_BITFIELDS(&ch[chn].ao.test2_3,
				TEST2_3_TEST2W, 1, TEST2_3_TEST2R, 0, TEST2_3_TEST1, 0);
	}

	dramc_engine2_compare(chn, wr);

	udelay(1);
	result = read32(&ch[chn].nao.cmp_err);
	SET32_BITFIELDS(&ch[chn].ao.test2_3,
			TEST2_3_TEST2W, 0, TEST2_3_TEST2R, 0, TEST2_3_TEST1, 0);

	return result;
}

static void dramc_engine2_end(u8 chn, u32 dummy_rd)
{
	clrbits_le32(&ch[chn].ao.test2_4, 0x1 << 17);
	write32(&ch[chn].ao.dummy_rd, dummy_rd);
}

static bool dramc_find_gating_window(u32 result_r, u32 result_f, u32 *debug_cnt,
		u8 dly_coarse_large, u8 dly_coarse_0p5t, u8 *pass_begin, u8 *pass_count,
		u8 *pass_count_1, u8 *dly_fine_xt, u8 *dqs_high, u8 *dqs_done)
{
	bool find_tune = false;
	u16 debug_cnt_perbyte, current_pass = 0, pass_byte_cnt = 0;

	for (u8 dqs = 0; dqs < DQS_NUMBER; dqs++) {
		u8 dqs_result_r = (u8) ((result_r >> (8 * dqs)) & 0xff);
		u8 dqs_result_f = (u8) ((result_f >> (8 * dqs)) & 0xff);

		if (pass_byte_cnt & (1 << dqs))
			continue;
		current_pass = 0;

		debug_cnt_perbyte = (u16) debug_cnt[dqs];
		if (dqs_result_r == 0 && dqs_result_f == 0 &&
			debug_cnt_perbyte == GATING_GOLDEND_DQSCNT)
			current_pass = 1;

		if (current_pass) {
			if (pass_begin[dqs] == 0) {
				pass_begin[dqs] = 1;
				pass_count_1[dqs] = 0;
				dramc_dbg("[Byte %d]First pass (%d, %d, %d)\n",
					dqs, dly_coarse_large, dly_coarse_0p5t, *dly_fine_xt);
			}

			if (pass_begin[dqs] == 1)
				pass_count_1[dqs]++;

			if (pass_begin[dqs] == 1 &&
			    pass_count_1[dqs] * DQS_GW_FINE_STEP > DQS_GW_FINE_END) {
				dqs_high[dqs] = 0;
				dqs_done[dqs] = 1;
			}

			if (pass_count_1[0] * DQS_GW_FINE_STEP > DQS_GW_FINE_END &&
			    pass_count_1[1] * DQS_GW_FINE_STEP > DQS_GW_FINE_END) {
				dramc_dbg("All bytes gating window > 1 coarse_tune, Early break\n");
				*dly_fine_xt = DQS_GW_FINE_END;
				find_tune = true;
			}
		} else {
			if (pass_begin[dqs] != 1)
				continue;

			dramc_dbg("[Byte %d] pass_begin[dqs]:%d, pass_count[dqs]:%d,pass_count_1:%d\n",
				dqs, pass_begin[dqs], pass_count[dqs], pass_count_1[dqs]);

			pass_begin[dqs] = 0;
			if (pass_count_1[dqs] > pass_count[dqs]) {
				pass_count[dqs] = pass_count_1[dqs];
				if (pass_count_1[dqs] * DQS_GW_FINE_STEP > 32 &&
				    pass_count_1[dqs] * DQS_GW_FINE_STEP < 96)
					pass_byte_cnt |= (1 << dqs);
				if (pass_byte_cnt == 3) {
					*dly_fine_xt = DQS_GW_FINE_END;
					find_tune = true;
				}
			}
		}
	}

	return find_tune;
}

static void dramc_find_dly_tune(u8 chn, u8 dly_coarse_large, u8 dly_coarse_0p5t,
		u8 dly_fine_xt, u8 *dqs_high, u8 *dly_coarse_large_cnt,
		u8 *dly_coarse_0p5t_cnt, u8 *dly_fine_tune_cnt, u8 *dqs_trans, u8 *dqs_done)
{
	for (size_t dqs = 0; dqs < DQS_NUMBER; dqs++) {
		u32 dqs_cnt = read32(&ch[chn].phy_nao.misc_phy_stben_b[dqs]);
		dqs_cnt = (dqs_cnt >> 16) & 3;

		if (dqs_done[dqs] == 1)
			continue;

		if (dqs_cnt == 3)
			dqs_high[dqs]++;

		if (dqs_high[dqs] * DQS_GW_FINE_STEP <= 16)
			continue;

		switch (dqs_cnt) {
		case 3:
			dly_coarse_large_cnt[dqs] = dly_coarse_large;
			dly_coarse_0p5t_cnt[dqs] = dly_coarse_0p5t;
			dly_fine_tune_cnt[dqs] = dly_fine_xt;
			dqs_trans[dqs] = 1;
			break;
		case 2:
		case 1:
			if (dqs_trans[dqs] == 1)
				dramc_dbg("[Byte %ld] Lead/lag falling Transition"
					  " (%d, %d, %d)\n",
					  dqs, dly_coarse_large_cnt[dqs],
					  dly_coarse_0p5t_cnt[dqs], dly_fine_tune_cnt[dqs]);
			dqs_trans[dqs]++;
			break;
		case 0:
			dramc_dbg("[Byte %ld] Lead/lag Transition tap number (%d)\n",
				dqs, dqs_trans[dqs]);
			dqs_high[dqs] = 0;
			break;
		}
	}
}

static void dram_phy_reset(u8 chn)
{
	SET32_BITFIELDS(&ch[chn].ao.ddrconf0, DDRCONF0_RDATRST, 1);
	SET32_BITFIELDS(&ch[chn].phy.misc_ctrl1, MISC_CTRL1_R_DMPHYRST, 1);
	clrbits_le32(&ch[chn].phy.b[0].dq[9], (1 << 4) | (1 << 0));
	clrbits_le32(&ch[chn].phy.b[1].dq[9], (1 << 4) | (1 << 0));

	udelay(1);
	setbits_le32(&ch[chn].phy.b[1].dq[9], (1 << 4) | (1 << 0));
	setbits_le32(&ch[chn].phy.b[0].dq[9], (1 << 4) | (1 << 0));
	SET32_BITFIELDS(&ch[chn].phy.misc_ctrl1, MISC_CTRL1_R_DMPHYRST, 0);
	SET32_BITFIELDS(&ch[chn].ao.ddrconf0, DDRCONF0_RDATRST, 0);
}

static void dramc_set_gating_mode(u8 chn, bool mode)
{
	u8 vref = 0, burst = 0;

	if (mode) {
		vref = 2;
		burst = 1;
	}

	for (size_t b = 0; b < 2; b++) {
		clrsetbits_le32(&ch[chn].phy.b[b].dq[6], 0x3 << 14, vref << 14);
		setbits_le32(&ch[chn].phy.b[b].dq[9], 0x1 << 5);
	}

	clrsetbits_le32(&ch[chn].ao.stbcal1, 0x1 << 5, burst << 5);
	setbits_le32(&ch[chn].ao.stbcal, 0x1 << 30);

	clrbits_le32(&ch[chn].phy.b[0].dq[9], (0x1 << 4) | (0x1 << 0));
	clrbits_le32(&ch[chn].phy.b[1].dq[9], (0x1 << 4) | (0x1 << 0));
	udelay(1);
	setbits_le32(&ch[chn].phy.b[1].dq[9], (0x1 << 4) | (0x1 << 0));
	setbits_le32(&ch[chn].phy.b[0].dq[9], (0x1 << 4) | (0x1 << 0));
}

static void dramc_rx_dqs_gating_cal_pre(u8 chn, u8 rank)
{
	SET32_BITFIELDS(&ch[chn].ao.refctrl0, REFCTRL0_PBREFEN, 0);

	dramc_hw_gating_onoff(chn, false);

	SET32_BITFIELDS(&ch[chn].ao.stbcal1, STBCAL1_STBENCMPEN, 1);
	SET32_BITFIELDS(&ch[chn].ao.stbcal1, STBCAL1_STBCNT_LATCH_EN, 1);
	SET32_BITFIELDS(&ch[chn].ao.ddrconf0, DDRCONF0_DM4TO1MODE, 0);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_DQSGCNTEN, 1);

	udelay(4);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_DQSGCNTRST, 1);
	udelay(1);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_DQSGCNTRST, 0);
	SET32_BITFIELDS(&ch[chn].phy.misc_ctrl1, MISC_CTRL1_R_DMSTBENCMP_RK,
			rank);

}

static void set_selph_gating_value(uint32_t *addr, u8 dly, u8 dly_p1)
{
	clrsetbits_le32(addr, 0x77777777,
		(dly << 0) | (dly << 8) | (dly << 16) | (dly << 24) |
		(dly_p1 << 4) | (dly_p1 << 12) | (dly_p1 << 20) | (dly_p1 << 28));
}

static void dramc_write_dqs_gating_result(u8 chn, u8 rank,
		u8 *best_coarse_tune2t, u8 *best_coarse_tune0p5t,
		u8 *best_coarse_tune2t_p1, u8 *best_coarse_tune0p5t_p1)
{
	u8 best_coarse_rodt[DQS_NUMBER], best_coarse_0p5t_rodt[DQS_NUMBER];
	u8 best_coarse_rodt_p1[DQS_NUMBER];
	u8 best_coarse_0p5t_rodt_p1[DQS_NUMBER];

	dramc_rx_dqs_isi_pulse_cg_switch(chn, true);

	clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dqsg0,
		0x77777777,
		(best_coarse_tune2t[0] << 0) | (best_coarse_tune2t[1] << 8) |
		(best_coarse_tune2t_p1[0] << 4) | (best_coarse_tune2t_p1[1] << 12));
	clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dqsg1,
		0x77777777,
		(best_coarse_tune0p5t[0] << 0) | (best_coarse_tune0p5t[1] << 8) |
		(best_coarse_tune0p5t_p1[0] << 4) | (best_coarse_tune0p5t_p1[1] << 12));

	for (size_t dqs = 0; dqs < DQS_NUMBER; dqs++) {
		u8 tmp_value = (best_coarse_tune2t[dqs] << 3)
			+ best_coarse_tune0p5t[dqs];

		if (tmp_value >= 11) {
			tmp_value -= 11;
			best_coarse_rodt[dqs] = tmp_value >> 3;
			best_coarse_0p5t_rodt[dqs] =
				tmp_value - (best_coarse_rodt[dqs] << 3);

			tmp_value = (best_coarse_tune2t_p1[dqs] << 3) +
				best_coarse_tune0p5t_p1[dqs] - 11;
			best_coarse_rodt_p1[dqs] = tmp_value >> 3;
			best_coarse_0p5t_rodt_p1[dqs] =
				tmp_value - (best_coarse_rodt_p1[dqs] << 3);

			dramc_dbg("Best RODT dly(2T, 0.5T) = (%d, %d)\n",
				  best_coarse_rodt[dqs],
				  best_coarse_0p5t_rodt[dqs]);
		} else {
			best_coarse_rodt[dqs] = 0;
			best_coarse_0p5t_rodt[dqs] = 0;
			best_coarse_rodt_p1[dqs] = 4;
			best_coarse_0p5t_rodt_p1[dqs] = 4;
			dramc_dbg("RxdqsGatingCal error: best_coarse_tune2t:%zd"
				  " is already 0. RODT cannot be -1 coarse\n",
				  dqs);
		}
	}

	clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_odten0,
		0x77777777,
		(best_coarse_rodt[0] << 0) | (best_coarse_rodt[1] << 8) |
		(best_coarse_rodt_p1[0] << 4) | (best_coarse_rodt_p1[1] << 12));
	clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_odten1,
		0x77777777,
		(best_coarse_0p5t_rodt[0] << 0) | (best_coarse_0p5t_rodt[1] << 8) |
		(best_coarse_0p5t_rodt_p1[0] << 4) | (best_coarse_0p5t_rodt_p1[1] << 12));
}

static void dramc_rx_dqs_gating_cal_partial(u8 chn, u8 rank,
	u32 coarse_start, u32 coarse_end, u8 freqDiv,
	u8 *pass_begin, u8 *pass_count, u8 *pass_count_1, u8 *dqs_done,
	u8 *dqs_high, u8 *dqs_transition, u8 *dly_coarse_large_cnt,
	u8 *dly_coarse_0p5t_cnt, u8 *dly_fine_tune_cnt)
{
	u8 dqs;
	u32 debug_cnt[DQS_NUMBER];

	for (u32 coarse_tune = coarse_start; coarse_tune < coarse_end;
	     coarse_tune++) {
		u32 dly_coarse_large_rodt = 0, dly_coarse_0p5t_rodt = 0;
		u32 dly_coarse_large_rodt_p1 = 4, dly_coarse_0p5t_rodt_p1 = 4;

		u8 dly_coarse_large = coarse_tune / RX_DQS_CTL_LOOP;
		u8 dly_coarse_0p5t = coarse_tune % RX_DQS_CTL_LOOP;
		u32 dly_coarse_large_p1 = (coarse_tune + freqDiv) / RX_DQS_CTL_LOOP;
		u32 dly_coarse_0p5t_p1 = (coarse_tune + freqDiv) % RX_DQS_CTL_LOOP;
		u32 value = (dly_coarse_large << 3) + dly_coarse_0p5t;

		if (value >= 11) {
			value -= 11;
			dly_coarse_large_rodt = value >> 3;
			dly_coarse_0p5t_rodt =
				value - (dly_coarse_large_rodt << 3);
			value = (dly_coarse_large << 3) + dly_coarse_0p5t - 11;
			dly_coarse_large_rodt_p1 = value >> 3;
			dly_coarse_0p5t_rodt_p1 =
				value - (dly_coarse_large_rodt_p1 << 3);
		}

		set_selph_gating_value(&ch[chn].ao.shu[0].rk[rank].selph_dqsg0,
			dly_coarse_large, dly_coarse_large_p1);
		set_selph_gating_value(&ch[chn].ao.shu[0].rk[rank].selph_dqsg1,
			dly_coarse_0p5t, dly_coarse_0p5t_p1);
		set_selph_gating_value(&ch[chn].ao.shu[0].rk[rank].selph_odten0,
			dly_coarse_large_rodt, dly_coarse_large_rodt_p1);
		set_selph_gating_value(&ch[chn].ao.shu[0].rk[rank].selph_odten1,
			dly_coarse_0p5t_rodt, dly_coarse_0p5t_rodt_p1);

		for (u8 dly_fine_xt = 0; dly_fine_xt < DQS_GW_FINE_END;
		     dly_fine_xt += 4) {
			dramc_set_gating_mode(chn, 0);
			WRITE32_BITFIELDS(&ch[chn].ao.shu[0].rk[rank].dqsien,
					  SHURK_DQSIEN_DQS0IEN, dly_fine_xt,
					  SHURK_DQSIEN_DQS1IEN, dly_fine_xt);

			dram_phy_reset(chn);
			SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_DQSGCNTRST, 1);
			udelay(1);
			SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_DQSGCNTRST, 0);

			dramc_engine2_run(chn, TE_OP_READ_CHECK);

			u32 result_r = READ32_BITFIELD(
					&ch[chn].phy.misc_stberr_rk0_r,
					MISC_STBERR_RK_R_STBERR_RK_R);
			u32 result_f = READ32_BITFIELD(
					&ch[chn].phy.misc_stberr_rk0_f,
					MISC_STBERR_RK_F_STBERR_RK_F);

			debug_cnt[0] = read32(&ch[chn].nao.dqsgnwcnt[0]);
			debug_cnt[1] = (debug_cnt[0] >> 16) & 0xffff;
			debug_cnt[0] &= 0xffff;

			dramc_set_gating_mode(chn, 1);
			dramc_engine2_run(chn, TE_OP_READ_CHECK);

			dramc_find_dly_tune(chn, dly_coarse_large,
				dly_coarse_0p5t, dly_fine_xt, dqs_high,
				dly_coarse_large_cnt, dly_coarse_0p5t_cnt,
				dly_fine_tune_cnt, dqs_transition, dqs_done);

			dramc_dbg("%d %d %d |", dly_coarse_large,
				  dly_coarse_0p5t, dly_fine_xt);
			for (dqs = 0; dqs < DQS_NUMBER; dqs++)
				dramc_dbg("%X ", debug_cnt[dqs]);

			dramc_dbg(" |");
			for (dqs = 0; dqs < DQS_NUMBER; dqs++) {
				dramc_dbg("(%X %X)",
				(result_f >> (DQS_BIT_NUMBER * dqs)) & 0xff,
				(result_r >> (DQS_BIT_NUMBER * dqs)) & 0xff);
			}

			dramc_dbg("\n");
			if (dramc_find_gating_window(result_r, result_f, debug_cnt,
			    dly_coarse_large, dly_coarse_0p5t, pass_begin,
			    pass_count, pass_count_1, &dly_fine_xt,
			    dqs_high, dqs_done))
				coarse_tune = coarse_end;
		}
	}
}

static void dramc_rx_dqs_gating_cal(u8 chn, u8 rank, u8 freq_group,
		const struct sdram_params *params, const bool fast_calib,
		const struct mr_value *mr)
{
	u8 dqs, fsp, freqDiv = 4;
	u8 pass_begin[DQS_NUMBER] = {0}, pass_count[DQS_NUMBER] = {0},
		pass_count_1[DQS_NUMBER] = {0}, dqs_done[DQS_NUMBER] = {0};
	u8 min_coarse_tune2t[DQS_NUMBER], min_coarse_tune0p5t[DQS_NUMBER],
		min_fine_tune[DQS_NUMBER];
	u8 best_fine_tune[DQS_NUMBER], best_coarse_tune0p5t[DQS_NUMBER],
		best_coarse_tune2t[DQS_NUMBER];
	u8 best_coarse_tune0p5t_p1[DQS_NUMBER], best_coarse_tune2t_p1[DQS_NUMBER];
	u8 dqs_high[DQS_NUMBER] = {0}, dqs_transition[DQS_NUMBER] = {0};
	u8 dly_coarse_large_cnt[DQS_NUMBER] = {0}, dly_coarse_0p5t_cnt[DQS_NUMBER] = {0},
		dly_fine_tune_cnt[DQS_NUMBER] = {0};
	u32 coarse_start, coarse_end;

	struct reg_value regs_bak[] = {
		{&ch[chn].ao.stbcal},
		{&ch[chn].ao.stbcal1},
		{&ch[chn].ao.ddrconf0},
		{&ch[chn].ao.spcmd},
		{&ch[chn].ao.refctrl0},
		{&ch[chn].phy.b[0].dq[6]},
		{&ch[chn].phy.b[1].dq[6]},
	};
	for (size_t i = 0; i < ARRAY_SIZE(regs_bak); i++)
		regs_bak[i].value = read32(regs_bak[i].addr);

	fsp = get_freq_fsq(freq_group);
	dramc_rx_dqs_isi_pulse_cg_switch(chn, false);

	dramc_mode_reg_write_by_rank(chn, rank, 0x1, mr->MR01Value[fsp] | 0x80);
	dramc_rx_dqs_gating_cal_pre(chn, rank);

	u32 dummy_rd_backup = read32(&ch[chn].ao.dummy_rd);
	dramc_engine2_init(chn, rank, TEST2_1_CAL, 0xaa000023, true);

	switch (freq_group) {
	case LP4X_DDR1600:
		coarse_start = 18;
		break;
	case LP4X_DDR2400:
		coarse_start = 25;
		break;
	case LP4X_DDR3200:
		coarse_start = 25;
		break;
	case LP4X_DDR3600:
		coarse_start = 21;
		break;
	default:
		die("Invalid DDR frequency group %u\n", freq_group);
		return;
	}
	coarse_end = coarse_start + 12;

	dramc_dbg("[Gating]\n");

	if (!fast_calib) {
		dramc_rx_dqs_gating_cal_partial(chn, rank,
			coarse_start, coarse_end,
			freqDiv, pass_begin, pass_count, pass_count_1, dqs_done,
			dqs_high, dqs_transition, dly_coarse_large_cnt,
			dly_coarse_0p5t_cnt, dly_fine_tune_cnt);
		dramc_engine2_end(chn, dummy_rd_backup);
	}

	for (dqs = 0; dqs < DQS_NUMBER; dqs++) {
		if  (fast_calib) {
			dramc_dbg("[bypass Gating params] dqs: %d\n", dqs);
			pass_count[dqs] = params->gating_pass_count[chn][rank][dqs];
			min_fine_tune[dqs] = params->gating_fine_tune[chn][rank][dqs];
			min_coarse_tune0p5t[dqs] = params->gating05T[chn][rank][dqs];
			min_coarse_tune2t[dqs] = params->gating2T[chn][rank][dqs];
		} else {
			pass_count[dqs] = dqs_transition[dqs];
			min_fine_tune[dqs] = dly_fine_tune_cnt[dqs];
			min_coarse_tune0p5t[dqs] = dly_coarse_0p5t_cnt[dqs];
			min_coarse_tune2t[dqs] = dly_coarse_large_cnt[dqs];
		}
		u8 tmp_offset = pass_count[dqs] * DQS_GW_FINE_STEP / 2;
		u8 tmp_value = min_fine_tune[dqs] + tmp_offset;
		best_fine_tune[dqs] = tmp_value % RX_DLY_DQSIENSTB_LOOP;

		tmp_offset = tmp_value / RX_DLY_DQSIENSTB_LOOP;
		tmp_value = min_coarse_tune0p5t[dqs] + tmp_offset;
		best_coarse_tune0p5t[dqs] = tmp_value % RX_DQS_CTL_LOOP;

		tmp_offset = tmp_value / RX_DQS_CTL_LOOP;
		best_coarse_tune2t[dqs] = min_coarse_tune2t[dqs] + tmp_offset;

		tmp_value = best_coarse_tune0p5t[dqs] + freqDiv;
		best_coarse_tune0p5t_p1[dqs] = tmp_value % RX_DQS_CTL_LOOP;

		tmp_offset = tmp_value / RX_DQS_CTL_LOOP;
		best_coarse_tune2t_p1[dqs] =
			best_coarse_tune2t[dqs] + tmp_offset;
	}

	for (dqs = 0; dqs < DQS_NUMBER; dqs++)
		dramc_dbg("Best DQS%d dly(2T, 0.5T, fine tune)"
			  " = (%d, %d, %d)\n", dqs, best_coarse_tune2t[dqs],
			  best_coarse_tune0p5t[dqs], best_fine_tune[dqs]);

	for (dqs = 0; dqs < DQS_NUMBER; dqs++)
		dramc_dbg("Best DQS%d P1 dly(2T, 0.5T, fine tune)"
			  " = (%d, %d, %d)\n", dqs, best_coarse_tune2t_p1[dqs],
			  best_coarse_tune0p5t_p1[dqs], best_fine_tune[dqs]);

	for (size_t i = 0; i < ARRAY_SIZE(regs_bak); i++)
		write32(regs_bak[i].addr, regs_bak[i].value);

	dramc_mode_reg_write_by_rank(chn, rank, 0x1, mr->MR01Value[fsp]);

	dramc_write_dqs_gating_result(chn, rank, best_coarse_tune2t,
		best_coarse_tune0p5t, best_coarse_tune2t_p1, best_coarse_tune0p5t_p1);

	WRITE32_BITFIELDS(&ch[chn].ao.shu[0].rk[rank].dqsien,
			  SHURK_DQSIEN_DQS0IEN, best_fine_tune[0],
			  SHURK_DQSIEN_DQS1IEN, best_fine_tune[1]);

	dram_phy_reset(chn);
}

static void dramc_rx_rd_dqc_init(u8 chn, u8 rank)
{
	const u8 *lpddr_phy_mapping = phy_mapping[chn];
	u16 temp_value = 0;

	for (size_t b = 0; b < 2; b++)
		clrbits_le32(&ch[chn].phy.shu[0].b[b].dq[7], 0x1 << 7);

	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSRK, rank);
	SET32_BITFIELDS(&ch[chn].ao.mpc_option, MPC_OPTION_MPCRKEN, 1);

	for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++)
		temp_value |= ((0x5555 >> bit) & 0x1) << lpddr_phy_mapping[bit];

	u16 mr15_golden_value = temp_value & 0xff;
	u16 mr20_golden_value = (temp_value >> 8) & 0xff;
	SET32_BITFIELDS(&ch[chn].ao.mr_golden,
			MR_GOLDEN_MR15_GOLDEN, mr15_golden_value,
			MR_GOLDEN_MR20_GOLDEN, mr20_golden_value);
}

static u32 dramc_rx_rd_dqc_run(u8 chn)
{
	u32 loop = 0;
	SET32_BITFIELDS(&ch[chn].ao.spcmdctrl, SPCMDCTRL_RDDQCDIS, 1);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_RDDQCEN, 1);

	while (!wait_us(10, read32(&ch[chn].nao.spcmdresp) & (0x1 << 7))) {
		if (loop++ > 10)
			dramc_dbg("[RDDQC] resp fail (time out)\n");
	}

	u32 result = read32(&ch[chn].nao.rdqc_cmp);
	SET32_BITFIELDS(&ch[chn].ao.spcmd, SPCMD_RDDQCEN, 0);
	SET32_BITFIELDS(&ch[chn].ao.spcmdctrl, SPCMDCTRL_RDDQCDIS, 0);

	return result;
}

static void dramc_rx_rd_dqc_end(u8 chn)
{
	SET32_BITFIELDS(&ch[chn].ao.mrs, MRS_MRSRK, 0);
}

static void dramc_rx_vref_pre_setting(u8 chn)
{
	SET32_BITFIELDS(&ch[chn].phy.b[0].dq[5], B0_DQ5_RG_RX_ARDQ_VREF_EN_B0, 1);
	SET32_BITFIELDS(&ch[chn].phy.b[1].dq[5], B1_DQ5_RG_RX_ARDQ_VREF_EN_B1, 1);
}

static void dramc_set_rx_vref(u8 chn, u8 vref)
{
	for (size_t b = 0; b < 2; b++)
		SET32_BITFIELDS(&ch[chn].phy.shu[0].b[b].dq[5],
				SHU1_BX_DQ5_RG_RX_ARDQ_VREF_SEL_B0, vref);
	dramc_dbg("set rx vref :%d\n", vref);
}

static void dramc_set_tx_vref(u8 chn, u8 rank, u8 value)
{
	dramc_mode_reg_write_by_rank(chn, rank, 14, value);
}

static void dramc_set_vref(u8 chn, u8 rank, enum CAL_TYPE type, u8 vref)
{
	if (type == RX_WIN_TEST_ENG)
		dramc_set_rx_vref(chn, vref);
	else
		dramc_set_tx_vref(chn, rank, vref);
}

static void dramc_transfer_dly_tune(u8 chn, u32 dly, u32 adjust_center,
		struct tx_dly_tune *dly_tune)
{
	u8 tune = 3, fine_tune = 0;
	u16 tmp;

	fine_tune = dly & (TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP - 1);
	tmp = (dly / TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP) << 1;

	if (adjust_center) {
		if (fine_tune < 10) {
			fine_tune += TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP >> 1;
			tmp--;
		} else if (fine_tune > TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP - 10) {
			fine_tune -= TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP >> 1;
			tmp++;
		}
	}

	dly_tune->fine_tune = fine_tune;
	dly_tune->coarse_tune_small = tmp - ((tmp >> tune) << tune);
	dly_tune->coarse_tune_large = tmp >> tune;

	tmp -= 3;
	dly_tune->coarse_tune_small_oen = tmp - ((tmp >> tune) << tune);
	dly_tune->coarse_tune_large_oen = tmp >> tune;
}

static void dramc_set_rx_dly_factor(u8 chn, u8 rank, enum RX_TYPE type, u32 val)
{
	int b, dq;

	switch (type) {
	case RX_DQ:
		for (dq = 2; dq < 6; dq++)
			for (b = 0; b < 2; b++)
				WRITE32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[b].dq[dq],
						  SHU1_R0_B0_DQ6_RK0_RX_ARDQM0_F_DLY_B0, val,
						  SHU1_R0_B0_DQ6_RK0_RX_ARDQM0_R_DLY_B0, val,
						  SHU1_R0_B0_DQ6_RK0_RX_ARDQS0_F_DLY_B0, val,
						  SHU1_R0_B0_DQ6_RK0_RX_ARDQS0_R_DLY_B0, val);
		break;

	case RX_DQM:
		for (b = 0; b < 2; b++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[b].dq[6],
					SHU1_R0_B0_DQ6_RK0_RX_ARDQM0_F_DLY_B0, val,
					SHU1_R0_B0_DQ6_RK0_RX_ARDQM0_R_DLY_B0, val);
		break;

	case RX_DQS:
		for (b = 0; b < 2; b++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[b].dq[6],
					SHU1_R0_B0_DQ6_RK0_RX_ARDQS0_F_DLY_B0, val,
					SHU1_R0_B0_DQ6_RK0_RX_ARDQS0_R_DLY_B0, val);
		break;
	default:
		dramc_err("error calibration type: %d\n", type);
		break;
	}
}

static void dramc_set_tx_dly_factor(u8 chn, u8 rk,
		enum CAL_TYPE type, u8 *dq_small_reg, u32 dly)
{
	struct tx_dly_tune dly_tune = {0};
	u32 dly_large = 0, dly_large_oen = 0, dly_small = 0, dly_small_oen = 0;
	u32 adjust_center = 0;

	dramc_transfer_dly_tune(chn, dly, adjust_center, &dly_tune);

	for (u8 i = 0; i < 4; i++) {
		dly_large += dly_tune.coarse_tune_large << (i * 4);
		dly_large_oen += dly_tune.coarse_tune_large_oen << (i * 4);
		dly_small += dly_tune.coarse_tune_small << (i * 4);
		dly_small_oen += dly_tune.coarse_tune_small_oen << (i * 4);
	}

	if (type == TX_WIN_DQ_DQM)
		dramc_dbg("%3d |%d  %d  %2d | [0]",
			dly, dly_tune.coarse_tune_large,
			dly_tune.coarse_tune_small, dly_tune.fine_tune);

	if (*dq_small_reg != dly_tune.coarse_tune_small) {
		if (type == TX_WIN_DQ_DQM || type == TX_WIN_DQ_ONLY) {
			clrsetbits_le32(&ch[chn].ao.shu[0].rk[rk].selph_dq[0],
				0x77777777, dly_large | (dly_large_oen << 16));
			clrsetbits_le32(&ch[chn].ao.shu[0].rk[rk].selph_dq[2],
				0x77777777, dly_small | (dly_small_oen << 16));
		}

		if (type == TX_WIN_DQ_DQM) {
			/* Large coarse_tune setting */
			clrsetbits_le32(&ch[chn].ao.shu[0].rk[rk].selph_dq[1],
				0x77777777, dly_large | (dly_large_oen << 16));
			clrsetbits_le32(&ch[chn].ao.shu[0].rk[rk].selph_dq[3],
				0x77777777, dly_small | (dly_small_oen << 16));
		}
	}
	*dq_small_reg = dly_tune.coarse_tune_small;

	if (type == TX_WIN_DQ_DQM || type == TX_WIN_DQ_ONLY) {
		for (size_t b = 0; b < 2; b++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rk].b[b].dq[7],
					FINE_TUNE_DQ, dly_tune.fine_tune);
	}
	if (type == TX_WIN_DQ_DQM) {
		for (size_t b = 0; b < 2; b++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rk].b[b].dq[7],
					FINE_TUNE_DQM, dly_tune.fine_tune);
	}
}

static u32 dramc_get_smallest_dqs_dly(u8 chn, u8 rank, const struct sdram_params *params)
{
	const u8 mck = 3;
	u32 min_dly = 0xffff, virtual_delay = 0;
	u32 tx_dly = read32(&ch[chn].ao.shu[0].selph_dqs0);
	u32 dly = read32(&ch[chn].ao.shu[0].selph_dqs1);
	u32 tmp;

	for (size_t dqs = 0; dqs < DQS_NUMBER; dqs++) {
		tmp = ((tx_dly >> (dqs << 2) & 0x7) << mck) +
		      (dly >> (dqs << 2) & 0x7);
		virtual_delay = (tmp << 5) + params->wr_level[chn][rank][dqs];
		min_dly = MIN(min_dly, virtual_delay);
	}

	return min_dly;
}

static void dramc_get_dly_range(u8 chn, u8 rank, enum CAL_TYPE type,
		u8 freq_group, u16 *pre_cal, s16 *begin, s16 *end,
		const struct sdram_params *params)
{
	u16 pre_dq_dly;
	switch (type) {
	case RX_WIN_RD_DQC:
	case RX_WIN_TEST_ENG:
		switch (freq_group) {
		case LP4X_DDR1600:
			*begin = -48;
			break;
		case LP4X_DDR2400:
			*begin = -30;
			break;
		case LP4X_DDR3200:
		case LP4X_DDR3600:
			*begin = -26;
			break;
		default:
			die("Invalid DDR frequency group %u\n", freq_group);
			return;
		}

		*end = MAX_RX_DQDLY_TAPS;
		break;

	case TX_WIN_DQ_DQM:
		*begin = dramc_get_smallest_dqs_dly(chn, rank, params);
		*end = *begin + 256;
		break;

	case TX_WIN_DQ_ONLY:
		pre_dq_dly = MIN(pre_cal[0], pre_cal[1]);
		pre_dq_dly = (pre_dq_dly > 24) ? (pre_dq_dly - 24) : 0;
		*begin = pre_dq_dly;
		*end = *begin + 64;
		break;
	default:
		dramc_err("error calibration type: %d\n", type);
		break;
	}
}

static int dramc_check_dqdqs_win(struct win_perbit_dly *perbit_dly,
	s16 dly, s16 dly_end, bool fail_bit)
{
	int pass_win = 0;

	if (perbit_dly->first_pass == PASS_RANGE_NA) {
		if (!fail_bit) /* compare correct: pass */
			perbit_dly->first_pass = dly;
	} else if (perbit_dly->last_pass == PASS_RANGE_NA) {
		if (fail_bit) /* compare error: fail */
			perbit_dly->last_pass = dly - 1;
		else if (dly == dly_end)
			perbit_dly->last_pass = dly;

		if (perbit_dly->last_pass != PASS_RANGE_NA) {
			pass_win = perbit_dly->last_pass - perbit_dly->first_pass;
			int best_pass_win = perbit_dly->best_last - perbit_dly->best_first;
			if (pass_win >= best_pass_win) {
				perbit_dly->best_last = perbit_dly->last_pass;
				perbit_dly->best_first = perbit_dly->first_pass;
			}

			/* Clear to find the next pass range if it has */
			perbit_dly->first_pass = PASS_RANGE_NA;
			perbit_dly->last_pass = PASS_RANGE_NA;
		}
	}

	return pass_win;
}

static void dramc_set_vref_dly(struct vref_perbit_dly *vref_dly, struct win_perbit_dly delay[])
{
	struct win_perbit_dly *perbit_dly = vref_dly->perbit_dly;

	for (u8 bit = 0; bit < DQ_DATA_WIDTH; bit++) {
		delay[bit].win_center = (delay[bit].best_first + delay[bit].best_last) >> 1;

		perbit_dly[bit].best_first = delay[bit].best_first;
		perbit_dly[bit].best_last = delay[bit].best_last;
		perbit_dly[bit].win_center = delay[bit].win_center;
		perbit_dly[bit].best_dqdly = delay[bit].best_dqdly;
	}
}

static bool dramk_calc_best_vref(enum CAL_TYPE type, u8 vref,
		struct vref_perbit_dly *vref_dly, struct win_perbit_dly delay[],
		u32 *win_min_max)
{
	u32 win_size, min_bit = 0xff, min_winsize = 0xffff, tmp_win_sum = 0;

	switch (type) {
	case RX_WIN_RD_DQC:
	case RX_WIN_TEST_ENG:
		for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++) {
			win_size = delay[bit].best_last - delay[bit].best_first;

			if (win_size < min_winsize) {
				min_bit = bit;
				min_winsize = win_size;
			}
			tmp_win_sum += win_size;
		}
		dramc_dbg("type:%d vref:%d Min Bit=%d, min_winsize=%d, win sum:%d\n",
			type, vref, min_bit, min_winsize, tmp_win_sum);

		if (tmp_win_sum > vref_dly->max_win_sum) {
			*win_min_max = min_winsize;
			vref_dly->max_win_sum = tmp_win_sum;

			/* best vref */
			vref_dly->best_vref = vref;
		}
		dramc_dbg("type:%d vref:%d, win_sum_total:%d, tmp_win_sum:%d)\n",
			type, vref, vref_dly->max_win_sum, tmp_win_sum);
		dramc_set_vref_dly(vref_dly, delay);

		if (tmp_win_sum < vref_dly->max_win_sum * 95 / 100) {
			dramc_dbg("type:%d best vref found[%d], early break! (%d < %d)\n",
				type, vref_dly->best_vref, tmp_win_sum,
				vref_dly->max_win_sum * 95 / 100);
			return true;
		}

		break;
	case TX_WIN_DQ_ONLY:
	case TX_WIN_DQ_DQM:
		for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++) {
			win_size = delay[bit].best_last - delay[bit].best_first;

			if (win_size < min_winsize) {
				min_bit = bit;
				min_winsize = win_size;
			}
			tmp_win_sum += win_size;
		}
		dramc_dbg("type:%d vref:%d Min Bit=%d, min_winsize=%d, win sum:%d\n",
			type, vref, min_bit, min_winsize, tmp_win_sum);

		if (min_winsize > *win_min_max ||
		    (min_winsize == *win_min_max &&
		     tmp_win_sum > vref_dly->max_win_sum)) {
			*win_min_max = min_winsize;
			vref_dly->max_win_sum = tmp_win_sum;

			/* best vref */
			vref_dly->best_vref = vref;
		}
		dramc_dbg("type:%d vref:%d, win_sum_total:%d, tmp_win_sum:%d)\n",
			type, vref, vref_dly->max_win_sum, tmp_win_sum);
		dramc_set_vref_dly(vref_dly, delay);

		if (tmp_win_sum < vref_dly->max_win_sum * 95 / 100) {
			dramc_dbg("type:%d best vref found[%d], early break! (%d < %d)\n",
				type, vref_dly->best_vref, tmp_win_sum,
				vref_dly->max_win_sum * 95 / 100);
			return true;
		}

		break;

	default:
		dramc_err("error calibration type: %d\n", type);
		break;
	}

	return false;
}

static void dramc_set_rx_dqdqs_dly(u8 chn, u8 rank, s32 dly)
{
	if (dly <= 0) {
		/* Set DQS delay */
		dramc_set_rx_dly_factor(chn, rank, RX_DQS, -dly);
		dram_phy_reset(chn);
	} else {
		/* Setup time calibration */
		dramc_set_rx_dly_factor(chn, rank, RX_DQM, dly);
		dram_phy_reset(chn);
		dramc_set_rx_dly_factor(chn, rank, RX_DQ, dly);
	}
}

static void dramc_set_tx_best_dly_factor(u8 chn, u8 rank_start, u8 type,
		struct per_byte_dly *tx_perbyte_dly, u16 *dq_precal_dly,
		u8 use_delay_cell, u32 *byte_dly_cell)
{
	u32 dq_large = 0, dq_large_oen = 0, dq_small = 0, dq_small_oen = 0, adjust_center = 1;
	u32 dqm_large = 0, dqm_large_oen = 0, dqm_small = 0, dqm_small_oen = 0;
	u16 dq_oen[DQS_NUMBER] = {0}, dqm_oen[DQS_NUMBER] = {0};
	struct tx_dly_tune dqdly_tune[DQS_NUMBER] = {0};
	struct tx_dly_tune dqmdly_tune[DQS_NUMBER] = {0};

	for (size_t i = 0; i < DQS_NUMBER; i++) {
		dramc_transfer_dly_tune(chn, tx_perbyte_dly[i].final_dly,
			adjust_center, &dqdly_tune[i]);
		dramc_transfer_dly_tune(chn, dq_precal_dly[i],
			adjust_center, &dqmdly_tune[i]);

		dq_large += dqdly_tune[i].coarse_tune_large << (i * 4);
		dq_large_oen += dqdly_tune[i].coarse_tune_large_oen << (i * 4);
		dq_small += dqdly_tune[i].coarse_tune_small << (i * 4);
		dq_small_oen += dqdly_tune[i].coarse_tune_small_oen << (i * 4);

		dqm_large += dqmdly_tune[i].coarse_tune_large << (i * 4);
		dqm_large_oen += dqmdly_tune[i].coarse_tune_large_oen << (i * 4);
		dqm_small += dqmdly_tune[i].coarse_tune_small << (i * 4);
		dqm_small_oen += dqmdly_tune[i].coarse_tune_small_oen << (i * 4);

		dq_oen[i] = (dqdly_tune[i].coarse_tune_large_oen << 3) +
			(dqdly_tune[i].coarse_tune_small_oen << 5) + dqdly_tune[i].fine_tune;
		dqm_oen[i] = (dqmdly_tune[i].coarse_tune_large_oen << 3) +
			(dqmdly_tune[i].coarse_tune_small_oen << 5) +
			dqmdly_tune[i].fine_tune;
	}

	for (size_t rank = rank_start; rank < RANK_MAX; rank++) {
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dq[0],
			0x77777777, dq_large | (dq_large_oen << 16));
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dq[2],
			0x77777777, dq_small | (dq_small_oen << 16));
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dq[1],
			0x77777777, dqm_large | (dqm_large_oen << 16));
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dq[3],
			0x77777777, dqm_small | (dqm_small_oen << 16));

		for (size_t byte = 0; byte < 2; byte++)
			SET32_BITFIELDS(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[7],
					FINE_TUNE_DQ, dqdly_tune[byte].fine_tune,
					FINE_TUNE_DQM, dqmdly_tune[byte].fine_tune);

		if (use_delay_cell == 1) {
			for (size_t byte = 0; byte < DQS_NUMBER; byte++)
				write32(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[0],
					byte_dly_cell[byte]);
		}

		if (type != TX_WIN_DQ_ONLY)
			continue;

		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].fine_tune, 0x3f3f3f3f,
			(dqdly_tune[0].fine_tune << 8) | (dqdly_tune[1].fine_tune << 0) |
			(dqmdly_tune[0].fine_tune << 24) | (dqmdly_tune[1].fine_tune << 16));

		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].dqs2dq_cal1, 0x7ff | (0x7ff << 16),
			(dqdly_tune[0].fine_tune << 0) | (dqdly_tune[1].fine_tune << 16));
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].dqs2dq_cal2, 0x7ff | (0x7ff << 16),
			(dqdly_tune[0].fine_tune << 0) | (dqdly_tune[1].fine_tune << 16));
		clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].dqs2dq_cal5, 0x7ff | (0x7ff << 16),
			(dqmdly_tune[0].fine_tune << 0) | (dqmdly_tune[1].fine_tune << 16));
	}
}

static void dramc_set_rx_best_dly_factor(u8 chn, u8 rank,
		struct win_perbit_dly *dly, s32 *dqsdly_byte, s32 *dqmdly_byte)
{
	u32 value;

	/* set dqs delay, (dqm delay) */
	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		value = (dqsdly_byte[byte] << 24) | (dqsdly_byte[byte] << 16) |
			(dqmdly_byte[byte] << 8) | (dqmdly_byte[byte] << 0);
		clrsetbits_le32(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[6], 0x7f7f3f3f, value);
	}
	dram_phy_reset(chn);

	/* set dq delay */
	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		for (u8 bit = 0; bit < DQS_BIT_NUMBER; bit += 2) {
			u8 index = bit + byte * DQS_BIT_NUMBER;
			u8 dq_num = 2 + bit / 2;
			value = (dly[index + 1].best_dqdly << 24) |
				(dly[index + 1].best_dqdly << 16) |
				(dly[index].best_dqdly << 8) | (dly[index].best_dqdly << 0);

			clrsetbits_le32(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[dq_num],
				0x3f3f3f3f, value);
		}
	}
}

static void dramc_set_dqdqs_dly(u8 chn, u8 rank, enum CAL_TYPE type, u8 *small_value, s32 dly)
{
	if (type == RX_WIN_RD_DQC || type == RX_WIN_TEST_ENG)
		dramc_set_rx_dqdqs_dly(chn, rank, dly);
	else
		dramc_set_tx_dly_factor(chn, rank, type, small_value, dly);
}

static void dramc_set_tx_dly_center(struct per_byte_dly *center_dly,
				    const struct win_perbit_dly *vref_dly)
{
	int index;
	struct per_byte_dly *dly;

	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		dly = &center_dly[byte];
		dly->min_center = 0xffff;
		dly->max_center = 0;

		for (u8 bit = 0; bit < DQS_BIT_NUMBER; bit++) {
			index = bit + 8 * byte;
			if (vref_dly[index].win_center < dly->min_center)
				dly->min_center = vref_dly[index].win_center;
			if (vref_dly[index].win_center > dly->max_center)
				dly->max_center = vref_dly[index].win_center;
		}
		dramc_dbg("center_dly[%d].min_center = %d, "
			  "center_dly[%d].max_center = %d\n",
			  byte, center_dly[byte].min_center,
			  byte, center_dly[byte].max_center);
	}
}

static void dramc_set_tx_best_dly(u8 chn, u8 rank, bool bypass_tx,
		struct win_perbit_dly *vref_dly, enum CAL_TYPE type, u8 freq_group,
		u16 *tx_dq_precal_result, u16 dly_cell_unit, const struct sdram_params *params,
		const bool fast_calib)
{
	int index, clock_rate;
	u8 use_delay_cell;
	u32 byte_dly_cell[DQS_NUMBER] = {0};
	struct per_byte_dly center_dly[DQS_NUMBER];
	u16 tune_diff, dq_delay_cell[DQ_DATA_WIDTH];

	switch (freq_group) {
	case LP4X_DDR1600:
		clock_rate = 800;
		break;
	case LP4X_DDR2400:
		clock_rate = 1200;
		break;
	case LP4X_DDR3200:
		clock_rate = 1600;
		break;
	case LP4X_DDR3600:
		clock_rate = 1866;
		break;
	default:
		die("Invalid DDR frequency group %u\n", freq_group);
		return;
	}

	if (type == TX_WIN_DQ_ONLY && get_freq_fsq(freq_group) == FSP_1)
		use_delay_cell = 1;
	else
		use_delay_cell = 0;

	if (fast_calib && bypass_tx) {
		dramc_dbg("bypass TX\n");
		for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
			center_dly[byte].min_center = params->tx_center_min[chn][rank][byte];
			center_dly[byte].max_center = params->tx_center_max[chn][rank][byte];
			for (u8 bit = 0; bit < DQS_BIT_NUMBER; bit++) {
				index = bit + 8 * byte;
				vref_dly[index].win_center =
					params->tx_win_center[chn][rank][index];
				vref_dly[index].best_first =
					params->tx_first_pass[chn][rank][index];
				vref_dly[index].best_last =
					params->tx_last_pass[chn][rank][index];
			}
		}
	} else {
		dramc_set_tx_dly_center(center_dly, vref_dly);
	}

	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		if (use_delay_cell == 0) {
			center_dly[byte].final_dly = (center_dly[byte].min_center +
				center_dly[byte].max_center) >> 1;
			tx_dq_precal_result[byte] = center_dly[byte].final_dly;
		} else {
			center_dly[byte].final_dly = center_dly[byte].min_center;
			tx_dq_precal_result[byte] = (center_dly[byte].min_center
						+ center_dly[byte].max_center) >> 1;

			for (u8 bit = 0; bit < DQS_BIT_NUMBER; bit++) {
				index = bit + 8 * byte;
				tune_diff = vref_dly[index].win_center -
					center_dly[byte].min_center;
				dq_delay_cell[index] = ((tune_diff * 100000000) /
					(clock_rate / 2 * 64)) / dly_cell_unit;
				byte_dly_cell[byte] |= (dq_delay_cell[index] << (bit * 4));
			}
		}
	}

	dramc_set_tx_best_dly_factor(chn, rank, type, center_dly, tx_dq_precal_result,
			use_delay_cell, byte_dly_cell);
}

static int dramc_set_rx_best_dly(u8 chn, u8 rank, struct win_perbit_dly *perbit_dly)
{
	u8 bit_first, bit_last;
	u16 u2TmpDQMSum;
	s32 dqsdly_byte[DQS_NUMBER] = {0x0}, dqm_dly_byte[DQS_NUMBER] = {0x0};

	for (u8 byte = 0; byte < DQS_NUMBER; byte++) {
		u2TmpDQMSum = 0;

		bit_first = DQS_BIT_NUMBER * byte;
		bit_last = DQS_BIT_NUMBER * byte + DQS_BIT_NUMBER - 1;
		dqsdly_byte[byte] = 64;

		for (u8 bit = bit_first; bit <= bit_last; bit++) {
			if (perbit_dly[bit].win_center < dqsdly_byte[byte])
				dqsdly_byte[byte] = perbit_dly[bit].win_center;
		}
		dqsdly_byte[byte] = (dqsdly_byte[byte] > 0) ?  0 : -dqsdly_byte[byte];

		for (u8 bit = bit_first; bit <= bit_last; bit++) {
			perbit_dly[bit].best_dqdly = dqsdly_byte[byte] +
				perbit_dly[bit].win_center;
			u2TmpDQMSum += perbit_dly[bit].best_dqdly;
		}

		dqm_dly_byte[byte] = u2TmpDQMSum / DQS_BIT_NUMBER;
	}

	dramc_set_rx_best_dly_factor(chn, rank, perbit_dly, dqsdly_byte, dqm_dly_byte);
	return 0;
}

static void dramc_get_vref_prop(u8 rank, enum CAL_TYPE type, u8 fsp,
		u8 *vref_scan_en, u8 *vref_begin, u8 *vref_end)
{
	if (type == RX_WIN_TEST_ENG && rank == RANK_0) {
		*vref_scan_en = 1;
		if (fsp == FSP_0)
			*vref_begin = 0x18;
		else
			*vref_begin = 0;
		*vref_end = RX_VREF_END;
	} else if (type == TX_WIN_DQ_ONLY) {
		*vref_scan_en = 1;
		if (fsp == FSP_0) {
			*vref_begin = 27 - 5;
			*vref_end = 27 + 5;
		} else {
			*vref_begin = TX_VREF_BEGIN;
			*vref_end = TX_VREF_END;
		}
	} else {
		*vref_scan_en = 0;
		*vref_begin = 0;
		*vref_end = 1;
	}
}

static u32 dram_k_perbit(u8 chn, enum CAL_TYPE type)
{
	u32 err_value;

	if (type == RX_WIN_RD_DQC) {
		err_value = dramc_rx_rd_dqc_run(chn);
	} else if (type == RX_WIN_TEST_ENG) {
		err_value = dramc_engine2_run(chn, TE_OP_WRITE_READ_CHECK);
	} else {
		dramc_engine2_setpat(chn, true);
		err_value = dramc_engine2_run(chn, TE_OP_WRITE_READ_CHECK);
		dramc_engine2_setpat(chn, false);
		err_value |= dramc_engine2_run(chn, TE_OP_WRITE_READ_CHECK);
	}
	return err_value;
}

static void dramc_window_perbit_cal_partial(u8 chn, u8 rank,
		s16 dly_begin, s16 dly_end, s16 dly_step,
		enum CAL_TYPE type, u8 *small_value, u8 vref_scan_enable,
		struct win_perbit_dly *win_perbit)
{
	u32 finish_bit = 0;

	for (s16 dly = dly_begin; dly < dly_end; dly += dly_step) {
		dramc_set_dqdqs_dly(chn, rank, type, small_value, dly);

		u32 err_value = dram_k_perbit(chn, type);
		if (!vref_scan_enable)
			dramc_dbg("%d ", dly);

		for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++) {
			bool bit_fail = (err_value & ((u32) 1 << bit)) != 0;

			/* pass window bigger than 7,
			   * consider as real pass window.
			   */
			if (dramc_check_dqdqs_win(&(win_perbit[bit]),
						  dly, dly_end, bit_fail) > 7)
				finish_bit |= (1 << bit);

			if (vref_scan_enable)
				continue;

			dramc_dbg("%s", bit_fail ? "x" : "o");
			if (bit % DQS_BIT_NUMBER == 7)
				dramc_dbg(" ");
		}

		if (!vref_scan_enable)
			dramc_dbg(" [MSB]\n");

		if (finish_bit == 0xffff && (err_value & 0xffff) == 0xffff) {
			dramc_dbg("all bits window found, "
				  "early break! delay=%#x\n", dly);
			break;
		}
	}
}

static u8 dramc_window_perbit_cal(u8 chn, u8 rank, u8 freq_group,
		enum CAL_TYPE type, const struct sdram_params *params,
		const bool fast_calib)
{
	u8 vref = 0, vref_begin = 0, vref_end = 1, vref_step = 1, vref_use = 0;
	u8 vref_scan_enable = 0, small_reg_value = 0xff;
	s16 dly_begin = 0, dly_end = 0, dly_step = 1;
	u32 dummy_rd_bak_engine2 = 0, finish_bit, win_min_max = 0;
	static u16 dq_precal_result[DQS_NUMBER];
	struct vref_perbit_dly vref_dly;
	struct win_perbit_dly win_perbit[DQ_DATA_WIDTH];
	u16 dly_cell_unit = params->delay_cell_unit;

	u8 fsp = get_freq_fsq(freq_group);
	u8 vref_range = !fsp;
	bool bypass_tx = !fsp;

	dramc_get_vref_prop(rank, type, fsp,
		&vref_scan_enable, &vref_begin, &vref_end);
	dramc_get_dly_range(chn, rank, type, freq_group, dq_precal_result,
		&dly_begin, &dly_end, params);

	if (fast_calib) {
		if (type == RX_WIN_TEST_ENG && vref_scan_enable == 1) {
			vref_begin = params->rx_vref[chn];
			vref_end = vref_begin + 1;
			dramc_dbg("bypass RX vref: %d\n", vref_begin);
		} else if (type == TX_WIN_DQ_ONLY) {
			vref_begin = params->tx_vref[chn][rank];
			vref_end = vref_begin + 1;
			dramc_dbg("bypass TX vref: %d\n", vref_begin);
		}
		vref_dly.best_vref = vref_begin;
	}

	if ((type == RX_WIN_RD_DQC || type == RX_WIN_TEST_ENG) && fsp == FSP_0)
		dly_step = 2;

	dramc_dbg("[channel %d] [rank %d] type: %d, vref_enable: %d, vref range[%d : %d]\n",
		chn, rank, type, vref_scan_enable, vref_begin, vref_end);

	if (type == TX_WIN_DQ_ONLY || type == TX_WIN_DQ_DQM) {
		for (size_t byte = 0; byte < 2; byte++) {
			write32(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[0], 0);
			clrbits_le32(&ch[chn].phy.shu[0].rk[rank].b[byte].dq[1],
				0xf);
		}
		setbits_le32(&ch[chn].phy.misc_ctrl1, 0x1 << 7);
		setbits_le32(&ch[chn].ao.dqsoscr, 0x1 << 7);
		if (fsp == FSP_1)
			vref_step = 2;
	}

	if (fast_calib && bypass_tx &&
	    (type == TX_WIN_DQ_ONLY || type == TX_WIN_DQ_DQM)) {
		dramc_set_tx_best_dly(chn, rank, true, vref_dly.perbit_dly,
			type, freq_group, dq_precal_result, dly_cell_unit,
			params, fast_calib);

		if (vref_scan_enable)
			dramc_set_vref(chn, rank, type, vref_dly.best_vref);
		return 0;
	}

	if (type == RX_WIN_RD_DQC) {
		dramc_rx_rd_dqc_init(chn, rank);
	} else {
		if (type == RX_WIN_TEST_ENG)
			dramc_rx_vref_pre_setting(chn);
		dummy_rd_bak_engine2 = read32(&ch[chn].ao.dummy_rd);
		dramc_engine2_init(chn, rank, TEST2_1_CAL, TEST2_2_CAL, false);
	}

	vref_dly.max_win_sum = 0;
	for (vref = vref_begin; vref < vref_end; vref += vref_step) {
		small_reg_value = 0xff;
		finish_bit = 0;
		if (type == TX_WIN_DQ_ONLY)
			vref_use = vref | (vref_range << 6);
		else
			vref_use = vref;

		for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++) {
			win_perbit[bit].first_pass = PASS_RANGE_NA;
			win_perbit[bit].last_pass = PASS_RANGE_NA;
			win_perbit[bit].best_first = PASS_RANGE_NA;
			win_perbit[bit].best_last = PASS_RANGE_NA;
		}

		if (vref_scan_enable)
			dramc_set_vref(chn, rank, type, vref_use);

		if (type == RX_WIN_RD_DQC || type == RX_WIN_TEST_ENG) {
			dramc_set_rx_dly_factor(chn, rank,
				RX_DQM, FIRST_DQ_DELAY);
			dramc_set_rx_dly_factor(chn, rank,
				RX_DQ, FIRST_DQ_DELAY);
		}

		if (fast_calib &&
		    (type == RX_WIN_RD_DQC || type == RX_WIN_TEST_ENG)) {
			dramc_dbg("bypass RX params\n");
			for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++) {
				win_perbit[bit].best_first =
					params->rx_firspass[chn][rank][bit];
				win_perbit[bit].best_last =
					params->rx_lastpass[chn][rank][bit];
			}
		} else {
			dramc_window_perbit_cal_partial(chn, rank,
				dly_begin, dly_end, dly_step,
				type, &small_reg_value,
				vref_scan_enable, win_perbit);
		}

		for (size_t bit = 0; bit < DQ_DATA_WIDTH; bit++)
			dramc_dbg("Dq[%zd] win width (%d ~ %d)  %d\n", bit,
				win_perbit[bit].best_first, win_perbit[bit].best_last,
				win_perbit[bit].best_last - win_perbit[bit].best_first);

		if (dramk_calc_best_vref(type, vref_use, &vref_dly,
			win_perbit, &win_min_max))
			break;
	}

	if (type == RX_WIN_RD_DQC)
		dramc_rx_rd_dqc_end(chn);
	else
		dramc_engine2_end(chn, dummy_rd_bak_engine2);

	if (vref_scan_enable && type == RX_WIN_TEST_ENG)
		dramc_set_vref(chn, rank, type, vref_dly.best_vref);

	if (type == RX_WIN_RD_DQC || type == RX_WIN_TEST_ENG)
		dramc_set_rx_best_dly(chn, rank, vref_dly.perbit_dly);
	else
		dramc_set_tx_best_dly(chn, rank, false,
			vref_dly.perbit_dly, type, freq_group,
			dq_precal_result, dly_cell_unit, params, fast_calib);

	if (vref_scan_enable && type == TX_WIN_DQ_ONLY)
		dramc_set_vref(chn, rank, type, vref_dly.best_vref);

	return 0;
}

static void dramc_dle_factor_handler(u8 chn, u8 val, u8 freq_group)
{
	u8 start_ext2 = 0, start_ext3 = 0, last_ext2 = 0, last_ext3 = 0;

	val = MAX(val, 2);
	SET32_BITFIELDS(&ch[chn].ao.shu[0].conf[1],
			SHU_CONF1_DATLAT, val,
			SHU_CONF1_DATLAT_DSEL, val - 2,
			SHU_CONF1_DATLAT_DSEL_PHY, val - 2);

	if (freq_group == LP4X_DDR3200 || freq_group == LP4X_DDR3600)
		start_ext2 = 1;

	if (val >= 24)
		last_ext2 = last_ext3 = 1;
	else if (val >= 18)
		last_ext2 = 1;

	SET32_BITFIELDS(&ch[chn].ao.shu[0].pipe,
			SHU_PIPE_READ_START_EXTEND1, 1,
			SHU_PIPE_DLE_LAST_EXTEND1, 1,
			SHU_PIPE_READ_START_EXTEND2, start_ext2,
			SHU_PIPE_DLE_LAST_EXTEND2, last_ext2,
			SHU_PIPE_READ_START_EXTEND3, start_ext3,
			SHU_PIPE_DLE_LAST_EXTEND3, last_ext3);
	dram_phy_reset(chn);
}

static u8 dramc_rx_datlat_cal(u8 chn, u8 rank, u8 freq_group,
		const struct sdram_params *params, const bool fast_calib,
		bool *test_passed)
{
	u32 datlat, begin = 0, first = 0, sum = 0, best_step;
	u32 datlat_start = 7;

	*test_passed = true;
	best_step = READ32_BITFIELD(&ch[chn].ao.shu[0].conf[1], SHU_CONF1_DATLAT);

	dramc_dbg("[DATLAT] start. CH%d RK%d DATLAT Default: 0x%x\n",
		chn, rank, best_step);

	u32 dummy_rd_backup = read32(&ch[chn].ao.dummy_rd);
	dramc_engine2_init(chn, rank, TEST2_1_CAL, TEST2_2_CAL, false);

	if (fast_calib) {
		best_step = params->rx_datlat[chn][rank];
		dramc_dbg("bypass DATLAT, best_step: %d\n", best_step);
	} else {
		for (datlat = datlat_start; datlat < DATLAT_TAP_NUMBER; datlat++) {
			dramc_dle_factor_handler(chn, datlat, freq_group);

			u32 err = dramc_engine2_run(chn, TE_OP_WRITE_READ_CHECK);
			if (err == 0) {
				if (begin == 0) {
					first = datlat;
					begin = 1;
				}
				if (begin == 1) {
					sum++;
					if (sum > 4)
						break;
				}
			} else {
				if (begin == 1)
					begin = 0xff;
			}

			dramc_dbg("Datlat=%2d, err_value=0x%4x, sum=%d\n", datlat, err, sum);
		}

		dramc_engine2_end(chn, dummy_rd_backup);

		*test_passed = (sum != 0);
		if (!*test_passed) {
			dramc_err("DRAM memory test failed\n");
			return 0;
		}

		if (sum <= 3)
			best_step = first + (sum >> 1);
		else
			best_step = first + 2;
		dramc_dbg("First_step=%d, total pass=%d, best_step=%d\n",
			begin, sum, best_step);
	}

	dramc_dle_factor_handler(chn, best_step, freq_group);

	SET32_BITFIELDS(&ch[chn].ao.padctrl,
			PADCTRL_DQIENQKEND, 1, PADCTRL_DQIENLATEBEGIN, 1);

	return (u8) best_step;
}

static void dramc_dual_rank_rx_datlat_cal(u8 chn, u8 freq_group, u8 datlat0, u8 datlat1)
{
	u8 final_datlat = MAX(datlat0, datlat1);
	dramc_dle_factor_handler(chn, final_datlat, freq_group);
}

static void dramc_rx_dqs_gating_post_process(u8 chn, u8 freq_group)
{
	s8 dqsinctl;
	u32 read_dqsinctl, rankinctl_root, reg_tx_dly_dqsgated_min = 3;
	u8 txdly_cal_min = 0xff, txdly_cal_max = 0, tx_dly_dqs_gated = 0;
	u32 best_coarse_tune2t[RANK_MAX][DQS_NUMBER];
	u32 best_coarse_tune2t_p1[RANK_MAX][DQS_NUMBER];

	if (freq_group == LP4X_DDR3200 || freq_group == LP4X_DDR3600)
		reg_tx_dly_dqsgated_min = 2;
	else
		reg_tx_dly_dqsgated_min = 1;

	/* get TXDLY_Cal_min and TXDLY_Cal_max value */
	for (size_t rank = 0; rank < RANK_MAX; rank++) {
		u32 dqsg0 = read32(&ch[chn].ao.shu[0].rk[rank].selph_dqsg0);
		for (size_t dqs = 0; dqs < DQS_NUMBER; dqs++) {
			best_coarse_tune2t[rank][dqs] = (dqsg0 >> (dqs * 8)) & 0x7;
			best_coarse_tune2t_p1[rank][dqs] = (dqsg0 >> (dqs * 8 + 4)) & 0x7;
			dramc_dbg("Rank%zd best DQS%zd dly(2T,(P1)2T)=(%d, %d)\n",
				rank, dqs, best_coarse_tune2t[rank][dqs],
				best_coarse_tune2t_p1[rank][dqs]);

			tx_dly_dqs_gated = best_coarse_tune2t[rank][dqs];
			txdly_cal_min = MIN(txdly_cal_min, tx_dly_dqs_gated);

			tx_dly_dqs_gated = best_coarse_tune2t_p1[rank][dqs];
			txdly_cal_max = MAX(txdly_cal_max, tx_dly_dqs_gated);
		}
	}

	dqsinctl = reg_tx_dly_dqsgated_min - txdly_cal_min;
	dramc_dbg("Dqsinctl:%d, dqsgated_min %d, txdly_cal_min %d, txdly_cal_max %d\n",
		dqsinctl, reg_tx_dly_dqsgated_min, txdly_cal_min, txdly_cal_max);

	if (dqsinctl != 0) {
		txdly_cal_min += dqsinctl;
		txdly_cal_max += dqsinctl;

		for (size_t rank = 0; rank < RANK_MAX; rank++) {
			dramc_dbg("Rank: %zd\n", rank);
			for (size_t dqs = 0; dqs < DQS_NUMBER; dqs++) {
				best_coarse_tune2t[rank][dqs] += dqsinctl;
				best_coarse_tune2t_p1[rank][dqs] += dqsinctl;

				dramc_dbg("Best DQS%zd dly(2T) = (%d)\n",
					dqs, best_coarse_tune2t[rank][dqs]);
				dramc_dbg("Best DQS%zd P1 dly(2T) = (%d)\n",
					dqs, best_coarse_tune2t_p1[rank][dqs]);
			}

			clrsetbits_le32(&ch[chn].ao.shu[0].rk[rank].selph_dqsg0,
				0x77777777,
				(best_coarse_tune2t[rank][0] << 0) |
				(best_coarse_tune2t[rank][1] << 8) |
				(best_coarse_tune2t_p1[rank][0] << 4) |
				(best_coarse_tune2t_p1[rank][1] << 12));
		}
	}

	read_dqsinctl = READ32_BITFIELD(&ch[chn].ao.shu[0].rk[0].dqsctl,
					SHURK_DQSCTL_DQSINCTL) - dqsinctl;
	rankinctl_root = (read_dqsinctl >= 2) ? (read_dqsinctl - 2) : 0;

	SET32_BITFIELDS(&ch[chn].ao.shu[0].rk[0].dqsctl, SHURK_DQSCTL_DQSINCTL, read_dqsinctl);
	SET32_BITFIELDS(&ch[chn].ao.shu[0].rk[1].dqsctl, SHURK_DQSCTL_DQSINCTL, read_dqsinctl);
	clrsetbits_le32(&ch[chn].ao.shu[0].rankctl,
		(0xf << 28) | (0xf << 20) | (0xf << 24) | 0xf,
		(read_dqsinctl << 28) | (rankinctl_root << 20) |
		(rankinctl_root << 24) | rankinctl_root);

	u8 ROEN = read32(&ch[chn].ao.shu[0].odtctrl) & 0x1;
	clrsetbits_le32(&ch[chn].ao.shu[0].rodtenstb, (0xffff << 8) | (0x3f << 2) | (0x1),
			(0xff << 8) | (0x9 << 2) | ROEN);
}

int dramc_calibrate_all_channels(const struct sdram_params *pams, u8 freq_group,
				 const struct mr_value *mr)
{
	bool fast_calib;
	switch (pams->source) {
	case DRAMC_PARAM_SOURCE_SDRAM_CONFIG:
		fast_calib = false;
		break;
	case DRAMC_PARAM_SOURCE_FLASH:
		fast_calib = true;
		break;
	default:
		die("Invalid DRAM param source %u\n", pams->source);
		return -1;
	}

	bool test_passed;
	u8 rx_datlat[RANK_MAX] = {0};
	for (u8 chn = 0; chn < CHANNEL_MAX; chn++) {
		for (u8 rk = RANK_0; rk < RANK_MAX; rk++) {
			dramc_dbg("Start K: freq=%d, ch=%d, rank=%d\n",
				  freq_group, chn, rk);
			dramc_cmd_bus_training(chn, rk, freq_group, pams,
				fast_calib);
			dramc_write_leveling(chn, rk, freq_group, pams->wr_level);
			dramc_auto_refresh_switch(chn, true);

			dramc_rx_dqs_gating_cal(chn, rk, freq_group, pams,
				fast_calib, mr);
			dramc_window_perbit_cal(chn, rk, freq_group,
				RX_WIN_RD_DQC, pams, fast_calib);
			dramc_window_perbit_cal(chn, rk, freq_group,
				TX_WIN_DQ_DQM, pams, fast_calib);
			dramc_window_perbit_cal(chn, rk, freq_group,
				TX_WIN_DQ_ONLY, pams, fast_calib);
			rx_datlat[rk] = dramc_rx_datlat_cal(chn, rk, freq_group,
				pams, fast_calib, &test_passed);
			if (!test_passed)
				return -2;
			dramc_window_perbit_cal(chn, rk, freq_group,
				RX_WIN_TEST_ENG, pams, fast_calib);
			dramc_auto_refresh_switch(chn, false);
		}

		dramc_rx_dqs_gating_post_process(chn, freq_group);
		dramc_dual_rank_rx_datlat_cal(chn, freq_group, rx_datlat[0], rx_datlat[1]);
	}
	return 0;
}
