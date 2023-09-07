/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2016-2017 Intel Corporation.
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

#include <bootmode.h>
#include <bootstate.h>
#include <fsp/api.h>
#include <arch/acpi.h>
#include <console/console.h>
#include <device/device.h>
#include <device/pci_ids.h>
#include <fsp/util.h>
#include <intelblocks/cfg.h>
#include <intelblocks/itss.h>
#include <intelblocks/lpc_lib.h>
#include <intelblocks/mp_init.h>
#include <intelblocks/pcie_rp.h>
#include <intelblocks/xdci.h>
#include <intelblocks/p2sb.h>
#include <intelpch/lockdown.h>
#include <romstage_handoff.h>
#include <soc/acpi.h>
#include <soc/intel/common/vbt.h>
#include <soc/interrupt.h>
#include <soc/iomap.h>
#include <soc/irq.h>
#include <soc/itss.h>
#include <soc/pci_devs.h>
#include <soc/ramstage.h>
#include <soc/systemagent.h>
#include <string.h>

#include "chip.h"

static const struct pcie_rp_group pch_lp_rp_groups[] = {
	{ .slot = PCH_DEV_SLOT_PCIE,	.count = 8 },
	{ .slot = PCH_DEV_SLOT_PCIE_1,	.count = 4 },
	{ 0 }
};

static const struct pcie_rp_group pch_h_rp_groups[] = {
	{ .slot = PCH_DEV_SLOT_PCIE,	.count = 8 },
	{ .slot = PCH_DEV_SLOT_PCIE_1,	.count = 8 },
	/* Sunrise Point PCH-H actually only has 4 ports in the
	   third group. But that would require a runtime check
	   and probing 4 non-existent ports shouldn't hurt. */
	{ .slot = PCH_DEV_SLOT_PCIE_2,	.count = 8 },
	{ 0 }
};

void soc_init_pre_device(void *chip_info)
{
	/* Snapshot the current GPIO IRQ polarities. FSP is setting a
	 * default policy that doesn't honor boards' requirements. */
	itss_snapshot_irq_polarities(GPIO_IRQ_START, GPIO_IRQ_END);

	/* Perform silicon specific init. */
	fsp_silicon_init(romstage_handoff_is_resume());

	/*
	 * Keep the P2SB device visible so it and the other devices are
	 * visible in coreboot for driver support and PCI resource allocation.
	 * There is no UPD setting for this.
	 */
	p2sb_unhide();

	/* Restore GPIO IRQ polarities back to previous settings. */
	itss_restore_irq_polarities(GPIO_IRQ_START, GPIO_IRQ_END);

	/* swap enabled PCI ports in device tree if needed */
	if (CONFIG(SKYLAKE_SOC_PCH_H))
		pcie_rp_update_devicetree(pch_h_rp_groups);
	else
		pcie_rp_update_devicetree(pch_lp_rp_groups);
}

void soc_fsp_load(void)
{
	fsps_load(romstage_handoff_is_resume());
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
	.write_acpi_tables	= &northbridge_write_acpi_tables,
	.acpi_name		= &soc_acpi_name,
#endif
};

static struct device_operations cpu_bus_ops = {
	.read_resources   = DEVICE_NOOP,
	.set_resources    = DEVICE_NOOP,
	.enable_resources = DEVICE_NOOP,
	.init             = DEVICE_NOOP,
#if CONFIG(HAVE_ACPI_TABLES)
	.acpi_fill_ssdt_generator = generate_cpu_entries,
#endif
};

static void soc_enable(struct device *dev)
{
	/* Set the operations if it is a special bus type */
	if (dev->path.type == DEVICE_PATH_DOMAIN)
		dev->ops = &pci_domain_ops;
	else if (dev->path.type == DEVICE_PATH_CPU_CLUSTER)
		dev->ops = &cpu_bus_ops;
}

struct chip_operations soc_intel_skylake_ops = {
	CHIP_NAME("Intel 6th Gen")
	.enable_dev	= &soc_enable,
	.init		= &soc_init_pre_device,
};

