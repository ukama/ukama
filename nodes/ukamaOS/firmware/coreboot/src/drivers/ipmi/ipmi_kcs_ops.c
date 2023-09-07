/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation; version 2 of
 * the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

/*
 * Place in devicetree.cb:
 *
 * chip drivers/ipmi
 *   device pnp ca2.0 on end         # IPMI KCS
 * end
 */

#include <console/console.h>
#include <device/device.h>
#include <device/pnp.h>
#if CONFIG(HAVE_ACPI_TABLES)
#include <arch/acpi.h>
#include <arch/acpigen.h>
#endif
#if CONFIG(GENERATE_SMBIOS_TABLES)
#include <smbios.h>
#endif
#include <version.h>
#include <delay.h>
#include <timer.h>
#include "ipmi_kcs.h"
#include "chip.h"

/* 4 bit encoding */
static u8 ipmi_revision_major = 0x1;
static u8 ipmi_revision_minor = 0x0;

static int ipmi_get_device_id(struct device *dev, struct ipmi_devid_rsp *rsp)
{
	int ret;

	ret = ipmi_kcs_message(dev->path.pnp.port, IPMI_NETFN_APPLICATION, 0,
			     IPMI_BMC_GET_DEVICE_ID, NULL, 0, (u8 *)rsp,
			     sizeof(*rsp));
	if (ret < sizeof(struct ipmi_rsp) || rsp->resp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s command failed (ret=%d resp=0x%x)\n",
		       __func__, ret, rsp->resp.completion_code);
		return 1;
	}
	if (ret != sizeof(*rsp)) {
		printk(BIOS_ERR, "IPMI: %s response truncated\n", __func__);
		return 1;
	}
	return 0;
}

static int ipmi_get_bmc_self_test_result(struct device *dev, struct ipmi_selftest_rsp *rsp)
{
	int ret;

	ret = ipmi_kcs_message(dev->path.pnp.port, IPMI_NETFN_APPLICATION, 0,
				 IPMI_BMC_GET_SELFTEST_RESULTS, NULL, 0, (u8 *)rsp,
				 sizeof(*rsp));

	if (ret < sizeof(struct ipmi_rsp) || rsp->resp.completion_code) {
		printk(BIOS_ERR, "IPMI: %s command failed (ret=%d resp=0x%x)\n",
		       __func__, ret, rsp->resp.completion_code);
		return 1;
	}
	if (ret != sizeof(*rsp)) {
		printk(BIOS_ERR, "IPMI: %s response truncated\n", __func__);
		return 1;
	}

	return 0;
}

