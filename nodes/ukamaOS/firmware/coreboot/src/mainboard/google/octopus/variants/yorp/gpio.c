/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Google LLC
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
#include <gpio.h>
#include <soc/gpio.h>

/* GPIOs needed prior to ramstage. */
static const struct pad_config early_gpio_table[] = {
	PAD_CFG_GPI(GPIO_190, NONE, DEEP), /* PCH_WP_OD */
	/* GSPI0_INT */
	PAD_CFG_GPI_APIC_IOS(GPIO_63, NONE, DEEP, LEVEL, INVERT, TxDRxE,
		DISPUPD), /* H1_PCH_INT_ODL */
	/* GSPI0_CLK */
	PAD_CFG_NF(GPIO_79, NONE, DEEP, NF1), /* H1_SLAVE_SPI_CLK_R */
	/* GSPI0_CS# */
	PAD_CFG_NF(GPIO_80, NONE, DEEP, NF1), /* H1_SLAVE_SPI_CS_L_R */
	/* GSPI0_MISO */
	PAD_CFG_NF(GPIO_82, NONE, DEEP, NF1), /* H1_SLAVE_SPI_MISO */
	/* GSPI0_MOSI */
	PAD_CFG_NF(GPIO_83, NONE, DEEP, NF1), /* H1_SLAVE_SPI_MOSI_R */

	/* Enable power to wifi early in bootblock and de-assert PERST#. */
	PAD_CFG_GPO(GPIO_178, 1, DEEP), /* EN_PP3300_WLAN */
	PAD_CFG_GPO(GPIO_164, 0, DEEP), /* WLAN_PE_RST */

	/*
	 * ESPI_IO1 acts as ALERT# (which is open-drain) and requies a weak
	 * pull-up for proper operation. Since there is no external pull present
	 * on this platform, configure an internal weak pull-up.
	 */
	PAD_CFG_NF_IOSSTATE_IOSTERM(GPIO_151, UP_20K, DEEP, NF2, HIZCRx1,
				    ENPU), /* ESPI_IO1 */
};

const struct pad_config *variant_early_gpio_table(size_t *num)
{
	*num = ARRAY_SIZE(early_gpio_table);
	return early_gpio_table;
}
