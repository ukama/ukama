/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corporation.
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

#include <console/console.h>
#include <device/device.h>
#include <device/pci.h>
#include <fsp/api.h>
#include <fsp/util.h>
#include <intelblocks/lpss.h>
#include <intelblocks/xdci.h>
#include <soc/intel/common/vbt.h>
#include <soc/pci_devs.h>
#include <soc/ramstage.h>
#include <soc/soc_chip.h>
#include <string.h>
#include <intelblocks/mp_init.h>
#include <fsp/ppi/mp_service_ppi.h>

static void parse_devicetree(FSP_S_CONFIG *params)
{
	const struct soc_intel_icelake_config *config;
	config = config_of_soc();

	for (int i = 0; i < CONFIG_SOC_INTEL_I2C_DEV_MAX; i++)
		params->SerialIoI2cMode[i] = config->SerialIoI2cMode[i];

	for (int i = 0; i < CONFIG_SOC_INTEL_COMMON_BLOCK_GSPI_MAX; i++) {
		params->SerialIoSpiMode[i] = config->SerialIoGSpiMode[i];
		params->SerialIoSpiCsMode[i] = config->SerialIoGSpiCsMode[i];
		params->SerialIoSpiCsState[i] = config->SerialIoGSpiCsState[i];
	}

	for (int i = 0; i < CONFIG_SOC_INTEL_UART_DEV_MAX; i++)
		params->SerialIoUartMode[i] = config->SerialIoUartMode[i];
}

static const pci_devfn_t serial_io_dev[] = {
	PCH_DEVFN_I2C0,
	PCH_DEVFN_I2C1,
	PCH_DEVFN_I2C2,
	PCH_DEVFN_I2C3,
	PCH_DEVFN_I2C4,
	PCH_DEVFN_I2C5,
	PCH_DEVFN_GSPI0,
	PCH_DEVFN_GSPI1,
	PCH_DEVFN_GSPI2,
	PCH_DEVFN_UART0,
	PCH_DEVFN_UART1,
	PCH_DEVFN_UART2
};

