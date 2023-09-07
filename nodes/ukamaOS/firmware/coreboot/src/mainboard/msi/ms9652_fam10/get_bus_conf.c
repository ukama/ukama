/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007 AMD
 * Written by Yinghai Lu <yinghailu@amd.com> for AMD.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <console/console.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <string.h>
#include <stdint.h>
#include <cpu/amd/multicore.h>

#include <cpu/amd/amdfam10_sysconf.h>

#include <stdlib.h>
#include "mb_sysconf.h"

// Global variables for MB layouts and these will be shared by irqtable mptable and acpi_tables
struct mb_sysconf_t mb_sysconf;

void get_bus_conf(void)
{
	unsigned int apicid_base;
	struct mb_sysconf_t *m;

	struct device *dev;
	int i;


	sysconf.mb = &mb_sysconf;

	m = sysconf.mb;
	memset(m, 0, sizeof(struct mb_sysconf_t));

	get_default_pci1234(32);

	sysconf.sbdn = (sysconf.hcdn[0] & 0xff); // first byte of first chain
	m->bus_mcp55[0] = (sysconf.pci1234[0] >> 12) & 0xff;

		/* MCP55 */
		dev = dev_find_slot(m->bus_mcp55[0], PCI_DEVFN(sysconf.sbdn + 0x06,0));
		if (dev) {
			m->bus_mcp55[1] = pci_read_config8(dev, PCI_SECONDARY_BUS);
		}
		else {
			printk(BIOS_DEBUG, "ERROR - could not find PCI 1:%02x.0, using defaults\n", sysconf.sbdn + 0x06);
		}

		for (i = 2; i < 8; i++) {
			dev = dev_find_slot(m->bus_mcp55[0], PCI_DEVFN(sysconf.sbdn + 0x0a + i - 2, 0));
			if (dev) {
				m->bus_mcp55[i] = pci_read_config8(dev, PCI_SECONDARY_BUS);
			}
			else {
				printk(BIOS_DEBUG, "ERROR - could not find PCI %02x:%02x.0, using defaults\n", m->bus_mcp55[0], sysconf.sbdn + 0x0a + i - 2);
			}
		}

/*I/O APICs:	APIC ID	Version	State		Address*/
	if (CONFIG(LOGICAL_CPUS)) {
		apicid_base = get_apicid_base(1);
		printk(BIOS_SPEW, "CONFIG_LOGICAL_CPUS == 1: apicid_base: %08x\n", apicid_base);
	} else {
		apicid_base = CONFIG_MAX_PHYSICAL_CPUS;
		printk(BIOS_SPEW, "CONFIG_LOGICAL_CPUS == 0: apicid_base: %08x\n", apicid_base);
	}
	m->apicid_mcp55 = apicid_base+0;
}
