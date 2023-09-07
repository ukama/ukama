/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Qualcomm Inc.
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

#include <device/mmio.h>
#include <types.h>
#include <commonlib/helpers.h>
#include <assert.h>
#include <soc/symbols.h>
#include <soc/clock.h>

#define DIV(div) (2*div - 1)

#define AOP_LOADED_SIGNAL_FLAG 0x11223344

struct clock_config qup_cfg[] = {
	{
		.hz = 7372800,
		.src = SRC_GPLL0_EVEN_300MHZ,
		.div = DIV(1),
		.m = 384,
		.n = 15625,
		.d_2 = 15625,
	},
	{
		.hz = 19200*KHz,
		.src = SRC_XO_19_2MHZ,
		.div = DIV(1),
	}
};

struct clock_config qspi_core_cfg[] = {
	{
		.hz = 19200*KHz,
		.src = SRC_XO_19_2MHZ,
		.div = DIV(0),
	},
	{
		.hz = 100*MHz,
		.src = SRC_GPLL0_MAIN_600MHZ,
		.div = DIV(6),
	},
	{
		.hz = 150*MHz,
		.src = SRC_GPLL0_MAIN_600MHZ,
		.div = DIV(4),
	},
	{
		.hz = 300*MHz,
		.src = SRC_GPLL0_MAIN_600MHZ,
		.div = DIV(2),
	}
};

static int clock_configure_gpll0(void)
{
	/* Keep existing GPLL0 configuration, in RUN mode @600Mhz. */
	setbits_le32(&gcc->gpll0.user_ctl,
			1 << CLK_CTL_GPLL_PLLOUT_EVEN_SHFT |
			1 << CLK_CTL_GPLL_PLLOUT_MAIN_SHFT |
			1 << CLK_CTL_GPLL_PLLOUT_ODD_SHFT);
	return 0;
}

static int clock_configure_mnd(struct sdm845_clock *clk, uint32_t m, uint32_t n,
				uint32_t d_2)
{
	setbits_le32(&clk->rcg.cfg,
			RCG_MODE_DUAL_EDGE << CLK_CTL_CFG_MODE_SHFT);

	write32(&clk->m, m & CLK_CTL_RCG_MND_BMSK);
	write32(&clk->n, ~(n-m) & CLK_CTL_RCG_MND_BMSK);
	write32(&clk->d_2, ~(d_2) & CLK_CTL_RCG_MND_BMSK);

	return 0;
}

static int clock_configure(struct sdm845_clock *clk,
				struct clock_config *clk_cfg,
				uint32_t hz, uint32_t num_perfs)
{
	uint32_t reg_val;
	uint32_t idx;

	for (idx = 0; idx < num_perfs; idx++)
		if (hz <= clk_cfg[idx].hz)
			break;

	assert(hz == clk_cfg[idx].hz);

	reg_val = (clk_cfg[idx].src << CLK_CTL_CFG_SRC_SEL_SHFT) |
			(clk_cfg[idx].div << CLK_CTL_CFG_SRC_DIV_SHFT);

	/* Set clock config */
	write32(&clk->rcg.cfg, reg_val);

	if (clk_cfg[idx].m != 0)
		clock_configure_mnd(clk, clk_cfg[idx].m, clk_cfg[idx].n,
				clk_cfg[idx].d_2);

	/* Commit config to RCG*/
	setbits_le32(&clk->rcg.cmd, BIT(CLK_CTL_CMD_UPDATE_SHFT));

	return 0;
}

static bool clock_is_off(u32 *cbcr_addr)
{
	return (read32(cbcr_addr) & CLK_CTL_CBC_CLK_OFF_BMSK);
}

static int clock_enable_vote(void *cbcr_addr, void *vote_addr,
				uint32_t vote_bit)
{

	/* Set clock vote bit */
	setbits_le32(vote_addr, BIT(vote_bit));

	/* Ensure clock is enabled */
	while (clock_is_off(cbcr_addr))
		;

	return 0;
}

static int clock_enable(void *cbcr_addr)
{

	/* Set clock enable bit */
	setbits_le32(cbcr_addr, BIT(CLK_CTL_CBC_CLK_EN_SHFT));

	/* Ensure clock is enabled */
	while (clock_is_off(cbcr_addr))
		;

	return 0;
}

void clock_reset_aop(void)
{
	/* Bring AOP out of RESET */
	uint32_t *mailbox;
	mailbox = (uint32_t *)_aop_ss_msg_ram_drv15;
	*mailbox = AOP_LOADED_SIGNAL_FLAG;
}

void clock_configure_qspi(uint32_t hz)
{
	clock_configure((struct sdm845_clock *)&gcc->qspi_core,
			qspi_core_cfg, hz,
			ARRAY_SIZE(qspi_core_cfg));
	clock_enable(&gcc->qspi_cnoc_ahb_cbcr);
	clock_enable(&gcc->qspi_core_cbcr);
}

int clock_reset_bcr(void *bcr_addr, bool reset)
{
	struct sdm845_bcr *bcr = bcr_addr;

	if (reset)
		setbits_le32(bcr, BIT(CLK_CTL_BCR_BLK_ARES_SHFT));
	else
		clrbits_le32(bcr, BIT(CLK_CTL_BCR_BLK_ARES_SHFT));

	return 0;
}

void clock_configure_qup(int qup, uint32_t hz)
{
	int s = qup % QUP_WRAP0_S7;
	struct sdm845_qupv3_clock *qup_clk = qup < QUP_WRAP1_S0 ?
	(struct sdm845_qupv3_clock *)&gcc->qup_wrap0_s[s] :
	(struct sdm845_qupv3_clock *)&gcc->qup_wrap1_s[s];
	clock_configure(&qup_clk->clk, qup_cfg, hz, ARRAY_SIZE(qup_cfg));
}

void clock_enable_qup(int qup)
{
	int s = qup % QUP_WRAP0_S7;
	int clk_en_off = qup < QUP_WRAP1_S0 ?
		QUPV3_WRAP0_CLK_ENA_S(s) : QUPV3_WRAP1_CLK_ENA_S(s);
	struct sdm845_qupv3_clock *qup_clk = qup < QUP_WRAP1_S0 ?
		&gcc->qup_wrap0_s[s] : &gcc->qup_wrap1_s[s];

	clock_enable_vote(&qup_clk->clk, &gcc->apcs_clk_br_en1,
			  clk_en_off);

}

void clock_init(void)
{

	clock_configure_gpll0();

	clock_enable_vote(&gcc->qup_wrap0_core_2x_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP0_CORE_2X_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap0_core_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP0_CORE_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap0_m_ahb_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP_0_M_AHB_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap0_s_ahb_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP_0_S_AHB_CLK_ENA);

	clock_enable_vote(&gcc->qup_wrap1_core_2x_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP1_CORE_2X_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap1_core_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP1_CORE_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap1_m_ahb_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP_1_M_AHB_CLK_ENA);
	clock_enable_vote(&gcc->qup_wrap1_s_ahb_cbcr,
				&gcc->apcs_clk_br_en1,
				QUPV3_WRAP_1_S_AHB_CLK_ENA);
}
