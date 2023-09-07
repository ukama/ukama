/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2015 Damien Zammit <damien@zamaudio.com>
 * Copyright (C) 2017 Arthur Heymans <arthur@aheymans.xyz>
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

#include <bootblock_common.h>
#include <device/pnp_ops.h>
#include <northbridge/intel/x4x/x4x.h>
#include <southbridge/intel/i82801gx/i82801gx.h>
#include <superio/nuvoton/common/nuvoton.h>
#include <superio/nuvoton/nct6776/nct6776.h>
#include <superio/winbond/common/winbond.h>
#include <superio/winbond/w83627dhg/w83627dhg.h>

#define SERIAL_DEV_R2 PNP_DEV(0x2e, NCT6776_SP1)
#define SERIAL_DEV_R1 PNP_DEV(0x2e, W83627DHG_SP1)
#define SUPERIO_DEV PNP_DEV(0x2e, 0)

void bootblock_mainboard_early_init(void)
{
	/* Set GPIOs on superio, enable UART */
	if (CONFIG(SUPERIO_NUVOTON_NCT6776)) {
		nuvoton_pnp_enter_conf_state(SERIAL_DEV_R2);
		pnp_set_logical_device(SERIAL_DEV_R2);

		pnp_write_config(SERIAL_DEV_R2, 0x1c, 0x80);
		pnp_write_config(SERIAL_DEV_R2, 0x27, 0x80);
		pnp_write_config(SERIAL_DEV_R2, 0x2a, 0x60);

		nuvoton_pnp_exit_conf_state(SERIAL_DEV_R2);
		nuvoton_enable_serial(SERIAL_DEV_R2, CONFIG_TTYS0_BASE);
	} else {
		winbond_enable_serial(SERIAL_DEV_R1, CONFIG_TTYS0_BASE);
	}
	/* IRQ routing */
	RCBA16(D31IR) = 0x0132;
	RCBA16(D29IR) = 0x0237;
}

void mb_get_spd_map(u8 spd_map[4])
{
	spd_map[0] = 0x50;
	spd_map[2] = 0x52;
}
