/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2016 Intel Corp.
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

#include <arch/romstage.h>
#include <arch/symbols.h>
#include <assert.h>
#include <cpu/x86/msr.h>
#include <cpu/x86/smm.h>
#include <cbmem.h>
#include <console/console.h>
#include <device/pci_def.h>
#include <fsp/util.h>
#include <intelblocks/cpulib.h>
#include <intelblocks/pmclib.h>
#include <memory_info.h>
#include <smbios.h>
#include <soc/intel/common/smbios.h>
#include <soc/msr.h>
#include <soc/pci_devs.h>
#include <soc/pm.h>
#include <soc/romstage.h>
#include <soc/systemagent.h>
#include <string.h>
#include <security/vboot/vboot_common.h>

#include "../chip.h"

#define FSP_SMBIOS_MEMORY_INFO_GUID	\
{	\
	0xd4, 0x71, 0x20, 0x9b, 0x54, 0xb0, 0x0c, 0x4e,	\
	0x8d, 0x09, 0x11, 0xcf, 0x8b, 0x9f, 0x03, 0x23	\
}

/* Memory Channel Present Status */
enum {
	CHANNEL_NOT_PRESENT,
	CHANNEL_DISABLED,
	CHANNEL_PRESENT
};

/* Save the DIMM information for SMBIOS table 17 */
static void save_dimm_info(void)
{
	int channel, dimm, dimm_max, index;
	size_t hob_size;
	uint8_t ddr_type;
	const CONTROLLER_INFO *ctrlr_info;
	const CHANNEL_INFO *channel_info;
	const DIMM_INFO *src_dimm;
	struct dimm_info *dest_dimm;
	struct memory_info *mem_info;
	const MEMORY_INFO_DATA_HOB *memory_info_hob;
	const uint8_t smbios_memory_info_guid[16] =
			FSP_SMBIOS_MEMORY_INFO_GUID;

	/* Locate the memory info HOB, presence validated by raminit */
	memory_info_hob =
		fsp_find_extension_hob_by_guid(smbios_memory_info_guid,
						&hob_size);
	if (memory_info_hob == NULL || hob_size == 0) {
		printk(BIOS_ERR, "SMBIOS MEMORY_INFO_DATA_HOB not found\n");
		return;
	}

	/*
	 * Allocate CBMEM area for DIMM information used to populate SMBIOS
	 * table 17
	 */
	mem_info = cbmem_add(CBMEM_ID_MEMINFO, sizeof(*mem_info));
	if (mem_info == NULL) {
		printk(BIOS_ERR, "CBMEM entry for DIMM info missing\n");
		return;
	}
	memset(mem_info, 0, sizeof(*mem_info));

	/* Describe the first N DIMMs in the system */
	index = 0;
	dimm_max = ARRAY_SIZE(mem_info->dimm);
	ctrlr_info = &memory_info_hob->Controller[0];
	for (channel = 0; channel < MAX_CH && index < dimm_max; channel++) {
		channel_info = &ctrlr_info->ChannelInfo[channel];
		if (channel_info->Status != CHANNEL_PRESENT)
			continue;
		for (dimm = 0; dimm < MAX_DIMM && index < dimm_max; dimm++) {
			src_dimm = &channel_info->DimmInfo[dimm];
			dest_dimm = &mem_info->dimm[index];

			if (src_dimm->Status != DIMM_PRESENT)
				continue;

			switch (memory_info_hob->MemoryType) {
			case MRC_DDR_TYPE_DDR4:
				ddr_type = MEMORY_TYPE_DDR4;
				break;
			case MRC_DDR_TYPE_DDR3:
				ddr_type = MEMORY_TYPE_DDR3;
				break;
			case MRC_DDR_TYPE_LPDDR3:
				ddr_type = MEMORY_TYPE_LPDDR3;
				break;
			default:
				ddr_type = MEMORY_TYPE_UNKNOWN;
				break;
			}
			u8 memProfNum = memory_info_hob->MemoryProfile;

			/* Populate the DIMM information */
			dimm_info_fill(dest_dimm,
				src_dimm->DimmCapacity,
				ddr_type,
				memory_info_hob->ConfiguredMemoryClockSpeed,
				src_dimm->RankInDimm,
				channel_info->ChannelId,
				src_dimm->DimmId,
				(const char *)src_dimm->ModulePartNum,
				sizeof(src_dimm->ModulePartNum),
				src_dimm->SpdSave + SPD_SAVE_OFFSET_SERIAL,
				memory_info_hob->DataWidth,
				memory_info_hob->VddVoltage[memProfNum],
				memory_info_hob->EccSupport,
				src_dimm->MfgId,
				src_dimm->SpdModuleType);
			index++;
		}
	}
	mem_info->dimm_cnt = index;
	printk(BIOS_DEBUG, "%d DIMMs found\n", mem_info->dimm_cnt);
}

