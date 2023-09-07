/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Google LLC
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

#include <baseboard/variants.h>
#include <boardid.h>
#include <console/console.h>
#include <ec/google/chromeec/ec.h>
#include <smbios.h>
#include <string.h>
#include <stdint.h>

#define SKU_UNKNOWN		0xFFFFFFFF
#define SKU_MAX			255

uint32_t get_board_sku(void)
{
	static uint32_t sku_id = SKU_UNKNOWN;

	if (sku_id != SKU_UNKNOWN)
		return sku_id;

	if (google_chromeec_cbi_get_sku_id(&sku_id))
		sku_id = SKU_UNKNOWN;

	return sku_id;
}

const char *smbios_system_sku(void)
{
	static char sku_str[7]; /* sku{0..255} */
	uint32_t sku_id = get_board_sku();

	if ((sku_id == SKU_UNKNOWN) || (sku_id > SKU_MAX)) {
		printk(BIOS_ERR, "%s: Unexpected SKU ID %u\n",
			__func__, sku_id);
		return "";
	}

	snprintf(sku_str, sizeof(sku_str), "sku%u", sku_id);

	return sku_str;
}

const char *smbios_mainboard_manufacturer(void)
{
	static char oem_name[32];
	static const char *manuf;

	if (manuf)
		return manuf;

	if (google_chromeec_cbi_get_oem_name(&oem_name[0],
				ARRAY_SIZE(oem_name)) < 0) {
		printk(BIOS_INFO, "Couldn't obtain OEM name from CBI\n");
		manuf = CONFIG_MAINBOARD_SMBIOS_MANUFACTURER;
	} else {
		manuf = &oem_name[0];
	}

	return manuf;
}
