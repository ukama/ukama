/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Google LLC
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
#include <console/console.h>
#include <ec/google/chromeec/ec.h>
#include <gpio.h>
#include <memory_info.h>
#include <soc/cnl_memcfg_init.h>
#include <soc/romstage.h>
#include <string.h>
#include <variant/gpio.h>

/*
 * GPIO_MEM_CH_SEL is set to 1 for single channel skus
 * and 0 for dual channel skus.
 */
#define GPIO_MEM_CH_SEL		GPP_F2

int __weak variant_memory_sku(void)
{
	const gpio_t spd_gpios[] = {
		GPIO_MEM_CONFIG_0,
		GPIO_MEM_CONFIG_1,
		GPIO_MEM_CONFIG_2,
		GPIO_MEM_CONFIG_3,
	};

	return gpio_base2_value(spd_gpios, ARRAY_SIZE(spd_gpios));
}

void mainboard_memory_init_params(FSPM_UPD *memupd)
{
	struct cnl_mb_cfg memcfg;
	int mem_sku;
	int is_single_ch_mem;

	variant_memory_params(&memcfg);
	mem_sku = variant_memory_sku();
	/*
	 * GPP_F2 is the MEM_CH_SEL gpio, which is set to 1 for single
	 * channel skus and 0 for dual channel skus.
	 */
	is_single_ch_mem = gpio_get(GPIO_MEM_CH_SEL);

	/*
	 * spd[0]-spd[3] map to CH0D0, CH0D1, CH1D0, CH1D1 respectively.
	 * Dual-DIMM memory is not used in hatch family, so we only
	 * fill in spd_info for CH0D0 and CH1D0 here.
	 */
	memcfg.spd[0].read_type = READ_SPD_CBFS;
	memcfg.spd[0].spd_spec.spd_index = mem_sku;
	if (!is_single_ch_mem) {
		memcfg.spd[2].read_type = READ_SPD_CBFS;
		memcfg.spd[2].spd_spec.spd_index = mem_sku;
	}

	cannonlake_memcfg_init(&memupd->FspmConfig, &memcfg);
}

void mainboard_get_dram_part_num(const char **part_num, size_t *len)
{
	static char part_num_store[DIMM_INFO_PART_NUMBER_SIZE];
	static enum {
		PART_NUM_NOT_READ,
		PART_NUM_AVAILABLE,
		PART_NUM_NOT_IN_CBI,
	} part_num_state = PART_NUM_NOT_READ;

	if (part_num_state == PART_NUM_NOT_READ) {
		if (google_chromeec_cbi_get_dram_part_num(&part_num_store[0],
						sizeof(part_num_store)) < 0) {
			printk(BIOS_ERR, "No DRAM part number in CBI!\n");
			part_num_state = PART_NUM_NOT_IN_CBI;
		} else {
			part_num_state = PART_NUM_AVAILABLE;
		}
	}

	if (part_num_state == PART_NUM_NOT_IN_CBI)
		return;

	*part_num = &part_num_store[0];
	*len = strlen(part_num_store) + 1;
}
