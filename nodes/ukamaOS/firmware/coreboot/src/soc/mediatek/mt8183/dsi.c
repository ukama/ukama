/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 MediaTek Inc.
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
#include <device/mmio.h>
#include <delay.h>
#include <soc/dsi.h>
#include <soc/pll.h>

void mtk_dsi_configure_mipi_tx(int data_rate, u32 lanes)
{
	unsigned int txdiv, txdiv0, txdiv1;
	u64 pcw;

	if (data_rate >= 2000) {
		txdiv = 1;
		txdiv0 = 0;
		txdiv1 = 0;
	} else if (data_rate >= 1000) {
		txdiv = 2;
		txdiv0 = 1;
		txdiv1 = 0;
	} else if (data_rate >= 500) {
		txdiv = 4;
		txdiv0 = 2;
		txdiv1 = 0;
	} else if (data_rate > 250) {
		/* Be aware that 250 is a special case that must use txdiv=4. */
		txdiv = 8;
		txdiv0 = 3;
		txdiv1 = 0;
	} else {
		/* MIN = 125 */
		assert(data_rate >= MTK_DSI_DATA_RATE_MIN_MHZ);
		txdiv = 16;
		txdiv0 = 4;
		txdiv1 = 0;
	}

	clrbits_le32(&mipi_tx->pll_con4, BIT(11) | BIT(10));
	setbits_le32(&mipi_tx->pll_pwr, AD_DSI_PLL_SDM_PWR_ON);
	udelay(30);
	clrbits_le32(&mipi_tx->pll_pwr, AD_DSI_PLL_SDM_ISO_EN);

	pcw = (u64)data_rate * (1 << txdiv0) * (1 << txdiv1);
	pcw <<= 24;
	pcw /= CLK26M_HZ / MHz;

	write32(&mipi_tx->pll_con0, pcw);
	clrsetbits_le32(&mipi_tx->pll_con1, RG_DSI_PLL_POSDIV, txdiv0 << 8);
	udelay(30);
	setbits_le32(&mipi_tx->pll_con1, RG_DSI_PLL_EN);

	/* BG_LPF_EN / BG_CORE_EN */
	write32(&mipi_tx->lane_con, 0x3fff0180);
	udelay(40);
	write32(&mipi_tx->lane_con, 0x3fff00c0);

	/* Switch OFF each Lane */
	clrbits_le32(&mipi_tx->d0_sw_ctl_en, DSI_SW_CTL_EN);
	clrbits_le32(&mipi_tx->d1_sw_ctl_en, DSI_SW_CTL_EN);
	clrbits_le32(&mipi_tx->d2_sw_ctl_en, DSI_SW_CTL_EN);
	clrbits_le32(&mipi_tx->d3_sw_ctl_en, DSI_SW_CTL_EN);
	clrbits_le32(&mipi_tx->ck_sw_ctl_en, DSI_SW_CTL_EN);

	setbits_le32(&mipi_tx->ck_ckmode_en, DSI_CK_CKMODE_EN);
}

void mtk_dsi_reset(void)
{
	write32(&dsi0->dsi_force_commit,
		DSI_FORCE_COMMIT_USE_MMSYS | DSI_FORCE_COMMIT_ALWAYS);
	write32(&dsi0->dsi_con_ctrl, 1);
	write32(&dsi0->dsi_con_ctrl, 0);
}