static void ipmi_kcs_init(struct device *dev)
{
	struct ipmi_devid_rsp rsp;
	uint32_t man_id = 0, prod_id = 0;
	struct drivers_ipmi_config *conf = NULL;
	struct ipmi_selftest_rsp selftestrsp;
	uint8_t retry_count;

	if (!dev->enabled)
		return;

	printk(BIOS_DEBUG, "IPMI: PNP KCS 0x%x\n", dev->path.pnp.port);

	if (dev->chip_info)
		conf = dev->chip_info;

	/* Get IPMI version for ACPI and SMBIOS */
	if (conf && conf->wait_for_bmc && conf->bmc_boot_timeout) {
		struct stopwatch sw;
		stopwatch_init_msecs_expire(&sw, conf->bmc_boot_timeout * 1000);
		printk(BIOS_DEBUG, "IPMI: Waiting for BMC...\n");

		while (!stopwatch_expired(&sw)) {
			if (inb(dev->path.pnp.port) != 0xff)
				break;
			mdelay(100);
		}
		if (stopwatch_expired(&sw)) {
			printk(BIOS_INFO, "IPMI: Waiting for BMC timed out\n");
			/* Don't write tables if communication failed */
			dev->enabled = 0;
			return;
		}
	}

	printk(BIOS_INFO, "Get BMC self test result...");
	for (retry_count = 0; retry_count < conf->bmc_boot_timeout; retry_count++) {
		if (!ipmi_get_bmc_self_test_result(dev, &selftestrsp))
			break;

		mdelay(1000);
	}

	switch (selftestrsp.result) {
	case IPMI_APP_SELFTEST_NO_ERROR: /* 0x55 */
		printk(BIOS_DEBUG, "No Error\n");
		break;
	case IPMI_APP_SELFTEST_NOT_IMPLEMENTED: /* 0x56 */
		printk(BIOS_DEBUG, "Function Not Implemented\n");
		break;
	case IPMI_APP_SELFTEST_ERROR: /* 0x57 */
		printk(BIOS_ERR, "BMC: Corrupted or inaccessible data or device\n");
		/* Don't write tables if communication failed */
		dev->enabled = 0;
		break;
	case IPMI_APP_SELFTEST_FATAL_HW_ERROR: /* 0x58 */
		printk(BIOS_ERR, "BMC: Fatal Hardware Error\n");
		/* Don't write tables if communication failed */
		dev->enabled = 0;
		break;
	case IPMI_APP_SELFTEST_RESERVED: /* 0xFF */
		printk(BIOS_DEBUG, "Reserved\n");
		break;

	default: /* Other Device Specific Hardware Error */
		printk(BIOS_ERR, "BMC: Device Specific Error\n");
		/* Don't write tables if communication failed */
		dev->enabled = 0;
		break;
	}

	if (!ipmi_get_device_id(dev, &rsp)) {
		/* Queried the IPMI revision from BMC */
		ipmi_revision_minor = IPMI_IPMI_VERSION_MINOR(rsp.ipmi_version);
		ipmi_revision_major = IPMI_IPMI_VERSION_MAJOR(rsp.ipmi_version);

		memcpy(&man_id, rsp.manufacturer_id,
		       sizeof(rsp.manufacturer_id));

		memcpy(&prod_id, rsp.product_id, sizeof(rsp.product_id));

		printk(BIOS_INFO, "IPMI: Found man_id 0x%06x, prod_id 0x%04x\n",
		       man_id, prod_id);

		printk(BIOS_INFO, "IPMI: Version %01x.%01x\n",
		       ipmi_revision_major, ipmi_revision_minor);
	} else {
		/* Don't write tables if communication failed */
		dev->enabled = 0;
	}
}

#if CONFIG(HAVE_ACPI_TABLES)
static uint32_t uid_cnt = 0;

static unsigned long
ipmi_write_acpi_tables(struct device *dev, unsigned long current,
		       struct acpi_rsdp *rsdp)
{
	struct drivers_ipmi_config *conf = NULL;
	struct acpi_spmi *spmi;
	s8 gpe_interrupt = -1;
	u32 apic_interrupt = 0;
	acpi_addr_t addr = {
		.space_id = ACPI_ADDRESS_SPACE_IO,
		.access_size = ACPI_ACCESS_SIZE_BYTE_ACCESS,
		.addrl = dev->path.pnp.port,
		.bit_width = 8,
	};

	switch (CONFIG_IPMI_KCS_REGISTER_SPACING) {
	case 4:
		addr.bit_offset = 32;
		break;
	case 16:
		addr.bit_offset = 128;
		break;
	default:
		printk(BIOS_ERR, "IPMI: Unsupported register spacing for SPMI\n");
		/* fall through */
	case 1:
		addr.bit_offset = 8;
		break;
	}

	current = ALIGN_UP(current, 8);
	printk(BIOS_DEBUG, "ACPI:    * SPMI at %lx\n", current);
	spmi = (struct acpi_spmi *)current;

	if (dev->chip_info)
		conf = dev->chip_info;

	if (conf) {
		if (conf->have_gpe)
			gpe_interrupt = conf->gpe_interrupt;
		if (conf->have_apic)
			apic_interrupt = conf->apic_interrupt;
	}

	/* Use command to get UID from ipmi_ssdt */
	acpi_create_ipmi(dev, spmi, (ipmi_revision_major << 8) |
			 (ipmi_revision_minor << 4), &addr,
			 IPMI_INTERFACE_KCS, gpe_interrupt, apic_interrupt,
			 dev->command);

	acpi_add_table(rsdp, spmi);

	current += spmi->header.length;

	return current;
}

