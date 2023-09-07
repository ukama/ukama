/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2007-2010 coresystems GmbH
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

#include <bootblock_common.h>
#include <stdint.h>
#include <device/pnp_ops.h>
#include <cpu/x86/lapic.h>
#include <device/pnp_def.h>
#include <northbridge/intel/i945/i945.h>
#include <southbridge/intel/i82801gx/i82801gx.h>
#include <superio/winbond/common/winbond.h>
#include <superio/winbond/w83627ehg/w83627ehg.h>

#define SERIAL_DEV PNP_DEV(0x4e, W83627EHG_SP1)
#define SUPERIO_DEV PNP_DEV(0x4e, 0)

/* This box has one superio
 * Also set up the GPIOs from the beginning. This is the "no schematic
 * but safe anyways" method.
 */
void bootblock_mainboard_early_init(void)
{
	pnp_devfn_t dev;

	dev = SUPERIO_DEV;
	pnp_enter_conf_state(dev);

	pnp_write_config(dev, 0x24, 0xc4); // PNPCVS

	pnp_write_config(dev, 0x29, 0x01); // GPIO settings
	pnp_write_config(dev, 0x2a, 0x40); // GPIO settings should be fc but gets set to 02
	pnp_write_config(dev, 0x2b, 0xc0); // GPIO settings?
	pnp_write_config(dev, 0x2c, 0x03); // GPIO settings?
	pnp_write_config(dev, 0x2d, 0x20); // GPIO settings?

	dev = PNP_DEV(0x4e, W83627EHG_SP1);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);
	pnp_set_iobase(dev, PNP_IDX_IO0, 0x3f8);
	pnp_set_irq(dev, PNP_IDX_IRQ0, 4);
	pnp_set_enable(dev, 1);

	dev = PNP_DEV(0x4e, W83627EHG_SP2);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);
	pnp_set_iobase(dev, PNP_IDX_IO0, 0x2f8);
	pnp_set_irq(dev, PNP_IDX_IRQ0, 3);
	// pnp_write_config(dev, 0xf1, 4); // IRMODE0
	pnp_set_enable(dev, 1);

	dev = PNP_DEV(0x4e, W83627EHG_KBC); // Keyboard
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);
	pnp_set_iobase(dev, PNP_IDX_IO0, 0x60);
	pnp_set_iobase(dev, PNP_IDX_IO1, 0x64);
	//pnp_write_config(dev, 0xf0, 0x82);
	pnp_set_enable(dev, 1);

	dev = PNP_DEV(0x4e, W83627EHG_GPIO2);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 1); // Just enable it

	dev = PNP_DEV(0x4e, W83627EHG_GPIO3);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);
	pnp_write_config(dev, 0xf0, 0xfb); // GPIO bit 2 is output
	pnp_write_config(dev, 0xf1, 0x00); // GPIO bit 2 is 0
	pnp_write_config(dev, 0x30, 0x03); // Enable GPIO3+4. pnp_set_enable is not sufficient

	dev = PNP_DEV(0x4e, W83627EHG_FDC);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);

	dev = PNP_DEV(0x4e, W83627EHG_PP);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);

	/* Enable HWM */
	dev = PNP_DEV(0x4e, W83627EHG_HWM);
	pnp_set_logical_device(dev);
	pnp_set_enable(dev, 0);
	pnp_set_iobase(dev, PNP_IDX_IO0, 0xa00);
	pnp_set_enable(dev, 1);

	pnp_exit_conf_state(dev);
}

void mainboard_late_rcba_config(void)
{
	/* Device 1f interrupt pin register */
	RCBA32(D31IP) = 0x00042210;
	/* Device 1d interrupt pin register */
	RCBA32(D28IP) = 0x00214321;

	/* dev irq route register */
	RCBA16(D31IR) = 0x0132;
	RCBA16(D30IR) = 0x0146;
	RCBA16(D29IR) = 0x0237;
	RCBA16(D28IR) = 0x3201;
	RCBA16(D27IR) = 0x0146;

	/* Enable PCIe Root Port Clock Gate */
	// RCBA32(0x341c) = 0x00000001;
}
