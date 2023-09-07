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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <arch/io.h>
#include <endian.h>
#include <device/pnp_ops.h>
#include <device/pnp_def.h>

#include "bootblock.h"

#define PNP_CFG_IDX	0x2e
#define PNP_LDN_SERIAL	0x0d

static void pnp_enter_conf_state(pnp_devfn_t dev)
{
	outb(0x55, PNP_CFG_IDX);
	outb(0x55, PNP_CFG_IDX);
}

static void pnp_exit_conf_state(pnp_devfn_t dev)
{
	outb(0xaa, PNP_CFG_IDX);
}

static void wilco_ec_serial_init(void)
{
	pnp_devfn_t dev = PNP_DEV(PNP_CFG_IDX, PNP_LDN_SERIAL);

	pnp_enter_conf_state(dev);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 1);
	pnp_set_iobase(dev, PNP_IDX_IO1, cpu_to_be16(CONFIG_TTYS0_BASE));
	pnp_write_config(dev, PNP_IDX_IO0, 1);
	pnp_exit_conf_state(dev);
}

void wilco_ec_early_init(void)
{
	if (CONFIG(DRIVERS_UART_8250IO))
		wilco_ec_serial_init();
}