static void ipmi_ssdt(struct device *dev)
{
	const char *scope = acpi_device_scope(dev);
	struct drivers_ipmi_config *conf = NULL;

	if (!scope) {
		printk(BIOS_ERR, "IPMI: Missing ACPI scope for %s\n",
		       dev_path(dev));
		return;
	}

	if (dev->chip_info)
		conf = dev->chip_info;

	/* Use command to pass UID to ipmi_write_acpi_tables */
	dev->command = uid_cnt++;

	/* write SPMI device */
	acpigen_write_scope(scope);
	acpigen_write_device("SPMI");
	acpigen_write_name_string("_HID", "IPI0001");
	acpigen_write_name_string("_STR", "IPMI_KCS");
	acpigen_write_name_byte("_UID", dev->command);
	acpigen_write_STA(0xf);
	acpigen_write_name("_CRS");
	acpigen_write_resourcetemplate_header();
	acpigen_write_io16(dev->path.pnp.port, dev->path.pnp.port, 1, 1, 1);
	acpigen_write_io16(dev->path.pnp.port + CONFIG_IPMI_KCS_REGISTER_SPACING,
			   dev->path.pnp.port + CONFIG_IPMI_KCS_REGISTER_SPACING, 1, 1, 1);

	if (conf) {
		// FIXME: is that correct?
		if (conf->have_apic)
			acpigen_write_irq(1 << conf->apic_interrupt);
	}

	acpigen_write_resourcetemplate_footer();

	acpigen_write_method("_IFT", 0);
	acpigen_write_return_byte(1);	// KCS
	acpigen_pop_len();

	acpigen_write_method("_SRV", 0);
	acpigen_write_return_integer((ipmi_revision_major << 8) |
				     (ipmi_revision_minor << 4));
	acpigen_pop_len();

	acpigen_pop_len(); /* pop device */
	acpigen_pop_len(); /* pop scope */
}
#endif

#if CONFIG(GENERATE_SMBIOS_TABLES)
static int ipmi_smbios_data(struct device *dev, int *handle,
			    unsigned long *current)
{
	struct drivers_ipmi_config *conf = NULL;
	u8 nv_storage = 0xff;
	u8 i2c_address = 0;
	u8 register_spacing;

	int len = 0;

	if (dev->chip_info)
		conf = dev->chip_info;

	if (conf) {
		if (conf->have_nv_storage)
			nv_storage = conf->nv_storage_device_address;
		i2c_address = conf->bmc_i2c_address;
	}

	switch (CONFIG_IPMI_KCS_REGISTER_SPACING) {
	case 4:
		register_spacing = 1 << 6;
		break;
	case 16:
		register_spacing = 2 << 6;
		break;
	default:
		printk(BIOS_ERR, "IPMI: Unsupported register spacing for SMBIOS\n");
		/* fall through */
	case 1:
		register_spacing = 0 << 6;
		break;
	}

	// add IPMI Device Information
	len += smbios_write_type38(
		current, handle,
		SMBIOS_BMC_INTERFACE_KCS,
		ipmi_revision_minor | (ipmi_revision_major << 4),
		i2c_address, // I2C address
		nv_storage, // NV storage
		dev->path.pnp.port | 1, // IO interface
		register_spacing,
		0); // no IRQ

	return len;
}
#endif

static void ipmi_set_resources(struct device *dev)
{
	struct resource *res;

	for (res = dev->resource_list; res; res = res->next) {
		if (!(res->flags & IORESOURCE_ASSIGNED))
			continue;

		res->flags |= IORESOURCE_STORED;
		report_resource_stored(dev, res, "");
	}
}

static void ipmi_read_resources(struct device *dev)
{
	struct resource *res = new_resource(dev, 0);
	res->base = dev->path.pnp.port;
	res->size = 2;
	res->flags = IORESOURCE_IO | IORESOURCE_ASSIGNED | IORESOURCE_FIXED;
}

static struct device_operations ops = {
	.read_resources   = ipmi_read_resources,
	.set_resources    = ipmi_set_resources,
	.enable_resources = DEVICE_NOOP,
	.init             = ipmi_kcs_init,
#if CONFIG(HAVE_ACPI_TABLES)
	.write_acpi_tables = ipmi_write_acpi_tables,
	.acpi_fill_ssdt_generator = ipmi_ssdt,
#endif
#if CONFIG(GENERATE_SMBIOS_TABLES)
	.get_smbios_data = ipmi_smbios_data,
#endif
};

static void enable_dev(struct device *dev)
{
	if (dev->path.type != DEVICE_PATH_PNP)
		printk(BIOS_ERR, "%s: Unsupported device type\n",
		       dev_path(dev));
	else if (dev->path.pnp.port & 1)
		printk(BIOS_ERR, "%s: Base address needs to be aligned to 2\n",
		       dev_path(dev));
	else
		dev->ops = &ops;
}

struct chip_operations drivers_ipmi_ops = {
	CHIP_NAME("IPMI KCS")
	.enable_dev = enable_dev,
};
