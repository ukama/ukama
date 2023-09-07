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

#ifndef EC_KONTRON_KEMPLD_INTERNAL_H
#define EC_KONTRON_KEMPLD_INTERNAL_H

#include <device/device.h>

/* i/o ports */
#define KEMPLD_IDX		0xa80
#define  KEMPLD_MUTEX_KEY	0x80
#define KEMPLD_DAT		0xa81

/* indexed registers */
#define KEMPLD_SPEC			0x06
#define  KEMPLD_SPEC_GET_MINOR(x)	(x & 0x0f)
#define  KEMPLD_SPEC_GET_MAJOR(x)	(x >> 4 & 0x0f)

#define KEMPLD_CFG			0x37
#define  KEMPLD_CFG_GPIO_I2C_MUX	(1 << 0)

#define KEMPLD_UART_0			0x30
#define KEMPLD_UART_1			0x31
#define  KEMPLD_UART_IRQ_SHIFT		0
#define  KEMPLD_UART_IRQ_MASK		(0xf << KEMPLD_UART_IRQ_SHIFT)
#define  KEMPLD_UART_IO_SHIFT		4
#define  KEMPLD_UART_IO_MASK		(0x3 << KEMPLD_UART_IO_SHIFT)
#define  KEMPLD_UART_ENABLE		0x80

#define KEMPLD_CLK 33333333 /* 33MHz */

void kempld_i2c_device_init(struct device *const dev);

#endif /* EC_KONTRON_KEMPLD_INTERNAL_H */
