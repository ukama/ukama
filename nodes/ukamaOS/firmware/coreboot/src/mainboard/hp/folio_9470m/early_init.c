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
#include <device/pci_ops.h>
#include <northbridge/intel/sandybridge/sandybridge.h>
#include <northbridge/intel/sandybridge/raminit_native.h>
#include <southbridge/intel/bd82x6x/pch.h>
#include <ec/hp/kbc1126/ec.h>

const struct southbridge_usb_port mainboard_usb_ports[] = {
	{ 1, 1, 0 }, /* SSP1: dock */
	{ 1, 1, 0 }, /* SSP2: left, EHCI Debug */
	{ 1, 1, 1 }, /* SSP3: right back side */
	{ 1, 1, 1 }, /* SSP4: right front side */
	{ 1, 0, 2 }, /* B0P5 */
	{ 1, 0, 2 }, /* B0P6: wlan USB */
	{ 0, 0, 3 }, /* B0P7 */
	{ 1, 1, 3 }, /* B0P8: smart card reader */
	{ 1, 1, 4 }, /* B1P1: fingerprint reader */
	{ 0, 0, 4 }, /* B1P2: (EHCI Debug, not connected) */
	{ 1, 1, 5 }, /* B1P3: Camera */
	{ 0, 0, 5 }, /* B1P4 */
	{ 1, 1, 6 }, /* B1P5: wwan USB */
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
	read_spd(&spd[0], 0x50, id_only);
	read_spd(&spd[2], 0x52, id_only);
}
