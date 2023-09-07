/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2015 Intel Corporation.
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

#include <console/console.h>
#include <fsp/api.h>
#include <soc/gpio.h>
#include <soc/pm.h>
#include <soc/romstage.h>
#include "spd/spd.h"
#include <spd_bin.h>


void mainboard_memory_init_params(FSPM_UPD *mupd)
{
	FSP_M_CONFIG *mem_cfg = &mupd->FspmConfig;

	struct spd_block blk = {
		.addr_map = { 0x50, 0x52, },
	};

	get_spd_smbus(&blk);
	dump_spd_info(&blk);
	printk(BIOS_SPEW, "spd block length: 0x%08x\n", blk.len);

	mem_cfg->MemorySpdPtr00 = (UINT32) blk.spd_array[0];
	mem_cfg->MemorySpdPtr10 = (UINT32) blk.spd_array[1];
	printk(BIOS_SPEW, "0x%08x: SpdDataBuffer_0_0\n", mem_cfg->MemorySpdPtr00);
	printk(BIOS_SPEW, "0x%08x: SpdDataBuffer_1_0\n", mem_cfg->MemorySpdPtr10);

	/*
	* Configure the DQ/DQS settings if required. In general the settings
	* should be set in the FSP flash image and should not need to be
	* changed.
	*/
	mainboard_fill_dq_map_data(&mem_cfg->DqByteMapCh0, &mem_cfg->DqByteMapCh1);
	mainboard_fill_dqs_map_data(&mem_cfg->DqsMapCpu2DramCh0, &mem_cfg->DqsMapCpu2DramCh1);
	mainboard_fill_rcomp_res_data(&mem_cfg->RcompResistor);
	mainboard_fill_rcomp_strength_data(&mem_cfg->RcompTarget);

	/* update spd length*/
	mem_cfg->MemorySpdDataLen = blk.len;
	mem_cfg->DqPinsInterleaved = TRUE;
}
