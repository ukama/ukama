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

#include <baseboard/gpio.h>
#include <baseboard/variants.h>
#include <commonlib/helpers.h>

/* Pad configuration in ramstage*/
static const struct pad_config gpio_table[] = {
/* I2S2_SCLK */
PAD_CFG_GPI(GPP_A7, NONE, PLTRST),
/* I2S2_RXD */
PAD_CFG_GPI(GPP_A10, NONE, PLTRST),
/* TCH_PNL2_RST_N */
PAD_CFG_GPO(GPP_A13, 1, DEEP),
/* ONBOARD_X4_PCIE_SLOT1_PWREN_N */
PAD_CFG_GPO(GPP_A14, 0, DEEP),
/* TCH_PNL2_INT_N */
PAD_CFG_GPI_APIC_EDGE_LOW(GPP_B3, NONE, PLTRST),
/* TC_RETIMER_FORCE_PWR */
PAD_CFG_GPO(GPP_B4, 0, DEEP),
/* FPS_RST_N */
PAD_CFG_GPO(GPP_B14, 1, DEEP),
/* WIFI_RF_KILL_N */
PAD_CFG_GPO(GPP_B15, 1, PLTRST),
/* M2_SSD_PWREN_N */
PAD_CFG_GPO(GPP_B16, 1, DEEP),
/* WWAN_PERST_N */
PAD_CFG_GPO(GPP_B17, 1, DEEP),
/* BT_RF_KILL_N */
PAD_CFG_GPO(GPP_B18, 1, PLTRST),
/* CRD_CAM_PWREN_1 */
PAD_CFG_GPO(GPP_B23, 1, PLTRST),
/* WF_CAM_CLK_EN */
PAD_CFG_GPO(GPP_C2, 1, PLTRST),
/* ONBOARD_X4_PCIE_SLOT1_RESET_N */
PAD_CFG_GPO(GPP_C5, 1, DEEP),
/* TCH_PAD_INT_N */
PAD_CFG_GPI_APIC_EDGE_LOW(GPP_C8, NONE, PLTRST),
/* WWAN_RST_N */
PAD_CFG_GPO(GPP_C10, 1, DEEP),
/* WWAN_FCP_OFF_N */
PAD_CFG_GPO(GPP_C11, 1, DEEP),
/* CODEC_INT_N */
PAD_CFG_GPI_APIC_LOW(GPP_C12, NONE, PLTRST),
/* SPKR_PD_N */
PAD_CFG_GPO(GPP_C13, 1, PLTRST),
/* WF_CAM_RST_N */
PAD_CFG_GPO(GPP_C15, 1, PLTRST),
/* CRD_CAM_STROBE_1 */
PAD_CFG_GPO(GPP_C22, 0, PLTRST),
/* CRD_CAM_PRIVACY_LED_1 */
PAD_CFG_GPO(GPP_C23, 0, PLTRST),
/* FLASH_DES_SEC_OVERRIDEs */
PAD_CFG_GPO(GPP_D13, 0, DEEP),
/* TCH_PAD_LS_EN */
PAD_CFG_GPO(GPP_D14, 1, PLTRST),
/* ONBOARD_X4_PCIE_SLOT1_DGPU_SEL */
PAD_CFG_GPO(GPP_D15, 0, DEEP),
/* MFR_MODE_DET_STRAP */
PAD_CFG_GPI(GPP_D16, NONE, PLTRST),
/* TBT_CIO_PWR_EN */
PAD_CFG_GPO(GPP_E0, 1, DEEP),
/* FPS_INT */
PAD_CFG_GPI_APIC(GPP_E3, NONE, PLTRST, LEVEL, NONE),
/* EC_SLP_S0_CS_N */
PAD_CFG_GPO(GPP_E6, 1, DEEP),
/* EC_SMI_N */
PAD_CFG_GPI_SMI(GPP_E7, NONE, DEEP, LEVEL, NONE),
/* TBT_CIO_PLUG_EVENT_N */
PAD_CFG_GPI_SCI(GPP_E17, NONE, DEEP, EDGE_SINGLE, NONE),
/* DISP_AUX_P_BIAS_GPIO */
PAD_CFG_GPO(GPP_E22, 0, PLTRST),
/* DISP_AUX_N_BIAS_GPIO */
PAD_CFG_GPO(GPP_E23, 1, DEEP),
/* SATA_HDD_PWREN */
PAD_CFG_GPO(GPP_F4, 1, PLTRST),
/* BIOS_REC */
PAD_CFG_GPI(GPP_F5, NONE, PLTRST),
/* SD_CD# */
PAD_CFG_NF(GPP_G5, UP_20K, PWROK, NF1),
/* SD_WP */
PAD_CFG_NF(GPP_G7, DN_20K, PWROK, NF1),
/* M2_SSD_RST_N */
PAD_CFG_GPO(GPP_H0, 1, DEEP),
};

/* Early pad configuration in bootblock */
static const struct pad_config early_gpio_table[] = {

};

const struct pad_config *__attribute__((weak)) variant_gpio_table(size_t *num)
{
	*num = ARRAY_SIZE(gpio_table);
	return gpio_table;
}

const struct pad_config *__attribute__((weak))
	variant_early_gpio_table(size_t *num)
{
	*num = ARRAY_SIZE(early_gpio_table);
	return early_gpio_table;
}

static const struct cros_gpio cros_gpios[] = {
	CROS_GPIO_REC_AL(CROS_GPIO_VIRTUAL, CROS_GPIO_DEVICE_NAME),
};

const struct cros_gpio *__attribute__((weak)) variant_cros_gpios(size_t *num)
{
	*num = ARRAY_SIZE(cros_gpios);
	return cros_gpios;
}
