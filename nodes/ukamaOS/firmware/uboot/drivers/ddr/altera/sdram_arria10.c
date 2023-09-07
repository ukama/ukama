// SPDX-License-Identifier: GPL-2.0
/*
 * Copyright (C) 2017 Intel Corporation <www.intel.com>
 */

#include <common.h>
#include <errno.h>
#include <fdtdec.h>
#include <malloc.h>
#include <wait_bit.h>
#include <watchdog.h>
#include <asm/io.h>
#include <asm/arch/fpga_manager.h>
#include <asm/arch/misc.h>
#include <asm/arch/reset_manager.h>
#include <asm/arch/sdram.h>
#include <linux/kernel.h>

DECLARE_GLOBAL_DATA_PTR;

static void sdram_mmr_init(void);
static u64 sdram_size_calc(void);

/* FAWBANK - Number of Bank of a given device involved in the FAW period. */
#define ARRIA10_SDR_ACTIVATE_FAWBANK	(0x1)

#define ARRIA_DDR_CONFIG(A, B, C, R) \
	(((A) << 24) | ((B) << 16) | ((C) << 8) | (R))
#define DDR_CONFIG_ELEMENTS	ARRAY_SIZE(ddr_config)
#define DDR_REG_SEQ2CORE        0xFFD0507C
#define DDR_REG_CORE2SEQ        0xFFD05078
#define DDR_READ_LATENCY_DELAY	40
#define DDR_SIZE_2GB_HEX	0x80000000
#define DDR_MAX_TRIES		0x00100000

#define IO48_MMR_DRAMSTS	0xFFCFA0EC
#define IO48_MMR_NIOS2_RESERVE0	0xFFCFA110
#define IO48_MMR_NIOS2_RESERVE1	0xFFCFA114
#define IO48_MMR_NIOS2_RESERVE2	0xFFCFA118

#define SEQ2CORE_MASK		0xF
#define CORE2SEQ_INT_REQ	0xF
#define SEQ2CORE_INT_RESP_BIT	3

static const struct socfpga_ecc_hmc *socfpga_ecc_hmc_base =
		(void *)SOCFPGA_SDR_ADDRESS;
static const struct socfpga_noc_ddr_scheduler *socfpga_noc_ddr_scheduler_base =
		(void *)SOCFPGA_SDR_SCHEDULER_ADDRESS;
static const struct socfpga_noc_fw_ddr_mpu_fpga2sdram
		*socfpga_noc_fw_ddr_mpu_fpga2sdram_base =
		(void *)SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS;
static const struct socfpga_noc_fw_ddr_l3 *socfpga_noc_fw_ddr_l3_base =
		(void *)SOCFPGA_SDR_FIREWALL_L3_ADDRESS;
static const struct socfpga_io48_mmr *socfpga_io48_mmr_base =
		(void *)SOCFPGA_HMC_MMR_IO48_ADDRESS;

/* The following are the supported configurations */
static u32 ddr_config[] = {
	/* Chip - Row - Bank - Column Style */
	/* All Types */
	ARRIA_DDR_CONFIG(0, 3, 10, 12),
	ARRIA_DDR_CONFIG(0, 3, 10, 13),
	ARRIA_DDR_CONFIG(0, 3, 10, 14),
	ARRIA_DDR_CONFIG(0, 3, 10, 15),
	ARRIA_DDR_CONFIG(0, 3, 10, 16),
	ARRIA_DDR_CONFIG(0, 3, 10, 17),
	/* LPDDR x16 */
	ARRIA_DDR_CONFIG(0, 3, 11, 14),
	ARRIA_DDR_CONFIG(0, 3, 11, 15),
	ARRIA_DDR_CONFIG(0, 3, 11, 16),
	ARRIA_DDR_CONFIG(0, 3, 12, 15),
	/* DDR4 Only */
	ARRIA_DDR_CONFIG(0, 4, 10, 14),
	ARRIA_DDR_CONFIG(0, 4, 10, 15),
	ARRIA_DDR_CONFIG(0, 4, 10, 16),
	ARRIA_DDR_CONFIG(0, 4, 10, 17),	/* 14 */
	/* Chip - Bank - Row - Column Style */
	ARRIA_DDR_CONFIG(1, 3, 10, 12),
	ARRIA_DDR_CONFIG(1, 3, 10, 13),
	ARRIA_DDR_CONFIG(1, 3, 10, 14),
	ARRIA_DDR_CONFIG(1, 3, 10, 15),
	ARRIA_DDR_CONFIG(1, 3, 10, 16),
	ARRIA_DDR_CONFIG(1, 3, 10, 17),
	ARRIA_DDR_CONFIG(1, 3, 11, 14),
	ARRIA_DDR_CONFIG(1, 3, 11, 15),
	ARRIA_DDR_CONFIG(1, 3, 11, 16),
	ARRIA_DDR_CONFIG(1, 3, 12, 15),
	/* DDR4 Only */
	ARRIA_DDR_CONFIG(1, 4, 10, 14),
	ARRIA_DDR_CONFIG(1, 4, 10, 15),
	ARRIA_DDR_CONFIG(1, 4, 10, 16),
	ARRIA_DDR_CONFIG(1, 4, 10, 17),
};

