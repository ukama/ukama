/*
 * This file is part of the coreboot project.
 *
 * Copyright 2015 Google Inc.
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

#include <device/mmio.h>
#include <assert.h>
#include <bootstate.h>
#include <console/console.h>
#include <ec/google/chromeec/ec.h>
#include <elog.h>
#include <reset.h>
#include <symbols.h>

#include "chromeos.h"
#include "symbols.h"

#define WATCHDOG_TOMBSTONE_MAGIC	0x9d2f41a7

DECLARE_OPTIONAL_REGION(watchdog_tombstone);

static void elog_handle_watchdog_tombstone(void *unused)
{
	bool flag = false;

	if (CONFIG(CHROMEOS_USE_EC_WATCHDOG_FLAG))
		flag |= google_chromeec_get_ap_watchdog_flag();

	if (REGION_SIZE(watchdog_tombstone)) {
		flag |= (read32(_watchdog_tombstone) ==
			 WATCHDOG_TOMBSTONE_MAGIC);
		write32(_watchdog_tombstone, 0);
	}

	if (flag)
		elog_add_event(ELOG_TYPE_ASYNC_HW_TIMER_EXPIRED);
}

BOOT_STATE_INIT_ENTRY(BS_POST_DEVICE, BS_ON_ENTRY,
		      elog_handle_watchdog_tombstone, NULL);

void mark_watchdog_tombstone(void)
{
	assert(REGION_SIZE(watchdog_tombstone));
	write32(_watchdog_tombstone, WATCHDOG_TOMBSTONE_MAGIC);
}

void reboot_from_watchdog(void)
{
	printk(BIOS_INFO, "Last reset was watchdog, reboot again to reset TPM!\n");
	mark_watchdog_tombstone();
	board_reset();
}
