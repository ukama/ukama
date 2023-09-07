/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
 * Copyright (C) 2015 secunet Security Networks AG
 * Copyright (C) 2011 Google Inc
 * Copyright (C) 2018 Patrick Rudolph <patrick.rudolph@9elements.com>
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

#include <stdlib.h>
#include <console/console.h>
#include <arch/io.h>
#include <device/mmio.h>
#include <device/device.h>
#include <device/pci_ops.h>
#include <device/pci_def.h>
#include <pc80/mc146818rtc.h>
#include <romstage_handoff.h>
#include <types.h>

#include "sandybridge.h"

static void systemagent_vtd_init(void)
{
	const u32 capid0_a = pci_read_config32(PCI_DEV(0, 0, 0), CAPID0_A);
	if (capid0_a & (1 << 23))
		return;

	/* setup BARs */
	MCHBAR32(0x5404) = IOMMU_BASE1 >> 32;
	MCHBAR32(0x5400) = IOMMU_BASE1 | 1;
	MCHBAR32(0x5414) = IOMMU_BASE2 >> 32;
	MCHBAR32(0x5410) = IOMMU_BASE2 | 1;

	/* lock policies */
	write32((void *)(IOMMU_BASE1 + 0xff0), 0x80000000);

	const struct device *const azalia = pcidev_on_root(0x1b, 0);
	if (azalia && azalia->enabled) {
		write32((void *)(IOMMU_BASE2 + 0xff0), 0x20000000);
		write32((void *)(IOMMU_BASE2 + 0xff0), 0xa0000000);
	} else {
		write32((void *)(IOMMU_BASE2 + 0xff0), 0x80000000);
	}
}

static void enable_pam_region(void)
{
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM0, 0x30);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM1, 0x33);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM2, 0x33);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM3, 0x33);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM4, 0x33);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM5, 0x33);
	pci_write_config8(PCI_DEV(0, 0x00, 0), PAM6, 0x33);
}

static void sandybridge_setup_bars(void)
{
	printk(BIOS_DEBUG, "Setting up static northbridge registers...");
	/* Set up all hardcoded northbridge BARs */
	pci_write_config32(PCI_DEV(0, 0x00, 0), EPBAR, DEFAULT_EPBAR | 1);
	pci_write_config32(PCI_DEV(0, 0x00, 0), EPBAR + 4, (0LL+DEFAULT_EPBAR) >> 32);
	pci_write_config32(PCI_DEV(0, 0x00, 0), MCHBAR, (uintptr_t)DEFAULT_MCHBAR | 1);
	pci_write_config32(PCI_DEV(0, 0x00, 0), MCHBAR + 4, (0LL+(uintptr_t)DEFAULT_MCHBAR) >> 32);
	pci_write_config32(PCI_DEV(0, 0x00, 0), DMIBAR, (uintptr_t)DEFAULT_DMIBAR | 1);
	pci_write_config32(PCI_DEV(0, 0x00, 0), DMIBAR + 4, (0LL+(uintptr_t)DEFAULT_DMIBAR) >> 32);

	printk(BIOS_DEBUG, " done\n");
}

static void sandybridge_setup_graphics(void)
{
	u32 reg32;
	u16 reg16;
	u8 reg8;
	u8 gfxsize;

	reg16 = pci_read_config16(PCI_DEV(0,2,0), PCI_DEVICE_ID);
	switch (reg16) {
	case 0x0102: /* GT1 Desktop */
	case 0x0106: /* GT1 Mobile */
	case 0x010a: /* GT1 Server */
	case 0x0112: /* GT2 Desktop */
	case 0x0116: /* GT2 Mobile */
	case 0x0122: /* GT2 Desktop >=1.3GHz */
	case 0x0126: /* GT2 Mobile >=1.3GHz */
	case 0x0152: /* IvyBridge */
	case 0x0156: /* IvyBridge */
	case 0x0162: /* IvyBridge */
	case 0x0166: /* IvyBridge */
	case 0x016a: /* IvyBridge */
		break;
	default:
		printk(BIOS_DEBUG, "Graphics not supported by this CPU/chipset.\n");
		return;
	}

	printk(BIOS_DEBUG, "Initializing Graphics...\n");

	if (get_option(&gfxsize, "gfx_uma_size") != CB_SUCCESS) {
		/* Setup IGD memory by setting GGC[7:3] = 1 for 32MB */
		gfxsize = 0;
	}
	reg16 = pci_read_config16(PCI_DEV(0,0,0), GGC);
	reg16 &= ~0x00f8;
	reg16 |= (gfxsize + 1) << 3;
	/* Program GTT memory by setting GGC[9:8] = 2MB */
	reg16 &= ~0x0300;
	reg16 |= 2 << 8;
	/* Enable VGA decode */
	reg16 &= ~0x0002;
	pci_write_config16(PCI_DEV(0,0,0), GGC, reg16);

	/* Enable 256MB aperture */
	reg8 = pci_read_config8(PCI_DEV(0, 2, 0), MSAC);
	reg8 &= ~0x06;
	reg8 |= 0x02;
	pci_write_config8(PCI_DEV(0, 2, 0), MSAC, reg8);

	/* Erratum workarounds */
	reg32 = MCHBAR32(0x5f00);
	reg32 |= (1 << 9)|(1 << 10);
	MCHBAR32(0x5f00) = reg32;

	/* Enable SA Clock Gating */
	reg32 = MCHBAR32(0x5f00);
	MCHBAR32(0x5f00) = reg32 | 1;

	/* GPU RC6 workaround for sighting 366252 */
	reg32 = MCHBAR32(0x5d14);
	reg32 |= (1 << 31);
	MCHBAR32(0x5d14) = reg32;

	/* VLW */
	reg32 = MCHBAR32(0x6120);
	reg32 &= ~(1 << 0);
	MCHBAR32(0x6120) = reg32;

	reg32 = MCHBAR32(0x5418);
	reg32 |= (1 << 4) | (1 << 5);
	MCHBAR32(0x5418) = reg32;
}

