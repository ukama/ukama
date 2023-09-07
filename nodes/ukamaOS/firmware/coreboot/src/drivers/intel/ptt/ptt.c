/*
 * This file is part of the coreboot project.
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

#include <soc/pci_devs.h>
#include <device/pci_ops.h>
#include <console/console.h>
#include <timer.h>

#include "ptt.h"

#define PCI_ME_HFSTS4 0x64
#define PTT_ENABLE (1 << 19)

/* Dump Intel ME register */
static uint32_t read_register(int reg_addr)
{
	if (!PCH_DEV_CSE)
		return 0xFFFFFFFF;

	return pci_read_config32(PCH_DEV_CSE, reg_addr);
}

/*
 * ptt_active()
 *
 * Check if PTT Flag is set - so that PTT is active.
 *
 * Return true if active, false otherwise.
 */
bool ptt_active(void)
{
	uint32_t fwsts4 = read_register(PCI_ME_HFSTS4);

	if (fwsts4 == 0xFFFFFFFF)
		return false;

	if ((fwsts4 & PTT_ENABLE) == 0) {
		printk(BIOS_DEBUG, "Intel ME Establishment bit not valid.\n");
		return false;
	}

	return true;
}
