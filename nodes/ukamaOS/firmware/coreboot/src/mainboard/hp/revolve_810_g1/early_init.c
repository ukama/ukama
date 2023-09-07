/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2008-2009 coresystems GmbH
 * Copyright (C) 2014 Vladimir Serbinenko
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

#include <bootblock_common.h>
#include <string.h>
#include <cbfs.h>
#include <console/console.h>
#include <device/pci_ops.h>
#include <northbridge/intel/sandybridge/sandybridge.h>
#include <northbridge/intel/sandybridge/raminit_native.h>
#include <southbridge/intel/bd82x6x/pch.h>
#include <ec/hp/kbc1126/ec.h>

const struct southbridge_usb_port mainboard_usb_ports[] = {
	{ 1, 1, 0 },
	{ 1, 0, 0 },
	{ 1, 1, 1 },
	{ 0, 1, 1 },
	{ 0, 0, 2 },
	{ 1, 0, 2 },
	{ 0, 0, 3 },
	{ 0, 0, 3 },
	{ 1, 0, 4 }, /* B1P1: Digitizer */
	{ 1, 0, 4 }, /* B1P2: wlan USB, EHCI debug */
	{ 1, 1, 5 }, /* B1P3: Camera */
	{ 0, 0, 5 }, /* B1P4 */
	{ 1, 0, 6 }, /* B1P5: wwan USB */
	{ 0, 0, 6 }, /* B1P6 */
};

void bootblock_mainboard_early_init(void)
{
	kbc1126_enter_conf();
	kbc1126_mailbox_init();
	kbc1126_kbc_init();
	kbc1126_ec_init();
	kbc1126_pm1_init();
	kbc1126_exit_conf();
}

void mainboard_get_spd(spd_raw_data *spd, bool id_only)
{
	/* C1S0 is a soldered RAM with no real SPD. Use stored SPD.  */
	size_t spd_file_len = 0;
	void *spd_file = cbfs_boot_map_with_leak("spd.bin", CBFS_TYPE_SPD,
		&spd_file_len);

	if (!spd_file || spd_file_len < sizeof(spd_raw_data))
		die("SPD data for C1S0 not found.");

	read_spd(&spd[0], 0x50, id_only);
	memcpy(&spd[2], spd_file, spd_file_len);
}
