/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2015 Intel Corporation.
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

#include <bootstate.h>
#include <cbmem.h>
#include <console/console.h>
#include <device/mmio.h>
#include <device/pci_ops.h>
#include <stdint.h>
#include <elog.h>
#include <intelblocks/pmclib.h>
#include <intelblocks/xhci.h>
#include <soc/pci_devs.h>
#include <soc/pm.h>
#include <soc/smbus.h>

static void pch_log_gpio_gpe(u32 gpe0_sts, u32 gpe0_en, int start)
{
	int i;

	gpe0_sts &= gpe0_en;

	for (i = 0; i <= 31; i++) {
		if (gpe0_sts & (1 << i))
			elog_add_event_wake(ELOG_WAKE_SOURCE_GPIO, i + start);
	}
}

struct pme_status_info {
#ifdef __SIMPLE_DEVICE__
	pci_devfn_t dev;
#else
	struct device *dev;
#endif
	uint8_t reg_offset;
	uint32_t elog_event;
};

#define PME_STS_BIT		(1 << 15)

#ifdef __SIMPLE_DEVICE__
static void pch_log_add_elog_event(const struct pme_status_info *info,
				   pci_devfn_t dev)
#else
static void pch_log_add_elog_event(const struct pme_status_info *info,
				   struct device *dev)
#endif
{
	/*
	 * If wake source is XHCI, check for detailed wake source events on
	 * USB2/3 ports.
	 */
	if ((info->dev == PCH_DEV_XHCI) &&
			pch_xhci_update_wake_event(soc_get_xhci_usb_info()))
		return;

	elog_add_event_wake(info->elog_event, 0);
}

static void pch_log_pme_internal_wake_source(void)
{
	size_t i;
#ifdef __SIMPLE_DEVICE__
	pci_devfn_t dev;
#else
	struct device *dev;
#endif
	uint16_t val;
	bool dev_found = false;

	struct pme_status_info pme_status_info[] = {
		{ PCH_DEV_HDA, 0x54, ELOG_WAKE_SOURCE_PME_HDA },
		{ PCH_DEV_GBE, 0xcc, ELOG_WAKE_SOURCE_PME_GBE },
		{ PCH_DEV_SATA, 0x74, ELOG_WAKE_SOURCE_PME_SATA },
		{ PCH_DEV_CSE, 0x54, ELOG_WAKE_SOURCE_PME_CSE },
		{ PCH_DEV_XHCI, 0x74, ELOG_WAKE_SOURCE_PME_XHCI },
		{ PCH_DEV_USBOTG, 0x84, ELOG_WAKE_SOURCE_PME_XDCI },
	};

	for (i = 0; i < ARRAY_SIZE(pme_status_info); i++) {
		dev = pme_status_info[i].dev;
		if (!dev)
			continue;

		val = pci_read_config16(dev, pme_status_info[i].reg_offset);

		if ((val == 0xFFFF) || !(val & PME_STS_BIT))
			continue;

		pch_log_add_elog_event(&pme_status_info[i], dev);
		dev_found = true;
	}

	/*
	 * If device is still not found, but the wake source is internal PME,
	 * try probing XHCI ports to see if any of the USB2/3 ports indicate
	 * that it was the wake source. This path would be taken in case of GSMI
	 * logging with S0ix where the pci_pm_resume_noirq runs and clears the
	 * PME_STS_BIT in controller register.
	 */
	if (!dev_found)
		dev_found = pch_xhci_update_wake_event(soc_get_xhci_usb_info());

	if (!dev_found)
		elog_add_event_wake(ELOG_WAKE_SOURCE_PME_INTERNAL, 0);
}

