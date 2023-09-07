/*
 * This file is part of the coreboot project.
 *
 * Copyright 2013 Google Inc.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; version 2 of
 * the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <console/console.h>
#include <device/pci_def.h>
#include <device/pci_ops.h>
#include <security/vboot/vbnv.h>
#include <pc80/mc146818rtc.h>
#include <elog.h>
#include "pmutil.h"
#include "rtc.h"

/* PCI Configuration Space (D31:F0): LPC */
#if defined(__SIMPLE_DEVICE__)
#define PCH_LPC_DEV	PCI_DEV(0, 0x1f, 0)
#else
#define PCH_LPC_DEV	pcidev_on_root(0x1f, 0)
#endif

int rtc_failure(void)
{
	return !!(pci_read_config8(PCH_LPC_DEV, D31F0_GEN_PMCON_3)
		  & RTC_BATTERY_DEAD);
}

void sb_rtc_init(void)
{
	int rtc_failed = rtc_failure();

	if (rtc_failed) {
		if (CONFIG(ELOG))
			elog_add_event(ELOG_TYPE_RTC_RESET);
		pci_update_config8(PCH_LPC_DEV, D31F0_GEN_PMCON_3,
				   ~RTC_BATTERY_DEAD, 0);
	}

	printk(BIOS_DEBUG, "RTC: failed = 0x%x\n", rtc_failed);

	cmos_init(rtc_failed);
}

int vbnv_cmos_failed(void)
{
	return rtc_failure();
}
