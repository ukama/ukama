/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2014 Google Inc.
 * Copyright (C) 2017 Siemens AG
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

#define __SIMPLE_DEVICE__

#include <types.h>
#include <string.h>
#include <device/pci.h>
#include <cpu/x86/lapic.h>
#include <cpu/x86/mp.h>
#include <cpu/x86/mtrr.h>
#include <cpu/x86/smm.h>
#include <cpu/intel/em64t101_save_state.h>
#include <cpu/intel/smm_reloc.h>
#include <console/console.h>
#include <device/pci_ops.h>
#include <soc/msr.h>
#include <soc/pci_devs.h>
#include <soc/smm.h>
#include <soc/broadwell_de.h>

/* This gets filled in and used during relocation. */
static struct smm_relocation_params smm_reloc_params;

static inline void write_smrr(struct smm_relocation_params *relo_params)
{
	printk(BIOS_DEBUG, "Writing SMRR. base = 0x%08x, mask=0x%08x\n",
			relo_params->smrr_base.lo, relo_params->smrr_mask.lo);
	wrmsr(IA32_SMRR_PHYS_BASE, relo_params->smrr_base);
	wrmsr(IA32_SMRR_PHYS_MASK, relo_params->smrr_mask);
}

static inline void write_prmrr(struct smm_relocation_params *relo_params)
{
	printk(BIOS_DEBUG, "Writing PRMRR. base = 0x%08x, mask=0x%08x\n",
			relo_params->prmrr_base.lo, relo_params->prmrr_mask.lo);
	wrmsr(MSR_PRMRR_PHYS_BASE, relo_params->prmrr_base);
	wrmsr(MSR_PRMRR_PHYS_MASK, relo_params->prmrr_mask);
}

static void update_save_state(int cpu, uintptr_t curr_smbase,
				uintptr_t staggered_smbase,
				struct smm_relocation_params *relo_params)
{
	u32 smbase;
	u32 iedbase;

	/* The relocated handler runs with all CPUs concurrently. Therefore
	   stagger the entry points adjusting SMBASE downwards by save state
	   size * CPU num. */
	smbase = staggered_smbase;
	iedbase = relo_params->ied_base;

	printk(BIOS_DEBUG, "New SMBASE=0x%08x IEDBASE=0x%08x\n",
		smbase, iedbase);

	/*
	 *  All threads need to set IEDBASE and SMBASE to the relocated
	 * handler region. However, the save state location depends on the
	 * smm_save_state_in_msrs field in the relocation parameters. If
	 * smm_save_state_in_msrs is non-zero then the CPUs are relocating
	 * the SMM handler in parallel, and each CPUs save state area is
	 * located in their respective MSR space. If smm_save_state_in_msrs
	 * is zero then the SMM relocation is happening serially so the
	 * save state is at the same default location for all CPUs.
	 */
	if (relo_params->smm_save_state_in_msrs) {
		msr_t smbase_msr;
		msr_t iedbase_msr;

		smbase_msr.lo = smbase;
		smbase_msr.hi = 0;

		/* According the BWG the IEDBASE MSR is in bits 63:32. It's
		   not clear why it differs from the SMBASE MSR. */
		iedbase_msr.lo = 0;
		iedbase_msr.hi = iedbase;

		wrmsr(SMBASE_MSR, smbase_msr);
		wrmsr(IEDBASE_MSR, iedbase_msr);
	} else {
		em64t101_smm_state_save_area_t *save_state;

		save_state = (void *)(curr_smbase + SMM_DEFAULT_SIZE -
					sizeof(*save_state));
		save_state->smbase = smbase;
		save_state->iedbase = iedbase;
	}
}

