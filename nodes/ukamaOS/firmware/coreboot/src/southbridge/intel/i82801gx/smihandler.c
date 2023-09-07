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

#include <types.h>
#include <console/console.h>
#include <cpu/x86/smm.h>
#include <device/pci_def.h>
#include <pc80/mc146818rtc.h>
#include <southbridge/intel/common/pmutil.h>
#include "i82801gx.h"

/* I945 */
#define SMRAM		0x9d
#define   D_OPEN	(1 << 6)
#define   D_CLS		(1 << 5)
#define   D_LCK		(1 << 4)
#define   G_SMRANE	(1 << 3)
#define   C_BASE_SEG	((0 << 2) | (1 << 1) | (0 << 0))

#include "nvs.h"

/* While we read PMBASE dynamically in case it changed, let's
 * initialize it with a sane value
 */
u16 pmbase = DEFAULT_PMBASE;
u8 smm_initialized = 0;

/* GNVS needs to be updated by an 0xEA PM Trap (B2) after it has been located
 * by coreboot.
 */
global_nvs_t *gnvs = (global_nvs_t *)0x0;

void southbridge_update_gnvs(u8 apm_cnt, int *smm_done)
{
	gnvs = *(global_nvs_t **)0x500;
	*smm_done = 1;
}

int southbridge_io_trap_handler(int smif)
{
	switch (smif) {
	case 0x32:
		printk(BIOS_DEBUG, "OS Init\n");
		/* gnvs->smif:
		 *  On success, the IO Trap Handler returns 0
		 *  On failure, the IO Trap Handler returns a value != 0
		 */
		gnvs->smif = 0;
		return 1; /* IO trap handled */
	}

	/* Not handled */
	return 0;
}

void southbridge_smi_monitor(void)
{
#define IOTRAP(x) (trap_sts & (1 << x))
	u32 trap_sts, trap_cycle;
	u32 data, mask = 0;
	int i;

	trap_sts = RCBA32(0x1e00); // TRSR - Trap Status Register
	RCBA32(0x1e00) = trap_sts; // Clear trap(s) in TRSR

	trap_cycle = RCBA32(0x1e10);
	for (i = 16; i < 20; i++) {
		if (trap_cycle & (1 << i))
			mask |= (0xff << ((i - 16) << 2));
	}


	/* IOTRAP(3) SMI function call */
	if (IOTRAP(3)) {
		if (gnvs && gnvs->smif)
			io_trap_handler(gnvs->smif); // call function smif
		return;
	}

	/* IOTRAP(2) currently unused
	 * IOTRAP(1) currently unused */

	/* IOTRAP(0) SMIC: currently unused  */

	printk(BIOS_DEBUG, "  trapped io address = 0x%x\n", trap_cycle & 0xfffc);
	for (i = 0; i < 4; i++)
		if (IOTRAP(i))
			printk(BIOS_DEBUG, "  TRAP = %d\n", i);
	printk(BIOS_DEBUG, "  AHBE = %x\n", (trap_cycle >> 16) & 0xf);
	printk(BIOS_DEBUG, "  MASK = 0x%08x\n", mask);
	printk(BIOS_DEBUG, "  read/write: %s\n", (trap_cycle & (1 << 24)) ? "read" : "write");

	if (!(trap_cycle & (1 << 24))) {
		/* Write Cycle */
		data = RCBA32(0x1e18);
		printk(BIOS_DEBUG, "  iotrap written data = 0x%08x\n", data);
	}
#undef IOTRAP
}

void southbridge_finalize_all(void)
{
}
