/*
 * This file is part of the coreboot project.
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

#include <arch/early_variables.h>
#include <console/console.h>
#include <security/tpm/tis.h>
#include <arch/acpigen.h>
#include <device/device.h>
#include <drivers/intel/ptt/ptt.h>

#include "tpm.h"
#include "chip.h"

static unsigned tpm_is_open CAR_GLOBAL;

static const struct {
	uint16_t vid;
	uint16_t did;
	const char *device_name;
} dev_map[] = {
	{0x1ae0, 0x0028, "CR50"},
	{0xa13a, 0x8086, "Intel iTPM"}
};

static const char *tis_get_dev_name(struct tpm2_info *info)
{
	int i;

	for (i = 0; i < ARRAY_SIZE(dev_map); i++)
		if ((dev_map[i].vid == info->vendor_id) && (dev_map[i].did == info->device_id))
			return dev_map[i].device_name;
	return "Unknown";
}


int tis_open(void)
{
	if (car_get_var(tpm_is_open)) {
		printk(BIOS_ERR, "%s called twice.\n", __func__);
		return -1;
	}

	if (CONFIG(HAVE_INTEL_PTT)) {
		if (!ptt_active()) {
			printk(BIOS_ERR, "%s: Intel PTT is not active.\n", __func__);
			return -1;
		}
		printk(BIOS_DEBUG, "%s: Intel PTT is active.\n", __func__);
	}

	return 0;
}

int tis_close(void)
{
	if (car_get_var(tpm_is_open)) {

		/*
		 * Do we need to do something here, like waiting for a
		 * transaction to stop?
		 */
		car_set_var(tpm_is_open, 0);
	}

	return 0;
}

int tis_init(void)
{
	struct tpm2_info info;

	// Wake TPM up (if necessary)
	if (tpm2_init() != 0)
		return -1;

	tpm2_get_info(&info);

	printk(BIOS_INFO, "Initialized TPM device %s revision %d\n", tis_get_dev_name(&info),
	       info.revision);

	return 0;
}


int tis_sendrecv(const uint8_t *sendbuf, size_t sbuf_size, uint8_t *recvbuf, size_t *rbuf_len)
{
	int len = tpm2_process_command(sendbuf, sbuf_size, recvbuf, *rbuf_len);

	if (len == 0)
		return -1;

	*rbuf_len = len;

	return 0;
}

#ifdef __RAMSTAGE__

static void crb_tpm_fill_ssdt(struct device *dev)
{
	const char *path = acpi_device_path(dev);
	if (!path) {
		path = "\\_SB_.TPM";
		printk(BIOS_DEBUG, "Using default TPM2 ACPI path: '%s'\n", path);
	}

	/* Device */
	acpigen_write_device(path);

	acpigen_write_name_string("_HID", "MSFT0101");
	acpigen_write_name_string("_CID", "MSFT0101");

	acpigen_write_name_integer("_UID", 1);

	acpigen_write_STA(ACPI_STATUS_DEVICE_ALL_ON);

	/* Resources */
	acpigen_write_name("_CRS");
	acpigen_write_resourcetemplate_header();
	acpigen_write_mem32fixed(1, TPM_CRB_BASE_ADDRESS, 0x5000);

	acpigen_write_resourcetemplate_footer();

	acpigen_pop_len(); /* Device */
}

static const char *crb_tpm_acpi_name(const struct device *dev)
{
	return "TPM";
}

static struct device_operations crb_ops = {
	.read_resources = DEVICE_NOOP,
	.set_resources = DEVICE_NOOP,
#if CONFIG(HAVE_ACPI_TABLES)
	.acpi_name = crb_tpm_acpi_name,
	.acpi_fill_ssdt_generator = crb_tpm_fill_ssdt,
#endif

};

static void enable_dev(struct device *dev)
{
	dev->ops = &crb_ops;
}

struct chip_operations drivers_crb_ops = {CHIP_NAME("CRB TPM").enable_dev = enable_dev};

#endif /* __RAMSTAGE__ */
