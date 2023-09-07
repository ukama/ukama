/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2006 Uwe Hermann <uwe@hermann-uwe.de>
 * Copyright (C) 2007 Philipp Degler <pdegler@rumms.uni-mannheim.de>
 * Copyright (C) 2017 Gergely Kiss <mail.gery@gmail.com>
 * Copyright (C) 2019 Protectli
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

#include <device/device.h>
#include <device/pnp.h>
#include <pc80/keyboard.h>
#include <superio/conf_mode.h>
#include <superio/ite/common/env_ctrl.h>

#include "chip.h"
#include "it8613e.h"

static void it8613e_init(struct device *dev)
{
	const struct superio_ite_it8613e_config *conf = dev->chip_info;
	const struct resource *res;

	if (!dev->enabled)
		return;

	switch (dev->path.pnp.device) {
	case IT8613E_EC:
		res = find_resource(dev, PNP_IDX_IO0);
		if (!conf || !res)
			break;
		ite_ec_init(res->base, &conf->ec);
		break;
	case IT8613E_KBCK:
		pc_keyboard_init(NO_AUX_DEVICE);
		break;
	case IT8613E_KBCM:
		break;
	}
}

static struct device_operations ops = {
	.read_resources		= pnp_read_resources,
	.set_resources		= pnp_set_resources,
	.enable_resources	= pnp_enable_resources,
	.enable				= pnp_alt_enable,
	.init				= it8613e_init,
	.ops_pnp_mode		= &pnp_conf_mode_870155_aa,
};

static struct pnp_info pnp_dev_info[] = {
	/* Serial Port 1 */
	{ NULL, IT8613E_SP1, PNP_IO0 | PNP_IRQ0 | PNP_MSC0, 0x0ff8, },
	/* Environmental Controller */
	{ NULL, IT8613E_EC, PNP_IO0 | PNP_IO1 | PNP_IRQ0, 0x0ff8, 0x0ff8, },
	/* KBC Keyboard */
	{ NULL, IT8613E_KBCK, PNP_IO0 | PNP_IO1 | PNP_IRQ0 | PNP_MSC0,
		0x0fff, 0x0fff, },
	/* KBC Mouse */
	{ NULL, IT8613E_KBCM, PNP_IRQ0 | PNP_MSC0, },
	/* GPIO */
	{ NULL, IT8613E_GPIO, PNP_IO0 | PNP_IO1 | PNP_IRQ0, 0x0ffc, 0x0ff8, },
	/* Consumer Infrared */
	{ NULL, IT8613E_CIR, PNP_IO0 | PNP_IRQ0 | PNP_MSC0, 0x0ff8, },
};

static void enable_dev(struct device *dev)
{
	pnp_enable_devices(dev, &ops, ARRAY_SIZE(pnp_dev_info), pnp_dev_info);
}

struct chip_operations superio_ite_it8613e_ops = {
	CHIP_NAME("ITE IT8613E Super I/O")
	.enable_dev = enable_dev,
};
