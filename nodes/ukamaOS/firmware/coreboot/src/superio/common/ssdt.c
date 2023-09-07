/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <superio/common/ssdt.h>

#include <device/device.h>
#include <device/pnp.h>
#include <arch/acpigen.h>
#include <arch/acpi.h>
#include <device/pnp_def.h>
#include <console/console.h>
#include <types.h>

struct superio_dev {
	const char *acpi_hid;
	u16 io_base[4];
	u8 irq[2];
};

static const struct superio_dev superio_devs[] = {
	{ACPI_HID_FDC, {0x3f0, 0x3f2, 0x3f7}, {6, } },
	{ACPI_HID_KEYBOARD, {60, 64, }, {1, } },
	{ACPI_HID_MOUSE, {60, 64, }, {12, } },
	{ACPI_HID_COM, {0x3f8, 0x2f8, 0x3e8, 0x2e8}, {4, 3} },
	{ACPI_HID_LPT, {0x378, }, {7, } },
};

static const u8 io_idx[] = {PNP_IDX_IO0, PNP_IDX_IO1, PNP_IDX_IO2, PNP_IDX_IO3};
static const u8 irq_idx[] = {PNP_IDX_IRQ0, PNP_IDX_IRQ1};

static const struct superio_dev *superio_guess_function(struct device *dev)
{
	for (size_t i = 0; i < ARRAY_SIZE(io_idx); i++) {
		struct resource *res = probe_resource(dev, io_idx[i]);
		if (!res || !res->base)
			continue;

		for (size_t j = 0; j < ARRAY_SIZE(superio_devs); j++) {
			for (size_t k = 0; k < 4; k++) {
				if (!superio_devs[j].io_base[k])
					continue;
				if (superio_devs[j].io_base[k] == res->base)
					return &superio_devs[j];
			}
		}
	}
	for (size_t i = 0; i < ARRAY_SIZE(irq_idx); i++) {
		struct resource *res = probe_resource(dev, irq_idx[i]);
		if (!res || !res->size)
			continue;
		for (size_t j = 0; j < ARRAY_SIZE(superio_devs); j++) {
			for (size_t k = 0; k < 2; k++) {
				if (!superio_devs[j].irq[k])
					continue;
				if (superio_devs[j].irq[k] == res->base)
					return &superio_devs[j];
			}
		}
	}
	return NULL;
}

/* Return true if there are resources to report */
static bool has_resources(struct device *dev)
{
	for (size_t i = 0; i < ARRAY_SIZE(io_idx); i++) {
		struct resource *res = probe_resource(dev, io_idx[i]);
		if (!res || !res->base || !res->size)
			continue;
		return 1;
	}
	for (size_t i = 0; i < ARRAY_SIZE(irq_idx); i++) {
		struct resource *res = probe_resource(dev, irq_idx[i]);
		if (!res || !res->size || res->base > 16)
			continue;
		return 1;
	}
	return 0;
}

/* Add IO and IRQ resources for _CRS or _PRS */
static void ldn_gen_resources(struct device *dev)
{
	uint16_t irq = 0;
	for (size_t i = 0; i < ARRAY_SIZE(io_idx); i++) {
		struct resource *res = probe_resource(dev, io_idx[i]);
		if (!res || !res->base)
			continue;
		resource_t base = res->base;
		resource_t size = res->size;
		while (size > 0) {
			resource_t sz = size > 255 ? 255 : size;
			/* TODO: Needs test with regions >= 256 bytes */
			acpigen_write_io16(base, base, 1, sz, 1);
			size -= sz;
			base += sz;
		}
	}
	for (size_t i = 0; i < ARRAY_SIZE(irq_idx); i++) {
		struct resource *res = probe_resource(dev, irq_idx[i]);
		if (!res || !res->size || res->base >= 16)
			continue;
		irq |= 1 << res->base;
	}
	if (irq)
		acpigen_write_irq(irq);

}

