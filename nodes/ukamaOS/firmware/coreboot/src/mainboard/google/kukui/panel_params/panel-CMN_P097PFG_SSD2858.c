/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Huaqin Telecom Inc.
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

struct panel_serializable_data P097PFG_SSD2858 = {
	.edid = {
		.ascii_string = "P097PFG",
		.manufacturer_name = "CMN",
		.panel_bits_per_color = 8,
		.panel_bits_per_pixel = 24,
		.mode = {
			.pixel_clock = 211660,
			.lvds_dual_channel = 0,
			.refresh = 60,
			.ha = 1536, .hbl = 160, .hso = 140, .hspw = 10,
			.va = 2048, .vbl = 32, .vso = 20, .vspw = 2,
			.phsync = '-', .pvsync = '-',
			.x_mm = 147, .y_mm = 196,
		},
	},
	.orientation = LB_FB_ORIENTATION_NORMAL,
	.init = {
		INIT_GENERIC_CMD(0xff, 0x00),
		/* LOCKCNT=0x1f4, MRX=0, POSTDIV=1 (/2} }, MULT=0x49
		 * 27 Mhz => 985.5 Mhz */
		INIT_GENERIC_CMD(0x00, 0x08, 0x01, 0xf4, 0x01, 0x49),
		/* MTXDIV=1, SYSDIV=3 (=> 4) */
		INIT_GENERIC_CMD(0x00, 0x0c, 0x00, 0x00, 0x00, 0x03),
		/* MTXVPF=24bpp, MRXLS=4 lanes, MRXVB=bypass, MRXECC=1,
		 * MRXEOT=1, MRXEE=1 */
		INIT_GENERIC_CMD(0x00, 0x14, 0x0c, 0x3d, 0x80, 0x0f),
		INIT_GENERIC_CMD(0x00, 0x20, 0x15, 0x92, 0x56, 0x7d),
		INIT_GENERIC_CMD(0x00, 0x24, 0x00, 0x00, 0x30, 0x00),

		INIT_GENERIC_CMD(0x10, 0x08, 0x01, 0x20, 0x08, 0x45),
		INIT_GENERIC_CMD(0x10, 0x1c, 0x00, 0x00, 0x00, 0x00),
		INIT_GENERIC_CMD(0x20, 0x0c, 0x00, 0x00, 0x00, 0x04),
		/* Pixel clock 985.5 Mhz * 0x49/0x4b = 959 Mhz */
		INIT_GENERIC_CMD(0x20, 0x10, 0x00, 0x4b, 0x00, 0x49),
		INIT_GENERIC_CMD(0x20, 0xa0, 0x00, 0x00, 0x00, 0x00),
		/* EOT=1, LPE = 0, LSOUT=4 lanes, LPD=25 */
		INIT_GENERIC_CMD(0x60, 0x08, 0x00, 0xd9, 0x00, 0x08),
		INIT_GENERIC_CMD(0x60, 0x14, 0x01, 0x00, 0x01, 0x06),
		/* DSI0 enable (default: probably not needed) */
		INIT_GENERIC_CMD(0x60, 0x80, 0x00, 0x00, 0x00, 0x0f),
		/* DSI1 enable */
		INIT_GENERIC_CMD(0x60, 0xa0, 0x00, 0x00, 0x00, 0x0f),

		/* HSA=0x18, VSA=0x02, HBP=0x50, VBP=0x0c */
		INIT_GENERIC_CMD(0x60, 0x0c, 0x0c, 0x50, 0x02, 0x18),
		/* VACT= 0x800 (2048} }, VFP= 0x14, HFP=0x50 */
		INIT_GENERIC_CMD(0x60, 0x10, 0x08, 0x00, 0x14, 0x50),
		/* HACT=0x300 (768) */
		INIT_GENERIC_CMD(0x60, 0x84, 0x00, 0x00, 0x03, 0x00),
		INIT_GENERIC_CMD(0x60, 0xa4, 0x00, 0x00, 0x03, 0x00),

		/* Take panel out of sleep. */
		INIT_GENERIC_CMD(0xff, 0x01),
		INIT_DCS_CMD(0x11),
		INIT_DELAY_CMD(120),
		INIT_DCS_CMD(0x29),
		INIT_DELAY_CMD(20),
		INIT_GENERIC_CMD(0xff, 0x00),

		INIT_DELAY_CMD(120),
		INIT_DCS_CMD(0x11),
		INIT_DELAY_CMD(120),
		INIT_DCS_CMD(0x29),
		INIT_DELAY_CMD(20),
		INIT_END_CMD,
	},
};