#define RP_PME_STS_BIT		(1 << 16)
static void pch_log_rp_wake_source(void)
{
	size_t i, maxports;
#ifdef __SIMPLE_DEVICE__
	pci_devfn_t dev;
#else
	struct device *dev;
#endif
	uint32_t val;

	struct pme_status_info pme_status_info[] = {
		{ PCH_DEV_PCIE1, 0x60, ELOG_WAKE_SOURCE_PME_PCIE1 },
		{ PCH_DEV_PCIE2, 0x60, ELOG_WAKE_SOURCE_PME_PCIE2 },
		{ PCH_DEV_PCIE3, 0x60, ELOG_WAKE_SOURCE_PME_PCIE3 },
		{ PCH_DEV_PCIE4, 0x60, ELOG_WAKE_SOURCE_PME_PCIE4 },
		{ PCH_DEV_PCIE5, 0x60, ELOG_WAKE_SOURCE_PME_PCIE5 },
		{ PCH_DEV_PCIE6, 0x60, ELOG_WAKE_SOURCE_PME_PCIE6 },
		{ PCH_DEV_PCIE7, 0x60, ELOG_WAKE_SOURCE_PME_PCIE7 },
		{ PCH_DEV_PCIE8, 0x60, ELOG_WAKE_SOURCE_PME_PCIE8 },
		{ PCH_DEV_PCIE9, 0x60, ELOG_WAKE_SOURCE_PME_PCIE9 },
		{ PCH_DEV_PCIE10, 0x60, ELOG_WAKE_SOURCE_PME_PCIE10 },
		{ PCH_DEV_PCIE11, 0x60, ELOG_WAKE_SOURCE_PME_PCIE11 },
		{ PCH_DEV_PCIE12, 0x60, ELOG_WAKE_SOURCE_PME_PCIE12 },
		{ PCH_DEV_PCIE13, 0x60, ELOG_WAKE_SOURCE_PME_PCIE13 },
		{ PCH_DEV_PCIE14, 0x60, ELOG_WAKE_SOURCE_PME_PCIE14 },
		{ PCH_DEV_PCIE15, 0x60, ELOG_WAKE_SOURCE_PME_PCIE15 },
		{ PCH_DEV_PCIE16, 0x60, ELOG_WAKE_SOURCE_PME_PCIE16 },
		{ PCH_DEV_PCIE17, 0x60, ELOG_WAKE_SOURCE_PME_PCIE17 },
		{ PCH_DEV_PCIE18, 0x60, ELOG_WAKE_SOURCE_PME_PCIE18 },
		{ PCH_DEV_PCIE19, 0x60, ELOG_WAKE_SOURCE_PME_PCIE19 },
		{ PCH_DEV_PCIE20, 0x60, ELOG_WAKE_SOURCE_PME_PCIE20 },
		{ PCH_DEV_PCIE21, 0x60, ELOG_WAKE_SOURCE_PME_PCIE21 },
		{ PCH_DEV_PCIE22, 0x60, ELOG_WAKE_SOURCE_PME_PCIE22 },
		{ PCH_DEV_PCIE23, 0x60, ELOG_WAKE_SOURCE_PME_PCIE23 },
		{ PCH_DEV_PCIE24, 0x60, ELOG_WAKE_SOURCE_PME_PCIE24 },
	};

	maxports = min(CONFIG_MAX_ROOT_PORTS, ARRAY_SIZE(pme_status_info));

	for (i = 0; i < maxports; i++) {
		dev = pme_status_info[i].dev;

		if (!dev)
			continue;

		val = pci_read_config32(dev, pme_status_info[i].reg_offset);

		if ((val == 0xFFFFFFFF) || !(val & RP_PME_STS_BIT))
			continue;

		/*
		 * Linux kernel uses PME STS bit information. So do not clear
		 * this bit.
		 */
		pch_log_add_elog_event(&pme_status_info[i], dev);
	}
}

