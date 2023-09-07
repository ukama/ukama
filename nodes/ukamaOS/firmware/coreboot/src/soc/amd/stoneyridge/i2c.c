/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 Google
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

#include <device/mmio.h>
#include <arch/acpi.h>
#include <console/console.h>
#include <delay.h>
#include <drivers/i2c/designware/dw_i2c.h>
#include <amdblocks/acpimmio.h>
#include <soc/iomap.h>
#include <soc/pci_devs.h>
#include <soc/southbridge.h>
#include <soc/i2c.h>
#include "chip.h"

#define I2C_BUS_ADDRESS(x) (I2C_BASE_ADDRESS + I2C_DEVICE_SIZE * (x))
#define I2CA_BASE_ADDRESS (I2C_BUS_ADDRESS(0))
#define I2CB_BASE_ADDRESS (I2C_BUS_ADDRESS(1))
#define I2CC_BASE_ADDRESS (I2C_BUS_ADDRESS(2))
#define I2CD_BASE_ADDRESS (I2C_BUS_ADDRESS(3))

/* Global to provide access to chip.c */
const char *i2c_acpi_name(const struct device *dev);

static const uintptr_t i2c_bus_address[] = {
	I2CA_BASE_ADDRESS,
	I2CB_BASE_ADDRESS,
	I2CC_BASE_ADDRESS,
	I2CD_BASE_ADDRESS,
};

uintptr_t dw_i2c_base_address(unsigned int bus)
{
	return bus < I2C_DEVICE_COUNT ? i2c_bus_address[bus] : 0;
}

static const struct soc_amd_stoneyridge_config *get_soc_config(void)
{
	const struct device *dev = pcidev_path_on_root(GNB_DEVFN);

	if (!dev || !dev->chip_info) {
		printk(BIOS_ERR, "%s: Could not find SoC devicetree config!\n",
			__func__);
		return NULL;
	}

	return dev->chip_info;
}

const struct dw_i2c_bus_config *dw_i2c_get_soc_cfg(unsigned int bus)
{
	const struct soc_amd_stoneyridge_config *config;

	if (bus >= ARRAY_SIZE(i2c_bus_address))
		return NULL;

	config = get_soc_config();
	if (config == NULL)
		return NULL;

	return &config->i2c[bus];
}

const char *i2c_acpi_name(const struct device *dev)
{
	switch (dev->path.mmio.addr) {
	case I2CA_BASE_ADDRESS:
		return "I2CA";
	case I2CB_BASE_ADDRESS:
		return "I2CB";
	case I2CC_BASE_ADDRESS:
		return "I2CC";
	case I2CD_BASE_ADDRESS:
		return "I2CD";
	default:
		return NULL;
	}
}

int dw_i2c_soc_dev_to_bus(struct device *dev)
{
	switch (dev->path.mmio.addr) {
	case I2CA_BASE_ADDRESS:
		return 0;
	case I2CB_BASE_ADDRESS:
		return 1;
	case I2CC_BASE_ADDRESS:
		return 2;
	case I2CD_BASE_ADDRESS:
		return 3;
	}
	return -1;
}

static void dw_i2c_soc_init(bool is_early_init)
{
	size_t i;
	const struct soc_amd_stoneyridge_config *config;

	config = get_soc_config();

	if (config == NULL)
		return;

	for (i = 0; i < ARRAY_SIZE(config->i2c); i++) {
		const struct dw_i2c_bus_config *cfg  = &config->i2c[i];

		if (cfg->early_init != is_early_init)
			continue;

		if (dw_i2c_init(i, cfg))
			printk(BIOS_ERR, "Failed to init i2c bus %zd\n", i);
	}
}

void i2c_soc_early_init(void)
{
	dw_i2c_soc_init(true);
}

void i2c_soc_init(void)
{
	dw_i2c_soc_init(false);
}

struct device_operations stoneyridge_i2c_mmio_ops = {
	/* TODO(teravest): Move I2C resource info here. */
	.read_resources = DEVICE_NOOP,
	.set_resources = DEVICE_NOOP,
	.enable_resources = DEVICE_NOOP,
	.scan_bus = scan_smbus,
	.acpi_name = i2c_acpi_name,
	.acpi_fill_ssdt_generator = dw_i2c_acpi_fill_ssdt,
};

