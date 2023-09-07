/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 MediaTek Inc.
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

#include <boardid.h>
#include <cbfs.h>
#include <console/console.h>
#include <soc/emi.h>

static const char *const sdram_configs[] = {
	[1] = "sdram-lpddr4x-H9HCNNNCPMALHR-4GB",
	[2] = "sdram-lpddr4x-MT53E1G32D4NQ-4GB",
	[3] = "sdram-lpddr4x-KMDH6001DA-B422-4GB",
	[4] = "sdram-lpddr4x-KMDP6001DA-B425-4GB",
	[5] = "sdram-lpddr4x-MT29VZZZAD8DQKSL-4GB",
	[6] = "sdram-lpddr4x-KMDV6001DA-B620-4GB",
	[7] = "sdram-lpddr4x-SDADA4CR-128G-4GB",
	[8] = "sdram-lpddr4x-K4UBE3D4AA-MGCL-4GB",
};

static struct sdram_params params;

const struct sdram_params *get_sdram_config(void)
{
	uint32_t ramcode = ram_code();

	if (ramcode >= ARRAY_SIZE(sdram_configs) ||
	    cbfs_boot_load_file(sdram_configs[ramcode], &params, sizeof(params),
				CBFS_TYPE_STRUCT) != sizeof(params))
		die("Cannot load SDRAM parameter file!");

	return &params;
}
