/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007 - 2009 coresystems GmbH
 * Copyright (C) 2013 Google Inc.
 * Copyright (C) 2014 - 2017 Intel Corporation.
 * Copyright (C) 2018 Online SAS
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 */

#include <arch/acpi.h>
#include <arch/acpigen.h>
#include <arch/cpu.h>
#include <arch/smp/mpspec.h>
#include <cpu/x86/smm.h>
#include <string.h>
#include <device/pci.h>
#include <device/pci_ops.h>
#include <cbmem.h>
#include <console/console.h>
#include <intelblocks/acpi.h>
#include <soc/acpi.h>
#include <soc/cpu.h>
#include <soc/soc_util.h>
#include <soc/pmc.h>
#include <soc/systemagent.h>

#define MWAIT_RES(state, sub_state)                         \
	{                                                   \
		.addrl = (((state) << 4) | (sub_state)),    \
		.space_id = ACPI_ADDRESS_SPACE_FIXED,       \
		.bit_width = ACPI_FFIXEDHW_VENDOR_INTEL,    \
		.bit_offset = ACPI_FFIXEDHW_CLASS_MWAIT,    \
		.access_size = ACPI_FFIXEDHW_FLAG_HW_COORD, \
	}

#define CSTATE_RES(address_space, width, offset, address)		\
	{								\
	.space_id = address_space,					\
	.bit_width = width,						\
	.bit_offset = offset,						\
	.addrl = address,						\
	}

static acpi_cstate_t cstate_map[] = {
	{
		/* C1 */
		.ctype = 1,		/* ACPI C1 */
		.latency = 2,
		.power = 1000,
		.resource = MWAIT_RES(0, 0),
	},
	{
		.ctype = 2,		/* ACPI C2 */
		.latency = 10,
		.power = 10,
		.resource = CSTATE_RES(ACPI_ADDRESS_SPACE_IO, 8, 0,
				       ACPI_BASE_ADDRESS + 0x14),
	},
	{
		.ctype = 3,		/* ACPI C3 */
		.latency = 50,
		.power = 10,
		.resource = CSTATE_RES(ACPI_ADDRESS_SPACE_IO, 8, 0,
				       ACPI_BASE_ADDRESS + 0x15),
	}
};

void acpi_init_gnvs(global_nvs_t *gnvs)
{
	/* CPU core count */
	gnvs->pcnt = dev_count_cpu();

	/* Top of Low Memory (start of resource allocation) */
	gnvs->tolm = (uintptr_t)cbmem_top();

#if CONFIG(CONSOLE_CBMEM)
	/* Update the mem console pointer. */
	gnvs->cbmc = (u32)cbmem_find(CBMEM_ID_CONSOLE);
#endif

	/* MMIO Low/High & TSEG base and length */
	gnvs->mmiob = (u32)get_top_of_low_memory();
	gnvs->mmiol = (u32)(get_pciebase() - 1);
	gnvs->mmiohb = (u64)get_top_of_upper_memory();
	gnvs->mmiohl = (u64)(((u64)1 << CONFIG_CPU_ADDR_BITS) - 1);
	gnvs->tsegb = (u32)get_tseg_memory();
	gnvs->tsegl = (u32)(get_top_of_low_memory() - get_tseg_memory());
}

uint32_t soc_read_sci_irq_select(void)
{
	struct device *dev = get_pmc_dev();

	if (!dev)
		return 0;

	return pci_read_config32(dev, PMC_ACPI_CNT);
}

acpi_cstate_t *soc_get_cstate_map(size_t *entries)
{
	*entries = ARRAY_SIZE(cstate_map);
	return cstate_map;
}

unsigned long acpi_fill_mcfg(unsigned long current)
{
	u32 pciexbar_reg;
	int max_buses;

	pciexbar_reg = get_pciebase();
	max_buses = get_pcielength();

	if (!pciexbar_reg)
		return current;

	current += acpi_create_mcfg_mmconfig((acpi_mcfg_mmconfig_t *)current,
					     pciexbar_reg, 0x0, 0x0,
					     (u8)(max_buses - 1));

	return current;
}

__attribute__ ((weak)) void motherboard_fill_fadt(acpi_fadt_t *fadt)
{
}