static void start_peg_link_training(void)
{
	u32 tmp;
	u32 deven;

	/* PEG on IvyBridge+ needs a special startup sequence.
	 * As the MRC has its own initialization code skip it. */
	if (((pci_read_config16(PCI_DEV(0, 0, 0), PCI_DEVICE_ID) &
			BASE_REV_MASK) != BASE_REV_IVB) ||
		CONFIG(HAVE_MRC))
		return;

	deven = pci_read_config32(PCI_DEV(0, 0, 0), DEVEN);

	if (deven & DEVEN_PEG10) {
		tmp = pci_read_config32(PCI_DEV(0, 1, 0), 0xC24) & ~(1 << 16);
		pci_write_config32(PCI_DEV(0, 1, 0), 0xC24, tmp | (1 << 5));
	}

	if (deven & DEVEN_PEG11) {
		tmp = pci_read_config32(PCI_DEV(0, 1, 1), 0xC24) & ~(1 << 16);
		pci_write_config32(PCI_DEV(0, 1, 1), 0xC24, tmp | (1 << 5));
	}

	if (deven & DEVEN_PEG12) {
		tmp = pci_read_config32(PCI_DEV(0, 1, 2), 0xC24) & ~(1 << 16);
		pci_write_config32(PCI_DEV(0, 1, 2), 0xC24, tmp | (1 << 5));
	}

	if (deven & DEVEN_PEG60) {
		tmp = pci_read_config32(PCI_DEV(0, 6, 0), 0xC24) & ~(1 << 16);
		pci_write_config32(PCI_DEV(0, 6, 0), 0xC24, tmp | (1 << 5));
	}
}

void systemagent_early_init(void)
{
	u32 capid0_a;
	u32 deven;
	u8 reg8;

	/* Device ID Override Enable should be done very early */
	capid0_a = pci_read_config32(PCI_DEV(0, 0, 0), 0xe4);
	if (capid0_a & (1 << 10)) {
		const size_t is_mobile = get_platform_type() == PLATFORM_MOBILE;

		reg8 = pci_read_config8(PCI_DEV(0, 0, 0), 0xf3);
		reg8 &= ~7; /* Clear 2:0 */

		if (is_mobile)
			reg8 |= 1; /* Set bit 0 */

		pci_write_config8(PCI_DEV(0, 0, 0), 0xf3, reg8);
	}

	/* Setup all BARs required for early PCIe and raminit */
	sandybridge_setup_bars();

	/* Set C0000-FFFFF to access RAM on both reads and writes */
	enable_pam_region();

	/* Setup IOMMU BARs */
	systemagent_vtd_init();

	/* Device Enable, don't touch PEG bits */
	deven = pci_read_config32(PCI_DEV(0, 0, 0), DEVEN) | DEVEN_IGD;
	pci_write_config32(PCI_DEV(0, 0, 0), DEVEN, deven);

	sandybridge_setup_graphics();

	/* Write magic value to start PEG link training.
	 * This should be done in PCI device enumeration, but
	 * the PCIe specification requires to wait at least 100msec
	 * after reset for devices to come up.
	 * As we don't want to increase boot time, enable it early and
	 * assume the PEG is up as soon as PCI enumeration starts.
	 * TODO: use time stamps to ensure the timings are met */
	start_peg_link_training();
}

void northbridge_romstage_finalize(int s3resume)
{
	MCHBAR16(SSKPD) = 0xCAFE;

	romstage_handoff_init(s3resume);
}