void mainboard_romstage_entry(void)
{
	bool s3wake;
	struct chipset_power_state *ps;

	/* Program MCHBAR, DMIBAR, GDXBAR and EDRAMBAR */
	systemagent_early_init();

	ps = pmc_get_power_state();
	s3wake = pmc_fill_power_state(ps) == ACPI_S3;
	fsp_memory_init(s3wake);
	pmc_set_disb();
	if (!s3wake)
		save_dimm_info();
}

static void cpu_flex_override(FSP_M_CONFIG *m_cfg)
{
	msr_t flex_ratio;
	m_cfg->CpuRatioOverride = 1;
	/*
	 * Set cpuratio to that value set in bootblock, This will ensure FSPM
	 * knows the intended flex ratio.
	 */
	flex_ratio = rdmsr(MSR_FLEX_RATIO);
	m_cfg->CpuRatio = (flex_ratio.lo >> 8) & 0xff;
}

static void soc_peg_init_params(FSP_M_CONFIG *m_cfg,
			FSP_M_TEST_CONFIG *m_t_cfg,
			const struct soc_intel_skylake_config *config)
{
	const struct device *dev;
	/*
	 * To enable or disable the corresponding PEG root port you need to
	 * add to the devicetree.cb:
	 *
	 *     device pci 01.0 on  end # enable PEG0 root port
	 *     device pci 01.1 off end # do not configure PEG1
	 *
	 * If PEG port is not defined in the device tree, it will be disabled
	 * in FSP
	 */
	dev = pcidev_on_root(SA_DEV_SLOT_PEG, 0); /* PEG 0:1:0 */
	if (!dev || !dev->enabled)
		m_cfg->Peg0Enable = 0;
	else if (dev->enabled) {
		m_cfg->Peg0Enable = dev->enabled;
		m_cfg->Peg0MaxLinkWidth = config->Peg0MaxLinkWidth;
		/* Use maximum possible link speed */
		m_cfg->Peg0MaxLinkSpeed = 0;
		/* Power down unused lanes based on the max possible width */
		m_cfg->Peg0PowerDownUnusedLanes = 1;
		/* Set [Auto] for options to enable equalization methods */
		m_t_cfg->Peg0Gen3EqPh2Enable = 2;
		m_t_cfg->Peg0Gen3EqPh3Method = 0;
	}

	dev = pcidev_on_root(SA_DEV_SLOT_PEG, 1); /* PEG 0:1:1 */
	if (!dev || !dev->enabled)
		m_cfg->Peg1Enable = 0;
	else if (dev->enabled) {
		m_cfg->Peg1Enable = dev->enabled;
		m_cfg->Peg1MaxLinkWidth = config->Peg1MaxLinkWidth;
		m_cfg->Peg1MaxLinkSpeed = 0;
		m_cfg->Peg1PowerDownUnusedLanes = 1;
		m_t_cfg->Peg1Gen3EqPh2Enable = 2;
		m_t_cfg->Peg1Gen3EqPh3Method = 0;
	}

	dev = pcidev_on_root(SA_DEV_SLOT_PEG, 2); /* PEG 0:1:2 */
	if (!dev || !dev->enabled)
		m_cfg->Peg2Enable = 0;
	else if (dev->enabled) {
		m_cfg->Peg2Enable = dev->enabled;
		m_cfg->Peg2MaxLinkWidth = config->Peg2MaxLinkWidth;
		m_cfg->Peg2MaxLinkSpeed = 0;
		m_cfg->Peg2PowerDownUnusedLanes = 1;
		m_t_cfg->Peg2Gen3EqPh2Enable = 2;
		m_t_cfg->Peg2Gen3EqPh3Method = 0;
	}
}

