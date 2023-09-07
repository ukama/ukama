/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Intel Corp.
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

#include <intelblocks/gpio.h>
#include <intelblocks/pcr.h>
#include <soc/pcr_ids.h>
#include <soc/pmc.h>

static const struct reset_mapping rst_map[] = {
	{ .logical = PAD_CFG0_LOGICAL_RESET_RSMRST, .chipset = 0U << 30 },
	{ .logical = PAD_CFG0_LOGICAL_RESET_DEEP, .chipset = 1U << 30 },
	{ .logical = PAD_CFG0_LOGICAL_RESET_PLTRST, .chipset = 2U << 30 },
};

static const struct reset_mapping rst_map_com0[] = {
	{ .logical = PAD_CFG0_LOGICAL_RESET_PWROK, .chipset = 0U << 30 },
	{ .logical = PAD_CFG0_LOGICAL_RESET_DEEP, .chipset = 1U << 30 },
	{ .logical = PAD_CFG0_LOGICAL_RESET_PLTRST, .chipset = 2U << 30 },
	{ .logical = PAD_CFG0_LOGICAL_RESET_RSMRST, .chipset = 3U << 30 },
};

/*
 * The GPIO driver for Cannonlake on Windows/Linux expects 32 GPIOs per pad
 * group, regardless of whether or not there is a physical pad for each
 * exposed GPIO number.
 *
 * This results in the OS having a sparse GPIO map, and devices that need
 * to export an ACPI GPIO must use the OS expected number.
 *
 * Not all pins are usable as GPIO and those groups do not have a pad base.
 *
 * This layout matches the Linux kernel pinctrl map for CNL-H at:
 * linux/drivers/pinctrl/intel/pinctrl-cannonlake.c
 */
static const struct pad_group cnl_community0_groups[] = {
	INTEL_GPP_BASE(GPP_A0, GPP_A0, GPIO_RSVD_0, 0),		/* GPP_A */
	INTEL_GPP_BASE(GPP_A0, GPP_B0, GPIO_RSVD_2, 32),	/* GPP_B */
};

static const struct pad_group cnl_community1_groups[] = {
	INTEL_GPP_BASE(GPP_C0, GPP_C0, GPP_C23, 64),		/* GPP_C */
	INTEL_GPP_BASE(GPP_C0, GPP_D0, GPP_D23, 96),		/* GPP_D */
	INTEL_GPP_BASE(GPP_C0, GPP_G0, GPP_G7, 128),		/* GPP_G */
	INTEL_GPP(GPP_C0, GPIO_RSVD_3, GPIO_RSVD_10),		/* AZA */
	INTEL_GPP_BASE(GPP_C0, CNV_BTEN, vISH_UART1_RTS_B, 160),/* VGPIO_0 */
	INTEL_GPP(GPP_C0, vCNV_BT_I2S_BCLK, vSSP2_RXD),		/* VGPIO_1 */
};

/* This community is not visible to the OS */
static const struct pad_group cnl_community2_groups[] = {
	INTEL_GPP(GPD0, GPD0, GPD11),				/* GPD */
};

static const struct pad_group cnl_community3_groups[] = {
	INTEL_GPP_BASE(GPP_K0, GPP_K0, GPP_K23, 192),		/* GPP_K */
	INTEL_GPP_BASE(GPP_K0, GPP_H0, GPP_H23, 224),		/* GPP_H */
	INTEL_GPP_BASE(GPP_K0, GPP_E0, GPP_E12, 256),		/* GPP_E */
	INTEL_GPP_BASE(GPP_K0, GPP_F0, GPP_F23, 288),		/* GPP_F */
	INTEL_GPP(GPP_K0, GPIO_RSVD_11, GPIO_RSVD_19),		/* SPI */
};

static const struct pad_group cnl_community4_groups[] = {
	INTEL_GPP(GPIO_RSVD_20, GPIO_RSVD_20, GPIO_RSVD_30),	/* CPU */
	INTEL_GPP(GPIO_RSVD_20, GPIO_RSVD_31, GPIO_RSVD_39),	/* JTAG */
	INTEL_GPP_BASE(GPIO_RSVD_20, GPP_I0, GPP_I14, 320),	/* GPP_I */
	INTEL_GPP_BASE(GPIO_RSVD_20, GPP_J0, GPP_J11, 352),	/* GPP_J */
};

