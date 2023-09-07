/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2013 Google Inc.
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

#include <stddef.h>
#include <arch/io.h>
#include <arch/romstage.h>
#include <device/mmio.h>
#include <device/pci_ops.h>
#include <console/console.h>
#include <cbmem.h>
#if CONFIG(EC_GOOGLE_CHROMEEC)
#include <ec/google/chromeec/ec.h>
#endif
#include <elog.h>
#include <romstage_handoff.h>
#include <string.h>
#include <timestamp.h>
#include <vendorcode/google/chromeos/chromeos.h>
#include <soc/gpio.h>
#include <soc/iomap.h>
#include <soc/lpc.h>
#include <soc/msr.h>
#include <soc/pci_devs.h>
#include <soc/pmc.h>
#include <soc/romstage.h>
#include <soc/spi.h>

static void program_base_addresses(void)
{
	uint32_t reg;
	const uint32_t lpc_dev = PCI_DEV(0, LPC_DEV, LPC_FUNC);

	/* Memory Mapped IO registers. */
	reg = PMC_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, PBASE, reg);
	reg = IO_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, IOBASE, reg);
	reg = ILB_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, IBASE, reg);
	reg = SPI_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, SBASE, reg);
	reg = MPHY_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, MPBASE, reg);
	reg = PUNIT_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, PUBASE, reg);
	reg = RCBA_BASE_ADDRESS | 1;
	pci_write_config32(lpc_dev, RCBA, reg);

	/* IO Port Registers. */
	reg = ACPI_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, ABASE, reg);
	reg = GPIO_BASE_ADDRESS | 2;
	pci_write_config32(lpc_dev, GBASE, reg);
}

static void spi_init(void)
{
	u32 *scs = (u32 *)(SPI_BASE_ADDRESS + SCS);
	u32 *bcr = (u32 *)(SPI_BASE_ADDRESS + BCR);
	uint32_t reg;

	/* Disable generating SMI when setting WPD bit. */
	write32(scs, read32(scs) & ~SMIWPEN);
	/*
	 * Enable caching and prefetching in the SPI controller. Disable
	 * the SMM-only BIOS write and set WPD bit.
	 */
	reg = (read32(bcr) & ~SRC_MASK) | SRC_CACHE_PREFETCH | BCR_WPD;
	reg &= ~EISS;
	write32(bcr, reg);
}

static struct chipset_power_state power_state;

static void migrate_power_state(int is_recovery)
{
	struct chipset_power_state *ps_cbmem;
	struct chipset_power_state *ps_car;

	ps_car = &power_state;
	ps_cbmem = cbmem_add(CBMEM_ID_POWER_STATE, sizeof(*ps_cbmem));

	if (ps_cbmem == NULL) {
		printk(BIOS_DEBUG, "Not adding power state to cbmem!\n");
		return;
	}
	memcpy(ps_cbmem, ps_car, sizeof(*ps_cbmem));
}
ROMSTAGE_CBMEM_INIT_HOOK(migrate_power_state)

static struct chipset_power_state *fill_power_state(void)
{
	struct chipset_power_state *ps = &power_state;

	ps->pm1_sts = inw(ACPI_BASE_ADDRESS + PM1_STS);
	ps->pm1_en = inw(ACPI_BASE_ADDRESS + PM1_EN);
	ps->pm1_cnt = inl(ACPI_BASE_ADDRESS + PM1_CNT);
	ps->gpe0_sts = inl(ACPI_BASE_ADDRESS + GPE0_STS);
	ps->gpe0_en = inl(ACPI_BASE_ADDRESS + GPE0_EN);
	ps->tco_sts = inl(ACPI_BASE_ADDRESS + TCO_STS);
	ps->prsts = read32((u32 *)(PMC_BASE_ADDRESS + PRSTS));
	ps->gen_pmcon1 = read32((u32 *)(PMC_BASE_ADDRESS + GEN_PMCON1));
	ps->gen_pmcon2 = read32((u32 *)(PMC_BASE_ADDRESS + GEN_PMCON2));

	printk(BIOS_DEBUG, "pm1_sts: %04x pm1_en: %04x pm1_cnt: %08x\n",
		ps->pm1_sts, ps->pm1_en, ps->pm1_cnt);
	printk(BIOS_DEBUG, "gpe0_sts: %08x gpe0_en: %08x tco_sts: %08x\n",
		ps->gpe0_sts, ps->gpe0_en, ps->tco_sts);
	printk(BIOS_DEBUG, "prsts: %08x gen_pmcon1: %08x gen_pmcon2: %08x\n",
		ps->prsts, ps->gen_pmcon1, ps->gen_pmcon2);

	return ps;
}

/* Return 0, 3, or 5 to indicate the previous sleep state. */
static int chipset_prev_sleep_state(struct chipset_power_state *ps)
{
	/* Default to S0. */
	int prev_sleep_state = ACPI_S0;

	if (ps->pm1_sts & WAK_STS) {
		switch (acpi_sleep_from_pm1(ps->pm1_cnt)) {
		case ACPI_S3:
			if (CONFIG(HAVE_ACPI_RESUME))
				prev_sleep_state = ACPI_S3;
			break;
		case ACPI_S5:
			prev_sleep_state = ACPI_S5;
			break;
		}
		/* Clear SLP_TYP. */
		outl(ps->pm1_cnt & ~(SLP_TYP), ACPI_BASE_ADDRESS + PM1_CNT);
	}

	if (ps->gen_pmcon1 & (PWR_FLR | SUS_PWR_FLR)) {
		prev_sleep_state = ACPI_S5;
	}

	return prev_sleep_state;
}

/* Entry from cpu/intel/car/romstage.c */
void mainboard_romstage_entry(void)
{
	struct chipset_power_state *ps;
	int prev_sleep_state;
	struct mrc_params mp;

	program_base_addresses();

	tco_disable();

	if (CONFIG(ENABLE_BUILTIN_COM1))
		byt_config_com1_and_enable();

	console_init();

	spi_init();

	set_max_freq();

	punit_init();

	gfx_init();

	memset(&mp, 0, sizeof(mp));
	mainboard_fill_mrc_params(&mp);

	timestamp_add_now(TS_BEFORE_INITRAM);

	ps = fill_power_state();
	prev_sleep_state = chipset_prev_sleep_state(ps);

	printk(BIOS_DEBUG, "prev_sleep_state = S%d\n", prev_sleep_state);

	elog_boot_notify(prev_sleep_state == ACPI_S3);

	/* Initialize RAM */
	raminit(&mp, prev_sleep_state);

	timestamp_add_now(TS_AFTER_INITRAM);

	romstage_handoff_init(prev_sleep_state == ACPI_S3);
}