static int match_ddr_conf(u32 ddr_conf)
{
	int i;

	for (i = 0; i < DDR_CONFIG_ELEMENTS; i++) {
		if (ddr_conf == ddr_config[i])
			return i;
	}
	return 0;
}

/* Check whether SDRAM is successfully Calibrated */
static int is_sdram_cal_success(void)
{
	return readl(&socfpga_ecc_hmc_base->ddrcalstat);
}

static unsigned char ddr_get_bit(u32 ereg, unsigned char bit)
{
	u32 reg = readl(ereg);

	return (reg & BIT(bit)) ? 1 : 0;
}

static unsigned char ddr_wait_bit(u32 ereg, u32 bit,
			   u32 expected, u32 timeout_usec)
{
	u32 tmr;

	for (tmr = 0; tmr < timeout_usec; tmr += 100) {
		udelay(100);
		WATCHDOG_RESET();
		if (ddr_get_bit(ereg, bit) == expected)
			return 0;
	}

	return 1;
}

static int emif_clear(void)
{
	u32 i = DDR_MAX_TRIES;
	u8 ret = 0;

	writel(0, DDR_REG_CORE2SEQ);

	do {
		ret = !wait_for_bit_le32((u32 *)DDR_REG_SEQ2CORE,
				   SEQ2CORE_MASK, 1, 50, 0);
	} while (ret && (--i > 0));

	return !i;
}

static int emif_reset(void)
{
	u32 c2s, s2c;

	c2s = readl(DDR_REG_CORE2SEQ);
	s2c = readl(DDR_REG_SEQ2CORE);

	debug("c2s=%08x s2c=%08x nr0=%08x nr1=%08x nr2=%08x dst=%08x\n",
	     c2s, s2c, readl(IO48_MMR_NIOS2_RESERVE0),
	     readl(IO48_MMR_NIOS2_RESERVE1),
	     readl(IO48_MMR_NIOS2_RESERVE2),
	     readl(IO48_MMR_DRAMSTS));

	if ((s2c & SEQ2CORE_MASK) && emif_clear()) {
		debug("failed emif_clear()\n");
		return -EPERM;
	}

	writel(CORE2SEQ_INT_REQ, DDR_REG_CORE2SEQ);

	if (ddr_wait_bit(DDR_REG_SEQ2CORE, SEQ2CORE_INT_RESP_BIT, 0, 1000000)) {
		debug("emif_reset failed to see interrupt acknowledge\n");
		return -EPERM;
	} else {
		debug("emif_reset interrupt acknowledged\n");
	}

	if (emif_clear()) {
		debug("emif_clear() failed\n");
		return -EPERM;
	}
	debug("emif_reset interrupt cleared\n");

	debug("nr0=%08x nr1=%08x nr2=%08x\n",
	     readl(IO48_MMR_NIOS2_RESERVE0),
	     readl(IO48_MMR_NIOS2_RESERVE1),
	     readl(IO48_MMR_NIOS2_RESERVE2));

	return 0;
}

static int ddr_setup(void)
{
	int i, j, ddr_setup_complete = 0;

	/* Try 3 times to do a calibration */
	for (i = 0; (i < 3) && !ddr_setup_complete; i++) {
		WATCHDOG_RESET();

		/* A delay to wait for calibration bit to set */
		for (j = 0; (j < 10) && !ddr_setup_complete; j++) {
			mdelay(500);
			ddr_setup_complete = is_sdram_cal_success();
		}

		if (!ddr_setup_complete)
			if (emif_reset())
				puts("Error: Failed to reset EMIF\n");
	}

	/* After 3 times trying calibration */
	if (!ddr_setup_complete) {
		puts("Error: Could Not Calibrate SDRAM\n");
		return -EPERM;
	}

	return 0;
}

