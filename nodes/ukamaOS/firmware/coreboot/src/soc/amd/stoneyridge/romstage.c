/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2015-2016 Advanced Micro Devices, Inc.
 * Copyright (C) 2015 Intel Corp.
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

#include <device/pci_ops.h>
#include <arch/cpu.h>
#include <arch/romstage.h>
#include <arch/acpi.h>
#include <cpu/x86/msr.h>
#include <cpu/x86/mtrr.h>
#include <cpu/x86/smm.h>
#include <cpu/amd/mtrr.h>
#include <cbmem.h>
#include <commonlib/helpers.h>
#include <console/console.h>
#include <device/device.h>
#include <program_loading.h>
#include <romstage_handoff.h>
#include <elog.h>
#include <amdblocks/agesawrapper.h>
#include <amdblocks/agesawrapper_call.h>
#include <soc/northbridge.h>
#include <soc/romstage.h>
#include <soc/southbridge.h>
#include <amdblocks/psp.h>

#include "chip.h"

void __weak mainboard_romstage_entry_s3(int s3_resume)
{
	/* By default, don't do anything */
}

static void load_smu_fw1(void)
{
	u32 base, limit, cmd;

	/* Open a posted hole from 0x80000000 : 0xfed00000-1 */
	base = (0x80000000 >> 8) | MMIO_WE | MMIO_RE;
	limit = (ALIGN_DOWN(HPET_BASE_ADDRESS - 1, 64 * KiB) >> 8);
	pci_write_config32(SOC_ADDR_DEV, D18F1_MMIO_LIMIT0_LO, limit);
	pci_write_config32(SOC_ADDR_DEV, D18F1_MMIO_BASE0_LO, base);

	/* Preload a value into "BAR3" and enable it */
	pci_write_config32(SOC_PSP_DEV, PSP_MAILBOX_BAR, PSP_MAILBOX_BAR3_BASE);
	pci_write_config32(SOC_PSP_DEV, PSP_BAR_ENABLES, PSP_MAILBOX_BAR_EN);

	/* Enable memory access and master */
	cmd = pci_read_config32(SOC_PSP_DEV, PCI_COMMAND);
	cmd |= PCI_COMMAND_MEMORY | PCI_COMMAND_MASTER;
	pci_write_config32(SOC_PSP_DEV, PCI_COMMAND, cmd);

	psp_load_named_blob(MBOX_BIOS_CMD_SMU_FW, "smu_fw");
}

static void agesa_call(void)
{
	post_code(0x37);
	do_agesawrapper(AMD_INIT_RESET, "amdinitreset");

	post_code(0x38);
	/* APs will not exit amdinitearly */
	do_agesawrapper(AMD_INIT_EARLY, "amdinitearly");
}

static void bsp_agesa_call(void)
{
	set_ap_entry_ptr(agesa_call); /* indicate the path to the AP */
	agesa_call();
}

