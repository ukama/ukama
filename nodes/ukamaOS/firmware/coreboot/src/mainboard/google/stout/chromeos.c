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

#include <bootmode.h>
#include <boot/coreboot_tables.h>
#include <device/pci_ops.h>
#include <console/console.h>
#include <device/device.h>
#include <device/pci.h>

#include <southbridge/intel/bd82x6x/pch.h>
#include <southbridge/intel/common/gpio.h>
#include <vendorcode/google/chromeos/chromeos.h>
#include "ec.h"
#include <ec/quanta/it8518/ec.h>

void fill_lb_gpios(struct lb_gpios *gpios)
{
	struct lb_gpio chromeos_gpios[] = {
		/* Write Protect: GPIO7 */
		{7, ACTIVE_LOW, !get_write_protect_state(), "write protect"},

		/* Lid Switch: Virtual switch */
		{-1, ACTIVE_HIGH, get_lid_switch(), "lid"},

		/* Power Button: Virtual switch */
		/* Hard-code value to de-asserted */
		{-1, ACTIVE_HIGH, 0, "power"},

		/* Was VGA Option ROM loaded? */
		/* -1 indicates that this is a pseudo GPIO */
		{-1, ACTIVE_HIGH, gfx_get_init_done(), "oprom"},

		/* EC is in RW mode when it isn't in recovery mode. */
		{-1, ACTIVE_HIGH, !get_recovery_mode_switch(), "ec_in_rw"}
	};
	lb_add_gpios(gpios, chromeos_gpios, ARRAY_SIZE(chromeos_gpios));
}

int get_write_protect_state(void)
{
	return !get_gpio(7);
}

int get_lid_switch(void)
{
	/* hard-code to open */
	return 1;
}

/*
 * The recovery-switch is virtual on Stout and is handled via the EC.
 * Stout recovery mode is only valid if RTC_PWR_STS is set and the EC
 * indicated the recovery keys were pressed. We use a global flag for
 * rec_mode to be used after RTC_POWER_STS has been cleared. This function
 * is complicated by romstage support, which can't use a global variable.
 * Note, rec_mode is the only time the EC is in RO mode, otherwise, RW.
 */
int get_recovery_mode_switch(void)
{
	MAYBE_STATIC_BSS int ec_in_rec_mode = 0;
	MAYBE_STATIC_BSS int ec_rec_flag_good = 0;

	if (ec_rec_flag_good)
		return ec_in_rec_mode;

	pci_devfn_t dev = PCI_DEV(0, 0x1f, 0);
	u8 reg8 = pci_s_read_config8(dev, GEN_PMCON_3);

	u8 ec_status = ec_read(EC_STATUS_REG);

	printk(BIOS_SPEW,"%s:  EC status:%#x   RTC_BAT: %x\n",
			__func__, ec_status, reg8 & RTC_BATTERY_DEAD);

	ec_in_rec_mode = (((reg8 & RTC_BATTERY_DEAD) != 0) &&
				((ec_status & 0x3) == EC_IN_RECOVERY_MODE));
	ec_rec_flag_good = 1;
	return ec_in_rec_mode;
}

static const struct cros_gpio cros_gpios[] = {
	CROS_GPIO_REC_AH(CROS_GPIO_VIRTUAL, CROS_GPIO_DEVICE_NAME),
	CROS_GPIO_REC_AH(CROS_GPIO_VIRTUAL, CROS_GPIO_DEVICE_NAME),
	CROS_GPIO_WP_AL(7, CROS_GPIO_DEVICE_NAME),
};

void mainboard_chromeos_acpi_generate(void)
{
	chromeos_acpi_gpio_generate(cros_gpios, ARRAY_SIZE(cros_gpios));
}
