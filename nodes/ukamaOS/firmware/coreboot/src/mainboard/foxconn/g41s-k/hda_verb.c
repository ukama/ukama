/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2015 Damien Zammit <damien@zamaudio.com>
 * Copyright (C) 2017 Samuel Holland <samuel@sholland.org>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; either version 2 of
 * the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <device/azalia_device.h>

#if CONFIG(BOARD_FOXCONN_G41S_K)
const u32 cim_verb_data[] = {
	/* coreboot specific header */
	0x10ec0888, /* Vendor ID */
	0x105b0dda, /* Subsystem ID */
	0x0000000e, /* Number of entries */

	/* Pin Widget Verb Table */

	AZALIA_PIN_CFG(0, 0x11, 0x99430140),
	AZALIA_PIN_CFG(0, 0x12, 0x411111f0),
	AZALIA_PIN_CFG(0, 0x14, 0x01014410),
	AZALIA_PIN_CFG(0, 0x15, 0x411111f0),
	AZALIA_PIN_CFG(0, 0x16, 0x411111f0),
	AZALIA_PIN_CFG(0, 0x17, 0x411111f0),
	AZALIA_PIN_CFG(0, 0x18, 0x01a19c50),
	AZALIA_PIN_CFG(0, 0x19, 0x02a19c60),
	AZALIA_PIN_CFG(0, 0x1a, 0x0181345f),
	AZALIA_PIN_CFG(0, 0x1b, 0x02214c20),
	AZALIA_PIN_CFG(0, 0x1c, 0x411111f0),
	AZALIA_PIN_CFG(0, 0x1d, 0x4004c601),
	AZALIA_PIN_CFG(0, 0x1e, 0x01441130),
	AZALIA_PIN_CFG(0, 0x1f, 0x411111f0),
};
#else /* CONFIG_BOARD_FOXCONN_G41M */
const u32 cim_verb_data[] = {
	/* coreboot specific header */
	0x10ec0888, /* Vendor ID */
	0x105b0dc0, /* Subsystem ID */
	0x0000000e, /* Number of entries */

	/* Pin Widget Verb Table */

	AZALIA_PIN_CFG(2, 0x11, 0x01441140),
	AZALIA_PIN_CFG(2, 0x12, 0x411111f0),
	AZALIA_PIN_CFG(2, 0x14, 0x01014410),
	AZALIA_PIN_CFG(2, 0x15, 0x01011412),
	AZALIA_PIN_CFG(2, 0x16, 0x01016411),
	AZALIA_PIN_CFG(2, 0x17, 0x01012414),
	AZALIA_PIN_CFG(2, 0x18, 0x01a19c50),
	AZALIA_PIN_CFG(2, 0x19, 0x02a19c60),
	AZALIA_PIN_CFG(2, 0x1a, 0x0181345f),
	AZALIA_PIN_CFG(2, 0x1b, 0x02014c20),
	AZALIA_PIN_CFG(2, 0x1c, 0x593301f0),
	AZALIA_PIN_CFG(2, 0x1d, 0x4007f603),
	AZALIA_PIN_CFG(2, 0x1e, 0x99430130),
	AZALIA_PIN_CFG(2, 0x1f, 0x411111f0),
};
#endif

const u32 pc_beep_verbs[0] = {};

const u32 pc_beep_verbs_size = ARRAY_SIZE(pc_beep_verbs);
const u32 cim_verb_data_size = ARRAY_SIZE(cim_verb_data);
