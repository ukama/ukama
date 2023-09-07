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

#include <console/console.h>
#include <device/device.h>

#include "chip.h"
#include "kempld.h"
#include "kempld_internal.h"

static void kempld_uart_read_resources(struct device *dev)
{
	static const unsigned int io_addr[] = { 0x3f8, 0x2f8, 0x3e8, 0x2e8 };

	const struct ec_kontron_kempld_config *const config = dev->chip_info;

	struct resource *const res_io = new_resource(dev, 0);
	struct resource *const res_irq = new_resource(dev, 1);
	const unsigned int uart = dev->path.generic.subid;
	if (!config || !res_io || !res_irq || uart >= KEMPLD_NUM_UARTS)
		return;

	const enum kempld_uart_io io = config->uart[uart].io;
	if (io >= ARRAY_SIZE(io_addr)) {
		printk(BIOS_ERR, "KEMPLD: Bad io value '%d' for UART#%u\n.",
		       io, uart);
		dev->enabled = false;
		return;
	}

	const int irq = config->uart[uart].irq;
	if (irq >= 16) {
		printk(BIOS_ERR, "KEMPLD: Bad irq value '%d' for UART#%u\n.",
		       irq, uart);
		dev->enabled = false;
		return;
	}

	res_io->base = io_addr[io];
	res_io->size = 8;
	res_io->flags = IORESOURCE_IO | IORESOURCE_FIXED |
			IORESOURCE_STORED | IORESOURCE_ASSIGNED;
	res_irq->base = irq;
	res_irq->size = 1;
	res_irq->flags = IORESOURCE_IO | IORESOURCE_FIXED |
			 IORESOURCE_STORED | IORESOURCE_ASSIGNED;

	if (kempld_get_mutex(100) < 0)
		return;

	const uint8_t reg = uart ? KEMPLD_UART_1 : KEMPLD_UART_0;
	const uint8_t val = kempld_read8(reg);
	kempld_write8(reg, (val & ~(KEMPLD_UART_IO_MASK | KEMPLD_UART_IRQ_MASK))
			   | io << KEMPLD_UART_IO_SHIFT
			   | irq << KEMPLD_UART_IRQ_SHIFT);

	kempld_release_mutex();
}

static void kempld_uart_enable_resources(struct device *dev)
{
	if (kempld_get_mutex(100) < 0)
		return;

	const unsigned int uart = dev->path.generic.subid;
	const uint8_t reg = uart ? KEMPLD_UART_1 : KEMPLD_UART_0;
	kempld_write8(reg, kempld_read8(reg) | KEMPLD_UART_ENABLE);

	kempld_release_mutex();
}

static struct device_operations kempld_uart_ops = {
	.read_resources   = kempld_uart_read_resources,
	.enable_resources = kempld_uart_enable_resources,
};

static void kempld_enable_dev(struct device *const dev)
{
	if (dev->path.type == DEVICE_PATH_GENERIC) {
		switch (dev->path.generic.id) {
		case 0:
			if (dev->path.generic.subid < KEMPLD_NUM_UARTS) {
				dev->ops = &kempld_uart_ops;
				break;
			}
			/* Fall through. */
		case 1:
			if (dev->path.generic.subid == 0) {
				kempld_i2c_device_init(dev);
				break;
			}
			/* Fall through. */
		default:
			printk(BIOS_WARNING,
			       "KEMPLD: Spurious device %s.\n",
			       dev_path(dev));
			break;
		}
	}
}

struct chip_operations ec_kontron_kempld_ops = {
	CHIP_NAME("Kontron KEMPLD")
	.enable_dev = kempld_enable_dev,
};
