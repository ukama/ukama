/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
 * Copyright (C) 2015 Google Inc.
 * Copyright (C) 2015 Intel Corporation
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

#include <baseboard/variant.h>
#include <ec/google/chromeec/ec.h>
#include <gpio.h>
#include <soc/romstage.h>
#include <variant/ec.h>
#include <variant/gpio.h>

#include "spd/spd_util.h"
#include "spd/spd.h"


void mainboard_memory_init_params(FSPM_UPD *mupd)
{
	FSP_M_CONFIG *mem_cfg = &mupd->FspmConfig;

#ifdef EC_ENABLE_KEYBOARD_BACKLIGHT
	/* Turn on keyboard backlight to indicate we are booting */
	const FSPM_ARCH_UPD *arch_upd = &mupd->FspmArchUpd;
	if (arch_upd->BootMode != FSP_BOOT_ON_S3_RESUME)
		google_chromeec_kbbacklight(25);
#endif

	/* Get SPD index */
	const gpio_t spd_gpios[] = {
		GPIO_MEM_CONFIG_0,
		GPIO_MEM_CONFIG_1,
		GPIO_MEM_CONFIG_2,
		GPIO_MEM_CONFIG_3,
	};
	const int spd_idx = gpio_base2_value(spd_gpios, ARRAY_SIZE(spd_gpios));

	mem_cfg->MemorySpdDataLen = SPD_LEN;
	mem_cfg->DqPinsInterleaved = FALSE;

	spd_memory_init_params(mupd, spd_idx);
	variant_memory_init_params(mupd, spd_idx);
}
