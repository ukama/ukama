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
#include <boardid.h>
#include <gpio.h>
#include <soc/gpio.h>

static const struct pad_config default_override_table[] = {

	PAD_NC(GPIO_52, UP_20K),
	PAD_NC(GPIO_53, UP_20K),
	PAD_NC(GPIO_67, UP_20K),
	PAD_NC(GPIO_117, UP_20K),

	PAD_NC(GPIO_143, UP_20K),
	PAD_NC(GPIO_144, UP_20K),
	PAD_NC(GPIO_145, UP_20K),

	/* EN_PP3300_TOUCHSCREEN */
	PAD_CFG_GPO_IOSSTATE_IOSTERM(GPIO_146, 0, DEEP, NONE, Tx0RxDCRx0,
					DISPUPD),

	PAD_NC(GPIO_161, UP_20K),

	PAD_NC(GPIO_213, DN_20K),
	PAD_NC(GPIO_214, DN_20K),
};

const struct pad_config *variant_override_gpio_table(size_t *num)
{
	*num = ARRAY_SIZE(default_override_table);

	return default_override_table;
}