/* Function to startup the SDRAM*/
static int sdram_startup(void)
{
	/* Release NOC ddr scheduler from reset */
	socfpga_reset_deassert_noc_ddr_scheduler();

	/* Bringup the DDR (calibration and configuration) */
	return ddr_setup();
}

static u64 sdram_size_calc(void)
{
	u32 dramaddrw = readl(&socfpga_io48_mmr_base->dramaddrw);

	u64 size = BIT(((dramaddrw &
		IO48_MMR_DRAMADDRW_CFG_CS_ADDR_WIDTH_MASK) >>
		IO48_MMR_DRAMADDRW_CFG_CS_ADDR_WIDTH_SHIFT) +
		((dramaddrw &
		IO48_MMR_DRAMADDRW_CFG_BANK_GROUP_ADDR_WIDTH_MASK) >>
		IO48_MMR_DRAMADDRW_CFG_BANK_GROUP_ADDR_WIDTH_SHIFT) +
		((dramaddrw &
		IO48_MMR_DRAMADDRW_CFG_BANK_ADDR_WIDTH_MASK) >>
		IO48_MMR_DRAMADDRW_CFG_BANK_ADDR_WIDTH_SHIFT) +
		((dramaddrw &
		IO48_MMR_DRAMADDRW_CFG_ROW_ADDR_WIDTH_MASK) >>
		IO48_MMR_DRAMADDRW_CFG_ROW_ADDR_WIDTH_SHIFT) +
		(dramaddrw & IO48_MMR_DRAMADDRW_CFG_COL_ADDR_WIDTH_MASK));

	size *= (2 << (readl(&socfpga_ecc_hmc_base->ddrioctrl) &
		       ALT_ECC_HMC_OCP_DDRIOCTRL_IO_SIZE_MSK));

	debug("SDRAM size=%llu", size);

	return size;
}

