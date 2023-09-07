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

#include <device/device.h>
#include <device/pnp.h>
#include <arch/acpigen.h>
#include <console/console.h>

static void generic_set_resources(struct device *dev)
{
	struct resource *res;

	for (res = dev->resource_list; res; res = res->next) {
		if (!(res->flags & IORESOURCE_ASSIGNED))
			continue;

		res->flags |= IORESOURCE_STORED;
		report_resource_stored(dev, res, "");
	}
}

static void generic_read_resources(struct device *dev)
{
	struct resource *res = new_resource(dev, 0);
	res->base = dev->path.pnp.port;
	res->size = 2;
	res->flags = IORESOURCE_IO | IORESOURCE_ASSIGNED | IORESOURCE_FIXED;
}

#if CONFIG(HAVE_ACPI_TABLES)
static void generic_ssdt(struct device *dev)
{
	const char *scope = acpi_device_scope(dev);
	const char *name = acpi_device_name(dev);

	if (!scope || !name) {
		printk(BIOS_ERR, "%s: Missing ACPI path/scope\n",
		       dev_path(dev));
		return;
	}

	/* Device */
	acpigen_write_scope(scope);
	acpigen_write_device(name);

	printk(BIOS_DEBUG, "%s.%s: %s\n", scope, name, dev_path(dev));

	acpigen_write_name_string("_HID", "PNP0C02");
	acpigen_write_name_string("_DDN", dev_name(dev));

	/* OperationRegion("IOID", SYSTEMIO, port, 2) */
	struct opregion opreg = OPREGION("IOID", SYSTEMIO, dev->path.pnp.port, 2);
	acpigen_write_opregion(&opreg);

	struct fieldlist l[] = {
		FIELDLIST_OFFSET(0),
		FIELDLIST_NAMESTR("INDX", 8),
		FIELDLIST_NAMESTR("DATA", 8),
	};

	/* Field (IOID, AnyAcc, NoLock, Preserve)
	 * {
	 *  Offset (0),
	 *  INDX,   8,
	 *  DATA,   8,
	 * } */
	acpigen_write_field(opreg.name, l, ARRAY_SIZE(l), FIELD_BYTEACC | FIELD_NOLOCK |
			    FIELD_PRESERVE);

	struct fieldlist i[] = {
		FIELDLIST_OFFSET(0x07),
		FIELDLIST_NAMESTR("LDN", 8),
		FIELDLIST_OFFSET(0x21),
		FIELDLIST_NAMESTR("SCF1", 8),
		FIELDLIST_NAMESTR("SCF2", 8),
		FIELDLIST_NAMESTR("SCF3", 8),
		FIELDLIST_NAMESTR("SCF4", 8),
		FIELDLIST_NAMESTR("SCF5", 8),
		FIELDLIST_NAMESTR("SCF6", 8),
		FIELDLIST_NAMESTR("SCF7", 8),
		FIELDLIST_OFFSET(0x29),
		FIELDLIST_NAMESTR("CKCF", 8),
		FIELDLIST_OFFSET(0x2F),
		FIELDLIST_NAMESTR("SCFF", 8),
		FIELDLIST_OFFSET(0x30),
		FIELDLIST_NAMESTR("ACT0", 1),
		FIELDLIST_NAMESTR("ACT1", 1),
		FIELDLIST_NAMESTR("ACT2", 1),
		FIELDLIST_NAMESTR("ACT3", 1),
		FIELDLIST_NAMESTR("ACT4", 1),
		FIELDLIST_NAMESTR("ACT5", 1),
		FIELDLIST_NAMESTR("ACT6", 1),
		FIELDLIST_NAMESTR("ACT7", 1),
		FIELDLIST_OFFSET(0x60),
		FIELDLIST_NAMESTR("IOH0", 8),
		FIELDLIST_NAMESTR("IOL0", 8),
		FIELDLIST_NAMESTR("IOH1", 8),
		FIELDLIST_NAMESTR("IOL1", 8),
		FIELDLIST_NAMESTR("IOH2", 8),
		FIELDLIST_NAMESTR("IOL2", 8),
		FIELDLIST_NAMESTR("IOH3", 8),
		FIELDLIST_NAMESTR("IOL3", 8),
		FIELDLIST_OFFSET(0x70),
		/* Interrupt level 0 (IRQ number) */
		FIELDLIST_NAMESTR("ITL0", 4),
		FIELDLIST_OFFSET(0x71),
		/* Interrupt type 0 */
		FIELDLIST_NAMESTR("ITT0", 2),
		FIELDLIST_OFFSET(0x72),
		/* Interrupt level 1 (IRQ number) */
		FIELDLIST_NAMESTR("ITL1", 4),
		FIELDLIST_OFFSET(0x73),
		/* Interrupt type 1 */
		FIELDLIST_NAMESTR("ITT1", 2),
		FIELDLIST_OFFSET(0x74),
		FIELDLIST_NAMESTR("DMCH", 8),
		FIELDLIST_OFFSET(0xE0),
		FIELDLIST_NAMESTR("RGE0", 8),
		FIELDLIST_NAMESTR("RGE1", 8),
		FIELDLIST_NAMESTR("RGE2", 8),
		FIELDLIST_NAMESTR("RGE3", 8),
		FIELDLIST_NAMESTR("RGE4", 8),
		FIELDLIST_NAMESTR("RGE5", 8),
		FIELDLIST_NAMESTR("RGE6", 8),
		FIELDLIST_NAMESTR("RGE7", 8),
		FIELDLIST_NAMESTR("RGE8", 8),
		FIELDLIST_NAMESTR("RGE9", 8),
		FIELDLIST_NAMESTR("RGEA", 8),
		FIELDLIST_OFFSET(0xF0),
		FIELDLIST_NAMESTR("OPT0", 8),
		FIELDLIST_NAMESTR("OPT1", 8),
		FIELDLIST_NAMESTR("OPT2", 8),
		FIELDLIST_NAMESTR("OPT3", 8),
		FIELDLIST_NAMESTR("OPT4", 8),
		FIELDLIST_NAMESTR("OPT5", 8),
		FIELDLIST_NAMESTR("OPT6", 8),
		FIELDLIST_NAMESTR("OPT7", 8),
		FIELDLIST_NAMESTR("OPT8", 8),
		FIELDLIST_NAMESTR("OPT9", 8),
	};

	acpigen_write_indexfield("INDX", "DATA", i, ARRAY_SIZE(i), FIELD_BYTEACC |
				 FIELD_NOLOCK | FIELD_PRESERVE);

	acpigen_pop_len(); /* Device */
	acpigen_pop_len(); /* Scope */
}

static const char *generic_acpi_name(const struct device *dev)
{
	return "SIO0";
}
#endif

static struct device_operations ops = {
	.read_resources   = generic_read_resources,
	.set_resources    = generic_set_resources,
	.enable_resources = DEVICE_NOOP,
#if CONFIG(HAVE_ACPI_TABLES)
	.acpi_fill_ssdt_generator = generic_ssdt,
	.acpi_name = generic_acpi_name,
#endif
};

static void enable_dev(struct device *dev)
{
	if (dev->path.type != DEVICE_PATH_PNP)
		printk(BIOS_ERR, "%s: Unsupported device type\n", dev_path(dev));
	else if (!dev->path.pnp.port)
		printk(BIOS_ERR, "%s: Base address not set\n", dev_path(dev));
	else
		dev->ops = &ops;

	/*
	 * Need to call enable_dev() on the devices "behind" the Generic Super I/O.
	 * coreboot's generic allocator doesn't expect them behind PnP devices.
	 */
	enable_static_devices(dev);
}

struct chip_operations superio_common_ops = {
	CHIP_NAME("Generic Super I/O")
	.enable_dev = enable_dev,
};