static void pch_log_wake_source(struct chipset_power_state *ps)
{
	/* Power Button */
	if (ps->pm1_sts & PWRBTN_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_PWRBTN, 0);

	/* RTC */
	if (ps->pm1_sts & RTC_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_RTC, 0);

	/* PCI Express (TODO: determine wake device) */
	if (ps->pm1_sts & PCIEXPWAK_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_PCIE, 0);

	/*
	 * PCIE Root Port .
	 * This should be done when PCIEXPWAK_STS bit is set.
	 * In SPT, this bit isn't getting set due to known bug.
	 * So scan all PCIe RP for PME status bit.
	 */
	pch_log_rp_wake_source();

	/* PME (TODO: determine wake device) */
	if (ps->gpe0_sts[GPE_STD] & PME_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_PME, 0);

	/* Internal PME (TODO: determine wake device) */
	if (ps->gpe0_sts[GPE_STD] & PME_B0_STS)
		pch_log_pme_internal_wake_source();

	/* SMBUS Wake */
	if (ps->gpe0_sts[GPE_STD] & SMB_WAK_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_SMBUS, 0);

	/* Log GPIO events in set 1-3 */
	pch_log_gpio_gpe(ps->gpe0_sts[GPE_31_0], ps->gpe0_en[GPE_31_0], 0);
	pch_log_gpio_gpe(ps->gpe0_sts[GPE_63_32], ps->gpe0_en[GPE_63_32], 32);
	pch_log_gpio_gpe(ps->gpe0_sts[GPE_95_64], ps->gpe0_en[GPE_95_64], 64);
	/* Treat the STD as an extension of GPIO to obtain visibility. */
	pch_log_gpio_gpe(ps->gpe0_sts[GPE_STD], ps->gpe0_en[GPE_STD], 96);
}

static void pch_log_power_and_resets(struct chipset_power_state *ps)
{
	bool deep_sx;

	/*
	 * Platform entered deep Sx if:
	 * 1. Prev sleep state was Sx and deep_sx_enabled() is true
	 * 2. SUS well power was lost
	 */
	deep_sx = ((((ps->prev_sleep_state == ACPI_S3) && deep_s3_enabled()) ||
		   ((ps->prev_sleep_state == ACPI_S5) && deep_s5_enabled())) &&
		   (ps->gen_pmcon_b & SUS_PWR_FLR));

	/* Thermal Trip */
	if (ps->gblrst_cause[0] & GBLRST_CAUSE0_THERMTRIP)
		elog_add_event(ELOG_TYPE_THERM_TRIP);

	/* PWR_FLR Power Failure */
	if (ps->gen_pmcon_b & PWR_FLR)
		elog_add_event(ELOG_TYPE_POWER_FAIL);

	/* SUS Well Power Failure */
	if (ps->gen_pmcon_b & SUS_PWR_FLR) {
		/* Do not log SUS_PWR_FLR if waking from deep Sx */
		if (!deep_sx)
			elog_add_event(ELOG_TYPE_SUS_POWER_FAIL);
	}

	/* TCO Timeout */
	if (ps->prev_sleep_state != ACPI_S3 &&
	    ps->tco2_sts & TCO_STS_SECOND_TO)
		elog_add_event(ELOG_TYPE_TCO_RESET);

	/* Power Button Override */
	if (ps->pm1_sts & PRBTNOR_STS)
		elog_add_event(ELOG_TYPE_POWER_BUTTON_OVERRIDE);

	/* RTC reset */
	if (ps->gen_pmcon_b & RTC_BATTERY_DEAD)
		elog_add_event(ELOG_TYPE_RTC_RESET);

	/* Host Reset Status */
	if (ps->gen_pmcon_b & HOST_RST_STS)
		elog_add_event(ELOG_TYPE_SYSTEM_RESET);

	/* ACPI Wake Event */
	if (ps->prev_sleep_state != ACPI_S0) {
		if (deep_sx)
			elog_add_event_byte(ELOG_TYPE_ACPI_DEEP_WAKE,
					    ps->prev_sleep_state);
		else
			elog_add_event_byte(ELOG_TYPE_ACPI_WAKE,
					    ps->prev_sleep_state);
	}
}

static void pch_log_state(void *unused)
{
	struct chipset_power_state *ps = cbmem_find(CBMEM_ID_POWER_STATE);

	if (ps == NULL) {
		printk(BIOS_ERR,
			"Not logging power state information. Power state not found in cbmem.\n");
		return;
	}

	/* Power and Reset */
	pch_log_power_and_resets(ps);

	/* Wake Sources */
	if (ps->prev_sleep_state > 0)
		pch_log_wake_source(ps);
}

BOOT_STATE_INIT_ENTRY(BS_DEV_INIT, BS_ON_EXIT, pch_log_state, NULL);

void elog_gsmi_cb_platform_log_wake_source(void)
{
	struct chipset_power_state ps;
	pmc_fill_pm_reg_info(&ps);
	pch_log_wake_source(&ps);
}