static const struct pad_community cnl_communities[] = {
	/* GPP A, B */
	[COMM_0] = {
		.port = PID_GPIOCOM0,
		.first_pad = GPP_A0,
		.last_pad = GPIO_RSVD_2,
		.num_gpi_regs = NUM_GPIO_COM0_GPI_REGS,
		.pad_cfg_base = PAD_CFG_BASE,
		.host_own_reg_0 = HOSTSW_OWN_REG_0,
		.gpi_int_sts_reg_0 = GPI_INT_STS_0,
		.gpi_int_en_reg_0 = GPI_INT_EN_0,
		.gpi_smi_sts_reg_0 = GPI_SMI_STS_0,
		.gpi_smi_en_reg_0 = GPI_SMI_EN_0,
		.max_pads_per_group = GPIO_MAX_NUM_PER_GROUP,
		.name = "GPP_AB",
		.acpi_path = "\\_SB.PCI0.GPIO",
		.reset_map = rst_map_com0,
		.num_reset_vals = ARRAY_SIZE(rst_map_com0),
		.groups = cnl_community0_groups,
		.num_groups = ARRAY_SIZE(cnl_community0_groups),
	},
	/* GPP C, D, G */
	[COMM_1] = {
		.port = PID_GPIOCOM1,
		.first_pad = GPP_C0,
		.last_pad = vSSP2_RXD,
		.num_gpi_regs = NUM_GPIO_COM1_GPI_REGS,
		.pad_cfg_base = PAD_CFG_BASE,
		.host_own_reg_0 = HOSTSW_OWN_REG_0,
		.gpi_int_sts_reg_0 = GPI_INT_STS_0,
		.gpi_int_en_reg_0 = GPI_INT_EN_0,
		.gpi_smi_sts_reg_0 = GPI_SMI_STS_0,
		.gpi_smi_en_reg_0 = GPI_SMI_EN_0,
		.max_pads_per_group = GPIO_MAX_NUM_PER_GROUP,
		.name = "GPP_CDG",
		.acpi_path = "\\_SB.PCI0.GPIO",
		.reset_map = rst_map,
		.num_reset_vals = ARRAY_SIZE(rst_map),
		.groups = cnl_community1_groups,
		.num_groups = ARRAY_SIZE(cnl_community1_groups),
	},
	/* GPD */
	[COMM_2] = {
		.port = PID_GPIOCOM2,
		.first_pad = GPD0,
		.last_pad = GPD11,
		.num_gpi_regs = NUM_GPIO_COM2_GPI_REGS,
		.pad_cfg_base = PAD_CFG_BASE,
		.host_own_reg_0 = HOSTSW_OWN_REG_0,
		.gpi_int_sts_reg_0 = GPI_INT_STS_0,
		.gpi_int_en_reg_0 = GPI_INT_EN_0,
		.gpi_smi_sts_reg_0 = GPI_SMI_STS_0,
		.gpi_smi_en_reg_0 = GPI_SMI_EN_0,
		.max_pads_per_group = GPIO_MAX_NUM_PER_GROUP,
		.name = "GPD",
		.acpi_path = "\\_SB.PCI0.GPIO",
		.reset_map = rst_map,
		.num_reset_vals = ARRAY_SIZE(rst_map),
		.groups = cnl_community2_groups,
		.num_groups = ARRAY_SIZE(cnl_community2_groups),
	},
	/* GPP K, H, E, F */
	[COMM_3] = {
		.port = PID_GPIOCOM3,
		.first_pad = GPP_K0,
		.last_pad = GPIO_RSVD_19,
		.num_gpi_regs = NUM_GPIO_COM3_GPI_REGS,
		.pad_cfg_base = PAD_CFG_BASE,
		.host_own_reg_0 = HOSTSW_OWN_REG_0,
		.gpi_int_sts_reg_0 = GPI_INT_STS_0,
		.gpi_int_en_reg_0 = GPI_INT_EN_0,
		.gpi_smi_sts_reg_0 = GPI_SMI_STS_0,
		.gpi_smi_en_reg_0 = GPI_SMI_EN_0,
		.max_pads_per_group = GPIO_MAX_NUM_PER_GROUP,
		.name = "GPP_KHEF",
		.acpi_path = "\\_SB.PCI0.GPIO",
		.reset_map = rst_map,
		.num_reset_vals = ARRAY_SIZE(rst_map),
		.groups = cnl_community3_groups,
		.num_groups = ARRAY_SIZE(cnl_community3_groups),
	},
	/* GPP I, J */
	[COMM_4] = {
		.port = PID_GPIOCOM4,
		.first_pad = GPIO_RSVD_20,
		.last_pad = GPP_J11,
		.num_gpi_regs = NUM_GPIO_COM4_GPI_REGS,
		.pad_cfg_base = PAD_CFG_BASE,
		.host_own_reg_0 = HOSTSW_OWN_REG_0,
		.gpi_int_sts_reg_0 = GPI_INT_STS_0,
		.gpi_int_en_reg_0 = GPI_INT_EN_0,
		.gpi_smi_sts_reg_0 = GPI_SMI_STS_0,
		.gpi_smi_en_reg_0 = GPI_SMI_EN_0,
		.max_pads_per_group = GPIO_MAX_NUM_PER_GROUP,
		.name = "GPP_IJ",
		.acpi_path = "\\_SB.PCI0.GPIO",
		.reset_map = rst_map,
		.num_reset_vals = ARRAY_SIZE(rst_map),
		.groups = cnl_community4_groups,
		.num_groups = ARRAY_SIZE(cnl_community4_groups),
	}
};

const struct pad_community *soc_gpio_get_community(size_t *num_communities)
{
	*num_communities = ARRAY_SIZE(cnl_communities);
	return cnl_communities;
}

const struct pmc_to_gpio_route *soc_pmc_gpio_routes(size_t *num)
{
	static const struct pmc_to_gpio_route routes[] = {
		{ PMC_GPP_A, GPP_A },
		{ PMC_GPP_B, GPP_B },
		{ PMC_GPP_C, GPP_C },
		{ PMC_GPP_D, GPP_D },
		{ PMC_GPP_E, GPP_E },
		{ PMC_GPP_F, GPP_F },
		{ PMC_GPP_G, GPP_G },
		{ PMC_GPP_H, GPP_H },
		{ PMC_GPP_I, GPP_I },
		{ PMC_GPP_J, GPP_J },
		{ PMC_GPP_K, GPP_K },
		{ PMC_GPD, GPD },
	};
	*num = ARRAY_SIZE(routes);
	return routes;
}
