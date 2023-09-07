/*
 * This file is part of the coreboot project.
 *
 * Copyright 2015 Intel Corporation.
 * Copyright 2019 Google LLC
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
#include <device/pci_ops.h>
#include <stdint.h>
#include <elog.h>
#include <intelblocks/pmclib.h>
#include <intelblocks/xhci.h>
#include <soc/pci_devs.h>
#include <soc/pm.h>

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

static void pch_log_add_elog_event(const struct pme_status_info *info)
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
		/*
		 * The power management control/status register is not
		 * listed in the cannonlake PCH EDS. We have been told
		 * that the PMCS register is at offset 0xCC.
		 */
		{ PCH_DEV_CNViWIFI, 0xcc, ELOG_WAKE_SOURCE_PME_WIFI },
	};

	for (i = 0; i < ARRAY_SIZE(pme_status_info); i++) {
		dev = pme_status_info[i].dev;
		if (!dev)
			continue;

		val = pci_read_config16(dev, pme_status_info[i].reg_offset);

		if ((val == 0xFFFF) || !(val & PME_STS_BIT))
			continue;

		pch_log_add_elog_event(&pme_status_info[i]);
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

static void pch_log_gpio_gpe(u32 gpe0_sts, u32 gpe0_en, int start)
{
	int i;

	gpe0_sts &= gpe0_en;

	for (i = 0; i <= 31; i++) {
		if (gpe0_sts & (1 << i))
			elog_add_event_wake(ELOG_WAKE_SOURCE_GPIO, i + start);
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

	/* PME (TODO: determine wake device) */
	if (ps->gpe0_sts[GPE_STD] & PME_STS)
		elog_add_event_wake(ELOG_WAKE_SOURCE_PME, 0);

	/* XHCI - "Power Management Event Bus 0" events include XHCI */
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
	/* Thermal Trip */
	if (ps->gblrst_cause[0] & GBLRST_CAUSE0_THERMTRIP)
		elog_add_event(ELOG_TYPE_THERM_TRIP);

	/* PWR_FLR Power Failure */
	if (ps->gen_pmcon_a & PWR_FLR)
		elog_add_event(ELOG_TYPE_POWER_FAIL);

	/* SUS Well Power Failure */
	if (ps->gen_pmcon_a & SUS_PWR_FLR)
		elog_add_event(ELOG_TYPE_SUS_POWER_FAIL);

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
	if (ps->gen_pmcon_a & HOST_RST_STS)
		elog_add_event(ELOG_TYPE_SYSTEM_RESET);

	/* ACPI Wake Event */
	if (ps->prev_sleep_state != ACPI_S0)
		elog_add_event_byte(ELOG_TYPE_ACPI_WAKE, ps->prev_sleep_state);
}

static void pch_log_state(void *unused)
{
	struct chipset_power_state *ps = pmc_get_power_state();

	if (!ps) {
		printk(BIOS_ERR, "chipset_power_state not found!\n");
		return;
	}

	/* Power and Reset */
	pch_log_power_and_resets(ps);

	/* Wake Sources */
	if (ps->prev_sleep_state > ACPI_S0)
		pch_log_wake_source(ps);
}

BOOT_STATE_INIT_ENTRY(BS_DEV_INIT, BS_ON_EXIT, pch_log_state, NULL);

void elog_gsmi_cb_platform_log_wake_source(void)
{
	struct chipset_power_state ps;
	pmc_fill_pm_reg_info(&ps);
	pch_log_wake_source(&ps);
}
