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

#include <device/device.h>
#include <soc/emi.h>
#include <soc/md_ctrl.h>
#include <soc/mmu_operations.h>
#include <soc/rtc.h>
#include <soc/sspm.h>
#include <symbols.h>

static void soc_read_resources(struct device *dev)
{
	ram_resource(dev, 0, (uintptr_t)_dram / KiB, sdram_size() / KiB);
}

static void soc_init(struct device *dev)
{
	mtk_mmu_disable_l2c_sram();
	mtk_md_early_init();
	mt6358_dcxo_disable_unused();
	sspm_init();
}

static struct device_operations soc_ops = {
	.read_resources = soc_read_resources,
	.init = soc_init,
};

static void enable_soc_dev(struct device *dev)
{
	dev->ops = &soc_ops;
}

struct chip_operations soc_mediatek_mt8183_ops = {
	CHIP_NAME("SOC Mediatek MT8183")
	.enable_dev = enable_soc_dev,
};
