/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

#include <console/console.h>
#include <fsp/util.h>

asmlinkage void chipset_teardown_car_main(void)
{
	FSP_INFO_HEADER *fih;
	uint32_t status;
	FSP_TEMP_RAM_EXIT temp_ram_exit;
	struct prog fsp = PROG_INIT(PROG_REFCODE, "fsp.bin");

	if (prog_locate(&fsp)) {
		die("Unable to locate fsp.bin\n");
	} else {
		/* This leaks a mapping which this code assumes is benign as
		 * the flash is memory mapped CPU's address space. */

		/* FIXME: the implementation of find_fsp is utter garbage
		   as it casts error values to FSP_INFO_HEADER pointers.
		   Checking for return values can only be done sanely once
		   that is fixed. */
		fih = find_fsp((uintptr_t)rdev_mmap_full(prog_rdev(&fsp)));
	}

	temp_ram_exit = (FSP_TEMP_RAM_EXIT)(fih->TempRamExitEntryOffset +
						fih->ImageBase);
	printk(BIOS_DEBUG, "Calling TempRamExit: %p\n", temp_ram_exit);
	status = temp_ram_exit(NULL);

	if (status != FSP_SUCCESS) {
		printk(BIOS_CRIT, "TempRamExit returned 0x%08x\n", status);
		die("TempRamExit returned an error!\n");
	}
}