/* Function to initialize SDRAM MMR and NOC DDR scheduler*/
static void sdram_mmr_init(void)
{
	u32 update_value, io48_value;
	u32 ctrlcfg0 = readl(&socfpga_io48_mmr_base->ctrlcfg0);
	u32 ctrlcfg1 = readl(&socfpga_io48_mmr_base->ctrlcfg1);
	u32 dramaddrw = readl(&socfpga_io48_mmr_base->dramaddrw);
	u32 caltim0 = readl(&socfpga_io48_mmr_base->caltiming0);
	u32 caltim1 = readl(&socfpga_io48_mmr_base->caltiming1);
	u32 caltim2 = readl(&socfpga_io48_mmr_base->caltiming2);
	u32 caltim3 = readl(&socfpga_io48_mmr_base->caltiming3);
	u32 caltim4 = readl(&socfpga_io48_mmr_base->caltiming4);
	u32 caltim9 = readl(&socfpga_io48_mmr_base->caltiming9);
	u32 ddrioctl;

	/*
	 * Configure the DDR IO size [0xFFCFB008]
	 * niosreserve0: Used to indicate DDR width &
	 *	bit[7:0] = Number of data bits (0x20 for 32bit)
	 *	bit[8]   = 1 if user-mode OCT is present
	 *	bit[9]   = 1 if warm reset compiled into EMIF Cal Code
	 *	bit[10]  = 1 if warm reset is on during generation in EMIF Cal
	 * niosreserve1: IP ADCDS version encoded as 16 bit value
	 *	bit[2:0] = Variant (0=not special,1=FAE beta, 2=Customer beta,
	 *			    3=EAP, 4-6 are reserved)
	 *	bit[5:3] = Service Pack # (e.g. 1)
	 *	bit[9:6] = Minor Release #
	 *	bit[14:10] = Major Release #
	 */
	if ((socfpga_io48_mmr_base->niosreserve1 >> 6) & 0x1FF) {
		update_value = readl(&socfpga_io48_mmr_base->niosreserve0);
		writel(((update_value & 0xFF) >> 5),
		       &socfpga_ecc_hmc_base->ddrioctrl);
	}

	ddrioctl = readl(&socfpga_ecc_hmc_base->ddrioctrl);

	/* Set the DDR Configuration [0xFFD12400] */
	io48_value = ARRIA_DDR_CONFIG(
			((ctrlcfg1 &
			IO48_MMR_CTRLCFG1_ADDR_ORDER_MASK) >>
			IO48_MMR_CTRLCFG1_ADDR_ORDER_SHIFT),
			((dramaddrw &
			IO48_MMR_DRAMADDRW_CFG_BANK_ADDR_WIDTH_MASK) >>
			IO48_MMR_DRAMADDRW_CFG_BANK_ADDR_WIDTH_SHIFT) +
			((dramaddrw &
			IO48_MMR_DRAMADDRW_CFG_BANK_GROUP_ADDR_WIDTH_MASK) >>
			IO48_MMR_DRAMADDRW_CFG_BANK_GROUP_ADDR_WIDTH_SHIFT),
			(dramaddrw &
			IO48_MMR_DRAMADDRW_CFG_COL_ADDR_WIDTH_MASK),
			((dramaddrw &
			IO48_MMR_DRAMADDRW_CFG_ROW_ADDR_WIDTH_MASK) >>
			IO48_MMR_DRAMADDRW_CFG_ROW_ADDR_WIDTH_SHIFT));

	update_value = match_ddr_conf(io48_value);
	if (update_value)
		writel(update_value,
		&socfpga_noc_ddr_scheduler_base->ddr_t_main_scheduler_ddrconf);

	/*
	 * Configure DDR timing [0xFFD1240C]
	 *  RDTOMISS = tRTP + tRP + tRCD - BL/2
	 *  WRTOMISS = WL + tWR + tRP + tRCD and
	 *    WL = RL + BL/2 + 2 - rd-to-wr ; tWR = 15ns  so...
	 *  First part of equation is in memory clock units so divide by 2
	 *  for HMC clock units. 1066MHz is close to 1ns so use 15 directly.
	 *  WRTOMISS = ((RL + BL/2 + 2 + tWR) >> 1)- rd-to-wr + tRP + tRCD
	 */
	u32 ctrlcfg0_cfg_ctrl_burst_len =
		(ctrlcfg0 & IO48_MMR_CTRLCFG0_CTRL_BURST_LENGTH_MASK) >>
		IO48_MMR_CTRLCFG0_CTRL_BURST_LENGTH_SHIFT;

	u32 caltim0_cfg_act_to_rdwr = caltim0 &
		IO48_MMR_CALTIMING0_CFG_ACT_TO_RDWR_MASK;

	u32 caltim0_cfg_act_to_act =
		(caltim0 & IO48_MMR_CALTIMING0_CFG_ACT_TO_ACT_MASK) >>
		IO48_MMR_CALTIMING0_CFG_ACT_TO_ACT_SHIFT;

	u32 caltim0_cfg_act_to_act_db =
		(caltim0 &
		IO48_MMR_CALTIMING0_CFG_ACT_TO_ACT_DIFF_BANK_MASK) >>
		IO48_MMR_CALTIMING0_CFG_ACT_TO_ACT_DIFF_BANK_SHIFT;

	u32 caltim1_cfg_rd_to_wr =
		(caltim1 & IO48_MMR_CALTIMING1_CFG_RD_TO_WR_MASK) >>
		IO48_MMR_CALTIMING1_CFG_RD_TO_WR_SHIFT;

	u32 caltim1_cfg_rd_to_rd_dc =
		(caltim1 & IO48_MMR_CALTIMING1_CFG_RD_TO_RD_DC_MASK) >>
		IO48_MMR_CALTIMING1_CFG_RD_TO_RD_DC_SHIFT;

	u32 caltim1_cfg_rd_to_wr_dc =
		(caltim1 & IO48_MMR_CALTIMING1_CFG_RD_TO_WR_DIFF_CHIP_MASK) >>
		IO48_MMR_CALTIMING1_CFG_RD_TO_WR_DIFF_CHIP_SHIFT;

	u32 caltim2_cfg_rd_to_pch =
		(caltim2 & IO48_MMR_CALTIMING2_CFG_RD_TO_PCH_MASK) >>
		IO48_MMR_CALTIMING2_CFG_RD_TO_PCH_SHIFT;

	u32 caltim3_cfg_wr_to_rd =
		(caltim3 & IO48_MMR_CALTIMING3_CFG_WR_TO_RD_MASK) >>
		IO48_MMR_CALTIMING3_CFG_WR_TO_RD_SHIFT;

	u32 caltim3_cfg_wr_to_rd_dc =
		(caltim3 & IO48_MMR_CALTIMING3_CFG_WR_TO_RD_DIFF_CHIP_MASK) >>
		IO48_MMR_CALTIMING3_CFG_WR_TO_RD_DIFF_CHIP_SHIFT;

	u32 caltim4_cfg_pch_to_valid =
		(caltim4 & IO48_MMR_CALTIMING4_CFG_PCH_TO_VALID_MASK) >>
		IO48_MMR_CALTIMING4_CFG_PCH_TO_VALID_SHIFT;

	u32 caltim9_cfg_4_act_to_act = caltim9 &
		IO48_MMR_CALTIMING9_CFG_WR_4_ACT_TO_ACT_MASK;

	update_value = (caltim2_cfg_rd_to_pch +  caltim4_cfg_pch_to_valid +
			caltim0_cfg_act_to_rdwr -
			(ctrlcfg0_cfg_ctrl_burst_len >> 2));

	io48_value = ((((socfpga_io48_mmr_base->dramtiming0 &
		      ALT_IO48_DRAMTIME_MEM_READ_LATENCY_MASK) + 2 + 15 +
		      (ctrlcfg0_cfg_ctrl_burst_len >> 1)) >> 1) -
		      /* Up to here was in memory cycles so divide by 2 */
		      caltim1_cfg_rd_to_wr + caltim0_cfg_act_to_rdwr +
		      caltim4_cfg_pch_to_valid);

	writel(((caltim0_cfg_act_to_act <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_ACTTOACT_LSB) |
		(update_value <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_RDTOMISS_LSB) |
		(io48_value <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_WRTOMISS_LSB) |
		((ctrlcfg0_cfg_ctrl_burst_len >> 2) <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_BURSTLEN_LSB) |
		(caltim1_cfg_rd_to_wr <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_RDTOWR_LSB) |
		(caltim3_cfg_wr_to_rd <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_WRTORD_LSB) |
		(((ddrioctl == 1) ? 1 : 0) <<
			ALT_NOC_MPU_DDR_T_SCHED_DDRTIMING_BWRATIO_LSB)),
		&socfpga_noc_ddr_scheduler_base->
			ddr_t_main_scheduler_ddrtiming);

	/* Configure DDR mode [0xFFD12410] [precharge = 0] */
	writel(((ddrioctl ? 0 : 1) <<
		ALT_NOC_MPU_DDR_T_SCHED_DDRMOD_BWRATIOEXTENDED_LSB),
		&socfpga_noc_ddr_scheduler_base->ddr_t_main_scheduler_ddrmode);

	/* Configure the read latency [0xFFD12414] */
	writel(((socfpga_io48_mmr_base->dramtiming0 &
		ALT_IO48_DRAMTIME_MEM_READ_LATENCY_MASK) >> 1) +
		DDR_READ_LATENCY_DELAY,
		&socfpga_noc_ddr_scheduler_base->
			ddr_t_main_scheduler_readlatency);

	/*
	 * Configuring timing values concerning activate commands
	 * [0xFFD12438] [FAWBANK alway 1 because always 4 bank DDR]
	 */
	writel(((caltim0_cfg_act_to_act_db <<
			ALT_NOC_MPU_DDR_T_SCHED_ACTIVATE_RRD_LSB) |
		(caltim9_cfg_4_act_to_act <<
			ALT_NOC_MPU_DDR_T_SCHED_ACTIVATE_FAW_LSB) |
		(ARRIA10_SDR_ACTIVATE_FAWBANK <<
			ALT_NOC_MPU_DDR_T_SCHED_ACTIVATE_FAWBANK_LSB)),
		&socfpga_noc_ddr_scheduler_base->ddr_t_main_scheduler_activate);

	/*
	 * Configuring timing values concerning device to device data bus
	 * ownership change [0xFFD1243C]
	 */
	writel(((caltim1_cfg_rd_to_rd_dc <<
			ALT_NOC_MPU_DDR_T_SCHED_DEVTODEV_BUSRDTORD_LSB) |
		(caltim1_cfg_rd_to_wr_dc <<
			ALT_NOC_MPU_DDR_T_SCHED_DEVTODEV_BUSRDTOWR_LSB) |
		(caltim3_cfg_wr_to_rd_dc <<
			ALT_NOC_MPU_DDR_T_SCHED_DEVTODEV_BUSWRTORD_LSB)),
		&socfpga_noc_ddr_scheduler_base->ddr_t_main_scheduler_devtodev);

	/* Enable or disable the SDRAM ECC */
	if (ctrlcfg1 & IO48_MMR_CTRLCFG1_CTRL_ENABLE_ECC) {
		setbits_le32(&socfpga_ecc_hmc_base->eccctrl,
			     (ALT_ECC_HMC_OCP_ECCCTL_AWB_CNT_RST_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL_CNT_RST_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL_ECC_EN_SET_MSK));
		clrbits_le32(&socfpga_ecc_hmc_base->eccctrl,
			     (ALT_ECC_HMC_OCP_ECCCTL_AWB_CNT_RST_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL_CNT_RST_SET_MSK));
		setbits_le32(&socfpga_ecc_hmc_base->eccctrl2,
			     (ALT_ECC_HMC_OCP_ECCCTL2_RMW_EN_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL2_AWB_EN_SET_MSK));
	} else {
		clrbits_le32(&socfpga_ecc_hmc_base->eccctrl,
			     (ALT_ECC_HMC_OCP_ECCCTL_AWB_CNT_RST_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL_CNT_RST_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL_ECC_EN_SET_MSK));
		clrbits_le32(&socfpga_ecc_hmc_base->eccctrl2,
			     (ALT_ECC_HMC_OCP_ECCCTL2_RMW_EN_SET_MSK |
			      ALT_ECC_HMC_OCP_ECCCTL2_AWB_EN_SET_MSK));
	}
}

