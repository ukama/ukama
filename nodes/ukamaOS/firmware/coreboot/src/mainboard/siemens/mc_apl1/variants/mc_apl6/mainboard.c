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

#include <arch/io.h>
#include <bootstate.h>
#include <cf9_reset.h>
#include <console/console.h>
#include <device/pci_def.h>
#include <device/pci_ids.h>
#include <device/pci_ops.h>
#include <gpio.h>
#include <hwilib.h>
#include <intelblocks/lpc_lib.h>
#include <intelblocks/pcr.h>
#include <soc/pcr_ids.h>
#include <timer.h>
#include <timestamp.h>
#include <baseboard/variants.h>
#include <types.h>

#define TX_DWORD3	0xa8c

void variant_mainboard_final(void)
{
	struct device *dev = NULL;
	uint16_t cmd = 0;

	/* PIR6 register mapping for PCIe root ports
	 * INTA#->PIRQD#, INTB#->PIRQA#, INTC#->PIRQB#, INTD#-> PIRQC#
	 */
	pcr_write16(PID_ITSS, 0x314c, 0x2103);

	/* Enable CLKRUN_EN for power gating LPC */
	lpc_enable_pci_clk_cntl();

	/*
	 * Enable LPC PCE (Power Control Enable) by setting IOSF-SB port 0xD2
	 * offset 0x341D bit3 and bit0.
	 * Enable LPC CCE (Clock Control Enable) by setting IOSF-SB port 0xD2
	 * offset 0x341C bit [3:0].
	 */
	pcr_or32(PID_LPC, PCR_LPC_PRC, (PCR_LPC_CCE_EN | PCR_LPC_PCE_EN));

	/* Set Master Enable for on-board PCI device. */
	dev = dev_find_device(PCI_VENDOR_ID_SIEMENS, 0x403e, 0);
	if (dev) {
		cmd = pci_read_config16(dev, PCI_COMMAND);
		cmd |= PCI_COMMAND_MASTER;
		pci_write_config16(dev, PCI_COMMAND, cmd);

		/* Disable clock outputs 0-3 (CLKOUT) for upstream
		 * XIO2001 PCIe to PCI Bridge.
		 */
		struct device *parent = dev->bus->dev;
		if (parent && parent->device == PCI_DEVICE_ID_TI_XIO2001)
			pci_write_config8(parent, 0xd8, 0x0F);
	}

	/* Disable clock outputs 2-5 (CLKOUT) for another XIO2001 PCIe to PCI
	 * Bridge on this mainboard.
	 */
	dev = dev_find_device(PCI_VENDOR_ID_SIEMENS, 0x403f, 0);
	if (dev) {
		struct device *parent = dev->bus->dev;
		if (parent && parent->device == PCI_DEVICE_ID_TI_XIO2001)
			pci_write_config8(parent, 0xd8, 0x3c);
	}

	/* Set Full Reset Bit in Reset Control Register (I/O port CF9h).
	 * When Bit 3 is set to 1 and then the reset button is pressed the PCH
	 * will drive SLP_S3 active (low). SLP_S3 is then used on the mainboard
	 * to generate the right reset timing.
	 */
	outb(FULL_RST, RST_CNT);
}

static void wait_for_legacy_dev(void *unused)
{
	uint32_t legacy_delay, us_since_boot;
	struct stopwatch sw;

	/* Open main hwinfo block. */
	if (hwilib_find_blocks("hwinfo.hex") != CB_SUCCESS)
		return;

	/* Get legacy delay parameter from hwinfo. */
	if (hwilib_get_field(LegacyDelay, (uint8_t *) &legacy_delay,
			      sizeof(legacy_delay)) != sizeof(legacy_delay))
		return;

	us_since_boot = get_us_since_boot();
	/* No need to wait if the time since boot is already long enough.*/
	if (us_since_boot > legacy_delay)
		return;
	stopwatch_init_msecs_expire(&sw, (legacy_delay - us_since_boot) / 1000);
	printk(BIOS_NOTICE, "Wait remaining %d of %d us for legacy devices...",
			legacy_delay - us_since_boot, legacy_delay);
	stopwatch_wait_until_expired(&sw);
	printk(BIOS_NOTICE, "done!\n");
}

static void finalize_boot(void *unused)
{
	/* Set coreboot ready LED. */
	gpio_output(CNV_RGI_DT, 1);
}

BOOT_STATE_INIT_ENTRY(BS_DEV_ENUMERATE, BS_ON_ENTRY, wait_for_legacy_dev, NULL);
BOOT_STATE_INIT_ENTRY(BS_PAYLOAD_BOOT, BS_ON_ENTRY, finalize_boot, NULL);
