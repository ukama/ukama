/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2013 Google Inc.
 * Copyright (C) 2015 Intel Corp.
 * Copyright (C) 2018-2019 Eltan B.V.
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

#include <build.h>
#include <cbfs.h>
#include <console/console.h>
#include <chip.h>
#include <device/pci_ops.h>
#if CONFIG(VENDORCODE_ELTAN_MBOOT)
#include <mboot.h>
#endif
#include <soc/lpc.h>
#include <soc/pci_devs.h>
#include <soc/romstage.h>
#include <spd_bin.h>
#include <stdint.h>

#include "cpld.h"

void mainboard_memory_init_params(struct romstage_params *params,
				  MEMORY_INIT_UPD *memory_params)
{
	struct region_device spd_rdev;
	u8 spd_index = 0;

	if (!CONFIG(ONBOARD_SAMSUNG_MEM)) {
		if (cpld_read_pcb_version() <= 7)
			spd_index = 1;
		else
			spd_index = 2;
	}

	if (get_spd_cbfs_rdev(&spd_rdev, spd_index) < 0)
		die("spd.bin not found\n");

	memory_params->PcdMemoryTypeEnable = MEM_DDR3;
	memory_params->PcdMemorySpdPtr = (uintptr_t)rdev_mmap_full(&spd_rdev);
	memory_params->PcdMemChannel0Config = 1; /* Memory down */
	memory_params->PcdMemChannel1Config = 2; /* Disabled */
}

void mainboard_after_memory_init(void)
{
	printk(BIOS_DEBUG, "%s/%s called\n", __FILE__, __func__);

	/* Disable the Braswell UART hardware for COM1. */
	pci_write_config32(PCI_DEV(0, LPC_DEV, 0), UART_CONT, 0);
}

#if CONFIG(VENDORCODE_ELTAN_MBOOT)
/**
 * mb_crtm
 *
 * Measures the crtm version. This consists of a string than can be defined
 * using make menuconfig and automatically generated version information.
 *
 * @param[in]  activePcr	bitmap of the support
 *
 * @retval TPM_SUCCESS		Operation completed successfully.
 * @retval TPM_E_IOERROR	Unexpected device behavior.
 */

static const uint8_t crtm_version[] =
	CONFIG_VENDORCODE_ELTAN_CRTM_VERSION_STRING COREBOOT_VERSION COREBOOT_EXTRA_VERSION
	" " COREBOOT_BUILD;

int mb_crtm(void)
{
	int status = TPM_E_IOERROR;
	TCG_PCR_EVENT2_HDR tcgEventHdr;

	/* Use FirmwareVersion string to represent CRTM version. */
	printk(BIOS_DEBUG, "%s: Measure CRTM Version\n", __func__);
	memset(&tcgEventHdr, 0, sizeof(tcgEventHdr));
	tcgEventHdr.pcrIndex = MBOOT_PCR_INDEX_0;
	tcgEventHdr.eventType = EV_S_CRTM_VERSION;
	tcgEventHdr.eventSize = sizeof(crtm_version);
	printk(BIOS_DEBUG, "%s: EventSize - %u\n", __func__, tcgEventHdr.eventSize);

	status = mboot_hash_extend_log(0, (uint8_t *)crtm_version,
				       tcgEventHdr.eventSize, &tcgEventHdr,
				       (uint8_t *)crtm_version);
	if (status) {
		printk(BIOS_DEBUG, "Measure CRTM Version returned 0x%x\n", status);
	}

	return status;
}
#endif
