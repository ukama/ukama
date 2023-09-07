/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2008-2009 coresystems GmbH
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
#include <cpu/x86/smm.h>
#include <device/pci_ops.h>
#include <southbridge/intel/ibexpeak/nvs.h>
#include <southbridge/intel/common/pmutil.h>
#include <northbridge/intel/nehalem/nehalem.h>
#include <ec/acpi/ec.h>

void mainboard_smi_gpi(u32 gpi_sts)
{
}

int mainboard_smi_apmc(u8 data)
{
	u8 tmp;
	switch (data) {
	case APM_CNT_ACPI_ENABLE:
		tmp = pci_read_config8(PCI_DEV(0, 0x1f, 0), 0xbb);
		tmp &= ~0x03;
		tmp |= 0x02;
		pci_write_config8(PCI_DEV(0, 0x1f, 0), 0xbb, tmp);
		break;
	case APM_CNT_ACPI_DISABLE:
		tmp = pci_read_config8(PCI_DEV(0, 0x1f, 0), 0xbb);
		tmp &= ~0x03;
		tmp |= 0x01;
		pci_write_config8(PCI_DEV(0, 0x1f, 0), 0xbb, tmp);
		break;
	default:
		break;
	}
	return 0;
}
