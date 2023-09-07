/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Google LLC
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/acpi.h>
#include <arch/acpi_device.h>
#include <arch/acpigen.h>
#include <bootstate.h>
#include <cbmem.h>
#include <device/pnp.h>
#include <ec/acpi/ec.h>
#include <pc80/keyboard.h>
#include <stdint.h>
#include <stdlib.h>

#include "commands.h"
#include "ec.h"
#include "chip.h"

/*
 * The UCSI fields are defined in the UCSI specification at
 * https://www.intel.com/content/www/us/en/io/universal-serial-bus/usb-type-c-ucsi-spec.html
 * https://www.intel.com/content/www/us/en/io/universal-serial-bus/bios-implementation-of-ucsi.html
 */

static struct fieldlist ucsi_region_fields[] = {
	FIELDLIST_NAMESTR("VER0", 8),
	FIELDLIST_NAMESTR("VER1", 8),
	FIELDLIST_NAMESTR("RSV0", 8),
	FIELDLIST_NAMESTR("RSV1", 8),
	FIELDLIST_NAMESTR("CCI0", 8),
	FIELDLIST_NAMESTR("CCI1", 8),
	FIELDLIST_NAMESTR("CCI2", 8),
	FIELDLIST_NAMESTR("CCI3", 8),
	FIELDLIST_NAMESTR("CTL0", 8),
	FIELDLIST_NAMESTR("CTL1", 8),
	FIELDLIST_NAMESTR("CTL2", 8),
	FIELDLIST_NAMESTR("CTL3", 8),
	FIELDLIST_NAMESTR("CTL4", 8),
	FIELDLIST_NAMESTR("CTL5", 8),
	FIELDLIST_NAMESTR("CTL6", 8),
	FIELDLIST_NAMESTR("CTL7", 8),
	FIELDLIST_NAMESTR("MGI0", 8),
	FIELDLIST_NAMESTR("MGI1", 8),
	FIELDLIST_NAMESTR("MGI2", 8),
	FIELDLIST_NAMESTR("MGI3", 8),
	FIELDLIST_NAMESTR("MGI4", 8),
	FIELDLIST_NAMESTR("MGI5", 8),
	FIELDLIST_NAMESTR("MGI6", 8),
	FIELDLIST_NAMESTR("MGI7", 8),
	FIELDLIST_NAMESTR("MGI8", 8),
	FIELDLIST_NAMESTR("MGI9", 8),
	FIELDLIST_NAMESTR("MGIA", 8),
	FIELDLIST_NAMESTR("MGIB", 8),
	FIELDLIST_NAMESTR("MGIC", 8),
	FIELDLIST_NAMESTR("MGID", 8),
	FIELDLIST_NAMESTR("MGIE", 8),
	FIELDLIST_NAMESTR("MGIF", 8),
	FIELDLIST_NAMESTR("MGO0", 8),
	FIELDLIST_NAMESTR("MGO1", 8),
	FIELDLIST_NAMESTR("MGO2", 8),
	FIELDLIST_NAMESTR("MGO3", 8),
	FIELDLIST_NAMESTR("MGO4", 8),
	FIELDLIST_NAMESTR("MGO5", 8),
	FIELDLIST_NAMESTR("MGO6", 8),
	FIELDLIST_NAMESTR("MGO7", 8),
	FIELDLIST_NAMESTR("MGO8", 8),
	FIELDLIST_NAMESTR("MGO9", 8),
	FIELDLIST_NAMESTR("MGOA", 8),
	FIELDLIST_NAMESTR("MGOB", 8),
	FIELDLIST_NAMESTR("MGOC", 8),
	FIELDLIST_NAMESTR("MGOD", 8),
	FIELDLIST_NAMESTR("MGOE", 8),
	FIELDLIST_NAMESTR("MGOF", 8),
};
static const size_t ucsi_region_len = ARRAY_SIZE(ucsi_region_fields);

static void wilco_ec_post_complete(void *unused)
{
	wilco_ec_send(KB_BIOS_PROGRESS, BIOS_PROGRESS_POST_COMPLETE);
}
BOOT_STATE_INIT_ENTRY(BS_PAYLOAD_LOAD, BS_ON_EXIT,
		      wilco_ec_post_complete, NULL);

static void wilco_ec_post_memory_init(void *unused)
{
	wilco_ec_send(KB_BIOS_PROGRESS, BIOS_PROGRESS_MEMORY_INIT);
}
BOOT_STATE_INIT_ENTRY(BS_PRE_DEVICE, BS_ON_EXIT,
		      wilco_ec_post_memory_init, NULL);

