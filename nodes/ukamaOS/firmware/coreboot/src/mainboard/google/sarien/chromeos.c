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

#include <arch/acpi.h>
#include <boot/coreboot_tables.h>
#include <gpio.h>
#include <soc/gpio.h>
#include <variant/gpio.h>
#include <vendorcode/google/chromeos/chromeos.h>
#include <security/tpm/tss.h>
#include <device/device.h>
#include <intelblocks/pmclib.h>

enum rec_mode_state {
	REC_MODE_UNINITIALIZED,
	REC_MODE_NOT_REQUESTED,
	REC_MODE_REQUESTED,
};

void fill_lb_gpios(struct lb_gpios *gpios)
{
	struct lb_gpio chromeos_gpios[] = {
		{GPIO_PCH_WP, ACTIVE_HIGH, get_write_protect_state(),
		 "write protect"},
		{-1, ACTIVE_HIGH, get_lid_switch(), "lid"},
		{-1, ACTIVE_HIGH, 0, "power"},
		{-1, ACTIVE_HIGH, gfx_get_init_done(), "oprom"},
		{-1, ACTIVE_HIGH, 0, "EC in RW"},
	};
	lb_add_gpios(gpios, chromeos_gpios, ARRAY_SIZE(chromeos_gpios));
}

static int cros_get_gpio_value(int type)
{
	const struct cros_gpio *cros_gpios;
	size_t i, num_gpios = 0;

	cros_gpios = variant_cros_gpios(&num_gpios);

	for (i = 0; i < num_gpios; i++) {
		const struct cros_gpio *gpio = &cros_gpios[i];
		if (gpio->type == type) {
			int state = gpio_get(gpio->gpio_num);
			if (gpio->polarity == CROS_GPIO_ACTIVE_LOW)
				return !state;
			else
				return state;
		}
	}
	return 0;
}

void mainboard_chromeos_acpi_generate(void)
{
	const struct cros_gpio *cros_gpios;
	size_t num_gpios = 0;

	cros_gpios = variant_cros_gpios(&num_gpios);

	chromeos_acpi_gpio_generate(cros_gpios, num_gpios);
}

int get_write_protect_state(void)
{
	return cros_get_gpio_value(CROS_GPIO_WP);
}

int get_recovery_mode_switch(void)
{
	static enum rec_mode_state saved_rec_mode = REC_MODE_UNINITIALIZED;
	enum rec_mode_state state = REC_MODE_NOT_REQUESTED;
	uint8_t cr50_state = 0;

	/* Check cached state, since TPM will only tell us the first time */
	if (saved_rec_mode != REC_MODE_UNINITIALIZED)
		return saved_rec_mode == REC_MODE_REQUESTED;

	/*
	 * Read one-time recovery request from cr50 in verstage only since
	 * the TPM driver won't be set up in time for other stages like romstage
	 * and the value from the TPM would be wrong anyway since the verstage
	 * read would have cleared the value on the TPM.
	 *
	 * The TPM recovery request is passed between stages through vboot data
	 * or cbmem depending on stage.
	 */
	if (ENV_VERSTAGE &&
	    tlcl_cr50_get_recovery_button(&cr50_state) == TPM_SUCCESS &&
	    cr50_state)
		state = REC_MODE_REQUESTED;

	/* Read state from the GPIO controlled by servo. */
	if (cros_get_gpio_value(CROS_GPIO_REC))
		state = REC_MODE_REQUESTED;

	/* Store the state in case this is called again in verstage. */
	saved_rec_mode = state;

	return state == REC_MODE_REQUESTED;
}

int get_lid_switch(void)
{
	return 1;
}

void mainboard_prepare_cr50_reset(void)
{
	/* Ensure system powers up after CR50 reset */
	if (ENV_RAMSTAGE)
		pmc_soc_set_afterg3_en(true);
}