/* UPD parameters to be initialized before SiliconInit */
void platform_fsp_silicon_init_params_cb(FSPS_UPD *supd)
{
	FSP_S_CONFIG *params = &supd->FspsConfig;
	FSP_S_TEST_CONFIG *tconfig = &supd->FspsTestConfig;
	struct soc_intel_skylake_config *config;
	struct device *dev;
	uintptr_t vbt_data = (uintptr_t)vbt_get();
	int i;

	config = config_of_soc();

	mainboard_silicon_init_params(params);
	/* Set PsysPmax if it is available from DT */
	if (config->psys_pmax) {
		/* PsysPmax is in unit of 1/8 Watt */
		tconfig->PsysPmax = config->psys_pmax * 8;
		printk(BIOS_DEBUG, "psys_pmax = %d\n", tconfig->PsysPmax);
	}

	params->GraphicsConfigPtr = (u32) vbt_data;

	for (i = 0; i < ARRAY_SIZE(config->usb2_ports); i++) {
		params->PortUsb20Enable[i] =
				config->usb2_ports[i].enable;
		params->Usb2OverCurrentPin[i] =
				config->usb2_ports[i].ocpin;
		params->Usb2AfePetxiset[i] =
				config->usb2_ports[i].pre_emp_bias;
		params->Usb2AfeTxiset[i] =
				config->usb2_ports[i].tx_bias;
		params->Usb2AfePredeemp[i] =
				config->usb2_ports[i].tx_emp_enable;
		params->Usb2AfePehalfbit[i] =
				config->usb2_ports[i].pre_emp_bit;
	}

	for (i = 0; i < ARRAY_SIZE(config->usb3_ports); i++) {
		params->PortUsb30Enable[i] = config->usb3_ports[i].enable;
		params->Usb3OverCurrentPin[i] = config->usb3_ports[i].ocpin;
		if (config->usb3_ports[i].tx_de_emp) {
			params->Usb3HsioTxDeEmphEnable[i] = 1;
			params->Usb3HsioTxDeEmph[i] =
				config->usb3_ports[i].tx_de_emp;
		}
		if (config->usb3_ports[i].tx_downscale_amp) {
			params->Usb3HsioTxDownscaleAmpEnable[i] = 1;
			params->Usb3HsioTxDownscaleAmp[i] =
				config->usb3_ports[i].tx_downscale_amp;
		}
	}

	memcpy(params->SataPortsEnable, config->SataPortsEnable,
	       sizeof(params->SataPortsEnable));
	memcpy(params->SataPortsDevSlp, config->SataPortsDevSlp,
	       sizeof(params->SataPortsDevSlp));
	memcpy(params->SataPortsHotPlug, config->SataPortsHotPlug,
	       sizeof(params->SataPortsHotPlug));
	memcpy(params->SataPortsSpinUp, config->SataPortsSpinUp,
	       sizeof(params->SataPortsSpinUp));
	memcpy(params->PcieRpClkReqSupport, config->PcieRpClkReqSupport,
	       sizeof(params->PcieRpClkReqSupport));
	memcpy(params->PcieRpClkReqNumber, config->PcieRpClkReqNumber,
	       sizeof(params->PcieRpClkReqNumber));
	memcpy(params->PcieRpAdvancedErrorReporting,
		config->PcieRpAdvancedErrorReporting,
			sizeof(params->PcieRpAdvancedErrorReporting));
	memcpy(params->PcieRpLtrEnable, config->PcieRpLtrEnable,
	       sizeof(params->PcieRpLtrEnable));
	memcpy(params->PcieRpHotPlug, config->PcieRpHotPlug,
	       sizeof(params->PcieRpHotPlug));

	/*
	 * PcieRpClkSrcNumber UPD is set to clock source number(0-6) for
	 * all the enabled PCIe root ports, invalid(0x1F) is set for
	 * disabled PCIe root ports.
	 */
	for (i = 0; i < CONFIG_MAX_ROOT_PORTS; i++) {
		if (config->PcieRpClkReqSupport[i])
			params->PcieRpClkSrcNumber[i] =
				config->PcieRpClkSrcNumber[i];
		else
			params->PcieRpClkSrcNumber[i] = 0x1F;
	}

	/* disable Legacy PME */
	memset(params->PcieRpPmSci, 0, sizeof(params->PcieRpPmSci));

	/* Legacy 8254 timer support */
	params->Early8254ClockGatingEnable = !CONFIG_USE_LEGACY_8254_TIMER;

	memcpy(params->SerialIoDevMode, config->SerialIoDevMode,
	       sizeof(params->SerialIoDevMode));

	params->PchCio2Enable = config->Cio2Enable;
	params->SaImguEnable = config->SaImguEnable;
	params->Heci3Enabled = config->Heci3Enabled;

	params->LogoPtr = config->LogoPtr;
	params->LogoSize = config->LogoSize;

	params->CpuConfig.Bits.VmxEnable = CONFIG(ENABLE_VMX);

	params->PchPmWoWlanEnable = config->PchPmWoWlanEnable;
	params->PchPmWoWlanDeepSxEnable = config->PchPmWoWlanDeepSxEnable;
	params->PchPmLanWakeFromDeepSx = config->WakeConfigPcieWakeFromDeepSx;

	params->PchLanEnable = config->EnableLan;
	if (config->EnableLan) {
		params->PchLanLtrEnable = config->EnableLanLtr;
		params->PchLanK1OffEnable = config->EnableLanK1Off;
		params->PchLanClkReqSupported = config->LanClkReqSupported;
		params->PchLanClkReqNumber = config->LanClkReqNumber;
	}
	params->SataSalpSupport = config->SataSalpSupport;
	params->SsicPortEnable = config->SsicPortEnable;
	params->ScsEmmcEnabled = config->ScsEmmcEnabled;
	params->ScsEmmcHs400Enabled = config->ScsEmmcHs400Enabled;
	params->ScsSdCardEnabled = config->ScsSdCardEnabled;

	if (!!params->ScsEmmcHs400Enabled && !!config->EmmcHs400DllNeed) {
		params->PchScsEmmcHs400DllDataValid =
			!!config->EmmcHs400DllNeed;
		params->PchScsEmmcHs400RxStrobeDll1 =
			config->ScsEmmcHs400RxStrobeDll1;
		params->PchScsEmmcHs400TxDataDll =
			config->ScsEmmcHs400TxDataDll;
	}

	/* If ISH is enabled, enable ISH elements */
	dev = pcidev_path_on_root(PCH_DEVFN_ISH);
	params->PchIshEnable = dev ? dev->enabled : 0;

	params->PchHdaEnable = config->EnableAzalia;
	params->PchHdaVcType = config->PchHdaVcType;
	params->PchHdaIoBufferOwnership = config->IoBufferOwnership;
	params->PchHdaDspEnable = config->DspEnable;
	params->Device4Enable = config->Device4Enable;
	params->SataEnable = config->EnableSata;
	params->SataMode = config->SataMode;
	params->SataSpeedLimit = config->SataSpeedLimit;
	params->SataPwrOptEnable = config->SataPwrOptEnable;
	params->EnableTcoTimer = !config->PmTimerDisabled;

	tconfig->PchLockDownGlobalSmi = config->LockDownConfigGlobalSmi;
	tconfig->PchLockDownRtcLock = config->LockDownConfigRtcLock;
	tconfig->PowerLimit4 = config->PowerLimit4;
	tconfig->SataTestMode = config->SataTestMode;
	/*
	 * To disable HECI, the Psf needs to be left unlocked
	 * by FSP till end of post sequence. Based on the devicetree
	 * setting, we set the appropriate PsfUnlock policy in FSP,
	 * do the changes and then lock it back in coreboot during finalize.
	 */
	tconfig->PchSbAccessUnlock = (config->HeciEnabled == 0) ? 1 : 0;
	if (get_lockdown_config() == CHIPSET_LOCKDOWN_COREBOOT) {
		tconfig->PchLockDownBiosInterface = 0;
		params->PchLockDownBiosLock = 0;
		params->PchLockDownSpiEiss = 0;
		/*
		 * Skip Spi Flash Lockdown from inside FSP.
		 * Making this config "0" means FSP won't set the FLOCKDN bit
		 * of SPIBAR + 0x04 (i.e., Bit 15 of BIOS_HSFSTS_CTL).
		 * So, it becomes coreboot's responsibility to set this bit
		 * before end of POST for security concerns.
		 */
		params->SpiFlashCfgLockDown = 0;
	}
	/* only replacing preexisting subsys ID defaults when non-zero */
	if (CONFIG_SUBSYSTEM_VENDOR_ID != 0) {
		params->DefaultSvid = CONFIG_SUBSYSTEM_VENDOR_ID;
		params->PchSubSystemVendorId = CONFIG_SUBSYSTEM_VENDOR_ID;
	}

	if (CONFIG_SUBSYSTEM_DEVICE_ID != 0) {
		params->DefaultSid = CONFIG_SUBSYSTEM_DEVICE_ID;
		params->PchSubSystemId = CONFIG_SUBSYSTEM_DEVICE_ID;
	}

	params->PchPmWolEnableOverride = config->WakeConfigWolEnableOverride;
	params->PchPmPcieWakeFromDeepSx = config->WakeConfigPcieWakeFromDeepSx;
	params->PchPmDeepSxPol = config->PmConfigDeepSxPol;
	params->PchPmSlpS0Enable = config->s0ix_enable;
	params->PchPmSlpS3MinAssert = config->PmConfigSlpS3MinAssert;
	params->PchPmSlpS4MinAssert = config->PmConfigSlpS4MinAssert;
	params->PchPmSlpSusMinAssert = config->PmConfigSlpSusMinAssert;
	params->PchPmSlpAMinAssert = config->PmConfigSlpAMinAssert;
	params->PchPmLpcClockRun = config->PmConfigPciClockRun;
	params->PchPmSlpStrchSusUp = config->PmConfigSlpStrchSusUp;
	params->PchPmPwrBtnOverridePeriod =
				config->PmConfigPwrBtnOverridePeriod;
	params->PchPmPwrCycDur = config->PmConfigPwrCycDur;

	/* Indicate whether platform supports Voltage Margining */
	params->PchPmSlpS0VmEnable = config->PchPmSlpS0VmEnable;

	params->PchSirqEnable = config->serirq_mode != SERIRQ_OFF;
	params->PchSirqMode = config->serirq_mode == SERIRQ_CONTINUOUS;

	params->CpuConfig.Bits.SkipMpInit = !CONFIG_USE_INTEL_FSP_MP_INIT;

	for (i = 0; i < ARRAY_SIZE(config->i2c_voltage); i++)
		params->SerialIoI2cVoltage[i] = config->i2c_voltage[i];

	for (i = 0; i < ARRAY_SIZE(config->domain_vr_config); i++)
		fill_vr_domain_config(params, i, &config->domain_vr_config[i]);

	/* Show SPI controller if enabled in devicetree.cb */
	dev = pcidev_path_on_root(PCH_DEVFN_SPI);
	params->ShowSpiController = dev ? dev->enabled : 0;

	/* Enable xDCI controller if enabled in devicetree and allowed */
	dev = pcidev_path_on_root(PCH_DEVFN_USBOTG);
	if (dev) {
		if (!xdci_can_enable())
			dev->enabled = 0;
		params->XdciEnable = dev->enabled;
	} else {
		params->XdciEnable = 0;
	}

	/* Enable or disable Gaussian Mixture Model in devicetree */
	dev = pcidev_path_on_root(SA_DEVFN_GMM);
	params->GmmEnable = dev ? dev->enabled : 0;

	/*
	 * Send VR specific mailbox commands:
	 * 000b - no VR specific command sent
	 * 001b - VR mailbox command specifically for the MPS IMPV8 VR
	 *	  will be sent
	 * 010b - VR specific command sent for PS4 exit issue
	 * 100b - VR specific command sent for MPS VR decay issue
	 */
	params->SendVrMbxCmd1 = config->SendVrMbxCmd;

	/*
	 * Activates VR mailbox command for Intersil VR C-state issues.
	 * 0 - no mailbox command sent.
	 * 1 - VR mailbox command sent for IA/GT rails only.
	 * 2 - VR mailbox command sent for IA/GT/SA rails.
	 */
	params->IslVrCmd = config->IslVrCmd;

	/* Acoustic Noise Mitigation */
	params->AcousticNoiseMitigation = config->AcousticNoiseMitigation;
	params->SlowSlewRateForIa = config->SlowSlewRateForIa;
	params->SlowSlewRateForGt = config->SlowSlewRateForGt;
	params->SlowSlewRateForSa = config->SlowSlewRateForSa;
	params->FastPkgCRampDisableIa = config->FastPkgCRampDisableIa;
	params->FastPkgCRampDisableGt = config->FastPkgCRampDisableGt;
	params->FastPkgCRampDisableSa = config->FastPkgCRampDisableSa;

	/* Enable PMC XRAM read */
	tconfig->PchPmPmcReadDisable = config->PchPmPmcReadDisable;

	/* Enable/Disable EIST */
	tconfig->Eist = config->eist_enable;

	/* Set TccActivationOffset */
	tconfig->TccActivationOffset = config->tcc_offset;

	/* Enable VT-d and X2APIC */
	if (!config->ignore_vtd && soc_is_vtd_capable()) {
		params->VtdBaseAddress[0] = GFXVT_BASE_ADDRESS;
		params->VtdBaseAddress[1] = VTVC0_BASE_ADDRESS;
		params->X2ApicOptOut = 0;
		tconfig->VtdDisable = 0;

		params->PchIoApicBdfValid = 1;
		params->PchIoApicBusNumber = V_P2SB_IBDF_BUS;
		params->PchIoApicDeviceNumber = V_P2SB_IBDF_DEV;
		params->PchIoApicFunctionNumber = V_P2SB_IBDF_FUN;
	}

	dev = pcidev_path_on_root(SA_DEVFN_IGD);
	if (CONFIG(RUN_FSP_GOP) && dev && dev->enabled)
		params->PeiGraphicsPeimInit = 1;
	else
		params->PeiGraphicsPeimInit = 0;

	soc_irq_settings(params);
}

/* Mainboard GPIO Configuration */
__weak void mainboard_silicon_init_params(FSP_S_CONFIG *params)
{
	printk(BIOS_DEBUG, "WEAK: %s/%s called\n", __FILE__, __func__);
}
