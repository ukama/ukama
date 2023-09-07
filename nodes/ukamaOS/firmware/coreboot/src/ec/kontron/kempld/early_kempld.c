/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2017 secunet Security Networks AG
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
#include <arch/io.h>
#include <delay.h>

#include "chip.h"
#include "kempld.h"
#include "kempld_internal.h"

void kempld_write8(const uint8_t addr, const uint8_t data)
{
	outb(addr, KEMPLD_IDX);
	outb(data, KEMPLD_DAT);
}

uint8_t kempld_read8(const uint8_t addr)
{
	outb(addr, KEMPLD_IDX);
	return inb(KEMPLD_DAT);
}

int kempld_get_mutex(int timeout_ms)
{
	while (inb(KEMPLD_IDX) & KEMPLD_MUTEX_KEY && timeout_ms--)
		mdelay(1);
	return timeout_ms;
}

void kempld_release_mutex(void)
{
	outb(KEMPLD_MUTEX_KEY, KEMPLD_IDX);
}

void kempld_enable_uart_for_console(void)
{
	if (!CONFIG(CONSOLE_SERIAL))
		return;

	if (kempld_get_mutex(100) < 0)
		return;

	switch (CONFIG_UART_FOR_CONSOLE) {
	case 0:
		kempld_write8(KEMPLD_UART_0,
			      KEMPLD_UART_ENABLE |
				KEMPLD_UART_3F8 << KEMPLD_UART_IO_SHIFT);
		break;
	case 1:
		kempld_write8(KEMPLD_UART_1,
			      KEMPLD_UART_ENABLE |
				KEMPLD_UART_2F8 << KEMPLD_UART_IO_SHIFT);
		break;
	default:
		break;
	}

	kempld_release_mutex();
}