void soc_fill_fadt(acpi_fadt_t *fadt)
{
	u16 pmbase = get_pmbase();

	/* System Management */
	if (!CONFIG(HAVE_SMI_HANDLER)) {
		fadt->smi_cmd = 0x00;
		fadt->acpi_enable = 0x00;
		fadt->acpi_disable = 0x00;
	}

	/* Power Control */
	fadt->pm2_cnt_blk = pmbase + PM2_CNT;
	fadt->pm_tmr_blk = pmbase + PM1_TMR;
	fadt->gpe1_blk = 0;

	/* Control Registers - Length */
	fadt->pm2_cnt_len = 1;
	fadt->pm_tmr_len = 4;
	fadt->gpe0_blk_len = 8;
	fadt->gpe1_blk_len = 0;
	fadt->gpe1_base = 0;
	fadt->cst_cnt = 0;
	fadt->p_lvl2_lat = ACPI_FADT_C2_NOT_SUPPORTED;
	fadt->p_lvl3_lat = ACPI_FADT_C3_NOT_SUPPORTED;
	fadt->flush_size = 0;   /* set to 0 if WBINVD is 1 in flags */
	fadt->flush_stride = 0; /* set to 0 if WBINVD is 1 in flags */
	fadt->duty_offset = 1;
	fadt->duty_width = 0;

	/* RTC Registers */
	fadt->day_alrm = 0x0D;
	fadt->mon_alrm = 0x00;
	fadt->century = 0x00;
	fadt->iapc_boot_arch = ACPI_FADT_LEGACY_DEVICES | ACPI_FADT_8042;

	fadt->flags = ACPI_FADT_WBINVD | ACPI_FADT_C1_SUPPORTED |
		      ACPI_FADT_C2_MP_SUPPORTED | ACPI_FADT_SLEEP_BUTTON |
		      ACPI_FADT_RESET_REGISTER | ACPI_FADT_SLEEP_TYPE |
		      ACPI_FADT_S4_RTC_WAKE | ACPI_FADT_PLATFORM_CLOCK;

	/* Reset Register */
	fadt->reset_reg.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->reset_reg.bit_width = 8;
	fadt->reset_reg.bit_offset = 0;
	fadt->reset_reg.access_size = ACPI_ACCESS_SIZE_BYTE_ACCESS;
	fadt->reset_reg.addrl = 0xCF9;
	fadt->reset_reg.addrh = 0x00;
	fadt->reset_value = 6;

	/* PM1 Status & PM1 Enable */
	fadt->x_pm1a_evt_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm1a_evt_blk.bit_width = 32;
	fadt->x_pm1a_evt_blk.bit_offset = 0;
	fadt->x_pm1a_evt_blk.access_size = ACPI_ACCESS_SIZE_DWORD_ACCESS;
	fadt->x_pm1a_evt_blk.addrl = fadt->pm1a_evt_blk;
	fadt->x_pm1a_evt_blk.addrh = 0x00;

	fadt->x_pm1b_evt_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm1b_evt_blk.bit_width = 0;
	fadt->x_pm1b_evt_blk.bit_offset = 0;
	fadt->x_pm1b_evt_blk.access_size = 0;
	fadt->x_pm1b_evt_blk.addrl = fadt->pm1b_evt_blk;
	fadt->x_pm1b_evt_blk.addrh = 0x00;

	/* PM1 Control Registers */
	fadt->x_pm1a_cnt_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm1a_cnt_blk.bit_width = 16;
	fadt->x_pm1a_cnt_blk.bit_offset = 0;
	fadt->x_pm1a_cnt_blk.access_size = ACPI_ACCESS_SIZE_WORD_ACCESS;
	fadt->x_pm1a_cnt_blk.addrl = fadt->pm1a_cnt_blk;
	fadt->x_pm1a_cnt_blk.addrh = 0x00;

	fadt->x_pm1b_cnt_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm1b_cnt_blk.bit_width = 0;
	fadt->x_pm1b_cnt_blk.bit_offset = 0;
	fadt->x_pm1b_cnt_blk.access_size = 0;
	fadt->x_pm1b_cnt_blk.addrl = fadt->pm1b_cnt_blk;
	fadt->x_pm1b_cnt_blk.addrh = 0x00;

	/* PM2 Control Registers */
	fadt->x_pm2_cnt_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm2_cnt_blk.bit_width = 8;
	fadt->x_pm2_cnt_blk.bit_offset = 0;
	fadt->x_pm2_cnt_blk.access_size = ACPI_ACCESS_SIZE_BYTE_ACCESS;
	fadt->x_pm2_cnt_blk.addrl = fadt->pm2_cnt_blk;
	fadt->x_pm2_cnt_blk.addrh = 0x00;

	/* PM1 Timer Register */
	fadt->x_pm_tmr_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_pm_tmr_blk.bit_width = 32;
	fadt->x_pm_tmr_blk.bit_offset = 0;
	fadt->x_pm_tmr_blk.access_size = ACPI_ACCESS_SIZE_DWORD_ACCESS;
	fadt->x_pm_tmr_blk.addrl = fadt->pm_tmr_blk;
	fadt->x_pm_tmr_blk.addrh = 0x00;

	/*  General-Purpose Event Registers */
	fadt->x_gpe0_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_gpe0_blk.bit_width = 64; /* EventStatus + EventEnable */
	fadt->x_gpe0_blk.bit_offset = 0;
	fadt->x_gpe0_blk.access_size = ACPI_ACCESS_SIZE_DWORD_ACCESS;
	fadt->x_gpe0_blk.addrl = fadt->gpe0_blk;
	fadt->x_gpe0_blk.addrh = 0x00;

	fadt->x_gpe1_blk.space_id = ACPI_ADDRESS_SPACE_IO;
	fadt->x_gpe1_blk.bit_width = 0;
	fadt->x_gpe1_blk.bit_offset = 0;
	fadt->x_gpe1_blk.access_size = 0;
	fadt->x_gpe1_blk.addrl = fadt->gpe1_blk;
	fadt->x_gpe1_blk.addrh = 0x00;

	motherboard_fill_fadt(fadt);
}