static void soc_memory_init_params(FSP_M_CONFIG *m_cfg,
			const struct soc_intel_skylake_config *config)
{
	int i;
	uint32_t mask = 0;

	m_cfg->MmioSize = 0x800; /* 2GB in MB */
	m_cfg->TsegSize = CONFIG_SMM_TSEG_SIZE;
	m_cfg->IedSize = CONFIG_IED_REGION_SIZE;
	m_cfg->ProbelessTrace = config->ProbelessTrace;
	m_cfg->SaGv = config->SaGv;
	m_cfg->UserBd = BOARD_TYPE_ULT_ULX;
	m_cfg->RMT = config->Rmt;
	m_cfg->CmdTriStateDis = config->CmdTriStateDis;
	m_cfg->DdrFreqLimit = config->DdrFreqLimit;
	m_cfg->VmxEnable = CONFIG(ENABLE_VMX);
	m_cfg->PrmrrSize = get_prmrr_size();
	for (i = 0; i < ARRAY_SIZE(config->PcieRpEnable); i++) {
		if (config->PcieRpEnable[i])
			mask |= (1<<i);
	}
	m_cfg->PcieRpEnableMask = mask;

	cpu_flex_override(m_cfg);

	if (!config->ignore_vtd) {
		m_cfg->PchHpetBdfValid = 1;
		m_cfg->PchHpetBusNumber = V_P2SB_HBDF_BUS;
		m_cfg->PchHpetDeviceNumber = V_P2SB_HBDF_DEV;
		m_cfg->PchHpetFunctionNumber = V_P2SB_HBDF_FUN;
	}
	m_cfg->HyperThreading = CONFIG(FSP_HYPERTHREADING);
}

static void soc_primary_gfx_config_params(FSP_M_CONFIG *m_cfg,
				const struct soc_intel_skylake_config *config)
{
	const struct device *dev;

	dev = pcidev_path_on_root(SA_DEVFN_IGD);
	if (!dev || !dev->enabled) {
		/*
		 * If iGPU is disabled or not defined in the devicetree.cb,
		 * the FSP does not initialize this device
		 */
		m_cfg->InternalGfx = 0;
		m_cfg->IgdDvmt50PreAlloc = 0;
	} else {
		m_cfg->InternalGfx = 1;
		/*
		 * Set IGD stolen size to 64MB.  The FBC hardware for skylake
		 * does not have access to the bios_reserved range so it always
		 * assumes 8MB is used and so the kernel will avoid the last
		 * 8MB of the stolen window. With the default stolen size of
		 * 32MB(-8MB) there is not enough space for FBC to work with
		 * a high resolution panel
		 */
		m_cfg->IgdDvmt50PreAlloc = 2;
	}
	m_cfg->PrimaryDisplay = config->PrimaryDisplay;
}

void platform_fsp_memory_init_params_cb(FSPM_UPD *mupd, uint32_t version)
{
	const struct soc_intel_skylake_config *config;
	FSP_M_CONFIG *m_cfg = &mupd->FspmConfig;
	FSP_M_TEST_CONFIG *m_t_cfg = &mupd->FspmTestConfig;

	config = config_of_soc();

	soc_memory_init_params(m_cfg, config);
	soc_peg_init_params(m_cfg, m_t_cfg, config);

	/* Skip creating Management Engine MBP HOB */
	m_t_cfg->SkipMbpHob = 0x01;

	/* Enable DMI Virtual Channel for ME */
	m_t_cfg->DmiVcm = 0x01;

	/* Enable Sending DID to ME */
	m_t_cfg->SendDidMsg = 0x01;
	m_t_cfg->DidInitStat = 0x01;

	/* DCI and TraceHub configs */
	m_t_cfg->PchDciEn = config->PchDciEn;
	m_cfg->EnableTraceHub = config->EnableTraceHub;
	m_cfg->TraceHubMemReg0Size = config->TraceHubMemReg0Size;
	m_cfg->TraceHubMemReg1Size = config->TraceHubMemReg1Size;

	/* Enable SMBus controller based on config */
	m_cfg->SmbusEnable = config->SmbusEnable;

	/* Set primary graphic device */
	soc_primary_gfx_config_params(m_cfg, config);
	m_t_cfg->SkipExtGfxScan = config->SkipExtGfxScan;

	mainboard_memory_init_params(mupd);
}

void soc_update_memory_params_for_mma(FSP_M_CONFIG *memory_cfg,
		struct mma_config_param *mma_cfg)
{
	/* Boot media is memory mapped for Skylake and Kabylake (SPI). */
	assert(CONFIG(BOOT_DEVICE_MEMORY_MAPPED));

	memory_cfg->MmaTestContentPtr =
			(uintptr_t) rdev_mmap_full(&mma_cfg->test_content);
	memory_cfg->MmaTestContentSize =
			region_device_sz(&mma_cfg->test_content);
	memory_cfg->MmaTestConfigPtr =
			(uintptr_t) rdev_mmap_full(&mma_cfg->test_param);
	memory_cfg->MmaTestConfigSize =
			region_device_sz(&mma_cfg->test_param);
	memory_cfg->MrcFastBoot = 0x00;
	memory_cfg->SaGv = 0x02;
}

__weak void mainboard_memory_init_params(FSPM_UPD *mupd)
{
	/* Do nothing */
}
