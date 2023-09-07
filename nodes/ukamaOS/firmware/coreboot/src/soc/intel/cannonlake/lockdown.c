/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 Intel Corporation.
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

#include <device/mmio.h>
#include <intelblocks/cfg.h>
#include <intelpch/lockdown.h>
#include <soc/pm.h>

static void pmc_lock_pmsync(void)
{
	uint8_t *pmcbase;
	uint32_t pmsyncreg;

	pmcbase = pmc_mmio_regs();

	pmsyncreg = read32(pmcbase + PMSYNC_TPR_CFG);
	pmsyncreg |= PCH2CPU_TPR_CFG_LOCK;
	write32(pmcbase + PMSYNC_TPR_CFG, pmsyncreg);
}

static void pmc_lock_abase(void)
{
	uint8_t *pmcbase;
	uint32_t reg32;

	pmcbase = pmc_mmio_regs();

	reg32 = read32(pmcbase + GEN_PMCON_B);
	reg32 |= (SLP_STR_POL_LOCK | ACPI_BASE_LOCK);
	write32(pmcbase + GEN_PMCON_B, reg32);
}

static void pmc_lock_smi(void)
{
	uint8_t *pmcbase;
	uint8_t reg8;

	pmcbase = pmc_mmio_regs();

	reg8 = read8(pmcbase + GEN_PMCON_B);
	reg8 |= SMI_LOCK;
	write8(pmcbase + GEN_PMCON_B, reg8);
}

static void pmc_lockdown_cfg(int chipset_lockdown)
{
	/* PMSYNC */
	pmc_lock_pmsync();
	/* Lock down ABASE and sleep stretching policy */
	pmc_lock_abase();

	if (chipset_lockdown == CHIPSET_LOCKDOWN_COREBOOT)
		pmc_lock_smi();
}

void soc_lockdown_config(int chipset_lockdown)
{
	/* PMC lock down configuration */
	pmc_lockdown_cfg(chipset_lockdown);
}