/* Add resource base and size for additional SuperIO code */
static void ldn_gen_resources_use(struct device *dev)
{
	char name[5];
	for (size_t i = 0; i < ARRAY_SIZE(io_idx); i++) {
		struct resource *res = probe_resource(dev, io_idx[i]);
		if (!res || !res->base || !res->size)
			continue;

		snprintf(name, sizeof(name), "IO%zXB", i);
		name[4] = '\0';
		acpigen_write_name_integer(name, res->base);

		snprintf(name, sizeof(name), "IO%zXS", i);
		name[4] = '\0';
		acpigen_write_name_integer(name, res->size);
	}
}

const char *superio_common_ldn_acpi_name(const struct device *dev)
{
	u8 ldn = dev->path.pnp.device & 0xff;
	u8 vldn = (dev->path.pnp.device >> 8) & 0x7;
	static char name[5];

	snprintf(name, sizeof(name), "L%02X%01X", ldn, vldn);

	name[4] = '\0';

	return name;
}

static const char *name_from_hid(const char *hid)
{
	static const struct {
		const char *hid;
		const char *name;
	} lookup[] = {
		{ACPI_HID_FDC, "FDC" },
		{ACPI_HID_KEYBOARD, "PS2 Keyboard" },
		{ACPI_HID_MOUSE, "PS2 Mouse"},
		{ACPI_HID_COM, "COM port" },
		{ACPI_HID_LPT, "LPT" },
		{ACPI_HID_PNP, "Generic PNP device" },
	};

	for (size_t i = 0; hid && i < ARRAY_SIZE(lookup); i++) {
		if (strcmp(hid, lookup[i].hid) == 0)
			return lookup[i].name;
	}
	return "Generic device";
}

void superio_common_fill_ssdt_generator(struct device *dev)
{
	const char *scope = acpi_device_scope(dev);
	const char *name = acpi_device_name(dev);
	const u8 ldn = dev->path.pnp.device & 0xff;
	const u8 vldn = (dev->path.pnp.device >> 8) & 0x7;
	const char *hid;

	if (!scope || !name) {
		printk(BIOS_ERR, "%s: Missing ACPI path/scope\n", dev_path(dev));
		return;
	}
	if (vldn) {
		printk(BIOS_DEBUG, "%s: Ignoring virtual LDN\n", dev_path(dev));
		return;
	}

	printk(BIOS_DEBUG, "%s.%s: %s\n", scope, name, dev_path(dev));

	/* Scope */
	acpigen_write_scope(scope);

	/* Device */
	acpigen_write_device(name);

	acpigen_write_name_byte("_UID", 0);
	acpigen_write_name_byte("LDN", ldn);
	acpigen_write_name_byte("VLDN", vldn);

	acpigen_write_STA(dev->enabled ? 0xf : 0);

	if (!dev->enabled) {
		acpigen_pop_len(); /* Device */
		acpigen_pop_len(); /* Scope */
		return;
	}

	if (has_resources(dev)) {
		/* Resources - _CRS */
		acpigen_write_name("_CRS");
		acpigen_write_resourcetemplate_header();
		ldn_gen_resources(dev);
		acpigen_write_resourcetemplate_footer();

		/* Resources - _PRS */
		acpigen_write_name("_PRS");
		acpigen_write_resourcetemplate_header();
		ldn_gen_resources(dev);
		acpigen_write_resourcetemplate_footer();

		/* Resources base and size for 3rd party ACPI code */
		ldn_gen_resources_use(dev);
	}

	hid = acpi_device_hid(dev);
	if (!hid) {
		printk(BIOS_ERR, "%s: SuperIO driver doesn't provide a _HID\n", dev_path(dev));
		/* Try to guess it... */
		const struct superio_dev *sdev = superio_guess_function(dev);
		if (sdev && sdev->acpi_hid) {
			hid = sdev->acpi_hid;
			printk(BIOS_WARNING, "%s: Guessed _HID is '%s'\n", dev_path(dev), hid);
		} else {
			hid = ACPI_HID_PNP;
			printk(BIOS_ERR, "%s: Failed to guessed _HID\n", dev_path(dev));
		}
	}

	acpigen_write_name_string("_HID", hid);
	acpigen_write_name_string("_DDN", name_from_hid(hid));

	acpigen_pop_len(); /* Device */
	acpigen_pop_len(); /* Scope */
}
