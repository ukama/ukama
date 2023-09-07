/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 The ChromiumOS Authors.  All rights reserved.
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

#include <boot_device.h>
#include <cbmem.h>
#include <console/cbmem_console.h>
#include <fmap.h>
#include <reset.h>
#include <stddef.h>
#include <security/vboot/gbb.h>
#include <security/vboot/vboot_common.h>
#include <security/vboot/vbnv.h>
#include <vb2_api.h>

int vboot_named_region_device(const char *name, struct region_device *rdev)
{
	return fmap_locate_area_as_rdev(name, rdev);
}

int vboot_named_region_device_rw(const char *name, struct region_device *rdev)
{
	return fmap_locate_area_as_rdev_rw(name, rdev);
}

/* Check if it is okay to enable USB Device Controller (UDC). */
int vboot_can_enable_udc(void)
{
	/* Always disable if not in developer mode */
	if (!vboot_developer_mode_enabled())
		return 0;
	/* Enable if GBB flag is set */
	if (gbb_is_flag_set(VB2_GBB_FLAG_ENABLE_UDC))
		return 1;
	/* Enable if VBNV flag is set */
	if (vbnv_udc_enable_flag())
		return 1;
	/* Otherwise disable */
	return 0;
}

/* ============================ VBOOT REBOOT ============================== */
void __weak vboot_platform_prepare_reboot(void)
{
}

void vboot_reboot(void)
{
	if (CONFIG(CONSOLE_CBMEM_DUMP_TO_UART))
		cbmem_dump_console();
	vboot_platform_prepare_reboot();
	board_reset();
}
