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

#include <arch/stages.h>
#include <console/console.h>
#include <fmap.h>
#include <soc/dramc_param.h>
#include <soc/emi.h>
#include <soc/mmu_operations.h>
#include <soc/mt6358.h>
#include <soc/pll.h>
#include <soc/rtc.h>

#include "early_init.h"

/* This must be defined in chromeos.fmd in same name and size. */
#define CALIBRATION_REGION		"RW_DDR_TRAINING"
#define CALIBRATION_REGION_SIZE		0x2000

_Static_assert(sizeof(struct dramc_param) <= CALIBRATION_REGION_SIZE,
	       "sizeof(struct dramc_param) exceeds " CALIBRATION_REGION);

static bool read_calibration_data_from_flash(struct dramc_param *dparam)
{
	const size_t length = sizeof(*dparam);
	size_t ret = fmap_read_area(CALIBRATION_REGION, dparam, length);
	printk(BIOS_DEBUG, "%s: ret=%#lx, length=%#lx\n",
	       __func__, ret, length);

	return ret == length;
}

static bool write_calibration_data_to_flash(const struct dramc_param *dparam)
{
	const size_t length = sizeof(*dparam);
	size_t ret = fmap_overwrite_area(CALIBRATION_REGION, dparam, length);
	printk(BIOS_DEBUG, "%s: ret=%#lx, length=%#lx\n",
	       __func__, ret, length);

	return ret == length;
}

/* dramc_param is ~2K and too large to fit in stack. */
static struct dramc_param dramc_parameter;

static struct dramc_param_ops dparam_ops = {
	.param = &dramc_parameter,
	.read_from_flash = &read_calibration_data_from_flash,
	.write_to_flash = &write_calibration_data_to_flash,
};

void platform_romstage_main(void)
{
	/* This will be done in verstage if CONFIG_VBOOT is enabled. */
	if (!CONFIG(VBOOT))
		mainboard_early_init();

	mt6358_init();
	/* Adjust VSIM2 down to 2.7V because it is shared with IT6505. */
	pmic_set_vsim2_cali(2700);
	mt_pll_raise_ca53_freq(1989 * MHz);
	pmic_init_scp_voltage();
	rtc_boot();
	mt_mem_init(&dparam_ops);
	mtk_mmu_after_dram();
}
