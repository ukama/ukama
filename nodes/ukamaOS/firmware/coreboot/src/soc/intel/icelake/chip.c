/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corp.
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

#include <device/device.h>
#include <device/pci.h>
#include <fsp/api.h>
#include <fsp/util.h>
#include <intelblocks/acpi.h>
#include <intelblocks/cfg.h>
#include <intelblocks/itss.h>
#include <intelblocks/xdci.h>
#include <romstage_handoff.h>
#include <soc/intel/common/vbt.h>
#include <soc/itss.h>
#include <soc/pci_devs.h>
#include <soc/ramstage.h>
#include <soc/soc_chip.h>

#if CONFIG(HAVE_ACPI_TABLES)
const char *soc_acpi_name(const struct device *dev)
{
	if (dev->path.type == DEVICE_PATH_DOMAIN)
		return "PCI0";

	if (dev->path.type != DEVICE_PATH_PCI)
		return NULL;

	switch (dev->path.pci.devfn) {
	case SA_DEVFN_ROOT:	return "MCHC";
	case SA_DEVFN_IGD:	return "GFX0";
	case PCH_DEVFN_ISH:	return "ISHB";
	case PCH_DEVFN_XHCI:	return "XHCI";
	case PCH_DEVFN_USBOTG:	return "XDCI";
	case PCH_DEVFN_THERMAL:	return "THRM";
	case PCH_DEVFN_I2C0:	return "I2C0";
	case PCH_DEVFN_I2C1:	return "I2C1";
	case PCH_DEVFN_I2C2:	return "I2C2";
	case PCH_DEVFN_I2C3:	return "I2C3";
	case PCH_DEVFN_CSE:	return "CSE1";
	case PCH_DEVFN_CSE_2:	return "CSE2";
	case PCH_DEVFN_CSE_IDER:	return "CSED";
	case PCH_DEVFN_CSE_KT:	return "CSKT";
	case PCH_DEVFN_CSE_3:	return "CSE3";
	case PCH_DEVFN_SATA:	return "SATA";
	case PCH_DEVFN_UART2:	return "UAR2";
	case PCH_DEVFN_I2C4:	return "I2C4";
	case PCH_DEVFN_I2C5:	return "I2C5";
	case PCH_DEVFN_PCIE1:	return "RP01";
	case PCH_DEVFN_PCIE2:	return "RP02";
	case PCH_DEVFN_PCIE3:	return "RP03";
	case PCH_DEVFN_PCIE4:	return "RP04";
	case PCH_DEVFN_PCIE5:	return "RP05";
	case PCH_DEVFN_PCIE6:	return "RP06";
	case PCH_DEVFN_PCIE7:	return "RP07";
	case PCH_DEVFN_PCIE8:	return "RP08";
	case PCH_DEVFN_PCIE9:	return "RP09";
	case PCH_DEVFN_PCIE10:	return "RP10";
	case PCH_DEVFN_PCIE11:	return "RP11";
	case PCH_DEVFN_PCIE12:	return "RP12";
	case PCH_DEVFN_PCIE13:	return "RP13";
	case PCH_DEVFN_PCIE14:	return "RP14";
	case PCH_DEVFN_PCIE15:	return "RP15";
	case PCH_DEVFN_PCIE16:	return "RP16";
	case PCH_DEVFN_PCIE17:	return "RP17";
	case PCH_DEVFN_PCIE18:	return "RP18";
	case PCH_DEVFN_PCIE19:	return "RP19";
	case PCH_DEVFN_PCIE20:	return "RP20";
	case PCH_DEVFN_PCIE21:	return "RP21";
	case PCH_DEVFN_PCIE22:	return "RP22";
	case PCH_DEVFN_PCIE23:	return "RP23";
	case PCH_DEVFN_PCIE24:	return "RP24";
	case PCH_DEVFN_UART0:	return "UAR0";
	case PCH_DEVFN_UART1:	return "UAR1";
	case PCH_DEVFN_GSPI0:	return "SPI0";
	case PCH_DEVFN_GSPI1:	return "SPI1";
	case PCH_DEVFN_GSPI2:	return "SPI2";
	case PCH_DEVFN_EMMC:	return "EMMC";
	case PCH_DEVFN_SDCARD:	return "SDXC";
	/* Keeping ACPI device name coherent with ec.asl */
	case PCH_DEVFN_ESPI:	return "LPCB";
	case PCH_DEVFN_P2SB:	return "P2SB";
	case PCH_DEVFN_PMC:	return "PMC_";
	case PCH_DEVFN_HDA:	return "HDAS";
	case PCH_DEVFN_SMBUS:	return "SBUS";
	case PCH_DEVFN_SPI:	return "FSPI";
	case PCH_DEVFN_GBE:	return "IGBE";
	case PCH_DEVFN_TRACEHUB:return "THUB";
	}

	return NULL;
}
#endif

/* SoC rotine to fill GPIO PM mask and value for GPIO_MISCCFG register */
static void soc_fill_gpio_pm_configuration(void)
{
	uint8_t value[TOTAL_GPIO_COMM];
	const config_t *config = config_of_soc();

	if (config->gpio_override_pm)
		memcpy(value, config->gpio_pm, sizeof(uint8_t) *
			TOTAL_GPIO_COMM);
	else
		memset(value, MISCCFG_ENABLE_GPIO_PM_CONFIG, sizeof(uint8_t) *
			TOTAL_GPIO_COMM);

	gpio_pm_configure(value, TOTAL_GPIO_COMM);
}

void soc_init_pre_device(void *chip_info)
{
	/* Snapshot the current GPIO IRQ polarities. FSP is setting a
	 * default policy that doesn't honor boards' requirements. */
	itss_snapshot_irq_polarities(GPIO_IRQ_START, GPIO_IRQ_END);

	/* Perform silicon specific init. */
	fsp_silicon_init(romstage_handoff_is_resume());

	 /* Display FIRMWARE_VERSION_INFO_HOB */
	fsp_display_fvi_version_hob();

	/* Restore GPIO IRQ polarities back to previous settings. */
	itss_restore_irq_polarities(GPIO_IRQ_START, GPIO_IRQ_END);

	soc_fill_gpio_pm_configuration();
}

static void pci_domain_set_resources(struct device *dev)
{
	assign_resources(dev->link_list);
}

static struct device_operations pci_domain_ops = {
	.read_resources   = &pci_domain_read_resources,
	.set_resources    = &pci_domain_set_resources,
	.scan_bus         = &pci_domain_scan_bus,
	#if CONFIG(HAVE_ACPI_TABLES)
	.acpi_name        = &soc_acpi_name,
	#endif
};

static struct device_operations cpu_bus_ops = {
	.read_resources   = DEVICE_NOOP,
	.set_resources    = DEVICE_NOOP,
	.enable_resources = DEVICE_NOOP,
	.init             = DEVICE_NOOP,
	.acpi_fill_ssdt_generator = generate_cpu_entries,
};

static void soc_enable(struct device *dev)
{
	/* Set the operations if it is a special bus type */
	if (dev->path.type == DEVICE_PATH_DOMAIN)
		dev->ops = &pci_domain_ops;
	else if (dev->path.type == DEVICE_PATH_CPU_CLUSTER)
		dev->ops = &cpu_bus_ops;
}

struct chip_operations soc_intel_icelake_ops = {
	CHIP_NAME("Intel Icelake")
	.enable_dev	= &soc_enable,
	.init		= &soc_init_pre_device,
};