static acpi_tstate_t denverton_tss_table[] = {
	{ 100, 1000, 0, 0x00, 0 },
	{ 88, 875, 0, 0x1e, 0 },
	{ 75, 750, 0, 0x1c, 0 },
	{ 63, 625, 0, 0x1a, 0 },
	{ 50, 500, 0, 0x18, 0 },
	{ 38, 375, 0, 0x16, 0 },
	{ 25, 250, 0, 0x14, 0 },
	{ 13, 125, 0, 0x12, 0 },
};

acpi_tstate_t *soc_get_tss_table(int *entries)
{
	*entries = ARRAY_SIZE(denverton_tss_table);
	return denverton_tss_table;
}

void soc_power_states_generation(int core_id, int cores_per_package)
{
	generate_p_state_entries(core_id, cores_per_package);

	generate_t_state_entries(core_id, cores_per_package);
}

int soc_madt_sci_irq_polarity(int sci)
{
	if (sci >= 20)
		return MP_IRQ_POLARITY_LOW;
	else
		return MP_IRQ_POLARITY_HIGH;
}

unsigned long southcluster_write_acpi_tables(struct device *device,
					     unsigned long current,
					     struct acpi_rsdp *rsdp)
{
	acpi_header_t *ssdt2;

	current = acpi_write_hpet(device, current, rsdp);
	current = (ALIGN(current, 16));

	ssdt2 = (acpi_header_t *)current;
	memset(ssdt2, 0, sizeof(acpi_header_t));
	acpi_create_serialio_ssdt(ssdt2);
	if (ssdt2->length) {
		current += ssdt2->length;
		acpi_add_table(rsdp, ssdt2);
		printk(BIOS_DEBUG, "ACPI:     * SSDT2 @ %p Length %x\n", ssdt2,
		       ssdt2->length);
		current = (ALIGN(current, 16));
	} else {
		ssdt2 = NULL;
		printk(BIOS_DEBUG, "ACPI:     * SSDT2 not generated.\n");
	}

	printk(BIOS_DEBUG, "current = %lx\n", current);

	return current;
}

void southcluster_inject_dsdt(struct device *device)
{
	global_nvs_t *gnvs;

	gnvs = cbmem_find(CBMEM_ID_ACPI_GNVS);
	if (!gnvs) {
		gnvs = cbmem_add(CBMEM_ID_ACPI_GNVS, sizeof(*gnvs));
		if (gnvs)
			memset(gnvs, 0, sizeof(*gnvs));
	}

	if (gnvs) {
		acpi_create_gnvs(gnvs);
		/* And tell SMI about it */
		smm_setup_structures(gnvs, NULL, NULL);

		/* Add it to DSDT.  */
		acpigen_write_scope("\\");
		acpigen_write_name_dword("NVSA", (u32)gnvs);
		acpigen_pop_len();
	}
}

__weak void acpi_create_serialio_ssdt(acpi_header_t *ssdt) {}
