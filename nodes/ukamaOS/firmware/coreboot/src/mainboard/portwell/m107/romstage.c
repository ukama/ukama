/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2013 Google Inc.
 * Copyright (C) 2015 Intel Corp.
 * Copyright (C) 2018-2019 Eltan B.V.
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

#include <cbfs.h>
#include <console/console.h>
#include <chip.h>
#include <device/pci_ops.h>
#include <soc/lpc.h>
#include <soc/pci_devs.h>
#include <soc/romstage.h>
#include <spd_bin.h>
#include <stdint.h>

void mainboard_memory_init_params(struct romstage_params *params,
	MEMORY_INIT_UPD *memory_params)
{
	struct region_device spd_rdev;
	u8 spd_index = 0;

	if (CONFIG(ONBOARD_MEM_MICRON))
		spd_index = 1;
	else if (CONFIG(ONBOARD_MEM_KINGSTON))
		spd_index = 2;

	if (get_spd_cbfs_rdev(&spd_rdev, spd_index) < 0)
		die("spd.bin not found\n");

	memory_params->PcdMemoryTypeEnable = MEM_DDR3;
	memory_params->PcdMemorySpdPtr = (uintptr_t)rdev_mmap_full(&spd_rdev);
	memory_params->PcdMemChannel0Config = 1; /* Memory down */
	memory_params->PcdMemChannel1Config = 2; /* Disabled */
}

void mainboard_after_memory_init(void)
{
	printk(BIOS_DEBUG, "%s/%s called\n", __FILE__, __func__);

	/* Disable the Braswell UART hardware for COM1. */
	pci_write_config32(PCI_DEV(0, LPC_DEV, 0), UART_CONT, 0);
}