/* Returns 1 if SMM MSR save state was set. */
static int bsp_setup_msr_save_state(struct smm_relocation_params *relo_params)
{
	msr_t smm_mca_cap;

	smm_mca_cap = rdmsr(SMM_MCA_CAP_MSR);
	if (smm_mca_cap.hi & SMM_CPU_SVRSTR_MASK) {
		uint32_t smm_feature_control;
		pci_devfn_t dev = PCI_DEV(QPI_BUS, SMM_DEV, SMM_FUNC);

		/*
		 * SMM_FEATURE_CONTROL on Broadwell-DE is not located in
		 * MSR range but in PCI config space. The used PCI device is
		 * located on bus 0xff, which has no root bridge and hence is
		 * not scanned by PCI scan. Use MMIO config access to read the
		 * needed 32 bit register.
		 */
		smm_feature_control = pci_read_config32(dev,
							SMM_FEATURE_CONTROL);
		smm_feature_control |= SMM_CPU_SAVE_EN;
		pci_write_config32(dev,
				   SMM_FEATURE_CONTROL, smm_feature_control);
		relo_params->smm_save_state_in_msrs = 1;
	}
	return relo_params->smm_save_state_in_msrs;
}

/*
 * The relocation work is actually performed in SMM context, but the code
 * resides in the ramstage module. This occurs by trampolining from the default
 * SMRAM entry point to here.
 */
void smm_relocation_handler(int cpu, uintptr_t curr_smbase,
				uintptr_t staggered_smbase)
{
	msr_t mtrr_cap;
	struct smm_relocation_params *relo_params = &smm_reloc_params;

	printk(BIOS_DEBUG, "In relocation handler: CPU %d\n", cpu);

	/* Determine if the processor supports saving state in MSRs. If so,
	   enable it before the non-BSPs run so that SMM relocation can occur
	   in parallel in the non-BSP CPUs. */
	if (cpu == 0) {
		/*
		 * If smm_save_state_in_msrs is 1 then that means this is the
		 * 2nd time through the relocation handler for the BSP.
		 * Parallel SMM handler relocation is taking place. However,
		 * it is desired to access other CPUs save state in the real
		 * SMM handler. Therefore, disable the SMM save state in MSRs
		 * feature.
		 */
		if (relo_params->smm_save_state_in_msrs) {
			uint32_t smm_feature_control;
			pci_devfn_t dev = PCI_DEV(QPI_BUS, SMM_DEV, SMM_FUNC);

			/*
			 * SMM_FEATURE_CONTROL on Broadwell-DE is not located in
			 * MSR range but in PCI config space. The used PCI
			 * device is located on bus 0xff, which has no root
			 * bridge and hence is not scanned by PCI scan.
			 * Use MMIO config access to read the needed 32 bit
			 * register.
			 */
			smm_feature_control = pci_read_config32(dev,
							SMM_FEATURE_CONTROL);
			smm_feature_control &= ~SMM_CPU_SAVE_EN;
			pci_write_config32(dev, SMM_FEATURE_CONTROL,
						smm_feature_control);
		} else if (bsp_setup_msr_save_state(relo_params))
			/*
			 * Just return from relocation handler if MSR save
			 * state is enabled. In that case the BSP will come
			 * back into the relocation handler to setup the new
			 * SMBASE as well disabling SMM save state in MSRs.
			 */
			return;
	}

	/* Make appropriate changes to the save state map. */
	update_save_state(cpu, curr_smbase, staggered_smbase, relo_params);
	/* Write PRMRR and SMRR MSRs based on indicated support. */
	mtrr_cap = rdmsr(MTRR_CAP_MSR);
	if (mtrr_cap.lo & SMRR_SUPPORTED)
		write_smrr(relo_params);

	if (mtrr_cap.lo & PRMRR_SUPPORTED)
		write_prmrr(relo_params);
}

