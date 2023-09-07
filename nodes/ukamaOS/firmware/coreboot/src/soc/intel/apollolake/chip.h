/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2015 Intel Corp.
 * Copyright (C) 2017 - 2018 Siemens AG
 * (Written by Alexandru Gagniuc <alexandrux.gagniuc@intel.com> for Intel Corp.)
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

#ifndef _SOC_APOLLOLAKE_CHIP_H_
#define _SOC_APOLLOLAKE_CHIP_H_

#include <commonlib/helpers.h>
#include <intelblocks/cfg.h>
#include <intelblocks/gspi.h>
#include <soc/gpe.h>
#include <soc/gpio.h>
#include <intelblocks/lpc_lib.h>
#include <device/i2c_simple.h>
#include <drivers/i2c/designware/dw_i2c.h>
#include <soc/pm.h>
#include <soc/usb.h>

#define MAX_PCIE_PORTS			6
#define CLKREQ_DISABLED		0xf

enum pnp_settings {
	PNP_PERF,
	PNP_POWER,
	PNP_PERF_POWER,
};

struct soc_intel_apollolake_config {

	/* Common structure containing soc config data required by common code*/
	struct soc_intel_common_config common_soc_config;

	/*
	 * Mapping from PCIe root port to CLKREQ input on the SOC. The SOC has
	 * four CLKREQ inputs, but six root ports. Root ports without an
	 * associated CLKREQ signal must be marked with "CLKREQ_DISABLED"
	 */
	uint8_t pcie_rp_clkreq_pin[MAX_PCIE_PORTS];

	/* Enable/disable hot-plug for root ports (0 = disable, 1 = enable). */
	uint8_t pcie_rp_hotplug_enable[MAX_PCIE_PORTS];

	/* De-emphasis enable configuration for each PCIe root port */
	uint8_t pcie_rp_deemphasis_enable[MAX_PCIE_PORTS];

	/* [14:8] DDR mode Number of dealy elements.Each = 125pSec.
	 * [6:0] SDR mode Number of dealy elements.Each = 125pSec.
	 */
	uint32_t emmc_tx_cmd_cntl;

	/* [14:8] HS400 mode Number of dealy elements.Each = 125pSec.
	 * [6:0] SDR104/HS200 mode Number of dealy elements.Each = 125pSec.
	 */
	uint32_t emmc_tx_data_cntl1;

	/* [30:24] SDR50 mode Number of dealy elements.Each = 125pSec.
	 * [22:16] DDR50 mode Number of dealy elements.Each = 125pSec.
	 * [14:8] SDR25/HS50 mode Number of dealy elements.Each = 125pSec.
	 * [6:0] SDR12/Compatibility mode Number of dealy elements.
	 *       Each = 125pSec.
	 */
	uint32_t emmc_tx_data_cntl2;

	/* [30:24] SDR50 mode Number of dealy elements.Each = 125pSec.
	 * [22:16] DDR50 mode Number of dealy elements.Each = 125pSec.
	 * [14:8] SDR25/HS50 mode Number of dealy elements.Each = 125pSec.
	 * [6:0] SDR12/Compatibility mode Number of dealy elements.
	 *       Each = 125pSec.
	 */
	uint32_t emmc_rx_cmd_data_cntl1;

	/* [14:8] HS400 mode 1 Number of dealy elements.Each = 125pSec.
	 * [6:0] HS400 mode 2 Number of dealy elements.Each = 125pSec.
	 */
	uint32_t emmc_rx_strobe_cntl;

	/* [13:8] Auto Tuning mode Number of dealy elements.Each = 125pSec.
	 * [6:0] SDR104/HS200 Number of dealy elements.Each = 125pSec.
	 */
	uint32_t emmc_rx_cmd_data_cntl2;

	/* Select the eMMC max speed allowed. */
	uint8_t emmc_host_max_speed;

	/* Specifies on which IRQ the SCI will internally appear. */
	uint8_t sci_irq;

	/* Configure serial IRQ (SERIRQ) line. */
	enum serirq_mode serirq_mode;

	uint8_t gpe0_dw1; /* GPE0_63_32 STS/EN */
	uint8_t gpe0_dw2; /* GPE0_95_64 STS/EN */
	uint8_t gpe0_dw3; /* GPE0_127_96 STS/EN */

	/* Configure LPSS S0ix Enable */
	uint8_t lpss_s0ix_enable;

	/* Enable DPTF support */
	int dptf_enable;

	/* TCC activation offset value in degrees Celsius */
	int tcc_offset;

	/* PL1 override value in mW for APL */
	uint16_t tdp_pl1_override_mw;
	/* PL2 override value in mW for APL */
	uint16_t tdp_pl2_override_mw;

	/* Configure Audio clk gate and power gate
	 * IOSF-SB port ID 92 offset 0x530 [5] and [3]
	 */
	uint8_t hdaudio_clk_gate_enable;
	uint8_t hdaudio_pwr_gate_enable;
	uint8_t hdaudio_bios_config_lockdown;

	/* SLP S3 minimum assertion width. */
	int slp_s3_assertion_width_usecs;

	/* GPIO pin for PERST_0 */
	uint16_t prt0_gpio;

	/* USB2 eye diagram settings per port */
	struct usb2_eye_per_port usb2eye[APOLLOLAKE_USB2_PORT_MAX];

	/* GPIO SD card detect pin */
	unsigned int sdcard_cd_gpio;

	/* Select PNP Settings.
	 * (0) Performance,
	 * (1) Power
	 * (2) Power & Performance */
	enum pnp_settings pnp_settings;

	/* PMIC PCH_PWROK delay configuration - IPC Configuration
	 * Upd for changing PCH_PWROK delay configuration : I2C_Slave_Address
	 * (31:24) + Register_Offset (23:16) + OR Value (15:8) + AND Value (7:0)
	 */
	uint32_t PmicPmcIpcCtrl;

	/* Options to disable XHCI Link Compliance Mode. Default is FALSE to not
	 * disable Compliance Mode. Set TRUE to disable Compliance Mode.
	 * 0:FALSE(Default), 1:True.
	 */
	uint8_t DisableComplianceMode;

	/* Options to change USB3 ModPhy setting for the Integrated Filter (IF)
	 * value. Default is 0 to not changing default IF value (0x12). Set
	 * value with the range from 0x01 to 0xff to change IF value.
	 */
	uint8_t ModPhyIfValue;

	/* Options to bump USB3 LDO voltage. Default is FALSE to not increasing
	 * LDO voltage. Set TRUE to increase LDO voltage with 40mV.
	 * 0:FALSE (default), 1:True.
	 */
	uint8_t ModPhyVoltageBump;

	/* Options to adjust PMIC Vdd2 voltage. Default is 0 to not adjusting
	 * the PMIC Vdd2 default voltage 1.20v. Upd for changing Vdd2 Voltage
	 * configuration: I2C_Slave_Address (31:23) + Register_Offset (23:16)
	 * + OR Value (15:8) + AND Value (7:0) through BUCK5_VID[3:2]:
	 * 00=1.10v, 01=1.15v, 10=1.24v, 11=1.20v (default).
	 */
	uint32_t PmicVdd2Voltage;

	/* Option to enable VTD feature. Default is 0 which disables VTD
	 * capability in FSP. Setting this option to 1 in devicetree will enable
	 * the Upd parameter VtdEnable.
	 */
	uint8_t enable_vtd;
};

typedef struct soc_intel_apollolake_config config_t;

#endif	/* _SOC_APOLLOLAKE_CHIP_H_ */