static void wilco_ec_post_video_init(void *unused)
{
	wilco_ec_send(KB_BIOS_PROGRESS, BIOS_PROGRESS_VIDEO_INIT);
}
BOOT_STATE_INIT_ENTRY(BS_DEV_INIT, BS_ON_EXIT,
		      wilco_ec_post_video_init, NULL);

static void wilco_ec_post_logo_displayed(void *unused)
{
	wilco_ec_send(KB_BIOS_PROGRESS, BIOS_PROGRESS_LOGO_DISPLAYED);
}
BOOT_STATE_INIT_ENTRY(BS_POST_DEVICE, BS_ON_EXIT,
		      wilco_ec_post_logo_displayed, NULL);

static void wilco_ec_resume(void *unused)
{
	wilco_ec_send_noargs(KB_RESTORE);
}
BOOT_STATE_INIT_ENTRY(BS_OS_RESUME, BS_ON_ENTRY, wilco_ec_resume, NULL);

static void wilco_ec_init(struct device *dev)
{
	if (!dev->enabled)
		return;

	/* Disable S0ix support in EC RAM with ACPI EC interface */
	if (!acpi_is_wakeup_s3()) {
		ec_set_ports(CONFIG_EC_BASE_ACPI_COMMAND,
			     CONFIG_EC_BASE_ACPI_DATA);
		ec_write(EC_RAM_S0IX_SUPPORT, 0);
	}

	/* Print EC firmware information */
	wilco_ec_print_all_info();

	/* Initialize keyboard, ignore emulated PS/2 mouse */
	pc_keyboard_init(NO_AUX_DEVICE);

	/* Direct power button to the host for processing */
	wilco_ec_send(KB_POWER_BUTTON_TO_HOST, 1);

	/* Unmute speakers */
	wilco_ec_send(KB_HW_MUTE_CONTROL, AUDIO_UNMUTE_125MS);

	/* Enable WiFi radio */
	wilco_ec_radio_control(RADIO_WIFI, 1);

	/* Turn on camera power */
	wilco_ec_send(KB_CAMERA, CAMERA_ON);
}

static void wilco_ec_resource(struct device *dev, int index,
			      size_t base, size_t size)
{
	struct resource *res = new_resource(dev, index);
	res->flags = IORESOURCE_IO | IORESOURCE_ASSIGNED | IORESOURCE_FIXED;
	res->base = base;
	res->size = size;
}

static void wilco_ec_read_resources(struct device *dev)
{
	/* ACPI command and data regions */
	wilco_ec_resource(dev, 0, CONFIG_EC_BASE_ACPI_DATA, 8);

	/* Host command and data regions */
	wilco_ec_resource(dev, 1, CONFIG_EC_BASE_HOST_DATA, 8);

	/* Packet region */
	wilco_ec_resource(dev, 2, CONFIG_EC_BASE_PACKET, 16);
}

static void wilco_ec_fill_ssdt_generator(struct device *dev)
{
	struct opregion opreg;
	void *region_ptr;

	if (!dev->enabled)
		return;

	region_ptr = cbmem_add(CBMEM_ID_ACPI_UCSI, ucsi_region_len);
	if (!region_ptr)
		return;
	memset(region_ptr, 0, ucsi_region_len);

	opreg.name = "UCSM";
	opreg.regionspace = SYSTEMMEMORY;
	opreg.regionoffset = (uintptr_t)region_ptr;
	opreg.regionlen = ucsi_region_len;

	acpigen_write_scope(acpi_device_path_join(dev, "UCSI"));
	acpigen_write_name("_CRS");
	acpigen_write_resourcetemplate_header();
	acpigen_write_mem32fixed(1, (uintptr_t)region_ptr, ucsi_region_len);
	acpigen_write_resourcetemplate_footer();
	acpigen_write_opregion(&opreg);
	acpigen_write_field(opreg.name, ucsi_region_fields, ucsi_region_len,
			    FIELD_ANYACC | FIELD_LOCK | FIELD_PRESERVE);
	acpigen_pop_len(); /* Scope */
}

static const char *wilco_ec_acpi_name(const struct device *dev)
{
	return "EC0";
}

static struct device_operations ops = {
	.init				= wilco_ec_init,
	.read_resources			= wilco_ec_read_resources,
	.enable_resources		= DEVICE_NOOP,
	.set_resources			= DEVICE_NOOP,
	.acpi_fill_ssdt_generator	= wilco_ec_fill_ssdt_generator,
	.acpi_name			= wilco_ec_acpi_name,
};

static struct pnp_info info[] = {
	{ NULL, 0, 0, 0, }
};

static void wilco_ec_enable_dev(struct device *dev)
{
	pnp_enable_devices(dev, &ops, ARRAY_SIZE(info), info);
}

struct chip_operations ec_google_wilco_ops = {
	CHIP_NAME("Google Wilco EC")
	.enable_dev = wilco_ec_enable_dev,
};