struct firewall_entry {
	const char *prop_name;
	const u32 cfg_addr;
	const u32 en_addr;
	const u32 en_bit;
};
#define FW_MPU_FPGA_ADDRESS \
	((const struct socfpga_noc_fw_ddr_mpu_fpga2sdram *)\
	SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS)

#define SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(ADDR) \
		(SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS + \
		offsetof(struct socfpga_noc_fw_ddr_mpu_fpga2sdram, ADDR))

#define SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(ADDR) \
		(SOCFPGA_SDR_FIREWALL_L3_ADDRESS + \
		offsetof(struct socfpga_noc_fw_ddr_l3, ADDR))

const struct firewall_entry firewall_table[] = {
	{
		"mpu0",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(mpuregion0addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_MPUREG0EN_SET_MSK
	},
	{
		"mpu1",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS +
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(mpuregion1addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_MPUREG1EN_SET_MSK
	},
	{
		"mpu2",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(mpuregion2addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_MPUREG2EN_SET_MSK
	},
	{
		"mpu3",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(mpuregion3addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_MPUREG3EN_SET_MSK
	},
	{
		"l3-0",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion0addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG0EN_SET_MSK
	},
	{
		"l3-1",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion1addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG1EN_SET_MSK
	},
	{
		"l3-2",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion2addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG2EN_SET_MSK
	},
	{
		"l3-3",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion3addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG3EN_SET_MSK
	},
	{
		"l3-4",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion4addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG4EN_SET_MSK
	},
	{
		"l3-5",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion5addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG5EN_SET_MSK
	},
	{
		"l3-6",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion6addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG6EN_SET_MSK
	},
	{
		"l3-7",
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(hpsregion7addr),
		SOCFPGA_SDR_FIREWALL_L3_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_HPSREG7EN_SET_MSK
	},
	{
		"fpga2sdram0-0",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram0region0addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR0REG0EN_SET_MSK
	},
	{
		"fpga2sdram0-1",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram0region1addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR0REG1EN_SET_MSK
	},
	{
		"fpga2sdram0-2",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram0region2addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR0REG2EN_SET_MSK
	},
	{
		"fpga2sdram0-3",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram0region3addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR0REG3EN_SET_MSK
	},
	{
		"fpga2sdram1-0",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram1region0addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR1REG0EN_SET_MSK
	},
	{
		"fpga2sdram1-1",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram1region1addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR1REG1EN_SET_MSK
	},
	{
		"fpga2sdram1-2",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram1region2addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR1REG2EN_SET_MSK
	},
	{
		"fpga2sdram1-3",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram1region3addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR1REG3EN_SET_MSK
	},
	{
		"fpga2sdram2-0",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram2region0addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR2REG0EN_SET_MSK
	},
	{
		"fpga2sdram2-1",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram2region1addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR2REG1EN_SET_MSK
	},
	{
		"fpga2sdram2-2",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram2region2addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR2REG2EN_SET_MSK
	},
	{
		"fpga2sdram2-3",
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET
		(fpga2sdram2region3addr),
		SOCFPGA_SDR_FIREWALL_MPU_FPGA_ADDRESS_OFFSET(enable),
		ALT_NOC_FW_DDR_SCR_EN_F2SDR2REG3EN_SET_MSK
	},

};

static int of_sdram_firewall_setup(const void *blob)
{
	int child, i, node, ret;
	u32 start_end[2];
	char name[32];

	node = fdtdec_next_compatible(blob, 0, COMPAT_ALTERA_SOCFPGA_NOC);
	if (node < 0)
		return -ENXIO;

	child = fdt_first_subnode(blob, node);
	if (child < 0)
		return -ENXIO;

	/* set to default state */
	writel(0, &socfpga_noc_fw_ddr_mpu_fpga2sdram_base->enable);
	writel(0, &socfpga_noc_fw_ddr_l3_base->enable);


	for (i = 0; i < ARRAY_SIZE(firewall_table); i++) {
		sprintf(name, "%s", firewall_table[i].prop_name);
		ret = fdtdec_get_int_array(blob, child, name,
					   start_end, 2);
		if (ret) {
			sprintf(name, "altr,%s", firewall_table[i].prop_name);
			ret = fdtdec_get_int_array(blob, child, name,
						   start_end, 2);
			if (ret)
				continue;
		}

		writel((start_end[0] & ALT_NOC_FW_DDR_ADDR_MASK) |
		       (start_end[1] << ALT_NOC_FW_DDR_END_ADDR_LSB),
		       firewall_table[i].cfg_addr);
		setbits_le32(firewall_table[i].en_addr,
			     firewall_table[i].en_bit);
	}

	return 0;
}

int ddr_calibration_sequence(void)
{
	WATCHDOG_RESET();

	/* Check to see if SDRAM cal was success */
	if (sdram_startup()) {
		puts("DDRCAL: Failed\n");
		return -EPERM;
	}

	puts("DDRCAL: Success\n");

	WATCHDOG_RESET();

	/* initialize the MMR register */
	sdram_mmr_init();

	/* assigning the SDRAM size */
	u64 size = sdram_size_calc();

	/*
	 * If size is less than zero, this is invalid/weird value from
	 * calculation, use default Config size.
	 * Up to 2GB is supported, 2GB would be used if more than that.
	 */
	if (size <= 0)
		gd->ram_size = PHYS_SDRAM_1_SIZE;
	else if (DDR_SIZE_2GB_HEX <= size)
		gd->ram_size = DDR_SIZE_2GB_HEX;
	else
		gd->ram_size = (u32)size;

	/* setup the dram info within bd */
	dram_init_banksize();

	if (of_sdram_firewall_setup(gd->fdt_blob))
		puts("FW: Error Configuring Firewall\n");

	return 0;
}

void dram_bank_mmu_setup(int bank)
{
	bd_t *bd = gd->bd;
	int	i;

	debug("%s: bank: %d\n", __func__, bank);
	for (i = bd->bi_dram[bank].start >> 20;
	     i < (bd->bi_dram[bank].start + bd->bi_dram[bank].size) >> 20;
	     i++) {
#if defined(CONFIG_SYS_ARM_CACHE_WRITETHROUGH)
		set_section_dcache(i, DCACHE_WRITETHROUGH);
#else
		set_section_dcache(i, DCACHE_WRITEBACK);
#endif
	}

	/* same as above but just that we would want cacheable for ocram too */
	i = CONFIG_SYS_INIT_RAM_ADDR >> 20;
#if defined(CONFIG_SYS_ARM_CACHE_WRITETHROUGH)
	set_section_dcache(i, DCACHE_WRITETHROUGH);
#else
	set_section_dcache(i, DCACHE_WRITEBACK);
#endif
}