static void fill_in_relocation_params(struct smm_relocation_params *params)
{
	uintptr_t tseg_base;
	size_t tseg_size;
	u32 prmrr_base;
	u32 prmrr_size;
	int phys_bits;
	/* All range registers are aligned to 4KiB */
	const u32 rmask = ~((1 << 12) - 1);

	/* Some of the range registers are dependent on the number of physical
	   address bits supported. */
	phys_bits = cpuid_eax(0x80000008) & 0xff;
	/*
	 * The range bounded by the TSEG_BASE and TSEG_LIMIT registers
	 * encompasses the SMRAM range as well as the IED range.
	 * However, the SMRAM available to the handler is 4MiB since the IEDRAM
	 * lives TSEG_BASE + 4MiB.
	 *
	 * Note that address bits 19:0 are ignored and not compared.
	 * The result is that BASE[19:0] is effectively 00000h and LIMIT is
	 * effectively FFFFFh.
	 */

	smm_region(&tseg_base, &tseg_size);

	/* SMRR has 32-bits of valid address aligned to 4KiB. */
	params->smrr_base.lo = (tseg_base & rmask) | MTRR_TYPE_WRBACK;
	params->smrr_base.hi = 0;
	params->smrr_mask.lo = (~(tseg_size - 1) & rmask) | MTRR_PHYS_MASK_VALID;
	params->smrr_mask.hi = 0;

	smm_subregion(SMM_SUBREGION_CHIPSET, &params->ied_base, &params->ied_size);

	/* The PRMRR is at IEDBASE + 2MiB */
	prmrr_base = (params->ied_base + (2 << 20)) & rmask;
	prmrr_size = params->ied_size - (2 << 20);

	/* PRMRR has 46 bits of valid address aligned to 4KiB. It's dependent
	   on the number of physical address bits supported. */
	params->prmrr_base.lo = prmrr_base | MTRR_TYPE_WRBACK;
	params->prmrr_base.hi = 0;
	params->prmrr_mask.lo = (~(prmrr_size - 1) & rmask)
		| MTRR_PHYS_MASK_VALID;
	params->prmrr_mask.hi = (1 << (phys_bits - 32)) - 1;
}

static void setup_ied_area(struct smm_relocation_params *params)
{
	char *ied_base;

	struct ied_header ied = {
		.signature = "INTEL RSVD",
		.size = params->ied_size,
		.reserved = {0},
	};

	ied_base = (void *)params->ied_base;

	/* Place IED header at IEDBASE. */
	memcpy(ied_base, &ied, sizeof(ied));

	/* Zero out 32KiB at IEDBASE + 1MiB */
	memset(ied_base + (1 << 20), 0, (32 << 10));
}

void smm_info(uintptr_t *perm_smbase, size_t *perm_smsize,
				size_t *smm_save_state_size)
{
	printk(BIOS_DEBUG, "Setting up SMI for CPU\n");

	fill_in_relocation_params(&smm_reloc_params);

	smm_subregion(SMM_SUBREGION_HANDLER, perm_smbase, perm_smsize);

	setup_ied_area(&smm_reloc_params);

	*smm_save_state_size = sizeof(em64t101_smm_state_save_area_t);
}

void smm_initialize(void)
{
	/* Clear the SMM state in the southbridge. */
	smm_southbridge_clear_state();

	/* Run the relocation handler for on the BSP to check and set up
	   parallel SMM relocation. */
	smm_initiate_relocation();

	if (smm_reloc_params.smm_save_state_in_msrs)
		printk(BIOS_DEBUG, "Doing parallel SMM relocation.\n");
}

/*
 * The default SMM entry can happen in parallel or serially. If the
 * default SMM entry is done in parallel the BSP has already setup
 * the saving state to each CPU's MSRs. At least one save state size
 * is required for the initial SMM entry for the BSP to determine if
 * parallel SMM relocation is even feasible.
 */
void smm_relocate(void)
{
	/*
	 * If smm_save_state_in_msrs is non-zero then parallel SMM relocation
	 * shall take place. Run the relocation handler a second time on the
	 * BSP to do the final move. For APs, a relocation handler always
	 * needs to be run.
	 */
	if (smm_reloc_params.smm_save_state_in_msrs)
		smm_initiate_relocation_parallel();
	else if (!boot_cpu())
		smm_initiate_relocation();
}
