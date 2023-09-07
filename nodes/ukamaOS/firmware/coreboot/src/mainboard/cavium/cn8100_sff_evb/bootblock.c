/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Facebook, Inc.
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
#include <soc/soc.h>
#include <soc/spi.h>
#include <soc/uart.h>

void bootblock_mainboard_early_init(void)
{
	if (CONFIG(BOOTBLOCK_CONSOLE)) {
		if (!uart_is_enabled(CONFIG_UART_FOR_CONSOLE))
			uart_setup(CONFIG_UART_FOR_CONSOLE, CONFIG_TTYS0_BAUD);
	}
}

static void configure_spi_flash(void)
{
	/* FIXME: Only tested on EM100 Pro */
	spi_init_custom(0, // bus
			25000000, // speed Hz
			0, // idle low disabled
			0, // zero idle cycles between transfers
			0, // MSB first
			0, // Chip select 0
			1); // assert is high
}

void bootblock_mainboard_init(void)
{
	configure_spi_flash();
}
