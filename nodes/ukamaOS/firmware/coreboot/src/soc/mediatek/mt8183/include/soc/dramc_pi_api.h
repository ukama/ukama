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

#ifndef _DRAMC_PI_API_MT8183_H
#define _DRAMC_PI_API_MT8183_H

#include <types.h>
#include <soc/emi.h>
#include <console/console.h>

#define dramc_err(_x_...) printk(BIOS_ERR, _x_)
#define dramc_show(_x_...) printk(BIOS_INFO, _x_)
#if CONFIG(DEBUG_DRAM)
#define dramc_dbg(_x_...) printk(BIOS_DEBUG, _x_)
#else
#define dramc_dbg(_x_...)
#endif

#define DATLAT_TAP_NUMBER 32

#define DRAMC_BROADCAST_ON 0x1f
#define DRAMC_BROADCAST_OFF 0x0
#define TX_DQ_COARSE_TUNE_TO_FINE_TUNE_TAP 64

#define IMP_LP4X_TERM_VREF_SEL		0x1b
#define IMP_DRVP_LP4X_UNTERM_VREF_SEL	0x1a
#define IMP_DRVN_LP4X_UNTERM_VREF_SEL	0x16
#define IMP_TRACK_LP4X_UNTERM_VREF_SEL	0x1a

enum dram_te_op {
	TE_OP_WRITE_READ_CHECK = 0,
	TE_OP_READ_CHECK
};

enum {
	PASS_RANGE_NA = 0x7fff
};

enum {
	GATING_PATTERN_NUM = 0x23,
	GATING_GOLDEND_DQSCNT = 0x4646
};

enum {
	IMPCAL_STAGE_DRVP = 0x1,
	IMPCAL_STAGE_DRVN,
	IMPCAL_STAGE_TRACKING
};

enum {
	DQS_GW_COARSE_STEP = 1,
	DQS_GW_FINE_END = 32,
	DQS_GW_FINE_STEP = 4,
	RX_DQS_CTL_LOOP = 8,
	RX_DLY_DQSIENSTB_LOOP = 32
};

enum {
	DLL_MASTER = 0,
	DLL_SLAVE,
};

struct reg_value {
	u32 *addr;
	u32 value;
};

#define _SELPH_DQS_BITS(l, h) ((l << 0) | (l << 4) | (l << 8) | (l << 12) | \
			       (h << 16) | (h << 20) | (h << 24) | (h << 28))

enum {
	DQ_DIV_SHIFT = 3,
	DQ_DIV_MASK = BIT(DQ_DIV_SHIFT) - 1,
	OEN_SHIFT = 16,

	SELPH_DQS0 = _SELPH_DQS_BITS(0x3, 0x3),
	SELPH_DQS1 = _SELPH_DQS_BITS(0x4, 0x1),
	SELPH_DQS0_1600 = _SELPH_DQS_BITS(0x2, 0x1),
	SELPH_DQS1_1600 = _SELPH_DQS_BITS(0x1, 0x6),
	SELPH_DQS0_2400 = _SELPH_DQS_BITS(0x3, 0x2),
	SELPH_DQS1_2400 = _SELPH_DQS_BITS(0x1, 0x6),
	SELPH_DQS0_3600 = _SELPH_DQS_BITS(0x4, 0x3),
	SELPH_DQS1_3600 = _SELPH_DQS_BITS(0x1, 0x6),
};

void dramc_get_rank_size(u64 *dram_rank_size);
void dramc_runtime_config(void);
void dramc_set_broadcast(u32 onoff);
u32 dramc_get_broadcast(void);
u8 get_freq_fsq(u8 freq_group);
void dramc_init(const struct sdram_params *params, u8 freq_group,
		struct dram_shared_data *shared);
void dramc_sw_impedance_save_reg(u8 freq_group,
				 const struct dram_impedance *impedance);
void dramc_sw_impedance_cal(const struct sdram_params *params, u8 term_option,
			    struct dram_impedance *impedance);
void dramc_apply_config_before_calibration(u8 freq_group);
void dramc_apply_config_after_calibration(const struct mr_value *mr);
int dramc_calibrate_all_channels(const struct sdram_params *pams,
				 u8 freq_group, const struct mr_value *mr);
void dramc_hw_gating_onoff(u8 chn, bool onoff);
void dramc_enable_phy_dcm(bool bEn);
void dramc_mode_reg_write(u8 chn, u8 mr_idx, u8 value);
void dramc_cke_fix_onoff(u8 chn, bool fix_on, bool fix_off);

#endif /* _DRAMC_PI_API_MT8183_H */
