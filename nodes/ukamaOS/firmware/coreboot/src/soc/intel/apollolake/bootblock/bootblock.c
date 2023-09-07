/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2016-2018 Intel Corp.
 * (Written by Andrey Petrov <andrey.petrov@intel.com> for Intel Corp.)
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

#include <bootblock_common.h>
#include <cpu/x86/pae.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <intelblocks/cpulib.h>
#include <intelblocks/fast_spi.h>
#include <intelblocks/lpc_lib.h>
#include <intelblocks/p2sb.h>
#include <intelblocks/pcr.h>
#include <intelblocks/rtc.h>
#include <intelblocks/systemagent.h>
#include <intelblocks/pmclib.h>
#include <intelblocks/tco.h>
#include <intelblocks/uart.h>
#include <soc/iomap.h>
#include <soc/cpu.h>
#include <soc/gpio.h>
#include <soc/systemagent.h>
#include <soc/pci_devs.h>
#include <soc/pm.h>
#include <spi-generic.h>

static const struct pad_config tpm_spi_configs[] = {
#if CONFIG(SOC_INTEL_GLK)
	PAD_CFG_NF(GPIO_81, NATIVE, DEEP, NF3),	/* FST_SPI_CS2_N */
#else
	PAD_CFG_NF(GPIO_106, NATIVE, DEEP, NF3),	/* FST_SPI_CS2_N */
#endif
};

static void tpm_enable(void)
{
	/* Configure gpios */
	gpio_configure_pads(tpm_spi_configs, ARRAY_SIZE(tpm_spi_configs));
}

asmlinkage void bootblock_c_entry(uint64_t base_timestamp)
{
	pci_devfn_t dev;

	bootblock_systemagent_early_init();

	p2sb_enable_bar();
	p2sb_configure_hpet();

	/* Decode the ACPI I/O port range for early firmware verification.*/
	dev = PCH_DEV_PMC;
	pci_write_config16(dev, PCI_BASE_ADDRESS_4, ACPI_BASE_ADDRESS);
	pci_write_config16(dev, PCI_COMMAND,
				PCI_COMMAND_IO | PCI_COMMAND_MASTER);

	enable_rtc_upper_bank();

	/* Call lib/bootblock.c main */
	bootblock_main_with_basetime(base_timestamp);
}

static void enable_pmcbar(void)
{
	pci_devfn_t pmc = PCH_DEV_PMC;

	/* Set PMC base addresses and enable decoding. */
	pci_write_config32(pmc, PCI_BASE_ADDRESS_0, PMC_BAR0);
	pci_write_config32(pmc, PCI_BASE_ADDRESS_1, 0);	/* 64-bit BAR */
	pci_write_config32(pmc, PCI_BASE_ADDRESS_2, PMC_BAR1);
	pci_write_config32(pmc, PCI_BASE_ADDRESS_3, 0);	/* 64-bit BAR */
	pci_write_config16(pmc, PCI_BASE_ADDRESS_4, ACPI_BASE_ADDRESS);
	pci_write_config16(pmc, PCI_COMMAND,
				PCI_COMMAND_IO | PCI_COMMAND_MEMORY |
				PCI_COMMAND_MASTER);
}

void bootblock_soc_early_init(void)
{
	enable_pmcbar();

	/* Clear global reset promotion bit */
	pmc_global_reset_enable(0);

	/* Prepare UART for serial console. */
	if (CONFIG(INTEL_LPSS_UART_FOR_CONSOLE))
		uart_bootblock_init();
	if (CONFIG(DRIVERS_UART_8250IO))
		lpc_io_setup_comm_a_b();

	if (CONFIG(TPM_ON_FAST_SPI))
		tpm_enable();

	enable_pm_timer_emulation();

	fast_spi_early_init(SPI_BASE_ADDRESS);

	fast_spi_cache_bios_region();

	/* Initialize GPE for use as interrupt status */
	pmc_gpe_init();

	/* Program TCO Timer Halt */
	tco_configure();

	/* Use Nx and paging to prevent the frontend from writing back dirty
	 * cache-as-ram lines to backing store that doesn't exist when the L1I
	 * speculatively fetches a line that is sitting in the L1D. */
	if (CONFIG(PAGING_IN_CACHE_AS_RAM)) {
		paging_set_nxe(1);
		paging_set_default_pat();
		paging_enable_for_car("pdpt", "pt");
	}
}
