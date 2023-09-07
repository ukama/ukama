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

#include <device/mmio.h>
#include <console/console.h>
#include <soc/wdt.h>
#include <vendorcode/google/chromeos/chromeos.h>

int mtk_wdt_init(void)
{
	uint32_t wdt_sta;

	/* Write Mode register will clear status register */
	wdt_sta = read32(&mtk_wdt->wdt_status);

	printk(BIOS_INFO, "WDT: Last reset was ");
	if (wdt_sta & MTK_WDT_STA_HW_RST) {
		printk(BIOS_INFO, "hardware watchdog\n");
		mark_watchdog_tombstone();
	} else if (wdt_sta & MTK_WDT_STA_SW_RST)
		printk(BIOS_INFO, "normal software reboot\n");
	else if (wdt_sta & MTK_WDT_STA_SPM_RST)
		printk(BIOS_INFO, "SPM reboot\n");
	else if (!wdt_sta)
		printk(BIOS_INFO, "cold boot\n");
	else
		printk(BIOS_INFO, "unexpected reset type: %#.8x\n", wdt_sta);

	/* Config watchdog reboot mode:
	 * Clearing bits:
	 * DUAL_MODE & IRQ: trigger reset instead of irq then reset.
	 * EXT_POL: select watchdog output signal as active low.
	 * ENABLE: disable watchdog on initialization.
	 * Setting bit EXTEN to enable watchdog output.
	 */
	clrsetbits_le32(&mtk_wdt->wdt_mode,
			MTK_WDT_MODE_DUAL_MODE | MTK_WDT_MODE_IRQ |
			MTK_WDT_MODE_EXT_POL | MTK_WDT_MODE_ENABLE,
			MTK_WDT_MODE_EXTEN | MTK_WDT_MODE_KEY);

	return wdt_sta;
}
