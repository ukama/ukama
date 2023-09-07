/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2019 Patrick Rudolph <siro@das-labor.org>
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
#include <bootmode.h>
#include <timer.h>
#include <delay.h>

#include "h8.h"

/**
 * HACK: Use Fn-Key as recovery mode switch.
 * Wait for sense register ready and read Fn-Key state.
 */
int get_recovery_mode_switch(void)
{
	struct stopwatch sw;

	if (!CONFIG(H8_FN_KEY_AS_VBOOT_RECOVERY_SW))
		return 0;

	/* Tests showed that it takes:
	 *  - 700msec on Lenovo T500 from AC power on
	 *  - less than 150msec on Lenovo T520 from AC power on
	 */
	stopwatch_init_msecs_expire(&sw, 1000);
	while (!stopwatch_expired(&sw) && !h8_get_sense_ready())
		mdelay(1);

	if (!h8_get_sense_ready())
		return 0;

	return h8_get_fn_key();
}

/**
 * Only used if CONFIG_CHROMEOS is set.
 * Always zero as the #WP pin of the flash is tied high.
 */
int get_write_protect_state(void)
{
	return 0;
}