/* UPD parameters to be initialized before SiliconInit */
void platform_fsp_silicon_init_params_cb(FSPS_UPD *supd)
{
	int i;
	FSP_S_CONFIG *params = &supd->FspsConfig;

	struct device *dev;
	struct soc_intel_icelake_config *config;
	config = config_of_soc();

	/* Parse device tree and enable/disable devices */
	parse_devicetree(params);

	/* Load VBT before devicetree-specific config. */
	params->GraphicsConfigPtr = (uintptr_t)vbt_get();

	/* Set USB OC pin to 0 first */
	for (i = 0; i < ARRAY_SIZE(params->Usb2OverCurrentPin); i++)
		params->Usb2OverCurrentPin[i] = 0;

	for (i = 0; i < ARRAY_SIZE(params->Usb3OverCurrentPin); i++)
		params->Usb3OverCurrentPin[i] = 0;

	if (CONFIG(USE_INTEL_FSP_TO_CALL_COREBOOT_PUBLISH_MP_PPI)) {
		params->CpuMpPpi = (uintptr_t) mp_fill_ppi_services_data();
		params->SkipMpInit = 0;
	} else {
		params->SkipMpInit = !CONFIG_USE_INTEL_FSP_MP_INIT;
	}

	mainboard_silicon_init_params(params);

	dev = pcidev_path_on_root(SA_DEVFN_IGD);
	if (CONFIG(RUN_FSP_GOP) && dev && dev->enabled)
		params->PeiGraphicsPeimInit = 1;
	else
		params->PeiGraphicsPeimInit = 0;
	if (dev && dev->enabled) {
		params->GtFreqMax = 2;
		params->CdClock = 3;
	}

	/* Unlock upper 8 bytes of RTC RAM */
	params->PchLockDownRtcMemoryLock = 0;

	params->CnviBtAudioOffload = config->CnviBtAudioOffload;
	/* SATA */
	dev = pcidev_on_root(PCH_DEV_SLOT_SATA, 0);
	if (!dev)
		params->SataEnable = 0;
	else {
		params->SataEnable = dev->enabled;
		params->SataMode = config->SataMode;
		params->SataSalpSupport = config->SataSalpSupport;
		memcpy(params->SataPortsEnable, config->SataPortsEnable,
				sizeof(params->SataPortsEnable));
		memcpy(params->SataPortsDevSlp, config->SataPortsDevSlp,
				sizeof(params->SataPortsDevSlp));
	}

	/* Lan */
	dev = pcidev_on_root(PCH_DEV_SLOT_ESPI, 6);
	if (!dev)
		params->PchLanEnable = 0;
	else
		params->PchLanEnable = dev->enabled;

	/* Audio */
	params->PchHdaDspEnable = config->PchHdaDspEnable;
	params->PchHdaAudioLinkHda = config->PchHdaAudioLinkHda;
	params->PchHdaAudioLinkDmic0 = config->PchHdaAudioLinkDmic0;
	params->PchHdaAudioLinkDmic1 = config->PchHdaAudioLinkDmic1;
	params->PchHdaAudioLinkSsp0 = config->PchHdaAudioLinkSsp0;
	params->PchHdaAudioLinkSsp1 = config->PchHdaAudioLinkSsp1;
	params->PchHdaAudioLinkSsp2 = config->PchHdaAudioLinkSsp2;
	params->PchHdaAudioLinkSndw1 = config->PchHdaAudioLinkSndw1;
	params->PchHdaAudioLinkSndw2 = config->PchHdaAudioLinkSndw2;
	params->PchHdaAudioLinkSndw3 = config->PchHdaAudioLinkSndw3;
	params->PchHdaAudioLinkSndw4 = config->PchHdaAudioLinkSndw4;

	/* disable Legacy PME */
	memset(params->PcieRpPmSci, 0, sizeof(params->PcieRpPmSci));

	/* Legacy 8254 timer support */
	params->Enable8254ClockGating = !CONFIG_USE_LEGACY_8254_TIMER;
	params->Enable8254ClockGatingOnS3 = 1;

	/* S0ix */
	params->PchPmSlpS0Enable = config->s0ix_enable;

	/* USB */
	for (i = 0; i < ARRAY_SIZE(config->usb2_ports); i++) {
		params->PortUsb20Enable[i] =
			config->usb2_ports[i].enable;
		params->Usb2OverCurrentPin[i] =
			config->usb2_ports[i].ocpin;
		params->Usb2PhyPetxiset[i] =
			config->usb2_ports[i].pre_emp_bias;
		params->Usb2PhyTxiset[i] =
			config->usb2_ports[i].tx_bias;
		params->Usb2PhyPredeemp[i] =
			config->usb2_ports[i].tx_emp_enable;
		params->Usb2PhyPehalfbit[i] =
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

	/* Enable xDCI controller if enabled in devicetree and allowed */
	dev = pcidev_on_root(PCH_DEV_SLOT_XHCI, 1);
	if (!xdci_can_enable())
		dev->enabled = 0;
	params->XdciEnable = dev->enabled;

	/* PCI Express */
	for (i = 0; i < ARRAY_SIZE(config->PcieClkSrcUsage); i++) {
		if (config->PcieClkSrcUsage[i] == 0)
			config->PcieClkSrcUsage[i] = PCIE_CLK_NOTUSED;
	}
	memcpy(params->PcieClkSrcUsage, config->PcieClkSrcUsage,
	       sizeof(config->PcieClkSrcUsage));
	memcpy(params->PcieClkSrcClkReq, config->PcieClkSrcClkReq,
	       sizeof(config->PcieClkSrcClkReq));

	/* eMMC */
	dev = pcidev_on_root(PCH_DEV_SLOT_STORAGE, 0);
	if (!dev)
		params->ScsEmmcEnabled = 0;
	else {
		params->ScsEmmcEnabled = dev->enabled;
		params->ScsEmmcHs400Enabled = config->ScsEmmcHs400Enabled;
		params->EmmcUseCustomDlls = config->EmmcUseCustomDlls;
		if (config->EmmcUseCustomDlls == 1) {
			params->EmmcTxCmdDelayRegValue =
					config->EmmcTxCmdDelayRegValue;
			params->EmmcTxDataDelay1RegValue =
					config->EmmcTxDataDelay1RegValue;
			params->EmmcTxDataDelay2RegValue =
					config->EmmcTxDataDelay2RegValue;
			params->EmmcRxCmdDataDelay1RegValue =
					config->EmmcRxCmdDataDelay1RegValue;
			params->EmmcRxCmdDataDelay2RegValue =
					config->EmmcRxCmdDataDelay2RegValue;
			params->EmmcRxStrobeDelayRegValue =
					config->EmmcRxStrobeDelayRegValue;
		}
	}

	/* SD */
	dev = pcidev_on_root(PCH_DEV_SLOT_XHCI, 5);
	if (!dev)
		params->ScsSdCardEnabled = 0;
	else {
		params->ScsSdCardEnabled = dev->enabled;
		params->SdCardPowerEnableActiveHigh =
				config->SdCardPowerEnableActiveHigh;
	}

	params->Heci3Enabled = config->Heci3Enabled;
	params->Device4Enable = config->Device4Enable;
}

/* Mainboard GPIO Configuration */
__weak void mainboard_silicon_init_params(FSP_S_CONFIG *params)
{
	printk(BIOS_DEBUG, "WEAK: %s/%s called\n", __FILE__, __func__);
}

/* Return list of SOC LPSS controllers */
const pci_devfn_t *soc_lpss_controllers_list(size_t *size)
{
	*size = ARRAY_SIZE(serial_io_dev);
	return serial_io_dev;
}
