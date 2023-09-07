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

// Use simple device model for this file even in ramstage
#define __SIMPLE_DEVICE__

#include <device/pci_ops.h>
#include <arch/romstage.h>
#include <cbmem.h>
#include <cpu/x86/mtrr.h>
#include <program_loading.h>
#include "e7505.h"

void *cbmem_top_chipset(void)
{
	pci_devfn_t mch = PCI_DEV(0, 0, 0);
	uintptr_t tolm;

	/* This is at 128 MiB boundary. */
	tolm = pci_read_config16(mch, TOLM) >> 11;
	tolm <<= 27;

	return (void *)tolm;
}

void northbridge_write_smram(u8 smram);

void northbridge_write_smram(u8 smram)
{
	pci_devfn_t mch = PCI_DEV(0, 0, 0);
	pci_write_config8(mch, SMRAMC, smram);
}

void fill_postcar_frame(struct postcar_frame *pcf)
{
	uintptr_t top_of_ram;

	/*
	 * Choose to NOT set ROM as WP cacheable here.
	 * Timestamps indicate the CPU this northbridge code is
	 * connected to, performs better for memcpy() and un-lzma
	 * operations when source is left as UC.
	 */

	pcf->skip_common_mtrr = 1;

	/* Cache RAM as WB from 0 -> CACHE_TMP_RAMTOP. */
	postcar_frame_add_mtrr(pcf, 0, CACHE_TMP_RAMTOP, MTRR_TYPE_WRBACK);

	/* Cache CBMEM region as WB. */
	top_of_ram = (uintptr_t)cbmem_top();
	postcar_frame_add_mtrr(pcf, top_of_ram - 8*MiB, 8*MiB,
		MTRR_TYPE_WRBACK);
}