/*
 * I2C pins are open drain with external pull up, so in order to bit bang them
 * all, SCL pins must become GPIO inputs with no pull, then they need to be
 * toggled between input-no-pull and output-low. This table is for the initial
 * conversion of all SCL pins to input with no pull.
 */
static const struct soc_amd_gpio i2c_2_gpi[] = {
	PAD_GPI(I2C0_SCL_PIN, PULL_NONE),
	PAD_GPI(I2C1_SCL_PIN, PULL_NONE),
	PAD_GPI(I2C2_SCL_PIN, PULL_NONE),
	PAD_GPI(I2C3_SCL_PIN, PULL_NONE),
};
#define saved_pins_count ARRAY_SIZE(i2c_2_gpi)

/*
 * To program I2C pins without destroying their programming, the registers
 * that will be changed need to be saved first.
 */
static void save_i2c_pin_registers(uint8_t gpio,
					struct soc_amd_i2c_save *save_table)
{
	uint32_t *gpio_ptr;

	gpio_ptr = (uint32_t *)gpio_get_address(gpio);
	save_table->mux_value = iomux_read8(gpio);
	save_table->control_value = read32(gpio_ptr);
}

static void restore_i2c_pin_registers(uint8_t gpio,
					struct soc_amd_i2c_save *save_table)
{
	uint32_t *gpio_ptr;

	gpio_ptr = (uint32_t *)gpio_get_address(gpio);
	iomux_write8(gpio, save_table->mux_value);
	iomux_read8(gpio);
	write32(gpio_ptr, save_table->control_value);
	read32(gpio_ptr);
}

/* Slaves to be reset are controlled by devicetree register i2c_scl_reset */
void sb_reset_i2c_slaves(void)
{
	const struct soc_amd_stoneyridge_config *cfg;
	const struct device *dev = pcidev_path_on_root(GNB_DEVFN);
	struct soc_amd_i2c_save save_table[saved_pins_count];
	uint8_t i, j, control;

	if (!dev || !dev->chip_info)
		return;
	cfg = dev->chip_info;
	control = cfg->i2c_scl_reset & GPIO_I2C_MASK;
	if (control == 0)
		return;

	/* Save and reprogram I2C SCL pins */
	for (i = 0; i < saved_pins_count; i++)
		save_i2c_pin_registers(i2c_2_gpi[i].gpio, &save_table[i]);
	program_gpios(i2c_2_gpi, saved_pins_count);

	/*
	 * Toggle SCL back and forth 9 times under 100KHz. A single read is
	 * needed after the writes to force the posted write to complete.
	 */
	for (j = 0; j < 9; j++) {
		if (control & GPIO_I2C0_SCL)
			write32((uint32_t *)GPIO_I2C0_ADDRESS, GPIO_SCL_LOW);
		if (control & GPIO_I2C1_SCL)
			write32((uint32_t *)GPIO_I2C1_ADDRESS, GPIO_SCL_LOW);
		if (control & GPIO_I2C2_SCL)
			write32((uint32_t *)GPIO_I2C2_ADDRESS, GPIO_SCL_LOW);
		if (control & GPIO_I2C3_SCL)
			write32((uint32_t *)GPIO_I2C3_ADDRESS, GPIO_SCL_LOW);

		read32((uint32_t *)GPIO_I2C3_ADDRESS); /* Flush posted write */
		udelay(4); /* 4usec gets 85KHz for 1 pin, 70KHz for 4 pins */

		if (control & GPIO_I2C0_SCL)
			write32((uint32_t *)GPIO_I2C0_ADDRESS, GPIO_SCL_HIGH);
		if (control & GPIO_I2C1_SCL)
			write32((uint32_t *)GPIO_I2C1_ADDRESS, GPIO_SCL_HIGH);
		if (control & GPIO_I2C2_SCL)
			write32((uint32_t *)GPIO_I2C2_ADDRESS, GPIO_SCL_HIGH);
		if (control & GPIO_I2C3_SCL)
			write32((uint32_t *)GPIO_I2C3_ADDRESS, GPIO_SCL_HIGH);

		read32((uint32_t *)GPIO_I2C3_ADDRESS); /* Flush posted write */
		udelay(4);
	}

	/* Restore I2C pins. */
	for (i = 0; i < saved_pins_count; i++)
		restore_i2c_pin_registers(i2c_2_gpi[i].gpio, &save_table[i]);
}
