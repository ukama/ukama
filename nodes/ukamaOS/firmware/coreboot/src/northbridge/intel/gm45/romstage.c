/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
 * Copyright (C) 2011 The ChromiumOS Authors.  All rights reserved.
 * Copyright (C) 2014 Vladimir Serbinenko
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

#include <cbmem.h>
#include <romstage_handoff.h>
#include <console/console.h>
#include <device/pci_ops.h>
#include <arch/acpi.h>
#include <cpu/x86/lapic.h>
#include <arch/romstage.h>
#include <northbridge/intel/gm45/gm45.h>
#include <southbridge/intel/i82801ix/i82801ix.h>
#include <southbridge/intel/common/gpio.h>
#include <southbridge/intel/common/pmclib.h>
#include <string.h>

#define LPC_DEV PCI_DEV(0, 0x1f, 0)
#define MCH_DEV PCI_DEV(0, 0, 0)

void __weak mb_setup_superio(void)
{
}

void __weak mb_pre_raminit_setup(sysinfo_t *sysinfo)
{
}

void __weak mb_post_raminit_setup(void)
{
}

/* Platform has no romstage entry point under mainboard directory,
 * so this one is named with prefix mainboard.
 */
void mainboard_romstage_entry(void)
{
	sysinfo_t sysinfo;
	int s3resume = 0;
	int cbmem_initted;
	u16 reg16;

	/* basic northbridge setup, including MMCONF BAR */
	gm45_early_init();

	enable_lapic();

	/* First, run everything needed for console output. */
	i82801ix_early_init();
	setup_pch_gpios(&mainboard_gpio_map);

	reg16 = pci_read_config16(LPC_DEV, D31F0_GEN_PMCON_3);
	pci_write_config16(LPC_DEV, D31F0_GEN_PMCON_3, reg16);
	if ((MCHBAR16(SSKPD_MCHBAR) == 0xCAFE) && !(reg16 & (1 << 9))) {
		printk(BIOS_DEBUG, "soft reset detected, rebooting properly\n");
		gm45_early_reset();
	}

	/* ASPM related setting, set early by original BIOS. */
	DMIBAR16(0x204) &= ~(3 << 10);

	/* Check for S3 resume. */
	s3resume = southbridge_detect_s3_resume();

	/* RAM initialization */
	enter_raminit_or_reset();
	memset(&sysinfo, 0, sizeof(sysinfo));
	get_mb_spd_addrmap(sysinfo.spd_map);
	const struct device *dev;
	dev = pcidev_on_root(2, 0);
	if (dev)
		sysinfo.enable_igd = dev->enabled;
	dev = pcidev_on_root(1, 0);
	if (dev)
		sysinfo.enable_peg = dev->enabled;
	get_gmch_info(&sysinfo);

	mb_pre_raminit_setup(&sysinfo);

	raminit(&sysinfo, s3resume);

	mb_post_raminit_setup();

	const u32 deven = pci_read_config32(MCH_DEV, D0F0_DEVEN);
	/* Disable D4F0 (unknown signal controller). */
	pci_write_config32(MCH_DEV, D0F0_DEVEN, deven & ~0x4000);

	init_pm(&sysinfo, 0);

	i82801ix_dmi_setup();
	gm45_late_init(sysinfo.stepping);
	i82801ix_dmi_poll_vc1();

	MCHBAR16(SSKPD_MCHBAR) = 0xCAFE;

	init_iommu();

	cbmem_initted = !cbmem_recovery(s3resume);

	romstage_handoff_init(cbmem_initted && s3resume);

	printk(BIOS_SPEW, "exit main()\n");
}
