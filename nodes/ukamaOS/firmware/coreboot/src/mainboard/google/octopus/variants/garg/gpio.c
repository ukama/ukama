/*
 * This file is part of the coreboot project.
 *
 * Copyright 2019 Google LLC
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
 * GNU General Public License for more details.
 */

#include <baseboard/gpio.h>
#include <baseboard/variants.h>
#include <boardid.h>
#include <gpio.h>
#include <soc/gpio.h>

enum {
	SKU_1_2A2C = 1,
	SKU_9_HDMI = 9,
	SKU_17_LTE = 17,
	SKU_18_LTE_TS = 18,
	SKU_37_2A2C_360 = 37,
};

static const struct pad_config default_override_table[] = {
	PAD_NC(GPIO_104, UP_20K),

	/* EN_PP3300_TOUCHSCREEN */
	PAD_CFG_GPO_IOSSTATE_IOSTERM(GPIO_146, 0, DEEP, NONE, Tx0RxDCRx0,
				     DISPUPD),

	PAD_NC(GPIO_213, DN_20K),
};

static const struct pad_config hdmi_override_table[] = {
	PAD_NC(GPIO_104, UP_20K),

	/* HV_DDI1_DDC_SDA */
	PAD_CFG_NF_IOSSTATE_IOSTERM(GPIO_126, NONE, DEEP, NF1, HIZCRx1,
				    DISPUPD),
	/* HV_DDI1_DDC_SCL */
	PAD_CFG_NF_IOSSTATE_IOSTERM(GPIO_127, NONE, DEEP, NF1, HIZCRx1,
				    DISPUPD),

	/* EN_PP3300_TOUCHSCREEN */
	PAD_CFG_GPO_IOSSTATE_IOSTERM(GPIO_146, 0, DEEP, NONE, Tx0RxDCRx0,
				     DISPUPD),

	PAD_NC(GPIO_213, DN_20K),

};

static const struct pad_config lte_override_table[] = {
	/* Default override table. */
	PAD_NC(GPIO_104, UP_20K),

	/* EN_PP3300_TOUCHSCREEN */
	PAD_CFG_GPO_IOSSTATE_IOSTERM(GPIO_146, 0, DEEP, NONE, Tx0RxDCRx0,
				     DISPUPD),

	PAD_NC(GPIO_213, DN_20K),

	/* Be specific to LTE SKU */
	/* UART2-CTS_B -- EN_PP3300_DX_LTE_SOC */
	PAD_CFG_GPO(GPIO_67, 1, PWROK),

	/* PCIE_WAKE1_B -- FULL_CARD_POWER_OFF */
	PAD_CFG_GPO(GPIO_117, 1, PWROK),

	/* AVS_I2S1_MCLK -- PLT_RST_LTE_L */
	PAD_CFG_GPO(GPIO_161, 1, DEEP),
};

const struct pad_config *variant_override_gpio_table(size_t *num)
{
	uint32_t sku_id;
	sku_id = get_board_sku();

	switch (sku_id) {
	case SKU_9_HDMI:
		*num = ARRAY_SIZE(hdmi_override_table);
		return hdmi_override_table;
	case SKU_17_LTE:
	case SKU_18_LTE_TS:
		*num = ARRAY_SIZE(lte_override_table);
		return lte_override_table;
	default:
		*num = ARRAY_SIZE(default_override_table);
		return default_override_table;
	}
}

static const struct pad_config lte_early_override_table[] = {
	/* UART2-CTS_B -- EN_PP3300_DX_LTE_SOC */
	PAD_CFG_GPO(GPIO_67, 1, PWROK),

	/* PCIE_WAKE1_B -- FULL_CARD_POWER_OFF */
	PAD_CFG_GPO(GPIO_117, 1, PWROK),

	/* AVS_I2S1_MCLK -- PLT_RST_LTE_L */
	PAD_CFG_GPO(GPIO_161, 0, DEEP),
};

const struct pad_config *variant_early_override_gpio_table(size_t *num)
{
	*num = ARRAY_SIZE(lte_early_override_table);

	return lte_early_override_table;
}
