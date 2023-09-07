/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2018 Patrick Rudolph <siro@das-labor.org>
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

#include <ec/acpi/ec.h>

#include "h8.h"

/**
 * Return the EC sense status register state.
 *
 * Observations showed the sense registers are all zero until the EC populates
 * them after some time. Likely the EC sets all bits to it's valid state at
 * once, but there's no prove as the firmware isn't available.
 *
 * Wait for any register having at least one bit set.
 * Unlikely that all register will be zero after booting has finished.
 *
 * @return 1 if the EC provides valid data in sense status registers
 */
int h8_get_sense_ready(void)
{
	static const u8 regs[] = { H8_STATUS0, H8_STATUS1, H8_STATUS2,
			H8_STATUS3};

	for (size_t i = 0; i < ARRAY_SIZE(regs); i++) {
		if (ec_read(regs[i]))
			return 1;
	}

	return 0;
}

/**
 * Return the state of Fn key.
 * Only valid if h8_get_sense_ready (see above) returns true.
 *
 * @return 1 if the key is pressed.
 */
int h8_get_fn_key(void)
{
	return ec_read(H8_STATUS0) & H8_STATUS0_FN_KEY_DOWN;
}
