/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2015 Google Inc.
 * Copyright (C) 2015 Intel Corporation
 * Copyright (C) 2017 Advanced Micro Devices, Inc.
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

#include <stdint.h>
#include <amdblocks/gpio_banks.h>
#include <amdblocks/acpimmio.h>
#include <soc/gpio.h>
#include <soc/smi.h>

static const struct soc_amd_event gpio_event_table[] = {
	{ GPIO_1, GEVENT_19 },
	{ GPIO_2, GEVENT_8 },
	{ GPIO_3, GEVENT_2 },
	{ GPIO_4, GEVENT_4 },
	{ GPIO_5, GEVENT_7 },
	{ GPIO_6, GEVENT_10 },
	{ GPIO_7, GEVENT_11 },
	{ GPIO_8, GEVENT_23 },
	{ GPIO_9, GEVENT_22 },
	{ GPIO_11, GEVENT_18 },
	{ GPIO_13, GEVENT_21 },
	{ GPIO_14, GEVENT_6 },
	{ GPIO_15, GEVENT_20 },
	{ GPIO_16, GEVENT_12 },
	{ GPIO_17, GEVENT_13 },
	{ GPIO_18, GEVENT_14 },
	{ GPIO_21, GEVENT_5 },
	{ GPIO_22, GEVENT_3 },
	{ GPIO_23, GEVENT_16 },
	{ GPIO_24, GEVENT_15 },
	{ GPIO_65, GEVENT_0 },
	{ GPIO_66, GEVENT_1 },
	{ GPIO_68, GEVENT_9 },
	{ GPIO_69, GEVENT_17 },
};

void soc_route_sci(uint8_t event)
{
	smi_write8(SMI_SCI_MAP(event), event);
}

void soc_get_gpio_event_table(const struct soc_amd_event **table, size_t *items)
{
	*table = gpio_event_table;
	*items = ARRAY_SIZE(gpio_event_table);
}

void soc_gpio_hook(uint8_t gpio, uint8_t mux)
{
	/* Always program Gevent when WAKE_L_AGPIO2 is configured as WAKE_L */
	if ((gpio == 2) && !(mux & AMD_GPIO_MUX_MASK))
		soc_route_sci(GPIO_2_EVENT);
}
