/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2011-2012 The ChromiumOS Authors.  All rights reserved.
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
#include <bootmode.h>
#include <boot/coreboot_tables.h>
#include <device/device.h>
#include <device/pci.h>

#include <southbridge/intel/bd82x6x/pch.h>
#include <southbridge/intel/common/gpio.h>
#include <ec/quanta/ene_kb3940q/ec.h>
#include <vendorcode/google/chromeos/chromeos.h>
#include "ec.h"

#define WP_GPIO		6
#define DEVMODE_GPIO	54
#define FORCE_RECOVERY_MODE	0

void fill_lb_gpios(struct lb_gpios *gpios)
{
	struct lb_gpio chromeos_gpios[] = {
		/* Write Protect: GPIO active Low */
		{WP_GPIO, ACTIVE_LOW, !get_write_protect_state(),
		 "write protect"},

		/* lid switch value from EC */
		{-1, ACTIVE_HIGH, get_lid_switch(), "lid"},

		/* Power Button - Hardcode Low as power button may still be
		 * pressed when read here.*/
		{-1, ACTIVE_HIGH, 0, "power"},

		/* Was VGA Option ROM loaded? */
		/* -1 indicates that this is a pseudo GPIO */
		{-1, ACTIVE_HIGH, gfx_get_init_done(), "oprom"},
	};
	lb_add_gpios(gpios, chromeos_gpios, ARRAY_SIZE(chromeos_gpios));
}

int get_write_protect_state(void)
{
	return !get_gpio(WP_GPIO);
}

int get_lid_switch(void)
{
	return (ec_mem_read(EC_HW_GPI_STATUS) >> EC_GPI_LID_STAT_BIT) & 1;
}

int get_recovery_mode_switch(void)
{
	int ec_rec_mode = 0;

	if (FORCE_RECOVERY_MODE) {
		printk(BIOS_DEBUG, "FORCING RECOVERY MODE.\n");
		return 1;
	}

	if (ENV_RAMSTAGE) {
		if (ec_mem_read(EC_CODE_STATE) == EC_COS_EC_RO)
			ec_rec_mode = 1;

		printk(BIOS_DEBUG, "RECOVERY MODE FROM EC: %x\n", ec_rec_mode);
	}

	return ec_rec_mode;
}

static const struct cros_gpio cros_gpios[] = {
	CROS_GPIO_REC_AH(CROS_GPIO_VIRTUAL, CROS_GPIO_DEVICE_NAME),
	CROS_GPIO_WP_AL(WP_GPIO, CROS_GPIO_DEVICE_NAME),
};

void mainboard_chromeos_acpi_generate(void)
{
	chromeos_acpi_gpio_generate(cros_gpios, ARRAY_SIZE(cros_gpios));
}
