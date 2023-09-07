/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Bitland Tech Inc.
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

#include "../panel.h"

struct panel_serializable_data AUO_B101UAN08_3 = {
	.edid = {
		.ascii_string = "B101UAN08.3",
		.manufacturer_name = "AUO",
		.panel_bits_per_color = 8,
		.panel_bits_per_pixel = 24,
		.mode = {
			.pixel_clock = 159192,
			.lvds_dual_channel = 0,
			.refresh = 60,
			.ha = 1200, .hbl = 144, .hso = 60, .hspw = 4,
			.va = 1920, .vbl = 60, .vso = 34, .vspw = 2,
			.phsync = '-', .pvsync = '-',
			.x_mm = 135, .y_mm = 216,
		},
	},
	.orientation = LB_FB_ORIENTATION_LEFT_UP,
	.init = {
		INIT_DELAY_CMD(24),
		INIT_DCS_CMD(0xB0, 0x01),
		INIT_DCS_CMD(0xC0, 0x48),
		INIT_DCS_CMD(0xC1, 0x48),
		INIT_DCS_CMD(0xC2, 0x47),
		INIT_DCS_CMD(0xC3, 0x47),
		INIT_DCS_CMD(0xC4, 0x46),
		INIT_DCS_CMD(0xC5, 0x46),
		INIT_DCS_CMD(0xC6, 0x45),
		INIT_DCS_CMD(0xC7, 0x45),
		INIT_DCS_CMD(0xC8, 0x64),
		INIT_DCS_CMD(0xC9, 0x64),
		INIT_DCS_CMD(0xCA, 0x4F),
		INIT_DCS_CMD(0xCB, 0x4F),
		INIT_DCS_CMD(0xCC, 0x40),
		INIT_DCS_CMD(0xCD, 0x40),
		INIT_DCS_CMD(0xCE, 0x66),
		INIT_DCS_CMD(0xCF, 0x66),
		INIT_DCS_CMD(0xD0, 0x4F),
		INIT_DCS_CMD(0xD1, 0x4F),
		INIT_DCS_CMD(0xD2, 0x41),
		INIT_DCS_CMD(0xD3, 0x41),
		INIT_DCS_CMD(0xD4, 0x48),
		INIT_DCS_CMD(0xD5, 0x48),
		INIT_DCS_CMD(0xD6, 0x47),
		INIT_DCS_CMD(0xD7, 0x47),
		INIT_DCS_CMD(0xD8, 0x46),
		INIT_DCS_CMD(0xD9, 0x46),
		INIT_DCS_CMD(0xDA, 0x45),
		INIT_DCS_CMD(0xDB, 0x45),
		INIT_DCS_CMD(0xDC, 0x64),
		INIT_DCS_CMD(0xDD, 0x64),
		INIT_DCS_CMD(0xDE, 0x4F),
		INIT_DCS_CMD(0xDF, 0x4F),
		INIT_DCS_CMD(0xE0, 0x40),
		INIT_DCS_CMD(0xE1, 0x40),
		INIT_DCS_CMD(0xE2, 0x66),
		INIT_DCS_CMD(0xE3, 0x66),
		INIT_DCS_CMD(0xE4, 0x4F),
		INIT_DCS_CMD(0xE5, 0x4F),
		INIT_DCS_CMD(0xE6, 0x41),
		INIT_DCS_CMD(0xE7, 0x41),
		INIT_DELAY_CMD(150),
		INIT_END_CMD,
	},
};
