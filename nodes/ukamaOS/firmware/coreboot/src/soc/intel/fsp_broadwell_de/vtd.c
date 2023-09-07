/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Siemens AG
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

#include <device/pci.h>
#include <device/pci_ops.h>
#include <device/pci_ids.h>
#include <soc/pci_devs.h>
#include <soc/acpi.h>
#include <soc/vtd.h>
#include <soc/broadwell_de.h>

#if ENV_RAMSTAGE

static void vtd_read_resources(struct device *dev)
{
	uint32_t vtbar;

	/* Add fixed MMIO resource for VT-d which was set up by the FSP. */
	vtbar = pci_read_config32(dev, VTBAR_OFFSET);
	if (vtbar & VTBAR_ENABLED) {
		mmio_resource(dev, VTBAR_OFFSET,
				(vtbar & VTBAR_MASK) / KiB, VTBAR_SIZE / KiB);
	}
}

static struct device_operations vtd_ops = {
	.read_resources		= vtd_read_resources,
	.set_resources		= DEVICE_NOOP,
	.write_acpi_tables	= vtd_write_acpi_tables,
};

static const struct pci_driver vtd_driver __pci_driver = {
	.ops    = &vtd_ops,
	.vendor = PCI_VENDOR_ID_INTEL,
	.device = VTD_DEVID,
};

#endif

uint8_t get_busno1(void)
{
	uint32_t reg32;

	/* Figure out what bus number is assigned for CPUBUSNO(1) */
	reg32 = pci_mmio_read_config32(VTD_PCI_DEV, VTD_CPUBUSNO);
	return ((reg32 >> VTD_CPUBUSNO_BUS1_SHIFT) & VTD_CPUBUSNO_BUS1_MASK);
}