asmlinkage void car_stage_entry(void)
{
	struct postcar_frame pcf;
	uintptr_t top_of_ram;
	msr_t base, mask;
	msr_t mtrr_cap = rdmsr(MTRR_CAP_MSR);
	int vmtrrs = mtrr_cap.lo & MTRR_CAP_VCNT;
	int s3_resume = acpi_s3_resume_allowed() && acpi_is_wakeup_s3();
	int i;

	console_init();

	if (CONFIG(SOC_AMD_PSP_SELECTABLE_SMU_FW))
		load_smu_fw1();

	mainboard_romstage_entry_s3(s3_resume);
	elog_boot_notify(s3_resume);

	bsp_agesa_call();

	if (!s3_resume) {
		post_code(0x40);
		do_agesawrapper(AMD_INIT_POST, "amdinitpost");

		post_code(0x41);
		/*
		 * TODO: This is a hack to work around current AGESA behavior.
		 *       AGESA needs to change to reflect that coreboot owns
		 *       the MTRRs.
		 *
		 * After setting up DRAM, AGESA also completes the configuration
		 * of the MTRRs, setting regions to WB.  Anything written to
		 * memory between now and and when CAR is dismantled will be
		 * in cache and lost.  For now, set the regions UC to ensure
		 * the writes get to DRAM.
		 */
		for (i = 0 ; i < vmtrrs ; i++) {
			base = rdmsr(MTRR_PHYS_BASE(i));
			mask = rdmsr(MTRR_PHYS_MASK(i));
			if (!(mask.lo & MTRR_PHYS_MASK_VALID))
				continue;

			if ((base.lo & 0x7) == MTRR_TYPE_WRBACK) {
				base.lo &= ~0x7;
				base.lo |= MTRR_TYPE_UNCACHEABLE;
				wrmsr(MTRR_PHYS_BASE(i), base);
			}
		}
		/* Disable WB from to region 4GB-TOM2. */
		msr_t sys_cfg = rdmsr(SYSCFG_MSR);
		sys_cfg.lo &= ~SYSCFG_MSR_TOM2WB;
		wrmsr(SYSCFG_MSR, sys_cfg);
	} else {
		printk(BIOS_INFO, "S3 detected\n");
		post_code(0x60);
		do_agesawrapper(AMD_INIT_RESUME, "amdinitresume");

		post_code(0x61);
	}

	post_code(0x42);
	psp_notify_dram();

	post_code(0x43);
	if (cbmem_recovery(s3_resume))
		printk(BIOS_CRIT, "Failed to recover cbmem\n");
	if (romstage_handoff_init(s3_resume))
		printk(BIOS_ERR, "Failed to set romstage handoff data\n");

	if (CONFIG(SMM_TSEG))
		smm_list_regions();

	post_code(0x44);
	if (postcar_frame_init(&pcf, 0))
		die("Unable to initialize postcar frame.\n");

	/*
	 * We need to make sure ramstage will be run cached. At this point exact
	 * location of ramstage in cbmem is not known. Instruct postcar to cache
	 * 16 megs under cbmem top which is a safe bet to cover ramstage.
	 */
	top_of_ram = (uintptr_t) cbmem_top();
	postcar_frame_add_mtrr(&pcf, top_of_ram - 16*MiB, 16*MiB,
		MTRR_TYPE_WRBACK);

	/* Cache the memory-mapped boot media. */
	postcar_frame_add_romcache(&pcf, MTRR_TYPE_WRPROT);

	/* Cache the TSEG region */
	postcar_enable_tseg_cache(&pcf);

	post_code(0x45);
	run_postcar_phase(&pcf);

	post_code(0x50);  /* Should never see this post code. */
}

void SetMemParams(AMD_POST_PARAMS *PostParams)
{
	const struct soc_amd_stoneyridge_config *cfg;
	const struct device *dev = pcidev_path_on_root(GNB_DEVFN);

	if (!dev || !dev->chip_info) {
		printk(BIOS_ERR, "ERROR: Cannot find SoC devicetree config\n");
		/* In case of a BIOS error, only attempt to set UMA. */
		PostParams->MemConfig.UmaMode = CONFIG(GFXUMA) ?
					UMA_AUTO : UMA_NONE;
		return;
	}

	cfg = dev->chip_info;

	PostParams->MemConfig.EnableMemClr = cfg->dram_clear_on_reset;

	switch (cfg->uma_mode) {
	case UMAMODE_NONE:
		PostParams->MemConfig.UmaMode = UMA_NONE;
		break;
	case UMAMODE_SPECIFIED_SIZE:
		PostParams->MemConfig.UmaMode = UMA_SPECIFIED;
		/* 64 KiB blocks. */
		PostParams->MemConfig.UmaSize = cfg->uma_size / (64 * KiB);
		break;
	case UMAMODE_AUTO_LEGACY:
		PostParams->MemConfig.UmaMode = UMA_AUTO;
		PostParams->MemConfig.UmaVersion = UMA_LEGACY;
		break;
	case UMAMODE_AUTO_NON_LEGACY:
		PostParams->MemConfig.UmaMode = UMA_AUTO;
		PostParams->MemConfig.UmaVersion = UMA_NON_LEGACY;
		break;
	}
}

void soc_customize_init_early(AMD_EARLY_PARAMS *InitEarly)
{
	const struct soc_amd_stoneyridge_config *cfg;
	const struct device *dev = pcidev_path_on_root(GNB_DEVFN);
	struct _PLATFORM_CONFIGURATION *platform;

	if (!dev || !dev->chip_info) {
		printk(BIOS_WARNING, "Warning: Cannot find SoC devicetree"
					" config, STAPM unchanged\n");
		return;
	}
	cfg = dev->chip_info;
	platform = &InitEarly->PlatformConfig;
	if ((cfg->stapm_percent) && (cfg->stapm_time_ms) &&
				    (cfg->stapm_power_mw)) {
		platform->PlatStapmConfig.CfgStapmScalar = cfg->stapm_percent;
		platform->PlatStapmConfig.CfgStapmTimeConstant =
							cfg->stapm_time_ms;
		platform->PkgPwrLimitDC = cfg->stapm_power_mw;
		platform->PkgPwrLimitAC = cfg->stapm_power_mw;
		platform->PlatStapmConfig.CfgStapmBoost = StapmBoostEnabled;
	}
}
