/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 or (at your option)
 * any later version of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <device/device.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <device/pci_ids.h>
#include <elog.h>
#include <smbios.h>
#include <string.h>
#include "chip.h"
#include "drivers/wifi/generic_wifi.h"

#define PMCS_DR 0xcc
#define PME_STS (1 << 15)

#if CONFIG(GENERATE_SMBIOS_TABLES)
static int smbios_write_wifi(struct device *dev, int *handle,
			     unsigned long *current)
{
	struct smbios_type_intel_wifi {
		u8 type;
		u8 length;
		u16 handle;
		u8 str;
		u8 eos[2];
	} __packed;

	struct smbios_type_intel_wifi *t =
		(struct smbios_type_intel_wifi *)*current;
	int len = sizeof(struct smbios_type_intel_wifi);

	memset(t, 0, sizeof(struct smbios_type_intel_wifi));
	t->type = 0x85;
	t->length = len - 2;
	t->handle = *handle;
	/*
	 * Intel wifi driver expects this string to be in the table 0x85
	 * with PCI IDs enumerated below.
	 */
	t->str = smbios_add_string(t->eos, "KHOIHGIUCCHHII");

	len = t->length + smbios_string_table_len(t->eos);
	*current += len;
	*handle += 1;
	return len;
}
#endif

#if CONFIG(HAVE_ACPI_TABLES)
static void intel_wifi_fill_ssdt(struct device *dev)
{
	struct drivers_intel_wifi_config *config = dev->chip_info;
	struct generic_wifi_config generic_config;

	if (config) {
		generic_config.wake = config->wake;
		/* By default, all intel wifi chips wake from S3 */
		generic_config.maxsleep = 3;
	}
	generic_wifi_fill_ssdt(dev, config ? &generic_config : NULL);
}
#endif

static void wifi_pci_dev_init(struct device *dev)
{
	pci_dev_init(dev);

	if (CONFIG(ELOG)) {
		uint32_t val;
		val = pci_read_config16(dev, PMCS_DR);
		if (val & PME_STS)
			elog_add_event_wake(ELOG_WAKE_SOURCE_PME_WIFI, 0);
	}
}

static struct pci_operations pci_ops = {
	.set_subsystem = pci_dev_set_subsystem,
};

struct device_operations device_ops = {
	.read_resources           = pci_dev_read_resources,
	.set_resources            = pci_dev_set_resources,
	.enable_resources         = pci_dev_enable_resources,
	.init                     = wifi_pci_dev_init,
#if CONFIG(GENERATE_SMBIOS_TABLES)
	.get_smbios_data          = smbios_write_wifi,
#endif
	.ops_pci                  = &pci_ops,
#if CONFIG(HAVE_ACPI_TABLES)
	.acpi_name                = generic_wifi_acpi_name,
	.acpi_fill_ssdt_generator = intel_wifi_fill_ssdt,
#endif
};

static const unsigned short pci_device_ids[] = {
	PCI_DEVICE_ID_1000_SERIES_WIFI,
	PCI_DEVICE_ID_6005_SERIES_WIFI,
	PCI_DEVICE_ID_6005_I_SERIES_WIFI,
	PCI_DEVICE_ID_1030_SERIES_WIFI,
	PCI_DEVICE_ID_6030_I_SERIES_WIFI,
	PCI_DEVICE_ID_6030_SERIES_WIFI,
	PCI_DEVICE_ID_6150_SERIES_WIFI,
	PCI_DEVICE_ID_2030_SERIES_WIFI,
	PCI_DEVICE_ID_2000_SERIES_WIFI,
	PCI_DEVICE_ID_0135_SERIES_WIFI,
	PCI_DEVICE_ID_0105_SERIES_WIFI,
	PCI_DEVICE_ID_6035_SERIES_WIFI,
	PCI_DEVICE_ID_5300_SERIES_WIFI,
	PCI_DEVICE_ID_5100_SERIES_WIFI,
	PCI_DEVICE_ID_6000_SERIES_WIFI,
	PCI_DEVICE_ID_6000_I_SERIES_WIFI,
	PCI_DEVICE_ID_5350_SERIES_WIFI,
	PCI_DEVICE_ID_5150_SERIES_WIFI,
	/* Wilkins Peak 2 */
	PCI_DEVICE_ID_WP_7260_SERIES_1_WIFI,
	PCI_DEVICE_ID_WP_7260_SERIES_2_WIFI,
	/* Stone Peak 2 */
	PCI_DEVICE_ID_SP_7265_SERIES_1_WIFI,
	PCI_DEVICE_ID_SP_7265_SERIES_2_WIFI,
	/* Stone Field Peak */
	PCI_DEVICE_ID_SFP_8260_SERIES_1_WIFI,
	PCI_DEVICE_ID_SFP_8260_SERIES_2_WIFI,
	/* Windstorm Peak */
	PCI_DEVICE_ID_WSP_8275_SERIES_1_WIFI,
	/* Jefferson Peak */
	PCI_DEVICE_ID_JP_9000_SERIES_1_WIFI,
	PCI_DEVICE_ID_JP_9000_SERIES_2_WIFI,
	PCI_DEVICE_ID_JP_9000_SERIES_3_WIFI,
	/* Thunder Peak 2 */
	PCI_DEVICE_ID_TP_9260_SERIES_WIFI,
	/* Harrison Peak */
	PCI_DEVICE_ID_HrP_9560_SERIES_1_WIFI,
	PCI_DEVICE_ID_HrP_9560_SERIES_2_WIFI,
	PCI_DEVICE_ID_HrP_9560_SERIES_3_WIFI,
	PCI_DEVICE_ID_HrP_9560_SERIES_4_WIFI,
	0
};

static const struct pci_driver pch_intel_wifi __pci_driver = {
	.ops	 = &device_ops,
	.vendor	 = PCI_VENDOR_ID_INTEL,
	.devices = pci_device_ids,
};

static void intel_wifi_enable(struct device *dev)
{
	dev->ops = &device_ops;
}

struct chip_operations drivers_intel_wifi_ops = {
	CHIP_NAME("Intel WiFi")
	.enable_dev = intel_wifi_enable
};
